#!/bin/bash
source /etc/profile
STATE=`bin/mysql -h$1 -u$2 -p$3 -e "show global variables like 'read_only'" | grep read_only | awk -F ' ' '{print $2}'`

if [ $STATE == 'ON' ];then
   echo 'Slave'
else 
   echo 'Master'
fi

exit
