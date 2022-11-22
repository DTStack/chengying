import { DashBoardActions } from '@/constants/actionTypes';
import { DashBoardStore } from '@/stores/modals';
import { cleanup } from '@testing-library/react';
import dashBoardReducer from '../dashBoardReducer';

afterEach(cleanup);

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

describe('dashboard reducer', () => {
  // 初始化
  test('initial state', () => {
    const state = dashBoardReducer(undefined, { type: '' });
    expect(state).toEqual(initialState);
  });
  // 仪表盘列表
  test('UPDATE_DASH_LIST', () => {
    const payload = {
      dashboards: [
        {
          title: 'General',
          id: 0,
          list: [
            {
              checked: false,
              id: 17,
              isStarred: false,
              tags: ['Batch', 'DTApp', 'DTinsight'],
              title: 'Batch_Overview',
              type: 'dash-db',
              uid: 'XB0Gua1izBatch',
              uri: 'db/batch_overview',
              url: '/grafana/d/XB0Gua1izBatch/batch_overview',
            },
          ],
          isOpen: true,
          checked: false,
        },
      ],
      tags: ['Batch', 'DTApp', 'DTinsight'],
    };
    const state = dashBoardReducer(initialState, {
      type: DashBoardActions.UPDATE_DASH_LIST,
      payload: {
        dashboards: payload.dashboards,
        tags: payload.tags,
      },
    });
    expect(state.dashboards).toEqual(payload.dashboards);
    expect(state.tags).toEqual(payload.tags);
  });
});
