#!/bin/env bash

#################################
#Script: mysql_exec.sh
#Author: 亚东(lvyadong@dtstack.com)
#Data: 2021-06-30
#Desciption: 完成数栈产品升级mysql数据库的操作，如执行sql文件，备份库，回滚库
#Usage: mysql_exec [--dump database] [--exec filename --database db_name ] [--rollback dumpfile --database db_name ] COMMAND
#################################

#定义args
model=$1

#定义datetime变量
log_date_hms=$(date +"%Y-%m-%d-%H:%M:%S")
log_date=$(date +"%Y-%m-%d")
datetime=$(date +%Y%m%d_%H%M%S_%N | cut -b1-20)
date=$(date +%Y%m%d)
#日志文件
log_directory=/data/my3306/exec_log/log_${date}
dump_directory=/data/my3306/dump

#定义日志输出
function logger() {
    echo -e "${log_date_hms}:\e[32m${1}\e[0m" >> ${log_directory}/log_${model}_${log_date}.log
}

#定义mysql_home
MYSQL_BIN=/opt/dtstack/Mysql/mysql/bin

#mysql_home/bin/mysql
if [ -d ${MYSQL_BIN} ]; then
	MYSQL="sudo ${MYSQL_BIN}/mysql -S /data/my3306/run/mysql.sock -vvv --show-warnings "
else
	printf "directory does not exist : /opt/dtstack/DTBase/mysql/bin\n"
fi

#mysql_home/bin/mysqldump
if [ -d ${MYSQL_BIN} ]; then
	MYSQLDUMP="sudo ${MYSQL_BIN}/mysqldump -S /data/my3306/run/mysql.sock --add-drop-table --master-data=2 --single-transaction --max_allowed_packet=2G  "
else
	printf "directory does not exist : /opt/dtstack/DTBase/mysql/bin\n"
fi

#创建日志目录
function create_log_directory() {

	#输出datetime
	#创建目录 使用mkdir命令
	#首先判断目录是否存在，如果不存在则创建，存在则不再创建
	if [ ! -d "${log_directory}" ]; then
		#echo "目录不存在"
		mkdir -p ${log_directory}
	fi

}

#创建备份目录
function create_dump_directory() {

        #创建目录 使用mkdir命令
        #首先判断目录是否存在，如果不存在则创建，存在则不再创建
        if [ ! -d "${dump_directory}" ]; then
                #echo "目录不存在"
                mkdir -p ${dump_directory}
        fi

}


create_log_directory

create_dump_directory


#help帮助文档
function print_usage() {
	echo "Usage: mysql_exec [--dump --database db_name] [--exec filename --database db_name ] [--rollback dumpfile --database db_name ] COMMAND"
	echo "       where COMMAND is one of:"
	echo "  dump                 备份mysql数据库,可指定database,默认为all"
	echo "  exec                 执行指定sql文件"
	echo "  rollback             回滚指定数据库，需要指定dump文件"
	echo ""
	echo "大多数命令在不带参数的情况下调用时打印帮助."
}

#检测mysql连接以及是否初始化
check_conection_and_initialization() {
	$MYSQL -e "select * from mysql.user;" > test_conn.tmp
	if [ $? = 0 ]; then
		printf "Connection successful:${MYSQL}\n"
		logger "Connection successful:${MYSQL}"
	else
		printf "\e[31m Connection fail:${MYSQL} \e[0m \n"
		logger "Connection fail:${MYSQL}"
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1
	fi
}

##############################
#Desciption：执行sql文件
#args: $1:sql文件名,$2数据库名
##############################
function implement_sql() {
 # check_conection_and_initialization
	local sql=$1
	local database=$2
	printf "start implement ${sql}\n"
	logger "start implement ${sql}"
	nohup $MYSQL --tee=${log_directory}/log_${model}_${log_date}.log ${database} < ${sql} >> ${log_directory}/log_${model}_${log_date}.log 2>&1
	if [ $? == 0 ]; then
                printf "\e[32m ${sql} implement success \e[0m \n"
		logger "${log_date}:${sql} implement success"
	else
		printf "\e[31m ${sql} implement faild!!!!! \e[0m \n"
                printf "\e[31m 请查看日志文件:${log_directory}/log_${model}_${log_date}.log \e[0m \n"
		logger "${sql} implement faild"
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1
	fi
        
        #printf "MYSQL --tee=${log_directory}/log_${model}_${log_date}.log ${database} < ${sql} >> ${log_directory}/log_${model}_${log_date}.log"

}

