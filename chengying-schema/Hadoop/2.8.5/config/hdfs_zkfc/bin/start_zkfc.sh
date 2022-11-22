#!/bin/bash  

source /etc/profile
export HADOOP_OPTS="$HADOOP_OPTS -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.local.only=false -Dcom.sun.management.jmxremote.port=9618 -javaagent:./dtstack/prometheus/jmx_prometheus_javaagent-0.3.1.jar=9518:./dtstack/prometheus/zkfc.yml"
mkdir -p /data/hadoop/hdfs

CMD_PATH=`dirname $0`
CMD_HOME=`cd "$CMD_PATH"/../; pwd`
bin=`dirname "${BASH_SOURCE-$0}"`
bin=`cd "$bin"; pwd`
hadoopScript="$HADOOP_HOME"/bin/hadoop



function get_log_dir() {

    EXTERNAL_LOG_DIR={{.external_log_dir}}
    if [ "${EXTERNAL_LOG_DIR}" == "" ]; then
    HADOOP_LOG_DIR="$CMD_HOME"/logs
    export HADOOP_LOG_DIR="$CMD_HOME"/logs
    else
    HADOOP_LOG_DIR="${EXTERNAL_LOG_DIR}"/hdfs
    export HADOOP_LOG_DIR="${EXTERNAL_LOG_DIR}"/hdfs
    fi

}

#获取日志路径
get_log_dir



command=zkfc

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

if [ -f "${HADOOP_CONF_DIR}/hadoop-env.sh" ]; then
  . "${HADOOP_CONF_DIR}/hadoop-env.sh"
fi


if [ "$HADOOP_IDENT_STRING" = "" ]; then
  export HADOOP_IDENT_STRING="$USER"
fi


# get log directory
if [ "$HADOOP_LOG_DIR" = "" ]; then
  export HADOOP_LOG_DIR="$HADOOP_HOME/logs"
fi

if [ ! -w "$HADOOP_LOG_DIR" ] ; then
  mkdir -p "$HADOOP_LOG_DIR"
  chown $HADOOP_IDENT_STRING $HADOOP_LOG_DIR
fi

if [ "$HADOOP_PID_DIR" = "" ]; then
  HADOOP_PID_DIR=/tmp
fi

# Set default scheduling priority
if [ "$HADOOP_NICENESS" = "" ]; then
    export HADOOP_NICENESS=0
fi

# some variables
export HADOOP_LOGFILE=hadoop-$HADOOP_IDENT_STRING-$command-$HOSTNAME.log
export HADOOP_ROOT_LOGGER=${HADOOP_ROOT_LOGGER:-"INFO,RFA"}
export HADOOP_SECURITY_LOGGER=${HADOOP_SECURITY_LOGGER:-"INFO,RFAS"}
export HDFS_AUDIT_LOGGER=${HDFS_AUDIT_LOGGER:-"INFO,NullAppender"}
log=$HADOOP_LOG_DIR/hadoop-$HADOOP_IDENT_STRING-$command-$HOSTNAME.out
pid=$HADOOP_PID_DIR/hadoop-$HADOOP_IDENT_STRING-$command.pid
HADOOP_STOP_TIMEOUT=${HADOOP_STOP_TIMEOUT:-5}

hadoop_rotate_log $log
hdfsScript="$HADOOP_HOME"/bin/hdfs

nice -n $HADOOP_NICENESS $hdfsScript --config $HADOOP_CONF_DIR $command "$@" > "$log" 2>&1 < /dev/null 


