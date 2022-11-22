import apis from '@/constants/apis';
import * as http from '@/utils/http';

const { installGuide } = apis;

export default {
  getProductStepOneList(params: any) {
    return http[installGuide.getProductStepOneList.method](
      installGuide.getProductStepOneList.url,
      params
    );
  },
  checkMySqlAddr(params: any) {
    const {cluster_id, final_upgrade, ip} = params;
    return http[installGuide.checkMySqlAddr(params).method](
      installGuide.checkMySqlAddr(params).url,
      {
        cluster_id,
        final_upgrade,
        ip,
      }
    );
  },
  getGlobalAutoConfig(params: any) {
    return http[installGuide.getGlobalAutoConfig(params).method](
      installGuide.getGlobalAutoConfig(params).url,
      params
    );
  },
  getAutoConfig(params: any) {
    return http[installGuide.getAutoConfig(params).method](
      installGuide.getAutoConfig(params).url,
      params
    );
  },
  getProductNames(params?: any) {
    return http[installGuide.getFilters.method](
      installGuide.getFilters.url,
      params
    );
  },
  getProductPackagelist(params: any) {
    return http[installGuide.getProductPackageList.method](
      installGuide.getProductPackageList.url,
      params
    );
  },
  getProductPackageServices(params: any) {
    return http[installGuide.getProductPackageServices(params).method](
      installGuide.getProductPackageServices(params).url,
      { relynamespace: params.relynamespace }
    );
  },
  deleteProducyPackage(params: any) {
    return http[installGuide.deleteProductPackage(params).method](
      installGuide.deleteProductPackage(params).url
    );
  },
  getInstallCMD(params?: any) {
    return http[installGuide.getInstallCmd.method](
      installGuide.getInstallCmd.url,
      params
    );
  },
  getProductServicesInfo(params: any) {
    const { namespace, unSelectService, relynamespace, clusterId, upgrade_mode } = params;
    return http[installGuide.getProductServicesInfo(params).method](
      installGuide.getProductServicesInfo(params).url,
      {
        namespace,
        unchecked_services: unSelectService,
        relynamespace,
        clusterId,
        upgrade_mode,
      }
    );
  },
  // 获取集群列表 - 复用集群列表
  getInstallClusterList(params?: any) {
    return http[installGuide.getInstallClusterList.method](
      installGuide.getInstallClusterList.url,
      params
    );
  },
  updateServiceHostList(params: any) {
    return http[installGuide.updateServiceHostList(params).method](
      installGuide.updateServiceHostList(params).url
    );
  },
  setParamConfigFieldValue(
    params: { productName: string; serviceName: string },
    conf: any
  ) {
    return http[installGuide.setParamConfigFieldValue(params).method](
      installGuide.setParamConfigFieldValue(params).url,
      conf
    );
  },
  setParamConfigFieldValueTotal(
    params: { productName: string; serviceName: string },
    p: any
  ) {
    return http[installGuide.setParamConfigFieldValueTotal(params).method](
      installGuide.setParamConfigFieldValueTotal(params).url,
      p
    );
  },
  // 运行部署配置文件
  modifyMultiSingleField(
    params: { productName: string; serviceName: string },
    p: any
  ) {
    return http[installGuide.modifyMultiSingleField(params).method](
      installGuide.modifyMultiSingleField(params).url,
      p
    );
  },
  resetParamConfigFieldValue(params: {
    field_path: string;
    productName: string;
    serviceName: string;
    pid: number;
    product_version: string;
    type: string;
    hosts: string;
  }) {
    if (params.type == '1') {
      const p = {
        field_path: params.field_path,
        pid: params.pid,
        product_version: params.product_version,
      };
      return http[installGuide.resetParamConfigFieldValue(params).method](
        installGuide.resetParamConfigFieldValue(params).url,
        p
      );
    } else {
      const p = {
        field_path: params.field_path,
        pid: params.pid,
        product_version: params.product_version,
        hosts: params.hosts,
      };
      return http[installGuide.resetMultiConfigFieldValue(params).method](
        installGuide.resetMultiConfigFieldValue(params).url,
        p
      );
    }
  },
  startDeploy(params: any) {
    return http[installGuide.deploy(params).method](
      installGuide.deploy(params).url,
      params
    );
  },
  checkDeployStatus(params: any) {
    return http[installGuide.checkDeployStatus(params).method](
      installGuide.checkDeployStatus(params).url
    );
  },
  getDeployTaskList(params: any) {
    return http[installGuide.getDeployTaskList(params).method](
      installGuide.getDeployTaskList(params).url
    );
  },
  stopDeploy(params: any) {
    return http[installGuide.stopDeploy(params).method](
      installGuide.stopDeploy(params).url,
      params
    );
  },
  stopAutoDeploy(params: any) {
    return http[installGuide.stopAutoDeploy.method](
      installGuide.stopAutoDeploy.url,
      params
    );
  },
  modifyProductConfigAll(params: any, p: any) {
    return http[installGuide.modifyProductConfigAll(params).method](
      installGuide.modifyProductConfigAll(params).url,
      p
    );
  },
  setIp(params: any, p: any) {
    return http[installGuide.setIp(params).method](
      installGuide.setIp(params).url,
      p
    );
  },
  getUncheckedService(params: any) {
    const { clusterId, namespace } = params;
    return http[installGuide.getUncheckedService(params).method](
      installGuide.getUncheckedService(params).url,
      { clusterId, namespace }
    );
  },
  serviceUpdate(params: any, p: any) {
    return http[installGuide.serviceUpdate(params).method](
      installGuide.serviceUpdate(params).url,
      p
    );
  },
  // 创建命名空间
  createNamespace(urlParams: any, params: any) {
    return http[installGuide.createNamespace(urlParams).method](
      installGuide.createNamespace(urlParams).url,
      params
    );
  },
  // 获取命名空间列表
  getNamespaceList(params: any) {
    return http[installGuide.getNamespaceList(params).method](
      installGuide.getNamespaceList(params).url
    );
  },
  // 产品依赖检测
  getBaseClusterList(params: any) {
    return http[installGuide.getBaseClusterList(params).method](
      installGuide.getBaseClusterList(params).url
    );
  },
  // 检查默认仓库是否存在
  checkDefaultImageStore(params: any) {
    return http[installGuide.checkDefaultImageStore(params).method](
      installGuide.checkDefaultImageStore(params).url
    );
  },
  // 自动编排
  autoOrchestration(params: any) {
    return http[installGuide.autoOrchestration.method](
      installGuide.autoOrchestration.url,
      params
    );
  },
  // 自动部署
  autoDeploy(params: any) {
    return http[installGuide.autoDeploy.method](
      installGuide.autoDeploy.url,
      params
    );
  },
  getHostRoleMap(params: any) {
    return http[installGuide.getHostRoleMap.method](
      installGuide.getHostRoleMap.url,
      params
    );
  },
  refreshServicesInfoForAutoDeploy(params: any) {
    return http[installGuide.refreshServicesInfoForAutoDeploy.method](
      installGuide.refreshServicesInfoForAutoDeploy.url,
      params
    );
  },

  getOldHostInfo(urlParams: any, params: any) {
    return http[installGuide.getOldHostInfo(urlParams).method](
      installGuide.getOldHostInfo(urlParams).url,
      params
    );
  },

  getServiceGroupFile(urlParams: any, params: any) {
    return http[installGuide.getServiceGroupFile(urlParams).method](
      installGuide.getServiceGroupFile(urlParams).url,
      params
    );
  },

  // 保存升级记录
  saveUpgrade(urlParams: any, params: any) {
    return http[installGuide.saveUpgrade(urlParams).method](
      installGuide.saveUpgrade(urlParams).url,
      params
    );
  },
  // 检查部署条件接口
deployCondition(params: any) {
  return http[installGuide.deployCondition.method](
    installGuide.deployCondition.url,
    params
  );
 }
};