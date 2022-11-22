#!/bin/bash 
set -e

cur=`pwd`
case "`uname`" in
    Linux)
        bin_path=$(readlink -f $(dirname $0))
        ;;
    *)
        bin_path=`cd $(dirname $0); pwd`
        ;;
esac


HADOOP_HOME=/data/hadoop_base
USER=`whoami`

#sed -i s/\${user\\.name}/$USER/g etc/hadoop/core-site.xml

if [ -d $HADOOP_HOME ];then
    unlink $HADOOP_HOME
fi
ln -s $bin_path $HADOOP_HOME

cat <<EOF | sudo tee /etc/profile.d/hadoop_base.sh
export HADOOP_HOME=$HADOOP_HOME
export PATH="$PATH:$JAVA_HOME/bin:$HADOOP_HOME/bin:$HADOOP_HOME/sbin"
export HADOOP_USER_NAME=$USER
export HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop
export YARN_CONF_DIR=$HADOOP_HOME/etc/hadoop
EOF

sudo chmod go+r /etc/profile.d/hadoop_base.sh

function check_corexml() {


sed -i s/"hadoop.proxyuser.{{.config_user_name}}.hosts"/"hadoop.proxyuser.${USER}.hosts"/g $HADOOP_HOME/etc/hadoop/core-site.xml
sed -i s/"hadoop.proxyuser.{{.config_user_name}}.groups"/"hadoop.proxyuser.${USER}.groups"/g $HADOOP_HOME/etc/hadoop/core-site.xml


}

check_corexml

mkdir -p /tmp/spark-events