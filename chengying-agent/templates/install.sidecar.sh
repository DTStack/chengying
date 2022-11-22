#!/bin/sh
#-----------------------------------------------#
# Copyright (c) 2017 DTStack Inc.               #
# @desc     auto install easymanage agent       #
# @author   fangyan@dtstack.com                 #
# @ver      0.1                                 #
# @date     2017/8/14                           #
#-----------------------------------------------#
export PATH=/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin

#管控server ip
SERVER_IP_ADDRESS='{{.SERVER_IP_ADDRESS}}'
#管控server port
SERVER_PORT={{.RPC_PORT}}
#安装包下载host#ip:port
EASYAGENT_DOWNLOAD_HOST='{{.callback}}'

UUID='{{.UUID}}'
INSTALL_TYPE='{{.INSTALL_TYPE}}'
CALLBACK_TO_PROD='{{.CALLBACK_TO_PROD}}'
SIDECAR_INSTALL_PATH='{{.SIDECAR_INSTALL_PATH}}'
DEBUG_MODE='{{.DEBUG_MODE}}'
TARGET_OS='{{.TARGET_OS}}'

#安装包下载地址
SIDECAR_DOWNLOAD_URL="http://${EASYAGENT_DOWNLOAD_HOST}/easyagent/easyagent.tar.gz"

SIDECAR_DIR="/opt/dtstack/easyagent"

if [ -n "${SIDECAR_INSTALL_PATH}" ];then
    SIDECAR_DIR="${SIDECAR_INSTALL_PATH}"
else
    SIDECAR_INSTALL_PATH="$SIDECAR_DIR"
fi

#特殊查询符号
CHECK_FLG='{{.create_time}}'

#安装完成的回调
CALLBACK_ADDRESS="http://${EASYAGENT_DOWNLOAD_HOST}/api/v1/deploy/sidecar/install/callback"

CALLBACK_URL="${CALLBACK_ADDRESS}?check_flg=${CHECK_FLG}&install_type=${INSTALL_TYPE}&client_id=${UUID}"

# success or failed
INSTALL_RES=''

# install comount msg
MSG=''

#set -e
#trap '[ "$?" -eq 0 ] || read -p "Looks like something went wrong in step ´$STEP´"' EXIT

Color_Text()
{
  echo -e " \e[0;$2m$1\e[0m"
}

Echo_Red()
{
  echo $(Color_Text "$1" "31")
}

Echo_Green()
{
  echo $(Color_Text "$1" "32")
}

Echo_Yellow()
{
  echo $(Color_Text "$1" "33")
}

Echo_Blue()
{
  echo $(Color_Text "$1" "34")
}

Unset_Requiretty()
{
    if [ $(id -u) != "0" ];then
        echo 'Defaults    !requiretty' | sudo tee /etc/sudoers.d/dtstack
    fi
}

Clean_Install_Pkg()
{
  [ -f /tmp/easyagent.tar.gz ] && rm -f /tmp/easyagent.tar.gz
  [ -d /tmp/easyagent ] && rm -rf /tmp/easyagent
}

