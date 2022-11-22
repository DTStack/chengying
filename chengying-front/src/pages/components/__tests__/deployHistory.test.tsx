import * as React from 'react';
import DeployHistory from '../deployHistory';
import { renderWithRedux } from '@/utils/test';
import reducer from '@/stores';
import { Service, servicePageService, deployService } from '@/services';
import {
  authorityList,
  DeployMock,
  ServiceMock,
  ServicePageMock,
} from '@/mocks';
import {
  fireEvent,
  RenderResult,
  waitForElementToBeRemoved,
} from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';

jest.mock('@/services');

const searchParams = {
  clusterId: '46',
  parentProductName: 'DTEM', // 产品名称
  productName: [], // 组件名称
  productVersion: undefined,
  shouldNameSpaceShow: false,
};

describe('deploy history list render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (Service.getProductUpdateRecords as any).mockResolvedValue({
      data: ServiceMock.getProductUpdateRecords,
    });
    wrapper = renderWithRedux(<DeployHistory {...searchParams} />, reducer, {
      UserCenterStore: { authorityList },
    });
  });

  // 表格列
  test('table columns', async () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    const columnKeys = [
      '组件名称',
      '组件版本号',
      '安装包类型',
      '部署人',
      '部署时间',
      '状态',
      '部署快照',
      '查看',
    ];
    for (const i in columns) {
      if (columns[i] instanceof HTMLElement) {
        expect(columns[i].textContent).toBe(columnKeys[i]);
      }
    }
    // 取表格第一行数据
    const rows = wrapper.container.getElementsByClassName('ant-table-tbody')[0];
    const row = rows.firstChild as HTMLElement;
    const rowInfo = ServiceMock.getProductUpdateRecords.data.list[0];

    /* 状态 */
    const statusCol = row.children.item(5);
    const status = statusCol.firstChild.firstChild as HTMLSpanElement;
    const statusClass = status.className;
    const statusIcon = (status.firstChild as HTMLElement).className;
    const statusText = status.textContent;
    switch (rowInfo.status) {
      case 'undeployed':
        expect(statusClass).toBe('deploy-status-green');
        expect(statusIcon).toMatch(/(check-circle)/);
        expect(statusText).toBe('卸载成功');
        break;
      case 'deploying':
        expect(statusClass).toBe('deploy-status-orange');
        expect(statusIcon).toMatch(/(exclamation-circle)/);
        expect(statusText).toBe('部署中');
        break;
      case 'deployed':
        expect(statusClass).toBe('deploy-status-green');
        expect(statusIcon).toMatch(/(check-circle)/);
        expect(statusText).toBe('部署成功');
        break;
      case 'deploy fail':
        expect(statusClass).toBe('deploy-status-red');
        expect(statusIcon).toMatch(/(close-circle)/);
        expect(statusText).toBe('部署失败');
        break;
      case 'undeploying':
        expect(statusClass).toBe('deploy-status-orange');
        expect(statusIcon).toMatch(/(exclamation-circle)/);
        expect(statusText).toBe('卸载中');
        break;
      case 'undeploy fail':
        expect(statusClass).toBe('deploy-status-red');
        expect(statusIcon).toMatch(/(close-circle)/);
        expect(statusText).toBe('卸载失败');
        break;
    }

    /* 部署快照 */
    const actionCol = row.children.item(6);
    const action = actionCol.firstChild;
    // 点击
    const deployShotApi = (Service.getDeployShot as any).mockResolvedValue({
      data: ServiceMock.getDeployShot,
    });
    fireEvent.click(action);
    // 部署快照弹框渲染
    await deployShotApi();
    const deployShotModal =
      wrapper.baseElement.getElementsByClassName('ant-modal')[0];
    const deployShotModalTitle =
      deployShotModal.getElementsByClassName('ant-modal-title')[0];
    const closeBtn =
      deployShotModal.getElementsByClassName('ant-modal-close')[0];
    expect(deployShotModalTitle.textContent).toBe('部署快照');
    // 关闭弹窗
    fireEvent.click(closeBtn);
    await waitForElementToBeRemoved(deployShotModal);

    /* 查看 */
    const viewCol = row.lastChild;
    const showLog = viewCol.lastChild;
    // 查看日志
    const groupApi = (
      servicePageService.getServiceGroup as any
    ).mockResolvedValue({ data: ServicePageMock.getServiceGroup });
    const logApi = (deployService.searchDeployLog as any).mockResolvedValue({
      data: DeployMock.searchDeployLog,
    });
    fireEvent.click(showLog);
    await groupApi();
    await logApi();
    const logTextModal =
      wrapper.baseElement.getElementsByClassName('ant-modal')[0];
    const logTextModalTitle =
      logTextModal.getElementsByClassName('ant-modal-title')[0];
    expect(logTextModalTitle.textContent).toBe('部署日志');
  });
});
