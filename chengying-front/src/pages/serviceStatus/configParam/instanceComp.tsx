import * as React from 'react';
import * as CryptoJS from 'crypto-js';
import {
  Input,
  Form,
  Modal,
  Tooltip,
  Icon,
  message,
  Switch,
  Progress,
  Typography,
} from 'antd';
import * as ServiceListActions from '@/actions/serviceAction';
import { servicePageService } from '@/services';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { ServiceStore } from '@/stores/modals';
import { formItemLayout, formItemCheckLayout } from '../constant';
import { FormComponentProps } from 'antd/lib/form';

const confirm = Modal.confirm;
const FormItem = Form.Item;
const SWITCH_TYPE = {
  download: '下载',
}; // 开关类型映射
const PROCESS_INTERVAL = 5000;
const STATUS_MAP = {
  RUNNING: 'active',
  FAIL: 'exception',
  SUCCESS: 'success',
};

const mapStateToProps = (state: AppStoreTypes) => ({
  ServiceStore: state.ServiceStore,
});
const mapDispatchToProps = (dispatch: any) => ({
  actions: bindActionCreators(Object.assign({}, ServiceListActions), dispatch),
});

interface Prop extends FormComponentProps {
  ServiceStore?: ServiceStore;
  updateServiceConfig?: any;
  actions?: any;
  pname: string;
  sname: string;
  pid: number;
  pversion: string;
}

