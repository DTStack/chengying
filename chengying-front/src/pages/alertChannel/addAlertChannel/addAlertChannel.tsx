import * as React from 'react';
import {
  Form,
  Icon,
  Radio,
  Button,
  Checkbox,
  Input,
  Select,
  Upload,
  message,
  Tooltip,
  Breadcrumb,
} from 'antd';
import { Link } from 'react-router-dom';
import channelMode from './channelCode';
import { alertChannelService } from '@/services';
import './style.scss';
const FormItem = Form.Item;
const RadioGroup = Radio.Group;
const TextArea = Input.TextArea;
const Option = Select.Option;

const DefaultUrls = [
  'http://dtalert:8085/gate/alert/sms_sync',
  'http://dtalert:8085/gate/alert/mail_sync',
  'http://dtalert:8085/gate/alert/ding_sync',
  'http://dtalert:8085/gate/alert/custom_sync',
  '', //这个千万别删
  'http://dtalert:8085/gate/alert/wechat_sync',
];
const DefaultFormData = {
  alertGateType: 1,
  alertGateCode: '',
  alertGateSource: '',
  alertGateName: '',
  alertGateJson: '',
  alertTemplate: '',
  alertTemplateId: '',
  webhook: '',
  filePath: '',
  phones: '',
  emails: '',
  dings: '',
  isDefault: false,
  uploadImage: false,
  url: DefaultUrls[0],
  subject: '',
};
const ChannelName = [
  'dtstack_sms',
  'dtstack_mail',
  'dtstack_ding',
  'dtstack_custom',
  '', //这个也不要删
  'dtstack_wechat',
];

interface P {
  form: any;
  location?: any;
}
interface S {
  gid?: number;
  tid?: number;
  formData: any;
  testActive: boolean;
  placeholder: {
    jar: string;
    api: string;
    sms_yp: string;
    mail_dt: string;
  };
  channels: any[];
  isDisable: boolean;
}
class AddAlertChannel extends React.Component<P, S> {
  constructor(props: any) {
    super(props);
  }

  state: S = {
    gid: this.getParamsFromUrl(window.location.href).gid || null,
    tid: this.getParamsFromUrl(window.location.href).tid || null,
    formData: DefaultFormData,
    testActive: false,
    placeholder: {
      jar: '{"classname":"com.dtstack.sender.sms.xxxsender"}',
      api: '{\n"cookiestore": false,\n"configs": [{\n"url": "",\n"method": "get",\n"header": {},\n"body": {}\n}],\n"context": {}\n} ',
      sms_yp: '请按照此格式输入配置信息：\n{"yp_api_key":"xxxxxx"}',
      mail_dt:
        '{\n"mail.smtp.host":"smtp.yeah.net",\n"mail.smtp.port":"25",\n"mail.smtp.ssl.enable":"false",\n"mail.smtp.username":"daishu@dtstack.com",\n"mail.smtp.password":"xxxx",\n"mail.smtp.from":"daishu@dtstack.com"\n}',
    },
    channels: [],
    isDisable: false,
  };

  componentDidMount() {
    const { gid, tid } = this.state;
    this.getAlertChannels();
    if (gid) {
      this.getChannelData(gid + '', tid + '');
    } else {
      this.setDefaultFormData(Object.assign({}, DefaultFormData));
    }
  }

  getParamsFromUrl(url: string): any {
    var obj = {};
    var keyvalue = [];
    var key = '';
    var value = '';
    var paraString = url.substring(url.indexOf('?') + 1, url.length).split('&');
    for (var i in paraString) {
      keyvalue = paraString[i].split('=');
      key = keyvalue[0];
      value = keyvalue[1];
      obj[key] = value;
    }
    return obj;
  }

  // 获取所有告警通道用于重复判断
  getAlertChannels = () => {
    const self = this;
    const channels: any = [];
    alertChannelService
      .dtstackAlertChannelList({ page: 1, size: 100 })
      .then((rst: any) => {
        if (rst.data.success) {
          alertChannelService.getAlertNotifications().then((res: any) => {
            for (const c of rst.data.data.data) {
              for (const g of res) {
                if (c.alertGateName === g.name) {
                  channels.push({
                    ...c,
                    gid: g.id,
                    type: g.type,
                    isDefault: g.isDefault,
                  });
                }
              }
            }
            self.setState({
              channels: channels,
            });
          });
        } else {
          message.error('获取通道失败！');
        }
      });
  };

