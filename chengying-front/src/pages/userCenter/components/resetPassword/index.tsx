import * as React from 'react';
import { Form, Input, Modal, message, Alert } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import './style.scss';

const FormItem = Form.Item;

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 6 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16 },
  },
};

interface IProps extends FormComponentProps {
  userInfo?: any;
  isCheckedResetPwd?: boolean;
  visible: boolean;
  onCancel: Function;
  onSubmit: Function;
}

class ResetPassword extends React.PureComponent<IProps, any> {
  // 密码格式校验
  validPass = (rules: any, value: any, callback: any) => {
    const { form } = this.props;
    if (value && value === form.getFieldValue('oldPass')) {
      callback('新密码不能与旧密码一致');
    } else if (
      value &&
      !/^(?![A-z0-9]+$)(?![A-z~!@#$%^&*]+$)(?![0-9~!@#$%^&*]+$)([0-9A-Za-z~!@#$%^&*]{8,20})$/.test(
        value
      )
    ) {
      callback('密码须包含数字、字母和特殊字符(支持~!@#$%^&*)，不低于8位');
    }
    callback();
  };

  // 重置密码
  handleResetSubmit = () => {
    const { form } = this.props;
    form.validateFields((err: any, value: any) => {
      if (err) {
        return;
      }
      if (value.newPass !== value.ensurePass) {
        message.error('两次密码填写不一致');
        return;
      }
      this.props.onSubmit(value);
    });
  };

  // 关闭弹窗
  handleCancel = () => {
    this.props.onCancel();
  };

  render() {
    const { form, userInfo, visible, isCheckedResetPwd } = this.props;
    const { getFieldDecorator } = form;
    return (
      <Modal
        className="reset-password"
        maskClosable={false}
        destroyOnClose={true}
        visible={visible}
        title="重置密码"
        onOk={this.handleResetSubmit}
        onCancel={this.handleCancel}>
        {isCheckedResetPwd && (
          <Alert
            type="warning"
            showIcon
            message="管理员开启了初始密码修改，您需要重置初始密码，否则将退出登录影响平台功能使用。"
          />
        )}
        <Form>
          {userInfo?.username ? (
            <FormItem {...formItemLayout} label="账号">
              {userInfo.username || ''}
            </FormItem>
          ) : (
            <FormItem {...formItemLayout} label="旧密码">
              {getFieldDecorator('oldPass', {
                rules: [{ required: true, message: '新密码不能为空' }],
              })(<Input.Password placeholder="请输入旧密码" />)}
            </FormItem>
          )}
          <FormItem {...formItemLayout} label="新密码">
            {getFieldDecorator('newPass', {
              rules: [
                { required: true, message: '新密码不能为空' },
                { validator: this.validPass },
              ],
            })(<Input.Password placeholder="请输入新密码" />)}
          </FormItem>
          <FormItem {...formItemLayout} label="确认新密码">
            {getFieldDecorator('ensurePass', {
              rules: [
                { required: true, message: '确认新密码不能为空' },
                { validator: this.validPass },
              ],
            })(<Input.Password placeholder="请再次输入新密码" />)}
          </FormItem>
        </Form>
      </Modal>
    );
  }
}
export default Form.create<IProps>()(ResetPassword);
