import { message } from 'antd';
import {
  deployActionTypes,
  unDeployActionTypes,
} from '@/constants/actionTypes';
import { deployService } from '@/services';
import utils from '@/utils/utils';

export const updateProductConfig = (params: any) => {
  return (dispatch: any) => {
    deployService
      .getProductConfig({
        productName: params.product_name,
        productVersion: params.product_version,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          dispatch({
            type: deployActionTypes.UPDATE_PRD_CONFIG,
            payload: data.data.list[0],
          });
        } else {
          message.error(data.msg);
        }
      });
  };
};

// 获取平滑升级下的组件名称
export const saveForcedUpgrade = (params: any) => {
  return {
    type: deployActionTypes.SAVE_FORCED_UPGRADE,
    payload: params,
  };
};

// 保存升级模式
export const getUpgradeType = (params: any) => {
  return {
    type: deployActionTypes.UPGRADE_TYPE,
    payload: params,
  };
};

// 判断是是第一次平滑升级
export const getIsFirstSmooth = (params: any) => {
  return {
    type: deployActionTypes.GET_FIRST_SMOOTH,
    payload: params,
  };
};

export const modifyServiceInput = (params: any) => {
  if (params.type === 0) {
    // Instance
    return {
      type: deployActionTypes.MODIFY_INSTANCE_CONFIG,
      payload: params,
    };
  } else {
    return {
      type: deployActionTypes.MODIFY_CONFIG_CONFIG,
      payload: params,
    };
  }
};

export const modifyInstanceConfig = (params: any) => {
  return {
    type: deployActionTypes.MODIFY_INSTANCE_CONFIG,
    payload: params,
  };
};
export const modifyConfigConfig = (params: any) => {
  return {
    type: deployActionTypes.MODIFY_CONFIG_CONFIG,
    payload: params,
  };
};

export const modifyServiceConfig = (params: any) => {
  return (dispatch: any) => {
    deployService
      .motifyServiceConfig(
        { productName: params.product_name, serviceName: params.service_name },
        {
          field_path: params.field_path,
          field: params.field,
        }
      )
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          if (data.data) {
            message.warning(data.data);
          }
          dispatch(updateProductConfig(params));
        } else {
          message.error(data.msg);
        }
      });
  };
};
export const resetServiceConfig = (params: any) => {
  const service_name = params.service_name;
  return (dispatch: any) => {
    deployService
      .resetSchemaField(
        {
          productName: params.product_name,
          serviceName: params.service_name,
        },
        {
          field_path: params.field_path,
          product_version: params.product_version,
        }
      )
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          deployService
            .getProductConfig({
              productName: params.product_name,
              productVersion: params.product_version,
            })
            .then((data: any) => {
              data = data.data;
              if (data.code === 0) {
                dispatch({
                  type: deployActionTypes.UPDATE_PRD_CONFIG,
                  payload: data.data.list[0],
                });
                deployService
                  .getIp({
                    productName: params.product_name,
                    serviceName: params.service_name,
                  })
                  .then((data: any) => {
                    data = data.data;
                    if (data.code === 0) {
                      dispatch({
                        type: deployActionTypes.GET_SERVICE_IPS,
                        payload: {
                          service_name: service_name,
                          ips: data.data.ip,
                        },
                      });
                    } else {
                      message.error(data.msg);
                    }
                  });
              } else {
                message.error(data.msg);
              }
            });
        } else {
          message.error(data.msg);
        }
      });
  };
};
// 获取服务的ip列表
export const getServiceIpList = (params: any) => {
  const service_name = params.service_name;
  return (dispatch: any) => {
    deployService
      .getIp({
        productName: params.product_name,
        serviceName: params.service_name,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          dispatch({
            type: deployActionTypes.GET_SERVICE_IPS,
            payload: {
              service_name: service_name,
              ips: data.data.ip || [],
            },
          });
        } else {
          message.error(data.msg);
        }
      });
  };
};
// ip字段input change事件
export const modifyIpConfig = (pindex: any, index: any, value: any) => {
  return {
    type: deployActionTypes.UPDATE_IP_CONFIG,
    payload: {
      pindex: pindex,
      index: index,
      ip: value,
    },
  };
};
// ip字段input blur事件
export const handleSetServiceIp = (params: any) => {
  return (dispatch: any) => {
    deployService
      .setIp(
        {
          serviceName: params.service_name,
          productName: params.product_name,
        },
        { ip: params.ip }
      )
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          // debugger;
          dispatch({
            type: deployActionTypes.GET_SERVICE_IPS,
            payload: {
              service_name: params.service_name,
              ips: params.ip ? params.ip.split(',') : [],
            },
          });
          dispatch(updateProductConfig(params));
        } else {
          message.error(data.msg);
        }
      });
  };
};
// 点击添加ip按钮
export const addIpToConfig = (pindex: any) => {
  return {
    type: deployActionTypes.ADD_IP_TO_CONFIG,
    payload: {
      pindex: pindex,
    },
  };
};
// update ip by service
export const updateIpsByService = (pindex: any, ips: any) => {
  return {
    type: deployActionTypes.UPDATE_IP_BY_SERVICE,
    payload: {
      service: pindex,
      ips: ips,
    },
  };
};
export const startProductDeploy = (params: any, callback: any) => {
  return (dispatch: any) => {
    deployService
      .deploy({
        productName: params.product_name,
        productVersion: params.product_version,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          dispatch({
            type: deployActionTypes.START_PRD_DEPLOY,
            payload: {
              deploy_uuid: data.data.deploy_uuid,
              deploy_status: true,
            },
          });
          callback();
        } else {
          message.error(data.msg);
        }
      });
  };
};

