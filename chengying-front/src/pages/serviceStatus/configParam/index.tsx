import * as React from 'react';
import { message, Collapse, Input, Modal, Icon, Select } from 'antd';
import {
  servicePageService,
  productService,
  userCenterService,
} from '@/services';
import { CurrentProductionParams } from '@/services/ServicePageService';
import ClassNames from 'classnames';
import Config from './config';
import ConfigComp from './configComp';
import InstanceConfigComp from './instanceComp';
import ConfigToolBar from './configToolBar';
import CodeDiff from '@/components/codeDiff';
import utils from '@/utils/utils';
import isEqual from 'lodash/isEqual';
import * as Http from '@/utils/http';

const Panel = Collapse.Panel;
const confirm = Modal.confirm;
const Option = Select.Option;

interface IState {
  dropList: any[];
  currentConfigFile: string;
  currentDiffFile: string;
  fileContent: string;
  oldParamArray: any[];
  newValue: string;
  newAddParamArray: any[];
  deleteConfig: string;
  configQuery: string;
  noHosts: string;
  visibleConfigInfo: boolean;
  visibleConfig: boolean;
  repeatParams: any[];
  allHostList: any[];
  currentHosts: string[];
  diffData: any;
  selectHost: string;
  encryptInfo: any;
}

export default class ConfigParams extends React.PureComponent<any, IState> {
  state: IState = {
    dropList: [], // 配置文件列表
    currentConfigFile: '', // 文件路径
    currentDiffFile: '', // diff 文件
    fileContent: '', // 文件内容
    oldParamArray: [], // 原有参数
    newValue: '',
    newAddParamArray: [],
    deleteConfig: '',
    configQuery: '',
    noHosts: '',
    visibleConfigInfo: false,
    visibleConfig: false,
    repeatParams: [],
    allHostList: [],
    currentHosts: [],
    diffData: { before: '', after: '' },
    selectHost: '',
    encryptInfo: {},
  };

  componentDidMount() {
    this.getConfigFile();
    // 获取主机列表
    this.getServiceInfo();
    this.getPublicKey();
  }

  componentDidUpdate(prevProps, prevState) {
    if (!isEqual(this.props.cur_service_info, prevProps.cur_service_info)) {
      this.getConfigFile();
      this.getServiceInfo();
    }
  }

  /**
   * 获取配置文件
   * @returns
   */
  getConfigFile = () => {
    const {
      cur_service_info: { service_name },
      ServiceStore: {
        cur_product: { product_name, product_version },
        configFile,
      },
    } = this.props;
    if (!service_name || !product_name || !product_version) {
      return;
    }
    servicePageService
      .getParamDropList({
        productName: product_name,
        productVersion: product_version,
        serviceName: service_name,
      })
      .then((res) => {
        res = res.data;
        // 配置文件列表
        const dropList = res?.data?.list || [];
        const defaultFile = configFile === 'all' ? undefined : configFile;
        // 当前配置文件
        const currentConfigFile = defaultFile || dropList[0] || '';
        this.setState({
          dropList,
          currentConfigFile,
          currentDiffFile: currentConfigFile,
          fileContent: '',
        });
        if (dropList.length) {
          this.handleConfigFileClick(currentConfigFile);
        }
      });
  };

  getServiceInfo = () => {
    const { cur_service_info, cur_product_info } = this.props;
    const { getHostsList } = this.props.actions;
    getHostsList(
      {
        product_name: cur_product_info?.product_name,
        service_name: cur_service_info?.service_name,
        namespace: utils.k8sNamespace,
      },
      (use_cloud, list: any) => {
        const hosts = list.map((o) => o.ip);
        this.setState({
          currentHosts: hosts,
          selectHost: hosts[0],
        });
      }
    );
  };

  getPublicKey = async () => {
    const { data } = await userCenterService.getPublicKey();
    if (data.code !== 0) {
      return;
    }
    this.setState({
      encryptInfo: data.data,
    });
  };

