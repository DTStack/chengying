import * as React from 'react';
import {
  Row,
  Col,
  Modal,
  Input,
  Icon,
  message,
  Checkbox,
  Select,
  Tag,
  Button,
  Tooltip,
} from 'antd';
import InputBox from './inputBox';
import { isEqual, cloneDeep } from 'lodash';
import Form, { FormComponentProps } from 'antd/lib/form';
import * as CryptoJS from 'crypto-js';
import FormItem from 'antd/lib/form/FormItem';
import { connect } from 'react-redux';
import { Dispatch, bindActionCreators } from 'redux';
import { AppStoreTypes } from '@/stores';
import { encryptSM } from '@/utils/password';
import * as installGuideAction from '@/actions/installGuideAction';
const { Option } = Select;

const formItemCheckLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 6 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16 },
  },
};

interface Prop extends FormComponentProps {
  config: any;
  saveParamValue: Function;
  resetParamFieldvalue: Function;
  handleParamBlur: Function;
  setConfig: Function;
  setAllConfig: Function;
  saveAllConfig: Function;
  hostsSelectList: any;
  noHosts: string;
  repeatParams: any[];
  installGuideProp: any;
  changeHosts: Function;
  getHostsList: Function;
  runtimeState?: string;
  actions?: any;
  encryptInfo: {
    encrypt_public_key: string;
    encrypt_type: string;
  };
  configFileList: string[];
  configFile: string;
  selectedConfigFile: string;
  handleFileChange: (value: string) => void;
}
interface State {
  changeRecorder: string[];
  localState: any;
  checkVisible: boolean;
  iconType: any;
  nowTitle: any;
  password_key: any;
  hostVisible: boolean;
  linkIndex: number;
  linkName: string;
  changeHosts: any[];
  configEditState: 'edit' | 'normal';
  beforeEditConfig: any;
  editFiled: string[];
  isShowErro: boolean;
  isUnlinkHost: boolean;
}
const mapStateToProps = (state: AppStoreTypes) => ({
  runtimeState: state.InstallGuideStore.runtimeState,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(
    Object.assign(
      {},
      { editRuntimeState: installGuideAction.editRuntimeState }
    ),
    dispatch
  ),
});
@(connect(mapStateToProps, mapDispatchToProps) as any)
class ParamCompConfig extends React.Component<Prop, State> {
  disabledLink = [];
  state: State = {
    changeRecorder: [],
    localState: {},
    checkVisible: false,
    iconType: 0,
    nowTitle: '',
    password_key: '',
    hostVisible: false,
    linkIndex: undefined,
    linkName: '',
    changeHosts: [],
    configEditState: 'normal',
    beforeEditConfig: {},
    editFiled: [],
    isShowErro: false,
    isUnlinkHost: false,
  };

  componentDidMount() {
    this.initLocalState();
  }

  componentDidUpdate(prevProps: Prop) {
    if (!isEqual(prevProps.config, this.props.config)) {
      this.initLocalState();
    }
  }

  initLocalState = () => {
    const { config } = this.props;
    const localState = {};
    for (const title in config) {
      if (
        (config[title].Value || config[title].Value === '') &&
        typeof config[title].Value === 'string'
      ) {
        if (!localState[title]) {
          localState[title] = `${config[title].Value}`;
        }
      }
    }
    this.setState({ localState });
  };

  // 添加一行
  addItem = (key: string, index: number, e: any) => {
    const { config } = this.props;

    // console.log(config[key].Value, 'config[key].value-----', config[key], 'index', index)
    this.setState({
      changeRecorder: [...this.state.changeRecorder, key],
    });
    if (typeof config[key].Value === 'string') {
      const addItem = [
        {
          field: config[key].Value,
          hosts: '',
        },
        {
          field: config[key].Value,
          hosts: '',
        },
      ];
      this.props.setConfig(key, addItem);
    } else {
      config[key].Value.push({
        field: config[key].Value[0].field,
        hosts: '',
      });
      this.setState({
        isShowErro: true,
      });
      this.props.setConfig(key, config[key].Value);
    }
  };

