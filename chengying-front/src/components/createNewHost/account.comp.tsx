import * as React from 'react';
import {
  Form,
  Input,
  Radio,
  Button,
  message,
  Icon,
  Tooltip,
  Select,
  Modal,
  Checkbox,
} from 'antd';
import { createHostService, Service, userCenterService } from '@/services';
import { FormComponentProps } from 'antd/lib/form';
import { formLayout } from './constant';
import { encryptStr, encryptSM } from '@/utils/password';

const FormItem = Form.Item;
const RadioGroup = Radio.Group;
const CheckboxGroup = Checkbox.Group;

interface State {
  showDetermine: boolean;
  canGoOn: boolean;
  fileName: string;
  method: number;
  testLoading: boolean;
  // hostArr: any[];
  showErrorHost: boolean;
  canInstallHost: any[];
  installLoading: boolean;
  hostGroups: any[];
  showResultModel: boolean;
  isCreateGroup: boolean;
  addHostData: any;
  encryptInfo: any;
}

interface Prop extends FormComponentProps {
  afterInstall?: () => void;
  //   refreshHost?: () => void;
  //   refreshGroup?: () => void;
  refresh?: (params) => void;
  clusterInfo?: any;
  onCancel?: () => void;
}
class StepAccount extends React.Component<Prop, State> {
  state: State = {
    canGoOn: false,
    fileName: '',
    method: 2,
    testLoading: false,
    // hostArr: [],
    showErrorHost: false,
    canInstallHost: [],
    installLoading: false,
    hostGroups: [],
    showResultModel: false,
    isCreateGroup: false,
    addHostData: {},
    showDetermine: true,
    encryptInfo: {},
  };

  secretKey = React.createRef();

  componentDidMount() {
    this.loadAllGroups();
    this.getPublicKey();
  }

  getPublicKey = async () => {
    const { data } = await userCenterService.getPublicKey();
    if (data.code !== 0) {
      return;
    }
    this.setState({
      encryptInfo: data.data,
    });
  };

  loadAllGroups = () => {
    const { clusterInfo = {} } = this.props;
    const { type, id } = clusterInfo;
    Service.getClusterhostGroupLists({
      type: type,
      id: id,
    }).then((res: any) => {
      console.log(res);
      res = res.data;
      if (res.code === 0) {
        this.setState({
          hostGroups: res.data || [],
        });
      }
    });
  };

  installHost = () => {
    const { form, clusterInfo = {} } = this.props;
    const {
      encryptInfo: { encrypt_type, encrypt_public_key },
    } = this.state;
    const groupName = form.getFieldValue('group');
    form.validateFields((err: any, values: any) => {
      if (!err) {
        console.log(values);
        const hosts = this.state.canInstallHost.toString();

        this.setState({
          installLoading: true,
        });
        const params = {
          host: hosts,
          port: values.port,
          user: values.user,
          group: values.group,
          cluster_type: clusterInfo.type,
          cluster_id: clusterInfo.id,
          role: values.role ? values.role.join(',') : undefined,
        };
        if (this.state.method === 1) {
          Object.assign(params, {
            type: 'pk',
            pk: values.pk,
          });
        } else {
          Object.assign(params, {
            type: 'pwd',
            password:
              encrypt_type === 'sm2'
                ? encryptSM(values.password, encrypt_public_key)
                : encryptStr(values.password, encrypt_public_key),
          });
        }
        createHostService
          .installHost(params, this.state.method)
          .then((res: any) => {
            res = res.data;
            if (res.code !== 0) {
              message.error(res.msg);
              this.setState({
                installLoading: false,
              });
            } else {
              setTimeout(() => {
                // this.props.refreshGroup();
                // this.props.refreshHost();
                this.props.refresh(groupName);
              }, 1000);
              this.props.afterInstall();
              this.setState({
                installLoading: false,
              });
            }
          });
      }
    });
  };

