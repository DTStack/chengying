import * as React from 'react';
import {
  Table,
  message,
  Modal,
  Divider,
  Input,
  Tooltip,
  Icon,
  Button,
} from 'antd';
import { userCenterService, clusterManagerService } from '@/services';
import MemberModal from './memberModal';
import ResetPasswordModal from '@/pages/userCenter/components/resetPassword';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import utils from '@/utils/utils';
import { encryptStr, encryptSM } from '@/utils/password';
import './style.scss';

interface IProps {
  authorityList: any;
}

interface State {
  membersList: any[];
  showModal: boolean;
  eddittingAccount: any;
  pagiParam: {
    current: number;
    size: number;
  };
  reqParams: {
    status: string;
    role_id: string;
    'sort-by': string;
    'sort-dir': string;
  };
  total: number;
  searchUsername: string;
  memberModalVisible: boolean;
  clusterList: any[];
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});

@(connect(mapStateToProps, undefined) as any)
export default class Members extends React.Component<IProps, State> {
  state: State = {
    membersList: [],
    showModal: false,
    eddittingAccount: null,
    pagiParam: {
      current: 1,
      size: 10,
    },
    reqParams: {
      status: undefined,
      role_id: undefined,
      'sort-by': 'create_time',
      'sort-dir': 'desc',
    },
    total: 0,
    searchUsername: '',
    memberModalVisible: false,
    clusterList: [],
  };

  componentDidMount() {
    this.initData();
    this.getClusterLists();
  }

  getClusterLists = () => {
    const params = {
      type: '',
      'sort-by': 'id',
      'sort-dir': 'desc',
      limit: 0,
      start: 0,
    };
    clusterManagerService.getClusterLists(params).then((res: any) => {
      res = res.data;
      const { code, data = {}, msg } = res;
      if (code === 0) {
        this.setState({
          clusterList: data.clusters || [],
        });
      } else {
        message.error(msg);
      }
    });
  };

