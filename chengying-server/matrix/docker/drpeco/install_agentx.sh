#!/bin/bash
set -ex

############################
#Copyright (c) 2017 DTStack Inc.
#Version  0.1
############################

agent_zip='{{.AGENT_ZIP}}'

#此参数值被程序使用，修改时请@huanxiong
app_dir='{{.AGENT_DIR}}'
agent_bin='{{.AGENT_BIN}}'
run_user='{{.RUN_USER}}'
data_dir='{{.DATA_DIR}}'
unzip_tmp_dir='/opt/dtstack/tmp'

#安装包下载地址
DOWNLOAD_URL='{{.AGENT_DOWNLOAD_URL}}'

trap '[ "$?" -eq 0 ] || read -p "Looks like something went wrong in step ´$STEP´"' EXIT

##install the filebeat##
install_agent() {
    mkdir -p "$app_dir"
    unzip -o "$unzip_tmp_dir/$agent_zip" -d "$app_dir"  >/dev/null 2>&1
}

##download and installed##
install(){
    STEP='install agent'
    echo "Use the curl download and install Please Waiting..."
    if [ ! -d "$unzip_tmp_dir" ];then
        mkdir -p $unzip_tmp_dir
    fi

    cd "$unzip_tmp_dir" && curl -L -O -s "$DOWNLOAD_URL"
    install_agent

    if [ ! -f "$agent_bin" ];then
        echo "cmd: $agent_bin not found!"
        exit 1
    fi
}

##chown##
chowns(){
    STEP='chown'
    if [ -n "$run_user" ];then
        if ! id -u $run_user >/dev/null 2>&1; then
          sudo useradd $run_user
          sudo echo "$run_user ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
        fi
        sudo chown -R "$run_user:$run_user" "$app_dir"
        if [ -n "$data_dir" ];then
            for path in `echo $data_dir`
            do
                sudo mkdir -p $path
                sudo chown -R "$run_user:$run_user" $path
            done
        fi
    else
        if [ -n "$data_dir" ];then
            for path in `echo $data_dir`
            do
                sudo mkdir -p $path
            done
        fi
    fi
}

##delete filebeat pkg##
delete(){
    STEP='delete'
    cd "$unzip_tmp_dir" && sudo rm -f "$agent_zip"
}

install
chowns
delete