  // 减少一行
  removeItem = (key: any, index: any, e: any) => {
    const { config } = this.props;
    if (config[key].Value.length > 2) {
      const temp = config[key].Value;
      temp.splice(index, 1);
      this.props.setConfig(key, temp);
      // console.log(config[key].Value, 'config[key].value---temp--', temp, 'index', index)
    } else {
      const removeChange = config[key].Value[0].field;
      this.props.setConfig(key, removeChange);
      this.props.handleParamBlur(config[key].Value, `Config.${key}.Value`);
      this.props.changeHosts();
    }
  };

  // 单配置输入框的值
  handleInputChange = (newValue: any, oldValue: string, field: string) => {
    if (newValue.target.value !== oldValue) {
      const key = `Config.${field}.Value`;
      const editFiled = [...this.state.editFiled];
      if (editFiled.indexOf(key) === -1) {
        editFiled.push(key);
        this.setState({
          editFiled,
        });
      }
      this.setState({
        changeRecorder: [...this.state.changeRecorder, field],
      });
      this.props.saveParamValue(newValue.target.value, key);
    }
  };

  // 多数组输入框修改信息
  handleInputArrayChange = (
    e: any,
    field: string,
    index: number,
    isBlur?: boolean
  ) => {
    const { config } = this.props;
    let newArr = [...config[field].Value].map((item) => item.field);
    if (!isBlur) {
      if (newArr.length > 0 && newArr.includes(e.target.value)) {
        this.setState({ isShowErro: true });
      } else {
        this.setState({ isShowErro: false });
      }
    }
    const spliceItem = {
      field: e.target.value,
      hosts: config[field].Value[index].hosts
        ? config[field].Value[index].hosts
        : '',
    };
    config[field].Value.splice(index, 1, spliceItem);
    this.setState({
      changeRecorder: [...this.state.changeRecorder, field],
    });
    this.props.setConfig(field, config[field].Value);
    this.props.saveParamValue(config[field].Value, `Config.${field}.Value`);
    if (isBlur) {
      this.props.handleParamBlur(config[field].Value, `Config.${field}.Value`);
    }
  };

  replacePwd = (newValue: any, field: string) => {
    this.setState({
      changeRecorder: [...this.state.changeRecorder, field],
    });
    this.props.saveParamValue(newValue, `Config.${field}.Value`);
  };

  updateIconType = (field: string) => {
    this.setState({
      changeRecorder: [...this.state.changeRecorder, field],
    });
    this.props.saveParamValue(1, `Config.${field}.iconType`);
  };

  resetIcontype = (field: string) => {
    this.setState({
      changeRecorder: [...this.state.changeRecorder, field],
    });
    this.props.saveParamValue(0, `Config.${field}.iconType`);
  };

