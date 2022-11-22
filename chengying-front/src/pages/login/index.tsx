import * as React from 'react';
import { Form, Input, message, Button } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { bindActionCreators, Dispatch } from 'redux';
import * as ServiceActions from '@/actions/serviceAction';
import { connect } from 'react-redux';
import { userCenterService, Service, servicePageService } from '@/services';
import { navData } from '@/constants/navData';
import { encryptStr, encryptSM } from '@/utils/password';
import * as Http from '@/utils/http';

import './style.scss';

import * as Cookie from 'js-cookie';

declare var APP: any;

const FormItem = Form.Item;
interface State {
  validCode: string;
  needValidCode: boolean;
}
interface IProps extends FormComponentProps {
  location: any;
  history: any;
  actions: any;
}
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, ServiceActions), dispatch),
});
class Login extends React.Component<IProps, State> {
  componentDidMount() {
    this.initValidCode();
  }

  state: State = {
    validCode: '',
    needValidCode: true,
  };

  initValidCode = () => {
    userCenterService.getValidCode().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({
          validCode: res.data,
          needValidCode: res.data.length > 0,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  validHasDeployedProduct = (cb: Function) => {
    const { history } = this.props;
    Service.getParentProductList().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        if (res.data.length === 0) {
          history.push(
            '/deploycenter/cluster/list?from=' + this.props.location.pathname
          );
        } else {
          // cb()
          const promises: any[] = [];

          res.data.map((o: any) => {
            promises.push(
              new Promise((resolve: any, reject: any) => {
                servicePageService
                  .getProductList({ limit: 0, parentProductName: o })
                  .then((res: any) => {
                    res = res.data;
                    if (res.code === 0) {
                      resolve(res.data.list);
                    } else {
                      message.error(res.msg);
                    }
                  });
              })
            );
          });
          let flag = true;
          Promise.all(promises).then((res: any) => {
            console.log(res);
            res.map((o: any) => {
              o.map((p: any) => {
                if (p.is_current_version === 1) {
                  flag = false;
                }
              });
            });

            console.log('结果是' + flag);
            if (flag) {
              history.push(
                '/deploycenter/cluster/list?from=' +
                  this.props.location.pathname
              );
            } else {
              cb();
            }
          });
        }
      } else {
        message.error(res.msg);
      }
    });
  };

  validCodeValue = (
    p: { captchaId: string; captchaSolution: string },
    f: Function
  ) => {
    userCenterService.checkValidCode(p).then((res: any) => f(res));
  };

  goLogin = async (p: { username: string; password: string }) => {
    const publicKeyRes = await userCenterService.getPublicKey();
    if (publicKeyRes.data.code !== 0) {
      return;
    }
    const { encrypt_type, encrypt_public_key } = publicKeyRes.data.data;
    userCenterService
      .login({
        username: p.username,
        password: encrypt_type === 'sm2' ? encryptSM(p.password, encrypt_public_key) : encryptStr(p.password, encrypt_public_key),
      })
      .then((res: any) => {
        Cookie.set('em_token', res.headers.em_token);
        Cookie.set('em_username', p.username);
        Cookie.set('em_admin', res.headers.em_admin);

        res = res.data;
        if (res.code === 0) {
          this.getAuthorityRouter();
          this.validHasDeployedProduct(() => {
            window.location.href = '/';
          });
          this.getRestartService();
        } else {
          message.error(res.msg);
          this.initValidCode();
        }
      });
  };

