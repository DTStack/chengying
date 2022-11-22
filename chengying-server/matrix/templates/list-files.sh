#!/bin/sh

current=`date "+%Y-%m-%d %H:%M:%S"`
timeStamp=`date -d "$current" +%s`

tmp="/tmp/heapdumps_{SERVICENAME}_$timeStamp"

find -L /opt/dtstack/{PRODUCTNAME}/{SERVICENAME}/*/heapdump.hprof -maxdepth 5 -size -4096M -type f -mmin -5 -print > $tmp 2>/dev/null

if [ -f "$tmp" ];then
for i in `cat $tmp`
do
echo '{"file_name":{"desc":"JavaHeapDump文件名称","value":"'$i'"},"product_name":{"desc":"所属组件","value":"{PRODUCTNAME}"},"service_name":{"desc":"服务名称","value":"{SERVICENAME}"},"host":{"desc":"主机IP","value":"{HOSTIP}"},"generate_time":{"desc":"生成时间","value":"'$current'"},"action":{"desc":"操作","value":"下载","instance":"{INSTANCEID}","path":"'$i'"}}'
break
done
fi

rm -f $tmp