  testConnection = () => {
    const {
      encryptInfo: { encrypt_type, encrypt_public_key },
    } = this.state;
    this.setState({
      showErrorHost: false,
    });
    const { form } = this.props;
    form.validateFields((err: any, value: any) => {
      if (!err) {
        value.role = value.role ? value.role.join(',') : undefined;
        value.password =
          encrypt_type === 'sm2'
            ? encryptSM(value.password, encrypt_public_key)
            : encryptStr(value.password, encrypt_public_key);
        this.setState({
          testLoading: true,
          canGoOn: false,
          canInstallHost: [],
        });
        const canInstallHost =
          value?.host?.indexOf(',') > -1
            ? value?.host.split(',')
            : [value.host];
        createHostService
          .testConnection(value, this.state.method)
          .then((res: any) => {
            const { code, data, msg } = res.data;
            console.log(data, 'res');
            if (code !== 0) {
              message.error(msg);
              this.setState({
                testLoading: false,
              });
              return;
            }

            this.setState(
              {
                testLoading: false,
                canInstallHost,
              },
              () => {
                const {
                  connectErrorIps,
                  currentClusterExistIp,
                  otherClusterExistIp,
                } = data;
                if (
                  connectErrorIps?.length === 0 &&
                  currentClusterExistIp?.length === 0 &&
                  otherClusterExistIp?.length === 0
                ) {
                  message.success('测试连通性通过');
                  this.setState({
                    showDetermine: false,
                    showResultModel: false,
                  });
                } else {
                  this.setState(() => ({
                    showResultModel: true,
                    addHostData: data,
                  }));
                }
              }
            );
          });
      }
    });
  };

  handleMethodChange = (e: any) => {
    console.log(e);
    this.setState({
      method: e.target.value,
    });
  };

  uploadSecret = () => {
    const { form } = this.props;

    const file = (this.secretKey as any).files[0];
    if (!file) {
      message.warn('请选择上传的秘钥');
      return;
    }
    const flag = file.size / 1024 >= 100;
    if (flag) {
      message.warn('请选择小于100Kb文件上传！');
      return;
    }
    if ((window as any).FileReader) {
      var reader = new FileReader();
      reader.readAsText(file);
      reader.onload = (e: any) => {
        form.setFieldsValue({ pk: e.target.result });
      };
      this.setState({
        fileName: file.name,
      });
      (this.secretKey as any).value = '';
    } else {
      message.warn('浏览器不支持');
    }
  };

  handleGroupBlur(e: any) {
    console.log(e);
  }

  handleGroupBlurChange = (e: any) => {
    console.log(e);
    if (e === 'CREATENEWGROUP') {
      this.setState(
        {
          isCreateGroup: true,
        },
        () => {
          this.props.form.setFieldsValue({
            group: '',
          });
        }
      );
    }
  };

