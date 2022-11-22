import * as React from 'react';
import Command from '../index';
import { render, RenderResult, cleanup } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';

const defaultProps = {
  history: {},
  location: {},
  match: {},
};

/**
 * 测试header组件的reducer
 */
describe('command unit test', () => {
  let wrapper: RenderResult;

  beforeEach(() => {
    wrapper = render(<Command {...defaultProps} />);
  });

  afterEach(cleanup);

  test('command page render search', () => {
    const searchBox = wrapper.getByPlaceholderText('请输入搜索关键词');
    expect(searchBox).toBeInTheDocument();
  });

  test('command list render', () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    // columns render
    for (const i in columns) {
      const titles = ['', '对象', '操作', '状态', '开始时间', '持续时间'];
      const column = columns[i];
      if (column instanceof HTMLElement) {
        expect(column.textContent).toBe(titles[i]);
      }
    }
  });

  // rerender
  test('command snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });
});
