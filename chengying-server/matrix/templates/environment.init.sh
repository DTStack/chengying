#!/bin/sh
set -x
STATIC_HOST='{{.STATIC_HOST}}'
NTP_SERVER='{{.NTP_SERVER}}'

cd `dirname $0`
PWD=`pwd`
#source ../../conf/host.conf
#IPS=`echo $all_host | tr ',' '\n'`
red(){
    CONTENT=$1
#    echo -e "\033[31m${CONTENT}\033[0m"
    echo "${CONTENT}"
    exit 1
}

green(){
    CONTENT=$1
#    echo -e "\033[32m${CONTENT}\033[0m"
}

yellow(){
    CONTENT=$1
#    echo -e "\033[33m${CONTENT}\033[0m"
}

check_data_dir(){
    sudo mkdir -p /data
    sudo chown `whoami`:`whoami` /data
}

check_User(){
    sudo cut -d : -f 1 /etc/passwd |grep $1 >/dev/null 2>&1
    if [ $? != 0 ];then
        sudo groupadd $1
        sudo useradd -g $1 $1
        if [ $? != 0 ];then
            red "New user $1 failed";exit
        fi
    else
        green "User $1 already exits"
    fi
}


install_arthas(){
 # TEMP_ARTHAS_FILE="./as.sh.$$"
  mkdir /opt/dtstack/
  curl -L -O -s "$STATIC_HOST/easyagent/arthas.tar.gz"
  tar -xf "arthas.tar.gz" -C /opt/dtstack/
#  curl -L -O -s "$STATIC_HOST/easyagent/

}
installyum(){
    yellow "\n\n===========================Yum install=====================================\n"
cat <<EOF | sudo tee /etc/yum.repos.d/dtstack.repo > /dev/null
[dtstack]
name=dtstack
cost=2000
priority=99
baseurl=$STATIC_HOST/easyagent/dtstack-repos
enable=1
gpgcheck=0
EOF

    sudo yum --disablerepo=\* --enablerepo=dtstack -y clean all >/dev/null 2>&1
    sudo yum --disablerepo=\* --enablerepo=dtstack -y makecache >/dev/null 2>&1
    sudo yum --disablerepo=\* --enablerepo=dtstack -y install curl unzip python-devel libaio perl nc ntp ntpdate epel-release python2-pip rsync psmisc unixODBC libicu kubectl openssh >/dev/null 2>&1
    exit_status=$?
    if [ $exit_status -eq 0 ]; then
        green "yum install success"
    else
        sudo yum -y install curl unzip python-devel libaio* perl nc ntp ntpdate epel-release python2-pip rsync psmisc unixODBC libicu node_exporter>/dev/null 2>&1
        exit_status=$?
        if [ $exit_status -eq 0 ]; then
            green "yum install success"
        else
            red "yum install failed"
        fi
    fi
}

check_umask(){
    yellow "\n\n===========================Check umask=====================================\n"
    Umask=`umask`
    if [ $Umask == 0022 ];then
      yellow "umask is 0022"
    else
      umask 0022
      echo "umask 0022" |sudo tee -a /etc/profile > /dev/null
    fi

}

check_iptables(){
    yellow "\n\n===========================iptabls已关闭？=====================================\n"
    # iptables
    sudo service iptables stop >/dev/null 2>&1
    sudo chkconfig iptables off >/dev/null 2>&1
    sudo service iptables status >/dev/null 2>&1
    exit_status=$?
    if [ $exit_status -ne 0 ]; then
        green "iptables is closed"
    else
        red "iptables is not closed"
    fi
}

check_firewalld(){
    yellow "\n\n===========================iptabls已关闭？=====================================\n"
    # iptables
    sudo systemctl stop firewalld >/dev/null 2>&1
    sudo systemctl is-enabled firewalld >/dev/null 2>&1
    if [ $? == 0 ]; then
      sudo systemctl disable firewalld >/dev/null 2>&1
    fi
}

check_selinux(){
    yellow "\n\n===========================selinux已关闭？=====================================\n"
    sudo /usr/sbin/getenforce | grep -i "disabled" >/dev/null 2>&1
    exit_status=$?
    if [ $exit_status -eq 0 ]; then
      green "selinux is disabled"
    else
      sudo /usr/sbin/setenforce 0
      yellow "selinux is enabled"
      sudo sed -i 's/^SELINUX=\(.*\)/SELINUX=disabled/g' /etc/selinux/config
    fi

}

