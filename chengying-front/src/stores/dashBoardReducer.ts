import { DashBoardActions } from '@/constants/actionTypes';
import { ActionType, DashBoardStore } from '@/stores/modals';

export type DashBoardStoreTypes = DashBoardStore;

const initialState: DashBoardStore = {
  dashboards: [],
  tags: [],
  folders: [
    {
      title: 'General',
      id: 0,
    },
  ],
};
export default (state = initialState, action: ActionType) => {
  const { type, payload } = action;
  switch (type) {
    case DashBoardActions.UPDATE_DASH_LIST:
      return { ...state, dashboards: payload.dashboards, tags: payload.tags };
    default:
      return state;
  }
};
