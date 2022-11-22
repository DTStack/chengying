import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators, Dispatch } from 'redux';
import * as DashboardAction from '@/actions/dashBoardAction';
import { dashboardService } from '@/services';
import { AppStoreTypes } from '@/stores';
import { FormComponentProps } from 'antd/lib/form/Form';
import { Input, Form, Modal, Icon, Upload, Radio, message } from 'antd';

const FormItem = Form.Item;
const RadioButton = Radio.Button;
const RadioGroup = Radio.Group;
const Dragger = Upload.Dragger;
const { TextArea } = Input;
const mapStateToProps = (state: AppStoreTypes) => ({
  dashboard: state.DashBoardStore,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, DashboardAction), dispatch),
});

interface ImportDashModalProp {
  dashboard?: any;
  onClose?: () => void;
  visible?: boolean;
}
interface ImportDashModalState {
  importMod: string;
  importTxt: string;
  importData: object;
  uploadCfg: object;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
class ImportDashModal extends React.Component<
  ImportDashModalProp & FormComponentProps,
  ImportDashModalState
> {
  constructor(props: any) {
    super(props);
  }

  uploadOnChange = (info: any) => {
    const { setFieldsValue } = this.props.form;
    console.log('change:', info);
    if (info.file.status !== 'uploading') {
      console.log(info.file, info.fileList);
    }
    if (info.file.response) {
      if (info.file.response.code !== 0) {
        message.error(info.file.response.msg);
      } else {
        message.success(`${info.file.name} 上传成功！`);
        if (info.file.response.data) {
          this.setState({
            importData: JSON.parse(info.file.response.data),
          });
          setFieldsValue({
            dashname: JSON.parse(info.file.response.data).title,
          });
        }
      }
    }
  };

  state = {
    importMod: '1',
    importTxt: '',
    importData: {
      title: '',
    },
    uploadCfg: {
      name: 'file',
      multiple: false,
      showUploadList: false,
      action: '/api/v2/common/file2text',
      headers: {
        authorization: 'authorization-text',
      },
      onChange: this.uploadOnChange,
    },
  };

  handleModalClose = () => {
    this.setState({
      importMod: '1',
      importTxt: '',
      importData: {
        title: '',
      },
    });
    this.props.onClose();
  };

  showSwitchPanel = () => {
    let panel = null;
    const { importMod, uploadCfg, importTxt } = this.state;
    switch (importMod) {
      case '1':
        panel = (
          <TextArea
            rows={8}
            value={importTxt}
            onChange={this.handleImportTxtChange}
          />
        );
        break;
      case '2':
        panel = (
          <Dragger {...uploadCfg} className="c-dragger_ant-upload">
            <p className="ant-upload-drag-icon">
              <img src={require('public/imgs/icon_upload.svg')}></img>
            </p>
            <p className="ant-upload-text">点击或将文件拖拽到此处上传</p>
          </Dragger>
        );
        break;
    }
    return panel;
  };

  handleImportModChange = (e: any) => {
    this.setState({
      importMod: e.target.value,
    });
  };

  handleImportTxtChange = (e: any) => {
    this.setState({
      importTxt: e.target.value,
    });
  };

  handleDashNameChange = (e: any) => {
    const { importData } = this.state;
    const { setFieldsValue } = this.props.form;
    importData.title = e.target.value;
    this.setState({
      importData,
    });
    setFieldsValue({
      dashname: e.target.value,
    });
  };

  validateDashName = (rule: any, value: any, callback: any) => {
    const { dashboards } = this.props.dashboard;
    // const { getFieldValue } = this.props.form;
    let isDup = false;
    for (const f of dashboards) {
      for (const d of f.list) {
        if (d.title === value) {
          isDup = true;
        }
      }
    }
    return isDup ? callback('仪表盘已经存在') : callback();
  };

  handleSubmitImport = () => {
    const self = this;
    const { importMod, importTxt, importData } = this.state;
    let dashFormData = null;
    switch (importMod) {
      case '1':
        dashFormData = JSON.parse(importTxt);
        break;
      case '2':
        dashFormData = importData;
        break;
    }
    dashboardService
      .importDashboard({
        dashboard: dashFormData,
        inputs: [
          {
            name: 'DS_PROMETHEUS',
            type: 'datasource',
            pluginId: 'prometheus',
            value: 'prometheus',
          },
        ],
        overwrite: true,
      })
      .then((res: any) => {
        self.handleModalClose();
      })
      .catch((err: any) => {
        message.error(err.message);
      });
  };

  render() {
    const { visible } = this.props;
    const { getFieldDecorator } = this.props.form;
    const { importMod } = this.state;

    return (
      <Modal
        title="导入仪表盘"
        width={520}
        visible={visible}
        onOk={this.handleSubmitImport}
        onCancel={this.handleModalClose}>
        <div className="import-data-panel">
          <RadioGroup
            size="small"
            value={importMod}
            onChange={this.handleImportModChange}>
            <RadioButton value="1">
              <Icon type="copy" /> Paste JSON
            </RadioButton>
            <RadioButton value="2">
              <Icon type="cloud-upload" /> Upload .json File
            </RadioButton>
          </RadioGroup>
        </div>
        <div style={{ marginBottom: 10 }}>{this.showSwitchPanel()}</div>
        <Form>
          <FormItem label="名称" required>
            {getFieldDecorator('dashname', {
              rules: [
                {
                  required: true,
                  message: '仪表盘名称不能为空！',
                },
                {
                  validator: this.validateDashName,
                  message: '仪表盘已经存在',
                },
              ],
            })(
              <Input
                style={{ width: 480 }}
                onChange={this.handleDashNameChange}
              />
            )}
          </FormItem>
        </Form>
      </Modal>
    );
  }
}
export default Form.create<any>()(ImportDashModal);
