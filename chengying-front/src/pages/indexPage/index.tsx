import * as React from 'react';
import moment from 'moment';
import { bindActionCreators, Dispatch } from 'redux';
import { connect } from 'react-redux';
import * as HeaderAction from '@/actions/headerAction';
import { HeaderStateTypes } from '@/stores/headerReducer';
import { AppStoreTypes } from '@/stores';
import { servicePageService, alertRuleService } from '@/services';
import { message, Icon, List, Avatar } from 'antd';
// import {RuleListData} from '@/pages/alertRule/__mocks__/ruleList.mock';
import FullScreen from '@/components/fullScreen';
import ServiceOverview from '@/components/serviceOverview';
import utils from '@/utils/utils';
import './style.scss';

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
});
const mapDispatchToProps = (
  dispatch: Dispatch<{ type: string; payload?: any }>
) => ({
  actions: bindActionCreators(HeaderAction, dispatch),
});

interface Props {
  HeaderStore?: HeaderStateTypes;
  actions?: HeaderAction.HeaderActionTypes;
  location?: any;
  history?: any;
}

interface State {
  homeurl: string;
  playUrl: string;
  serviceList: any[];
  collapseKey: string[] | string;
  showRight: boolean;
  serviceProduct: any[];
  serviceGroup: string[];
  list: any[];
  cur_tag: number;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
class IndexPage extends React.Component<Props, State> {
  constructor(props: any) {
    super(props);
  }

  state: State = {
    homeurl: '-1',
    playUrl: '',
    serviceList: [],
    collapseKey: [],
    showRight: false,
    serviceProduct: [],
    serviceGroup: [],
    list: [],
    cur_tag: 0,
  };

  componentDidMount() {
    this.getProductHomeUrl(this.props.HeaderStore);
    this.getProductComponents();
    this.handleGetAlertRuleList({
      state: 'not_ok',
      query: '',
      dashboardTag: this.props.HeaderStore.cur_parent_product,
    });
  }

  componentDidUpdate(prevProp: Props) {
    if (
      this.props.HeaderStore.cur_parent_product !==
      prevProp.HeaderStore.cur_parent_product
    ) {
      this.getProductHomeUrl(this.props.HeaderStore);
      this.getProductComponents();
      this.handleGetAlertRuleList({
        state: 'not_ok',
        query: '',
        dashboardTag: this.props.HeaderStore.cur_parent_product,
      });
    }
  }

  handlerClickList = (e: any) => {
    const { HeaderStore } = this.props;
    utils.setNaviKey('menu_deploy_center', 'sub_menu_dashboard');
    this.props.history.push(
      `/deploycenter/monitoring/dashdetail?url=${e.url}&panelId=${e.panelId}&var-cluster=${HeaderStore.cur_parent_cluster.name}`
    );
  };

  // 根据产品名称获取对应dashboard的url，默认product_name + _index作为tag
  getProductHomeUrl = (prop: HeaderStateTypes) => {
    const self = this;
    servicePageService
      .getServiceDashPlaylists()
      .then((res: any) => {
        res = res.data;
        if (res) {
          const playlists =
            res.filter(
              (item) => item.name === `${prop.cur_parent_product}_index`
            ) || [];
          self.setState({
            homeurl: playlists.length
              ? `${window.location.protocol}//${window.location.hostname}:${window.APPCONFIG.GRAFANA_PORT}/grafana/playlists/play/${playlists[0].id}`
              : '-1',
          });
        } else {
          self.setState({ homeurl: '-1' });
        }
      })
      .catch((err: any) => {
        message.error(err);
      });
  };

  handleShow = () => {
    const { showRight } = this.state;
    if (!showRight) {
      document.getElementById('container_left').style.display = 'none';
      document.getElementById('container_right').className +=
        ' container_right_all';
    } else {
      document.getElementById('container_left').style.display = 'inline-block';
      document.getElementById('container_right').className = 'container_right';
    }
    this.setState({ showRight: !showRight });
  };

