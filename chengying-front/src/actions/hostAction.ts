import { Dispatch } from 'redux';
import { Service } from '@/services';
import { HostActions } from '@/constants/actionTypes';
import { message } from 'antd';

export const setFilterSelectedItem = (value: string) => {
  return {
    type: HostActions.HOST_SET_FILTER,
    payload: value,
  };
};
export const getFilterData = (params: any) => {
  // /{"tenantId":1}/;
  return (dispatch: Dispatch) => {
    Service.getFilterData(params).then((data: any) => {
      dispatch({
        type: HostActions.HOST_GET_FILTER,
        payload: data.result.data,
      });
    });
  };
};

// export const changePanelLoading = (value: boolean) => {
//     return (dispatch: Dispatch) => {
//         dispatch({
//             type: HostActions.CHANGE_PANEL_LOADING,
//             payload: value
//         })
//     }
// }
export const getHostList = (params: any, type: string, f?: Function) => {
  console.log(params);
  return (dispatch: Dispatch) => {
    Service.getClusterHostList(params, type).then((res: any) => {
      if (res.data.code === 0) {
        dispatch({
          type: HostActions.UPDATE_HOST_LIST,
          payload: res.data,
        });
        f && f(res.data.data.hosts); // tslint:disable-line
      } else {
        message.error(res.data.msg);
      }
    });
  };
};
/**
 * 获取主机集群下的主机列表
 */
export const getClusterHostList = (params: any, f?: Function) => {
  return (dispatch: Dispatch) => {
    const type = params.cluster_type;
    if (type === 'hosts') {
      delete params.role;
    }
    delete params.cluster_type;

    Service.getClusterHostList(params, type).then((res: any) => {
      if (res.data.code === 0) {
        dispatch({
          type: HostActions.GET_CLUSTER_HOST_LIST,
          payload: res.data,
        });
        f && f(res.data.data); // tslint:disable-line
      } else {
        message.error(res.data.msg);
      }
    });
  };
};
/**
 * 获取主机分组列表
 */
export const gethostGroupLists = (params: any, cb?: any) => {
  return (dispatch: Dispatch) => {
    Service.getClusterhostGroupLists(params).then((res: any) => {
      if (res.data.code === 0) {
        dispatch({
          type: HostActions.GET_HOSTGROUP_LISTS,
          payload: res.data.data || [],
        });
        if (res.data.data && cb) {
          cb(res.data.data);
        }
      } else {
        message.error(res.data.msg);
      }
    });
  };
};
/**
 * 获取集群下主机分组列表
 */
export const getClusterhostGroupLists = (params: any) => {
  return (dispatch: Dispatch) => {
    Service.getClusterhostGroupLists(params).then((res: any) => {
      if (res.data.code === 0) {
        dispatch({
          type: HostActions.GET_CLUSTER_HOSTGROUP_LISTS,
          payload: res.data.data || [],
        });
      } else {
        message.error(res.data.msg);
      }
    });
  };
};

export const getTableData = (params: any) => {
  return (dispatch: Dispatch) => {
    Service.getTableData({ ...params, pageSize: 20 }).then((data: any) => {
      dispatch({
        type: HostActions.HOST_GET_LIST,
        payload: data.result,
      });
    });
  };
};
// 暂时先不发送请求，认为只要主机添加传完参数过来就是添加ip成功的数据
export const updateHostList = (param: any) => {
  return {
    type: HostActions.HOST_GET_LIST,
    payload: param,
  };
};
export const setTableSelectRows = (value: string) => {
  return {
    type: HostActions.HOST_SET_SELECTED_ROWS,
    payload: value,
  };
};
export const setSearchValue = (value: string) => {
  return {
    type: HostActions.HOST_SET_SEARCH_VALUE,
    payload: value,
  };
};
export const setPager = (value: string) => {
  return {
    type: HostActions.HOST_SET_PAGER,
    payload: value,
  };
};
export const showInstallInfo = (params: any) => {
  return Service.showInstallInfo(params)
    .then((data: any) => {
      return data;
    })
    .catch((err: any) => console.error(err));
};

export const getHostInstance = (params: any) => {
  return (dispatch: Dispatch) => {
    Service.getHostInstance(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        dispatch({
          type: HostActions.UPDATE_HOST_INSTANCES_LIST,
          payload: res.data,
        });
      }
    });
  };
};

export const updateSelectedHostInfo = (params: any, i: string) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: HostActions.UPDATE_SELECT_HOST_INFO,
      payload: { params: params, index: i },
    });
  };
};

// import D from 'public/json/selectedHostServices'
export const updateServicesList = (params: any) => {
  console.log(params);
  return (dispatch: Dispatch) => {
    Service.getHostServicesList({
      pid_list: params.pid_list,
      ip: params.ip,
    }).then((res: any) => {
      res = res.data;
      let result = [];
      if (res.code === 0) {
        result = res.data;
      } else {
        message.error(res.msg);
      }
      dispatch({
        type: HostActions.UPDATE_HOST_SERVICES_LIST,
        payload: result,
      });
    });

    // dispatch({
    //     type: HostActions.UPDATE_HOST_SERVICES_LIST,
    //     payload:D.data
    // })
  };
};

export const updateHostsList = (params: any) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: HostActions.UPDATE_HOST_LIST,
      payload: params,
    });
  };
};

// 重置主机列表数据
export const resetHostList = () => {
  return {
    type: HostActions.RESET_HOST_LIST,
    payload: true,
  };
};
export interface HostActionTypes {
  setFilterSelectedItem: Function;
  getFilterData: Function;
  getHostList: Function;
  getTableData: Function;
  updateHostList: Function;
  updateHostsList: Function;
  setTableSelectRows: Function;
  setSearchValue: Function;
  setPager: Function;
  showInstallInfo: Function;
  updateSelectedHostInfo: Function;
  updateServicesList: Function;
  gethostGroupLists: Function;
  getClusterhostGroupLists: Function;
  getClusterHostList: Function;
  resetHostList: Function;
}
