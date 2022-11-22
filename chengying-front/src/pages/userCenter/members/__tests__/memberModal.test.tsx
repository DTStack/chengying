import * as React from 'react';
import {
  cleanup,
  fireEvent,
  RenderResult,
  render,
} from '@testing-library/react';
import MemberModal from '../memberModal';
import { userCenterService } from '@/services';
import '@testing-library/jest-dom/extend-expect';

jest.mock('@/services');

const defaultProps = {
  visible: true,
  onCancel: jest.fn(),
  memberInfo: null,
};

const userInfo = {
  username: 'test',
  fullName: '',
  email: 'test@dtstack.com',
  phone: '',
  company: '',
  roleId: 2,
  userId: undefined,
};
describe('member modal unit test', () => {
  let wrapper: RenderResult;
  let container: HTMLElement; // 外层
  let okBtn: HTMLButtonElement; // 确定按钮
  let cancelBtn: HTMLButtonElement;
  beforeEach(() => {
    wrapper = render(<MemberModal {...defaultProps} />);
    okBtn = wrapper.getByText('确 认').parentElement as HTMLButtonElement;
    cancelBtn = wrapper.getByText('取 消').parentElement as HTMLButtonElement;
    container = okBtn.parentElement.parentElement;
  });

  afterEach(cleanup);

  // 创建流程校验
  test('create member', async () => {
    // 步骤
    const step = container.getElementsByClassName('ant-steps')[0];
    expect(step).toBeInTheDocument();
    const steps = step.getElementsByClassName('ant-steps-item-title');
    expect(steps[0].textContent).toBe('编辑信息');
    expect(steps[1].textContent).toBe('生成账号');

    /* form 表单 */
    const username = wrapper.getByLabelText('账号');
    const email = wrapper.getByLabelText('邮箱');
    const roles = container.querySelectorAll('input[type=radio]');
    // 第一步
    // 输入表单必填信息
    fireEvent.change(username, { target: { value: userInfo.username } });
    fireEvent.change(email, { target: { value: userInfo.email } });
    fireEvent.click(roles[0]);
    // 点击确认 - 下一步
    const reigstAPI = (userCenterService.regist as any).mockResolvedValue({
      data: {
        code: 0,
        data: { username: 'test', password: '123456' },
        msg: 'ok',
      },
    });
    fireEvent.click(okBtn);
    expect(userCenterService.regist).toHaveBeenCalled();
    expect(userCenterService.regist).toHaveBeenCalledTimes(1);
    expect(userCenterService.regist).toHaveBeenCalledWith(userInfo);
    await reigstAPI();
    // 第二步
    const userInfoContent =
      container.getElementsByClassName('new-user-info')[0];
    expect(userInfoContent).toBeInTheDocument();
    expect(okBtn.textContent).toBe('确认复制');
    fireEvent.click(cancelBtn);
  });

  // 编辑流程校验
  test('edit member', () => {
    const rerenderProps = {
      ...defaultProps,
      memberInfo: {
        ...userInfo,
        role_id: userInfo.roleId,
        full_name: userInfo.fullName,
      },
    };
    wrapper.rerender(<MemberModal {...rerenderProps} />);

    // 没有步骤流程
    const step = container.getElementsByClassName('ant-steps')[0];
    expect(step).toBeUndefined();
    // 表单
    const username = wrapper.getByLabelText('账号') as HTMLInputElement;
    const email = wrapper.getByLabelText('邮箱') as HTMLInputElement;
    // 账号不可编辑
    expect(username.getAttribute('disabled')).not.toBeNull();
    expect(username.value).toBe(userInfo.username);
    expect(email.value).toBe(userInfo.email);

    // 修改
    const newEmail = 'test11@dtstack.com';
    fireEvent.change(email, { target: { value: newEmail } });
    (userCenterService.modifyInfoByAdmin as any).mockResolvedValue({
      data: { code: 0, data: true, msg: 'ok' },
    });
    fireEvent.click(okBtn);
    expect(userCenterService.modifyInfoByAdmin).toHaveBeenCalledWith({
      ...userInfo,
      email: newEmail,
    });
  });
});
