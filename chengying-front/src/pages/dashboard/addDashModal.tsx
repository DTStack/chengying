import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators, Dispatch } from 'redux';
import * as DashboardAction from '@/actions/dashBoardAction';
import { dashboardService } from '@/services';
import { AppStoreTypes } from '@/stores';

import { FormComponentProps } from 'antd/lib/form/Form';
import {
  Select,
  Input,
  Form,
  Modal,
  Button,
  Icon,
  Tag,
  message,
  Tooltip,
} from 'antd';

const FormItem = Form.Item;
const Option = Select.Option;

const mapStateToProps = (state: AppStoreTypes) => ({
  dashboard: state.DashBoardStore,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, DashboardAction), dispatch),
});
// declare type validateStatus = "success" | "error" | "warning" | "validating";
interface AddDashModalProp {
  dashboard?: any;
  finish: boolean;
  onClose: () => void;
  visible: boolean;
}
interface AddDashModalState {
  finish: boolean;
  isNewFolder: boolean;
  newFolderOk: boolean;
  folderId: any;
  folder: string;
  dashName: string;
  tags: any[];
  tagInputVisible: boolean;
  tagInputValue: string;
  url: string;
  //   nameValidateStatus: string;
  //   helpMsg: string;
  //   folderHelpMsg: string;
  //   folderValidateStatus: string;
}
@(connect(mapStateToProps, mapDispatchToProps) as any)
class AddDashModal extends React.Component<
  AddDashModalProp & FormComponentProps,
  AddDashModalState
