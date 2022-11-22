import { EditClusterActions } from '@/constants/actionTypes';
import { ActionType, EditClusterStore } from './modals';

export type EditClusterStoreTypes = EditClusterStore;
const initialState: EditClusterStore = {
  clusterInfo: {
    id: -1,
    type: 'hosts',
    mode: 0,
  },
};
export default (state: EditClusterStore = initialState, action: ActionType) => {
  const { type, payload } = action;
  switch (type) {
    case EditClusterActions.SET_CLUSTER_INFO:
      return {
        ...state,
        clusterInfo: Object.assign({}, state.clusterInfo, payload),
      };
    case EditClusterActions.RESET_CLUSTER_INFO:
      return {
        ...state,
        clusterInfo: initialState.clusterInfo,
      };
    default:
      return state;
  }
};