  // 渲染节点
  renderNode = (state) => {
    const { isShowErro } = this.state;
    const { config, noHosts, installGuideProp = {}, repeatParams } = this.props;
    const { installType } = installGuideProp;
    let key = 0;
    const p: any[] = [];
    const disabled = state === 'normal';
    for (const title in config) {
      if (
        (config[title].Value || config[title].Value === '') &&
        typeof config[title].Value === 'string'
      ) {
        key++;
        p.push(
          <Row style={{ marginBottom: 20 }} key={key}>
            <Col span={10} className="param-input-label">
              <span>{title}：</span>
            </Col>
            <Col span={14} style={{ textAlign: 'left' }}>
              <InputBox
                tooltipOnClick={() =>
                  this.props.resetParamFieldvalue({
                    field_path: `Config.${title}.Value`,
                    type: '1',
                  })
                }
                // onBlur={(e) => this.handleBlur(e, title)}
                onChange={(e) =>
                  this.handleInputChange(e, config[title].Value, title)
                }
                defaultvalue={
                  typeof config[title].Value === 'string'
                    ? config[title].Value
                    : config[title].Value[0].field
                }
                title={title}
                iconType={config[title].iconType}
                inputDisabled={
                  disabled || /password/.test(title.toLowerCase())
                    ? !config[title].iconType
                    : false
                }
              />
              {/password/.test(title.toLowerCase())
                ? this.renderIcon(title)
                : null}
              {!/password/.test(title.toLowerCase()) &&
                installType !== 'kubernetes' && (
                  <span className="operate-icon">
                    <Icon
                      className="plus-icon"
                      type="plus"
                      onClick={this.addItem.bind(this, title, 0)}
                    />
                  </span>
                )}
            </Col>
          </Row>
        );
      } else if (Array.isArray(config[title].Value)) {
        p.push(
          <div>
            {config[title].Value.map((item, index) => {
              if (config[title].Value.length > 1) {
                return (
                  <Row style={{ marginBottom: 20 }} key={index}>
                    <Col span={10} className="param-input-label">
                      {index === 0 && <span>{title}：</span>}
                      {index === 0 && noHosts === title && (
                        <span
                          style={{
                            fontSize: '12px',
                            color: '#ff5f5c',
                            float: 'left',
                            marginLeft: '-130px',
                            marginRight: '10px',
                            marginBottom: '-40px',
                          }}>
                          存在主机未关联该参数
                        </span>
                      )}
                    </Col>
                    <Col span={14} style={{ textAlign: 'left' }}>
                      <InputBox
                        tooltipOnClick={() =>
                          this.props.resetParamFieldvalue({
                            field_path: `Config.${title}.Value`,
                            type: '2',
                            hosts: item.hosts,
                          })
                        }
                        // onBlur={(e) => this.handleArrayBlur(e, title, index)}
                        onChange={(e) =>
                          this.handleInputArrayChange(e, title, index)
                        }
                        onBlur={(e) => {
                          this.handleInputArrayChange(e, title, index, true);
                        }}
                        defaultvalue={
                          typeof config[title].Value === 'string'
                            ? config[title].Value
                            : config[title].Value[index].field
                        }
                        title={title}
                        iconType={config[title].iconType}
                        inputDisabled={
                          disabled || /password/.test(title.toLowerCase())
                            ? !config[title].iconType
                            : false
                        }
                      />
                      {/password/.test(title.toLowerCase())
                        ? this.renderIcon(title)
                        : null}
                      {!/password/.test(title.toLowerCase()) && (
                        <span className="operate-icon">
                          <Icon
                            className="plus-icon"
                            type="plus"
                            onClick={this.addItem.bind(this, title, index)}
                          />
                          <Icon
                            type="link"
                            className="link-icon"
                            onClick={this.linkItem.bind(this, title, index)}
                          />
                          {index > 0 && (
                            <Icon
                              className="dynamic-delete-button"
                              type="delete"
                              onClick={this.removeItem.bind(this, title, index)}
                            />
                          )}
                        </span>
                      )}
                      {repeatParams &&
                        index === 0 &&
                        repeatParams.some((x) => x === title) && (
                          <p style={{ fontSize: '12px', color: '#ff5f5c' }}>
                            参数值重复,请重新输入{isShowErro}1
                          </p>
                        )}

                      {isShowErro && index === 0 && (
                        <p style={{ fontSize: '12px', color: '#ff5f5c' }}>
                          参数值重复,请重新输入{isShowErro}
                        </p>
                      )}
                      <Row
                        style={{
                          marginTop:
                            item.hosts && item.hosts.split(',').length > 0
                              ? '20px'
                              : '0',
                        }}>
                        {item.hosts &&
                          item.hosts.split(',').map((e, e1) => {
                            return (
                              <div key={e1}>
                                <Col
                                  span={5}
                                  key={e1}
                                  style={{ marginRight: '12px' }}>
                                  <Tag
                                    visible
                                    closable
                                    onClose={this.onCloseHost.bind(
                                      this,
                                      title,
                                      e,
                                      index
                                    )}>
                                    {e}
                                  </Tag>
                                </Col>
                                {/* 使用react官方文档方法可以解析换行且不影响删除tag操作，但是这个官方文档解释类似于innerHtml操作，存在xss攻击风险 */}
                                {/* {e1 % 3 === 2 && <div dangerouslySetInnerHTML={{ __html: '&nbsp;' }} />} */}
                                {e1 % 3 === 2 && (
                                  <Row className="splice-row"></Row>
                                )}
                              </div>
                            );
                          })}
                      </Row>
                    </Col>
                  </Row>
                );
              } else {
                return (
                  <Row style={{ marginBottom: 20 }} key={index}>
                    <Col span={10} className="param-input-label">
                      <span>{title}：</span>
                    </Col>
                    <Col span={14} style={{ textAlign: 'left' }}>
                      <InputBox
                        tooltipOnClick={() =>
                          this.props.resetParamFieldvalue({
                            field_path: `Config.${title}.Value`,
                            type: '1',
                          })
                        }
                        // onBlur={(e) => this.handleArrayBlur(e, title, index)}
                        onChange={(e) =>
                          this.handleInputChange(e, config[title].Value, title)
                        }
                        defaultvalue={
                          typeof config[title].Value === 'string'
                            ? config[title].Value
                            : config[title].Value[index].field
                        }
                        title={title}
                        iconType={config[title].iconType}
                        inputDisabled={
                          disabled || /password/.test(title.toLowerCase())
                            ? !config[title].iconType
                            : false
                        }
                      />
                      {/password/.test(title.toLowerCase())
                        ? this.renderIcon(title)
                        : null}
                      {!/password/.test(title.toLowerCase()) && (
                        <span className="operate-icon">
                          <Icon
                            className="plus-icon"
                            type="plus"
                            onClick={this.addItem.bind(this, title, index)}
                          />
                        </span>
                      )}
                    </Col>
                  </Row>
                );
              }
            })}
          </div>
        );
      }
    }

    if (p.length === 0) {
      p.push(<p key="1">无</p>);
    }

    return p;
  };

