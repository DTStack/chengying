#!/bin/bash -ex

cp ../../dt-alert/sql/alert.ddl alert.ddl.sql
docker build --no-cache -t registry.cn-hangzhou.aliyuncs.com/dtstack-dev/manage.sql:debug .
docker push registry.cn-hangzhou.aliyuncs.com/dtstack-dev/manage.sql:debug
