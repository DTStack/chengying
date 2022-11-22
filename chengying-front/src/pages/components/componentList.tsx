import * as React from 'react';
import { connect } from 'react-redux';
import {
  Table,
  Select,
  Modal,
  message,
  Icon,
  Popconfirm,
  Divider,
  Tooltip,
} from 'antd';
import { bindActionCreators, Dispatch } from 'redux';
import { Service, deployService } from '@/services';
import * as DeployAction from '@/actions/deployAction';
import { AppStoreTypes } from '@/stores';
import { deployStatusFilter } from '@/constants/const';
import DeployProgressModal from './deployProgressModal';
import ComponentConfigModal from './componentConfigModal';
import { QueryParams } from './container';
import UpgradeModal from './components/upgradeModal';
import utils from '@/utils/utils';
const Option = Select.Option;

interface Props extends QueryParams {
  history?: any;
  location?: any;
  unDeployActions?: any;
  clusterList: any[];
  getParentClustersList: Function;
  authorityList?: any;
  shouldNameSpaceShow: boolean;
  mode: number;
  deployProp?: any;
}
interface State {
  /** 组件数据 */
  componentData: {
    list: any[];
    count: number;
  };
  loading: boolean;
  searchParam: QueryParams;
  isShowUnDeploy: boolean; // 卸载弹框
  progressModalType: string;
  unDeployRecord: any;
  activeKey: any;
  selectedRecord: any;
  modalStatus: string;
  showModal: boolean;
  currentPage: number;
  upgradeModalVisible: boolean; // 升级弹窗
  modalType: string;
  versionLists: any[];
  record: any; // 当前产品
  expandedKeys: any[];
}

interface VersionType {
  id: string;
  product_version: string;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
  deployProp: state.DeployStore,
});

const mapDispatchToProps = (dispatch: Dispatch) => ({
  unDeployActions: bindActionCreators(
    Object.assign({}, DeployAction),
    dispatch
  ),
});

@(connect(mapStateToProps, mapDispatchToProps) as any)
class ProductPackage extends React.Component<Props, State> {
  state: State = {
    componentData: {
      list: [],
      count: 0,
    },
    loading: false,
    searchParam: {
      productName: undefined,
      parentProductName: undefined,
      clusterId: undefined,
      productVersion: undefined,
      deploy_status: '',
      'sort-by': 'deploy_time',
      'sort-dir': 'desc',
      limit: 10,
      start: 0,
    },
    isShowUnDeploy: false,
    progressModalType: '',
    unDeployRecord: '',
    activeKey: '',
    modalStatus: '',
    selectedRecord: undefined,
    currentPage: 1,
    showModal: false,
    upgradeModalVisible: false,
    modalType: '',
    versionLists: [],
    record: null,
    expandedKeys: [],
  };

  static getDerivedStateFromProps(props, state) {
    const { parentProductName, clusterId, productVersion, productName } =
      state.searchParam;
    if (
      props.clusterId !== clusterId ||
      props.parentProductName !== parentProductName ||
      props.productVersion !== productVersion ||
      props.productName !== productName
    ) {
      return {
        currentPage: 1,
        searchParam: Object.assign({}, state.searchParam, {
          clusterId: props.clusterId,
          parentProductName: props.parentProductName,
          productVersion: props.productVersion,
          productName: props.productName,
          start: 0,
        }),
      };
    }
    return null;
  }

  componentDidMount() {
    this.getDataList();
  }

  // 权限不足消息提醒
  errorMsg = () => {
    message.error('权限不足，请联系管理员！');
  };

