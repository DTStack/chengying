import * as React from 'react';
import { connect } from 'react-redux';
import { Select, message, Input, Icon, Button } from 'antd';
import { uniqBy, uniq } from 'lodash';
import { Service, servicePageService } from '@/services';
import { AppStoreTypes } from '@/stores';
import ComponentList from './componentList';
import DeployHistory from './deployHistory';
import UpdatePatchHistory from './updataPatchHistory';
import DownloadModal from './components/downloadModal';
import UpdatePatchModal from './components/updatePatchModal';
import utils from '@/utils/utils';
import * as Cookie from 'js-cookie';
import './style.scss';

const Option = Select.Option;
const Search = Input.Search;
// const TabPane = Tabs.TabPane;

interface Props {
  location?: any;
  history?: any;
  extraContent?: any;
  authorityList?: any;
  HeaderStore?: any;
}

export interface QueryParams {
  clusterId?: string;
  parentProductName?: string;
  productName?: string[];
  productVersion?: string;
  deploy_status?: any;
  start?: number;
  limit?: number;
  'sort-by'?: string;
  'sort-dir'?: string;
}

interface State {
  clusterList: any[];
  productList: any[]; // 产品列表
  /** 组件列表 */
  componentList: any[];
  loading: boolean;
  searchParam: QueryParams;
  downModalVisible: boolean; // 下载产品部署内容
  updateModalVisible: boolean; // 上传补丁包
  defaultValue: any;
  clusterType: string;
  deployType: string;
  uploadPack: number;
  bodyWidth: number;
}
const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
  HeaderStore: state.HeaderStore,
});
@(connect(mapStateToProps) as any)
class ComponentContainer extends React.Component<Props, State> {
  state: State = {
    clusterList: [],
    productList: [],
    componentList: [],
    loading: false,
    searchParam: {
      productName: undefined, // 组件名称
      clusterId: undefined,
      parentProductName: undefined, // 产品名称
      productVersion: undefined,
    },
    downModalVisible: false,
    updateModalVisible: false,
    defaultValue: undefined,
    clusterType: '',
    deployType: '',
    uploadPack: Date.now(),
    bodyWidth: document.body.clientWidth,
  };

  componentDidMount() {
    const { location, HeaderStore } = this.props;
    const { cur_parent_cluster } = HeaderStore;
    let localSession =
      JSON.parse(sessionStorage.getItem('service_object')) || {};
    let clusterId =
      cur_parent_cluster?.id > 0
        ? cur_parent_cluster?.id
        : +Cookie.get('em_current_cluster_id');
    this.getParentClustersList(clusterId, localSession.parentProductName);
    // sessionStorage.removeItem('service_object')
    let routerArrs = location?.pathname.split('/');
    let deployType = routerArrs?.length
      ? routerArrs[routerArrs?.length - 1]
      : '';
    if (deployType) {
      this.setState({ deployType });
    }
  }
 
  // transfer list for the reason of introduce parameter namespace
  transferList = (cluster: any, inputParentProductName?: string) => {
    const defaultProductName = 'DTinsight';
    function getParentProductName(productList) {
      let isParentProductNameExist = false;
      let isDefaultParentProductNameExist = false;
      productList.forEach((item) => {
        item === inputParentProductName && (isParentProductNameExist = true);
        item === defaultProductName && (isDefaultParentProductNameExist = true);
      });
      return isParentProductNameExist
        ? inputParentProductName
        : isDefaultParentProductNameExist
        ? defaultProductName
        : productList[0];
    }
    let productList = [];
    const subdomain = cluster.subdomain || {};
    if (cluster.mode === 0) {
      productList = subdomain.products || [];
    } else {
      Object.keys(subdomain).forEach((namespace) => {
        subdomain[namespace].forEach((p) => productList.push(p));
      });
    }
    return {
      productList: uniq(productList),
      productName: getParentProductName(productList),
    };
  };

