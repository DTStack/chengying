import { combineReducers } from 'redux';
import HeaderStore, { HeaderStateTypes } from './headerReducer';
import ServiceStore, { ServiceStoreTypes } from './serviceReducer';
import HostStore, { HostStoreTypes } from './hostReducer';
import TestStore, { TestStoreTypes } from './test.reducer';
import DashBoardStore, { DashBoardStoreTypes } from './dashBoardReducer';
import InstallGuideStore, {
  InstallGuideStoreTypes,
} from './installGuideReducer';
import addHostStore, { AddHostStoreType } from './addHostReducer';
import DeployStore, { DeployStoreType } from './deployReducer';
import UnDeployStore, { UnDeployStoreType } from './unDeployReducer';
import UserCenterStore, { UserCenterStoreTypes } from './userCenterReducer';
import editClusterStore, { EditClusterStoreTypes } from './editClusterReducer';
export interface AppStoreTypes {
  HeaderStore: HeaderStateTypes;
  ServiceStore: ServiceStoreTypes;
  HostStore: HostStoreTypes;
  TestStore: TestStoreTypes;
  DashBoardStore: DashBoardStoreTypes;
  InstallGuideStore: InstallGuideStoreTypes;
  addHostStore: AddHostStoreType;
  DeployStore: DeployStoreType;
  UnDeployStore: UnDeployStoreType;
  UserCenterStore: UserCenterStoreTypes;
  editClusterStore: EditClusterStoreTypes;
}
const rootStore = combineReducers({
  HeaderStore,
  ServiceStore,
  HostStore,
  TestStore,
  DashBoardStore,
  InstallGuideStore,
  addHostStore,
  DeployStore,
  UnDeployStore,
  UserCenterStore,
  editClusterStore,
});

export default rootStore;
