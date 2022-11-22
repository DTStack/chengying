import * as React from 'react';
import { shallow, mount } from 'enzyme';
// import toJson from 'enzyme-to-json';
import { AlertRulePage } from '../rule';
import HeaderStore from '@/stores/headerReducer';
import * as HeaderActions from '@/actions/headerAction';
import { MemoryRouter } from 'react-router-dom'; // 避免组件中使用Link的情况
import { AlertRuleProps } from '../__mocks__/ruleList.mock';
// import { Alert } from 'antd';

const setup = () => {
  const props = AlertRuleProps;
  const sRuleList = shallow(
    <MemoryRouter>
      <AlertRulePage {...props} />
    </MemoryRouter>
  );
  const mRuleList = mount(
    <MemoryRouter>
      <AlertRulePage {...props} />
    </MemoryRouter>
  );
  return {
    props,
    sRuleList,
    mRuleList,
  };
};
/**
 * 测试header组件的reducer
 */
describe('header nav reducer', () => {
  it('should handle initial state', () => {
    expect(HeaderStore(undefined, { type: null })).toEqual({
      cur_product: {
        product_id: -1,
        product_name: '选择产品',
      },
      products: [],
      cur_parent_product: '选择产品', // 当前产品
      parentProducts: [],
      cur_parent_cluster: {
        id: -1,
        name: '选择集群',
        type: 'hosts',
        mode: 0,
      },
      parentClusters: [],
    });
  });

  it('get product list', async () => {
    await HeaderActions.getParentProductList();
    // console.log('parent product list:', res);
  });
});

/**
 * 测试告警规则列表是否显示正常
 */
describe('rule list', () => {
  const { props, mRuleList } = setup();
  // console.log(AlertRulePage.prototype);
  const updateAlertRuleList = jest.spyOn(
    AlertRulePage.prototype,
    'updateAlertRuleList'
  );
  const tree = shallow(<AlertRulePage {...props} />);
  // test snapshot
  // test('get a snapshot', () => {
  //     expect(tree).toMatchSnapshot();
  // });
  // 测试搜索和状态筛选是否正确渲染
  test('rule page render correct', () => {
    expect(mRuleList.find('.ant-input-search').exists()).toBeTruthy();
    expect(mRuleList.find('.ant-select').exists()).toBe(true);
  });

  // 测试条件搜索功能
  test('search can work well', () => {
    tree.setState({
      query: '',
    });
    const input = tree.find('Search');
    input.simulate('search', 'aaa');
    expect(tree.state('query')).toBe('aaa');
    expect(updateAlertRuleList).toHaveBeenCalled();
  });
  // 测试状态筛选功能
  test('state filter can work well', () => {
    tree.setState({
      state: '',
    });
    const select = tree.find('Select').last();
    select.simulate('change', 'ok');
    expect(tree.state('state')).toBe('ok');
    expect(updateAlertRuleList).toHaveBeenCalled();
  });
});