  handleCheckNameReuse = (rule: any, value: string, callback: any) => {
    const { channels } = this.state;
    let isExist = false;
    for (const c of channels) {
      if (c.alertGateName === value) {
        isExist = true;
      }
    }
    if (isExist) {
      // form.validateFields(['alertGateName'], { force: true });
      callback('告警通道名称已经存在！');
    } else {
      callback();
    }
  };

  getChannelData = (gid: string, tid: string) => {
    const self = this;
    alertChannelService
      .getDtstackAlertDetail({
        id: tid,
      })
      .then((res: any) => {
        if (res.data.success) {
          let dtAlert = res.data.data;
          alertChannelService
            .getGrafanaAlertDetail({
              alertId: gid,
            })
            .then((rst: any) => {
              // debugger;

              const gAlert = rst.data;
              dtAlert = {
                ...dtAlert,
                phones: gAlert.settings.phones,
                emails: gAlert.settings.emails,
                dings: gAlert.settings.dings,
                isDefault: gAlert.isDefault,
                uploadImage: gAlert.settings.uploadImage,
                url: gAlert.settings.url,
                webhook: gAlert.settings.webhook,
                subject: gAlert.settings.subject,
                created: gAlert.created,
                updated: gAlert.updated,
              };
              // debugger;
              self.setDefaultFormData(dtAlert);
            });
        } else {
          message.error(res.message);
        }
      });
  };

  /**
   * 重置表单初始值
   * @memberof AddAlertChannel
   */
  setDefaultFormData = (formData: any) => {
    const { setFieldsValue } = this.props.form;
    this.setState(
      {
        formData: formData,
      },
      () => {
        setFieldsValue(formData);
      }
    );
  };

  /**
   * 表单input、checkbox、textarea输入
   * @memberof AddAlertChannel
   */
  handleFieldChange = (e: any) => {
    const { formData } = this.state;
    const field = e.target.name;
    formData[field] = e.target.value;
    this.setState({
      formData,
    });
  };

  /**
   * 表单checkbox选择
   * @memberof AddAlertChannel
   */
  handleCheckboxChange = (e: any) => {
    // debugger;
    const field = e.target.name;
    const { formData } = this.state;
    formData[field] = e.target.checked;
    this.setState({
      formData,
    });
  };

  /**
   * 告警通道选择
   * @memberof AddAlertChannel
   */
  handleGateTypeChange = (e: any) => {
    const { formData } = this.state;
    const { setFieldsValue } = this.props.form;
    formData.alertGateType = e.target.value;
    formData.alertGateCode = '';
    formData.url = DefaultUrls[e.target.value - 1];
    this.setState({
      formData,
    });
    setFieldsValue(formData);
  };

  /**
   * 通道模式选择
   * @memberof AddAlertChannel
   */
  handleGateCodeChange = (value: any) => {
    const { formData } = this.state;
    formData.alertGateCode = value;
    this.setState({
      formData,
    });
  };

  /**
   * jar包上传
   * @memberof AddAlertChannel
   */
  fileUploadChange = (info: any) => {
    const { formData } = this.state;
    if (info.file.status !== 'uploading') {
    }
    if (info.file.status === 'done') {
      message.success('上传成功');
      if (info.file.response.success) {
        formData.filePath = info.file.response.data;
        this.setState({
          formData,
        });
      } else {
        message.error(info.file.response.message);
      }
    } else if (info.file.status === 'error') {
      message.error('上传失败');
    }
  };

