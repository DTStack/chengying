import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators, Dispatch } from 'redux';
import {
  Layout,
  Input,
  Select,
  Form,
  Row,
  Col,
  Radio,
  Icon,
  Button,
  message,
  Spin,
  Tooltip,
} from 'antd';
import { AppStoreTypes } from '@/stores';
import { ServiceStore, HeaderStoreType } from '@/stores/modals';
import { FormComponentProps } from 'antd/lib/form/Form';
import { servicePageService } from '@/services';
import * as ServiceActions from '@/actions/serviceAction';
import TextMark from './textMark';
import { debounce } from 'lodash';
import './style.scss';

const { Content } = Layout;
const Search = Input.Search;

const MAX_LOG_RANGE = 300;

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
  ServiceStore: state.ServiceStore,
  authorityList: state.UserCenterStore.authorityList,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, ServiceActions), dispatch),
});

interface Prop {
  HeaderStore: HeaderStoreType;
  ServiceStore: ServiceStore;
  authorityList: any;
  actions: any;
  match?: any;
  location?: any;
}
interface DashBoardState {
  sourceRange: any[];
  sourceText: any[];
  tex: string;
  products: any[];
  serviceGroup: any[];
  hostGroup: any[];
  fileGroup: any[];
  cur_product: string;
  cur_service: string;
  cur_host: number;
  cur_file: string;
  action: string;
  log_service_id: number;
  logPaths: any[];
  fileType: string;
  fileLoading: boolean;
  cur_type: string;
  markText: string;
  downUrl: string;
  isFullScreen: boolean;
}
@(connect(mapStateToProps, mapDispatchToProps) as any)
class Log extends React.Component<Prop & FormComponentProps, DashBoardState> {
  constructor(props: any) {
    super(props);
  }

  state: DashBoardState = {
    products: [],
    serviceGroup: [],
    hostGroup: [],
    fileGroup: [],
    cur_product: '',
    cur_service: '',
    cur_host: null,
    log_service_id: null,
    action: 'latest',
    cur_file: '',
    logPaths: [],
    fileType: 'txt',
    sourceRange: [], // 当前日志行数区间
    sourceText: [],
    tex: '',
    fileLoading: false,
    cur_type: '',
    markText: '',
    downUrl: '',
    isFullScreen: false,
  };

  componentDidMount = () => {
    const { state } = this.props.location;
    console.log('11111,', state);
    if (state) {
      this.setState(
        {
          cur_product: state.product_name,
          cur_service: state.service_name,
          cur_host: state.log_service_id,
        },
        () => {
          this.handleSelectProduct(state.product_name);
          this.handSelectService(state.service_name);
          this.handSelectHost(state.log_service_id);
        }
      );
      this.props.form.setFieldsValue({
        component: state.product_name,
        service: state.service_name,
        host: state.log_service_id,
      });
    }
    const { cur_parent_product } = this.props.HeaderStore;
    this.getProductList(cur_parent_product);
    console.log(this.props.location);

    const textPre = document.getElementById('textPre');
    if (this.state.fileType === 'txt') {
      textPre.scrollTop = textPre.scrollHeight - textPre.clientHeight - 30;
      textPre.addEventListener('scroll', this.onScrollEvent);
    }
    this.watchFullScreen();
  };

  componentDidUpdate = (prevProps) => {
    if (prevProps.HeaderStore !== this.props.HeaderStore) {
      this.props.form.setFieldsValue({
        component: null,
        service: null,
        host: null,
      });
      this.setState({
        serviceGroup: [],
        hostGroup: [],
        fileGroup: [],
        cur_product: null,
        cur_host: null,
        cur_service: '',
        sourceText: [],
        sourceRange: [],
      });
    }
  };

