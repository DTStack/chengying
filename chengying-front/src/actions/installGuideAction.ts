import { Dispatch } from 'redux';
import { InstallGuideActions } from '@/constants/actionTypes';
import { difference, cloneDeep } from 'lodash';
import { installGuideService, deployService } from '@/services';
import { message, notification, Modal } from 'antd';

export const nextStep = () => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.NEXT_STEP,
      payload: {},
    });
  };
};
export const lastStep = () => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.LAST_STEP,
      payload: {},
    });
  };
};
export const saveInstallInfo = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_INSTALL_INFO,
      payload: params,
    });
  };
};

// 保存部署方式
export const saveInstallType = (installType: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_INSTALL_TYPE,
      payload: installType,
    });
  };
};

export const saveSelectCluster = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_SELECTED_CLUSTER,
      payload: params,
    });
  };
};

// 保存命名空间，store and api
export const saveSelectNamespace = (
  urlParams: any,
  reqParams: any,
  isNewSpace: boolean,
  callback?: Function
) => {
  return async (dispatch: Dispatch) => {
    if (isNewSpace) {
      const response = await installGuideService.createNamespace(
        urlParams,
        reqParams
      );
      const res = response.data;
      if (res.code) {
        Modal.confirm({
          title: '创建命名空间发生错误，是否继续使用该命名空间部署应用？',
          content: res.msg,
          icon: 'close-circle',
          className: 'modal-icon-error',
          okType: 'danger',
          okText: '继续',
          cancelText: '取消',
          onOk: () => {
            dispatch({
              type: InstallGuideActions.SAVE_SELECTED_NAMESPACE,
              payload: reqParams.namespace,
            });
            callback && callback();
          },
        });
        return;
      }
    }
    dispatch({
      type: InstallGuideActions.SAVE_SELECTED_NAMESPACE,
      payload: reqParams.namespace,
    });
    callback && callback();
  };
};

// 保存依赖集群
export const saveSelectBaseCluster = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_SELECTED_BASECLUSTER,
      payload: params,
    });
  };
};

export const saveSelectedService = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_SELECTED_SEERVICE,
      payload: params,
    });
  };
};

export const saveUnSelectedService = (params: any) => {
  console.log('我是服务action');
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_UNSELECTED_SEERVICE,
      payload: params,
    });
  };
};

export const getUncheckedService = (params: any) => {
  return (dispatch: Dispatch) => {
    installGuideService.getUncheckedService(params).then((res: any) => {
      if (res.data.code === 0) {
      }
    });
  };
};

// 获取产品包的依赖集群列表
export const getBaseClusterList = (params: any, callback?: Function) => {
  return (dispath: Dispatch) => {
    installGuideService.getBaseClusterList(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const { candidates, targets, message, hasDepends } = res.data;
        const cluster = candidates.find(
          (item: any) => item.clusterId === targets.clusterId
        );
        const payload = {
          baseClusterInfo: {
            baseClusterList: candidates,
            hasDepends: hasDepends,
            dependMessage: message,
          },
          baseClusterId: cluster ? targets.relynamespace : -1,
        };
        dispath({
          type: InstallGuideActions.UPDATE_PRODUCT_BASECLUSTER_INFO,
          payload: payload,
        });
        callback && callback(payload);
      } else {
        message.error(res.msg);
      }
    });
  };
};

// 获取之前未选择的服务
export const getUncheckedServices = (
  params: any,
  defaultCheckboxValue,
  callback?: Function
) => {
  return (dispatch: Dispatch) => {
    installGuideService.getUncheckedService(params).then((res: any) => {
      if (res.data.code === 0) {
        // 将全部服务与未选择服务做对比筛选出已选择服务
        const unSelectService = res.data.data || [];
        const selectService = difference(defaultCheckboxValue, unSelectService);
        dispatch({
          type: InstallGuideActions.SAVE_SELECTED_SEERVICE,
          payload: selectService,
        });
        dispatch({
          type: InstallGuideActions.SAVE_UNSELECTED_SEERVICE,
          payload: unSelectService || [],
        });
        callback && callback();
      }
    });
  };
};

