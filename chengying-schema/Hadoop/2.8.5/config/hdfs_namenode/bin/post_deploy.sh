#!/bin/bash 
set -e

source /etc/profile

mkdir -p /data/hadoop/hdfs
if [ -d /data/hadoop/hdfs/name ]; then
    echo "success"  
    exit 0
fi

init=`cat namenode_init`
if [ "$init" == "yes" ]; then
    $HADOOP_HOME/bin/hadoop namenode -format -nonInteractive
fi
if [ -z "$init" ];then
    sleep 5
    $HADOOP_HOME/bin/hadoop namenode -bootstrapStandby -nonInteractive
fi
