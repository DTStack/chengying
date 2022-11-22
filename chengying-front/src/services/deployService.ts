import apis from '@/constants/apis';
import * as http from '@/utils/http';
import { SearchDeployLogs } from '@/model/apis';

const { deployAPI, unDeployAPI, installGuide } = apis;

export default {
  getProductConfig(params: any) {
    return http[deployAPI.getProductConfig(params).method](
      deployAPI.getProductConfig(params).url
    );
  },
  motifyServiceConfig(params: any, e: any) {
    return http[installGuide.setParamConfigFieldValue(params).method](
      installGuide.setParamConfigFieldValue(params).url,
      e
    );
  },
  resetSchemaField(params: any, e: any) {
    return http[installGuide.resetParamConfigFieldValue(params).method](
      installGuide.resetParamConfigFieldValue(params).url,
      e
    );
  },
  getIp(params: any) {
    return http[deployAPI.getIp(params).method](deployAPI.getIp(params).url);
  },
  setIp(params: any, e: any) {
    return http[installGuide.setIp(params).method](
      installGuide.setIp(params).url,
      e
    );
  },
  deploy(params: any) {
    return http[deployAPI.deploy(params).method](deployAPI.deploy(params).url);
  },
  getDeployList(params: any) {
    return http[deployAPI.getDeployList(params).method](
      deployAPI.getDeployList(params).url
    );
  },
  cancelDeploy(params: any) {
    return http[deployAPI.cancelDeploy(params).method](
      deployAPI.cancelDeploy(params).url
    );
  },
  getLeftIp(params: any) {
    return http[installGuide.updateServiceHostList(params).method](
      installGuide.updateServiceHostList(params).url
    );
  },
  // 卸载服务
  unDeployService(urlParams: any, params: any) {
    return http[unDeployAPI.unDeploy(urlParams).method](
      unDeployAPI.unDeploy(urlParams).url,
      params
    );
  },
  getUnDeployList(params: any) {
    return http[unDeployAPI.getUnDeployList(params).method](
      unDeployAPI.getUnDeployList(params).url
    );
  },
  forceStop(params: any) {
    return http[unDeployAPI.forceStop(params).method](
      unDeployAPI.forceStop(params).url
    );
  },
  forceUninstall(params: any) {
    return http[unDeployAPI.forceUninstall(params).method](
      unDeployAPI.forceUninstall(params).url
    );
  },
  getUnDeployLog(params: any) {
    return http[unDeployAPI.getUnDeployLog(params).method](
      unDeployAPI.getUnDeployLog(params).url
    );
  },
  searchDeployLog(params: SearchDeployLogs) {
    const req = deployAPI.searchDeployLog(params);
    return http[req.method](req.url);
  },

  // 停止部署
  stopUndeploy(params: any) {
    return http[unDeployAPI.stopUndeploy(params).method](
      unDeployAPI.stopUndeploy(params).url
    );
  },
  // 获取自动部署记录
  getOrchestrationHistory(params: any) {
    return http[deployAPI.getOrchestrationHistory.method](
      deployAPI.getOrchestrationHistory.url,
      params
    );
  },

  // 备份
  handleBackUp(urlParams: any, params: any) {
    return http[deployAPI.handleBackUp(urlParams).method](
      deployAPI.handleBackUp(urlParams).url,
      params
    );
  },


  // 备份
  handleRollBack(urlParams: any, params: any) {
    return http[deployAPI.handleRollBack(urlParams).method](
      deployAPI.handleRollBack(urlParams).url,
      params
    );
  },

  // 回滚版本列表
  getRollBackList(urlParams: any, params: any) {
    return http[deployAPI.getRollBackList(urlParams).method](
      deployAPI.getRollBackList(urlParams).url,
      params
    );
  },

  getBackupTimes(urlParams: any, params: any) {
    return http[deployAPI.getBackupTimes(urlParams).method](
      deployAPI.getBackupTimes(urlParams).url,
      params
    );
  },
};
