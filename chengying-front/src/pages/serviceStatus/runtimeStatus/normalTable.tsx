import * as React from 'react';
import { Table, Button, Tooltip, Badge, Modal, Icon } from 'antd';
import { ServiceProp } from '@/pages/serviceStatus';
import utils from '@/utils/utils';
const { confirm } = Modal;

interface IProps extends ServiceProp {
  sHosts: any[];
  isRestart: boolean;
  haRoleClass: string; // 是否隐藏ha_role列
  selectedRow: string | number;
  ipRender: (text: string, record: any) => React.ReactNode; // ip地址渲染
  cur_service_info: any;
  getServiceGroup: any;
  handleOpenHostConfig: (record: any) => void; // 查看配置
  // handleStartHost: (record: any) => void; // 启动
  // handleStopHost: (record: any) => void; // 停止
  handleClickRow: (record: any) => void; // 停止
}
export default class NormalTable extends React.PureComponent<IProps, any> {
  // 配置
  handleTableConfig = (text, record) => {
    const { authorityList } = this.props;
    if (!authorityList.service_view) {
      return '--';
    }
    if ('schema' in record && record.schema) {
      const { Instance } = JSON.parse(record.schema);
      return (
        <>
          <a
            style={{
              display: Instance.ConfigPaths ? 'inline' : 'none',
              paddingRight: '5px',
            }}
            onClick={(e) => {
              e.stopPropagation();
              this.props.handleOpenHostConfig(record);
            }}>
            配置
          </a>
          <a
            style={{ display: Instance.Logs ? 'inline' : 'none' }}
            onClick={(e) => {
              e.stopPropagation();
              this.handleToLog(record);
            }}>
            日志
          </a>
        </>
      );
    } else {
      return '--';
    }
  };

  // columns 初始化
  initColumns = () => {
    const {
      isRestart,
      authorityList,
      haRoleClass,
      HeaderStore: { cur_parent_cluster },
      ipRender,
    } = this.props;
    const CAN_SERVER_ACTION = authorityList?.service_start_stop;
    const isKubernetes = cur_parent_cluster?.type === 'kubernetes';
    const columns = [
      {
        title: 'IP地址',
        dataIndex: 'ip',
        key: 'ip',
        render: ipRender,
        width: '16.5%',
      },
      {
        title: '角色',
        dataIndex: 'ha_role',
        key: 'ha_role',
        // width: '9.5%',
        className: haRoleClass,
        render: (text: any) => {
          return text || ' - ';
        },
      },
      {
        title: '运行状态',
        dataIndex: 'status',
        key: 'status',
        // width: '10%',
        render: (text: any, record: any) => {
          return (
            <Tooltip title={record.status_message}>
              <span>{text}</span>
            </Tooltip>
          );
        },
      },
      {
        title: '健康状态',
        dataIndex: 'health_state',
        key: 'health_state',
        width: '10.5%',
        render: (text: any, record: any) => {
          let status = '-';
          let status_cls: any = {};
          if (record.status === 'running') {
            switch (record.health_state) {
              case -2:
                status = '等待';
                status_cls = {
                  color: '#333',
                };
                break;
              case 0:
                status = '不健康';
                status_cls = {
                  color: '#FF5F5C',
                };
                break;
              case 1:
                status = '健康';
                status_cls = {
                  color: '#12BC6A',
                };
                break;
            }
            status_cls.display = 'inline-block';
          } else {
            status_cls.display = 'none';
          }
          if (record.health_state === -1) {
            status = '未设置';
            status_cls.display = 'none';
          }

          return (
            <>
              {status_cls.display !== 'none' && (
                <Badge color={status_cls.color} />
              )}
              <span style={{ paddingRight: 8 }}>{status}</span>
            </>
          );
        },
      },
      {
        title: '组件版本',
        dataIndex: 'product_version',
        key: 'product_version',
        // width: "10%",
      },
      {
        title: '更新时间',
        dataIndex: 'update_time',
        key: 'update_time',
        width: '21.1%',
      },
      {
        title: '查看',
        dataIndex: 'config',
        key: 'config',
        // width: '10%',
        render: (text: any, record: any) =>
          this.handleTableConfig(text, record),
      },
      {
        title: '停止/启动',
        dataIndex: 'option',
        key: 'option',
        render: (text: any, record: any) => {
          // let state = false;
          const startCls = {
            display: 'none',
          };
          const stopCls = {
            display: 'none',
          };
          switch (record.status) {
            case 'stopped':
            case 'run fail':
              // state = false;
              startCls.display = 'inline-block';
              stopCls.display = 'none';
              break;
            case 'running':
            case 'health-checked':
            case 'stop fail':
            case 'health-check fail':
              // state = true;
              startCls.display = 'none';
              stopCls.display = 'inline-block';
              break;
          }
          return (
            <span>
              {CAN_SERVER_ACTION ? (
                <span>
                  <Button
                    disabled={isRestart || record.isDisable}
                    onClick={(e: any) => {
                      e.stopPropagation();
                      this.handleStartHost(record);
                    }}
                    type="primary"
                    ghost
                    style={{ ...startCls, width: 64, height: 26, padding: 0 }}>
                    启动
                  </Button>
                  <Button
                    disabled={isRestart || record.isDisable}
                    onClick={(e: any) => {
                      e.stopPropagation();
                      this.handleStopHost(record);
                    }}
                    type="primary"
                    ghost
                    style={{ ...stopCls, width: 64, height: 26, padding: 0 }}>
                    停止
                  </Button>
                </span>
              ) : (
                '--'
              )}
            </span>
          );
        },
      },
    ];
    if (isKubernetes) {
      columns.pop();
    }

    if (haRoleClass === 'hide') {
      const roleIdx = columns.findIndex(
        (item: any) => item.dataIndex === 'ha_role'
      );
      roleIdx > -1 && columns.splice(roleIdx, 1);
    }
    return columns;
  };

