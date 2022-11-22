#!/bin/sh
set -e

############################
#Copyright (c) 2017 DTStack Inc.
#Version  0.1
############################

script_wrapper_app='script_wrapper.tar.gz'
app_dir='/opt/dtstack/easyagent'
agent_bin="$app_dir/script_wrapper/script-wrapper"

#安装包下载地址
DOWNLOAD_URL='http://{{.easyagent_server}}/easyagent/script_wrapper.tar.gz'

set -e
trap '[ "$?" -eq 0 ] || read -p "Looks like something went wrong in step ´$STEP´"' EXIT

##install the script-wrapper##
install_script_wrapper() {
    STEP='install_script_wrapper'
    mkdir -p "$app_dir"
    tar -zxf "/tmp/$script_wrapper_app" -C "$app_dir"  >/dev/null 2>&1

     if [ -f "$agent_bin" ];then
         echo -e "...The script_wrapper install...  \033[34;1m[Success]\033[0m"
     else
         echo -e "...The script_wrapper install...  \033[31;1m[Failed]\033[0m"
         exit 1
     fi
}

##download and installed##
install(){
    STEP='install'
    echo -e " Use the curl download and install,\033[32;1m Please Waiting...\033[0m"
    cd /tmp/ && curl -L -O -s "$DOWNLOAD_URL"
    install_script_wrapper

    if [ ! -f "$agent_bin" ];then
        echo -e "\033[31;1mThe script_wrapper is faile to install....\033[0m"
        exit 1
    fi
}

##delete $script_wrapper pakg##
delete(){
    STEP='delete'
    cd /tmp/ && rm -f "$script_wrapper_app"
}

install
delete