  // 获取集群列表
  getParentClustersList = (clusterId?: number, parentProductName?: string) => {
    Service.getClusterProductList().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        if (res.data && res.data.length > 0) {
          const cluster =
            (clusterId &&
              res.data.find((item: any) => item.clusterId === clusterId)) ||
            res.data[0] ||
            {};
          const { productList, productName } = this.transferList(
            cluster,
            parentProductName
          );
          this.setState(
            {
              clusterType: cluster.clusterType,
              clusterList: res.data,
              productList,
              searchParam: {
                clusterId: clusterId || cluster.clusterId,
                parentProductName: productName,
                productName: undefined,
                productVersion: undefined,
              },
            },
            () => {
              this.getProductComponents({
                clusterId: clusterId || cluster.clusterId,
                parentProductName: productName,
              });
            }
          );
        } else {
          this.resetPages();
        }
      } else {
        message.error(res.msg);
      }
    });
  };

  // 重置页面数据
  resetPages = () => {
    this.setState({
      clusterList: [],
      productList: [],
      componentList: [],
      loading: false,
      searchParam: {
        productName: undefined, // 组件名称
        clusterId: undefined,
        parentProductName: undefined, // 产品名称
        productVersion: undefined,
      },
    });
  };

  // 获取已部署产品组件列表
  getProductComponents = (params?: any) => {
    const reqParams = Object.assign({}, this.state.searchParam, params);
    reqParams.mode = this.currentCluster.mode;
    servicePageService.getProductName({clusterId: params?.clusterId, mode: 0, parentProductName: params?.parentProductName}).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const data = res.data.list;
        let arr = [];
        arr = uniqBy(data, 'product_name');
        this.setState({
          componentList: arr,
          searchParam: reqParams,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 选择集群
  handleClusterChange = (clusterId: string) => {
    const { clusterList } = this.state;
    const cluster =
      clusterList.find((item) => item.clusterId === clusterId) || {};
    const { productList, productName: parentProductName } =
      this.transferList(cluster);
    this.setState(
      {
        searchParam: Object.assign({}, this.state.searchParam, {
          clusterId,
          parentProductName,
          productName: undefined,
          productVersion: undefined,
          componentList: [],
        }),
        productList,
        clusterType: cluster.clusterType,
      },
      () => {
        this.getProductComponents();
      }
    );
  };

  // 选择产品
  handleSelectChange = (value: string) => {
    const newState = Object.assign({}, this.state.searchParam, {
      parentProductName: value,
      productName: undefined,
      productVersion: undefined,
      componentList: [],
    });
    this.setState({ searchParam: newState }, () => {
      this.getProductComponents();
    });
  };

  handleComponentSearch = (value: any) => {
    const newState = Object.assign({}, this.state.searchParam, {
      productVersion: value,
    });
    this.setState({ searchParam: newState });
  };

  onComponentChange = (key) => {
    this.setState({
      searchParam: Object.assign({}, this.state.searchParam, {
        productName: key,
        deploy_status: '',
      }),
    });
  };

  // 下载产品部署日志
  handleDownModalShow = (): any => {
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'package_download')) {
      return false;
    }
    this.setState({
      downModalVisible: !this.state.downModalVisible,
    });
  };

  closeModal = () => {
    this.setState({
      defaultValue: undefined,
      updateModalVisible: false,
    });
  };

  openModal = (): any => {
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'patches_update')) {
      return false;
    }
    this.setState({
      updateModalVisible: true,
    });
  };

  changeDefaultValue = (value) => {
    this.setState(
      {
        defaultValue: value,
      },
      () => this.openModal()
    );
  };

  resetKey = () => {
    this.setState({
      uploadPack: Date.now(),
    });
  };

  get currentCluster() {
    const { clusterList, searchParam } = this.state;
    const { clusterId } = searchParam;
    let ret;
    clusterList.forEach((item) => {
      if (item.clusterId === clusterId) {
        ret = item;
      }
    });
    return ret || {};
  }

  get shouldNameSpaceShow() {
    return this.currentCluster.mode === 1;
  }

  render() {
    const { HeaderStore } = this.props;
    const { cur_parent_cluster } = HeaderStore;
    const {
      componentList,
      searchParam,
      productList,
      clusterList,
      downModalVisible,
      updateModalVisible,
      defaultValue,
      clusterType,
      uploadPack,
      deployType,
      bodyWidth,
    } = this.state;
    const paneKey = `${searchParam.clusterId}-${searchParam.parentProductName}-${searchParam.productName}-${searchParam.productVersion}-${uploadPack}`;
    // console.log('searchParam', clusterType, HeaderStore)
    return (
      <div className="content-container">
        <div className="mb-12 clearfix">
          <div
            style={{
              marginBottom: `${bodyWidth <= 1440 ? '10px' : '0px'}`,
              display: 'inline-block',
            }}>
            <span className="mr-20">
              集群：
              <Select
                className="dt-form-shadow-bg"
                style={{ width: 180 }}
                size="default"
                placeholder="选择集群"
                disabled
                value={
                  clusterList.includes(searchParam.clusterId)
                    ? searchParam.clusterId
                    : cur_parent_cluster.name
                } // 兼容无部署组件的集群显示
                onChange={this.handleClusterChange}>
                {Array.isArray(clusterList) &&
                  clusterList.map((item: any, index: number) => (
                    <Option
                      data-testid={`cluster-option-${item.clusterName}`}
                      key={`${index}`}
                      value={item.clusterId}>
                      <Icon type="appstore-o" style={{ marginRight: '6px' }} />
                      {item.clusterName}
                    </Option>
                  ))}
              </Select>
            </span>
            <span className="mr-20">
              产品：
              <Select
                className="dt-form-shadow-bg"
                style={{ width: 180, fontSize: 12 }}
                size="default"
                placeholder="选择产品"
                value={searchParam.parentProductName}
                onChange={this.handleSelectChange}>
                {Array.isArray(productList) &&
                  productList.map((item: any, index: number) => (
                    <Option key={`${index}`} value={item}>
                      <Icon type="appstore-o" style={{ marginRight: '6px' }} />
                      {item}
                    </Option>
                  ))}
              </Select>
            </span>
            <span className="mr-20">
              组件名称：
              <Select
                className="dt-form-shadow-bg"
                mode="multiple"
                size="default"
                placeholder="选择组件"
                value={this.state.searchParam.productName}
                style={{ width: 180 }}
                onChange={this.onComponentChange}>
                {Array.isArray(componentList) &&
                  componentList.map((o: any) => (
                    <Option key={o.product_name}>
                      {o.product_name}
                    </Option>
                  ))}
              </Select>
            </span>
            <span className="mr-20">
              <Search
                className="dt-form-shadow-bg"
                style={{ width: 264 }}
                placeholder="按组件版本号搜索"
                onSearch={this.handleComponentSearch}
              />
            </span>
          </div>
          {cur_parent_cluster?.type === 'hosts' && (
            <Button
              className={'fl-r'}
              type="primary"
              style={{ marginLeft: '20px' }}
              onClick={this.openModal}>
              上传补丁包
            </Button>
          )}
          <Button
            className={
              // bodyWidth === 1440 ?
              // "mr-20" :
              'fl-r'
            }
            type="primary"
            icon="download"
            onClick={this.handleDownModalShow}>
            下载产品部署内容
          </Button>
        </div>
        {deployType === 'deployed' && (
          <ComponentList
            key={'component-list' + paneKey}
            {...this.state.searchParam}
            {...this.props}
            shouldNameSpaceShow={this.shouldNameSpaceShow}
            mode={this.currentCluster.mode}
            clusterList={clusterList}
            getParentClustersList={this.getParentClustersList}
          />
        )}
        {deployType === 'history' && (
          <DeployHistory
            shouldNameSpaceShow={this.shouldNameSpaceShow}
            key={'history' + paneKey}
            {...this.state.searchParam}
          />
        )}
        {deployType === 'patchHistory' &&
          !(
            cur_parent_cluster?.type === 'kubernetes' &&
            cur_parent_cluster?.mode === 0
          ) && (
            <UpdatePatchHistory
              key={'patches-update-history' + paneKey}
              {...this.state.searchParam}
              changeDefaultValue={this.changeDefaultValue}
              shouldNameSpaceShow={this.shouldNameSpaceShow}
            />
          )}
        {downModalVisible && (
          <DownloadModal
            visible={downModalVisible}
            defaultValue={{
              clusterId: searchParam.clusterId
                ? Number(searchParam.clusterId)
                : undefined,
              parentProductName: searchParam.parentProductName,
            }}
            data={clusterList}
            onCancel={this.handleDownModalShow}
          />
        )}
        {
          <UpdatePatchModal
            visible={updateModalVisible}
            onCancel={this.closeModal}
            data={componentList}
            defaultValue={defaultValue}
            resetKey={this.resetKey}
          />
        }
      </div>
    );
  }
}

export default ComponentContainer;