// 指定产品包下的服务
// import B from "public/json/getProductPackageServices";
export const getProductPackageServices = (params: any, callback?: Function) => {
  return (dispatch: Dispatch) => {
    const {
      productName,
      productVersion,
      pid,
      clusterId,
      relynamespace,
      namespace,
    } = params;
    installGuideService
      .getProductPackageServices({
        productName,
        productVersion,
        relynamespace: relynamespace === -1 ? undefined : relynamespace,
      })
      .then((res: any) => {
        if (res.data.code === 0) {
          const defaultCheckboxValue = [];
          res.data.data.map((item) =>
            defaultCheckboxValue.push(item.serviceName)
          );
          let copyData = cloneDeep(res.data.data);
          const arrlist = [];
          const arrbaseProductlist = [];
          for (const i in copyData) {
            if (copyData[i].baseProduct === '') {
              arrlist.push(copyData[i]);
            } else {
              arrbaseProductlist.push(copyData[i]);
            }
          }
          copyData = [...arrlist, ...arrbaseProductlist];
          dispatch({
            type: InstallGuideActions.UPDATE_PRODUCT_PACEAGE_SERVICES_LIST,
            payload: copyData,
          });
          // dispatch({
          //     type: InstallGuideActions.SAVE_SELECTED_SEERVICE,
          //     payload: defaultCheckboxValue
          //   })
          // dispatch({
          // type: InstallGuideActions.SAVE_UNSELECTED_SEERVICE,
          // payload: []
          // })
          // 获取之前未选择的服务
          installGuideService
            .getUncheckedService({ pid, clusterId, namespace })
            .then((res: any) => {
              if (res.data.code === 0) {
                // 将全部服务与未选择服务做对比筛选出已选择服务
                const unSelectService = res.data.data || [];
                const selectService = difference(
                  defaultCheckboxValue,
                  unSelectService
                );
                dispatch({
                  type: InstallGuideActions.SAVE_SELECTED_SEERVICE,
                  payload: selectService,
                });
                dispatch({
                  type: InstallGuideActions.SAVE_UNSELECTED_SEERVICE,
                  payload: unSelectService || [],
                });
                callback && callback();
              }
            });
        } else {
          message.error(res.data.msg);
        }
      });
    // dispatch({
    //   type: InstallGuideActions.UPDATE_PRODUCT_PACEAGE_SERVICES_LIST,
    //   payload: B.data
    // });
  };
};

// 产品包
// import installGuideService from '@/services/installGuideService';
// import { message } from 'antd';
// import A from 'public/json/getProductPaackageList'
export const getProductPackageList = (param: any, callback?: Function) => {
  return (dispatch: Dispatch) => {
    installGuideService.getProductPackagelist(param).then((res: any) => {
      if (res.data.code === 0) {
        dispatch({
          type: InstallGuideActions.UPDATE_PRODUCT_PACEAGE_LIST,
          payload: res.data.data.list,
        });
        callback && callback();
      } else {
        message.error(res.data.msg);
      }
    });
    // dispatch({
    //   type: InstallGuideActions.UPDATE_PRODUCT_PACEAGE_LIST,
    //   payload: A.data.list
    // });
  };
};

export const getProductStepOneList = (
  param: any,
  status?: boolean,
  callback?: Function
) => {
  return (dispatch: Dispatch) => {
    installGuideService.getProductStepOneList(param).then((res: any) => {
      if (res.data.code === 0) {
        res.data.data.list.map((item) => {
          item.isOpen = status;
          item.list = item.children;
          delete item.children;
          return item;
        });
        dispatch({
          type: InstallGuideActions.UPDATE_PRODUCT_PACEAGE_LIST,
          payload: res.data.data.list,
        });
        callback && callback();
      } else {
        message.error(res.data.msg);
      }
    });
    // dispatch({
    //   type: InstallGuideActions.UPDATE_PRODUCT_PACEAGE_LIST,
    //   payload: A.data.list
    // });
  };
};

export const resetInstallGuideConfig = (params?: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.RESET_INSTALL_CONFIG,
      payload: {},
    });
  };
};

export const quitGuide = () => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.QUIT_GUIDE,
      payload: {},
    });
  };
};

// 第三步获取产品下的服务组信息
// import C from "public/json/productServicesInfo";
export const getProductServicesInfo = (
  params: any,
  callback?: Function,
  forcedUpgrade?: any[]
) => {
  return (dispatch: Dispatch) => {
    installGuideService.getProductServicesInfo(params).then((res: any) => {
      if (res.data.code === 0) {
        callback && callback(res.data.data); // tslint:disable-line
        if (forcedUpgrade?.length > 0) {
          let list = res.data.data;
          for (let i of Object.values(list)) {
            for (let values of Object.keys(i)) {
              if (forcedUpgrade.includes(values)) {
                dispatch({
                  type: InstallGuideActions.SET_SMOOTH_SELECT_SERVICE,
                  payload: i[values],
                });
              }
            }
          }
        }
        dispatch({
          type: InstallGuideActions.UPDATE_PRODUCT_SERVICES_INFO,
          payload: res.data.data,
        });
      } else {
        message.error(res.data.msg);
      }
    });
  };
};
// 查看自动编排结果信息
export const getProductServicesInfoForAutoDeploy = (
  params: any,
  callback?: Function
) => {
  return (dispatch: Dispatch) => {
    installGuideService.autoOrchestration(params).then((res: any) => {
      if (res.data.code === 0) {
        callback && callback(res.data.data); // tslint:disable-line
        dispatch({
          type: InstallGuideActions.UPDATE_PRODUCT_SERVICES_INFO,
          payload: res.data.data || [],
        });
      } else {
        message.error(res.data.msg);
      }
    });
  };
};

