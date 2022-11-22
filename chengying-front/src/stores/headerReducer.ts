// import { fromJS } from 'immutable';
import { HeaderActions } from '@/constants/actionTypes';
import { ActionType, HeaderStoreType } from './modals';

export type HeaderStateTypes = HeaderStoreType;
/**
 * init state
 */
const initialState: HeaderStoreType = {
  cur_product: {
    product_id: -1,
    product_name: '选择产品',
  },
  products: [],
  cur_parent_product: '选择产品', // 当前产品
  parentProducts: [],
  cur_parent_cluster: {
    id: -1,
    name: '选择集群',
    type: 'hosts',
    mode: 0,
  },
  parentClusters: [],
};
export default (state: HeaderStoreType = initialState, action: ActionType) => {
  const { type, payload } = action;
  switch (type) {
    case HeaderActions.GET_S_G_LIST:
      // 方法未使用
      return { ...state, sgList: payload };
    case HeaderActions.GET_INSTANCE_LIST:
      // 方法未使用
      return { ...state, instanceList: payload };
    case HeaderActions.SET_PRODUCT_LIST:
      return { ...state, products: payload };
    case HeaderActions.SET_CUR_PRODUCT:
      // 方法被调用，但 headerStore 的 cur_product 没有用过，用的都是serviceStore的（？？）
      return { ...state, cur_product: payload };
    case HeaderActions.GET_PARENT_PROD_LIST:
      return {
        ...state,
        parentProducts: payload.list,
        cur_parent_product: payload.name,
        cur_parent_cluster: payload.cluster,
      };
    case HeaderActions.SET_CUR_PARENT_PROD:
      return { ...state, cur_parent_product: payload };
    case HeaderActions.GET_PARENT_CLUSTER_LIST:
      return { ...state, parentClusters: payload };
    case HeaderActions.SET_CUR_PARENT_CLUSTER:
      return { ...state, cur_parent_cluster: payload };
    default:
      return state;
  }
};
