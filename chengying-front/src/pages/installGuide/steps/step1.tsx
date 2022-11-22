import * as React from 'react';
import * as Cookie from 'js-cookie';
import {
  Select,
  Icon,
  Form,
  message,
  Radio,
  Tooltip,
  Button,
  Input,
} from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { RadioChangeEvent } from 'antd/lib/radio';
import { formItemCenterLayout } from '@/constants/formLayout';

const RadioGroup = Radio.Group;
const Option = Select.Option;
const FormItem = Form.Item;

interface IProps extends FormComponentProps {
  installGuideProp: any;
  actions: any;
  location: any;
}
interface State {
  clusterList: any[];
  clusterName: string;
  isNewSpace: boolean;
  mode: number;
}

class StepOne extends React.Component<IProps, State> {
  constructor(props: IProps) {
    super(props);
    this.state = {
      clusterList: [],
      clusterName: '',
      isNewSpace: false,
      mode: 0,
    };
  }

  componentDidMount() {
    // const { installGuideProp = {} } = this.props;
    // const { installType } = installGuideProp;
    // this.props.actions.getInstallClusterList({
    //     limit: 0,
    //     type: installType
    // });
  }

  handleSelect = (clusterId) => {
    const { installGuideProp = {}, form } = this.props;
    const { installType, clusterList } = installGuideProp;
    const cluster =
      Array.isArray(clusterList) &&
      clusterList.find((item) => item.id === clusterId);
    Cookie.set('em_current_cluster_id', clusterId);
    Cookie.set('em_current_cluster_type', cluster.type);
    this.props.actions.saveSelectCluster(clusterId);
    // 保存mode信息
    this.setState({
      mode: cluster.mode,
    });
    // 获取第二步的产品包列表
    this.props.actions.getProductStepOneList(
      {
        product_line_name: '',
        product_line_version: '',
        deploy_status: '',
        product_type: installType !== 'kubernetes' ? 0 : 1,
      },
      false
    );
    // 获取namespace
    if (installType === 'kubernetes') {
      this.props.actions.getNamespaceList({ clusterId });
      this.setState({ isNewSpace: false });
      form.setFieldsValue({
        namespace: undefined,
        isNewSpace: false,
      });
    }
  };

  // 集群模糊搜索
  handleSearch = (clusterName: string) => {
    const { clusterList = [] } = this.props.installGuideProp;
    this.setState({
      clusterName,
      clusterList: clusterList.filter(
        (item) => item.name && item.name.indexOf(clusterName) !== -1
      ),
    });
  };

  // 选择命名空间
  handleNamespaceChange = (namespace: string) => {
    if (this.props.installGuideProp.clusterId === -1) {
      message.error('请选择集群');
      setTimeout(() => {
        this.props.form.setFieldsValue({ namespace: undefined });
      }, 0);
      return;
    }
    if (namespace === 'create_new') {
      this.setState(
        {
          isNewSpace: true,
        },
        () => {
          this.props.form.setFieldsValue({
            namespace: '',
            isNewSpace: true,
          });
        }
      );
    } else {
      this.props.actions.saveSelectNamespace(null, { namespace }, false);
    }
  };

  // 新建命名空间原则校验
  validateNamespace = (rule: any, value: any, callback) => {
    if (value && this.state.isNewSpace) {
      const pattern = /^dtstack-[a-z0-9-]{0,54}[a-z0-9]$/;
      if (!pattern.test(value)) {
        callback(
          '请以“dtstack-”开头, 并使用小写字母、数字、中划线进行命名, 不超过63个字符'
        );
      }
    }
    callback();
  };

  // 部署方式變更
  handleInstallTypeChange = (e: RadioChangeEvent) => {
    const installType = e.target.value;
    this.setState({
      isNewSpace: false,
    });
    this.props.actions.saveInstallType(installType);
    this.props.actions.getInstallClusterList({
      limit: 0,
      type: installType,
    });
    this.props.form.setFieldsValue({
      clusterId: undefined,
      namespace: undefined,
      isNewSpace: false,
    });
  };

