import { HostActions } from '@/constants/actionTypes';
import { ActionType, HostStore } from './modals';

export type HostStoreTypes = HostStore;
const initialState: HostStore = {
  filterSelectedItem: [],
  filterData: [],
  hostList: [],
  clusterHostList: [],
  pager: {
    total: 0,
    pageSize: 20,
  },
  selectRows: [],
  searchValue: '',
  selectedHost: {},
  selectedHostServices: [],
  selectedIndex: 0,
  hostGroupLists: [], // 主机分组
  clusterHostGroupList: [], // 集群下主机分组
};
export default (state: HostStore = initialState, action: ActionType) => {
  const { type, payload } = action;
  switch (type) {
    case HostActions.HOST_HOME:
      return state + payload;
    case HostActions.HOST_SET_FILTER:
      return { ...state, filterSelectedItem: payload };
    case HostActions.GET_HOSTGROUP_LISTS:
      return { ...state, hostGroupLists: payload };
    case HostActions.GET_CLUSTER_HOSTGROUP_LISTS:
      return { ...state, clusterHostGroupList: payload };
    case HostActions.HOST_GET_FILTER:
      return { ...state, filterData: payload };
    case HostActions.UPDATE_HOST_LIST:
      return {
        ...state,
        hostList: payload.data.hosts,
        pager: { total: payload.data.count, pageSize: 20 },
      };
    case HostActions.GET_CLUSTER_HOST_LIST:
      return {
        ...state,
        clusterHostList: payload.data.hosts,
        pager: { total: payload.data.count, pageSize: 20 },
      };
    case HostActions.HOST_GET_LIST:
      return {
        ...state,
        hostList: payload.data,
        rootHostList: payload.data,
        pager: { ...state.pager, total: payload.total, current: payload.page },
      };
    case HostActions.HOST_SET_SELECTED_ROWS:
      return { ...state, selectRows: payload };
    case HostActions.HOST_SET_SEARCH_VALUE:
      return { ...state, searchValue: payload };
    case HostActions.HOST_SET_PAGER:
      return { ...state, pager: payload };
    case HostActions.UPDATE_SELECT_HOST_INFO:
      return {
        ...state,
        selectedHost: payload.params,
        selectedIndex: payload.index,
      };
    case HostActions.UPDATE_HOST_SERVICES_LIST:
      return { ...state, selectedHostServices: payload };
    case HostActions.RESET_HOST_LIST:
      return { ...state, clusterHostGroupList: [], clusterHostList: [] };
    default:
      return state;
  }
};