  // 监听滚动条
  onScrollEvent = debounce(() => {
    const textPre = document.getElementById('textPre');
    const scrollTop = textPre.scrollTop;
    const clientHeight = textPre.clientHeight;
    const scrollHeight = textPre.scrollHeight;
    console.log('滚动了', scrollTop, clientHeight, scrollHeight);

    const { sourceText, sourceRange } = this.state;
    if (!sourceText.length) {
      return;
    }
    // 向上滚动时 sourceRange = 1 不请求
    if (scrollTop < 1 && sourceRange[0] === 1) {
      return;
    }
    if (scrollTop + clientHeight > scrollHeight - 0.5) {
      console.log('到底了');
      this.onLogFileContent('down');
    } else if (scrollTop < 1) {
      console.log('到头了');
      this.onLogFileContent('up');
    } else {
      console.log('不触发');
    }
  }, 800);

  // 获得组件
  getProductList = (curParentProduct: any) => {
    const self = this;
    servicePageService
      .getProductName({ parentProductName: curParentProduct })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0 && res.data?.list) {
          const list = res.data?.list;
          self.setState({
            products: list,
          });
        }
      });
  };

  // 切换组件
  handleSelectProduct = (item: any) => {
    const serviceGroup = [];
    this.setState({
      serviceGroup: [],
      hostGroup: [],
      fileGroup: [],
      cur_product: item,
      cur_host: null,
      cur_service: '',
      cur_file: '',
      sourceText: [],
      sourceRange: [],
    });
    this.props.form.setFieldsValue({
      component: item,
      service: '',
      host: null,
      file: undefined,
    });
    servicePageService
      .getServiceGroup({ product_name: item })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          if (data.data.count > 0) {
            for (const i in data.data.groups) {
              serviceGroup.push({
                name: i,
                list: data.data.groups[i],
              });
              this.setState({ serviceGroup });
            }
          } else {
            message.error(data.msg);
          }
        }
      });
  };

  // 切换服务
  handSelectService = (item: any) => {
    this.setState({
      cur_service: item,
      cur_host: null,
      cur_file: '',
      hostGroup: [],
      fileGroup: [],
      sourceText: [],
      sourceRange: [],
    });
    this.props.form.setFieldsValue({
      service: item,
      host: null,
      file: undefined,
    });
    const hostGroup = [];
    servicePageService
      .getServiceHostsList({
        product_name: this.state.cur_product,
        service_name: item,
      })
      .then((data: any) => {
        data = data.data;
        if (data.code === 0) {
          if (data.data.count > 0) {
            data.data.list.forEach((item: any) => {
              hostGroup.push({ label: item.id, value: item.ip });
            });
            this.setState({ hostGroup });
            console.log('host,,,,', this.state.hostGroup);
          }
        } else {
          message.error(data.msg);
        }
      });
  };

  // 切换主机
  handSelectHost = (item: any) => {
    this.setState({
      cur_host: item,
      cur_file: '',
      fileGroup: [],
      sourceText: [],
      sourceRange: [],
    });
    this.props.form.setFieldsValue({
      host: item,
      file: undefined,
    });
    const { fileType } = this.state;
    servicePageService
      .getLogFileList({ instanceId: item, type: fileType })
      .then((data: any) => {
        data = data.data;
        console.log('日志文件：', data);
        if (data.code === 0) {
          this.setState({ fileGroup: data.data.list });
        } else {
          message.error(data.msg);
        }
      });
  };

  // 切换日志文件
  handSelectFile = (item: any) => {
    this.setState({
      cur_file: item,
      sourceText: [],
      sourceRange: [],
    });
    if (this.state.fileType === 'txt') {
      this.setState({ cur_file: item }, () => {
        this.onLogFileContent('latest');
      });
    }
  };

  // 获取下一次请求的区间位置
  getLogReqRange = (type: string) => {
    const {
      sourceRange: [startOf, endOf],
    } = this.state;
    let start, end;
    if (type === 'up') {
      const nextStart = startOf - MAX_LOG_RANGE + 1;
      start = nextStart < 1 ? 1 : nextStart;
      end = startOf - 1;
    } else if (type === 'down') {
      start = endOf + 1;
      end = endOf + MAX_LOG_RANGE;
    }
    return { start, end };
  };

  // 获取新的sourceText
  getNewSourceText = (type: string, list: any[]) => {
    if (type === 'latest') {
      return list;
    }
    const { sourceText } = this.state;
    const newLoadedText =
      type === 'up' ? [...list, ...sourceText] : [...sourceText, ...list];
    return newLoadedText;
  };

  // 获取当前日志区间范围
  getSourceRange = (type: string, total: number, sourceText: any[]) => {
    const {
      sourceRange: [startOf, endOf],
    } = this.state;
    let start, end;
    if (type === 'up') {
      const _start = startOf || total;
      start = _start < MAX_LOG_RANGE ? 1 : _start - MAX_LOG_RANGE + 1;
      end = endOf || total;
    } else if (type === 'latest') {
      const _start = total;
      start = _start < MAX_LOG_RANGE ? 1 : _start - MAX_LOG_RANGE + 1;
      end = total;
    } else {
      const _end = endOf + MAX_LOG_RANGE;
      start = startOf;
      end = _end > total ? total : _end;
    }
    return [start, end];
  };

  // 获取日志内容
  onLogFileContent = (type: string) => {
    const { start, end } = this.getLogReqRange(type);
    const config = {
      logfile: this.state.cur_file,
      action: type,
      start,
      end,
    };
    const textPre = document.getElementById('textPre');
    const regionScrollHeight = textPre.scrollHeight;
    this.setState({ fileLoading: true, cur_type: type });
    servicePageService
      .getLogFile({ instanceId: this.state.cur_host }, config)
      .then((response) => {
        const { code, data, msg } = response.data;
        if (code === 0) {
          const list = data?.list || [];
          const total = data?.total || 0;
          const newSourceText = this.getNewSourceText(type, list);
          this.setState(
            {
              sourceText: newSourceText,
              fileLoading: false,
              sourceRange: this.getSourceRange(type, total, newSourceText),
            },
            () => {
              if (type === 'up') {
                // textPre.scrollTop = 30
                const timer = setInterval(() => {
                  const newScrollHeight = textPre.scrollHeight;
                  if (newScrollHeight > regionScrollHeight) {
                    textPre.scrollTop = newScrollHeight - regionScrollHeight;
                    clearInterval(timer);
                  }
                }, 1000);
              } else if (type === 'latest') {
                textPre.scrollTop =
                  textPre.scrollHeight - textPre.clientHeight - 30;
              } else {
                if (list.length === 0) {
                  textPre.scrollTop =
                    textPre.scrollHeight - textPre.clientHeight - 50;
                }
              }
              textPre.addEventListener('scroll', this.onScrollEvent);
            }
          );
        } else {
          this.setState({ fileLoading: false });
          message.error(msg);
        }
      });
  };

  // 下载日志文件
  handleLogFileDown = () => {
  const { authorityList } = this.props;
  if (!authorityList.log_download) {
    message.error('权限不足，请联系管理员！')
    return
  }
    const { cur_host, cur_file } = this.state;
    this.setState({
      downUrl: `/api/v2/instance/${cur_host}/logdown?logfile=${cur_file}`,
    });
  };

  handleSearchChange = (e: any) => {
    this.setState({ markText: e });
  };

  onRadioChange = (e) => {
    this.props.form.setFieldsValue({
      file: undefined,
    });
    this.setState(
      {
        fileType: e.target.value,
        sourceText: [],
        sourceRange: [],
        cur_file: '',
        fileGroup: [],
      },
      () => {
        if (this.state.fileType === 'zip') {
          this.setState({
            sourceText: [],
            sourceRange: [],
          });
        }
        if (
          this.state.cur_product &&
          this.state.cur_host &&
          this.state.cur_service
        ) {
          this.handSelectHost(this.state.cur_host);
        }
      }
    );
  };

  handleFullScreen = () => {
    console.log('fullscreen:', this.state.isFullScreen);
    if (!this.state.isFullScreen) {
      this.requestFullScreen();
    } else {
      this.setState({ isFullScreen: false });
      this.exitFullscreen();
    }
  };

  // 进入全屏
  requestFullScreen = () => {
    console.log('requestFullScreen');
    this.setState({ isFullScreen: true });
    const textPre = document.getElementById('textPre');
    const scrollTop = textPre.scrollTop;
    const scrollHeight = textPre.scrollHeight;
    const clientHeight = textPre.clientHeight;
    const proportion = scrollTop / (scrollHeight - clientHeight);
    textPre.className += ' log_table_fullContent';
    if (this.state.fileType === 'txt') {
      if (
        textPre.scrollTop + textPre.clientHeight >
        textPre.scrollHeight - 0.5
      ) {
        textPre.scrollTop =
          (textPre.scrollHeight - textPre.clientHeight) * proportion - 5;
      } else if (textPre.scrollTop < 1) {
        textPre.scrollTop =
          (textPre.scrollHeight - textPre.clientHeight) * proportion + 5;
      } else {
        textPre.scrollTop =
          (textPre.scrollHeight - textPre.clientHeight) * proportion;
      }
    }
    const de: any = document.getElementById('monitorDataFlowContainer');
    if (de.requestFullscreen) {
      de.requestFullscreen();
    } else if (de.mozRequestFullScreen) {
      de.mozRequestFullScreen();
    } else if (de.webkitRequestFullScreen) {
      de.webkitRequestFullScreen();
    }
  };

  // 退出全屏
  exitFullscreen = () => {
    const textPre = document.getElementById('textPre');
    const scrollTop = textPre.scrollTop;
    const scrollHeight = textPre.scrollHeight;
    const clientHeight = textPre.clientHeight;
    const proportion = scrollTop / (scrollHeight - clientHeight);
    textPre.className = 'log_table_content';
    if (this.state.fileType === 'txt') {
      if (
        textPre.scrollTop + textPre.clientHeight >
        textPre.scrollHeight - 0.5
      ) {
        textPre.scrollTop =
          (textPre.scrollHeight - textPre.clientHeight) * proportion - 5;
      } else if (textPre.scrollTop < 1) {
        textPre.scrollTop =
          (textPre.scrollHeight - textPre.clientHeight) * proportion + 5;
      } else {
        textPre.scrollTop =
          (textPre.scrollHeight - textPre.clientHeight) * proportion;
      }
    }
    var de: any = document;
    if (de.exitFullscreen) {
      de.exitFullscreen();
    } else if (de.mozCancelFullScreen) {
      de.mozCancelFullScreen();
    } else if (de.webkitCancelFullScreen) {
      de.webkitCancelFullScreen();
    }
  };

  // 监听事件
  watchFullScreen = () => {
    const _self = this;
    const de: any = document;
    document.addEventListener(
      'fullscreenchange',
      function () {
        _self.setState({
          isFullScreen: de.fullscreen,
        });
        if (!de.fullscreen) {
          document.getElementById('textPre').className = 'log_table_content';
        }
      },
      false
    );
    document.addEventListener(
      'webkitfullscreenchange',
      function () {
        _self.setState({
          isFullScreen: de.webkitIsFullScreen,
        });
        if (!de.webkitIsFullScreen) {
          document.getElementById('textPre').className = 'log_table_content';
        }
      },
      false
    );
    document.addEventListener(
      'mozfullscreenchange',
      function () {
        _self.setState({
          isFullScreen: de.mozFullScreen,
        });
        if (!de.mozFullScreen) {
          document.getElementById('textPre').className = 'log_table_content';
        }
      },
      false
    );
  };

  render() {
    const { authorityList, form } = this.props;
    const { getFieldDecorator } = form;
    const { Option, OptGroup } = Select;
    const {
      products,
      serviceGroup,
      hostGroup,
      fileLoading,
      cur_type,
      sourceText,
      sourceRange,
      markText,
      fileType,
      cur_file,
      downUrl,
      fileGroup,
      isFullScreen,
    } = this.state;
    const formItemLayout2 = {
      labelCol: { span: 6 },
      wrapperCol: { span: 18 },
    };
    const formItemLayout3 = {
      wrapperCol: { span: 22 },
    };
    return (
      <Layout id="monitorDataFlowContainer">
        <Content>
          <div className="div_fullscreen">
            {!isFullScreen ? (
              <span onClick={this.handleFullScreen}>
                <img
                  src={require('./img/fullScreen.png')}
                  style={{ width: '40px' }}
                />
              </span>
            ) : (
              <span onClick={this.handleFullScreen}>
                <img
                  src={require('./img/smallScreen.png')}
                  style={{ width: '40px' }}
                />
              </span>
            )}
          </div>
          <div className="dash-page">
            <div className="top-navbar clearfix">
              <Form className="ant-advanced-search-form">
                <div style={{ minWidth: '1000px', display: 'flex' }}>
                  <Form.Item label="组件：">
                    {getFieldDecorator('component', {
                      rules: [{ required: true, message: '请选择组件' }],
                    })(
                      <Select
                        className="dt-form-shadow-bg"
                        style={{ width: 264 }}
                        onChange={this.handleSelectProduct}
                        getPopupContainer={(e: any) => e.parentNode}>
                        {products
                          ? products.map((item, index) => (
                              <Option key={index} value={item.product_name}>
                                {item.product_name}
                              </Option>
                            ))
                          : ''}
                      </Select>
                    )}
                  </Form.Item>
                  <Form.Item label="服务：">
                    {getFieldDecorator('service', {
                      rules: [{ required: true, message: '请选择服务' }],
                    })(
                      <Select
                        className="dt-form-shadow-bg"
                        style={{ width: 264 }}
                        onChange={this.handSelectService}
                        getPopupContainer={(e: any) => e.parentNode}>
                        {serviceGroup
                          ? serviceGroup.map((item, inde) => {
                              return (
                                <OptGroup label={item.name} key={inde}>
                                  {item.list.map((i, index) => (
                                    <Option key={index} value={i.service_name}>
                                      {i.service_name}
                                    </Option>
                                  ))}
                                </OptGroup>
                              );
                            })
                          : ''}
                      </Select>
                    )}
                  </Form.Item>
                  <Form.Item label="主机：">
                    {getFieldDecorator('host', {
                      rules: [{ required: true, message: '请选择主机' }],
                    })(
                      <Select
                        className="dt-form-shadow-bg"
                        style={{ width: 264 }}
                        onChange={this.handSelectHost}
                        getPopupContainer={(e: any) => e.parentNode}>
                        {hostGroup
                          ? hostGroup.map((item) => (
                              <Option key={item.label} value={item.label}>
                                {item.value}
                              </Option>
                            ))
                          : []}
                      </Select>
                    )}
                  </Form.Item>
                </div>
                <div style={{ minWidth: '1280px', display: 'flex' }}>
                  <Form.Item label="日志文件" {...formItemLayout2}>
                    <Radio.Group onChange={this.onRadioChange} value={fileType}>
                      <Radio value="txt">文本文件</Radio>
                      <Radio value="zip">压缩文件</Radio>
                    </Radio.Group>
                  </Form.Item>
                  <Form.Item {...formItemLayout3}>
                    {getFieldDecorator('file', {
                      rules: [{ required: true, message: '请选择日志文件' }],
                    })(
                      <Select
                        placeholder="请选择日志文件"
                        className="dt-form-shadow-bg"
                        style={{ width: 264 }}
                        getPopupContainer={(e: any) => e.parentNode}
                        onChange={this.handSelectFile}
                        showSearch>
                        {Array.isArray(fileGroup) &&
                          fileGroup.map((item, index) => (
                            <Option key={index} value={item}>
                              <Tooltip
                                placement="right"
                                title={item}
                                key={index}>
                                {item.length > 35
                                  ? item.slice(0, 35) + '...'
                                  : item}
                              </Tooltip>
                            </Option>
                          ))}
                      </Select>
                    )}
                  </Form.Item>
                  <Form.Item {...formItemLayout3}>
                    <Search
                      placeholder="请输入日志关键字"
                      className="dt-form-shadow-bg"
                      style={{ width: 264 }}
                      onSearch={this.handleSearchChange}
                    />
                  </Form.Item>
                </div>
              </Form>
            </div>
            <div className="log_table box-shadow-style">
              <Row className="log_table_header">
                <Col span={3} style={{ padding: '6px' }}>
                  日志内容
                </Col>
                <Col span={21} style={{ textAlign: 'right' }}>
                  <Button
                    type="primary"
                    style={{ marginRight: '5px' }}
                    ghost
                    disabled={!(fileType === 'txt' && cur_file)}
                    onClick={() => this.onLogFileContent('latest')}>
                    <Icon type="step-forward" /> 跳至日志最新行
                  </Button>
                  {authorityList.log_download && (
                    <Button
                      type="primary"
                      ghost
                      disabled={!cur_file}
                      onClick={this.handleLogFileDown}>
                      <a href={downUrl} download={cur_file}>
                        <Icon type="download" /> 下载完整日志
                      </a>
                    </Button>
                  )}
                </Col>
              </Row>
              <Spin
                spinning={fileLoading && cur_type === 'latest'}
                tip="日志加载中..."
                style={{ width: '100%' }}>
                <div className="log_table_content" id="textPre">
                  {sourceText.length > 0 && (
                    <Spin
                      spinning={fileLoading && cur_type === 'up'}
                      style={{ position: 'relative', paddingLeft: '49%' }}
                    />
                  )}
                  {sourceRange.length > 0 && sourceRange[0] === 1 && (
                    <div className="log-to-top">
                      <span>日志到顶啦</span>
                      <div className="line"></div>
                    </div>
                  )}
                  {fileType === 'zip' ? (
                    <div className="log_content" id="log_content">
                      <p className="log_table_content_p1">
                        <img
                          src={require('./img/file.png')}
                          className="textPre_img"
                        />
                      </p>
                      <p className="log_table_content_p2">
                        不支持压缩文件日志内容预览，请下载查看!
                      </p>
                    </div>
                  ) : sourceText.length ? (
                    sourceText.map((item: string, index: number) => {
                      const lineStart = sourceRange.length ? sourceRange[0] : 1;
                      if (item.includes('ERROR')) {
                        return (
                          <p
                            className="log_table_content__line"
                            key={index.toString()}
                            style={{ color: 'red' }}>
                            <span className="line_number">
                              {lineStart + index}
                            </span>
                            <TextMark text={item} markText={markText} />
                          </p>
                        );
                      } else {
                        return (
                          <p
                            className="log_table_content__line"
                            key={index.toString()}>
                            <span className="line_number">
                              {lineStart + index}
                            </span>
                            <TextMark text={item} markText={markText} />
                          </p>
                        );
                      }
                    })
                  ) : (
                    cur_file &&
                    !fileLoading && (
                      <div className="log_content" id="log_content">
                        <p className="log_table_content_p1">
                          <img
                            src={require('./img/file.png')}
                            className="textPre_img"
                          />
                        </p>
                        <p className="log_table_content_p2">
                          当前日志文件下暂无具体日志!
                        </p>
                      </div>
                    )
                  )}
                  {sourceText.length ? (
                    <Spin
                      spinning={fileLoading && cur_type === 'down'}
                      style={{ position: 'relative', paddingLeft: '49%' }}
                    />
                  ) : (
                    ''
                  )}
                </div>
              </Spin>
            </div>
          </div>
        </Content>
      </Layout>
    );
  }
}
export default Form.create<any>()(Log);