  render() {
    const { clusterInfo = {} } = this.props;
    const { getFieldDecorator } = this.props.form;
    const { addHostData } = this.state;
    return (
      <div className="step-account">
        <Form style={{ marginBottom: 20 }}>
          <FormItem
            style={{ marginBottom: 20 }}
            {...formLayout}
            label={
              <span>
                主机分组
                <Tooltip title="主机分组是对主机的划分，当集群中主机过多时，可对主机进行分组，方便管理。">
                  <Icon type="question-circle" />
                </Tooltip>
              </span>
            }>
            {getFieldDecorator(
              'group',
              {}
            )(
              this.state.isCreateGroup ? (
                <Input placeholder="输入新的分组名称" />
              ) : (
                <Select
                  onChange={this.handleGroupBlurChange}
                  onBlur={this.handleGroupBlur}>
                  <Select.Option value={'CREATENEWGROUP'}>
                    新建主机分组
                  </Select.Option>
                  {this.state.hostGroups &&
                    this.state.hostGroups.map((o: any, i: number) => (
                      <Select.Option value={o} key={i + ''}>
                        {o}
                      </Select.Option>
                    ))}
                </Select>
              )
            )}
          </FormItem>
          <FormItem style={{ marginBottom: 0 }} {...formLayout} label="机器IP">
            {getFieldDecorator('host', {
              rules: [{ required: true, message: '机器IP不能为空' }],
            })(
              <Input.TextArea
                rows={4}
                placeholder="支持输入多个IP及IP段。多个IP用英文都逗号分隔，如172.16.10,171.12.10.9。IP段用英文‘-’连接，如172.16.10-80。"
              />
            )}
          </FormItem>
          {/* {
                        this.state.errorHost.length > 0 &&
                        <span className="error-host-box">
                        下面主机测试失败:
                        {this.state.errorHost.map((o:any) => (
                            <p>{o}</p>
                        ))}
                        </span>
                        } */}
          <FormItem style={{ marginTop: 20 }} {...formLayout} label="SSH接口">
            {getFieldDecorator('port', {
              rules: [{ required: true, message: 'SSH端口不能为空' }],
              initialValue: '22',
            })(<Input />)}
          </FormItem>
          {clusterInfo.type === 'kubernetes' && (
            <FormItem
              {...formLayout}
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
              {getFieldDecorator('role', {
                rules: [{ required: true, message: '主机角色不能为空' }],
                initialValue: ['Etcd', 'Control', 'Worker'],
              })(
                <CheckboxGroup>
                  <Checkbox value="Etcd">Etcd</Checkbox>
                  <Checkbox value="Control">Control</Checkbox>
                  <Checkbox value="Worker">Worker</Checkbox>
                </CheckboxGroup>
              )}
            </FormItem>
          )}
          <FormItem {...formLayout} label="登录方式">
            <RadioGroup
              value={this.state.method}
              onChange={this.handleMethodChange}>
              <Radio value={1}>密钥</Radio>
              <Radio value={2}>密码</Radio>
            </RadioGroup>
          </FormItem>
          <FormItem
            {...formLayout}
            label="用户名"
            extra="需要sudo+NOPASSWD权限">
            {getFieldDecorator('user', {
              rules: [{ required: true, message: '用户名不能为空' }],
              initialValue: 'admin',
            })(<Input />)}
          </FormItem>
          {this.state.method === 2 ? (
            <FormItem {...formLayout} label="密码">
              {getFieldDecorator('password', {
                rules: [{ required: true, message: '密码不能为空' }],
              })(<Input.Password />)}
            </FormItem>
          ) : (
            <FormItem {...formLayout} label="密钥">
              {getFieldDecorator('pk', {
                rules: [{ required: true, message: '密钥不可为空!' }],
              })(<Input style={{ display: 'none' }} />)}
              <label className="ant-btn">
                上传秘钥文件
                <input
                  ref={(input: any) => (this.secretKey = input)}
                  type="file"
                  onChange={this.uploadSecret}
                  style={{ display: 'none' }}
                />
              </label>
              <span>{this.state.fileName}</span>
            </FormItem>
          )}
        </Form>
        {this.state.showResultModel ? (
          <div style={{ marginBottom: 20 }}>
            <p>
              <Icon
                type="close-circle"
                style={{ color: 'red', fontSize: 14, marginRight: 6 }}
                theme="filled"
              />
              异常IP列表:
            </p>
            <div className="error-desc">
              {addHostData?.connectErrorIps?.length ? (
                <p>
                  - 连通失败IP：
                  {addHostData?.connectErrorIps?.join(',') || '无'}
                </p>
              ) : (
                ''
              )}
              {addHostData?.currentClusterExistIp?.length ? (
                <p>
                  - 以下IP在当前集群已存在：
                  {addHostData?.currentClusterExistIp?.join(',') || '无'}
                </p>
              ) : (
                ''
              )}
              {addHostData?.otherClusterExistIp?.length ? (
                <p>
                  - 以下IP在其他集群已存在：
                  {addHostData?.otherClusterExistIp?.join(',') || '无'}
                </p>
              ) : (
                ''
              )}
            </div>
          </div>
        ) : (
          ''
        )}
        <div className="btn-container">
          <Button
            loading={this.state.testLoading}
            type="primary"
            onClick={this.testConnection}>
            测试连通性
          </Button>
          <div>
            <Button
              onClick={this.props.onCancel}
              type="default"
              className="mr-8"
              // loading={this.state.installLoading}
              // disabled={this.state.canInstallHost.length === 0}
            >
              取消
            </Button>
            <Button
              onClick={this.installHost}
              type="primary"
              loading={this.state.installLoading}
              disabled={this.state.showResultModel || this.state.showDetermine}>
              确定
            </Button>
          </div>
        </div>
        {/* <Modal
          className="host-test-modal"
          onCancel={() => this.setState({ showResultModel: false })}
          visible={this.state.showResultModel}
          title={
            <span style={{ color: '#333333', fontSize: 14, marginLeft: 9 }}>
              <Icon
                type="close-circle"
                theme="filled"
                style={{ color: '#E6432C', fontSize: 14, marginRight: 6 }}
              />
              主机测试失败列表
            </span>
          }
          footer={
            <span className="btn-box">
              <Button
                className="btn"
                onClick={() => this.setState({ showResultModel: false })}>
                取消
              </Button>
              <Button
                type="primary"
                className="btn"
                onClick={() => this.setState({ showResultModel: false })}>
                关闭
              </Button>
            </span>
          }>
          {this.state.errorHost.map((o: any) => {
            return (
              <p style={{ marginLeft: 24, marginBottom: 8 }} key={o}>
                {o}
              </p>
            );
          })}
        </Modal> */}
      </div>
    );
  }
}
export default Form.create<Prop>()(StepAccount);
