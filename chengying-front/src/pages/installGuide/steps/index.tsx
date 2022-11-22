import * as React from 'react';
import { Steps, Button, Modal, message, notification, Icon } from 'antd';
import { connect } from 'react-redux';
import { Dispatch, bindActionCreators } from 'redux';
import { AppStoreTypes } from '@/stores';
import { DeployStore, InstallGuideStore } from '@/stores/modals';
import * as HeaderAction from '@/actions/headerAction';
import { HeaderStateTypes } from '@/stores/headerReducer';
import * as installGuideAction from '@/actions/installGuideAction';
import StepOne from './step1';
import StepTwo from './step2';
import StepThree from './step3';
import StepFour from './step4';
import * as Cookie from 'js-cookie';
import { deployService, installGuideService } from '@/services';
import { EnumDeployMode } from './types';

import '../style.scss';
import { UserCenterStoreTypes } from '@/stores/userCenterReducer';
import { alertModal } from '@/utils/modal';
import utils from '@/utils/utils';

const mapStateToProps = (state: AppStoreTypes) => {
  return {
    installGuideProp: state.InstallGuideStore,
    deployProp: state.DeployStore,
    userCenterProp: state.UserCenterStore,
    headerStore: state.HeaderStore,
  };
};
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(
    Object.assign({}, installGuideAction, HeaderAction),
    dispatch
  ),
});

const Step = Steps.Step;

interface StepIndexProp {
  installGuideProp: InstallGuideStore;
  deployProp: DeployStore;
  actions: installGuideAction.InstallGuideActionTypes &
    HeaderAction.HeaderActionTypes;
  router?: any;
  location?: any;
  history?: any;
  userCenterProp: UserCenterStoreTypes;
  headerStore: HeaderStateTypes;
}
interface StepIndexState {
  urlParams: any;
  gobackLocation: string;
  canUpgrade: boolean;
  regionPID: number;
  autoProductList: any[];
  deployMode: EnumDeployMode;
  autoExpendRowKeys: string[];
  autoSelectedProducts: any[];
  deployUUID: string;
  upgradeStep: number;
  isUpgrade: boolean;
}

const parser = (url: string): any => {
  if (!url) return {};
  return url
    .split('?')[1]
    .split('&')
    .reduce((temp, cur) => {
      const [key, value] = cur.split('=');
      temp[key] = value;
      return temp;
    }, {});
};

@(connect(mapStateToProps, mapDispatchToProps) as any)
class StepIndex extends React.Component<StepIndexProp, StepIndexState> {
  constructor(props: any) {
    super(props);
  }

  state: StepIndexState = {
    urlParams: {},
    gobackLocation: '',
    canUpgrade: true,
    regionPID: -1,
    // 自动部署
    autoProductList: [], // 自动部署产品包收集
    deployMode: EnumDeployMode.AUTO,
    autoExpendRowKeys: [],
    autoSelectedProducts: [],
    deployUUID: '',
    upgradeStep: 0,
    isUpgrade: false,
  };

  private stepTwoForm: any = null;
  componentDidMount() {
    this.handleUrlParamSearch();
    // console.log(this.props.location.search);
    const qo = parser(this.props.location.search);
    if (
      new RegExp('/deploycenter').test(qo.from) &&
      qo.product_name !== undefined
    ) {
      this.setState({
        deployMode: EnumDeployMode.MANUAL,
      });
    } else {
      this.setState(
        {
          deployMode: EnumDeployMode.AUTO,
        },
        () => {
          this.getOrchestrationHistory();
        }
      );
    }
  }

  componentWillUnmount() {
    this.props.actions.initInstallGuide();
  }

  static getDerivedStateFromProps(
    nextProps: StepIndexProp,
    prevState: StepIndexState
  ) {
    const { selectedProduct, clusterId, namespace, step } =
      nextProps.installGuideProp;
    const { urlParams = {}, regionPID } = prevState;
    const { new_version, redeploy, product_name, product_version, id } =
      urlParams;
    // 升级下，产品包选择有变更，不允许下一步
    if (new_version && step === 1 && Object.keys(selectedProduct).length) {
      const selectedPID = selectedProduct?.id;
      // 判断产品包是否一致
      const productCheck =
        `${selectedPID}` === id &&
        selectedProduct?.product_name === product_name &&
        selectedProduct?.product_name === product_version;
      // 判断集群是否一致
      const clusterCheck =
        `${clusterId}` === redeploy &&
        (namespace || undefined) === urlParams?.namespace;
      if (!(productCheck && clusterCheck)) {
        if (selectedPID !== regionPID) {
          message.error('与之前选择不一致，请重新选择');
          return {
            regionPID: selectedPID,
            canUpgrade: false,
          };
        }
        return null;
      }
    }
    return {
      canUpgrade: true,
    };
  }

