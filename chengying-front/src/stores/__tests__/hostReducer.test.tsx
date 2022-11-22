import { HostActions } from '@/constants/actionTypes';
import HostReducer from '../hostReducer';
// import { cleanup } from '@testing-library/react';

const initialState = {
  clusterHostGroupList: [],
  clusterHostList: [],
  filterData: [],
  filterSelectedItem: [],
  hostGroupLists: [],
  hostList: [],
  rootHostList: [],
  pager: {
    current: 0,
    total: 0,
    pageSize: 20,
  },
  searchValue: '',
  selectRows: [],
  selectedHost: {},
  selectedHostServices: [],
  selectedIndex: 0,
};

describe('host reducer', () => {
  test('initial host state', () => {
    const state = HostReducer(initialState, { type: '' });
    expect(state).toEqual(initialState);
  });

  test('get host list', () => {
    const payload = {
      page: 1,
      total: 1,
      data: [],
    };
    const state = HostReducer(initialState, {
      type: HostActions.HOST_GET_LIST,
      payload,
    });
    expect(state.hostList).toEqual(payload.data);
  });
});
