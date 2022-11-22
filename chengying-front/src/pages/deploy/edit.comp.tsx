import * as React from 'react';
import { Input, Icon } from 'antd';

const REG_DIGIT = new RegExp('^[0-9]*$');

interface State {
  default_input_value: string;
  has_config_value: boolean;
}
// tslint:disable
export class EditPanel extends React.Component<any, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      default_input_value: '',
      has_config_value: false,
    };
  }

  /**
   * 构建实例配置项
   * @param {Obj} data
   */
  buildInstanceEditByNode(data: any, level: any, initKey?: any) {
    const _this = this;
    const nextLevel = level + 1;
    let nodecls: any = null;
    level === 1 ? (nodecls = 'node-item top-floor') : (nodecls = 'node-item'); // console.log('level:', level, nextLevel);
    return Object.keys(data).map(function (keyName, keyIndex) {
      const key = initKey ? initKey + '.' + keyName : keyName;
      // 如果是数组则不显示title
      const isnum = REG_DIGIT.test(keyName);
      if (typeof data[keyName] === 'object') {
        return (
          <li key={key} className={nodecls}>
            <span className="node-item-title">{keyName}: </span>
            <ul style={{ display: 'inline-block', verticalAlign: 'top' }}>
              {_this.buildInstanceEditByNode(data[keyName], nextLevel, key)}
            </ul>
          </li>
        );
      } else if (keyName === 'MaxReplica' || keyName === 'StartAfterInstall') {
        return (
          <li key={key} className={nodecls}>
            <span className="node-item-title">{keyName}: </span>
            <span>{data[keyName]}</span>
          </li>
        );
      } else if (keyName !== 'UseCloud') {
        return (
          <li key={key} className={nodecls}>
            <span
              className="node-item-title"
              style={isnum ? { display: 'none' } : { display: 'inline-block' }}>
              {keyName}:{' '}
            </span>
            <Input
              defaultValue={data[keyName]}
              value={data[keyName]}
              addonAfter={
                <Icon
                  type="reload"
                  style={{ fontSize: 14 }}
                  onClick={_this.handleReloadConfig.bind(_this, key)}
                />
              }
              onChange={_this.handleInputChange.bind(_this, key, 0)}
              onFocus={_this.handleInputFocus.bind(_this, key)}
              onBlur={_this.handleInputBlur.bind(_this, key)}
            />
          </li>
        );
      } else {
        return '';
      }
    });
  }

  /**
   * 构建Config配置项，只对下面value是string类型的节点编辑
   */
  buildConfigEditByNode(data: any, initKey?: any) {
    const _this = this;
    let hasvalue = false;
    const configNodes = Object.keys(data).map(function (keyName, keyIndex) {
      const key = initKey ? initKey + '.' + keyName : keyName;
      return Object.keys(data[keyName]).map(function (keyName2, keyIndex2) {
        const key2 = key + '.' + keyName2;
        if (
          keyName2 === 'Value' &&
          typeof data[keyName][keyName2] === 'string'
        ) {
          hasvalue = true;
          return (
            <li key={key2} className="node-item top-floor">
              <span className="node-item-title">{key}: </span>
              <Input
                defaultValue={data[keyName][keyName2]}
                value={data[keyName][keyName2]}
                addonAfter={
                  <Icon
                    type="reload"
                    style={{ fontSize: 14 }}
                    onClick={_this.handleReloadConfig.bind(_this, key2)}
                  />
                }
                onChange={_this.handleInputChange.bind(_this, key2, 1)}
                onFocus={_this.handleInputFocus.bind(_this, key2)}
                onBlur={_this.handleInputBlur.bind(_this, key2)}
              />
            </li>
          );
        } else {
          return '';
        }
      });
    });
    return hasvalue ? (
      configNodes
    ) : (
      <span style={{ margin: '1px 1%' }}>无</span>
    );
  }

  /**
   * 统一处理input的change事件
   * @param {Event} e
   */
  handleInputChange(key: any, type: any, e: any) {
    this.props.modifyServiceConfig({
      pname: this.props.pindex,
      type: type,
      path: key,
      value: e.target.value,
    });
  }

  handleInputFocus(key: any, e: any) {
    this.setState({
      default_input_value: e.target.value,
    });
  }

  /**
   * 统一处理input的blur事件
   * @param {Event} e
   */
  handleInputBlur(key: any, e: any) {
    if (this.state.default_input_value !== e.target.value) {
      this.props.onBlur({
        pindex: this.props.pindex,
        path: this.props.type + '.' + key,
        // path: 'product.Service.' + this.props.pindex + '.' + this.props.type + '.' + key,
        value: e.target.value,
      });
    }
  }

  // 单个配置字段恢复默认值
  handleReloadConfig(key: any, e: any) {
    this.props.resetServiceConfig({
      pindex: this.props.pindex,
      field_path: this.props.type + '.' + key,
    });
  }

  render() {
    const { data, type, dpd } = this.props;
    // if (data) {
    if (type === 'Instance') {
      if (data) {
        return (
          <div
            className="edit-panel config-box"
            style={data.UseCloud ? { display: 'none' } : { display: 'block' }}>
            <div className="mod-title">实例配置:</div>
            {JSON.stringify(data) !== '{}' ? (
              <ul className="config-tree">
                {this.buildInstanceEditByNode(data, 1)}
              </ul>
            ) : (
              <p style={{ lineHeight: '20px' }}>无</p>
            )}
          </div>
        );
      } else {
        return (
          <div className="edit-panel config-box" style={{ display: 'none' }}>
            <div className="mod-title">实例配置:</div>
            <p style={{ lineHeight: '20px' }}>无</p>
          </div>
        );
      }
    } else {
      if (data) {
        return (
          <div className="edit-panel config-box">
            <div className="mod-title">Config配置:</div>
            <ul className="config-tree">{this.buildConfigEditByNode(data)}</ul>
            <div className="mod-title">依赖服务:</div>
            <p style={{ lineHeight: '20px' }}>{dpd.toString() || '无'}</p>
          </div>
        );
      } else {
        return (
          <div className="edit-panel config-box">
            <div className="mod-title">Config配置:</div>
            <p style={{ lineHeight: '20px' }}>无</p>
          </div>
        );
      }
    }
  }
}
