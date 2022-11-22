import * as React from 'react';
import { connect } from 'react-redux';
import { Table, Select, Modal, message, Icon, Divider } from 'antd';
import { bindActionCreators, Dispatch } from 'redux';
import { AppStoreTypes } from '@/stores';
import installGuideService from '@/services/installGuideService';
import * as DeployAction from '@/actions/deployAction';
import ComponentConfigModal from '@/pages/components/componentConfigModal';
import Utils from '@/utils/utils';

import { QueryParams } from './productPackage';

const Option = Select.Option;

interface Props extends QueryParams {
  location: any;
  history?: any;
  unDeployActions?: any;
  getDataList: Function;
  getParentProductsList: Function;
  componentData: {
    list: any[];
    count: number;
  };
  authorityList?: any;
}
interface State {
  /** 组件数据 */
  loading: boolean;
  searchParam: QueryParams;
  unDeployRecord: any;
  activeKey: any;
  selectedRecord: any;
  modalStatus: string;
  showModal: boolean;
  currentPage: number;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
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
    loading: false,
    searchParam: {
      clusterId: 0,
      productName: undefined,
      parentProductName: undefined,
      productVersion: undefined,
      deploy_status: '',
      'sort-by': 'create_time',
      'sort-dir': 'desc',
      limit: 10,
      start: 0,
    },
    unDeployRecord: '',
    activeKey: '',
    modalStatus: '',
    selectedRecord: undefined,
    currentPage: 1,
    showModal: false,
  };

  static getDerivedStateFromProps(props, state) {
    const { parentProductName, productVersion, productName } =
      state.searchParam;
    if (
      props.parentProductName !== parentProductName ||
      props.productVersion !== productVersion ||
      props.productName !== productName
    ) {
      return {
        currentPage: 1,
        searchParam: Object.assign({}, state.searchParam, {
          parentProductName: props.parentProductName,
          productVersion: props.productVersion,
          productName: props.productName,
          start: 0,
        }),
      };
    }
    return null;
  }

  handleDelProduct = (e: any, canDelete) => {
    if (!canDelete) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    const {
      componentData: { list },
    } = this.props;
    Modal.confirm({
      title: '确定删除此安装包？',
      content: '删除后，再次部署时需要再次上传。',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        // 搁浅
        installGuideService
          .deleteProducyPackage({
            productName: e.product_name,
            productVersion: e.product_version,
          })
          .then((res: any) => {
            res = res.data;
            if (res.code === 0) {
              list.length > 1
                ? this.props.getDataList()
                : this.props.getParentProductsList();
            } else {
              message.error(res.msg);
            }
          });
      },
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
    this.setState({ searchParam, currentPage: pagination.current }, () =>
      this.props.getDataList(searchParam)
    );
  };

  goToDeplay = (record, canDeploy) => {
    if (!canDeploy) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.props.unDeployActions.saveForcedUpgrade([])
    this.props.unDeployActions.getIsFirstSmooth(false)
    this.props.unDeployActions.getUpgradeType('')
    const url =
      '/deploycenter/appmanage/installs?type=' +
      (record.product_type === 0 ? 'hosts' : 'kubernetes') +
      '&product_version=' +
      record.product_version +
      '&product_name=' +
      record.product_name +
      '&id=' +
      record.id +
      '&from=' +
      this.props.location.pathname;
    Utils.setNaviKey('menu_deploy_center', 'sub_menu_product_deploy');
    this.props.history.push(url);
  };

  getTableColumns = () => {
    const { authorityList } = this.props;
    const CAN_VIEW = authorityList.package_view;
    const CAN_DEPLOY = authorityList.sub_menu_product_deploy;
    const CAN_DELETE = authorityList.package_upload_delete;
    const tableCol = [
      {
        title: '组件名称',
        dataIndex: 'product_name_display',
      },
      {
        title: '组件版本号',
        dataIndex: 'product_version',
        key: 'product_version',
        sorter: true,
        render(productVersion: string, record: any) {
          return (
            <span>
              {productVersion}
              {record.is_current_version === 1 && (
                <Icon style={{ marginLeft: 3 }} type="star" />
              )}
            </span>
          );
        },
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
        title: '上传时间',
        key: 'create_time',
        dataIndex: 'create_time',
        sorter: true,
      },
      {
        title: '上传人',
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
        dataIndex: 'action',
        render: (text: string, record: any) => {
          return (
            <React.Fragment>
              <React.Fragment>
                <a onClick={() => this.goToDeplay(record, CAN_DEPLOY)}>部署</a>
                <Divider type="vertical" />
              </React.Fragment>
              {!record.is_current_version && (
                <a
                  onClick={this.handleDelProduct.bind(
                    this,
                    record,
                    CAN_DELETE
                  )}>
                  删除
                </a>
              )}
            </React.Fragment>
          );
        },
      },
    ];
    return tableCol;
  };

  render = () => {
    const { loading, searchParam, currentPage } = this.state;
    const { componentData } = this.props;
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
          rowKey="id"
          className="dt-table-fixed-base"
          style={{ height: 'calc(100vh - 136px)' }}
          columns={this.getTableColumns()}
          dataSource={componentData.list}
          scroll={{ y: true }}
          pagination={pagination}
          loading={loading}
          onChange={this.handleTableChange}
        />
        <ComponentConfigModal
          componentData={this.state.selectedRecord}
          onClose={() => this.setState({ selectedRecord: undefined })}
          shouldGetAll={true}
        />
      </React.Fragment>
    );
  };
}

export default ProductPackage;
