import * as React from 'react';
import {
  cleanup,
  fireEvent,
  render,
  RenderResult,
} from '@testing-library/react';
import SecurityAudit from '../index';
import '@testing-library/jest-dom/extend-expect';
import { SecurityAuditService } from '@/services';
import { SecurityAuditMock } from '@/mocks';

jest.mock('@/services');

describe('security audit test', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (SecurityAuditService.getAuditModule as any).mockResolvedValue({
      data: SecurityAuditMock.getAuditModule,
    });
    (SecurityAuditService.getSafetyAudit as any).mockResolvedValue({
      data: SecurityAuditMock.getSafetyAudit,
    });
    wrapper = render(<SecurityAudit />);
  });

  afterEach(cleanup);

  // 操作模块选择
  test('operate modal change', () => {
    const operateSelect = wrapper.getByText('操作模块：').nextElementSibling;
    const modules = SecurityAuditMock.getAuditModule.data.list;
    fireEvent.click(operateSelect);
    // 点击过滤操作模块
    // 同时请求其下对应动作
    const option = wrapper.getByText(modules[0]);
    (SecurityAuditService.getAuditOperation as any).mockResolvedValue({
      data: SecurityAuditMock.getAuditOperation,
    });
    (SecurityAuditService.getSafetyAudit as any).mockResolvedValue({
      data: SecurityAuditMock.getSafetyAudit,
    });
    fireEvent.click(option);
    expect(SecurityAuditService.getAuditOperation).toHaveBeenCalledTimes(1);
    expect(SecurityAuditService.getAuditOperation).toHaveBeenCalledWith({
      module: modules[0],
    });
  });
});
