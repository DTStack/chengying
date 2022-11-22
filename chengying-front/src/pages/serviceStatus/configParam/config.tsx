import * as React from 'react';
import { Select, Modal, Form, Input, Tooltip } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/mode-javascript';
import 'ace-builds/src-noconflict/theme-xcode';
import '../style.scss';

const Option = Select.Option;

interface Prop extends FormComponentProps {
  visible?: boolean;
  onCancel?: (e: React.MouseEvent<any>) => void;
  dropList: any[];
  value: string;
  onChange: (value: any) => void;
  handleContentChange: (value: any, doc: any) => void;
  fileContent: string;
  handleSetParamSubmit: any;
  newAddParamArray: any[];
}

class ParamConfig extends React.Component<Prop, {}> {
  handleSubmit = (e) => {
    this.props.form.validateFields((err, values) => {
      if (!err) {
        this.props.handleSetParamSubmit();
      }
    });
  };

  render() {
    const {
      visible,
      onCancel,
      dropList,
      onChange,
      value,
      fileContent,
      handleContentChange,
      newAddParamArray,
      form,
    } = this.props;
    const { getFieldDecorator } = form;
    const formItemLayout = {
      labelCol: { span: 3 },
      wrapperCol: { span: 10 },
    };
    return (
      <Modal
        title="配置参数"
        width={900}
        visible={visible}
        onCancel={onCancel}
        onOk={this.handleSubmit}
        okText="确认"
        cancelText="取消"
        okButtonProps={{ disabled: !dropList }}>
        <div style={{ marginBottom: '10px' }}>
          <Select style={{ minWidth: 160 }} onChange={onChange} value={value}>
            {dropList != null
              ? dropList.map((l, i) => {
                  return (
                    <Option key={`${i}`} value={l}>
                      <Tooltip title={l}>{l}</Tooltip>
                    </Option>
                  );
                })
              : ''}
          </Select>
        </div>
        <div
          style={{
            backgroundColor: '#F5F5F5',
            border: '1px dashed #DDDDDD',
            height: '400px',
            width: '845px',
          }}
          className="service-config-pre service-config">
          {dropList ? (
            <AceEditor
              mode="javascript"
              theme="xcode"
              onChange={handleContentChange}
              value={fileContent}
              style={{
                width: '100%',
                height: '100%',
                fontSize: '12px',
                lineHeight: '20px',
              }}
            />
          ) : (
            ''
          )}
        </div>
        <div style={{ marginTop: '10px' }}>
          <span style={{ fontWeight: 'bold', fontSize: '12px' }}>
            参数赋值：
          </span>
          <span style={{ fontSize: '12px', color: '#A7A7A7' }}>
            为新增参数赋默认值，若无默认值，输入空格即可
          </span>
        </div>
        <Form>
          {newAddParamArray.map((item, index) => {
            const newLabel =
              item && item.length > 10 ? item.slice(0, 10) + '...' : item;
            return (
              <Form.Item
                {...formItemLayout}
                style={{ fontSize: '16px', marginBottom: '10px' }}
                label={<Tooltip title={item}>{newLabel}</Tooltip>}
                key={index}>
                {getFieldDecorator(`${item}`, {
                  rules: [
                    {
                      required: true,
                      message: '参数默认值不能为空',
                    },
                    {
                      message: '参数默认值长度不能大于512个字符',
                      max: 512,
                    },
                  ],
                })(<Input placeholder="请输入默认值" id={item} />)}
              </Form.Item>
            );
          })}
        </Form>
      </Modal>
    );
  }
}
export default Form.create<Prop>()(ParamConfig);