// 自动部署状态下，获取产品下服务组信息
export const refreshServicesInfoForAutoDeploy = (
  params: any,
  callback?: Function
) => {
  return (dispatch: Dispatch) => {
    installGuideService
      .refreshServicesInfoForAutoDeploy(params)
      .then((res: any) => {
        if (res.data.code === 0) {
          callback && callback(res.data.data); // tslint:disable-line
          console.log('------侧边栏服务列表', res.data.data);
          dispatch({
            type: InstallGuideActions.UPDATE_PRODUCT_SERVICES_INFO,
            payload: res.data.data || [],
          });
        } else {
          message.error(res.data.msg);
        }
      });
  };
};

// 存储可选的集群列表
export const getInstallClusterList = (params: any) => {
  return (dispatch: Dispatch) => {
    installGuideService.getInstallClusterList(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const clusters =
          Array.isArray(res.data.clusters) &&
          res.data.clusters.filter(
            (item) => item.status === 'Running' || item.status === 'Error'
          );
        dispatch({
          type: InstallGuideActions.UPDATE_CLUSTER_LIST,
          payload: clusters,
        });
      } else {
        message.error(res.msg);
      }
    });
  };
};

// 存储可选命名空间
export const getNamespaceList = (params: any) => {
  return (dispatch: Dispatch) => {
    installGuideService.getNamespaceList(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const cluster: any = res.data || {};
        const namespaceList =
          Array.isArray(cluster.namespaces) &&
          cluster.namespaces.map((item) => item.namespace);
        dispatch({
          type: InstallGuideActions.UPDATE_NAMESPACE_LIST,
          payload: namespaceList,
        });
      } else {
        message.error(res.msg);
      }
    });
  };
};

export const setSelectedConfigService = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SET_SELECTED_CONFIG_SERVICE,
      payload: params,
    });
  };
};

// import G from "public/json/serviceHostList";
// import installGuideService from "@/services/installGuideService";
// import { message } from "antd";
// 获取服务下可安装的主机
export const updateServiceHostList = (params: any) => {
  return (dispatch: Dispatch) => {
    installGuideService.updateServiceHostList(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        dispatch({
          type: InstallGuideActions.UPDATE_SERVICE_HOST_LIST,
          payload: res.data.hosts,
        });
      } else {
        message.error(res.msg);
      }
    });
    // dispatch({
    //   type: InstallGuideActions.UPDATE_SERVICE_HOST_LIST,
    //   payload: G.data.hosts
    // });
  };
};

export const saveResourceState = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_RESOURCE_STATE,
      payload: params,
    });
  };
};
export const saveSmoothSelected = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SET_SMOOTH_SELECT_SERVICE,
      payload: params,
    });
  };
};
export const saveParamsFieldConfigState = (params: {
  key: string;
  value: string;
}) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SAVE_PARAMS_FIELD_CONFIG_STATE,
      payload: params,
    });
  };
};

// import H from "public/json/startDeploy";

export const startDeploy = (params: any, callback: Function) => {
  return (dispatch: Dispatch) => {
    installGuideService.startDeploy(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        dispatch({
          type: InstallGuideActions.START_DEPLOY,
          payload: res.data.deploy_uuid,
        });
        callback();
      } else {
        // message.error(res.msg)
        // callback();  // test
        notification.error({
          message: '提示',
          description: res.msg,
          duration: 5,
        });
      }
    });
  };
};

// 自动部署，调用接口
export const startAutoDeploy = (params: any, callback: Function) => {
  console.log(params);
  return (dispatch: Dispatch) => {
    installGuideService.autoDeploy(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        dispatch({
          type: InstallGuideActions.START_DEPLOY,
          payload: res.data.deploy_uuid,
        });
        callback(res.data);
      } else {
        // message.error(res.msg)
        // callback();  // test
        notification.error({
          message: '提示',
          description: res.msg,
          duration: 5,
        });
      }
    });
  };
};

