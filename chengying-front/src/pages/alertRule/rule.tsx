import * as React from 'react';
import { Layout, Input, Select, Table, message, Divider, Button } from 'antd';
import moment from 'moment';
import { alertRuleService, Service } from '@/services';
import { AppStoreTypes } from '@/stores';
import { HeaderStoreType } from '@/stores/modals';
import { connect } from 'react-redux';
import './style.scss';
import { isArray } from 'lodash';

const { Content } = Layout;
const Search = Input.Search;
const Option = Select.Option;

interface RuleProp {
  HeaderStore?: HeaderStoreType;
  history?: any;
  authorityList?: any;
}

interface RuleState {
  list: any[];
  query: string;
  state: string;
  cur_parent_product: string;
  cur_parent_cluster: {
    id: number;
  };
  productList: any[]; // 产品列表
  page: number;
  size: number;
  total: number;
  selectedRowKeys: any[]; //批量启停
}

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
  authorityList: state.UserCenterStore.authorityList,
});

export class AlertRulePage extends React.Component<RuleProp, RuleState> {
  constructor(props: any) {
    super(props);
  }

  state: RuleState = {
    list: [],
    query: undefined,
    state: 'all',
    cur_parent_product: undefined,
    cur_parent_cluster: {
      id: -1,
    },
    productList: [],
    page: 1,
    size: 20,
    total: 0,
    selectedRowKeys: [],
  };
  private time;
  componentDidMount() {
    const { state, query, cur_parent_product, page, size } = this.state;
    this.updateAlertRuleList({
      state: state,
      query: query,
      dashboardTag: cur_parent_product,
      page: page,
      size: size,
    });
    this.getProductList();
  }

