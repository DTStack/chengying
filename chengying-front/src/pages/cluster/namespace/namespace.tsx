import * as React from 'react';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import {
  Input,
  Table,
  Button,
  Badge,
  Divider,
  message,
  Modal,
  Icon,
  Drawer,
} from 'antd';
import { PaginationConfig, SorterResult } from 'antd/lib/table';
import { ClusterNamespaceService, imageStoreService } from '@/services';
import ProgressBar from '@/components/progressBar';
import NamespaceModal from './components/namespaceModal';
import IpModal from './components/ipModal';
import NamespaceInfo from './components/namespaceInfo';
import './style.scss';

const Search = Input.Search;

interface IProps {
  cur_parent_cluster: any;
  authorityList: any;
}
interface IState {
  imageStoreList: any[];
  namespaceList: any[];
  reqParams: {
    desc: boolean;
    status: string;
    type: string;
    namespace: string;
  };
  tableLoading: boolean;
  namespaceVisible: boolean;
  record: any;
  ipVisible: boolean;
  drawerVisible: boolean;
}
interface StatusItemType {
  text: string;
  value: string;
  color: string;
}
const mapStateToProps = (state: AppStoreTypes) => ({
  cur_parent_cluster: state.HeaderStore.cur_parent_cluster,
  authorityList: state.UserCenterStore.authorityList,
});

@(connect(mapStateToProps, undefined) as any)
export default class NamespacePage extends React.PureComponent<IProps, IState> {
  state: IState = {
    imageStoreList: [],
    namespaceList: [],
    reqParams: {
      desc: true,
      status: undefined,
      type: undefined,
      namespace: undefined,
    },
    tableLoading: false,
    namespaceVisible: false,
    record: {},
    ipVisible: false,
    drawerVisible: false,
  };

  private timer: any = null;

  componentDidMount() {
    this.getNamespaceList();
    this.getImageStoreList();
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  // 获取命名空间列表
  getNamespaceList = async () => {
    this.setState({ tableLoading: true });
    this.getNamespaceListApi(() => {
      this.cycleNamespaceList();
    });
    this.setState({ tableLoading: false });
  };

  // 获取命名空间列表数据接口
  getNamespaceListApi = async (callback?: Function) => {
    const { reqParams } = this.state;
    const response = await ClusterNamespaceService.getNamespaceList(reqParams);
    const res = response.data;
    const { code, data, msg } = res;
    if (code === 0) {
      this.setState({
        namespaceList: data || [],
      });
      callback && callback();
    } else {
      message.error(msg);
    }
  };

  // 命名空间列表接口轮询
  cycleNamespaceList = () => {
    clearInterval(this.timer);
    this.timer = setInterval(async () => {
      await this.getNamespaceListApi();
    }, 60 * 1000);
  };

  // 获取镜像仓库列表
  getImageStoreList = () => {
    const { cur_parent_cluster } = this.props;
    imageStoreService
      .getImageStoreList({
        cluster_id: cur_parent_cluster.id,
      })
      .then((response: any) => {
        const res = response.data;
        if (res.code === 0) {
          this.setState({
            imageStoreList: res.data ? res.data.list : [],
          });
        } else {
          message.error(res.msg);
        }
      });
  };

  // 搜索
  handleNamespaceSearch = (namespace: string) => {
    const reqParams = {
      ...this.state.reqParams,
      namespace: namespace || undefined,
    };
    this.setState({ reqParams }, this.getNamespaceList);
  };

  // 添加 | 编辑 命名空间
  namespaceModalShow = (record?: any) => {
    const { authorityList } = this.props;
    const CAN_EDIT = authorityList.ns_edit;
    if (!CAN_EDIT) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.setState({
      namespaceVisible: !this.state.namespaceVisible,
      record: record || {},
    });
  };

  // 校验是否可以删除
  handleConfirmDelete = async (record: any) => {
    const response = await ClusterNamespaceService.confirmDelete({
      namespace: record.namespace,
    });
    const { code, data, msg } = response.data;
    if (code === 0) {
      if (data.status) {
        this.handleDeleteNamespace(record);
      } else {
        Modal.confirm({
          title: '该命名空间已部署产品，不支持删除',
          icon: <Icon type="exclamation-circle" theme="filled" />,
          okText: '确定',
          cancelText: '取消',
          onOk: () => {},
          onCancel: () => {},
        });
      }
    } else {
      message.error(msg);
    }
  };

  // 删除
  handleDeleteNamespace = (record: any) => {
    Modal.confirm({
      title: '确定删除该命名空间吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: async () => {
        const response = await ClusterNamespaceService.deleteNamespace({
          namespace: record.namespace,
        });
        const { code, msg } = response.data;
        if (code === 0) {
          message.success('执行成功');
          this.getNamespaceList();
        } else {
          message.error(msg);
        }
      },
    });
  };

  // IP配置
  handleIpConnect = (record?: any) => {
    this.setState({
      ipVisible: !this.state.ipVisible,
      record: record || {},
    });
  };

  // 表格分页，筛选等
  handleTableChange = (
    pagination: PaginationConfig,
    filters: Partial<Record<string | number | symbol, string[]>>,
    sorter: SorterResult<any>
  ) => {
    const reqParams = {
      ...this.state.reqParams,
    };
    if (Object.keys(filters).length > 0) {
      if (Array.isArray(filters.status)) {
        reqParams.status = filters.status.join(',') || undefined;
      }
      if (Array.isArray(filters.type)) {
        reqParams.type = filters.type.join(',') || undefined;
      }
    }
    if (sorter.field) {
      reqParams.desc = sorter.order !== 'ascend';
    }
    this.setState({ reqParams }, this.getNamespaceList);
  };

