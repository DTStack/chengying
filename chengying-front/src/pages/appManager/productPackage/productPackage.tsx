import * as React from 'react';
import {
  Select,
  message,
  Input,
  Icon,
  Upload,
  Button,
  Dropdown,
  Menu,
  Modal,
  Form,
  Badge,
  Progress,
  Row,
  Col,
  List,
} from 'antd';
import { uniqBy } from 'lodash';
import { Service, servicePageService } from '@/services';
import ComponentList from './componentList';
import '../style.scss';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import utils from '@/utils/utils';
import { FormComponentProps } from 'antd/lib/form';
import axios, { AxiosResponse } from 'axios';
const Option = Select.Option;
const Search = Input.Search;
let time = null;
const innerAxios = {};
// interface CreateNoticeModalProps extends FormComponentProps {
//     isShow: boolean
//     onCancel: any
//     onOk: any
// }
interface Props extends FormComponentProps {
  location?: any;
  history?: any;
  extraContent?: any;
  authorityList: any;
  keycode?: any;
}

export interface QueryParams {
  clusterId?: number;
  parentProductName?: string;
  productName?: string[];
  productVersion?: string;
  deploy_status?: any;
  start?: number;
  limit?: number;
  'sort-by'?: string;
  'sort-dir'?: string;
}

interface State {
  productList: any[]; // 产品列表
  /** 组件列表 */
  componentList: any[];
  loading: boolean;
  searchParam: QueryParams;
  componentData: {
    list: any[];
    count: number;
  };
  visible: boolean;
  fileList: any;
  isShowPag?: boolean;
  interFileList: any;
  // innerAxios: any;
  // inputText:string;
}
const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});
@(connect(mapStateToProps, undefined) as any)
class ComponentContainer extends React.Component<Props, State> {
  state: State = {
    productList: [],
    componentList: [],
    loading: false,
    searchParam: {
      clusterId: 0,
      productName: undefined,
      parentProductName: undefined,
      productVersion: undefined,
      deploy_status: '',
      'sort-by': 'create_time',
      'sort-dir': 'desc',
      limit: 10,
      start: 0,
    },
    componentData: {
      list: [],
      count: 0,
    },
    visible: false,
    fileList: [],
    isShowPag: false,
    interFileList: [],
    // innerAxios: {}
    // inputText:'asdfghjk',
  };

  componentDidMount() {
    this.getParentProductsList();
  }