  /**
   * 抽取grafana告警通道参数
   * @memberof AddAlertChannel
   */
  getGrafanaParams = () => {
    // debugger
    const { formData, gid } = this.state;
    let params: any = gid
      ? {
          id: parseInt(gid + ''),
          created: formData.created,
          updated: formData.updated,
        }
      : {};
    switch (formData.alertGateType) {
      case 1:
        params = {
          ...params,
          name: formData.alertGateName,
          type: ChannelName[formData.alertGateType - 1],
          isDefault: formData.isDefault,
          settings: {
            httpMethod: 'POST',
            autoResolve: true,
            uploadImage: formData.uploadImage,
            url: formData.url,
            phones: formData.phones,
            subject: formData.subject,
            source: formData.alertGateSource,
          },
        };
        break;
      case 2:
        params = {
          ...params,
          name: formData.alertGateName,
          type: ChannelName[formData.alertGateType - 1],
          isDefault: formData.isDefault,
          settings: {
            httpMethod: 'POST',
            autoResolve: true,
            uploadImage: formData.uploadImage,
            url: formData.url,
            emails: formData.emails,
            subject: formData.subject,
            source: formData.alertGateSource,
          },
        };
        break;
      case 3:
        params = {
          ...params,
          name: formData.alertGateName,
          type: ChannelName[formData.alertGateType - 1],
          isDefault: formData.isDefault,
          settings: {
            httpMethod: 'POST',
            autoResolve: true,
            uploadImage: formData.uploadImage,
            url: formData.url,
            dings: formData.dings,
            subject: formData.subject,
            source: formData.alertGateSource,
          },
        };
        break;
      case 4:
        params = {
          ...params,
          name: formData.alertGateName,
          type: ChannelName[formData.alertGateType - 1],
          isDefault: formData.isDefault,
          settings: {
            httpMethod: 'POST',
            autoResolve: true,
            uploadImage: formData.uploadImage,
            url: formData.url,
            identification: formData.identification,
            subject: formData.subject,
            source: formData.alertGateSource,
          },
        };
        break;
      case 6:
        params = {
          ...params,
          name: formData.alertGateName,
          type: ChannelName[formData.alertGateType - 1],
          isDefault: formData.isDefault,
          settings: {
            httpMethod: 'POST',
            autoResolve: true,
            uploadImage: formData.uploadImage,
            url: formData.url,
            webhook: formData.webhook,
            subject: formData.subject,
            source: formData.alertGateSource,
          },
        };
        break;
    }
    return params;
  };

  /**
   * 测试通道
   * @memberof AddAlertChannel
   */
  handleTestChannel = () => {
    const params = this.getGrafanaParams();
    alertChannelService.grafanaAlertChannelTest(params).then((res: any) => {
      if (res.data.message === 'Test notification sent') {
        message.success('测试告警发送成功！');
      }
    });
  };

  /**
   * 保存告警通道
   * @memberof AddAlertChannel
   */
  handleFormSubmit = (e: any) => {
    e.preventDefault();
    const self = this;
    const { formData, gid, tid } = this.state;
    const dtstackParams = {
      alertGateType: formData.alertGateType,
      alertGateCode: formData.alertGateCode,
      alertGateSource: formData.alertGateSource,
      alertGateName: formData.alertGateName,
      alertGateJson: formData.alertGateJson,
      alertTemplate: formData.alertTemplate,
      alertTemplateId: formData.alertTemplateId,
      filePath: formData.filePath,
      id: tid,
    };
    const grafanaParams = this.getGrafanaParams();
    if (gid) {
      // 编辑
      this.props.form.validateFields((errors: any) => {
        if (!errors) {
          alertChannelService
            .grafanaAlertChannelUpdate(grafanaParams)
            .then((rst: any) => {
              rst = rst.data;
              if (rst.id) {
                alertChannelService
                  .dtstackAlertChannelSave(dtstackParams)
                  .then((res: any) => {
                    res = res.data;
                    if (res.success) {
                      message.success('保存成功！');
                      self.setState({
                        testActive: true,
                      });
                    } else {
                      message.error(res.message);
                    }
                  });
              } else {
                message.error(rst.message);
              }
            });
        }
      });
    } else {
      // 创建
      this.props.form.validateFields((errors: any) => {
        if (!errors) {
          alertChannelService
            .grafanaAlertChannelSave(grafanaParams)
            .then((rst: any) => {
              rst = rst.data;
              if (rst.id) {
                alertChannelService
                  .dtstackAlertChannelSave(dtstackParams)
                  .then((res: any) => {
                    res = res.data;
                    if (res.success) {
                      message.success('保存成功！');
                      self.setState({
                        testActive: true,
                        isDisable: true,
                      });
                    } else {
                      message.error(res.message);
                    }
                  });
              } else {
                message.error(rst.message);
              }
            });
        }
      });
    }
  };

