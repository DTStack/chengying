export const clusterInfo = {
    id: 46,
    type: 'kubernetes',
    mode: 0,
    desc: '这是测试啊啊啊啊啊啊',
    name: 'DTLogger',
    tags: '日志中心,test',
    network_plugin: 'flannel',
    version: 'v1.16.3',
    yaml: 'nodes:\n- address: 172.16.8.135\n  port: \"22\"\n  internal_address: \"\"\n  role:\n  - etcd\n  - controlplane\n  - worker\n  hostname_override: \"\"\n  user: docker\n  docker_socket: \"\"\n  ssh_key: \"\"\n  ssh_key_path: ~/.ssh/id_rsa\n  ssh_cert: \"\"\n  ssh_cert_path: \"\"\n  labels: {}\n  taints: []\n- address: 172.16.10.152\n  port: \"22\"\n  internal_address: \"\"\n  role:\n  - etcd\n  - controlplane\n  - worker\n  hostname_override: \"\"\n  user: docker\n  docker_socket: \"\"\n  ssh_key: \"\"\n  ssh_key_path: ~/.ssh/id_rsa\n  ssh_cert: \"\"\n  ssh_cert_path: \"\"\n  labels: {}\n  taints: []\nservices:\n  etcd:\n    image: \"\"\n    extra_args: {}\n    extra_binds: []\n    extra_env: []\n    external_urls: []\n    ca_cert: \"\"\n    cert: \"\"\n    key: \"\"\n    path: \"\"\n    uid: 0\n    gid: 0\n    snapshot: null\n    retention: \"\"\n    creation: \"\"\n    backup_config: null\n  kube-api:\n    image: \"\"\n    extra_args: {}\n    extra_binds: []\n    extra_env: []\n    service_cluster_ip_range: \"\"\n    service_node_port_range: \"\"\n    pod_security_policy: false\n    always_pull_images: false\n    secrets_encryption_config: null\n    audit_log: null\n    admission_configuration: null\n    event_rate_limit: null\n  kube-controller:\n    image: \"\"\n    extra_args: {}\n    extra_binds: []\n    extra_env: []\n    cluster_cidr: \"\"\n    service_cluster_ip_range: \"\"\n  scheduler:\n    image: \"\"\n    extra_args: {}\n    extra_binds: []\n    extra_env: []\n  kubelet:\n    image: \"\"\n    extra_args: {}\n    extra_binds: []\n    extra_env: []\n    cluster_domain: \"\"\n    infra_container_image: \"\"\n    cluster_dns_server: \"\"\n    fail_swap_on: false\n    generate_serving_certificate: false\n  kubeproxy:\n    image: \"\"\n    extra_args: {}\n    extra_binds: []\n    extra_env: []\nnetwork:\n  plugin: flannel\n  options: {}\n  mtu: 0\n  node_selector: {}\n  update_strategy: null\nauthentication:\n  strategy: \"rbac\"\n  sans: []\n  webhook: null\naddons: \"\"\naddons_include: []\nsystem_images:\n  etcd: \"\"\n  alpine: \"\"\n  nginx_proxy: \"\"\n  cert_downloader: \"\"\n  kubernetes_services_sidecar: \"\"\n  kubedns: \"\"\n  dnsmasq: \"\"\n  kubedns_sidecar: \"\"\n  kubedns_autoscaler: \"\"\n  coredns: \"\"\n  coredns_autoscaler: \"\"\n  nodelocal: \"\"\n  kubernetes: \"\"\n  flannel: \"\"\n  flannel_cni: \"\"\n  calico_node: \"\"\n  calico_cni: \"\"\n  calico_controllers: \"\"\n  calico_ctl: \"\"\n  calico_flexvol: \"\"\n  canal_node: \"\"\n  canal_cni: \"\"\n  canal_flannel: \"\"\n  canal_flexvol: \"\"\n  weave_node: \"\"\n  weave_cni: \"\"\n  pod_infra_container: \"\"\n  ingress: \"\"\n  ingress_backend: \"\"\n  metrics_server: \"\"\n  windows_pod_infra_container: \"\"\nssh_key_path: \"\"\nssh_cert_path: \"\"\nssh_agent_auth: false\nauthorization:\n  mode: none\n  options: {}\nignore_docker_version: true\nkubernetes_version: v1.16.3\nprivate_registries: []\ningress:\n  provider: nginx\n  options: {}\n  node_selector: {}\n  extra_args: {}\n  dns_policy: \"\"\n  extra_envs: []\n  extra_volumes: []\n  extra_volume_mounts: []\n  update_strategy: null\ncluster_name: shaseng_k8s\ncloud_provider:\n  name: \"\"\nprefix_path: \"\"\naddon_job_timeout: 0\nbastion_host:\n  address: \"\"\n  port: \"\"\n  user: \"\"\n  ssh_key: \"\"\n  ssh_key_path: \"\"\n  ssh_cert: \"\"\n  ssh_cert_path: \"\"\nmonitoring:\n  provider: metrics-server\n  options: {}\n  node_selector: {}\n  update_strategy: null\n  replicas: null\nrestore:\n  restore: false\n  snapshot_name: \"\"\ndns: null\n'
}

