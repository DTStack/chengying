export default {
    msg: 'ok',
    code: 0,
    data: [
        {
            product_name: 'DTBase',
            group: 'default',
            service_name_list: 'prometheus,pushgateway,redis,zookeeper'
        },
        {
            product_name: 'DTBase',
            group: 'mysql',
            service_name_list: 'mysql,mysql_slave'
        },
        {
            product_name: 'Hadoop',
            group: 'carbondata',
            service_name_list: 'carbon_hivemetastore,carbon_thriftserver'
        },
        {
            product_name: 'Hadoop',
            group: 'default',
            service_name_list: 'hadoop_pkg'
        },
        {
            product_name: 'Hadoop',
            group: 'flink',
            service_name_list: 'historyserver,yarn_session'
        },
        {
            product_name: 'Hadoop',
            group: 'hdfs',
            service_name_list:
        'hdfs_datanode,hdfs_journalnode,hdfs_namenode,hdfs_zkfc'
        },
        {
            product_name: 'Hadoop',
            group: 'spark',
            service_name_list: 'hivemetastore,thriftserver'
        },
        {
            product_name: 'Hadoop',
            group: 'yarn',
            service_name_list: 'jobhistory,yarn_nodemanager,yarn_resourcemanager'
        },
        {
            product_name: 'DTinsight',
            group: 'default',
            service_name_list:
        'DTinsight_Analytics,DTinsight_API,DTinsight_Batch,DTinsight_Console,DTinsight_Engine,DTinsight_Gateway,DTinsight_Schedule,DTinsight_Stream,DTinsight_Valid,tengine'
        }
    ]
};
