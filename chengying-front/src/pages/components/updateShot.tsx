import * as React from 'react';
import { Modal, Table, Progress, Tooltip, Icon } from 'antd';

interface Prop {
  showModal: boolean;
  dataList: any[];
  closeModal: () => void;
  handleTableChange: () => void;
  status: string;
  shotPagination: any;
  title: string;
}

class UpdateShot extends React.Component<Prop, any> {
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
        title: '更新进度',
        dataIndex: 'progressbar',
        render: (e: any, record: any) => {
          const isSuccess =
            (record.progress === 100 || record.progress === 50) &&
            record.status_message === '';
          return (
            <Progress
              percent={record.progress}
              status={isSuccess ? 'success' : 'exception'}
            />
          );
        },
      },
      {
        title: '更新状态',
        dataIndex: 'status',
        filters: [
          {
            text: 'success',
            value: 'success',
          },
          {
            text: 'fail',
            value: 'fail',
          },
          {
            text: 'update',
            value: 'update',
          },
        ],
        render: (text: any, record: any) => {
          let serviceStatus = {};
          const isSuccess =
            (record.progress === 100 || record.progress === 50) &&
            record.status_message === '';
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
    const { dataList, showModal, shotPagination, title } = this.props;
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
        title={title}
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

export default UpdateShot;
