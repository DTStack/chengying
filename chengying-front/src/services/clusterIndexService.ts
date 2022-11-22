import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { clusterIndex } = apis;

export default {
  // 获取主机集群总览信息
  getHostsClusterOverview(params: any) {
    return http[clusterIndex.getHostsClusterOverview(params).method](
      clusterIndex.getHostsClusterOverview(params).url
    );
  },

  // 获取kubernetes集群总览信息
  getKubernetesClusterOverview(params: any) {
    return http[clusterIndex.getKubernetesClusterOverview(params).method](
      clusterIndex.getKubernetesClusterOverview(params).url
    );
  },

  // 获取CPU,MEMERRY等容量趋势图 - hosts
  getHostsClusterPerformance(id: number, params: any) {
    return http[clusterIndex.getHostsClusterPerformance(id).method](
      clusterIndex.getHostsClusterPerformance(id).url,
      params
    );
  },

  // 获取CPU,MEMERRY等容量趋势图 - k8s
  getKubernetesClusterPerformance(id: number, params: any) {
    return http[clusterIndex.getKubernetesClusterPerformance(id).method](
      clusterIndex.getKubernetesClusterPerformance(id).url,
      params
    );
  },

  // 趋势图
  getClusterPerformance(id: number, params: any, type: string) {
    return type === 'hosts'
      ? this.getHostsClusterPerformance(id, params)
      : this.getKubernetesClusterPerformance(id, params);
  },

  // 总览信息
  getClusterOverview(params: any, type: string) {
    return type === 'hosts'
      ? this.getHostsClusterOverview(params)
      : this.getKubernetesClusterOverview(params);
  },
};