##备份库
function dumpDatabase() {
  check_conection_and_initialization
	local database=$1
	#sudo mkdir -p $sqlbak
	logger "start dump $db"
	printf "start dump\n"
	$MYSQLDUMP $database > ${dump_directory}/${database}_${log_date_hms}.dump
	if [ $? == 0 ]; then
		logger "dump $database complete"
                printf "dump $database complete\n"
                logger "dump file : ${dump_directory}/${database}_${log_date_hms}.dump"
                printf "dump file : ${dump_directory}/${database}_${log_date_hms}.dump\n"
	else
		printf "dump $database faild\n"
		logger "dump $database faild"
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1
	fi
}

##回滚库
function rollbackDatabase() {
  check_conection_and_initialization
	local database=$1
	local dumpfile=$2
	logger "start rollback ${database}"
	printf "start rollback ${database}\n"
        printf "正在回滚，详情请查看: ${log_directory}/log_${model}_${log_date}.log \n"
	nohup $MYSQL $database < ${dumpfile} >> ${log_directory}/log_${model}_${log_date}.log 2>&1
	if [ $? == 0 ]; then
		logger "rollback $dumpfile complete for $database"
	else
		printf "rollback $database faild\n"
		logger "rollback $database faild"
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1
	fi
}

	while true; do
		case "$model" in

		-exec | --exec | exec)
			model=exec
			shift
			break
			;;

		-dump | --dump | dump)
			model=dump
			shift
			break
			;;

		-rollback | --rollback | rollback)
			model=rollback
			shift
			break
			;;

		*)
			printf "Unknown option %s\n" "$1"
			print_usage
logger "------------------------------------------华丽的分割线---------------------------------------------"
			exit 1
			;;
		esac
	done

#printf "model is : ${model}\n"

if [ "$model" == "dump" ]; then

	if [ "$1" == "--database" ] || [ "$1" == "-database" ] || [ "$1" == "database" ]; then

		if [ "$2" == "" ]; then

			printf "未指定库名\n"
		else
                        printf "model  : ${model}\n"
                        printf "database  : $2\n"
                        logger "model : ${model}"
                        logger "database  : $2"
			printf "start dump\n"
			sleep 2s
			dumpDatabase $2
			printf "end dump\n"
		fi

	else

		printf "Unknown option : ${1} \n"
                print_usage
		logger "Unknown option : ${1} "
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1

	fi

fi

if [ "$model" == "exec" ]; then

	if [[ ! $1 =~ (.*)sql ]]; then

		printf "指定的sq文件无效，请指定后缀为.sql的文件\n"
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1

	fi

	if [ "$2" == "--database" ] || [ "$2" == "-database" ] || [ "$2" == "database" ]; then

		if [ "$3" == "" ]; then

			printf "未指定库名\n"
		else
                        printf "model : ${model}\n"
                        printf "database : $3\n"
                        logger "model : ${model}"
                        logger "database : $3"
			printf "start exec\n"
			implement_sql $1 $3
			printf "end exec\n"
		fi

	else

		printf "Unknown option : ${2} \n"
                print_usage
                logger "Unknown option : ${2} "
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1

	fi
fi

if [ "$model" == "rollback" ]; then

	if [[ ! $1 =~ (.*)dump ]]; then

		printf "指定的sq文件无效，请指定后缀为.dump的文件:${1}\n"
                logger "指定的sq文件无效，请指定后缀为.dump的文件:${1}"
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1

	fi

	if [ "$2" == "--database" ] || [ "$2" == "-database" ] || [ "$2" == "database" ]; then

		if [ "$3" == "" ]; then

			printf "未指定库名\n"
		else

                        read -p "请确认回滚操作[yes/no]" node
                        if [ "$node" != "yes" ]; then
                      
                        printf "放弃回滚，操作结束\n"
                        logger "放弃回滚，操作结束"
logger "------------------------------------------华丽的分割线---------------------------------------------"
                        exit 1
                        fi
                        printf "model  : ${model}\n"
                        printf "database  : $3\n"
                        logger "model  : ${model}"
                        logger "database  : $3"
			printf "start rollback: ${1} to ${3}\n"
                        logger "start rollback: ${1} to ${3}"
                        rollbackDatabase $3 $1
			printf "end rollback\n"
		fi

	else

		printf "Unknown option : ${2} \n"
                print_usage
		logger "Unknown option : ${2} "
logger "------------------------------------------华丽的分割线---------------------------------------------"
		exit 1

	fi
fi

logger "------------------------------------------华丽的分割线---------------------------------------------"
printf "logfile : /data/my3306/exec_log/log_${date}/log_${model}_${log_date}.log\n"
