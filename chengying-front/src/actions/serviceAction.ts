import { message } from 'antd';
import { servicePageService } from '@/services';
import { ServiceActions } from '@/constants/actionTypes';
import { Dispatch } from 'redux';

export const getServiceGroup = (
  params: any,
  cb: Function,
  defaultService: string,
  config?: any
) => {
  return (dispatch: Dispatch) => {
    servicePageService.getServiceGroup(params, config).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        // debugger
        let count = 0;
        let first_service = ''; // tslint:disable-line
        let flag = []; // openKey
        let selectedKey = 'sub0_0';
        if (defaultService && defaultService !== '') {
          for (const s in data.data.groups) {
            data.data.groups[s].forEach((o: any, i: number) => {
              if (o.service_name === defaultService) {
                first_service = o;
                selectedKey = `sub${count}_${i}`;
              }
              flag.push(`f${count}`);
            });
            count++;
          }
        } else {
          let session_default = sessionStorage.getItem('service_object');
          Object.keys(data.data.groups).map((item, index) => {
            if (index === 0 && !session_default) {
              // flag = [`f${index}`];
              first_service = data.data.groups[item][0];
            }
            if (session_default) {
              const session_service: {
                group: string;
                objectValue: string;
              } =
                typeof session_default === 'string'
                  ? JSON.parse(session_default)
                  : {
                      group: '',
                      objectValue: '',
                    };
              if (item === session_service.group) {
                data.data.groups[item].map((o, i) => {
                  // debugger
                  if (o.service_name === session_service.objectValue) {
                    first_service = o;
                    selectedKey = `sub${index}_${i}`;
                  }
                });
                console.log(first_service);
              }
              sessionStorage.removeItem('service_object');
            }
            flag.push(`f${index}`);
          });
        }

        // urlParams.service = undefined 的情况
        if (
          !first_service &&
          flag.length === 1 &&
          flag[0] === 'f0' &&
          selectedKey === 'sub0_0'
        ) {
          const paramList = Object.keys(data.data.groups);
          const openKey = paramList[0];
          first_service =
            (openKey &&
              data.data.groups[openKey] &&
              data.data.groups[openKey][0]) ||
            '';
        }
        cb(first_service, flag, selectedKey);
        dispatch({
          type: ServiceActions.GET_SERVICE_GROUP,
          payload: data.data.groups,
        });
      } else {
        message.error(data.msg);
      }
    });
  };
};

export const setServiceGroupStart = (params: any, cb: Function) => {
  return (dispatch: Dispatch) => {
    servicePageService.setServiceGroupStart(params).then((data: any) => {
      serviceCallback(data, '', cb, function () {
        dispatch({
          type: ServiceActions.UPDATE_SERVICE_LIST,
          payload: data.data,
        });
      });
    });
  };
};

/**
 * 获取健康检查列表
 * @param params
 * @returns
 */
export const getHealthCheck = (params: any) => {
  return (dispatch: Dispatch) => {
    servicePageService.getHealthCheck(params).then((res: any) => {
      let data = res.data;
      if (data.code === 0) {
        dispatch({
          type: ServiceActions.GET_HEALTH_LIST,
          payload: data.data,
        });
      } else {
        message.error(data.msg);
      }
    });
  };
};

export const setServiceGroupStop = (params: any, cb: Function) => {
  return (dispatch: Dispatch) => {
    servicePageService.setServiceGroupStop(params).then((data: any) => {
      serviceCallback(data, '', cb, function () {
        dispatch({
          type: ServiceActions.UPDATE_SERVICE_LIST,
          payload: data.data,
        });
      });
    });
  };
};
export const operateExtension = (params: any, cb: Function) => {
  return (dispatch: Dispatch) => {
    servicePageService.operateExtension(params);
  };
};
// EM提醒
export const setResartServiceList = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: ServiceActions.GET_RESART_SERVICE,
      payload: params,
    });
  };
};

export const clearGroupList = () => {
  return {
    type: ServiceActions.CLEAR_SERVICE_LIST,
  };
};

export const getServiceList = (params: any, callback: Function) => {
  return (dispatch: Dispatch) => {
    servicePageService.getServiceList(params).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        dispatch({
          type: ServiceActions.GET_SERVICES,
          payload: data.data.services,
        });
        callback(data.data.services[0]);
      } else {
        message.error(data.msg);
      }
    });
  };
};

export const getHostsList = (params: any, callback: Function) => {
  const service_name = params.service_name;
  return (dispatch: Dispatch) => {
    servicePageService.getServiceHostsList(params).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        dispatch({
          type: ServiceActions.GET_HOSTS,
          payload: {
            use_cloud: data.data.use_cloud,
            hosts: data.data.list,
            service: service_name,
          },
        });
        callback && callback(data.data.use_cloud, data.data.list); // tslint:disable-line
      } else {
        message.error(data.msg);
      }
    });
  };
};

export const getHostConfig = (params: any, callback: Function) => {
  return (dispatch: Dispatch) => {
    servicePageService.getServiceHostConfig(params).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        dispatch({
          type: ServiceActions.GET_HOST_CONFIG,
          payload: data.data.result,
        });
        callback();
      } else {
        message.error(data.msg);
      }
    });
  };
};

