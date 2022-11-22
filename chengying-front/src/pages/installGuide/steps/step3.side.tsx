import * as React from 'react';
import { Select, Menu, Icon } from 'antd';
import { alertModal } from '@/utils/modal';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { InstallGuideActionTypes } from '@/actions/installGuideAction';
import { EnumDeployMode } from './types';
import classnames from 'classnames';
import '../style.scss';

// const Search = Input.Search;
const Option = Select.Option;
const MenuItem = Menu.Item;
const SubMenu = Menu.SubMenu;
interface Prop {
  actions?: InstallGuideActionTypes;
  productServicesInfo: any[];
  selectedProduct: any;
  setSelectedConfigService: Function;
  updateServiceHostList: Function;
  selectedService: any;
  clusterId: number;
  width: number;
  runtimeState?: string;
  deployState?: string;
  deployMode: EnumDeployMode;
  saveInstallInfo: Function;
  upgradeType?: string;
  namespace?: string;
  productName?: string;
  smoothSelectService?: any;
}

interface State {
  menuList: any[];
  openKeys: any[];
  allServices: any[];
  searchValue: string;
  searchData: any[];
  searchFlag: boolean;
  productNameFilter: string;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  runtimeState: state.InstallGuideStore.runtimeState,
  deployState: state.InstallGuideStore.deployState,
  namespace: state.InstallGuideStore.namespace,
  upgradeType: state.DeployStore.upgradeType,
  smoothSelectService: state.InstallGuideStore.smoothSelectService
});
@(connect(mapStateToProps) as any)
class StepThreeSide extends React.Component<Prop, State> {
  private container = React.createRef<HTMLDivElement>();

  state: State = {
    menuList: [],
    openKeys: [],
    allServices: [],
    searchValue: '',
    searchData: [],
    searchFlag: false,
    productNameFilter: '',
  };

  componentDidMount() {
    for (const o in this.props.productServicesInfo) {
      this.state.openKeys.indexOf(o) === -1 && this.state.openKeys.push(o); // tslint:disable-line
    }
    this.setState({
      openKeys: this.state.openKeys,
    });
    this.getAllService(this.props);
  }

  firstRender = true;

  componentDidUpdate() {
    const { productServicesInfo, deployMode } = this.props;
    // TODO:临时处理，待优化
    if (
      deployMode === EnumDeployMode.AUTO &&
      this.firstRender &&
      Array.isArray(productServicesInfo)
    ) {
      const selectedProduct = productServicesInfo[0];
      this.props.saveInstallInfo({
        product_name: selectedProduct.productName,
        product_version: selectedProduct.version,
      });
      this.firstRender = false;
    }

    const a = this.container.current.getElementsByClassName(
      'ant-menu-item-selected'
    );
    if (a.length > 0 && this.state.searchFlag) {
      this.container.current.scrollTo(
        0,
        (
          this.container.current.getElementsByClassName(
            'ant-menu-item-selected'
          )[0] as any
        ).offsetTop
      );
      this.setState({
        searchFlag: false,
      });
    }
  }

  serviceSelected = (params: any) => {
    const { runtimeState, deployState, clusterId, upgradeType, productName, smoothSelectService } = this.props;
    if (!alertModal(runtimeState, deployState)) {
      return;
    }
    let final_upgrade = false
    if (smoothSelectService?.ServiceAddr) {
      if (smoothSelectService?.ServiceAddr?.UnSelect) {
        final_upgrade = false
      } else {
        final_upgrade = true
      }
    }
    if (params.ServiceDisplay === 'mysql' && upgradeType === 'smooth') {
      this.props.actions.setSqlErro({
        product_name: productName,
        cluster_id: clusterId,
        final_upgrade: final_upgrade,
        ip: params.ServiceAddr.IP.toString()
      })
    } 
    this.props.updateServiceHostList({
      productName: this.props.selectedProduct.product_name,
      serviceName: params.serviceKey,
      clusterId,
    });
    this.props.setSelectedConfigService(params);
  };

