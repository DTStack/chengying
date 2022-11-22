import * as React from 'react';
import {
  render,
  cleanup,
  fireEvent,
  RenderResult,
} from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import StepOne from '../stepOne';
import { clusterInfo } from '@/mocks/clusterMnagerMock';

afterEach(cleanup);

const defaultProps = {
  isEdit: true,
  clusterInfo: clusterInfo,
  location: window.location,
  action: null,
};

// 基本信息字段规则校验
describe('base params rule validate', () => {
  let input_tags;
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<StepOne {...defaultProps} />);
    input_tags = wrapper.container.getElementsByClassName(
      'ant-select-selection__choice'
    );
  });

  // 集群名称
  test('cluster name', () => {
    const input_name = wrapper.getByLabelText('集群名称');
    expect(input_name).toBeInTheDocument();
    expect(input_name.getAttribute('placeholder')).toBe('请输入集群名称');
    expect((input_name as HTMLInputElement).value).toMatch(
      /^[A-Za-z0-9_]{1,64}$/
    );
  });

  // 集群描述
  test('cluster desc', () => {
    const input_desc = wrapper.getByLabelText('集群描述');
    expect(input_desc).toBeInTheDocument();
    expect(input_desc.getAttribute('placeholder')).toBe('请输入集群描述');
    expect(
      (input_desc as HTMLTextAreaElement).value.length
    ).toBeLessThanOrEqual(200);
  });

  // 集群标签
  test('cluster tag', () => {
    const tags = defaultProps.clusterInfo.tags.split(',');
    for (const i in input_tags) {
      if (input_tags[i] instanceof HTMLElement) {
        const tag = tags[i];
        expect(tag).toMatch(/^\S{1,32}$/);
        expect(input_tags[i].textContent).toBe(tag);
      }
    }
    expect(input_tags.length).toEqual(tags.length);
    expect(input_tags.length).toBeLessThanOrEqual(5);
  });
});

// k8s信息
describe('kubernetes params value validate', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<StepOne {...defaultProps} />);
  });

  // 模式切换
  test('mode change', () => {
    const btn_change = wrapper.getByTestId('btn-mode-change');
    const text_before = btn_change.textContent;
    expect(text_before).toBe('编辑YMAL');

    fireEvent.click(btn_change);
    const text_after = btn_change.textContent;
    expect(text_after).toBe('表单输入');
  });
});
