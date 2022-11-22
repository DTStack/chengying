#! /bin/bash
export PATH=/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin
basedir=$(cd `dirname $0`/..; pwd)

source /etc/profile
. /lib/lsb/init-functions

#### dtstack updated
export HBASE_OPTS="$HBASE_OPTS -javaagent:${basedir}/dtstack/prometheus/jmx_prometheus_javaagent-0.12.0.jar=9530:${basedir}/dtstack/prometheus/hbase_master.yml -Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.local.only=false"
export HBASE_LOG_DIR=${basedir}/logs
export HBASE_PID_DIR=${basedir}/run
#### dtstack updated

$HBASE_HOME/bin/hbase-daemon.sh stop master
sleep 1
$HBASE_HOME/bin/hbase-daemon.sh start master

tail -f /dev/null

