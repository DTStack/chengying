import * as React from 'react';
import { connect } from 'react-redux';
import {
  hostAddAction,
  installHost,
  checkInstall,
} from '@/actions/addHostAction';
import { Table, Icon, message, Popover, Progress, Tooltip } from 'antd';
type progressStatus = 'exception' | 'active' | 'success';

const progress = (status: any) => {
  switch (status) {
    case 0:
      return {
        percent: 0,
        title: '等待安装',
        color: 'host-status-wait-install',
      };
    case -1:
      return {
        percent: 30,
        status: 'exception' as progressStatus,
        title: '管控安装失败',
        color: 'host-status-control-failed',
      };
    case 1:
      return {
        percent: 30,
        status: 'active' as progressStatus,
        title: '管控安装成功',
        color: 'host-status-control-successed',
      };
    case -2:
      return {
        percent: 60,
        status: 'exception' as progressStatus,
        title: 'script安装失败',
        color: 'host-status-script-failed',
      };
    case 2:
      return {
        percent: 60,
        status: 'active' as progressStatus,
        title: 'script安装成功',
        color: 'host-status-script-successed',
      };
    case -3:
      return {
        percent: 100,
        status: 'exception' as progressStatus,
        title: '主机初始化失败',
        color: 'host-status-init-failed',
      };
    case 3:
      return {
        percent: 100,
        status: 'success' as progressStatus,
        title: '主机初始化成功',
        color: 'host-status-init-successed',
      };
    default:
      return { percent: 100, status: 'exception' as progressStatus };
  }
};

// interface State { }

export class StepThree extends React.Component<any> {
  constructor(props: any) {
    super(props);
    // this.reInstall = props.reInstall
  }

  private timeoutObject: any = {};
  private currentHostArr: any[] = [];

  componentWillUnmount() {
    // clear timeout
    const timeoutObject = this.timeoutObject;
    Object.keys(timeoutObject).forEach((objectKey, index) => {
      clearTimeout(timeoutObject[objectKey]);
    });
  }

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
      // this.setState({ installMsg }) //验证连接成功后，更新该组件state， 使page3页面能调用receiveProps()更新
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

  handleSubmit(host: any) {
    // debugger;
    const { hostArr } = this.props;
    // console.log(this.currentHostArr);
    if (host) {
      let hitem = null;
      for (const h of hostArr) {
        if (h.host === host) {
          hitem = h;
        }
      }
      if (hitem) {
        this.installProgress(hitem);
      }
    } else {
      this.currentHostArr.forEach((item: any, index: any) => {
        this.installProgress(item);
      });
    }
  }

  render() {
    // let datasource = this.state.installMsg;
    const { installMsg } = this.props;
    console.log('install msg:', this.props.installMsg);
    const col = [
      {
        title: 'IP',
        key: 'ip',
        width: '25%',
        dataIndex: 'host',
      },
      {
        title: '安装进度',
        width: '40%',
        key: 'progress',
        dataIndex: 'progress',
        render: (text: any, records: any) => {
          // console.log(records)
          const progressResult = progress(records.status);
          return (
            <Progress
              className="c-host-add-progress"
              percent={progressResult.percent}
              strokeWidth={5}
              status={progressResult.status as any}
            />
          );
        },
      },
      {
        title: '状态',
        width: '25%',
        dataIndex: 'status',
        key: 'status',
        render: (text: any, records: any) => {
          const statusError =
            records.status === -1 ||
            records.status === -2 ||
            records.status === -3;
          return (
            <span className={`${progress(records.status).color} install-opt`}>
              {progress(records.status).title}
              {statusError && records.msg && (
                <Tooltip placement="bottom" title={records.msg}>
                  <Icon className="icon-tips" type="question-circle" />
                </Tooltip>
              )}
              {records.aid && statusError && (
                <a
                  style={{ marginLeft: 10 }}
                  onClick={() => this.reInstall(records)}>
                  重试
                </a>
              )}
            </span>
          );
        },
      },
      {
        title: '详情',
        dataIndex: 'detail',
        key: 'detail',
        render: (text: any, records: any) => {
          return (
            <Popover
              placement="left"
              content={
                <div style={{ maxHeight: 300, overflowY: 'auto' }}>{text}</div>
              }>
              {
                <a href="#" style={{ color: '#aaf' }}>
                  查看详情
                </a>
              }
            </Popover>
          );
        },
      },
    ];
    return (
      <div className="steps-main access-table">
        <Table
          scroll={{ y: 500 }}
          columns={col}
          dataSource={installMsg}
          pagination={false}
        />
      </div>
    );
  }
}

const mapStateToProps = (state: any) => ({
  current: state.hostAdd.current,
  disabled: state.hostAdd.disabled,
  forms: state.hostAdd.forms,
  hostArr: state.hostAdd.hostArr,
  installMsg: state.hostAdd.installMsg,
});

const mapDispatchToProps = {
  jumpPage: (v: any) => hostAddAction.jumpPage(v),
  updateInstallMsg: (v: any) => hostAddAction.updateInstallMsg(v),
  updateHostArr: (v: any) => hostAddAction.updateHostArr(v),
  setDisabled: (v: any) => hostAddAction.setDisabled(v),
  getNewState: (v: any) => hostAddAction.getNewState(v),
  resetState: () => hostAddAction.resetState(),
};

export default connect(mapStateToProps, mapDispatchToProps)(StepThree);