export default {
    getClusterList: {
        msg: 'ok',
        code: 0,
        data: {
            clusters: [
                {
                    cpu_core_size_display: '6core',
                    cpu_core_used_display: '0.92core',
                    create_time: '0001-01-01T00:00:00Z',
                    create_user: 'admin@dtstack.com',
                    desc: 'EM2.0',
                    disk_size_display: '247.65GB',
                    disk_used_display: '103.01GB',
                    id: 1,
                    mem_size_display: '11.34GB',
                    mem_used_display: '4.92GB',
                    mode: 0,
                    name: 'dtstack',
                    nodes: 2,
                    status: 'Error',
                    tags: '',
                    type: 'hosts',
                    update_time: '2020-11-10T13:59:34+08:00',
                    update_user: 'admin@dtstack.com',
                    version: ''
                },
                {
                    cpu_core_size_display: '',
                    cpu_core_used_display: '',
                    create_time: '2020-06-19T19:31:39+08:00',
                    create_user: 'admin',
                    desc: '',
                    id: 11,
                    mem_size_display: '',
                    mem_used_display: '',
                    mode: 1,
                    name: 'k8s_1_12_9_import',
                    nodes: 0,
                    pod_size_display: '0个',
                    pod_used_display: '0个',
                    status: 'Running',
                    tags: '',
                    type: 'kubernetes',
                    update_time: '2020-11-10T13:59:34+08:00',
                    update_user: 'admin',
                    version: 'v1.12.9'
                },
                {
                    cpu_core_size_display: '12.00core',
                    cpu_core_used_display: '1.03core',
                    create_time: '2020-08-11T17:55:52+08:00',
                    create_user: 'admin',
                    desc: '',
                    id: 58,
                    mem_size_display: '23.15GB',
                    mem_used_display: '1.73GB',
                    mode: 0,
                    name: 'shaseng_k8s',
                    nodes: 2,
                    pod_size_display: '220个',
                    pod_used_display: '23个',
                    status: 'Running',
                    tags: '',
                    type: 'kubernetes',
                    update_time: '2020-11-10T13:59:34+08:00',
                    update_user: 'admin',
                    version: 'v1.16.3'
                },
                {
                    cpu_core_size_display: '',
                    cpu_core_used_display: '',
                    create_time: '2020-08-12T11:56:28+08:00',
                    create_user: 'admin',
                    desc: '',
                    id: 59,
                    mem_size_display: '',
                    mem_used_display: '',
                    mode: 1,
                    name: 'wx_test2',
                    nodes: 0,
                    pod_size_display: '0个',
                    pod_used_display: '0个',
                    status: 'Waiting',
                    tags: '',
                    type: 'kubernetes',
                    update_time: '2020-11-10T13:59:34+08:00',
                    update_user: 'admin',
                    version: 'v1.16.3'
                }
            ],
            counts: 20
        }
    },
    getClusterInfo: {
        msg: 'ok',
        code: 0,
        data: clusterInfo
    },
    getKubernetesAvaliable: {
        msg: 'ok',
        code: 0,
        data: [
            {
                properties: {
                    network_plugin: ['flannel']
                },
                version: 'v1.16.3'
            }
        ]
    },
    clusterSubmitOperate: {
        msg: 'ok',
        code: 0,
        data: clusterInfo
    }
}
