import * as React from 'react';
import { bindActionCreators, Dispatch } from 'redux';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { addHostService } from '@/services';
import * as HostAction from '@/actions/hostAction';
import classnames from 'classnames';
import {
  Button,
  Select,
  Input,
  Table,
  Icon,
  Modal,
  Tooltip,
  Popover,
  Progress,
  Collapse,
  Alert,
  Form,
  Tag,
  message,
  Empty,
} from 'antd';
import { hostStatusInfoMap } from '@/constants/const';
import CreateHostComp from '@/components/createNewHost';
import FileLogShow from '@/components/fileLogShow';
import EditableCell from './tableCell';
import { formItemCenterLayout } from '@/constants/formLayout';
import CommandPopbox from '@/pages/command/details/popbox';
import './host.scss';
const Option = Select.Option;
const Search = Input.Search;
const Panel = Collapse.Panel;
const FormItem = Form.Item;

const pageSize = 10;

type progressStatus = 'exception' | 'active' | 'success';

const mapStateToProps = (state: AppStoreTypes) => ({
  filterSelectedItem: state.HostStore.filterSelectedItem,
  filterData: state.HostStore.filterData,
  selectRows: state.HostStore.selectRows,
  searchValue: state.HostStore.searchValue,
  clusterHostList: state.HostStore.clusterHostList,
  pager: state.HostStore.pager,
  clusterHostGroupList: state.HostStore.clusterHostGroupList,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, HostAction), dispatch),
});
interface HostProps {
  filterSelectedItem: any[];
  filterData: any[];
  selectRows: any[];
  searchValue: string;
  clusterHostList: any[];
  pager: any;
  actions: any;
  location?: any;
  form: any;
  clusterHostGroupList: any[];
  clusterInfo: any;
}
interface HostState {
  addedHost: any;
  queryParams: {
    'sort-by': any;
    'sort-dir': any;
    q_host: string;
    hostGroupName: string;
    limit: any;
    start: any;
    currentPage: number;
    is_running: [];
    status: [];
    roles: string;
  };
  aHostDetail: {
    status: string;
    status_msg: string;
  };
  showModal: boolean;
  rootLink: any[];
  activeKey: any; // 联动分组下拉菜单和点击panel key
  selectedRowKeys: any[];
  selectedRows: any[];
  hostMoveModal: boolean;
  isNewGroup: boolean;
  tableLoading: boolean;
  clusterHostList: any[];
  visibleInfo: any;
  total: number;
}
@(connect(mapStateToProps, mapDispatchToProps) as any)
class HostPage extends React.Component<HostProps, HostState> {
  constructor(props: object) {
    super(props as HostProps);
  }

  componentDidMount() {
    this.loadData();
    this.getHostGroups(); // 暂不执行(请求分组列表)
  }

  componentDidUpdate() {
    const { clusterHostGroupList } = this.props;
    if (
      this.state.activeKey === '' &&
      Array.isArray(clusterHostGroupList) &&
      clusterHostGroupList.length
    ) {
      const defaultKey = this.props.clusterHostGroupList[0] || '';
      this.setState(
        {
          activeKey: defaultKey,
          queryParams: Object.assign(this.state.queryParams, {
            hostGroupName: defaultKey,
          }),
        },
        () => {
          this.loadData();
          this.getHostGroups();
        }
      );
    }
  }

  componentWillUnmount() {
    clearTimeout(this.timer);
  }

  public state: HostState = {
    addedHost: null,
    queryParams: {
      'sort-by': 'id',
      'sort-dir': 'desc',
      q_host: '',
      hostGroupName: '',
      limit: pageSize,
      start: 0,
      currentPage: 1,
      is_running: [],
      status: [],
      roles: '',
    },
    aHostDetail: {
      status: '',
      status_msg: '',
    },
    showModal: false,
    rootLink: [
      '/deploycenter/cluster/list',
      '/deploycenter/appmanage/products',
      '/installguide',
      '/deploycenter/appmanage/installs',
    ],
    activeKey: '',
    selectedRowKeys: [], // 选择表格列表数据项
    selectedRows: [],
    hostMoveModal: false,
    isNewGroup: false,
    tableLoading: false,
    clusterHostList: [],
    visibleInfo: {
      visible: false,
      title: '日志查看',
      type: 'log',
      execId: '',
    },
    total: 0,
  };

