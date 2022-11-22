import * as React from 'react';
import ClusterLayout from '../index';
import reducer from '@/stores';
import { renderWithRedux } from '@/utils/test';
import { BrowserRouter } from 'react-router-dom';
import { clusterNavData } from '@/constants/navData';
import { fireEvent, RenderResult, cleanup } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';

const defaultProps = {
  HeaderStore: {
    cur_product: {
      product_id: -1,
      product_name: '选择产品',
    },
    products: [],
    cur_parent_product: '选择产品', // 当前产品
    parentProducts: [],
    cur_parent_cluster: {
      id: 1,
      mode: 0,
      name: 'dtstack',
      type: 'hosts',
    },
    parentClusters: [],
  },
  location: {},
  history: {},
  match: {},
};

describe('test clusterlayout navbar', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = renderWithRedux(
      <BrowserRouter>
        <ClusterLayout {...defaultProps} />
      </BrowserRouter>,
      reducer,
      { HeaderStore: defaultProps.HeaderStore }
    );
  });

  afterEach(cleanup);

  test('test tabs render', () => {
    const tabsNav = wrapper.container.getElementsByClassName('ant-tabs-nav')[0];
    const tabs = tabsNav.getElementsByClassName('ant-tabs-tab');
    const tabpanes =
      wrapper.container.getElementsByClassName('ant-tabs-tabpane');
    const clusterTypeMaps = {
      hosts: {
        0: ['imagestore', 'namespace'],
      },
      kubernetes: {
        0: ['namespace', 'patchHistory', 'echoList'], // 自建
        1: ['index', 'host', 'patchHistory', 'echoList'], // 导入
      },
    };
    const {
      HeaderStore: {
        cur_parent_cluster: { mode, type },
      },
    } = defaultProps;
    const clusterMaps = clusterTypeMaps[type];
    // 过滤条件
    const filterConditions = clusterMaps[mode];
    // 过滤不同集群类型菜单
    const realClusterData = clusterNavData.filter(
      (nav) => nav.isShow && !filterConditions.includes(nav.key)
    );
    for (const i in tabs) {
      if (tabs[i] instanceof HTMLElement) {
        expect(tabs[i].textContent).toBe(realClusterData[i].title);
      }
    }
    expect(tabs.length).toBe(tabpanes.length);
  });

  test('test cluster render', () => {
    const a = wrapper.container.getElementsByClassName(
      'ant-dropdown-trigger'
    )[0];
    fireEvent.click(a);
    expect(a).toBeInTheDocument();
    const title = a.getElementsByClassName('title')[0];
    expect(title.textContent).toBe('dtstack');
  });

  test('test snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });
});
