import * as React from 'react';
import CustomTransfer from '../index';
import {
  fireEvent,
  render,
  RenderResult,
  cleanup,
} from '@testing-library/react';

const defaultProps = {
  dataSource: [
    { id: '172.16.8.59', ip: '172.16.8.59', key: '172.16.8.59' },
    { id: '172.16.8.76', ip: '172.16.8.76', key: '172.16.8.76' },
    { id: '172.16.8.57', ip: '172.16.8.57', key: '172.16.8.57' },
  ],
  filterOption: () => null,
  footerLeft: (props) => null,
  footerRight: (props) => null,
  leftPageCurrent: 1, // 左边当前页码
  rightPageCurrent: 1, // 右边当前页码
  realTargetKeys: ['172.16.8.57'], // 全部已经选择的keys
  onChange: jest.fn(),
  listStyle: {},
  selectedKeys: [],
  showSearch: true,
  rowKey: () => null,
  targetKeys: ['172.16.8.57'],
  totalData: [
    { id: 2, ip: '172.16.8.59' },
    { id: 3, ip: '172.16.8.76' },
  ], // 所有数据（已选+未选）
  isRightPanelCheckoutAll: false,
  render: () => null,
};

describe('CustomTransfer unit test', () => {
  let wrapper: RenderResult;

  beforeEach(() => {
    wrapper = render(<CustomTransfer {...defaultProps} />);
  });

  afterEach(cleanup);

  test('CustomTransfer table render', () => {
    const left = wrapper.getByTestId('left');
    expect(left).not.toBe('10 / page');
    const tr = left.getElementsByTagName('tr');
    expect(tr.length).toBe(3);
    for (const i in tr) {
      if (tr[i] instanceof HTMLElement) {
        expect(tr[i].getAttribute('data-row-key')).toBe(
          defaultProps.dataSource[i].ip
        );
      }
    }
  });

  // rerender
  test('CustomTransfer snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });
});
