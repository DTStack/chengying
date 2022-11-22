import * as React from 'react';
import { Tabs } from 'antd';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { Dispatch, bindActionCreators } from 'redux';
import * as HostAction from '@/actions/hostAction';
import { Service } from '@/services';
import SideNav from './sideNav';
import HostDetail from './hostDetail';
import HostDash from './hostDash';
import './style.scss';
import { isEqual } from 'lodash';

const TabPane = Tabs.TabPane;

const mapStateToProps = (state: AppStoreTypes) => ({
  hostList: state.HostStore.hostList,
  selectedHostServices: state.HostStore.selectedHostServices,
  selectedHost: state.HostStore.selectedHost,
  cur_parent_product: state.HeaderStore.cur_parent_product,
  cur_parent_cluster: state.HeaderStore.cur_parent_cluster,
  selectedIndex: state.HostStore.selectedIndex,
  hostGroupLists: state.HostStore.hostGroupLists,
  pager: state.HostStore.pager,
});

const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, HostAction), dispatch),
});

interface State {
  sideNavHeight: number;
  searchText: string;
  sideNavLoading: boolean;
  selectedGroup: string;
  currentHost: any;
}
interface Prop {
  hostList: any[];
  hostGroupLists: any[];
  pager: any;
  selectedHost: object & { ip: string };
  actions: HostAction.HostActionTypes;
  cur_parent_product: any;
  cur_parent_cluster: any;
  selectedHostServices: any[];
  selectedIndex: number;
  history?: any;
}

interface HostProp {
  hostIp: string;
  [key: string]: any;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
class HostStatus extends React.Component<Prop, State> {
  state: State = {
    sideNavHeight: document.body.clientHeight - 128,
    searchText: '',
    sideNavLoading: false,
    selectedGroup: '',
    currentHost: {},
  };

  componentDidMount() {
    this.getHostGroups();
  }

  componentDidUpdate(prevProps, prevState) {
    const { currentHost } = this.state;
    const { hostList, cur_parent_cluster } = this.props;
    if (!isEqual(cur_parent_cluster, prevProps.cur_parent_cluster)) {
      this.getHostGroups();
    }
    if (!isEqual(hostList, prevProps.hostList)) {
      // 默认选中第一个
      const param = {
        host: hostList[0],
        index: 0,
      };
      if (currentHost?.hostIp !== param.host) {
        hostList.forEach((o, i) => {
          if (o.ip === currentHost.ip) {
            param.host = o;
            param.index = i;
          }
        });
      }
      this.handleSideHostClick(param.host, param.index);
    }
  }

  componentWillUnmount() {
    // 重置更新hostslist
    this.props.actions.updateHostsList({
      data: {
        hosts: [],
        count: 0,
      },
    });
  }

  // 获取主机集群下主机列表
  getHostList = (pageParams) => {
    const { id, type } = this.props.cur_parent_cluster;
    const params = Object.assign(
      {
        parent_product_name: this.props.cur_parent_product,
        host_or_ip: this.state.searchText,
        limit: 0,
        start: 0,
        cluster_id: id,
      },
      pageParams
    );
    this.setState({ sideNavLoading: true });
    // /api/v2/cluster/hosts/hosts
    this.props.actions.getHostList(params, type, () => {
      this.setState({ sideNavLoading: false });
    });
    const loadTime = setTimeout(() => {
      clearTimeout(loadTime);
      if (this.state.sideNavLoading) {
        this.setState({ sideNavLoading: false });
      }
    }, 2000);
  };

  // 获取主机分组
  getHostGroups = (group?: any) => {
    const { cur_parent_product, actions, cur_parent_cluster } = this.props;
    const { type, id } = cur_parent_cluster;
    if (id <= 0 || !cur_parent_product) {
      return;
    }
    // /api/v2/cluster/hostgroups
    actions.gethostGroupLists(
      {
        type,
        id,
        host_or_ip: this.state.searchText,
        parent_product_name: cur_parent_product,
      },
      async (hostsGroup) => {
        const sessionHost: HostProp = JSON.parse(
          sessionStorage.getItem('service_object')
        );
        const res = await Service.getClusterHostList(
          {
            parent_product_name: cur_parent_product,
            host_or_ip: sessionHost?.hostIp || '',
            limit: 0,
            start: 0,
            cluster_id: id,
          },
          type
        );
        const currentHost = res?.data?.data?.hosts[0];
        if (currentHost) {
          this.setState(
            {
              selectedGroup: sessionHost ? currentHost?.group : hostsGroup[0],
              currentHost,
            },
            () => {
              sessionStorage.removeItem('service_object');
            }
          );
        }
        this.getHostList(
          group || {
            limit: 0,
            start: 0,
            group: sessionHost ? currentHost?.group : hostsGroup[0],
          }
        );
      }
    );
  };

  /**
   * 选中主机
   * @param {*} e
   * @param {number} i
   * @memberof Host
   */
  handleSideHostClick = (e: any, i: number) => {
    // 更新选中主机
    this.props.actions.updateSelectedHostInfo(e, i);
    // 获取选中主机的服务组件
    // this.props.actions.updateServicesList(e);
  };

  /**
   *
   * @param e
   * @param group
   */
  handleSideSearch = (e, group) => {
    this.setState(
      {
        searchText: e,
      },
      () => {
        this.getHostGroups(group);
        // this.getHostList(group)
      }
    );
  };

  render() {
    const {
      actions,
      history,
      selectedIndex,
      hostList,
      hostGroupLists,
      pager,
      selectedHost,
      selectedHostServices,
      cur_parent_cluster,
    } = this.props;
    return (
      <div
        className="host-container"
        style={{ height: this.state.sideNavHeight }}>
        <SideNav
          loading={this.state.sideNavLoading}
          selectedGroup={this.state.selectedGroup}
          selectedIndex={selectedIndex}
          handleSearch={this.handleSideSearch}
          clickGroup={(groupValue) => {
            this.getHostList(groupValue);
          }}
          pager={pager}
          handleHostItemClick={(e, i) => this.handleSideHostClick(e, i)}
          style={{ height: this.state.sideNavHeight }}
          hostList={hostList}
          hostGroupLists={hostGroupLists}
        />
        <div
          className="detail-container"
          style={{ height: this.state.sideNavHeight, overflow: 'auto' }}>
          <Tabs>
            <TabPane tab="运行状态" key="1">
              <HostDetail
                detailData={selectedHost}
                runningServices={selectedHostServices}
                getServicesList={actions.updateServicesList}
                cur_parent_cluster={cur_parent_cluster}
                history={history}
                selectedHost={selectedHost}
              />
            </TabPane>
            <TabPane tab="仪表盘" key="2">
              <HostDash
                searchStr={`?host=${selectedHost?.ip}`}
                cur_parent_cluster={cur_parent_cluster}
              />
            </TabPane>
          </Tabs>
        </div>
      </div>
    );
  }
}

export default HostStatus;
