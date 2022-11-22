#!/bin/bash

source ~/.bash_profile

#MYSQL_BIN=/opt/dtstack/DTBase/mysql/bin
MYSQL_BIN={{.mysql_path}}
database=$4


##create database
function create_database() {

 if [ ! $($MYSQL -e "use ${database}") ]; then

		printf "create database "${database}"\n"
    $MYSQL -e "CREATE DATABASE IF NOT EXISTS ${database} DEFAULT CHARSET latin1;"
		if [ $? != 0 ]; then
			printf "database create faild:${database}!\n"
			exit 1

		fi

	fi

}

#执行sql文件，执行失败即退出，返回错误：1
function import_sql() {

	tables=$($MYSQL -e "show tables from ${database}")

  #判断是否是第一次部署
	if [ "$tables" == "" ]; then

printf "init ${database}\n"
$MYSQL -P3306 metastore < sql/hive-schema-2.3.0.mysql-add-txn.sql
$MYSQL -P3306 metastore < sql/hive-fixed-2.3.0.sql
	else
		printf "${database} has been initialized\n"
		exit 0
	fi
}

MYSQL="${MYSQL_BIN}/mysql -h$3 -u$1 -p$2"

create_database
import_sql
