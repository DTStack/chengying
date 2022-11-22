import * as React from 'react';
import {
  Card,
  Input,
  Button,
  Icon,
  message,
  Tooltip,
  Table,
  Switch,
  Pagination,
  Modal,
  Badge,
} from 'antd';
import { RESULT_STATUS, RESULT_FILTER } from '../const';
import { AppStoreTypes } from '@/stores';
import { connect } from 'react-redux';

import SettingModal from './settingModal';
import PreviewScript from './previewScript';
import TaskHistory from './taskHistory';
import UploadScript from './uploadScript';
import { scriptManager } from '@/services';
import { EllipsisText } from 'dt-react-component';

import './style.scss';

const { Search } = Input;
const { confirm } = Modal;
interface IProps {
  authorityList: any;
}

let timer = null;

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});

const renderStatus = (type: any) => {
  switch (type) {
    case RESULT_STATUS.NORMAL:
      return (
        <span>
          <Badge color="#12BC6A" /> 正常
        </span>
      );
    case RESULT_STATUS.UNNORMAL:
      return (
        <span>
          <Badge color="#FF5F5C" /> 异常
        </span>
      );
    case RESULT_STATUS.RUN:
      return (
        <span>
          <Badge color="#3F87FF" /> 运行中
        </span>
      );
    case RESULT_STATUS.UNRUN:
      return (
        <span>
          <Badge color="#BFBFBF" /> 未运行
        </span>
      );
  }
};

interface IState {
  searchName: string;
  selectedRowKeys: number[];
  data: any[];
  total: number;
  visible: boolean;
  isTimeSetting: boolean;
  previewVisible: boolean;
  previewContent: string;
  scriptName: string;
  visibleHistory: boolean;
  spec: string; // cron表达式
  hosts: any[];
  taskId: number;
  tableLoading: boolean;
  start: number;
  execStatus: any;
  times: number;
  limit: number;
  showUploadScript: boolean;
  uploadTitle: string;
  detailInfo: any;
}
@(connect(mapStateToProps) as any)
export default class ScriptManager extends React.PureComponent<IProps, IState> {
  state: IState = {
    searchName: '',
    selectedRowKeys: [],
    data: [],
    total: 0,
    visible: false,
    isTimeSetting: false,
    previewVisible: false,
    previewContent: '',
    scriptName: '',
    visibleHistory: false,
    spec: '',
    hosts: [],
    taskId: 0,
    tableLoading: false,
    start: 0,
    execStatus: '',
    times: 0,
    limit: 10,
    showUploadScript: false,
    uploadTitle: '上传脚本',
    detailInfo: null,
  };

