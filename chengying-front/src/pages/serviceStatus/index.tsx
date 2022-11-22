import { connect } from 'react-redux';
import * as ServiceActions from '@/actions/serviceAction';
import { bindActionCreators, Dispatch } from 'redux';
import * as React from 'react';
import { Layout, Modal, Icon, Tabs, message, Dropdown, Menu } from 'antd';
import Logtail from '@/components/logtail';
import { servicePageService } from '@/services';
import { AppStoreTypes } from '@/stores';
import { HeaderStoreType, ServiceStore } from '@/stores/modals';
import * as Cookie from 'js-cookie';
import ConfigParam from './configParam';
import SideNav from './sideNav';
import RuntimeStatus from './runtimeStatus';
import utils from '@/utils/utils';
import {
  ProductionListParams,
  CurrentProductionParams,
} from '@/services/ServicePageService';
import './style.scss';

const MenuItem = Menu.Item;
const { Content } = Layout;
const TabPane = Tabs.TabPane;
const NOTICE_INTERVAL: number = 5000; // EM提醒、异常服务轮询间隔5s

export interface ServiceProp {
  HeaderStore: HeaderStoreType;
  actions: any;
  ServiceStore: ServiceStore;
  authorityList: any;
  location?: any;
  history?: any;
}

interface ServiceState {
  cur_parent_product: string;
  products: any[];
  cur_product_info: {
    // 当前product, 实质上感觉和cur_product没有区别，除了init
    product_name: string;
    product_version?: string;
    product_id: number;
  };
  cur_service_info: any;
  log_modal_visible: boolean;
  logpaths: any[];
  log_service_id: number;
  activeKey: string;
  dashurl: string;
  services: any[];
  urlParams: any;
  dashId: string;
  menuKey: {
    openKeys: any[];
    selectedKeys: any[];
  };
  redService: {
    count: number;
    list: any[];
  };
  redVisible: boolean;
  noHosts: string;
  repeatParams: any[];
  // allHostList: any[];
}

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
  ServiceStore: state.ServiceStore,
  authorityList: state.UserCenterStore.authorityList,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, ServiceActions), dispatch),
});
class ServiceList extends React.Component<ServiceProp, ServiceState> {
  state: ServiceState = {
    cur_parent_product: this.props.HeaderStore.cur_parent_product || '',
    products: [],
    cur_product_info: {
      product_name: Cookie.get('em_product_name') || '',
      product_version: '',
      product_id: parseInt(Cookie.get('em_product_id'), 10) || -1,
    },
    cur_service_info: {},
    log_modal_visible: false,
    logpaths: [],
    log_service_id: -1,
    activeKey: '0',
    dashurl: '',
    services: null,
    urlParams: {},
    dashId: '',
    menuKey: {
      openKeys: ['f0'],
      selectedKeys: ['sub0_0'],
    },
    redService: {
      count: 0,
      list: [],
    },
    redVisible: false,
    noHosts: '',
    repeatParams: [],
    // allHostList: [],
  };

  private haInterval: any = null;
  private noticeInterval: any = null;
  private parentProduct: string = '';

  componentDidMount() {
    const {
      HeaderStore: { cur_parent_product },
    } = this.props;
    this.checkUrlParam(this.initServices);
    if (this.noticeInterval) {
      clearInterval(this.noticeInterval);
    }
    this.noticeInterval = setInterval(() => {
      const curParentProduct =
        cur_parent_product === '选择产品'
          ? Cookie.get('em_current_parent_product')
          : cur_parent_product;
      this.getRestartService();
      this.getRedService(curParentProduct);
    }, NOTICE_INTERVAL);
  }

  componentWillUnmount() {
    clearInterval(this.haInterval);
    clearInterval(this.noticeInterval);
  }

  componentDidUpdate(prevProps, prevState) {
    const {
      location: { search },
      HeaderStore,
      ServiceStore,
    } = this.props;
    const { cur_parent_product, cur_parent_cluster } = HeaderStore;
    const { redService } = ServiceStore;
    const { products, cur_product_info } = prevState;
    if (search !== prevProps.location.search) {
      this.checkUrlParam();
    }
    if (
      cur_parent_product !== '选择产品' &&
      cur_parent_product !== '' &&
      (this.parentProduct === '' ||
        this.parentProduct !== cur_parent_product ||
        cur_parent_cluster.id !== HeaderStore.cur_parent_cluster.id) &&
      products.length === 0
    ) {
      this.getProductList(cur_parent_product);
      this.parentProduct = cur_parent_product;
    }
    if (redService.count !== prevProps.ServiceStore.redService.count) {
      this.getCurrentProduct(
        cur_product_info.product_name,
        cur_product_info.product_version,
        true
      );
    }

    if (
      prevProps?.HeaderStore?.cur_parent_product !== cur_parent_product &&
      cur_parent_product !== '选择产品'
    ) {
      this.getRedService(cur_parent_product);
    }
  }

