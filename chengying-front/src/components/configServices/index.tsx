import * as React from 'react';
import { Tabs, Collapse, Modal, Icon, message, Select } from 'antd';
import { Prompt } from 'react-router-dom';
import Config from './config.resource';
import ParamConfig from './config.param.config';
import InstanceConfig from './config.param.instance';
import * as Http from '@/utils/http';
import {
  installGuideService,
  servicePageService,
  userCenterService,
} from '@/services';
import CodeDiff from '@/components/codeDiff';
import { cloneDeep, isEqual } from 'lodash';
import ClassNames from 'classnames';
import './style.scss';

const TabPane = Tabs.TabPane;
const Panel = Collapse.Panel;
const Option = Select.Option;

interface Prop {
  defaultActiveKey?: 'resource' | 'param';
  existIp?: []; // 已经存在的主机（区分自定ip和transfer ip）
  hostList: any[]; // 主机列表
  serviceData: any; // 勾选的服务数据
  saveResourceState: Function; // 保存主机勾选结果
  resourceState: any; // 主机列表包含targetKeys(已选择)和selectedKeys(已勾选)
  afterParamFieldChange: Function; // 配置参数项发生改变后调用
  resetParamFieldvalue: Function; // 点击重置参数项后调用
  commitParamChange: Function; // 点击保存修改后调用-param config
  handleCancel?: Function; // 点击取消后调用
  handleResourceSubmit: Function; // 点击保存调用-host
  paramChangeBlur: Function; // 参数配置项失去焦点后调用
  selectedProduct: any;
  installGuideProp: any;
  actions: any;
  isKubernetes: boolean; // 部署包类型
  sname: string;
  pname: string;
  noHosts: string; // 没有绑定主机多配置的情况
  repeatParams: any[]; /// 输入值重复的总输入框
  changeHosts: Function; // 修改删除为一行事没有未关联主机提示
  getProductServicesInfo: Function;
  DeployProp: any;
}

interface State {
  selectedHost: any;
  showBtnBox: boolean;
  serviceInfo: any; // 用本地变量代替props变量
  changeFields: any; // 修改过的字段存放在这里
  resetField: string; // 重置的字段
  cancel: boolean; // 是否是取消
  isCommitHost: boolean; // 是否是提交主机变化导致的prop更新
  activeTabpane: string;
  Instate: any;
  hostsSelectList: any[];
  activeKey: string;
  visibleConfig: boolean;
  diffData: any;
  currentHosts: string[];
  selectHost: string;
  configFileList: string[];
  configFile: string;
  selectedConfigFile: string; // 配置文件
  encryptInfo: any;
  isRerro: boolean;
  isLerro: boolean;
  smoothTarget: any[];
  isChange: boolean;
}

class ConfigServices extends React.Component<Prop, State> {
  state: State = {
    selectedHost: [],
    showBtnBox: false,
    serviceInfo: {},
    changeFields: {},
    resetField: '',
    cancel: false,
    isCommitHost: false,
    activeTabpane: this.props.defaultActiveKey,
    Instate: {},
    hostsSelectList: [],
    activeKey: null,
    visibleConfig: false,
    diffData: {},
    currentHosts: [],
    selectHost: '',
    configFileList: [],
    configFile: '',
    selectedConfigFile: '',
    encryptInfo: {},
    isRerro: false,
    isLerro: false,
    smoothTarget: [],
    isChange: false,
  };

  componentDidMount() {
    this.setState({
      serviceInfo: { ...this.configData(this.props.serviceData) },
    });
    this.computeTransferIpList();
    this.getPublicKey();
    this.getConfigFile();
    this.initErro();
  }

  componentDidUpdate(prevProps, prevState) {
    if (
      !isEqual(
        this.props.serviceData.serviceKey,
        prevProps.serviceData.serviceKey
      ) || this.state.isChange
    ) {
      const Instate = this.state.serviceInfo.Instance;
      this.setState(
        {
          serviceInfo: { ...this.configData(this.props.serviceData) },
          selectedConfigFile: '',
          Instate: Instate,
          cancel: false,
          isCommitHost: false,
        },
        () => {
          this.computeTransferIpList();
          this.getConfigFile();
          this.initErro();
          this.setState({isChange: false})
        }
      );
    }
    if (this.state.resetField !== '') {
      // 重置更新的prop
      this.changeServiceInfo(
        this.state.resetField,
        this.getFieldValue(this.props, this.state.resetField)
      );
      this.setState({
        resetField: '',
      });
    }
  }