  private timer: any = null;

  progressProp = (status: number) => {
    const { clusterInfo = {} } = this.props;
    switch (status) {
      case 0:
        return {
          percent: 0,
          title: '等待安装',
          color: 'host-status-wait-install',
          // color: 'green-dot dot'
        };
      case -1:
        return {
          percent: 30,
          status: 'exception' as progressStatus,
          title: '管控安装失败',
          // color: "host-status-control-failed"
          color: 'red-dot dot',
        };
      case 1:
        return {
          percent: 30,
          status: 'active' as progressStatus,
          title: '管控安装成功',
          // color: "host-status-control-successed"
          color: 'green-dot dot',
        };
      case -2:
        return {
          percent: clusterInfo.type === 'hosts' ? 60 : 40,
          status: 'exception' as progressStatus,
          title: 'script安装失败',
          // color: "host-status-script-failed"
          color: 'red-dot dot',
        };
      case 2:
        return {
          percent: clusterInfo.type === 'hosts' ? 60 : 40,
          status: 'active' as progressStatus,
          title: 'script安装成功',
          // color: "host-status-script-successed"
          color: 'green-dot dot',
        };
      case -3:
        return {
          percent: clusterInfo.type === 'hosts' ? 100 : 50,
          status: 'exception' as progressStatus,
          title: '主机初始化失败',
          // color: "host-status-init-failed"
          color: 'red-dot dot',
        };
      case 3:
        return {
          percent: clusterInfo.type === 'hosts' ? 100 : 50,
          status: (clusterInfo.type === 'hosts'
            ? 'success'
            : 'active') as progressStatus,
          title: '主机初始化成功',
          // color: "host-status-init-successed"
          color: 'green-dot dot',
        };
      case -4:
        return {
          percent: 50,
          status: 'exception' as progressStatus,
          title: '主机下线',
          color: 'red-dot dot',
        };
      case -5:
        return {
          percent: 60,
          status: 'exception' as progressStatus,
          title: 'K8S DOCKER初始化失败',
          // color: "host-status-init-failed"
          color: 'red-dot dot',
        };
      case 5:
        return {
          percent: 60,
          status: 'active' as progressStatus,
          title: 'K8S DOCKER初始化成功',
          // color: "host-status-init-successed"
          color: 'green-dot dot',
        };
      case -6:
        return {
          percent: 80,
          status: 'exception' as progressStatus,
          title: 'K8S NODE初始化失败',
          // color: "host-status-init-failed"
          color: 'red-dot dot',
        };
      case 6:
        return {
          percent: 80,
          status: 'active' as progressStatus,
          title: 'K8S NODE初始化成功',
          // color: "host-status-init-successed"
          color: 'green-dot dot',
        };
      case -7:
        return {
          percent: 100,
          status: 'exception' as progressStatus,
          title: 'K8S NODE部署失败',
          color: 'red-dot dot',
        };
      case 7:
        return {
          percent: 100,
          status: 'success' as progressStatus,
          title: 'K8S NODE部署成功',
          color: 'green-dot dot',
        };
      default:
        return { percent: 100, status: 'exception' as progressStatus };
    }
  };

  // 按主机ip或名称搜索
  handleSearchByHost = (value: string) => {
    const { queryParams } = this.state;
    this.setState(
      {
        selectedRowKeys: [],
        queryParams: Object.assign(queryParams, {
          q_host: value,
          start: 0,
          currentPage: 1,
        }),
      },
      () => {
        this.getHostGroups();
        this.loadData();
      }
    );
  };
  loadData = () => {
    const { clusterInfo } = this.props;
    const { queryParams } = this.state;
    this.setState({ tableLoading: true });
    const is_running: string = queryParams?.is_running.length
      ? queryParams?.is_running.join(',')
      : '';
    const status: string = queryParams?.status.length
      ? queryParams?.status.join(',')
      : '';
    this.props.actions.getClusterHostList(
      {
        'sort-by': queryParams['sort-by'],
        'sort-dir': queryParams['sort-dir'],
        limit: queryParams.limit,
        start: queryParams.start,
        host_or_ip: queryParams.q_host,
        is_running,
        status,
        group: queryParams.hostGroupName,
        cluster_id: clusterInfo?.id,
        cluster_type: clusterInfo?.type,
        role: clusterInfo?.type === 'hosts' ? undefined : queryParams.roles,
      },
      (res: any) => {
        const { hosts, count } = res;
        console.log(res);
        this.setState(
          {
            clusterHostList: hosts,
            total: count,
            tableLoading: false,
          },
          () => {
            clearTimeout(this.timer);
            this.timer = setTimeout(() => {
              this.loadData();
              this.getHostGroups();
            }, 3000);
          }
        );
      }
    );
  };

