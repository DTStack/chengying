#!/bin/bash

zookeeper_host=$1
export LANG=en_US.UTF-8
mkdir -p /data/zookeeper/data
id=`cat conf/zoo.cfg | grep ${zookeeper_host} | awk -F '=' '{print $1}' | awk -F '.' '{print $2}'`
echo $id > /data/zookeeper/myid
