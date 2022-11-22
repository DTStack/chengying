import { message } from 'antd';
import { userCenterService } from '@/services';
import { UserCenterActions } from '@/constants/actionTypes';
import { navData } from '@/constants/navData';
import { Dispatch } from 'redux';

export const saveUserName = (userName: string) => {
  return (dispatch: Dispatch) => {
    dispatch({
      type: UserCenterActions.SAVE_USER_INFO,
      payload: userName,
    });
  };
};

export const getAuthorityRouter = (routers: any[], list: any[], codeList) => {
  list.forEach((item: any) => {
    if (codeList[item.code]) {
      routers.push(...item.routers);
    }
    if (item.children.length) {
      getAuthorityRouter(routers, item.children, codeList);
    }
  });
  return routers;
};

export const setRoleAuthorityList = () => {
  return (dispatch: any) => {
    userCenterService.getRoleCodes().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const obj = {};
        const routers = [];
        Array.isArray(res.data) &&
          res.data.forEach((item: string) => {
            obj[item] = true;
          });
        getAuthorityRouter(routers, navData, obj);
        dispatch({
          type: UserCenterActions.SET_ROLE_AUTHORITY_LIST,
          payload: {
            authorityList: obj,
            authorityRouter: routers,
          },
        });
      } else {
        message.error(res.msg);
      }
    });
  };
};

export interface UserCenterActionTypes {
  // tslint:disable-line
  setRoleAuthorityList: Function;
}