Install_Sidecar(){
    STEP='Install_Sidecar'
    #side_car 安装目录SIDECAR_DIR
    Echo_Blue "[+] Installing agent (easyagent sidecar)..."

    sudo mkdir -p ${SIDECAR_DIR}
    sudo chown -R `whoami`:`whoami` ${SIDECAR_DIR}

    # chown for product directory
    sudo chown -R `whoami`:`whoami` /opt/dtstack

    cd /tmp

    if [ ! -f easyagent.tar.gz ];then
        curl -L -O -s "${SIDECAR_DOWNLOAD_URL}"
    fi

    if [ ! -f easyagent.tar.gz  ];then
        Echo_Red "sidecar source file not found ..."
        MSG=$( echo -n "sidecar agent 压缩包下载出错" |base64 )
        INSTALL_RES='failed'
        curl -s -g "${CALLBACK_URL}&msg=${MSG}&install_res=${INSTALL_RES}" >/dev/null

        Clean_Install_Pkg
        exit 1
    fi

    tar -xzf easyagent.tar.gz
    if [ ! -d easyagent ] || [ ! -f easyagent/easyagent-sidecar ];then
        Echo_Red "sidecar source file decompression error ..."
        MSG=$( echo -n "sidecar agent 压缩包解压出错" |base64 )
        INSTALL_RES='failed'
        curl -s -g "${CALLBACK_URL}&msg=${MSG}&install_res=${INSTALL_RES}" >/dev/null

        Clean_Install_Pkg
        exit 1
    fi
    if [ ! -f ${SIDECAR_DIR}/easyagent-sidecar ];then
        mv easyagent/* ${SIDECAR_DIR}/
        Clean_Install_Pkg
    else
        if [ -f ${SIDECAR_DIR}/easyagent.sh ];then
            sh ${SIDECAR_DIR}/easyagent.sh stop
        fi
        cp -f easyagent/easyagent-sidecar ${SIDECAR_DIR}/
        Clean_Install_Pkg
    fi
}

Config_Sidecar(){
    STEP='Config_Sidecar'
    Echo_Blue "[+] config EasyManage agent (easyagent sidecar)..."

    if [ -d ${SIDECAR_DIR} ] && [ -f ${SIDECAR_DIR}/sidecar.yml ];then
        installed=`grep "CallBack=" ${SIDECAR_DIR}/sidecar.yml`
        if [ -z "$installed" ];then
            sed -i "s|uuid:.*|uuid: ${UUID}|" ${SIDECAR_DIR}/sidecar.yml
        else
            Echo_Blue "[+] easyagent-sidecar is already installed..."
        fi
        if [ -n "${CALLBACK_TO_PROD}" ];then
            sed -i "/  - .*CallBack=/"d ${SIDECAR_DIR}/sidecar.yml
            echo "  - ${CALLBACK_TO_PROD}" >> ${SIDECAR_DIR}/sidecar.yml
        fi
        if [ -n "${SIDECAR_INSTALL_PATH}" ];then
            sed -i "s|dir:.*|dir: ${SIDECAR_INSTALL_PATH}/logs|" ${SIDECAR_DIR}/sidecar.yml
        fi
        cp -pfr ${SIDECAR_DIR}/sidecar.yml ${SIDECAR_DIR}/sidecar.yml.ori.${CHECK_FLG}
        sed -i "s|server:.*|server: ${SERVER_IP_ADDRESS}|" ${SIDECAR_DIR}/sidecar.yml
        sed -i "s|port:.*|port: ${SERVER_PORT}|" ${SIDECAR_DIR}/sidecar.yml
    fi
}

Config_Sidecar_Docker(){
    STEP='Config_Sidecar_Docker'
    Echo_Blue "[+] config EasyManage agent (easyagent sidecar)..."

    SIDECAR_DIR=/sidecar
    sed -i "s|uuid:.*|uuid: ${UUID}|" ${SIDECAR_DIR}/conf/sidecar.yml
    if [ -n "${CALLBACK_TO_PROD}" ];then
        sed -i "/  - .*CallBack=/"d ${SIDECAR_DIR}/conf/sidecar.yml
        echo "  - ${CALLBACK_TO_PROD}" >> ${SIDECAR_DIR}/conf/sidecar.yml
    fi
    sed -i "s|server:.*|server: ${SERVER_IP_ADDRESS}|" ${SIDECAR_DIR}/conf/sidecar.yml
    sed -i "s|port:.*|port: ${SERVER_PORT}|" ${SIDECAR_DIR}/conf/sidecar.yml
}

Setting_Sidecar(){
    STEP='Setting_Sidecar'
    Echo_Blue "[+] setting EasyManage agent (easyagent sidecar)..."

cat >${SIDECAR_DIR}/easyagent.sh <<'EOF'
#!/bin/sh
export PATH=/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin

CMD_PATH=`dirname $0`
CMD_HOME=`cd "$CMD_PATH"/; pwd`
APP_HOME=`dirname $CMD_HOME`

AGENT_HOME=$CMD_HOME

APP_NAME=easyagent-sidecar
AGENT_CONF=$AGENT_HOME/sidecar.yml
AGENT_FILE_CONF=$AGENT_HOME/agents-file.yml
AGENT_PIDFILE=$AGENT_HOME/easyagent-sidecar.pid
KILL_ON_STOP_TIMEOUT=1

Color_Text()
{
  echo -e " \e[0;$2m$1\e[0m"
}

Echo_Red()
{
  echo $(Color_Text "$1" "31")
}

Echo_Green()
{
  echo $(Color_Text "$1" "32")
}

Echo_Yellow()
{
  echo $(Color_Text "$1" "33")
}

Echo_Blue()
{
  echo $(Color_Text "$1" "34")
}


check_init() {
  if [ ! -d $APP_HOME ]; then
    Echo_Blue "$APP_HOME is not exists.";
    exit 1
  fi

  if [ ! -f $AGENT_CONF ]; then
    Echo_Blue "$AGENT_CONF is not exists.";
    exit 1
  fi
}

status(){
   if [ -f "$AGENT_PIDFILE" ] ; then
    pid=`cat "$AGENT_PIDFILE"`
    if kill -0 $pid > /dev/null 2> /dev/null ; then
      # process by this pid is running.
      # It may not be our pid, but that's what you get with just pidfiles.
      # TODO(sissel): Check if this process seems to be the same as the one we
      # expect. It'd be nice to use flock here, but flock uses fork, not exec,
      # so it makes it quite awkward to use in this case.
      return 0
    else
      return 2 # program is dead but pid file exists
    fi
  else
    return 3 # program is not running
  fi
}

start(){
    Echo_Blue "Starting $APP_NAME"
    nohup $AGENT_HOME/$APP_NAME -c $AGENT_CONF --ags $AGENT_FILE_CONF RUNMODE 1>$AGENT_HOME/easyagent.info 2>&1  &
    echo $! > $AGENT_PIDFILE
    ret=$?
    if [ $ret -eq 0 ] ;then
        Echo_Blue "started"
    else
        Echo_Yellow "start failed"
    fi
}

stop(){
    Echo_Blue "Stopping $APP_NAME"
    # Try a few times to kill TERM the program
  if status ; then
    pid=`cat "$AGENT_PIDFILE"`
    echo $pid
    echo "Killing $APP_NAME (pid $pid) with SIGTERM"
    kill -TERM $pid
    # Wait for it to exit.
    for i in {1..40}; do
      echo "Waiting $APP_NAME (pid $pid) to die..."
      status || break
      sleep 1
    done
    if status ; then
      if [ "$KILL_ON_STOP_TIMEOUT" == "1" ] ; then
        echo "Timeout reached. Killing $APP_NAME (pid $pid) with SIGKILL. This may result in data loss."
        kill -KILL $pid
        echo "$APP_NAME killed with SIGKILL."
      else
        echo "$APP_NAME stop failed; still running."
        return 1 # stop timed out and not forced
      fi
    else
      echo -n "$APP_NAME stopped "
    fi
    rm -rf $AGENT_PIDFILE
  fi
    ret=$?
    if [ $ret -eq 0 ] ;then
        Echo_Blue "stopped"
    else
        Echo_Yellow "stop failed"
    fi
}

case "$1" in
  start)
    check_init
    status
    code=$?
    if [ $code -eq 0 ]; then
      echo "$APP_NAME is already running "
    else
      start
      code=$?
    fi
    exit $code
    ;;

  stop)
    check_init
    stop
    ;;

  restart)
    check_init

    app_pid=`cat "$AGENT_PIDFILE" 2>/dev/null`
    if [ -n "$app_pid" ] && ps -p $app_pid > /dev/null 2>&1; then
       stop
    fi
     sleep 5
     start
    ;;

  status)
   status
    code=$?
    if [ $code -eq 0 ] ; then
      echo "$APP_NAME is running "
    else
      echo "$APP_NAME is not running "
    fi
    exit $code
   ;;

  *)
    echo "Usage: $0 {start|stop|reload|restart|status}"
    exit 1
    ;;
esac

EOF

cat >${SIDECAR_DIR}/cron_easyagent.sh <<'EOF'
#!/bin/sh

export PATH=/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin

CMD_PATH=`dirname $0`
CMD_HOME=`cd "$CMD_PATH"/; pwd`
APP_HOME=`dirname $CMD_HOME`
APP_NAME=easyagent-sidecar

curtime()
{
    date '+%Y-%m-%d %H:%M:%S'
}

start() {
check=`ps -ef | grep $APP_NAME | grep -v grep`

if [ ! -n "$check" ]; then
 echo 'restart' >> $CMD_HOME/cron_easyagent.log
 echo $(curtime) >> $CMD_HOME/cron_easyagent.log
 sh $CMD_HOME/easyagent.sh restart
fi
}

start

EOF

}

Set_Debug(){
    is_in=`grep "RUNMODE" ${SIDECAR_DIR}/easyagent.sh`
    if [ ! -z "$is_in" ];then
        sed -i "s/RUNMODE/$DEBUG_MODE/" ${SIDECAR_DIR}/easyagent.sh
    fi
}

AddCrontab() {
    crontab -l > $SIDECAR_INSTALL_PATH/crontab.bak
    is_in=`grep cron_easyagent.sh $SIDECAR_INSTALL_PATH/crontab.bak`
    if [ -z "$is_in" ];then
        echo "*/3 * * * * sh $SIDECAR_INSTALL_PATH/cron_easyagent.sh >> $SIDECAR_INSTALL_PATH/dtcrontab.log 2>&1"  >> $SIDECAR_INSTALL_PATH/crontab.bak
        crontab $SIDECAR_INSTALL_PATH/crontab.bak
    fi
}