> {
  constructor(props: any) {
    super(props);
  }

  state: AddDashModalState = {
    finish: this.props.finish || false,
    isNewFolder: false,
    newFolderOk: false,
    folderId: null,
    folder: '',
    dashName: '',
    tags: [],
    tagInputVisible: false,
    tagInputValue: '',
    url: '',
    // nameValidateStatus: "validating",
    // helpMsg: "",
    // folderHelpMsg: '',
    // folderValidateStatus: 'validating'
  };

  UNSAFE_componentWillReceiveProps(nextProps: AddDashModalProp) {
    if (nextProps.finish !== this.state.finish) {
      this.setState({
        finish: nextProps.finish,
        folderId: -1,
        folder: '',
        dashName: '',
        tags: [],
        tagInputVisible: false,
        tagInputValue: '',
      });
    }
  }

  handleAddNewFolder = () => {
    const self = this;
    const { folder } = this.state;
    dashboardService
      .createDashFolder({
        uid: null,
        title: folder,
      })
      .then((res: any) => {
        if (res.id) {
          self.setState({
            newFolderOk: true,
            folderId: res.id,
          });
        }
      });
  };

  handleSetExistFolder = () => {
    this.setState({
      isNewFolder: false,
      folder: '',
    });
  };

  handleFolderSelectChange = (value: string) => {
    if (value === '-1') {
      this.setState({
        isNewFolder: true,
        folder: '',
        newFolderOk: false,
        folderId: null,
      });
    } else {
      const folder = JSON.parse(value);
      this.setState({
        folder: folder.title,
        folderId: folder.id,
      });
    }
  };

  handleFolderInputChange = (e: any) => {
    this.setState({
      folder: e.target.value,
    });
  };

  handleDashNameChange = (e: any) => {
    // debugger;
    if (this.validateForm(e.target.value)) {
      return true;
    } else {
      return false;
    }
  };

  validDashName = (rule: any, value: any, callback: any) => {
    if (value) {
      console.log(this.validateForm(value));
      if (this.validateForm(value)) {
        callback();
      } else {
        callback('仪表盘已存在');
      }
    } else {
      callback();
    }
  };

  handleRemoveTag = (removedTag: any) => {
    const tags = this.state.tags.filter((tag) => tag !== removedTag);
    this.setState({ tags });
  };

  showTagInput = () => {
    this.setState({ tagInputVisible: true });
  };

  handleTagInputChange = (e: any) => {
    this.setState({
      tagInputValue: e.target.value,
    });
  };

  handleTagInputConfirm = () => {
    const state = this.state;
    const tagInputValue = state.tagInputValue;
    let tags = state.tags;
    if (tagInputValue && tags.indexOf(tagInputValue) === -1) {
      tags = [...tags, tagInputValue];
    }
    console.log(tags);
    this.setState({
      tags,
      tagInputVisible: false,
      tagInputValue: '',
    });
  };

  validateForm = (dashName: string) => {
    console.log('检测仪表盘名字');
    const { folderId } = this.state;
    const { dashboards } = this.props.dashboard;
    let isValidate = true;
    let matchDashs = null;
    for (const f of dashboards) {
      if (f.id === folderId) {
        matchDashs = f.list;
      }
    }
    if (matchDashs) {
      for (const d of matchDashs) {
        if (d.title === dashName) {
          isValidate = false;
        }
      }
    }
    return isValidate;
  };

  handleFormSubmit = (e: any) => {
    e.preventDefault();
    const self = this;
    const { folderId, dashName, tags } = this.state;
    const { form } = this.props;
    form.validateFields((err: boolean, value: any) => {
      if (err) {
        return;
      }
      dashboardService
        .createDashboard({
          dashboard: {
            id: null,
            uid: null,
            title: value.dashName_form,
            tags: tags,
            timezone: 'browser',
            version: '0',
          },
          folderId: JSON.parse(value.folderId_form).id,
          overwrite: false,
        })
        .then((rst: any) => {
          // debugger;
          if (rst.data.id) {
            self.setState({
              finish: true,
              url: rst.data.url,
              dashName: value.dashName_form,
            });
          } else {
            message.error(rst.message);
          }
        })
        .catch((err: any) => {
          // console.log(err)
          if (err.status === 'name-exists') {
            message.error('已存在相同名称的仪表盘');
          }
          // message.error(err);
        });
    });
    console.log(folderId, dashName);
    // if (this.validateForm(dashName)) {
    //   this.setState({
    //     nameValidateStatus: "validating",
    //     folderValidateStatus: 'validating',
    //     folderHelpMsg: '',
    //     helpMsg: ""
    //   });
    // } else {
    //   this.setState({
    //     nameValidateStatus: "error",
    //     helpMsg: "仪表盘已存在！"
    //   });
    // }
  };

  handleModalClose = () => {
    this.setState({
      finish: false,
      isNewFolder: false,
      folderId: null,
      folder: '',
      dashName: '',
      tags: [],
      tagInputVisible: false,
      tagInputValue: '',
    });
    this.props.onClose();
  };

  gotoEdit() {
    window.location.href =
      '/deploycenter/monitoring/dashdetail?url=' +
      encodeURIComponent(this.state.url);
  }

  render() {
    const { visible, onClose, dashboard } = this.props;
    const { getFieldDecorator } = this.props.form;
    const {
      finish,
      isNewFolder,
      folder,
      tags,
      tagInputVisible,
      tagInputValue,
      newFolderOk,
      //   nameValidateStatus,
      //   helpMsg
    } = this.state;
    const formItemLayout = {
      labelCol: {
        xs: { span: 24 },
        sm: { span: 4 },
      },
      wrapperCol: {
        xs: { span: 24 },
        sm: { span: 20 },
      },
    };
    return (
      <Modal
        title="新增仪表盘"
        visible={visible}
        onCancel={onClose}
        footer={
          finish ? (
            <React.Fragment>
              <Button type="default" onClick={this.handleModalClose}>
                关闭
              </Button>
              <Button type="primary" onClick={this.gotoEdit.bind(this)}>
                去编辑
              </Button>
            </React.Fragment>
          ) : (
            <React.Fragment>
              <Button type="default" onClick={this.handleModalClose}>
                取消
              </Button>
              <Button type="primary" onClick={this.handleFormSubmit}>
                创建
              </Button>
            </React.Fragment>
          )
        }>
        <div className="add-dash-modal">
          {finish ? (
            <div className="add-success">
              <i className="emicon emicon-health_service" />
              <p className="success-word">
                仪表盘“{this.state.dashName}”创建成功!现在去编辑。
              </p>
            </div>
          ) : (
            <Form>
              <FormItem {...formItemLayout} label="Folder" required>
                {getFieldDecorator('folderId_form', {
                  rules: [{ required: true, message: '文件夹不能为空' }],
                })(
                  isNewFolder ? (
                    <div>
                      <Input
                        value={folder}
                        onChange={this.handleFolderInputChange}
                        style={{ width: 340 }}
                      />
                      {!newFolderOk && (
                        <Button
                          type="default"
                          style={{ marginLeft: 10 }}
                          onClick={this.handleAddNewFolder}>
                          <Icon type="check" />
                        </Button>
                      )}
                      <Button type="default" style={{ marginLeft: 10 }}>
                        <Icon
                          type="close"
                          onClick={this.handleSetExistFolder}
                        />
                      </Button>
                    </div>
                  ) : (
                    <Select
                      style={{ width: 340 }}
                      value={folder}
                      onChange={this.handleFolderSelectChange}>
                      {/* <Option value="-1">- New Folder -</Option> */}
                      {dashboard.dashboards.map((folder: any, index: any) => {
                        return (
                          <Option key={index} value={JSON.stringify(folder)}>
                            {folder.title}
                          </Option>
                        );
                      })}
                    </Select>
                  )
                )}
              </FormItem>
              <FormItem {...formItemLayout} label="Name" required>
                {getFieldDecorator('dashName_form', {
                  rules: [
                    { required: true, message: '仪表盘名称不能为空' },
                    { validator: this.validDashName },
                  ],
                })(
                  <Input
                    placeholder="输入仪表盘名称"
                    onChange={this.handleDashNameChange}
                    style={{ width: 340 }}
                  />
                )}
              </FormItem>
              <FormItem {...formItemLayout} label="Tags">
                <div>
                  {tags.map((tag) => {
                    const isLongTag = tag.length > 10;
                    const tagElem = (
                      <Tag
                        key={tag}
                        color="blue"
                        closable
                        afterClose={() => this.handleRemoveTag}>
                        {isLongTag ? `${tag.slice(0, 10)}...` : tag}
                      </Tag>
                    );
                    return isLongTag ? (
                      <Tooltip title={tag} key={tag}>
                        {tagElem}
                      </Tooltip>
                    ) : (
                      tagElem
                    );
                  })}
                  {tagInputVisible ? (
                    <Input
                      // ref={this.saveInputRef}
                      type="text"
                      size="small"
                      style={{ width: 78 }}
                      value={tagInputValue}
                      onChange={this.handleTagInputChange}
                      onBlur={this.handleTagInputConfirm}
                      onPressEnter={this.handleTagInputConfirm}
                    />
                  ) : (
                    <Tag
                      onClick={this.showTagInput}
                      style={{ background: '#fff', borderStyle: 'dashed' }}>
                      <Icon type="plus" /> New Tag
                    </Tag>
                  )}
                </div>
              </FormItem>
            </Form>
          )}
        </div>
      </Modal>
    );
  }
}
export default Form.create<any>()(AddDashModal);
