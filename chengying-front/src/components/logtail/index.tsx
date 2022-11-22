import * as React from 'react';
import { message, Select, Tooltip, Input, Icon } from 'antd';
import * as Http from '@/utils/http';

const Option = Select.Option;

interface Prop {
  logs: any;
  serviceid: any;
  isreset: any;
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
}

export default class Logtail extends React.Component<Prop, State> {
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
  };

  UNSAFE_componentWillReceiveProps(nextProps: Prop) {
    // if (nextProps.isreset) {
    //     this.setState({
    //         first_logfile: '',
    //         second_logfile: '',
    //         show_second_select: false,
    //         second_logs: [],
    //         sourceText: ''
    //     });
    //     this.setLogMessage('', true)
    // }
  }

  handleMenu1Click = (path: any) => {
    const ctx = this;
    // let path = log.item.props.children;
    const id = this.props.serviceid;
    if (path.indexOf('*') > -1) {
      Http.get(`/api/v2/instance/${id}/log`, {
        logfile: path.replace(/\s+/g, ''),
        is_match: false,
      }).then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          ctx.setState({
            first_logfile: path,
            second_logs: data?.data?.result?.replace(/\s+/g, '').split(','),
            show_second_select: true,
          });
        } else {
          message.error(data.msg);
        }
      });
    } else {
      Http.get(`/api/v2/instance/${id}/log`, {
        logfile: path,
        is_match: true,
      }).then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          // console.log(data.data.result.repalce('/n', '<br/>'));
          ctx.setState({
            first_logfile: path,
            show_second_select: false,
            sourceText: data.data.result,
          });
          this.setLogMessage(data.data.result, true);
        } else {
          message.error(data.msg);
        }
      });
    }
  };

  handleMenu2Click = (path: any) => {
    const ctx = this;
    // let path = log.item.props.children.trim();
    const id = this.props.serviceid;
    Http.get(`/api/v2/instance/${id}/log`, {
      logfile: path,
      is_match: true,
      tail_num: this.state.lineNum,
    }).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        ctx.setState({
          second_logfile: path,
          show_second_select: true,
          sourceText: data.data.result,
        });
        this.setLogMessage(data.data.result);
      } else {
        message.error(data.msg);
      }
    });
  };

  // 刷新日志
  handleRefreshLogs = () => {
    const id = this.props.serviceid;
    let path = this.state.first_logfile;
    if (
      parseInt(this.state.lineNum) > 1000 ||
      parseInt(this.state.lineNum) < 0
    ) {
      message.error('行数只能位于0-1000之间');
    }
    if (this.state.second_logfile) {
      path = this.state.second_logfile;
    }
    Http.get(`/api/v2/instance/${id}/log`, {
      logfile: path,
      is_match: true,
      tail_num: this.state.lineNum,
    }).then((data: any) => {
      data = data.data;
      if (data.code === 0) {
        // console.log(data.data.result.repalce('/n', '<br/>'));
        this.setState({
          sourceText: data.data.result,
        });
        this.setLogMessage(data.data.result, true);
      } else {
        message.error(data.msg);
      }
    });
  };

  handleLineNumChange = (e: any) => {
    this.setState({
      lineNum: e.target.value,
    });
  };

  // tslint:disable
  handleSearchChange = (e: any) => {
    // 杀了我算了
    const p = e.toLowerCase();
    this.setState({
      searchText: e,
    });
    if (e === '') {
      this.setLogMessage(this.state.sourceText);
      return;
    }
    // 用小写的text获取小写e的坐标
    const a = this.state.sourceText.toLocaleLowerCase().split(p);
    const _a: any = [];
    a.forEach((o: any) => {
      _a.push(o.length);
    });

    const _b: any = [];
    let __b = 0;
    _a.forEach((o: any) => {
      _b.push(__b + o);
      __b += o + e.length;
      console.log('KK咋回事' + this.state.sourceText.charAt(__b));
    });
    console.log(_b);

    let r = '';
    let f = '';
    const d = this.state.sourceText;
    let q = 0;
    // debugger;
    for (let i = 0; i < _b.length; i++) {
      f = d.slice(_b[i], _b[i] + e.length);
      f = `<span style="color: red">${f}</span>`;

      r += d.slice(q, _b[i]) + f;
      // console.log(r)
      q = _b[i] + e.length;
    }
    console.log(r);
    this.setLogMessage(r);
  };

  setLogMessage = (text: any, goBottom?: boolean) => {
    document.getElementById('textPre').innerHTML = text;
    if (goBottom) {
      document.getElementsByClassName('log-message-box').length > 0 &&
        document
          .getElementsByClassName('log-message-box')[0]
          .scrollTo(0, document.getElementById('textPre').offsetHeight);
    }
    this.setState({
      log_message: text,
    });
  };

  render() {
    const { logs } = this.props;
    const { first_logfile, second_logfile, second_logs = [] } = this.state;
    return (
      <div className="logtail-box">
        <div className="option-bar" style={{ display: 'flex' }}>
          <Select
            style={{ height: 26, minWidth: 150, marginBottom: 10 }}
            value={first_logfile}
            onChange={this.handleMenu1Click}>
            {logs.map((l: any, i: number) => {
              return (
                <Option key={i + ''} value={l}>
                  <Tooltip key={i + ''} title={l}>
                    {l}
                  </Tooltip>
                </Option>
              );
            })}
          </Select>
          <div
            style={
              this.state.show_second_select
                ? { display: 'inline-block', marginLeft: 8 }
                : { display: 'none' }
            }>
            <Select
              style={{
                height: 26,
                minWidth: 150,
                maxWidth: 300,
                marginBottom: 10,
              }}
              value={second_logfile}
              onChange={this.handleMenu2Click}>
              {second_logs.map((f: any, j: number) => {
                return (
                  <Option key={j + ''} value={f}>
                    <Tooltip key={j + ''} title={f.trim()}>
                      {f.trim()}
                    </Tooltip>
                  </Option>
                );
              })}
            </Select>
          </div>
          <Input
            placeholder="输入日志行数，“0-1000”之间"
            onPressEnter={this.handleRefreshLogs}
            style={{ height: 26, width: 230, marginLeft: 8, marginBottom: 10 }}
            onChange={this.handleLineNumChange}
            addonAfter={
              <Icon
                style={{ cursor: 'pointer' }}
                onClick={this.handleRefreshLogs}
                type="reload"
              />
            }
          />

          <Input.Search
            placeholder="输入关键词"
            style={{ height: 26, width: 179, marginLeft: 10, marginBottom: 10 }}
            onSearch={this.handleSearchChange}
          />
          {/* <Button style={{ marginLeft: '10px' }} onClick={this.handleRefreshLogs.bind(this)}>刷新</Button> */}
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
          <span>
            <pre id="textPre" style={{ paddingLeft: 10 }}></pre>
          </span>
        </div>
      </div>
    );
  }
}
