import * as React from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Dispatch, bindActionCreators } from 'redux';
import { Layout, Menu, Icon, message, Modal } from 'antd';
import { AppStoreTypes } from '@/stores';
import * as UserCenterAction from '@/actions/userCenterAction';
import * as ServiceActions from '@/actions/serviceAction';
import * as HeaderAction from '@/actions/headerAction';
import { HeaderStateTypes } from '@/stores/headerReducer';
import { navData } from '@/constants/navData';
import * as Cookie from 'js-cookie';
import * as ClipboardJS from 'clipboard';
import Utils from '@/utils/utils';
import * as Http from '@/utils/http';
import SideNav from '@/components/sideNav';
import RootHeader from './header';
import ClassNames from 'classnames';
import './style.scss';

const { Content, Sider } = Layout;

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
  authorityRouter: state.UserCenterStore.authorityRouter,
  ServiceStore: state.ServiceStore,
  HeaderStore: state.HeaderStore,
});
const mapDispatchToProps = (
  dispatch: Dispatch<{ type: string; payload?: any }>
) => ({
  actions: bindActionCreators(
    Object.assign({}, UserCenterAction, ServiceActions, HeaderAction),
    dispatch
  ),
});
interface Props {
  actions: any;
  location?: any;
  history: any;
  authorityList?: any;
  authorityRouter?: any[];
  ServiceStore: any;
  HeaderStore: HeaderStateTypes;
}
interface State {
  lookMoreVisible: boolean;
  tooltipVisible: boolean;
  popupVisible: boolean;
}

let timeEvent;
@(connect(mapStateToProps, mapDispatchToProps) as any)
class RootLayout extends React.Component<Props, any> {
  state: State = {
    lookMoreVisible: false,
    tooltipVisible: false,
    popupVisible: false,
  };

  componentDidMount() {
    if (!Cookie.get('em_token') && window.APPCONFIG.userCenter) {
      this.props.location.pathname = '/login';
    }
    this.redirectRoute();
    this.getRoleCodes();
    const clipboard = new ClipboardJS('#J_CopyBtn');
    clipboard.on('success', (e: any) => {
      message.success('复制成功！');
    });
    clipboard.on('error', (e: any) => {
      message.error('复制失败！');
    });
    this.getRestartService();
  }

  componentDidUpdate(prevProps: Props) {
    const { location } = this.props;
    if (this.props.location.pathname != prevProps.location.pathname) {
      this.redirectRoute();
    }
    if (
      location?.pathname.indexOf('opscenter') > -1 &&
      sessionStorage.getItem('firstLevelNav') != 'menu_ops_center'
    ) {
      Utils.setNaviKey(
        'menu_ops_center',
        sessionStorage.getItem('siderLevelNav')
      );
    }
    if (
      location?.pathname.indexOf('deploycenter') > -1 &&
      sessionStorage.getItem('firstLevelNav') != 'menu_deploy_center'
    ) {
      Utils.setNaviKey(
        'menu_deploy_center',
        sessionStorage.getItem('siderLevelNav')
      );
    }
  }

  componentWillUnmount() {
    clearTimeout(timeEvent);
  }

  // 获取需要依赖组件重新配置服务列表——EM提醒
  getRestartService = () => {
    Http.get('/api/v2/cluster/restartServices', {}).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.props.actions.setResartServiceList({
          count: res.data?.count,
          list: res.data?.list,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 无权限时路由重置
  redirectRoute = () => {
    // 路由权限
    const { authorityRouter } = this.props;
    if (
      Array.isArray(authorityRouter) &&
      authorityRouter.length &&
      !authorityRouter.includes(location.pathname)
    ) {
      this.props.history.push('/nopermission');
    }
  };

  // 获取权限表
  getRoleCodes = () => {
    this.props.actions.setRoleAuthorityList();
  };
  
  getMenuList = () => {
    const isK8s = Cookie.get('em_current_cluster_type') === 'kubernetes';
    const urlParams = this.props.location.pathname.split('/');
    const menu = navData.find((item: any) => item.url === `/${urlParams[1]}`);
    let menuList = menu?.children || [];
    // 兼容k8s
    if (isK8s && menu?.code === 'menu_ops_center') {
      menuList = menuList.filter((m) => isK8s && m.code === 'sub_menu_service');
    }
    return menuList;
  };

  render() {
    const { ServiceStore, authorityList, location, history, HeaderStore, actions } = this.props;
    const pathname = this.props.location.pathname;
    const { restartService } = ServiceStore;
    const restarList =
      restartService.list &&
      restartService.list.map((item, index) => {
        return (
          <span className="attention-title" key={index}>
            {item.depend_product_name}——{item.depend_service_name}
            配置变更后,相关依赖服务
            <span className="red-attention">
              {item.product_name}——{item.service_name}
            </span>
            需要重新“滚动重启”。
          </span>
        );
      });
    return (
      <Layout className="root-layout">
        <RootHeader 
          authorityList={authorityList} 
          location={location} 
          HeaderStore={HeaderStore} 
          history={history}
          actions={actions}
        />
        <Layout style={{ marginTop: 48, position: 'fixed', width: '100%' }}>
          {pathname != '/installguide' && (
            <Sider className="dt-layout-sider">
              <SideNav
                {...this.props}
                menuList={this.getMenuList()}
                style={{ minHeight: `calc(100vh - 48px)`, height: '100%' }}
              />
            </Sider>
          )}
          <Content
            className={ClassNames('root-content', {
              'has-alert__gobal': restartService.count > 0,
            })}>
            {this.props.children}
          </Content>
        </Layout>
        {restartService.count > 0 && (
          <div className="root-attention">
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
                <span className="attention-more">{restarList}</span>
                {restartService.count >= 3 && (
                  <a
                    onClick={() => {
                      this.setState({ lookMoreVisible: true });
                    }}>
                    查看更多
                  </a>  
                )}
              </div>
            </div>
          </div>
        )}
        <Modal
          visible={this.state.lookMoreVisible}
          title="组件变更提示"
          onCancel={() => {
            this.setState({ lookMoreVisible: false });
          }}
          footer={null}>
          <ul>
            {restartService.list &&
              restartService.list.map((item, index) => {
                return (
                  <li key={index}>
                    {item.depend_product_name}——{item.depend_service_name}
                    配置变更后,相关依赖服务
                    <span style={{ color: '#ff5f5c' }}>
                      {item.product_name}——{item.service_name}
                    </span>
                    需要重新“滚动重启”。
                  </li>
                );
              })}
          </ul>
        </Modal>
      </Layout>
    );
  }
}

export default withRouter<any, any>(RootLayout);
