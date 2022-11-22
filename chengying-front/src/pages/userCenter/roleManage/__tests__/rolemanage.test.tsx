import * as React from 'react';
import { act, cleanup, RenderResult } from '@testing-library/react';
import { renderWithRouter } from '@/utils/renderWithRouter';
import RoleManage from '../index';
import '@testing-library/jest-dom/extend-expect';
import { userCenterService } from '@/services';
import { UserCenterMock } from '@/mocks';

jest.mock('@/services');

describe('role manage render', () => {
  let wrapper: RenderResult;
  beforeEach(async () => {
    await act(async () => {
      (userCenterService.getRoleList as any).mockResolvedValue({
        data: UserCenterMock.getRoleList,
      });
      wrapper = await renderWithRouter(<RoleManage />, {
        route: '/usercenter/role',
      });
    });
  });

  afterEach(cleanup);

  test('table render', async () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    const columnNames = ['角色名称', '角色描述', '最近修改时间', '操作'];
    for (const i in columns) {
      if (columns[i] instanceof HTMLElement) {
        expect(columns[i].textContent).toBe(columnNames[i]);
      }
    }
  });
});
