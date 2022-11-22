import * as React from 'react';
import ComponentList from '../componentList';
import { createMemoryHistory } from 'history';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { authorityList, ServiceMock } from '@/mocks';
import { fireEvent, RenderResult, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import { Router } from 'react-router-dom';

jest.mock('@/services');

const history = createMemoryHistory();
history.push('/appmanage/productpackage');

const { list, count } = ServiceMock.getAllProducts.data;
const defaultProps = {
  location: history.location,
  getDataList: jest.fn(),
  getParentProductsList: jest.fn(),
  componentData: {
    list,
    count,
  },
  authorityList,
};

const reqParams = {
  clusterId: 0,
  productName: undefined,
  parentProductName: undefined,
  productVersion: undefined,
  deploy_status: undefined,
  'sort-by': 'create_time',
  'sort-dir': 'desc',
  limit: 10,
  start: 0,
};

describe('component list render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = renderWithRedux(
      <Router history={history}>
        <ComponentList {...defaultProps} />
      </Router>,
      reducer,
      { UserCenterStore: { authorityList } }
    );
  });

  // 表格列
  test('table columns', () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    const columnsKeys = [
      '组件名称',
      '组件版本号',
      '安装包类型',
      '上传时间',
      '上传人',
      '查看',
      '操作',
    ];
    for (const i in columns) {
      if (columns[i] instanceof HTMLElement) {
        const text = columns[i].textContent;
        expect(text).toBe(columnsKeys[i]);
      }
    }
    // 排序 - 组件版本号为例 product_version
    const productVersion = columns[1].parentElement.parentElement.parentElement;
    fireEvent.click(productVersion);
    expect(defaultProps.getDataList).toHaveBeenCalled();
    expect(defaultProps.getDataList).toHaveBeenCalledWith({
      ...reqParams,
      'sort-by': 'product_version',
      'sort-dir': 'asc',
    });
  });

  // 表格行
  test('table data render', async () => {
    const row = wrapper.container.getElementsByClassName('ant-table-row')[0];
    const result = list[0];
    const { product_type } = result;

    // 安装包类型
    const productType = row.children.item(2);
    const type = product_type === 1 ? 'Kubernetes包' : '传统包';
    expect(productType.textContent).toBe(type);

    // 操作
    const actions = row.lastChild;
    const deployBtn = actions.firstChild;
    expect(deployBtn.textContent).toBe('部署');
    // 没有部署的包 可以删除
    if (!(result as any).is_current_version) {
      const deleteBtn = actions.lastChild;
      expect(deleteBtn.textContent).toBe('删除');
      // 删除
      fireEvent.click(deleteBtn);
      // 二次确认弹框
      const confirmModal = screen.getByText('确定删除此安装包？');
      expect(confirmModal).toBeInTheDocument();
    }
  });
});
