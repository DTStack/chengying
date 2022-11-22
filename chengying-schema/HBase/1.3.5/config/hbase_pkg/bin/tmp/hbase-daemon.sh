#!/usr/bin/env bash
#
#/**
# * Licensed to the Apache Software Foundation (ASF) under one
# * or more contributor license agreements.  See the NOTICE file
# * distributed with this work for additional information
# * regarding copyright ownership.  The ASF licenses this file
# * to you under the Apache License, Version 2.0 (the
# * "License"); you may not use this file except in compliance
# * with the License.  You may obtain a copy of the License at
# *
# *     http://www.apache.org/licenses/LICENSE-2.0
# *
# * Unless required by applicable law or agreed to in writing, software
# * distributed under the License is distributed on an "AS IS" BASIS,
# * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# * See the License for the specific language governing permissions and
# * limitations under the License.
# */
#
# Runs a Hadoop hbase command as a daemon.
#
# Environment Variables
#
#   HBASE_CONF_DIR   Alternate hbase conf dir. Default is ${HBASE_HOME}/conf.
#   HBASE_LOG_DIR    Where log files are stored.  PWD by default.
#   HBASE_PID_DIR    The pid files are stored. /tmp by default.
#   HBASE_IDENT_STRING   A string representing this instance of hadoop. $USER by default
#   HBASE_NICENESS The scheduling priority for daemons. Defaults to 0.
#   HBASE_STOP_TIMEOUT  Time, in seconds, after which we kill -9 the server if it has not stopped.
#                        Default 1200 seconds.
#
# Modelled after $HADOOP_HOME/bin/hadoop-daemon.sh

usage="Usage: hbase-daemon.sh [--config <conf-dir>]\
 [--autostart-window-size <window size in hours>]\
 [--autostart-window-retry-limit <retry count limit for autostart>]\
 (start|stop|restart|autostart|autorestart|foreground_start) <hbase-command> \
 <args...>"

