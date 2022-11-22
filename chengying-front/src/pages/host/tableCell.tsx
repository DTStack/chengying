import * as React from 'react';
import { Input, Icon, message } from 'antd';
import './tableCell.scss';
interface TableProps {
  value: any;
  onChange: any;
  isShowEditIcon: boolean;
}
interface TableState {
  value: any;
  editable: boolean;
}

class EditableCell extends React.Component<TableProps, TableState> {
  state = {
    value: this.props.value,
    editable: false,
  };

  handleChange = (e: any) => {
    const value = e.target.value;
    this.setState({ value });
  };

  /**
   * 父组件调用执行
   */
  successClose = () => {
    return () => {
      this.setState({
        editable: false,
      });
    };
  };

  check = () => {
    const checkGroupName = (name) => {
      const reg = /^[a-zA-Z0-9_@*#/]+$/;
      if (name === '') {
        message.error('分组名称不可为空！');
        return false;
      } else if (name.length > 20) {
        message.error('分组名称不得超过20字符');
        return false;
      } else if (!reg.test(name)) {
        message.error('分组名称输入有误，请重新输入！');
        return false;
      }
      return true;
    };
    if (this.props.onChange && checkGroupName(this.state.value)) {
      this.props.onChange(this.state.value, this.successClose());
    }
  };

  close = () => {
    this.setState({ editable: false, value: this.props.value });
  };

  edit = () => {
    this.setState({ editable: true });
  };

  render() {
    const { value, editable } = this.state;
    return (
      <div className="editable-cell">
        {editable ? (
          <div className="editable-cell-input-wrapper">
            <Input
              value={value}
              style={{ width: '350px', fontSize: '12px' }}
              placeholder="只支持字母、数字及_ / @ * #特殊字符，不超过20个字符"
              onChange={this.handleChange}
              onPressEnter={this.check}
            />
            <Icon
              type="check"
              className="editable-cell-icon-check"
              onClick={this.check}
            />
            <Icon
              type="close"
              className="editable-cell-icon-close"
              onClick={this.close}
            />
          </div>
        ) : this.props.isShowEditIcon ? (
          <div className="editable-cell-text-wrapper">
            {value || ' '}
            <Icon
              type="edit"
              className="editable-cell-icon"
              onClick={this.edit}
            />
          </div>
        ) : null}
      </div>
    );
  }
}
export default EditableCell;