  // 过滤已选主机
  filterHosts(data: any[]) {
    const disabledArray = [];
    Array.isArray(data) &&
      data.map((item, index) => {
        const a = item.hosts && item.hosts.split(',');
        if (a && a.length >= 1) {
          for (const i in a) {
            if (a[i] !== '' && typeof a[i] === 'string') {
              disabledArray.push(a[i]);
            }
          }
        }
        return disabledArray;
      });
    this.disabledLink = disabledArray;
  }

  // 叉掉关联的tag标签
  onCloseHost = (key: string, item: string, index: number, e: any) => {
    const { config } = this.props;
    const itemHosts = config[key].Value[index].hosts.split(',');
    const removeIndex = itemHosts.indexOf(item);
    itemHosts.splice(removeIndex, 1);
    const hosts = itemHosts.join(',');

    this.setState({
      changeHosts: hosts.split(','),
    });

    const removeHostsItem = {
      field:
        typeof config[key].Value === 'string'
          ? config[key].Value
          : config[key].Value[index].field,
      hosts: hosts,
    };
    config[key].Value.splice(index, 1, removeHostsItem);
    this.props.handleParamBlur(config[key].Value, `Config.${key}.Value`);
    this.props.setConfig(key, config[key].Value);
  };

  // 关联主机
  linkItem = (key: string, index: number, e: any) => {
    const { config } = this.props;
    this.props.getHostsList();
    this.filterHosts(config[key].Value);
    this.setState({
      hostVisible: true,
      linkIndex: index,
      linkName: key.split('.')[0],
      changeHosts: [],
    });
  };

