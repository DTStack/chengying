/*
 * @Description: This is desc
 * @Author: wulin
 * @Date: 2021-05-10 11:26:40
 * @LastEditors: wulin
 * @LastEditTime: 2021-06-24 14:34:41
 */
import * as React from 'react';
import ImageStore from '../imageStore';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { cur_parent_cluster, authorityList, ImageStoreMock } from '@/mocks';
import { fireEvent, RenderResult, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import { imageStoreService } from '@/services';

jest.mock('@/services');

const defaultProps = {
  cur_parent_cluster,
  authorityList,
};
describe('iamge store render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (imageStoreService.getImageStoreList as any).mockResolvedValue({
      data: ImageStoreMock.getImageStoreList,
    });
    wrapper = renderWithRedux(<ImageStore {...defaultProps} />, reducer, {
      HeaderStore: {
        cur_parent_cluster,
      },
      UserCenterStore: {
        authorityList,
      },
    });
  });
  // 表格列渲染
  test('table render', () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    const columnNames = ['', '仓库名称', '仓库别名', '仓库地址', '操作', ''];
    for (const i in columns) {
      if (columns[i] instanceof HTMLElement) {
        expect(columns[i].textContent).toBe(columnNames[i]);
      }
    }
  });
  // 添加仓库
  test('add store', () => {
    const addBtn = wrapper.getByTestId('add-btn');
    expect(addBtn.textContent).toBe('添加仓库');
  });
  // 编辑
  test('edit store', async () => {
    const tableRow = wrapper.container.getElementsByClassName('ant-table-row');
    const firstRow = tableRow[0];
    const action = firstRow.children.item(4).firstChild;
    // 点击编辑
    const getInfoApi = (
      imageStoreService.getImageStoreInfo as any
    ).mockResolvedValue({ data: ImageStoreMock.getImageStoreInfo });
    fireEvent.click(action);
    // 校验模态框
    await getInfoApi();
    const editModal = screen.getByText('编辑仓库');
    const inputName = screen.getByLabelText('仓库名称') as HTMLInputElement;
    expect(editModal).toBeInTheDocument();
    expect(inputName.value).toBe('yuwan_test');
  });

  // 设置默认仓库
  test('set default store', () => {
    const setDefaultBtns =
      wrapper.container.getElementsByClassName('set-default-box');
    const firstDefaultBtn = setDefaultBtns[0].firstChild;
    const secondDefaultBtn = setDefaultBtns[1].firstChild;
    expect(firstDefaultBtn.nodeName).toBe('SPAN');
    expect(firstDefaultBtn.textContent).toBe('默认仓库');
    expect(secondDefaultBtn.nodeName).toBe('BUTTON');
    expect(secondDefaultBtn.textContent).toBe('设为默认仓库');

    // 设置
    (imageStoreService.setDefaultStore as any).mockResolvedValue({
      data: ImageStoreMock.setDefaultStore,
    });
    fireEvent.click(secondDefaultBtn);
    expect(imageStoreService.setDefaultStore).toHaveBeenCalledTimes(1);
  });
});
