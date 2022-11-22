import * as React from 'react';
import { Layout, Menu, Icon, Dropdown, message } from 'antd';
import { HeaderStateTypes } from '@/stores/headerReducer';
import { navData } from '@/constants/navData';
import { userCenterService } from '@/services';
import * as Cookie from 'js-cookie';
import Utils from '@/utils/utils';
import { encryptStr, encryptSM } from '@/utils/password';
import HeaderNamespace from '@/components/headerNamespace';
import ResetPassword from '@/pages/userCenter/components/resetPassword';
import InstallGuideModal from '@/components/installGuide';
declare var APP: any;
const logoPng = require('public/imgs/logo_chengying@2x.png'); // tslint:disable-line

const MenuItem = Menu.Item;
const { Header } = Layout;

interface Props {
  location?: any;
  history: any;
  authorityList?: any;
  HeaderStore: HeaderStateTypes;
  actions: any;
}

interface State {
  selfInfo: any;
  isCheckedResetPwd: boolean;
  showModal: boolean;
  showInstallGuide: boolean;
}

class RootHeader extends React.Component<Props, State> {
  state: State = {
    selfInfo: {},
    isCheckedResetPwd: false,
    showModal: false,
    showInstallGuide: false
  };

  componentDidMount() {
    this.getUserInfo();
    if (this.props.location.pathname === '/') {
      this.handleMenuSelected({ key: 'menu_deploy_center' });
    }
  }