  getOrchestrationHistory = async () => {
    const { clusterId } = this.props.installGuideProp;
    const res = await deployService.getOrchestrationHistory({
      cluster_id: clusterId,
    });
    const { data = [], code } = res.data;

    if (code === 0) {
      // 整合
      const nextAutoSelectedProducts = Array.isArray(data)
        ? data.map((product) => {
            const sortedData = product.service.all_service.sort(
              (a, b) => a.baseProduct.length - b.baseProduct.length
            );
            const disabled = sortedData
              .filter((item) => item.baseProduct !== '')
              .map((item) => item.serviceName);
            const { check_service = [], uncheck_service = [] } =
              product.service;
            const checked = check_service.map((item) => item.service_name);
            const unChecked = uncheck_service.map((item) => item.service_name);

            return {
              ID: product.pid,
              productName: product.product_name,
              service: {
                all: sortedData,
                checked,
                unChecked,
                disabled,
              },
            };
          })
        : [];
      this.setState({
        autoSelectedProducts: nextAutoSelectedProducts,
      });
    }
  };

  refreshDeployService = (callback?) => {
    const { autoSelectedProducts } = this.state;
    const params = {
      cluster_id: this.props.installGuideProp.clusterId,
      product_info: autoSelectedProducts.map((item) => {
        return {
          id: item.ID,
          name: item.productName,
          service_list: item.service.checked,
          uncheck_service_list: item.service.unChecked,
        };
      }),
    };
    this.props.actions.refreshServicesInfoForAutoDeploy(params, callback);
  };

  // 自动部署更新service信息
  updateAutoDeployService = (productID, checkedService) => {
    const { autoSelectedProducts } = this.state;
    const nextAutoSelectedProduct = autoSelectedProducts.map((product) => {
      if (productID === product.ID) {
        const unChecked = [];
        product.service.all.forEach((service) => {
          if (
            checkedService.findIndex(
              (checkedServiceName) => checkedServiceName === service.serviceName
            ) === -1
          ) {
            unChecked.push(service.serviceName);
          }
        });
        return {
          ...product,
          service: {
            all: product.service.all,
            checked: checkedService,
            unChecked,
            disabled: product.service.disabled,
          },
        };
      }
      return product;
    });
    this.setState({
      autoSelectedProducts: nextAutoSelectedProduct,
    });
  };

  initStep = () => {
    const {
      query_str,
      cluster_id,
      product_name,
      product_version,
      id,
      redeploy,
      type,
      namespace,
      new_version,
    } = this.state.urlParams;
    const installType = type || 'hosts';
    this.props.actions.initInstallGuide();
    this.props.actions.saveInstallType(installType);
    if (query_str && query_str.length > 0) {
      if (installType === 'kubernetes') {
        this.props.actions.saveSelectCluster(Number(cluster_id));
        this.props.actions.saveSelectNamespace({}, { namespace }, false);
      }
      // 第四种情况只需要名字和版本即可
      this.props.actions.saveInstallInfo({
        product_name: product_name,
        product_version: product_version,
        id: id,
      });
      this.props.actions.setDeployUUID(query_str);
      this.props.actions.goToStep(3);
    } else {
      this.props.actions.getInstallClusterList({
        limit: 0,
        type: installType,
      });

      if (redeploy) {
        Cookie.set('em_current_cluster_id', redeploy);
        Cookie.set('em_current_cluster_type', installType);
        this.props.actions.saveSelectCluster(Number(redeploy));
        // 获取namespace
        if (installType === 'kubernetes') {
          this.props.actions.getNamespaceList({ clusterId: Number(redeploy) });
          namespace &&
            this.props.actions.saveSelectNamespace({}, { namespace }, false);
        }
        // 设置为第二步
        // 如果是升级，直接设为第三步
        if (new_version) {
          this.setState({ isUpgrade: true }, () => {
            this.goToStepThree();
          });
          return;
        }
        this.props.actions.goToStep(1);
        // 获取第二步的产品包列表
        this.props.actions.getProductStepOneList(
          {
            product_line_name: '',
            product_line_version: '',
            product_type: installType === 'hosts' ? 0 : 1,
            deploy_status: '',
          },
          false
        );
      } else {
        sessionStorage.removeItem('upgradeType');
        sessionStorage.removeItem('forcedUpgrade');
        sessionStorage.removeItem('isFirstSmooth');
      }
    }
  };

