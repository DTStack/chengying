import { unDeployActionTypes } from '@/constants/actionTypes';
import { ActionType, UnDeployStore } from './modals';
export type UnDeployStoreType = UnDeployStore;

const initialState: UnDeployStore = {
  deploy_uuid: '',
  autoRefresh: false,
  complete: 'undeploying',
  unDeployList: [],
  start: 0,
  count: 0,
  unDeployLog: '',
};
export default (state = initialState, action: ActionType) => {
  const newState = Object.assign({}, state);
  const { type, payload } = action;
  switch (type) {
    case unDeployActionTypes.START_UNDEPLOY: {
      Object.assign(newState, {
        deploy_uuid: payload.deploy_uuid,
        autoRefresh: payload.autoRefresh,
        complete: payload.complete || 'undeploying'
      });
      return newState;
    }
    case unDeployActionTypes.UPDATE_UNDEPLOY_LIST: {
      Object.assign(newState, {
        unDeployList: payload.unDeployList,
        complete: payload.complete,
        start: payload.start,
        count: payload.count,
      });
      return newState;
    }
    case unDeployActionTypes.UPDATE_CURRENT_UNDEPLOY_LIST: {
      Object.assign(newState, {
        unDeployList: payload.unDeployList,
        complete: payload.complete,
        count: payload.count,
      });
      return newState;
    }
    case unDeployActionTypes.GET_UNDEPLOY_LOG: {
      // 没调用
      newState.unDeployLog = payload.unDeployLog;
      return newState;
    }
    default:
      return state;
  }
};
