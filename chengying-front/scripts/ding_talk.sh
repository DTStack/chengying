#!/bin/bash
#参考钉钉文档 https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq

# Notify sonar report to dingding
sonarreport=$(curl -s http://172.16.100.198:8082/?projectname=dt-em-front)
curl -s "https://oapi.dingtalk.com/robot/send?access_token=d6824e5686b3d0d84ce7efc616e1fa32376fa7ecf42f0c70f8950efa4e404507" \
    -H "Content-Type: application/json" \
    -d "{
     \"msgtype\": \"markdown\",
     \"markdown\": {
         \"title\":\"sonar代码质量\",
         \"text\": \"## sonar代码质量报告: \n
> [sonar地址](http://172.16.100.198:9000/dashboard?id=dt-em-front) \n
> ${sonarreport} \n\"
     }
 }"