  // 跳转到第三步的请求处理
  goToStepThree = () => {
    const { type } = this.state.urlParams;
    const installType = type || 'hosts';
    // 获取第二步的产品包列表，并获取选中的包，及其服务
    this.props.actions.getProductPackageList(
      {
        limit: 0,
        product_type: installType === 'hosts' ? 0 : 1,
      },
      this.getSelectedProduct
    );
  };

  // 获取选中的产品包
  getSelectedProduct = () => {
    const { product_version, product_name, id } = this.state.urlParams;
    const { productPackageList = [] } = this.props.installGuideProp;
    productPackageList.forEach((o: any) => {
      if (
        product_name === o.ProductName &&
        product_version === o.ProductVersion &&
        id === o.ID.toString()
      ) {
        this.handleProductSelected({
          product_name: o.ProductName,
          product_version: o.ProductVersion,
          id: o.ID,
          status: o.Status,
          product_type: o.ProductType,
          namespace: o.Namespace,
          parent_product_name: o.ParentProductName,
          product_name_display: o.ProductNameDisplay,
          deploy_uuid: o.DeployUUID,
        });
      }
    });
  };

  // 选中的产品包下服务
  handleProductSelected = (record: any) => {
    this.getProductPackageServices(record);
    // 保存选中的产品包信息
    this.props.actions.saveInstallInfo(record);
  };

  // 获取选中的服务以及未选中的服务
  getProductPackageServices = (
    selectedProduct: any,
    baseClusterId?: number
  ) => {
    this.props.actions.getProductPackageServices(
      {
        productName:
          selectedProduct.product_name ?? selectedProduct.ProductName,
        productVersion:
          selectedProduct.product_version ?? selectedProduct.ProductVersion,
        pid: selectedProduct.id ?? selectedProduct.ID,
        clusterId: this.props.installGuideProp.clusterId,
        baseClusterId: baseClusterId,
      },
      () => {
        // 获取第三部的服务组信息
        this.getProductServicesInfo();
        // 跳转第三步
        this.props.actions.goToStep(2);
      }
    );
  };

  // 获取服务组信息
  getProductServicesInfo = () => {
    const { deployMode, urlParams } = this.state;
    if (deployMode === EnumDeployMode.AUTO) return;
    const {
      baseClusterId,
      selectedProduct,
      unSelectedServiceList,
      namespace,
      clusterId,
    } = this.props.installGuideProp;
    const { upgradeType, forcedUpgrade } = this.props.deployProp;
    // 产品包确认后，获取第三步左侧选择栏的数据（产品下的服务组信息）
    if (!selectedProduct.product_name || !selectedProduct.product_version) {
      return;
    }
    let params = {
      productName: selectedProduct.product_name,
      productVersion: selectedProduct.product_version,
      unSelectService: unSelectedServiceList,
      relynamespace: baseClusterId === -1 ? undefined : baseClusterId,
      namespace,
      clusterId,
    };
    if (upgradeType === 'smooth') {
      Object.assign(params, { upgrade_mode: 'smooth' });
    }
    this.props.actions.getProductServicesInfo(params, null, forcedUpgrade);
  };

  // componentWillReceiveProps(n: StepIndexProp) { }

  handleLastStep = () => {
    const { runtimeState, deployState } = this.props.installGuideProp;
    const { setSelectedConfigService } = this.props.actions;
    if (!alertModal(runtimeState, deployState)) return;
    this.props.actions.lastStep();
    setSelectedConfigService({});
  };

