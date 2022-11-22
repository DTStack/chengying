import * as React from 'react';
import {
  Button,
  Upload,
  Table,
  Icon,
  Radio,
  message,
  Modal,
  Checkbox,
  Select,
  Tooltip,
} from 'antd';
import { difference, isEqual } from 'lodash';
import { installGuideService, productLine } from '@/services';
import { deployStatusFilter } from '@/constants/const';
import Step2SetModal from './step2SetModal';
import UploadProductLine from './step2UploadProductLine';
import { EnumDeployMode } from './types';

const Option = Select.Option;

interface State {
  productPackageList: any[];
  fileList: any[];
  filters: any[];
  deploy_status: any;
  modalShow: any;
  productLine: any[];
  selectProductLine: any;
  showProductLine: boolean;
  selectProductVersion: string;
  deployProcess: any[];
  autoExpandedRowKeys: any[];
  manualExpandedRowKeys: any[];
}
const CheckboxGroup = Checkbox.Group;

const renderDeployStatus = (status: any, isShow?: boolean) => {
  let state: React.ReactNode = '';
  switch (status) {
    case 'deploying':
      if (isShow) {
        state = <span>{'部署中'}</span>;
      } else {
        state = <span>{'部署中'}</span>;
      }
      break;
    case 'deployed':
      if (isShow) {
        state = <span>{'已部署'}</span>;
      } else {
        state = <span>{'已部署'}</span>;
      }
      break;
    case 'deploy fail':
      if (isShow) {
        state = <span>{'部署失败'}</span>;
      } else {
        state = <span>{'部署失败'}</span>;
      }
      break;
    case 'undeploying':
      if (isShow) {
        state = <span>{'卸载中'}</span>;
      } else {
        state = <span>{'卸载中'}</span>;
      }
      break;
    case 'undeploy fail':
      if (isShow) {
        state = <span>{'卸载失败'}</span>;
      } else {
        state = <span>{'卸载失败'}</span>;
      }
      break;
    case 'undeployed':
    default:
      if (isShow) {
        state = '';
      } else {
        state = <span>{'未部署'}</span>;
      }
      break;
  }
  return <span>{state}</span>;
};
class StepTwo extends React.Component<any, State> {
  setmodal: any;
  constructor(props: any) {
    super(props);
  }

  componentDidMount() {
    const { isK8s, deployMode } = this.props;
    this.props.actions.getProductStepOneList(
      {
        product_line_name: '',
        product_line_version: '',
        product_type: !this.props.isKubernetes ? 0 : 1,
        deploy_status: '',
      },
      false
    );
    if (deployMode === EnumDeployMode.AUTO) {
      this.getProductLine();
    } else {
      this.initProductPackageList();
    }
    this.props.getOrchestrationHistory();
    installGuideService
      .getProductNames({ product_type: !this.props.isKubernetes ? 0 : 1 })
      .then((res: any) => {
        res = res.data;
        if (res.code !== 0) {
          message.error(res.msg);
          return;
        }
        res.data.forEach((o: string) => {
          this.state.filters.push({
            text: o,
            value: o,
          });
        });
        if (isK8s) {
          this.props.updateParentState({
            deployMode: EnumDeployMode.MANUAL,
          });
        }
        this.setState({
          filters: this.state.filters,
        });
      });
  }

