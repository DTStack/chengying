#!/bin/sh

current=`date "+%Y-%m-%d %H:%M:%S"`
timeStamp=`date -d "$current" +%s`

agentId={AGENTID}
interval={INTERVAL}

#agentId=cd3d4341-e9ad-446d-9932-8884a034d1cb
#interval=28

logs=/opt/dtstack/easymanager/easyagent/logs/agent.log

cdate=`date -d "-$interval hour" +"%Y/%m/%d %H"`

user=`whoami`
s=0

for i in $(cat $logs |grep "$cdate"|grep "exit(exit status 1"|grep "$agentId"| awk '{print$1"#"$2}')
do
ret=`echo $(echo $(echo "$i"|sed "s/\//-/g")| sed "s/AGENT-DEBUG://g" )|sed "s/#/ /g"`
results[s]='{"start_time":{"desc":"启动时间","value":"'$ret'"},"service_name":{"desc":"启动服务","value":"{SERVICENAME}"},"host":{"desc":"启动主机","value":"{HOSTIP}"},"product_name":{"desc":"所属组件","value":"{PRODUCTNAME}"},"run_user":{"desc":"启动用户","value":"'$user'"}}'
s=$[$s+1];
done

len=${#results[@]}
for ((i=$len - 1;i>=0;i--))
do
    echo ${results[$i]}
done