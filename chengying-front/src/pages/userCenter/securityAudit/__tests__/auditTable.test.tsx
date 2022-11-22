import * as React from 'react';
import { cleanup, render, RenderResult } from '@testing-library/react';
import AuditTable from '../auditTable';
import '@testing-library/jest-dom/extend-expect';
import { SecurityAuditService } from '@/services';
import { SecurityAuditMock } from '@/mocks';

jest.mock('@/services');

const reqParams = {
  from: undefined,
  to: undefined,
  module: undefined,
  operation: undefined,
  content: undefined,
  operator: undefined,
  ip: undefined,
};

describe('audit table render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    (SecurityAuditService.getSafetyAudit as any).mockResolvedValue({
      data: SecurityAuditMock.getSafetyAudit,
    });
    wrapper = render(<AuditTable reqParams={reqParams} />);
  });

  afterEach(cleanup);

  test('table render', () => {
    const columns = wrapper.container.getElementsByClassName(
      'ant-table-column-title'
    );
    const columnNames = [
      '时间',
      '操作人',
      '操作模块',
      '动作',
      '详细内容',
      '来源IP',
    ];
    for (const i in columns) {
      if (columns[i] instanceof HTMLElement) {
        expect(columns[i].textContent).toBe(columnNames[i]);
      }
    }
  });
});