  componentDidMount() {
    this.getScriptList();
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.state.times !== prevState.times) {
      clearTimeout(timer);
      timer = setTimeout(() => {
        this.getScriptList();
      }, 5000);
    }
  }

  componentWillUnmount() {
    clearTimeout(timer);
  }

  getScriptList = () => {
    const { searchName, start, execStatus, times, limit } = this.state;
    this.setState({ tableLoading: true });
    const params = {
      name: searchName,
      'sort-by': 'create_time',
      'sort-dir': 'desc',
      limit,
      start,
      'exec-status': execStatus,
    };
    scriptManager.getList(params).then((res: any) => {
      res = res.data;
      const { code, data = {}, msg } = res;
      if (code === 0) {
        data.list?.map((item) => {
          item.key = item.id;
          return item;
        });
        this.setState({
          tableLoading: false,
          data: data.list || [],
          times: times + 1,
          total: data.count,
        });
      } else {
        this.setState({ tableLoading: false });
        message.error(msg);
      }
    });
  };

  handleUploadFile = (isEdit?: boolean) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.setState({ showUploadScript: true });
  };

  // 按脚本名称搜索
  doSearch = (value) => {
    this.setState({ searchName: value, start: 0, execStatus: '' }, () => {
      this.getScriptList();
    });
  };

  // 改变页码数
  onChangePage = (page: number) => {
    this.setState({ start: (page - 1) * 10, tableLoading: true }, () => {
      this.getScriptList();
    });
  };

  // 点击脚本名称
  clickName = (record) => {
    this.setState({
      taskId: record.id,
      previewVisible: true,
      scriptName: record.name,
    });
  };

  // 点击执行历史
  previewHistory = (record, isEdit) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.setState({
      taskId: record.id,
      scriptName: record.name,
      visibleHistory: true,
    });
  };

  // 表格配置
  initColumns: any = () => {
    return [
      {
        title: '脚本名称',
        dataIndex: 'name',
        key: 'name',
        render: (text: any, record) => {
          return (
            <div
              onClick={() => {
                this.clickName(record);
              }}
              style={{ cursor: 'pointer', color: '#0091FF' }}>
              {text.length > 20 ? (
                <EllipsisText value={text} maxWidth={200} />
              ) : (
                <span>{text}</span>
              )}
            </div>
          );
        },
      },
      {
        title: '脚本描述',
        dataIndex: 'describe',
        key: 'describe',
        width: 200,
        render: (text: string) => {
          return (
            <div>
              {text.length > 20 ? (
                <EllipsisText value={text} maxWidth={200} />
              ) : (
                <span>{text || '--'}</span>
              )}
            </div>
          );
        },
      },
      {
        title: '上传时间',
        dataIndex: 'create_time',
        key: 'create_time',
        width: 150,
      },
      {
        title: '最近执行时间',
        dataIndex: 'end_time',
        key: 'end_time',
        width: 150,
        render: (text: string) => text || '--',
      },
      {
        title: '最近一次执行结果',
        dataIndex: 'exec_status',
        key: 'exec_status',
        width: 135,
        filters: RESULT_FILTER,
        filterMultiple: false,
        render: renderStatus,
      },
      {
        title: '定时状态',
        dataIndex: 'status', // 0 关闭 1 开启
        key: 'status',
        render: (text: any, record: any) => {
          return (
            <Switch
              checked={text === 1}
              defaultChecked={text === 1}
              onChange={() => {
                this.changeSwitch(record);
              }}
            />
          );
        },
      },
      {
        title: '操作',
        dataIndex: 'operator',
        width: '280px',
        render: this.renderOperator,
      },
    ];
  };

  // 表格选择
  handleSelectChange = (selectedRowKeys: number[]) => {
    this.setState({ selectedRowKeys });
  };

  // 表格筛选
  handleChangeTable = (pagination, filters) => {
    let data = filters.exec_status;
    this.setState({ execStatus: data[0] }, () => {
      this.getScriptList();
    });
  };

  // 开关变化
  changeSwitch = (record) => {
    this.changeStatus(String(record.id), record.status === 1 ? 0 : 1);
  };

  // 批量操作
  switchAll = (status: number, isEdit: boolean) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    const { selectedRowKeys } = this.state;
    if (selectedRowKeys.length === 0) {
      message.error('请先选择脚本');
      return;
    }
    let ids = selectedRowKeys.join(',');
    this.changeStatus(ids, status);
  };

  // 修改定时状态
  changeStatus = (ids: string, status: number) => {
    const param = {
      task_id: ids,
      status,
    };
    scriptManager.getTaskStatus(param).then((res: any) => {
      if (res.data.code == 0) {
        this.getScriptList();
        this.setState({ selectedRowKeys: [] });
      } else {
        message.error(res.data.msg);
      }
    });
  };

  // 关闭设置弹框
  close = () => {
    this.setState({ visible: false });
  };

  // 打开定时设置弹框
  settingTime = (record, isEdit) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.setState({
      visible: true,
      taskId: record.id,
      hosts: record?.hosts,
      spec: record?.spec,
      isTimeSetting: true,
    });
  };

  // 打开手动执行
  handleSetting = (record, isEdit) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    this.setState({
      visible: true,
      isTimeSetting: false,
      hosts: record?.hosts,
      spec: record?.spec,
      taskId: record.id,
    });
  };

  // 关闭脚本预览
  closePreviewScript = () => {
    this.setState({ previewVisible: false });
  };

  // 关闭执行历史
  closeTaskHistory = () => {
    this.setState({ visibleHistory: false });
  };

  // 删除脚本
  deleteTask = (id) => {
    scriptManager.deleteTask({ id }).then((res: any) => {
      if (res.data.code == 0) {
        message.success(res.data.msg);
        this.getScriptList();
      } else {
        message.error(res.data.msg);
      }
    });
  };

  // 删除脚本
  deleteScript = (record, isEdit) => {
    if (!isEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    confirm({
      title: `确认删除脚本 ${record.name}？`,
      icon: (
        <Icon type="close-circle" theme="filled" style={{ color: '#FF5F5C' }} />
      ),
      content: '删除脚本将同时清除执行历史。',
      okText: '删除',
      cancelText: '取消',
      onOk: () => {
        this.deleteTask(record.id);
      },
      onCancel() {},
    });
  };

  // 关闭上传脚本弹框
  closeUpload = () => {
    this.setState({
      showUploadScript: false,
      detailInfo: null,
      uploadTitle: '上传脚本',
    });
  };

  // 上传脚本后成功回调
  callUpload = () => {
    this.closeUpload();
    this.getScriptList();
  };

  // 编辑脚本
  showUploadModal = (info, id) => {
    this.setState({ detailInfo: info, previewVisible: false }, () => {
      this.setState({
        uploadTitle: '编辑',
        showUploadScript: true,
        taskId: id,
      });
    });
  };

  // 表格操作按钮
  renderOperator = (text: any, record: any) => {
    const { authorityList } = this.props;
    const isView = authorityList.sub_menu_platform_manager_scriptManager_view;
    const isEdit = authorityList.sub_menu_configuration_platformsecurity_edit;

    if (!isView) return <></>;
    return (
      <div className="btnList">
        <Button
          type="link"
          onClick={() => {
            this.settingTime(record, isEdit);
          }}>
          定时设置
        </Button>
        <div className="tableLine">
          <span></span>
        </div>
        <Button
          type="link"
          disabled={record?.run_status === 1}
          onClick={() => {
            this.handleSetting(record, isEdit);
          }}>
          手动执行
        </Button>
        <div className="tableLine">
          <span></span>
        </div>
        <Button
          type="link"
          onClick={() => {
            this.previewHistory(record, isEdit);
          }}>
          执行历史
        </Button>
        <div className="tableLine">
          <span></span>
        </div>
        <Button
          type="link"
          onClick={() => {
            this.deleteScript(record, isEdit);
          }}>
          删除
        </Button>
      </div>
    );
  };

  render() {
    const { authorityList } = this.props;
    const isView = authorityList.sub_menu_platform_manager_scriptManager_view;
    const isEdit = authorityList.sub_menu_configuration_platformsecurity_edit
      ? true
      : false;
    const {
      tableLoading,
      previewVisible,
      previewContent,
      selectedRowKeys,
      visibleHistory,
      data,
      total,
      visible,
      scriptName,
      limit,
      spec,
      hosts,
      taskId,
      isTimeSetting,
      start,
      showUploadScript,
      uploadTitle,
      detailInfo,
    } = this.state;

    const rowSelection = {
      selectedRowKeys,
      onChange: this.handleSelectChange,
    };
    // 计算当前页
    const current: number = start / limit + 1;

    return (
      <div className="scriptManager">
        {/* 上传脚本弹框 */}
        {showUploadScript && (
          <UploadScript
            id={taskId}
            onClose={this.closeUpload}
            onOk={this.callUpload}
            title={uploadTitle}
            showVisible={showUploadScript}
            detailInfo={detailInfo}
          />
        )}
        {/* 定时/手动操作弹框 */}
        {visible && (
          <SettingModal
            title={isTimeSetting ? '定时设置' : '手动执行'}
            visible={visible}
            isTimeSetting={isTimeSetting}
            close={this.close}
            spec={spec}
            hosts={hosts}
            id={taskId}
            sucCall={this.getScriptList}
          />
        )}
        {/* 查看脚本 */}
        {previewVisible && (
          <PreviewScript
            showUploadModal={this.showUploadModal}
            id={taskId}
            title={`脚本查看/${scriptName}`}
            visible={previewVisible}
            content={previewContent}
            close={this.closePreviewScript}
          />
        )}
        {/* 执行历史 */}
        {visibleHistory && (
          <TaskHistory
            title={`执行历史/${scriptName}`}
            visible={visibleHistory}
            id={taskId}
            close={this.closeTaskHistory}
          />
        )}
        <Card bordered={false} className="card-box">
          <div className="card-content">
            <div className="scriptManager-nav">
              <div className="scriptManage-searchBox">
                <Search
                  placeholder="请输入脚本名称/描述"
                  onSearch={this.doSearch}
                  style={{
                    width: 340,
                    boxShadow: '0px 2px 6px 0px rgba(0, 0, 0, 0.08)',
                  }}
                />
              </div>
              <div className="scriptManage-uploadBox">
                {isView && (
                  <Button
                    type="primary"
                    onClick={() => this.handleUploadFile(isEdit)}>
                    上传脚本
                  </Button>
                )}
                <Tooltip placement="topRight" title="仅支持 .py，.sh 格式文件">
                  <div style={{ width: 35, paddingLeft: '7px', marginTop: 8 }}>
                    <Icon
                      type="question-circle"
                      style={{
                        fontSize: 16,
                        color: '#999',
                      }}
                    />
                  </div>
                </Tooltip>
              </div>
            </div>
            <div className="tabList">
              <Table
                scroll={{ y: 'calc(100vh - 274px)' }}
                style={{ height: 'calc(100vh - 230px)' }}
                loading={tableLoading}
                rowSelection={rowSelection}
                columns={this.initColumns()}
                dataSource={data}
                onChange={this.handleChangeTable}
                pagination={false}
              />
            </div>
            <div className="tabFooter">
              {isView && (
                <div className="tabFooter-btnList">
                  <Button
                    type="primary"
                    style={{ marginRight: 10 }}
                    onClick={() => this.switchAll(1, isEdit)}>
                    开启
                  </Button>
                  <Button onClick={() => this.switchAll(0, isEdit)}>
                    关闭
                  </Button>
                </div>
              )}
              <div>
                <Pagination
                  current={current}
                  size="small"
                  total={total}
                  onChange={this.onChangePage}
                  showTotal={(total) => (
                    <span>
                      共<span style={{ color: '#3F87FF' }}>{total}</span>
                      条数据，每页显示
                      <span style={{ color: '#3F87FF' }}>10</span>条
                    </span>
                  )}
                />
              </div>
            </div>
          </div>
        </Card>
      </div>
    );
  }
}
