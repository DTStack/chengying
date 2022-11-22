#!/bin/bash

source /etc/profile
#export HADOOP_HEAPSIZE=$1
export HADOOP_CLIENT_OPTS=" $1 -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.local.only=false -Dcom.sun.management.jmxremote.port=9606 -javaagent:/opt/dtstack/Hive/hive_pkg/dtstack/prometheus/jmx_prometheus_javaagent-0.3.1.jar=9506:./dtstack/prometheus/metastore.yml"
if [ ! -d "{{.hive_meta_logs}}" ];then
  mkdir -p {{.hive_meta_logs}}
fi
/opt/dtstack/Hive/hive_pkg/bin/hive --service metastore >> {{.hive_meta_logs}}/metastore-`date +%Y-%m-%d`.log