  // 保存
  handleSave = (values: any, callback?: Function) => {
    ClusterNamespaceService.saveNamespace(values).then((response: any) => {
      const res = response.data;
      const { code, msg } = res;
      callback && callback();
      if (code === 0) {
        message.success('保存成功');
        this.namespaceModalShow();
        this.getNamespaceList();
      } else {
        message.error(msg);
      }
    });
  };

  // 查看详情
  infoDrawerShow = (record?: any) => {
    this.setState({
      drawerVisible: !this.state.drawerVisible,
      record: record || {},
    });
  };

  // 表格列
  initColums = () => {
    const statusList: StatusItemType[] = [
      { text: '正常', value: 'valid', color: '#12BC6A' },
      { text: '异常', value: 'invalid', color: '#FF5F5C' },
      { text: '未连接', value: 'not_connect', color: '#FFB310' },
    ];
    const cloumns = [
      {
        title: 'namespace',
        dataIndex: 'namespace',
        key: 'namespace',
        render: (value: string, record: any) =>
          value ? (
            <a onClick={this.infoDrawerShow.bind(this, record)}>{value}</a>
          ) : (
            '--'
          ),
      },
      {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        filters: statusList,
        filterMultiple: true,
        render: (value: string) => {
          const status: StatusItemType = statusList.find(
            (item: any) => item.value === value
          );
          return <Badge color={status.color} text={status.text} />;
        },
      },
      {
        title: '添加方式',
        dataIndex: 'type',
        key: 'type',
        filters: [
          { text: 'kubeconfig', value: 'kubeconfig' },
          { text: 'agent', value: 'agent' },
        ],
        filterMultiple: true,
        render: (value: string) => value || '--',
      },
      {
        title: 'cpu使用',
        dataIndex: 'cpu',
        key: 'cpu',
        render: (value: string, record: any) => {
          return (
            <ProgressBar
              now={record.cpu_used}
              total={record.cpu_total}
              percent={record.cpu_percent}
            />
          );
        },
      },
      {
        title: '内存使用',
        dataIndex: 'memory',
        key: 'memory',
        render: (value: string, record: any) => {
          return (
            <ProgressBar
              now={record.memory_used}
              total={record.memory_total}
              percent={record.memory_percent}
            />
          );
        },
      },
      {
        title: '最近修改人',
        dataIndex: 'user',
        key: 'user',
        render: (value: string) => value || '--',
      },
      {
        title: '最近修改时间',
        dataIndex: 'update_time',
        key: 'update_time',
        sorter: true,
      },
      {
        title: '操作',
        dataIndex: 'actions',
        key: 'actions',
        render: (value: any, record: any) => {
          const { authorityList } = this.props;
          const CAN_EDIT = authorityList.ns_edit;
          if (!CAN_EDIT) {
            return '--';
          }
          return (
            <React.Fragment>
              <a onClick={this.namespaceModalShow.bind(this, record)}>编辑</a>
              <React.Fragment>
                <Divider type="vertical" />
                {record.isDeployed ? (
                  <span style={{ color: '#999' }}>删除</span>
                ) : (
                  <a onClick={this.handleConfirmDelete.bind(this, record)}>
                    删除
                  </a>
                )}
              </React.Fragment>
              {record.status === 'not_connect' && (
                <React.Fragment>
                  <Divider type="vertical" />
                  <a onClick={this.handleIpConnect.bind(this, record)}>
                    IP配置
                  </a>
                </React.Fragment>
              )}
            </React.Fragment>
          );
        },
      },
    ];
    return cloumns;
  };

  render() {
    const {
      imageStoreList,
      namespaceList,
      tableLoading,
      namespaceVisible,
      record,
      ipVisible,
      drawerVisible,
    } = this.state;
    const { authorityList } = this.props;
    return (
      <div className="cluster-page-body cluster-namespace-page">
        <div className="clearfix mb-12">
          <Search
            className="dt-form-shadow-bg"
            style={{ width: 264 }}
            placeholder="按namespace搜索"
            onSearch={this.handleNamespaceSearch}
          />
          <Button
            className="fl-r"
            type="primary"
            onClick={() => this.namespaceModalShow()}>
            添加命名空间
          </Button>
        </div>
        <Table
          rowKey="id"
          className="dt-table-fixed-base"
          style={{ height: 'calc(100vh - 230px)' }}
          scroll={{ y: true }}
          loading={tableLoading}
          columns={this.initColums()}
          dataSource={namespaceList}
          pagination={false}
          onChange={this.handleTableChange}
        />
        {namespaceVisible && (
          <NamespaceModal
            namespace={record.namespace}
            authorityList={authorityList}
            visible={namespaceVisible}
            imageStoreList={imageStoreList}
            handleSave={this.handleSave}
            handleCancel={() => this.namespaceModalShow()}
          />
        )}
        {ipVisible && (
          <IpModal
            namespace={record.namespace}
            visible={ipVisible}
            handleCancel={() => this.handleIpConnect()}
            getTableList={this.getNamespaceList}
          />
        )}
        <Drawer
          className="c-namespace-info_ant-drawer"
          title={record.namespace}
          key={record.id + record.namespace}
          // getContainer={false}
          width={'70%'}
          placement="right"
          onClose={() => this.infoDrawerShow()}
          visible={drawerVisible}>
          {drawerVisible && (
            <NamespaceInfo {...this.props} namespace={record.namespace} />
          )}
        </Drawer>
      </div>
    );
  }
}
