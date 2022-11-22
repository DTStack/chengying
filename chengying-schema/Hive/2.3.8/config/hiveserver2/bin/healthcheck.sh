#!/bin/bash
source /etc/profile
#usage ./port_status.sh localhost 80 8080 8088
#parameter 1 is host
#parameters after 1 are list of ports to check
#all ports are ok echo 1, else echo 0
source /etc/profile

if [ $# -lt 2 ] ; then
  echo "USAGE: $0 host [ports]"
  exit
fi


BEELINE="/opt/dtstack/Hive/hive_pkg/bin/beeline"
address=$1
ret=0
shift
i=$@

# AUTH_MODE=`cat /opt/dtstack/Hive/hive_pkg/conf/hive-site.xml | grep -w 'KERBEROS'`
# 默认输出hive认证模式为简单模式，不开启kerberos认证
sleep 3

$BEELINE -u "jdbc:hive2://$address:$i" -e "select 1;"
if [ $? -eq 0 ] ; then
  ret=$(( $ret + 0 ))
else
  ret=$(( $ret + 1 ))
fi


exit $ret


#for i in "$@"; do
#  if command -v nc >/dev/null 2>&1; then
# # echo "exists nc"
#   #  echo  $address $i
#     nc -w 10  $address $i  < /dev/null >/dev/null 2>&1
#
#   #  echo status ======= $?
#     if [ $? -eq 0 ] ; then
#        ret=$(( $ret + 0 ))
#     else
#        ret=$(( $ret + 1 ))
#     fi
#  else
#    ret=2
#  fi
#done

