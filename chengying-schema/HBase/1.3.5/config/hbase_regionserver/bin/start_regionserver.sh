#!/usr/bin/env bash
basedir=$(cd `dirname $0`/..; pwd)

source /etc/profile

#### dtstack updated
export HBASE_OPTS="$HBASE_OPTS -javaagent:${basedir}/dtstack/prometheus/jmx_prometheus_javaagent-0.12.0.jar=9531:${basedir}/dtstack/prometheus/hbase_regionserver.yml -Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.local.only=false"
export HBASE_LOG_DIR=${basedir}/logs
export HBASE_PID_DIR=${basedir}/run
#### dtstack updated

$HBASE_HOME/bin/hbase-daemon.sh stop regionserver
sleep 1
$HBASE_HOME/bin/hbase-daemon.sh start regionserver

tail -f /dev/null

