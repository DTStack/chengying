import * as React from 'react';
import { Modal, Form, Input, Tooltip, Icon, message } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { imageStoreService } from '@/services';
import '../style.scss';

const FormItem = Form.Item;

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 6 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 15 },
  },
};
interface IProps extends FormComponentProps {
  handleCancel: () => void;
  getImageStoreList: Function;
  imageStoreInfo: any;
  clusterId: number;
  isEdit: boolean;
}
interface IState {
  modalLoading: boolean;
}

class EditModal extends React.PureComponent<IProps, IState> {
  state: IState = {
    modalLoading: false,
  };

  // 提交
  handleSubmit = () => {
    const { form, clusterId, isEdit } = this.props;
    form.validateFields((err, values) => {
      if (!err) {
        this.setState({ modalLoading: true });
        const operate = isEdit ? 'updateImageStore' : 'createImageStore';
        imageStoreService[operate]({
          ...values,
          clusterId,
        }).then((res: any) => {
          res = res.data;
          this.setState({ modalLoading: false });
          if (res.code === 0) {
            message.success('执行成功');
            this.props.getImageStoreList();
            this.props.handleCancel();
          } else {
            message.error(res.msg);
          }
        });
      }
    });
  };

  render() {
    const { modalLoading } = this.state;
    const { imageStoreInfo, isEdit, form } = this.props;
    const { getFieldDecorator } = form;
    return (
      <Modal
        title={`${isEdit ? '编辑' : '添加'}仓库`}
        visible={true}
        confirmLoading={modalLoading}
        onOk={this.handleSubmit}
        onCancel={this.props.handleCancel}>
        <Form>
          <FormItem label="仓库名称" {...formItemLayout}>
            {getFieldDecorator('name', {
              initialValue: isEdit ? imageStoreInfo.name : '',
              rules: [
                { required: true, message: '仓库名称不能为空' },
                {
                  pattern: /^\S{1,64}$/,
                  message: '除空格外，支持其他字符，最大长度不超过64',
                },
              ],
            })(<Input placeholder="请输入仓库名称" />)}
          </FormItem>
          <FormItem label="仓库别名" {...formItemLayout}>
            {getFieldDecorator('alias', {
              initialValue: isEdit ? imageStoreInfo.alias : '',
              rules: [
                { required: true, message: '仓库别名不能为空' },
                {
                  pattern: /^[a-z0-9]([-a-z0-9]{0,62}[a-z0-9])?$/,
                  message:
                    '支持小写字母、数字、中划线，且不以中划线开头或结尾，最大长度不超过64',
                },
              ],
            })(<Input placeholder="请输入仓库别名" />)}
          </FormItem>
          <FormItem
            label={
              <span>
                仓库地址{' '}
                <Tooltip title="http://docker.io">
                  <Icon type="info-circle" />
                </Tooltip>
              </span>
            }
            {...formItemLayout}>
            {getFieldDecorator('address', {
              initialValue: isEdit ? imageStoreInfo.address : '',
              rules: [
                { required: true, message: '仓库地址不能为空' },
                {
                  pattern: /^[^\u4e00-\u9fa5\s]+$/,
                  message: '不支持中文、空格',
                },
              ],
            })(<Input placeholder="请输入仓库地址" />)}
          </FormItem>
          <FormItem label="用户名" {...formItemLayout}>
            {getFieldDecorator('username', {
              initialValue: isEdit ? imageStoreInfo.username : '',
              rules: [{ required: true, message: '用户名不能为空' }],
            })(<Input placeholder="请输入用户名" />)}
          </FormItem>
          <FormItem label="密码" {...formItemLayout}>
            {getFieldDecorator('password', {
              initialValue: isEdit ? imageStoreInfo.password : '',
              rules: [{ required: true, message: '密码不能为空' }],
            })(<Input type="password" placeholder="请输入密码" />)}
          </FormItem>
          <FormItem label="邮箱" {...formItemLayout}>
            {getFieldDecorator('email', {
              initialValue: isEdit ? imageStoreInfo.email : '',
              rules: [
                { pattern: /^\S+@\S+$/, message: '请输入正确的邮箱格式' },
              ],
            })(<Input type="email" placeholder="请输入邮箱" />)}
          </FormItem>
          <FormItem style={{ display: 'none' }}>
            {getFieldDecorator('id', {
              initialValue: isEdit ? imageStoreInfo.id : undefined,
            })(<Input />)}
          </FormItem>
        </Form>
      </Modal>
    );
  }
}
export default Form.create<IProps>()(EditModal);
