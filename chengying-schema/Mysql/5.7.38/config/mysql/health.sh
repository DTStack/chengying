#!/bin/bash

#usage ./port_status.sh localhost 80 8080 8088 false
#parameter 1 is host
#parameters after 1 are list of ports to check
#all ports are ok echo 1, else echo 0


#数栈是否对接第三方mysql
is_other=$4

if [ "${is_other}" == "false" ]; then

if [ $# -lt 2 ] ; then
  echo "USAGE: $0 host [ports]" 
  exit
fi

MYSQL="./bin/mysql -h$1 -u$2 -p$3"
if [ ! -f "./success.txt" ]; then
$MYSQL -e "set GLOBAL max_connections=5000;"
touch success.txt
fi
./bin/mysqlshow -h$1 -u$2 -p$3


else

  printf "对接第三方mysql,不做任何操作\n"

fi