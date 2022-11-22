#!/bin/bash



source /etc/profile

HADOOP_CONF_DIR=/data/hadoop_base/etc/hadoop

CONF_FILE=$HADOOP_CONF_DIR/yarn-site.xml

#IP=`ip addr | grep inet | grep -v inet6 | grep -v 127 |awk -F ' ' '{print $2}'| awk -F '/' '{print $1}'`
IP=$1
RM1=`cat $CONF_FILE | grep yarn.resourcemanager.address.rm1 -A 1 | grep value | awk -F ':' '{print $1}' | awk -F '>' '{print $2}'`

RM2=`cat $CONF_FILE | grep yarn.resourcemanager.address.rm2 -A 1 | grep value | awk -F ':' '{print $1}' | awk -F '>' '{print $2}'`

#echo $IP $RM1 $RM2
RM1_STATUS=`yarn rmadmin -getServiceState rm1`
RM2_STATUS=`yarn rmadmin -getServiceState rm2`

state=`sed '/^STATE=/!d;s/.*=//' gcconfig`

if [ $RM1_STATUS != 'active' ] && [ $RM2_STATUS != 'active' ]; then
 ((state=$state+1))
 sed -i 's#^STATE=[0-9]*#STATE='"${state}"'#g' gcconfig
 if [ $state -gt 60 ]; then
   sed -i 's#^STATE=[0-9]*#STATE=1#g' gcconfig
   ps -ef |grep "org.apache.hadoop.yarn.server.resourcemanager.ResourceManager"| grep yarn_resourcemanager | grep -v grep |grep -v "start_resourcemanager"|grep -v "show_resourcemanager"|awk '{print $2}'|xargs kill -9
 fi
else
  sed -i 's#^STATE=[0-9]*#STATE=0#g' gcconfig
fi

nowtime=`date +%Y%m%d%H%M`
if [ ! -f marktime.txt ]; then
  echo $nowtime > marktime.txt
fi

filetime=`cat marktime.txt`
time01=$(($nowtime-$filetime))
if [ $time01 -gt 10 ]; then
  rm -f rm_state.txt
  echo $nowtime > marktime.txt
fi

if [ ! -f rm_state.txt ]
then
if [ "$IP" == "$RM1" ];then
  yarn rmadmin -getServiceState rm1 > rm_state.txt 2>&1
else
  yarn rmadmin -getServiceState rm2 > rm_state.txt 2>&1
fi
else
  cat rm_state.txt
fi
exit
