import * as React from 'react';
import { message, Button, Row, Col } from 'antd';
import { userCenterService } from '@/services';
import MemberInfoForm from '@/pages/userCenter/components/memberInfoForm';
import ResetPasswordModal from '@/pages/userCenter/components/resetPassword';
import { encryptStr, encryptSM } from '@/utils/password';
import './style.scss';

interface State {
  selfInfo: any;
  canEdit: boolean;
  showModal: boolean;
}

export default class SelfInfo extends React.Component<any, State> {
  state: State = {
    selfInfo: {},
    canEdit: false,
    showModal: false,
  };

  private memberInfoForm: any = React.createRef();

  componentDidMount() {
    this.getUserInfo();
  }

  // 获取个人信息
  getUserInfo = () => {
    userCenterService.getLoginedUserInfo().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({
          selfInfo: res.data.info,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 重置密码
  handleResetSubmit = async (value: any) => {
    const publicKeyRes = await userCenterService.getPublicKey();
    if (publicKeyRes.data.code !== 0) {
      return;
    }
    const { encrypt_type, encrypt_public_key } = publicKeyRes.data.data;
    const p = {
      old_password: encrypt_type === 'sm2' ? encryptSM(value.oldPass, encrypt_public_key) : encryptStr(value.oldPass, encrypt_public_key),
      new_password: encrypt_type === 'sm2' ? encryptSM(value.newPass, encrypt_public_key) :encryptStr(value.newPass, encrypt_public_key),
    };
    userCenterService.resetPasswordSelf(p).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        message.success('修改成功');
        this.setState({ showModal: false }, this.getUserInfo);
      } else {
        message.error(res.msg);
      }
    });
  };

  // 重置密码弹窗关闭
  resetPwdModalShow = () => {
    this.setState({ showModal: !this.state.showModal });
  };

  // 保存或开启信息编辑
  handleEditClick = () => {
    const { canEdit } = this.state;
    if (!canEdit) {
      this.setState({ canEdit: true });
      return;
    }
    this.memberInfoForm.props.form.validateFields(
      ['email', 'fullName', 'phone', 'company'],
      (err: any, value: any) => {
        if (err) {
          return;
        }
        const p = {
          ...value,
        };
        userCenterService.motifyUserInfo(p).then((res: any) => {
          res = res.data;
          if (res.code === 0) {
            message.success('修改成功');
            this.setState({
              canEdit: false,
            });
            this.getUserInfo();
          } else {
            message.error(res.msg);
          }
        });
      }
    );
  };

  // 取消
  handleEditCancelClick = () => {
    const { selfInfo } = this.state;
    this.memberInfoForm.props.form.setFieldsValue({
      email: selfInfo.email,
      phone: selfInfo.phone,
      fullName: selfInfo.full_name,
      company: selfInfo.company,
    });
    this.setState({
      canEdit: false,
    });
  };

  render() {
    const { selfInfo, canEdit } = this.state;

    return (
      <div className="selfinfo-container box-shadow-style">
        <MemberInfoForm
          memberInfo={selfInfo}
          canEdit={canEdit}
          wrappedComponentRef={(form) => (this.memberInfoForm = form)}
          formLayout={{
            labelCol: {
              xs: { span: 24 },
              sm: { span: 9 },
            },
            wrapperCol: {
              xs: { span: 24 },
              sm: { span: 10 },
            },
          }}
        />
        <Row>
          <Col span={9}></Col>
          <Col span={15}>
            {canEdit && (
              <Button className="mr-20" onClick={this.handleEditCancelClick}>
                取消
              </Button>
            )}
            <Button className="mr-20" onClick={this.handleEditClick}>
              {canEdit ? '保存' : '编辑资料'}
            </Button>
            <Button type="primary" onClick={this.resetPwdModalShow}>
              重置密码
            </Button>
          </Col>
        </Row>
        <ResetPasswordModal
          visible={this.state.showModal}
          onCancel={this.resetPwdModalShow}
          onSubmit={this.handleResetSubmit}
        />
      </div>
    );
  }
}