  getPublicKey = async () => {
    const { data } = await userCenterService.getPublicKey();
    if (data.code !== 0) {
      return;
    }
    this.setState({
      encryptInfo: data.data,
    });
  };

  // 获取当前服务下配置的所有主机列表
  getHostsList = () => {
    const { sname, pname, installGuideProp } = this.props;
    if (!sname || !pname) {
      return;
    }
    Http.get(
      `/api/v2/product/${pname}/service/${sname}/selected_hosts?clusterId=${installGuideProp.clusterId}`,
      {}
    ).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        this.setState({
          hostsSelectList: data.data.hosts,
        });
      } else {
        message.error(data.msg);
      }
    });
  };

  /**
   * 获取配置文件
   * @returns
   */
  getConfigFile = () => {
    const {
      pname,
      selectedProduct: { product_version },
      sname,
    } = this.props;
    servicePageService
      .getParamDropList({
        productName: pname,
        productVersion: product_version,
        serviceName: sname,
      })
      .then((res) => {
        res = res.data;
        // 配置文件列表
        const configFileList = res?.data?.list || [];
        const defaultFile = configFileList.length ? configFileList[0] : '';
        this.setState({
          configFileList,
          configFile: defaultFile,
        });
      });
  };

  configData = (stateData) => {
    if (stateData) {
      const data = { ...stateData };
      for (const title in data.Config) {
        data.Config[title] = {
          ...data.Config[title],
          iconType: 0,
        };
      }
      return data;
    } else {
      return {};
    }
  };

  handleResourceSubmit = (e: any) => {
    this.setState({
      isCommitHost: true,
    });
    this.props.handleResourceSubmit(e);
  };

  getFieldValue = (prop: Prop, f: string) => {
    const location = f.split('.');
    let q: any = null;

    location.forEach((o: any) => {
      !q ? (q = prop.serviceData[o]) : (q = q[o]);
    });

    return q;
  };

  afterParamFieldChange = (v: any, f: string, type: string) => {
    if (type === 'Instance') {
      this.changeServiceInfo(f, v);
      this.setState({
        changeFields: Object.assign({}, this.state.changeFields, {
          [f]: v,
        }),
      });
      return;
    }
    if (Array.isArray(v)) {
      this.setConfig(f.split('.')[1], v);
      this.setState({
        changeFields: Object.assign({}, this.state.changeFields, {
          [f]: v,
        }),
      });
    } else {
      this.changeServiceInfo(f, v);
      this.setState({
        changeFields: Object.assign({}, this.state.changeFields, {
          [f]: v,
        }),
      });
    }
  };

  handleParamsBlur = (v: any, f: any) => {
    this.props.paramChangeBlur(f, v);
  };

  // 更改本地的service state
  changeServiceInfo = (f: string, v: string) => {
    // const { serviceInfo } = this.state;
    const serviceInfo = cloneDeep(this.state.serviceInfo);
    const location = f.split('.');

    if (location.length > 2) {
      // 数组或json
      if (serviceInfo[location[0]][location[1]] instanceof Array) {
        // 数组--切割替换
        serviceInfo[location[0]][location[1]].splice(
          parseInt(location[2]),
          1,
          v
        );
      } else {
        // json
        serviceInfo[location[0]][location[1]][location[2]] = v;
      }
    } else {
      // string
      serviceInfo[location[0]][location[1]] = v;
    }
    this.setState({
      serviceInfo,
    });
  };

  resetParamFieldvalue = (e: any) => {
    Modal.confirm({
      title: '确认恢复此字段默认值吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        this.props.resetParamFieldvalue(e, (res: any) => {
          if (res.code === 0) {
            if (
              Object.keys(this.state.changeFields).indexOf(e.field_path) !== -1
            ) {
              delete this.state.changeFields[e.field_path];
            }

            this.setState({
              resetField: e.field_path,
              showBtnBox: true,
            });
          }
        });
      },
    });
  };

  syncHost = (e: any[]) => {
    this.setState({
      selectedHost: e,
    });
  };

  handleCancel = () => {
    this.setState({
      changeFields: {},
      showBtnBox: false,
      cancel: true,
    });
    this.props.handleCancel();
  };

  commitParamChange = () => {
    this.props.commitParamChange(this.state.changeFields);
  };

  computeTransferIpList = () => {
    const { ServiceAddr = {} } = this.props.serviceData;
    const { upgradeType, isFirstSmooth } = this.props.DeployProp;
    const filterRepeat: any = [];
    const result: any = [];
    const hostL = this.props.hostList;
    if (upgradeType !== 'smooth') {
      hostL.forEach((o: any) => {
        if (filterRepeat.indexOf(o.ip) === -1 && o.ip) {
          result.push({
            id: o.ip,
            ip: o.ip,
            key: o.ip,
          });
          filterRepeat.push(o.ip);
        }
      });
    }
    if (
      !this.props.serviceData.Instance ||
      !this.props.serviceData.Instance.UseCloud
    ) {
      const selected = ServiceAddr.Select ? ServiceAddr.Select : [];
      selected.map((item) => {
        if (ServiceAddr?.IP) {
          if (ServiceAddr?.IP.includes(item.IP)) {
            item.disabled = true;
          }
        }
        return item;
      });
      const unSelected = ServiceAddr.UnSelect ? ServiceAddr.UnSelect : [];
      unSelected.map((item) => {
        item.disabled = false;
        return item;
      });
      if (upgradeType === 'smooth') {
        const allHostsList = [...selected, ...unSelected];
        allHostsList.forEach((o: any) => {
          result.push({
            id: o.IP,
            ip: o.IP,
            key: o.IP,
            disabled: !isFirstSmooth ? o.disabled : false,
          });
        });
        this.setState({
          currentHosts: result,
          selectHost: selected[0]?.IP,
          smoothTarget: [...selected].map((item) => item.IP),
        });
        return result;
      }
      const allHosts = [...selected, ...unSelected].map((item) => item.IP);
      allHosts.forEach((o: any) => {
        if (filterRepeat.indexOf(o) === -1 && o) {
          result.push({
            id: o,
            ip: o,
            key: o,
          });
          filterRepeat.push(o);
        }
      });
      this.setState({
        currentHosts: result,
        selectHost: selected[0]?.IP,
      });
      return result;
    } else {
      this.setState({
        currentHosts: result,
        selectHost: Array.isArray(ServiceAddr?.IP) ? ServiceAddr?.IP[0] : '',
      });
      return result;
    }
  };

  // 修改config
  setConfig = (key: string, values: any[] | string, index?: number) => {
    const { serviceInfo } = this.state;
    if (index) {
      serviceInfo.Config[key].Value[index] = values;
    } else {
      serviceInfo.Config[key].Value = values;
    }

    this.setState({
      serviceInfo,
    });
  };

  // 切换Collapse
  changeCollapse = (e: any) => {
    this.setState({
      activeKey: e,
    });
  };

  // 却换tabs
  changeTabs = (e: any) => {
    this.initErro();
    this.setState({ activeTabpane: e });
    if (e == 'param') {
      this.getHostsList()
      this.setState({ isChange: true })
    } else {
      this.setState({ isChange: false })
    }
  };

  // 判断初始状态值
  initErro = () => {
    const { selectedService } = this.props.installGuideProp;
    const { ServiceAddr } = selectedService;
    const { forcedUpgrade, isFirstSmooth } = this.props.DeployProp;
    // 是否属于可平滑升级的组件
    if (forcedUpgrade.includes(selectedService.serviceKey) && isFirstSmooth) {
      if (!ServiceAddr.UnSelect) {
        this.setState({ isLerro: true });
      } else {
        this.setState({ isLerro: false });
      }
      if (!ServiceAddr.Select) {
        this.setState({ isRerro: true });
      } else {
        this.setState({ isRerro: false });
      }
    } else {
      this.setState({ isLerro: false, isRerro: false });
    }
  };

  // 获取diff config
  getConfDiff = () => {
    const {
      pname,
      selectedProduct: { product_version },
      sname,
    } = this.props;
    const { selectHost, configFile } = this.state;
    const param = {
      file: configFile || '',
      ip: selectHost || '',
      product_name: pname,
      product_version: product_version,
      service_name: sname,
    };
    servicePageService.getConfDiff(param).then((res: any) => {
      const { code, data, msg } = res.data;
      let result = {
        before: '',
        after: '',
      };
      if (code === 0) {
        result = { ...data };
      } else {
        message.error(msg);
      }
      this.setState({
        diffData: result,
      });
    });
  };

  handleDispath = () => {
    if (!this.state.visibleConfig) {
      this.getConfDiff();
    }
    this.setState({
      visibleConfig: !this.state.visibleConfig,
    });
  };

  // 前端修改运行配置
  setAllConfig = (key: string, config: any) => {
    const { serviceInfo } = this.state;
    serviceInfo[key] = config;
    console.log(key, config);
    console.log('serviceInfo', serviceInfo);
    this.setState({
      serviceInfo,
    });
  };

  recursionObjKeyValue = (preKeyString, obj, arr) => {
    if (Array.isArray(obj)) {
      obj.forEach((item: any, index: null) => {
        arr.push({
          field_path: `${preKeyString}.${index}`,
          field: item,
        });
      });
    } else {
      for (const key in obj) {
        const value = obj[key];
        if (typeof value === 'string') {
          arr.push({
            field_path: `${preKeyString}.${key}`,
            field: value,
          });
        } else {
          this.recursionObjKeyValue(`${preKeyString}.${key}`, value, arr);
        }
      }
    }
  };

  // 配置文件选择
  handleFileChange = (value: string) => {
    this.setState(
      {
        configFile: value,
      },
      () => {
        if (this.state.visibleConfig) {
          this.getConfDiff();
        }
      }
    );
  };

  handleConfigFileChange = (value: string) => {
    this.setState(
      {
        selectedConfigFile: value,
      },
      () => {
        this.getCurrentFileConf(value);
      }
    );
  };

  getCurrentFileConf = (file) => {
    const {
      clusterId,
      selectedProduct: { product_version },
    } = this.props.installGuideProp;
    const { pname, sname } = this.props;
    installGuideService
      .getServiceGroupFile(
        {
          productName: pname,
          productVersion: product_version,
        },
        {
          servicename: sname,
          clusterId,
          file,
        }
      )
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.setAllConfig('Config', res.data);
        } else {
          message.error(res.msg);
        }
      });
  };

  // 保存全部配置
  saveAllConfig = (
    type: string,
    config: any,
    editFiled?: any,
    isAlert?: boolean
  ) => {
    const keyValueArr = [];
    for (const field in config) {
      if (type === 'Config') {
        const value = config[field].Value;
        if (typeof value === 'string') {
          keyValueArr.push({
            field_path: `Config.${field}.Value`,
            field: value,
          });
        } else {
          this.recursionObjKeyValue(
            `Config.${field}.Value`,
            value,
            keyValueArr
          );
        }
      } else {
        const value = config[field];
        if (typeof value === 'string') {
          keyValueArr.push({
            field_path: `Instance.${field}`,
            field: value,
          });
        } else {
          this.recursionObjKeyValue(`Instance.${field}`, value, keyValueArr);
        }
      }
    }
    const { namespace, clusterId } = this.props.installGuideProp;
    const { pname, sname } = this.props;
    const toServiceArr = [];
    keyValueArr.forEach((item: any) => {
      if (editFiled.includes(item.field_path)) {
        toServiceArr.push({
          ...item,
          namespace: namespace,
          clusterId,
        });
      }
    });

    installGuideService
      .setParamConfigFieldValueTotal(
        {
          productName: pname,
          serviceName: sname,
        },
        toServiceArr
      )
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          isAlert && message.success('保存成功');
          this.props.getProductServicesInfo(false);
          // 保存修改后获取diff
          isAlert && this.handleDispath();
        } else {
          message.error(res.msg);
          this.props.getProductServicesInfo(false);
        }
      });
  };

  // 当资源配置不做显示时，Tab默认选中参数配置
  defaultTabKey = (visibleTabResource: boolean) => {
    if (visibleTabResource) {
      return 'resource';
    } else {
      return 'param';
    }
  };

  handleDiffHost = (value) => {
    this.setState(
      {
        selectHost: value,
      },
      () => {
        this.getConfDiff();
      }
    );
  };

  render() {
    const {
      serviceData: { ServiceAddr },
    } = this.props;
    const {
      activeKey,
      activeTabpane,
      visibleConfig,
      diffData,
      currentHosts,
      selectHost,
      configFile,
      configFileList,
      encryptInfo,
      selectedConfigFile,
      isRerro,
      isLerro,
      smoothTarget,
    } = this.state;
    // const visibleTabResource = !this.props.isKubernetes || (this.props.isKubernetes && !this.props.serviceData.Instance);
    const visibleTabResource = true;

    const extraSider = (
      <div className="extra-sider">
        <Select
          value={selectHost}
          onChange={this.handleDiffHost}
          placeholder="选择主机"
          size="small">
          {ServiceAddr?.IP?.map((item) => (
            <Option key={item}>{item}</Option>
          ))}
        </Select>
        <h2>配置文件</h2>
        <ul
          className="extra-ul"
          onClick={(e) => this.handleFileChange(e.target.title)}>
          {(configFileList || []).map((file: string) => (
            <li
              title={file}
              key={file}
              className={ClassNames('', {
                'extra-selected': file === configFile,
              })}>
              {file}
            </li>
          ))}
        </ul>
      </div>
    );
    return (
      <div className="COMP_CONFIG_SERVICE" style={{ position: 'relative' }}>
        <Prompt message="" />
        <Tabs
          onChange={(e: any) => this.changeTabs(e)}
          activeKey={activeTabpane || this.defaultTabKey(visibleTabResource)}>
          {visibleTabResource && (
            <TabPane key="resource" tab="资源分配">
              <Config
                hasInstance={!!this.props.serviceData.Instance}
                maxSelected={
                  this.props.serviceData.Instance
                    ? this.props.serviceData.Instance.MaxReplica
                    : 9999
                }
                handleCancel={this.handleCancel}
                handleResourceSubmit={(e: any) => this.handleResourceSubmit(e)}
                isCloud={
                  this.props.serviceData.Instance &&
                  this.props.serviceData.Instance.UseCloud
                    ? this.props.serviceData.Instance.UseCloud
                    : false
                }
                existIp={
                  this.props.serviceData.ServiceAddr
                    ? this.props.serviceData.ServiceAddr.IP || []
                    : []
                }
                serviceKey={this.props.serviceData.serviceKey}
                selectedKeys={this.props.resourceState.selectedKeys || []}
                targetKeys={this.props.resourceState.targetKeys || []}
                syncHost={this.props.saveResourceState}
                hostList={currentHosts}
                selectedProduct={this.props.selectedProduct}
                Instance={this.state.Instate}
                installGuideProp={this.props.installGuideProp}
                DeployProp={this.props.DeployProp}
                actions={this.props.actions}
                smoothTarget={smoothTarget}
                isKubernetes={this.props.isKubernetes}
                is_Lerro={isLerro}
                is_Rerro={isRerro}
              />
            </TabPane>
          )}
          <TabPane key="param" tab="参数配置">
            <Collapse
              activeKey={activeKey}
              className="mb-20"
              accordion
              onChange={(e: any) => this.changeCollapse(e)}>
              <Panel header="运行配置" key="1">
                <ParamConfig
                  resetParamFieldvalue={this.resetParamFieldvalue}
                  saveParamValue={this.afterParamFieldChange}
                  handleParamBlur={this.handleParamsBlur}
                  config={this.state.serviceInfo.Config || {}}
                  setConfig={this.setConfig}
                  saveAllConfig={this.saveAllConfig}
                  setAllConfig={this.setAllConfig}
                  hostsSelectList={this.state.hostsSelectList}
                  noHosts={this.props.noHosts}
                  repeatParams={this.props.repeatParams}
                  changeHosts={this.props.changeHosts}
                  getHostsList={this.getHostsList}
                  installGuideProp={this.props.installGuideProp}
                  encryptInfo={encryptInfo}
                  configFileList={configFileList}
                  configFile={configFile}
                  selectedConfigFile={selectedConfigFile}
                  handleFileChange={this.handleConfigFileChange}
                />
              </Panel>
              <Panel header="部署配置" key="2">
                <InstanceConfig
                  setAllConfig={this.setAllConfig}
                  saveAllConfig={this.saveAllConfig}
                  instanceData={this.state.serviceInfo.Instance || {}}
                  resetParamFieldvalue={this.resetParamFieldvalue}
                  handleParamBlur={this.handleParamsBlur}
                  saveParamValue={this.afterParamFieldChange}
                  isCloud={
                    this.props.serviceData.Instance &&
                    this.props.serviceData.Instance.UseCloud
                      ? this.props.serviceData.Instance.UseCloud
                      : false
                  }
                />
              </Panel>
              <Panel header="依赖服务" key="3">
                {this.state.serviceInfo.DependsOn
                  ? this.state.serviceInfo.DependsOn.toString()
                  : '无'}
              </Panel>
            </Collapse>
          </TabPane>
        </Tabs>
        {visibleConfig && (
          <CodeDiff
            title="配置确认（请确认保存后的配置，第一次部署的服务不存在当前版本）"
            visible={visibleConfig}
            data={[diffData.before, diffData.after]}
            extra={extraSider}
            handleCancle={this.handleDispath}
            handleSubmit={this.handleDispath}
          />
        )}
      </div>
    );
  }
}

export default ConfigServices;
