import * as React from 'react';
import { renderWithRedux } from '@/utils/test';
import { Host } from '../host';
import { Service } from '@/services';
import { ServiceMock, cur_parent_cluster } from '@/mocks';
import { cleanup, fireEvent, RenderResult } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import reducer from '@/stores';

jest.mock('@/services');

const reqParams = {
  cluster_id: cur_parent_cluster.id,
  limit: 10,
  start: 0,
  'sort-by': 'id',
  'sort-dir': 'desc',
  host_or_ip: '',
  is_running: '',
  status: '',
  group: '',
  role: '',
};
const defaultProps = {
  cur_parent_cluster,
  history: null,
  authorityList: [],
};
describe('host modal test', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (Service.getClusterhostGroupLists as any).mockResolvedValue({
      data: ServiceMock.getClusterhostGroupLists,
    });
    (Service.getClusterHostList as any).mockResolvedValue({
      data: ServiceMock.getClusterHostList,
    });
    wrapper = renderWithRedux(<Host {...defaultProps} />, reducer, {
      HeaderStore: {
        cur_parent_cluster,
      },
    });
  });
  afterEach(() => {
    cleanup();
  });
  // 搜索框单测
  test('search input test', () => {
    const searchInput =
      wrapper.container.getElementsByClassName('ant-input-search')[0];
    const input = searchInput.firstChild as HTMLInputElement;
    const icon = searchInput.lastChild.firstChild;
    expect(searchInput.getAttribute('style')).toBe('width: 264px;');
    expect(input.getAttribute('placeholder')).toBe('按主机ip或名称搜索');

    // 搜索
    fireEvent.change(input, {
      target: { value: 'test' },
    });
    fireEvent.click(icon);
    // 算上最开始初始化的一次，依次叠加
    expect(Service.getClusterHostList).toHaveBeenCalledTimes(2);
    expect(Service.getClusterHostList).toHaveBeenLastCalledWith(
      { ...reqParams, host_or_ip: 'test' },
      cur_parent_cluster.type
    );
  });
  // 表格渲染
  test('table render', () => {
    // 表格列
    if (cur_parent_cluster.type === 'kubernetes') {
      const roleColumns = wrapper.getByText('角色');
      const podColumns = wrapper.getByText('POD');
      expect(roleColumns).toBeInTheDocument();
      expect(podColumns).toBeInTheDocument();
    } else {
      const diskColumns = wrapper.getByText('磁盘');
      expect(diskColumns).toBeInTheDocument();
    }

    // 全选
    const checkAll = wrapper.getByTestId('check-all-btn') as HTMLInputElement;
    // 点击全选
    fireEvent.click(checkAll);
    expect(checkAll.checked).toBeTruthy();
    // 所有行被选中
    const tableBody =
      wrapper.container.getElementsByClassName('ant-table-tbody')[0];
    const checkboxs = tableBody.getElementsByClassName('ant-checkbox-input');
    for (const i in checkboxs) {
      if (checkboxs[i] instanceof HTMLInputElement) {
        const checked = (checkboxs[i] as HTMLInputElement).checked;
        expect(checked).toBeTruthy();
      }
    }
  });
});
