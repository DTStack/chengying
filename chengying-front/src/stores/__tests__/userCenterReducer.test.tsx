import { getAuthorityRouter } from '@/actions/userCenterAction';
import { UserCenterActions } from '@/constants/actionTypes';
import { navData } from '@/constants/navData';
import { UserCenterMock } from '@/mocks';
import { UserCenterStore } from '../modals';
import UserCenterReducer from '../userCenterReducer';
import { cleanup } from '@testing-library/react';

afterEach(cleanup);

const initialState: UserCenterStore = {
  userName: '',
  userEmail: '', // 用户邮箱
  company: '', // 公司
  phone: '', // 手机号
  createTime: '', // 创建时间
  status: -1, // 状态
  role: '', // 角色
  authorityList: {}, // 权限表
  authorityRouter: [], // 权限路由
};

describe('usercenter reducer', () => {
  test('INITIAL_STATE', () => {
    const state = UserCenterReducer(undefined, { type: '' });
    expect(state).toEqual(initialState);
  });
  test('SAVE_USER_INFO', () => {
    const payload = 'username';
    const state = UserCenterReducer(initialState, {
      type: UserCenterActions.SAVE_USER_INFO,
      payload: payload,
    });
    expect(state.userName).toEqual(payload);
  });
  test('SET_ROLE_AUTHORITY_LIST', () => {
    const codes = UserCenterMock.getRoleCodes.data;
    const obj = {};
    const routers = [];
    codes.forEach((item: string) => {
      obj[item] = true;
    });
    getAuthorityRouter(routers, navData, obj);
    const payload = {
      authorityList: obj,
      authorityRouter: routers,
    };
    const state = UserCenterReducer(initialState, {
      type: UserCenterActions.SET_ROLE_AUTHORITY_LIST,
      payload,
    });
    expect(state.authorityList).toEqual(payload.authorityList);
    expect(state.authorityRouter).toEqual(payload.authorityRouter);
  });
});
