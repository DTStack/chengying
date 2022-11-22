import * as React from 'react';
import { connect } from 'react-redux';
import { Button, Table, message, Modal, Tag, Icon } from 'antd';
import './style.scss';
import { AppStoreTypes } from '@/stores';
import { imageStoreService } from '@/services';
import * as Cookie from 'js-cookie';
import EditModal from './editModal';
import utils from '@/utils/utils';

interface IProps {
  cur_parent_cluster: any;
  authorityList: any;
  location?: any;
}
interface IState {
  reqParams: {
    current: number;
    pageSize: number;
  };
  imageStores: {
    list: any[];
    count: number;
  };
  tableLoading: boolean;
  selectedRowKeys: number[];
  showModal: boolean;
  imageStoreInfo: any;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  cur_parent_cluster: state.HeaderStore.cur_parent_cluster,
  authorityList: state.UserCenterStore.authorityList,
});

@(connect(mapStateToProps, undefined) as any)
export default class ImageStore extends React.PureComponent<IProps, IState> {
  state: IState = {
    reqParams: {
      current: 1,
      pageSize: 10,
    },
    imageStores: {
      list: [],
      count: 0,
    },
    tableLoading: false,
    selectedRowKeys: [],
    showModal: false,
    imageStoreInfo: {},
  };

  componentDidMount() {
    this.getImageStoreList();
  }

  // 获取镜像仓库列表
  getImageStoreList = () => {
    const { cur_parent_cluster } = this.props;
    this.setState({ tableLoading: true });
    const cluster_id =
      cur_parent_cluster?.id > 0
        ? cur_parent_cluster?.id
        : +Cookie.get('em_current_cluster_id');
    console.log('cur_parent_cluster', cur_parent_cluster);
    imageStoreService.getImageStoreList({ cluster_id }).then((res: any) => {
      res = res.data;
      this.setState({ tableLoading: false });
      if (res.code === 0) {
        this.setState({
          imageStores: res.data,
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 表格分页，
  handleTableChange = (pagination, filters, sorter) => {
    const reqParams = {
      ...this.state.reqParams,
      current: pagination.current,
    };
    this.setState({ reqParams });
  };

  // 权限控制
  authorityControl = (action: string, code: string, record?: any) => {
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, code)) {
      return;
    }
    switch (action) {
      case 'add':
        this.handleShowModal();
        break;
      case 'delete':
        this.handleDelete();
        break;
      default:
        break;
    }
  };

  // 模态框是否显示
  handleShowModal = () => {
    const { showModal, imageStoreInfo } = this.state;
    this.setState({
      showModal: !showModal,
      imageStoreInfo: !showModal ? imageStoreInfo : {},
    });
  };

  // 编辑仓库
  editImageStore = (record: any, e: any) => {
    e.preventDefault();
    imageStoreService
      .getImageStoreInfo({ store_id: record.id })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.setState({ imageStoreInfo: res.data }, () => {
            this.handleShowModal();
          });
        } else {
          message.error(res.msg);
        }
      });
  };

  // 表格勾选项
  selectChange = (selectedRowKeys: number[]) => {
    this.setState({ selectedRowKeys });
  };

  // 删除镜像仓库
  handleDelete = () => {
    const { selectedRowKeys = [] } = this.state;
    Modal.confirm({
      title: '确定要删除该镜像仓库吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      // okType: "danger",
      onOk: async () => {
        const res = await imageStoreService.deleteImageStore({
          id: selectedRowKeys,
        });
        const { code, msg } = res.data;
        if (code === 0) {
          message.success('执行成功');
          this.getImageStoreList();
        } else {
          message.error(msg);
        }
      },
      onCancel: () => {},
    });
  };

  // 设置默认仓库
  handleSetDefaultStore = (record: any) => {
    const { cur_parent_cluster } = this.props;
    imageStoreService
      .setDefaultStore({
        clusterId: cur_parent_cluster.id,
        id: record.id,
      })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          message.success('更改默认仓库成功');
          this.getImageStoreList();
        } else {
          message.error('更改默认仓库失败');
        }
      });
  };

  // 表格列信息
  getColumns = () => {
    const { authorityList } = this.props;
    const columns = [
      {
        title: '仓库名称',
        dataIndex: 'name',
      },
      {
        title: '仓库别名',
        dataIndex: 'alias',
      },
      {
        title: '仓库地址',
        dataIndex: 'address',
        render: (text, record) => (
          <div>
            <p>{record.address}</p>
            <p className="text-gray">用户名：{record.username}</p>
          </div>
        ),
      },
      {
        title: '操作',
        dataIndex: 'func',
        render: (text, record) => (
          <React.Fragment>
            {authorityList.image_store_edit ? (
              <a onClick={this.editImageStore.bind(this, record)}>编辑</a>
            ) : (
              '--'
            )}
          </React.Fragment>
        ),
      },
      {
        title: '',
        dataIndex: 'is_default',
        render: (text, record) => (
          <div className="set-default-box">
            {record.is_default ? (
              <Tag color="rgba(63, 135, 255, 0.1)">默认仓库</Tag>
            ) : (
              <React.Fragment>
                {authorityList.image_store_edit ? (
                  <Button
                    type="default"
                    onClick={this.handleSetDefaultStore.bind(this, record)}>
                    设为默认仓库
                  </Button>
                ) : (
                  '--'
                )}
              </React.Fragment>
            )}
          </div>
        ),
      },
    ];
    return columns;
  };

  render() {
    const {
      reqParams,
      imageStores,
      tableLoading,
      selectedRowKeys = [],
      showModal,
      imageStoreInfo,
    } = this.state;
    const { cur_parent_cluster } = this.props;
    const columns = this.getColumns();
    const pagination = {
      pageSize: reqParams.pageSize,
      total: imageStores.count,
      current: reqParams.current,
      size: 'small',
      showTotal: (total) => (
        <span style={{ color: '#666666' }}>
          共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示
          {reqParams.pageSize}条
        </span>
      ),
    };
    const rowSelection = {
      selectedRowKeys,
      onChange: this.selectChange,
    };
    return (
      <div className="cluster-page-body">
        <div className="page-imagestore-header mb-12">
          <Button
            data-testid="add-btn"
            type="primary"
            onClick={() => this.authorityControl('add', 'image_store_edit')}>
            添加仓库
          </Button>
        </div>
        <Table
          rowKey="id"
          className="dt-table-fixed-contain-footer c-cluster__table"
          style={{ height: 'calc(100vh - 230px)' }}
          columns={columns}
          dataSource={imageStores.list}
          pagination={pagination}
          loading={tableLoading}
          scroll={{ y: true }}
          onChange={this.handleTableChange}
          rowSelection={rowSelection}
          footer={() =>
            !selectedRowKeys.length ? (
              <Button disabled>
                <i className="emicon emicon-delete" />
                删除
              </Button>
            ) : (
              <Button
                type="danger"
                onClick={() =>
                  this.authorityControl('delete', 'image_store_edit')
                }>
                <i className="emicon emicon-delete" />
                删除
              </Button>
            )
          }
        />
        {showModal && (
          <EditModal
            handleCancel={this.handleShowModal}
            getImageStoreList={this.getImageStoreList}
            imageStoreInfo={imageStoreInfo}
            clusterId={cur_parent_cluster.id}
            isEdit={!!Object.keys(imageStoreInfo).length}
          />
        )}
      </div>
    );
  }
}