  // 获取个人信息
  getUserInfo = () => {
    userCenterService.getLoginedUserInfo().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const { info, reset_password } = res.data;
        const showModal = info.reset_password_status === 0 && reset_password;
        this.setState({
          selfInfo: res.data.info,
          isCheckedResetPwd: showModal,
          showModal,
        });
        localStorage.setItem('em_role_id', info?.role_id);
      } else {
        message.error(res.msg);
      }
    });
  };

  // 构建菜单
  initMenuList = (list: any[]) => {
    const { authorityList } = this.props;
    // const pathname = location.pathname;
    const menu = list.map((item: any) => {
      const { title, code } = item;
      const IS_USER_MANAGE = code === 'menu_user_manage';
      // 过滤
      if (
        !authorityList[code] ||
        IS_USER_MANAGE ||
        item.level !== 'first' ||
        code === 'menu_deploy_guide' ||
        code === 'menu_system_configuration'
      ) {
        return null;
      }
      return <MenuItem key={code}>{title}</MenuItem>;
    });
    return menu;
  };

  // 点击
  handleMenuSelected = (e) => {
    const { HeaderStore } = this.props;
    const isK8s = HeaderStore.cur_parent_cluster?.type === 'kubernetes';
    const current: any = navData.find((nav) => nav.code === e.key) || {};
    if (current?.code && current?.children?.length) {
      const { children, code } = current;
      let menuList = children;
      // 兼容k8s
      if (isK8s && current?.code === 'menu_ops_center') {
        menuList = menuList.filter(
          (m) => isK8s && m.code === 'sub_menu_service'
        );
      }

      const firstLink = menuList[0]?.children?.length
        ? menuList[0]?.children[0]
        : menuList[0];
      Utils.setNaviKey(code, firstLink.code);
      this.props.history.push(firstLink.url);
    } else {
      // 兼容部署向导
      this.setState({showInstallGuide: true})
    }
  };

  handleLogout = () => {
    userCenterService.logout().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        try {
          Cookie.remove('em_token');
          Cookie.remove('em_current_cluster_id');
        } catch (error) {
          window.location.href = 'login';
        }
        window.location.href = 'login';
      } else {
        message.error(res.msg);
        window.location.href = 'login';
      }
    });
  };

  downloadInfo = async() => {
    let res = await userCenterService.generate({em_version: APP.VERSION})
    if (res.data.code == 0) {
      window.open(
        `/api/v2/common/deployInfo/download`,
        '_self'
      );
    };
  }

  getUserHelp = (e) => {
    const { children } = e.item.props;
    const { key } = e;
    if (children == '部署信息下载') {
      this.downloadInfo()
    }
    if (key === '/logout') {
      this.handleLogout();
    } else {
      const naviKey =
        key === '/usercenter/members'
          ? 'sub_menu_user_manage'
          : 'menu_security_audit';
      Utils.setNaviKey('', naviKey);
      this.props.history.push(key);
    }
  };

  // 成员管理的下拉菜单列表
  initUserMenuList = () => {
    const { authorityList } = this.props;
    const userMenu = (
      <Menu onClick={this.getUserHelp}>
        {Cookie.get('em_admin') === 'true' &&
          authorityList.sub_menu_user_manage && (
            <MenuItem key="/usercenter/members">用户管理</MenuItem>
          )}
        {authorityList.menu_deploy_info_markdown_download && (
          <MenuItem key="/usercenter/members">部署信息下载</MenuItem>
        )}
        {authorityList.menu_security_audit && (
          <MenuItem key="/usercenter/audit">安全审计</MenuItem>
        )}
        <MenuItem key="/logout">退出</MenuItem>
      </Menu>
    );
    return userMenu;
  };

  // 重置密码
  handleResetSubmit = async (value: any) => {
    const publicKeyRes = await userCenterService.getPublicKey();
    if (publicKeyRes.data.code !== 0) {
      return;
    }
    const { encrypt_type, encrypt_public_key } = publicKeyRes.data.data;
    const p = {
      old_password:
        encrypt_type === 'sm2'
          ? encryptSM(value.oldPass, encrypt_public_key)
          : encryptStr(value.oldPass, encrypt_public_key),
      new_password:
        encrypt_type === 'sm2'
          ? encryptSM(value.newPass, encrypt_public_key)
          : encryptStr(value.newPass, encrypt_public_key),
    };
    userCenterService.resetPasswordSelf(p).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        message.success('修改成功');
        this.setState({ showModal: false, isCheckedResetPwd: false });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 重置密码弹窗关闭
  resetPwdModalShow = () => {
    this.setState({ showModal: !this.state.showModal }, () => {
      if (!this.state.showModal) {
        this.handleLogout();
      }
    });
  };

  // 关闭部署向导弹框
  closeInstallGuideShow = () => {
    this.setState({showInstallGuide: false})
  }

  render() {
    const { authorityList, history } = this.props;
    const { selfInfo, showModal, isCheckedResetPwd, showInstallGuide } = this.state;
    const firstLevelNav = sessionStorage.getItem('firstLevelNav');
    return (
      <Header className="root-nav">
        <div className="logo dt-header-log-wrapper">
          <img src={logoPng} />
          <span className="dt-header-logo-name">ChengYing</span>
        </div>
        <div className="nav-container">
          <div className="header-ops-namespace">
            {firstLevelNav === 'menu_ops_center' && (
              <HeaderNamespace {...this.props} />
            )}
          </div>
          <Menu
            selectedKeys={[firstLevelNav]}
            mode="horizontal"
            theme="dark"
            onClick={this.handleMenuSelected}>
            {this.initMenuList(navData)}
          </Menu>
        </div>
        <div className="nav-setting">
          {authorityList['menu_deploy_guide'] && (
            <a
              className="em-guide"
              onClick={() =>
                this.handleMenuSelected({ key: 'menu_deploy_guide' })
              }>
              <i
                className={`emicon emicon-yunhangbushuxiangdao`}
                style={{ verticalAlign: 'bottom', marginRight: 8 }}
              />
              运行部署向导
            </a>
          )}
          &emsp;&emsp;&emsp;
          {authorityList['menu_system_configuration'] && (
            <a
              className="em-guide"
              onClick={() =>
                this.handleMenuSelected({ key: 'menu_system_configuration' })
              }>
              <Icon type="setting" style={{ marginRight: 8 }} />
              系统配置
            </a>
          )}
          {window.APPCONFIG.userCenter && (
            <Dropdown
              overlayClassName="dt-em-dropdown-black root-drop"
              overlay={this.initUserMenuList()}
              className="dropdown">
              <a onClick={(e: any) => e.preventDefault()}>
                {selfInfo?.username || Cookie.get('em_username')}
                <Icon style={{ marginLeft: 1 }} type="down" />
              </a>
            </Dropdown>
          )}
        </div>
        <ResetPassword
          visible={showModal}
          isCheckedResetPwd={isCheckedResetPwd}
          onCancel={this.resetPwdModalShow}
          onSubmit={this.handleResetSubmit}
        />
        <InstallGuideModal
          onClose={this.closeInstallGuideShow}
          history={history}
          visible={showInstallGuide}
        />
      </Header>
    );
  }
}

export default RootHeader;
