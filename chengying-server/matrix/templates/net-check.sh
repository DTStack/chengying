#!/bin/sh

#usage ./en-check.sh emip 8889 8890
#parameter 1 is em host
#parameters after 1 are list of ports to check

if [ $# -lt 2 ] ; then
  echo "USAGE: $0 em-host [ports]"
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

em_host=$1
shift
i=$@

for i in "$@"; do
    if command -v nc >/dev/null 2>&1; then
        nc -w 10  $em_host $i  < /dev/null >/dev/null 2>&1
        if [ $? -eq 0 ] ; then
            green "to em $i is ok!"
        else
            red "to em $i is not ok!"
        fi
    else
        red "command nc is not installed!"
    fi
done


check_iptables(){
    yellow "\n\n===========================iptabls已关闭？=====================================\n"
    sudo service iptables status >/dev/null 2>&1
    exit_status=$?
    if [ $exit_status -ne 0 ]; then
        green "iptables is closed"
    else
        red "iptables is not closed"
    fi
}

check_firewalld(){
    yellow "\n\n===========================firewalld已关闭？=====================================\n"
    sudo systemctl is-enabled firewalld >/dev/null 2>&1
    if [ $? == 0 ]; then
        red "firewalld is not closed"
    else
        green "firewalld is closed"
    fi
}

check_selinux(){
    yellow "\n\n===========================selinux已关闭？=====================================\n"
    sudo /usr/sbin/getenforce | grep -E "disabled|ermissive" >/dev/null 2>&1
    exit_status=$?
    if [ $exit_status -eq 0 ]; then
      green "selinux is disabled"
    else
      red "selinux is not closed"
    fi

}

check_to_em(){
    yellow "\n\n===========================to em host network ok？=====================================\n"
    for i in "$@"; do
        if command -v nc >/dev/null 2>&1; then
            nc -w 10  $em_host $i  < /dev/null >/dev/null 2>&1
            if [ $? -eq 0 ] ; then
                green "to em $i is ok!"
            else
                red "to em $i is not ok!"
            fi
        else
            red "command nc is not installed!"
        fi
    done
    green "to em host network is ok"
}

check_to_em
check_iptables
check_firewalld
check_selinux