// import I from "public/json/deployStatus";
export const checkDeployStatus = (params: any, callback: Function) => {
  return (dispatch: Dispatch) => {
    console.log('****check deploy status');
    installGuideService.checkDeployStatus(params).then((res: any) => {
      console.log('回调0');
      console.log(res);
      res = res.data.data.list;
      console.log('回调1');
      console.log(res);
      callback(res);
    });
    // let res = I.data.list;
    // callback(res)
  };
};

// import J from 'public/json/deployTaskList'
// import installGuideService from "@/services/installGuideService";
// import { message } from "antd";
/**
 * 增加滚屏数据处理()
 */
export const getDeployList = (params: any) => {
  return (dispatch: Dispatch) => {
    installGuideService.getDeployTaskList(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const start = res.data.count - (res.data.list && res.data.list.length); // 计算下次start开始数
        console.log(res.data);
        dispatch({
          type: InstallGuideActions.UPDATE_DEPLOY_LIST,
          payload: {
            deployList: res.data.list,
            complete: res.data.complete,
            start: start,
            count: res.data.count,
            deployType: res.data.deploy_type,
          },
        });
      } else {
        message.error(res.msg);
      }
    });
  };
};
// 刷新当前start
export const getCurrentDeployList = (params: any) => {
  return (dispatch: Dispatch) => {
    installGuideService.getDeployTaskList(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        dispatch({
          type: InstallGuideActions.UPDATE_CURRENT_DEPLOY_LIST,
          payload: {
            deployList: res.data.list,
            complete: res.data.complete,
            count: res.data.count,
          },
        });
      } else {
        message.error(res.msg);
      }
    });
  };
};
export const stopDeploy = (
  isKubernetes: boolean,
  params: any,
  callback?: any
) => {
  return (dispatch: Dispatch) => {
    if (isKubernetes) {
      const { productName, version, namespace, clusterId } = params;
      const urlParams = {
        productName: productName,
        productVersion: version,
      };
      deployService
        .unDeployService(urlParams, { namespace, clusterId })
        .then((res: any) => {
          const data = res.data;
          if (data.code === 0) {
            console.log('data.data.deploy_uuid', data.data.deploy_uuid);
            dispatch({
              type: InstallGuideActions.STOP_DEPLOY,
              payload: {
                stopDeployBySelf: true,
                deployUUID: data.data.deploy_uuid,
              },
            });
          } else {
            message.error(data.msg);
          }
        });
    } else {
      installGuideService.stopDeploy(params).then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          dispatch({
            type: InstallGuideActions.STOP_DEPLOY,
            payload: {
              stopDeployBySelf: true,
            },
          });
          //  停止部署仍然需要继续发请求
          // callback()
        } else {
          message.error(res.msg);
        }
      });
    }
    // dispatch({
    //   type: InstallGuideActions.STOP_DEPLOY,
    //   payload: ""
    // })
  };
};

export const stopAutoDeploy = (
  isKubernetes: boolean,
  params: any,
  callback?: any
) => {
  return (dispatch: Dispatch) => {
    if (isKubernetes) {
      const { productName, version, namespace, clusterId } = params;
      const urlParams = {
        productName: productName,
        productVersion: version,
      };
      deployService
        .unDeployService(urlParams, { namespace, clusterId })
        .then((res: any) => {
          const data = res.data;
          if (data.code === 0) {
            console.log('data.data.deploy_uuid', data.data.deploy_uuid);
            dispatch({
              type: InstallGuideActions.STOP_DEPLOY,
              payload: {
                stopDeployBySelf: true,
                deployUUID: data.data.deploy_uuid,
              },
            });
          } else {
            message.error(data.msg);
          }
        });
    } else {
      installGuideService.stopAutoDeploy(params).then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          dispatch({
            type: InstallGuideActions.STOP_DEPLOY,
            payload: {
              stopDeployBySelf: true,
            },
          });
          //  停止部署仍然需要继续发请求
          // callback()
        } else {
          message.error(res.msg);
        }
      });
    }
  };
};

/**
 *
 * @param isFinished isFinished用于判断是否是正常部署完成
 * 整个流程判断状态有：
 * 1,stopDeployBySelf(用户手动停止)
 * 2,deployUUID为-1和其他，-1代表没有正在部署的产品
 * 3,deployFinished代表部署已经完成，用于验证完成按钮是否可点击
 *
 * 完成按钮只能在部署完成（deployFinished === true）的时候可点击
 *
 * 停止部署和点击完成按钮之后再进入即回退到第二步，所以点击它们后调用initInstallGuide方法进行store初始化。
 * 在前三步点击退出按钮进行二次提示，退出时initInstallGuideStore，最后一步点击退出后不进行操作，用户再次进入仍进入最后一步。
 */
