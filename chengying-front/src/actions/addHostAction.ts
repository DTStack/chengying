import { AddHostActions } from '@/constants/actionTypes';
import { message } from 'antd';
import { addHostService } from '@/services';

export const hostAddAction = {
  updateHostArr: (arr: any) => {
    return {
      type: AddHostActions.update_hostArr,
      payload: arr,
    };
  },
  updateInstallMsg: (arr: any) => {
    return {
      type: AddHostActions.update_installMsg,
      payload: arr,
    };
  },
  setDisabled: (v: any) => {
    return {
      type: AddHostActions.set_disabled,
      payload: v,
    };
  },
  jumpPage: (page: any) => {
    return {
      type: AddHostActions.to_page,
      payload: page,
    };
  },
  getNewState: (param: any) => {
    return (dispatch: any) => {
      dispatch({
        type: 'getnewinstallMsg',
        payload: param,
      });
    };
  },

  polling: (param: any) => {
    return (dispatch: any) => {
      addHostService.pwdInstallUrl(param).then((data: any) => {
        dispatch({
          type: AddHostActions.install_host_by_pass_word,
          payload: data.data,
        });
      });
    };
  },
  // 卸载阶段重置所有的state
  resetState: () => {
    return {
      type: AddHostActions.reset_state,
    };
  },
};

export function connectHost(param: any, callback?: Function) {
  addHostService.connectHost(param.item, param.type).then((res: any) => {
    return connectError(res, callback);
  });
}

export function installHost(param: any) {
  return addHostService.installHost(param.item, param.type).then((res: any) => {
    return connectError(res);
  });
}
export function checkInstall(aid: any) {
  return addHostService.checkInstallUrl({ aid: aid }).then((res: any) => {
    return connectError(res);
  });
}

function connectError(res: any, callback?: Function) {
  if (res === undefined) {
    message.error('服务器连接不上');
    return false;
  }
  if (callback) {
    return callback(res);
  }
  return res;
}