  // 获取需要依赖组件重新配置服务列表——EM提醒
  getRestartService = () => {
    Http.get('/api/v2/cluster/restartServices', {}).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        console.log(this.props.actions);
        this.props.actions?.setResartServiceList({
          count: res?.data?.count || 0,
          list: res?.data?.list || [],
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 根据code拿到路由
  authorityRouterFilter = (routers: any[], list: any[], codeList) => {
    list.forEach((item: any) => {
      if (codeList.includes(item.code)) {
        routers.push(...item.routers);
      }
      if (item.children.length) {
        this.authorityRouterFilter(routers, item.children, codeList);
      }
    });
    return routers;
  };

  // 获取权限路由
  getAuthorityRouter = () => {
    userCenterService.getRoleCodes().then((response: any) => {
      const { code, data } = response.data;
      if (code === 0) {
        const routers = [];
        this.authorityRouterFilter(routers, navData, data);
        localStorage.setItem('authorityRouter', JSON.stringify(routers));
      }
    });
  };

  handleLogin = () => {
    const { form } = this.props;
    form.validateFields((err: any, value: any) => {
      if (!err) {
        const p = {
          username: value.username,
          password: value.password,
        };
        if (this.state.needValidCode) {
          this.validCodeValue(
            {
              captchaId: this.state.validCode,
              captchaSolution: value.validCode,
            },
            (res: any) => {
              res = res.data;
              if (res.code !== 0) {
                message.error(res.msg);
                this.initValidCode();
                form.setFieldsValue([{ validCode: undefined }]);
                return;
              }
              this.goLogin(p);
            }
          );
        } else {
          this.goLogin(p);
        }
      }
    });
  };

  render() {
    const { getFieldDecorator } = this.props.form;
    return (
      <div
        className="login-container"
        style={{ height: document.body.clientHeight }}>
        <div className="header">
          <a href="/" className="header-logo-wrapper">
            <img src={require('public/imgs/logo_chengying@2x.png')} />
            <span className="header-logo-name">ChengYing</span>
          </a>
        </div>
        <div style={{ height: document.body.clientHeight - 40, width: '100%' }}>
          <img className="bg-img" src={require('public/imgs/BG.png')} />
          <div
            className="content"
            style={{ minHeight: this.state.needValidCode ? 330 : 275 }}>
            <p style={{ fontSize: 18, color: '#333333', marginBottom: 31 }}>
              欢迎登录ChengYing
            </p>
            <Form>
              <FormItem>
                {getFieldDecorator('username', {
                  rules: [
                    { required: true, message: '用户名不能为空' },
                    // {
                    //     pattern: /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/,
                    //     message: '请输入正确的账户名称'
                    // }
                  ],
                })(
                  <Input
                    style={{ width: 330, height: 40 }}
                    placeholder="请输入账号信息"
                  />
                )}
              </FormItem>
              <FormItem>
                {getFieldDecorator('password', {
                  rules: [{ required: true, message: '密码不能为空' }],
                })(
                  <Input
                    style={{ width: 330, height: 40 }}
                    placeholder="请输入登录密码"
                    type="password"
                  />
                )}
              </FormItem>
              {this.state.needValidCode && (
                <FormItem>
                  {getFieldDecorator('validCode', {
                    rules: [{ required: true, message: '验证码不能为空' }],
                  })(
                    <span
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                      }}>
                      <Input
                        style={{ width: 200, height: 40 }}
                        placeholder="验证码"
                      />
                      <img
                        title="点击更换验证码"
                        style={{ width: 120, height: 36, cursor: 'pointer' }}
                        onClick={this.initValidCode}
                        src={
                          this.state.validCode
                            ? `/api/v2/user/showCaptcha/${this.state.validCode}.png`
                            : ''
                        }
                      />
                    </span>
                  )}
                </FormItem>
              )}

              <Button
                htmlType="submit"
                onClick={this.handleLogin}
                type="primary"
                style={{
                  width: '100%',
                  height: 40,
                  fontSize: 14,
                  marginBottom: 10,
                }}>
                登录
              </Button>
            </Form>
            {/* <p
                            style={{
                                color: "#666666",
                                fontSize: 12,
                                textAlign: "right",
                                margin: 0
                            }}
                        >
                            没有账号？<a href="regist">免费注册</a>
                        </p> */}
          </div>
        </div>
        <div
          className="version"
          style={{
            position: 'fixed',
            bottom: 10,
            left: '50%',
            transform: 'translate(-50%,0)',
          }}>
          ChengYing@V{APP.VERSION}
        </div>
      </div>
    );
  }
}
export default connect(undefined, mapDispatchToProps)(Form.create()(Login));