  getProductComponents = (params?: any) => {
    const parentProduct = this.props.HeaderStore.cur_parent_product;
    const serviceProduct = new Set();
    servicePageService
      .getProductList({ limit: 0, parentProductName: parentProduct })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          const data: any[] = res.data.list||[];
          data.forEach((item) => {
            serviceProduct.add(item.product_name);
          });
          this.setState(
            { serviceProduct: Array.from(serviceProduct), cur_tag: 0 },
            () => {
              this.handleTagsChange(0);
            }
          );
        } else {
          message.error(res.msg);
        }
      });
  };

  handleTagsChange = (e) => {
    const { HeaderStore } = this.props;
    const params = {
      parentProductName: HeaderStore.cur_parent_product,
      limit: 0,
    };
    this.setState({ cur_tag: Number(e), serviceList: [] }, () => {
      servicePageService.getServiceStatus(params).then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          const serviceList = res.data.list;
          const collapseKey = [];
          const group = new Set();
          let serviceGroup = [];
          if (serviceList != null && serviceList != []) {
            serviceList.forEach((item, index) => {
              Array.isArray(item.service_list) &&
                item.service_list.forEach((service) => {
                  group.add(service.group);
                });
              if (item.service_count) {
                collapseKey.push(item.product_name);
              }
            });
            serviceGroup = Array.from(group);
            this.setState({ serviceGroup });
            this.setState({ serviceList: res.data.list, collapseKey });
          } else {
            this.setState({ serviceList: [], serviceGroup, collapseKey: [] });
          }
        } else {
          message.error(res.msg);
          this.setState({ serviceList: [] });
        }
      });
    });
  };

  onTableRowClick = (e, service) => {
    // const { cur_tag, serviceProduct } = this.state
    // const parent = serviceProduct[cur_tag]
    utils.setNaviKey('menu_ops_center', 'sub_menu_service');
    this.props.history.push(
      `/opscenter/service?component=${service.product_name}&service_group=${e.group}&service=${e.service_name}`
    );
  };

  handleGetAlertRuleList = (params: any) => {
    const self = this;
    alertRuleService.getAlertsByDashId(params).then((res: any) => {
      self.setState({
        list: Array.isArray(res.data) ? res.data : [],
      });
    });
  };

  render() {
    const {
      showRight,
      serviceGroup,
      serviceList = [],
      list,
      homeurl,
    } = this.state;
    list.sort((a, b) => {
      return a.newStateDate < b.newStateDate ? 1 : -1;
    });
    return (
      <div id="index_page" className="index-page-container">
        <div className="div_fullscreen">
          <FullScreen idName="index_page" />
        </div>
        <div
          className="arrow_container"
          style={{ height: document.body.clientHeight - 88 }}>
          <div className="arrow_left" onClick={this.handleShow}>
            {showRight ? (
              <img src={require('./img/right.png')} className="arrow_img" />
            ) : (
              <img src={require('./img/left.png')} className="arrow_img" />
            )}
          </div>
        </div>
        <div id="container_left">
          <ServiceOverview
            className="c-indexpage__overview ant-collapse-no-border box-shadow-style mb-20"
            activeKey={this.state.collapseKey}
            onChange={(collapseKey) => this.setState({ collapseKey })}
            serviceList={serviceList}
            serviceGroup={serviceGroup}
            onTableRowClick={this.onTableRowClick}
          />
          <div>
            <p className="text-title-bold">Alert</p>
            <List
              itemLayout="horizontal"
              dataSource={list}
              className="page_container_list box-shadow-style"
              size="small"
              renderItem={(item) =>
                item.state === 'alerting' ? (
                  <List.Item onClick={(e) => this.handlerClickList(item)}>
                    <List.Item.Meta
                      avatar={
                        <Avatar
                          style={{
                            backgroundColor: '#fff',
                            fontSize: '20px',
                            color: 'red',
                          }}>
                          <Icon type="bell" theme="filled" />
                        </Avatar>
                      }
                      title={item.name}
                      description={
                        <p>
                          <span style={{ color: 'red', marginRight: '5px' }}>
                            {item.state.toUpperCase()}
                          </span>
                          {moment(item.newStateDate).fromNow()}
                        </p>
                      }
                    />
                  </List.Item>
                ) : (
                  <List.Item onClick={(e) => this.handlerClickList(item)}>
                    <List.Item.Meta
                      avatar={
                        <Avatar
                          style={{
                            backgroundColor: '#fff',
                            fontSize: '20px',
                            color: '#F5A841',
                          }}>
                          <Icon type="exclamation-circle" theme="filled" />
                        </Avatar>
                      }
                      title={item.name}
                      description={
                        <p>
                          <span
                            style={{ color: '#F5A841', marginRight: '5px' }}>
                            {item.state.toUpperCase()}
                          </span>{' '}
                          {moment(item.newStateDate).fromNow()}
                        </p>
                      }
                    />
                  </List.Item>
                )
              }
            />
          </div>
        </div>
        <div
          id="container_right"
          className="container_right"
          style={{ height: document.body.clientHeight + 108 }}>
          {homeurl === '-1' ? (
            <div
              style={{
                minHeight: document.body.clientHeight - 108,
                position: 'relative',
              }}>
              <div className="container_right_content">
                <p style={{ textAlign: 'center' }}>
                  <img
                    src={require('./img/file.png')}
                    style={{ width: 100, height: 100 }}
                  />
                </p>
                <p className="holder">
                  暂无Dashboard，前往“Playlists-Edit”进行
                  <a
                    href={`${window.location.protocol}//${window.location.hostname}:${window.APPCONFIG.GRAFANA_PORT}/grafana/playlists/create`}>
                    设置
                  </a>
                </p>
              </div>
            </div>
          ) : (
            <iframe
              style={{ minHeight: '100%', width: '100%' }}
              src={homeurl}
              frameBorder="0"></iframe>
          )}
        </div>
      </div>
    );
  }
}

export default IndexPage;
