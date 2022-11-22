import * as React from 'react';
import { connect } from 'react-redux';
import {
  hostAddAction,
  connectHost,
  installHost,
  checkInstall,
} from '@/actions/addHostAction';

import * as Http from '@/utils/http';
import { AppStoreTypes } from '@/stores';

import StepOne from './StepOne';
import { StepTwo } from './StepTwo';
import { StepThree } from './StepThree';
import { Button, Input, Steps, message, Tabs } from 'antd';

const Step = Steps.Step;
const TabPane = Tabs.TabPane;
const { TextArea } = Input;
const steps = [
  {
    title: '添加机器',
  },
  {
    title: '安装内容说明',
  },
  {
    title: '执行安装',
  },
];

// let count = 6; //轮询次数

const mapStateToProps = (state: AppStoreTypes) => ({
  current: state.addHostStore.current,
  disabled: state.addHostStore.disabled,
  forms: state.addHostStore.forms,
  hostArr: state.addHostStore.hostArr,
  installMsg: state.addHostStore.installMsg,
});

const mapDispatchToProps = {
  jumpPage: (v: any) => hostAddAction.jumpPage(v),
  updateInstallMsg: (v: any) => hostAddAction.updateInstallMsg(v),
  updateHostArr: (v: any) => hostAddAction.updateHostArr(v),
  setDisabled: (v: any) => hostAddAction.setDisabled(v),
  getNewState: (v: any) => hostAddAction.getNewState(v),
  resetState: () => hostAddAction.resetState(),
};

interface Prop {
  current: any;
  disabled: any;
  foprms: any;
  hostArr: any;
  installMsg: any;
  jumpPage: (e: any) => void;
  updateInstallMsg: (e: any) => void;
  updateHostArr: (e: any) => void;
  setDisabled: (v: any) => void;
  getNewState: (e: any) => void;
  resetState: () => void;
}

interface State {
  installMsg: any[];
  shellCMD: string;
  cmdVisible: boolean;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
class HostAdd extends React.Component<Prop, State> {
  constructor(props: any) {
    super(props);
    this.timeoutObject = {}; // 轮询安装状态
    this.currentHostArr = []; // 第一页新添加的主机字段
    // this.updateInstallItemStatus = false
    this.state = {
      installMsg: [],
      shellCMD: '',
      cmdVisible: false,
    };
  }

  private timeoutObject: any = {};
  private currentHostArr: any = [];
  private HostAddForm: any = React.createRef();
  // private updateInstallItemStatus: boolean = false;

  componentDidMount() {
    const ctx = this;
    Object.prototype.hasOwnProperty.call(sessionStorage, 'batch_installInfo') &&
      delete sessionStorage.batch_installInfo;
    Http.get('/api/v2/agent/install/installCmd', {}).then((res: any) => {
      if (res.code === 0) {
        ctx.setState({
          shellCMD: res.data,
        });
      } else {
        message.error(res.msg);
      }
    });
    // console.log(this.refs);
  }

  componentWillUnmount() {
    this.props.resetState();
    Object.prototype.hasOwnProperty.call(sessionStorage, 'batch_installInfo') &&
      delete sessionStorage.batch_installInfo;
    // clear timeout
    const timeoutObject = this.timeoutObject;
    Object.keys(timeoutObject).forEach((objectKey, index) => {
      clearTimeout(timeoutObject[objectKey]);
    });
  }

  goValid() {
    const { updateInstallMsg } = this.props;
    const formdata = this.HostAddForm.getFieldsValue();
    // 验证是否表单必填字段是否全部填完
    for (var key in formdata) {
      if (formdata[key] === undefined || formdata[key] === '') {
        return;
      }
    }
    let installVelocity: any[] = [];
    let hostArr = [];
    formdata.host.indexOf(',') > -1
      ? (hostArr = formdata.host.split(','))
      : hostArr.push(formdata.host);
    // 重组hostArr 并对多个主机同时请求
    hostArr.forEach((v: any, i: any) => {
      const V = {
        host: v,
        port: formdata.port,
        user: formdata.user,
        password: formdata.password || '',
        pk: formdata.pk || '',
      };
      const current = {
        host: v,
        percent: 0,
        status: 0,
        msg: '',
      };
      installVelocity.push(current);
      this.linkProgress(V, current);
    });
    installVelocity = this.concatSessionStore(installVelocity);
    updateInstallMsg(installVelocity);
  }