  // 渲染Icon
  renderIcon(title) {
    const { configEditState } = this.state;
    const { config } = this.props;
    if (config[title].iconType) {
      return (
        <Icon
          type="unlock"
          theme="filled"
          style={{ fontSize: 20, marginLeft: '20px' }}
          onClick={this.closeKey.bind(this, title)}
        />
      );
    } else {
      if (configEditState == 'normal') {
        return (
          <Tooltip placement="top" title="请先开放编辑">
            <Icon
              type="lock"
              theme="filled"
              style={{ fontSize: 20, cursor: 'pointer', marginLeft: '20px' }}
            />
          </Tooltip>
        );
      }
      return (
        <Icon
          type="lock"
          theme="filled"
          style={{ fontSize: 20, cursor: 'pointer', marginLeft: '20px' }}
          onClick={this.saveKey.bind(this, title)}
        />
      );
    }
  }

  // 开锁 => Modal => 记录Key
  saveKey(title) {
    this.setState({
      checkVisible: true,
      nowTitle: title,
    });
  }

  // 关闭锁重新加密
  closeKey(nowTitle) {
    // const { nowTitle } = this.state
    const { config } = this.props;
    const { setFieldsValue } = this.props.form;
    const { password_key } = this.state;
    let encryptWord = config[nowTitle].Value;
    if (encryptWord == '') {
      encryptWord = '45d3161f8aa54588ab3901fa954cc01a';
    }
    // 密钥必须先MD5加密
    const securityKey16 = CryptoJS.MD5(password_key);

    // iv偏移写死 和后端保持一致
    const iv = '1234567890123456';

    // 所有的config
    const closeConfig = {
      securityKey: securityKey16,
      iv: iv,
    };
    const result = this.encrypt(encryptWord, closeConfig);

    this.replacePwd(result, nowTitle);
    setTimeout(() => {
      this.resetIcontype(nowTitle);
    }, 0);
    setFieldsValue({ input_password: '' });
  }

  // 解密密码
  checkPassWord() {
    const { getFieldValue, setFieldsValue } = this.props.form;
    const { nowTitle } = this.state;
    const { config } = this.props;

    const password = getFieldValue('input_password');

    // 密钥必须先MD5加密
    const securityKey16 = CryptoJS.MD5(password);

    // iv偏移写死 和后端保持一致
    const iv = '1234567890123456';

    // 所有的config
    const checkConfig = {
      securityKey: securityKey16,
      iv: iv,
    };
    let result = this.decrypt(config[nowTitle].Value, checkConfig);
    if (result == '45d3161f8aa54588ab3901fa954cc01a') {
      result = '';
      this.replacePwd(result, nowTitle);
      this.setState(
        {
          checkVisible: false,
          iconType: 1,
          password_key: password,
        },
        () => {
          setFieldsValue({ input_password: '' });
          this.updateIconType(nowTitle);
        }
      );
    } else if (!result) {
      message.error('密码输入错误！请重新输入');
    } else {
      // this.replacePasswd(key, result)
      this.replacePwd(result, nowTitle);
      this.setState(
        {
          checkVisible: false,
          iconType: 1,
          password_key: password,
        },
        () => {
          setFieldsValue({ input_password: '' });
          this.updateIconType(nowTitle);
        }
      );
    }
  }

  // 加密
  encrypt(plainText, config) {
    const securityKey = CryptoJS.enc.Utf8.parse(config.securityKey);
    const iv = CryptoJS.enc.Utf8.parse(config.iv);
    const encrypted = CryptoJS.AES.encrypt(plainText, securityKey, {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    });
    return encrypted.toString();
  }

  // 解密
  decrypt(cipherText, config) {
    const securityKey = CryptoJS.enc.Utf8.parse(config.securityKey);
    const iv = CryptoJS.enc.Utf8.parse(config.iv);
    const decrypted = CryptoJS.AES.decrypt(cipherText, securityKey, {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    });
    return decrypted.toString(CryptoJS.enc.Utf8);
  }

