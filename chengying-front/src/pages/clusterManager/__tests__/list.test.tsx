import * as React from 'react';
import { waitFor, cleanup } from '@testing-library/react';
import List from '../list';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { clusterManagerService } from '@/services';
import { clusterMnagerMock } from '@/mocks';

afterEach(cleanup);

jest.mock('@/services');

const defaultProps = {
  history: null,
  HeaderStore: null,
  actions: null,
  authorityList: [],
};

describe('Cluster Manage List', () => {
  let cluster_container;
  let cards;
  let hosts: HTMLElement;
  let kubernetes: HTMLElement;
  let kubernetes_self: HTMLElement;
  beforeEach(async () => {
    (clusterManagerService.getClusterLists as any).mockResolvedValue({
      data: clusterMnagerMock.getClusterList,
    });
    const { container } = renderWithRedux(
      <List {...defaultProps} />,
      reducer,
      {}
    );
    // 等待列表渲染，获取集群
    const cluster_container_array = await waitFor(() =>
      container.getElementsByClassName('cluster-list-page')
    );
    cluster_container = cluster_container_array[0];
    // 集群显示信息校验
    cards = cluster_container.getElementsByClassName('cluster-page-card');
    hosts = cards[0];
    kubernetes = cards[1];
    kubernetes_self = cards[2];
  });

  // 每行3个
  test('each row card can is at most 3', () => {
    const row = cluster_container.getElementsByClassName('ant-row')[0];
    expect(row.children.length).toBeLessThanOrEqual(3);
  });

  // 状态校验
  test('status validate', () => {
    const statusMap = {
      Running: 'status-running',
      Waiting: 'status-waiting',
      Error: 'status-error',
    };
    for (const i in cards) {
      if (cards[i] && cards[i] instanceof HTMLElement) {
        const status = cards[i].getElementsByClassName(
          'card-extra-status'
        )[0] as HTMLDivElement;
        const status_class = status.className;
        const status_text = status.textContent;
        expect(status_class).toBe(
          'card-extra-status ' + statusMap[status_text]
        );
      }
    }
  });

  // 集群模式
  test('cluster mode', () => {
    expect(hosts.getElementsByTagName('p')[0].textContent).toMatch(/主机集群/);
    expect(kubernetes.getElementsByTagName('p')[0].textContent).toMatch(
      /Kubernetes集群（导入已有集群）/
    );
    expect(kubernetes_self.getElementsByTagName('p')[0].textContent).toMatch(
      /Kubernetes集群（自建集群）/
    );
  });

  // 主机集群 - 磁盘， kubernetes - PODS
  test('cluster render info validate', () => {
    expect(
      hosts.getElementsByClassName('data-view-item')[2].firstChild.nextSibling
        .textContent
    ).toBe('磁盘');
    expect(
      kubernetes.getElementsByClassName('data-view-item')[2].firstChild
        .nextSibling.textContent
    ).toBe('PODS');
  });
});
