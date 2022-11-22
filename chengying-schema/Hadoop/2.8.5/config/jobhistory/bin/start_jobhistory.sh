#!/bin/bash

source /etc/profile
export HADOOP_OPTS="$HADOOP_OPTS -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.local.only=false -Dcom.sun.management.jmxremote.port=9619 -javaagent:./dtstack/prometheus/jmx_prometheus_javaagent-0.3.1.jar=9519:./dtstack/prometheus/historyserver.yml"

bin=`dirname "${BASH_SOURCE-$0}"`
bin=`cd "$bin"; pwd`

DEFAULT_LIBEXEC_DIR="$bin"/../libexec
HADOOP_LIBEXEC_DIR=${HADOOP_LIBEXEC_DIR:-$DEFAULT_LIBEXEC_DIR}
if [ -e ${HADOOP_LIBEXEC_DIR}/mapred-config.sh ]; then
  . $HADOOP_LIBEXEC_DIR/mapred-config.sh
fi


command=historyserver

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
    HADOOP_MAPRED_LOG_DIR="$CMD_HOME"/logs
    export HADOOP_MAPRED_LOG_DIR="$CMD_HOME"/logs
    else
    HADOOP_MAPRED_LOG_DIR="${EXTERNAL_LOG_DIR}"/mapred
    export HADOOP_MAPRED_LOG_DIR="${EXTERNAL_LOG_DIR}"/mapred
    fi

}

#获取日志路径
get_log_dir


if [ "$HADOOP_MAPRED_IDENT_STRING" = "" ]; then
  export HADOOP_MAPRED_IDENT_STRING="$USER"
fi

export HADOOP_MAPRED_HOME=${HADOOP_MAPRED_HOME:-${HADOOP_PREFIX}}
export HADOOP_MAPRED_LOGFILE=mapred-$HADOOP_MAPRED_IDENT_STRING-$command-$HOSTNAME.log
export HADOOP_MAPRED_ROOT_LOGGER=${HADOOP_MAPRED_ROOT_LOGGER:-INFO,RFA}
export HADOOP_JHS_LOGGER=${HADOOP_JHS_LOGGER:-INFO,JSA}

if [ -f "${HADOOP_CONF_DIR}/mapred-env.sh" ]; then
  . "${HADOOP_CONF_DIR}/mapred-env.sh"
fi

if [ ! -w "$HADOOP_MAPRED_LOG_DIR" ]; then
  mkdir -p "$HADOOP_MAPRED_LOG_DIR"
  chown "$HADOOP_MAPRED_IDENT_STRING" "$HADOOP_MAPRED_LOG_DIR"
fi

if [ "$HADOOP_MAPRED_PID_DIR" = "" ]; then
  HADOOP_MAPRED_PID_DIR=/tmp
fi

HADOOP_OPTS="$HADOOP_OPTS -Dhadoop.id.str=$HADOOP_MAPRED_IDENT_STRING"

log=$HADOOP_MAPRED_LOG_DIR/mapred-$HADOOP_MAPRED_IDENT_STRING-$command-$HOSTNAME.out
pid=$HADOOP_MAPRED_PID_DIR/mapred-$HADOOP_MAPRED_IDENT_STRING-$command.pid

HADOOP_MAPRED_STOP_TIMEOUT=${HADOOP_MAPRED_STOP_TIMEOUT:-5}

# Set default scheduling priority
if [ "$HADOOP_MAPRED_NICENESS" = "" ]; then
  export HADOOP_MAPRED_NICENESS=0
fi

nice -n $HADOOP_MAPRED_NICENESS "$HADOOP_HOME"/bin/mapred --config $HADOOP_CONF_DIR $command "$@" > "$log" 2>&1 < /dev/null

