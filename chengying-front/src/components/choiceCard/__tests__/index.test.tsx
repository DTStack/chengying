import * as React from 'react';
import ChoiceCard from '../index';
import { fireEvent, render, RenderResult } from '@testing-library/react';

const defaultProps = {
  title: 'title',
  content: 'content',
  imgSrc: 'require("public/imgs/cluster_hosts.png")',
};

describe('card unit test', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<ChoiceCard {...defaultProps} />);
  });

  test('card render', () => {
    const p = wrapper.container.getElementsByTagName('p');
    expect(p[0].textContent).toBe(defaultProps.title);
    expect(p[1].textContent).toBe(defaultProps.content);
  });

  // rerender
  test('test not required params', () => {
    const params = Object.assign({}, defaultProps, {
      title: 'test',
      className: 'cluster-box-style',
      handleTypeClick: jest.fn(),
      onFocus: jest.fn(),
      onBlur: jest.fn(),
    });
    wrapper.rerender(<ChoiceCard {...params} />);
    const card =
      wrapper.container.getElementsByClassName('c-card__container')[0];
    // class
    expect(card.className).toMatch(params.className);
    // event
    fireEvent.click(card);
    expect(params.handleTypeClick).toHaveBeenCalledTimes(1);
    fireEvent.focus(card);
    expect(params.onFocus).toHaveBeenCalledTimes(1);
  });
});
