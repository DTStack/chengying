import * as React from 'react';
import { act, cleanup, RenderResult } from '@testing-library/react';
import { renderWithRouter } from '@/utils/renderWithRouter';
import RoleInfo from '../roleInfo';
import '@testing-library/jest-dom/extend-expect';
import { userCenterService } from '@/services';
import { UserCenterMock } from '@/mocks';

jest.mock('@/services');

describe('role info render', () => {
  let wrapper: RenderResult;
  beforeEach(async () => {
    await act(async () => {
      (userCenterService.getAuthorityTree as any).mockResolvedValue({
        data: UserCenterMock.getAuthorityTree,
      });
      wrapper = await renderWithRouter(
        <RoleInfo location={window.location} />,
        { route: '/usercenter/role/view?id=1' }
      );
    });
  });

  afterEach(cleanup);

  test('form item render', () => {
    const labels = wrapper.container.getElementsByClassName(
      'ant-form-item-label'
    );
    const labelNames = ['角色名称', '角色描述', '功能权限'];
    for (const i in labels) {
      if (labels[i] instanceof HTMLElement) {
        expect(labels[i].textContent).toBe(labelNames[i]);
      }
    }
    expect(1).toBe(1);
  });
});