  productPackageValidator = () => {
    const {
      installType,
      selectedProduct,
      selectedServiceList,
      baseClusterId,
      baseClusterInfo,
    } = this.props.installGuideProp;
    const { autoSelectedProducts, deployMode } = this.state;
    const validateStatus = {
      status: true,
      message: '',
    };
    if (deployMode === EnumDeployMode.AUTO) {
      // 自动部署，校验autoProductList
      if (autoSelectedProducts.length === 0) {
        validateStatus.status = false;
        validateStatus.message = '请选择要安装的产品';
      } else if (
        autoSelectedProducts.some(
          (product) => product.service.checked.length === 0
        )
      ) {
        validateStatus.status = false;
        validateStatus.message = '请选择要安装的服务';
      }
    } else {
      if (JSON.stringify(selectedProduct) === '{}') {
        validateStatus.status = false;
        validateStatus.message = '请选择要安装的产品';
      } else if (
        installType === 'kubernetes' &&
        baseClusterId === -1 &&
        baseClusterInfo.hasDepends
      ) {
        validateStatus.status = false;
        validateStatus.message = '请选择依赖集群';
      } else if (selectedServiceList.length === 0) {
        validateStatus.status = false;
        validateStatus.message = '请选择要安装的服务';
      }
    }
    return validateStatus;
  };

  handleNextStep = async () => {
    const {
      step,
      installType,
      selectedProduct,
      unSelectedServiceList,
      namespace,
      clusterId,
      baseClusterId,
      deployState,
      runtimeState,
      smoothSelectService,
    } = this.props.installGuideProp;
    const { upgradeType } = this.props.deployProp;
    const { autoSelectedProducts, deployMode, upgradeStep, isUpgrade } =
      this.state;
    // 选择集群的校验处理
    if (step === 0) {
      return this.stepTwoForm.props.form.validateFields(
        async (err: any, values: any) => {
          if (err) {
            return;
          }
          if (installType !== 'kubernetes') {
            this.props.actions.nextStep();
            return;
          }
          const urlParams = { clusterId };
          const reqParams = { clusterId, namespace: values.namespace };
          // 在保存命名空间，成功后下一步
          this.props.actions.saveSelectNamespace(
            urlParams,
            reqParams,
            values.isNewSpace, // 是新建还是从原有的取
            () => {
              // values.isNewSpace && this.props.actions.getNamespaceList(urlParams);
              // 是否有默认镜像仓库
              this.props.actions.checkDefaultImageStore(
                urlParams,
                (exist: boolean) => {
                  if (!exist) {
                    notification.error({
                      message: '该集群中无默认镜像仓库',
                      description: (
                        <span>
                          请
                          <a href="deploycenter/cluster/detail/imagestore">
                            前往集群管理
                          </a>
                          设置默认镜像仓库
                        </span>
                      ),
                      className: 'check-default-notification',
                      duration: null,
                    });
                    return;
                  }
                  this.props.actions.nextStep();
                }
              );
            }
          );
        }
      );
    } else if (step === 1) {
      // 产品包确认后，获取第三步左侧选择栏的数据（产品下的服务组信息）
      const validate = this.productPackageValidator();
      if (!validate.status) return message.error(validate.message);
      if (deployMode === EnumDeployMode.MANUAL) {
        const res = await installGuideService.deployCondition({
          cluster_id: this.props.installGuideProp.clusterId,
          auto_deploy: false,
          product_name:
            this.props.installGuideProp.selectedProduct.product_name,
          product_type: installType === 'hosts' ? 0 : 1,
        });
        if (res.data.code !== 0) {
          if (res.data.msg.indexOf('存在平滑升级版本') !== -1) {
            Modal.error({
              title: '提示',
              content: res.data.msg,
            });
            return;
          }
          message.error(res.data.msg);
          return;
        }
        this.getProductServicesInfo();
      } else if (deployMode === EnumDeployMode.AUTO) {
        const res = await installGuideService.deployCondition({
          cluster_id: this.props.installGuideProp.clusterId,
          auto_deploy: true,
          product_line_name:
            this.props.installGuideProp.selectProductLine?.product_line_name,
          product_line_version:
            this.props.installGuideProp.selectProductLine?.product_line_version,
          product_type: installType === 'hosts' ? 0 : 1,
        });
        if (res.data.code !== 0) {
          if (res.data.msg.indexOf('存在平滑升级版本') !== -1) {
            Modal.error({
              title: '提示',
              content: res.data.msg,
            });
            return;
          }
          message.error(res.data.msg);
          return;
        }
        const params = {
          product_line_name:
            this.props.installGuideProp.selectProductLine?.product_line_name,
          product_line_version:
            this.props.installGuideProp.selectProductLine?.product_line_version,
          cluster_id: this.props.installGuideProp.clusterId,
          product_info: autoSelectedProducts.map((item) => {
            return {
              id: item.ID,
              name: item.productName,
              service_list: item.service.checked,
              uncheck_service_list: item.service.unChecked,
            };
          }),
        };
        this.props.actions.getProductServicesInfoForAutoDeploy(params, () => {
          this.props.actions.nextStep();
        });
        return;
      }
    } else if (step === 2 || (upgradeStep === 0 && isUpgrade)) {
      const { userCenterProp } = this.props;
      const { authorityList } = userCenterProp;
      // 权限校验
      if (!alertModal(runtimeState, deployState)) return;
      if (!authorityList.package_upload_deploy) {
        return message.error('权限不足，请联系管理员！');
      }
      // 判断是否产品包升级
      if (isUpgrade) {
        // 产品包升级信息
        this.saveUpgrade();
      }
      if (deployMode === EnumDeployMode.AUTO) {
        // 调用自动部署的接口
        const params = {
          product_line_name:
            this.props.installGuideProp.selectProductLine?.product_line_name,
          product_line_version:
            this.props.installGuideProp.selectProductLine?.product_line_version,
          cluster_id: this.props.installGuideProp.clusterId,
          product_info: autoSelectedProducts.map((item) => {
            return {
              id: item.ID,
              name: item.productName,
              service_list: item.service.checked,
              uncheck_service_list: item.service.unChecked,
            };
          }),
        };
        this.props.actions.startAutoDeploy(params, (data) => {
          this.setState({
            deployUUID: data.deploy_uuid,
          });
          this.props.actions.nextStep();
        });
      } else {
        let info = JSON.parse(sessionStorage.getItem('product_backup_info'));
        let params: any = {
          productName: selectedProduct.product_name,
          version: selectedProduct.product_version,
          unchecked_services: unSelectedServiceList,
          deployMode: this.state.urlParams.new_version ? 1 : undefined,
          clusterId: this.props.installGuideProp.clusterId,
        };
        if (upgradeType === 'smooth') {
          let final_upgrade = false;
          if (smoothSelectService?.ServiceAddr) {
            if (smoothSelectService?.ServiceAddr?.UnSelect) {
              final_upgrade = false;
            } else {
              final_upgrade = true;
            }
          }
          Object.assign(params, {
            deployMode: 3,
            source_version: info?.source_version,
            final_upgrade: final_upgrade,
          });
        }
        if (installType === 'kubernetes') {
          params = Object.assign({}, params, {
            namespace: namespace,
            relynamespace: baseClusterId === -1 ? undefined : baseClusterId,
            clusterId: clusterId,
            pid: selectedProduct.id,
          });
        }
        this.props.actions.startDeploy(
          params,
          !isUpgrade
            ? this.props.actions.nextStep
            : () => {
                this.setState({ upgradeStep: 1 });
              }
        );
      }
      return;
    }
    this.props.actions.nextStep();
  };