export const returnToStepOne = () => {
  return {
    type: deployActionTypes.RETURN_TO_CONFIG,
  };
};

export const getDeployList = (params: any, callback: any) => {
  return (dispatch: any) => {
    deployService
      .getProductConfig({
        productName: params.product_name,
        productVersion: params.product_version,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          dispatch({
            type: deployActionTypes.UPDATE_DEPLOY_STATUS,
            payload: data.data.list[0].status,
          });
          deployService
            .getDeployList({
              deployUuid: params.deploy_uuid,
            })
            .then((data: any) => {
              if (data.code === 0) {
                return dispatch({
                  type: deployActionTypes.UPDATE_DEPLOY_LIST,
                  payload: {
                    list: data.data.list,
                    complete: data.data.complete,
                  },
                });
              } else {
                message.error(data.msg);
              }
            });
          if (data.data.list[0].status !== 'deploying') {
            callback(data);
          }
        } else {
          message.error(data.msg);
        }
      });
  };
};

export const resetDeployStatus = () => {
  return {
    type: deployActionTypes.RESET_DEPLOY_STATUS,
  };
};

export const switchUseCloudByProduct = (services: any) => {
  return {
    type: deployActionTypes.SWITH_USE_CLOUD,
    payload: services,
  };
};

export const cancelDeploy = (params: any) => {
  return (dispatch: any) => {
    deployService
      .cancelDeploy({
        productName: params.product_name,
        productVersion: params.product_version,
      })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          // 停止部署成功
          message.success('已停止部署！');
        } else {
          message.error(res.msg);
        }
      });
  };
};

/** 卸载action
 * @param params: 请求参数
 * @param callback: 回调list接口
 */
export const startUnDeployService = (params: any) => {
  return (dispatch: any) => {
    deployService
      .unDeployService(
        {
          productName: params.product_name,
          productVersion: params.product_version,
        },
        {
          clusterId: params.clusterId,
          namespace: params.namespace,
        }
      )
      .then((res: any) => {
        const data = res.data;
        if (data.code === 0) {
          console.log('data.data.deploy_uuid', data.data.deploy_uuid);
          dispatch({
            type: unDeployActionTypes.START_UNDEPLOY,
            payload: {
              deploy_uuid: data.data.deploy_uuid,
              autoRefresh: true, // 开启自动刷新
            },
          });
        } else {
          message.error(data.msg);
        }
      });
    // callback(); // test
  };
};

export const getUndeploy = ({
  deploy_uuid,
  autoRefresh,
  complete,
}: {
  deploy_uuid: string;
  autoRefresh: boolean;
  complete?: string;
}) => {
  return {
    type: unDeployActionTypes.START_UNDEPLOY,
    payload: {
      deploy_uuid,
      autoRefresh, // 开启自动刷新
      complete,
    },
  };
};

