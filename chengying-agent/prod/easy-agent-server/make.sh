#!/bin/sh

gox -os=linux -arch=amd64
mv easy-agent-server_linux_amd64 easy-agent-server
docker build -t registry.cn-hangzhou.aliyuncs.com/dtstack-dev/easy-agent-server:debug .
docker push registry.cn-hangzhou.aliyuncs.com/dtstack-dev/easy-agent-server:debug