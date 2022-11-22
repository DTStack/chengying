import * as React from 'react';
import Error from '@/components/error/error';
import { render, RenderResult, cleanup } from '@testing-library/react';

describe('error unit test', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<Error />);
  });

  afterEach(cleanup);

  test('should Snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });
});
