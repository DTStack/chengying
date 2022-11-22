#!/bin/bash

cur=`pwd`
case "`uname`" in
    Linux)
        bin_path=$(readlink -f $(dirname $0))
        ;;
    *)
        bin_path=`cd $(dirname $0); pwd`
        ;;
esac

mkdir -p /data/db_backup

BASE=$(dirname $bin_path)/mysql_slave
conf=${BASE}/my.cnf
log=${BASE}/conf/logback.xml
pfile=$BASE/bin/mysql.pid
datadir=/data/my3306
MYSQL="./bin/mysql -uroot -S /data/my3306/run/mysql.sock"
logfile=$datadir/log/mysql.log
errorlog=$datadir/log/error.log
master_host=$1
master_user=$4
master_password=$5
repl_user=$2
repl_password=$3
master_port=3306
check_num=$6

export LANG=en_US.UTF-8
export BASE=$BASE


install_mysql(){
who=`whoami`

if [ ! -d  $datadir ];then
  mkdir -p $datadir/{run,data,log,tmp}
  bin/mysqld --defaults-file=$conf --initialize-insecure --console --user=$who
  if [ $? -ne 0 ]; then echo "command failed"; exit 1; fi
fi
}

start_mysql(){
stop_mysql
./bin/mysqld_safe --defaults-file=$conf --user=$who &
}


checki=1
check_mysql() {
        sleep 10
        $MYSQL -e "select * from mysql.user;"
        status=$?
        if [[ $status = 1 ]] && [[ ! $checki = $check_num ]]; then
                let checki=$checki+1
                check_mysql
        elif [[ $status = 1 ]] && [ $checki = $check_num ]; then
                echo "mysql check faild..."
                stop_mysql
                exit 1
        fi
}


init_mysql() {
  # 是否已经做过数据同步
$MYSQL -e "select * from metastore.VERSION"
if [ $? -ne 0 ]; then

# 同步数据
./bin/mysqldump --single-transaction --master-data=2 -A -R -E -h$master_host -u$master_user -p$master_password > ${BASE}/all.sql
if [ $? -ne 0 ]; then echo "command failed"; exit 1; fi
log_file=`less ${BASE}/all.sql | grep "MASTER_LOG_FILE='*'" | awk -F ',' '{print $1}'|awk -F "'" '{print $2}'`
log_pos=`less ${BASE}/all.sql | grep  "MASTER_LOG_FILE='*'" |awk -F '=' '{print $NF}'`

# 导入数据
$MYSQL -e "source ${BASE}/all.sql;"
if [ $? -ne 0 ]; then echo "source command failed"; exit 1; fi

#检测slave是否启动
#if [ $Slave_IO_Running == "Yes" ] || [ $Slave_SQL_Running == "Yes" ];then
#    $MYSQL -e "stop slave;"
#    if [ $? -ne 0 ]; then echo "stop slave failed"; exit 1; fi
#fi

$MYSQL -e "stop slave;"

#主从同步
$MYSQL -e "change master to master_host='$master_host',master_user='$repl_user',master_password='$repl_password',master_port=$master_port,master_log_file='$log_file',master_log_pos=$log_pos"
if [ $? -ne 0 ]; then echo "change command failed"; exit 1; fi

cat <<EOF | tee -a ~/.bash_profile
  export MYSQL_BIN=/opt/dtstack/Mysql/mysql_slave/bin
EOF

$MYSQL -e "start slave;"

seelp 10

# 是否同步成功
Slave_IO_Running=`$MYSQL -e "show slave status\G;" | grep Slave_IO_Running|awk -F ' ' '{print $NF}' |head -1`
Slave_SQL_Running=`$MYSQL -e "show slave status\G;" | grep Slave_SQL_Running|awk -F ' ' '{print $NF}'|head -1`

if [ $Slave_IO_Running != "Yes" ] || [ $Slave_SQL_Running != "Yes" ];then
   echo "Slave status is not right."
   exit
fi

fi
}

stop_mysql(){
pid=`ps aux|grep "mysqld"|grep -v grep|awk '{print $2}'`
for i in $pid
do
kill $i
done
}


install_mysql
start_mysql
check_mysql
init_mysql
stop_mysql