  handleQuit = () => {
    // this.props.actions.quitGuide();
    const { deployState, runtimeState } = this.props.installGuideProp;
    const { upgradeStep, isUpgrade } = this.state;
    if (!alertModal(runtimeState, deployState)) return;
    if (
      (this.props.installGuideProp.step !== 3 && !isUpgrade) ||
      (isUpgrade && upgradeStep === 0)
    ) {
      Modal.confirm({
        title: '退出会失去已改动的内容，确认退出吗?',
        icon: <Icon type="exclamation-circle" theme="filled" />,
        okType: 'danger',
        onOk: () => {
          this.props.actions.initInstallGuide();
          this.props.history.push(this.getGoBackLocation());
          console.log('this.props:::', this.props);
        },
      });
    } else if (this.props.installGuideProp.complete === 'deployed') {
      this.props.actions.initInstallGuide();
      this.props.history.push(this.getGoBackLocation());
    } else {
      this.props.history.push(this.getGoBackLocation());
    }
  };

  // 跳转返回的路径
  getGoBackLocation = () => {
    const { gobackLocation } = this.state;
    if (!gobackLocation) {
      utils.setNaviKey('sub_menu_product_deploy', 'sub_menu_cluster_list');
      return '/deploycenter/cluster/detail/deployed';
    }
    if (gobackLocation.indexOf('/appmanage/') !== -1) {
      utils.setNaviKey('sub_menu_product_deploy', 'sub_menu_cluster_list');
      return '/deploycenter/cluster/list';
    }
    utils.setNaviKey('sub_menu_product_deploy', 'sub_menu_cluster_list');
    return '/deploycenter/cluster/detail/deployed';
  };

