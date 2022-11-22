import * as React from 'react';
import EditModal from '../editModal';
import {
  fireEvent,
  render,
  RenderResult,
  screen,
} from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import { cur_parent_cluster, imageStoreInfo, ImageStoreMock } from '@/mocks';
import { imageStoreService } from '@/services';

jest.mock('@/services');

const defaultProps = {
  handleCancel: () => {},
  getImageStoreList: () => {},
  imageStoreInfo,
  clusterId: cur_parent_cluster.id,
  isEdit: !!Object.keys(imageStoreInfo).length,
};
describe('edit image store modal', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    wrapper = render(<EditModal {...defaultProps} />);
  });
  // 命名规则校验
  test('params rule test', () => {
    // 仓库名称
    const inputName = wrapper.getByLabelText('仓库名称') as HTMLInputElement;
    expect(inputName).toBeInTheDocument();
    expect(inputName.getAttribute('placeholder')).toBe('请输入仓库名称');
    expect(inputName.value).toMatch(/^\S{1,64}$/);
    // 仓库别名
    const inputAlias = wrapper.getByLabelText('仓库别名') as HTMLInputElement;
    expect(inputAlias.getAttribute('placeholder')).toBe('请输入仓库别名');
    expect(inputAlias.value).toMatch(/^[a-z0-9]([-a-z0-9]{0,62}[a-z0-9])?$/);
    // 仓库地址
    const inputAddr = wrapper.getByLabelText('仓库地址') as HTMLInputElement;
    expect(inputAddr.getAttribute('placeholder')).toBe('请输入仓库地址');
    expect(inputAddr.value).toMatch(/^[^\u4e00-\u9fa5\s]+$/);
    // 用户名
    const inputUserName = wrapper.getByLabelText('用户名') as HTMLInputElement;
    expect(inputUserName.getAttribute('placeholder')).toBe('请输入用户名');
    // 密码
    const inputPassword = wrapper.getByLabelText('密码') as HTMLInputElement;
    expect(inputPassword.getAttribute('placeholder')).toBe('请输入密码');
    expect(inputPassword.type).toBe('password');
    // 邮箱
    const inputEmail = wrapper.getByLabelText('邮箱') as HTMLInputElement;
    expect(inputEmail.getAttribute('placeholder')).toBe('请输入邮箱');
    expect(inputEmail.type).toBe('email');
  });
  // 模态框基础样式
  test('modal render', async () => {
    // title
    const title =
      wrapper.baseElement.getElementsByClassName('ant-modal-title')[0];
    expect(title.textContent).toBe(
      `${defaultProps.isEdit ? '编辑' : '添加'}仓库`
    );
    // 底部按钮
    const bottom =
      wrapper.baseElement.getElementsByClassName('ant-modal-footer')[0];
    const btns = bottom.getElementsByTagName('button');
    const submitBtn = btns[1];
    expect(btns.length).toEqual(2);

    // 提交
    let submitApi;
    if (defaultProps.isEdit) {
      submitApi = (imageStoreService.updateImageStore as any).mockResolvedValue(
        { data: ImageStoreMock.updateImageStore }
      );
    } else {
      submitApi = (imageStoreService.createImageStore as any).mockResolvedValue(
        { data: ImageStoreMock.createImageStore }
      );
    }
    fireEvent.click(submitBtn);
    if (defaultProps.isEdit) {
      expect(imageStoreService.updateImageStore).toHaveBeenCalledTimes(1);
      expect(imageStoreService.updateImageStore).toHaveBeenLastCalledWith({
        ...imageStoreInfo,
        clusterId: cur_parent_cluster.id,
      });
    } else {
      expect(imageStoreService.createImageStore).toHaveBeenCalledTimes(1);
    }
    // 提交成功
    await submitApi();
    const message = screen.queryByText('执行成功');
    expect(message).toBeInTheDocument();
  });
});
