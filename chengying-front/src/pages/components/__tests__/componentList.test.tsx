import * as React from 'react';
import ComponentList from '../componentList';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { Service } from '@/services';
import { authorityList, ServiceMock } from '@/mocks';
import { RenderResult } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';

jest.mock('@/services');

const defaultProps = {
  location: null,
  clusterList: [],
  shouldNameSpaceShow: false,
  mode: 0,
  getParentClustersList: jest.fn(),
  HeaderStore: { cur_parent_cluster: {} },
};
describe('componet list render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (Service.getAllProducts as any).mockResolvedValue({
      data: ServiceMock.getAllProducts,
    });
    wrapper = renderWithRedux(<ComponentList {...defaultProps} />, reducer, {
      UserCenterStore: { authorityList },
    });
  });

  // 表格列
  test('table columns', () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    const columnKeys = [
      '组件名称',
      '组件版本号',
      '安装包类型',
      '部署状态',
      '部署时间',
      '部署人',
      '查看',
      '操作',
    ];
    for (const i in columns) {
      if (columns[i] instanceof HTMLElement) {
        expect(columns[i].textContent).toMatch(columnKeys[i]);
      }
    }
  });
});