  handleUrlParamSearch = () => {
    // if (this.props.location.search === '') { return; }

    const urlParams: any = {};
    const search = this.props.location?.search.slice(1);
    if (search) {
      search.split('&').forEach((o: string) => {
        urlParams[o.split('=')[0]] = o.split('=')[1];
      });
    }
    this.setState(
      {
        urlParams,
        gobackLocation: urlParams.from || '/',
      },
      () => {
        this.initStep();
      }
    );
  };

  /**
   * 保存版本升级记录
   * @returns
   */
  saveUpgrade = async () => {
    const { selectedProduct = {}, oldHostInfo } = this.props.installGuideProp;
    const { upgradeType } = this.props.deployProp;
    let extra = {};
    if (upgradeType === 'smooth') {
      Object.assign(extra, { upgrade_mode: 'smooth' });
    }
    const backUpInfo =
      JSON.parse(sessionStorage.getItem('product_backup_info')) || {};
    const { data } = await installGuideService.saveUpgrade(
      {
        productName: selectedProduct.product_name,
      },
      {
        ...backUpInfo,
        ...oldHostInfo,
        ...extra,
      }
    );
    // this.setState({upgradeStep: 1})
    if (data.code === 0 && data.data?.upgrade_id) {
      return true;
    }
    return false;
  };

