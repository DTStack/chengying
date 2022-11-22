import * as React from 'react';
import { Input, Button, message, Checkbox, Form, Tooltip, Icon } from 'antd';
import { installGuideService } from '@/services';
import { formLayout } from './constant';

const CheckboxGroup = Checkbox.Group;
const FormItem = Form.Item;

interface State {
  canGoOn: boolean;
  cmdContent: string;
  roles: string[];
  formItemMsg: any;
}

interface IProps {
  clusterInfo: any;
  onCancel?: () => void;
}

export default class StepShell extends React.Component<IProps, State> {
  state: State = {
    canGoOn: false,
    cmdContent: '',
    roles: ['Etcd', 'Control', 'Worker'],
    formItemMsg: {
      status: 'success',
      help: '',
    },
  };

  componentDidMount() {
    this.fetchInstallCmd();
  }

  fetchInstallCmd = () => {
    const { clusterInfo = {} } = this.props;
    const { type, id } = clusterInfo;
    installGuideService
      .getInstallCMD({
        type: type,
        clusterId: id,
        role: this.state.roles.join(','),
      })
      .then((res: any) => {
        if (res.data.code === 0) {
          this.setState({
            cmdContent: res.data.data,
          });
        } else {
          message.error(res.data.msg);
        }
      });
  };

  // 切换角色
  handleChange = (roles: string[]) => {
    this.setState({ roles }, () => {
      this.fetchInstallCmd();
      this.handleCheck();
    });
  };

  // 检测角色是否选择
  handleCheck = () => {
    const { roles } = this.state;
    const formItemMsg = {
      status: 'error',
      help: '请选择主机角色',
    };
    if (roles.length) {
      formItemMsg.status = 'success';
      formItemMsg.help = '';
    }
    this.setState({ formItemMsg });
    return formItemMsg.status;
  };

  // 复制代码，并校验
  handleCopy = () => {
    const { clusterInfo = {} } = this.props;
    const status =
      clusterInfo.type === 'kubernetes' ? this.handleCheck() : 'success';
    if (status === 'success') {
      const e = document.createEvent('MouseEvents');
      e.initEvent('click', true, true);
      document.getElementById('J_CopyBtn').dispatchEvent(e);
    }
  };

  render() {
    const { clusterInfo = {} } = this.props;
    const { formItemMsg, roles } = this.state;
    return (
      <div className="step-shell">
        {clusterInfo.type === 'kubernetes' && (
          <Form>
            <FormItem
              {...formLayout}
              required
              validateStatus={formItemMsg.status}
              help={formItemMsg.help}
              label={
                <Tooltip
                  title={
                    <div>
                      <p>主机角色说明：</p>
                      <p>
                        Etcd角色节点：运行etcd组件，用于存储Kubernetes集群配置数据。
                      </p>
                      <p>
                        Control角色节点：运行Kubernetes主组件（kube-apiserver，kube-scheduler，kube-controller-manager和cloud-controller-manager）。
                      </p>
                      <p>
                        Worker角色节点：运行Kubernetes
                        kubelet，kube-proxy，Container runtime组件。
                      </p>
                    </div>
                  }>
                  <span>主机角色</span>
                  <Icon type="question-circle" />
                </Tooltip>
              }>
              <CheckboxGroup value={roles} onChange={this.handleChange}>
                <Checkbox value="Etcd">Etcd</Checkbox>
                <Checkbox value="Control">Control</Checkbox>
                <Checkbox value="Worker">Worker</Checkbox>
              </CheckboxGroup>
            </FormItem>
          </Form>
        )}
        <div className="clearfix">
          <div className="ant-col ant-form-item-control-wrapper ant-col-xs-24 ant-col-sm-18 ant-col-offset-3">
            <p style={{ fontSize: 12, color: '#333', marginBottom: 8 }}>
              复制下面命令到命令行执行：
            </p>
            <Input.TextArea
              id="J_CMDCode"
              rows={7}
              value={this.state.cmdContent}
            />
            <Button className="mt-10 mb-20" onClick={this.handleCopy}>
              复制代码
            </Button>
            <Button
              style={{ visibility: 'hidden' }}
              id="J_CopyBtn"
              data-clipboard-action="cut"
              data-clipboard-target="#J_CMDCode"></Button>
          </div>
        </div>

        <div style={{ textAlign: 'right', marginTop: '28px' }}>
          <Button className="mr-8" type="default" onClick={this.props.onCancel}>
            取消
          </Button>
          <Button type="primary" onClick={this.props.onCancel}>
            确定
          </Button>
        </div>
      </div>
    );
  }
}