  render() {
    const { getFieldDecorator, getFieldsValue } = this.props.form;
    const { formData, placeholder, testActive, gid, isDisable } = this.state;
    const formItemLayout = {
      labelCol: {
        xs: { span: 24 },
        sm: { span: 8 },
      },
      wrapperCol: {
        xs: { span: 24 },
        sm: { span: 10 },
      },
    };
    const uploadProp = {
      name: 'file',
      action: '/gate/alert/jarUpload',
      onChange: this.fileUploadChange,
    };

    let channelModeList = [];
    switch (getFieldsValue().alertGateType) {
      case 1:
        channelModeList = channelMode.sms;
        break;
      case 2:
        channelModeList = channelMode.mail;
        break;
      case 3:
        channelModeList = channelMode.dingTalk;
        break;
      case 4:
        channelModeList = channelMode.identification;
        break;
      case 6:
        channelModeList = channelMode.webhook;
        break;
    }
    const upload = (
      <FormItem {...formItemLayout} label="上传文件">
        <Upload {...uploadProp}>
          <a href="javascript:;">选择jar文件</a>
        </Upload>
      </FormItem>
    );
    return (
      <div className="add-alert-channel">
        <Breadcrumb>
          <Breadcrumb.Item>
            <Link to="/deploycenter/monitoring/alert?activeKey=channel">
              告警通道
            </Link>
          </Breadcrumb.Item>
          <Breadcrumb.Item>
            {gid ? '编辑告警通道' : '新增告警通道'}
          </Breadcrumb.Item>
        </Breadcrumb>
        <div className="box-shadow-style overflow-scroll">
          <Form onSubmit={this.handleFormSubmit}>
            <FormItem {...formItemLayout} label="告警类型">
              {getFieldDecorator('alertGateType', {
                rules: [{ required: true, message: '请选择告警通道' }],
                initialValue: formData.alertGateType,
              })(
                <RadioGroup
                  name="alertGateType"
                  onChange={this.handleGateTypeChange}>
                  <Radio value={1}>短信通道</Radio>
                  <Radio value={2}>邮件通道</Radio>
                  <Radio value={3}>钉钉通道</Radio>
                  <Radio value={6}>企业微信通道</Radio>
                  <Radio value={4}>自定义通道</Radio>
                </RadioGroup>
              )}
            </FormItem>
            <FormItem {...formItemLayout} label="通道模式">
              {getFieldDecorator('alertGateCode', {
                rules: [
                  {
                    required: true,
                    message: '通道模式不可为空',
                  },
                ],
                initialValue: formData.alertGateCode,
              })(
                <Select onChange={this.handleGateCodeChange}>
                  {channelModeList.map((o, index) => {
                    return (
                      <Option value={o.value} key={`${index}`}>
                        {o.title}
                      </Option>
                    );
                  })}
                </Select>
              )}
            </FormItem>
            {getFieldsValue().alertGateCode &&
              getFieldsValue().alertGateCode.indexOf('jar') !== -1 &&
              upload}
            <FormItem {...formItemLayout} label="使用场景">
              {getFieldDecorator('alertGateSource', {
                rules: [
                  { required: true, message: '产品不能为空' },
                  { max: 32, message: '场景名称不能超过32个字符' },
                ],
                initialValue: formData.alertGateSource,
              })(
                <Input
                  placeholder="使用场景一般指产品名称，请输入产品名称，不超过32个字符"
                  name="alertGateSource"
                  onChange={this.handleFieldChange}
                />
              )}
            </FormItem>
            <FormItem {...formItemLayout} label="通道名称">
              {getFieldDecorator('alertGateName', {
                rules: [
                  { required: true, message: '通道名称不能为空' },
                  { max: 32, message: '通道名称不能超过32个字符' },
                ],
                initialValue: formData.alertGateName,
              })(
                <Input
                  name="alertGateName"
                  disabled={!!gid}
                  placeholder="请输入通道名称，不超过32个字符"
                  onChange={this.handleFieldChange}
                />
              )}
            </FormItem>
            <FormItem {...formItemLayout} label="通道配置信息">
              {getFieldDecorator('alertGateJson', {
                rules: [{ required: true, message: '配置信息不能为空' }],
                initialValue: formData.alertGateJson,
              })(
                <TextArea
                  name="alertGateJson"
                  placeholder={
                    getFieldsValue().alertGateCode &&
                    (getFieldsValue().alertGateCode === 'sms_yp'
                      ? placeholder.sms_yp
                      : getFieldsValue().alertGateCode === 'mail_dt'
                      ? placeholder.mail_dt
                      : getFieldsValue().alertGateCode.indexOf('jar') !== -1
                      ? placeholder.jar
                      : getFieldsValue().alertGateCode.indexOf('api') !== -1
                      ? placeholder.api
                      : '')
                  }
                  rows={6}
                  onChange={this.handleFieldChange}
                />
              )}
            </FormItem>
            <FormItem {...formItemLayout} label="通知消息模版">
              {getFieldDecorator('alertTemplate', {
                rules: [{ required: true, message: '消息模版不能为空' }],
                initialValue: formData.alertTemplate,
              })(
                <TextArea
                  name="alertTemplate"
                  placeholder={
                    '请按照此格式填写："【企业名称】$' +
                    `{message}，请登录EasyManage处理"，如【${window.APPCONFIG.company}】$` +
                    '{message}，请登录EasyManage处理'
                  }
                  rows={4}
                  onChange={this.handleFieldChange}
                />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label={
                <span>
                  适用于所有告警
                  <Tooltip title="对所有告警使用此通知">
                    <Icon style={{ marginLeft: 2 }} type="question-circle" />
                  </Tooltip>
                </span>
              }>
              {getFieldDecorator('isDefault', {
                initialValue: formData.isDefault,
                valuePropName: 'checked',
              })(
                <Checkbox
                  name="isDefault"
                  onChange={this.handleCheckboxChange}></Checkbox>
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label={
                <span>
                  包含图片
                  <Tooltip title="捕获监控图片并将其包含在告警信息中">
                    <Icon style={{ marginLeft: 2 }} type="question-circle" />
                  </Tooltip>
                </span>
              }>
              {getFieldDecorator('uploadImage', {
                initialValue: formData.uploadImage,
                valuePropName: 'checked',
              })(
                <Checkbox
                  name="uploadImage"
                  onChange={this.handleCheckboxChange}></Checkbox>
              )}
            </FormItem>
            <FormItem {...formItemLayout} label="URL">
              {getFieldDecorator('url', {
                rules: [{ required: true, message: 'URL不能为空' }],
                initialValue: formData.url,
              })(<Input name="url" onChange={this.handleFieldChange} />)}
            </FormItem>
            {formData.alertGateType === 2 ? (
              <FormItem {...formItemLayout} label="Email">
                {getFieldDecorator('emails', {
                  rules: [{ required: true, message: 'emails不能为空' }],
                  initialValue: formData.emails,
                })(
                  <Input
                    name="emails"
                    placeholder="填写告警邮箱，多个用逗号分隔"
                    onChange={this.handleFieldChange}
                  />
                )}
              </FormItem>
            ) : (
              ''
            )}
            {formData.alertGateType === 1 ? (
              <FormItem {...formItemLayout} label="Phones">
                {getFieldDecorator('phones', {
                  rules: [{ required: true, message: 'phones不能为空' }],
                  initialValue: formData.phones,
                })(
                  <Input
                    name="phones"
                    placeholder="填写告警手机号，多个用逗号分隔"
                    onChange={this.handleFieldChange}
                  />
                )}
              </FormItem>
            ) : (
              ''
            )}
            {formData.alertGateType === 3 ? (
              <FormItem {...formItemLayout} label="Dings">
                {getFieldDecorator('dings', {
                  rules: [{ required: true, message: 'dings不能为空' }],
                  initialValue: formData.dings,
                })(
                  <Input
                    name="dings"
                    placeholder="填写告警钉钉，多个用逗号分隔"
                    onChange={this.handleFieldChange}
                  />
                )}
              </FormItem>
            ) : (
              ''
            )}
            {console.log(formData)}
            {formData.alertGateType === 6 ? (
              <FormItem {...formItemLayout} label="Webhook">
                {getFieldDecorator('webhook', {
                  rules: [{ required: true, message: 'webhook不能为空' }],
                  initialValue: formData.webhook,
                })(
                  <Input
                    name="webhook"
                    placeholder="多个用逗号分隔"
                    onChange={this.handleFieldChange}
                  />
                )}
              </FormItem>
            ) : (
              ''
            )}
            {formData.alertGateType !== 1 && (
              <FormItem {...formItemLayout} label="Subject">
                {getFieldDecorator('subject', {
                  rules: [{ required: true, message: 'subject不能为空' }],
                  initialValue: formData.subject,
                })(<Input name="subject" onChange={this.handleFieldChange} />)}
              </FormItem>
            )}
            <FormItem className="add-alert-channel-bottom">
              <Button
                type="primary"
                disabled={!testActive}
                style={{ marginRight: 10 }}
                onClick={this.handleTestChannel}>
                发送测试
              </Button>
              <Button type="primary" disabled={isDisable} htmlType="submit">
                保存
              </Button>
            </FormItem>
          </Form>
        </div>
      </div>
    );
  }
}
export default Form.create()(AddAlertChannel);

// import * as React from 'react';

// class A extends React.Component{
//   render(){
//     return (
//       <p>HOP</p>
//     )
//   }
// }
// export default A
