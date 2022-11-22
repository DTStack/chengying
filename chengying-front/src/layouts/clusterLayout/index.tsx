import * as React from 'react';
import { bindActionCreators, Dispatch } from 'redux';
import { connect } from 'react-redux';
import { Layout, Breadcrumb, Tabs, Menu, Dropdown, Icon, Tooltip } from 'antd';
import { Link } from 'react-router-dom';
import * as Cookie from 'js-cookie';
import * as HeaderAction from '@/actions/headerAction';
import { HeaderStateTypes } from '@/stores/headerReducer';
import { AppStoreTypes } from '@/stores';
import { clusterNavData } from '@/constants/navData';
import './style.scss';
import { isEqual } from 'lodash';
const { Content } = Layout;
const { TabPane } = Tabs;

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
});
const mapDispatchToProps = (
  dispatch: Dispatch<{ type: string; payload?: any }>
) => ({
  actions: bindActionCreators(HeaderAction, dispatch),
});

interface IProps {
  HeaderStore?: HeaderStateTypes & {
    cur_parent_cluster: {
      status?: string;
    };
  };
  actions?: HeaderAction.HeaderActionTypes;
  location?: any;
  history: any;
  match?: any;
}

interface IState {
  activeKey: string;
}

const clusterTypeMaps = {
  hosts: {
    0: ['imagestore', 'namespace'],
  },
  kubernetes: {
    0: ['namespace', 'patchHistory', 'echoList'], // 自建
    1: ['index', 'host', 'patchHistory', 'echoList'], // 导入
  },
};
@(connect(mapStateToProps, mapDispatchToProps) as any)
export default class ClusterLayout extends React.PureComponent<IProps, any> {
  state: IState = {
    activeKey: 'index',
  };

  componentDidMount() {
    const { location, actions } = this.props;
    const { pathname } = location;
    actions.getClusterList();
    if (pathname) {
      let arrs = pathname.split('/');
      this.setState({
        activeKey: arrs[arrs.length - 1],
      });
    }
  }

  componentDidUpdate(nextProps, nextState) {
    if (isEqual(nextState.activeKey, this.state.activeKey)) {
      const { location } = this.props;
      let arrs = location.pathname.split('/');
      this.setState({
        activeKey: arrs[arrs.length - 1],
      });
    }
  }

  changeActive = (activeKey) => {
    const { history } = this.props;
    this.setState(
      {
        activeKey,
      },
      () => {
        history.push(`/deploycenter/cluster/detail/${activeKey}`);
      }
    );
  };

  handleSwitchProduct = (e: any) => {
    const { key } = e;
    const { parentClusters } = this.props.HeaderStore;
    const cluster = parentClusters.find((item) => item.id === +key);
    console.log('cluster', cluster);
    this.props.actions.setCurrentParentProduct(cluster);
    Cookie.set('em_current_cluster_id', cluster?.id);
    Cookie.set('em_current_cluster_type', cluster?.type);
    window.location.href =
      cluster?.mode === 0
        ? '/deploycenter/cluster/detail/index'
        : '/deploycenter/cluster/detail/imagestore';
  };

  render() {
    const { children, HeaderStore } = this.props;
    const { parentClusters, cur_parent_cluster } = HeaderStore;
    const { type, mode } = cur_parent_cluster;
    const { activeKey } = this.state;
    const clusterMaps = clusterTypeMaps[type];
    const productMenuOverlay = (
      <Menu onClick={this.handleSwitchProduct}>
        {parentClusters.map((item) => (
          <Menu.Item data-testid={`header-nav-${item.name}`} key={item.id}>
            {item.name}
          </Menu.Item>
        ))}
      </Menu>
    );
    // 过滤条件
    const filterConditions = clusterMaps[mode];
    // 过滤不同集群类型菜单
    const realClusterData = clusterNavData.filter(
      (nav) => nav.isShow && !filterConditions.includes(nav.key)
    );

    return (
      <Layout className="container">
        <div className="container-header">
          <Breadcrumb className="bread">
            <Breadcrumb.Item>
              <Link to="/deploycenter/cluster/list">集群管理</Link>
            </Breadcrumb.Item>
            <Breadcrumb.Item>集群详情</Breadcrumb.Item>
          </Breadcrumb>
          <Dropdown
            trigger={['click']}
            overlay={productMenuOverlay}
            overlayStyle={{
              maxHeight: 220,
              overflowY: 'auto',
              boxShadow: '0px 2px 8px 0px rgb(6 14 26 / 8%)',
            }}>
            <a
              className="ant-dropdown-link"
              onClick={(e) => e.preventDefault()}>
              <Tooltip title={cur_parent_cluster?.name}>
                <span className="title">
                  {cur_parent_cluster?.name || '主机集群名称筛选'}
                </span>{' '}
                <Icon type="down" />
              </Tooltip>
            </a>
          </Dropdown>
        </div>
        <Tabs activeKey={activeKey} onChange={this.changeActive}>
          {realClusterData.map((nav) => (
            <TabPane tab={nav.title} key={nav.key} />
          ))}
        </Tabs>
        <Content>{children}</Content>
      </Layout>
    );
  }
}
