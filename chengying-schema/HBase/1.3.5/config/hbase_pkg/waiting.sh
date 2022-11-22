#!/bin/bash
source /etc/profile
[ ! -d "${HADOOP_CONF_DIR}" ]&& exit 1
\cp  -f  ${HADOOP_CONF_DIR}/{core-site.xml,hdfs-site.xml}  conf/

tail -f /dev/null

