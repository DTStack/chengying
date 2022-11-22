import React, { useState, useEffect } from 'react';
import { Table, Switch, Modal, Badge, message, Button } from 'antd';
import { servicePageService } from '@/services';

interface HealthCheckProps {
  currentHostDetail: any;
  curProduct: any;
  curService: any;
}

interface RecordProps {
  record_id: number;
  cluster_id: number;
  product_name: string;
  service_name: string;
  script_name: string;
  script_name_display: string;
  exec_status: number;
  error_message: string;
  auto_exec: boolean;
  ip: string;
  start_time: string;
  end_time: string;
  create_time: string;
  update_time: string;
}

const HealthStatusMaps: { [key: number]: string[] } = {
  0: ['default', '未就绪'],
  1: ['processing', '进行中'],
  2: ['success', '正常'],
  3: ['error', '异常'],
};

let timer = null;

const HealthCheck: React.FC<HealthCheckProps> = (props) => {
  const { curProduct, curService, currentHostDetail } = props;
  const [healths, getHealth] = useState<Array<RecordProps>>([]);

  useEffect(() => {
    getHealthCheck();
    return () => clearTimeout(timer);
  }, [currentHostDetail]);

  const changeAutoCheck = (record: RecordProps) => {
    const auto_exec = !record?.auto_exec;
    servicePageService
      .setAutoexecSwitch({
        record_id: record.record_id,
        auto_exec,
        product_name: record?.product_name,
        service_name: record?.service_name,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          message.success(data.msg);
          // clearTimeout(timer);
          getHealthCheck();
        }
      });
  };

  const manualCheck = (record: RecordProps) => {
    servicePageService
      .manualExecution({
        record_id: record.record_id,
        product_name: record?.product_name,
        service_name: record?.service_name,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          message.success(data.msg);
          // clearTimeout(timer);
          getHealthCheck();
        }
      });
  };

  const getHealthCheck = () => {
    servicePageService
      .getHealthCheck({
        product_name: curProduct?.product_name,
        service_name: curService?.service_name,
        ip: currentHostDetail?.ip,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          getHealth(data?.data?.list);
          clearTimeout(timer);
          timer = setTimeout(() => {
            getHealthCheck();
          }, 5000);
        } else {
          message.error(data.msg);
        }
      });
  };

  const handleManual = (record: RecordProps) => {
    Modal.confirm({
      title: '手动执行',
      content: '确认立即执行该脚本？状态将在执行完成后更新，请耐心等待。',
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        manualCheck(record);
      },
    });
  };

  return (
    <Table
      rowKey="record_id"
      className="dt-table-fixed-base"
      dataSource={healths}
      pagination={false}
      scroll={{ y: true }}
      style={{ height: 485 }}
      size="middle">
      <Table.Column title="脚本名称" dataIndex="script_name_display" />
      <Table.Column
        title="状态"
        dataIndex="exec_status"
        render={(text: string): React.ReactNode => {
          const [status, desc] = HealthStatusMaps[text];
          return (
            <span>
              <Badge status={status} />
              {desc}
            </span>
          );
        }}
      />
      <Table.Column
        title="检查时间"
        dataIndex="start_time"
        render={(text: string) => text || '--'}
      />
      <Table.Column
        title="详情"
        dataIndex="error_message"
        ellipsis
        render={(text: string) => {
          return <span>{text || '--'}</span>;
        }}
      />
      <Table.Column
        title="定时执行"
        dataIndex="auto_exec"
        render={(text: boolean, record: RecordProps): React.ReactNode => (
          <Switch
            onChange={() => changeAutoCheck(record)}
            checkedChildren="开"
            unCheckedChildren="关"
            checked={text}
          />
        )}
      />
      <Table.Column
        title="操作"
        width={120}
        render={(record: RecordProps): React.ReactNode => (
          <Button
            type="link"
            style={{ padding: 0 }}
            disabled={record.exec_status === 1}
            onClick={() => handleManual(record)}>
            手动执行
          </Button>
        )}
      />
    </Table>
  );
};

export default HealthCheck;
