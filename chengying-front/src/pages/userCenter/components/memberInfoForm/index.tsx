import * as React from 'react';
import { Form, Input, Radio, Transfer } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import * as Cookie from 'js-cookie';

const FormItem = Form.Item;
const RadioGroup = Radio.Group;

const formLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 6 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 14 },
  },
};

interface IProps extends FormComponentProps {
  memberInfo: any;
  canEdit?: boolean;
  formLayout?: any;
  clusterList?: any[];
}

class MemberInfo extends React.PureComponent<IProps, any> {
  state = {
    targetKeys: [],
  };

  componentDidMount() {
    const { memberInfo } = this.props;
    if (memberInfo?.cluster_list) {
      const targetKeys = memberInfo?.cluster_list.map(
        (item) => item.cluster_id
      );
      this.setState({
        targetKeys,
      });
    }
  }

  renderItem = (item) => {
    const customLabel = (
      <span className="custom-item" key={item.id}>
        {item.name}
      </span>
    );
    return {
      label: customLabel, // for displayed item
      value: item.name, // for title and filter matching
    };
  };

  handleChange = (targetKeys, direction, moveKeys) => {
    console.log(targetKeys, direction, moveKeys);
    this.setState({ targetKeys });
  };

  render() {
    const { form, memberInfo, canEdit, clusterList } = this.props;
    const { targetKeys } = this.state;
    const { getFieldDecorator } = form;
    const formItemLayout = this.props.formLayout || formLayout;
    console.log(targetKeys);
    return (
      <Form>
        <FormItem label="账号" {...formItemLayout}>
          {getFieldDecorator('username', {
            initialValue: memberInfo ? memberInfo.username : '',
            rules: [
              { required: true, message: '账号信息不可为空' },
              {
                pattern: /^\S{1,32}$/,
                message: '请输入除空格外的其他字符，并不超过32个字符',
              },
            ],
          })(
            <Input
              placeholder="请输入账号，允许中英文"
              disabled={!!memberInfo}
            />
          )}
        </FormItem>
        <FormItem label="姓名" {...formItemLayout}>
          {getFieldDecorator('fullName', {
            initialValue: memberInfo ? memberInfo.full_name : '',
            rules: [
              {
                pattern: /^\S{1,32}$/,
                message: '请输入除空格外的其他字符，并不超过32个字符',
              },
            ],
          })(
            <Input
              placeholder="请输入姓名，允许中英文"
              disabled={canEdit !== undefined && !canEdit}
            />
          )}
        </FormItem>
        <FormItem label="邮箱" {...formItemLayout}>
          {getFieldDecorator('email', {
            initialValue: memberInfo ? memberInfo.email : '',
            rules: [
              { required: true, message: '邮箱信息不可为空' },
              {
                pattern: /^[a-zA-Z0-9._-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/,
                message: '请输入正确的邮箱格式',
              },
              { max: 32, message: '长度不超过32字符' },
            ],
          })(
            <Input
              placeholder="请输入邮箱"
              disabled={canEdit !== undefined && !canEdit}
            />
          )}
        </FormItem>
        <FormItem label="手机号" {...formItemLayout}>
          {getFieldDecorator('phone', {
            initialValue: memberInfo ? memberInfo.phone : '',
            rules: [
              {
                pattern: /^0?1[35784][0-9][0-9]{8}$/,
                message: '请输入正确的手机号格式',
              },
            ],
          })(
            <Input
              placeholder="请输入手机号"
              disabled={canEdit !== undefined && !canEdit}
            />
          )}
        </FormItem>
        <FormItem label="公司名称" {...formItemLayout}>
          {getFieldDecorator('company', {
            initialValue: memberInfo ? memberInfo.company : '',
            rules: [
              {
                pattern: /^\S{1,128}$/,
                message: '请输入除空格外的其他字符，并不超过128个字符',
              },
            ],
          })(
            <Input
              placeholder="请输入公司名称"
              disabled={canEdit !== undefined && !canEdit}
            />
          )}
        </FormItem>
        <FormItem label="用户角色" {...formItemLayout}>
          {getFieldDecorator('roleId', {
            initialValue: memberInfo ? memberInfo.role_id : 2,
            rules: [{ required: true, message: '请选择用户角色' }],
          })(
            <RadioGroup disabled={canEdit !== undefined}>
              {canEdit !== undefined && <Radio value={1}>Administrator</Radio>}
              <Radio value={2}>Cluster Operator</Radio>
              <Radio value={3}>Cluster Reader</Radio>
            </RadioGroup>
          )}
        </FormItem>
        {window.location.pathname !== '/usercenter/selfinfo' &&
          Cookie.get('em_admin') === 'true' && (
            <FormItem label="集群权限" {...formItemLayout}>
              用户自主添加的集群默认具备该集群操作权限
              {getFieldDecorator('clusterList', {
                initialValue: memberInfo ? targetKeys : [],
              })(
                <Transfer
                  rowKey={(record) => record.id}
                  titles={['未选', '已选']}
                  dataSource={clusterList}
                  showSearch
                  targetKeys={this.state.targetKeys}
                  onChange={this.handleChange}
                  render={this.renderItem}
                />
              )}
            </FormItem>
          )}
        <FormItem>
          {getFieldDecorator('userId', {
            initialValue: memberInfo ? memberInfo.id : undefined,
          })(<Input style={{ display: 'none' }} />)}
        </FormItem>
      </Form>
    );
  }
}
export default Form.create<IProps>()(MemberInfo);
