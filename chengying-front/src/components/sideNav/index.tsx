import * as React from 'react';
import { AppStoreTypes } from '@/stores';
import { connect } from 'react-redux';
import { isEqual } from 'lodash';
import { Menu } from 'antd';
import Utils from '@/utils/utils';
import './style.scss';
const SubMenu = Menu.SubMenu;

interface SideNavProps {
  history?: any;
  match?: any;
  menuList: any[];
  defaultOpenKey?: string[];
  authorityList?: any;
  style?: any;
}

interface SideNavState {
  openKeys: any[];
  selectKeys: string[];
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});

@(connect(mapStateToProps, undefined) as any)
export default class SideNav extends React.Component<
  SideNavProps,
  SideNavState
> {
  constructor(props: SideNavProps) {
    super(props);
    this.state = {
      openKeys: props.defaultOpenKey,
      selectKeys: [sessionStorage.getItem('siderLevelNav')],
    };
  }

  componentDidUpdate(prevProps, prevState) {
    const selectKeys = [sessionStorage.getItem('siderLevelNav')];
    if (selectKeys?.length && !isEqual(prevState.selectKeys, selectKeys)) {
      this.setState({
        selectKeys,
      });
    }
  }

  // 获取菜单列表
  getMenuList = (menuList: any[]) => {
    const list = menuList.map((item) => {
      const { icon, title, children = [], code, level } = item;
      if (!this.props.authorityList[code]) {
        return null;
      }
      const menuItem = (
        <React.Fragment>
          {level === 'second' && (
            <i
              className={`emicon ${icon}`}
              style={{ verticalAlign: 'bottom', marginRight: 8 }}
            />
          )}
          {title}
        </React.Fragment>
      );
      return children.length ? (
        <SubMenu key={code} title={menuItem}>
          {this.getMenuList(children)}
        </SubMenu>
      ) : (
        <Menu.Item key={code} {...item}>
          {menuItem}
        </Menu.Item>
      );
    });
    return list;
  };

  handleMenuSelected = (e) => {
    const { history } = this.props;
    const { key, item } = e;
    if (item.props?.url) {
      this.setState(
        {
          selectKeys: key,
        },
        () => {
          const top = sessionStorage.getItem('firstLevelNav');
          Utils.setNaviKey(top, key);
          history.push(item.props?.url);
        }
      );
    }
  };

  render() {
    const { menuList, style } = this.props;
    const { selectKeys, openKeys } = this.state;
    const sideStyle = {
      minHeight: 'calc(100vh - 40px)',
      height: '100%',
      ...style,
    };
    // console.log('menuList', menuList)
    return (
      <div className="dt-sidemenu-dark" style={sideStyle}>
        <Menu
          mode="inline"
          theme="dark"
          openKeys={openKeys}
          onOpenChange={(e) => this.setState({ openKeys: e })}
          selectedKeys={selectKeys}
          onClick={this.handleMenuSelected}
          style={{ width: 200 }}>
          {this.getMenuList(menuList)}
        </Menu>
      </div>
    );
  }
}
