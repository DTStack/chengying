import * as React from 'react';
import { Input, message, Modal, Transfer } from 'antd';
import { deployService } from '@/services';
import './style.scss';

interface State {
  originIps: any;
  existIPs: any;
  leftIPs: any[];
  sourceIps: any[];
  targetKeys: any[];
  selectedKeys: any[];
  ipStr: string;
  max: number;
  showModal: boolean;
}
export class IPEditPanel extends React.Component<any, State> {
  constructor(props: any) {
    super(props);
    let mx = 1;
    if (!this.props.instance || this.props.instance.MaxReplica === '') {
      mx = 9999;
    } else {
      mx = this.props.instance.MaxReplica;
    }
    const ext = [];
    const ips =
      this.props.data.IP && this.props.data.IP.length ? this.props.data.IP : [];
    for (const i of ips) {
      ext.push({
        key: i,
        title: i,
        description: i,
        disable: false,
      });
    }
    this.state = {
      originIps: ext,
      existIPs: ext,
      leftIPs: [],
      sourceIps: [],
      targetKeys: ips || [],
      selectedKeys: [],
      ipStr: '',
      max: mx,
      showModal: false,
    };
  }

  componentDidMount() {
    this.getLeftIps();
  }

  UNSAFE_componentWillReceiveProps(nextProps: any) {
    const { instance, data } = nextProps;
    let mx = 1;
    if (!instance || instance.MaxReplica === '') {
      mx = 9999;
    } else {
      mx = parseInt(instance.MaxReplica);
    }
    const ext = [];
    const ips =
      this.props.data.IP && this.props.data.IP.length ? this.props.data.IP : [];
    for (const i of ips) {
      ext.push({
        key: i,
        title: i,
        description: i,
        disable: false,
      });
    }

    this.setState({
      originIps: ext,
      existIPs: ext,
      targetKeys: ips,
      ipStr: data,
      max: mx,
    });
  }

  getLeftIps() {
    const ctx = this;
    // let { existIPs } = _this.state;
    const productName = ctx.props.pname;
    const serviceName = ctx.props.sname;
    deployService
      .getLeftIp({
        productName: productName,
        serviceName: serviceName,
      })
      .then((data: any) => {
        data = data.data;
        const lips = [];
        if (data.code === 0) {
          if (data.data.hosts) {
            for (const h of data.data.hosts) {
              let exist = false;
              for (const p of ctx.state.existIPs) {
                if (p === h.ip) {
                  exist = true;
                }
              }
              if (!exist) {
                lips.push({
                  key: h.ip,
                  title: h.ip,
                  description: h.ip,
                  disable: false,
                });
              }
            }
            ctx.setState({
              sourceIps: ctx.state.existIPs.concat(lips),
              leftIPs: lips,
            });
          }
        } else {
          message.error(data.msg);
        }
      });
  }

  handleSetIps = () => {
    // debugger;
    const { targetKeys, max } = this.state;
    // let total = targetKeys.concat(s);
    if (targetKeys.length > max) {
      message.error('IP数量限制' + max + '，目前超出限制！');
    } else {
      this.props.setip(this.props.pindex, targetKeys.join(','));
      this.setState({
        showModal: false,
      });
    }
  };

  handleShowEditModal = () => {
    const { originIps } = this.state;
    const tk = [];
    for (const t of originIps) {
      tk.push(t.key);
    }
    this.setState({
      showModal: true,
      targetKeys: tk,
    });
  };

  handleHideEditModal = () => {
    const { originIps } = this.state;
    const tk = [];
    for (const t of originIps) {
      tk.push(t.key);
    }
    this.setState({
      showModal: false,
      existIPs: originIps,
      targetKeys: tk,
    });
  };

  handleTransferIpChange = (t: any, d: any, m: any) => {
    const ext = [];
    for (const i of t) {
      ext.push({
        key: i,
        title: i,
        description: i,
        disable: false,
      });
    }
    this.setState({
      targetKeys: t,
      existIPs: ext,
    });
  };

  handleTransferSelectChange = (s: any, t: any) => {
    const { targetKeys, max } = this.state;
    const total = targetKeys.concat(s);
    if (total.length > max) {
      message.error('IP数量限制' + max + '，目前超出限制！');
    }
  };

  handleIpInputChange(e: any) {
    const { updateip } = this.props;
    updateip(this.props.pindex, e.target.value);
  }

  updateIpInputChange(e: any) {
    this.props.setip(this.props.pindex, e.target.value);
  }

  render() {
    const { instance, data } = this.props;
    if (instance) {
      return (
        <div className="edit-panel" style={{ borderBottom: 0 }}>
          <div
            className="mod-title"
            style={
              instance.UseCloud ? { display: 'none' } : { display: 'block' }
            }>
            <a onClick={this.handleShowEditModal}>编辑ip地址</a>：
            <span style={{ padding: 10, lineHeight: '20px' }}>
              {data.IP ? data.IP.join(',') : ''}
            </span>
          </div>
          <div
            className="mod-title"
            style={
              instance.UseCloud ? { display: 'block' } : { display: 'none' }
            }>
            <span style={{ width: '8%', display: 'inline-block' }}>
              编辑ip地址:{data.IP}{' '}
            </span>
            <Input
              style={{ width: '90%' }}
              placeholder="多个ip用逗号分开"
              value={data.IP ? data.IP.join(',') : ''}
              onChange={this.handleIpInputChange.bind(this)}
              onBlur={this.updateIpInputChange.bind(this)}
            />
          </div>
          <Modal
            visible={this.state.showModal}
            width={476}
            onOk={this.handleSetIps}
            onCancel={this.handleHideEditModal}
            className="ip-transfer-modal">
            <Transfer
              titles={['未选', '已选']}
              dataSource={this.state.sourceIps}
              showSearch
              targetKeys={this.state.targetKeys}
              onChange={this.handleTransferIpChange}
              onSelectChange={this.handleTransferSelectChange}
              render={(item) => item.title}
            />
          </Modal>
        </div>
      );
    } else {
      return (
        <div className="edit-panel" style={{ borderBottom: 0 }}>
          <div className="mod-title">
            <span style={{ width: '8%', display: 'inline-block' }}>
              编辑ip地址:{' '}
            </span>
            <Input
              style={{ width: '90%' }}
              placeholder="多个ip用逗号分开"
              value={data.IP ? data.IP.join(',') : ''}
              onChange={this.handleIpInputChange.bind(this)}
              onBlur={this.updateIpInputChange.bind(this)}
            />
          </div>
        </div>
      );
    }
  }
}