  /**
   * TODO 当筛选产品组件，主机ip时，都要更新主机列表，
   * 只显示含有当前信息的分组列表(产品设计是有问题的~~)
   */
  getHostGroups = () => {
    const { clusterInfo = {} } = this.props;
    this.props.actions.getClusterhostGroupLists({
      type: clusterInfo.type,
      id: clusterInfo.id,
    });
  };

  goToAdd = () => {
    this.setState({
      showModal: true,
    });
  };

  afterHostInstall = () => {
    this.setState({
      showModal: false,
    });
    this.loadData();
  };

  hostDelete = () => {
    const { selectedRows } = this.state;
    const aid = [];
    selectedRows.forEach((item) => {
      if (!item.is_running) {
        aid.push(item.id);
      }
    });
    const config = { aid: aid };
    if (aid.length > 0) {
      Modal.confirm({
        title: '只可下架agent已停止的主机，确定下架这些主机吗？',
        content:
          '主机下架将会将主机上的数据库信息、服务信息清除掉，请谨慎操作！',
        icon: <Icon type="exclamation-circle" theme="filled" />,
        okText: '确定',
        cancelText: '取消',
        onCancel: () => {
          console.log('Cancel');
        },
        onOk: () => {
          addHostService.deleteHost(config).then((response) => {
            const data = response.data;
            if (data.code === 0) {
              message.success('下架成功！');
              this.getHostGroups();
              this.loadData();
            } else {
              message.error(data.msg);
            }
          });
        },
      });
    } else {
      Modal.confirm({
        title: '当前无agent停止的主机，不支持下架',
        icon: <Icon type="exclamation-circle" theme="filled" />,
        okText: '确定',
      });
    }
  };

  handleEvent = (record) => {
    const { visibleInfo } = this.state;
    this.setState({
      visibleInfo: {
        ...visibleInfo,
        visible: !visibleInfo.visible,
        execId: record ? record.exec_id : '',
      },
    });
  };

