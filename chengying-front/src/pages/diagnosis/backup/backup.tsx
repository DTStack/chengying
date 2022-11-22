import * as React from 'react';
import { AppStoreTypes } from '@/stores';
import { connect } from 'react-redux';
import { HeaderStateTypes } from '@/stores/headerReducer';
import {
  Layout,
  Select,
  Table,
  message,
  Form,
  Button,
  Empty,
  Modal,
  Icon,
} from 'antd';
import * as Http from '@/utils/http';
import './style.scss';

const { Content } = Layout;
const Option = Select.Option;

interface IProps {
  authorityList?: any;
  HeaderStore: HeaderStateTypes;
}

interface SideNavState {
  originalData: any[];
  backupList: any[];
  selectedRows: any[];
  selectedRowKeys: any[];
  components: string;
  componentList: any[];
  serviceList: any[];
  services: string;
  hostsList: any[];
  hosts: string;
  page: number;
  size: number;
  total: number;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
  HeaderStore: state.HeaderStore,
});

@(connect(mapStateToProps, undefined) as any)
export default class SideNav extends React.Component<IProps, SideNavState> {
  constructor(props: IProps) {
    super(props);
    this.state = {
      originalData: [], // 原来从接口调用获取到的数据
      backupList: [], // 处理后的表格数据
      selectedRows: [], // 批量删除选中的数据
      selectedRowKeys: [], // 批量选中的key
      components: '', // 组件
      componentList: [], // 组件下拉选项
      serviceList: [], // 服务下拉选项
      services: '', // 服务选择
      hostsList: [], // 组件下拉选项
      hosts: '', // 组件
      page: 1,
      size: 10,
      total: 0,
    };
  }

  componentDidMount() {
    this.getBackupList();
  }

