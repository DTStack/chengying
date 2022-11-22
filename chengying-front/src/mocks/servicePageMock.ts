export default {
    getServiceGroup: {
        msg: 'ok',
        code: 0,
        data: {
            count: 0,
            groups: {}
        }
    },
    getCurrentProduct: {
        msg: 'ok',
        code: 0,
        data: {
            create_time: '2020-09-04 12:51:07',
            deploy_time: '2020-09-04 17:00:35',
            id: 4,
            is_current_version: 1,
            product: {
                ParentProductName: 'DTEM',
                ProductName: 'DTK8S',
                ProductNameDisplay: '',
                ProductVersion: '1.0.5',
                Service: {
                    docker: {
                        ServiceDisplay: '',
                        Version: '18.06.1-ce',
                        Instance: {
                            ConfigPaths: [
                                'config/daemon.json'
                            ],
                            RunUser: '',
                            Cmd: './bin/start.sh',
                            PostDeploy: 'sh post_deploy.sh docker',
                            PostUndeploy: '',
                            PrometheusPort: ''
                        },
                        Group: 'default',
                        Config: {
                            docker_user: {
                                Default: 'docker',
                                Desc: 'internal',
                                Type: 'internal',
                                Value: 'docker'
                            },
                            log_driver: {
                                Default: 'json-file',
                                Desc: 'internal',
                                Type: 'internal',
                                Value: 'json-file'
                            }
                        },
                        BaseProduct: '',
                        BaseProductVersion: '',
                        BaseService: '',
                        BaseParsed: false,
                        BaseAtrribute: '',
                        ServiceAddr: {
                            Host: [
                                '172-16-100-233'
                            ],
                            IP: [
                                '172.16.100.233'
                            ],
                            NodeId: 1,
                            SingleIndex: 0
                        }
                    },
                    rke: {
                        ServiceDisplay: '',
                        Version: '1.0.0',
                        Instance: {
                            ConfigPaths: [
                                'config/cluster.yml',
                                'post_deploy.sh',
                                'config/node_list'
                            ],
                            Logs: [
                                'logs/*.log'
                            ],
                            RunUser: '',
                            Cmd: './bin/waiting.sh',
                            PostDeploy: './post_deploy.sh docker',
                            PostUndeploy: '',
                            PrometheusPort: '',
                            MaxReplica: '1'
                        },
                        Group: 'default',
                        DependsOn: [
                            'docker'
                        ],
                        Config: {
                            cluster_ip: {
                                Default: {
                                    Host: [
                                        '172-16-100-233'
                                    ],
                                    IP: [
                                        '172.16.100.233'
                                    ],
                                    NodeId: 1,
                                    SingleIndex: 0
                                },
                                Desc: 'internal',
                                Type: 'internal',
                                Value: {
                                    Host: [
                                        '172-16-100-233'
                                    ],
                                    IP: [
                                        '172.16.100.233'
                                    ],
                                    NodeId: 1,
                                    SingleIndex: 0
                                }
                            },
                            docker_user: {
                                Default: 'docker',
                                Desc: 'internal',
                                Type: 'internal',
                                Value: 'docker'
                            }
                        },
                        BaseProduct: '',
                        BaseProductVersion: '',
                        BaseService: '',
                        BaseParsed: false,
                        BaseAtrribute: '',
                        ServiceAddr: {
                            Host: [
                                '172-16-100-233'
                            ],
                            IP: [
                                '172.16.100.233'
                            ],
                            NodeId: 1,
                            SingleIndex: 0
                        }
                    }
                }
            },
            product_name: 'DTK8S',
            product_version: '1.0.5',
            status: 'deployed'
        }
    }
}
