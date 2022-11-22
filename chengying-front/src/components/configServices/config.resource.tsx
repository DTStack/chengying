import * as React from 'react';
import { connect } from 'react-redux';
import { Switch, Input, Row, message, Button, Modal, Icon, Tag } from 'antd';
import CustomTransfer from '@/components/customTransfer';
import installGuideService from '@/services/installGuideService';
import { AppStoreTypes } from '@/stores';
import { isEqual } from 'lodash';
declare const window: any;

interface Prop {
  hasInstance: boolean;
  maxSelected: number;
  handleCancel: () => void;
  handleResourceSubmit: Function;
  isCloud: boolean;
  existIp: any[];
  serviceKey: string;
  selectedKeys: any[];
  targetKeys: any[];
  syncHost: Function;
  hostList: any[];
  selectedProduct: any;
  Instance: any;
  installGuideProp: any;
  actions: any;
  isKubernetes?: boolean;
  is_Rerro?: boolean;
  is_Lerro?: boolean;
  DeployProp?: any;
  smoothTarget?: any[];
  initData?: any;
}
interface State {
  targetKeys: any[];
  dataList: any[];
  selectedKeys: any[];
  isCloud_state: boolean;
  cloudHost: string;
  showBtn: boolean;
  isTransferChange: boolean;
  showAutoBtn: boolean;
  isRerro: boolean;
  isLerro: boolean;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  initData: state.InstallGuideStore,
});
@(connect(mapStateToProps) as any)
class Resource extends React.Component<Prop, State> {
  constructor(props: Prop) {
    super(props);
  }

  state: State = {
    targetKeys: [],
    dataList: [],
    selectedKeys: [],
    isCloud_state: this.props.isCloud,
    cloudHost: '',
    showBtn: false,
    isTransferChange: false,
    showAutoBtn: false,
    isRerro: this.props.is_Rerro,
    isLerro: this.props.is_Lerro,
  };

  componentDidMount() {
    // if(this.state.selectedKeys.length === 0 && this.state.targetKeys.length === 0)
    const { upgradeType } = this.props.DeployProp;

    if (this.props.hasInstance) {
      this.setState(
        {
          targetKeys:
            upgradeType == 'smooth'
              ? this.props?.smoothTarget
              : this.props.targetKeys,
          selectedKeys: this.props.selectedKeys,
        },
        () => {
          this.initDataRight();
        }
      );
    } else {
      this.setState({
        cloudHost: this.props.existIp.toString() || '',
      });
    }
  }

  componentDidUpdate(prevProps: Prop) {
    const { upgradeType } = this.props.DeployProp;
    if (!isEqual(this.props, prevProps)) {
      // 没有instance的时候就没有transferchange
      if (this.state.isTransferChange && this.props.hasInstance) {
        this.setState({
          isTransferChange: false,
        });
        return;
      }
      this.computeRealTargetKeys(this.props);
      if (!this.props.hasInstance) {
        this.setState({
          cloudHost: this.props.existIp.toString() || '',
        });
        return;
      }
      this.setState(
        {
          dataList: [],
          selectedKeys: [],
        },
        () => {
          // this.initDataRight();
          this.setState({
            isCloud_state: this.props.isCloud,
            cloudHost: this.props.isCloud ? this.props.existIp.toString() : '',
          });
        }
      );
      if (this.props?.smoothTarget !== prevProps?.smoothTarget) {
        this.setState({
          targetKeys:
            upgradeType == 'smooth'
              ? this.props?.smoothTarget
              : this.props.targetKeys,
        });
      }
      if (this.props.is_Lerro !== prevProps.is_Lerro) {
        this.setState({ isLerro: this.props.is_Lerro });
      }
      if (this.props.is_Rerro !== prevProps.is_Rerro) {
        this.setState({ isRerro: this.props.is_Rerro });
      }
      this.initData_left(this.props.hostList);
    }
  }

  computeRealTargetKeys = (nextProp: Prop) => {
    const { upgradeType } = this.props.DeployProp;
    const { selectedService } = this.props.initData;
    const f: any[] = [];
    let existIp: any[] = [];
    if (upgradeType === 'smooth') {
      existIp = selectedService.ServiceAddr.Select
        ? [...selectedService.ServiceAddr.Select].map((item) => item.IP)
        : [];
    } else {
      existIp = nextProp.existIp;
    }
    // tslint:disable-next-line
    nextProp.isCloud === this.state.isCloud_state &&
      existIp.forEach((o: any) => {
        this.props.hostList.forEach((q: any) => {
          if (q.ip === o) {
            f.push(q.id);
          }
        });
      });

    this.setState(
      {
        targetKeys: upgradeType == 'smooth' ? this.props?.smoothTarget : f,
      },
      () => {
        this.initDataRight(f);
      }
    );
  };

