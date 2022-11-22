export default {
    getClusterHostList: {
        msg: 'ok',
        code: 0,
        data: {
            count: 2,
            hosts: [
                {
                    id: 20,
                    sid: 'f06edaad-3f20-4c79-a901-f890b704949c',
                    hostname: '172-16-8-135',
                    ip: '172.16.8.135',
                    status: 7,
                    steps: 6,
                    errorMsg: 'K8S NODE部署成功',
                    updated: '2020-09-04 19:43:20',
                    created: '2020-08-11 17:32:28',
                    group: 'default',
                    product_name_list: 'DTCommon,DTConsole,DTFront,DTSchedule,DTStream',
                    product_name_display_list: 'DTCommon,DTConsole,DTFront,DTSchedule,DTStream',
                    pid_list: '140,156,159,160,162',
                    mem_size: 8201072640,
                    mem_usage: 1738571776,
                    disk_usage: {
                        String: '[{"mountPoint":"/","usedSpace":9247059968,"totalSpace":37688381440},{"mountPoint":"/boot","usedSpace":151760896,"totalSpace":1063256064},{"mountPoint":"/data","usedSpace":24361082880,"totalSpace":316931899392},{"mountPoint":"/var / lib / kubelet / pods / 41330959 - 2778 - 437b- 9e62 - 896f659cc9b3/ volume - subpaths / volume - dtinsight - dtfront - portalfront / portalfront / 0","usedSpace":9247059968,"totalSpace":37688381440},{"mountPoint":" /var/lib/kubelet / pods / 41330959 - 2778 - 437b- 9e62 - 896f659cc9b3/ volume - subpaths / volume - dtinsight - dtfront - portalfront / portalfront / 1","usedSpace":9247059968,"totalSpace":37688381440},{"mountPoint":"/var/lib / kubelet / pods / 41330959 - 2778 - 437b- 9e62 - 896f659cc9b3 / volume - subpaths / volume - dtinsight - dtfront - portalfront / portalfront / 2","usedSpace":9247059968,"totalSpace":37688381440},{"mountPoint":" /var/lib/kubelet / pods / 41330959 - 2778 - 437b - 9e62 - 896f659cc9b3 / volume - subpaths / volume - dtinsight - dtfront - portalfront / portalfront / 3","usedSpace":9247059968,"totalSpace":37688381440},{"mountPoint":" /var/lib/kubelet / pods / 41330959 - 2778 - 437b - 9e62 - 896f659cc9b3 / volume - subpaths / volume - dtinsight - dtfront - portalfront / portalfront / 4","usedSpace":9247059968,"totalSpace":37688381440}]',
                        Valid: true
                    },
                    net_usage: {
                        String: '[{"ifName":"cni0","ifIp":["10.42.0.1 / 24"],"bytesSent":11725873102,"bytesRecv":5968664938},{"ifName":"vethb1779a4e","bytesSent":97094125,"bytesRecv":98054758},{"ifName":"veth4b049412","bytesSent":542948562,"bytesRecv":85734760},{"ifName":"vethd758b876","bytesSent":1085647,"bytesRecv":1933833},{"ifName":"eth0","ifIp":["172.16.8.135 / 24"],"bytesSent":30733833681,"bytesRecv":71540251259},{"ifName":"veth95beba8e","bytesSent":3177385,"bytesRecv":315347},{"ifName":"veth99f45d61","bytesSent":1826354952,"bytesRecv":616914014},{"ifName":"veth9b8750cf","bytesSent":1176},{"ifName":"vethe2701f18","bytesSent":24771878,"bytesRecv":29305489},{"ifName":"veth898cf0b8","bytesSent":1333003124,"bytesRecv":345009458},{"ifName":"veth12cf3921","bytesSent":17189819,"bytesRecv":16102434},{"ifName":"flannel.1","ifIp":["10.42.0.0 / 32"],"bytesSent":204213296,"bytesRecv":114715774},{"ifName":"docker0","ifIp":["172.17.0.1 / 16"],"bytesSent":1346,"bytesRecv":776}]',
                        Valid: true
                    },
                    mem_size_display: '7.64GB',
                    mem_used_display: '1.62GB',
                    disk_size_display: '506.76GB',
                    disk_used_display: '74.50GB',
                    file_size_display: '471.66GB',
                    file_used_display: '65.89GB',
                    cpu_core_size_display: '4core',
                    cpu_core_used_display: '0.14core',
                    is_running: false,
                    cpu_usage_pct: 3.59,
                    mem_usage_pct: 0.212,
                    disk_usage_pct: 0.25,
                    pod_used_display: '0个',
                    pod_size_display: '0个',
                    pod_usage_pct: 0,
                    roles: {
                        Control: true,
                        Etcd: true,
                        Worker: true
                    },
                    run_user: ''
                }
            ]
        }
    },
    getClusterhostGroupLists: {
        msg: 'ok',
        code: 0,
        data: ['default']
    },
    getParentProductList: {
        msg: 'ok',
        code: 0,
        data: ['DTinsight', 'DTEM']
    },
    getClusterProductList: {
        msg: 'ok',
        code: 0,
        data: [{
            clusterId: 46,
            clusterName: 'DTLogger',
            clusterType: 'hosts',
            mode: 0,
            subdomain: {
                products: [
                    'DTEM'
                ]
            }
        },
        {
            clusterId: 58,
            clusterName: 'shaseng_k8s',
            clusterType: 'kubernetes',
            mode: 1,
            subdomain: {
                'DTinsight': [
                    'DTEM',
                    'DTUic'
                ]
            }
        }]
    },
    getAllProducts: {
        msg: 'ok',
        code: 0,
        data: {
            count: 1,
            list: [
                {
                    create_time: '2020-08-14 11:16:46',
                    deploy_time: '',
                    deploy_uuid: '',
                    id: 162,
                    namespace: '',
                    product_name: 'DTCommon',
                    product_name_display: 'DTCommon',
                    product_type: 1,
                    product_version: '4.0.15_110_beta',
                    status: '',
                    username: ''
                }
            ]
        }
    },
    getProductUpdateRecords: {
        msg: 'ok',
        code: 0,
        data: {
            count: 22,
            list: [
                {
                    create_time: '2020-07-30 15:42:36',
                    deploy_uuid: '24ccb9ec-c155-48d5-9e65-bf49a6a95edb',
                    id: 862,
                    namespace: '',
                    product_name: 'DTLogger',
                    product_name_display: '',
                    product_type: 0,
                    product_version: '4.0.4',
                    status: 'deployed',
                    username: 'admin@dtstack.com'
                },
                {
                    create_time: '2020-07-30 14:42:51',
                    deploy_uuid: '8004636b-60df-44ca-84c7-4b6d864f7337',
                    id: 860,
                    namespace: '',
                    product_name: 'DTLogger',
                    product_name_display: '',
                    product_type: 0,
                    product_version: '4.0.4',
                    status: 'deploy fail',
                    username: 'admin@dtstack.com'
                }
            ]
        }
    },
    getDeployShot: {
        msg: 'ok',
        code: 0,
        data: {
            complete: 'deployed',
            count: 1,
            list: [
                {
                    create_time: '2020-07-30 15:42:36',
                    deploy_uuid: '24ccb9ec-c155-48d5-9e65-bf49a6a95edb',
                    group: 'default',
                    id: 2969,
                    instance_id: 2352,
                    ip: '172.16.101.121',
                    product_name: 'DTLogger',
                    product_version: '4.0.4',
                    progress: 100,
                    schema: '{}',
                    service_name: 'flinkx',
                    service_version: '6.1.2',
                    sid: 'd8aa39a3-96ed-48a0-b777-0aaf4d16bc99',
                    status: 'health-checked',
                    status_message: '',
                    update_time: '2020-07-30 15:42:58'
                }
            ]
        }
    }
}
