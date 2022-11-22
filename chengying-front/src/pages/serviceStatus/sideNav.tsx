import * as React from 'react';
import { Layout, Dropdown, Menu, message, Modal, Icon } from 'antd';
import { servicePageService } from '@/services';
import * as Cookie from 'js-cookie';
import isEqual from 'lodash/isEqual';
import utils from '@/utils/utils';
import AutoTest from './autoTest';

const { Sider } = Layout;
const SubMenu = Menu.SubMenu;

interface IProps {
  products: any[];
  cur_product_info: {
    product_id: number;
    product_name: string;
    product_version?: string;
  };
  menuKey: {
    openKeys: any[];
    selectedKeys: any[];
  };
  isKubernetes: boolean;
  ServiceStore: any;
  HeaderStore: any;
  actions: any;
  authorityList: any;
  handleSwitchService: Function;
  setCurrentService: Function;
  setSideNavState: Function; // setState
  getCurrentProduct: Function;
  getRestartService: Function;
}

interface IState {
  openKeys: any[];
  selectedKeys: any[];
  visible: boolean;
  smokeInfo: {
    cluster_id?: number;
    product_name?: string;
    create_time?: string;
    end_time?: string;
    auto_test?: boolean;
    exec_status?: number;
    report_url?: string;
  };
  smokeMsg: string;
}

export default class SideNav extends React.PureComponent<IProps, IState> {
  private timeInterval;

  state: IState = {
    openKeys: ['f0'],
    selectedKeys: ['sub0_0'],
    visible: false,
    smokeInfo: {},
    smokeMsg: '',
  };

  componentDidUpdate(prevProps, prevState) {
    // 用户角色 1、2、3
    const em_roleId = +localStorage.getItem('em_role_id') || 1;
    if (!isEqual(this.props.cur_product_info, prevProps.cur_product_info) && em_roleId !== 3) {
      // console.log('componentDidUpdate')
      clearTimeout(this.timeInterval);
      this.getAutoTestStatus();
    }
    if (
      !isEqual(this.props.menuKey?.openKeys, prevState?.openKeys) ||
      !isEqual(this.props.menuKey?.selectedKeys, prevState?.selectedKeys)
    ) {
      this.setState({
        openKeys: this.props.menuKey?.openKeys,
        selectedKeys: this.props.menuKey?.selectedKeys,
      });
    }
  }

  componentWillUnmount() {
    clearTimeout(this.timeInterval);
  }

