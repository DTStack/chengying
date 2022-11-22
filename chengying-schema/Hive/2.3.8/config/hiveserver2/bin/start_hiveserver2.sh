#!/bin/bash

source /etc/profile
#export HADOOP_HEAPSIZE=1024
export HADOOP_CLIENT_OPTS=" $1 "

if [ ! -d "{{.hive_logs}}" ]; then
  mkdir -p {{.hive_logs}}
fi

/opt/dtstack/Hive/hive_pkg/bin/hive --service hiveserver2 >> {{.hive_logs}}/server2.log 2>&1