  // 表格初始化
  getColumns = () => {
    const { authorityList } = this.props;
    return [
      {
        title: '告警名称',
        dataIndex: 'name',
        key: 'name',
        ellipsis: true,
      },
      {
        title: '状态',
        dataIndex: 'state',
        key: 'state',
        ellipsis: true,
        render: (_, record) => {
          let icon = null;
          switch (record.state) {
            case 'ok':
              icon = (
                <span>
                  <span className="color-ok status-icon"></span>
                  {record.state}
                </span>
              );
              break;
            case 'pending':
              icon = (
                <span>
                  <span className="color-pending status-icon"></span>
                  {record.state}
                </span>
              );
              break;
            case 'alerting':
              icon = (
                <span>
                  <span className="color-alerting status-icon"></span>
                  {record.state}
                </span>
              );
              break;
            case 'no_data':
              icon = (
                <span>
                  <span className="color-no-data status-icon"></span>
                  {record.state}
                </span>
              );
              break;
            case 'paused':
              icon = (
                <span>
                  <span className="color-paused status-icon"></span>
                  {record.state}
                </span>
              );
              break;
            default:
              icon = (
                <span>
                  <span className="color-no-data status-icon"></span>
                  {record.state}
                </span>
              );
          }
          return icon;
        },
      },
      {
        title: '仪表盘名称',
        dataIndex: 'dashboardSlug',
        key: 'dashboardSlug',
        ellipsis: true,
      },
      {
        title: '仪表盘标题',
        dataIndex: 'panelTitle',
        key: 'panelTitle',
        ellipsis: true,
      },
      {
        title: '更新时间',
        dataIndex: 'newStateDate',
        key: 'newStateDate',
        render: (date) => moment(date).format('YYYY-MM-DD HH:mm:ss'),
        ellipsis: true,
      },
      {
        title: '操作',
        key: 'opreation',
        width: 120,
        render: (_, record) => {
          const CAN_EDIT = authorityList.alarm_record_open_close;
          const link =
            '/deploycenter/monitoring/dashdetail?url=' +
            encodeURIComponent(record.url) +
            '&panelId=' +
            record.panelId;
          return (
            <span>
              <a
                onClick={this.handleSwitchPauseState.bind(
                  this,
                  record,
                  CAN_EDIT
                )}>
                {record.state === 'paused' ? (
                  <span>开始</span>
                ) : (
                  <span>暂停</span>
                )}
              </a>
              <Divider type="vertical" style={{ backgroundColor: '#3F87FF' }} />
              <a onClick={this.handleLink.bind(this, link, CAN_EDIT)}>
                <span>编辑</span>
              </a>
            </span>
          );
        },
      },
    ];
  };
  //告警批量启停
  onSelectChange = (selectedRowKeys) => {
    console.log('selectedRowKeys changed: ', selectedRowKeys);
    this.setState({ selectedRowKeys });
  };
  // 暂停操作
  handleSwitchPauseState = (rule: any, isEdit: boolean) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    if (rule?.state && rule?.state === 'paused') {
      rule.paused = false;
    } else if (rule?.state) {
      rule.paused = true;
    }
    alertRuleService
      .switchGrafanaAlertPause({
        alertId: isArray(rule?.selectedRowKeys)
          ? rule?.selectedRowKeys.join(',')
          : `${rule?.id}`,
        paused: rule.paused,
      })
      .then((res: any) => {
        const { state, query, cur_parent_product, page, size } = this.state;
        this.setState(
          {
            selectedRowKeys: [],
          },
          () => {
            this.updateAlertRuleList({
              state: state,
              query: query,
              dashboardTag: cur_parent_product,
              page: page,
              size: size,
            });
          }
        );
      });
  };

  handleLink = (url: string, isEdit: boolean) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.props.history.push(url);
  };

  // 获取产品列表
  getProductList = () => {
    Service.getParentProductList().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({ productList: res.data });
      }
    });
  };
  // 更新告警规则列表
  updateAlertRuleList(params: any) {
    alertRuleService.getAlertRuleList(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState(
          {
            list: Array.isArray(res.data.data) ? res.data.data : [],
            total: res.data.total,
          },
          () => {
            let is = this.state.list.some((item) => {
              return item.state === 'pending';
            });
            const { state, query, cur_parent_product, page, size } = this.state;
            clearTimeout(this.time);
            if (is) {
              this.time = setTimeout(() => {
                this.updateAlertRuleList({
                  state: state,
                  query: query,
                  dashboardTag: cur_parent_product,
                  page: page,
                  size: size,
                });
              }, 3000);
            }
          }
        );
      } else {
        message.error(res.msg);
      }
    });
  }

  onChangePage = (page: any, pageSize: any) => {
    const { state, cur_parent_product, query } = this.state;
    this.setState({
      page: page,
      size: pageSize,
    });
    this.updateAlertRuleList({
      state: state,
      query: query,
      dashboardTag: cur_parent_product,
      page: page,
      size: pageSize,
    });
  };

  // 告警规则输入搜索内容
  handleFilterAlertByQuery = (value: any) => {
    const { state, cur_parent_product, page, size } = this.state;
    this.setState({
      query: value,
    });
    this.updateAlertRuleList({
      state: state,
      query: value,
      dashboardTag: cur_parent_product,
      page: page,
      size: size,
    });
  };

  // 状态 onChange
  handleFilterAlertByState = (value: any) => {
    const { cur_parent_product, size, query } = this.state;
    this.setState({
      state: value,
      page: 1,
    });
    this.updateAlertRuleList({
      state: value,
      query: query,
      dashboardTag: cur_parent_product,
      page: 1,
      size: size,
    });
  };

  // 选择产品
  handleProductChange = (value: string) => {
    const { state, size, query } = this.state;
    this.setState({ cur_parent_product: value, page: 1 });
    this.updateAlertRuleList({
      state: state,
      query: query,
      dashboardTag: value,
      page: 1,
      size,
    });
  };
  batchStop = () => {
    const { authorityList } = this.props;
    const params = {
      selectedRowKeys: this.state.selectedRowKeys,
      paused: true,
    };
    this.handleSwitchPauseState(params, authorityList.alarm_record_open_close);
  };
  batchStart = () => {
    const { authorityList } = this.props;
    const params = {
      selectedRowKeys: this.state.selectedRowKeys,
      paused: false,
    };
    this.handleSwitchPauseState(params, authorityList.alarm_record_open_close);
  };
  render() {
    const {
      list,
      state,
      productList = [],
      size,
      total,
      page,
      selectedRowKeys,
    } = this.state;
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
    const rowSelection = {
      selectedRowKeys,
      onChange: this.onSelectChange,
    };
    return (
      <Layout style={{ height: 'calc(100vh - 140px)' }}>
        <Content>
          <div className="alert-rules-page">
            <div className="clearfix mb-12">
              <span>产品：</span>
              <Select
                className="dt-form-shadow-bg mr-20"
                style={{ width: 264 }}
                placeholder="请选择产品"
                onChange={this.handleProductChange}>
                {Array.isArray(productList) &&
                  productList.map((item: any) => (
                    <Option key={item} value={item}>
                      {item}
                    </Option>
                  ))}
              </Select>
              <Search
                className="dt-form-shadow-bg mr-20"
                style={{ width: 264 }}
                placeholder="请输入告警名称或仪表盘名称"
                onSearch={this.handleFilterAlertByQuery}
              />
              <span>状态：</span>
              <Select
                className="dt-form-shadow-bg"
                style={{ width: 264 }}
                defaultValue="all"
                value={state}
                onChange={this.handleFilterAlertByState}>
                <Option value="all">All</Option>
                <Option value="ok">OK</Option>
                <Option value="not_ok">Not OK</Option>
                <Option value="alerting">Alerting</Option>
                <Option value="no_data">No Data</Option>
                <Option value="paused">Paused</Option>
                <Option value="pending">Pending</Option>
              </Select>
            </div>
            {/* 告警规则表格 */}
            <Table
              rowKey="id"
              className="dt-table-fixed-base"
              style={{ height: 'calc(100vh - 180px' }}
              scroll={{ y: true }}
              columns={this.getColumns()}
              rowSelection={rowSelection}
              dataSource={list}
              pagination={pagination}
            />
            <div
              style={{
                position: 'absolute',
                bottom: 35,
                left: 40,
                zIndex: 1000,
              }}>
              <Button onClick={() => this.batchStop()}>批量暂停</Button>
              &emsp;&emsp;
              <Button onClick={() => this.batchStart()}>批量开始</Button>
            </div>
          </div>
        </Content>
      </Layout>
    );
  }
}

export default connect(mapStateToProps)(AlertRulePage);
