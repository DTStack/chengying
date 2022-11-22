import * as React from 'react';
import { Modal, Button, Steps, Icon, message } from 'antd';
import { userCenterService } from '@/services';
import MemberInfoForm from '@/pages/userCenter/components/memberInfoForm';
import './style.scss';
const ClipboardJS = require('clipboard');

const Step = Steps.Step;

interface IProps {
  visible: boolean;
  onCancel: Function;
  memberInfo: any;
  clusterList: any[];
}

interface IState {
  step: number;
  loading: boolean;
  newUser: {
    username: string;
    password: string;
  };
}

export default class MemberModal extends React.PureComponent<IProps, IState> {
  state: IState = {
    step: 0,
    loading: false,
    newUser: {
      username: '',
      password: '',
    },
  };

  private stepOneForm: any = React.createRef();
  private clipboard = new ClipboardJS('#Copy_Btn');

  componentDidMount() {
    this.clipboard.on('success', (e: any) => {
      message.success('复制成功！');
      this.handleCancel();
    });
    this.clipboard.on('error', (e: any) => {
      message.error('复制失败！');
    });
  }

  componentWillUnmount() {
    this.clipboard.destroy();
  }

  // 取消
  handleCancel = () => {
    this.props.onCancel();
  };

  // 提交表单
  handleSubmit = () => {
    const { memberInfo } = this.props;
    this.stepOneForm.props.form.validateFields((err: any, values: any) => {
      if (!err) {
        this.setState({ loading: true });
        memberInfo ? this.modifyMemberInfo(values) : this.createMember(values);
      }
    });
  };

  // 创建成员
  createMember = (values: any) => {
    userCenterService.regist(values).then((res: any) => {
      this.setState({ loading: false });
      res = res.data;
      if (res.code === 0) {
        this.setState({
          step: this.state.step + 1,
          newUser: res.data,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 修改成员信息
  modifyMemberInfo = (values: any) => {
    userCenterService.modifyInfoByAdmin(values).then((res: any) => {
      this.setState({ loading: false });
      res = res.data;
      if (res.code === 0) {
        message.success('修改成功');
        this.handleCancel();
      } else {
        message.error(res.msg);
      }
    });
  };

  getFooter = () => {
    const { step, newUser } = this.state;
    return (
      <React.Fragment>
        <Button onClick={this.handleCancel}>取消</Button>
        {step === 0 ? (
          <Button type="primary" onClick={this.handleSubmit}>
            确认
          </Button>
        ) : (
          <Button
            type="primary"
            id="Copy_Btn"
            data-clipboard-text={`登录账户:${newUser.username} 登录密码:${newUser.password}`}>
            确认复制
          </Button>
        )}
      </React.Fragment>
    );
  };

  render() {
    const { step, loading, newUser } = this.state;
    const { visible, memberInfo, clusterList } = this.props;
    return (
      <Modal
        className="member-modal"
        visible={visible}
        confirmLoading={loading}
        footer={this.getFooter()}
        onCancel={this.handleCancel}
        width={740}>
        <div>
          {!memberInfo ? (
            <React.Fragment>
              <p className="title">创建账号</p>
              <Steps size="small" current={step}>
                <Step title="编辑信息" />
                <Step title="生成账号" />
              </Steps>
            </React.Fragment>
          ) : (
            <p className="title">编辑账号</p>
          )}
          {step === 0 ? (
            <MemberInfoForm
              memberInfo={memberInfo}
              clusterList={clusterList}
              wrappedComponentRef={(form) => (this.stepOneForm = form)}
            />
          ) : (
            <div className="new-user-info">
              <div className="success-prompt">
                <Icon
                  type="check-circle"
                  style={{ color: '#00A755', fontSize: 50, marginBottom: 10 }}
                />
                <p>创建成功</p>
              </div>
              <div className="info">
                <p>
                  登录地址：
                  {window.location.protocol +
                    '//' +
                    window.location.host +
                    '/#/login'}
                </p>
                <p>登录账户：{newUser.username}</p>
                <p>登录密码：{newUser.password}</p>
              </div>
            </div>
          )}
        </div>
      </Modal>
    );
  }
}
