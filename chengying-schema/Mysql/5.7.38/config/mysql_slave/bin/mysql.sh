#!/bin/bash

stop_mysql(){
  pid=`ps aux|grep "mysqld"|grep -v grep|awk '{print $2}'`
  for i in $pid
  do
    kill $i
  done
  sleep 10
}

stop_mysql
./mysql_exporter/mysqld_exporter --config.my-cnf=./mysql_exporter/my-exporter.cnf &
./bin/mysqld_safe --defaults-file=my.cnf --user=`whoami`
