import * as React from 'react';
import Header from '@/components/headerNamespace';
import { cur_parent_cluster, clusterMnagerMock } from '@/mocks';
import {
  fireEvent,
  render,
  RenderResult,
  waitFor,
} from '@testing-library/react';

const parentClusters = clusterMnagerMock.getClusterList.data.clusters;
const defaultProps = {
  HeaderStore: {
    parentClusters,
    cur_parent_product: 'DTinsight',
    cur_parent_cluster,
  },
  location: window.location,
  handleSwitchProduct: jest.fn(),
  getRedService: jest.fn(),
  actions: {
    setCurrentParentProduct: jest.fn(),
    getClusterProductList: jest.fn(),
  },
};

describe('header unit test', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<Header {...defaultProps} />);
  });

  test('should Snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });

  // 集群 | 产品 切换
  test('dropdown menu select', async () => {
    // 当前选择项
    const dropKey = wrapper.container.getElementsByClassName(
      'header-cascader-title'
    )[0];
    expect(dropKey.textContent).toBe('DTLogger/DTinsight');
    /* 下拉菜单切换 */
    // click 渲染下拉框
    // fireEvent.click(dropKey);
    // const dropItem = await waitFor(() => document.body.getElementsByClassName('ant-cascader-menu-item'))[0];
    // // 点击切换
    // fireEvent.click(dropItem);
    // expect(defaultProps.handleSwitchProduct).toHaveBeenCalledTimes(1);
  });
});
