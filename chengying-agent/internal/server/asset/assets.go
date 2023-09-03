// Code generated for package asset by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/easyagent_install_4win.ps1
// templates/install.script.wrapper.sh
// templates/install.sidecar.sh
// templates/install.sidecar.win.ps1
package asset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesEasyagent_install_4winPs1 = []byte(`#管控server ip
$SERVER_IP_ADDRESS='{{.SERVER_IP_ADDRESS}}'
#管控server port
$SERVER_PORT=9890
#安装包下载host#ip:port
$EASYAGENT_DOWNLOAD_HOST='{{.callback}}'

$UUID='{{.UUID}}'
$CALLBACK_TO_PROD='{{.CALLBACK_TO_PROD}}'
$SIDECAR_INSTALL_PATH='{{.SIDECAR_INSTALL_PATH}}'
$DEBUG_MODE='{{.DEBUG_MODE}}'
$TARGET_OS='{{.TARGET_OS}}'

$SIDECAR_DIR="c:\dtstack\easyagent"

If(![String]::IsNullOrEmpty($SIDECAR_INSTALL_PATH))
{
	$SIDECAR_DIR=$SIDECAR_INSTALL_PATH
}

$SIDECAR_DOWNLOAD_URL="http://$EASYAGENT_DOWNLOAD_HOST/easyagent/easyagent_for_win_x86/easyagent.zip"
$INSTALL_ZIP="easyagent.zip"
$SIDECAR_APP_DIR="easyagent"

# Is this a 64 bit process
function Test-Win64() {
    return [IntPtr]::size -eq 8
}

# Is this a 32 bit process
function Test-Win32() {
    return [IntPtr]::size -eq 4
}

if(Test-Win64)
{
    $SIDECAR_DOWNLOAD_URL="http://$EASYAGENT_DOWNLOAD_HOST/easyagent/easyagent_for_win_x64/easyagent.zip"
}

$INSTALL_TYPE='sidecar'
$CHECK_FLG='{{.create_time}}'
$CALLBACK_URL="http://$EASYAGENT_DOWNLOAD_HOST/api/v1/deploy/sidecar/install/callback?check_flg=$CHECK_FLG&install_type=$INSTALL_TYPE&client_id=$UUID"

# success or failed
$INSTALL_RES=''

# install comount msg
$MSG=''


Function DownloadFile([string]$url, [string]$dest)
{
    Write-Host "-->Download File: $url, save to $dest"
	$client = new-object System.Net.WebClient
	$client.DownloadFile($url, $dest)
}

Function Clean_Install_Pkg
{
	If(Test-Path "$SIDECAR_DIR\$INSTALL_ZIP")
	{
		Write-Host "-->clean $SIDECAR_DIR\$INSTALL_ZIP "
		del "$SIDECAR_DIR\$INSTALL_ZIP"
	}
}

# Convert string to base64 string
function ConvertTo-Base64String([string]$str)
{
    $byteArray = [System.Text.UnicodeEncoding]::Unicode.GetBytes($str)
    [Convert]::ToBase64String( $byteArray )
}

# Convert base64 string to string
function ConvertFrom-Base64String([string]$string)
{
    $byteArray = [Convert]::FromBase64String($string)
    [System.Text.UnicodeEncoding]::Unicode.GetString($byteArray)
}


Function Unzip-File
{
    param([string]$ZipFile,[string]$TargetFolder)
    if(!(Test-Path $TargetFolder))
    {
		mkdir $TargetFolder
    }
    $shellApp = New-Object -ComObject Shell.Application
    $files = $shellApp.NameSpace($ZipFile).Items()
	$files|%{Remove-Item ("$TargetFolder/{0}" -f  $_.name ) -Force -Recurse -ErrorAction SilentlyContinue}
    $shellApp.NameSpace($TargetFolder).CopyHere($files)
}

Function Install_Sidecar
{
	Write-Host '-->Install easyagent sidecar...'
	if(!(Test-Path $SIDECAR_DIR))
	{
		md $SIDECAR_DIR
	}
    cd "$SIDECAR_DIR"

	If(!(Test-Path "$SIDECAR_DIR\$INSTALL_ZIP"))
	{
		Write-Host "[+]download $INSTALL_ZIP from $SIDECAR_DOWNLOAD_URL"
		DownloadFile -url "$SIDECAR_DOWNLOAD_URL" -dest "$SIDECAR_DIR\$INSTALL_ZIP"
	}
	If(!(Test-Path "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent-sidecar.exe"))
	{
	    Write-Host "-->unzip file $SIDECAR_DIR\$INSTALL_ZIP"
	    Unzip-File -ZipFile "$SIDECAR_DIR\$INSTALL_ZIP" -TargetFolder "$SIDECAR_DIR"
	}
	If(!(Test-Path "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent-sidecar.exe"))
	{
		Write-Host "[+]unzip file $SIDECAR_DIR\$INSTALL_ZIP failed"
		Clean_Install_Pkg
		$MSG = ConvertTo-Base64String -str "sidecar agent unzip file $SIDECAR_DIR\$INSTALL_ZIP failed"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
		exit 1
	}
}

Function Config_Sidecar
{
    Write-Host "-->config EasyManage agent (easyagent sidecar)..."

	if((Test-Path $SIDECAR_DIR) -and (Test-Path $SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent.yml))
	{
		$config = "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent.yml"
		$is_in = select-string $config -pattern 'Callback='
		if(!$is_in)
		{
			(type $config) -replace ('\A\s{1,10}uuid:.*$',"  uuid: $UUID")|out-file $config
		}
		Else
		{
			Write-Host '[+] easyagent-sidecar is already installed...'
		}
		if($CALLBACK_TO_PROD)
		{
			Write-Host '[+]Config CALLBACK_TO_PROD'
			(type $config) -replace ('\A.{1,100}CallBack=.*$',"")|out-file $config
			"  - $CALLBACK_TO_PROD"|Out-File -Append $config
		}
		if($SIDECAR_DIR)
		{
			Write-Host '[+]Config SIDECAR_INSTALL_PATH'
			(type $config) -replace ('\A.{1,10}dir:.*$',"  dir: $SIDECAR_DIR\$SIDECAR_APP_DIR")|out-file $config
		}
		if($SERVER_IP_ADDRESS)
		{
			Write-Host '[+]Config SERVER_IP_ADDRESS'
			(type $config) -replace ('\A.{1,10}server:.*$',"  server: $SERVER_IP_ADDRESS")|out-file $config
		}
		if($SERVER_PORT)
		{
			Write-Host '[+]Config SERVER_PORT'
			(type $config) -replace ('\A.{1,10}port:.*$',"  port: $SERVER_PORT")|out-file $config
		}
	}
}

#----------------------------begin start servcie------------------------
function StartTheService()
{
        param([string]$service_name)
		$Service = Get-Service $service_name
		$status = $Service.Status
		Write-Host "service $service_name status $status"
        if($Service.Status -eq "Running")
        {
                Write-Host $Service "The service the running now,will restart the service"
                Write-Host $Service "Starting service..."
                Restart-Service -Name $service_name
        }
        elseif($Service.Status -eq "Disabled")
        {
                #Set the status of service, cause the cmd Start-Service can not start the service if the status is disabled
                Set-Service -Name $service_name -StartupType automatic Manual
                Write-Host $status "Set the status to manual successful"
                Write-Host $Service "Starting service..."
                Start-Service -Name $service_name
        }
        else
        {
                Write-Host $Service.Name "The service the stopped now,will start the service"
                Write-Host $Service.Name "Starting service..."
                Start-Service -Name $service_name
                $Service = Get-Service $service_name
                $step = 1
                do
                {
                    $step = $step + 1
                    Write-Host $service_name "Starting service..."
                    Start-Service -Name $service_name
                    $Service = Get-Service $service_name
                }while($Service.Status -ne "Running" -and $step -le 3)
        }
}
#----------------------------end start servcie------------------------
#---------------------------------------------------------------------
#----------------------------begin stop servcie-----------------------
function StopTheService()
{
		param([string]$service_name)
		$Service = Get-Service $service_name
        if($Service.Status -eq "Running")
        {
                do
                {
                        Write-Host $Service "Stopping service..."
                        Stop-Service -Name $Service.Name
                }while($_.Status -eq "Running") #need to wait until it stopped
                if($_.Status -eq "Running")
                {
                        # stop the servcie failed,send the error mail
                        #sendMail here

                }
                else
                {
                        Write-Host $Service "Stop the servcie successful"
                }
        }
        else
        {
                Write-Host $Service.Status "The service is already stopped now"
        }
}

Function Install_Service
{
	$service_name = "dtstackeasyagent"
	$workdir = "$SIDECAR_DIR\$SIDECAR_APP_DIR"
	# create new service
	if(Check_Service)
	{
	    Write-Host $service_name "service already exist!"
	}
	else
	{
	    New-Service -name $service_name ` + "`" + `
        		-displayName $service_name ` + "`" + `
        		-Description "dtstack easyagent service" ` + "`" + `
        		-StartupType Automatic ` + "`" + `
        		-binaryPathName "` + "`" + `"$workdir\easyagent-sidecar.exe` + "`" + `" -c ` + "`" + `"$workdir\easyagent.yml` + "`" + `" --agents ` + "`" + `"$workdir\agents-file.yml` + "`" + `""
	}

	$Service = Get-Service $service_name
	if([String]::IsNullOrEmpty($Service))
	{
        Write-Host "New service for easyagent failed, service_name: $service_name"
		$MSG = ConvertTo-Base64String -str "New service for easyagent failed, service_name: $service_name"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
		exit 1
	}
	Write-Host "The dtstackeasyagent starting!"
	Start-Sleep -Seconds 10
	StartTheService($service_name)
	$service_check = Get-Service $service_name
	if($service_check.Status -eq "Running")
    {
        Write-Host $service_check "The easyagent install and start success!"
		$MSG = ConvertTo-Base64String -str "The easyagent install and start success!"
        $INSTALL_RES='success'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
    }
	Else
	{
		Write-Host $service_check "The easyagent servcie start error!"
		$MSG = ConvertTo-Base64String -str "The easyagent servcie start error!"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
	}
}

Function CheckAdmin
{
	$currentWi = [Security.Principal.WindowsIdentity]::GetCurrent()
    $currentWp = [Security.Principal.WindowsPrincipal]$currentWi
	if( -not $currentWp.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator))
	{
		return $false
	}
	return $true
}

Function Check_Service
{
    $service_name = "dtstackeasyagent"
    # delete service if it already exists = Get-Service $service_name
    $Service = Get-Service $service_name -ErrorAction SilentlyContinue
    if(![String]::IsNullOrEmpty($Service))
    {
        return $true
    }
    return $false
}

#if(!CheckAdmin)
#{
#    Write-Host "Need administrator privileges to run this install script!"
#    exit 1;
#}

Check_Service
Install_Sidecar
Config_Sidecar
Install_Service`)