  initColumns = () => {
    const { clusterInfo } = this.props;
    const { queryParams } = this.state;
    const statusInfo = hostStatusInfoMap[clusterInfo.type];
    const staticColumn = [
      {
        title: '主机IP',
        dataIndex: 'ip',
        key: 'ip',
        sorter: true,
        width: 150,
        fixed: 'left' as 'left',
      },
      {
        title: '主机名称',
        dataIndex: 'hostname',
        key: 'hostname',
        render: (e: string, record: any) => (
          <Popover content={<span>sid:{record.sid}</span>}>{e}</Popover>
        ),
      },
      {
        title: '最近心跳时间',
        dataIndex: 'updated',
        key: 'updated',
        sorter: true,
      },
      {
        title: '主机创建时间',
        dataIndex: 'created',
        key: 'created',
        sorter: true,
      },
      {
        title: '启动用户',
        dataIndex: 'run_user',
        key: 'run_user',
        render: (text: string) => (
          <Tooltip placement="right" title={text}>
            {text && text.length > 12 ? text.slice(0, 13) + '...' : text}
          </Tooltip>
        ),
      },
      {
        title: 'agent状态',
        dataIndex: 'is_running',
        key: 'is_running',
        width: 150,
        filters: [
          { text: '运行中', value: 'true' },
          { text: '已停止', value: 'false' },
        ],
        filteredValue: queryParams?.is_running,
        filterMultiple: true,
        render: (value: number) => {
          const cls = value ? 'green-dot dot' : 'red-dot dot';
          return (
            <span className="status-box">
              <span className={cls} />
              <span>{value ? '运行中' : '已停止'}</span>
            </span>
          );
        },
      },
      {
        title: (
          <span>
            初始化状态{' '}
            <Tooltip title={statusInfo.title}>
              <Icon type="info-circle" />
            </Tooltip>
          </span>
        ),
        dataIndex: 'errorMsg',
        key: 'errorMsg',
        filters: statusInfo.filters,
        filteredValue: queryParams?.status,
        filterMultiple: true,
        render: (data: any, record: any) => {
          const finalStatus = statusInfo.finalStatus;
          return (
            <span>
              {([-3, -2, -1, -4, -5, -6, -7, finalStatus].includes(
                record.status
              ) && (
                <span className="status-box">
                  <span className={this.progressProp(record.status).color} />
                  <span>{this.progressProp(record.status).title}</span>

                  {record.status !== finalStatus && (
                    <Tooltip title={record.errorMsg}>
                      &nbsp;
                      <Icon type="info-circle" />
                    </Tooltip>
                  )}
                </span>
              )) || (
                <Progress
                  className="status-progress"
                  style={{ width: '80%' }}
                  format={(percent: number) => `${percent}% ${record.errorMsg}`}
                  {...this.progressProp(record.status)}
                />
              )}
            </span>
          );
        },
      },
      {
        title: '查看',
        render: (text: string, record: any) => (
          <a onClick={(e) => this.handleEvent(record)}>日志</a>
        ),
      },
    ];
    const roleColumn = [
      {
        title: '主机角色',
        dataIndex: 'roles',
        width: 250,
        filters: [
          { text: 'Etcd', value: 'Etcd' },
          { text: 'Control', value: 'Control' },
          { text: 'Worker', value: 'Worker' },
          { text: '全部', value: 'all' },
        ],
        render: (text, record) => {
          const roles = Object.keys(record.roles || {});
          return roles.map(
            (item) =>
              record.roles[item] && (
                <Tag className="c-cluster__tag" key={item}>
                  {item}
                </Tag>
              )
          );
        },
      },
    ];
    return clusterInfo.type === 'hosts'
      ? staticColumn
      : [...staticColumn, ...roleColumn];
  };

  onCellChange = (key) => {
    return (value, fn) => {
      const target = this.props.clusterHostGroupList.find(
        (item) => item === key
      );
      addHostService
        .updateGroupName({
          old: target,
          new: value,
        })
        .then((res: any) => {
          if (res.data.code === 0) {
            message.success('分组名称修改成功');
            fn();

            this.setState(
              {
                queryParams: Object.assign(this.state.queryParams, {
                  hostGroupName: value,
                }),
                activeKey: value,
              },
              () => {
                this.getHostGroups();
                this.loadData();
              }
            );
          } else {
            message.error(res.data.msg);
          }
        });
      //   拿到新分组名称和分组信息
      // 请求更新接口成功之后更新分组信息
    };
  };

  /**
   * movehost
   */
  tableFooter = (currentPageData) => {
    return (
      <div>
        <Button className="mr-10" onClick={this.clickMoveHost} icon="swap">
          移动到
        </Button>
        <Button
          type="danger"
          className="btn_dele"
          onClick={this.hostDelete}
          icon="delete">
          下架
        </Button>
      </div>
    );
  };

  clickMoveHost = () => {
    this.setState({
      hostMoveModal: true,
    });
  };