  // 处理数据，分页左面板
  initData_left = (list: any) => {
    this.setState({
      selectedKeys: [],
      dataList: list,
    });
  };

  // 处理数据，分页右面板数据
  initDataRight = (targetKeys?: string[]) => {
    this.setState({
      selectedKeys: [],
      targetKeys,
    });
  };

  handleChange = (targetKeys: any, direction: any, moveKeys: any) => {
    const { selectedService } = this.props.installGuideProp;
    const { forcedUpgrade, isFirstSmooth } = this.props.DeployProp;
    // 是否属于可平滑升级的组件
    if (forcedUpgrade.includes(selectedService.serviceKey) && isFirstSmooth) {
      if (this.props.hostList?.length === targetKeys.length) {
        this.setState({ isRerro: false, isLerro: true });
      } else {
        if (targetKeys.length === 0) {
          this.setState({ isLerro: false, isRerro: true });
        } else {
          this.setState({ isLerro: false, isRerro: false });
        }
      }
    }
    if (targetKeys.length > this.props.maxSelected) {
      message.error(`IP数量限制${this.props.maxSelected},目前超出限制！`);
      const flagP = targetKeys.filter((tar) => !moveKeys.includes(tar));
      this.setState({
        targetKeys: flagP,
      });
      return;
    }
    // 保存至后端
    this.props.handleResourceSubmit({
      isCloud: this.state.isCloud_state,
      hosts:
        this.state.isCloud_state || !this.props.hasInstance
          ? this.state.cloudHost
          : targetKeys,
    });
    this.setState(
      {
        showBtn: true,
        isTransferChange: false,
      },
      () => {
        this.initDataRight(targetKeys);
      }
    );
  };

  handleSelect = (e: any, p: any) => {
    this.setState({
      selectedKeys: [...e, ...p],
    });
  };

  // 后端需要在切换调用(将true, false传入)具体逻辑他也不知道。
  handleChangeSwitch = () => {
    this.props.handleResourceSubmit({
      isCloud: this.state.isCloud_state,
      // hosts: [],
      hosts:
        this.state.isCloud_state || !this.props.hasInstance
          ? this.state.cloudHost
          : this.state.targetKeys,
    });
  };

  useCloudChange = (e: boolean) => {
    this.setState(
      {
        isCloud_state: e,
        showBtn: true,
      },
      () => {
        this.handleChangeSwitch();
        this.computeRealTargetKeys(this.props);
      }
    );
  };

  handleCloudHostChange = (e: any) => {
    this.setState({
      cloudHost: e.target.value,
      showBtn: true,
    });
  };

  handleCancel = (e: any) => {
    this.props.handleCancel();

    const f: any = [];
    this.props.hostList.forEach((q: any) => {
      if (this.props.existIp.includes(q.ip)) {
        f.push(q.id);
      }
    });
    this.setState({
      selectedKeys: [],
      targetKeys: f,
      isCloud_state: this.props.isCloud,
    });
  };

  handleInputBlur = () => {
    console.log('失去焦点了哦');
    this.props.handleResourceSubmit({
      isCloud: this.state.isCloud_state,
      hosts:
        this.state.isCloud_state || !this.props.hasInstance
          ? this.state.cloudHost
          : this.state.targetKeys,
    });
  };

  filterOption = (inputValue: any, option: any) =>
    option.ip.indexOf(inputValue) > -1;

  handleGetAutoConfig = () => {
    installGuideService
      .getAutoConfig({
        productName: this.props.selectedProduct.product_name,
        serviceName: this.props.serviceKey,
        productVersion: this.props.selectedProduct.product_version,
      })
      .then((res) => {
        console.log(res);
        if (res.data.code === 0) {
          // this.setState({showLoading:'block'})
          const baseClusterId = this.props.installGuideProp.baseClusterId;
          this.props.actions.getProductServicesInfo(
            {
              productName:
                this.props.installGuideProp.selectedProduct.product_name,
              productVersion:
                this.props.installGuideProp.selectedProduct.product_version,
              unSelectService:
                this.props.installGuideProp.unSelectedServiceList,
              baseClusterId: baseClusterId === -1 ? undefined : baseClusterId,
              clusterId: this.props.installGuideProp.clusterId,
            },
            (res: any) => {
              for (const fk in res) {
                for (const sk in res[fk]) {
                  if (
                    sk ===
                    this.props.installGuideProp.selectedService.serviceKey
                  ) {
                    this.props.actions.setSelectedConfigService({
                      ...res[fk][sk],
                      serviceKey: sk,
                    });
                  }
                }
              }
              this.props.actions.updateServiceHostList({
                productName:
                  this.props.installGuideProp.selectedProduct.product_name,
                serviceName:
                  this.props.installGuideProp.selectedService.serviceKey,
              });
            }
          );
          this.setState({ showAutoBtn: false });
          message.success('配置完成！');
        } else {
          message.error('配置失败！');
        }
      });
  };

