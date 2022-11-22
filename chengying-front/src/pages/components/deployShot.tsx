import * as React from 'react';
import { Modal, Table, Progress, Tooltip, Icon } from 'antd';

interface Prop {
  showModal: boolean;
  dataList: any[];
  closeModal: () => void;
  handleTableChange: () => void;
  status: string;
  shotPagination: any;
}

class DeployShot extends React.Component<Prop, any> {
  initColumns = () => {
    return [
      {
        title: '执行时间',
        dataIndex: 'update_time',
      },
      {
        title: '服务名称',
        dataIndex: 'service_name',
      },
      {
        title: '服务版本号',
        dataIndex: 'service_version',
      },
      {
        title: '主机IP',
        dataIndex: 'ip',
      },
      {
        title: '组件版本号',
        dataIndex: 'product_version',
      },
      {
        title: '启动及初始化',
        dataIndex: 'progressbar',
        render: (e: any, record: any) => {
          // 区分强制卸载，强制停止
          // const isForce = record.status_message === "force uninstall" || record.status_message === "force stop"
          const isSuccess =
            record.progress === 100 && record.status_message === '';
          return (
            <Progress
              percent={record.progress}
              status={isSuccess ? 'success' : 'exception'}
            />
          );
        },
      },
      {
        title: '启动状态',
        dataIndex: 'status',
        filters: [
          {
            text: 'installing',
            value: 'installing',
          },
          {
            text: 'installed',
            value: 'installed',
          },
          {
            text: 'install fail',
            value: 'install fail',
          },
          {
            text: 'uninstalling',
            value: 'uninstalling',
          },
          {
            text: 'uninstalled',
            value: 'uninstalled',
          },
          {
            text: 'uninstall fail',
            value: 'uninstall fail',
          },
          {
            text: 'running',
            value: 'running',
          },
          {
            text: 'run fail',
            value: 'run fail',
          },
          {
            text: 'health-checked',
            value: 'health-checked',
          },
          {
            text: 'health-check fail',
            value: 'health-check fail',
          },
          {
            text: 'health-check cancelled',
            value: 'health-check cancelled',
          },
          {
            text: 'stopped',
            value: 'stopped',
          },
          {
            text: 'stopping',
            value: 'stopping',
          },
          {
            text: 'stop fail',
            value: 'stop fail',
          },
        ],
        render: (text: any, record: any) => {
          let serviceStatus = {};
          const isSuccess =
            record.progress === 100 && record.status_message === '';
          if (isSuccess) {
            serviceStatus = {
              color: '#12BC6A',
            };
          } else {
            serviceStatus = {
              color: '#FF5F5C',
            };
          }
          return (
            <div>
              <span style={serviceStatus}>{text}</span>
              <Tooltip title={record.status_message}>
                {record.status_message ? (
                  <Icon style={{ marginLeft: 3 }} type="info-circle" />
                ) : null}
              </Tooltip>
            </div>
          );
        },
      },
    ];
  };

  render() {
    const columns = this.initColumns();
    const { dataList, showModal, shotPagination } = this.props;
    const pagination = {
      size: 'small',
      current: shotPagination.current,
      total: shotPagination.count,
      pageSize: shotPagination.limit,
    };
    return (
      <Modal
        footer={null}
        width={950}
        destroyOnClose={true}
        visible={showModal}
        title="部署快照"
        onOk={() => this.props.closeModal()}
        onCancel={() => this.props.closeModal()}>
        <Table
          rowKey="id"
          size="middle"
          className="dt-pagination-lower"
          columns={columns}
          dataSource={dataList}
          onChange={this.props.handleTableChange}
          pagination={pagination}
        />
      </Modal>
    );
  }
}

export default DeployShot;
