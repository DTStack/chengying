#!/bin/bash


#################################
#Script: post_deploy.sh
#Author: 亚东(lvyadong@dtstack.com)
#Data: 2021-06-01
#Desciption: 1.完成mysql的安装操作
#Usage: ./post_deploy.sh ${user} ${password} ${repl_user} ${repl_password} ${check_num}
#################################

cur=$(pwd)
case "$(uname)" in
Linux)
	bin_path=$(readlink -f $(dirname $0))
	;;
*)
	bin_path=$(
		cd $(dirname $0)
		pwd
	)
	;;
esac

#定义脚本变量

user=$1
password=$2
repl_user=$3
repl_password=$4
#默认值为3次，一次等于10s，可根据真实情况调整
check_num=$5
is_other=$6
BASE=$(dirname $bin_path)/mysql
conf=${BASE}/my.cnf
log=${BASE}/conf/logback.xml
datadir=/data/my3306/
MYSQL="./bin/mysql -uroot -S /data/my3306/run/mysql.sock --connect_timeout=5"
logfile=$datadir/log/mysql.log
errorlog=$datadir/log/error.log

export LANG=en_US.UTF-8
export BASE=$BASE

#初始化操作
install_mysql() {
	who=$(whoami)
	if [ ! -d $datadir ]; then
		mkdir -p $datadir/{run,data,log,tmp}
		bin/mysqld --defaults-file=$conf --initialize-insecure --console --user=$who
		if [ $? -ne 0 ]; then
			echo "command failed"
			exit 1
		fi
	fi
}

#mysql检查是否初始化 and 启动成功
checki=1
check_connection_and_initialization() {
        sleep 10
        $MYSQL -e "select * from mysql.user;"
        status=$?
        #判断mysql连接失败 and 不超过自定义检测次数
        if [[ $status = 1 ]] && [[ ! $checki = $check_num ]]; then
                let checki=$checki+1
                #递归
                check_connection_and_initialization
         #判断mysql连接失败 and 等于自定义检测次数
        elif [[ $status = 1 ]] && [ $checki = $check_num ]; then
                echo "mysql check faild..."
                stop_mysql
                exit 1
        fi
}

#启动mysql守护进程
start_mysql() {
  stop_mysql
	./bin/mysqld_safe --defaults-file=$conf --user=$who &
}


#初始化mysql用户
init_mysql() {
	$MYSQL -e "grant all privileges on *.* to '$user'@'%' identified by '$password';"
	$MYSQL -e "grant all privileges on *.* to monitor@'%' identified by 'Abc123Admini';"
	$MYSQL -e "grant REPLICATION SLAVE,REPLICATION CLIENT on *.* to '$repl_user'@'%' identified by '$repl_password';"
	if [ $? -ne 0 ]; then
		echo "grant privileges failed"
		exit 1
	fi
	$MYSQL -e "delete from mysql.user where user='';"
	if [ $? -ne 0 ]; then
		echo "delete user failed"
		exit 1
	fi
	$MYSQL -e "flush privileges;"
	if [ $? -ne 0 ]; then
		echo "flush privileges failed"
		exit 1
	fi

	cat <<EOF | tee -a ~/.bash_profile
  export MYSQL_BIN=/opt/dtstack/Mysql/mysql/bin
EOF

}

#停止mysql进程
stop_mysql() {
	pid=$(ps aux | grep "mysqld" | grep -v grep | awk '{print $2}')
	for i in $pid; do
		kill -9 $i
	done
}

if [ "${is_other}" == "false" ]; then

install_mysql
start_mysql
check_connection_and_initialization
init_mysql
stop_mysql

else
  printf "对接第三方mysql,不做安装操作\n"
  echo "对接第三方mysql,不做安装操作" > is_other
  exit 0
fi
