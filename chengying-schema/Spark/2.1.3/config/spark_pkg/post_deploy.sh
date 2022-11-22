#!/bin/bash

source /etc/profile


SPARK_HOME=/opt/dtstack/Spark/spark_pkg

if [[ ! -f /etc/profile.d/spark_base.sh ]];then
cat <<EOF | sudo tee /etc/profile.d/spark_base.sh
export SPARK_HOME=$SPARK_HOME
export SPARK_CONF_DIR=$SPARK_HOME/conf
export PATH="\$PATH:$JAVA_HOME/bin:$SPARK_HOME/bin:$SPARK_HOME/sbin"
EOF

sudo chmod go+r /etc/profile.d/spark_base.sh
fi
. "${SPARK_HOME}/sbin/spark-config.sh"
. "${SPARK_HOME}/bin/load-spark-env.sh"

if [ "$SPARK_IDENT_STRING" = "" ]; then
  export SPARK_IDENT_STRING="$USER"
fi

if [ "$SPARK_LOG_DIR" = "" ]; then
  export SPARK_LOG_DIR="$SPARK_HOME/logs"
fi
mkdir -p "$SPARK_LOG_DIR"
touch "$SPARK_LOG_DIR"/.spark_test > /dev/null 2>&1

TEST_LOG_DIR=$?
if [ "${TEST_LOG_DIR}" = "0" ]; then
  rm -f "$SPARK_LOG_DIR"/.spark_test
else
  chown "$SPARK_IDENT_STRING" "$SPARK_LOG_DIR"
fi