  // 切换产品
  handleSelectProduct = (item: any) => {
    const { products } = this.props;
    const pn = products[item.key].product_name;
    const pv = products[item.key].product_version;
    if (pn === '') {
      return;
    }
    servicePageService
      .getCurrentProduct({ product_name: pn })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          Cookie.set('em_product_name', data.data.product_name);
          Cookie.set('em_product_id', data.data.id);
          this.props.actions.setCurrentProduct(data.data);
          this.props.actions.getServiceGroup(
            {
              product_name: pn,
            },
            (firstService: any, key: Array<string>, selectedKey: any) => {
              this.props.setSideNavState({
                cur_product_info: {
                  product_id: data.data.id,
                  product_name: pn,
                  product_version: pv,
                },
                services: data.data.product.Service,
                cur_service: firstService,
                menuKey: {
                  selectedKeys: selectedKey,
                  openKeys: [...key],
                },
              });
              this.props.handleSwitchService(firstService, pn);
              this.props.setCurrentService(
                data.data.product.Service[firstService.service_name],
                firstService.service_name
              );
            },
            '',
            { cloud: true }
          );
        } else {
          message.error(data.msg);
        }
      });
  };

  handleUseOption = (e) => {
    if (e.key != '2') {
      this.handleSelectSSOption(e);
    } else {
      this.handleToggle();
    }
  };

  handleToggle = () => {
    let { smokeMsg, visible } = this.state;
    if (visible) {
      smokeMsg = '';
    } else {
      // 打开弹窗
      clearTimeout(this.timeInterval);
      this.getAutoTestStatus();
    }
    this.setState({
      visible: !visible,
      smokeMsg,
    });
  };

  // 获取自动化测试状态
  getAutoTestStatus = () => {
    // debugger
    const { cur_product_info, HeaderStore } = this.props;
    const { cur_parent_cluster } = HeaderStore;
    const params = {
      product_name: cur_product_info?.product_name,
      cluster_id: cur_parent_cluster?.id,
    };
    if (!params.product_name) return;
    servicePageService.getAutoTestStatus(params).then((res) => {
      const { code, data } = res.data;
      if (code === 0) {
        this.setState(
          {
            smokeInfo: data,
          },
          () => {
            this.timeInterval = setTimeout(() => {
              this.getAutoTestStatus();
            }, 5000);
          }
        );
      }
    });
  };

  // 触发自动化测试
  handleTest = () => {
    const { cur_product_info, HeaderStore } = this.props;
    const { cur_parent_cluster } = HeaderStore;
    const { smokeInfo } = this.state;
    const params = {
      product_name: cur_product_info?.product_name,
      cluster_id: cur_parent_cluster?.id,
    };
    // 执行测试
    const autoTest = () => {
      servicePageService.startAutoTestTest(params).then((res) => {
        const { code, msg } = res.data;
        if (code != 0) {
          smokeInfo.exec_status = 3;
          this.setState({
            smokeInfo,
            smokeMsg: msg,
          });
        } else {
          clearTimeout(this.timeInterval);
          this.getAutoTestStatus();
        }
      });
    };

    this.setState(
      {
        smokeInfo: {
          ...smokeInfo,
          exec_status: 1,
        },
      },
      autoTest
    );
  };

  // 组件启停
  handleSelectSSOption = (item: any) => {
    const self = this;
    const {
      products,
      cur_product_info: { product_name, product_version },
      authorityList,
    } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'service_product_start_stop')) {
      return;
    }
    let pid = -1;
    for (const p of products) {
      if (p.product_name === product_name) {
        pid = p.id;
      }
    }
    const stopContent = (
      <div>
        <p>1. 为防止误操作停止，若配置了监控告警，平台将正常发送告警；</p>
        <p>2. 组件停止后，若服务器发生断电恢复的情况，服务将自动重启；</p>
        <p>3. 可通过点击“启动”按钮恢复组件服务。</p>
      </div>
    );
    const ssmodal = Modal.confirm({
      title:
        item.key === '0'
          ? '停止组件将停止其下所有服务，请注意并谨慎操作：'
          : '启动该组件的所有服务吗？',
      content: item.key === '0' ? stopContent : '',
      okType: item.key === '0' ? 'danger' : 'primary',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okText: item.key === '0' ? '停止' : '启动',
      className: 'ss-confirm-modal',
      width: item.key === '0' ? 550 : 416,
      onOk() {
        if (item.key === '0') {
          ssmodal.update({
            content: <p>{product_name}服务正在停止中，请耐心等待...</p>,
            cancelButtonProps: {
              disabled: true,
            },
          });
          return servicePageService
            .stopAllServiceByProduct({ pid: pid })
            .then((res: any) => {
              res = res.data;
              if (res.code === 0) {
                self.props.getCurrentProduct(product_name, product_version);
                self.props.getRestartService();
              } else {
                message.error(res.msg);
                ssmodal.destroy();
              }
            })
            .catch(() => {
              ssmodal.destroy();
            });
        } else {
          ssmodal.update({
            content: <p>{product_name}服务正在启动中，请耐心等待...</p>,
            cancelButtonProps: {
              disabled: true,
            },
          });
          return servicePageService
            .startAllServiceByProduct({ pid: pid })
            .then((res: any) => {
              res = res.data;
              if (res.code === 0) {
                self.props.getCurrentProduct(product_name, product_version);
                self.props.getRestartService();
              } else {
                message.error(res.msg);
                ssmodal.destroy();
              }
            })
            .catch(() => {
              ssmodal.destroy();
            });
        }
      },
    });
  };

  // submenu子项折叠
  handleMenuOpenChange = (openKeys: any) => {
    this.setState(
      {
        openKeys: openKeys,
      },
      () => {
        this.props.setSideNavState({
          menuKey: {
            selectedKeys: this.state.selectedKeys,
            openKeys: openKeys,
          },
        });
      }
    );
  };

  // 切换服务
  handleChangeService(e: any) {
    const {
      cur_product_info: { product_name },
    } = this.props;
    const service_group: any[] = [];
    for (const s in this.props.ServiceStore.services) {
      service_group.push({
        name: s,
        sub: this.props.ServiceStore.services[s],
      });
    }
    const pos = e.key.substr(3).split('_');
    // 切换服务需要更新服务列表的状态
    this.props.handleSwitchService(
      service_group[pos[0]].sub[pos[1]],
      product_name
    );
    this.props.setSideNavState({
      cur_service: service_group[pos[0]].sub[pos[1]],
      menuKey: {
        selectedKeys: [e.key],
        openKeys: this.state.openKeys,
      },
    });
    this.setState({
      selectedKeys: [e.key],
    });
  }

  render() {
    const { products, cur_product_info, isKubernetes, ServiceStore } =
      this.props;
    const { openKeys, selectedKeys, visible, smokeInfo, smokeMsg } = this.state;
    // console.log('products', products)
    const productMenu = (
      <Menu onClick={this.handleSelectProduct}>
        {products.map((p, i) => {
          return <Menu.Item key={i}>{p.product_name}</Menu.Item>;
        })}
      </Menu>
    );
    const productSSMenu = (
      <Menu onClick={this.handleUseOption}>
        <Menu.Item key="0">停止</Menu.Item>
        <Menu.Item key="1">启动</Menu.Item>
        {smokeInfo?.auto_test && <Menu.Divider />}
        {smokeInfo?.auto_test && (
          <Menu.Item key="2">
            自动化测试
            {smokeInfo?.exec_status === 1 && (
              <Icon type="loading" style={{ marginLeft: 10 }} />
            )}
          </Menu.Item>
        )}
      </Menu>
    );
    const service_group = [];

    for (const s in ServiceStore.services) {
      service_group.push({
        name: s,
        sub: ServiceStore.services[s],
      });
    }
    return (
      <Sider
        width={240}
        className="box-shadow-style"
        style={{
          marginBottom: 20,
          minHeight: document.body.clientHeight - 88 - 40,
          overflow: 'hidden',
        }}>
        <div className="side-product">
          <div className="product-filter-wrapper clearfix">
            <Dropdown
              overlayClassName="dt-em-dropdown-light service-dt-em-dropdown-black"
              overlay={productMenu}
              trigger={['click']}>
              <a className="ant-dropdown-link" href="#">
                {cur_product_info.product_name} <Icon type="down" />
              </a>
            </Dropdown>
            {!isKubernetes && (
              <Dropdown
                overlayClassName="dt-em-dropdown-light service-dt-em-dropdown-setting"
                overlay={productSSMenu}
                trigger={['click']}>
                <a
                  className="fl-r"
                  style={{ color: '#3f87ff' }}
                  onClick={(e) => e.preventDefault()}>
                  <Icon type="setting" />
                </a>
              </Dropdown>
            )}
          </div>
        </div>
        <Menu
          className="c-sidenav__menu"
          mode="inline"
          theme="light"
          selectedKeys={selectedKeys}
          openKeys={openKeys}
          onOpenChange={this.handleMenuOpenChange}
          onClick={this.handleChangeService.bind(this)}>
          {service_group.length > 1
            ? service_group.map((n, index) => {
                return (
                  <SubMenu key={'f' + index} title={<span>{n.name}</span>}>
                    {n.sub.map((s: any, j: any) => {
                      const {
                        cur_product: { product },
                      } = ServiceStore;
                      const service = product.Service[s.service_name];
                      return (
                        <Menu.Item key={'sub' + index + '_' + j}>
                          <i
                            className="dot-cls"
                            style={
                              s.alert
                                ? { display: 'inline-block' }
                                : { display: 'none' }
                            }
                          />
                          <div className="clearfix">
                            <span
                              className="c-text-ellipsis"
                              style={{ maxWidth: 100 }}>
                              {s.service_name_display}
                            </span>
                            {service && service.Version && (
                              <div className="service-version">
                                {service.Version.split('-')[0].substring(
                                  0,
                                  1
                                ) !== 'v' && 'v'}
                                {service.Version.split('-')[0]}
                              </div>
                            )}
                          </div>
                        </Menu.Item>
                      );
                    })}
                  </SubMenu>
                );
              })
            : service_group.map((n: any, i: number) => {
                return n.sub.map((s: any, j: number) => {
                  const {
                    cur_product: { product },
                  } = ServiceStore;
                  const service = product.Service[s.service_name];
                  return (
                    <Menu.Item key={'sub' + 0 + '_' + j}>
                      <i
                        className="dot-cls"
                        style={
                          s.alert
                            ? { display: 'inline-block' }
                            : { display: 'none' }
                        }
                      />
                      <div className="clearfix">
                        <span
                          className="c-text-ellipsis"
                          style={{ maxWidth: 100 }}>
                          {s.service_name_display}
                        </span>
                        {service && service.Version && (
                          <div className="service-version">
                            {service.Version.split('-')[0].substring(0, 1) !==
                              'v' && 'v'}
                            {service.Version.split('-')[0]}
                          </div>
                        )}
                      </div>
                    </Menu.Item>
                  );
                });
              })}
        </Menu>
        {visible && (
          <AutoTest
            visible={visible}
            autoTestInfo={smokeInfo}
            autoTestMsg={smokeMsg}
            handleClose={this.handleToggle}
            handleTest={this.handleTest}
          />
        )}
      </Sider>
    );
  }
}
