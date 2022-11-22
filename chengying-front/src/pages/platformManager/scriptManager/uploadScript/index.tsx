import React from 'react';
import {
  InputNumber,
  Form,
  Button,
  message,
  Modal,
  Upload,
  Icon,
  Input,
} from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { scriptManager } from '@/services';
import './style.scss';
import axios from 'axios';
const FormItem = Form.Item;
const { TextArea } = Input;
const innerAxios = {};

interface Prop extends FormComponentProps {
  title: string;
  showVisible: boolean;
  onClose: () => void;
  onOk: () => void;
  detailInfo?: any;
  id?: number;
}
interface State {
  visible: boolean;
  fileList: any;
  describe: string;
  exec_timeout: number;
  log_retention: number;
  uploadFile: any;
}

class UploadScript extends React.Component<Prop, State> {
  state: State = {
    visible: this.props.showVisible,
    uploadFile: null,
    fileList: [],
    describe: this.props?.detailInfo?.describe,
    exec_timeout: this.props?.detailInfo
      ? this.props?.detailInfo.execTimeout
      : 60,
    log_retention: this.props?.detailInfo
      ? this.props?.detailInfo?.logRetention
      : 3,
  };

  handleCancle = () => {
    this.setState({ visible: false });
    this.props.onClose();
  };

  doUploadScript = () => {
    const {
      uploadFile: file,
      describe,
      exec_timeout,
      log_retention,
    } = this.state;
    const formData = new FormData();
    formData.append('file', file);
    formData.append('describe', describe ?? '');
    formData.append('exec_timeout', exec_timeout.toString());
    formData.append('log_retention', log_retention.toString());
    var CancelToken = axios.CancelToken;
    let source = CancelToken.source();
    axios
      .post('/api/v2/task/upload', formData, {})
      .then(({ data }): any => {
        if (data?.code !== 0) {
          return message.error(data.msg);
        } else {
          this.props.onOk();
        }
      })
      .catch((onError) => {});
    innerAxios[file.uid] = {
      abort() {
        source.cancel();
      },
    };
  };

  editScriptInfo = () => {
    const { id } = this.props;
    const { describe, exec_timeout, log_retention } = this.state;
    const obj: any = {
      id,
      describe,
      exec_timeout,
      log_retention,
    };
    scriptManager.editScript(obj).then(({ data }): any => {
      if (data?.code !== 0) {
        return message.error(data.msg);
      } else {
        this.props.onOk();
      }
    });
  };

  handleSubmit = (e) => {
    this.props.form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        const { detailInfo } = this.props;
        if (detailInfo) {
          this.editScriptInfo();
        } else {
          this.doUploadScript();
        }
      } else {
        e.preventDefault();
        if (!values.exec_timeout.toString) {
          this.props.form.setFieldsValue({ exec_timeout: 60 });
        }
        if (!values.log_retention.toString) {
          this.props.form.setFieldsValue({ log_retention: 3 });
        }
      }
    });
  };

  handleUploadFile = (file: any) => {
    if (file) {
      const isAccept =
        file.type === 'text/x-sh' || file.type === 'text/x-python-script';
      if (!isAccept) {
        message.error('仅支持 .py，.sh 格式文件!');
        return;
      }
    }
  };

  customRequest = (options) => {
    const isAccept =
      options.file.type === 'text/x-sh' ||
      options.file.type === 'text/x-python-script';
    if (isAccept) {
      this.setState({
        fileList: [options.file],
        uploadFile: options.file,
      });
    } else {
      this.setState({ fileList: [], uploadFile: null });
    }
  };

  onRemove = (file) => {
    this.setState({ fileList: [], uploadFile: null });
  };

  changeExecTime = (value) => {
    this.setState({ exec_timeout: value });
  };

  changeDescribe = (e) => {
    this.setState({ describe: e.target.value });
  };

  changeLog = (value) => {
    this.setState({ log_retention: value });
  };

  render() {
    const { getFieldDecorator } = this.props.form;
    const { title, detailInfo } = this.props;
    const { visible, fileList, describe, exec_timeout, log_retention } =
      this.state;

    const formLayout = {
      labelCol: {
        xs: { span: 12 },
        sm: { span: 6 },
      },
      wrapperCol: {
        xs: { span: 24 },
        sm: { span: 15 },
      },
    };

    const uploadFileProps: any = {
      name: 'file',
      customRequest: (options) => this.customRequest(options),
      beforeUpload: (file: any) => {
        this.handleUploadFile(file);
      },
      onRemove: (file) => this.onRemove(file),
      fileList: [],
    };

    return (
      <div className="uploadScript">
        <Modal
          className="uploadScriptBox"
          title={title}
          width={520}
          visible={visible}
          onCancel={this.handleCancle}
          onOk={this.handleSubmit}>
          <Form>
            {!detailInfo && (
              <FormItem label="上传文件 " {...formLayout}>
                <Upload {...uploadFileProps} fileList={fileList}>
                  <Button>
                    <Icon type="upload" /> 上传脚本
                  </Button>
                  <span className="noticeTxt">仅支持.py，.sh 格式文件</span>
                </Upload>
              </FormItem>
            )}
            <FormItem label="脚本描述 " {...formLayout}>
              <TextArea
                value={describe}
                maxLength={200}
                onChange={this.changeDescribe}
                placeholder="请输入脚本描述，长度限制在200个字符以内"
                style={{ width: '280px', height: 56 }}
              />
            </FormItem>
            <FormItem label="超时设置 " {...formLayout}>
              {getFieldDecorator('exec_timeout', {
                initialValue: exec_timeout,
                rules: [
                  { required: true, message: '请输入超时时间' },
                  {
                    pattern: /^[1-9]*[1-9][0-9]*$/,
                    message: '请输入大于0的正整数',
                  },
                ],
              })(
                <div>
                  <InputNumber
                    value={exec_timeout}
                    step={1}
                    onChange={this.changeExecTime}
                  />{' '}
                  s
                </div>
              )}
            </FormItem>
            <FormItem label="执行历史保存周期 " {...formLayout}>
              {getFieldDecorator('log_retention', {
                initialValue: log_retention,
                rules: [
                  { required: true, message: '请输入保存周期' },
                  { pattern: /^[1-7]$/, message: '请输入7以内正整数' },
                ],
              })(
                <div>
                  <InputNumber
                    value={log_retention}
                    step={1}
                    onChange={this.changeLog}
                  />{' '}
                  天 <span className="noticeTxt">最长不超过7天</span>
                </div>
              )}
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create<Prop>()(UploadScript);