export const initInstallGuide = (isFinished?: boolean) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.INIT_INSTALLGUIDE,
      payload: {
        // complete: 'deploy fail'
      },
    });
  };
};

export const deployFinished = () => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.DEPLOY_FINISHED,
      payload: '',
    });
  };
};

export const quitConfig = (type: 'param' | 'resource') => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.QUIT_CONFIG,
      payload: type,
    });
  };
};

export const goToStep = (step: number) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.GO_TO_STEP,
      payload: step,
    });
  };
};

export const setDeployUUID = (param: string) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.SET_DEPLOY_UUID,
      payload: param,
    });
  };
};

// 第一步时校验默认仓库是否存在
export const checkDefaultImageStore = (params: any, callback?: Function) => {
  return (dispath: Dispatch) => {
    installGuideService.checkDefaultImageStore(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        callback && callback(res.data.exist);
      } else {
        message.error(res.msg);
      }
    });
  };
};

// 修改运行配置编辑状态
export const editRuntimeState = (param: string) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.EDIT_RUNTIME_STATE,
      payload: param,
    });
  };
};

// 修改部署配置编辑状态
export const editDeployState = (param: string) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: InstallGuideActions.EDIT_DEPLOY_STATE,
      payload: param,
    });
  };
};

export const setOldHostInfo = (params: any) => {
  return {
    type: InstallGuideActions.SET_OLD_HOST_INFO,
    payload: params,
  };
};

// MySQL地址校验
export const setSqlErro = (params: any) => {
  return (dispatch: Dispatch) => {
    installGuideService.checkMySqlAddr(params).then((res: any) => {
      res = res.data;
      if (res.code == 0) {
        dispatch({
          type: InstallGuideActions.SET_SQL_ERRO,
          payload: res.data,
        });
      }
    });
  };
};

// 设置选中产品线
export const setProductLine = (params: any) => {
  return {
    type: InstallGuideActions.SET_SELECT_PRODUCTLINE,
    payload: params,
  };
};

export interface InstallGuideActionTypes {
  // tslint:disable-line
  nextStep: Function;
  lastStep: Function;
  saveInstallInfo: Function; // 记录key-value到store
  saveSelectedService: Function; // 存储第二步选中服务列表
  saveUnSelectedService: Function; // 存储第二步未选中服务列表
  getProductPackageServices: Function; // 获取第二步勾选上的产品下的服务信息
  getBaseClusterList: Function; // 获取第二步产品包的依赖集群
  getProductPackageList: Function; // 获取第二步产品列表
  resetInstallGuideConfig: Function; // 重置勾选的产品下的服务信息
  quitGuide: Function; // 退出
  getProductServicesInfo: Function; // 第三步获取产品下的服务组信息
  getProductServicesInfoForAutoDeploy: Function; // 第三步获取产品下的服务组信息
  refreshServicesInfoForAutoDeploy: Function; // 第三步获取产品下的服务组信息
  getHostIntallToList: Function; // 第三步-获取所有ip
  setSelectedConfigService: Function; // 第三部选中左侧导航的服务事件
  updateServiceHostList: Function; // 第三部获取勾选服务下可安装的主机
  saveResourceState: Function; // 第三部保存主机穿梭框选择值
  saveParamsFieldConfigState: Function; // 第三步保存手风琴参数编辑状态
  startDeploy: Function; // 第三部点击执行部署开始部署并返回uuid
  startAutoDeploy: Function; // 自动部署
  checkDeployStatus: Function; // 第四步轮询检查安装状态
  getDeployList: Function;
  getCurrentDeployList: Function; // 刷新当前start数据
  stopDeploy: Function;
  stopAutoDeploy: Function;
  initInstallGuide: Function;
  deployFinished: Function;
  quitConfig: Function;
  goToStep: Function;
  setDeployUUID: Function;
  getInstallClusterList: Function; // 获取第一步的可选的集群列表
  getNamespaceList: Function; // 获取第一步的可选命名空间
  saveInstallType: Function; // 存储部署方式
  saveSelectCluster: Function; // 存储第一步保存的集群信息
  saveSelectNamespace: Function; // 存储第一步保存的命名空间
  saveSelectBaseCluster: Function; // 存储第二步的依赖集群
  checkDefaultImageStore: Function; // 第一步时校验默认仓库是否存在
  setOldHostInfo: Function;
  setSqlErro: Function;
  saveSmoothSelected: Function;
  getProductStepOneList: Function;
  setProductLine: Function;
  getUncheckedService: Function;
}