  // 实例启停
  handleStartHost = (record: any) => {
    let index = null;
    const { sHosts } = this.props;
    const { cur_product } = this.props.ServiceStore;
    const { disableInstance, startInstance, getHostsList } = this.props.actions;
    // const hosts = cur_product.product.Service[cur_service_info.service_name].hosts;
    for (const i in sHosts) {
      if (sHosts[i].id === record.id) {
        index = i;
      }
    }
    disableInstance({
      instance_index: index,
      service_name: record.service_name,
    });
    startInstance(
      {
        agent_id: record.agent_id,
        instance_index: index,
        product_name: cur_product.product_name,
        service_name: record.service_name,
      },
      () => {
        // getServiceGroup();
        getHostsList({
          product_name: cur_product.product_name,
          service_name: record.service_name,
        });
      }
    );
  };

  // 实例停止
  handleStopHost = (record: any) => {
    let index: any = null;
    const { sHosts } = this.props;
    const { cur_product } = this.props.ServiceStore;
    const { disableInstance, stopInstance, getHostsList } = this.props.actions;
    // const hosts = cur_product.product.Service[cur_service_info.service_name].hosts || {};
    for (const i in sHosts) {
      if (sHosts[i].id === record.id) {
        index = i;
      }
    }

    confirm({
      title: '确定要停用该实例吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okButtonProps: {
        disabled: record.isDisable,
      },
      onOk() {
        disableInstance({
          instance_index: index,
          service_name: record.service_name,
        });
        stopInstance(
          {
            agent_id: record.agent_id,
            instance_index: index,
            product_name: cur_product.product_name,
            service_name: record.service_name,
          },
          () => {
            // getServiceGroup();
            getHostsList({
              product_name: cur_product.product_name,
              service_name: record.service_name,
            });
          }
        );
      },
    });
  };

  // 跳转到日志页面
  handleToLog = (item: any) => {
    const schema = JSON.parse(item.schema);
    utils.setNaviKey('menu_ops_center', 'sub_menu_diagnose_log');
    this.props.history.push('/opscenter/diagnosis/log', {
      product_name: item.product_name,
      host_ip: item.ip,
      service_name: item.service_name,
      logpaths: schema.Instance.Logs,
      log_service_id: item.id,
    });
  };

  render() {
    const { sHosts, handleClickRow, selectedRow } = this.props;
    return (
      <Table
        rowKey="id"
        scroll={{ y: 49 * 10, x: false }}
        size="middle"
        className="dt-pagination-lower box-shadow-style"
        pagination={false}
        columns={this.initColumns()}
        dataSource={sHosts}
        onRow={(record) => {
          return {
            onClick: (event) => handleClickRow(record), // 点击行
          };
        }}
        rowClassName={(record, index) =>
          record.ip === selectedRow ? 'recordSelectedRow' : ''
        }
      />
    );
  }
}
