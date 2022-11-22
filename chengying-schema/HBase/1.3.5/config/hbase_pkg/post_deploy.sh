#!/bin/bash
set -e

source /etc/profile

#[ -d "${HADOOP_CONF_DIR}" ]&& exit 0 || exit 1

cat <<EOF | sudo tee /etc/profile.d/dthadoop_hbase.sh

export HBASE_HOME=/opt/dtstack/HBase/hbase_pkg
export HBASE_CONF_DIR=/opt/dtstack/HBase/hbase_pkg/conf

export PATH="\$PATH:\$HBASE_HOME/bin"

EOF