  // 保存关联主机配置
  saveHosts() {
    const { changeHosts, linkIndex, linkName } = this.state;
    const { config } = this.props;

    // 判断有相同参数值提醒错误
    let newArr = [...config[linkName].Value].map((item) => item.field);
    const res = new Map();
    newArr = newArr.filter((a) => !res.has(a) && res.set(a, 1));
    console.log(newArr);
    if (config[linkName].Value.length !== newArr.length) {
      this.setState({ isShowErro: true });
    } else {
      this.setState({ isShowErro: false });
    }
    const saveHostsItem = {
      field:
        typeof config[linkName].Value[linkIndex] === 'string'
          ? config[linkName].Value[linkIndex]
          : config[linkName].Value[linkIndex].field,
      hosts:
        changeHosts.length !== 0
          ? config[linkName].Value[linkIndex].hosts.length > 1
            ? config[linkName].Value[linkIndex].hosts +
              ',' +
              changeHosts.join(',')
            : changeHosts.join(',')
          : config[linkName].Value[linkIndex].hosts,
    };

    config[linkName].Value.splice(linkIndex, 1, saveHostsItem);
    this.props.handleParamBlur(
      config[linkName].Value,
      `Config.${linkName}.Value`
    );
    this.props.setConfig(linkName, config[linkName].Value);
    this.cancelLink();
  }

  // 关闭关联主机
  cancelLink() {
    this.setState({
      hostVisible: false,
      linkIndex: null,
      linkName: '',
      changeHosts: [],
    });
  }

  // 运行配置编辑
  runtimeEdit = () => {
    this.setState({
      configEditState: 'edit',
      beforeEditConfig: cloneDeep(this.props.config),
    });
    this.props.actions.editRuntimeState('edit');
  };

  // 运行配置取消
  runtimeCancel = () => {
    this.props.setAllConfig('Config', this.state.beforeEditConfig);
    const editObj = this.contrastAndReturnEditData(
      this.props.config,
      this.state.beforeEditConfig
    );
    this.setState({
      configEditState: 'normal',
    });
    this.props.actions.editRuntimeState('normal');
    if (JSON.stringify(editObj) === '{}') return;
    this.props.saveAllConfig('Config', editObj, this.state.editFiled);
  };

  // 运行配置保存
  runtimeSave = () => {
    const { config, encryptInfo, hostsSelectList } = this.props;
    const { password_key, isShowErro } = this.state;
    let selectHost = [];
    if (isShowErro) {
      message.warning('参数值重复,请重新输入');
      return;
    }
    for (const title in config) {
      if (Array.isArray(config[title].Value)) {
        config[title].Value.forEach((item) => {
          if (item.hosts) {
            selectHost.push(item.hosts);
          }
        });
        if (selectHost.length !== hostsSelectList.length) {
          message.warning('存在主机未关联改参数');
          return;
        }
      }
      if (config[title].iconType) {
        message.warning('请先上锁');
        return;
      }
    }
    this.setState({
      configEditState: 'normal',
    });
    const editObj = this.contrastAndReturnEditData(
      this.state.beforeEditConfig,
      this.props.config
    );
    this.props.actions.editRuntimeState('normal');
    if (JSON.stringify(editObj) === '{}') return;
    // 支持 sm2
    if (password_key && encryptInfo.encrypt_type === 'sm2') {
      // 密钥必须先MD5加密
      const securityKey16 = CryptoJS.MD5(password_key);
      // iv偏移写死 和后端保持一致
      const iv = '1234567890123456';
      // 所有的config
      const checkConfig = {
        securityKey: securityKey16,
        iv: iv,
      };
      for (const key in editObj) {
        if (key.indexOf('password') > -1) {
          // 先支持aes解密
          const defaultPwd = this.decrypt(editObj[key].Default, checkConfig);
          const valuePwd = this.decrypt(editObj[key].Value, checkConfig);
          // 再进行sm2加密
          editObj[key].Default = encryptSM(
            defaultPwd,
            encryptInfo.encrypt_public_key
          );
          editObj[key].Value = encryptSM(
            valuePwd,
            encryptInfo.encrypt_public_key
          );
        }
      }
    }
    this.props.saveAllConfig('Config', editObj, this.state.editFiled, true);
  };

