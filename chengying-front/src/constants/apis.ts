import {
  ServiceTree,
  ServiceFile,
  DistributeServiceConfig,
  SearchDeployLogs,
} from '@/model/apis';

export default {
  event: {
    getEventCount: {
      url: '/api/v2/instance/event/statistics',
      method: 'get',
    },
    getEventType: {
      url: '/api/v2/instance/event/typeList',
      method: 'get',
    },
    getEventEcharts: {
      url: '/api/v2/instance/event/coordinate',
      method: 'get',
    },
    getEventList: {
      url: '/api/v2/instance/event/list',
      method: 'get',
    },
    getEventProductRank: (params: any) => {
      const { product_or_service } = params;
      return {
        url: `api/v2/instance/event/${product_or_service}/rank`,
        method: 'get',
      };
    },
  },
  config: {
    getConfigAlertGroups: {
      url: '/api/v2/product/configAlterGroups',
      method: 'get',
    },
    getConfigAlertaction: {
      url: '/api/v2/product/configAlteration',
      method: 'get',
    },
  },
  product: {
    Verifylink: {
      url: '/api/v2/product/check_param',
      method: 'post',
    },
    // 上传安装包来自网络
    addNetworkPackage: {
      url: '/api/v2/product/uploadAsync',
      method: 'post',
    },
    // 获取正在上传产品包列表
    getProductpackageList: {
      url: '/api/v2/product/in_progress',
      method: 'get',
    },
    // 删除正在上传的产品包
    delProductpackageItem: {
      url: '/api/v2/product/cancel_upload',
      method: 'postFormData',
    },
    getParentProductList: {
      url: '/api/v2/product/parentProduct',
      method: 'get',
    },
    // 获取产品列表
    getClusterProductList: {
      url: '/api/v2/cluster/products',
      method: 'get',
    },
    getProductList: {
      url: '/api/v2/product',
      method: 'get',
    },
    getProductName: {
      url: '/api/v2/product/product_name_list',
      method: 'get',
    },
    // 获取更新目录
    getPatchPath: {
      url: '/api/v2/product/patchpath',
      method: 'get',
    },
    updatePatchPath: {
      url: '/api/v2/product/patchupload',
      method: 'postFormData',
    },
    updatePatchStatus: {
      url: '/api/v2/product/patchupdate',
      method: 'post',
    },
    getCurrentProduct: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/product/${product_name}/current`,
        method: 'get',
      };
    },
    getProductUpdateRecords: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/product/${product_name}/history`,
        method: 'get',
      };
    },
    getProductUpdate: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/product/${product_name}/updatehistory`,
        method: 'get',
      };
    },
    // 获取产品包升级目标版本列表
    getProductVersionList: (params: any) => {
      const { product_name, product_version } = params;
      return {
        url: `/api/v2/product/${product_name}/version/${product_version}/upgrade_candidate`,
        method: 'get',
      };
    },
    getProductDeployRecords: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/product/${product_name}`,
        method: 'get',
      };
    },
    delProduct: (params: any) => {
      const { product_name, product_version } = params;
      return {
        url: `/api/v2/product/${product_name}/version/${product_version}`,
        method: 'dele',
      };
    },
    modifyProductConfigAll: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/modifyAll`,
        method: 'post',
      };
    },
    //  配置信息关联主机
    modifyMultiAllHosts: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/modifyMultiAll`,
        method: 'post',
      };
    },
    resetServiceConfig: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/reset_schema_field`,
        method: 'postFormData',
      };
    },
    resetMultiServiceConfig: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/reset_multi_schema_field`,
        method: 'postFormData',
      };
    },
    getDeployShot: (params: any) => {
      const { uuid, limit, start, status } = params;
      return {
        url: `/api/v2/instance/${uuid}/list?limit=${limit}&start=${start}&status=${status}`,
        method: 'get',
      };
    },
    getPatchUpdateShot: (params: any) => {
      const { uuid, limit, start, status } = params;
      return {
        url: `/api/v2/instance/${uuid}/listupdate?limit=${limit}&start=${start}&status=${status}`,
        method: 'get',
      };
    },

    /**
     * 全部产品包-查看-配置-选择一个文件查看内容
     */
    getServiceTree: (params: ServiceTree) => {
      const { productName, productVersion, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/service/${serviceName}/serviceTree`,
        method: 'get',
      };
    },

    /**
     * 全部产品包-查看-配置-选择一个服务后返回文件列表
     */
    getServiceFile: (params: ServiceFile) => {
      const { productName, productVersion, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/service/${serviceName}/serviceFile`,
        method: 'get',
      };
    },
    /**
     * 服务配置下发
     */
    distributeServiceConfig: (params: DistributeServiceConfig) => {
      const { productId, serviceName } = params;
      return {
        url: `/api/v2/service/${productId}/${serviceName}/config_update`,
        method: 'post',
      };
    },
  },
  service: {
    getConfDiff: (params: any) => {
      const { product_name, product_version, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/version/${product_version}/service/${service_name}/serviceConfigDiff`,
        method: 'get',
      };
    },
    // 切换自动执行健康检查
    setAutoexecSwitch: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/instance/product/${product_name}/service/${service_name}/healthCheck/setAutoexecSwitch`,
        method: 'post',
      };
    },
    // 手动执行健康检查
    manualExecution: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/instance/product/${product_name}/service/${service_name}/healthCheck/manualExecution`,
        method: 'post',
      };
    },
    getHealthCheck: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/instance/product/${product_name}/service/${service_name}/healthCheck`,
        method: 'get',
      };
    },
    getHostAlert: {
      url: `/api/v2/cluster/hosts/alert`,
      method: 'get',
    },
    // 获取异常服务
    getRedService: (params: any) => {
      return {
        url: '/api/v2/product/anomalyService',
        method: 'get',
      };
    },
    getRestartService: (params: any) => {
      return {
        url: '/api/v2/cluster/restartServices',
        method: 'get',
      };
    },
    // 获取冒烟测试状态
    getAutoTestStatus: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/product/${product_name}/autoTest/history`,
        method: 'get',
      };
    },
    // 冒烟测试
    startAutoTestTest: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/product/${product_name}/autoTest/start`,
        method: 'post',
      };
    },
    getServiceStatus: {
      url: '/api/v2/product/status',
      method: 'get',
    },
    getLogFileList: (params: any) => {
      const { instanceId, type } = params;
      return {
        url: `/api/v2/instance/${instanceId}/logfiles?type=${type}`,
        method: 'get',
      };
    },
    getLogFile: (params: any) => {
      const { instanceId } = params;
      return {
        url: `/api/v2/instance/${instanceId}/logmore`,
        method: 'get',
      };
    },
    getLogFileDownLoad: (params: any) => {
      const { instanceId, logFile } = params;
      return {
        url: `/api/v2/instance/${instanceId}/logdown?logfile=${logFile}`,
        method: 'get',
      };
    },
    getParamDropList: (params: any) => {
      const { productName, productVersion, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/service/${serviceName}/serviceConfigFiles`,
        method: 'get',
      };
    },
    getParamContent: (params: any) => {
      const { productName, productVersion, serviceName, servicePath } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/service/${serviceName}/serviceConfigFile?file=${servicePath}`,
        method: 'get',
      };
    },
    setParamUpdate: (params: any) => {
      const { productName, productVersion, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/service/${serviceName}/configUpdate`,
        method: 'post',
      };
    },
    getServiceDashInfo: {
      url: '/api/search',
      method: 'get',
    },
    getServiceDashPlaylists: {
      url: '/api/playlists',
      method: 'get',
    },
    getServiceDashUrl: (params: any) => {
      const { id } = params;
      return {
        url: `/api/playlists/play/${id}`,
        method: 'get',
      };
    },
    getServiceGroup: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/instance/product/${product_name}/group_list`,
        method: 'get',
      };
    },
    getAllServiceGroup: (params: any) => {
      const { product_name, product_version } = params;
      return {
        url: `/api/v2/product/${product_name}/version/${product_version}/group_list`,
        method: 'get',
      };
    },
    setServiceGroupStart: (params: any) => {
      const { pid, group_name } = params;
      return {
        url: `/api/v2/group/${pid}/${group_name}/start`,
        method: 'post',
      };
    },
    setServiceGroupStop: (params: any) => {
      const { pid, group_name } = params;
      return {
        url: `/api/v2/group/${pid}/${group_name}/stop`,
        method: 'post',
      };
    },
    getServiceList: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/instance/product/${product_name}/service_list`,
        method: 'get',
      };
    },
    // 服务主机
    getServiceHostsList: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/instance/product/${product_name}/service/${service_name}`,
        method: 'get',
      };
    },
    // 服务 - 参数 - 运行配置
    getServiceConfig: (params: any) => {
      const { product_name, service_name, product_version } = params;
      return {
        url: `/api/v2/product/${product_name}/version/${product_version}/service/${service_name}/serviceConfig`,
        method: 'get',
      };
    },
    getServiceHostConfig: (params: any) => {
      const { instance_id } = params;
      return {
        url: `/api/v2/instance/${instance_id}/config`,
        method: 'get',
      };
    },
    startService: (params: any) => {
      const { pid, service_name } = params;
      return {
        url: `/api/v2/service/${pid}/${service_name}/start`,
        method: 'post',
      };
    },
    stopService: (params: any) => {
      const { pid, service_name } = params;
      return {
        url: `/api/v2/service/${pid}/${service_name}/stop`,
        method: 'post',
      };
    },
    startServiceInstance: (params: any) => {
      const { agent_id } = params;
      return {
        url: `/api/v2/instance/${agent_id}/start`,
        method: 'post',
      };
    },
    stopServiceInstance: (params: any) => {
      const { agent_id } = params;
      return {
        url: `/api/v2/instance/${agent_id}/stop`,
        method: 'post',
      };
    },
    setServiceRollRestart: (params: any) => {
      const { pid, service_name } = params;
      return {
        url: `/api/v2/service/${pid}/${service_name}/rolling_restart`,
        method: 'post',
      };
    },
    getServiceHARole: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/haRole/${service_name}`,
        method: 'get',
      };
    },
    startAllServiceByProduct: (params: any) => {
      const { pid } = params;
      return {
        url: `/api/v2/product/${pid}/start`,
        method: 'get',
      };
    },
    stopAllServiceByProduct: (params: any) => {
      const { pid } = params;
      return {
        url: `/api/v2/product/${pid}/stop`,
        method: 'get',
      };
    },
    // k8s 扩缩容
    instanceReplica: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/instance/${product_name}/${service_name}/replica`,
        method: 'post',
      };
    },
    operateExtension: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/extention_operation`,
        method: 'postFormData',
      };
    },
    operateSwitch: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/operate_switch`,
        method: 'postFormData',
      };
    },
    // 判断开关是否正在操作中
    getSwitchRecord: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/switch/record`,
        method: 'get',
      };
    },
    // 获取开关当前开启进度
    getSwitchDetail: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/product/${product_name}/service/${service_name}/switch/detail`,
        method: 'get',
      };
    },
  },
  instance: {
    getInstanceLog: (params: any) => {
      const { id } = params;
      return {
        url: `/api/v2/instance/${id}/log`,
        method: 'get',
      };
    },
  },
  host: {
    deleteHost: {
      url: 'api/v2/agent/hostdelete',
      method: 'post',
    },
    getFilterData: {
      url: '/api/host/condition',
      method: 'get',
    },
    getHostList: {
      url: '/api/v2/agent/hosts',
      method: 'get',
    },
    getTableData: {
      url: '/api/host/queryHost',
      method: 'get',
    },
    showInstallInfo: {
      url: '/api/v2/agent/install/checkinstallbysid',
      method: 'get',
    },
    getHostInstance: {
      url: 'HTTP/api/v2/agent/hosts',
      method: 'get',
    },
    getHostServicesList: {
      url: '/api/v2/agent/hostService',
      method: 'get',
    },
    confirmMoveHost: {
      url: '/api/v2/agent/hostmove', // 确定移动主机
      method: 'post',
    },
    updateGroupName: {
      url: '/api/v2/agent/hostgroup_rename', // 修改主机名称
      method: 'post',
    },
  },
  addHost: {
    pwdConnectUrl: {
      url: '/api/v2/agent/install/pwdconnect',
      method: 'post',
    },
    pwdInstallUrl: {
      url: '/api/v2/agent/install/pwdinstall',
      method: 'post',
    },
    pkConnectUrl: {
      url: '/api/v2/agent/install/pkconnect',
      method: 'post',
    },
    pkInstallUrl: {
      url: '/api/v2/agent/install/pkinstall',
      method: 'post',
    },
    checkInstallUrl: {
      url: '/api/v2/agent/install/checkinstall',
      method: 'post',
    },
    checkInstallAllUrl: {
      url: '/api/v2/agent/install/checkinstallall?sort-by=id&sort-dir=0&limit=0',
      method: 'get',
    },
  },
  dashboard: {
    createDashboard: {
      url: '/api/dashboards/db',
      method: 'post',
    },
    createDashFolder: {
      url: '/api/folders',
      method: 'post',
    },
    delDashByUid: {
      url: '/api/dashboards/uid/',
      method: 'dele',
    },
    delFolderByUid: {
      url: '/api/folders/',
      method: 'dele',
    },
    importDashboard: {
      url: '/api/dashboards/import',
      method: 'post',
    },
    getAlertNotifications: {
      url: '/api/alert-notifications',
      method: 'get',
    },
    exportDash: {
      url: '/api/v2/dashboard/export',
      method: 'get',
      contentType: 'application/zip',
    },
    delAlertNotification: {
      url: '/api/alert-notifications/',
      method: 'dele',
    },
    addGrafanaNotification: {
      url: '/api/alert-notifications',
      method: 'post',
    },
  },
  alertRule: {
    getAlertsByDashId: {
      url: '/api/alerts',
      method: 'get',
    },
    getAlertRuleList: {
      url: '/api/v2/dashboard/alerts',
      method: 'get',
    },
    getAlertsHistory: {
      url: '/api/annotations',
      method: 'get',
    },
    dtstackAlertChannelTest: {
      url: '/gate/alert/send',
      method: 'post',
    },
    dtstackAlertChannelSave: {
      url: '/gate/alert/edit',
      method: 'post',
    },
    dtstackAlertChannelList: {
      url: '/gate/alert/list',
      method: 'post',
    },
    dtstackAlertChannelDel: {
      url: '/gate/alert/delete',
      method: 'get',
    },
    grafanaAlertChannelTest: {
      url: '/api/alert-notifications/test',
      method: 'post',
    },
    grafanaAlertChannelSave: {
      url: '/api/alert-notifications',
      method: 'post',
    },
    switchGrafanaAlertPause: {
      url: `/api/v2/dashboard/alerts/pause`,
      method: 'post',
    },
  },
  alertChannel: {
    grafanaAlertChannelUpdate: (params: any) => {
      const { id } = params;
      return {
        url: `/api/alert-notifications/${id}`,
        method: 'put',
      };
    },
    getAlertNotifications: {
      url: '/api/alert-notifications',
      method: 'get',
    },
    delAlertNotification: {
      url: '/api/alert-notifications/',
      method: 'dele',
    },
    grafanaAlertChannelTest: {
      url: '/api/alert-notifications/test',
      method: 'post',
    },
    grafanaAlertChannelSave: {
      url: '/api/alert-notifications',
      method: 'post',
    },
    dtstackAlertChannelSave: {
      url: '/gate/alert/edit',
      method: 'post',
    },
    dtstackAlertChannelList: {
      url: '/gate/alert/list',
      method: 'post',
    },
    getDtstackAlertDetail: {
      url: '/gate/alert/get',
      method: 'get',
    },
    getGrafanaAlertDetail: (params: any) => {
      const { alertId } = params;
      return {
        url: `/api/alert-notifications/${alertId}`,
        method: 'get',
      };
    },
    dtstackAlertChannelDel: {
      url: '/gate/alert/delete',
      method: 'get',
    },
  },
  alert: {
    getAlertsByDashId: {
      url: '/api/alerts',
      method: 'get',
    },
    getAlertsHistory: (params: any) => {
      const { product_name, service_name } = params;
      return {
        url: `/api/v2/instance/product/${product_name}/service/${service_name}/alert`,
        method: 'get',
      };
    },
    dtstackAlertChannelTest: {
      url: '/gate/alert/send',
      method: 'post',
    },
    dtstackAlertChannelSave: {
      url: '/gate/alert/edit',
      method: 'post',
    },
    dtstackAlertChannelDel: {
      url: '/gate/alert/delete',
      method: 'get',
    },
    grafanaAlertChannelTest: {
      url: '/api/alert-notifications/test',
      method: 'post',
    },
    grafanaAlertChannelSave: {
      url: '/api/alert-notifications',
      method: 'post',
    },
  },
  installGuide: {
    getHostRoleMap: {
      url: '/api/v2/cluster/hosts/role_info',
      method: 'get',
    },
    refreshServicesInfoForAutoDeploy: {
      url: '/api/v2/cluster/hosts/auto_svcgroup',
      method: 'post',
    },
    checkMySqlAddr: (params: any) => {
      const { product_name } = params;
      return {
        url: `/api/v2/product/${product_name}/checkMysqlAddr`,
        method: 'post',
      };
    },
    getGlobalAutoConfig: (params: any) => {
      const { productName, productVersion, carbon_thriftserver } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/serviceGraphy?unchecked_services=${carbon_thriftserver}`,
        method: 'post',
      };
    },
    getAutoConfig: (params: any) => {
      const { productName, productVersion, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/service/${serviceName}/serviceGraphy`,
        method: 'post',
      };
    },
    getFilters: {
      url: '/api/v2/product/productName',
      method: 'get',
    },
    getProductPackageList: {
      url: '/api/v2/product/productList?sort-by=create_time&sort-dir=desc',
      method: 'get',
    },
    getProductStepOneList: {
      url: '/api/v2/product_line/product_list',
      method: 'get',
    },
    getProductPackageServices: (params: any) => {
      const { productName, productVersion } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/service`,
        method: 'get',
      };
    },
    deleteProductPackage: (params: any) => {
      const { productName, productVersion } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}`,
        method: 'dele',
      };
    },
    serviceUpdate: (params: any) => {
      const { ProductName } = params;
      return {
        url: `/api/v2/product/${ProductName}/serviceUpdate`,
        method: 'post',
      };
    },
    getInstallCmd: {
      url: '/api/v2/agent/install/installCmd',
      method: 'get',
    },

    // 获取集群列表
    getInstallClusterList: {
      url: '/api/v2/cluster/list',
      method: 'get',
    },
    getProductServicesInfo: (params: any) => {
      const { productName, productVersion, unSelectService } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/serviceGroup?unchecked_services=${unSelectService}`,
        method: 'get',
      };
    },
    updateServiceHostList: (params: any) => {
      const { productName, serviceName, clusterId } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/hosts?clusterId=${clusterId}`,
        method: 'get',
      };
    },
    setParamConfigFieldValue: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/modify_schema_field`,
        method: 'postFormData',
      };
    },
    setParamConfigFieldValueTotal: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/modify_schema_field_batch`,
        method: 'post',
      };
    },
    resetParamConfigFieldValue: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/reset_schema_field`,
        method: 'postFormData',
      };
    },
    resetMultiConfigFieldValue: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/reset_multi_schema_field`,
        method: 'postFormData',
      };
    },
    deploy: (params: any) => {
      const { productName, version } = params;
      return {
        url: `/api/v2/product/${productName}/version/${version}/deploy`,
        method: 'post',
      };
    },
    checkDeployStatus: (params: any) => {
      const { productName, version } = params;
      return {
        url: `/api/v2/product/${productName}/version/${version}`,
        method: 'get',
      };
    },
    getDeployTaskList: (params: any) => {
      const { uuid, start, status } = params;
      return {
        url: `/api/v2/instance/${uuid}/list?start=${start}&limit=20&status=${status}`,
        method: 'get',
      };
    },
    stopDeploy: (params: any) => {
      const { productName, version } = params;
      return {
        url: `/api/v2/product/${productName}/version/${version}/cancel`,
        method: 'postFormData',
      };
    },
    stopAutoDeploy: {
      url: '/api/v2/cluster/hosts/auto_deploy_cancel',
      method: 'post',
    },
    modifyProductConfigAll: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/modifyAll`,
        method: 'post',
      };
    },
    // 运行部署配置信息
    modifyMultiSingleField: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/modifyMultiSingleField`,
        method: 'postFormData',
      };
    },
    setIp: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/set_ip`,
        method: 'postFormData',
      };
    },
    getUncheckedService: (params: any) => {
      return {
        url: `/api/v2/product/${params.pid}/unchecked_services`,
        method: 'get',
      };
    },

    // 创建命名空间
    createNamespace: (params: any) => {
      return {
        url: `/api/v2/cluster/kubernetes/${params.clusterId}/namespace/create`,
        method: 'post',
      };
    },

    // 获取命名空间列表
    getNamespaceList: (params: any) => {
      return {
        url: `/api/v2/cluster/kubernetes/${params.clusterId}/namespace/list`,
        method: 'get',
      };
    },

    // 部署向导2-产品依赖检测
    getBaseClusterList: (params: any) => {
      const { clusterId, namespace, pid } = params;
      return {
        url: `/api/v2/cluster/kubernetes/${clusterId}/namespace/${namespace}/product/${pid}/depends`,
        method: 'get',
      };
    },

    // k8s部署 - 第一步 - 检测默认镜像仓库是否存在
    checkDefaultImageStore: (params: any) => {
      const { clusterId } = params;
      return {
        url: `/api/v2/cluster/kubernetes/imageStore/${clusterId}/checkDefault`,
        method: 'get',
      };
    },
    // 自动编排
    autoOrchestration: {
      url: '/api/v2/cluster/hosts/auto_orchestration',
      method: 'post',
    },
    // 自动部署
    autoDeploy: {
      url: '/api/v2/cluster/hosts/auto_deploy',
      method: 'post',
    },

    // 获取老版本服务主机编排与配置信息
    getOldHostInfo: (params: any) => {
      const { productName } = params;
      return {
        url: `/api/v2/product/${productName}/currentInfo`,
        method: 'post',
      };
    },

    getServiceGroupFile: (params: any) => {
      const { productName, productVersion } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/serviceGroupFile`,
        method: 'get',
      };
    },

    // 保存升级历史记录
    saveUpgrade: (params: any) => {
      const { productName } = params;
      return {
        url: `/api/v2/product/${productName}/saveUpgrade`,
        method: 'post',
      };
    },

    deployCondition: {
      url: `/api/v2/product/deployCondition`,
      method: 'post',
    },
  },
  deployAPI: {
    getProductConfig: (params: any) => {
      const { productName, productVersion } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}`,
        method: 'get',
      };
    },
    getIp: (params: any) => {
      const { productName, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/service/${serviceName}/get_ip`,
        method: 'get',
      };
    },
    deploy: (params: any) => {
      const { productName, productVersion } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/deploy`,
        method: 'post',
      };
    },
    getDeployList: (params: any) => {
      const { deployUuid } = params;
      return {
        url: `/api/v2/instance/${deployUuid}/list`,
        method: 'get',
      };
    },
    cancelDeploy: (params: any) => {
      const { productName, productVersion } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/cancel`,
        method: 'post',
      };
    },
    /**
     * 搜索部署日志
     */
    searchDeployLog: (params: SearchDeployLogs) => {
      const { productName, productVersion, deployId, serviceName } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/deployLogs?deploy_uuid=${deployId}&service=${serviceName}`,
        method: 'get',
      };
    },
    // 获取自动部署安装记录
    getOrchestrationHistory: {
      url: '/api/v2/cluster/hosts/orchestration_history',
      method: 'get',
    },
    // 处理备份
    handleBackUp: (params: SearchDeployLogs) => {
      const { productName } = params;
      return {
        url: `/api/v2/product/${productName}/backupDb`,
        method: 'post',
      };
    },

    // 处理回滚
    handleRollBack: (params: SearchDeployLogs) => {
      const { productName } = params;
      return {
        url: `/api/v2/product/${productName}/rollback`,
        method: 'post',
      };
    },

    // 获取回滚版本
    getRollBackList: (params) => {
      const { productName } = params;
      return {
        url: `/api/v2/product/${productName}/rollbackVersions`,
        method: 'post',
      };
    },

    // 获取备份记录时间
    getBackupTimes: (params) => {
      const { productName } = params;
      return {
        url: `/api/v2/product/${productName}/backupTimes`,
        method: 'post',
      };
    },
  },
  // 获取集群巡检
  // 集群巡检统计设置
  clusterInspectionApi: {
    getClusterInspectionStatisSet: {
      url: `/api/v2/platform/inspect/statisticsConfig/update`,
      method: 'post',
    },
    getSettingData: {
      url: `/api/v2/platform/inspect/statisticsConfig`,
      method: 'get',
    },
    getClusterInspectionBaseInfo: {
      url: `/api/v2/platform/inspect/baseInfo/status`,
      method: 'get',
    },
    getClusterInspectionTable: {
      url: `/api/v2/platform/inspect/form/data`,
      method: 'get',
    },
    getClusterrApplicationServerInfo: {
      url: `/api/v2/platform/inspect/graph/config`,
      method: 'get',
    },
    getClusterApplicationServerTableInfo: {
      url: `/api/v2/inspect/graph/data`,
      method: 'get',
    },
    getClusterBigDataInfo: {
      url: `/api/v2/platform/inspect/baseInfo/name_node`,
      method: 'get',
    },
  },
  unDeployAPI: {
    unDeploy: (params: any) => {
      const { productName, productVersion } = params;
      return {
        url: `/api/v2/product/${productName}/version/${productVersion}/undeploy`,
        method: 'post',
      };
    },
    getUnDeployList: (params: any) => {
      const { deployUuid, start, status, limit } = params;
      return {
        url: `/api/v2/instance/${deployUuid}/list?start=${start}&limit=${limit}&status=${status}`,
        method: 'get',
      };
    },
    forceStop: (params: any) => {
      return {
        url: `/api/v2/instance_record/${params}/force_stop`,
        method: 'post',
      };
    },
    forceUninstall: (params: any) => {
      return {
        url: `/api/v2/instance_record/${params}/force_uninstall`,
        method: 'post',
      };
    },
    getUnDeployLog: (params: any) => {
      const { id } = params;
      return {
        url: `/api/v2/instance/deploy/${id}/log`,
        method: 'get',
      };
    },
    // 停止卸载
    stopUndeploy: (params: any) => {
      const { clusterId, namespace, pid } = params;
      return {
        url: `/api/v2/cluster/kubernetes/${clusterId}/namespace/${namespace}/product/${pid}/stop`,
        method: 'post',
      };
    },
  },
  createHost: {
    testConnection_pwd: {
      url: '/api/v2/agent/install/pwdconnect',
      method: 'postJsonData',
    },
    testConnection_pk: {
      url: '/api/v2/agent/install/pkconnect',
      method: 'post',
    },
    pkInstallUrl: {
      url: '/api/v2/agent/install/pkinstall',
      method: 'post',
    },
    pwdInstallUrl: {
      url: '/api/v2/agent/install/pwdinstall',
      method: 'post',
    },
  },

  scriptEcho: {
    getEchoCount: (params: any) => {
      return {
        url: '/api/v2/cluster/currExecCount',
        method: 'get',
      };
    },
    getEchoOrder: (params: any) => {
      return {
        url: '/api/v2/cluster/orderList',
        method: 'get',
      };
    },
    getEchoSearchList: (params: any) => {
      return {
        url: '/api/v2/cluster/listObjectValue',
        method: 'get',
      };
    },
    getEchoOrderDetail: (params: any) => {
      return {
        url: '/api/v2/cluster/orderDetail',
        method: 'get',
      };
    },
    showShellLog: {
      url: '/api/v2/cluster/showShellLog',
      method: 'get',
    },
    showShellContent: {
      url: '/api/v2/cluster/previewShellContent',
      method: 'get',
    },
    downShellLog: {
      url: '/api/v2/cluster/orderLogDownload',
      method: 'get',
    },
    downShellContent: {
      url: '/api/v2/cluster/downLoadShellContent',
      method: 'get',
    },
  },

  userCenter: {
    login: {
      url: '/api/v2/user/login',
      method: 'postFormData',
    },
    getMembers: {
      url: '/api/v2/user/list',
      method: 'get',
    },
    removeMember: {
      url: '/api/v2/user/remove',
      method: 'postFormData',
    },
    resetPasswordSelf: {
      url: '/api/v2/user/modifyPwd',
      method: 'postFormData',
    },
    resetPassword: {
      url: '/api/v2/user/resetPwdByAdmin',
      method: 'postFormData',
    },
    regist: {
      url: '/api/v2/user/register',
      method: 'postFormData',
    },
    // 个人信息修改
    motifyUserInfo: {
      url: '/api/v2/user/modifyInfo',
      method: 'postFormData',
    },
    // 成员管理 - 编辑信息
    modifyInfoByAdmin: {
      url: '/api/v2/user/modifyInfoByAdmin',
      method: 'postFormData',
    },
    validEmail: {
      url: '',
      method: '',
    },
    logout: {
      url: '/api/v2/user/logout',
      method: 'postFormData',
    },
    getValidCode: {
      url: '/api/v2/user/getCaptcha',
      method: 'get',
    },
    checkValidCode: {
      url: '/api/v2/user/processCaptcha',
      method: 'postFormData',
    },
    enableUser: {
      url: '/api/v2/user/enable',
      method: 'postFormData',
    },
    disableUser: {
      url: '/api/v2/user/disable',
      method: 'postFormData',
    },
    getLoginedUserInfo: {
      url: '/api/v2/user/personal',
      method: 'post',
    },
    // 获取公钥
    getPublicKey: {
      url: '/api/v2/user/getPublicKey',
      method: 'get',
    },
    // 获取角色列表
    getRoleList: {
      url: '/api/v2/role/list',
      method: 'get',
    },
    // 获取权限树
    getAuthorityTree: (params: any) => {
      const { roleId } = params;
      return {
        url: `/api/v2/role/${roleId}/permissions`,
        method: 'get',
      };
    },
    // 获取权限code
    getRoleCodes: {
      url: '/api/v2/role/codes',
      method: 'get',
    },
    generate: {
      url: '/api/v2/common/deployInfo/generate',
      method: 'post',
    },
    downloadInfo: {
      url: '/api/v2/common/deployInfo/download',
      method: 'get',
    },
  },
  clusterManager: {
    // 获取集群列表
    getClusterLists: {
      url: '/api/v2/cluster/list',
      method: 'get',
    },
    /* -- 主机集群 -- */
    // 创建主机集群
    createHostCluster: {
      url: '/api/v2/cluster/hosts/create',
      method: 'post',
    },
    // 编辑主机集群
    updateHostCluster: {
      url: '/api/v2/cluster/hosts/update',
      method: 'post',
    },
    // 删除主机集群
    deleteHostCluster: (params: any) => {
      const { cluster_id } = params;
      return {
        url: `/api/v2/cluster/hosts/${cluster_id}/delete`,
        method: 'post',
      };
    },
    // 获取主机集群详情
    getHostClusterInfo: (params: any) => {
      const { cluster_id } = params;
      return {
        url: `/api/v2/cluster/hosts/${cluster_id}/info`,
        method: 'get',
      };
    },
    /* -- kubernetes集群 -- */
    // 创建kubernetes集群
    createKubernetesCluster: {
      url: '/api/v2/cluster/kubernetes/create',
      method: 'post',
    },
    // 更新kubernetes集群
    updateKubernetesCluster: {
      url: '/api/v2/cluster/kubernetes/update',
      method: 'post',
    },
    // 删除kubernetes集群
    deleteKubernetesCluster: (params: any) => {
      const { cluster_id } = params;
      return {
        url: `/api/v2/cluster/kubernetes/${cluster_id}/delete`,
        method: 'post',
      };
    },
    // 获取kubernetes集群详情
    getKubernetesClusterInfo: (params: any) => {
      const { cluster_id } = params;
      return {
        url: `/api/v2/cluster/kubernetes/${cluster_id}/info`,
        method: 'get',
      };
    },
    // 获取自建kubernetes版本，网络组件信息
    getKubernetesAvaliable: {
      url: '/api/v2/cluster/kubernetes/available',
      method: 'get',
    },
    // 获取导入kubernetes集群接入信息
    getKubernetesInstallCmd: {
      url: '/api/v2/cluster/kubernetes/installCmd',
      method: 'get',
    },
    // 获取yaml模板
    getKubernetesRketemplate: {
      url: '/api/v2/cluster/kubernetes/rketemplate',
      method: 'get',
    },
  },
  clusterHost: {
    // 主机删除，沿用之前接口
    deleteHost: {
      url: '/api/v2/agent/hostdelete',
      method: 'post',
    },
    // 获取集群下角色清单
    getRoleList: {
      url: 'api/v2/cluster/hosts/role_list',
      method: 'get',
    },
    // 删除集群下角色
    deleteRole: {
      url: '/api/v2/cluster/hosts/role_delete',
      method: 'post',
    },
    // 添加角色
    addRole: {
      url: '/api/v2/cluster/hosts/role_add',
      method: 'post',
    },
    modifyRole: {
      url: '/api/v2/cluster/hosts/role_rename',
      method: 'post',
    },
    // 主机角色编辑
    bindHostRoles: {
      url: '/api/v2/cluster/hosts/role',
      method: 'post',
    },
  },
  imageStore: {
    // 获取集群下镜像仓库列表
    getImageStoreList: (params: any) => {
      return {
        url: `/api/v2/cluster/kubernetes/imageStore/${params.cluster_id}/clusterInfo`,
        method: 'get',
      };
    },
    // 获取镜像仓库信息
    getImageStoreInfo: (params: any) => {
      return {
        url: `/api/v2/cluster/kubernetes/imageStore/${params.store_id}/info`,
        method: 'get',
      };
    },
    // 创建镜像仓库
    createImageStore: {
      url: '/api/v2/cluster/kubernetes/imageStore/create',
      method: 'post',
    },
    // 更新镜像仓库
    updateImageStore: {
      url: '/api/v2/cluster/kubernetes/imageStore/update',
      method: 'post',
    },
    // 删除镜像仓库
    deleteImageStore: {
      url: '/api/v2/cluster/kubernetes/imageStore/delete',
      method: 'post',
    },
    // 设置默认仓库
    setDefaultStore: {
      url: '/api/v2/cluster/kubernetes/imageStore/setDefault',
      method: 'post',
    },
  },
  clusterIndex: {
    // 获取主机集群总览信息
    getHostsClusterOverview: (params?: any) => {
      return {
        url: `/api/v2/cluster/hosts/${params.cluster_id}/overview`,
        method: 'get',
      };
    },
    // 获取Kubernetes集群总览信息
    getKubernetesClusterOverview: (params?: any) => {
      return {
        url: `/api/v2/cluster/kubernetes/${params.cluster_id}/overview`,
        method: 'get',
      };
    },
    // 获取CPU,MEMERRY等容量趋势图
    getHostsClusterPerformance: (id: number, params?: any) => {
      return {
        url: `/api/v2/cluster/hosts/${id}/performance`,
        method: 'get',
      };
    },
    // 获取CPU,MEMERRY等容量趋势图 - k8s
    getKubernetesClusterPerformance: (id: number, params?: any) => {
      return {
        url: `/api/v2/cluster/kubernetes/${id}/performance`,
        method: 'get',
      };
    },
  },
  clusterNamespace: {
    // 获取命名空间列表
    getNamespaceList: (params: any) => {
      const { namespace } = params;
      const url =
        '/api/v2/cluster/manage/namespaces' +
        (namespace ? '/' + namespace : '');
      return {
        url: url,
        method: 'get',
      };
    },
    // 测试连通性
    pingConnect: {
      url: '/api/v2/cluster/manage/namespace/ping',
      method: 'post',
    },
    // 获取yaml文件
    getYamlFile: {
      url: '/api/v2/cluster/manage/namespace/agent/generate',
      method: 'post',
    },
    // 保存namespace
    saveNamespace: {
      url: '/api/v2/cluster/manage/namespace/save',
      method: 'post',
    },
    // 删除前进行校验，是否可以删除
    confirmDelete: (params) => {
      const { namespace } = params;
      return {
        url: `/api/v2/cluster/manage/${namespace}/delete/confirm`,
        method: 'get',
      };
    },
    // 删除命名空间
    deleteNamespace: (params) => {
      const { namespace } = params;
      return {
        url: `/api/v2/cluster/manage/${namespace}/delete`,
        method: 'post',
      };
    },
    // 获取命名空间信息
    getNamespaceInfo: (params) => {
      const { namespace } = params;
      return {
        url: `/api/v2/cluster/manage/${namespace}/get`,
        method: 'get',
      };
    },
    // 获取服务列表
    getServiceLists: (params) => {
      const { namespace, parentProductName, productName, serviceName } = params;
      const url =
        '/api/v2/product/manage' +
        (namespace ? '/' + namespace : '') +
        (parentProductName ? '/' + parentProductName : '') +
        (productName ? '/' + productName : '') +
        (serviceName ? '/' + serviceName : '');
      return {
        url,
        method: 'get',
      };
    },
    // 获取事件查看
    getEventLists: (params) => {
      const { namespace } = params;
      return {
        url: `/api/v2/cluster/manage/${namespace}/events`,
        method: 'get',
      };
    },
  },
  cluster: {
    getClusterList: {
      url: '/api/v2/cluster/list',
      method: 'get',
    },
  },
  securityAudit: {
    // 获取安全审计列表
    getSafetyAudit: {
      url: '/api/v2/common/safetyAudit/list',
      method: 'get',
    },
    // 获取审计操作模块
    getAuditModule: {
      url: '/api/v2/common/safetyAudit/module',
      method: 'get',
    },
    // 获取审计操作模块
    getAuditOperation: {
      url: '/api/v2/common/safetyAudit/operation',
      method: 'get',
    },
  },
  common: {
    // 获取集群下主机分组
    getClusterhostGroupLists: {
      url: '/api/v2/cluster/hostgroups',
      method: 'get',
    },
    // 获取主机集群下主机列表
    getHostClusterHostList: {
      url: '/api/v2/cluster/hosts/hosts',
      method: 'get',
    },
    // 获取kubernetes集群下主机列表
    getKubernetesClusterHostList: {
      url: '/api/v2/cluster/kubernetes/hosts',
      method: 'get',
    },
  },
  backup: {
    //查询组建备份路径
    queryBuildBackupPath: {
      url: '/api/v2/product/backup/getconfig',
      method: 'get',
    },
    //设置组建备份路径
    SetUpBackupPath: {
      url: '/api/v2/product/backup/setconfig',
      method: 'post',
    },
  },
  security: {
    getSecurity: {
      url: '/api/v2/user/sys_config/platformSecurity',
      method: 'get',
    },

    setSecurity: {
      url: '/api/v2/user/sys_config/platformSecurity',
      method: 'post',
    },
  },

  scriptManager: {
    // 获取脚本列表
    getList: {
      url: '/api/v2/task',
      method: 'get',
    },
    // 获取脚本内容
    getTaskContent: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/task/${id}/content`,
        method: 'get',
      };
    },
    // 修改定时状态
    getTaskStatus: {
      url: '/api/v2/task/status',
      method: 'post',
    },
    // 修改定时设置
    getSettingStatus: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/task/${id}`,
        method: 'post',
      };
    },
    // cron表达式校验
    checkCronParse: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/task/${id}/cronParse`,
        method: 'post',
      };
    },
    // 手动执行脚本
    taskRun: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/task/${id}/run`,
        method: 'post',
      };
    },
    // 获取执行历史
    getTaskLog: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/task/${id}/log`,
        method: 'get',
      };
    },
    // 删除脚本
    deleteTask: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/task/${id}`,
        method: 'dele',
      };
    },
    editScript: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/task/${id}/edit`,
        method: 'post',
      };
    },
  },
  // 全局配置
  globalConfig: {
    // 查询全局配置
    getGlobalConfig: {
      url: '/api/v2/user/sys_config/globalConfig',
      method: 'get',
    },
    // 修改全局配置
    setGlobalConfig: {
      url: '/api/v2/user/sys_config/globalConfig',
      method: 'post',
    },
  },
  // 产品线相关接口
  productLine: {
    // 获取所有产品线信息接口
    getProductLine: {
      url: '/api/v2/product_line',
      method: 'get',
    },
    uploadProductLine: {
      url: '/api/v2/product_line/upload',
      method: 'post',
    },
    deleteProductLine: (params) => {
      const { id } = params;
      return {
        url: `/api/v2/product_line/${id}`,
        method: 'dele',
      };
    },
  },
};