# if no args specified, show usage
if [ $# -le 1 ]; then
  echo $usage
  exit 1
fi

# default autostart args value indicating infinite window size and no retry limit
AUTOSTART_WINDOW_SIZE=0
AUTOSTART_WINDOW_RETRY_LIMIT=0

bin=`dirname "${BASH_SOURCE-$0}"`
bin=`cd "$bin">/dev/null; pwd`

. "$bin"/hbase-config.sh
. "$bin"/hbase-common.sh

# get arguments
startStop=$1
shift

command=$1
shift


#if [ "$command" = "master" ]; then
#  HBASE_OPTS="$HBASE_OPTS -javaagent:$HBASE_HOME/dtstack/prometheus/jmx_prometheus_javaagent-0.12.0.jar=9530:$HBASE_HOME/dtstack/prometheus/hbase_jmx_config.yaml"
#elif [ "$command" = "regionserver" ]; then
#  HBASE_OPTS="$HBASE_OPTS -javaagent:$HBASE_HOME/dtstack/prometheus/jmx_prometheus_javaagent-0.12.0.jar=9531:$HBASE_HOME/dtstack/prometheus/hbase_jmx_config.yaml"
#fi

hbase_rotate_log ()
{
    log=$1;
    num=5;
    if [ -n "$2" ]; then
    num=$2
    fi
    if [ -f "$log" ]; then # rotate logs
    while [ $num -gt 1 ]; do
        prev=`expr $num - 1`
        [ -f "$log.$prev" ] && mv -f "$log.$prev" "$log.$num"
        num=$prev
    done
    mv -f "$log" "$log.$num";
    fi
}

cleanAfterRun() {
  if [ -f ${HBASE_PID} ]; then
    # If the process is still running time to tear it down.
    kill -9 `cat ${HBASE_PID}` > /dev/null 2>&1
    rm -f ${HBASE_PID} > /dev/null 2>&1
  fi

  if [ -f ${HBASE_ZNODE_FILE} ]; then
    if [ "$command" = "master" ]; then
      HBASE_OPTS="$HBASE_OPTS $HBASE_MASTER_OPTS" $bin/hbase master clear > /dev/null 2>&1
    else
      #call ZK to delete the node
      ZNODE=`cat ${HBASE_ZNODE_FILE}`
      HBASE_OPTS="$HBASE_OPTS $HBASE_REGIONSERVER_OPTS" $bin/hbase zkcli delete ${ZNODE} > /dev/null 2>&1
    fi
    rm ${HBASE_ZNODE_FILE}
  fi
}

check_before_start(){
    #ckeck if the process is not running
    mkdir -p "$HBASE_PID_DIR"
    if [ -f $HBASE_PID ]; then
      if kill -0 `cat $HBASE_PID` > /dev/null 2>&1; then
        echo $command running as process `cat $HBASE_PID`.  Stop it first.
        exit 1
      fi
    fi
}

wait_until_done ()
{
    p=$1
    cnt=${HBASE_SLAVE_TIMEOUT:-300}
    origcnt=$cnt
    while kill -0 $p > /dev/null 2>&1; do
      if [ $cnt -gt 1 ]; then
        cnt=`expr $cnt - 1`
        sleep 1
      else
        echo "Process did not complete after $origcnt seconds, killing."
        kill -9 $p
        exit 1
      fi
    done
    return 0
}

# get log directory
if [ "$HBASE_LOG_DIR" = "" ]; then
  export HBASE_LOG_DIR="$HBASE_HOME/logs"
fi
mkdir -p "$HBASE_LOG_DIR"

if [ "$HBASE_PID_DIR" = "" ]; then
  HBASE_PID_DIR=/tmp
fi

if [ "$HBASE_IDENT_STRING" = "" ]; then
  export HBASE_IDENT_STRING="$USER"
fi

# Some variables
# Work out java location so can print version into log.
if [ "$JAVA_HOME" != "" ]; then
  #echo "run java in $JAVA_HOME"
  JAVA_HOME=$JAVA_HOME
fi
if [ "$JAVA_HOME" = "" ]; then
  echo "Error: JAVA_HOME is not set."
  exit 1
fi

JAVA=$JAVA_HOME/bin/java
export HBASE_LOG_PREFIX=hbase-$HBASE_IDENT_STRING-$command-$HOSTNAME
export HBASE_LOGFILE=$HBASE_LOG_PREFIX.log

if [ -z "${HBASE_ROOT_LOGGER}" ]; then
export HBASE_ROOT_LOGGER=${HBASE_ROOT_LOGGER:-"INFO,RFA"}
fi

if [ -z "${HBASE_SECURITY_LOGGER}" ]; then
export HBASE_SECURITY_LOGGER=${HBASE_SECURITY_LOGGER:-"INFO,RFAS"}
fi

HBASE_LOGOUT=${HBASE_LOGOUT:-"$HBASE_LOG_DIR/$HBASE_LOG_PREFIX.out"}
HBASE_LOGGC=${HBASE_LOGGC:-"$HBASE_LOG_DIR/$HBASE_LOG_PREFIX.gc"}
HBASE_LOGLOG=${HBASE_LOGLOG:-"${HBASE_LOG_DIR}/${HBASE_LOGFILE}"}
HBASE_PID=$HBASE_PID_DIR/hbase-$HBASE_IDENT_STRING-$command.pid
export HBASE_ZNODE_FILE=$HBASE_PID_DIR/hbase-$HBASE_IDENT_STRING-$command.znode
export HBASE_AUTOSTART_FILE=$HBASE_PID_DIR/hbase-$HBASE_IDENT_STRING-$command.autostart

if [ -n "$SERVER_GC_OPTS" ]; then
  export SERVER_GC_OPTS=${SERVER_GC_OPTS/"-Xloggc:<FILE-PATH>"/"-Xloggc:${HBASE_LOGGC}"}
fi
if [ -n "$CLIENT_GC_OPTS" ]; then
  export CLIENT_GC_OPTS=${CLIENT_GC_OPTS/"-Xloggc:<FILE-PATH>"/"-Xloggc:${HBASE_LOGGC}"}
fi

# Set default scheduling priority
if [ "$HBASE_NICENESS" = "" ]; then
    export HBASE_NICENESS=0
fi

thiscmd="$bin/$(basename ${BASH_SOURCE-$0})"
args=$@


#check_before_start
hbase_rotate_log $HBASE_LOGOUT
hbase_rotate_log $HBASE_LOGGC
echo running $command, logging to $HBASE_LOGOUT



sleep 1; touch "${HBASE_LOGOUT}"
trap cleanAfterRun SIGHUP SIGINT SIGTERM EXIT
echo "`date` Starting $command on `hostname`" >> ${HBASE_LOGLOG}
echo "`ulimit -a`" >> "$HBASE_LOGLOG" 2>&1
# in case the parent shell gets the kill make sure to trap signals.
# Only one will get called. Either the trap or the flow will go through.
nice -n $HBASE_NICENESS "$HBASE_HOME"/bin/hbase \
    --config "${HBASE_CONF_DIR}" \
    $command "$@" start >> ${HBASE_LOGOUT} 2>&1
