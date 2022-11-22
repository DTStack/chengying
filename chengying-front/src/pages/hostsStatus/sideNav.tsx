import * as React from 'react';
import { Input, Menu, Spin, Layout } from 'antd';
const { Sider } = Layout;
const SubMenu = Menu.SubMenu;
const Search = Input.Search;
const pageSize = 0;

interface Prop {
  hostList: any[];
  selectedGroup: string;
  pager: any;
  hostGroupLists: any[];
  style: any;
  handleHostItemClick: (e: any, i: number) => void;
  handleSearch: (e: any, group: any) => void;
  clickGroup: any;
  selectedIndex: number;
  loading: boolean;
}

interface State {
  openKeys: any[];
  selectedKeys: string[];
  limit: number;
  start: number;
}

class SideNav extends React.Component<Prop, State> {
  public state: State = {
    openKeys: [],
    selectedKeys: [],
    limit: pageSize,
    start: 0,
  };

  componentDidUpdate(prevProps, prevState) {
    const { hostGroupLists, selectedIndex, selectedGroup } = this.props;
    if (prevProps.selectedIndex !== selectedIndex) {
      this.setState({
        openKeys:
          prevProps.selectedIndex !== selectedIndex
            ? this.state.openKeys
            : [`${hostGroupLists[0]}`],
        selectedKeys:
          prevProps.selectedIndex !== selectedIndex
            ? [`${this.state.openKeys[0]}-${selectedIndex}`]
            : [`${hostGroupLists[0]}-${selectedIndex}`],
      });
    }
    // 兼容集群命令跳转到主机的定位
    if (selectedGroup && prevProps.selectedGroup !== selectedGroup) {
      this.setState({
        openKeys: [selectedGroup],
        selectedKeys: [`${selectedGroup}-${selectedIndex}`],
      });
    }
  }

  onOpenChange = (openKeys) => {
    console.log(openKeys);
    const { hostList, handleHostItemClick, hostGroupLists, clickGroup } =
      this.props;
    const latestOpenKey = openKeys.find(
      (key) => this.state.openKeys.indexOf(key) === -1
    );
    if (hostGroupLists && hostGroupLists.indexOf(latestOpenKey) > -1) {
      openKeys = latestOpenKey ? [latestOpenKey] : [];
    }
    this.setState(
      {
        openKeys,
        selectedKeys: openKeys.length ? [`${openKeys[0]}-${0}`] : [],
      },
      () => {
        if (!openKeys.length) {
          return;
        }
        const defaultHost = hostList[0];
        clickGroup({
          limit: pageSize,
          start: 0,
          group: openKeys[0],
        });
        handleHostItemClick(defaultHost, 0);
      }
    );
  };

  render() {
    const { hostList, hostGroupLists, loading } = this.props;
    return (
      <Sider
        className="host-side-nav box-shadow-style"
        style={this.props.style}
        width={240}>
        <div className="search-box">
          <Search
            onSearch={(e) =>
              this.props.handleSearch(e, { group: this.state.openKeys[0] })
            }
            style={{ width: 200, margin: '12px 20px 8px' }}
            placeholder="输入主机ip"
          />
        </div>
        <Spin className="host-box" spinning={loading}>
          <Menu
            mode="inline"
            selectedKeys={this.state.selectedKeys}
            openKeys={this.state.openKeys}
            onOpenChange={this.onOpenChange}
            className="c-sidenav__menu">
            {hostGroupLists &&
              hostGroupLists.map((item: any) => {
                return (
                  <SubMenu
                    key={`${item}`}
                    title={<span>{item}</span>}
                    // onTitleClick={(params) => {
                    //   console.log(params)
                    //   this.props.clickGroup({
                    //     limit: pageSize,
                    //     start: 0,
                    //     group: params.key,
                    //   });
                    // }}
                  >
                    {hostList &&
                      hostList.map((o: any, index: number) => (
                        <Menu.Item
                          key={`${item}-${index}`}
                          onClick={(e: any) => {
                            console.log(e);
                            this.setState({ selectedKeys: e.keyPath });
                            this.props.handleHostItemClick(o, index);
                          }}>
                          <i
                            className="dot-cls"
                            style={
                              o.alert
                                ? { display: 'inline-block' }
                                : { display: 'none' }
                            }
                          />
                          {o.ip}
                        </Menu.Item>
                      ))}
                  </SubMenu>
                );
              })}
          </Menu>
        </Spin>
      </Sider>
    );
  }
}

export default SideNav;
