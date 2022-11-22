#!/bin/bash

is_other=$1


if [ "${is_other}" == "false" ]; then

./mysql_exporter/mysqld_exporter --config.my-cnf=./mysql_exporter/my-exporter.cnf &

if [ -f "./success.txt" ]; then
  rm -f ./success.txt
fi
./bin/mysqld_safe --defaults-file=my.cnf --user=`whoami`

else

  printf "对接第三方mysql,不做安装操作\n"
  echo "对接第三方mysql,不做启动操作" > is_other
#  nc -l 3306
  tail -f -
fi