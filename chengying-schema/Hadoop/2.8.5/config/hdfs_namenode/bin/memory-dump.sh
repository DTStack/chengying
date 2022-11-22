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

jmap=/opt/dtstack/java/bin/jmap
output_dir=/tmp/dump

app_name=namenode
app_pid=$(ps aux |grep java |grep org.apache.hadoop.hdfs.server.namenode.NameNode |awk '{print $2}')

output_file=${output_dir}/dump_${app_name}_`date +%Y%m%d%H%M`_`hostname`.hprof

mkdir -p $output_dir

echo "app_pid == $app_pid"
echo "output_file == $output_file"
$jmap -dump:format=b,file=${output_file} $app_pid