/**
 * @param params: 请求参数
 * @param rollCallback: 轮询list
 * @param clearRoll: 关闭list轮询
 * @param autoRefresh: 触发翻页，是否继续轮询（上一页关闭轮询，下一页触发轮询）
 */
export const getUnDepolyList = (params: any) => {
  return (dispatch: any) => {
    deployService
      .getUnDeployList({
        deployUuid: params.deployUuid,
        start: params.start,
        status: params.status,
        limit: params.limit,
      })
      .then((res: any) => {
        if (res.data.code === 0) {
          // 此start不是页码，是数据的开始条数(后端无理逻辑)
          // 首次请求需要计算start = count - 20  <0 = 0 (跳转最后一页，第一次见这种逻辑)
          const start =
            res.data.data.count - res.data.data.list &&
            res.data.data.list.length;
          dispatch({
            type: unDeployActionTypes.UPDATE_UNDEPLOY_LIST,
            payload: {
              unDeployList: res.data.data.list || [],
              complete: res.data.data.complete,
              start: start,
              count: res.data.data.count,
            },
          });
          // 不为undeploying，表示卸载已完成，此时clearSetInterval
          // 或者autoRefresh为false
          // if (res.data.complete != 'undeploying') {
          //     // clearRoll()
          //     message.success('卸载已完成！')
          // }
        } else {
          message.error(res.data.msg);
        }
      });
  };
};
/**
 * 刷新当前start，每次请求不改变start
 */
export const getCurrentUnDepolyList = (params: any) => {
  return (dispatch: any) => {
    deployService
      .getUnDeployList({
        deployUuid: params.deployUuid,
        start: params.start,
        status: params.status,
        limit: params.limit,
      })
      .then((res: any) => {
        if (res.data.code === 0) {
          // 此start不是页码，是数据的开始条数(后端无理逻辑)
          // 首次请求需要计算start = count - 20  <0 = 0 (跳转最后一页，第一次见这种逻辑)
          // const start = res.data.data.count - res.data.data.list && res.data.data.list.length
          dispatch({
            type: unDeployActionTypes.UPDATE_CURRENT_UNDEPLOY_LIST,
            payload: {
              unDeployList: res.data.data.list || [],
              complete: res.data.data.complete,
              count: res.data.data.count,
            },
          });
          // 不为undeploying，表示卸载已完成，此时clearSetInterval
          // 或者autoRefresh为false
          // if (res.data.complete != 'undeploying') {
          //     // clearRoll()
          //     message.success('卸载已完成！')
          // }
        } else {
          message.error(res.data.msg);
        }
      });
  };
};

/**
 * 强制停止
 */
export const forceStop = (params: any, callback: any) => {
  return (dispatch: any) => {
    deployService.forceStop(params).then((res: any) => {
      utils.serviceCallback(res, '强制卸载成功！', callback);
    });
  };
};
/**
 * 强制卸载
 */
export const forceUninstall = (params: any, callback: any) => {
  return (dispatch: any) => {
    deployService.forceUninstall(params).then((res: any) => {
      utils.serviceCallback(res, '强制卸载成功！', callback);
    });
  };
};

export const getUnDeployLog = (params: any) => {
  return (dispatch: any) => {
    deployService.getUnDeployLog(params).then((res: any) => {
      if (res.code === 0) {
        // message.success('强制卸载成功！')
        dispatch({
          type: unDeployActionTypes.GET_UNDEPLOY_LOG,
          payload: {
            unDeployLog: res.data || '',
          },
        });
      } else {
        message.error(res.msg);
      }
    });
  };
};

export interface DeployActionTypes {
  updateProductConfig: Function;
  modifyServiceInput: Function;
  modifyInstanceConfig: Function;
  modifyConfigConfig: Function;
  modifyServiceConfig: Function;
  cancelDeploy: Function;
  startUnDeployService: Function;
  getUndeploy: Function;
  getUnDepolyList: Function;
  getCurrentUnDepolyList: Function;
  forceStop: Function;
  forceUninstall: Function;
  getUnDeployLog: Function;
  saveForcedUpgrade: Function;
  getUpgradeType: Function;
  getIsFirstSmooth: Function;
}