check_ping(){
    for ip in `echo $1 | tr ',' '\n'`
    do
        ping -c 1 -w 5 $ip >> /dev/null
        exit_status=$?
        if [ $exit_status -eq 0 ]; then
            green "$2 $ip ping is enabled"
        else
            red "$2 $ip ping is disabled"
        fi
    done
}

check_all_ping(){
    yellow "\n\n===========================ping已开通？=====================================\n"
  check_ping $all_host  all_host
#  check_ping $ftp_host
  check_ping $redis_host  redis_host
  check_ping $rabbitmq_host  rabbitmq_host
  check_ping $db_host    db_host
  check_ping $web_host   web_host
  check_ping $es_host   es_host
  check_ping $jlogstash_host  jlogstash
  check_ping $kafka_host    kafka
  check_ping $zk_host   zk
  f_host=`echo $ftp_host  | awk -F '//' '{print $2}'`
  check_ping $f_host   ftp

}



check_sysctl(){
     yellow "\n\n===========================sysctl conf=====================================\n"
cat <<'EOF' | sudo tee /etc/sysctl.d/sidecar.conf > /dev/null
vm.max_map_count=655300
vm.overcommit_memory=1
vm.swappiness=0
net.ipv4.ip_forward=1
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.tcp_max_tw_buckets = 5000
EOF
sudo systemctl restart systemd-sysctl >/dev/null 2>&1
}


check_limits(){
    yellow "\n\n===========================ulimits已设置？=====================================\n"
cat <<'LIMIT' | sudo tee /etc/security/limits.d/90-nproc.conf > /dev/null
* soft nproc 655350
* hard nproc 655350
LIMIT

cat <<'DEF' | sudo tee /etc/security/limits.d/def.conf > /dev/null
* soft nproc 655350
* hard nproc 655350
DEF

cat <<'EOF' | sudo tee /etc/security/limits.conf > /dev/null
# End of file
* soft nproc 655350
* hard nproc 655350
* soft nofile 655350
* hard nofile 655350
* - memlock unlimited
EOF

    green "ulimit is success"

    sudo sed -i 's#^PATH\(.*\)#PATH="/opt/java/bin:/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin"#g' /etc/init.d/functions
}



check_java(){
    yellow "\n\n===========================java已安装？=====================================\n"

#    host=`echo $ftp_host  | awk -F '//' '{print $2}'`
    #check_port 21 FTP
    if [ -f "/opt/dtstack/java/bin/java" ];then
        green "java success"
    else
        sudo rpm -i "$STATIC_HOST/easyagent/dtstack-java.noarch.rpm"
        if [ ! -f "/opt/dtstack/java/bin/java" ];then
            red "java install fild"
        fi
    fi
}

check_pip_depends(){
    yellow "\n\n===========================python depends[elasticsearch/MySQL-python/redis/pykafka/kazoo/zk_shell]已安装？=====================================\n"

    which pip >/dev/null 2>&1
    if [ $? != 0 ]; then
        yellow "pip install failed"
    else
        sudo pip install --disable-pip-version-check --retries=1 --timeout=1 elasticsearch redis pykafka kazoo zk_shell MySQL-python >/dev/null 2>&1

        exit_status=$?
        if [ $exit_status -eq 0 ]; then
            green "pip depends install success"
        else
            yellow "pip install failed, try install python depends from local package"
            curl -L -O -s "$STATIC_HOST/easyagent/packages.tar.gz"
            tar -xf "packages.tar.gz"

            if [ -d packages ];then
                sudo pip install --no-index --find-links=./packages elasticsearch redis pykafka kazoo zk_shell MySQL-python >/dev/null 2>&1
                exit_status=$?
                if [ $exit_status -eq 0 ]; then
                    green "pip depends install from local success "
                else
                    yellow "pip depends install from local failed"
                fi
            else
                yellow "no local pip packages"
            fi
        fi
    fi
}