  /**
   * 校验分组名称
   */
  checkGroupName = (name) => {
    const reg = /^[a-zA-Z0-9_@*#/]+$/;
    if (name === '') {
      message.error('分组名称不可为空！');
      return false;
    } else if (name.length > 20) {
      message.error('分组名称不得超过20字符');
      return false;
    } else if (!reg.test(name)) {
      message.error('分组名称输入有误，请重新输入！');
      return false;
    }
    return true;
  };

  confirmMove = () => {
    const { selectedRowKeys, queryParams } = this.state;
    const newGroupName = this.props.form.getFieldValue('name');
    this.props.form.validateFields((err) => {
      if (!err && this.checkGroupName(newGroupName)) {
        addHostService
          .confirmMoveHost({
            aid: selectedRowKeys,
            to_group: newGroupName,
          })
          .then((res: any) => {
            if (res.data.code === 0) {
              message.success('移动分组成功');
              this.props.form.resetFields();
              this.setState(
                {
                  hostMoveModal: false,
                  isNewGroup: false,
                  activeKey: newGroupName,
                  selectedRowKeys: [],
                  queryParams: Object.assign(queryParams, {
                    hostGroupName: newGroupName,
                  }),
                },
                () => {
                  this.getHostGroups();
                  this.loadData(); // 返回当前组信息
                }
              );
            } else {
              message.error(res.data.msg);
            }
          });
      }
    });
  };

  /**
   * 主机接入成功之后刷新相应activeKey以及group, hostList
   * @param activeKey
   */

  refresh = (params: any) => {
    this.setState(
      {
        queryParams: Object.assign(this.state.queryParams, {
          hostGroupName: params,
        }),
        activeKey: params,
      },
      () => {
        // this.getHostGroups();
        this.loadData();
      }
    );
  };

  /**
   * 切换panel
   */
  handleClickPanel = (key) => {
    const { queryParams } = this.state;
    this.setState(
      {
        selectedRowKeys: [],
        queryParams: Object.assign(queryParams, {
          'sort-by': 'id',
          'sort-dir': 'desc',
          q_host: '',
          hostGroupName: key,
          limit: pageSize,
          start: 0,
          currentPage: 1,
          is_running: '',
          status: '',
          roles: '',
        }),
        activeKey: key,
      },
      () => {
        this.loadData();
      }
    );
  };

  onSelectChange = (selectedRowKeys, selectedRows) => {
    this.setState({
      selectedRowKeys,
      selectedRows,
    });
  };

  onTableChange = (pagination, filters, sorter) => {
    const { clusterInfo } = this.props;
    const queryParams = Object.assign(this.state.queryParams, {
      start: (pagination.current - 1) * pageSize,
      currentPage: pagination.current,
    });
    if (sorter) {
      const { field, order } = sorter;
      switch (field) {
        case 'ip': {
          queryParams['sort-by'] = 'ip';
          queryParams['sort-dir'] = order === 'descend' ? 'desc' : 'asc';
          break;
        }
        case 'updated': {
          queryParams['sort-by'] = 'updated';
          queryParams['sort-dir'] = order === 'descend' ? 'desc' : 'asc';
          break;
        }
        case 'created': {
          queryParams['sort-by'] = 'created';
          queryParams['sort-dir'] = order === 'descend' ? 'desc' : 'asc';
          break;
        }
      }
    }
    if (Object.keys(filters).length) {
      queryParams.is_running = filters.is_running || [];
      queryParams.status = filters.errorMsg || [];
      queryParams.roles =
        clusterInfo.type === 'hosts'
          ? undefined
          : (filters.roles || []).join(',');
    }

    this.setState(
      {
        queryParams,
        selectedRowKeys: [],
      },
      () => {
        this.loadData();
      }
    );
  };

  handleChangeGroup = (value: any) => {
    const { setFieldsValue } = this.props.form;
    if (value === 'NEWGROUP') {
      this.setState(
        {
          isNewGroup: true,
        },
        () => {
          setFieldsValue({ name: '' });
        }
      );
    }
  };

  render() {
    const { clusterHostGroupList, clusterInfo } = this.props;
    const {
      selectedRowKeys,
      queryParams,
      isNewGroup,
      clusterHostList,
      visibleInfo,
    } = this.state;
    const { getFieldDecorator } = this.props.form;
    const hasSelected = selectedRowKeys.length > 0;
    const rowSelection = {
      selectedRowKeys,
      onChange: this.onSelectChange,
    };
    const pagination = {
      size: 'small',
      current: queryParams.currentPage,
      pageSize: pageSize,
      total: this.state.total,
    };

    console.log(clusterHostList);
    return (
      <div className="host-page">
        <div className="host-option">
          <Search
            placeholder="按主机ip或名称搜索"
            style={{ width: 264 }}
            onSearch={this.handleSearchByHost}
          />
          <Button
            type="primary"
            style={{ float: 'right' }}
            onClick={this.goToAdd}>
            添加主机
          </Button>
        </div>
        {hasSelected ? (
          <Alert
            type="info"
            showIcon
            message={`已选择${selectedRowKeys.length}项`}
          />
        ) : null}
        {/* 折叠面板区域 */}
        {Array.isArray(clusterHostGroupList) && clusterHostGroupList.length ? (
          <>
            <Collapse
              // defaultActiveKey={[`${clusterHostGroupList[0]}`]}
              className="host-collapse"
              activeKey={this.state.activeKey}
              accordion
              onChange={this.handleClickPanel}
              style={{ borderRadius: '0', marginTop: '8px' }}>
              {clusterHostGroupList.map((item) => {
                return (
                  <Panel
                    style={{ borderRadius: '0' }}
                    header={
                      <div
                        onClick={(e) => {
                          e.stopPropagation();
                        }}
                        className="panel-table">
                        <Table
                          rowKey={(record, index) => index + ''}
                          showHeader={false}
                          className="select-table"
                          bordered={false}
                          pagination={false}
                          dataSource={[
                            {
                              name: item,
                            },
                          ]}
                          columns={[
                            {
                              title: 'name',
                              dataIndex: 'name',
                              render: (text, record) => (
                                <EditableCell
                                  isShowEditIcon={true}
                                  value={text}
                                  onChange={this.onCellChange(text)}
                                />
                              ),
                            },
                          ]}
                        />
                      </div>
                    }
                    key={`${item}`}>
                    <Table
                      rowKey="id"
                      size="middle"
                      className={classnames('dt-pagination-lower', {
                        'has-footer': hasSelected,
                      })}
                      // scroll={{ x: 1600, y: false }}
                      columns={this.initColumns()}
                      dataSource={clusterHostList}
                      rowSelection={rowSelection}
                      onChange={this.onTableChange}
                      pagination={pagination}
                      footer={hasSelected ? this.tableFooter : null}
                      loading={this.state.tableLoading}
                    />
                  </Panel>
                );
              })}
            </Collapse>
            {clusterInfo.type === 'kubernetes' && (
              <FileLogShow
                wsUrl={`ws://${window.location.host}/api/v2/cluster/kubernetes/${clusterInfo.id}/installLog`}
              />
            )}
          </>
        ) : (
          <Empty style={{ marginTop: '10%' }} />
        )}
        <Modal
          destroyOnClose={true}
          visible={this.state.showModal}
          title="添加主机"
          onCancel={() => {
            this.setState({ showModal: false });
          }}
          // onOk={this.handleAddHost}
          footer={null}>
          <CreateHostComp
            afterInstall={this.afterHostInstall}
            // refreshHost={this.loadData}
            // refreshGroup={this.getHostGroups}
            refresh={this.refresh}
            clusterInfo={clusterInfo}
            onCancel={() => {
              this.setState({ showModal: false });
            }}
          />
        </Modal>

        <Modal
          visible={this.state.hostMoveModal}
          width="650px"
          title="移动主机至"
          onCancel={() => {
            this.props.form.resetFields();
            this.setState({ hostMoveModal: false, isNewGroup: false });
          }}
          onOk={this.confirmMove}>
          <Form>
            <FormItem {...formItemCenterLayout} label="主机分组">
              {getFieldDecorator('name', {
                rules: [
                  {
                    required: true,
                    message: '请选择或输入主机分组',
                  },
                ],
              })(
                isNewGroup ? (
                  <Input
                    placeholder="请输入新建分组名称"
                    style={{ minWidth: 200 }}
                  />
                ) : (
                  <Select
                    placeholder="请选择主机分组"
                    style={{ minWidth: 200 }}
                    onChange={this.handleChangeGroup}>
                    <Option value="NEWGROUP">新建分组</Option>
                    {clusterHostGroupList &&
                      clusterHostGroupList.map((item: any) => {
                        return (
                          <Option value={item} key={item}>
                            {item}
                          </Option>
                        );
                      })}
                  </Select>
                )
              )}
            </FormItem>
          </Form>
        </Modal>
        {visibleInfo.visible && (
          <CommandPopbox {...visibleInfo} onColse={this.handleEvent} />
        )}
      </div>
    );
  }
}

export default Form.create()(HostPage);
