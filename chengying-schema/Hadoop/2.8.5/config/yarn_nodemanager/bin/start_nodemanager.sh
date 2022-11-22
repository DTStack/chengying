#!/bin/bash  
source /etc/profile
export YARN_OPTS="$HADOOP_OPTS -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.local.only=false -Dcom.sun.management.jmxremote.port=9604 -javaagent:./dtstack/prometheus/jmx_prometheus_javaagent-0.3.1.jar=9504:./dtstack/prometheus/nodemanager.yml"

CMD_PATH=`dirname $0`
CMD_HOME=`cd "$CMD_PATH"/../; pwd`
bin=`dirname "${BASH_SOURCE-$0}"`
bin=`cd "$bin"; pwd`
hadoopScript="$HADOOP_HOME"/bin/hadoop



CGROUP_HOME=/sys/fs/cgroup/cpu,cpuacct/hadoop-yarn
#export HADOOP_HEAPSIZE=1024
#export YARN_HEAPSIZE=1024
DEFAULT_LIBEXEC_DIR="$bin"/../libexec
HADOOP_LIBEXEC_DIR=${HADOOP_LIBEXEC_DIR:-$DEFAULT_LIBEXEC_DIR}

if [ -e ${HADOOP_LIBEXEC_DIR}/yarn-config.sh ]; then
  . $HADOOP_LIBEXEC_DIR/yarn-config.sh
else
    HADOOP_LIBEXEC_DIR=/opt/dtstack/Hadoop/hadoop_pkg/libexec
    . $HADOOP_LIBEXEC_DIR/yarn-config.sh
fi


#. $HADOOP_LIBEXEC_DIR/yarn-config.sh

command=nodemanager


check_cgroup() {
	if [ ! -d $CGROUP_HOME ]; then
		sudo mkdir -p /sys/fs/cgroup/cpu,cpuacct/hadoop-yarn/
		sudo chown -R $USER:$USER /sys/fs/cgroup/cpu,cpuacct/hadoop-yarn
	fi
}


check_cgroup
check_kernel(){
cpu_kernel=`uname -a | awk '{print $3}'`

if [ "$cpu_kernel" = "3.10.0-327.el7.x86_64" ];
then
sed -i s/"org.apache.hadoop.yarn.server.nodemanager.util.CgroupsLCEResourcesHandler"/"org.apache.hadoop.yarn.server.nodemanager.util.DefaultLCEResourcesHandler"/g $HADOOP_HOME/etc/hadoop/yarn-site.xml
sed -i s/"org.apache.hadoop.yarn.server.nodemanager.LinuxContainerExecutor"/"org.apache.hadoop.yarn.server.nodemanager.DefaultContainerExecutor"/g $HADOOP_HOME/etc/hadoop/yarn-site.xml
fi
}

check_kernel

function check_BinaryFile() {

   if [ "" == "`ls /usr/lib64/ | grep libcrypto.so.1.1`" ]; then

    sudo cp -rp $HADOOP_HOME/lib/native/libcrypto.so.1.1 /usr/lib64/

   fi


}



check_startUser(){

sed -i s/"allowed.system.users={{.allowed_system_users}}"/"allowed.system.users=${USER}"/g /etc/yarn-cgroup/container-executor.cfg
sed -i s/"yarn.nodemanager.linux-container-executor.group={{.yarn_nodemanager_linux_container_executor_group}}"/"yarn.nodemanager.linux-container-executor.group=${USER}"/g /etc/yarn-cgroup/container-executor.cfg
sed -i s/"<value>{{.yarn_nodemanager_linux_container_executor_group}}<\/value>"/"<value>${USER}<\/value>"/g $HADOOP_HOME/etc/hadoop/yarn-site.xml

}

check_startUser


hadoop_rotate_log ()
{
    log=$1;
    num=5;
    if [ -n "$2" ]; then
	num=$2
    fi
    if [ -f "$log" ]; then # rotate logs
	while [ $num -gt 1 ]; do
	    prev=`expr $num - 1`
	    [ -f "$log.$prev" ] && mv "$log.$prev" "$log.$num"
	    num=$prev
	done
	mv "$log" "$log.$num";
    fi
}

function get_log_dir() {

    EXTERNAL_LOG_DIR={{.external_log_dir}}
    if [ "${EXTERNAL_LOG_DIR}" == "" ]; then
    YARN_LOG_DIR="$CMD_HOME"/logs
    export YARN_LOG_DIR="$CMD_HOME"/logs
    else
    YARN_LOG_DIR="${EXTERNAL_LOG_DIR}"/yarn
    export YARN_LOG_DIR="${EXTERNAL_LOG_DIR}"/yarn
    fi

}

#获取日志路径
get_log_dir



if [ -f "${YARN_CONF_DIR}/yarn-env.sh" ]; then
  . "${YARN_CONF_DIR}/yarn-env.sh"
fi

if [ "$YARN_IDENT_STRING" = "" ]; then
  export YARN_IDENT_STRING="$USER"
fi

# get log directory
if [ "$YARN_LOG_DIR" = "" ]; then
  export YARN_LOG_DIR="$HADOOP_YARN_HOME/logs"
fi

if [ ! -w "$YARN_LOG_DIR" ] ; then
  mkdir -p "$YARN_LOG_DIR"
  chown $YARN_IDENT_STRING $YARN_LOG_DIR 
fi

# some variables
export YARN_LOGFILE=yarn-$YARN_IDENT_STRING-$command-$HOSTNAME.log
export YARN_ROOT_LOGGER=${YARN_ROOT_LOGGER:-INFO,RFA}
log=$YARN_LOG_DIR/yarn-$YARN_IDENT_STRING-$command-$HOSTNAME.out
YARN_STOP_TIMEOUT=${YARN_STOP_TIMEOUT:-5}

# Set default scheduling priority
if [ "$YARN_NICENESS" = "" ]; then
    export YARN_NICENESS=0
fi

hadoop_rotate_log $log
nice -n $YARN_NICENESS "$HADOOP_HOME"/bin/yarn --config $YARN_CONF_DIR $command "$@" > "$log" 2>&1 < /dev/null

