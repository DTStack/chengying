import apis from '@/constants/apis';
import * as http from '@/utils/http';

const { product, service, instance, host, alert, cluster, common } = apis;

export const Service = {
  addNetworkPackage(params: any) {
    return http[product.addNetworkPackage.method](
      product.addNetworkPackage.url,
      params
    );
  },
  getProductpackageList() {
    return http[product.getProductpackageList.method](
      product.getProductpackageList.url
    );
  },
  delProductPackageItem(params: any) {
    return http[product.delProductpackageItem.method](
      product.delProductpackageItem.url,
      params
    );
  },
  getAllProducts(params: any) {
    return http[product.getProductList.method](
      product.getProductList.url,
      params
    );
  },
  Verifylink(params: any) {
    return http[product.Verifylink.method](product.Verifylink.url, params);
  },
  // 获取更新目录
  getPatchPath(params: any) {
    return http[product.getPatchPath.method](product.getPatchPath.url, params);
  },
  // 上传更新接口
  updatePatchPath(params: any) {
    return http[product.updatePatchPath.method](
      product.updatePatchPath.url,
      params
    );
  },
  // 获取更新状态
  updatePatchStatus(params: any) {
    return http[product.updatePatchStatus.method](
      product.updatePatchStatus.url,
      params
    );
  },
  getParentProductList(params?: any) {
    return http[product.getParentProductList.method](
      product.getParentProductList.url,
      params
    );
  },
  // 获取产品列表
  getClusterProductList(params?: any) {
    return http[product.getClusterProductList.method](
      product.getClusterProductList.url,
      params
    );
  },
  getCurrentProduct(params: any) {
    return http[product.getCurrentProduct(params).method](
      product.getCurrentProduct(params).url,
      params
    );
  },
  getProductUpdateRecords(params: any) {
    return http[product.getProductUpdateRecords(params).method](
      product.getProductUpdateRecords(params).url,
      params
    );
  },
  // 获取补丁包更新历史
  getProductUpdate(params: any) {
    return http[product.getProductUpdate(params).method](
      product.getProductUpdate(params).url,
      params
    );
  },
  // 获取产品包升级目标版本列表
  getProductVersionList(params: any) {
    const { upgrade_mode } = params;
    return http[product.getProductVersionList(params).method](
      product.getProductVersionList(params).url,
      {
        upgrade_mode,
      }
    );
  },
  getProductDeployRecords(params: any) {
    return http[product.getProductDeployRecords(params).method](
      product.getProductDeployRecords(params).url,
      params
    );
  },
  getServiceDashInfo(params: any) {
    const url = service.getServiceDashInfo.url;
    const { product_name, service_name } = params;
    return http[service.getServiceDashInfo.method](
      `${url}?tag=${product_name}&tag=${service_name}`
    );
  },
  getInstanceLog(params: any) {
    return http[instance.getInstanceLog(params).method](
      instance.getInstanceLog(params).url,
      params
    );
  },
  getFilterData(params: any) {
    return http[host.getFilterData.method](host.getFilterData.url, params);
  },
  // 获取主机集群下的主机列表
  getHostClusterHostList(params: any) {
    return http[common.getHostClusterHostList.method](
      common.getHostClusterHostList.url,
      params
    );
  },
  // 获取kubernetes集群下的主机列表
  getKubernetesClusterHostList(params: any) {
    return http[common.getKubernetesClusterHostList.method](
      common.getKubernetesClusterHostList.url,
      params
    );
  },
  // 获取主机列表
  getClusterHostList(params: any, type: string) {
    return type === 'hosts'
      ? this.getHostClusterHostList(params)
      : this.getKubernetesClusterHostList(params);
  },
  // 获取集群下主机分组列表
  getClusterhostGroupLists(params: any) {
    return http[common.getClusterhostGroupLists.method](
      common.getClusterhostGroupLists.url,
      params
    );
  },
  getTableData(params: any) {
    return http[host.getTableData.method](host.getTableData.url, params);
  },
  showInstallInfo(params: any) {
    return http[host.showInstallInfo.method](host.showInstallInfo.url, params);
  },
  getAlertsHistory(params: any) {
    return http[alert.getAlertsHistory(params).method](
      alert.getAlertsHistory(params).url,
      params
    );
  },
  getHostInstance(params: any) {
    return http[host.getHostInstance.method](host.getHostInstance.url, params);
  },
  getHostServicesList(params: any) {
    return http[host.getHostServicesList.method](
      host.getHostServicesList.url,
      params
    );
  },
  getDeployShot(params: any) {
    return http[product.getDeployShot(params).method](
      product.getDeployShot(params).url
    );
  },
  getPatchUpdateShot(params: any) {
    return http[product.getPatchUpdateShot(params).method](
      product.getPatchUpdateShot(params).url
    );
  },
  delProduct(params: any) {
    return http[product.delProduct(params).method](
      product.delProduct(params).url
    );
  },
  getClusterList(params: any) {
    return http[cluster.getClusterList.method](
      cluster.getClusterList.url,
      params
    );
  },
};

export { default as productService } from './productService';
export { default as dashboardService } from './dashboardService';
export { default as alertRuleService } from './alertRuleService';
export { default as alertChannelService } from './alertChannelService';
export { default as installGuideService } from './installGuideService';
export { default as addHostService } from './addHostService';
export { default as deployService } from './deployService';
export { default as servicePageService } from './ServicePageService';
export { default as createHostService } from './createHostService';
export { default as userCenterService } from './userCenterService';
export { default as eventService } from './eventService';
export { default as configService } from './configService';
export { default as clusterManagerService } from './clusterManagerService';
export { default as clusterHostService } from './clusterHostService';
export { default as imageStoreService } from './imageStoreService';
export { default as SecurityAuditService } from './securityAuditService';
export { default as ClusterNamespaceService } from './clusterNamespaceService';
export { default as inspectionReportService } from './inspectionReportService';
export { default as echoService } from './echoService';
export { default as securityService } from './securityService';
export { default as scriptManager } from './scriptManager';
export { default as globalConfig } from './globalConfig';
export { default as productLine } from './productLine';