  // 获取productInfo信息
  // 自动部署状态下获取过滤后得到的产品，手动部署直接获取productServiceInfo
  getProductInfo = ({ deployMode, productNameFilter, productServicesInfo }) => {
    if (
      deployMode === EnumDeployMode.AUTO &&
      Array.isArray(productServicesInfo)
    ) {
      if (!productNameFilter) {
        const selectedProduct = productServicesInfo[0];
        return selectedProduct.content;
      } else {
        return productServicesInfo.find(
          (product) => product.productName === productNameFilter
        )?.content;
      }
    } else {
      return productServicesInfo;
    }
  };

  getAllService = (p: Prop) => {
    const { deployMode } = p;
    const { productNameFilter } = this.state;
    const productInfo = this.getProductInfo({
      deployMode,
      productNameFilter,
      productServicesInfo: p.productServicesInfo,
    });
    const allServices = [];
    for (const o in productInfo) {
      for (const q in productInfo[o]) {
        allServices.push(q);
      }
    }
    this.setState({
      allServices,
    });
  };

  UNSAFE_componentWillReceiveProps(nextProp: Prop) {
    const { deployMode, productServicesInfo } = nextProp;
    // 初始化获取productServicesInfo后，设置默认选中第一项
    if (
      deployMode === EnumDeployMode.AUTO &&
      Array.isArray(productServicesInfo)
    ) {
      if (this.props.productServicesInfo !== productServicesInfo) {
        this.setState(
          {
            productNameFilter:
              this.state.productNameFilter ||
              productServicesInfo[0].productName,
          },
          () => {
            this.getAllService(nextProp);
            // 获取默认展开的selectedKeys
            const productInfo = this.getProductInfo({
              deployMode: deployMode,
              productNameFilter: this.state.productNameFilter,
              productServicesInfo: nextProp.productServicesInfo,
            });
            const nextOpenKeys = this.getAllOpenKeys(productInfo);
            this.setState({
              openKeys: nextOpenKeys,
            });
          }
        );
      }
    } else {
      this.getAllService(nextProp);
      if (
        JSON.stringify(nextProp.selectedService) !== '{}' &&
        this.state.openKeys.indexOf(nextProp.selectedService.Group) === -1
      ) {
        const nextOpenKeys = this.getAllOpenKeys(productServicesInfo);
        nextOpenKeys.push(nextProp.selectedService.Group);
        this.setState({
          openKeys: nextOpenKeys,
        });
      } else {
        const nextOpenKeys = this.getAllOpenKeys(nextProp.productServicesInfo);
        this.setState({
          openKeys: nextOpenKeys,
        });
      }
    }
  }

  getAllOpenKeys = (productInfo) => {
    if (!productInfo) return [];
    const nextOpenKeys = [...this.state.openKeys];
    for (const o in productInfo) {
      nextOpenKeys.indexOf(o) === -1 && nextOpenKeys.push(o); // tslint:disable-line
    }
    return nextOpenKeys;
  };

  loadMenu = (p: any) => {
    if (!p) return null;
    const keys = Object.keys(p);
    const { deployMode } = this.props;
    return keys.map((k: any, i: number) => {
      return (
        <SubMenu
          className={classnames({
            'pl-20': deployMode === EnumDeployMode.AUTO,
          })}
          key={k}
          title={k}>
          {this.loadSonMenu(p[k], i)}
        </SubMenu>
      );
    });
  };

  loadSonMenu = (sonMenu: any, l: number) => {
    const keys = Object.keys(sonMenu);
    return keys.map((k: any, i: number) => {
      return (
        <MenuItem
          onClick={() => this.serviceSelected({ ...sonMenu[k], serviceKey: k })}
          key={k}>
          {sonMenu[k].ServiceDisplay}
        </MenuItem>
      );
    });
  };

  handleSearch = (e: any) => {
    let { searchData } = this.state;
    searchData = [];
    this.state.allServices.forEach((o: string) => {
      if (o.indexOf(e) !== -1) {
        searchData.push(o);
      }
    });
    this.setState({
      searchData,
    });
  };

  handleOpenChange = (e?: any) => {
    this.setState({
      openKeys: e,
    });
  };