Start_Sidecar(){
    STEP='Start_Sidecar'
#    if [ -f /etc/rc.d/rc.local  ];then
#        is_in=`grep ${SIDECAR_DIR}/easyagent.sh /etc/rc.d/rc.local`
#        if [ -z "$is_in" ];then
#            echo "sh ${SIDECAR_DIR}/easyagent.sh restart" | sudo tee -a /etc/rc.d/rc.local
#        fi
#        if [ $? -eq 0 ];then
#            Echo_Blue "Add to rc local service success!"
#        else
#            Echo_Red "Add easyagent to rc local service failed!"
#        fi
#    else
#        if [ -f /etc/rc.local ];then
#            is_in=`grep ${SIDECAR_DIR}/easyagent.sh /etc/rc.local`
#            if [ -z "$is_in" ];then
#                echo "sh ${SIDECAR_DIR}/easyagent.sh restart" | sudo tee -a /etc/rc.local
#            fi
#            if [ $? -eq 0 ];then
#                Echo_Blue "Add to rc local service success!"
#            else
#                Echo_Red "Add easyagent to rc local service failed!"
#            fi
#        fi
#    fi

    sh ${SIDECAR_DIR}/easyagent.sh restart

    if [ $? -eq 0 ];then
        Echo_Blue "start easyagent success!  ..."
        if [ -z "$installed" ];then
            MSG=$( echo -n "sidecar agent 安装成功并启动" |base64 )
            INSTALL_RES='success'
            curl -s -g "${CALLBACK_URL}&msg=${MSG}&install_res=${INSTALL_RES}" >/dev/null
            AddCrontab
            exit 0
        fi
    else
        Echo_Red "start easyagent failed!  ..."
        if [ -z "$installed" ];then
            MSG=$( echo -n "sidecar agent 启动失败" |base64 )
            INSTALL_RES='failed'
            curl -s -g "${CALLBACK_URL}&msg=${MSG}&install_res=${INSTALL_RES}" >/dev/null
            exit 1
        fi
    fi
}

CallBack_Sidecar_Docker(){
    STEP='CallBack_Sidecar_Docker'
    Echo_Blue "[+] callback EasyManage agent (easyagent sidecar)..."

    MSG=$( echo -n "sidecar agent 安装成功并启动" |base64 )
    INSTALL_RES='success'
    curl -s -g "${CALLBACK_URL}&msg=${MSG}&install_res=${INSTALL_RES}" >/dev/null
    exit 0
}

if [ "$TARGET_OS" != "docker" ];then
    Unset_Requiretty
    Install_Sidecar
    Config_Sidecar
    Setting_Sidecar
    Set_Debug
    Start_Sidecar
else
    Config_Sidecar_Docker
    CallBack_Sidecar_Docker
fi