  render() {
    const { clusterName, isNewSpace, mode } = this.state;
    const { installGuideProp, form } = this.props;
    const { namespaceList, installType } = installGuideProp;
    const { getFieldDecorator } = form;
    const clusterList = clusterName
      ? this.state.clusterList
      : installGuideProp.clusterList;

    // k8s 导入
    const isImport = mode === 1;

    return (
      <div className="step-content-container step-two-container">
        <div className="header-box" style={{ marginBottom: 24 }}>
          <p>
            <Button
              type="primary"
              onClick={() => {
                const newOpen = window.open('/deploycenter/cluster/create');
                newOpen.opener = null;
              }}
            >
              <Icon type="plus" />
              添加集群
            </Button>
            <span style={{ marginLeft: 6 }}>
              选择部署方式和部署集群，将在此集群上部署该产品
            </span>
          </p>
        </div>
        <Form>
          <FormItem required {...formItemCenterLayout} label="部署方式">
            <RadioGroup
              value={installType}
              onChange={this.handleInstallTypeChange}
            >
              {/* <Radio value="kubernetes">
                <span>
                  Kubernetes部署
                  <Tooltip title="应用部署在K8S集群上，适用于容器部署模式">
                    <Icon className="ml-10" type="info-circle" />
                  </Tooltip>
                </span>
              </Radio> */}
              <Radio value="hosts">
                <span>
                  物理/虚拟机部署
                  <Tooltip title="应用部署在物理/虚拟机上，适用于传统部署模式">
                    <Icon className="ml-10" type="info-circle" />
                  </Tooltip>
                </span>
              </Radio>
            </RadioGroup>
          </FormItem>
          <FormItem {...formItemCenterLayout} label="选择集群">
            {getFieldDecorator('clusterId', {
              initialValue:
                installGuideProp.clusterId === -1
                  ? undefined
                  : installGuideProp.clusterId,
              rules: [{ required: true, message: '集群不可为空' }],
            })(
              <Select
                showSearch
                filterOption={false}
                placeholder="请选择集群"
                style={{ width: 460 }}
                onChange={this.handleSelect}
                onSearch={this.handleSearch}
              >
                {Array.isArray(clusterList) &&
                  clusterList.map((item) => (
                    <Option key={item.id} value={item.id}>
                      {item.name}
                    </Option>
                  ))}
              </Select>
            )}
          </FormItem>
          {installType === 'kubernetes' && (
            <FormItem {...formItemCenterLayout} label="选择命名空间">
              {getFieldDecorator('namespace', {
                initialValue: installGuideProp.namespace || undefined,
                rules: [
                  { required: true, message: '命名空间不可为空' },
                  { validator: isNewSpace ? this.validateNamespace : null },
                ],
              })(
                isNewSpace ? (
                  <Input
                    placeholder='请输入新建命名空间，且前缀必须以"dtstack-"开头'
                    style={{ width: 460 }}
                  />
                ) : (
                  <Select
                    showSearch
                    placeholder="请选择命名空间"
                    style={{ width: 460 }}
                    onChange={this.handleNamespaceChange}
                  >
                    {!isImport && (
                      <Option value="create_new">
                        <a
                          onClick={(
                            e: React.MouseEvent<HTMLAnchorElement, MouseEvent>
                          ) => e.preventDefault()}
                        >
                          新建命名空间
                        </a>
                      </Option>
                    )}
                    {Array.isArray(namespaceList) &&
                      namespaceList.map((item) => (
                        <Option key={item} value={item}>
                          {item}
                        </Option>
                      ))}
                  </Select>
                )
              )}
              {isImport && (
                <a
                  className="ml-10"
                  target="_blank"
                  rel="noopener noreferrer"
                  href="/deploycenter/cluster/detail/namespace"
                >
                  新建命名空间
                </a>
              )}
            </FormItem>
          )}
          <FormItem>
            {getFieldDecorator('isNewSpace', {
              initialValue: false,
            })(<Input style={{ display: 'none' }} />)}
          </FormItem>
        </Form>
      </div>
    );
  }
}

export default Form.create<IProps>()(StepOne);