check_ssh_nopwd(){
    yellow "\n\n===========================ssh免密码登陆已设置？=====================================\n"

    mkdir -p ~/.ssh
    echo "StrictHostKeyChecking no" > ~/.ssh/config
    echo "UserKnownHostsFile /dev/null" >> ~/.ssh/config
cat <<'EOF' >> ~/.ssh/authorized_keys
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQQtDR6WayM5FgxW7oXRYXJ/g7U7LTGVjCvkJScb6Xf8Um6CwZuFymTd7oVwrY95895AY0mSige0Qn0Q1lGTgxOul+pOcEDdhYsYpFXZnnb/CHtikEsKHAOAPm738Rk9lCFmGuXYHg0ova83qTvgwtw/jh8agIpPU3wnv7lcH3PR20RXM6qF/DUOOgsxTXwug7y4OJwMse2nX/FvcxTp1GRdd8OhUWHosc6sQaaIkoudVPQpnYFIq1Z+QOel8/0Ic0CUSezLeaHX9h0TfdDz85WO4yyBTzOE/O85NragbSPxYJ+WNsyRBajaXRsRmYjT4ug69r1CPZXrVYpWe0FpZZ
EOF
cat <<'EOF' > ~/.ssh/id_rsa.pub
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQQtDR6WayM5FgxW7oXRYXJ/g7U7LTGVjCvkJScb6Xf8Um6CwZuFymTd7oVwrY95895AY0mSige0Qn0Q1lGTgxOul+pOcEDdhYsYpFXZnnb/CHtikEsKHAOAPm738Rk9lCFmGuXYHg0ova83qTvgwtw/jh8agIpPU3wnv7lcH3PR20RXM6qF/DUOOgsxTXwug7y4OJwMse2nX/FvcxTp1GRdd8OhUWHosc6sQaaIkoudVPQpnYFIq1Z+QOel8/0Ic0CUSezLeaHX9h0TfdDz85WO4yyBTzOE/O85NragbSPxYJ+WNsyRBajaXRsRmYjT4ug69r1CPZXrVYpWe0FpZZ
EOF

cat <<'EOF' > ~/.ssh/id_rsa
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0ELQ0elmsjORYMVu6F0WFyf4O1Oy0xlYwr5CUnG+l3/FJugs
Gbhcpk3e6FcK2PefPeQGNJkooHtEJ9ENZRk4MTrpfqTnBA3YWLGKRV2Z52/wh7Yp
BLChwDgD5u9/EZPZQhZhrl2B4NKL2vN6k74MLcP44fGoCKT1N8J7+5XB9z0dtEVz
Oqhfw1DjoLMU18LoO8uDicDLHtp1/xb3MU6dRkXXfDoVFh6LHOrEGmiJKLnVT0KZ
2BSKtWfkDnpfP9CHNAlEnsy3mh1/YdE33Q8/OVjuMsgU8zhPzvOTa2oG0j8WCflj
bMkQWo2l0bEZmI0+LoOva9Qj2V61WKVntBaWWQIDAQABAoIBAEA3Wgxf27q3C+y6
9CfSDC23tWMVu947wpLJ98lUKWCtlS0aCf8kSe09ta4OqNmeAQwkc4VkdJKHb8yX
OXoi/KCxea4DSviFDGDM3cXlxr8GVTSDhsJbdITAwPNEhfn1jNhD7nrFenbIdTde
PD2yLC8wbl29FgMUWkiHP5Oj6dtB/l+uDdBFacPwsE7i6tJPz+CeuPFnHbmm0IcB
qZC3vKCkCDYG7NME70tdIJdRvJDxrHpnC3Uh2ZjHBgWhEifQp1pNax6O8gQqjUMd
43yjECNYHJnPp1PveiGOk656lh6B6zZgVcmWq+unuVBIIKWVH1oLJIlaCcjvGrIE
Ls9y8MECgYEA9uQu+xpJqoDbTEvryZwgPz8uQgJQq5YVarJknvPXpM37+LP/lfXW
Y5HBlUXq7arg4i1nPblcAM7RCYerG9nlheQNVf8uxfkDYVJM0uYWNU7u44jJ9xKF
3d0t6o5TAggwYftmASYmejEzG1QI8Ke4R/8O1UJJe5DPebtdmDVdi0cCgYEA1/HH
p1/5cAXsuGAWoEwfFdOwMBiQkKBtWiykwc/+muw4U5ulrz2mYH8LcD5EMkKvBnAN
JdjTg0ervwC9h8uiP84zZqsK/CE/MSSol7UD0n/le1YgJ4+8mqwjcI10MuNKx5zs
Akbh9ZmZUAOeN1xjZEH4nrdxr77Rq3x1Z9+WYV8CgYEA8TTloXWBWw8hMV1LM2FR
L/luHBde8QRkUCWRoKnimKSV89kHb18K4aRZWJlCeIhlxRNOtkMN28wy9fiR91fe
ofy3+rig3dH2t3JMlA8uy04enjXZ+4UhPC0G2w3Jk0ak3lVaxuya0u6WW7CyO0vh
ApXxy/BDEnHcAbZILpiUl6kCgYAxakWK2p6E8QYjhvGGiwTOBNwblXN3dQ6xBOhK
5aLYpteY4lQ+zrfE+y2N6OHIMpREv91yXLTCJON7WNsGOSYOvcvrYZao7Z69Sys9
VBHk2JmV6LaA9xZsupl9hxEhF/wjw5gKSVgp0JVRxxIsjvq7lZEkGMxuMBdgy6E+
PN73twKBgQDUZz5ZuANb2y/xq2Yg8h/GkeoagsS5a0GvUDtVPtN5+oTz13zUD340
tIx2TLC0RdTBRz8GXCaSHDH2iSEQAAnSh51o7IFEYeU/mMjSXVaux8JEJ9JKThit
Xrx6wkZxipg0Y9jImPSY86oHsKRfVb+zPoR5Frfxw3DbE7h+oUMc7w==
-----END RSA PRIVATE KEY-----
EOF
    chmod 0600 ~/.ssh/id_rsa
}