  componentDidUpdate(prevProps: any) {
    const prevProductPackageList =
      prevProps.installGuideProp.productPackageList;
    const { productPackageList } = this.props.installGuideProp;
    if (!isEqual(prevProductPackageList, productPackageList)) {
      this.setState({ productPackageList });
    }
    if (
      JSON.stringify(this.props.installGuideProp.selectedProduct) === '{}' &&
      JSON.stringify(this.props.defaultSelectedProduct) !== '{}' &&
      this.props.deployMode === EnumDeployMode.MANUAL
    ) {
      this.initSelectedProduct();
    }
    if (prevProps.deployMode != this.props.deployMode) {
      this.setState({ selectProductLine: {}, deployProcess: [] });
      this.props.actions.setProductLine({});
      this.props.actions.getProductStepOneList(
        {
          product_line_name: '',
          product_line_version: '',
          product_type: !this.props.isKubernetes ? 0 : 1,
          deploy_status: '',
        },
        false,
        () => {
          if (this.props.deployMode === EnumDeployMode.AUTO) {
            this.setState({
              productPackageList:
                this.props.installGuideProp.productPackageList,
              autoExpandedRowKeys: [],
            });
          }
        }
      );
      if (this.props.deployMode === EnumDeployMode.AUTO) {
        this.getProductLine();
      } else {
        this.initProductPackageList();
        if (JSON.stringify(this.props.defaultSelectedProduct) !== '{}') {
          this.initSelectedProduct();
        }
      }
    }
  }

  state: State = {
    productPackageList: this.props.installGuideProp.productPackageList,
    fileList: [],
    filters: [],
    deploy_status: '',
    modalShow: false,
    showProductLine: false,
    productLine: [],
    selectProductLine: {},
    selectProductVersion: '',
    deployProcess: [],
    autoExpandedRowKeys: [],
    manualExpandedRowKeys: [],
  };

  uploadOnChange = async (info: any) => {
    const { selectProductLine } = this.state;
    if (info.file.response) {
      if (info.file.response.code !== 0) {
        message.error(info.file.response.msg);
      } else {
        message.success(`${info.file.name} 上传成功！`);
        // 刷新列表
        this.props.actions.getProductStepOneList(
          {
            product_line_name: selectProductLine.product_line_name,
            product_line_version: selectProductLine.product_line_version,
            product_type: !this.props.isKubernetes ? 0 : 1,
            deploy_status: '',
          },
          false
        );
        let res = await productLine.getProductLine();
        if (res.data.code == 0) {
          this.setState({
            deployProcess: [],
          });
          let arr = res.data.data.list.filter(
            (item) => item.id == selectProductLine.id
          );
          let isExitArr = arr[0]?.deploy_process.filter(
            (item) => item.is_exist === false
          );
          this.setState({
            productLine: res.data.data.list,
            deployProcess: arr[0]?.deploy_process,
            selectProductLine: {
              ...arr[0],
              isShowTip: isExitArr?.length > 0 ? true : false,
            },
          });
        }
        if (info.file.response.data) {
          message.warning(info.file.response.data);
        }
      }
    }
  };

  closeOpenStatus = (record) => {
    const { deployMode } = this.props;
    const { productPackageList } = this.state;
    const index = productPackageList.findIndex(
      (item) => record.product_name === item.product_name
    );
    let arr = productPackageList;
    arr[index].isOpen = false;
    this.setState({ productPackageList: arr });
    if (deployMode === EnumDeployMode.AUTO) {
      const openKeys = arr.filter((item) => item.isOpen === true);
      this.setState({
        autoExpandedRowKeys: [...openKeys].map((item) => item.product_name),
      });
    }
  };

  setNextSelectedProduct = (data, productSelect, record) => {
    const { autoSelectedProducts } = this.props;
    const checkedList = data.filter((service) => service.baseProduct === '');
    const unCheckedList = data.filter((service) => service.baseProduct !== '');
    const disabled = data
      .filter((product) => product.baseProduct !== '')
      .map((item) => item.serviceName);
    // 此处逻辑，当product版本冲突时，需将第一个版本移除
    const nextAutoSelectedProduct = autoSelectedProducts.filter(
      (product) => product.productName !== productSelect.product_name
    );
    nextAutoSelectedProduct.push({
      ID: productSelect.id,
      productName: record.product_name,
      productVersion: productSelect.product_version,
      service: {
        all: [...checkedList, ...unCheckedList],
        checked: this.props.installGuideProp.selectedServiceList,
        unChecked: this.props.installGuideProp.unSelectedServiceList,
        disabled,
      },
    });
    this.props.updateParentState({
      autoSelectedProducts: nextAutoSelectedProduct,
    });
  };

