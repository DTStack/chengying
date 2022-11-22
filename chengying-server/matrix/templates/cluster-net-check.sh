#!/bin/sh

#usage ./cluster_net_check.sh user emip [ip]
#parameter 1 is ssh user no password
#parameter 2 is em host
#parameters after 2 are list of host ip to check
# if not passed, some red error appears

if [ $# -lt 2 ] ; then
  echo "USAGE: $0 user [ip]"
  exit
fi

red(){
    CONTENT=$1
    echo -e "\033[31m${CONTENT}\033[0m"
    exit 1
}

green(){
    CONTENT=$1
    echo -e "\033[32m${CONTENT}\033[0m"
}

yellow(){
    CONTENT=$1
    echo -e "\033[33m${CONTENT}\033[0m"
}

user=$1
em_host=$2
shift
shift
i=$@

for i in "$@"; do
    yellow "\n\n===========================check $i？=====================================\n"
    scp net-check.sh $user@$i:/tmp/
    ssh $user@$i "sh /tmp/net-check.sh $em_host 8889 8890"
done