  // 获取配置文件内容
  handleConfigFileClick = (value: string) => {
    this.setState(
      {
        currentConfigFile: value,
      },
      () => {
        this.getParamContent(value);
      }
    );
  };

  handleDiffFileClick = (value: string) => {
    this.setState(
      {
        currentDiffFile: value,
      },
      () => {
        if (this.state.visibleConfig) {
          this.getConfDiff();
        }
      }
    );
  };

  getParamContent = (value: string) => {
    const {
      cur_service_info,
      ServiceStore: { cur_product },
    } = this.props;
    servicePageService
      .getParamContent({
        productName: cur_product.product_name,
        productVersion: cur_product.product_version,
        serviceName: cur_service_info.service_name,
        servicePath: value,
      })
      .then((res) => {
        this.handleFileChange(value);
        const fileContent = res.data.data;
        const oldParamArray = this.getParamArray(fileContent);
        this.setState({
          fileContent,
          oldParamArray: oldParamArray,
        });
      });
  };

  // 获取参数
  getParamArray = (value: string) => {
    const regex = /{{\.(\w+)}}/gm;
    let m = regex.exec(value);
    let n = [];
    while (m !== null) {
      if (m.index === regex.lastIndex) {
        regex.lastIndex++;
      }
      m.forEach((match, groupIndex) => {
        if (groupIndex === 1) {
          n.push(match);
        }
      });
      m = regex.exec(value);
    }
    n = Array.from(new Set(n));
    return n;
  };

  // 配置文件内容变更
  handleFocusOut = (value: any, doc: any) => {
    // this.setState({
    //   newValue: value,
    //   fileContent: value
    // })
    const { oldParamArray } = this.state;
    const deleteConfig = [];
    const n = this.getParamArray(value);
    oldParamArray.forEach((param: any) => {
      if (!n.includes(param)) {
        deleteConfig.push(param);
      } else {
        const index = n.findIndex((item: any) => item === param);
        n.splice(index, 1);
      }
    });
    this.setState({
      newValue: value,
      fileContent: value,
      newAddParamArray: n,
      deleteConfig: deleteConfig.toString(),
    });
    console.log('删除的：', deleteConfig.toString(), '新的：', n);
  };

  // 参数配置 - 提交
  handleSetParamSubmit = () => {
    const newAddValue = {};
    const {
      ServiceStore: { cur_product },
      cur_service_info,
      actions,
    } = this.props;
    const {
      newAddParamArray = [],
      currentConfigFile,
      newValue,
      deleteConfig,
    } = this.state;
    newAddParamArray.forEach((item) => {
      const input: HTMLInputElement = document.getElementById(item) as any;
      newAddValue[`${item}`] = input.value;
    });
    const configContent = {
      file: currentConfigFile,
      content: newValue,
      values: newAddValue,
      deleted: deleteConfig,
    };
    servicePageService
      .setParamUpdate(
        {
          productName: cur_product?.product_name,
          productVersion: cur_product?.product_version,
          serviceName: cur_service_info?.service_name,
        },
        configContent
      )
      .then((res) => {
        if (res.data.code === 0) {
          message.success('添加参数成功！');
          this.setState(
            {
              newAddParamArray: [],
              visibleConfigInfo: false,
            },
            () => {
              this.getConfigFile();
            }
          );
          // this.handleAddParamConfigShow(false);
          actions.refreshProductAndService({
            product_name: cur_product.product_name,
            product_version: cur_product.product_version,
            service_name: cur_service_info.service_name,
          });
        } else {
          message.error(res.data.msg);
        }
      });
  };

  onCloseModal = () => {
    this.setState({
      // dropList: null,
      // currentConfigFile: '',
      // fileContent: '',
      newAddParamArray: [],
    });
    this.handleAddParamConfigShow(false);
  };

