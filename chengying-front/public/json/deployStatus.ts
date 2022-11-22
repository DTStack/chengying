export default {
    msg: 'ok',
    code: 0,
    data: {
        count: 1,
        list: [
            {
                create_time: '2018-12-28 18:57:35',
                deploy_time: '2019-01-02 17:19:10',
                id: 4,
                is_current_version: 1,
                product: {
                    ParentProductName: 'DTinsight',
                    ProductName: 'DTinsight',
                    ProductVersion: '3.3.0',
                    Service: {
                        DTinsight_API: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/application.properties',
                                    'conf/service.properties'
                                ],
                                Logs: ['logs/*.log'],
                                Environment: { HO_HEAP_SIZE: '512m', HO_MIN_HEAP_SIZE: '512m' },
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 8087',
                                    Period: '20s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: './bin/base.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9510'
                            },
                            Group: 'default',
                            DependsOn: ['DTinsight_Gateway'],
                            Config: {
                                mail_from: {
                                    Default: '袋鼠云数据产品技术支持',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '袋鼠云数据产品技术支持test'
                                },
                                mail_host: {
                                    Default: 'smtp.mxhichina.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'smtp.mxhichina.com'
                                },
                                mail_password: {
                                    Default: '@@dtstack2018..',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '@@dtstack2018..'
                                },
                                mail_port: {
                                    Default: '25',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '25'
                                },
                                mail_subject: {
                                    Default: 'RDOS',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'RDOS'
                                },
                                mail_username: {
                                    Default: 'data_support@dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'data_support@dtstack.com'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                rdos_api_gateway_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_api_web_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_tengine_url: {
                                    Default: 'insight.dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'insight.dtstack.com'
                                },
                                redis_db: {
                                    Default: '1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '1'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                uic_host: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Analytics: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/application.properties',
                                    'conf/hive.properties',
                                    'conf/service.properties'
                                ],
                                Logs: ['logs/*.log'],
                                Environment: { HO_HEAP_SIZE: '512m', HO_MIN_HEAP_SIZE: '512m' },
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 9022',
                                    Period: '60s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 1
                                },
                                Cmd: './bin/base.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9521'
                            },
                            Group: 'default',
                            DependsOn: ['DTinsight_Console'],
                            Config: {
                                carbon_thriftserver_ip: {
                                    Default: {
                                        Host: [],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: [],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mail_from: {
                                    Default: '袋鼠云数据产品技术支持',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '袋鼠云数据产品技术支持'
                                },
                                mail_host: {
                                    Default: 'smtp.mxhichina.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'smtp.mxhichina.com'
                                },
                                mail_password: {
                                    Default: '@@dtstack2018..',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '@@dtstack2018..'
                                },
                                mail_port: {
                                    Default: '25',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '25'
                                },
                                mail_subject: {
                                    Default: 'RDOS',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'RDOS'
                                },
                                mail_username: {
                                    Default: 'data_support@dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'data_support@dtstack.com'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                rdos_console_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                uic_host: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Batch: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/application.properties',
                                    'conf/hive.properties',
                                    'conf/service.properties',
                                    'etc/hadoop/hdfs-site.xml',
                                    'etc/hadoop/core-site.xml'
                                ],
                                Logs: ['logs/*.log'],
                                Environment: { HO_HEAP_SIZE: '512m', HO_MIN_HEAP_SIZE: '512m' },
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 9020',
                                    Period: '40s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: './bin/base.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9513'
                            },
                            Group: 'default',
                            DependsOn: ['DTinsight_Engine', 'DTinsight_Console'],
                            Config: {
                                engine_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                ide_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                isHA_hdfs: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                mail_from: {
                                    Default: '袋鼠云数据产品技术支持',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '袋鼠云数据产品技术支持'
                                },
                                mail_host: {
                                    Default: 'smtp.mxhichina.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'smtp.mxhichina.com'
                                },
                                mail_password: {
                                    Default: '@@dtstack2018..',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '@@dtstack2018..'
                                },
                                mail_port: {
                                    Default: '25',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '25'
                                },
                                mail_subject: {
                                    Default: 'RDOS',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'RDOS'
                                },
                                mail_username: {
                                    Default: 'data_support@dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'data_support@dtstack.com'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_db: {
                                    Default: 'metastore',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'metastore'
                                },
                                mysql_metastore_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_metastore_port: {
                                    Default: '3306',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '3306'
                                },
                                mysql_metastore_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                namenode_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                rdos_console_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_db: {
                                    Default: '1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '1'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                },
                                thriftserver_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_host: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Console: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/application.properties',
                                    'etc/hadoop/hdfs-site.xml',
                                    'etc/hadoop/core-site.xml',
                                    'etc/hadoop/yarn-site.xml',
                                    'etc/hadoop/hive-site.xml'
                                ],
                                Logs: ['logs/*.log'],
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 8084',
                                    Period: '30s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: 'bin/base.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9520'
                            },
                            Group: 'default',
                            Config: {
                                engine_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                ha_rm_id1: {
                                    Default: 'rm1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm1'
                                },
                                ha_rm_id2: {
                                    Default: 'rm2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm2'
                                },
                                hive_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                isHA_hdfs: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                isHA_yarn: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_db: {
                                    Default: 'metastore',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'metastore'
                                },
                                mysql_metastore_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_metastore_port: {
                                    Default: '3306',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '3306'
                                },
                                mysql_metastore_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                namenode_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                redis_db: {
                                    Default: '1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '1'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                resourcemanager_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                },
                                uic_host: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                },
                                zk_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Engine: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/node.yml',
                                    'etc/hadoop/hdfs-site.xml',
                                    'etc/hadoop/core-site.xml',
                                    'etc/hadoop/yarn-site.xml',
                                    'etc/hadoop/hive-site.xml'
                                ],
                                Logs: ['logs/*.log'],
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 8090',
                                    Period: '40s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: './bin/base.sh',
                                PostDeploy:
                  '[ -d /data/flinkplugin ] || cp -rf flink/flinkplugin /data/',
                                PostUndeploy: '',
                                PrometheusPort: '9515'
                            },
                            Group: 'default',
                            Config: {
                                flink_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                ha_rm_id1: {
                                    Default: 'rm1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm1'
                                },
                                ha_rm_id2: {
                                    Default: 'rm2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm2'
                                },
                                historyserver_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                hive_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                isHA_hdfs: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                isHA_yarn: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                namenode_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                rdos_engine_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                resourcemanager_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                },
                                zk_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Gateway: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: ['conf/node.yml'],
                                Logs: ['logs/*.log'],
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 8086',
                                    Period: '20s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: './bin/base.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9512'
                            },
                            Group: 'default',
                            Config: {
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                rdos_api_web: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_db: {
                                    Default: '1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '1'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Schedule: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/application.properties',
                                    'conf/service.properties',
                                    'conf/hive.properties',
                                    'plugin/da/conf/application.properties',
                                    'plugin/dq/conf/application.properties',
                                    'plugin/analysis/conf/application.properties',
                                    'plugin/stream/conf/application.properties',
                                    'plugin/rdos/conf/application.properties',
                                    'plugin/rdos/conf/hadoop/hdfs-site.xml',
                                    'plugin/rdos/conf/hadoop/core-site.xml'
                                ],
                                Logs: ['logs/*.log'],
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 9030',
                                    Period: '20s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: './bin/base.sh 172.16.10.100',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9514'
                            },
                            Group: 'default',
                            DependsOn: ['DTinsight_Engine', 'DTinsight_Console'],
                            Config: {
                                engine_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                isHA_hdfs: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                mail_from: {
                                    Default: '袋鼠云数据产品技术支持',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '袋鼠云数据产品技术支持'
                                },
                                mail_host: {
                                    Default: 'smtp.mxhichina.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'smtp.mxhichina.com'
                                },
                                mail_password: {
                                    Default: '@@dtstack2018..',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '@@dtstack2018..'
                                },
                                mail_port: {
                                    Default: '25',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '25'
                                },
                                mail_subject: {
                                    Default: 'RDOS',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'RDOS'
                                },
                                mail_username: {
                                    Default: 'data_support@dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'data_support@dtstack.com'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_db: {
                                    Default: 'metastore',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'metastore'
                                },
                                mysql_metastore_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_metastore_port: {
                                    Default: '3306',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '3306'
                                },
                                mysql_metastore_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                namenode_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                prometheus_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                pushgateway_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_console_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_db: {
                                    Default: '1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '1'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_pwd: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                },
                                thriftserver_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Stream: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/application.properties',
                                    'conf/metric.properties',
                                    'conf/service.properties'
                                ],
                                Logs: ['logs/*.log'],
                                Environment: { HO_HEAP_SIZE: '512m', HO_MIN_HEAP_SIZE: '512m' },
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 9021',
                                    Period: '60s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 1
                                },
                                Cmd: './bin/base.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9522'
                            },
                            Group: 'default',
                            DependsOn: ['DTinsight_Console', 'DTinsight_Engine'],
                            Config: {
                                engine_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mail_from: {
                                    Default: '袋鼠云数据产品技术支持',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '袋鼠云数据产品技术支持'
                                },
                                mail_host: {
                                    Default: 'smtp.mxhichina.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'smtp.mxhichina.com'
                                },
                                mail_password: {
                                    Default: '@@dtstack2018..',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '@@dtstack2018..'
                                },
                                mail_port: {
                                    Default: '25',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '25'
                                },
                                mail_subject: {
                                    Default: 'RDOS',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'RDOS'
                                },
                                mail_username: {
                                    Default: 'data_support@dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'data_support@dtstack.com'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                prometheus_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                pushgateway_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_console_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                streamapp_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_host: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        DTinsight_Valid: {
                            Version: '3.3.0',
                            Instance: {
                                ConfigPaths: [
                                    'conf/application.properties',
                                    'conf/service.properties',
                                    'conf/hive.properties'
                                ],
                                Logs: ['logs/*.log'],
                                Environment: { HO_HEAP_SIZE: '512m', HO_MIN_HEAP_SIZE: '512m' },
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.100 8089',
                                    Period: '20s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: './bin/base.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9511'
                            },
                            Group: 'default',
                            DependsOn: ['DTinsight_Schedule'],
                            Config: {
                                mail_from: {
                                    Default: '袋鼠云数据产品技术支持',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '袋鼠云数据产品技术支持'
                                },
                                mail_host: {
                                    Default: 'smtp.mxhichina.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'smtp.mxhichina.com'
                                },
                                mail_password: {
                                    Default: '@@dtstack2018..',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '@@dtstack2018..'
                                },
                                mail_port: {
                                    Default: '25',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '25'
                                },
                                mail_subject: {
                                    Default: 'RDOS',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'RDOS'
                                },
                                mail_username: {
                                    Default: 'data_support@dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'data_support@dtstack.com'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_db: {
                                    Default: 'metastore',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'metastore'
                                },
                                mysql_metastore_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_metastore_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_metastore_port: {
                                    Default: '3306',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '3306'
                                },
                                mysql_metastore_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                rdos_schedule_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_tengine_url: {
                                    Default: 'insight.dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'insight.dtstack.com'
                                },
                                redis_db: {
                                    Default: '1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '1'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                thriftserver_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_host: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        carbon_thriftserver: {
                            Version: '2.1.0',
                            Group: 'default',
                            Config: {
                                carbon_thriftserver_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                dt_center_analysis_ip: {
                                    Default: '127.0.0.1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '172.16.10.100'
                                },
                                hive_analysis_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                zk_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: 'Hadoop',
                            BaseService: 'carbon_thriftserver',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: [],
                                IP: ['172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        dtuic: {
                            Version: '3.0.0',
                            Instance: {
                                ConfigPaths: ['conf/application-prod.properties'],
                                Logs: ['logs/*.log'],
                                Environment: { HO_HEAP_SIZE: '512m', HO_MIN_HEAP_SIZE: '512m' },
                                HealthCheck: {
                                    Shell: './bin/health.sh 172.16.10.47 8006',
                                    Period: '20s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 3
                                },
                                Cmd: 'bin/dtuic_java.sh',
                                PostDeploy: '',
                                PostUndeploy: '',
                                PrometheusPort: '9516'
                            },
                            Group: 'default',
                            Config: {
                                cookie_domain: {
                                    Default: '172.16.10.47',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                redis_host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                redis_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                tengine_uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node2'],
                                IP: ['172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        flink: {
                            Version: '1.5.0',
                            Group: 'default',
                            Config: {
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                ha_rm_id1: {
                                    Default: 'rm1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm1'
                                },
                                ha_rm_id2: {
                                    Default: 'rm2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm2'
                                },
                                jobmanager: {
                                    Default: '2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '2'
                                },
                                namenode_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                pushgateway_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                resourcemanager_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                },
                                taskmanager: {
                                    Default: '3',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '3'
                                },
                                zk_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: 'Hadoop',
                            BaseService: 'yarn_session',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node2'],
                                IP: ['172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        hdfs_namenode: {
                            Version: '2.7.6-3',
                            Group: 'default',
                            Config: {
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                isHA: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                journalnode_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                namenode_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                }
                            },
                            BaseProduct: 'Hadoop',
                            BaseService: 'hdfs_namenode',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node1', 'node2'],
                                IP: ['172.16.10.26', '172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        historyserver: {
                            Version: '1.5.0',
                            Group: 'default',
                            Config: {
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                }
                            },
                            BaseProduct: 'Hadoop',
                            BaseService: 'historyserver',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node1'],
                                IP: ['172.16.10.26'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        hive: {
                            Version: '2.1.1',
                            Group: 'default',
                            Config: {
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                mysql_db: {
                                    Default: 'metastore',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'metastore'
                                },
                                mysql_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                mysql_password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                mysql_port: {
                                    Default: '3306',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '3306'
                                },
                                mysql_user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                }
                            },
                            BaseProduct: 'Hadoop',
                            BaseService: 'hivemetastore',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node1'],
                                IP: ['172.16.10.26'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        mysql: {
                            Version: '5.6.35-2',
                            Group: 'default',
                            Config: {
                                db: {
                                    Default: 'metastore',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'metastore'
                                },
                                password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                },
                                port: {
                                    Default: '3306',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '3306'
                                },
                                repl_password: {
                                    Default: 'repl@123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'repl@123'
                                },
                                repl_user: {
                                    Default: 'repl',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'repl'
                                },
                                tengine_ip: {
                                    Default: '127.0.0.1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '172.16.10.100'
                                },
                                user: {
                                    Default: 'dtstack',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'dtstack'
                                }
                            },
                            BaseProduct: 'DTBase',
                            BaseService: 'mysql',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node1'],
                                IP: ['172.16.10.26'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        prometheus: {
                            Version: '2.3.2',
                            Group: 'default',
                            BaseProduct: 'DTBase',
                            BaseService: 'prometheus',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node2'],
                                IP: ['172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        pushgateway: {
                            Version: '0.4.0',
                            Group: 'default',
                            BaseProduct: 'DTBase',
                            BaseService: 'pushgateway',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node2'],
                                IP: ['172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        redis: {
                            Version: '3.2.5',
                            Group: 'default',
                            Config: {
                                db: {
                                    Default: '1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '1'
                                },
                                host: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                password: {
                                    Default: 'abc123',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'abc123'
                                }
                            },
                            BaseProduct: 'DTBase',
                            BaseService: 'redis',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node1'],
                                IP: ['172.16.10.26'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        tengine: {
                            Version: '2.1.1-3',
                            Instance: {
                                ConfigPaths: [
                                    'conf/conf.d/dtuic.conf',
                                    'conf/conf.d/rdos.conf',
                                    'rdos_front/dist/public/console/config/config.js',
                                    'rdos_front/dist/public/common/config.js',
                                    'rdos_front/dist/public/analyticsEngine/config/config.js',
                                    'rdos_front/dist/public/stream/config/config.js',
                                    'rdos_front/dist/public/main/config/config.js',
                                    'rdos_front/dist/public/rdos/config/config.js',
                                    'rdos_front/dist/public/dataQuality/config/config.js',
                                    'rdos_front/dist/public/dataLabel/config/config.js',
                                    'rdos_front/dist/public/dataApi/config/config.js'
                                ],
                                Logs: ['logs/*.log'],
                                HealthCheck: {
                                    Shell: './health.sh 172.16.10.100 443 82',
                                    Period: '20s',
                                    StartPeriod: '',
                                    Timeout: '',
                                    Retries: 1
                                },
                                Cmd: './nginx.sh',
                                PostDeploy: './post_deploy.sh',
                                PostUndeploy: '',
                                PrometheusPort: ''
                            },
                            Group: 'default',
                            DependsOn: [
                                'dtuic',
                                'DTinsight_Batch',
                                'DTinsight_Valid',
                                'DTinsight_API',
                                'DTinsight_Console'
                            ],
                            Config: {
                                analysis_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                console_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                dtuic_server_host: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                gateway_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_api_web_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_dq_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_server_host: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rdos_tengine_url: {
                                    Default: 'insight.dtstack.com',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'insight.dtstack.com'
                                },
                                rdos_web_port: {
                                    Default: '443',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '443'
                                },
                                schedule_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                stream_ip: {
                                    Default: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3'],
                                        IP: ['172.16.10.100'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                uic_port: {
                                    Default: '82',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '82'
                                },
                                uic_url: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: '',
                            BaseService: '',
                            BaseParsed: false,
                            ServiceAddr: {
                                Host: ['node3'],
                                IP: ['172.16.10.100'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        thriftserver: {
                            Version: '2.1.0',
                            Group: 'default',
                            Config: {
                                ha_namenode_id1: {
                                    Default: 'nn1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn1'
                                },
                                ha_namenode_id2: {
                                    Default: 'nn2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'nn2'
                                },
                                ha_rm_id1: {
                                    Default: 'rm1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm1'
                                },
                                ha_rm_id2: {
                                    Default: 'rm2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm2'
                                },
                                hive_ip: {
                                    Default: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1'],
                                        IP: ['172.16.10.26'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                namenode_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                nameservices: {
                                    Default: 'ns1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'ns1'
                                },
                                resourcemanager_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                rpc_address_port: {
                                    Default: '9000',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: '9000'
                                },
                                zk_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: 'Hadoop',
                            BaseService: 'thriftserver',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node1'],
                                IP: ['172.16.10.26'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        yarn_resourcemanager: {
                            Version: '2.7.6-3',
                            Group: 'default',
                            Config: {
                                ha_rm_id1: {
                                    Default: 'rm1',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm1'
                                },
                                ha_rm_id2: {
                                    Default: 'rm2',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'rm2'
                                },
                                isHA: {
                                    Default: 'true',
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: 'true'
                                },
                                jobhistory_ip: {
                                    Default: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node2'],
                                        IP: ['172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                resourcemanager_ip: {
                                    Default: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node1', 'node2'],
                                        IP: ['172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                },
                                zk_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: 'Hadoop',
                            BaseService: 'yarn_resourcemanager',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node1', 'node2'],
                                IP: ['172.16.10.26', '172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        },
                        zookeeper: {
                            Version: '3.4.8-2',
                            Group: 'default',
                            Config: {
                                zk_ip: {
                                    Default: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    },
                                    Desc: 'internal',
                                    Type: 'internal',
                                    Value: {
                                        Host: ['node3', 'node1', 'node2'],
                                        IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                        NodeId: 0,
                                        SingleIndex: 0
                                    }
                                }
                            },
                            BaseProduct: 'DTBase',
                            BaseService: 'zookeeper',
                            BaseParsed: true,
                            ServiceAddr: {
                                Host: ['node3', 'node1', 'node2'],
                                IP: ['172.16.10.100', '172.16.10.26', '172.16.10.47'],
                                NodeId: 0,
                                SingleIndex: 0
                            }
                        }
                    }
                },
                product_name: 'DTinsight',
                product_version: '3.3.0',
                status: 'deploying'
            }
        ]
    }
};
