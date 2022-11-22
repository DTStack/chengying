import React, { useEffect, useState } from 'react';
import { Table, Tooltip } from 'antd';
import moment from 'moment';
import { servicePageService } from '@/services';
import Utils from '@/utils/utils';

export interface AlertRecordProps {
  alert_name: string;
  dashboard_name: string;
  panel_title: string;
  state: string;
  time: string;
  url: string;
}

interface AlertServiceProps {
  dashId: string;
  history: any;
  currentHostDetail?: any;
  currentHosts: any;
  cur_service_info: any;
  ServiceStore: any;
  timeRange?: any;
}

const AlertService: React.FC<AlertServiceProps> = ({
  history,
  currentHostDetail,
  currentHosts,
  dashId,
  cur_service_info,
  ServiceStore,
  timeRange,
}) => {
  const { ip } = currentHostDetail;
  const {
    cur_product: { product_name },
  } = ServiceStore;
  const [alerts, getAlertData] = useState([]);

  useEffect(() => {
    if (dashId) {
      getServiceAlert();
    }
  }, [dashId, ip]);

  // 获取告警列表
  const getServiceAlert = () => {
    servicePageService
      .getAlertsHistory({
        dashboardId: dashId,
        ip: ip || currentHosts.map((o) => o.ip).join(','),
        product_name: product_name,
        service_name: cur_service_info?.service_name,
      })
      .then((res: any) => {
        const rst = res.data;
        if (rst.code === 0) {
          getAlertData(rst?.data?.data);
        }
      });
  };

  const showIconWithState = (state: any): React.ReactNode => {
    let icon = null;
    switch (state) {
      case 'ok':
        icon = (
          <Tooltip title={state}>
            <i className="emicon emicon-heart-fill color-ok"></i>
          </Tooltip>
        );
        break;
      case 'pending':
        icon = (
          <Tooltip title={state}>
            <i className="emicon emicon-warning-circle-fill color-pending"></i>
          </Tooltip>
        );
        break;
      case 'alerting':
        icon = (
          <Tooltip title={state}>
            <i className="emicon emicon-alert-fill color-alerting"></i>
          </Tooltip>
        );
        break;
      case 'no_data':
        icon = (
          <Tooltip title={state}>
            <i className="emicon emicon-question-circle-fill color-no-data"></i>
          </Tooltip>
        );
        break;
      case 'paused':
        icon = (
          <Tooltip title={state}>
            <i className="emicon emicon-poweroff-circle-fill color-paused"></i>
          </Tooltip>
        );
        break;
      default:
        icon = (
          <Tooltip title={state}>
            <i className="emicon emicon-API- color-no-data"></i>
          </Tooltip>
        );
    }
    return icon;
  };

  const getDashboard = (record: AlertRecordProps) => {
    let path: string = `/deploycenter/monitoring/dashdetail?url=${record.url}`;
    Utils.setNaviKey('menu_deploy_center', 'sub_menu_dashboard');
    history.push(path);
  };

  return (
    <Table<AlertRecordProps>
      rowKey={(record: AlertRecordProps) =>
        `${record.panel_title}_${record.time}`
      }
      className="dt-table-fixed-base"
      dataSource={alerts}
      pagination={false}
      scroll={{ y: true }}
      style={{ height: 485 }}
      size="middle">
      <Table.Column
        title="指标名称"
        key="panel_title"
        render={(record: AlertRecordProps): React.ReactNode => {
          return (
            <div>
              {showIconWithState(record.state)}
              &nbsp;&nbsp;
              {record.panel_title}
            </div>
          );
        }}
      />
      <Table.Column title="告警名称" dataIndex="alert_name" />
      <Table.Column
        title="来源仪表盘"
        dataIndex="dashboard_name"
        render={(text: string, record: AlertRecordProps): React.ReactNode => {
          return (
            <a onClick={(e: React.ReactNode) => getDashboard(record)}>{text}</a>
          );
        }}
      />
      <Table.Column
        title="告警时间"
        dataIndex="time"
        render={(text: string): string =>
          text ? moment(text).format('YYYY-MM-DD HH:mm:ss') : '--'
        }
      />
    </Table>
  );
};

export default AlertService;
