import * as React from 'react';
import { Form, Input, Select, message } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { clusterManagerService } from '@/services';
import * as EditClusterAction from '@/actions/editClusterAction';
import * as yamlJS from 'js-yaml';
import utils from '@/utils/utils';
import YamlEditor from '@/components/yamlEditor';
import '../style.scss';

const FormItem = Form.Item;
const TextArea = Input.TextArea;
const Option = Select.Option;

const formLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 8 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 9 },
  },
};

interface IProps extends FormComponentProps {
  isEdit: boolean;
  clusterInfo: any;
  location: any;
  action: EditClusterAction.EditClusterActionTypes;
}

interface IState {
  versions: string[];
  networkPlugins: string[];
  avaliable: any[];
  advancedMode: boolean;
}

class StepOne extends React.PureComponent<IProps, IState> {
  state: IState = {
    versions: [],
    networkPlugins: [],
    avaliable: [],
    advancedMode: false,
  };

  componentDidMount() {
    const { location } = this.props;
    const { type, mode } = utils.getParamsFromUrl(location.search || '') as any;
    if (type === 'kubernetes' && +mode === 0) {
      this.getKubernetesAvaliable(mode);
    }
  }

  // 獲取k8s版本號
  getKubernetesAvaliable = (mode: number) => {
    clusterManagerService.getKubernetesAvaliable({ mode }).then((res) => {
      res = res.data;
      if (res.code === 0) {
        const versions = res.data.map((item) => item.version);
        const { properties = {} } = res.data[0];
        const networkPlugins = properties.network_plugin;
        this.setState({
          versions,
          networkPlugins,
          avaliable: res.data,
        });
        this.props.action.setClusterInfo({
          version: versions[0],
          network_plugin: networkPlugins[0],
        });
      }
    });
  };

  // kubernetes版本变更
  handleVersionChange = (version: string) => {
    const { clusterInfo = {}, form } = this.props;
    const { avaliable = [] } = this.state;
    const { properties = {} } =
      avaliable.find((item) => item.version === version) || {};
    const network_plugin = properties.network_plugin[0];
    this.props.action.setClusterInfo({
      version,
      network_plugin,
    });
    clusterInfo.yaml && this.setYaml(version, network_plugin);
    form.setFieldsValue({ network_plugin });
    this.setState({
      networkPlugins: properties.network_plugin,
    });
  };

  // 网络组件切换
  handlePluginChange = (plugin: string) => {
    const { yaml, version } = this.props.clusterInfo;
    yaml && this.setYaml(version, plugin);
  };

  // obj => yaml 同步
  setYaml = (version: string, plugin: string) => {
    const configs = yamlJS.load(this.props.clusterInfo.yaml) || {};
    configs.kubernetes_version = version;
    configs.network.plugin = plugin;

    this.props.action.setClusterInfo({
      yaml: yamlJS.dump(configs),
    });
  };

  // 高级模式切换
  handleModeClick = (e) => {
    e.preventDefault();
    const { advancedMode } = this.state;
    const { clusterInfo, form } = this.props;
    // 第一次创建时需请求获取模板
    if (!advancedMode) {
      form.validateFields(
        ['name', 'version', 'network_plugin'],
        (err, values) => {
          if (!err) {
            !clusterInfo.yaml && this.props.action.getTemplate(values);
            this.setState({ advancedMode: !advancedMode });
          }
        }
      );
    } else {
      // 切回表单进行校验
      const result = this.validateYaml();
      if (result) {
        this.setState({ advancedMode: !advancedMode });
      }
    }
  };

  // yaml配置文件修改
  handleYamlChange = (yaml: string) => {
    this.props.action.setClusterInfo({
      yaml: yaml,
    });
  };

  // 校验yaml语法
  validateYaml = () => {
    const { yaml } = this.props.clusterInfo;
    const { versions, networkPlugins } = this.state;
    try {
      const configs: any = yamlJS.safeLoad(yaml);
      const { kubernetes_version, network = {} } = configs;
      // 版本，网络组件是否存在
      const hasVersion =
        kubernetes_version && versions.includes(kubernetes_version);
      const hasPlugin =
        network.plugin && networkPlugins.includes(network.plugin);
      if (hasVersion && hasPlugin) {
        this.props.action.setClusterInfo({
          version: kubernetes_version,
          network_plugin: network.plugin,
        });
        return true;
      } else {
        message.error('版本或网络组件不匹配');
      }
    } catch (err) {
      message.error(err.message);
    }
    return false;
  };

