import * as React from 'react';
import Pagination from '../index';
import { fireEvent, render, RenderResult } from '@testing-library/react';

const defaultProps = {
  handleClickTop: jest.fn(),
  handleClickUp: jest.fn(),
  handleClickDown: jest.fn(),
  handleClickNew: jest.fn(),
};
describe('pagination unit test', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<Pagination {...defaultProps} />);
  });

  test('actions test', () => {
    const buttons = wrapper.container.getElementsByTagName('button');
    fireEvent.click(buttons[0]);
    expect(defaultProps.handleClickTop).toHaveBeenCalled();
    expect(defaultProps.handleClickTop).toHaveBeenCalledTimes(1);
    fireEvent.click(buttons[1]);
    expect(defaultProps.handleClickUp).toHaveBeenCalled();
    fireEvent.click(buttons[2]);
    expect(defaultProps.handleClickDown).toHaveBeenCalled();
    fireEvent.click(buttons[3]);
    expect(defaultProps.handleClickNew).toHaveBeenCalled();
  });
});
