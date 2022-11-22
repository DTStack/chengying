import * as React from 'react';
import {
  cleanup,
  fireEvent,
  render,
  RenderResult,
} from '@testing-library/react';
import SelfInfo from '../index';
import '@testing-library/jest-dom/extend-expect';
import { userCenterService } from '@/services';
import { UserCenterMock } from '@/mocks';

jest.mock('@/services');

describe('self info render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (userCenterService.getLoginedUserInfo as any).mockResolvedValue({
      data: UserCenterMock.getLoginedUserInfo,
    });
    wrapper = render(<SelfInfo />);
  });

  afterEach(cleanup);

  // 编辑信息
  test('edit test', async () => {
    const editBtn = wrapper.getByText('编辑资料').parentElement;

    // 编辑前
    const fullName = wrapper.getByLabelText('姓名');
    expect(fullName.getAttribute('disabled')).not.toBeNull();

    // 编辑后
    fireEvent.click(editBtn);
    const cancelBtn = wrapper.getByText('取 消').parentElement;
    const saveBtn = wrapper.getByText('保 存').parentElement;
    expect(cancelBtn).toBeInTheDocument();
    expect(fullName.getAttribute('disabled')).toBeNull();

    // 取消
    fireEvent.click(cancelBtn);
    expect(fullName.getAttribute('disabled')).not.toBeNull();

    // 保存
    fireEvent.click(editBtn);
    const userInfo = UserCenterMock.getLoginedUserInfo.data;
    const params = {
      fullName: 'test11',
      email: userInfo.email,
      phone: userInfo.phone,
      company: userInfo.company,
    };
    fireEvent.change(fullName, { target: { value: params.fullName } });
    (userCenterService.motifyUserInfo as any).mockResolvedValue({
      data: UserCenterMock.motifyUserInfo,
    });
    (userCenterService.getLoginedUserInfo as any).mockResolvedValue({
      data: UserCenterMock.getLoginedUserInfo,
    });
    fireEvent.click(saveBtn);
    expect(userCenterService.motifyUserInfo).toHaveBeenCalledWith(params);
  });
});