  // 获取产品
  getParentProductsList = (parentProductName?: string) => {
    Service.getParentProductList({ limit: 0 }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const data = res.data;
        this.setState({
          productList: data,
        });
        if (data && data.length > 0) {
          this.getProductComponents({
            parentProductName: parentProductName || data[0],
          });
        } else {
          this.resetPages();
        }
      } else {
        message.error(res.msg);
      }
    });
  };

  // 重置页面数据
  resetPages = () => {
    this.setState({
      productList: [],
      componentList: [],
      searchParam: {
        clusterId: 0,
        productName: undefined,
        parentProductName: undefined,
        productVersion: undefined,
        deploy_status: '',
        'sort-by': 'create_time',
        'sort-dir': 'desc',
        limit: 10,
        start: 0,
      },
      componentData: {
        list: [],
        count: 0,
      },
    });
  };

  // 获取表格信息（产品包组件信息）
  getDataList = (searchParam?: QueryParams) => {
    const reqParams: any = Object.assign(
      {},
      searchParam && Object.keys(searchParam).length
        ? searchParam
        : this.state.searchParam
    );
    if (!reqParams.parentProductName) {
      return;
    }

    reqParams.start = reqParams.start * reqParams.limit;
    if (reqParams.deploy_status) {
      reqParams.deploy_status = reqParams.deploy_status.join(',');
    }
    if (reqParams.productName) {
      reqParams.productName = reqParams.productName.join(',');
    }
    this.setState({
      loading: true,
    });
    Service.getAllProducts(reqParams).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const data = res.data;
        this.setState({
          componentData: data,
        });
      } else {
        message.error(res.msg);
      }
      this.setState({
        loading: false,
      });
    });
  };

  getProductComponents = (params?: any) => {
    const reqParams = Object.assign({}, this.state.searchParam, params);
    servicePageService
      .getProductName({
        clusterId: 0,
        parentProductName: params.parentProductName,
      })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          const data = res.data.list;
          let arr = [];
          arr = uniqBy(data, 'product_name');
          this.setState(
            {
              componentList: arr,
              searchParam: reqParams,
            },
            () => {
              this.getDataList();
            }
          );
        } else {
          message.error(res.msg);
        }
      });
  };

  handleSelectChange = (value: string) => {
    debugger;
    const newState = Object.assign({}, this.state.searchParam, {
      parentProductName: value,
      productName: undefined,
      productVersion: undefined,
    });
    this.setState({ searchParam: newState }, () => {
      this.getProductComponents(newState);
    });
  };

  handleComponentSearch = (value: any) => {
    const newState = Object.assign({}, this.state.searchParam, {
      productVersion: value,
    });
    this.setState({ searchParam: newState }, this.getDataList);
  };

  onComponentChange = (key) => {
    this.setState(
      {
        searchParam: Object.assign({}, this.state.searchParam, {
          productName: key,
          deploy_status: '',
        }),
      },
      this.getDataList
    );
  };

  renderProducts = () => {
    const { productList, searchParam } = this.state;
    const options =
      productList &&
      productList.map((item: any, index: number) => (
        <Option key={`${index}`} value={item}>
          <Icon type="appstore-o" style={{ marginRight: '6px' }} />
          {item}
        </Option>
      ));
    return (
      <Select
        className="dt-form-shadow-bg"
        style={{ width: 180 }}
        size="default"
        placeholder="选择产品"
        value={searchParam.parentProductName}
        onChange={this.handleSelectChange}>
        {options}
      </Select>
    );
  };

  // 上传网络安装包
  showModal = () => {
    const { authorityList } = this.props;
    if (!utils.noAuthorityToDO(authorityList, 'package_upload_delete')) {
      this.setState({
        visible: true,
      });
    }
  };

  hideModal = () => {
    this.setState({
      visible: false,
    });
    this.props.form.resetFields();
  };

  // 校验链接
  Verifylink = (params: any) => {
    Service.Verifylink(params).then((res: any) => {
      console.log(res);
      const { code, msg } = res.data;
      if (code === 0) {
        this.hideModal();
        this.addNetworkPackages(params);
        this.props.form.resetFields();
      } else {
        this.showModal();
        message.warning(msg);
      }
    });
  };

  addNetworkPackages = (params: any) => {
    Service.addNetworkPackage(params).then((res: any) => {
      const { msg, code } = res.data;
      // let imsg = msg.split(/[\n]/)
      // console.log(imsg)
      if (code === 0) {
        message.success(msg);
      } else {
        message.error(msg);
      }
      this.getDataList();
    });

    this.getProductpackageList();
  };

  getProductpackageList = () => {
    time = setInterval(() => {
      Service.getProductpackageList().then((res: any) => {
        const data = res.data.data.data;
        this.setState({
          interFileList: data,
        });
        if (data.length === 0) {
          clearInterval(time);
        }
      });
    }, 500);
  };

  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        const website = values.TextArea.split('\n');
        // this.addNetworkPackages({
        //     name: website
        // })
        this.Verifylink({
          name: website,
        });
      }
    });
  };

  // 删除正在上传文件
  delPackage = (index) => {
    const { fileList } = this.state;
    fileList.map((item, i) => {
      if (index === i) {
        item.status = 'removed';
      }
      return item;
    });
    this.setState({
      fileList,
    });
  };

  // 删除网络上传文件
  delInterPackage = (index, interFileList) => {
    console.log(interFileList[index].id);
    Service.delProductPackageItem({
      record_id: interFileList[index].id,
    }).then((res: any) => {
      console.log(res);
      const data = res.data.data;
      console.log(data);
      this.setState({
        interFileList: data ?? [],
      });
    });
  };

  // 微标数展示进度
  istrueShow = () => {
    this.setState({
      isShowPag: !this.state.isShowPag,
    });
  };

  isfalseShow = () => {
    this.setState(
      {
        isShowPag: false,
      },
      () => {
        this.state.isShowPag !== true
          ? clearInterval(time)
          : this.getProductpackageList();
      }
    );
  };

  handleChange = ({ event, file, fileList }) => {
    // 此处拦截error状态，erorr是event为空， 因为upload的success、error都会再次调用onchange事件，
    if (!event) return;
    // 清除error数据
    fileList = fileList.filter((item) => item.status != 'error');
    // 删除移除的
    if (file && file.status === 'removed') {
      fileList = fileList.filter((item) => item.status != 'removed');
      if (file.status === 'removed') {
        innerAxios[file.uid].abort();
      }
    }
    // 删除成功的
    if (file && file.percent === '100.00') {
      // status状态有问题，因此使用percent
      fileList = fileList.filter((item) => item.percent != '100.00');
      setTimeout(() => {
        this.getDataList();
      }, 1000);
    }
    fileList = fileList.filter((item) => item.status !== 'done');
    this.setState({ fileList });
  };

  // 自定义进度条
  Iporgrss = (fileList, interFileList) => {
    return (
      <div>
        <div
          style={{
            position: 'absolute',
            top: 115,
            right: 14,
            display: `${this.state.isShowPag === true ? 'block' : 'none'}`,
          }}>
          {fileList.length || interFileList.length !== 0 ? (
            <List
              itemLayout="vertical"
              size="small"
              split={false}
              style={{
                width: 300,
                height: 500,
                maxHeight: 500,
                overflow: 'hidden',
                overflowY: 'auto',
                backgroundColor: '#fff',
                borderRadius: '3px',
                boxShadow: '0px 0px 8px #ccc',
                paddingTop: '10px',
                zIndex: 9999,
              }}>
              <div
                style={{
                  position: 'absolute',
                  top: 8,
                  right: 10,
                  zIndex: 999,
                  display: `${
                    this.state.isShowPag === true ? 'block' : 'none'
                  }`,
                }}>
                <Icon onClick={this.isfalseShow} type="close" />
              </div>
              <div
                style={{
                  width: 80 + '%',
                  height: 30 + 'px',
                  textAlign: 'center',
                  lineHeight: 30 + 'px',
                  border: '2px solid #ffe58f',
                  borderRadius: 5,
                  backgroundColor: '#fffbe6',
                  margin: '0px auto',
                }}>
                <Icon
                  type="exclamation-circle"
                  theme="filled"
                  style={{ color: '#ffe58f', fontSize: 16 }}
                />
                &emsp;<span>安装包上传中，请勿刷新页面</span>
              </div>

              {fileList.length === 0
                ? ''
                : fileList.map((item, index) => {
                    return (
                      <List.Item key={index} style={{ paddingLeft: 20 }}>
                        <Row justify="center">
                          <Col span={23}>
                            <span style={{ fontSize: 10 }}>{item.name}</span>
                          </Col>
                        </Row>
                        <Row justify="start">
                          <Col span={2}>
                            <Icon
                              theme="outlined"
                              style={{ color: '#1890ff' }}
                              type={
                                item.percent !== 100
                                  ? 'loading'
                                  : 'check-circle'
                              }
                            />
                          </Col>
                          <Col span={18}>
                            <Progress
                              strokeWidth={3}
                              size="small"
                              status="active"
                              percent={parseInt(item.percent)}
                            />
                          </Col>
                          <Col span={2}>
                            <Icon
                              theme="outlined"
                              style={{ color: '#1890ff', paddingLeft: 10 }}
                              onClick={() => this.delPackage(index)}
                              type="delete"
                            />
                          </Col>
                        </Row>
                      </List.Item>
                    );
                  })}
              {interFileList.length === 0
                ? ''
                : interFileList.map((item, index) => {
                    return (
                      <List.Item key={index} style={{ paddingLeft: 20 }}>
                        <Row justify="center">
                          <Col span={23}>
                            <span style={{ fontSize: 10 }}>{item.name}</span>
                          </Col>
                        </Row>
                        <Row justify="start">
                          <Col span={2}>
                            <Icon
                              theme="outlined"
                              style={{ color: '#1890ff' }}
                              type={
                                item.percent !== 100
                                  ? 'loading'
                                  : 'check-circle'
                              }
                            />
                          </Col>
                          <Col span={18}>
                            <Progress
                              strokeWidth={3}
                              size="small"
                              status="active"
                              percent={parseInt(item.percent)}
                            />
                          </Col>
                          <Col span={2}>
                            <Icon
                              theme="outlined"
                              style={{ color: '#1890ff', paddingLeft: 10 }}
                              onClick={() =>
                                this.delInterPackage(index, interFileList)
                              }
                              type="delete"
                            />
                          </Col>
                        </Row>
                      </List.Item>
                    );
                  })}
            </List>
          ) : (
            ''
          )}
        </div>
      </div>
    );
  };

  customRequest = (options) => {
    const {
      action,
      file,
      filename,
      headers,
      onError,
      onProgress,
      onSuccess,
      withCredentials,
    } = options;
    const formData = new FormData();
    formData.append(filename, file);
    var CancelToken = axios.CancelToken;
    var source = CancelToken.source();

    axios
      .post(action, formData, {
        withCredentials,
        headers,
        onUploadProgress: ({ total, loaded }) => {
          onProgress({
            percent: Math.round((loaded / total) * 100).toFixed(2),
            name: file.name,
            uid: file.uid,
          });
        },
        cancelToken: source.token,
      })
      .then(({ data }): any => {
        if (data?.code !== 0) {
          return message.error(data.msg);
        }
        onSuccess(data, file);
      })
      .catch((onError) => {});
    innerAxios[file.uid] = {
      abort() {
        source.cancel();
      },
    };
  };

  render() {
    const {
      componentList,
      searchParam,
      componentData,
      visible,
      fileList,
      interFileList,
    } = this.state;
    const { authorityList, history } = this.props;
    const paneKey = `${searchParam.parentProductName}-${searchParam.productName}-${searchParam.productVersion}`;
    const uploadCfg = {
      name: 'package',
      action: '/api/v2/product/upload',
      accept: '.tar',
      onChange: this.handleChange,
      showUploadList: false,
      openFileDialogOnClick: !!authorityList.package_upload_delete,
      customRequest: (options) => this.customRequest(options),
    };
    const menu = (
      <div>
        <Menu>
          <Menu.Item>
            <Upload {...uploadCfg}>
              <a
                onClick={() =>
                  utils.noAuthorityToDO(authorityList, 'package_upload_delete')
                }>
                上传本地包
              </a>
            </Upload>
          </Menu.Item>
          <Menu.Item>
            <a
              onClick={this.showModal}
              style={{ color: '#3F87FF' }}
              className="ant-upload ant-upload-select ant-upload-select-text">
              来自网络
            </a>
          </Menu.Item>
        </Menu>
      </div>
    );

    const { getFieldDecorator } = this.props.form;
    return (
      <div className="app-manager-container product-package-container">
        {/* <p className="title mb-12">安装包管理</p> */}
        <div className="clearfix mb-12">
          <React.Fragment>
            <span className="mr-20">产品：{this.renderProducts()}</span>
            <span className="mr-20">
              组件名称：
              <Select
                className="dt-form-shadow-bg"
                mode="multiple"
                size="default"
                placeholder="选择组件"
                value={this.state.searchParam.productName}
                style={{ width: 180 }}
                onChange={this.onComponentChange}>
                {componentList &&
                  componentList.map((o: any) => (
                    <Option
                      data-testid={`option-${o.product_name_display}`}
                      key={o.product_name}>
                      {o.product_name}
                    </Option>
                  ))}
              </Select>
            </span>
            <span className="mr-20">
              <Search
                className="dt-form-shadow-bg"
                style={{ width: 264 }}
                placeholder="按组件版本号搜索"
                onSearch={this.handleComponentSearch}
              />
            </span>
          </React.Fragment>
          <div className="fl-r">
            <a onClick={this.istrueShow}>
              <Badge
                count={
                  fileList.length || interFileList.length !== 0
                    ? fileList.length + interFileList.length
                    : 0
                }>
                <Dropdown overlay={menu} placement="bottomLeft">
                  <Button type="primary">上传组件安装包</Button>
                </Dropdown>
              </Badge>
            </a>
          </div>
        </div>
        <ComponentList
          key={'component-list' + paneKey}
          {...this.state.searchParam}
          location={location}
          history={history}
          componentData={componentData}
          getDataList={this.getDataList}
          getParentProductsList={this.getParentProductsList}
        />
        {this.Iporgrss(fileList, interFileList)}
        <Modal
          title="来自网络"
          visible={visible}
          onOk={this.handleSubmit}
          onCancel={this.hideModal}>
          <Form name="basic" labelCol={{ span: 4 }} wrapperCol={{ span: 20 }}>
            <Form.Item label="地址">
              {getFieldDecorator('TextArea', {
                rules: [{ required: true, message: '请输入地址' }],
              })(
                <Input.TextArea placeholder="请输入安装包地址，多个地址间需换行输入" />
              )}
            </Form.Item>
          </Form>
        </Modal>
      </div>
    );
  }
}
export default Form.create<Props>()(ComponentContainer);
