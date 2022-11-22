import * as React from 'react';
import { connect } from 'react-redux';
import { Table, message, Icon, Divider } from 'antd';
import { get, cloneDeep } from 'lodash';
import { SearchDeployLogs } from '@/model/apis';
import { Service, servicePageService, deployService } from '@/services';
import { AppStoreTypes } from '@/stores';
import { deployStatusFilter } from '@/constants/const';
import FileViewModal from '@/components/fileViewModal';
import DeployShotModal from './deployShot';
import ComponentConfigModal from './componentConfigModal';
import { QueryParams } from './container';

interface Props extends QueryParams {
  location?: any;
  authorityList?: any;
  shouldNameSpaceShow: boolean;
}

interface State {
  shotPagination: any; // 快照分页
  shotRecord: any;
  shotData: any[];

  deployData: {
    list: any[];
    count: number;
  };
  loading: boolean;
  searchParam: QueryParams;
  modalStatus: string;
  selectedRecord: any;
  selectedConfigModalRecord: any;
  selectedService: string;
  showModal: boolean;
  serviceGroup: any;
  modalContent: string;
  visibleDeployLog: boolean;
  currentPage: number;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});
@(connect(mapStateToProps, undefined) as any)
class DeployHistory extends React.Component<Props, State> {
  state: State = {
    shotRecord: '',
    shotData: [],
    shotPagination: {
      limit: 10,
      start: 0,
      count: 0,
      current: 1,
      status: '',
    },
    deployData: {
      list: [],
      count: 0,
    },
    loading: false,
    searchParam: {
      clusterId: undefined,
      productName: undefined,
      parentProductName: undefined,
      productVersion: undefined,
      deploy_status: '',
      'sort-by': 'create_time',
      'sort-dir': 'desc',
      limit: 10,
      start: 0,
    },
    selectedRecord: null,
    selectedConfigModalRecord: null,
    selectedService: null,
    modalStatus: '',
    showModal: false,
    visibleDeployLog: false,
    modalContent: null,
    serviceGroup: null,
    currentPage: 1,
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
    this.getProductHistory();
  }