  handleAutoConfig = () => {
    this.setState({ showAutoBtn: true });
    const that = this;
    Modal.confirm({
      title: '确定使用自动配置吗？',
      content: '系统将会为该服务自动分配主机。',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk() {
        that.handleGetAutoConfig();
      },
      onCancel() {
        console.log('Cancel');
        that.setState({ showAutoBtn: false });
      },
    });
  };

  trasnferItemRender = (ip: string) => {
    let target = window.hostRoleMap.find((item) => item.ip === ip);
    if (!target) {
      target = {};
    }
    const { role_info = [] } = target;
    return (
      <div className="tagList">
        <div style={{ width: '100px' }}>{ip}</div>
        <div style={{ marginLeft: '5px', flex: 1 }}>
          {role_info.map((role, index) => (
            <Tag
              key={index}
              style={{
                borderRadius: '11px',
              }}>
              {role.role_name}
            </Tag>
          ))}
        </div>
      </div>
    );
  };

  render() {
    const { dataList, isCloud_state, targetKeys, isRerro, isLerro } =
      this.state;
    const { sqlErro, selectedService } = this.props.installGuideProp;
    const { forcedUpgrade, upgradeType } = this.props.DeployProp;

    const autoBtn = (
      <div style={{ position: 'absolute', right: 30 }}>
        <Button
          icon="apartment"
          type="primary"
          ghost
          onClick={this.handleAutoConfig}
          disabled={this.state.showAutoBtn}>
          自动配置
        </Button>
      </div>
    );
    const {
      selectedProduct: { status },
      Instance = {},
    } = this.props;

    const autoBtnShow =
      !isCloud_state &&
      status === 'undeployed' &&
      Instance.MaxReplica != undefined;
    return (
      <div className="resource-container">
        <div>
          {this.props.hasInstance && (
            <Row
              style={{
                marginBottom: 18,
                display: 'flex',
                alignItems: 'center',
              }}>
              <div className="ml-10 mr-20" style={{ fontSize: 12 }}>
                使用外部主机:
              </div>
              <div>
                <Switch
                  style={{ marginTop: '-3px' }}
                  className="switch"
                  checked={this.state.isCloud_state}
                  onChange={(e) => {
                    this.useCloudChange(e);
                  }}
                />
              </div>
              {/* 注释自动配置功能 */}
              {/* {autoBtnShow && autoBtn} */}
            </Row>
          )}
        </div>
        {this.state.isCloud_state || !this.props.hasInstance ? (
          <div>
            <Row style={{ display: 'flex' }}>
              <div style={{ textAlign: 'right', paddingRight: 20, width: 116 }}>
                IP地址:
              </div>
              <div>
                <Input.TextArea
                  style={{ width: 500 }}
                  value={this.state.cloudHost}
                  onBlur={this.handleInputBlur}
                  onChange={(e) => this.handleCloudHostChange(e)}
                  placeholder="可填写多个IP地址，多个IP地址用英文逗号分割开，如172.10.16.2,172.10.20.6。"
                />
                {sqlErro && selectedService.serviceKey === 'mysql' && (
                  <div className="errRight" style={{ display: 'block' }}>
                    {sqlErro}
                  </div>
                )}
              </div>
            </Row>
          </div>
        ) : (
          !this.props.isKubernetes && (
            <Row>
              <CustomTransfer
                disabled={
                  forcedUpgrade?.includes(selectedService.serviceKey) ||
                  upgradeType !== 'smooth'
                    ? false
                    : true
                }
                rowKey={(record: any) => record.id}
                dataSource={dataList}
                showSearch
                onChange={this.handleChange}
                filterOption={this.filterOption}
                selectedKeys={this.state.selectedKeys}
                targetKeys={targetKeys}
                render={(item: any) => this.trasnferItemRender(item.ip)}
                onSelectChange={this.handleSelect}
                listStyle={{
                  flex: 1,
                  height: 450,
                }}
              />
              {isLerro && (
                <span className="errRight">
                  首次平滑升级本服务，请至少保留一台主机
                </span>
              )}
              {isRerro && <span className="errLeft">请至少选择一台主机</span>}
            </Row>
          )
        )}
      </div>
    );
  }
}
export default Resource;
