import { unDeployActionTypes } from '@/constants/actionTypes';
import { UnDeployStore } from '../modals';
import UnDeployReducer from '../unDeployReducer';
import { cleanup } from '@testing-library/react';
import { DeployMock } from '@/mocks';

afterEach(cleanup);

const initialState: UnDeployStore = {
  deploy_uuid: '',
  autoRefresh: false,
  complete: 'undeploying',
  unDeployList: [],
  start: 0,
  count: 0,
  unDeployLog: '',
};
describe('undeploy reducer', () => {
  test('INITIAL STATE', () => {
    const state = UnDeployReducer(undefined, { type: '' });
    expect(state).toEqual(initialState);
  });
  test('START_UNDEPLOY', () => {
    const payload = {
      deploy_uuid: DeployMock.unDeployService.data.deploy_uuid,
      autoRefresh: true,
    };
    const state = UnDeployReducer(initialState, {
      type: unDeployActionTypes.START_UNDEPLOY,
      payload,
    });
    expect(state.deploy_uuid).toEqual(payload.deploy_uuid);
    expect(state.autoRefresh).toEqual(payload.autoRefresh);
  });
  test('UPDATE_UNDEPLOY_LIST', () => {
    const data = DeployMock.getUnDeployList.data;
    const start = data.count - (data.list && data.list.length);
    const payload = {
      unDeployList: data.list || [],
      complete: data.complete,
      start: start,
      count: data.count,
    };
    const state = UnDeployReducer(initialState, {
      type: unDeployActionTypes.UPDATE_UNDEPLOY_LIST,
      payload,
    });
    expect(state.unDeployList).toEqual(payload.unDeployList);
    expect(state.complete).toEqual(payload.complete);
    expect(state.start).toEqual(payload.start);
    expect(state.count).toEqual(payload.count);
  });
  test('UPDATE_CURRENT_UNDEPLOY_LIST', () => {
    const data = DeployMock.getUnDeployList.data;
    const payload = {
      unDeployList: data.list || [],
      complete: data.complete,
      count: data.count,
    };
    const state = UnDeployReducer(initialState, {
      type: unDeployActionTypes.UPDATE_CURRENT_UNDEPLOY_LIST,
      payload,
    });
    expect(state.unDeployList).toEqual(payload.unDeployList);
    expect(state.complete).toEqual(payload.complete);
    expect(state.count).toEqual(payload.count);
  });
  test('GET_UNDEPLOY_LOG', () => {
    const payload = {
      unDeployLog: 'log',
    };
    const state = UnDeployReducer(initialState, {
      type: unDeployActionTypes.GET_UNDEPLOY_LOG,
      payload,
    });
    expect(state.unDeployLog).toEqual(payload.unDeployLog);
  });
});
