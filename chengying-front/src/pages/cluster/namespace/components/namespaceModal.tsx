import * as React from 'react';
import { Modal, Form, Input, Select, Radio, message } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { ClusterNamespaceService } from '@/services';
import DraggerUpload, { DraggerUploadProps } from './draggerUpload/index';
import YamlEditor, { YamlEditorProps } from '@/components/yamlEditor';

const FormItem = Form.Item;
const Option = Select.Option;
const RadioGroup = Radio.Group;

const DraggerUploadRef = React.forwardRef((props: DraggerUploadProps, ref) => (
  <DraggerUpload {...props} />
));
const YamlEditorRef = React.forwardRef((props: YamlEditorProps, ref) => (
  <YamlEditor {...props} />
));

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 6 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16 },
  },
};

interface IProps extends FormComponentProps {
  namespace: string;
  visible: boolean;
  imageStoreList: any[];
  authorityList: any;
  handleSave: (values: any, callback?: Function) => void;
  handleCancel: () => void;
}
const NamespaceModal: React.FC<IProps> = (props) => {
  const {
    getFieldDecorator,
    validateFields,
    setFieldsValue,
    getFieldsValue,
    getFieldValue,
  } = props.form;
  const {
    namespace,
    visible,
    imageStoreList,
    authorityList,
    handleSave,
    handleCancel,
  } = props;
  const isEdit = namespace !== undefined;
  const [namespaceInfo, setNamespaceInfo] = React.useState<any>({});
  const [fileName, setFileName] = React.useState<string>('');
  const [yaml, setYaml] = React.useState<string>('');
  const [loading, setLoading] = React.useState<boolean>(false);
  const currentValue = getFieldsValue();
  const CAN_EDIT_YAML = authorityList.yaml_edit;

  // 编辑状态下获取namespace信息
  React.useEffect(() => {
    if (namespace === undefined) {
      return;
    }
    ClusterNamespaceService.getNamespaceInfo({
      namespace,
    }).then((response: any) => {
      const res = response.data;
      const { code, data, msg } = res;
      if (code === 0) {
        setNamespaceInfo(data || {});
        setFileName(data.file_name || '');
        setYaml(data.yaml || '');
      } else {
        message.error(msg);
      }
    });
  }, [namespace]);

  // agent下，如果命名空间和镜像仓库改变，重新获取yaml
  React.useEffect(() => {
    if (currentValue.type === 'agent') {
      ClusterNamespaceService.getYamlFile({
        namespace: currentValue.namespace,
        registry_id: currentValue.registry_id,
      }).then((response: any) => {
        const res = response.data;
        const { code, data, msg } = res;
        if (code === 0) {
          setYaml(data.yaml);
        } else {
          message.error(msg);
        }
      });
    }
  }, [currentValue.namespace, currentValue.registry_id, currentValue.type]);

  // 设置文件内容
  const setFileContent = (fileList: any[]) => {
    if (fileList.length > 0) {
      const file = fileList[0];
      const fileReader = new FileReader();
      fileReader.readAsText(file);
      fileReader.onload = function () {
        setFieldsValue({ file: fileReader.result });
        setFileName(file.name);
      };
    } else {
      setFieldsValue({ file: '' });
      setFileName('');
    }
  };

  // 保存
  const handleOk = () => {
    validateFields((err, values) => {
      if (!err) {
        setLoading(true);
        const params = {
          ...values,
        };
        if (params.type === 'kubeconfig') {
          params.yaml = params.file;
          params.file_name = fileName;
          delete params.file;
        } else {
          params.yaml = yaml;
        }
        handleSave(params, () => {
          setLoading(false);
        });
      }
    });
  };

  return (
    <Modal
      title={(isEdit ? '编辑' : '添加') + '命名空间'}
      visible={visible}
      confirmLoading={loading}
      onOk={handleOk}
      onCancel={handleCancel}>
      <Form>
        <FormItem label="命名空间名称" {...formItemLayout}>
          {getFieldDecorator('namespace', {
            initialValue: isEdit ? namespaceInfo.namespace : '',
            rules: [
              { required: true, message: '请输入命名空间' },
              {
                pattern: /^\S{1,64}$/,
                message: '请输入除空格外的字符，不超过64个字符',
              },
            ],
          })(<Input placeholder="请输入命名空间名称" disabled={isEdit} />)}
        </FormItem>
        <FormItem label="选择镜像仓库" {...formItemLayout}>
          {getFieldDecorator('registry_id', {
            initialValue:
              isEdit && namespaceInfo.registry !== -1
                ? namespaceInfo.registry
                : undefined,
            rules: [{ required: true, message: '请选择镜像仓库' }],
          })(
            <Select placeholder="请选择镜像仓库">
              {Array.isArray(imageStoreList) &&
                imageStoreList.map((item: any) => (
                  <Option key={item.id} value={item.id}>
                    {item.name}
                  </Option>
                ))}
            </Select>
          )}
        </FormItem>
        <FormItem required label="添加方式" {...formItemLayout}>
          {getFieldDecorator('type', {
            initialValue: isEdit ? namespaceInfo.type : 'kubeconfig',
          })(
            <RadioGroup>
              <Radio value="kubeconfig">kubeconfig</Radio>
              <Radio value="agent">agent</Radio>
            </RadioGroup>
          )}
        </FormItem>
        {getFieldValue('type') !== 'agent' ? (
          <FormItem label="上传文件" {...formItemLayout}>
            {getFieldDecorator('file', {
              initialValue: isEdit ? namespaceInfo.yaml : '',
              rules: [{ required: true, message: '请上传文件' }],
            })(
              <DraggerUploadRef
                icon={<img src={require('public/imgs/icon_upload.svg')} />}
                defaultFileList={
                  isEdit && namespaceInfo.file_name
                    ? [
                        {
                          uid: '1',
                          name: namespaceInfo.file_name,
                          status: 'done',
                        },
                      ]
                    : []
                }
                onChange={setFileContent}
              />
            )}
          </FormItem>
        ) : (
          <FormItem>
            {getFieldDecorator('yaml', {
              initialValue: isEdit ? namespaceInfo.yaml : '',
            })(
              <YamlEditorRef
                yaml={yaml}
                readOnly={!CAN_EDIT_YAML}
                onChange={(value: string) => setYaml(value)}
              />
            )}
          </FormItem>
        )}
        <FormItem style={{ display: 'none' }}>
          {getFieldDecorator('id', {
            initialValue: isEdit ? namespaceInfo.id : undefined,
          })(<Input />)}
        </FormItem>
      </Form>
    </Modal>
  );
};
export default Form.create<IProps>()(NamespaceModal);
