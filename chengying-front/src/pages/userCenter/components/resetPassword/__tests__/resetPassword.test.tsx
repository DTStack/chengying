import * as React from 'react';
import {
  cleanup,
  fireEvent,
  RenderResult,
  render,
} from '@testing-library/react';
import ResetPwdModal from '../index';
import '@testing-library/jest-dom/extend-expect';

jest.mock('@/services');

const defaultProps = {
  userInfo: null,
  visible: true,
  onCancel: jest.fn(),
  onSubmit: jest.fn(),
};
const pwdValue = 'abc123';
describe('reset pwd unit test', () => {
  let wrapper: RenderResult;
  let okBtn: HTMLButtonElement;
  let cancelBtn: HTMLButtonElement;
  beforeEach(() => {
    wrapper = render(<ResetPwdModal {...defaultProps} />);
    okBtn = wrapper.getByText('OK').parentElement as HTMLButtonElement;
    cancelBtn = wrapper.getByText('Cancel').parentElement as HTMLButtonElement;
  });

  afterEach(cleanup);

  // 用户修改个人密码
  test('user reset own pwd', () => {
    // 没有账号信息，但需要输入旧密码
    const oldPwd = wrapper.getByLabelText('旧密码');
    expect(oldPwd).toBeInTheDocument();

    const newPwd = wrapper.getByLabelText('新密码');
    const confirmPwd = wrapper.getByLabelText('确认新密码');

    /* 重置密码 */
    fireEvent.change(oldPwd, { target: { value: pwdValue } });
    fireEvent.change(newPwd, { target: { value: pwdValue } });
    fireEvent.change(confirmPwd, { target: { value: pwdValue } });
    fireEvent.click(okBtn);
    expect(defaultProps.onSubmit).toHaveBeenCalled();
    expect(defaultProps.onSubmit).toHaveBeenCalledTimes(1);
  });

  // 管理员权限重置其他用户密码
  test('admin reset user pwd', () => {
    const userInfo = {
      username: 'yuwan',
    };
    wrapper.rerender(<ResetPwdModal {...defaultProps} userInfo={userInfo} />);

    // 有账号信息
    const unameLabel = wrapper.getByText('账号');
    const uname = unameLabel.parentElement.nextElementSibling;
    expect(uname.textContent).toBe(userInfo.username);

    // 重置密码
    const newPwd = wrapper.getByLabelText('新密码');
    const confirmPwd = wrapper.getByLabelText('确认新密码');
    fireEvent.change(newPwd, { target: { value: pwdValue } });
    fireEvent.change(confirmPwd, { target: { value: pwdValue } });
    fireEvent.click(okBtn);
    expect(defaultProps.onSubmit).toHaveBeenCalled();
    expect(defaultProps.onSubmit).toHaveBeenCalledTimes(2);
  });

  test('cancel', () => {
    fireEvent.click(cancelBtn);
    expect(defaultProps.onCancel).toHaveBeenCalled();
  });
});
