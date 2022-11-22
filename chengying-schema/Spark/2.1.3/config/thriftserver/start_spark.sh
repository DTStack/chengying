#!/bin/bash

source /etc/profile

command="org.apache.spark.sql.hive.thriftserver.HiveThriftServer2"
instance=1


spark_rotate_log ()
{
    log=$1;
    num=5;
    if [ -n "$2" ]; then
	num=$2
    fi
    if [ -f "$log" ]; then # rotate logs
	while [ $num -gt 1 ]; do
	    prev=`expr $num - 1`
	    [ -f "$log.$prev" ] && mv "$log.$prev" "$log.$num"
	    num=$prev
	done
	mv "$log" "$log.$num";
    fi
}


if [ -d {{.spark_thriftserver_logs_path}} ];then
  SPARK_LOG_DIR={{.spark_thriftserver_logs_path}}
else
  mkdir -p {{.spark_thriftserver_logs_path}}
fi

# some variables
logname="spark-$SPARK_IDENT_STRING-$command-$instance-$HOSTNAME"
log="$SPARK_LOG_DIR/spark-$SPARK_IDENT_STRING-$command-$instance-$HOSTNAME.log"

# Set default scheduling priority
if [ "$SPARK_NICENESS" = "" ]; then
    export SPARK_NICENESS=0
fi
spark_rotate_log "$log"

{{if eq (print .is_kerberos) "true"}}
nice -n "$SPARK_NICENESS" bash "${SPARK_HOME}"/bin/spark-submit  --class "$command"  --name 'Thrift JDBC/ODBC Server'  --conf "spark.driver.extraJavaOptions=-javaagent:./dtstack/prometheus/jmx_prometheus_javaagent-0.3.1.jar=9508:./dtstack/prometheus/spark-prometheus.yml" --conf spark.driver.memory={{.spark_driver_mem}} --executor-memory {{.executor_mem}} --num-executors {{.executor_num}} --name=ThriftServer02 --principal {{.principal}} --keytab {{.keytab}} --driver-java-options "-XX:HeapDumpPath=./logs/thriftserver_heap.hprof -XX:+HeapDumpOnOutOfMemoryError -XX:+ExitOnOutOfMemoryError -XX:+UseGCLogFileRotation -Xloggc:./logs/spark_thriftserver_gc.log -XX:NumberOfGCLogFiles=2 -XX:GCLogFileSize=256M -XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintHeapAtGC -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime -Dlog4j.configuration=file:/opt/dtstack/Spark/spark_pkg/conf/log4j.properties -Dapp.logging.name=$logname "   >> "${log}" 2>&1 < /dev/null
{{else}}
nice -n "$SPARK_NICENESS" bash "${SPARK_HOME}"/bin/spark-submit  --class "$command"  --name 'Thrift JDBC/ODBC Server' --conf "spark.driver.extraJavaOptions=-javaagent:./dtstack/prometheus/jmx_prometheus_javaagent-0.3.1.jar=9508:./dtstack/prometheus/spark-prometheus.yml" --conf spark.driver.memory={{.spark_driver_mem}} --executor-memory {{.executor_mem}} --num-executors {{.executor_num}} --name=ThriftServer02  --driver-java-options "-XX:HeapDumpPath=./logs/thriftserver_heap.hprof -XX:+HeapDumpOnOutOfMemoryError -XX:+ExitOnOutOfMemoryError -XX:+UseGCLogFileRotation -Xloggc:./logs/spark_thriftserver_gc.log -XX:NumberOfGCLogFiles=2 -XX:GCLogFileSize=256M -XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintHeapAtGC -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime -Dlog4j.configuration=file:/opt/dtstack/Spark/spark_pkg/conf/log4j.properties -Dapp.logging.name=$logname "   >> "${log}" 2>&1 < /dev/null
{{end}}
