import * as React from 'react';
import { connect } from 'react-redux';
import {
  Layout,
  Form,
  Menu,
  Collapse,
  Input,
  Tooltip,
  message,
  Spin,
  Empty,
  Icon,
} from 'antd';
import { AppStoreTypes } from '@/stores';
import { HeaderStoreType } from '@/stores/modals';
import { configService } from '@/services';
import moment from 'moment';
import * as Cookies from 'js-cookie';
import './style.scss';
const { SubMenu } = Menu;
const { Content, Sider } = Layout;
const { Panel } = Collapse;
const FormItem = Form.Item;

interface IProps {
  HeaderStore: HeaderStoreType;
}

interface IState {
  menuOpenKey: string[];
  menuSelectKey: string[];
  collapseActiveKey: string;
  menuList: any[];
  serviceList: any[];
  serviceLoading: boolean;
  menuLoading: boolean;
}

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 8 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16 },
  },
};

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
});

@(connect(mapStateToProps, undefined) as any)
export default class Config extends React.PureComponent<IProps, IState> {
  state: IState = {
    menuOpenKey: [''],
    menuSelectKey: [''],
    collapseActiveKey: 'redis',
    menuList: [],
    serviceList: [],
    serviceLoading: false,
    menuLoading: false,
  };

  componentDidMount() {
    this.getAlertGroup();
  }

  // 获取发生变更的组件及其下服务组
  getAlertGroup = () => {
    const { HeaderStore } = this.props;
    console.log(HeaderStore);
    const params = {
      limit: 0,
      parentProductName:
        HeaderStore?.cur_parent_product != '选择产品'
          ? HeaderStore?.cur_parent_product
          : Cookies.get('em_current_parent_product'),
    };
    this.setState({ menuLoading: true });
    configService.getConfigAlertGroups(params).then((res: any) => {
      this.setState({ menuLoading: false });
      const {
        code,
        msg,
        data: { count, list = [] },
      } = res.data;
      if (code === 0) {
        const params: any = {
          menuList: list,
        };
        if (count) {
          const { product_name, groups = [] } = list[0];
          params.menuOpenKey = [product_name];
          params.menuSelectKey = [`${product_name}_${groups[0]}`];
          this.getConfigAlertaction(product_name, groups[0]);
        }
        this.setState({
          ...params,
        });
      } else {
        message.error(msg);
      }
    });
  };

  // 获取服务组下发生变更的服务
  getConfigAlertaction = (ProductName, group) => {
    const { HeaderStore } = this.props;
    const params = {
      limit: 0,
      parentProductName: HeaderStore.cur_parent_product,
      ProductName,
      group,
    };
    this.setState({ serviceLoading: true });
    configService.getConfigAlertaction(params).then((res: any) => {
      this.setState({ serviceLoading: false });
      const {
        code,
        msg,
        data: { list = [] },
      } = res.data;
      if (code === 0) {
        const firstGroup = list[0] || {};
        this.setState({
          serviceList: list,
          collapseActiveKey: firstGroup.service_name,
        });
      } else {
        message.error(msg);
      }
    });
  };

  // 切换服务组
  handleServiceChange = (e) => {
    const keyArr = e.key.split('_');
    this.setState({
      menuSelectKey: [e.key],
    });
    this.getConfigAlertaction(keyArr[0], keyArr[1]);
  };

  // 切换服务
  handleCollapseChange = (e) => {
    this.setState({
      collapseActiveKey: e,
    });
  };

  render() {
    const {
      menuList = [],
      serviceList = [],
      menuLoading,
      serviceLoading,
      menuOpenKey,
      menuSelectKey,
      collapseActiveKey,
    } = this.state;
    return (
      <Layout>
        <Content className="config-page">
          <Spin spinning={menuLoading}>
            {menuList.length ? (
              <Layout className="config-content box-shadow-style">
                <Sider className="config-sider" width={250}>
                  <Menu
                    className="dt-em-menu"
                    mode="inline"
                    openKeys={menuOpenKey}
                    selectedKeys={menuSelectKey}
                    onOpenChange={(e) => this.setState({ menuOpenKey: e })}
                    onClick={this.handleServiceChange}>
                    {menuList.map((list) => (
                      <SubMenu
                        key={list.product_name}
                        title={
                          <div className="submenu-title-style">
                            <span>{list.product_name}</span>
                            <Icon
                              type="down"
                              className={
                                menuOpenKey.includes(list.product_name)
                                  ? 'up'
                                  : ''
                              }
                            />
                          </div>
                        }>
                        {Array.isArray(list.groups) &&
                          list.groups.map((group) => (
                            <Menu.Item key={`${list.product_name}_${group}`}>
                              {group}
                            </Menu.Item>
                          ))}
                      </SubMenu>
                    ))}
                  </Menu>
                </Sider>
                <Content>
                  <Spin spinning={serviceLoading}>
                    <Collapse
                      activeKey={collapseActiveKey}
                      bordered={false}
                      onChange={this.handleCollapseChange}>
                      {serviceList.map((item) => {
                        const { service_name, alteration = [] } = item;
                        return (
                          <Panel key={service_name} header={service_name}>
                            <Form>
                              {alteration.map((param) => {
                                const label = (
                                  <React.Fragment>
                                    {param.isnew && (
                                      <i className="emicon emicon-new icon-style" />
                                    )}
                                    {param.config.length > 20 ? (
                                      <Tooltip title={param.config}>
                                        <span className="text-ellipsis label">
                                          {param.config}
                                        </span>
                                      </Tooltip>
                                    ) : (
                                      <span>{param.config}</span>
                                    )}
                                  </React.Fragment>
                                );
                                return (
                                  <FormItem
                                    {...formItemLayout}
                                    label={label}
                                    key={param.config}>
                                    {param.current.length > 44 ? (
                                      <Tooltip title={param.current}>
                                        <Input
                                          className="text-ellipsis input"
                                          value={param.current}
                                          disabled
                                          addonAfter={
                                            <i
                                              className="emicon emicon-undo"
                                              style={{ fontSize: 14 }}
                                            />
                                          }
                                          style={{ width: 400 }}
                                        />
                                      </Tooltip>
                                    ) : (
                                      <Input
                                        value={param.current}
                                        disabled
                                        addonAfter={
                                          <i
                                            className="emicon emicon-undo"
                                            style={{ fontSize: 14 }}
                                          />
                                        }
                                        style={{ width: 400 }}
                                      />
                                    )}
                                    <Tooltip
                                      title={
                                        <React.Fragment>
                                          <p>
                                            最近修改时间：
                                            {moment(param.updated).format(
                                              'YYYY-MM-DD HH:mm:ss'
                                            )}
                                          </p>
                                          <p>默认值：{param.default}</p>
                                        </React.Fragment>
                                      }>
                                      <i className="emicon emicon-revise icon-style"></i>
                                    </Tooltip>
                                  </FormItem>
                                );
                              })}
                            </Form>
                          </Panel>
                        );
                      })}
                    </Collapse>
                  </Spin>
                </Content>
              </Layout>
            ) : (
              <Empty
                className="config-empty box-shadow-style"
                description="暂无数据变更"
              />
            )}
          </Spin>
        </Content>
      </Layout>
    );
  }
}
