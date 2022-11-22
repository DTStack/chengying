import * as React from 'react';
import Container from '../container';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { Service } from '@/services';
import { ServiceMock } from '@/mocks';
import { fireEvent, RenderResult, cleanup } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';

jest.mock('@/services');

describe('app deploy list', () => {
  let wrapper: RenderResult;
  // 初始数据
  const clusterList = ServiceMock.getClusterProductList.data;
  const firstCluster = clusterList[0];
  // const secondCluster = clusterList[1];
  beforeEach(() => {
    (Service.getClusterProductList as any).mockResolvedValue({
      data: ServiceMock.getClusterProductList,
    });
    (Service.getAllProducts as any).mockResolvedValue({
      data: ServiceMock.getAllProducts,
    });
    wrapper = renderWithRedux(<Container />, reducer, {});
  });

  afterEach(cleanup);

  test('const Snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });

  // 下拉筛选
  test('select unit test', () => {
    const selects = wrapper.container.getElementsByClassName(
      'ant-select-selection__rendered'
    );
    // render
    for (const i in selects) {
      const list = [
        {
          placeholder: '选择集群',
          value: '选择集群', // DTLogger
        },
        {
          placeholder: '选择产品',
          value: firstCluster.subdomain.products[0], // DTEM
        },
        {
          placeholder: '选择组件',
          value: '',
        },
      ];
      if (selects[i] instanceof HTMLElement) {
        expect(selects[i].firstChild.textContent).toBe(list[i].placeholder);
        list[i].value &&
          expect(selects[i].lastChild.textContent).toBe(list[i].value);
      }
    }
    /* 下拉选择 - 以集群为例 */
    // const clusterSelect = wrapper.container.getElementsByClassName('ant-select')[0];
    // fireEvent.click(clusterSelect);
    // 下拉菜单
    // const clusterOptions = wrapper.getByTestId(`cluster-option-${secondCluster.clusterName}`);
    // expect(clusterOptions).toBeInTheDocument();
    // 选择
    // fireEvent.click(clusterOptions);
    // expect(selects[0].lastChild.textContent).toBe(secondCluster.clusterName);
    expect(Service.getAllProducts).toHaveBeenCalled();
    // expect(Service.getAllProducts).toHaveBeenLastCalledWith({
    //     clusterId: secondCluster.clusterId,
    //     parentProductName: secondCluster.subdomain.DTinsight[0],
    //     limit: 0,
    //     mode: 1,
    //     productName: undefined,
    //     productVersion: undefined,
    //     componentList: []
    // });
  });

  // 版本号搜索
  test('search input', () => {
    const searchInput = wrapper.getByPlaceholderText('按组件版本号搜索');
    const icon = searchInput.nextSibling.firstChild;
    expect(searchInput).toBeInTheDocument();
    // 搜索
    fireEvent.change(searchInput, { target: { value: '_beta' } });
    fireEvent.click(icon);
    expect(Service.getAllProducts).toHaveBeenCalled();
  });
});
