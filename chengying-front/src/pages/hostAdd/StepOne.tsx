import * as React from 'react';
import { Form, Input, Radio, message } from 'antd';
import { formItemCenterLayout } from '@/constants/formLayout';

const FormItem = Form.Item;
const { TextArea } = Input;
const RadioGroup = Radio.Group;

interface State {
  value: string;
  fileName: string;
  password: string;
  ip: string;
  IPArr: any;
}

class StepOne extends React.Component<any, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      value: '密码',
      fileName: '',
      password: '',
      ip: '',
      IPArr: null,
    };
  }

  validIp = (rule: any, value: any, callback: any) => {
    const regx =
      /^(?:(?:^|,)(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])(?:\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3})+$/;
    if (value && !regx.test(value)) {
      callback('存在格式不正确的ip地址');
    } else {
      if (value && regx.test(value)) {
        const isRep = this.repeatIps(value);
        if (isRep) {
          callback(isRep);
        } else {
          callback();
        }
      } else {
        callback();
      }
    }
  };

  repeatIps = (values: any) => {
    let msg = '';
    const nary = values ? values.split(',').sort() : [];
    for (var i = 0; i < nary.length - 1; i++) {
      if (nary[i] === nary[i + 1]) {
        msg = `ip:${nary[i]}已存在，请检查输入！`;
      }
    }
    return msg;
  };

  checkIP = (e: any) => {
    const v = e.target.value;
    let IPArr = [];
    v.indexOf(',') > -1 ? (IPArr = v.split(',')) : IPArr.push(v);
    this.setState({ IPArr: IPArr });
  };

  onChange = (e: any) => {
    const checked = e.target.value;
    this.setState({
      value: checked,
      fileName: '',
      password: '',
    });
  };

  handleChange = (e: any) => {
    const { form } = this.props;
    const eValue = e.target.value;
    this.setState({
      password: eValue,
    });
    form.setFieldsValue({ password: eValue });
  };

  private inputDOM = React.createRef<HTMLInputElement>();
  uploadSecret = () => {
    const { form } = this.props;

    const file = (this.inputDOM as any).files[0];
    if (!file) {
      message.warn('请选择上传的秘钥');
      return;
    }
    const flag = file.size / 1024 >= 100;
    if (flag) {
      message.warn('请选择小于100Kb文件上传！');
      return;
    }
    if ((window as any).FileReader) {
      var reader = new FileReader();
      reader.readAsText(file);
      reader.onload = (e: any) => {
        form.setFieldsValue({ pk: e.target.result });
      };
      this.setState({
        fileName: file.name,
      });
      (this.inputDOM as any).value = '';
    } else {
      message.warn('浏览器不支持');
    }
    console.log(this.props.form);
  };

  render() {
    const { value, fileName, password } = this.state;
    const { getFieldDecorator } = this.props.form; // console.log(this.props.shellCMD);

    return (
      <div style={{ margin: '50px auto', width: '80%' }}>
        <Form className="login-form" style={{ width: '60%' }}>
          <FormItem {...formItemCenterLayout} label="机器IP">
            {getFieldDecorator('host', {
              rules: [
                { required: true, message: 'please input machine ip' },
                { validator: this.validIp },
              ],
            })(
              <TextArea
                onBlur={(e) => this.checkIP(e)}
                autosize={{ minRows: 2, maxRows: 6 }}
                rows={4}
                placeholder="输入机器ip，英文逗号分隔（如：172.1.1.1,172.1.1.2）"
              />
            )}
          </FormItem>
          <FormItem {...formItemCenterLayout} label="SSH端口">
            {getFieldDecorator('port', {
              initialValue: '22',
              rules: [
                { required: true, message: '必填项' },
                { max: 10, message: '端口最大不能超过10位!' },
                {
                  pattern: /^[1-9]{1}[0-9]*$/,
                  message: '端口号仅限为非0开头的数字!',
                },
              ],
            })(
              <Input autoComplete="true" type="port" placeholder="输入端口号" />
            )}
          </FormItem>
          <FormItem {...formItemCenterLayout} label="登陆方式">
            <RadioGroup onChange={(e) => this.onChange(e)} defaultValue="密码">
              <Radio value="密钥">密钥</Radio>
              <Radio value="密码">密码</Radio>
            </RadioGroup>
          </FormItem>
          <FormItem
            {...formItemCenterLayout}
            label="用户名"
            extra="需要sudo+NOPASSWD权限">
            {getFieldDecorator('user', {
              initialValue: 'admin',
              rules: [{ required: true }],
            })(<Input placeholder="输入用户名" type="user" />)}
          </FormItem>
          {value === '密钥' ? (
            <FormItem {...formItemCenterLayout} label="密钥">
              {getFieldDecorator('pk', {
                rules: [{ required: true, message: '密钥不可为空!' }],
              })(<Input style={{ display: 'none' }} />)}
              <label className="ant-btn" style={{ lineHeight: '28px' }}>
                上传秘钥文件
                <input
                  ref={(input: any) => (this.inputDOM = input)}
                  type="file"
                  onChange={this.uploadSecret}
                  style={{ display: 'none' }}
                />
              </label>
              <span>{fileName}</span>
            </FormItem>
          ) : (
            <FormItem {...formItemCenterLayout} label="密码">
              {getFieldDecorator('password', {
                rules: [{ required: true, message: '密码不可为空!' }],
              })(
                <Input
                  style={{
                    display: 'none',
                  }}
                  autoComplete="off"
                  placeholder="密码"
                />
              )}
              <Input
                type="password"
                value={password}
                onChange={this.handleChange}
                autoComplete="off"
                placeholder="密码"
              />
            </FormItem>
          )}
        </Form>
      </div>
    );
  }
}

export default Form.create({})(StepOne);