  async linkProgress(hostArrItem: any, current: any) {
    const { setDisabled, updateHostArr } = this.props;
    hostArrItem = this.pkOrPwd(hostArrItem);
    connectHost(hostArrItem, (res: any) => {
      if (res) {
        const { code, msg } = res;
        if (code === 0 || msg === 'ok') {
          message.success('ip连通性验证通过！');
          updateHostArr(hostArrItem.item);
          this.setState({ installMsg: [...this.state.installMsg, current] }); // 验证连接成功后，更新该组件state， 使page3页面能调用receiveProps()更新
          this.currentHostArr.push(hostArrItem.item);
          setDisabled(false);
        } else {
          message.error(hostArrItem.host + '连接失败' + res.msg);
          setDisabled(true);
        }
      }
    });
  }

  handleSubmit(host: any) {
    // debugger;
    // console.log(this.currentHostArr);
    if (host) {
      let hitem = null;
      for (const h of this.currentHostArr) {
        if (h.host === host) {
          hitem = h;
        }
      }
      if (hitem) {
        this.installProgress(hitem);
      }
    } else {
      this.currentHostArr.forEach((item: any) => {
        this.installProgress(item);
      });
    }
  }

  async installProgress(item: any) {
    const { installMsg, updateInstallMsg } = this.props;
    const v = this.pkOrPwd(item);
    const res = await installHost(v);

    if (res) {
      const { code, msg, data } = res;
      if (code === 0 || msg === 'ok') {
        // debugger;
        installMsg.forEach((v: any) => {
          if (v.host === item.host) {
            v.aid = data.aid;
            v.detail = `ip: ${item.host},安装信息：${data.result}`;
          }
        });
        // console.log('msg:', installMsg);
        updateInstallMsg(installMsg);
        this.installCheck(data.aid);
      } else {
        message.error(item.host + '安装失败' + msg);
      }
    }

    // console.log('install', this.props.installMsg)
  }

  async installCheck(aid: any) {
    const { installMsg, updateInstallMsg } = this.props;
    if (aid == null || aid === undefined || aid === '' || !aid) {
      return;
    }
    const res = await checkInstall(aid);
    const { data, code } = res;
    if (code === 0) {
      installMsg.forEach((v: any, i: any) => {
        if (v.aid === aid) {
          if (data.status === 1 || data.status === 2 || data.status === 0) {
            // 安装未完成，轮询
            this.timeoutObject[aid] = setTimeout(() => {
              // count--
              this.installCheck(aid);
              // count < 0 && (this.timeoutObject[aid] = null)
            }, 2000);
          }
          v.status = data.status || 1;
          v.detail = `${data.status_msg}`;
        }
      });

      updateInstallMsg(installMsg);
      this.setState({ installMsg }); // 验证连接成功后，更新该组件state， 使page3页面能调用receiveProps()更新
      // :[...this.state.installMsg,...installMsg]
    } else {
      // 失败
      installMsg.forEach((v: any, i: any) => {
        if (v.aid === aid) {
          v.status = data.result_code || -1;
          v.detail = res.msg;
        }
      });
      message.error(res.msg);
    }
  }

  pkOrPwd(item: any) {
    let type;
    if (item.pk) {
      type = 'pk';
      delete item.password;
    } else {
      type = 'pwd';
      delete item.pk;
    }
    return { type: type, item: item };
  }

  concatSessionStore = (installVelocity: any) => {
    if (
      Object.prototype.hasOwnProperty.call(
        sessionStorage,
        'batch_inaaaawdstallInfo'
      )
    ) {
      const store =
        JSON.parse(sessionStorage.getItem('batch_installInfo')) || [];
      installVelocity = installVelocity.concat(store);
    }
    return installVelocity;
  };

