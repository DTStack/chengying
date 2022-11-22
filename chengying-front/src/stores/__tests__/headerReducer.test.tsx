// import { fromJS } from 'immutable';
import { HeaderActions } from '@/constants/actionTypes';
import { HeaderStoreType } from '../modals';
import HeaderReducer from '../headerReducer';
import { clusterMnagerMock, ServicePageMock, ServiceMock } from '@/mocks';

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

describe('header reducer', () => {
  test('INITIAL_STATE', () => {
    const state = HeaderReducer(undefined, { type: '' });
    expect(state).toEqual(initialState);
  });
  test('SET_PRODUCT_LIST', () => {
    const payload = [];
    const state = HeaderReducer(initialState, {
      type: HeaderActions.SET_PRODUCT_LIST,
      payload,
    });
    expect(state.products).toEqual(payload);
  });
  test('SET_CUR_PRODUCT', () => {
    const payload = ServicePageMock.getCurrentProduct.data;
    const state = HeaderReducer(initialState, {
      type: HeaderActions.SET_CUR_PRODUCT,
      payload,
    });
    expect(state.cur_product).toEqual(payload);
  });
  test('GET_PARENT_PROD_LIST', () => {
    const list = ServiceMock.getClusterProductList.data;
    const cluster = list[0];
    const product = cluster.subdomain.products[0];
    const payload = {
      name: product || '选择产品',
      list: list,
      cluster: {
        id: cluster.clusterId || -1,
        name: cluster.clusterName || '选择集群',
        type: cluster.clusterType || 'hosts',
        mode: cluster.mode || 0,
      },
    };
    const state = HeaderReducer(initialState, {
      type: HeaderActions.GET_PARENT_PROD_LIST,
      payload,
    });
    expect(state.parentProducts).toEqual(payload.list);
    expect(state.cur_parent_product).toEqual(payload.name);
    expect(state.cur_parent_cluster).toEqual(payload.cluster);
  });
  test('SET_CUR_PARENT_PROD', () => {
    const payload = 'DTEM';
    const state = HeaderReducer(initialState, {
      type: HeaderActions.SET_CUR_PARENT_PROD,
      payload,
    });
    expect(state.cur_parent_product).toEqual(payload);
  });
  test('GET_PARENT_CLUSTER_LIST', () => {
    const payload = clusterMnagerMock.getClusterList.data.clusters;
    const state = HeaderReducer(initialState, {
      type: HeaderActions.GET_PARENT_CLUSTER_LIST,
      payload,
    });
    expect(state.parentClusters).toEqual(payload);
  });
  test('SET_CUR_PARENT_CLUSTER', () => {
    const payload = clusterMnagerMock.getClusterList.data.clusters[0];
    const state = HeaderReducer(initialState, {
      type: HeaderActions.SET_CUR_PARENT_CLUSTER,
      payload,
    });
    expect(state.cur_parent_cluster).toEqual(payload);
  });
});
