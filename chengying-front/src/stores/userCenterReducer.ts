import { UserCenterActions } from '@/constants/actionTypes';
import { ActionType, UserCenterStore } from './modals';

export type UserCenterStoreTypes = UserCenterStore;
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

export default function UserCenterReducer(
  state = initialState,
  action: ActionType
) {
  const newState = Object.assign({}, state);
  const { type, payload } = action;
  switch (type) {
    case UserCenterActions.SAVE_USER_INFO:
      newState.userName = payload;
      return newState;
    case UserCenterActions.SET_ROLE_AUTHORITY_LIST:
      newState.authorityList = payload.authorityList;
      newState.authorityRouter = payload.authorityRouter;
      return newState;
    default:
      return state;
  }
}