func templatesEasyagent_install_4winPs1Bytes() ([]byte, error) {
	return _templatesEasyagent_install_4winPs1, nil
}

func templatesEasyagent_install_4winPs1() (*asset, error) {
	bytes, err := templatesEasyagent_install_4winPs1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/easyagent_install_4win.ps1", size: 9674, mode: os.FileMode(420), modTime: time.Unix(1669097377, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesInstallScriptWrapperSh = []byte(`#!/bin/sh
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
`)

func templatesInstallScriptWrapperShBytes() ([]byte, error) {
	return _templatesInstallScriptWrapperSh, nil
}

func templatesInstallScriptWrapperSh() (*asset, error) {
	bytes, err := templatesInstallScriptWrapperShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/install.script.wrapper.sh", size: 1342, mode: os.FileMode(420), modTime: time.Unix(1669097377, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesInstallSidecarSh = []byte(`#!/bin/sh
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
    sudo chown -R ` + "`" + `whoami` + "`" + `:` + "`" + `whoami` + "`" + ` ${SIDECAR_DIR}

    # chown for product directory
    sudo chown -R ` + "`" + `whoami` + "`" + `:` + "`" + `whoami` + "`" + ` /opt/dtstack

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
        installed=` + "`" + `grep "CallBack=" ${SIDECAR_DIR}/sidecar.yml` + "`" + `
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

CMD_PATH=` + "`" + `dirname $0` + "`" + `
CMD_HOME=` + "`" + `cd "$CMD_PATH"/; pwd` + "`" + `
APP_HOME=` + "`" + `dirname $CMD_HOME` + "`" + `

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
    pid=` + "`" + `cat "$AGENT_PIDFILE"` + "`" + `
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
    pid=` + "`" + `cat "$AGENT_PIDFILE"` + "`" + `
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

    app_pid=` + "`" + `cat "$AGENT_PIDFILE" 2>/dev/null` + "`" + `
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

CMD_PATH=` + "`" + `dirname $0` + "`" + `
CMD_HOME=` + "`" + `cd "$CMD_PATH"/; pwd` + "`" + `
APP_HOME=` + "`" + `dirname $CMD_HOME` + "`" + `
APP_NAME=easyagent-sidecar

curtime()
{
    date '+%Y-%m-%d %H:%M:%S'
}

start() {
check=` + "`" + `ps -ef | grep $APP_NAME | grep -v grep` + "`" + `

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
    is_in=` + "`" + `grep "RUNMODE" ${SIDECAR_DIR}/easyagent.sh` + "`" + `
    if [ ! -z "$is_in" ];then
        sed -i "s/RUNMODE/$DEBUG_MODE/" ${SIDECAR_DIR}/easyagent.sh
    fi
}

AddCrontab() {
    crontab -l > $SIDECAR_INSTALL_PATH/crontab.bak
    is_in=` + "`" + `grep cron_easyagent.sh $SIDECAR_INSTALL_PATH/crontab.bak` + "`" + `
    if [ -z "$is_in" ];then
        echo "*/3 * * * * sh $SIDECAR_INSTALL_PATH/cron_easyagent.sh >> $SIDECAR_INSTALL_PATH/dtcrontab.log 2>&1"  >> $SIDECAR_INSTALL_PATH/crontab.bak
        crontab $SIDECAR_INSTALL_PATH/crontab.bak
    fi
}

Start_Sidecar(){
    STEP='Start_Sidecar'
#    if [ -f /etc/rc.d/rc.local  ];then
#        is_in=` + "`" + `grep ${SIDECAR_DIR}/easyagent.sh /etc/rc.d/rc.local` + "`" + `
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
#            is_in=` + "`" + `grep ${SIDECAR_DIR}/easyagent.sh /etc/rc.local` + "`" + `
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
`)

func templatesInstallSidecarShBytes() ([]byte, error) {
	return _templatesInstallSidecarSh, nil
}

func templatesInstallSidecarSh() (*asset, error) {
	bytes, err := templatesInstallSidecarShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/install.sidecar.sh", size: 11947, mode: os.FileMode(420), modTime: time.Unix(1669097377, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesInstallSidecarWinPs1 = []byte(`#管控server ip
$SERVER_IP_ADDRESS='{{.SERVER_IP_ADDRESS}}'
#管控server port
$SERVER_PORT=8890
#安装包下载host#ip:port
$EASYAGENT_DOWNLOAD_HOST='{{.callback}}'

$UUID='{{.UUID}}'
$CALLBACK_TO_PROD='{{.CALLBACK_TO_PROD}}'
$SIDECAR_INSTALL_PATH='{{.SIDECAR_INSTALL_PATH}}'
$DEBUG_MODE='{{.DEBUG_MODE}}'
$TARGET_OS='{{.TARGET_OS}}'

$SIDECAR_DIR="d:\dtstack\easyagent"

If(![String]::IsNullOrEmpty($SIDECAR_INSTALL_PATH))
{
	$SIDECAR_DIR=$SIDECAR_INSTALL_PATH
}

$SIDECAR_DOWNLOAD_URL="http://$EASYAGENT_DOWNLOAD_HOST/easyagent/easyagent_for_win_x86/easyagent.zip"
$INSTALL_ZIP="easyagent.zip"
$SIDECAR_APP_DIR="easyagent"

# Is this a 64 bit process
function Test-Win64() {
    return [IntPtr]::size -eq 8
}

# Is this a 32 bit process
function Test-Win32() {
    return [IntPtr]::size -eq 4
}

if(Test-Win64)
{
    $SIDECAR_DOWNLOAD_URL="http://$EASYAGENT_DOWNLOAD_HOST/easyagent/easyagent_for_win_x64/easyagent.zip"
}

$INSTALL_TYPE='sidecar'
$CHECK_FLG='{{.create_time}}'
$CALLBACK_URL="http://$EASYAGENT_DOWNLOAD_HOST/api/v1/deploy/sidecar/install/callback?check_flg=$CHECK_FLG&install_type=$INSTALL_TYPE&client_id=$UUID"

# success or failed
$INSTALL_RES=''

# install comount msg
$MSG=''


Function DownloadFile([string]$url, [string]$dest)
{
	$client = new-object System.Net.WebClient
	$client.DownloadFile($url, $dest)
}

Function Clean_Install_Pkg
{
	If(Test-Path "$SIDECAR_DIR\$INSTALL_ZIP")
	{
		Write-Host "-->clean $SIDECAR_DIR\$INSTALL_ZIP "
		del "$SIDECAR_DIR\$INSTALL_ZIP"
	}
}

# Convert string to base64 string
function ConvertTo-Base64String([string]$str)
{
    $byteArray = [System.Text.UnicodeEncoding]::Unicode.GetBytes($str)
    [Convert]::ToBase64String( $byteArray )
}

# Convert base64 string to string
function ConvertFrom-Base64String([string]$string)
{
    $byteArray = [Convert]::FromBase64String($string)
    [System.Text.UnicodeEncoding]::Unicode.GetString($byteArray)
}


Function Unzip-File
{
    param([string]$ZipFile,[string]$TargetFolder)
    if(!(Test-Path $TargetFolder))
    {
		mkdir $TargetFolder
    }
    $shellApp = New-Object -ComObject Shell.Application
    $files = $shellApp.NameSpace($ZipFile).Items()
	$files|%{Remove-Item ("$TargetFolder/{0}" -f  $_.name ) -Force -Recurse -ErrorAction SilentlyContinue}
    $shellApp.NameSpace($TargetFolder).CopyHere($files)
}

Function Install_Sidecar
{
	Write-Host '-->Install easyagent sidecar...'
	if(!(Test-Path $SIDECAR_DIR))
	{
		md $SIDECAR_DIR
	}
    cd "$SIDECAR_DIR"

	If(!(Test-Path "$SIDECAR_DIR\$INSTALL_ZIP"))
	{
		Write-Host "[+]download $INSTALL_ZIP from $SIDECAR_DOWNLOAD_URL"
		DownloadFile -url "$SIDECAR_DOWNLOAD_URL" -dest "$SIDECAR_DIR\$INSTALL_ZIP"
	}
	Write-Host "-->unzip file $SIDECAR_DIR\$INSTALL_ZIP"
	Unzip-File -ZipFile "$SIDECAR_DIR\$INSTALL_ZIP" -TargetFolder "$SIDECAR_DIR"
	If(!(Test-Path "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent-sidecar.exe"))
	{
		Write-Host "[+]unzip file $SIDECAR_DIR\$INSTALL_ZIP failed"
		Clean_Install_Pkg
		$MSG = ConvertTo-Base64String -str "sidecar agent unzip file $SIDECAR_DIR\$INSTALL_ZIP failed"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
		exit 1
	}
}

Function Config_Sidecar
{
    Write-Host "-->config EasyManage agent (easyagent sidecar)..."

	if((Test-Path $SIDECAR_DIR) -and (Test-Path $SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent.yml))
	{
		$config = "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent.yml"
		$is_in = select-string $config -pattern 'Callback='
		if(!$is_in)
		{
			(type $config) -replace ('\A\s{1,10}uuid:.*$',"  uuid: $UUID")|out-file $config
		}
		Else
		{
			Write-Host '[+] easyagent-sidecar is already installed...'
		}
		if($CALLBACK_TO_PROD)
		{
			Write-Host '[+]Config CALLBACK_TO_PROD'
			(type $config) -replace ('\A.{1,100}CallBack=.*$',"")|out-file $config
			"  - $CALLBACK_TO_PROD"|Out-File -Append $config
		}
		if($SIDECAR_DIR)
		{
			Write-Host '[+]Config SIDECAR_INSTALL_PATH'
			(type $config) -replace ('\A.{1,10}dir:.*$',"  dir: $SIDECAR_DIR\$SIDECAR_APP_DIR")|out-file $config
		}
		if($SERVER_IP_ADDRESS)
		{
			Write-Host '[+]Config SERVER_IP_ADDRESS'
			(type $config) -replace ('\A.{1,10}server:.*$',"  server: $SERVER_IP_ADDRESS")|out-file $config
		}
		if($SERVER_PORT)
		{
			Write-Host '[+]Config SERVER_PORT'
			(type $config) -replace ('\A.{1,10}port:.*$',"  port: $SERVER_PORT")|out-file $config
		}
	}
}

#----------------------------begin start servcie------------------------
function StartTheService()
{
        param([string]$service_name)
		$Service = Get-Service $service_name
		$status = $Service.Status
		Write-Host "service $service_name status $status"
        if($Service.Status -eq "Running")
        {
                Write-Host $Service "The service the running now,will restart the service"
                Write-Host $Service "Starting service..."
                Restart-Service -Name $service_name
        }
        elseif($Service.Status -eq "Disabled")
        {
                #Set the status of service, cause the cmd Start-Service can not start the service if the status is disabled
                Set-Service -Name $service_name -StartupType automatic Manual
                Write-Host $status "Set the status to manual successful"
                Write-Host $Service "Starting service..."
                Start-Service -Name $service_name
        }
        else
        {
                Write-Host $Service.Name "The service the stopped now,will start the service"
                Write-Host $Service "Starting service..."
                Start-Service -Name $service_name
        }
}
#----------------------------end start servcie------------------------
#---------------------------------------------------------------------
#----------------------------begin stop servcie-----------------------
function StopTheService()
{
		param([string]$service_name)
		$Service = Get-Service $service_name
        if($Service.Status -eq "Running")
        {
                do
                {
                        Write-Host $Service "Stopping service..."
                        Stop-Service -Name $Service.Name
                }while($_.Status -eq "Running") #need to wait until it stopped
                if($_.Status -eq "Running")
                {
                        # stop the servcie failed,send the error mail
                        #sendMail here

                }
                else
                {
                        Write-Host $Service "Stop the servcie successful"
                }
        }
        else
        {
                Write-Host $Service.Status "The service is already stopped now"
        }
}

Function Install_Service
{
	$service_name = "dtstackeasyagent"
	$workdir = "$SIDECAR_DIR\$SIDECAR_APP_DIR"
	# create new service
	New-Service -name $service_name ` + "`" + `
		-displayName $service_name ` + "`" + `
		-binaryPathName "` + "`" + `"$workdir\easyagent-sidecar.exe` + "`" + `" -c ` + "`" + `"$workdir\easyagent.yml` + "`" + `""
	#start new service
	$Service = Get-Service $service_name
	if([String]::IsNullOrEmpty($Service))
	{
        Write-Host "New service for easyagent failed, service_name: $service_name"
		$MSG = ConvertTo-Base64String -str "New service for easyagent failed, service_name: $service_name"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
		exit 1
	}
	StartTheService($service_name)
	$service_check = Get-Service $service_name
	if($service_check.Status -eq "Running")
    {
        Write-Host $service_check "The easyagent install and start success!"
		$MSG = ConvertTo-Base64String -str "The easyagent install and start success!"
        $INSTALL_RES='success'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
    }
	Else
	{
		Write-Host $service_check "The easyagent servcie start error!"
		$MSG = ConvertTo-Base64String -str "The easyagent servcie start error!"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
	}
}

Function CheckAdmin
{
	$currentWi = [Security.Principal.WindowsIdentity]::GetCurrent()
    $currentWp = [Security.Principal.WindowsPrincipal]$currentWi
	if( -not $currentWp.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator))
	{
		return $false
	}
	return $true
}

Function Check_Service
{
    $service_name = "dtstackeasyagent"
    # delete service if it already exists = Get-Service $service_name
    $Service = Get-Service $service_name -ErrorAction SilentlyContinue
    if(![String]::IsNullOrEmpty($Service))
    {
        Write-Host "dtstackeasyagent servcie already exists, please run sc delete $servcie_name to delete first"
        exit 1;
    }
}

#if(!CheckAdmin)
#{
#    Write-Host "Need administrator privileges to run this install script!"
#    exit 1;
#}

Check_Service
Install_Sidecar
Config_Sidecar
Install_Service`)

func templatesInstallSidecarWinPs1Bytes() ([]byte, error) {
	return _templatesInstallSidecarWinPs1, nil
}

func templatesInstallSidecarWinPs1() (*asset, error) {
	bytes, err := templatesInstallSidecarWinPs1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/install.sidecar.win.ps1", size: 8926, mode: os.FileMode(420), modTime: time.Unix(1669097377, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/easyagent_install_4win.ps1": templatesEasyagent_install_4winPs1,
	"templates/install.script.wrapper.sh":  templatesInstallScriptWrapperSh,
	"templates/install.sidecar.sh":         templatesInstallSidecarSh,
	"templates/install.sidecar.win.ps1":    templatesInstallSidecarWinPs1,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"easyagent_install_4win.ps1": &bintree{templatesEasyagent_install_4winPs1, map[string]*bintree{}},
		"install.script.wrapper.sh":  &bintree{templatesInstallScriptWrapperSh, map[string]*bintree{}},
		"install.sidecar.sh":         &bintree{templatesInstallSidecarSh, map[string]*bintree{}},
		"install.sidecar.win.ps1":    &bintree{templatesInstallSidecarWinPs1, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