  // 权限控制
  authorityControl = (
    action: string,
    code: string,
    record?: any,
    status?: boolean
  ) => {
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, code)) {
      return;
    }
    switch (action) {
      case 'add':
        this.memberModalShow();
        break;
      default:
        break;
    }
  };

  // 获取成员列表
  initData = () => {
    const { reqParams, pagiParam, searchUsername } = this.state;
    userCenterService
      .getMembers({
        start: (pagiParam.current - 1) * pagiParam.size,
        limit: pagiParam.size,
        username: searchUsername,
        ...reqParams,
      })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.setState({
            membersList: res.data.list,
            total: res.data.count,
          });
        } else {
          message.error(res.msg);
        }
      });
  };

  handleRemoveMember = (re: any, event: any) => {
    event.preventDefault();
    Modal.confirm({
      title: '确定移除该成员？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        const p = {
          targetUserId: re.id,
        };
        userCenterService.removeMember(p).then((res: any) => {
          res = res.data;
          if (res.code === 0) {
            message.success('移除成功');
            this.setState(
              {
                pagiParam: {
                  ...this.state.pagiParam,
                  current: 1,
                },
              },
              this.initData
            );
          } else {
            message.error(res.msg);
          }
        });
      },
    });
  };

  resetPassword = (re: any, event: any) => {
    event.preventDefault();
    this.setState({
      eddittingAccount: re,
      showModal: true,
    });
  };

  handleResetSubmit = async (value: any) => {
    const publicKeyRes = await userCenterService.getPublicKey();
    if (publicKeyRes.data.code !== 0) {
      return;
    }
    const { encrypt_type, encrypt_public_key } = publicKeyRes.data.data;
    const p = {
      password: encrypt_type === 'sm2' ? encryptSM(value.newPass, encrypt_public_key) : encryptStr(value.newPass, encrypt_public_key),
      targetUserId: this.state.eddittingAccount.id,
    };
    userCenterService.resetPassword(p).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        message.success('修改成功');
        this.setState({
          showModal: false,
          eddittingAccount: null,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 重置密码弹窗关闭
  handleResetPwdModalCancel = () => {
    this.setState({ showModal: false });
  };

  // 启用禁用
  handleTaggleStatus = (e: any, a: boolean, event: any) => {
    event.preventDefault();
    Modal.confirm({
      title: `确定${a ? '启用' : '禁用'}该成员？`,
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: async () => {
        const response = await userCenterService.taggleStatus(
          { targetUserId: e.id },
          a
        );
        const res = response.data;
        if (res.code === 0) {
          message.success('设置成功');
          this.initData();
        } else {
          message.error(res.msg);
        }
      },
    });
  };

  // 表格分页，筛选，排序
  handleTableChange = (pagination: any, filters: any, sorter: any) => {
    const { pagiParam } = this.state;
    const reqParams = {
      ...this.state.reqParams,
    };
    if (Object.keys(filters).length) {
      reqParams.status =
        filters.status && filters.status.length
          ? filters.status.join(',')
          : undefined;
      reqParams.role_id =
        filters.role_name && filters.role_name.length
          ? filters.role_name.join(',')
          : undefined;
    }
    if (sorter.field) {
      reqParams['sort-dir'] = sorter.order === 'ascend' ? 'asc' : 'desc';
      reqParams['sort-by'] =
        sorter.field === 'UpdateTimeFormat' ? 'create_time' : sorter.field;
    }
    this.setState(
      {
        pagiParam: Object.assign({}, pagiParam, {
          current: pagination.current,
        }),
        reqParams,
      },
      () => {
        this.initData();
      }
    );
  };

  handleSearch = (v: string) => {
    this.setState(
      {
        searchUsername: v,
        pagiParam: {
          ...this.state.pagiParam,
          current: 1,
        },
      },
      this.initData
    );
  };

  // 创建账号
  memberModalShow = () => {
    const { memberModalVisible, eddittingAccount } = this.state;
    if (memberModalVisible) {
      if (eddittingAccount && Object.keys(eddittingAccount).length) {
        this.initData();
      } else {
        this.setState(
          {
            pagiParam: {
              ...this.state.pagiParam,
              current: 1,
            },
          },
          this.initData
        );
      }
    }
    this.setState({
      memberModalVisible: !memberModalVisible,
      eddittingAccount: null,
    });
  };

  // 编辑账号
  editMemberInfo = (record: any) => {
    this.setState({
      memberModalVisible: true,
      eddittingAccount: record,
    });
  };

  initColumns = () => {
    const { authorityList } = this.props;
    const tableCOl = [
      {
        title: '账号',
        dataIndex: 'username',
        width: '15%',
        render: (e: string) => e || '--',
      },
      {
        title: '姓名',
        dataIndex: 'full_name',
        width: 100,
        render: (e: string, record: any) => (
          <Tooltip title={e.replace(/\s+/g, '') === '' ? '--' : e}>
            <span className="name">
              {e.replace(/\s+/g, '') === '' ? '--' : e}
            </span>
          </Tooltip>
        ),
      },
      {
        title: '邮箱',
        dataIndex: 'email',
        width: '15%',
        render: (text: string) => text || '--',
      },
      {
        title: '手机号',
        dataIndex: 'phone',
        render: (text: string) => text || '--',
      },
      {
        title: '账号状态',
        dataIndex: 'status',
        filters: [
          { text: '启用', value: '0' },
          { text: '禁用', value: '1' },
        ],
        render: (e: any, record: any) => (e === 0 ? '启用' : '禁用'),
      },
      {
        title: '角色',
        dataIndex: 'role_name',
        filters: [
          { text: 'Administrator', value: '1' },
          { text: 'Cluster Operator', value: '2' },
          { text: 'Cluster Reader', value: '3' },
        ],
      },
      {
        title: '创建时间',
        dataIndex: 'UpdateTimeFormat',
        width: '15%',
        sorter: true,
      },
      {
        title: '操作',
        dataIndex: 'action',
        width: '18%',
        render: (e: any, record: any) => {
          const SHOULD_USER_DO = authorityList.user_able_disable;
          const CAN_EDIT = authorityList.user_edit;
          const CAN_DELETE = authorityList.user_delete;
          const CAN_RESET_PWD = authorityList.user_reset_password;
          if (
            ((!SHOULD_USER_DO && !CAN_EDIT && !CAN_DELETE) ||
              record.role_name === 'Administrator') &&
            !CAN_RESET_PWD
          ) {
            return '--';
          }
          return (
            <span>
              {record.role_name !== 'Administrator' && (
                <span>
                  {SHOULD_USER_DO && (
                    <span>
                      {record.status === 0 ? (
                        <a
                          onClick={this.handleTaggleStatus.bind(
                            this,
                            record,
                            false
                          )}>
                          禁用
                        </a>
                      ) : (
                        <a
                          onClick={this.handleTaggleStatus.bind(
                            this,
                            record,
                            true
                          )}>
                          启用
                        </a>
                      )}
                      <Divider type="vertical" />
                    </span>
                  )}
                  {CAN_EDIT && (
                    <span>
                      <a onClick={this.editMemberInfo.bind(this, record)}>
                        编辑
                      </a>
                      <Divider type="vertical" />
                    </span>
                  )}
                  {CAN_DELETE && (
                    <span>
                      <a onClick={this.handleRemoveMember.bind(this, record)}>
                        删除
                      </a>
                      <Divider type="vertical" />
                    </span>
                  )}
                </span>
              )}
              {CAN_RESET_PWD && (
                <a onClick={this.resetPassword.bind(this, record)}>重置密码</a>
              )}
            </span>
          );
        },
      },
    ];
    return tableCOl;
  };

  render() {
    const {
      membersList,
      total,
      pagiParam,
      showModal,
      eddittingAccount,
      memberModalVisible,
      clusterList,
    } = this.state;
    const tableCOl = this.initColumns();
    const pagination = {
      size: 'small',
      pageSize: pagiParam.size,
      current: pagiParam.current,
      total: total,
      showTotal: (total) => (
        <span>
          共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示
          {pagiParam.size}条
        </span>
      ),
    };
    return (
      <div className="members-container" style={{ display: 'flex' }}>
        <div style={{ width: '100%', padding: '20px' }}>
          <div className="members-container-header">
            <Input.Search
              className="dt-form-shadow-bg mb-12"
              onSearch={this.handleSearch}
              style={{ width: 264 }}
              placeholder="按账号/姓名搜索"
            />
            <div className="header-r">
              <Button
                type="primary"
                onClick={() => this.authorityControl('add', 'user_add')}>
                创建账号
              </Button>
            </div>
          </div>
          <Table
            rowKey="username"
            onChange={this.handleTableChange}
            className="dt-table-fixed-base"
            style={{ height: 'calc(100vh - 135px)' }}
            columns={tableCOl}
            dataSource={membersList}
            scroll={{ y: true }}
            pagination={pagination}
          />
        </div>
        <ResetPasswordModal
          userInfo={eddittingAccount}
          visible={showModal}
          onCancel={this.handleResetPwdModalCancel}
          onSubmit={this.handleResetSubmit}
        />
        {memberModalVisible && (
          <MemberModal
            visible={memberModalVisible}
            onCancel={this.memberModalShow}
            memberInfo={eddittingAccount}
            clusterList={clusterList}
          />
        )}
      </div>
    );
  }
}
