export default {
    msg: 'ok',
    code: 0,
    data: {
        default: {
            SLB: {
                Version: 'slb-1.1',
                Group: 'default',
                Config: {
                    ip: {
                        Default: '127.0.0.2',
                        Desc: 'internal',
                        Type: 'internal',
                        Value: '127.0.0.2'
                    }
                },
                BaseProduct: '',
                BaseService: '',
                BaseParsed: false,
                ServiceAddr: {
                    Host: [],
                    IP: [
                        '172.16.10.16'
                    ],
                    NodeId: 0,
                    SingleIndex: 0
                }
            },
            dtlog_new: {
                Version: 'dtlog-new-2.3',
                Instance: {
                    UseCloud: true,
                    ConfigPaths: [
                        'conf/config.ini',
                        'conf/config-sp.ini'
                    ],
                    Logs: [
                        '*.log',
                        '/tmp/dtlog_post_undeploy.log'
                    ],
                    HealthCheck: {
                        Shell: 'echo test health',
                        Period: '5s',
                        StartPeriod: '',
                        Timeout: '',
                        Retries: 1
                    },
                    Cmd: './dtlog.sh conf/config.ini dtlog.log 172.16.10.107',
                    HARoleCmd: 'echo master',
                    Pseudo: true,
                    PostDeploy: 'echo dtlog deploy success >> dtlog_post_deploy.log && echo risk.sql',
                    PostUndeploy: 'echo dtlog undeploy success >> /tmp/dtlog_post_undeploy.log && echo risk.sh',
                    PrometheusPort: ''
                },
                Group: 'default',
                DependsOn: [
                    'dtuic',
                    'test3'
                ],

                ServiceAddr: {
                    Host: [],
                    IP: [
                        '172.16.10.16'
                    ],
                    NodeId: 0,
                    SingleIndex: 0
                },
                Config: {
                    config_path: {
                        Default: 'conf/config.ini',
                        Desc: 'internal',
                        Type: 'internal',
                        Value: 'conf/config.ini'
                    },
                    dtuic_ip: {
                        Default: {
                            Host: [
                                '172-16-10-107',
                                '172-16-10-108'
                            ],
                            IP: [
                                '172.16.10.110',
                                '172.16.10.101'
                            ],
                            NodeId: 0,
                            SingleIndex: 0
                        },
                        Desc: 'internal',
                        Type: 'internal',
                        Value: {
                            Host: [
                                '172-16-10-107',
                                '172-16-10-108'
                            ],
                            IP: [
                                '172.16.10.107',
                                '172.16.10.108'
                            ],
                            NodeId: 0,
                            SingleIndex: 0
                        }
                    },
                    log: {
                        Default: 'dtlog.log',
                        Desc: 'internal',
                        Type: 'internal',
                        Value: 'dtlog.log'
                    },
                    self_ip: {
                        Default: '${@dtlog_new}',
                        Desc: 'internal',
                        Type: 'internal',
                        Value: '${@dtlog_new}'
                    }
                },
                BaseProduct: '',
                BaseService: '',
                BaseParsed: false
            },
            test1: {
                Version: 'test1-1.0',
                Instance: {
                    HealthCheck: {
                        Shell: 'sleep 2',
                        Period: '3s',
                        StartPeriod: '',
                        Timeout: '',
                        Retries: 1
                    },
                    Cmd: './test1.sh',
                    PostDeploy: '',
                    PostUndeploy: '',
                    PrometheusPort: ''
                },
                Group: 'default',
                BaseProduct: '',
                BaseService: '',
                BaseParsed: false,
                ServiceAddr: {
                    Host: [
                        '172-16-10-107',
                        '172-16-10-108'
                    ],
                    IP: [
                        '172.16.10.107',
                        '172.16.10.108'
                    ],
                    NodeId: 0,
                    SingleIndex: 0
                }
            },
            test3: {
                Version: 'test3-1.0',
                Instance: {
                    UseCloud: true,
                    HealthCheck: {
                        Shell: 'sleep 2',
                        Period: '3s',
                        StartPeriod: '',
                        Timeout: '',
                        Retries: 1
                    },
                    Cmd: './test3.sh',
                    PostDeploy: '',
                    PostUndeploy: '',
                    PrometheusPort: ''
                },
                Group: 'default',
                BaseProduct: '',
                BaseService: '',
                BaseParsed: false,
                ServiceAddr: {
                    Host: [],
                    IP: [
                        '172.16.10.16'
                    ],
                    NodeId: 0,
                    SingleIndex: 0
                }
            }
        },
        'dtuic-group': {
            dtuic: {
                Version: 'dtuic-1.0',
                Instance: {
                    HealthCheck: {
                        Shell: 'sleep 2',
                        Period: '3s',
                        StartPeriod: '',
                        Timeout: '',
                        Retries: 1
                    },
                    Cmd: './dtuic.sh',
                    PostDeploy: '',
                    PostUndeploy: '',
                    PrometheusPort: '',
                    MaxReplica: '3',
                    StartAfterInstall: true
                },
                Group: 'dtuic-group',
                DependsOn: [
                    'test1'
                ],
                BaseProduct: '',
                BaseService: '',
                BaseParsed: false,
                ServiceAddr: {
                    Host: [
                        '172-16-10-107',
                        '172-16-10-108'
                    ],
                    IP: [
                        '172.16.10.107',
                        '172.16.10.108'
                    ],
                    NodeId: 0,
                    SingleIndex: 0
                }
            }
        }
    }
}
