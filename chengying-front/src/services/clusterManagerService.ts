import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { clusterManager } = apis;

export default {
  // 获取集群列表
  getClusterLists(params: any) {
    return http[clusterManager.getClusterLists.method](
      clusterManager.getClusterLists.url,
      params
    );
  },
  /* -- 主机集群 -- */
  // 创建主机集群
  createHostCluster(params: any) {
    return http[clusterManager.createHostCluster.method](
      clusterManager.createHostCluster.url,
      params
    );
  },
  // 编辑主机集群
  updateHostCluster(params: any) {
    return http[clusterManager.updateHostCluster.method](
      clusterManager.updateHostCluster.url,
      params
    );
  },
  // 删除主机集群
  deleteHostCluster(params: any) {
    return http[clusterManager.deleteHostCluster(params).method](
      clusterManager.deleteHostCluster(params).url
    );
  },
  // 获取主机集群详情
  getHostClusterInfo(params: any) {
    return http[clusterManager.getHostClusterInfo(params).method](
      clusterManager.getHostClusterInfo(params).url
    );
  },
  /* -- kubernetes集群 -- */
  // 创建kubernetes集群
  createKubernetesCluster(params: any) {
    return http[clusterManager.createKubernetesCluster.method](
      clusterManager.createKubernetesCluster.url,
      params
    );
  },
  // 更新kubernetes集群
  updateKubernetesCluster(params: any) {
    return http[clusterManager.updateKubernetesCluster.method](
      clusterManager.updateKubernetesCluster.url,
      params
    );
  },
  // 删除kubernetes集群
  deleteKubernetesCluster(params: any) {
    return http[clusterManager.deleteKubernetesCluster(params).method](
      clusterManager.deleteKubernetesCluster(params).url
    );
  },
  // 获取kubernetes集群详情
  getKubernetesClusterInfo(params: any) {
    return http[clusterManager.getKubernetesClusterInfo(params).method](
      clusterManager.getKubernetesClusterInfo(params).url
    );
  },
  // 获取自建kubernetes版本，网络组件信息
  getKubernetesAvaliable(params: any) {
    return http[clusterManager.getKubernetesAvaliable.method](
      clusterManager.getKubernetesAvaliable.url,
      params
    );
  },
  // 获取导入kubernetes集群接入信息
  getKubernetesInstallCmd(params: any) {
    return http[clusterManager.getKubernetesInstallCmd.method](
      clusterManager.getKubernetesInstallCmd.url,
      params
    );
  },
  // 获取yaml模板
  getKubernetesRketemplate(params: any) {
    return http[clusterManager.getKubernetesRketemplate.method](
      clusterManager.getKubernetesRketemplate.url,
      params
    );
  },
  /* -- 合并 -- */
  // 获取集群详情
  getClusterInfo(params: any, type: string) {
    return type === 'hosts'
      ? this.getHostClusterInfo(params)
      : this.getKubernetesClusterInfo(params);
  },
  // 删除集群
  deleteCluster(params: any, type: string) {
    return type === 'hosts'
      ? this.deleteHostCluster(params)
      : this.deleteKubernetesCluster(params);
  },
  // 集群创建，更新操作
  clusterSubmitOperate(params: any, type: string, isEdit: boolean) {
    if (type === 'hosts') {
      return isEdit
        ? this.updateHostCluster(params)
        : this.createHostCluster(params);
    }
    return isEdit
      ? this.updateKubernetesCluster(params)
      : this.createKubernetesCluster(params);
  },
};
