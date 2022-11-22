#!/bin/bash -ex

# cp -r ../dt-alert/alert-front/dist dt-alert
yarn build
# docker build -t 124.71.131.239:8889/zhijian/manage.front:log_v1  .
# docker push 124.71.131.239:8889/zhijian/manage.front:log_v1 
docker build -t 172.16.84.121/dtstack-dev/manage.front:debug  .
docker push 172.16.84.121/dtstack-dev/manage.front:debug

ssh root@172.16.82.176 " \
docker-compose -f /data/em-release-4.2.1-beta/easymanager/docker-compose.yml stop manage_front;\
docker-compose -f /data/em-release-4.2.1-beta/easymanager/docker-compose.yml rm -f manage_front;\
docker-compose -f /data/em-release-4.2.1-beta/easymanager/docker-compose.yml pull manage_front; \
docker-compose -f /data/em-release-4.2.1-beta/easymanager/docker-compose.yml up -d;"