import { deployActionTypes } from '@/constants/actionTypes';
import { ActionType, DeployStore } from './modals';

export type DeployStoreType = DeployStore;

const initialState: DeployStore = {
  upgradeType: sessionStorage.getItem('upgradeType') ??  '',
  isFirstSmooth: sessionStorage.getItem('isFirstSmooth') ? JSON.parse(sessionStorage.getItem('isFirstSmooth')) : false,
  forcedUpgrade: sessionStorage.getItem('forcedUpgrade') ? JSON.parse(sessionStorage.getItem('forcedUpgrade')) : [],
  versionList: [],
  product_name: '',
  product_version: '',
  status: '',
  product: {
    ProductName: '',
    ProductVersion: '',
    Service: [],
  },
  deploy: {
    deploy_uuid: '',
    deploy_status: false,
  },
  deploy_list: [],
  complete: 'deploying',
  deployType: '',
};
export default (state = initialState, action: ActionType) => {
  const newState = Object.assign({}, state);
  const { type, payload } = action;
  switch (type) {
    case deployActionTypes.UPDATE_PRD_CONFIG:
      newState.product_name = payload.product_name;
      newState.product_version = payload.product_version;
      newState.status = payload.status;
      newState.product = payload.product;
      return newState;
    case deployActionTypes.MODIFY_INSTANCE_CONFIG:
      const paths = payload.path.split('.');
      let v = newState.product.Service[payload.pname].Instance;
      for (const i in paths) {
        if (parseInt(i) === paths.length - 1) {
          v[paths[i]] = payload.value;
        } else {
          v = v[paths[i]];
        }
      }
      return newState;
    case deployActionTypes.MODIFY_CONFIG_CONFIG:
      const configPaths = payload.path.split('.');
      let configVersion = newState.product.Service[payload.pname].Config;
      for (const i in configPaths) {
        if (parseInt(i) === configPaths.length - 1) {
          configVersion[configPaths[i]] = payload.value;
        } else {
          configVersion = configVersion[configPaths[i]];
        }
      }
      return newState;
    case deployActionTypes.GET_SERVICE_IPS: // debugger;
      newState.product.Service[payload.service_name].Ips =
        payload.ips.toString();
      return newState;
    case deployActionTypes.UPDATE_IP_CONFIG:
      newState.product.Service[payload.pindex].Ips[payload.index] = payload.ip;
      return newState;
    case deployActionTypes.ADD_IP_TO_CONFIG:
      let ips = newState.product.Service[payload.pindex].Ips;
      if (ips) {
        ips.push('');
      } else {
        ips = [''];
      }
      newState.product.Service[payload.pindex].Ips = ips;
      return newState;
    case deployActionTypes.UPDATE_IP_BY_SERVICE:
      newState.product.Service[payload.service].Ips = payload.ips;
      return newState;
    case deployActionTypes.START_PRD_DEPLOY:
      newState.deploy.deploy_uuid = payload.deploy_uuid;
      newState.deploy.deploy_status = payload.deploy_status;
      newState.deploy_list = [];
      newState.complete = 'deploying';
      return newState;
    case deployActionTypes.UPDATE_DEPLOY_LIST:
      newState.deploy_list = payload.list;
      newState.deployType = payload.deployType;
      // newState.complete = payload.complete;
      return newState;
    case deployActionTypes.UPDATE_DEPLOY_STATUS:
      newState.complete = payload;
      return newState;
    case deployActionTypes.RETURN_TO_CONFIG:
      newState.deploy.deploy_status = false;
      return newState;
    case deployActionTypes.RESET_DEPLOY_STATUS:
      newState.deploy.deploy_status = false;
      newState.deploy_list = [];
      newState.deployType = '';
      return newState;
    case deployActionTypes.SWITH_USE_CLOUD: // debugger;
      newState.product.Service = payload;
      return newState;
    case deployActionTypes.SAVE_FORCED_UPGRADE:
      sessionStorage.setItem('forcedUpgrade', JSON.stringify(payload));
      newState.forcedUpgrade = payload;
      return newState;
    case deployActionTypes.UPGRADE_TYPE:
      sessionStorage.setItem('upgradeType', payload);
      newState.upgradeType = payload;
      return newState;
    case deployActionTypes.GET_FIRST_SMOOTH:
      sessionStorage.setItem('isFirstSmooth', JSON.stringify(payload));
      newState.isFirstSmooth = payload;
      return newState;
    default:
      return state;
  }
};