  // 关联主机复选框选择值后
  onChangeHosts = (checkedValues) => {
    this.setState({
      changeHosts: checkedValues,
    });
  };

  // 对比新老数据返回修改的数据
  contrastAndReturnEditData = (oldData, newData) => {
    const obj: any = {};
    for (const filed in newData) {
      const newValue = newData[filed]?.Value;
      const oldValue = oldData[filed]?.Value;
      if (newValue != oldValue) {
        obj[filed] = newData[filed];
      }
    }
    return obj;
  };

  render() {
    const { getFieldDecorator, setFieldsValue } = this.props.form;
    const { checkVisible, hostVisible, changeHosts, configEditState } =
      this.state;
    const {
      hostsSelectList,
      config,
      handleFileChange,
      selectedConfigFile,
      configFileList,
    } = this.props;
    const disabledLink = this.disabledLink;

    return (
      <div className="config-param_COMP">
        <div>
          配置文件：
          <Select
            placeholder="请选择配置文件"
            value={selectedConfigFile}
            style={{ width: '264px', marginRight: '-50px' }}
            onChange={(e) => handleFileChange(e)}>
            {configFileList?.length && <Option value={''}>全部</Option>}
            {(configFileList || []).map((file: string) => (
              <Option key={file}>{file}</Option>
            ))}
          </Select>
        </div>
        {JSON.stringify(config) === '{}' ? null : (
          <div className="handle-btn">
            {configEditState === 'normal' ? (
              <div>
                <Button type="link" onClick={this.runtimeEdit}>
                  编辑
                </Button>
              </div>
            ) : (
              <div>
                <Button type="link" onClick={this.runtimeSave}>
                  保存
                </Button>
                <Button type="link" onClick={this.runtimeCancel}>
                  取消
                </Button>
              </div>
            )}
          </div>
        )}
        <div>{this.renderNode(configEditState).map((o: any) => o)}</div>
        <Modal
          visible={checkVisible}
          title="解锁提醒"
          onOk={() => this.checkPassWord()}
          onCancel={() =>
            this.setState({ checkVisible: false }, () =>
              setFieldsValue({ input_password: '' })
            )
          }
          okText="确认"
          cancelText="取消">
          <form action="" autoComplete="off">
            <FormItem {...formItemCheckLayout} label="解锁密码">
              {getFieldDecorator(
                'input_password',
                {}
              )(<Input type="password" autoComplete="off" />)}
            </FormItem>
          </form>
        </Modal>
        <Modal
          title="关联主机"
          className="linkhosts-modal"
          visible={hostVisible}
          onOk={() => this.saveHosts()}
          onCancel={() => this.cancelLink()}>
          <div className="linkhosts-title">
            <Icon
              type="info-circle"
              theme="filled"
              style={{
                marginRight: '5px',
                color: '#3f87ff',
                verticalAlign: 'middle',
              }}
            />
            <span>已关联的主机需要删除后才能再次关联</span>
          </div>
          <div className="linkhosts-content">
            <Checkbox.Group
              onChange={this.onChangeHosts}
              style={{ width: '100%' }}
              value={changeHosts}>
              <Row>
                {hostsSelectList &&
                  hostsSelectList.map((item, index) => {
                    if (
                      disabledLink &&
                      disabledLink.findIndex((y) => y === item) > -1
                    ) {
                      return (
                        <Col span={8} key={index}>
                          <Checkbox
                            value={item}
                            disabled
                            style={{ marginBottom: '20px' }}>
                            {item}
                          </Checkbox>
                        </Col>
                      );
                    } else {
                      return (
                        <Col span={8} key={index}>
                          <Checkbox
                            value={item}
                            style={{ marginBottom: '20px' }}>
                            {item}
                          </Checkbox>
                        </Col>
                      );
                    }
                  })}
              </Row>
            </Checkbox.Group>
          </div>
        </Modal>
      </div>
    );
  }
}

export default Form.create<Prop>()(ParamCompConfig);