check_ntp(){
    yellow "\n\n===========================ntp已安装？=====================================\n"

    which ntpdate >/dev/null 2>&1
    if [ $? != 0 ]; then
        yellow "ntp install failed"
    else
        sudo ntpdate -u "$NTP_SERVER" >/dev/null 2>&1
        if [ $? != 0 ]; then
            red "ntpdate $NTP_SERVER failed"
        fi
        sudo mv /etc/ntp.conf /etc/ntp.conf.bk
        echo "server $NTP_SERVER iburst" | sudo tee /etc/ntp.conf > /dev/null
        sudo systemctl restart ntpd >/dev/null 2>&1
        if [ $? != 0 ]; then
            red "restart ntpd failed"
        fi

        sudo systemctl is-enabled ntpd >/dev/null 2>&1
        if [ $? != 0 ]; then
          sudo systemctl enable ntpd >/dev/null 2>&1
        fi
        sudo grep "ntpdate -u $NTP_SERVER" /etc/rc.d/rc.local >/dev/null 2>&1
        if [ $? != 0 ]; then
            echo "ntpdate -u $NTP_SERVER" | sudo tee -a /etc/rc.d/rc.local > /dev/null
        fi
        sudo chmod a+x /etc/rc.d/rc.local
    fi
}

check_node_exporter(){
    yellow "\n\n===========================node_exporter已安装？=====================================\n"

    which node_exporter >/dev/null 2>&1
    if [ $? != 0 ]; then
        sudo rpm -i "$STATIC_HOST/easyagent/dtstack-repos/node_exporter-0.15.2-1.el7.centos.x86_64.rpm"
        which node_exporter >/dev/null 2>&1
        if [ $? != 0 ]; then
           red "node_exporter install failed"
        fi
    else

        green "node_exporter success"
    fi

    sudo systemctl restart node_exporter >/dev/null 2>&1
    if [ $? != 0 ]; then
        red "restart node_exporter failed"
    fi

    sudo systemctl is-enabled node_exporter >/dev/null 2>&1
    if [ $? != 0 ]; then
      sudo systemctl enable node_exporter >/dev/null 2>&1
    fi
}

add_rc_local(){
   yellow "\n\n===========================ntp已安装？=====================================\n"

    sudo grep "cgroup.clone_children" /etc/rc.d/rc.local >/dev/null 2>&1
    if [ $? != 0 ]; then
        echo 1 | sudo tee /sys/fs/cgroup/cpu/cgroup.clone_children > /dev/null
        echo "echo 1 > /sys/fs/cgroup/cpu/cgroup.clone_children" | sudo tee -a /etc/rc.d/rc.local > /dev/null
    fi
    sudo chmod a+x /etc/rc.d/rc.local

}

check_data_dir
#check_User admin
installyum
check_umask
#check_iptables
check_firewalld
check_selinux
check_sysctl
#check_all_ping
check_limits
check_java
add_rc_local
check_pip_depends
check_ssh_nopwd
check_ntp
check_node_exporter

