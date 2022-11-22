#export HADOOP_HEAPSIZE=1024
#export HIVE_CONF_DIR=/opt/dtstack/Hadoop/hive/conf
#export HIVE_AUX_JARS_PATH=/opt/dtstack/Hadoop/hive/lib
#export HIVE_HOME=/opt/dtstack/Hadoop/hive

if [ "$SERVICE" = "metastore" ]; then

export HIVE_METASTORE_HADOOP_OPTS=" {{.hive_opts}} -XX:HeapDumpPath=./run -XX:+HeapDumpOnOutOfMemoryError  -XX:+ExitOnOutOfMemoryError -Xloggc:./logs/gc.log -XX:+UseGCLogFileRotation -XX:NumberOfGCLogFiles=5 -XX:GCLogFileSize=512M -XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintHeapAtGC -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime -XX:-OmitStackTraceInFastThrow $HIVE_METASTORE_HADOOP_OPTS"

fi


if [ "$SERVICE" = "hiveserver2" ]; then

export HIVE_SERVER2_HADOOP_OPTS=" -XX:HeapDumpPath=./logs/hiveserver2_heap.hprof  -XX:+HeapDumpOnOutOfMemoryError -XX:+ExitOnOutOfMemoryError -XX:+UseGCLogFileRotation -Xloggc:./logs/hiveserver2_gc.log -XX:NumberOfGCLogFiles=2 -XX:GCLogFileSize=256M -XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintHeapAtGC -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime"

export HADOOP_CLIENT_OPTS="$HADOOP_CLIENT_OPTS $HIVE_SERVER2_HADOOP_OPTS"
fi