interface State {
  default_input_value: string;
  has_config_value: boolean;
  checkVisible: boolean;
  key: any;
  statusVisible: boolean;
  statusDetail: any;
  currentOperate: any;
  currentOperateSwitch: any;
  service_name: any;
  product_name: any;
  product_version: any;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
class InstanceConfigComp extends React.Component<Prop, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      default_input_value: '',
      has_config_value: false,
      checkVisible: false,
      key: {},
      statusVisible: false,
      statusDetail: {},
      currentOperate: '',
      currentOperateSwitch: {},
      service_name: '',
      product_name: '',
      product_version: '',
    };
  }

  private processInterval: any = null;

  componentWillUnmount() {
    this.processInterval && clearInterval(this.processInterval);
  }

  // 开关的相关文件下载
  handleSwitchExtension(data: any) {
    const { sname, pname } = this.props;
    window.location.href = `/api/v2/product/${pname}/service/${sname}/extention_operation?type=${data.Type}&value=${data.Value}`;
  }

  // 开关操作
  handleSwitch(params: any, name: string) {
    const { sname, pname } = this.props;
    params.name = name;
    this.setState({
      currentOperateSwitch: name,
    });
    servicePageService
      .getSwitchRecord({
        product_name: pname,
        service_name: sname,
        name,
      })
      .then((res) => {
        const { data } = res;
        if (data.code === 0) {
          if (data && data.data && data.data.record_id) {
            // 进度中，弹出进度弹窗
            this.setState({
              statusVisible: true,
            });
            this.getSwitchDetail(data.data.record_id);
          } else {
            params.name = name;
            params.sname = sname;
            params.pname = pname;
            this.handleSwitchConfirm(params);
          }
        }
      });
  }

  // 开关操作的二次确认弹窗
  handleSwitchConfirm(params: any) {
    const { IsOn } = params;
    const _this = this;
    confirm({
      title: `确定${IsOn ? '关闭' : '开启'}Nginx吗？`,
      icon: <Icon type="exclamation-circle" theme="filled" />,
      content: `Nginx${IsOn ? '关闭' : '开启'}后，默认重启服务，确定${
        IsOn ? '关闭' : '开启'
      }吗？`,
      onOk() {
        _this.operateSwitch(params);
      },
    });
  }

  // 操作开关开启关闭
  operateSwitch(params: any) {
    const { pversion, pname, sname } = this.props;
    servicePageService
      .operateSwitch({
        product_name: pname,
        service_name: sname,
        name: params.name,
        type: params.IsOn ? 'off' : 'on',
        product_version: pversion,
      })
      .then((res) => {
        const { data } = res;
        if (data.code === 0) {
          if (data.data && data.data.record_id) {
            this.setState({
              statusVisible: true,
              statusDetail: {},
            });
            this.getSwitchDetail(data.data.record_id);
          }
        }
      });
  }

  // 根据开关id获取开关操作详情
  getSwitchDetail(record_id: any) {
    const { sname, pname } = this.props;
    servicePageService
      .getSwitchDetail({
        product_name: pname,
        service_name: sname,
        record_id,
      })
      .then((res) => {
        const { data } = res;
        if (data.code === 0) {
          if (data.data && data.data.status === 'RUNNING') {
            // 正在操作中，轮询
            this.setState({
              statusDetail: data.data,
            });
            if (this.processInterval) {
            } else {
              this.processInterval = setInterval(
                () => this.getSwitchDetail(record_id),
                PROCESS_INTERVAL
              );
            }
          } else {
            // 操作结束，关闭轮询
            this.setState({
              statusDetail: data.data,
            });
            this.processInterval && clearInterval(this.processInterval);
            this.processInterval = null;
            if (data.data.status === 'SUCCESS') {
              this.props.updateServiceConfig();
            }
          }
        }
      });
  }

  /**
   * 构建实例配置项
   * @param {Obj} data
   *
   */
  buildInstanceEditByNode(data: any, level: any, initKey?: any): any {
    const ctx = this;
    const nextLevel = level + 1;
    const instanceDataList = Object.keys(data);
    const envName = 'environment';
    const envIdx = instanceDataList.findIndex((item) => item === envName);
    if (envIdx !== -1) {
      instanceDataList.splice(envIdx, 1);
      instanceDataList.push(envName);
    }
    return instanceDataList.map((keyName, keyIndex) => {
      const key = initKey ? initKey + '.' + keyName : keyName;
      // 如果是数组则不显示title
      if (typeof data[keyName] === 'object') {
        if (keyName == 'Switch') {
          for (const key in data[keyName]) {
            return (
              <FormItem label={key} {...formItemLayout}>
                <FormItem
                  key={keyName + '-' + key}
                  style={{ display: 'inline-block' }}>
                  <Switch
                    checked={data[keyName].Nginx.IsOn}
                    onClick={ctx.handleSwitch.bind(
                      ctx,
                      data[keyName].Nginx,
                      'Nginx'
                    )}></Switch>
                </FormItem>
                {data[keyName].Nginx.IsOn && (
                  <a
                    style={{ marginLeft: '5px' }}
                    onClick={ctx.handleSwitchExtension.bind(
                      ctx,
                      data[keyName].Nginx.Extention
                    )}>
                    {SWITCH_TYPE[data[keyName].Nginx.Extention?.Type]}{' '}
                    {data[keyName].Nginx.Extention?.Value}
                  </a>
                )}
              </FormItem>
            );
          }
        } else {
          return ctx.buildInstanceEditByNode(data[keyName], nextLevel, key);
        }
      } else if (keyName === 'MaxReplica' || keyName === 'StartAfterInstall') {
        return (
          <FormItem {...formItemLayout} label={key} key={keyName + '-' + level}>
            <span>{data[keyName]}</span>
          </FormItem>
        );
      } else if (keyName !== 'UseCloud') {
        const addon = (
          <i className="emicon emicon-undo" style={{ fontSize: 14 }} />
        );
        let label =
          key.length > 20 ? (
            <Tooltip title={key}>{key.slice(0, 20)}...</Tooltip>
          ) : (
            key
          );
        // 如果是环境变量
        if (initKey === envName) {
          label = (
            <span>
              <i
                className="emicon emicon-env mr-8"
                style={{
                  position: 'absolute',
                  left: '-52px',
                  fontSize: 44,
                  color: '#00a755',
                }}
              />
              {label}
            </span>
          );
        }
        const { Paragraph } = Typography;
        const input = (
          <Tooltip
            trigger={data[keyName].length === 0 ? 'focus' : 'hover'}
            arrowPointAtCenter={true}
            title={() => (
              <Paragraph copyable style={{ color: '#fff' }}>
                {data[keyName]}
              </Paragraph>
            )}>
            <div>
              <Input
                defaultValue={data[keyName]}
                style={{ width: 400 }}
                value={
                  data[keyName].length > 40
                    ? data[keyName].slice(0, 40) + '...'
                    : data[keyName]
                }
                className="c-paramConfig__ant-input"
                disabled={true}
                addonAfter={addon}
                onChange={ctx.handleInputChange.bind(ctx, key)}
                onFocus={ctx.handleInputFocus.bind(ctx, key)}
              />
            </div>
          </Tooltip>
        );
        return (
          <FormItem
            {...formItemLayout}
            label={label}
            key={keyName + '-' + level}>
            {input}
            {/* {
                            // 当前版本暂时不做部署配置相关
                            /password/.test(key.toLowerCase()) ? ctx.renderIcon(data[keyName]) : null
                        } */}
          </FormItem>
        );
      }
    });
  }

  renderIcon(nowKey) {
    const { cur_service } = this.props.ServiceStore;
    const iconType = cur_service.Config[nowKey.config].iconType;
    const ctx = this;
    if (iconType) {
      return (
        <Icon
          type="unlock"
          theme="filled"
          style={{ fontSize: 20, marginLeft: '20px' }}
          onClick={ctx.closeKey.bind(ctx, nowKey.config)}
        />
      );
    } else {
      return (
        <Icon
          type="lock"
          theme="filled"
          style={{ fontSize: 20, cursor: 'pointer', marginLeft: '20px' }}
          onClick={ctx.saveKey.bind(ctx, nowKey.config)}
        />
      );
    }
  }

  // 开锁 => Modal => 记录Key
  saveKey(newKey) {
    this.setState({
      checkVisible: true,
      key: newKey,
    });
  }

  closeKey(newKey) {
    const { cur_service } = this.props.ServiceStore;
    const { getFieldValue } = this.props.form;
    const password = getFieldValue('input_password');
    // 密钥必须先MD5加密
    const securityKey16 = CryptoJS.MD5(password);

    // iv偏移写死 和后端保持一致
    const iv = '1234567890123456';

    // 所有的config
    const closeConfig = {
      securityKey: securityKey16,
      iv: iv,
    };
    const result = this.encrypt(
      cur_service.Config[newKey].current,
      closeConfig
    );

    // this.setState({ key: newKey })
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.Config);
    config[newKey].current = result;
    config[newKey].iconType = 0;
    cur_service.Config = config;
    setServiceConfigModify(cur_service);
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

  /**
   * 统一处理input的change事件
   * @param {Event} e
   */
  handleInputChange(key: any, e: any) {
    const { cur_service } = this.props.ServiceStore;
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.configModify);
    config['Instance.' + key] = e.target.value;
    if (key.indexOf('.') > -1) {
      cur_service.Instance[key.split('.')[0]][key.split('.')[1]] =
        e.target.value;
    } else {
      cur_service.Instance[key] = e.target.value;
    }
    cur_service.configModify = config;
    setServiceConfigModify(cur_service);
  }

  handleInputFocus(key: any, e: any) {
    this.setState({
      default_input_value: e.target.value,
    });
  }

  // 单个配置字段恢复默认值
  handleReloadConfig(key: any, e: any) {
    Modal.confirm({
      title: '确认恢复此字段的默认值吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        const { sname, pname, pid, pversion } = this.props;
        this.props.actions.resetServiceConfig({
          product_name: pname,
          service_name: sname,
          pid: pid,
          product_version: pversion,
          field_path: 'Instance.' + key,
        });
      },
    });
  }

  checkPassWord() {
    const { getFieldValue } = this.props.form;
    const { key } = this.state;

    const { cur_service } = this.props.ServiceStore;

    const password = getFieldValue('input_password');

    // 密钥必须先MD5加密
    const securityKey16 = CryptoJS.MD5(password);

    // iv偏移写死 和后端保持一致
    const iv = '1234567890123456';

    // 所有的config
    const config = {
      securityKey: securityKey16,
      iv: iv,
    };

    const result = this.decrypt(cur_service.Config[key].current, config);
    if (
      !result &&
      cur_service.Config[key].current != 'Qve2w+9noccWGql7x32JQg=='
    ) {
      message.error('密码输入错误！请重新输入');
    } else {
      this.replacePasswd(key, result);
      this.setState({
        checkVisible: false,
      });
    }
  }

  replacePasswd(key: any, passwd: any) {
    const { cur_service } = this.props.ServiceStore;
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.Config);
    config[key].current = passwd;
    config[key].iconType = 1;
    cur_service.Config = config;
    setServiceConfigModify(cur_service);
  }

  render() {
    const { cur_service } = this.props.ServiceStore;
    const Instance = cur_service.Instance;
    const { checkVisible, statusVisible, statusDetail, currentOperateSwitch } =
      this.state;
    const { getFieldDecorator } = this.props.form;

    if (Instance) {
      return (
        <div
          className="edit-panel config-box"
          style={
            Instance.UseCloud ? { display: 'none' } : { display: 'block' }
          }>
          {JSON.stringify(Instance) !== '{}' ? (
            <Form>{this.buildInstanceEditByNode(Instance, 1)}</Form>
          ) : (
            <p style={{ lineHeight: '20px' }}>无</p>
          )}
          <Modal
            title={`${
              statusDetail.switch_type === 'off' ? '关闭' : '开启'
            }Nginx`}
            visible={statusVisible}
            onOk={() => this.setState({ statusVisible: false })}
            onCancel={() => this.setState({ statusVisible: false })}
            okText="确认"
            cancelText="取消"
            bodyStyle={{ padding: 0 }}>
            {statusDetail.status == 'RUNNING' && (
              <div className="config-tips">
                <Icon
                  style={{ color: '#3F87FF' }}
                  type="exclamation-circle"
                  theme="filled"
                />
                <span style={{ marginLeft: '8px' }}>
                  关闭窗口后，进行中的任务将在后台运行
                </span>
              </div>
            )}
            {statusDetail.status !== 'SUCCESS' && (
              <div
                className="config-progress"
                style={{
                  width: statusDetail.status == 'FAIL' ? '350px' : '325px',
                }}>
                <span style={{ marginRight: '8px' }}>
                  {statusDetail.status == 'FAIL'
                    ? `${
                        statusDetail.switch_type === 'off' ? '关闭' : '开启'
                      }失败`
                    : `${
                        statusDetail.switch_type === 'off' ? '关闭' : '开启'
                      }中`}
                </span>
                <Progress
                  style={{ width: '267px' }}
                  percent={statusDetail.progress}
                  status={STATUS_MAP[statusDetail.status]}
                />
                {statusDetail.status == 'FAIL' && (
                  <span
                    style={{ color: '#3F87FF', cursor: 'pointer' }}
                    onClick={() => this.operateSwitch(currentOperateSwitch)}>
                    重试
                  </span>
                )}
              </div>
            )}
            {statusDetail.status !== 'SUCCESS' && (
              <div
                style={{
                  lineHeight: '32px',
                  height: '86px',
                  overflow: 'scroll',
                  whiteSpace: 'pre-wrap',
                  padding: '0 83px',
                  marginBottom: '28px',
                  overflowX: 'hidden',
                }}
                // dangerouslySetInnerHTML={{__html:statusDetail.status_message}}
              >
                {statusDetail.status_message}
              </div>
            )}
            {statusDetail.status === 'SUCCESS' && (
              <div style={{ padding: '24px 177px 28px', textAlign: 'center' }}>
                <img
                  src={require('public/imgs/kerberosSuccess@2x.png')}
                  style={{
                    display: 'block',
                    margin: '0 auto',
                    width: '80px',
                    height: '80px',
                  }}
                />
                <div
                  style={{
                    marginTop: '24px',
                    lineHeight: '28px',
                    fontSize: '20px',
                  }}>
                  Nginx{statusDetail.switch_type === 'off' ? '关闭' : '开启'}
                  成功
                </div>
              </div>
            )}
          </Modal>
          <Modal
            visible={checkVisible}
            title="解锁提醒"
            onOk={() => this.checkPassWord()}
            onCancel={() => this.setState({ checkVisible: false })}
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
        </div>
      );
    } else {
      return (
        <div className="edit-panel config-box" style={{ display: 'none' }}>
          <p style={{ lineHeight: '20px' }}>无</p>
        </div>
      );
    }
  }
}
export default Form.create<Prop>()(InstanceConfigComp);
