#!/bin/bash
export PATH=/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin

basedir=$(cd `dirname $0`/..; pwd)

checkUser() {
  if [ "`whoami`" == "root" ]; then
   echo "The root user is not allowed to run."
   exit 1
  fi
}
checkUser

jstat=/opt/dtstack/java/bin/jstat

app_name=spark_historyserver
app_pid=$(ps aux |grep java |grep org.apache.spark.deploy.history.HistoryServer |awk '{print $2}')

echo "app_pid == $app_pid"
$jstat -gcutil $app_pid 1000

