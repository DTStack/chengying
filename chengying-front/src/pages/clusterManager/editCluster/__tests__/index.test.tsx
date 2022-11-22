import * as React from 'react';
import {
  fireEvent,
  RenderResult,
  waitForElementToBeRemoved,
} from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import { Router, matchPath } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import EditCluster, { linkMap } from '../index';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { clusterManagerService, Service } from '@/services';
import { clusterMnagerMock, ServiceMock } from '@/mocks';
import { clusterInfo } from '@/mocks/clusterMnagerMock';

jest.mock('@/services');

// 测试步骤
describe('edit cluster step', () => {
  let wrapper: RenderResult;
  beforeEach(async () => {
    const { id, type, mode } = clusterInfo;
    const history = createMemoryHistory();
    const path = '/clustermanage/create/edit';
    const url = `${path}?id=${id}&type=${type}&mode=${mode}`;
    history.push(url);
    const defaultProps = {
      location: history.location,
      history: history,
      match: matchPath(path, {
        path: '/clustermanage/create/:type?',
        exact: false,
        strict: false,
      }),
      actions: null,
      clusterInfo,
    };
    (clusterManagerService.getClusterInfo as any).mockResolvedValue({
      data: clusterMnagerMock.getClusterInfo,
    });
    (clusterManagerService.getKubernetesAvaliable as any).mockResolvedValue({
      data: clusterMnagerMock.getKubernetesAvaliable,
    });
    wrapper = renderWithRedux(
      <Router history={history}>
        <EditCluster {...defaultProps} />
      </Router>,
      reducer,
      { editClusterStore: { clusterInfo } }
    );
  });

  // 检查label名
  test('check label', () => {
    const { id, type, mode } = clusterInfo;
    const navLink = wrapper.container.getElementsByClassName('nav-link')[0];
    const label = `${id !== -1 ? '编辑' : '添加'}集群（${
      linkMap[type][mode]
    }）`;
    expect(navLink.textContent).toBe(label);
  });

  // 步骤
  test('check step render', async () => {
    // 步骤按钮
    const stepBtns = wrapper.container
      .getElementsByClassName('edit-cluster-bottom')[0]
      .getElementsByClassName('ant-btn');

    // stepOne 渲染
    const stepOne = wrapper.getByTestId('ec-step-one');
    expect(stepOne).toBeInTheDocument();
    // 步骤按钮校验
    expect(stepBtns.length).toEqual(2);
    expect(stepBtns[0].textContent).toBe('取 消');
    expect(stepBtns[1].textContent).toBe('下一步');

    // 点击下一步
    (clusterManagerService.clusterSubmitOperate as any).mockResolvedValue({
      data: clusterMnagerMock.clusterSubmitOperate,
    });
    (Service.getClusterHostList as any).mockResolvedValue({
      data: ServiceMock.getClusterHostList,
    });
    (Service.getClusterhostGroupLists as any).mockResolvedValue({
      data: ServiceMock.getClusterhostGroupLists,
    });
    fireEvent.click(stepBtns[1]);
    await waitForElementToBeRemoved(wrapper.getByTestId('ec-step-one'));

    // stepFinal 渲染
    const stepFinal = wrapper.getByTestId('ec-step-final');
    expect(stepFinal).toBeInTheDocument();
    // 步骤按钮校验
    expect(stepBtns.length).toEqual(1);
    expect(stepBtns[0].textContent).toBe('完 成');
  });
});
