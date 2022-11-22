#!/bin/bash
source /etc/profile

HADOOP_CONF_DIR=/data/hadoop_base/etc/hadoop

CONF_FILE=$HADOOP_CONF_DIR/hdfs-site.xml

IP=$1
NN1=`cat $CONF_FILE | grep dfs.namenode.rpc-address.ns1.nn1 -A 1 | grep value | awk -F ':' '{print $1}' | awk -F '>' '{print $2}'`

NN2=`cat $CONF_FILE | grep dfs.namenode.rpc-address.ns1.nn1 -A 1 | grep value | awk -F ':' '{print $1}' | awk -F '>' '{print $2}'`

#echo $IP $RM1 $RM2

nowtime=`date +%Y%m%d%H%M`
if [ ! -f marktime.txt ]; then
  echo $nowtime > marktime.txt
fi

filetime=`cat marktime.txt`
time01=$(($nowtime-$filetime))
if [ $time01 -gt 10 ]; then
  rm -f nn_state.txt
  echo $nowtime > marktime.txt
fi


if [ ! -f nn_state.txt ]
then
if [ $IP == $NN1 ];then
  status1=`hdfs haadmin -getServiceState nn1`
  if [ $status1 == "active" ] || [ $status1 == "standby" ];
  then
  echo $status1 > nn_state.txt 2>&1
  else
  echo '获取异常' > nn_state.txt 2>&1
  fi
else
  status2=`hdfs haadmin -getServiceState nn2`
  if [ $status2 == active ] || [ $status2 == standby ];
  then
  echo $status2 > nn_state.txt 2>&1
  else
  echo '获取异常' > nn_state.txt 2>&1
  fi
fi
else
 cat nn_state.txt
fi


exit
