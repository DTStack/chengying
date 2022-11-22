import apis from '@/constants/apis';
import * as http from '@/utils/http';
import cloneDeep from 'lodash/cloneDeep';

const { product, service, alert } = apis;

export interface ProductionListParams {
  limit?: number;
  parentProductName?: string;
  namespace?: string;
  clusterId?: number;
  mode?: number;
}

export interface CurrentProductionParams {
  product_name: string;
  namespace?: string;
}

export interface InstanceReplicaParams {
  replica: number;
  namespace?: string;
}

export default {

  // 获取config diff
  getConfDiff(params: any) {
    const newParams = {
      file: params?.file,
      ip: params?.ip,
    };
    return http[service.getConfDiff(params).method](
      service.getConfDiff(params).url,
      newParams
    );
  },


  // 切换自动执行健康检查
  setAutoexecSwitch(params: any) {
    const newParams = {
      record_id: params?.record_id,
      auto_exec: params?.auto_exec,
    };
    return http[service.setAutoexecSwitch(params).method](
      service.setAutoexecSwitch(params).url,
      newParams
    );
  },
  // 手动执行健康检查
  manualExecution(params: any) {
    const newParams = {
      record_id: params?.record_id,
    };
    return http[service.manualExecution(params).method](
      service.manualExecution(params).url,
      newParams
    );
  },
  // 获取健康检查
  getHealthCheck(params: any) {
    const newParam = {
      ip: params.ip,
    };
    return http[service.getHealthCheck(params).method](
      service.getHealthCheck(params).url,
      newParam
    );
  },
  getServiceStatus(params: any) {
    return http[service.getServiceStatus.method](
      service.getServiceStatus.url,
      params
    );
  },

  getHostAlert(params: any) {
    return http[service.getHostAlert.method](service.getHostAlert.url, params);
  },

  getLogFileList(params: any) {
    return http[service.getLogFileList(params).method](
      service.getLogFileList(params).url
    );
  },
  getLogFile(params: any, reqParams: any) {
    return http[service.getLogFile(params).method](
      service.getLogFile(params).url,
      reqParams
    );
  },
  getLogFileDownLoad(params: any) {
    return http[service.getLogFileDownLoad(params).method](
      service.getLogFileDownLoad(params).url
    );
  },
  getParamDropList(params: any) {
    return http[service.getParamDropList(params).method](
      service.getParamDropList(params).url
    );
  },
  getParamContent(params: any) {
    return http[service.getParamContent(params).method](
      service.getParamContent(params).url
    );
  },
  setParamUpdate(params: any, config: any) {
    return http[service.setParamUpdate(params).method](
      service.setParamUpdate(params).url,
      config
    );
  },
  getServiceGroup(params: any, config?: any) {
    return http[service.getServiceGroup(params).method](
      service.getServiceGroup(params).url,
      config
    );
  },
  getAllServiceGroup(params: any, config?: any) {
    return http[service.getAllServiceGroup(params).method](
      service.getAllServiceGroup(params).url,
      config
    );
  },
  setServiceGroupStart(params: any) {
    return http[service.setServiceGroupStart(params).method](
      service.setServiceGroupStart(params).url
    );
  },

  setServiceGroupStop(params: any) {
    return http[service.setServiceGroupStop(params).method](
      service.setServiceGroupStop(params).url
    );
  },

  getServiceList(params: any) {
    return http[service.getServiceList(params).method](
      service.getServiceList(params).url
    );
  },

  getServiceHostsList(params: any) {
    const { namespace } = params;
    return namespace
      ? http[service.getServiceHostsList(params).method](
          service.getServiceHostsList(params).url,
          { namespace }
        )
      : http[service.getServiceHostsList(params).method](
          service.getServiceHostsList(params).url
        );
  },
  // 服务 - 运行配置
  getServiceConfig(params: any) {
    const request = {
      file: params.file === 'all' ? undefined : params.file,
    };
    return http[service.getServiceConfig(params).method](
      service.getServiceConfig(params).url,
      request
    );
  },

  getServiceHostConfig(params: any) {
    const { configfile } = params;
    return http[service.getServiceHostConfig(params).method](
      service.getServiceHostConfig(params).url,
      { configfile: configfile }
    );
  },

  startService(params: any) {
    return http[service.startService(params).method](
      service.startService(params).url
    );
  },

  stopService(params: any) {
    return http[service.stopService(params).method](
      service.stopService(params).url
    );
  },

  startServiceInstance(params: any) {
    return http[service.startServiceInstance(params).method](
      service.startServiceInstance(params).url
    );
  },

  stopServiceInstance(params: any) {
    return http[service.stopServiceInstance(params).method](
      service.stopServiceInstance(params).url
    );
  },

  getProductList(params: ProductionListParams) {
    return http[product.getProductList.method](
      product.getProductList.url,
      params
    );
  },

  getProductName(params: ProductionListParams) {
    return http[product.getProductName.method](
      product.getProductName.url,
      params
    );
  },

  resetServiceConfig(params: any) {
    const { field_path, pid, product_version } = params;
    return http[product.resetServiceConfig(params).method](
      product.resetServiceConfig(params).url,
      {
        field_path: field_path,
        pid,
        product_version,
      }
    );
  },
  resetMultiServiceConfig(params: any) {
    const { field_path, pid, product_version, hosts } = params;
    return http[product.resetMultiServiceConfig(params).method](
      product.resetMultiServiceConfig(params).url,
      {
        field_path: field_path,
        pid,
        product_version,
        hosts,
      }
    );
  },

  modifyProductConfigAll(params: any, config: any) {
    return http[product.modifyProductConfigAll(params).method](
      product.modifyProductConfigAll(params).url,
      config
    );
  },

  // 配置信息关联主机
  modifyMultiAllHosts(params: any, config: any) {
    return http[product.modifyMultiAllHosts(params).method](
      product.modifyMultiAllHosts(params).url,
      config
    );
  },
  getCurrentProduct(params: CurrentProductionParams) {
    return http[product.getCurrentProduct(params).method](
      product.getCurrentProduct(params).url,
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

  getHomeServiceDashInfo(params: any) {
    return http[service.getServiceDashInfo.method](
      service.getServiceDashInfo.url,
      params
    );
  },
  getServiceDashPlaylists() {
    return http[service.getServiceDashPlaylists.method](
      service.getServiceDashPlaylists.url
    );
  },
  getAlertsHistory(params: any) {
    const newParams = cloneDeep(params);
    delete newParams.product_name;
    delete newParams.service_name;
    return http[alert.getAlertsHistory(params).method](
      alert.getAlertsHistory(params).url,
      newParams
    );
  },

  setServiceRollRestart(params: any) {
    return http[service.setServiceRollRestart(params).method](
      service.setServiceRollRestart(params).url
    );
  },

  getServiceHARole(params: any) {
    return http[service.getServiceHARole(params).method](
      service.getServiceHARole(params).url
    );
  },

  startAllServiceByProduct(params: any) {
    return http[service.startAllServiceByProduct(params).method](
      service.startAllServiceByProduct(params).url
    );
  },
  stopAllServiceByProduct(params: any) {
    return http[service.stopAllServiceByProduct(params).method](
      service.stopAllServiceByProduct(params).url
    );
  },
  // k8s 扩缩容
  instanceReplica(params: any, conf: InstanceReplicaParams) {
    return http[service.instanceReplica(params).method](
      service.instanceReplica(params).url,
      conf
    );
  },
  // 服务异常数量列表
  getRedService(params: any) {
    return http[service.getRedService(params).method](
      service.getRedService(params).url,
      params
    );
  },
  // 重启服务数量列表
  getRestartService(params: any) {
    return http[service.getRestartService(params).method](
      service.getRestartService(params).url
    );
  },
  //
  getAutoTestStatus(params: any) {
    return http[service.getAutoTestStatus(params).method](
      service.getAutoTestStatus(params).url,
      {}
    );
  },
  //
  startAutoTestTest(params: any) {
    return http[service.startAutoTestTest(params).method](
      service.startAutoTestTest(params).url,
      params
    );
  },

  // 下载凭证文件
  operateExtension(params: any) {
    return http[service.operateExtension(params).method](
      service.operateExtension(params).url,
      {
        type: params.type,
        value: params.value,
      }
    );
  },
  // 操作开关启停
  operateSwitch(params: any) {
    return http[service.operateSwitch(params).method](
      service.operateSwitch(params).url,
      {
        type: params.type,
        name: params.name,
        product_version: params.product_version,
      }
    );
  },
  // 判断开关是否正在操作中
  getSwitchRecord: (params: any) => {
    return http[service.getSwitchRecord(params).method](
      service.getSwitchRecord(params).url,
      {
        name: params.name,
      }
    );
  },
  // 获取开关当前操作进度
  getSwitchDetail(params: any) {
    return http[service.getSwitchDetail(params).method](
      service.getSwitchDetail(params).url,
      {
        record_id: params.record_id,
      }
    );
  },
};
