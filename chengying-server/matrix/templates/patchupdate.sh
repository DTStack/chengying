#!/bin/bash
set -ex

############################
#Copyright (c) 2017 DTStack Inc.
#Version  0.1git
############################
##如：/opt/dtstack/DTBase/redis/bin/abc.jar
patch_path='{{.PATCHES_PATH}}'

#此参数值被程序使用，修改时请@huanxiong
#如：/opt/dtstack/DTBase/redis
app_dir='{{.AGENT_DIR}}'

#安装包下载地址
#(http://172.16.8.87:8864/easyagent/DTBase/2.0.9/redis/patches_package/abc.jar)
DOWNLOAD_URL='{{.AGENT_DOWNLOAD_URL}}'

trap '[ "$?" -eq 0 ] || read -p "Looks like something went wrong in step ´$STEP´"' EXIT

##判断备份目录是否存在
backup_dir(){
  STEP='create backup dir'
  if [ ! -d ${app_dir}/patch_backup ]
  then
      mkdir -p ${app_dir}/patch_backup
  fi
}

##下载补丁包
download_patch(){
    STEP='download patch'
    echo "Use the curl download and replace Please Waiting..."
    package_name=${patch_path##*/}
    cp ${patch_path} ${app_dir}/patch_backup/${package_name}.$(date +%s)
    cd /tmp/ && curl -L -O -s "$DOWNLOAD_URL"
    \mv /tmp/${DOWNLOAD_URL##*/} ${patch_path}
}

backup_dir
download_patch
