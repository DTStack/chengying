import { EditClusterActions } from '@/constants/actionTypes';
import { EditClusterStore } from '../modals';
import EditClusterReducer from '../editClusterReducer';
import { cleanup } from '@testing-library/react';

afterEach(cleanup);

const initialState: EditClusterStore = {
  clusterInfo: {
    id: -1,
    type: 'hosts',
    mode: 0,
  },
};

describe('edit cluster reducer', () => {
  test('INITIAL_STATE', () => {
    const state = EditClusterReducer(undefined, { type: '' });
    expect(state).toEqual(initialState);
  });
  test('SET_CLUSTER_INFO', () => {
    const payload = {
      desc: '',
      id: 2,
      name: 'RanRan_1',
      network_plugin: '',
      tags: '',
      version: 'v1.16.3',
      yaml: '',
    };
    const state = EditClusterReducer(initialState, {
      type: EditClusterActions.SET_CLUSTER_INFO,
      payload,
    });
    const clusterInfo = Object.assign({}, initialState.clusterInfo, payload);
    expect(state.clusterInfo).toEqual(clusterInfo);
  });

  test('RESET_CLUSTER_INFO', () => {
    const state = EditClusterReducer(initialState, {
      type: EditClusterActions.RESET_CLUSTER_INFO,
      payload: null,
    });
    expect(state.clusterInfo).toEqual(initialState.clusterInfo);
  });
});
