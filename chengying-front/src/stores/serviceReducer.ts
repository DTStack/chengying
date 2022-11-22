import { ServiceActions } from '@/constants/actionTypes';
import { ActionType, ServiceStore } from './modals';

export type ServiceStoreTypes = ServiceStore;
const initialState: ServiceStore = {
  products: [], // 组件列表
  cur_product: {}, // 当前组件
  cur_service: {}, // 当前服务
  services: [], // 当前组件服务组
  hosts: [], // 当前服务包含主机列表
  sgList: [],
  config: '',
  configModify: {},
  use_cloud: false,
  configFile: 'all', // 运行配置 - 配置文件
  restartService: {
    count: 0,
    list: [],
  }, // EM提醒
  redService: {
    count: 0,
    list: [],
  },
};

export default function ServiceListReducer(
  state = initialState,
  action: ActionType
) {
  const { type, payload } = action;
  const newState = Object.assign({}, state);
  switch (type) {
    case ServiceActions.GET_RESART_SERVICE:
      return { ...state, restartService: payload };
    case ServiceActions.SET_RED_SERVICE:
      return { ...state, redService: payload };
    case ServiceActions.GET_SERVICE_GROUP:
      return { ...state, services: payload };
    case ServiceActions.UPDATE_SERVICE_LIST:
      return { ...state, sgList: payload };
    case ServiceActions.CLEAR_SERVICE_LIST:
      return { ...state, sgList: [] };
    case ServiceActions.GET_SERVICES:
      newState.services = payload;
      return newState;
    case ServiceActions.GET_HOSTS:
      // const serviceHost = newState.cur_product.product.Service[payload.service] || {};
      newState.hosts = payload.hosts;
      newState.use_cloud = payload.use_cloud;
      return newState;
    case ServiceActions.GET_HOST_CONFIG:
      newState.config = payload;
      return newState;
    case ServiceActions.DISABLE_INSTANCE:
      newState.hosts[payload.instance_index].isDisable = true;
      return newState;
    case ServiceActions.START_INSTANCE:
      newState.hosts[payload.instance_index].status = 'running';
      newState.hosts[payload.instance_index].isDisable = false;
      return newState;
    case ServiceActions.STOP_INSTANCE:
      newState.hosts[payload.instance_index].status = 'stopped';
      newState.hosts[payload.instance_index].isDisable = false;
      return newState;
    case ServiceActions.ADD_HA_ROLE:
      // debugger;
      for (const h of newState.hosts) {
        for (const r in payload.roles) {
          if (h.agent_id === r) {
            h.haRole = payload.roles[r];
          }
        }
      }
      console.log(newState.hosts);
      return newState;
    case ServiceActions.SET_CONFIG_MODIFY:
      return { ...state, cur_service: payload.cur_service };
    /**
     * 重构计划，梳理产品服务store结构
     */
    case ServiceActions.GET_ALL_PRODUCTS:
      return { ...state, products: payload };
    case ServiceActions.SET_CURRENT_PRODUCT:
      return { ...state, cur_product: payload };
    case ServiceActions.SET_CURRENT_SERVICE:
      return { ...state, cur_service: payload };
    case ServiceActions.REFRESH_PROD_SERVICE:
      const { products, cur_product, cur_service } = payload;
      return {
        ...state,
        products: products,
        cur_product: cur_product,
        cur_service: cur_service,
      };
    case ServiceActions.SWITCH_SERVICE_RESTART:
      const { service_name, isRestart } = payload;
      newState.cur_product.product.Service[service_name].isRestart = isRestart;
      return newState;
    // 保存配置文件路径
    case ServiceActions.SET_CONFIG_FILE:
      newState.configFile = payload;
      return newState;
    default:
      return state;
  }
}
