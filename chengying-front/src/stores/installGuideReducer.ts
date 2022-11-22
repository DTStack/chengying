import { InstallGuideActions } from '@/constants/actionTypes';
import { ActionType, InstallGuideStore } from './modals';
// import { message } from "antd";

export type InstallGuideStoreTypes = InstallGuideStore;
const defaultStore: InstallGuideStore = {
  sqlErro: '',
  runtimeState: 'normal',
  deployState: 'normal',
  step: 0,
  selectedProduct: {},
  productServices: [],
  productPackageList: [],
  clusterList: [], // 集群列表
  namespaceList: [], // namespace列表
  smoothSelectService: {},
  productServicesInfo: {},
  hostInstallToList: [],
  selectedService: {
    DependsOn: [],
  },
  serviceHostList: [],
  resourceState: {},
  paramConfigState: [],
  deployUUID: '',
  deployList: [],
  start: 0,
  count: 0,
  complete: 'deploying',
  stopDeployBySelf: false,
  deployFinished: false,
  selectedServiceList: [], // 选中服务
  unSelectedServiceList: [], // 未选中服务
  installType: 'hosts', // 部署方式
  clusterId: -1, // 选中的集群
  namespace: '', // 选中的命名空间
  baseClusterId: -1, // 选中的依赖集群
  baseClusterInfo: {
    // 依赖集群信息
    baseClusterList: [],
    hasDepends: false,
    dependMessage: '',
  },
  oldHostInfo: {},
  selectProductLine: {}
};

const initialState: InstallGuideStore =
  JSON.parse(localStorage.getItem('installGuideStore')) || defaultStore;

const saveToStorage = (state: any) => {
  try {
    console.log('保存store到storage');
    localStorage.setItem('installGuideStore', JSON.stringify(state));
  } catch (error) {
    console.log(error);
    localStorage.removeItem('installGuideStore');
  }
};