  autoProductOpen = (record, productSelect) => {
    this.getProductServicePromise({
      productRecord: productSelect,
      isKubernetes: this.props.isKubernetes,
    }).then((data) => {
      const checkedList = data.filter((service) => service.baseProduct === '');
      const checkedService = [...checkedList].map((item) => item.serviceName);
      this.props.actions.getUncheckedServices(
        {
          pid: productSelect.id,
          clusterId: this.props.installGuideProp.clusterId,
          namespace: this.props.installGuideProp.namespace,
        },
        checkedService,
        () => {
          this.setNextSelectedProduct(data, productSelect, record);
        }
      );
    });
  };
  handleProductSelected = (record: any, isOpen, selectProduct?: any) => {
    const { deployMode } = this.props;
    const { productPackageList } = this.state;
    const index = productPackageList.findIndex(
      (item) => record.product_name === item.product_name
    );
    let arr = productPackageList;
    if (deployMode === EnumDeployMode.AUTO) {
      arr[index].isOpen = isOpen;
    }
    this.setState({ productPackageList: arr });
    const productSelect =
      selectProduct ?? record.list.filter((item) => item.is_default);
    this.setState({ selectProductVersion: productSelect[0].product_version });
    if (deployMode === EnumDeployMode.AUTO) {
      const openKeys = arr.filter((item) => item.isOpen === true);
      this.setState({
        autoExpandedRowKeys: [...openKeys].map((item) => item.product_name),
      });
      if (!selectProduct) {
        this.autoProductOpen(record, productSelect[0]);
      }
      return;
    }
    if (this.props.isKubernetes) {
      // 检测是否有依赖关系
      this.props.actions.getBaseClusterList(
        {
          clusterId: this.props.installGuideProp.clusterId,
          namespace: this.props.installGuideProp.namespace,
          pid: productSelect[0].id,
        },
        (state) => {
          if (!state.baseClusterInfo.hasDepends || state.baseClusterId !== -1) {
            this.props.actions.resetInstallGuideConfig();
            if (!state.baseClusterId) return;
            this.getProductPackageServices(
              productSelect[0],
              state.baseClusterId
            );
          }
        }
      );
    } else {
      this.getProductPackageServices(productSelect[0]);
    }
    // 如果是其他地方进来，那就执行获取未选择，更新未选择值，默认为空。  然后选中的值默认全部，需要再去派发一次，全部和未选择筛选
    // this.props.actions.getUncheckedService(record.ID)
    // 保存选中的产品包信息
    this.setState({ manualExpandedRowKeys: [productSelect[0].product_name] });
    this.props.actions.saveInstallInfo(productSelect[0]);
  };

  // 获取组件
  getProductPackageServices = (
    selectedProduct: any,
    relynamespace?: string | -1
  ) => {
    this.props.actions.getProductPackageServices({
      productName: selectedProduct.product_name,
      productVersion: selectedProduct.product_version,
      pid: selectedProduct.id,
      clusterId: this.props.installGuideProp.clusterId,
      relynamespace: relynamespace === -1 ? undefined : relynamespace,
      namespace: this.props.installGuideProp.namespace,
    });
  };

