#!/bin/bash 
set -e

source /etc/profile

mkdir -p /data/hadoop/hdfs
if [ -f /data/hadoop/hdfs/formatZK.success ]; then
    echo "success"  
    exit 0
fi

init=`cat zkfc_init`
if [ "$init" == "yes" ]; then
    $HADOOP_HOME/bin/hdfs zkfc -formatZK -nonInteractive
    touch /data/hadoop/hdfs/formatZK.success
fi