  initServices = () => {
    const {
      HeaderStore: { cur_parent_product },
    } = this.props;
    const { products } = this.state;
    if (
      cur_parent_product !== '选择产品' &&
      cur_parent_product !== '' &&
      this.parentProduct === '' &&
      !products.length
    ) {
      // 获取产品列表
      this.getProductList(cur_parent_product);
      this.getRedService(cur_parent_product);
      this.parentProduct = cur_parent_product;
    }
    this.getRestartService();
  };

  // 获取异常服务数量列表
  getRedService = (curParentProduct: any) => {
    servicePageService
      .getRedService({ parentProductName: curParentProduct })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.setState(
            {
              redService: {
                count: res.data.count,
                list: res.data.list,
              },
              redVisible: res.data.count > 0,
            },
            () => {
              this.props.actions.setRedService({
                count: res.data.count,
                list: res.data.list,
              });
            }
          );
        } else {
          message.error(res.msg);
        }
      });
  };

  // 获取需要依赖组件重新配置服务列表——EM提醒
  getRestartService = () => {
    servicePageService.getRestartService({}).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.props.actions.setResartServiceList({
          count: res.data.count,
          list: res.data.list,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  checkUrlParam = (cb?: any) => {
    // 注释参数为空判断，修复服务与概览主机诊断之间切换默认值问题
    const u: any = {};
    const s = this.props.location.search.slice(1);
    s.split('&').map((o: any) => {
      u[o.split('=')[0]] = o.split('=')[1];
    });
    this.setState(
      {
        urlParams: u,
      },
      () => {
        cb && cb();
      }
    );
  };

  // 获取产品包列表 DTBase...
  getProductList = (curParentProduct: any) => {
    const { cur_product_info } = this.state;
    const params: ProductionListParams = {
      parentProductName: curParentProduct,
    };
    // '/api/v2/product'
    servicePageService.getProductName(params).then((res: any) => {
      res = res.data;
      if (res.code === 0 && res.data?.list) {
        const list = res.data?.list
        if (!list.length) {
          cur_product_info.product_name = '';
          cur_product_info.product_id = -1;
        }
        this.setState(
          {
            products: list,
            cur_product_info,
          },
          () => {
            if (list.length) {
              this.getDefaultService(list);
            } else {
              Cookie.remove('em_product_id');
              Cookie.remove('em_product_name');
              this.props.actions.resetServices();
            }
          }
        );
      }
    });
  };

  // 获取默认选中展示服务
  getDefaultService = (products: any) => {
    const { cur_product_info, urlParams } = this.state;
    const { product_name, product_version } = cur_product_info;
    let pn = urlParams.component || product_name;
    let pv = product_version;
    let isExist = false;
    for (const p of products) {
      if (p.product_name === pn) {
        isExist = true;
        pv = p.product_version;
        pn = p.product_name;
        // 兼容跳转
        Cookie.set('em_product_id', p.id);
        //  pn处理这步之前没有这一步，这次加上了，因为发现切换url如果这步处理没来更改，pn保留的就是原来的url中的component值
      }
    }
    if (!isExist && products.length) {
      pn = products[0].product_name;
      pv = products[0].product_version;
      const sessionHostService = JSON.parse(
        sessionStorage.getItem('service_object')
      );
      if (sessionHostService && sessionHostService.productName) {
        pn = sessionHostService.productName;
      }
    }
    this.getCurrentProduct(pn, pv);
  };

  /**
   * 获取当前产品的信息
   * @param productName 当前产品名称
   * @param productVersion
   * @param notNeedUpdateService
   * @returns
   */
  getCurrentProduct = (
    productName: string,
    productVersion?: string,
    notNeedUpdateService?: boolean
  ) => {
    let { cur_service_info, urlParams, menuKey } = this.state;
    if (!productName) {
      return;
    }
    const params: CurrentProductionParams = { product_name: productName };
    utils.k8sNamespace && (params.namespace = utils.k8sNamespace);
    Cookie.remove('em_product_name');
    servicePageService.getCurrentProduct(params).then((res: any) => {
      const data = res.data;
      if (data.code === 0) {
        this.props.actions.setCurrentProduct(data.data);
        const serviceGroupConf: any = { cloud: true };
        utils.k8sNamespace && (serviceGroupConf.namespace = utils.k8sNamespace);
        this.props.actions.getServiceGroup(
          {
            product_name: productName,
          },
          (firstService: any, key: Array<string>, selectedKey: any) => {
            //
            if (!notNeedUpdateService) {
              menuKey.openKeys = [...key];
              menuKey.selectedKeys = [selectedKey];
              cur_service_info = firstService;
            }
            // 更新当前产品信息，所有服务，选中服务
            this.setState(
              {
                cur_product_info: {
                  product_id: data.data.id,
                  product_name: productName,
                  product_version: productVersion,
                },
                services: data.data.product.Service,
                cur_service_info,
                menuKey,
              },
              () => {
                this.handleSwitchService(
                  cur_service_info,
                  productName,
                  data.data
                );
              }
            );
          },
          urlParams.service || '',
          serviceGroupConf
        );
      } else {
        message.error(data.msg);
      }
    });
  };

  // 切换服务
  handleSwitchService = (service: any, productName?: any, curProduct?: any) => {
    const { ServiceStore, history, actions } = this.props;
    const { cur_product } = ServiceStore;
    const newProduct = curProduct || cur_product;
    let { services, dashurl, activeKey, dashId } = this.state;
    let hasDash = false; // 存在仪表盘
    // 重置选择的配置文件
    actions.setConfigFile('all');

    for (const s in services) {
      if (s === service.service_name) {
        if (services[s].Instance) {
          if (services[s].Instance.PrometheusPort !== '') {
            hasDash = true;
          }
        }
      }
    }

    if (this.haInterval) {
      clearInterval(this.haInterval);
    }

    this.setCurrentService(
      {
        ...newProduct.product.Service[service.service_name],
        configModify: {},
      },
      service.service_name
    );

    // 记录在url上
    history.replace(
      `/opscenter/service?component=${productName}&service_group=${
        (newProduct.product.Service[service.service_name] || {}).Group
      }&service=${service.service_name}`
    );
    // debugger;
    if (hasDash) {
      servicePageService
        .getServiceDashInfo({
          product_name: newProduct.product_name,
          service_name: service.service_name,
        })
        .then((res: any) => {
          res = res.data;
          if (res.length) {
            dashurl = res[0].url;
            dashId = res[0].id;
          } else {
            // activeKey = '0';
            dashurl = '';
            dashId = '';
          }
          this.setState({
            cur_service_info: service,
            dashurl,
            dashId,
            activeKey,
          });
        });
    } else {
      this.setState({
        cur_service_info: service,
        dashurl: '',
        // activeKey: '0',
        dashId: '',
      });
    }
  };

  // 获取组件列表
  getServiceGroup = () => {
    // debugger
    const { cur_product_info, cur_service_info } = this.state;
    const { product_name } = cur_product_info;
    const serviceGroupConf: any = { cloud: true };
    utils.k8sNamespace && (serviceGroupConf.namespace = utils.k8sNamespace);
    this.props.actions.getServiceGroup(
      {
        product_name,
      },
      (firstService: any, key: string[], selectedKey: any) => {
        this.handleSwitchService(cur_service_info, product_name);
      },
      '',
      serviceGroupConf
    );
  };

  handleLogModalCancel() {
    this.setState({
      log_modal_visible: false,
    });
  }

  handleTabsChange = (activeKey: any) => {
    this.setState({
      activeKey: activeKey,
    });
  };

  setCurrentService = (service, service_name) => {
    const { cur_product } = this.props.ServiceStore;
    this.props.actions.setCurrentService(service, {
      product_name: cur_product.product_name,
      service_name: service_name,
      product_version: cur_product.product_version,
    });
  };

  // error list navigator
  setMenuKey = (item) => {
    const { actions } = this.props;
    const { service_name, product_name } = item;
    // const self = this;
    servicePageService.getCurrentProduct({ product_name }).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        Cookie.set('em_product_name', data.data.product_name);
        Cookie.set('em_product_id', data.data.id);
        actions.setCurrentProduct(data.data);
        actions.getServiceGroup(
          {
            product_name,
          },
          (firstService: any, key: string[], selectedKey: any) => {
            this.setState({
              cur_product_info: {
                product_id: data.data.id,
                product_name,
              },
              services: data.data.product.Service,
              cur_service_info: firstService,
              menuKey: {
                selectedKeys: [selectedKey],
                openKeys: [...key],
              },
            });
            this.handleSwitchService(firstService, product_name);
            this.setCurrentService(
              data.data.product.Service[firstService.service_name],
              firstService.service_name
            );
          },
          service_name,
          { cloud: true }
        );
      } else {
        message.error(data.msg);
      }
    });
  };

  // 改变未关联主机情况
  changeHosts = () => {
    this.setState({
      noHosts: '',
      repeatParams: [],
    });
  };

  render() {
    const {
      ServiceStore: { restartService },
      ServiceStore,
      HeaderStore,
      authorityList,
    } = this.props;
    const {
      dashurl,
      activeKey,
      cur_service_info,
      cur_product_info,
      redService,
      redVisible,
      products,
    } = this.state;
    const CAN_DASHBOARD_VIEW = authorityList.service_dashboard_view;
    const isKubernetes: boolean =
      HeaderStore?.cur_parent_cluster?.type === 'kubernetes';
    const dashUrl: string =
      window.location.protocol +
      '//' +
      window.location.hostname +
      ':' +
      window.APPCONFIG.GRAFANA_PORT +
      dashurl +
      '?theme=light&no-feedback=true' +
      '&var-cluster=' +
      HeaderStore.cur_parent_cluster.name;
    const sidenavProps = {
      products,
      cur_product_info: cur_product_info,
      menuKey: this.state.menuKey,
      isKubernetes: isKubernetes,
      ServiceStore: ServiceStore,
      HeaderStore: HeaderStore,
      actions: this.props.actions,
      authorityList: authorityList,
    };

    const menu: React.ReactNode = (
      <Menu>
        {redService.list &&
          redService.list.map((item, index) => {
            return (
              <MenuItem key={index}>
                <a onClick={() => this.setMenuKey(item)}>
                  {item.product_name}——{item.service_name}
                </a>
              </MenuItem>
            );
          })}
      </Menu>
    );

    return (
      <div className="service-container">
        {restartService.count == 0 && redService.count > 0 && redVisible && (
          <div className="service-attention">
            <div className="attention-content">
              <Icon
                type="exclamation-circle"
                theme="filled"
                style={{
                  color: '#ffa941',
                  marginRight: '8px',
                  marginLeft: '8px',
                }}
              />
              <Dropdown overlay={menu} placement="bottomRight">
                <a onClick={(e) => e.preventDefault()} className="dropdown">
                  <span>一共有{redService.count}个异常服务</span>{' '}
                  <Icon type="down" />
                </a>
              </Dropdown>
            </div>
          </div>
        )}
        <Layout className="service-content">
          <SideNav
            {...sidenavProps}
            handleSwitchService={this.handleSwitchService}
            setCurrentService={this.setCurrentService}
            setSideNavState={(params: any) => this.setState(params)}
            getCurrentProduct={this.getCurrentProduct}
            getRestartService={this.getRestartService}
          />
          <Content>
            <Tabs
              className="service-tab c-tabs-padding"
              style={{ padding: '0 0 0 20px' }}
              activeKey={activeKey}
              onChange={this.handleTabsChange}>
              <TabPane tab="运行状态" key="0">
                <RuntimeStatus
                  {...this.props}
                  cur_service_info={cur_service_info}
                  cur_product_info={cur_product_info}
                  dashId={this.state.dashId}
                  getServiceGroup={this.getServiceGroup}
                  products={products}
                  setCurrentService={this.setCurrentService}
                />
              </TabPane>
              <TabPane tab="参数配置" key="1">
                <ConfigParam
                  {...this.props}
                  {...this.state}
                  isKubernetes={isKubernetes}
                  getCurrentProduct={this.getCurrentProduct}
                  // getCurrentHostsList={this.getCurrentHostsList}
                />
              </TabPane>
              {CAN_DASHBOARD_VIEW && this.state.dashurl !== '' && (
                <TabPane tab="仪表盘" key="2">
                  <iframe
                    className="style-iframe"
                    src={dashUrl}
                    frameBorder="0"
                    style={{ height: 'calc(100vh - 180px)' }}
                  />
                </TabPane>
              )}
            </Tabs>
          </Content>
        </Layout>
        <Modal
          destroyOnClose={true}
          title="执行日志"
          footer={null}
          width={900}
          visible={this.state.log_modal_visible}
          onCancel={this.handleLogModalCancel.bind(this)}>
          <Logtail
            logs={this.state.logpaths}
            serviceid={this.state.log_service_id}
            isreset={!this.state.log_modal_visible}
          />
        </Modal>
      </div>
    );
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceList);
