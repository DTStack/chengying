#!/bin/bash

source /etc/profile

CMD_PATH=`dirname $0`
CMD_HOME=`cd "$CMD_PATH"/./; pwd`

init=`cat thriftserver_init`

#判断是否是其他集群
is_other={{.is_others}}
#spark当前版本号
spark_version=`cat thriftserver_version`

path=/opt/dtstack/Spark/spark_pkg


if [ -f "/data/init_thriftserver" ];then
  old_version=`cat /data/init_thriftserver`
  if [ "$spark_version" != "$old_version" ];then
    if [ "$init" == "yes" ]; then
    hdfs dfs -mkdir -p /dtInsight/user/spark/client/
    hdfs dfs -mkdir -p /dtInsight/sparkjars/jars/
    hdfs dfs -mkdir -p /dtInsight/pythons/
    if [[ $is_other == "no" ]];then
    hdfs dfs -put $path/jars/* /dtInsight/sparkjars/jars/
    if [ $? != 0 ]; then
      exit 1
      fi
    hdfs dfs -put $path/python/lib/* /dtInsight/pythons/
    if [ $? != 0 ]; then
      exit 1
      fi
    hdfs dfs -put $path/sparkjar/spark-sql-proxy.jar /dtInsight/user/spark/client/
    if [ $? != 0 ]; then
      exit 1
      fi
    fi
    echo "$spark_version" > /data/init_thriftserver
   fi
   fi
elif [ ! -f "/data/init_thriftserver" ]; then
    if [ "$init" == "yes" ]; then
hdfs dfs -mkdir -p /dtInsight/user/spark/client/
hdfs dfs -mkdir -p /tmp/spark-yarn-logs/
hdfs dfs -chmod -R g+w /dtInsight/user/spark/client/
hdfs dfs -chmod -R 777 /dtInsight/user/spark/client/
hdfs dfs -mkdir -p /dtInsight/sparkjars/jars/
hdfs dfs -mkdir -p /dtInsight/pythons/
hdfs dfs -mkdir -p /dtInsight/user/dtscript/
if [[ $is_other == "no" ]];then
hdfs dfs -put $path/sparkjar/spark-sql-proxy.jar /dtInsight/user/spark/client/
if [ $? != 0 ]; then
      exit 1
      fi
hdfs dfs -put $path/jars/* /dtInsight/sparkjars/jars/
if [ $? != 0 ]; then
      exit 1
      fi
hdfs dfs -put $path/python/lib/* /dtInsight/pythons/
if [ $? != 0 ]; then
      exit 1
      fi
fi
touch /data/init_thriftserver
echo "$spark_version" > /data/init_thriftserver
fi
fi



add_crond(){

   CRONDFILF=/var/spool/cron/${USER}
   IFS=$'\n\n'

   for i in `cat crond`
   do
   if [[ "${i}" != "" ]] && [[ -z `sudo grep "${CMD_HOME}" "$CRONDFILF"` ]]; then
   echo ${i} | sudo tee -a ${CRONDFILF} > /dev/null
   fi
   done

}


add_crond