  // 获取备份包数据
  getBackupList() {
    const { id } = this.props.HeaderStore.cur_parent_cluster;

    Http.get('/api/v2/product/backup', { clusterId: id }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const { components, services, hosts } = this.state;
        this.setState(
          {
            originalData: res.data,
          },
          () => {
            if (components == '' && services == '' && hosts == '') {
              this.setBackupData(res.data); // 处理未进行任何下拉选择情况下的表格展示数据
            } else {
              this.filterData(); // 存在筛选条件的时候获取数据格式
            }
            this.setHostsList(); // 获取主机列表
            this.setComponentsArr(); // 获取组件列表
          }
        );
      } else {
        message.error(res.msg);
      }
    });
  }

  // 处理备份包data数据
  setBackupData = (data: any[]) => {
    const backupList = [];
    let id = -1;
    data.map((product) => {
      product.service.map((info) => {
        id++;
        backupList.push({
          id: id,
          componentName: product.product,
          componentVersion: product.version,
          serviceName: info.name,
          hostsInfo: info.host_info,
        });
      });
      return backupList;
    });
    this.setState({
      backupList: backupList,
      total: backupList.length,
    });
  };

  // 表格初始化
  getColumns = () => {
    const { authorityList } = this.props;
    return [
      {
        title: '组件名称',
        dataIndex: 'componentName',
        key: 'componentName',
        ellipsis: true,
      },
      {
        title: '组件版本号',
        dataIndex: 'componentVersion',
        key: 'componentVersion',
        ellipsis: true,
      },
      {
        title: '服务名称',
        dataIndex: 'serviceName',
        key: 'serviceName',
        ellipsis: true,
      },
      {
        title: '操作',
        key: 'opreation',
        width: 120,
        render: (_, record) => {
          return (
            <span>
              {authorityList.sub_menu_backup_manage_delete ? (
                <a
                  onClick={() => {
                    this.toDelete(record);
                  }}>
                  删除
                </a>
              ) : (
                '--'
              )}
            </span>
          );
        },
      },
    ];
  };

  // 删除操作
  toDelete = (record: any) => {
    Modal.confirm({
      title: '确认要删除该备份文件吗？',
      content: '删除后该备份将不可恢复！',
      icon: (
        <Icon type="close-circle" theme="filled" style={{ color: '#FF5F5C' }} />
      ),
      okType: 'danger',
      onOk: () => {
        const { hosts } = this.state;
        let hostsInfo = [];
        if (hosts !== '') {
          record.hostsInfo.map((info) => {
            if (info.ip == hosts) {
              hostsInfo.push(info);
            }
            return hostsInfo;
          });
        } else {
          hostsInfo = record.hostsInfo;
        }
        const params = [
          {
            product: record.componentName,
            version: record.componentVersion,
            service: [
              {
                name: record.serviceName,
                host_info: hostsInfo,
              },
            ],
          },
        ];
        Http.post('/api/v2/product/clean', params).then((res: any) => {
          res = res.data;
          if (res.code === 0) {
            this.getBackupList();
            // this.setState({
            // services: '',
            // components: '',
            // hosts: ''
            // })
            message.success('删除成功');
          } else {
            message.error(res.msg);
          }
        });
      },
      okText: '删除',
    });
  };

  // 分页
  onChangePage = (page: any, pageSize: any) => {
    this.setState({
      page: page,
      size: pageSize,
    });
  };

  // 组件数组下拉选择去重
  setComponentsArr = () => {
    const { originalData, hosts } = this.state;
    const componentList = [];
    if (hosts == '') {
      originalData.map((product) => {
        componentList.push(product.product);
        return componentList;
      });
    } else {
      originalData.map((product) => {
        product.service.map((info) => {
          info.host_info.map((item) => {
            if (item.ip == hosts) {
              componentList.push(product.product);
            }
          });
        });
        return componentList;
      });
    }
    this.setState({
      componentList: componentList.reduce(
        (prev, cur) => (prev.includes(cur) ? prev : [...prev, cur]),
        []
      ),
    });
  };

  // 选择组件
  onChangeComponent = (value: string) => {
    this.setServiceList(value);
    this.setState(
      {
        components: value,
      },
      () => {
        this.filterData();
        this.setHostsList();
      }
    );
  };

  // 选了组件后联动了服务那块的下拉选择
  setServiceList = (param: string) => {
    const serviceList = [];
    const { originalData } = this.state;
    // 处理选择组件后影响了的生成的服务下拉选项
    originalData.map((product) => {
      if (product.product === param) {
        product.service.map((info) => {
          serviceList.push(info.name);
        });
      }
      return serviceList;
    });
    this.setState({
      services: '',
      serviceList: serviceList.reduce(
        (prev, cur) => (prev.includes(cur) ? prev : [...prev, cur]),
        []
      ),
    });
  };

  // 选择服务
  onChangeService = (value: string) => {
    this.setState(
      {
        services: value,
      },
      () => {
        this.filterData();
      }
    );
  };

  // 获取主机下拉选择数组内容
  setHostsList = () => {
    const { originalData, components } = this.state;
    const hostsList = [];
    if (components == '') {
      originalData.map((product) => {
        product.service.map((info) => {
          info.host_info.map((item) => {
            hostsList.push(item.ip);
          });
        });
        return hostsList;
      });
    } else {
      originalData.map((product) => {
        if (product.product == components) {
          product.service.map((info) => {
            info.host_info.map((item) => {
              hostsList.push(item.ip);
            });
          });
        }
        return hostsList;
      });
    }

    this.setState({
      hostsList: hostsList.reduce(
        (prev, cur) => (prev.includes(cur) ? prev : [...prev, cur]),
        []
      ),
    });
  };

  // 选择主机
  onChangeHosts = (value: string) => {
    this.setState(
      {
        hosts: value,
      },
      () => {
        this.setComponentsArr();
        this.filterData();
      }
    );
  };

  // 筛选数据
  filterData = () => {
    const { originalData, components, services, hosts } = this.state;
    const tableData = [];
    let id = -1;
    originalData.map((product) => {
      if (product.product === components && services === '' && hosts === '') {
        // 只选了组件
        product.service.map((info) => {
          id++;
          tableData.push({
            id: id,
            componentName: product.product,
            componentVersion: product.version,
            serviceName: info.name,
            hostsInfo: info.host_info,
          });
        });
      } else if (
        product.product === components &&
        services !== '' &&
        hosts === ''
      ) {
        // 如果存在组件跟服务都选了，没选主机
        product.service.map((info) => {
          if (info.name === services) {
            id++;
            tableData.push({
              id: id,
              componentName: product.product,
              componentVersion: product.version,
              serviceName: info.name,
              hostsInfo: info.host_info,
            });
          }
        });
      } else if (
        product.product === components &&
        services !== '' &&
        hosts !== ''
      ) {
        // 如果存在组件、服务、主机三个都选了
        product.service.map((info) => {
          if (info.name === services) {
            info.host_info.map((item) => {
              if (item.ip === hosts) {
                id++;
                tableData.push({
                  id: id,
                  componentName: product.product,
                  componentVersion: product.version,
                  serviceName: info.name,
                  hostsInfo: info.host_info,
                });
              }
            });
          }
        });
      } else if (
        product.product === components &&
        services === '' &&
        hosts !== ''
      ) {
        // 如果只选了组件、主机，没选服务
        product.service.map((info) => {
          info.host_info.map((item) => {
            if (item.ip === hosts) {
              id++;
              tableData.push({
                id: id,
                componentName: product.product,
                componentVersion: product.version,
                serviceName: info.name,
                hostsInfo: info.host_info,
              });
            }
          });
        });
      } else if (components === '' && services === '' && hosts !== '') {
        // 如果只选了主机，没选组件、服务
        product.service.map((info) => {
          info.host_info.map((item) => {
            if (item.ip === hosts) {
              id++;
              tableData.push({
                id: id,
                componentName: product.product,
                componentVersion: product.version,
                serviceName: info.name,
                hostsInfo: info.host_info,
              });
            }
          });
        });
      }
      return tableData;
    });
    this.setState({
      backupList: tableData,
      total: tableData.length,
    });
  };

  // 批量删除
  onChangeTable = (selectedRowKeys, selectedRows) => {
    this.setState({
      selectedRows: selectedRows,
      selectedRowKeys: selectedRowKeys,
    });
  };

  // 批量删除
  toMoreDelete = () => {
    const { authorityList } = this.props;
    if (authorityList.sub_menu_backup_manage_delete) {
      Modal.confirm({
        title: '确认要批量删除备份文件吗？',
        content: '删除后备份将不可恢复！',
        icon: (
          <Icon
            type="close-circle"
            theme="filled"
            style={{ color: '#FF5F5C' }}
          />
        ),
        okType: 'danger',
        onOk: () => {
          const { selectedRows, hosts } = this.state;
          const paramsArray = [];
          selectedRows.map((item) => {
            if (hosts !== '') {
              item.hostsInfo.map((info) => {
                if (info.ip == hosts) {
                  paramsArray.push({
                    product: item.componentName,
                    version: item.componentVersion,
                    service: [
                      {
                        name: item.serviceName,
                        host_info: [info],
                      },
                    ],
                  });
                }
              });
            } else {
              paramsArray.push({
                product: item.componentName,
                version: item.componentVersion,
                service: [
                  {
                    name: item.serviceName,
                    host_info: item.hostsInfo,
                  },
                ],
              });
            }
            return paramsArray;
          });

          Http.post('/api/v2/product/clean', paramsArray).then((res: any) => {
            res = res.data;
            if (res.code === 0) {
              this.getBackupList();
              this.setState({
                selectedRows: [],
                selectedRowKeys: [],
                // services: '',
                // components: '',
                // hosts: ''
              });
              message.success('删除成功');
            } else {
              message.error(res.msg);
            }
          });
        },
        okText: '删除',
      });
    } else {
      message.error('尚无权限执行该操作');
    }
  };

  render() {
    const {
      componentList,
      components,
      hosts,
      serviceList,
      backupList,
      hostsList,
      size,
      total,
      page,
      services,
      selectedRows,
      selectedRowKeys,
    } = this.state;
    const rowSelection = {
      selectedRowKeys,
      onChange: this.onChangeTable,
    };
    const pagination = {
      size: 'small',
      pageSize: size,
      total: total,
      current: page,
      onChange: this.onChangePage,
      showTotal: (total) => (
        <span>
          共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示
          <span style={{ color: '#3F87FF' }}>{size}</span>条
        </span>
      ),
    };
    // const tableFooter = () => (
    //     <Row style={{ width: '100%' }}>
    //         <Col span={12}>
    //             <Button size="small" type="primary" onClick={this.toMoreDelete}>批量删除</Button>
    //         </Col>
    //         <Col span={12}>
    //             <Pagination
    //                 current={page}
    //                 pageSize={size}
    //                 size='small'
    //                 total={total}
    //                 onChange={this.onChangePage}
    //                 style={{ float: 'right' }}
    //                 showTotal={(total) => <span>
    //                     共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示10条
    //                 </span>}
    //             />
    //         </Col>
    //     </Row>
    // );
    return (
      <Layout id="monitorDataFlowContainer">
        <Content>
          <div className="backup-page">
            <div className="top-navbar clearfix">
              <Form className="ant-advanced-search-form">
                <div style={{ minWidth: '1000px', display: 'flex' }}>
                  <Form.Item label="组件：">
                    <Select
                      className="dt-form-shadow-bg mr-20"
                      style={{ width: 160 }}
                      placeholder="请选择组件"
                      onChange={this.onChangeComponent}
                      value={components || undefined}>
                      {Array.isArray(componentList) &&
                        componentList.map((item: any, index: number) => (
                          <Option key={index} value={item}>
                            {item}
                          </Option>
                        ))}
                    </Select>
                  </Form.Item>
                  <Form.Item label="服务：">
                    <Select
                      className="dt-form-shadow-bg mr-20"
                      style={{ width: 160 }}
                      placeholder="请选择服务"
                      onChange={this.onChangeService}
                      value={services || undefined}>
                      {Array.isArray(serviceList) &&
                        serviceList.map((item: any) => (
                          <Option key={item} value={item}>
                            {item}
                          </Option>
                        ))}
                    </Select>
                  </Form.Item>
                  <Form.Item label="主机：">
                    <Select
                      className="dt-form-shadow-bg mr-20"
                      style={{ width: 160 }}
                      placeholder="请选择主机"
                      onChange={this.onChangeHosts}
                      value={hosts || undefined}>
                      {Array.isArray(hostsList) &&
                        hostsList.map((item: any) => (
                          <Option key={item} value={item}>
                            {item}
                          </Option>
                        ))}
                    </Select>
                  </Form.Item>
                </div>
              </Form>
            </div>
            <div className="log_table box-shadow-style">
              <Table
                rowKey="id"
                className="dt-table-fixed-base"
                style={{ height: 'calc(100vh - 40px - 40px - 72px)' }}
                scroll={{ y: true }}
                locale={{ emptyText: <Empty /> }}
                rowSelection={rowSelection}
                columns={this.getColumns()}
                dataSource={backupList}
                pagination={pagination}
                // pagination={false}
                // footer={tableFooter}
              />
              <span className="batch-delete" style={{ zIndex: 100 }}>
                <Button
                  disabled={total == 0 || selectedRows.length == 0}
                  size="small"
                  type="primary"
                  onClick={this.toMoreDelete}>
                  批量删除
                </Button>
              </span>
            </div>
          </div>
        </Content>
      </Layout>
    );
  }
}
