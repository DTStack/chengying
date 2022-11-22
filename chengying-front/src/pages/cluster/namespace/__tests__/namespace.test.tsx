import * as React from 'react';
import { renderWithRedux } from '@/utils/test';
import Namespace from '../namespace';
import { ClusterNamespaceService, imageStoreService } from '@/services';
import { cleanup, RenderResult } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import reducer from '@/stores';
import {
  authorityList,
  cur_parent_cluster,
  ClusterNamespaceMock,
  ImageStoreMock,
} from '@/mocks';

jest.mock('@/services');

const defaultProps = {
  cur_parent_cluster,
  authorityList,
};
describe('namespace unit test', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (ClusterNamespaceService.getNamespaceList as any).mockResolvedValue({
      data: ClusterNamespaceMock.getNamespaceList,
    });
    (imageStoreService.getImageStoreList as any).mockResolvedValue({
      data: ImageStoreMock.getImageStoreList,
    });
    wrapper = renderWithRedux(<Namespace {...defaultProps} />, reducer, {
      HeaderStore: {
        cur_parent_cluster,
      },
      UserCenterStore: {
        authorityList,
      },
    });
  });

  afterEach(cleanup);

  // 列表渲染
  test('table render', () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    const colTitles = [
      'namespace',
      '状态',
      '添加方式',
      'cpu使用',
      '内存使用',
      '最近修改人',
      '最近修改时间',
      '操作',
    ];
    for (const i in columns) {
      const column = columns[i];
      if (column instanceof HTMLElement) {
        expect(column.textContent).toBe(colTitles[i]);
      }
    }
  });
});
