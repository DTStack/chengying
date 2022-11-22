import { FormComponentProps } from 'antd/lib/form';
import { formItemBaseLayout } from '@/constants/formLayout';
import {
  Button,
  Form,
  Icon,
  message,
  Modal,
  Select,
  TreeSelect,
  Upload,
} from 'antd';
import { cloneDeep, isEqual } from 'lodash';
import { Service } from '@/services';
import * as React from 'react';

const { TreeNode } = TreeSelect;
const FormItem = Form.Item;
const Option = Select.Option;

const uploadImg = require('public/imgs/icon_upload@3x.png');

interface IProps extends FormComponentProps {
  visible: boolean;
  onCancel?: (e?: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  resetKey?: (e?: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  data?: any[];
  defaultValue?: any;
}

interface IState {
  componentList: any[];
  category: any[];
  fileList: any[];
  btnLoading: boolean;
}

class UpdatePatchModal extends React.Component<IProps, IState> {
  state: IState = {
    componentList: [],
    fileList: [],
    category: [],
    btnLoading: false,
  };

  componentDidUpdate(prePros, nextProps) {
    const { defaultValue, form, data } = this.props;
    if (data && !isEqual(data, prePros.data)) {
      this.setState({ componentList: cloneDeep(data) });
    }
    if (defaultValue && !isEqual(defaultValue, prePros.defaultValue)) {
      form.setFieldsValue({
        id: defaultValue.product_id,
      });
      this.getCatalog(defaultValue.product_id);
    }
  }

  handleComponentChange = (id) => {
    const { componentList } = this.state;
    const { form } = this.props;
    form.resetFields(['updateTarget']);
    componentList.forEach((item) => {
      if (item.id === id) {
        this.getCatalog(id);
      }
    });
  };

  // 获取目录
  getCatalog = (id) => {
    const reqParams = Object.assign({}, { product_id: id });
    Service.getPatchPath(reqParams).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const path = res.data.path;
        this.setState({
          category: path,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 目录下拉框
  renderTreeNode = (data) => {
    return (
      Array.isArray(data) &&
      data.map((item: any, index: number) => {
        if (item.children && item.children.length) {
          return (
            <TreeNode
              value={item.path}
              disabled
              title={item.name}
              key={`${Math.random() + index}`}>
              {this.renderTreeNode(item.children)}
            </TreeNode>
          );
        } else {
          return (
            <TreeNode
              value={item.path}
              title={item.name}
              key={`${Math.random() * 10 + index}`}
            />
          );
        }
      })
    );
  };

  normFile = (e: any) => {
    if (Array.isArray(e)) {
      return e;
    }
    const fileList = e && e.fileList.slice(-1);
    this.setState({ fileList: fileList });
    return fileList;
  };

  handleUpdate = () => {
    const { validateFields } = this.props.form;
    const { componentList } = this.state;
    const { defaultValue, resetKey } = this.props;
    validateFields(async (err, value) => {
      if (value.fileList[0].originFileObj) {
        const flag = this.checkFile(value.fileList[0].originFileObj);
        if (!flag) return;
      }
      if (!err) {
        let parentProductName,
          product_name,
          version,
          product_type,
          package_name;
        componentList.forEach((item) => {
          if (item.id === value.id) {
            parentProductName = item.product.ParentProductName;
            product_name = item.product_name;
            version = item.product_version;
            product_type = item.product_type;
            package_name =
              value.fileList[0]?.originFileObj?.name ||
              defaultValue.package_name;
          }
        });
        const fileTypeCheck = this.checkFileType(
          package_name,
          value.updateTarget
        );
        if (!fileTypeCheck) {
          message.error('补丁包类型必须和更新目标后缀保持一致！');
          return;
        }
        this.setState({ btnLoading: true });
        const statusParams = {
          parentProductName,
          product_name,
          version,
          path: value.updateTarget,
          product_type,
          package_name,
        };
        const reqParams = {
          ...statusParams,
          package: value.fileList[0].originFileObj,
        };
        const res = await Service.updatePatchPath(reqParams);

        if (res && res.data.code === 0) {
          const uuid = res.data.data;
          this.clearFields();
          this.setState({ btnLoading: false });
          resetKey();
          const statusRes = await Service.updatePatchStatus({
            ...statusParams,
            uuid,
          });
          if (statusRes && statusRes.data.code === 0) {
            resetKey();
          } else {
            message.error(statusRes.data.msg);
          }
        } else {
          this.setState({ btnLoading: false });
          message.error(res.data.msg);
        }
      } else {
        message.error(err);
        this.setState({ btnLoading: false });
      }
    });
  };

  checkFileType = (fileName, pathName) => {
    const fileType = fileName.split('.').pop();
    const pathType = pathName.split('.').pop();
    if (fileType != pathType) {
      return false;
    } else {
      return true;
    }
  };

  beforeUpload = (file) => {
    this.checkFile(file);
    return false;
  };

  checkFile = (file) => {
    if (!file) return false;
    if (file.size > 1024 * 1024 * 1024) {
      message.error('上传文件大小不能超过1G！');
      return false;
    } else {
      return true;
    }
  };

  clearFields = () => {
    const { onCancel, form } = this.props;
    form.resetFields();
    this.setState({
      category: [],
      btnLoading: false,
    });
    onCancel();
  };

  render() {
    const { getFieldDecorator } = this.props.form;
    const { visible, defaultValue } = this.props;
    const { componentList, category, btnLoading } = this.state;
    const initialFileList = defaultValue?.package_name
      ? [
          {
            uid: -1,
            name: defaultValue.package_name,
          },
        ]
      : [];
    const disabled = !!defaultValue?.package_name;
    return (
      <Modal
        title="补丁包更新"
        className="patch-modal"
        visible={visible}
        footer={[
          <Button key="cancel" onClick={this.clearFields}>
            取消
          </Button>,
          <Button
            key="download"
            type="primary"
            onClick={this.handleUpdate}
            loading={btnLoading}>
            更新
          </Button>,
        ]}
        onCancel={this.clearFields}>
        <div className="update-tips">
          <Icon
            theme="filled"
            type="exclamation-circle"
            style={{ color: '#3F87FF', marginLeft: 20, marginRight: 8 }}
          />
          <div className="update-content">
            补丁包更新用于紧急更新，无需制作标准安装包，可对已部署组件进行快速更新
          </div>
        </div>
        <Form>
          <FormItem {...formItemBaseLayout} label="选择组件">
            {getFieldDecorator('id', {
              validateTrigger: 'onBlur',
              rules: [{ required: true, message: '请选择组件' }],
            })(
              <Select
                placeholder="请选择组件"
                onChange={this.handleComponentChange}
                disabled={disabled}
                showSearch
                filterOption={(input, option) => {
                  const temStr = option.props.children + '';
                  return (
                    temStr
                      .toLocaleLowerCase()
                      .indexOf(input.toLocaleLowerCase()) >= 0
                  );
                }}>
                {Array.isArray(componentList) &&
                  componentList.map((item: any) => (
                    <Option key={`${item.id}`} value={item.id}>
                      {item.product_name}（{`${item.product_version}`}）
                    </Option>
                  ))}
              </Select>
            )}
          </FormItem>
        </Form>
        <Form>
          <FormItem {...formItemBaseLayout} label="更新目标">
            {getFieldDecorator('updateTarget', {
              rules: [{ required: true, message: '请选择路径' }],
            })(
              <TreeSelect
                style={{ width: '100%' }}
                dropdownStyle={{ maxHeight: 400, overflow: 'auto' }}
                placeholder="请先选择更新组件"
                // treeDefaultExpandAll
                treeNodeFilterProp="title"
                showSearch>
                {this.renderTreeNode(category)}
              </TreeSelect>
            )}
          </FormItem>
        </Form>
        <Form>
          <FormItem {...formItemBaseLayout} label="上传补丁包">
            <div className="dropbox" style={{ height: 114 }}>
              {getFieldDecorator('fileList', {
                // validateTrigger: 'onBlur',
                rules: [
                  {
                    required: true,
                    message: '请选择上传补丁包!',
                  },
                ],
                valuePropName: 'fileList',
                getValueFromEvent: this.normFile,
                initialValue: initialFileList,
              })(
                <Upload.Dragger beforeUpload={this.beforeUpload}>
                  <div className="update-icon">
                    <img src={uploadImg} height="36px" width="36px" />
                  </div>
                  <div>点击或将文件拖拽到此处上传</div>
                </Upload.Dragger>
              )}
            </div>
          </FormItem>
        </Form>
      </Modal>
    );
  }
}
export default Form.create<IProps>()(UpdatePatchModal);