  // 标签个数校验
  validateTags = (rule: any, value: any, callback: any) => {
    if (value.length > 5) {
      callback('最多输入5个标签');
    }
    if (Array.isArray(value)) {
      const hasError = value.filter((item) => !/^\S{1,32}$/.test(item));
      hasError.length && callback('单个标签不超过32位且不包含空格');
    }
    callback();
  };

  render() {
    const { versions = [], networkPlugins = [], advancedMode } = this.state;
    const { isEdit, form, clusterInfo = {} } = this.props;
    const { getFieldDecorator } = form;
    const kubernetesOwn =
      clusterInfo.type === 'kubernetes' && clusterInfo.mode === 0; // 是k8s自建
    return (
      <Form data-testid="ec-step-one">
        <p className="c-title__color">基础信息</p>
        <div className="mt-20">
          <FormItem label="集群名称" {...formLayout}>
            {getFieldDecorator('name', {
              initialValue: isEdit ? clusterInfo.name : '',
              rules: [
                { required: true, message: '集群名称不可为空' },
                {
                  pattern: /^[A-Za-z0-9_]{1,64}$/,
                  message: '只支持字母、数字、下划线字符，最大长度不超过64',
                },
              ],
            })(<Input placeholder="请输入集群名称" />)}
          </FormItem>
          <FormItem label="集群描述" {...formLayout}>
            {getFieldDecorator('desc', {
              initialValue: isEdit ? clusterInfo.desc : '',
              rules: [
                {
                  min: 1,
                  max: 200,
                  message: '请输入集群描述，最大长度不超过200',
                },
              ],
            })(<TextArea rows={4} placeholder="请输入集群描述" />)}
          </FormItem>
          <FormItem label="集群标签" {...formLayout}>
            {getFieldDecorator('tags', {
              initialValue:
                isEdit && clusterInfo.tags ? clusterInfo.tags.split(',') : [],
              rules: [{ validator: this.validateTags }],
            })(
              <Select
                mode="tags"
                style={{ width: '100%' }}
                placeholder="请输入集群标签，按回车键输入下一个标签，最多5个标签">
                {/* <Option key={-1}>请选择</Option> */}
              </Select>
            )}
          </FormItem>
          <FormItem style={{ display: 'none' }}>
            {getFieldDecorator('id', {
              initialValue: isEdit ? clusterInfo.id : undefined,
            })(<Input />)}
          </FormItem>
        </div>
        {kubernetesOwn && (
          <React.Fragment>
            <p className="c-title__color clearfix">
              <span>Kubernetes选项</span>
              <a
                data-testid="btn-mode-change"
                className="fl-r mr-20"
                onClick={this.handleModeClick}>
                {advancedMode ? '表单输入' : '编辑YMAL'}
              </a>
            </p>
            <div className="mt-20 mb-20">
              {advancedMode ? (
                <YamlEditor
                  yaml={clusterInfo.yaml}
                  onChange={this.handleYamlChange}
                />
              ) : (
                <React.Fragment>
                  <FormItem label="Kubernetes版本" {...formLayout}>
                    {getFieldDecorator('version', {
                      initialValue: isEdit ? clusterInfo.version : versions[0],
                      rules: [
                        { required: true, message: '请选择Kubernetes版本' },
                      ],
                    })(
                      <Select
                        placeholder="请选择Kubernetes版本"
                        onChange={this.handleVersionChange}>
                        {versions.map((item) => (
                          <Option key={item} value={item}>
                            {item}
                          </Option>
                        ))}
                      </Select>
                    )}
                  </FormItem>
                  <FormItem label="网络组件" {...formLayout}>
                    {getFieldDecorator('network_plugin', {
                      initialValue: isEdit
                        ? clusterInfo.network_plugin
                        : networkPlugins[0],
                      rules: [{ required: true, message: '请选择网络组件' }],
                    })(
                      <Select
                        placeholder="请选择网络组件"
                        onChange={this.handlePluginChange}>
                        {networkPlugins.map((item) => (
                          <Option key={item} value={item}>
                            {item}
                          </Option>
                        ))}
                      </Select>
                    )}
                  </FormItem>
                </React.Fragment>
              )}
            </div>
          </React.Fragment>
        )}
      </Form>
    );
  }
}
export default Form.create<IProps>()(StepOne);
