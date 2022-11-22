import * as React from 'react';
import * as CryptoJS from 'crypto-js';
import {
  Input,
  Form,
  Tooltip,
  Modal,
  Icon,
  message,
  Tag,
  Checkbox,
  Row,
  Col,
} from 'antd';
import * as ServiceListActions from '@/actions/serviceAction';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { ServiceStore } from '@/stores/modals';
import { formItemLayout, formItemCheckLayout } from '../constant';
import moment from 'moment';
import { FormComponentProps } from 'antd/lib/form';
import * as Http from '@/utils/http';
import { encryptSM } from '@/utils/password';

const FormItem = Form.Item;
const mapStateToProps = (state: AppStoreTypes) => ({
  ServiceStore: state.ServiceStore,
});
const mapDispatchToProps = (dispatch: any) => ({
  actions: bindActionCreators(Object.assign({}, ServiceListActions), dispatch),
});

interface Prop extends FormComponentProps {
  ServiceStore?: ServiceStore;
  actions?: any;
  query: string;
  pname: string;
  sname: string;
  noHosts: string;
  repeatParams: any[];
  changeHosts: Function;
  isKubernetes: boolean;
  HeaderStore: any;
  pid: number;
  pversion: string;
  noEditAuthority: boolean;
  encryptInfo: {
    encrypt_public_key: string;
    encrypt_type: string;
  };
  canEditPwd: boolean;
  canEditHost: boolean;
  canReset: boolean;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
class ConfigComp extends React.Component<Prop, any> {
  disabledLink = [];
  constructor(props: any) {
    super(props);
    this.state = {
      default_input_value: '',
      has_config_value: false,
      checkVisible: false,
      iconType: 0,
      passwd: '',
      key: {},
      password_key: '',
      hostVisible: false,
      changeHosts: [],
      hostsList: [],
      linkIndex: null,
      linkName: '',
    };
  }

  componentDidMount() {
    this.getHostsList();
  }

