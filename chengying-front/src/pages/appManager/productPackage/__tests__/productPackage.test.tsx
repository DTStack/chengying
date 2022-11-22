import * as React from 'react';
import { Router } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import ProductPackage from '../productPackage';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { authorityList, ServiceMock } from '@/mocks';
import { fireEvent, RenderResult } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import { Service } from '@/services';

jest.mock('@/services');

const defaultProps = {
  authorityList,
};
const reqParams = {
  clusterId: 0,
  productName: undefined,
  parentProductName: undefined,
  productVersion: undefined,
  deploy_status: '',
  'sort-by': 'create_time',
  'sort-dir': 'desc',
  limit: 10,
  start: 0,
};

const INIT_PRODUCT = 'DTinsight';
const NEW_PRODUCT = 'DTEM';
const NEW_COMPONENT = 'DTCommon';

// 安装包管理界面
describe('package manager list', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (Service.getParentProductList as any).mockResolvedValue({
      data: ServiceMock.getParentProductList,
    });
    (Service.getAllProducts as any).mockResolvedValue({
      data: ServiceMock.getAllProducts,
    });
    wrapper = renderWithRedux(
      <Router history={createMemoryHistory()}>
        <ProductPackage {...defaultProps} />
      </Router>,
      reducer,
      { UserCenterStore: { authorityList } }
    );
  });

  // 筛选框
  test('select input unit test', () => {
    const selects = wrapper.container.getElementsByClassName('ant-select');
    const productSelect = selects[0];
    const componentSelect = selects[1];
    const productPlaceholder = productSelect.getElementsByClassName(
      'ant-select-selection__placeholder'
    )[0];
    const productValue = productSelect.getElementsByClassName(
      'ant-select-selection-selected-value'
    )[0];
    const componentPlaceholder = componentSelect.getElementsByClassName(
      'ant-select-selection__placeholder'
    )[0];
    const componentValue = componentPlaceholder.nextSibling; // ul
    expect(productPlaceholder.textContent).toBe('选择产品');
    expect(productValue.textContent).toBe(INIT_PRODUCT);
    expect(componentPlaceholder.textContent).toBe('选择组件');
    expect(componentValue.childNodes.length).toEqual(1);

    /* 产品筛选 */
    // 点击展开下拉框
    fireEvent.click(productSelect);
    // 点击切换
    const dropLi = wrapper.getByText(NEW_PRODUCT);
    fireEvent.click(dropLi);
    expect(productValue.textContent).toBe(NEW_PRODUCT);
    expect(Service.getAllProducts).toHaveBeenLastCalledWith({
      ...reqParams,
      parentProductName: NEW_PRODUCT,
      limit: 0,
    });
    /* 组件筛选 */
    // 点击展开下拉框
    fireEvent.click(componentSelect);
    // 点击切换
    const option = wrapper.getByTestId(`option-${NEW_COMPONENT}`);
    fireEvent.click(option);
    expect(componentValue.childNodes.length).toEqual(2);
    expect(Service.getAllProducts).toHaveBeenLastCalledWith({
      ...reqParams,
      parentProductName: NEW_PRODUCT,
      productName: NEW_COMPONENT,
      limit: 0,
    });
  });

  // 搜索
  test('search input unit test', () => {
    const input = wrapper.getByPlaceholderText('按组件版本号搜索');
    const icon = input.nextSibling.firstChild;
    expect(input).toBeInTheDocument();
    // 搜索
    fireEvent.change(input, { target: { value: '_beta' } });
    fireEvent.click(icon);
    expect(Service.getAllProducts).toHaveBeenLastCalledWith({
      ...reqParams,
      parentProductName: INIT_PRODUCT,
      productVersion: '_beta',
      limit: 0,
    });
  });

  // 上传组件安装包
  test('update unit test', () => {
    const updateBtn = wrapper.container.getElementsByTagName('button')[0];
    expect(updateBtn.textContent).toBe('上传组件安装包');
  });
});
