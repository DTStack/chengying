import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators, Dispatch } from 'redux';
import * as DashboardAction from '@/actions/dashBoardAction';
import { Layout, Button, Input, Select, Form, Modal, message } from 'antd';
import { dashboardService } from '@/services';
import { AppStoreTypes } from '@/stores';
import AddDashModal from './addDashModal';
import ImportDashModal from './importDashModal';
import DashListComp from './list';
import { DashBoardStore, HeaderStoreType } from '@/stores/modals';
import './style.scss';
import * as Cookies from 'js-cookie';

const { Content } = Layout;
const Search = Input.Search;
const Option = Select.Option;
const FormItem = Form.Item;

const mapStateToProps = (state: AppStoreTypes) => ({
  dashboard: state.DashBoardStore,
  headerStore: state.HeaderStore,
  authorityList: state.UserCenterStore.authorityList,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, DashboardAction), dispatch),
});

interface DashBoardProps {
  dashboard: DashBoardStore;
  authorityList: any;
  headerStore: HeaderStoreType;
  actions: DashboardAction.DashBoardActionsTypes;
  match?: any;
}
interface DashBoardState {
  selectTags: any[];
  addModalVisible: boolean;
  addFolderVisible: boolean;
  newFolderName: string;
  newFolderError: string;
  newFolderStatus: any;
  importModalVisible: boolean;
  addDashFinish: boolean;
  allTags: any[];
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
export default class DashListPage extends React.Component<
  DashBoardProps,
  DashBoardState
> {
  constructor(props: any) {
    super(props as DashBoardProps);
  }

  state: DashBoardState = {
    selectTags: [],
    addModalVisible: false,
    addFolderVisible: false,
    newFolderName: '',
    newFolderError: '',
    newFolderStatus: 'success',
    importModalVisible: false,
    addDashFinish: false,
    allTags: [],
  };

  componentDidMount() {
    const { headerStore } = this.props;
    const defaultProduct =
      headerStore?.cur_parent_product != '选择产品'
        ? headerStore?.cur_parent_product
        : Cookies.get('em_current_parent_product');
    this.setState(
      {
        selectTags: defaultProduct ? [defaultProduct] : [],
      },
      () => {
        this.getDashList();
      }
    );
  }

  /**
   * 更新仪表盘列表
   * @memberof DashListPage
   */
  getDashList = () => {
    this.props.actions.getDashboardList({
      mode: 'tree',
      query: '',
      skipRecent: true,
      skipStarred: true,
      starred: false,
      tags: this.state.selectTags,
    });
  };

  /**
   * 根据tag筛选仪表盘
   * @memberof DashListPage
   */
  handleTagFilterChange = (tagIndexs: any[]) => {
    this.setState(
      {
        selectTags: tagIndexs,
      },
      () => {
        this.getDashList();
      }
    );
  };

  /**
   * 根据dashboard名称搜索结果
   * @params(value) string 搜索输入框内容
   * @memberof DashListPage
   */
  handleNameFilterChange = (value: any) => {
    this.props.actions.getDashboardList({
      mode: 'tree',
      query: value,
      skipRecent: true,
      skipStarred: true,
      starred: false,
      tags: this.state.selectTags,
    });
  };

  /**
   * 新增dashboard，弹出弹窗
   * @memberof DashListPage
   */
  showAddModal = (canEdit) => {
    if (!canEdit) {
      message.error('权限不足,请联系管理员！')
      return
    }
    this.setState({
      addModalVisible: true,
    });
  };

  /**
   * 新增dashboard，关闭弹窗
   * @memberof DashListPage
   */
  hideAddModal = () => {
    this.setState({
      addModalVisible: false,
      addDashFinish: false,
    });
    this.getDashList();
  };

  // 新增文件夹，打开/关闭弹出窗口
  showAddFolderModal = () => {
    this.setState({
      addFolderVisible: true,
    });
  };

  hideAddFolderModal = () => {
    this.setState({
      addFolderVisible: false,
      newFolderName: '',
      newFolderError: '',
      newFolderStatus: '',
    });
  };

  handleAddFolderNameChange = (e: any) => {
    this.setState({
      newFolderName: e.target.value,
    });
    this.validateAddFolder(e.target.value);
  };

  validateAddFolder = (newFolderName: string) => {
    const { dashboards } = this.props.dashboard;
    let validate = true;
    let isExist = false;
    if (newFolderName === '') {
      this.setState({
        newFolderError: '请填写文件夹名称',
        newFolderStatus: 'error',
      });
      validate = false;
    } else {
      for (const f of dashboards) {
        if (f.title === newFolderName) {
          this.setState({
            newFolderError: '文件夹已经存在',
            newFolderStatus: 'error',
          });
          isExist = true;
          validate = false;
        }
      }
      if (!isExist) {
        this.setState({
          newFolderError: '',
          newFolderStatus: 'success',
        });
        validate = true;
      }
    }
    return validate;
  };

  handleAddFolderOk = () => {
    const self = this;
    const { newFolderName } = this.state;
    if (this.validateAddFolder(newFolderName)) {
      dashboardService
        .createDashFolder({
          uid: null,
          title: newFolderName,
        })
        .then((rst: any) => {
          self.getDashList();
          self.hideAddFolderModal();
        });
    }
  };

  // 导入dashboard，打开/关闭弹出窗口
  showImportModal = (canEdit) => {
    if (!canEdit) {
      message.error('权限不足, 请联系管理员！')
      return;
    }
    this.setState({
      importModalVisible: true,
    });
  };

  hideImportModal = () => {
    this.setState({
      importModalVisible: false,
    });
    this.getDashList();
  };

  render() {
    const { dashboard, authorityList } = this.props;
    const {
      addModalVisible,
      addFolderVisible,
      newFolderName,
      newFolderError,
      newFolderStatus,
      importModalVisible,
      addDashFinish,
    } = this.state;
    // const { cur_parent_product } = headerStore;
    const { tags } = dashboard;
    const CAN_EDIT = authorityList?.sub_menu_service_grafana_edit
    // const defaultProduct = cur_parent_product != "选择产品" ? cur_parent_product : (Cookies.get('em_current_parent_product')||undefined);
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
      <Layout style={{ minHeight: document.body.clientHeight - 88 }}>
        <Content>
          <div className="dash-page">
            <div className="top-navbar mb-12 clearfix">
              <Search
                className="dt-form-shadow-bg mr-20"
                placeholder="按名称搜索"
                onSearch={this.handleNameFilterChange}
                style={{ width: 264, float: 'left' }}
              />
              <div className="tag-filter clearfix">
                <div className="filter-title box-shadow-style">Tags</div>
                <Select
                  mode="multiple"
                  className="dt-form-shadow-bg"
                  style={{ minWidth: 264, float: 'left' }}
                  placeholder="按Tag筛选"
                  optionFilterProp="children"
                  // defaultValue={[defaultProduct]}
                  onChange={this.handleTagFilterChange}
                  filterOption={(input, option) =>
                    (option.props.children as any)
                      .toLowerCase()
                      .indexOf(input.toLowerCase()) >= 0
                  }>
                  {tags.map((tag: any, index: number) => {
                    return (
                      <Option key={`${index}`} value={tag}>
                        {tag}
                      </Option>
                    );
                  })}
                </Select>
              </div>
              {/* <Button type="primary" style={{ float: 'right' }} onClick={this.showAddFolderModal}>+ 文件夹</Button> */}
              <Button
                type="primary"
                style={{ float: 'right', marginLeft: 12 }}
                onClick={() => this.showImportModal(CAN_EDIT)}>
                + 导入
              </Button>
              <Button
                type="primary"
                style={{ float: 'right' }}
                onClick={() => {this.showAddModal(CAN_EDIT)}}>
                + 仪表盘
              </Button>
            </div>
            <DashListComp
              canEdit={CAN_EDIT}
              onChange={this.getDashList}
              list={dashboard.dashboards}
            />
            <AddDashModal
              visible={addModalVisible}
              onClose={this.hideAddModal}
              finish={addDashFinish}></AddDashModal>
            <ImportDashModal
              visible={importModalVisible}
              onClose={this.hideImportModal}></ImportDashModal>
            <Modal
              visible={addFolderVisible}
              title="创建文件夹"
              onOk={this.handleAddFolderOk}
              onCancel={this.hideAddFolderModal}>
              <Form>
                <FormItem
                  {...formItemLayout}
                  required
                  label="Name"
                  validateStatus={newFolderStatus}
                  help={newFolderError}>
                  <Input
                    value={newFolderName}
                    onChange={this.handleAddFolderNameChange}
                    style={{ width: 400 }}
                  />
                </FormItem>
              </Form>
            </Modal>
          </div>
        </Content>
      </Layout>
    );
  }
}
