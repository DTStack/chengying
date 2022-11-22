import { AddHostActions } from '@/constants/actionTypes';
import { AddHostStore, ActionType } from '@/stores/modals';

export type AddHostStoreType = AddHostStore;

const initialState: AddHostStore = {
  current: 0,
  disabled: true,
  hostArr: [],
  installMsg: [],
  hostPKArr: [],
  forms: {},
};

// tslint:disable
export default (state = initialState, action: ActionType) => {
  const { type, payload } = action;
  switch (type) {
    case 'getnewinstallMsg':
      if (payload) {
        return state.installMsg;
      }
      return state;
    case AddHostActions.update_hostArr:
      if (Array.isArray(payload)) {
        return { ...state, hostArr: payload };
      } else if (
        Object.prototype.toString.call(payload) === '[object Object]'
      ) {
        return { ...state, hostArr: [...state.hostArr, payload] };
      } else {
        console.error('hostArr参数传入有误');
        return state;
      }
    case AddHostActions.update_installMsg:
      // console.log('reducer:', payload);
      // debugger;
      if (Array.isArray(payload)) {
        return { ...state, installMsg: payload };
      } else if (
        Object.prototype.toString.call(payload) === '[object Object]'
      ) {
        return { ...state, installMsg: [...state.installMsg, payload] };
      } else {
        console.error('installMsg参数传入有误');
        return state;
      }

    case AddHostActions.set_disabled:
      return { ...state, disabled: payload };
    case AddHostActions.to_page:
      return { ...state, current: payload };
    // case AddHostActions.valid_host_by_pk:
    //     return { ...state, hostPKArr: payload };
    case AddHostActions.valid_host:
      return payload;

    case AddHostActions.install_host_by_pass_word:
      if (payload) {
        // let { code, msg, data } = res; let is_succ;
      }
      // console.log('pwd:', payload);
      return { ...state, installMsg: payload };
    case AddHostActions.install_host_by_pk:
      return { ...state, installPKArr: payload };
    case 'updateForm':
      return { ...state, forms: payload };
    case AddHostActions.reset_state:
      return {
        current: 0,
        disabled: true,
        hostArr: [],
        installMsg: [],
        hostPKArr: [],
      };
    default:
      return state;
  }
};