export const startService = (params: any, cb: Function) => {
  return () => {
    servicePageService.startService(params).then((data: any) => {
      serviceCallback(data, '服务启动成功！', cb);
    });
  };
};

export const stopService = (params: any, cb: Function) => {
  return () => {
    servicePageService.stopService(params).then((data: any) => {
      serviceCallback(data, '服务关闭成功！', cb);
    });
  };
};

function serviceCallback(
  res: any,
  successMsg?: string,
  callback?: Function,
  dispatch?: any
) {
  res = res.data;
  if (res.code === 0) {
    dispatch && dispatch();
    successMsg && message.success(successMsg);
  } else {
    message.error(res.msg);
  }
  callback && callback();
}

export const disableInstance = (params: any) => {
  const { instance_index, service_name } = params;
  return {
    type: ServiceActions.DISABLE_INSTANCE,
    payload: {
      instance_index: instance_index,
      service_name: service_name,
    },
  };
};

export const startInstance = (params: any, cb: Function) => {
  const service_name = params.service_name;
  return (dispatch: Dispatch) => {
    servicePageService.startServiceInstance(params).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        dispatch({
          type: ServiceActions.START_INSTANCE,
          payload: {
            instance_index: params.instance_index,
            service_name,
          },
        });
        cb && cb(); // tslint:disable-line
      } else {
        message.error(data.msg);
      }
      // servicePageService.getServiceHostsList(params).then((data: any) => {
      //   data = data.data;
      //   if (data.code == 0) {
      //     dispatch({
      //       type: ServiceActions.GET_HOSTS,
      //       payload: {
      //         hosts: data.data.list,
      //         service: service_name
      //       }
      //     });
      //   } else {
      //     message.error(data.msg);
      //   }
      // });
    });
  };
};
export const stopInstance = (params: any, cb: Function) => {
  const service_name = params.service_name;
  return (dispatch: Dispatch) => {
    servicePageService.stopServiceInstance(params).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        dispatch({
          type: ServiceActions.STOP_INSTANCE,
          payload: {
            instance_index: params.instance_index,
            service_name,
          },
        });
        cb();
      } else {
        message.error(data.msg);
      }
      // servicePageService.getServiceHostsList(params).then((data: any) => {
      //     data = data.data;
      //     if (data.code === 0) {
      //         dispatch({
      //             type: ServiceActions.GET_HOSTS,
      //             payload: {
      //                 hosts: data.data.list,
      //                 service: service_name,
      //                 use_cloud: true
      //             }
      //         });
      //     } else {
      //         message.error(data.msg);
      //     }
      // });
    });
  };
};

export const setRedService = (redService: any) => {
  return {
    type: ServiceActions.SET_RED_SERVICE,
    payload: redService,
  };
};

export const addHaRoleToHosts = (roles: any, service_name: any) => {
  // tslint:disable-line
  // debugger;
  return {
    type: ServiceActions.ADD_HA_ROLE,
    payload: {
      roles: roles,
      service_name: service_name,
    },
  };
};

export const getALLProducts = () => {
  return (dispatch: any) => {
    servicePageService.getProductList({ limit: 0 }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const list = [];
        for (const p of res.data.list) {
          if (p.is_current_version === 1) {
            list.push(p);
          }
        }
        dispatch({
          type: ServiceActions.GET_ALL_PRODUCTS,
          payload: list,
        });
      }
    });
  };
};

export const setCurrentProduct = (product: any) => {
  return {
    type: ServiceActions.SET_CURRENT_PRODUCT,
    payload: product,
  };
};

// 保存筛选的配置文件路径
export const setConfigFile = (file: string) => {
  return {
    type: ServiceActions.SET_CONFIG_FILE,
    payload: file,
  };
};

// 获取配置参数列表
export const getServiceConfigList = (params: any, new_service?: any) => {
  return (dispatch: Dispatch, getState: Function) => {
    const {
      ServiceStore: { cur_service, configFile },
    } = getState();
    const service = new_service || cur_service;
    // 获取Config
    servicePageService
      .getServiceConfig({
        ...params,
        file: params.file || configFile,
      })
      .then((response) => {
        const res = response.data;
        const { code, data, msg } = res;
        if (code === 0) {
          const config = {};
          const list = data?.list || [];
          list.forEach((item: any) => {
            config[item.config] = {
              ...item,
              iconType: 0,
            };
          });
          service.Config = config;
          dispatch({
            type: ServiceActions.SET_CURRENT_SERVICE,
            payload: service,
          });
        } else {
          message.error(msg);
        }
      });
  };
};

export const setCurrentService = (service: any, params?: any) => {
  // THERE
  if (params) {
    return (dispatch: any) => {
      // 获取Config
      dispatch(getServiceConfigList(params, service));
    };
  } else {
    return {
      type: ServiceActions.SET_CURRENT_SERVICE,
      payload: service,
    };
  }
};
// tslint:disable-next-line
export const setServiceConfigModify = (cur_service: any) => {
  return {
    type: ServiceActions.SET_CONFIG_MODIFY,
    payload: {
      cur_service: cur_service,
    },
  };
};