  // 获取diff config
  getConfDiff = () => {
    const {
      ServiceStore: { cur_product },
      cur_service_info,
    } = this.props;
    const { selectHost, currentDiffFile } = this.state;
    const { product_name, product_version } = cur_product;
    const { service_name } = cur_service_info;
    const param = {
      file: currentDiffFile,
      ip: selectHost || '',
      product_name,
      product_version,
      service_name,
    };
    servicePageService.getConfDiff(param).then((res: any) => {
      const { code, data, msg } = res.data;
      let result = {
        before: '',
        after: '',
      };
      if (code === 0) {
        result = { ...result, ...data };
      } else {
        message.error(msg);
      }
      this.setState(
        {
          diffData: result,
        },
        () => {
          console.log(this.state.diffData);
        }
      );
    });
  };

  handleServiceConfigSearch = (e: any) => {
    e.stopPropagation();
    this.setState({
      configQuery: e.target.value,
    });
  };

  // 获取当前服务下配置的所有主机列表
  getCurrentHostsList = () => {
    const { cur_product_info, cur_service_info, HeaderStore } = this.props;
    const { cur_parent_cluster } = HeaderStore;
    const { product_name } = cur_product_info;
    const { service_name } = cur_service_info;
    if (!product_name || !service_name) {
      return;
    }
    Http.get(
      `/api/v2/product/${product_name}/service/${service_name}/selected_hosts?clusterId=${cur_parent_cluster.id}`,
      {}
    ).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        this.setState(
          {
            allHostList: data.data.hosts,
          },
          () => {
            this.handleSaveMoreConfig(); // 含有多配置的时候先调用下获取全部主机，然后再去进行保存功能
          }
        );
      } else {
        message.error(data.msg);
      }
    });
  };

  // 保存存在主机关联的配置信息
  handleSaveMoreConfig = () => {
    const { ServiceStore, cur_service_info } = this.props;
    const {
      cur_product: { product_name },
      cur_service,
    } = ServiceStore;
    let noHosts = '';
    // 判断一下是否都满足了绑定了主机
    const canMore = Object.keys(cur_service.configModify).every((obj) => {
      // 判断关联主机的是否都是数组
      if (Array.isArray(cur_service.configModify[obj])) {
        const arr = cur_service.configModify[obj].map((key) => {
          return key.hosts;
        });
        const { allHostList } = this.state;
        const selectHosts = arr.join(',').split(',');
        // 遍历获取全部的每个参数下的hosts并以数组形式存储
        if (
          allHostList &&
          selectHosts.length === allHostList.length &&
          allHostList.some((host) => selectHosts.indexOf(host) !== -1)
        ) {
          return true;
        } else {
          noHosts = obj.split('.')[1];
          return false;
        }
      } else {
        return true;
      }
    });
    // 如果满足了
    if (canMore) {
      Object.keys(cur_service.configModify).map((config) => {
        const saveParams = {
          [config]: cur_service.configModify[config],
        };
        // 判断是数组吗，是的就走多配置接口。
        if (Array.isArray(cur_service.configModify[config])) {
          servicePageService
            .modifyMultiAllHosts(
              {
                product_name: product_name,
                service_name: cur_service_info.service_name,
              },
              saveParams
            )
            .then((res: any) => {
              res = res.data;
              if (res.code === 0) {
                // ServiceStore.cur_service.configModify = {};
                delete cur_service.configModify[config];
                this.changeHosts();
                message.success('关联成功!');
              } else if (res.code === 100 && res.msg.includes('输入值重复')) {
                const repeatArr = res.msg
                  .substring(1, res.msg.indexOf(')'))
                  .split(',');
                this.setState({
                  repeatParams: repeatArr,
                  noHosts: '',
                });
              } else {
                message.error(res.msg);
              }
            });
        } else {
          // 不是数组那条走老接口
          servicePageService
            .modifyProductConfigAll(
              {
                product_name: product_name,
                service_name: cur_service_info.service_name,
              },
              saveParams
            )
            .then((res: any) => {
              res = res.data;
              if (res.code === 0) {
                delete cur_service.configModify[config];
                message.success('保存完成!');
                this.setState({
                  noHosts: '',
                });
              } else {
                message.error(res.msg);
              }
            });
        }
        // 在这里处理如果保存接口都调用成功了，那么才会给他切换内容，刷新服务，修改按钮保存置灰，否则当前页面还是修改那个状态
        if (JSON.stringify(cur_service.configModify) === '{}') {
          cur_service.configModify = {};
          this.setCurrentService(cur_service, cur_service_info.service_name);
        }
      });
    } else {
      // 没满足找出来那个参数
      this.setState({
        noHosts: noHosts,
      });
    }
  };

  // 点击保存配置信息
  handleSaveServiceConfig = () => {
    const { authorityList, ServiceStore, cur_service_info } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'service_config_edit')) {
      return;
    }
    const {
      cur_product: { product_name },
      cur_service,
    } = ServiceStore;
    // 判断是否是数组
    if (Object.values(cur_service.configModify).some((x) => Array.isArray(x))) {
      // 判断数组下的配置是否是关联了主机
      if (
        Object.values(cur_service.configModify).some(
          (x) => Array.isArray(x) && x.some((y) => y.hosts === '')
        )
      ) {
        Object.keys(cur_service.configModify).some((x) => {
          if (
            Array.isArray(cur_service.configModify[x]) &&
            cur_service.configModify[x].some((y) => y.hosts === '')
          ) {
            // console.log(x, '配置里的主机没关联')
            this.setState({
              noHosts: x.split('.')[1],
            });
          }
        });
      } else {
        this.changeHosts(); // 清掉可能搞好了不需要的参数提示
        this.getCurrentHostsList(); // 含有多配置的时候去触发一下方法获取全部主机
      }
    } else {
      servicePageService
        .modifyProductConfigAll(
          {
            product_name: product_name,
            service_name: cur_service_info.service_name,
          },
          cur_service.configModify
        )
        .then((res: any) => {
          res = res.data;
          if (res.code === 0) {
            cur_service.configModify = {};
            // 当前服务
            this.setCurrentService(cur_service, cur_service_info.service_name);
            message.success('保存完成!');
            this.changeHosts();
          }
        });
    }
  };

  handleDiffHost = (selectHost) => {
    this.setState(
      {
        selectHost,
      },
      () => {
        this.getConfDiff();
      }
    );
  };

  // 配置文件选择
  handleFileChange = (value: string) => {
    const {
      ServiceStore: { cur_product },
      cur_service_info: { service_name },
      actions,
    } = this.props;
    const { product_name, product_version } = cur_product;
    if (!product_name || !product_version || !service_name) {
      return;
    }
    actions.getServiceConfigList({
      product_name: product_name,
      product_version: product_version,
      service_name: service_name,
      file: value,
    });
    actions.setConfigFile(value);
  };

  // 重置
  handleResetServiceConfig = () => {
    const { authorityList, ServiceStore, cur_service_info } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'service_config_edit')) {
      return;
    }
    const {
      cur_product: { product_name },
    } = ServiceStore;
    if (!product_name) {
      return;
    }
    const params: CurrentProductionParams = { product_name: product_name };
    utils.k8sNamespace && (params.namespace = utils.k8sNamespace);
    servicePageService.getCurrentProduct(params).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        let service = {};
        for (const s in data.data.product.Service) {
          if (cur_service_info.service_name === s) {
            service = data.data.product.Service[s];
          }
        }
        this.setCurrentService(service, cur_service_info.service_name);
        this.changeHosts();
      } else {
        message.error(data.msg);
      }
    });
  };

  setCurrentService = (service, service_name) => {
    const { cur_product } = this.props.ServiceStore;
    this.props.actions.setCurrentService(service, {
      product_name: cur_product.product_name,
      service_name: service_name,
      product_version: cur_product.product_version,
    });
  };
  // 参数配置弹框显示
  handleAddParamConfigShow = (visibleConfigInfo: boolean) => {
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'service_config_edit')) {
      return;
    }
    this.setState({ visibleConfigInfo });
  };

  handleDispath = () => {
    const { authorityList } = this.props;
    if (!authorityList.service_config_distribute) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    const { visibleConfig } = this.state;
    if (!visibleConfig) {
      this.getConfDiff();
    }
    this.setState({
      visibleConfig: !visibleConfig,
    });
  };

  /**
   * 配置下发
   */
  distributeServiceConfig = () => {
    const { authorityList, ServiceStore, cur_service_info } = this.props;
    const { visibleConfig, diffData } = this.state;
    if (utils.noAuthorityToDO(authorityList, 'service_config_distribute')) {
      return;
    }
    const { cur_product } = ServiceStore;
    //
    if (visibleConfig && diffData.before !== diffData.after) {
      confirm({
        title: '确定将运行配置下发至各主机吗？',
        icon: <Icon type="exclamation-circle" theme="filled" />,
        className: 'rollstart-comfirm-dialog',
        onOk: () => {
          productService
            .distributeServiceConfig({
              serviceName: cur_service_info.service_name,
              productId: cur_product.id,
            })
            .then((res) => {
              if (res.data.code === 0) {
                message.success('配置下发成功！');
              } else {
                message.error('配置下发失败！');
              }
            });
          this.handleDispath();
        },
        onCancel() {},
      });
    } else {
      this.handleDispath();
    }
  };

  // 改变未关联主机情况
  changeHosts = () => {
    this.setState({
      noHosts: '',
      repeatParams: [],
    });
  };

  // 滚动重启服务
  handleRestartServiceInTurn = () => {
    const {
      authorityList,
      ServiceStore: { cur_product },
      products,
      cur_service_info,
      actions,
    } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'service_roll_restart')) {
      return;
    }
    let pid = -1;
    for (const p of products) {
      if (p.product_name === cur_product.product_name) {
        pid = p.id;
      }
    }
    confirm({
      title: '确定要滚动重启该服务下的所有主机吗？',
      content:
        '重启后最新的运行配置将生效。重启过程中会出现部分主机上服务暂时停止的情况，但服务整体将正常运行。',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      className: 'rollstart-comfirm-dialog',
      onOk() {
        actions.setServiceRollRestartState(cur_service_info.service_name, true);
        return new Promise((resolve, reject) => {
          setTimeout(resolve, 500);
          servicePageService
            .setServiceRollRestart({
              pid: pid,
              service_name: cur_service_info.service_name,
            })
            .then((res: any) => {
              // debugger;
              const { code, msg } = res.data;
              code === 0
                ? message.success('滚动重启完成')
                : message.success(msg);
              // res = res.data;
              // debugger;
              actions.setServiceRollRestartState(
                cur_service_info.service_name,
                false
              );
            })
            .catch(() => {
              message.success('重启完成');
              actions.setServiceRollRestartState(
                cur_service_info.service_name,
                false
              );
            });
        });
      },
      onCancel() {
        actions.setServiceRollRestartState(
          cur_service_info.service_name,
          false
        );
      },
    });
  };

  render() {
    const {
      ServiceStore: { cur_product, configFile, cur_service },
      HeaderStore,
      services,
      isKubernetes,
      cur_product_info,
      authorityList,
      getCurrentProduct,
      cur_service_info,
    } = this.props;
    const {
      dropList,
      currentConfigFile,
      fileContent,
      visibleConfigInfo,
      currentDiffFile,
      visibleConfig,
      noHosts,
      currentHosts,
      configQuery,
      diffData,
      selectHost,
      repeatParams,
      encryptInfo,
    } = this.state;

    const couldSaveConfig =
      cur_service?.configModify &&
      Object.keys(cur_service?.configModify).length;
    let isRestart: any = null;
    if (
      cur_service &&
      cur_product?.product?.Service[cur_service.service_name]?.hosts
    ) {
      isRestart =
        cur_product?.product?.Service[cur_service.service_name].isRestart;
    }
    const CAN_CONFIG_EDIT = authorityList.service_config_edit;

    const extraSider = (
      <div className="extra-sider">
        <Select
          value={selectHost}
          onChange={this.handleDiffHost}
          placeholder="选择主机"
          size="small">
          {currentHosts.map((item) => (
            <Option key={item}>{item}</Option>
          ))}
        </Select>
        <h2>配置文件</h2>
        <ul
          className="extra-ul"
          onClick={(e) => this.handleDiffFileClick(e.target?.title)}>
          {(cur_service?.Instance?.ConfigPaths || []).map(
            (file: string, index: number) => (
              <li
                title={file}
                key={file}
                className={ClassNames('', {
                  'extra-selected': file === currentDiffFile,
                })}>
                {file}
              </li>
            )
          )}
        </ul>
      </div>
    );

    const settingQuery: React.ReactNode = (
      <div className="run-setting">
        运行配置
        <Input
          placeholder="输入参数名称搜索"
          value={configQuery}
          onClick={(e) => {
            e.stopPropagation();
          }}
          autoComplete="off"
          onChange={this.handleServiceConfigSearch}
        />
      </div>
    );

    return (
      <>
        <div style={{ maxHeight: 'calc(100vh - 180px)', overflowY: 'auto' }}>
          <Collapse
            className="ant-collapse-no-border box-shadow-style collapse-config"
            defaultActiveKey={['1']}>
            {/* 运行配置 */}
            <Panel key="1" header={settingQuery}>
              <ConfigToolBar
                cur_service={cur_service}
                configFile={configFile}
                handleFileChange={this.handleFileChange}
                isRestart={isRestart}
                handleRestartServiceInTurn={this.handleRestartServiceInTurn}
                couldSaveConfig={couldSaveConfig}
                handleSaveServiceConfig={this.handleSaveServiceConfig}
                distributeServiceConfig={this.handleDispath}
                handleAddParamConfigShow={this.handleAddParamConfigShow}
                handleResetServiceConfig={this.handleResetServiceConfig}
                isKubernetes={isKubernetes}
              />
              {!!cur_service_info && (
                <ConfigComp
                  canEditHost={authorityList?.sub_menu_scheme_host_associated}
                  canEditPwd={authorityList?.sub_menu_scheme_password_lock}
                  canReset={authorityList?.reset_service_config_edit}
                  sname={cur_service_info.service_name}
                  pname={cur_product_info.product_name}
                  query={configQuery}
                  pid={cur_product_info.product_id}
                  HeaderStore={HeaderStore}
                  pversion={cur_product_info.product_version}
                  noEditAuthority={!CAN_CONFIG_EDIT}
                  noHosts={noHosts}
                  changeHosts={this.changeHosts}
                  isKubernetes={isKubernetes}
                  repeatParams={repeatParams}
                  encryptInfo={encryptInfo}
                />
              )}
            </Panel>
            <Panel header="部署配置" key="2">
              {cur_service_info ? (
                <InstanceConfigComp
                  sname={cur_service_info.service_name}
                  pname={cur_product_info.product_name}
                  pid={cur_product_info.product_id}
                  pversion={cur_product_info.product_version}
                  updateServiceConfig={() =>
                    getCurrentProduct(
                      cur_product_info.product_name,
                      cur_product_info.product_version,
                      true
                    )
                  }
                />
              ) : null}
            </Panel>
            <Panel header="依赖组件" key="3">
              {cur_service_info &&
              services &&
              services[cur_service_info.service_name]?.DependsOn
                ? services[cur_service_info.service_name].DependsOn.join(',')
                : '无'}
            </Panel>
          </Collapse>
        </div>
        {visibleConfigInfo && (
          <Config
            visible={visibleConfigInfo}
            onCancel={this.onCloseModal}
            handleSetParamSubmit={this.handleSetParamSubmit}
            dropList={dropList}
            onChange={this.handleConfigFileClick}
            value={currentConfigFile}
            fileContent={fileContent}
            handleContentChange={this.handleFocusOut}
            newAddParamArray={this.state.newAddParamArray}
          />
        )}
        {visibleConfig && (
          <CodeDiff
            visible={visibleConfig}
            title="配置下发（请确认变更前后的配置，若配置前后无变化，则无需重启服务）"
            extra={extraSider}
            data={[diffData.before, diffData.after]}
            handleCancle={this.handleDispath}
            handleSubmit={this.distributeServiceConfig}
          />
        )}
      </>
    );
  }
}