  handleDeleteProductPackage = (e: any) => {
    Modal.confirm({
      title: '确定删除此安装包？',
      content: '删除后，再次部署时需要再次上传。',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        installGuideService
          .deleteProducyPackage({
            productName: e.product_name,
            productVersion: e.product_version,
          })
          .then((res: any) => {
            if (res.data.code === 0) {
              this.setState({});
              this.props.actions.getProductPackageList({
                limit: 0,
                product_type: !this.props.isKubernetes ? 0 : 1,
              });
            } else {
              message.error(res.data.msg);
            }
          });
      },
    });
  };

  initSelectedProduct = () => {
    const params = this.props.defaultSelectedProduct;
    this.state.productPackageList.forEach((o: any) => {
      if (params.product_name === o.product_name) {
        const selectProduct = o.list.filter(
          (item) => params.product_version == item.product_version
        );
        this.changeProductVersion(selectProduct[0], o);
      }
    });
  };

  initProductPackageList = () => {
    const { productPackageList } = this.state;
    let arr = productPackageList;
    arr.map((item) => {
      item.list.map((citem) => {
        citem.is_default = false;
        return citem;
      });
      return item;
    });
    this.setState({ productPackageList: arr });
  };

  changeCheckbox = (checkedValue) => {
    const { productServices } = this.props.installGuideProp;
    const allProductService = productServices.map(
      (item: any) => item.serviceName
    );
    const unSelect = difference(allProductService, checkedValue);
    this.props.actions.saveSelectedService(checkedValue);
    this.props.actions.saveUnSelectedService(unSelect);
  };

  getProductServicePromise = ({ isKubernetes, productRecord }) => {
    const params = {
      productName: productRecord.product_name,
      productVersion: productRecord.product_version,
      pid: productRecord.id,
      clusterId: this.props.installGuideProp.clusterId,
      baseClusterId: undefined,
    };

    if (isKubernetes) {
      return installGuideService
        .getBaseClusterList({
          clusterId: this.props.installGuideProp.clusterId,
          namespace: this.props.installGuideProp.namespace,
          pid: productRecord.id,
        })
        .then((res) => {
          res = res.data;
          if (res.code === 0) {
            const { candidates, targets, message, hasDepends } = res.data;
            const cluster = candidates.find(
              (item: any) => item.relynamespace === targets.relynamespace
            );
            const payload = {
              baseClusterInfo: {
                baseClusterList: candidates,
                hasDepends: hasDepends,
                dependMessage: message,
              },
              baseClusterId: cluster
                ? targets.relynamespace === 0
                  ? -1
                  : targets.relynamespace
                : -1,
            };
            return payload;
          } else {
            throw new Error(res.msg);
          }
        })
        .then(({ baseClusterId, baseClusterInfo }) => {
          if (!baseClusterInfo.hasDepends || baseClusterId !== -1) {
            const id = baseClusterId === -1 ? undefined : baseClusterId;
            params.baseClusterId = id;
            return installGuideService.getProductPackageServices(params);
          }
        })
        .then((res) => {
          const { data } = res.data;
          return data;
        });
    } else {
      return installGuideService
        .getProductPackageServices(params)
        .then((res) => {
          const { data } = res.data;
          return data;
        })
        .catch(() => {
          message.error('获取组件失败');
        });
    }
  };

  changeProductVersion = (item, record) => {
    if (this.props.deployMode === EnumDeployMode.MANUAL) {
      this.initProductPackageList();
    }
    const { productPackageList } = this.state;
    let arr = productPackageList;
    let selectProductPackageList: any = {};
    const index = productPackageList.findIndex(
      (items) => items.product_name === item.product_name
    );
    selectProductPackageList = productPackageList[index];
    selectProductPackageList.list.map((citem) => {
      if (citem.product_version == item.product_version) {
        citem.is_default = true;
      } else {
        citem.is_default = false;
      }
      return citem;
    });
    arr[index] = selectProductPackageList;
    this.setState(
      {
        productPackageList: arr,
        selectProductVersion: item.product_version,
      },
      () => {
        this.handleProductSelected(record, false, [item]);
      }
    );
  };

  tableColInit = () => {
    return [
      {
        title: '组件名称',
        dataIndex: 'product_name',
        ellipsis: true,
        render: (e: any, record: any) => {
          return (
            <div style={{ display: 'inline-block' }}>
              {this.props.deployMode === EnumDeployMode.AUTO && (
                <>
                  {record?.isOpen ? (
                    <Icon
                      type="down"
                      style={{
                        color: '#999999',
                        marginRight: 20,
                        cursor: 'pointer',
                      }}
                      onClick={() => this.closeOpenStatus(record)}
                    />
                  ) : (
                    <Icon
                      type="right"
                      style={{
                        color: '#999999',
                        marginRight: 20,
                        cursor: 'pointer',
                      }}
                      onClick={() => this.handleProductSelected(record, true)}
                    />
                  )}
                </>
              )}
              <span>{e}</span>
            </div>
          );
        },
      },
      {
        title: '版本号',
        dataIndex: 'list',
        render: (text, record) => {
          return (
            <div className="versionList">
              {text.map((item: any) => (
                <div
                  onClick={() => this.changeProductVersion(item, record)}
                  className={item.is_default ? 'activeItem' : 'versionItem'}>
                  {item?.product_version}{' '}
                  {renderDeployStatus(item?.status, true)}
                </div>
              ))}
            </div>
          );
        },
      },
      {
        title: '部署状态',
        dataIndex: 'status',
        filters: [
          ...deployStatusFilter,
          {
            text: '未部署',
            value: 'undeployed',
          },
        ],
        render: (text) => renderDeployStatus(text, false),
      },
      {
        title: '操作',
        dataIndex: 'action',
        render: (e: any, record: any) => {
          return (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
              }}>
              <a
                style={{
                  marginRight: '10px',
                }}
                onClick={(e) => {
                  record.list.map((item) => {
                    item.field = '';
                    return item;
                  });
                  e.preventDefault();
                  this.setState({
                    modalShow: record,
                  });
                }}>
                <Icon type="setting" />
              </a>
            </div>
          );
        },
      },
    ];
  };

  handleTableChange = (pagi: any, filter: any, sorter: any) => {
    const { selectProductLine } = this.state;
    this.props.actions.getProductStepOneList({
      product_line_name: selectProductLine?.product_line_name,
      product_line_version: selectProductLine?.product_line_version,
      product_type: !this.props.isKubernetes ? 0 : 1,
      deploy_status: filter.status && filter.status.join(',').toString(),
    });
  };

  serviceUpdate = (params: any) => {
    let propsParams = [];
    params.forEach((item: any) => {
      propsParams.push({
        product_version: item.product_version,
        field: item.field,
        field_path: 'Instance.RunUser',
      });
    });
    installGuideService
      .serviceUpdate(
        { ProductName: this.setmodal.props?.data.product_name },
        propsParams
      )
      .then((res: any) => {
        if (res.data.code === 0) {
          message.success('设置成功！');
          this.setState({
            modalShow: false,
          });
        } else {
          message.error(res.data.msg);
        }
      });
  };

  // 选择依赖
  handleBaseClusterChange = (relynamespace: string) => {
    const { selectedProduct } = this.props.installGuideProp;
    this.props.actions.saveSelectBaseCluster(relynamespace);
    this.getProductPackageServices(selectedProduct, relynamespace);
  };

  handleDeployModeChange = (e) => {
    this.props.updateParentState({
      deployMode: e.target.value,
    });
    if (e.target.value === EnumDeployMode.AUTO) {
      this.props.getOrchestrationHistory();
    }
  };

  // 产品线变化
  changeProductLine = (value) => {
    if (!value) {
      this.setState({
        deployProcess: [],
        selectProductLine: {},
      });
      this.props.actions.getProductStepOneList({
        product_line_name: '',
        product_line_version: '',
        product_type: !this.props.isKubernetes ? 0 : 1,
        deploy_status: '',
      });
      this.props.actions.setProductLine({});
      return;
    }
    const { productLine } = this.state;
    let arr = productLine.filter((item) => item.id == value);
    let isExitArr = arr[0]?.deploy_process.filter(
      (item) => item.is_exist === false
    );
    this.setState({
      deployProcess: arr[0]?.deploy_process,
      selectProductLine: {
        ...arr[0],
        isShowTip: isExitArr?.length > 0 ? true : false,
      },
    });
    this.props.actions.setProductLine({
      product_line_name: arr[0].product_line_name,
      product_line_version: arr[0].product_line_version,
    });
    this.props.actions.getProductStepOneList(
      {
        product_line_name: arr[0].product_line_name,
        product_line_version: arr[0].product_line_version,
        product_type: !this.props.isKubernetes ? 0 : 1,
        deploy_status: '',
      },
      false,
      () => {
        this.props.updateParentState({
          autoSelectedProducts: [],
        });
        this.setState({ autoExpandedRowKeys: [] });
        this.props.installGuideProp.productPackageList.map((item) => {
          const productSelect = item.list.filter((items) => items.is_default);
          this.getProductServicePromise({
            productRecord: productSelect[0],
            isKubernetes: this.props.isKubernetes,
          }).then((data) => {
            const checkedList = data.filter(
              (service) => service.baseProduct === ''
            );
            const checkedService = [...checkedList].map(
              (item) => item.serviceName
            );
            this.props.actions.getUncheckedServices(
              {
                pid: productSelect[0].id,
                clusterId: this.props.installGuideProp.clusterId,
                namespace: this.props.installGuideProp.namespace,
              },
              checkedService,
              () => {
                this.setNextSelectedProduct(data, productSelect[0], item);
              }
            );
          });
        });
      }
    );
  };
  // 打开产品线上传
  showUploadProductModal = () => {
    this.setState({ showProductLine: true });
  };

  // 关闭产品线上传弹框
  cancelProductLine = () => {
    this.setState({ showProductLine: false });
  };

  // 产品线列表回调
  callList = (selectId?: string) => {
    if (selectId === this.state.selectProductLine.id) {
      this.setState({ selectProductLine: {}, deployProcess: [] });
      this.props.actions.getProductStepOneList(
        {
          product_line_name: '',
          product_line_version: '',
          product_type: !this.props.isKubernetes ? 0 : 1,
          deploy_status: '',
        },
        false
      );
    }
    this.getProductLine();
  };

  // 获取产品线列表
  getProductLine = async () => {
    let res = await productLine.getProductLine();
    if (res.data.code == 0) {
      this.setState({ productLine: res.data.data.list });
    } else {
      message.error(res.data.msg);
    }
  };

  render() {
    const uploadCfg = {
      name: 'package',
      action: '/api/v2/product/upload',
      headers: {
        authorization: 'authorization-text',
      },
      onChange: this.uploadOnChange,
      accept: '.tar',
    };

    const { productServices, baseClusterId, baseClusterInfo } =
      this.props.installGuideProp;
    const { baseClusterList, hasDepends, dependMessage } = baseClusterInfo;
    const { deployMode } = this.props;
    const radioStyle = {
      display: 'block',
      marginBottom: '20px',
    };

    const {
      productLine,
      selectProductLine,
      showProductLine,
      deployProcess,
      autoExpandedRowKeys,
      manualExpandedRowKeys,
    } = this.state;

    return (
      <div
        className="step-one-container step-content-container"
        style={{ overflow: 'scroll' }}>
        <div>
          <Radio.Group
            value={deployMode}
            defaultValue={EnumDeployMode.AUTO}
            onChange={this.handleDeployModeChange}>
            {/* TODO: k8s不支持自动部署，部分代码逻辑需要删除 */}
            <Radio
              disabled={this.props.isK8s}
              style={radioStyle}
              value={EnumDeployMode.AUTO}>
              <span className="radio-main">自动部署</span>
              <span className="radio-tips">
                （基于产品线顺序部署多个组件包，支持主机角色自动编排）
              </span>
            </Radio>
            <Radio style={radioStyle} value={EnumDeployMode.MANUAL}>
              <span className="radio-main">手动部署</span>
              <span className="radio-tips">
                （基于单个组件包部署，支持自定义主机编排）
              </span>
            </Radio>
          </Radio.Group>
        </div>
        <div className="header-box">
          <Upload {...uploadCfg}>
            <Button type="primary">
              <Icon type="upload" />
              上传组件包
            </Button>
          </Upload>
          <span style={{ marginLeft: 12 }}>
            请在下方选择需要安装的组件包，可上传新的组件包，也可重新部署之前的组件包。
          </span>
        </div>
        {deployMode === EnumDeployMode.MANUAL && (
          <div>
            <Table
              scroll={{ x: '130%' }}
              pagination={false}
              className="dt-em-table dt-table-border dt-table-last-row-noborder stepTwoTable"
              rowKey="product_name"
              onChange={this.handleTableChange}
              dataSource={this.state.productPackageList}
              columns={this.tableColInit()}
              expandedRowRender={() => (
                <div style={{ overflow: 'auto' }}>
                  {hasDepends && (
                    <div
                      className="mt-10 mb-10"
                      style={{ margin: '12px 0 0 0', lineHeight: '32px' }}>
                      <Icon
                        type="exclamation-circle"
                        theme="filled"
                        className="mr-8"
                        style={{ color: '#f5a841', lineHeight: '20px' }}
                      />
                      <span className="mr-8" style={{ lineHeight: '20px' }}>
                        {dependMessage}
                      </span>
                      {baseClusterList.length > 0 && (
                        <Select
                          style={{ width: 264, height: '32px' }}
                          placeholder="请选择依赖集群"
                          value={
                            baseClusterId === -1 ? undefined : baseClusterId
                          }
                          onChange={this.handleBaseClusterChange}>
                          {/* TODO: 当前版本此处baseClusterId的意义已经更改为relynamespace,后续需要变更变量名称 */}
                          {baseClusterList.map((item) => (
                            <Option
                              key={item.relynamespace}
                              value={item.relynamespace}>
                              {item.relynamespace}
                            </Option>
                          ))}
                        </Select>
                      )}
                    </div>
                  )}
                  {!(!hasDepends || baseClusterId !== -1) && (
                    <p style={{ height: '12px' }}></p>
                  )}
                  {(!hasDepends || baseClusterId !== -1) && (
                    <React.Fragment>
                      <div
                        style={{ margin: '12px 0 0 36px', lineHeight: '20px' }}>
                        ChengYing安装程序将安装以下服务：
                      </div>
                      <CheckboxGroup
                        style={{
                          float: 'left',
                          marginBottom: '20px',
                          marginLeft: '36px',
                        }}
                        onChange={(checkedValue) => {
                          this.changeCheckbox(checkedValue);
                        }}
                        value={this.props.installGuideProp.selectedServiceList} // 传值类型及字段待后端确定
                      >
                        {productServices.map((o: any, index: number) => (
                          <p
                            key={`${o.serviceName}-${index}`}
                            style={{ marginTop: '12px', lineHeight: '20px' }}>
                            <Checkbox
                              value={`${o.serviceName}`}
                              disabled={!!o.baseProduct}>
                              <span style={{ marginRight: 10 }}>
                                {o.serviceName}
                              </span>
                              <span>
                                {o.serviceVersion ? o.serviceVersion : ''}
                              </span>
                              <span style={{ marginLeft: 10 }}>
                                {o.baseProduct}
                              </span>
                            </Checkbox>
                          </p>
                        ))}
                      </CheckboxGroup>
                    </React.Fragment>
                  )}
                </div>
              )}
              expandedRowKeys={manualExpandedRowKeys}
            />
          </div>
        )}

        {deployMode === EnumDeployMode.AUTO && (
          <div>
            <UploadProductLine
              callList={this.callList}
              visible={showProductLine}
              onCancel={this.cancelProductLine}
              dataList={productLine}
            />
            <div className="productLine">
              <div className="productNav">
                <div style={{ display: 'flex' }}>
                  <div className="navTopLine"></div>
                  <div className="navTopTxt">产品线</div>
                </div>
                <div className="navButton">
                  <Button
                    icon="upload"
                    type="primary"
                    size="small"
                    onClick={this.showUploadProductModal}>
                    上传产品线
                  </Button>
                </div>
              </div>
              <div className="productContent">
                <div className="productContentTop">
                  <div className="productContentName">产品线名称</div>
                  <div className="productContentProcess">部署流程</div>
                </div>
                <div className="productContentBottom">
                  {selectProductLine.isShowTip && (
                    <div className="productIcon">
                      <Tooltip
                        placement="top"
                        title="缺失组件包，请先“上传组件包”">
                        <Icon
                          type="exclamation-circle"
                          style={{ color: '#faad14' }}
                        />
                      </Tooltip>
                    </div>
                  )}
                  <div>
                    <Select
                      value={selectProductLine?.id}
                      allowClear
                      showSearch
                      style={{ width: 200 }}
                      placeholder="请选择一条产品线"
                      optionFilterProp="children"
                      onChange={this.changeProductLine}>
                      {productLine.map((item: any) => (
                        <Option value={item.id}>
                          {item.product_line_name} ({item.product_line_version})
                        </Option>
                      ))}
                    </Select>
                  </div>
                  {deployProcess.length > 0 && (
                    <div className="productContentList">
                      {deployProcess?.map((sitem: any, index: number) => (
                        <div style={{ display: 'flex' }}>
                          <div
                            className={
                              sitem.is_exist ? 'listItem' : 'disListItem'
                            }>
                            {sitem.product_name}
                          </div>
                          {index !== deployProcess.length - 1 && (
                            <>
                              <div className="listLine"></div>
                              <div className="listArrow"></div>
                            </>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
            <div className="productNav" style={{ margin: '20px 0 20px 0' }}>
              <div style={{ display: 'flex' }}>
                <div className="navTopLine"></div>
                <div className="navTopTxt">
                  组件{' '}
                  <span className="noticeTxt">
                    （选中产品线后，组件列表仅展示其包含的组件）
                  </span>
                </div>
              </div>
            </div>
            <Table
              scroll={{ x: '130%' }}
              pagination={false}
              className="dt-em-table dt-table-border dt-table-last-row-noborder stepTwoTable"
              rowKey="product_name"
              onChange={this.handleTableChange}
              dataSource={this.state.productPackageList}
              columns={this.tableColInit()}
              expandedRowRender={(record) => {
                const target = this.props.autoSelectedProducts.find(
                  (item) => item.productName === record.product_name
                );
                if (!target) return null;
                return (
                  <div>
                    <React.Fragment>
                      <div
                        style={{ margin: '12px 0 0 36px', lineHeight: '20px' }}>
                        ChengYing安装程序将安装以下服务：
                      </div>
                      <CheckboxGroup
                        style={{
                          float: 'left',
                          marginBottom: '20px',
                          marginLeft: '36px',
                        }}
                        onChange={(checkedValue) => {
                          this.props.updateAutoDeployService(
                            target.ID,
                            checkedValue
                          );
                        }}
                        value={[
                          ...target.service.checked,
                          ...target.service.disabled,
                        ]}>
                        {target.service.all.map(
                          (service: any, index: number) => (
                            <p
                              key={`${service.serviceName}-${index}`}
                              style={{ marginTop: '12px', lineHeight: '20px' }}>
                              <Checkbox
                                value={`${service.serviceName}`}
                                disabled={!!service.baseProduct}
                                checked={true}>
                                <span style={{ marginRight: 10 }}>
                                  {service.serviceName}
                                </span>
                                <span>
                                  {service.serviceVersion
                                    ? service.serviceVersion
                                    : ''}
                                </span>
                                <span style={{ marginLeft: 10 }}>
                                  {service.baseProduct}
                                </span>
                              </Checkbox>
                            </p>
                          )
                        )}
                      </CheckboxGroup>
                    </React.Fragment>
                  </div>
                );
              }}
              expandedRowKeys={autoExpandedRowKeys}
            />
          </div>
        )}
        <Modal
          visible={!!this.state.modalShow}
          title="设置"
          okText="保存"
          width={600}
          onOk={() => {
            this.serviceUpdate(this.setmodal.state.testData);
          }}
          onCancel={() => {
            this.setState({
              modalShow: false,
            });
          }}>
          <Step2SetModal
            data={this.state.modalShow}
            ref={(dom) => {
              this.setmodal = dom;
            }}
          />
        </Modal>
      </div>
    );
  }
}
export default StepTwo;