  // 获取当前服务下配置的所有主机列表
  getHostsList = () => {
    const { sname, pname, HeaderStore } = this.props;
    if (!sname || !pname) {
      return;
    }
    Http.get(
      `/api/v2/product/${pname}/service/${sname}/selected_hosts?clusterId=${HeaderStore.cur_parent_cluster?.id}`,
      {}
    ).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        this.setState({
          hostsList: data.data.hosts,
        });
      } else {
        message.error(data.msg);
      }
    });
  };

  // 构建Config配置项，只对下面value是string类型的节点编辑
  buildConfigEditByNode(data: any, initKey?: any) {
    const {
      query,
      noEditAuthority,
      ServiceStore,
      noHosts,
      isKubernetes,
      repeatParams,
    } = this.props;
    const couldSaveConfig =
      ServiceStore.cur_service &&
      ServiceStore.cur_service.configModify &&
      Object.keys(ServiceStore.cur_service.configModify).length;
    const ctx = this;
    let hasvalue = false;
    const configNodes = Object.keys(data).map((keyName, keyIndex) => {
      const key = initKey ? initKey + '.' + keyName : keyName;
      return Object.keys(data[keyName]).map((keyName2, keyIndex2) => {
        const key2 = key + '.' + keyName2;
        if (
          keyName2 === 'current' &&
          typeof data[keyName][keyName2] === 'string'
        ) {
          hasvalue = true;
          return key.indexOf(query) > -1 ? (
            <FormItem
              {...formItemLayout}
              labelAlign="right"
              label={
                <span className="labelBox">
                  {data[keyName].nouse && (
                    <i className="emicon emicon-wyy icon-style" />
                  )}
                  {data[keyName].isnew && (
                    <i className="emicon emicon-new icon-style" />
                  )}
                  {key.length > 20 ? (
                    <Tooltip title={key}>
                      <span>{key.slice(0, 20)}...</span>
                    </Tooltip>
                  ) : (
                    key
                  )}
                </span>
              }
              key={key}>
              {/password/.test(key.toLowerCase()) ? (
                <span style={{ display: 'inline-flex' }}>
                  <Input
                    type={
                      !data[keyName].iconType
                        ? /password/.test(key.toLowerCase())
                          ? 'password'
                          : 'text'
                        : 'text'
                    }
                    // disabled={noEditAuthority || !data[keyName].iconType}
                    disabled={
                      /password/.test(key.toLowerCase())
                        ? !data[keyName].iconType
                        : noEditAuthority
                    }
                    className={
                      noEditAuthority ? 'c-paramConfig__ant-input' : ''
                    }
                    defaultValue={data[keyName][keyName2]}
                    style={{ width: 400 }}
                    value={data[keyName][keyName2]}
                    addonAfter={
                      <Tooltip title="恢复默认值">
                        <i
                          className="emicon emicon-undo"
                          style={{ fontSize: 14 }}
                          onClick={
                            !couldSaveConfig
                              ? /password/.test(key.toLowerCase())
                                ? data[keyName].iconType
                                  ? noEditAuthority
                                    ? null
                                    : ctx.handleReloadConfig.bind(ctx, key2, 1)
                                  : null
                                : ctx.handleReloadConfig.bind(ctx, key2, 1)
                              : null
                          }
                        />
                      </Tooltip>
                    }
                    onChange={ctx.handleInputChange.bind(ctx, key2)}
                    onFocus={ctx.handleInputFocus.bind(ctx, key2)}
                  />
                </span>
              ) : (
                <span style={{ display: 'inline-flex' }}>
                  <Input.TextArea
                    disabled={
                      /password/.test(key.toLowerCase())
                        ? !data[keyName].iconType
                        : noEditAuthority
                    }
                    className={
                      noEditAuthority ? 'c-paramConfig__ant-input' : ''
                    }
                    defaultValue={data[keyName][keyName2]}
                    value={data[keyName][keyName2]}
                    style={{
                      width: 363,
                      height: 32,
                      outline: 'none',
                      resize: 'none',
                    }}
                    onChange={ctx.handleInputChange.bind(ctx, key2)}
                    onFocus={ctx.handleInputFocus.bind(ctx, key2)}
                  />
                  {
                    <span className="afteron">
                      <Tooltip title="恢复默认值">
                        <i
                          className="emicon emicon-undo"
                          style={{ fontSize: 14 }}
                          onClick={
                            !couldSaveConfig
                              ? /password/.test(key.toLowerCase())
                                ? data[keyName].iconType
                                  ? noEditAuthority
                                    ? null
                                    : ctx.handleReloadConfig.bind(ctx, key2, 1)
                                  : null
                                : ctx.handleReloadConfig.bind(ctx, key2, 1)
                              : null
                          }
                        />
                      </Tooltip>
                    </span>
                  }
                </span>
              )}
              {data[keyName].updated && (
                <Tooltip
                  title={
                    <React.Fragment>
                      <p>
                        最近修改时间：
                        {moment(data[keyName].updated)
                          .utc()
                          .zone(+6)
                          .format('YYYY-MM-DD HH:mm:ss')}
                      </p>
                      <p>默认值：{data[keyName].default}</p>
                    </React.Fragment>
                  }>
                  <i className="emicon emicon-revise icon-style"></i>
                </Tooltip>
              )}
              {!/password/.test(key.toLowerCase()) && !isKubernetes && (
                <span className="operate-icon">
                  <Icon
                    className="plus-icon"
                    type="plus"
                    onClick={ctx.addItem.bind(ctx, key2)}
                  />
                </span>
              )}
              {/password/.test(key.toLowerCase())
                ? ctx.renderIcon(data[keyName])
                : null}
            </FormItem>
          ) : (
            ''
          );
        } else {
          if (keyName2 === 'current' && Array.isArray(data[keyName].current)) {
            hasvalue = true;
            return key.indexOf(query) > -1 ? (
              <FormItem
                {...formItemLayout}
                label={
                  <React.Fragment>
                    {data[keyName].nouse && (
                      <i className="emicon emicon-wyy icon-style" />
                    )}
                    {data[keyName].isnew && (
                      <i className="emicon emicon-new icon-style" />
                    )}
                    {key.length > 20 ? (
                      <Tooltip title={key}>
                        <span>{key.slice(0, 20)}...</span>
                      </Tooltip>
                    ) : (
                      key
                    )}
                  </React.Fragment>
                }
                key={key}
                className="config-array">
                {data[keyName].current.map((item, index) => {
                  const hasRepeat = repeatParams.some((x) => x === keyName);
                  return (
                    <FormItem key={index}>
                      {data[keyName].current.length > 1 ? (
                        <Input
                          type={
                            !data[keyName].iconType
                              ? /password/.test(key.toLowerCase())
                                ? 'password'
                                : 'text'
                              : 'text'
                          }
                          // disabled={noEditAuthority || !data[keyName].iconType}
                          disabled={
                            /password/.test(key.toLowerCase())
                              ? !data[keyName].iconType
                              : noEditAuthority
                          }
                          className={
                            noEditAuthority ? 'c-paramConfig__ant-input' : ''
                          }
                          defaultValue={
                            typeof item === 'string' ? item : item.field
                          }
                          style={{ width: 400 }}
                          value={typeof item === 'string' ? item : item.field}
                          addonAfter={
                            <Tooltip title="恢复默认值">
                              <i
                                className="emicon emicon-undo"
                                style={{ fontSize: 14 }}
                                onClick={
                                  !couldSaveConfig
                                    ? /password/.test(key.toLowerCase())
                                      ? data[keyName].iconType
                                        ? noEditAuthority
                                          ? null
                                          : ctx.handleReloadConfig.bind(
                                              ctx,
                                              key2,
                                              2,
                                              item.hosts
                                            )
                                        : null
                                      : ctx.handleReloadConfig.bind(
                                          ctx,
                                          key2,
                                          2,
                                          item.hosts
                                        )
                                    : null
                                }
                              />
                            </Tooltip>
                          }
                          onChange={ctx.handleInputArrayChange.bind(
                            ctx,
                            key2,
                            index
                          )}
                          onFocus={ctx.handleInputFocus.bind(ctx, key2)}
                        />
                      ) : (
                        <Input
                          type={
                            !data[keyName].iconType
                              ? /password/.test(key.toLowerCase())
                                ? 'password'
                                : 'text'
                              : 'text'
                          }
                          // disabled={noEditAuthority || !data[keyName].iconType}
                          disabled={
                            /password/.test(key.toLowerCase())
                              ? !data[keyName].iconType
                              : noEditAuthority
                          }
                          className={
                            noEditAuthority ? 'c-paramConfig__ant-input' : ''
                          }
                          defaultValue={
                            typeof item === 'string' ? item : item.field
                          }
                          style={{ width: 400 }}
                          value={typeof item === 'string' ? item : item.field}
                          addonAfter={
                            <Tooltip title="恢复默认值">
                              <i
                                className="emicon emicon-undo"
                                style={{ fontSize: 14 }}
                                onClick={
                                  !couldSaveConfig
                                    ? /password/.test(key.toLowerCase())
                                      ? data[keyName].iconType
                                        ? noEditAuthority
                                          ? null
                                          : ctx.handleReloadConfig.bind(
                                              ctx,
                                              key2,
                                              2,
                                              item.hosts
                                            )
                                        : null
                                      : ctx.handleReloadConfig.bind(
                                          ctx,
                                          key2,
                                          2,
                                          item.hosts
                                        )
                                    : null
                                }
                              />
                            </Tooltip>
                          }
                          onChange={ctx.handleInputChange.bind(ctx, key2)}
                          onFocus={ctx.handleInputFocus.bind(ctx, key2)}
                        />
                      )}
                      {!/password/.test(key.toLowerCase()) && (
                        <span className="operate-icon">
                          <Icon
                            className="plus-icon"
                            type="plus"
                            onClick={ctx.addItem.bind(ctx, key2)}
                          />
                          {data[keyName][keyName2].length > 1 && (
                            <Icon
                              type="link"
                              className="link-icon"
                              onClick={ctx.linkItem.bind(ctx, key2, index)}
                            />
                          )}
                          {index > 0 ? (
                            <Icon
                              className="dynamic-delete-button"
                              type="delete"
                              onClick={ctx.removeItem.bind(ctx, key2, index)}
                            />
                          ) : null}
                        </span>
                      )}
                      {/password/.test(key.toLowerCase())
                        ? ctx.renderIcon(data[keyName])
                        : null}
                      {data[keyName].updated && (
                        <Tooltip
                          title={
                            <React.Fragment>
                              <p>
                                最近修改时间：
                                {moment(data[keyName].updated)
                                  .utc()
                                  .zone(+6)
                                  .format('YYYY-MM-DD HH:mm:ss')}
                              </p>
                              <p>默认值：{data[keyName].default}</p>
                            </React.Fragment>
                          }>
                          <i className="emicon emicon-revise icon-style"></i>
                        </Tooltip>
                      )}
                      {repeatParams && index === 0 && hasRepeat && (
                        <p className="repeatparams-span">
                          参数值重复,请重新输入
                        </p>
                      )}

                      {keyName === noHosts && index === 0 && (
                        <span className="nohosts-span">
                          存在主机未关联该参数
                        </span>
                      )}
                      <Row
                        style={
                          data[keyName].current.length > 1 && item.hosts
                            ? { marginTop: '12px' }
                            : {}
                        }>
                        {data[keyName].current.length > 1 && item.hosts
                          ? item.hosts.split(',').map((e, e1) => {
                              return (
                                <div key={e1}>
                                  <div
                                    key={e1}
                                    style={{
                                      marginRight: '4px',
                                      height: '28px',
                                      display: 'inline-block',
                                      float: 'left',
                                    }}>
                                    <Tag
                                      visible
                                      closable
                                      onClose={ctx.onCloseHost.bind(
                                        ctx,
                                        key2,
                                        e,
                                        index
                                      )}>
                                      {e}
                                    </Tag>
                                  </div>

                                  {/* 使用react官方文档方法可以解析换行且不影响删除tag操作，但是这个官方文档解释类似于innerHtml操作，存在xss攻击风险 */}
                                  {e1 % 3 === 2 &&
                                    e1 !== item.hosts.split(',').length - 1 && (
                                      <div
                                        dangerouslySetInnerHTML={{
                                          __html: '&nbsp;',
                                        }}
                                      />
                                    )}
                                </div>
                              );
                            })
                          : null}
                      </Row>
                    </FormItem>
                  );
                })}
              </FormItem>
            ) : (
              ''
            );
          } else {
            return '';
          }
        }
      });
    });
    return hasvalue ? (
      configNodes
    ) : (
      <span style={{ margin: '1px 1%' }}>无</span>
    );
  }

  // 叉掉关联的tag标签
  onCloseHost = (key: any, item: any, index: any, e: any) => {
    const { cur_service } = this.props.ServiceStore;
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.configModify);
    const spliceIndex = key.split('.')[0];
    const changeHosts =
      cur_service.Config[spliceIndex].current[index].hosts.split(',');

    const removeIndex = changeHosts.indexOf(item);
    changeHosts.splice(removeIndex, 1);
    const hosts = changeHosts.join(',');

    this.setState({
      changeHosts: hosts.split(','),
    });
    const removeHostsItem = {
      hosts: hosts,
      field:
        typeof cur_service.Config[spliceIndex].current[index] === 'string'
          ? cur_service.Config[spliceIndex].current[index]
          : cur_service.Config[spliceIndex].current[index].field,
    };

    cur_service.Config[spliceIndex].current.splice(index, 1, removeHostsItem);
    config['Config.' + key] = cur_service.Config[spliceIndex].current;

    cur_service.configModify = config;
    setServiceConfigModify(cur_service);
  };

  // 关联主机一行
  linkItem = (key: any, index: any, e: any) => {
    this.getHostsList();
    const { cur_service } = this.props.ServiceStore;
    this.filterHosts(cur_service.Config[key.split('.')[0]].current);
    this.setState({
      hostVisible: true,
      linkIndex: index,
      linkName: key.split('.')[0],
      changeHosts: [],
    });
  };

  // 减少一行
  removeItem = (key: any, index: any, e: any) => {
    const { cur_service } = this.props.ServiceStore;
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.configModify);
    const spliceIndex = key.split('.')[0];

    if (cur_service.Config[spliceIndex].current.length > 2) {
      cur_service.Config[spliceIndex].current.splice(index, 1);
      config['Config.' + key] = cur_service.Config[spliceIndex].current;
    } else {
      if (cur_service.Config[spliceIndex].current.length == 2) {
        this.props.changeHosts();
      }
      cur_service.Config[spliceIndex].current =
        cur_service.Config[spliceIndex].current[0].field;
      config['Config.' + key] = cur_service.Config[spliceIndex].current;
    }
    // console.log(config['Config.' + key], 'config------remove', config)

    // console.log(config, '删除后的config-----', cur_service.Config[spliceIndex].current)
    cur_service.configModify = config;
    setServiceConfigModify(cur_service);
  };

  // 增加一行
  addItem = (key: any, e: any) => {
    const { canEditHost } = this.props;
    if (!canEditHost) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    const { cur_service } = this.props.ServiceStore;
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.configModify);
    const arr = [];
    const splitIndex = key.split('.')[0];
    const pushItem = {
      hosts: '',
      field: Array.isArray(cur_service.Config[splitIndex].current)
        ? cur_service.Config[splitIndex].current[0].field
        : cur_service.Config[splitIndex].current,
    };
    if (key.indexOf('.') > -1) {
      if (Array.isArray(cur_service.Config[splitIndex].current)) {
        cur_service.Config[splitIndex].current.push(pushItem);
      } else {
        const addItem = {
          hosts: '',
          field: cur_service.Config[splitIndex].current,
        };
        arr.push(addItem, pushItem);
        cur_service.Config[splitIndex].current = arr;
      }
    } else {
      cur_service.Config[splitIndex].current = Array.isArray(
        cur_service.Config[splitIndex].current
      )
        ? cur_service.Config[splitIndex].current.push(pushItem)
        : arr.push(cur_service.Config[splitIndex].current, pushItem);
      // cur_service.Config[key] = e.target.value;
    }
    cur_service.configModify = config;
    setServiceConfigModify(cur_service);
  };

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
    const { canEditPwd } = this.props;
    if (!canEditPwd) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.setState({
      checkVisible: true,
      key: newKey,
    });
  }

  // 统一处理input的change事件 @param {Event} e
  handleInputChange(key: any, e: any) {
    const { ServiceStore, actions, encryptInfo } = this.props;
    const { password_key } = this.state;
    const { cur_service } = ServiceStore;
    const { setServiceConfigModify } = actions;
    const config = Object.assign({}, cur_service.configModify);

    // config['Config.' + key] = e.target.value;
    if (/password/.test(key.toLowerCase())) {
      if (encryptInfo.encrypt_type !== 'sm2') {
        const securityKey16 = CryptoJS.MD5(password_key);
        const iv = '1234567890123456';
        const closeConfig = {
          securityKey: securityKey16,
          iv: iv,
        };
        config['Config.' + key] = this.encrypt(e.target.value, closeConfig);
      } else {
        config['Config.' + key] = encryptSM(
          e.target.value,
          encryptInfo.encrypt_public_key
        );
      }
    } else {
      config['Config.' + key] = e.target.value;
    }

    if (key.indexOf('.') > -1) {
      cur_service.Config[key.split('.')[0]][key.split('.')[1]] = e.target.value;
    } else {
      cur_service.Config[key] = e.target.value;
    }

    // console.log(config, '单配置输入框config-----')
    cur_service.configModify = config;
    setServiceConfigModify(cur_service);
  }

  // 数组情况下输入框的值改变
  handleInputArrayChange(key: any, index: any, e: any) {
    const { cur_service } = this.props.ServiceStore;
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.configModify);

    const spliceIndex = key.split('.')[0];

    const spliceItem = {
      hosts: cur_service.Config[spliceIndex].current[index].hosts
        ? cur_service.Config[spliceIndex].current[index].hosts
        : '',
      field: e.target.value,
    };

    cur_service.Config[spliceIndex].current.splice(index, 1, spliceItem);
    config['Config.' + key] = cur_service.Config[spliceIndex].current;
    cur_service.configModify = config;

    // console.log(config['Config.' + key], '数组的改了之后config-----')
    setServiceConfigModify(cur_service);
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

  // 关闭锁重新加密
  closeKey(newKey) {
    const { form, encryptInfo, ServiceStore, actions } = this.props;
    const { cur_service } = ServiceStore;
    const { setFieldsValue } = form;
    const { setServiceConfigModify } = actions;
    const { password_key } = this.state;
    let encryptWord = cur_service.Config[newKey].current;
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

    // this.setState({key: newKey })
    const config = Object.assign({}, cur_service.Config);
    config[newKey].current = result;
    config[newKey].iconType = 0;
    cur_service.Config = config;
    setServiceConfigModify(cur_service);
    setFieldsValue({ input_password: '' });
  }

  handleInputFocus(key: any, e: any) {
    this.setState({
      default_input_value: e.target.value,
    });
  }

  // 单个配置字段恢复默认值
  handleReloadConfig(key: any, type: any, hosts: any) {
    const { canReset } = this.props;
    if (!canReset) {
      message.error('权限不足，请联系管理员！');
      return false;
    }
    // type: 1是绑定单台主机  2是绑定多台配置
    Modal.confirm({
      title: '确认恢复此字段的默认值吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        const { sname, pname, pid, pversion } = this.props;
        if (type == 1) {
          this.props.actions.resetServiceConfig({
            product_name: pname,
            service_name: sname,
            pid: pid,
            product_version: pversion,
            field_path: 'Config.' + key,
          });
        } else {
          this.props.actions.resetMultiServiceConfig({
            product_name: pname,
            service_name: sname,
            pid: pid,
            product_version: pversion,
            field_path: 'Config.' + key,
            hosts: hosts,
          });
        }
      },
    });
  }

  // 恢复默认设置
  handleServiceConfigReset(params: any) {
    const { sname, pname, pid, pversion } = this.props;
    this.props.actions.resetServiceConfig({
      product_name: pname,
      service_name: sname,
      pid: pid,
      product_version: pversion,
      field_path: params.field_path,
    });
  }

  // 检查密码
  checkPassWord() {
    const { getFieldValue, setFieldsValue } = this.props.form;
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

    let result = this.decrypt(cur_service.Config[key].current, config);
    if (result == '45d3161f8aa54588ab3901fa954cc01a') {
      result = '';
      this.replacePasswd(key, result);
      this.setState(
        {
          checkVisible: false,
          password_key: password,
        },
        () => setFieldsValue({ input_password: '' })
      );
    } else if (!result) {
      message.error('密码输入错误！请重新输入');
    } else {
      this.replacePasswd(key, result);
      this.setState(
        {
          checkVisible: false,
          password_key: password,
        },
        () => setFieldsValue({ input_password: '' })
      );
    }
  }

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

  // 关联主机复选框选择值后
  onChangeHosts = (checkedValues) => {
    this.setState({
      changeHosts: checkedValues,
    });
  };

  // 保存关联主机配置
  saveHosts() {
    const { changeHosts, linkIndex, linkName } = this.state;
    const { cur_service } = this.props.ServiceStore;
    const { setServiceConfigModify } = this.props.actions;
    const config = Object.assign({}, cur_service.configModify);

    // console.log(config, '关联前的config----')
    const saveHostsItem = {
      hosts:
        changeHosts.length !== 0
          ? cur_service.Config[linkName].current[linkIndex].hosts.length > 1
            ? cur_service.Config[linkName].current[linkIndex].hosts +
              ',' +
              changeHosts.join(',')
            : changeHosts.join(',')
          : cur_service.Config[linkName].current[linkIndex].hosts,
      field:
        typeof cur_service.Config[linkName].current[linkIndex] === 'string'
          ? cur_service.Config[linkName].current[linkIndex]
          : cur_service.Config[linkName].current[linkIndex].field,
    };

    // console.log(saveHostsItem, '关联主机保存-----', changeHosts, 'changeHosts.length === 0', changeHosts.length === 0)
    // const arr = []
    // if (Array.isArray(config['Config.' + linkName + '.current'])) {
    //     config['Config.' + linkName + '.current'].splice(linkIndex, 1, saveHostsItem);
    //     config['Config.' + linkName + '.current'] = config['Config.' + linkName + '.current'];
    // } else {
    //     const stringItem = {
    //         hosts: cur_service.Config[linkName].current[0].hosts ? cur_service.Config[linkName].current[0].hosts : '',
    //         field: cur_service.Config[linkName].current[0].field
    //     }
    //     arr.push(stringItem, saveHostsItem);
    //     config['Config.' + linkName + '.current'] = arr;
    // }

    cur_service.Config[linkName].current.splice(linkIndex, 1, saveHostsItem);
    config['Config.' + linkName + '.current'] =
      cur_service.Config[linkName].current;
    cur_service.configModify = config;

    // console.log(config, '关联后的config----')
    setServiceConfigModify(cur_service);
    // console.log(cur_service.configModify, ' cur_service.configModify-----', cur_service)
    this.cancelLink();
  }

  // 关闭关联主机
  cancelLink() {
    this.setState({
      hostVisible: false,
      linkIndex: null,
      linkName: '',
    });
  }

  // 过滤已选主机   cur_service.Config[key.split('.')[0]].current[index].hosts.split(',')
  filterHosts(data: any) {
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
      });
    this.disabledLink = disabledArray;
  }

  render() {
    const { cur_service = {} } = this.props.ServiceStore;
    const Config = cur_service.Config;
    const { checkVisible, hostsList, hostVisible, changeHosts } = this.state;
    const { getFieldDecorator, setFieldsValue } = this.props.form;
    const disabledLink = this.disabledLink;

    if (Config) {
      return (
        <div className="edit-panel config-box">
          <Form>{this.buildConfigEditByNode(Config)}</Form>
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
                  {hostsList &&
                    hostsList.map((item, index) => {
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
    } else {
      return (
        <div className="edit-panel config-box">
          <p style={{ lineHeight: '20px' }}>无</p>
        </div>
      );
    }
  }
}
export default Form.create<Prop>()(ConfigComp);
