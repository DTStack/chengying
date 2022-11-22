import * as React from 'react';
import { Input, Tooltip, Icon, Form, Typography } from 'antd';
import { FormComponentProps } from 'antd/lib/form';

interface Props extends FormComponentProps {
  tooltipOnClick: Function;
  onBlur?: Function;
  onChange?: Function;
  defaultvalue: string;
  disabled?: boolean;
  inputDisabled?: boolean;
  title?: string;
  iconType?: string;
}

class R extends React.Component<Props, any> {
  render() {
    const { defaultvalue, inputDisabled, title, iconType, disabled } =
      this.props;
    // eslint-disable-next-line no-constant-condition
    if (/password/.test(title.toLowerCase())) {
      return (
        <span style={{ display: 'inline-flex' }}>
          <Input
            type={
              !iconType
                ? /password/.test(title.toLowerCase())
                  ? 'password'
                  : 'textarea'
                : 'textarea'
            }
            className="param-input-textarea"
            // rows={1}
            value={defaultvalue}
            // onBlur={(e) => this.props.onBlur(e)}
            onChange={(e) => this.props.onChange(e)}
            disabled={inputDisabled}
          />
          {disabled ? null : (
            <span
              className="afteron"
              onClick={
                inputDisabled ? null : () => this.props.tooltipOnClick()
              }>
              <Tooltip title="恢复默认值">
                <Icon type="undo" />
              </Tooltip>
            </span>
          )}
        </span>
      );
    }
    const { Paragraph } = Typography;
    const mdisabled = inputDisabled ? 'focus' : 'hover';
    const pass = /password/.test(title.toLocaleLowerCase())
      ? mdisabled
      : 'hover';
    return (
      <Tooltip
        trigger={!defaultvalue ? 'focus' : pass}
        arrowPointAtCenter={true}
        title={() => (
          <Paragraph copyable style={{ color: '#fff' }}>
            {defaultvalue}
          </Paragraph>
        )}>
        <span style={{ display: 'inline-flex' }}>
          <Input.TextArea
            className="param-input-textarea"
            // rows={1}
            value={defaultvalue}
            onBlur={(e) => this.props.onBlur(e)}
            onChange={(e) => this.props.onChange(e)}
            disabled={inputDisabled}
            style={{ height: 32 }}
          />
          {disabled ? null : (
            <span
              className="afteron"
              onClick={
                inputDisabled ? null : () => this.props.tooltipOnClick()
              }>
              <Tooltip title="恢复默认值">
                <Icon type="undo" />
              </Tooltip>
            </span>
          )}
        </span>
      </Tooltip>
    );
  }
}
export default Form.create<Props>()(R);