  handleChange = (e: string) => {
    this.setState({ searchValue: e }, () => {
      if (this.state.allServices.indexOf(e) !== -1) {
        this.locateSearch(e);
      }
    });
  };

  locateSearch = (e: string) => {
    const { deployMode, productServicesInfo } = this.props;
    const { productNameFilter } = this.state;
    const list = this.getProductInfo({
      deployMode,
      productNameFilter,
      productServicesInfo,
    });
    const keys = Object.keys(list);
    keys.forEach((o: string) => {
      const sonKeys = Object.keys(list[o]);
      sonKeys.forEach((q: string) => {
        if (q === e) {
          this.serviceSelected({ ...list[o][q], serviceKey: q });
          this.setState({
            openKeys: this.state.openKeys.concat([o]),
            searchFlag: true,
          });
        }
      });
    });
  };

  handleProdutNameFilter = (filter) => {
    const { deployState, runtimeState } = this.props;
    if (!alertModal(runtimeState, deployState)) return;
    const { productServicesInfo } = this.props;
    const selectedProduct = productServicesInfo.find(
      (product) => product.productName === filter.productName
    );
    const nextOpenKeys = this.getAllOpenKeys(selectedProduct.content);
    this.setState(
      {
        productNameFilter: filter.productName,
        openKeys: nextOpenKeys,
      },
      () => {
        this.props.saveInstallInfo({
          product_name: selectedProduct.productName,
          product_version: selectedProduct.version,
        });
        this.props.setSelectedConfigService({});
      }
    );
  };

  render() {
    const { productServicesInfo, deployMode } = this.props;
    const { productNameFilter } = this.state;
    console.log(productNameFilter)
    console.log(productServicesInfo)
    const options = this.state.searchData.map((o: any) => (
      <Option key={o} value={o}>
        {o}
      </Option>
    ));
    return (
      <div style={{display: 'flex'}}>
        {deployMode === EnumDeployMode.AUTO && Array.isArray(productServicesInfo) && (
        <div className='stepThreeProductList'>
            <div className='stepThreeProductListTop'>组件</div>
            {productServicesInfo?.map(item => {
              return (
                <div 
                    onClick={() => this.handleProdutNameFilter(item)}
                    className={item.productName == productNameFilter ? 'stepThreeListItem activeItem' : 'stepThreeListItem'}
                    key={item.productName} > {item.productName} {item.version}
                </div>
              )
            })}
        </div> )}
      <div
        style={{ width: this.props.width }}
        className="step-three-side-container">
        <div
          className="search-box"
          style={{
            padding: '14px 20px 10px',
            borderRight: '1px solid #DDDDDD',
            position: 'relative',
          }}>
          <div
            style={{
              position: 'relative',
              width: '100%',
            }}>
            <Select
              showSearch
              value={this.state.searchValue || undefined}
              placeholder={'输入服务名'}
              style={{ fontSize: 12, width: '100%', height: 28 }}
              defaultActiveFirstOption={false}
              showArrow={false}
              filterOption={false}
              onSearch={this.handleSearch}
              onChange={this.handleChange}
              notFoundContent={null}>
              {options}
            </Select>
            <Icon
              style={{
                position: 'absolute',
                top: '8px',
                right: '5px',
                transform: 'translate(-50%,0)',
                color: '#999999',
              }}
              type="search"
            />
          </div>
        </div>
        <div
          ref={this.container}
          className="menu-box"
          style={{ overflow: 'scroll', borderRight: '1px solid #dddddd' }}>
          <Menu
            className="c-sidenav__menu"
            onOpenChange={(e) => this.handleOpenChange(e)}
            openKeys={this.state.openKeys}
            selectedKeys={[this.props.selectedService.serviceKey]}
            mode="inline">
            {deployMode === EnumDeployMode.MANUAL &&
              this.loadMenu(productServicesInfo)}
            {deployMode === EnumDeployMode.AUTO &&
              Array.isArray(productServicesInfo) &&
              this.loadMenu(
                productServicesInfo.find(
                  (product) => product.productName === productNameFilter
                )?.content
              )}
          </Menu>
        </div>
      </div>
      </div>
    );
  }
}

export default StepThreeSide;
