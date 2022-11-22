#!/bin/sh

gox -os=linux -arch=amd64
mv easymonitor_linux_amd64 easymonitor
docker build -t registry.cn-hangzhou.aliyuncs.com/dtstack/easymonitor:v1.0.3 .
docker tag registry.cn-hangzhou.aliyuncs.com/dtstack/easymonitor:v1.0.3 172.16.8.120:5443/dtstack-dev/easymonitor:v1.0.3
docker push registry.cn-hangzhou.aliyuncs.com/dtstack/easymonitor:v1.0.3
docker push 172.16.8.120:5443/dtstack-dev/easymonitor:v1.0.3