const switchIt = (type: any, payload: any, state: any) => {
  let newState = {};
  switch (type) {
    case InstallGuideActions.NEXT_STEP:
      console.log(1);
      if (state.step < 3) {
        // if (state.step === 0) {
        //   if (JSON.stringify(state.selectedProduct) === "{}") {
        //     message.error("请选择需要安装的产品包。");
        //     newState = state;
        //   }
        // }
        if (state.step === 2) {
          newState = { ...state, step: state.step + 1, deployList: [] };
        } else {
          newState = { ...state, step: state.step + 1 };
        }

        break;
      } else {
        newState = { ...state, step: state.step };
        break;
      }
    case InstallGuideActions.LAST_STEP:
      if (state.step > 0) {
        newState = { ...state, step: state.step - 1 };
        break;
      } else {
        newState = { ...state, step: state.step };
        break;
      }
    case InstallGuideActions.SAVE_INSTALL_INFO:
      newState = { ...state, selectedProduct: payload };
      break;
    case InstallGuideActions.SAVE_INSTALL_TYPE:
      newState = {
        ...defaultStore,
        installType: payload,
      };
      break;
    case InstallGuideActions.SAVE_SELECTED_CLUSTER:
      newState = {
        ...defaultStore,
        installType: state.installType,
        clusterId: payload,
        clusterList: state.clusterList,
        productPackageList: state.productPackageList,
      };
      break;
    case InstallGuideActions.SAVE_SELECTED_NAMESPACE:
      newState = { ...state, namespace: payload };
      break;
    case InstallGuideActions.SAVE_SELECTED_BASECLUSTER:
      newState = { ...state, baseClusterId: payload };
      break;
    case InstallGuideActions.SAVE_SELECTED_SEERVICE:
      newState = { ...state, selectedServiceList: payload };
      break;
    case InstallGuideActions.SAVE_UNSELECTED_SEERVICE:
      newState = { ...state, unSelectedServiceList: payload };
      break;
    case InstallGuideActions.UPDATE_PRODUCT_PACEAGE_LIST:
      // state.productServices = [];
      if (state.step === 0 || state.step === 1) {
        // 在第一步的时候清空已选择的产品包和服务列表
        newState = {
          ...state,
          selectedProduct: {},
          productPackageList: payload,
        };
      } else if (state.step === 2) {
        // 针对升级
        newState = { ...state, productPackageList: payload };
      } else {
        newState = {
          ...state,
        };
      }
      break;
    case InstallGuideActions.RESET_INSTALL_CONFIG:
      newState = { ...state, productServices: [] };
      break;
    case InstallGuideActions.UPDATE_PRODUCT_BASECLUSTER_INFO:
      newState = {
        ...state,
        baseClusterInfo: payload.baseClusterInfo,
        baseClusterId: payload.baseClusterId,
      };
      break;
    case InstallGuideActions.UPDATE_PRODUCT_PACEAGE_SERVICES_LIST:
      newState = { ...state, productServices: payload || [] };
      break;
    case InstallGuideActions.UPDATE_PRODUCT_SERVICES_INFO:
      // const mode = Array.isArray(payload) ? 'AUTO' : 'MANUAL';
      // if (mode === 'AUTO') {
      //     console.log(payload, state)
      //     const { ProductName } = state.selectedProduct;
      //     const { ServiceDisplay } = state.selectedService
      //     if (!ProductName || !ServiceDisplay) return { ...state, productServicesInfo: payload };
      //     const target = payload.find(item => item.productName === ProductName);
      //     const services = Object.keys(target.content).reduce((temp, key) => {
      //         return {
      //             ...temp,
      //             ...target.content[key],
      //         }
      //     }, {})
      //     const service = Object.keys(services)
      //         .map(key => {
      //             return services[key]
      //         })
      //         .find(item => item.ServiceDisplay === ServiceDisplay)
      //     console.log(service);
      //     const ServiceAddr = service.ServiceAddr;
      //     return {
      //         ...state,
      //         productServicesInfo: payload,
      //         selectedService: {
      //             ...service,
      //             ServiceAddr,
      //         },
      //     }
      // }
      newState = { ...state, productServicesInfo: payload };
      break;
    case InstallGuideActions.UPDATE_HOST_INSTALL_TO_LIST:
      newState = { ...state, hostInstallToList: payload };
      break;
    case InstallGuideActions.UPDATE_CLUSTER_LIST:
      newState = { ...state, clusterList: payload };
      break;
    case InstallGuideActions.UPDATE_NAMESPACE_LIST:
      newState = { ...state, namespaceList: payload };
      break;
    case InstallGuideActions.SET_SELECTED_CONFIG_SERVICE:
      newState = { ...state, selectedService: payload };
      break;
    case InstallGuideActions.SET_SMOOTH_SELECT_SERVICE:
      newState = { ...state, smoothSelectService: payload };
      break;
    case InstallGuideActions.UPDATE_SERVICE_HOST_LIST:
      newState = { ...state, serviceHostList: payload, resourceState: {} };
      break;
    case InstallGuideActions.SAVE_RESOURCE_STATE:
      newState = { ...state, resourceState: payload };
      break;
    case InstallGuideActions.SAVE_PARAMS_FIELD_CONFIG_STATE:
      const q = state.paramConfigState;
      q[payload.key] = payload.value;
      newState = { ...state, paramConfigState: q };
      break;
    case InstallGuideActions.START_DEPLOY:
      newState = { ...state, deployUUID: payload, deployFinished: false };
      saveToStorage(newState);
      break;
    case InstallGuideActions.UPDATE_DEPLOY_LIST:
      newState = {
        ...state,
        deployList: payload.deployList,
        complete: payload.complete,
        start: payload.start,
        count: payload.count,
      };
      saveToStorage(newState);
      break;
    case InstallGuideActions.UPDATE_CURRENT_DEPLOY_LIST:
      newState = {
        ...state,
        deployList: payload.deployList,
        complete: payload.complete,
        count: payload.count,
      };
      saveToStorage(newState);
      break;
    case InstallGuideActions.STOP_DEPLOY:
      newState = {
        ...state,
        ...payload,
        // deployUUID: '-1',
        // stopDeployBySelf: true
      };
      saveToStorage(newState);
      break;
    case InstallGuideActions.INIT_INSTALLGUIDE:
      newState = { ...defaultStore };
      saveToStorage(newState);
      break;
    case InstallGuideActions.DEPLOY_FINISHED:
      newState = { ...state, deployFinished: true };
      saveToStorage(newState);
      break;
    case InstallGuideActions.QUIT_CONFIG:
      if (type === 'param') {
        newState = { ...state, resourceState: {} };
      } else {
        newState = { ...state };
      }
      saveToStorage(newState);
      return newState;
    case InstallGuideActions.GO_TO_STEP:
      newState = { ...state, step: payload };
      return newState;
    case InstallGuideActions.SET_DEPLOY_UUID:
      newState = { ...state, deployUUID: payload };
      return newState;
    case InstallGuideActions.EDIT_RUNTIME_STATE:
      newState = { ...state, runtimeState: payload };
      return newState;
    case InstallGuideActions.EDIT_DEPLOY_STATE:
      newState = { ...state, deployState: payload };
      return newState;
    case InstallGuideActions.SET_OLD_HOST_INFO:
      newState = { ...state, oldHostInfo: payload };
      return newState;
    case InstallGuideActions.SET_SQL_ERRO:
      newState = { ...state, sqlErro: payload };
      return newState;
    case InstallGuideActions.SET_SELECT_PRODUCTLINE:
      newState = { ...state, selectProductLine: payload };
      return newState;
    default:
      newState = state;
      // localStorage.setItem('installGuideStore',JSON.stringify(newState));
      break;
  }
  return newState;
};

const red = (state: InstallGuideStore = initialState, action: ActionType) => {
  const { type, payload } = action;
  return switchIt(type, payload, state);
};

export default red;