  getProductHistory = (params?) => {
    const ctx = this;
    const reqParams: any = Object.assign({}, this.state.searchParam);
    if (!reqParams.parentProductName) {
      return;
    }

    this.setState({
      loading: true,
    });

    // 计算分页
    reqParams.start = reqParams.start * reqParams.limit;
    // ,分隔产品名称
    if (reqParams.productName) {
      reqParams.productName = reqParams.productName.join(',');
    }
    if (reqParams.deploy_status) {
      reqParams.deploy_status = reqParams.deploy_status.join(',');
    }
    Service.getProductUpdateRecords(reqParams).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        ctx.setState({
          deployData: res.data,
        });
      }
      this.setState({
        loading: false,
      });
    });
  };

  loadServiceGroup = (params: any) => {
    servicePageService.getServiceGroup(params).then((res: any) => {
      if (res.data.code === 0) {
        this.setState({
          serviceGroup: get(res, 'data.data.groups', {}),
        });
      }
    });
  };

  getDeployLog = (params: SearchDeployLogs) => {
    deployService.searchDeployLog(params).then((res: any) => {
      if (res.data.code === 0) {
        this.setState({
          modalContent: get(res, 'data.data.result', ''),
        });
      }
    });
  };

  handleShotTableChange = (pagination, filter, sorter) => {
    const { shotRecord } = this.state;
    const shotPagination = Object.assign(this.state.shotPagination, {
      start: (pagination.current - 1) * this.state.shotPagination.limit,
      current: pagination.current,
      status: filter.status ? filter.status.join(',').toString() : '',
    });
    this.setState(
      {
        shotPagination,
      },
      () => {
        this.handleViewDeployShot(shotRecord);
      }
    );
  };

  handleDeployHistoryTableChange = (pagination, filters, sorter) => {
    const searchParam = Object.assign(this.state.searchParam);
    if (filters.status) {
      searchParam.deploy_status = filters.status;
    }
    searchParam.start = pagination.current - 1;
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
    this.setState(
      { searchParam, currentPage: pagination.current },
      this.getProductHistory
    );
  };

  handleViewDeployShot = (e: any) => {
    this.setState({
      modalStatus: e.status,
    });
    const { shotPagination } = this.state;
    Service.getDeployShot({
      uuid: e.deploy_uuid,
      limit: shotPagination.limit,
      start: shotPagination.start,
      status: shotPagination.status,
    }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({
          showModal: true,
          shotData: res.data.list,
          shotPagination: Object.assign(this.state.shotPagination, {
            count: res.data.count,
          }),
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  onShowConfigInfo = (record) => {
    this.setState({
      selectedConfigModalRecord: record,
    });
  };

  onShowDeployLog = (record) => {
    if (record && record.product_name) {
      this.setState(
        {
          selectedRecord: record,
          visibleDeployLog: true,
        },
        () => {
          this.loadServiceGroup({
            product_name: record.product_name,
          });
          this.getDeployLog({
            serviceName: '',
            deployId: record.deploy_uuid,
            productName: record.product_name,
            productVersion: record.product_version,
          });
        }
      );
    }
  };

  onSelectedLogService = (service: string) => {
    const { selectedRecord } = this.state;
    this.setState({
      selectedService: service,
    });
    this.getDeployLog({
      deployId: selectedRecord.deploy_uuid,
      serviceName: service || '',
      productName: selectedRecord.product_name,
      productVersion: selectedRecord.product_version,
    });
  };

  onCloseModal = () => {
    this.setState({
      visibleDeployLog: false,
      selectedRecord: null,
      selectedService: null,
      modalContent: null,
    });
  };

  initTableColumns = () => {
    const { authorityList } = this.props;
    const CAN_VIEW = authorityList.installed_app_view;
    const cloneStatus = cloneDeep(deployStatusFilter);
    const deployStatusFiltered = cloneStatus.map((o: any) => {
      if (o.value === 'undeployed') {
        o.text = '卸载成功';
      }
      return o;
    });
    const columns = [
      {
        title: '组件名称',
        dataIndex: 'product_name_display',
        key: 'product_name_display',
        render: (e: string, record: any) => e || record.product_name,
      },
      {
        title: '组件版本号',
        dataIndex: 'product_version',
        key: 'product_version',
        sorter: true,
        render: (productVersion: string, record: any) => (
          <span>
            {productVersion}
            {record.is_current_version === 1 && (
              <Icon style={{ marginLeft: 3 }} type="star" />
            )}
          </span>
        ),
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
        title: '部署人',
        dataIndex: 'username',
        key: 'username',
      },
      {
        title: '部署时间',
        dataIndex: 'create_time',
        key: 'create_time',
        sorter: true,
      },
      {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        filters: [
          ...deployStatusFiltered,
          {
            text: '卸载成功',
            value: 'undeployed',
          },
        ],
        render: (text: string) => {
          let state: React.ReactNode = '';
          switch (text) {
            case 'undeployed':
              state = (
                <span className="deploy-status-green">
                  <Icon
                    style={{ fontSize: 12, color: '#12BC6A', marginRight: 6 }}
                    type="check-circle"
                    theme="filled"
                  />
                  {'卸载成功'}
                </span>
              );
              break;
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
        title: '部署快照',
        dataIndex: 'action',
        render: (e: any, record: any) => (
          <React.Fragment>
            {CAN_VIEW ? (
              <Icon
                type="camera"
                style={{ cursor: 'pointer', color: '#3f87ff' }}
                title="查看部署快照"
                onClick={() => {
                  this.setState(
                    {
                      shotRecord: record,
                    },
                    () => {
                      this.handleViewDeployShot(record);
                    }
                  );
                }}
              />
            ) : (
              '--'
            )}
          </React.Fragment>
        ),
      },
      {
        title: '查看',
        dataIndex: 'visit',
        render: (text: string, record: any) => (
          <React.Fragment>
            {CAN_VIEW ? (
              <React.Fragment>
                <a onClick={this.onShowConfigInfo.bind(this, record)}>配置</a>
                <Divider type="vertical" />
                <a onClick={this.onShowDeployLog.bind(this, record)}>
                  部署日志
                </a>
              </React.Fragment>
            ) : (
              '--'
            )}
          </React.Fragment>
        ),
      },
    ];
    if (!this.props.shouldNameSpaceShow) {
      columns.splice(2, 1);
    }
    return columns;
  };

  render = () => {
    const { deployData, loading, visibleDeployLog, currentPage } = this.state;

    const pageConf = {
      size: 'small',
      pageSize: 10,
      current: currentPage,
      total: deployData.count,
      showTotal: (total) => (
        <span>
          共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示
          {this.state.shotPagination.limit}条
        </span>
      ),
    };

    return (
      <React.Fragment>
        <Table
          rowKey="id"
          className="dt-table-fixed-base"
          style={{ height: 'calc(100vh - 260px)' }}
          pagination={pageConf}
          columns={this.initTableColumns()}
          loading={loading}
          onChange={this.handleDeployHistoryTableChange}
          dataSource={deployData.list}
          scroll={{ y: true }}
        />
        <DeployShotModal
          status={this.state.modalStatus}
          dataList={this.state.shotData}
          showModal={this.state.showModal}
          closeModal={() => {
            this.setState({
              showModal: false,
              shotPagination: Object.assign(this.state.shotPagination, {
                start: 0,
                current: 1,
                status: '',
              }),
            });
          }}
          shotPagination={this.state.shotPagination}
          handleTableChange={this.handleShotTableChange.bind(this)}
        />
        <FileViewModal
          key="deployModal"
          title="部署日志"
          visible={visibleDeployLog}
          content={this.state.modalContent}
          serviceData={this.state.serviceGroup}
          onCancel={this.onCloseModal}
          onSelectedService={this.onSelectedLogService}
        />
        <ComponentConfigModal
          componentData={this.state.selectedConfigModalRecord}
          onClose={() =>
            this.setState({ selectedConfigModalRecord: undefined })
          }
        />
      </React.Fragment>
    );
  };
}

export default DeployHistory;