  handleStopDeploy = () => {
    const {
      selectedProduct = {},
      clusterId,
      namespace,
      installType,
    } = this.props.installGuideProp;
    const isKubernetes = installType === 'kubernetes';
    const { deployType } = this.props.deployProp;
    Modal.confirm({
      title: '确定停止当前部署吗？',
      content: '停止后，服务将不在继续部署！',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        if (deployType === 'auto') {
          // 停止自动部署
          const cluster_id = this.props.installGuideProp.clusterId;
          const deploy_uuid = this.state.deployUUID;
          this.props.actions.stopAutoDeploy(isKubernetes, {
            cluster_id,
            deploy_uuid: deploy_uuid || this.props.installGuideProp.deployUUID,
          });
        } else {
          const { upgradeType } = this.props.deployProp;
          let deploy_mode = 0;
          if (this.state.isUpgrade) {
            if (upgradeType === 'smooth') {
              deploy_mode = 3;
            } else {
              deploy_mode = 1;
            }
          }
          // 停止部署成功之后继续发请求直至后端响应结束
          const isKubernetes = installType === 'kubernetes';
          let params = {
            productName: selectedProduct.product_name,
            version: selectedProduct.product_version,
            product_type: selectedProduct.product_type,
            deploy_mode,
          };
          if (isKubernetes) {
            params = Object.assign(params, {
              clusterId,
              namespace,
            });
          }
          this.props.actions.stopDeploy(isKubernetes, params);
        }
      },
    });
  };

  deployDone = () => {
    if (this.state.gobackLocation != '/login') {
      const { actions, headerStore, installGuideProp } = this.props;
      const { parentClusters } = headerStore;
      const { clusterId } = installGuideProp;
      const currentCluster = parentClusters.find(
        (cluster) => cluster.id === +clusterId
      );
      if (currentCluster) {
        actions.setCurrentParentCluster(currentCluster);
      }
    }
    utils.setNaviKey('menu_deploy_center', 'sub_menu_product_deploy');
    this.props.history.push(
      this.state.gobackLocation === '/login'
        ? '/'
        : '/deploycenter/appmanage/installs'
    );
  };

  render() {
    const { canUpgrade, urlParams, isUpgrade, upgradeStep } = this.state;
    const { step, complete, stopDeployBySelf, baseClusterInfo, installType } =
      this.props.installGuideProp;
    const { baseClusterList, hasDepends } = baseClusterInfo || {};
    const currentStep = step;
    const isKubernetes = installType === 'kubernetes';
    // k8s集群下的选择产品包下如果存在依赖但依赖包没有部署的情况下，无法下一步
    const cannotK8sNext =
      isKubernetes &&
      currentStep === 1 &&
      hasDepends &&
      !baseClusterList.length;
    // 当前是否是升级
    // const isUpgrade = currentStep === 1 && !!urlParams.new_version;
    const P = {
      installGuideProp: this.props.installGuideProp,
      actions: this.props.actions,
      location: this.props.location,
      defaultSelectedProduct: urlParams,
      DeployProp: this.props.deployProp,
      isKubernetes,
    };
    // console.log(this.props.installGuideProp, 'currentStep')
    return (
      <div
        className="step-container"
        style={{ height: document.body.clientHeight - 50 }}>
        <div className="step-main-container">
          {!isUpgrade && (
            <div className="na-page-content-box">
              <Steps current={currentStep}>
                <Step title="选择集群" />
                <Step title="选择产品包" />
                <Step title="配置服务" />
                <Step title="执行部署" />
              </Steps>
              <div className="step-triangle-box">
                <div
                  className={
                    currentStep === 0
                      ? 'step-triangle-item'
                      : 'step-triangle-item not-active-step'
                  }></div>
                <div
                  className={
                    currentStep === 1
                      ? 'step-triangle-item'
                      : 'step-triangle-item not-active-step'
                  }></div>
                <div
                  className={
                    currentStep === 2
                      ? 'step-triangle-item'
                      : 'step-triangle-item not-active-step'
                  }></div>
                <div
                  className={
                    currentStep === 3
                      ? 'step-triangle-item'
                      : 'step-triangle-item not-active-step'
                  }></div>
              </div>
              {
                // 这里流程对调一下，把选择集群放在第一步
                (currentStep === 0 && (
                  <StepOne
                    wrappedComponentRef={(form) => (this.stepTwoForm = form)}
                    {...P}
                  />
                )) ||
                  (currentStep === 1 && (
                    <StepTwo
                      autoProductList={this.state.autoProductList}
                      deployMode={this.state.deployMode}
                      autoSelectedProducts={this.state.autoSelectedProducts}
                      updateAutoDeployService={this.updateAutoDeployService}
                      getOrchestrationHistory={this.getOrchestrationHistory}
                      updateParentState={this.setState.bind(this)}
                      isK8s={
                        this.props.installGuideProp.installType === 'kubernetes'
                      }
                      {...P}
                    />
                  )) ||
                  (currentStep === 2 && (
                    <StepThree
                      productName={urlParams.product_name}
                      deployMode={this.state.deployMode}
                      refreshDeployService={this.refreshDeployService}
                      {...P}
                    />
                  )) ||
                  (currentStep === 3 && <StepFour {...P} />)
              }
            </div>
          )}
          {isUpgrade && (
            <div className="na-page-content-box">
              <Steps current={upgradeStep} style={{ width: 300 }}>
                <Step title="配置服务" />
                <Step title="执行部署" />
              </Steps>
              <div className="step-triangle-box" style={{ width: '100%' }}>
                {(upgradeStep === 0 && (
                  <StepThree
                    productName={urlParams.product_name}
                    deployMode={this.state.deployMode}
                    refreshDeployService={this.refreshDeployService}
                    {...P}
                  />
                )) ||
                  (upgradeStep === 1 && <StepFour {...P} />)}
              </div>
            </div>
          )}
          <div className="btn-container">
            <Button className="dt-em-btn" onClick={this.handleQuit}>
              退出
            </Button>
            {currentStep === 3 || upgradeStep === 1 ? (
              <span>
                <Button
                  ghost
                  onClick={() => this.handleStopDeploy()}
                  disabled={
                    complete !== 'deploying' ||
                    this.props.deployProp.deployType === ''
                  }
                  className="dt-em-btn"
                  type="danger">
                  停止部署
                </Button>
                <Button
                  onClick={this.deployDone}
                  disabled={complete !== 'deployed' && !stopDeployBySelf}
                  className="dt-em-btn"
                  type="primary">
                  继续部署
                </Button>
              </span>
            ) : (
              <span>
                {currentStep !== 0 && !isUpgrade && (
                  <Button className="dt-em-btn" onClick={this.handleLastStep}>
                    上一步
                  </Button>
                )}
                <Button
                  className="dt-em-btn"
                  type="primary"
                  disabled={cannotK8sNext || !canUpgrade}
                  onClick={this.handleNextStep}>
                  {currentStep !== 2 && !isUpgrade ? '下一步' : '执行部署'}
                </Button>
              </span>
            )}
          </div>
        </div>
      </div>
    );
  }
}

export default StepIndex;