  getDataList = () => {
    let isSearch: boolean = false;
    const reqParams: any = Object.assign({}, this.state.searchParam);
    if (!reqParams.parentProductName) {
      return;
    }

    reqParams.start = reqParams.start * reqParams.limit;
    if (reqParams.deploy_status) {
      isSearch = true;
      reqParams.deploy_status = reqParams.deploy_status.join(',');
    }
    if (reqParams.productName) {
      isSearch = true;
      reqParams.productName = reqParams.productName.join(',');
    }
    if (reqParams.productVersion) {
      isSearch = true;
    }
    reqParams.mode = this.props.mode;
    if (reqParams.mode) {
      isSearch = true;
    }
    this.setState({
      loading: true,
    });
    Service.getAllProducts(reqParams).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const data = res.data;
        let arr = [];
        res.data.list.map((item, index) => {
          item.key = `tr_${item.id}${index}`;
          if (isSearch) {
            arr.push(item.key);
          }
          if (item?.smooth_upgrade_product) {
            if (
              !item?.can_rollback &&
              !item?.can_upgrade &&
              item?.status !== 'deploying'
            ) {
              item.isShowNull = true;
            }
            let children = [];
            children.push(item?.smooth_upgrade_product);
            item.children = children;
            if (
              !(
                item?.smooth_upgrade_product?.can_rollback ||
                item?.smooth_upgrade_product?.can_upgrade ||
                item?.smooth_upgrade_product?.status == 'deploying'
              )
            ) {
              children[0].isShowNull = true;
            }
            children[0].isChildren = true;
            children[0].key = `tr_${children[0].id}`;
            return item;
          }
        });
        this.setState({
          expandedKeys: arr,
          componentData: data,
        });
      } else {
        message.error(res.msg);
      }
      this.setState({
        loading: false,
      });
    });
  };

  renderOption = (params: any[]) => {
    return params.map((o: any, i: number) => {
      return (
        <Option key={i + ''} value={o}>
          {o}
        </Option>
      );
    });
  };

  handleTableChange = (pagination, filters, sorter) => {
    const searchParam = Object.assign(this.state.searchParam);
    if (filters.type) {
      searchParam.type = filters.type[0];
    }
    searchParam.start = pagination.current - 1;
    searchParam.deploy_status = filters.status;
    if (filters.product_type) {
      searchParam.product_type = filters.product_type
        .map((item) => +item)
        .join(',');
    }
    if (sorter) {
      const { field, order } = sorter;
      if (order) {
        searchParam['sort-dir'] = order === 'descend' ? 'desc' : 'asc';
      }
      if (field) {
        searchParam['sort-by'] = field;
      }
    }
    this.setState({ searchParam, currentPage: pagination.current }, () => {
      this.getDataList();
    });
  };

  closeUnDeployModal = () => {
    const {
      componentData: { list },
    } = this.state;
    const { clusterId, parentProductName, getParentClustersList } = this.props;
    this.setState({
      isShowUnDeploy: false,
    });
    setTimeout(() => {
      if (list.length > 1) {
        this.getDataList();
      } else {
        getParentClustersList(clusterId, parentProductName);
      }
    }, 2000);
  };

  // 停止卸载
  stopUndeploy = (record: any, stopDelete) => {
    if (!stopDelete) {
      this.errorMsg();
      return;
    }
    Modal.confirm({
      title: '确认要停止卸载吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: async () => {
        const response = await deployService.stopUndeploy({
          clusterId: this.props.clusterId,
          namespace: record.namespace ?? '',
          pid: record.id,
        });
        const res = response.data;
        if (res.code === 0) {
          message.success('停止卸载成功');
          this.getDataList();
        } else {
          message.error(res.msg);
        }
      },
    });
  };

  // 产品升级弹窗取消
  handleUpgradeModalCancel = () => {
    this.setState({
      upgradeModalVisible: false,
      record: null,
      versionLists: [],
      modalType: '',
    });
  };

  // 跳转部署向导
  jumpToGuide = (version: string) => {
    const { record, versionLists } = this.state;
    const versions: VersionType =
      versionLists.find(
        (item: VersionType) => item.product_version === version
      ) || {};
    const url = this.initUrl(record, versions);
    utils.setNaviKey('menu_ops_center', 'sub_menu_cluster_list');
    this.props.history.push(url);
    this.handleUpgradeModalCancel();
  };

  handleOk = (param: any) => {
    const { modalType, record } = this.state;
    const { upgradeType } = this.props.deployProp;
    // if (upgradeType === 'smooth') {
    //   Object.assign(params, {upgrade_mode: 'smooth'})
    // }
    if (modalType === 'upgrade') {
      if (upgradeType === 'smooth') {
        this.props.unDeployActions.saveForcedUpgrade(record?.upgrade_service);
        this.props.unDeployActions.getIsFirstSmooth(
          record?.can_smooth_upgrade && !record?.smooth_upgrade_product
            ? true
            : false
        );
      } else {
        this.props.unDeployActions.saveForcedUpgrade([]);
        this.props.unDeployActions.getIsFirstSmooth(false);
      }
      this.jumpToGuide(param);
    } else {
      this.handleRollback(param);
    }
  };

  // 回滚
  handleRollback = async (param: any) => {
    const { record } = this.state;
    const response = await deployService.handleRollBack(
      {
        productName: record?.product_name,
      },
      param
    );
    const res = response.data;
    if (res.code === 0) {
      this.setState(
        {
          isShowUnDeploy: true,
          progressModalType: 'rollback',
          unDeployRecord: record,
        },
        () => {
          this.props.unDeployActions.getUndeploy({
            deploy_uuid: res.data.deploy_uuid,
            autoRefresh: true, // 开启自动刷新
            complete: 'deploying',
          });
          this.getDataList();
          this.handleUpgradeModalCancel();
        }
      );
    } else {
      message.error(res.msg);
    }
  };

  // 改变升级模式
  changeUpgradeType = (type: string) => {
    const { record } = this.state;
    this.props.unDeployActions.getUpgradeType(type);
    this.getProductVersionList(record, type);
  };

  // 产品升级
  handleProductUpgrade = (
    record: any,
    e: React.MouseEvent<HTMLAnchorElement, MouseEvent>
  ): any => {
    e.preventDefault();
    this.props.unDeployActions.getUpgradeType(
      record?.smooth_upgrade_product ? 'smooth' : ''
    );
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'package_upgrade')) {
      return false;
    }
    this.getProductVersionList(record);
    this.setState({
      upgradeModalVisible: true,
      modalType: 'upgrade',
      record,
    });
  };

  // 回滚弹窗
  handleProductRollback = (
    record: any,
    e: React.MouseEvent<HTMLAnchorElement, MouseEvent>
  ): any => {
    e.preventDefault();
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'package_upgrade')) {
      return false;
    }
    this.setState({
      upgradeModalVisible: true,
      modalType: 'rollback',
      record,
    });
  };

  // 获取产品包版本列表
  getProductVersionList = async (record: any, type?: string) => {
    let params = {
      product_name: record.product_name,
      product_version: record.product_version,
    };
    if (type === 'smooth' || record?.smooth_upgrade_product) {
      Object.assign(params, { upgrade_mode: 'smooth' });
    }
    try {
      const response = await Service.getProductVersionList(params);
      const { code, data, msg } = response.data;
      if (code === 0) {
        this.setState({ versionLists: data?.list || [] });
      } else {
        message.error(msg);
      }
    } catch (err) {
      message.error(err);
    }
  };

  // 初始化跳转路由
  initUrl = (record: any = {}, versions?: VersionType) => {
    const { clusterList = [], clusterId } = this.props;
    const cluster =
      clusterList.find((item) => item.clusterId === clusterId) || {};
    let url =
      '/deploycenter/appmanage/installs?type=' +
      cluster.clusterType +
      '&product_version=' +
      (versions?.product_version || record.product_version) +
      '&product_name=' +
      record.product_name +
      '&id=' +
      (versions?.id || record.id) +
      '&from=' +
      this.props.location.pathname;
    if (record.namespace) {
      url += '&namespace=' + record.namespace;
    }
    if (versions?.product_version) {
      // 新增
      url =
        url +
        '&new_version=' +
        versions.product_version +
        '&redeploy=' +
        cluster.clusterId;
    }
    return url;
  };

  handleLink = (url, txt) => {
    const { history, authorityList } = this.props;
    const CAN_REDEPLOY = authorityList.menu_deploy_guide;
    const CAN_PROGRESS = authorityList.sub_menu_service_view_progress;
    if (txt === '重新部署' || txt === '部署') {
      this.props.unDeployActions.saveForcedUpgrade([]);
      this.props.unDeployActions.getIsFirstSmooth(false);
      this.props.unDeployActions.getUpgradeType('');
      if (!CAN_REDEPLOY) {
        this.errorMsg();
        return;
      }
    }
    if (txt === '查看进度') {
      if (!CAN_PROGRESS) {
        this.errorMsg();
        return;
      }
    }
    history.push(url);
  };

  getTableColumns = () => {
    const {
      authorityList,
      HeaderStore: { cur_parent_cluster },
    } = this.props;
    const CAN_VIEW = authorityList.installed_app_view;
    const CAN_STOP_DELETE = authorityList.sub_menu_service_stop_delete;
    const CAN_DELETE = authorityList.sub_menu_component_delete;
    const CAN_FORCE_UNINSTALL = authorityList.sub_menu_service_stop_force;
    const CAN_FORCE_STOP = authorityList.sub_menu_service_delete_force;
    const tableCol = [
      {
        title: '组件名称',
        dataIndex: 'product_name_display',
        width: '10%',
        onCell: () => {
          return {
            style: {
              maxWidth: 160,
              overflow: 'hidden',
              whiteSpace: 'nowrap',
              textOverflow: 'ellipsis',
              cursor: 'pointer',
            },
          };
        },
        render(product_name_display: string) {
          return (
            <Tooltip placement="topLeft" title={product_name_display}>
              <span>{product_name_display}</span>
            </Tooltip>
          );
        },
      },
      {
        title: '组件版本号',
        dataIndex: 'product_version',
        key: 'product_version',
        width: '15%',
        sorter: true,
        onCell: () => {
          return {
            style: {
              maxWidth: 160,
              overflow: 'hidden',
              whiteSpace: 'nowrap',
              textOverflow: 'ellipsis',
              cursor: 'pointer',
            },
          };
        },
        render(productVersion: string, record: any) {
          return (
            <Tooltip
              placement="topLeft"
              title={() => (
                <span>
                  {productVersion}
                  {record.is_current_version === 1 && (
                    <Icon style={{ marginLeft: 3 }} type="star" />
                  )}
                </span>
              )}>
              <span>
                {productVersion}
                {record.is_current_version === 1 && (
                  <Icon style={{ marginLeft: 3 }} type="star" />
                )}
              </span>
            </Tooltip>
          );
        },
      },
      {
        title: '命名空间',
        dataIndex: 'namespace',
        key: 'namespace',
        render: (text) => text || '--',
      },
      {
        title: '安装包类型',
        dataIndex: 'product_type',
        key: 'product_type',
        filters: [
          { text: '传统包', value: '0' },
          { text: 'Kubernetes包', value: '1' },
        ],
        render: (text: number) => (
          <span>{text ? 'Kubernetes包' : '传统包'}</span>
        ),
      },
      {
        title: '部署状态',
        dataIndex: 'status',
        filters: deployStatusFilter,
        render: (e: string, record: any) => {
          let state: React.ReactNode = '';
          switch (e) {
            case 'deploying':
              state = (
                <span className="deploy-status-orange">
                  <Icon
                    style={{ fontSize: 12, color: '#FFB310', marginRight: 6 }}
                    type="exclamation-circle"
                    theme="filled"
                  />
                  {'部署中'}
                </span>
              );
              break;
            case 'deployed':
              state = (
                <span className="deploy-status-green">
                  <Icon
                    style={{ fontSize: 12, color: '#12BC6A', marginRight: 6 }}
                    type="check-circle"
                    theme="filled"
                  />
                  {'部署成功'}
                </span>
              );
              break;
            case 'deploy fail':
              state = (
                <span className="deploy-status-red">
                  <Icon
                    style={{ fontSize: 12, color: '#FF5F5C', marginRight: 6 }}
                    type="close-circle"
                    theme="filled"
                  />
                  {'部署失败'}
                </span>
              );
              break;
            case 'undeploying':
              state = (
                <span className="deploy-status-orange">
                  <Icon
                    style={{ fontSize: 12, color: '#FFB310', marginRight: 6 }}
                    type="exclamation-circle"
                    theme="filled"
                  />
                  {'卸载中'}
                </span>
              );
              break;
            case 'undeploy fail':
              state = (
                <span className="deploy-status-red">
                  <Icon
                    style={{ fontSize: 12, color: '#FF5F5C', marginRight: 6 }}
                    type="close-circle"
                    theme="filled"
                  />
                  {'卸载失败'}
                </span>
              );
              break;
          }
          return <span>{state}</span>;
        },
      },
      {
        title: '部署时间',
        key: 'deploy_time',
        dataIndex: 'deploy_time',
        sorter: true,
      },
      {
        title: '部署人',
        dataIndex: 'username',
        key: 'username',
        render: (text: string, record: any) => text || '--',
      },
      {
        title: '查看',
        dataIndex: 'visit',
        render: (text: string, record: any) => (
          <React.Fragment>
            {CAN_VIEW ? (
              <a onClick={() => this.setState({ selectedRecord: record })}>
                配置
              </a>
            ) : (
              '--'
            )}
          </React.Fragment>
        ),
      },
      {
        title: '操作',
        width: 200,
        dataIndex: 'action',
        render: (text: string, record: any) => {
          const { clusterId } = this.props;
          let isUndeploy = false;
          let txt = '部署';
          if (record?.smooth_upgrade_product || record?.isChildren) {
            txt = '';
          }
          let url = this.initUrl(record);
          switch (record.status) {
            case 'deploying':
              isUndeploy = false;
              txt = '查看进度';
              url +=
                '&query_str=' + record.deploy_uuid + '&cluster_id=' + clusterId;
              break;

            case 'deployed':
            case 'deploy fail':
            case 'undeploy fail':
              if (record?.smooth_upgrade_product || record?.isChildren) {
                break;
              }
              txt = '重新部署';
              isUndeploy = true;
              url += `&redeploy=${this.props.clusterId}`;
              utils.setNaviKey('menu_deploy_center', 'sub_menu_cluster_list');
              break;
            case 'undeploying':
              txt = '查看进度';
              isUndeploy = false;
              break;
          }
          return (
            <span>
              {record.can_upgrade && (
                <React.Fragment>
                  <a onClick={this.handleProductUpgrade.bind(this, record)}>
                    升级
                  </a>
                  <Divider type="vertical" />
                </React.Fragment>
              )}
              {record?.isShowNull && (
                <span style={{ color: '#3F87FF' }}>--</span>
              )}
              {record?.can_rollback && (
                <>
                  <a onClick={this.handleProductRollback.bind(this, record)}>
                    回滚
                  </a>
                  <Divider type="vertical" />
                </>
              )}
              {record.status !== 'undeploying' ? (
                <a onClick={() => this.handleLink(url, txt)}>{txt}</a>
              ) : null}
              {record.status === 'undeploying' && (
                <React.Fragment>
                  <a
                    onClick={() => {
                      if (!CAN_FORCE_UNINSTALL || !CAN_FORCE_STOP) {
                        this.errorMsg();
                        return;
                      }
                      this.setState(
                        {
                          progressModalType: 'unDeploy',
                          unDeployRecord: record,
                          isShowUnDeploy: true,
                        },
                        () => {
                          this.props.unDeployActions.getUndeploy({
                            deploy_uuid: record.deploy_uuid,
                            autoRefresh: true, // 开启自动刷新
                          });
                        }
                      );
                    }}>
                    {txt}
                  </a>
                  <Divider type="vertical" />
                  <a
                    onClick={this.stopUndeploy.bind(
                      this,
                      record,
                      CAN_STOP_DELETE
                    )}>
                    停止卸载
                  </a>
                </React.Fragment>
              )}
              {isUndeploy &&
                CAN_DELETE &&
                (!record.smooth_upgrade_product || !record?.isChildren) && (
                  <Popconfirm
                    placement="left"
                    title={
                      <div style={{ width: 240 }}>
                        卸载后该安装包下的服务将全部删除，请确认组件及所在集群信息，谨慎操作！
                        <p>集群：{cur_parent_cluster?.name}</p>
                        <p>组件：{record?.product_name}</p>
                      </div>
                    }
                    okText="卸载"
                    cancelText="取消"
                    onCancel={() => {
                      this.closeUnDeployModal();
                    }}
                    onConfirm={() => {
                      this.setState(
                        {
                          progressModalType: 'unDeploy',
                          unDeployRecord: record,
                          isShowUnDeploy: true,
                        },
                        () => {
                          // 开始卸载，在用户未操作情况下轮询list接口
                          this.props.unDeployActions.startUnDeployService(
                            {
                              product_name: record.product_name,
                              product_version: record.product_version,
                              clusterId: this.props.clusterId,
                              namespace: record.product_type
                                ? record.namespace
                                : undefined,
                            },
                            this.getDataList
                          );
                        }
                      );
                    }}>
                    <span>
                      <Divider type="vertical" />
                      <a>卸载</a>
                    </span>
                  </Popconfirm>
                )}
              {isUndeploy &&
                !CAN_DELETE &&
                (!record?.smooth_upgrade_product || !record?.isChildren) && (
                  <span
                    onClick={() => {
                      this.errorMsg();
                    }}>
                    <Divider type="vertical" />
                    <a>卸载</a>
                  </span>
                )}
            </span>
          );
        },
      },
    ];
    if (!this.props.shouldNameSpaceShow) {
      tableCol.splice(2, 1);
    }
    return tableCol;
  };

  onExpand = (expanded, record) => {
    const { expandedKeys } = this.state;
    if (expanded) {
      this.setState({ expandedKeys: [...expandedKeys, record.key] });
    } else {
      let arr = expandedKeys?.filter((v) => {
        return v !== record.key;
      });
      this.setState({ expandedKeys: arr });
    }
  };

  render = () => {
    const {
      componentData,
      loading,
      searchParam,
      currentPage,
      progressModalType,
      isShowUnDeploy,
      unDeployRecord,
      selectedRecord,
      upgradeModalVisible,
      modalType,
      versionLists,
      record,
      expandedKeys,
    } = this.state;
    const { deployProp } = this.props;
    const pagination = {
      size: 'small',
      pageSize: searchParam.limit,
      total: componentData.count,
      current: currentPage,
      showTotal: (total) => (
        <span>
          共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示
          {searchParam.limit}条
        </span>
      ),
    };
    return (
      <React.Fragment>
        <Table
          onExpand={this.onExpand}
          expandedRowKeys={expandedKeys}
          className="dt-table-fixed-base"
          style={{ height: 'calc(100vh - 260px)' }}
          columns={this.getTableColumns()}
          dataSource={componentData.list}
          scroll={{ y: true }}
          pagination={pagination}
          loading={loading}
          onChange={this.handleTableChange}
        />
        {isShowUnDeploy && (
          <DeployProgressModal
            visible={isShowUnDeploy}
            progressType={progressModalType}
            deployReacord={unDeployRecord}
            getDataList={this.getDataList}
            onClose={this.closeUnDeployModal}
          />
        )}

        <ComponentConfigModal
          componentData={selectedRecord}
          onClose={() => this.setState({ selectedRecord: undefined })}
        />
        {upgradeModalVisible && (
          <UpgradeModal
            isFirstSmooth={deployProp.isFirstSmooth}
            visible={upgradeModalVisible}
            type={modalType}
            options={versionLists}
            clusterId={this.props.clusterId}
            record={record}
            onOk={this.handleOk}
            changeUpgradeType={this.changeUpgradeType}
            onCancel={this.handleUpgradeModalCancel}
          />
        )}
      </React.Fragment>
    );
  };
}

export default ProductPackage;
