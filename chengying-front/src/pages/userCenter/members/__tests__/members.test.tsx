import * as React from 'react';
import { cleanup, fireEvent, RenderResult } from '@testing-library/react';
import { renderWithRedux } from '@/utils/test';
import Members from '../index';
import reducer from '@/stores';
import { authorityList, UserCenterMock, clusterMnagerMock } from '@/mocks';
import { userCenterService, clusterManagerService } from '@/services';
import '@testing-library/jest-dom/extend-expect';

jest.mock('@/services');

const defaultProps = {
  authorityList,
};
const reqParams = {
  start: 0,
  limit: 10,
  username: '',
  'sort-by': 'create_time',
  'sort-dir': 'desc',
};
describe('members render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (userCenterService.getMembers as any).mockResolvedValue({
      data: UserCenterMock.getMembers,
    });
    (clusterManagerService.getClusterLists as any).mockResolvedValue({
      data: clusterMnagerMock.getClusterList,
    });
    wrapper = renderWithRedux(<Members {...defaultProps} />, reducer, {
      UserCenterStore: { authorityList },
    });
  });

  afterEach(cleanup);

  // 搜索框
  test('search input unit test', () => {
    const searchInput = wrapper.getByPlaceholderText('按账号/姓名搜索');
    expect(searchInput).toBeInTheDocument();
    // 搜索
    const searchName = 'test';
    fireEvent.change(searchInput, { target: { value: searchName } });
    fireEvent.keyDown(searchInput, { keyCode: 13 });
    expect(userCenterService.getMembers).toHaveBeenCalledTimes(2);
    expect(userCenterService.getMembers).toHaveBeenLastCalledWith({
      ...reqParams,
      username: searchName,
    });
  });

  // 创建账号
  test('create member', () => {
    const createBtn = wrapper.getByRole('button');
    expect(createBtn.textContent).toBe('创建账号');
    // click
    fireEvent.click(createBtn);
    const userInput = wrapper.getByPlaceholderText('请输入账号，允许中英文');
    expect(userInput).toBeInTheDocument();
  });

  // 表格
  test('list render', () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    // columns render
    for (const i in columns) {
      const titles = [
        '账号',
        '姓名',
        '邮箱',
        '手机号',
        '账号状态',
        '角色',
        '创建时间',
        '操作',
      ];
      const column = columns[i];
      if (column instanceof HTMLElement) {
        expect(column.textContent).toBe(titles[i]);
      }
    }

    // 角色 & action
    const rows = wrapper.container.getElementsByClassName('ant-table-row');
    const listMock = UserCenterMock.getMembers.data.list;
    for (const i in rows) {
      const row = rows[i];
      if (row instanceof HTMLElement) {
        const tds = row.getElementsByTagName('td');
        // 角色
        const role = tds[5];
        expect(role.textContent).toBe(listMock[i].role_name);
        // 操作 管理员账号 - 仅可重置密码
        const action = tds[tds.length - 1];
        if (listMock[i].role_name === 'Administrator') {
          expect(action.textContent).toBe('重置密码');
        } else {
          // 以 operator 为例
          if (listMock[i].role_name === 'Cluster Operator') {
            testActions(action);
          }
        }
      }
    }

    // 操作 - 以禁用为例
    async function testActions(action: HTMLElement) {
      const actions = action.getElementsByTagName('a');
      const forbiddenBtn = actions[0];
      /* 禁用 */
      expect(forbiddenBtn).toBeInTheDocument();
      expect(forbiddenBtn.textContent).toBe('禁用');
      // 点击出弹框
      fireEvent.click(forbiddenBtn);
      const modalTitle = wrapper.getByText('确定禁用该成员？');
      expect(modalTitle).toBeInTheDocument();
      // 点击确认禁用
      const okBtn = modalTitle.parentElement.nextElementSibling.lastChild;
      (userCenterService.taggleStatus as any).mockResolvedValue({
        data: { code: 0, data: null, msg: 'ok' },
      });
      (userCenterService.getMembers as any).mockResolvedValue({
        data: UserCenterMock.getMembers,
      });
      fireEvent.click(okBtn);
      expect(userCenterService.getMembers).toHaveBeenLastCalledWith(reqParams);
    }
  });

  test('members Snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });
});
