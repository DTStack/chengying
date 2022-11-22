import * as React from 'react';
import ProgressBar from '../index';
import { render, RenderResult } from '@testing-library/react';

const defaultProps = {
  unit: 'GB',
  now: 5,
  total: 10,
};
describe('progressbar component test', () => {
  let wrapper: RenderResult;
  let text;
  let progress;
  beforeEach(() => {
    wrapper = render(<ProgressBar {...defaultProps} />);
    text = wrapper.container.firstChild;
    progress = wrapper.container.getElementsByClassName('ant-progress-bg')[0];
  });

  // unit
  test('unit exist render', () => {
    const { now, total, unit } = defaultProps;
    const percent = total ? (now / total) * 100 : 0;
    expect(text.textContent).toBe(`${now}${unit} / ${total}${unit}`);
    expect(progress.getAttribute('style')).toMatch(`width: ${percent}%`);
  });

  // percent
  test('percent exist render', () => {
    const newProps = {
      percent: 50,
      now: '5GB',
      total: '11GB',
    };
    const { now, total, percent } = newProps;
    // rerender
    wrapper.rerender(<ProgressBar {...newProps} />);
    // test
    expect(text.textContent).toBe(`${now} / ${total}`);
    expect(progress.getAttribute('style')).toMatch(`width: ${percent}%`);
  });
});
