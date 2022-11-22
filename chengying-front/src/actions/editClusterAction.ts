import { EditClusterActions } from '@/constants/actionTypes';
import { Dispatch } from 'redux';
import { clusterManagerService } from '@/services';
import { message } from 'antd';

export const setClusterInfo = (params: any) => {
  return (dispatch: any) => {
    dispatch({
      type: EditClusterActions.SET_CLUSTER_INFO,
      payload: params,
    });
  };
};

export const resetClusterInfo = () => {
  return {
    type: EditClusterActions.RESET_CLUSTER_INFO,
    payload: null,
  };
};

export const getTemplate = (params: any) => {
  return async (dispatch: Dispatch) => {
    const res: any = await clusterManagerService.getKubernetesRketemplate({
      version: params.version,
      network_plugin: params.network_plugin,
      cluster: params.name,
    });
    if (res.data.code === 0) {
      dispatch({
        type: EditClusterActions.SET_CLUSTER_INFO,
        payload: {
          yaml: res.data.data,
        },
      });
    } else {
      message.error(res.data.msg);
    }
  };
};

// 提交保存操作
export const clusterSumbmitOperate = (
  params: any,
  isEdit: boolean,
  callback: Function
) => {
  return (dispatch: any) => {
    const type = params.type;
    if (!isEdit) {
      delete params.id;
    }
    delete params.type;
    clusterManagerService
      .clusterSubmitOperate(params, type, isEdit)
      .then((res) => {
        res = res.data;
        if (res.code === 0) {
          dispatch(setClusterInfo(res.data));
          callback && callback();
        } else {
          message.error(res.msg);
        }
      });
  };
};

export interface EditClusterActionTypes {
  setClusterInfo: Function;
  resetClusterInfo: Function;
  clusterSumbmitOperate: Function;
  getTemplate: Function;
}
