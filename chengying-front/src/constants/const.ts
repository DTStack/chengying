// 部署状态
export const deployStatus = {
};

// 部署状态filter
export const deployStatusFilter = [
  {
    text: '部署成功',
    value: 'deployed',
  },
  {
    text: '部署失败',
    value: 'deploy fail',
  },
  {
    text: '部署中',
    value: 'deploying',
  },
  {
    text: '卸载失败',
    value: 'undeploy fail',
  },
  {
    text: '卸载中',
    value: 'undeploying',
  },
];

// 更新状态filter
export const updateStatusFilter = [
  {
    text: '更新成功',
    value: 'success',
  },
  {
    text: '更新失败',
    value: 'fail',
  },
  {
    text: '更新中',
    value: 'update',
  },
];

export const hostsInstallStatus = [
  { text: '管控安装成功', value: '管控安装成功' },
  { text: '管控安装失败', value: '管控安装失败' },
  { text: 'script安装成功', value: 'script安装成功' },
  { text: 'script安装失败', value: 'script安装失败' },
  { text: '主机初始化成功', value: '主机初始化成功' },
  { text: '主机初始化失败', value: '主机初始化失败' },
];

export const kubernetesInstallStatus = [
  { text: '管控安装失败', value: '管控安装失败', status: -1 },
  { text: '管控安装成功', value: '管控安装成功', status: 1 },
  { text: '主机初始化失败', value: '主机初始化失败', status: -3 },
  { text: '主机初始化成功', value: '主机初始化成功', status: 3 },
  { text: 'K8S DOCKER初始化失败', value: 'K8S DOCKER初始化失败', status: -5 },
  { text: 'K8S DOCKER初始化成功', value: 'K8S DOCKER初始化成功', status: 5 },
  { text: 'K8S NODE初始化失败', value: 'K8S NODE初始化失败', status: -6 },
  { text: 'K8S NODE初始化成功', value: 'K8S NODE初始化成功', status: 6 },
  { text: 'K8S NODE部署失败', value: 'K8S NODE部署失败', status: -7 },
  { text: 'K8S NODE部署成功', value: 'K8S NODE部署成功', status: 7 },
];

export const clusterTypeMap = {
  hosts: ['主机集群'],
  kubernetes: ['Kubernetes集群（自建集群）', 'Kubernetes集群（导入已有集群）'],
};

// k8s & hosts 主机初始化状态
export const hostStatusInfoMap = {
  hosts: {
    title:
      '主机安装会经过2个状态：管控安装-主机初始化，若某一步安装失败，请查看具体日志。',
    filters: hostsInstallStatus,
    finalStatus: 3,
  },
  kubernetes: {
    title:
      'Kubernetes集群主机安装会经过5个状态：管控安装-主机初始化-K8S Docker初始化-K8S Node初始化-K8S Node部署成功，若某一步安装失败，请查看具体日志',
    filters: kubernetesInstallStatus,
    finalStatus: 7,
  },
};

export const COOKIES = {
  NAMESPACE: 'em_current_k8s_namespace',
};