export const resetServiceConfig = (params: any) => {
  const product_name = params.product_name;
  const service_name = params.service_name;
  return (dispatch: any) => {
    servicePageService
      .resetServiceConfig({
        product_name: product_name,
        service_name: service_name,
        pid: params.pid,
        product_version: params.product_version,
        field_path: params.field_path,
      })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          dispatch(refreshProductAndService(params));
        } else {
          message.error(res.msg);
        }
      });
  };
};
export const resetMultiServiceConfig = (params: any) => {
  const product_name = params.product_name;
  const service_name = params.service_name;
  return (dispatch: any) => {
    servicePageService
      .resetMultiServiceConfig({
        product_name: product_name,
        service_name: service_name,
        pid: params.pid,
        product_version: params.product_version,
        field_path: params.field_path,
        hosts: params.hosts,
      })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          dispatch(refreshProductAndService(params));
        } else {
          message.error(res.msg);
        }
      });
  };
};

export const modifyProductConfigAll = (params: any, configModify: any) => {
  const product_name = params.product_name;
  const service_name = params.service_name;
  return (dispatch: any) => {
    servicePageService
      .modifyProductConfigAll(
        {
          product_name: product_name,
          service_name: service_name,
        },
        configModify
      )
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          message.success('保存成功');
          dispatch(refreshProductAndService(params));
        } else {
          message.error(res.msg);
        }
      });
  };
};

// 配置信息关联主机
export const modifyMultiAllHosts = (params: any, configModify: any) => {
  const product_name = params.product_name;
  const service_name = params.service_name;
  return (dispatch: any) => {
    servicePageService
      .modifyMultiAllHosts(
        {
          product_name: product_name,
          service_name: service_name,
        },
        configModify
      )
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          message.success('保存成功');
          dispatch(refreshProductAndService(params));
        } else {
          message.error(res.msg);
        }
      });
  };
};
// 部署文件配置信息
export const modifyMultiSingleField = (params: any, configModify: any) => {
  const product_name = params.product_name;
  const service_name = params.service_name;
  return (dispatch: any) => {
    servicePageService
      .modifyMultiAllHosts(
        {
          product_name: product_name,
          service_name: service_name,
        },
        configModify
      )
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          message.success('保存成功');
          dispatch(refreshProductAndService(params));
        } else {
          message.error(res.msg);
        }
      });
  };
};

export const refreshProductAndService = (params: any) => {
  const product_name = params.product_name;
  const service_name = params.service_name;
  return (dispatch: any, getState: Function) => {
    const {
      ServiceStore: { configFile },
    } = getState();
    servicePageService.getProductList({ limit: 0 }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const ps = res.data.list;
        let cur_prod = null; // tslint:disable-line
        let cur_service = null; // tslint:disable-line
        for (const p of ps) {
          if ('product_version' in params) {
            const product_version = params.product_version;
            if (
              p.product_name === product_name &&
              p.product_version === product_version
            ) {
              cur_prod = p;
            }
          } else {
            if (p.product_name === product_name) {
              cur_prod = p;
            }
          }
        }
        if (cur_prod) {
          for (const s in cur_prod.product.Service) {
            if (s === service_name) {
              cur_service = cur_prod.product.Service[s];
            }
          }
        }
        dispatch({
          type: ServiceActions.REFRESH_PROD_SERVICE,
          payload: {
            products: ps,
            cur_product: cur_prod,
            cur_service: cur_service,
          },
        });

        // 重新获取配置
        servicePageService
          .getServiceConfig({
            product_name,
            service_name,
            product_version: cur_prod.product_version,
            file: configFile,
          })
          .then((response) => {
            const { code, data = {} } = response.data;
            if (code === 0) {
              const { list = [] } = data;
              list.forEach((item: any) => {
                cur_service.Config[item.config] = item;
              });
              dispatch({
                type: ServiceActions.SET_CURRENT_SERVICE,
                payload: cur_service,
              });
            }
          });
      } else {
        message.error(res.msg);
      }
    });
  };
};

export const setServiceRollRestartState = (
  service_name: any, // tslint:disable-line
  isRestart: any
) => {
  return {
    type: ServiceActions.SWITCH_SERVICE_RESTART,
    payload: {
      service_name: service_name,
      isRestart: isRestart,
    },
  };
};

export const resetServices = () => {
  return {
    type: ServiceActions.GET_SERVICES,
    payload: [],
  };
};

export interface ServiceActionTypes {
  getServiceGroup: Function;
  setServiceGroupStart: Function;
  setServiceGroupStop: Function;
  setResartServiceList: Function;
  clearGroupList: Function;
  getServiceList: Function;
  getHostsList: Function;
  getHostConfig: Function;
  startService: Function;
  stopService: Function;
  startInstance: Function;
  stopInstance: Function;
  disableInstance: Function;
  resetServices: Function;
  getServiceConfigList: Function;
}
