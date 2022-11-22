import * as React from 'react';
import TimerInc from '../index';
import { RenderResult, cleanup, render } from '@testing-library/react';

describe('timerInc render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<TimerInc initial={1} />);
  });

  afterEach(cleanup);

  test('timerInc render', () => {
    const time = wrapper.container.getElementsByTagName('div')[0];
    expect(time.textContent).toBe('1s');
  });

  test('timerInc snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });
});
