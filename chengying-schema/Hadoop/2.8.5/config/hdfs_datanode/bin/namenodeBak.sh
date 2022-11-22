#!/bin/bash

##------------------------
##namenode备份逻辑
##每天备份一次
##备份目录：/data/hadoop/hdfs/namenodebak/
##------------------------
source /etc/profile


nowtime=`date +%Y%m%d`

BACK_TIME_FILE=/data/hadoop/hdfs/.namenobak_time.txt

 mkdir -p /data/hadoop/hdfs/namenodebak
 if [ $? != 0 ]; then
   printf "备份目录创建失败，请检查权限问题\n"
   exit 1
 fi

if [ ! -f ${BACK_TIME_FILE} ]; then

  echo $nowtime > ${BACK_TIME_FILE}

fi

filetime=`cat ${BACK_TIME_FILE}`


status=0


if [ $nowtime -ne $filetime ]
then
  hdfs dfsadmin -fetchImage  /data/hadoop/hdfs/namenodebak/
  cp -rap /data/hadoop/hdfs/name/current/seen_txid /data/hadoop/hdfs/namenodebak/
  cp -rap /data/hadoop/hdfs/name/current/VERSION /data/hadoop/hdfs/namenodebak/
  echo $nowtime > ${BACK_TIME_FILE}

  nn1=`hdfs haadmin -getServiceState nn1`
  if [ $? -ne 0 ]
  then
    status=1
  fi

  nn2=`hdfs haadmin -getServiceState nn2`
  if [ $? -ne 0 ]
  then
    status=1
  fi

  if [ $status -eq 0 ]
  then
    find /data/hadoop/hdfs/namenodebak  -mtime +7 -name "fsimage_*" -exec rm -f {} \;
  fi

fi