  reInstall = (record: any) => {
    // debugger;
    const { installMsg, updateInstallMsg } = this.props;
    let current = -1;
    installMsg.forEach((item: any, index: any) => {
      if (item.aid === record.aid) {
        current = index;
      }
    });
    if (current !== -1) {
      installMsg[current] = {
        ...installMsg[current],
        status: 0,
        detail: '等待安装',
      };
      updateInstallMsg(installMsg);
      // this.installProgress(installMsg[current])
      this.handleSubmit(installMsg[current].host);
    }
  };

  backToHostList = () => {
    window.location.href = '/host';
  };
  configInitialize = () => {
    const { jumpPage, installMsg } = this.props;
    this.currentHostArr = [];
    const current = 0;
    sessionStorage.setItem('batch_installInfo', JSON.stringify(installMsg));
    jumpPage(current);
  };

  prev() {
    const { current, jumpPage } = this.props;
    if (current === 0) {
      window.location.href = '/host';
    } else {
      jumpPage(current - 1);
    }
  }

  next() {
    const { current, jumpPage, setDisabled } = this.props;
    if (current === 1) {
      this.handleSubmit(null);
      setDisabled(true);
    }
    jumpPage(current + 1);
  }

  showCmdModal() {
    this.setState({
      cmdVisible: true,
    });
  }

  handleCloseNoCmdModal() {
    this.setState({
      cmdVisible: false,
    });
  }
  // componentWillUnmount() {
  //     this.props.resetState();
  // }

  render() {
    const btnStyle = {
      margin: '0 10px',
    };
    const { current, disabled } = this.props;
    console.log('im:', this.state.shellCMD);
    const params = { ...this.props, reInstall: this.reInstall };
    return (
      <div style={{ width: '100%', minHeight: '700px', overflow: 'auto' }}>
        <div
          style={{
            margin: '15px 170px 170px 170px',
            width: '80%',
            paddingBottom: '50px',
            backgroundColor: '#fff',
          }}>
          <div style={{ padding: '40px 80px' }}>
            <Steps current={current}>
              {steps.map((item) => (
                <Step key={item.title} title={item.title} />
              ))}
            </Steps>
            <div
              style={{
                margin: '20px 0',
                padding: '20px',
                width: '100%',
                backgroundColor: 'rgb(250,250,250)',
              }}>
              {current === 0 && (
                <Tabs defaultActiveKey="1">
                  <TabPane tab="账号接入" key="1">
                    <StepOne
                      ref={(form: any) => (this.HostAddForm = form)}
                      {...this.props}
                    />
                  </TabPane>
                  <TabPane tab="命令行接入" key="2">
                    <p>复制下面命令到命令行执行：</p>
                    <TextArea
                      id="J_CMDCode"
                      rows={4}
                      value={this.state.shellCMD}></TextArea>
                    <Button
                      id="J_CopyBtn"
                      style={{ marginTop: 15 }}
                      type="default"
                      data-clipboard-action="cut"
                      data-clipboard-target="#J_CMDCode">
                      复制代码
                    </Button>
                  </TabPane>
                </Tabs>
              )}
              {current === 1 && <StepTwo />}
              {current === 2 && <StepThree {...params} />}
            </div>

            <div
              className="steps-action"
              style={{
                marginRight: 0,
                padding: '0 20px',
                width: '100%',
                textAlign: 'right',
              }}>
              {current !== 2 && (
                <Button style={btnStyle} onClick={() => this.prev()}>
                  返回{' '}
                </Button>
              )}
              {current === 0 && (
                <Button
                  style={btnStyle}
                  type="primary"
                  onClick={() => this.goValid()}>
                  验证连通性
                </Button>
              )}
              {current !== 2 && (
                <Button
                  style={btnStyle}
                  type="primary"
                  disabled={disabled}
                  onClick={() => this.next()}>
                  下一步
                </Button>
              )}
              {current === 2 && (
                <Button style={btnStyle} onClick={this.backToHostList}>
                  进入主机列表
                </Button>
              )}
              {current === 2 && (
                <Button
                  style={btnStyle}
                  type="primary"
                  onClick={this.configInitialize}>
                  继续配置
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default HostAdd;
