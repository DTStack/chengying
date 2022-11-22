import * as React from 'react';
import { TreeSelect, Select, Tooltip, Input, Modal, Icon } from 'antd';

// tslint:disable:variable-name
// tslint:disable:prefer-for-of

const Option = Select.Option;
const { TreeNode } = TreeSelect;

interface Prop {
  fileList?: any[]; // 文件下拉列表
  serviceData: any; // 服务列表
  title: string; // 弹框标题
  maskClosable?: boolean; // 点击蒙层是否允许关闭
  visible?: boolean;
  content: string; // 内容
  selectedFile?: string; // 已选择的文件
  onCancel?: (e: React.MouseEvent<any>) => void;
  onSelectedService?: (service: any) => void; // 选择服务
  onSelectedFile?: (file: any) => void; // 选择文件
  fileGoBottom?: boolean; // 是否默认跳到最新
  loading?: boolean;
}

interface State {
  first_logfile: string;
  second_logfile: string;
  show_second_select: boolean;
  second_logs: any[];
  log_message: string;
  lineNum: string;
  searchText: string;
  sourceText: string;
  query: string; // 搜索关键词
}

export default class FileViewModal extends React.Component<Prop, State> {
  constructor(props: Prop) {
    super(props);
  }

  state: State = {
    first_logfile: '',
    second_logfile: '',
    show_second_select: false,
    second_logs: [],
    log_message: '',
    lineNum: '200',
    searchText: '',
    sourceText: '',
    query: '',
  };

  componentDidMount() {
    this.setLogMessage(this.props.content);
  }

  componentDidUpdate(prevProps, prevState) {
    // 首次进来没有拿到DOM实例
    const ele = document.getElementById('textPre');
    if (
      this.props.content === prevProps.content &&
      this.props.content &&
      !ele.innerHTML
    ) {
      this.setLogMessage(this.props.content);
    }
    if (this.props.content !== prevProps.content) {
      this.setLogMessage(this.props.content);
    }
  }

  onSearch = (query: string) => {
    this.setState(
      {
        query,
      },
      () => {
        this.setLogMessage(this.props.content, false);
      }
    );
  };

  // 文本内容（标红）
  setContent = (content: any) => {
    const queryContent = this.setQueryKey(content);

    /* -- ERROR自动标红 -- */
    const regex = new RegExp(/([\S ]*ERROR[\S ]*)/, 'gi');
    const result = queryContent.replace(
      regex,
      '<span style="color: red">$1</span>'
    );
    return result;
  };

  // 关键词搜索（标红）
  setQueryKey = (content: any) => {
    const { query } = this.state;
    if (query === '') {
      return content;
    }
    const matchTarget = query.toLowerCase();
    const regex = new RegExp(`(${matchTarget})`, 'gi');
    const result = content.replace(
      regex,
      '<span style="color: #2391F7">$1</span>'
    );
    return result;
  };

  // 将文本添加进pre元素
  setLogMessage = (text: any, goBottom?: boolean) => {
    const { fileGoBottom } = this.props;
    const ele = document.getElementById('textPre');
    if (ele) {
      ele.innerHTML = this.setContent(text || '');
      if (goBottom === undefined ? fileGoBottom : goBottom) {
        document.getElementsByClassName('log-message-box').length > 0 &&
          document
            .getElementsByClassName('log-message-box')[0]
            .scrollTo(0, ele.offsetHeight);
      }
    }
  };

  renderServiceSelect = () => {
    const { serviceData, onSelectedService } = this.props;
    const renderNodes = (groups: any) => {
      const nodes = [];
      for (const key in groups) {
        if (Object.prototype.hasOwnProperty.call(groups, key)) {
          const o = groups[key];
          nodes.push(
            <TreeNode value={key} title={key} key={key} selectable={false}>
              {o &&
                o.map((item: any) => (
                  <TreeNode
                    isLeaf={true}
                    value={item.service_name}
                    title={item.service_name_display}
                    key={key + '-' + item.service_name}
                  />
                ))}
            </TreeNode>
          );
        }
      }
      return nodes;
    };

    return (
      <TreeSelect
        showSearch
        style={{ width: 200, marginRight: 10 }}
        dropdownStyle={{ maxHeight: 400, overflow: 'auto' }}
        placeholder="按服务搜索"
        allowClear
        treeDefaultExpandAll
        onChange={onSelectedService}>
        {renderNodes(serviceData)}
      </TreeSelect>
    );
  };

  render() {
    const {
      visible,
      onCancel,
      title,
      maskClosable,
      onSelectedFile,
      fileList,
      content,
      selectedFile,
      serviceData,
      loading,
    } = this.props;

    return (
      <Modal
        className="logtail-box"
        destroyOnClose={true}
        title={title}
        footer={null}
        width={800}
        visible={visible}
        maskClosable={maskClosable}
        onCancel={onCancel}>
        <div className="logtail-box">
          <div className="option-bar" style={{ display: 'flex' }}>
            {serviceData && <span>{this.renderServiceSelect()}</span>}
            {fileList ? (
              <Select
                allowClear
                showSearch={true}
                placeholder="选择文件"
                value={selectedFile}
                style={{ width: 200, marginBottom: 10, marginRight: 10 }}
                onChange={onSelectedFile}>
                {fileList.map((f: any, j: number) => {
                  const name = f.trim();
                  return (
                    <Option key={`${name}-${j}`} value={name}>
                      <Tooltip title={name}>{name}</Tooltip>
                    </Option>
                  );
                })}
              </Select>
            ) : null}
            <Input.Search
              placeholder="输入关键词"
              style={{ width: 200, marginBottom: 10 }}
              onSearch={this.onSearch}
            />
          </div>
          <div
            className="log-message-box"
            style={{
              height: '400px',
              paddingTop: '10px',
              overflowY: 'auto',
              backgroundColor: '#F5F5F5',
              border: '1px dashed #DDDDDD',
            }}>
            {content ? (
              <pre
                id="textPre"
                style={{ paddingLeft: 10, height: 'auto' }}></pre>
            ) : null}
            {content === '' ? (
              <div
                style={{
                  lineHeight: '380px',
                  color: '#dddddd',
                  textAlign: 'center',
                }}>
                没有记录
              </div>
            ) : null}
            {content === null ? (
              <div
                style={{
                  lineHeight: '380px',
                  color: '#dddddd',
                  textAlign: 'center',
                }}>
                请选择上方过滤条件
              </div>
            ) : null}
            {loading && content && <Icon className="ml-10" type="loading" />}
          </div>
        </div>
      </Modal>
    );
  }
}
