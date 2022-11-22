import * as React from 'react';
import { bindActionCreators, Dispatch } from 'redux';
import { connect } from 'react-redux';
import moment from 'moment';
import * as Cookies from 'js-cookie';
import { chunk } from 'lodash';
import {
  Card,
  Tag,
  Dropdown,
  Menu,
  Icon,
  Button,
  Row,
  Col,
  Modal,
  message,
  Divider,
  Tooltip,
  Descriptions,
  Empty,
  Spin,
} from 'antd';
import * as HeaderAction from '@/actions/headerAction';
import { HeaderStateTypes } from '@/stores/headerReducer';
import { AppStoreTypes } from '@/stores';
import { clusterTypeMap } from '@/constants/const';
import { clusterManagerService } from '@/services';
import ProgressBar from '@/components/progressBar';
import utils from '@/utils/utils';
import './style.scss';

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
  authorityList: state.UserCenterStore.authorityList,
});

const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(HeaderAction, dispatch),
});

interface IProps {
  history: any;
  HeaderStore: HeaderStateTypes;
  actions: HeaderAction.HeaderActionTypes;
  authorityList: any;
}

interface IState {
  modalShow: boolean;
  clusterList: any[];
  info: any;
  loading: boolean;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
export default class ClusterList extends React.PureComponent<IProps, IState> {
  state: IState = {
    modalShow: false,
    info: {},
    clusterList: [],
    loading: false,
  };

  componentDidMount() {
    this.getClusterLists();
  }

  // 获取集群列表
  getClusterLists = () => {
    const params = {
      type: '',
      'sort-by': 'id',
      'sort-dir': 'desc',
      limit: 0,
      start: 0,
    };
    this.setState({ loading: true });
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
      this.setState({ loading: false });
    });
  };

  handleMenuClick = (item: any, e: any) => {
    const { authorityList } = this.props;
    e.domEvent.stopPropagation();
    // 如果是查看详情
    if (e.key === 'view') {
      if (utils.noAuthorityToDO(authorityList, 'cluster_view')) {
        return;
      }
      return this.handleModalShow(item);
    }
    // 其他权限统一处理
    if (utils.noAuthorityToDO(authorityList, 'cluster_edit')) {
      return;
    }
    switch (e.key) {
      case 'edit':
        this.props.history.push(
          `/deploycenter/cluster/create/edit?id=${item.id}&type=${item.type}&mode=${item.mode}`
        );
        break;
      case 'delete':
        this.handleDeleteCluster(item);
    }
  };

  // 查看详情
  handleModalShow = (info) => {
    const { modalShow } = this.state;

    this.setState({
      modalShow: !modalShow,
      info: info || {},
    });
  };

  // 删除集群
  handleDeleteCluster = (item: any) => {
    Modal.confirm({
      title: '确定要删除该集群吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      // okType: "danger",
      onOk: async () => {
        const response = await clusterManagerService.deleteCluster(
          { cluster_id: item.id },
          item.type
        );
        const res = response.data;
        if (res.code === 0) {
          message.success('集群删除成功');
          this.getClusterLists();
        } else {
          message.error(res.msg);
        }
      },
      onCancel: () => {},
    });
  };

  // 切换到集群视角下查看集群的信息
  handleClusterView = (cluster: any) => {
    const { actions, history } = this.props;
    actions.setCurrentParentCluster(cluster);
    //
    sessionStorage.setItem(
      'em_current_cluster_id',
      Cookies.get('em_current_cluster_id')
    );
    sessionStorage.setItem(
      'em_current_cluster_type',
      Cookies.get('em_current_cluster_type')
    );
    Cookies.set('em_current_cluster_id', cluster.id);
    Cookies.set('em_current_cluster_type', cluster.type);
    Cookies.set('em_current_k8s_namespace', '');
    Cookies.set('em_current_parent_product', '');
    // 跳转
    const url =
      cluster.mode === 0
        ? '/deploycenter/cluster/detail/index'
        : '/deploycenter/cluster/detail/imagestore';
    history.push(url);
  };

  // 卡片 - 更多
  getCardExtra = (item) => {
    const statusMap = {
      Running: 'status-running',
      Waiting: 'status-waiting',
      Error: 'status-error',
    };
    return (
      <React.Fragment>
        <div className={`card-extra-status ${statusMap[item.status]}`}>
          {item.status}
        </div>
        <Dropdown
          className="c-extra-more__dropdown ml-20"
          placement="bottomCenter"
          overlayClassName="c-extra-more__menuitem"
          overlay={
            <Menu onClick={this.handleMenuClick.bind(this, item)}>
              <Menu.Item key="view">
                <Icon type="profile" className="mr-10" />
                详情
              </Menu.Item>
              <Menu.Divider style={{ margin: 0 }} />
              <Menu.Item key="edit">
                <Icon type="edit" className="mr-10" />
                编辑
              </Menu.Item>
              {item.status === 'Waiting' && (
                <Menu.Item key="delete">
                  <Icon type="delete" className="mr-10" />
                  删除
                </Menu.Item>
              )}
            </Menu>
          }>
          <span onClick={(e) => e.stopPropagation()}> ··· </span>
        </Dropdown>
      </React.Fragment>
    );
  };

  // 添加集群
  handleAddCluster = () => {
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'cluster_edit')) {
      return;
    }
    utils.setNaviKey('menu_deploy_center', 'sub_menu_cluster_edit');
    this.props.history.push('/deploycenter/cluster/create');
  };

  render() {
    const { clusterList = [], info = {}, modalShow, loading } = this.state;
    return (
      <div className="cluster-manager-page">
        {clusterList.length > 0 ? (
          <div className="cluster-container cluster-list-page">
            {chunk(clusterList, 3).map((list, index) => (
              <Row key={index} gutter={20}>
                {list.map((item) => (
                  <Col className="mb-20" key={item.id} span={8}>
                    <Card
                      size="small"
                      className="cluster-page-card cluster-box-style"
                      onClick={this.handleClusterView.bind(this, item)}
                      title={item.name}
                      extra={this.getCardExtra(item)}>
                      <p>集群模式：{clusterTypeMap[item.type][item.mode]}</p>
                      <p>节点数：{item.nodes}</p>
                      <div className="mt-10">
                        {
                          item.tags ? (
                            item.tags.split(',').map((tag) => (
                              <Tag className="c-cluster__tag" key={tag}>
                                {tag}
                              </Tag>
                            ))
                          ) : (
                            <Tag style={{ border: 0 }}></Tag>
                          ) // 空tag占位
                        }
                      </div>
                      <Divider dashed style={{ margin: '20px 0' }} />
                      <Row className="mb-20">
                        <Col span={12} className="data-view-item">
                          <img
                            src={require('public/imgs/cpu.png')}
                            width={30}
                            height={30}
                          />
                          <span>CPU</span>
                          <ProgressBar
                            now={item.cpu_core_used_display}
                            total={item.cpu_core_size_display}
                          />
                        </Col>
                        <Col span={12} className="data-view-item">
                          <img
                            src={require('public/imgs/memory.png')}
                            width={30}
                            height={30}
                          />
                          <span>内存</span>
                          <ProgressBar
                            now={item.mem_used_display}
                            total={item.mem_size_display}
                          />
                        </Col>
                      </Row>
                      <Row>
                        {item.type === 'kubernetes' ? (
                          <Col span={12} className="data-view-item">
                            <img
                              src={require('public/imgs/pods.png')}
                              width={30}
                              height={30}
                            />
                            <span>PODS</span>
                            <ProgressBar
                              now={item.pod_used_display}
                              total={item.pod_size_display}
                            />
                          </Col>
                        ) : (
                          <Col span={12} className="data-view-item">
                            <img
                              src={require('public/imgs/disk.png')}
                              width={30}
                              height={30}
                            />
                            <span>磁盘</span>
                            <ProgressBar
                              now={item.disk_used_display}
                              total={item.disk_size_display}
                            />
                          </Col>
                        )}
                      </Row>
                    </Card>
                  </Col>
                ))}
              </Row>
            ))}
            <div className="cluster-bottom">
              <Tooltip title="添加集群">
                <span className="btn-box-shadow">
                  <i
                    className="emicon emicon-add"
                    onClick={this.handleAddCluster}
                  />
                </span>
              </Tooltip>
            </div>
          </div>
        ) : (
          <Spin spinning={loading}>
            <div
              style={{
                position: 'relative',
                height: 'calc(100vh - 40px)',
                width: '100%',
              }}>
              <div className="cluster-empty">
                <Empty
                  image={require('public/imgs/image_empty.png')}
                  imageStyle={{
                    height: 224,
                    width: 380,
                  }}
                  description={
                    <span
                      style={{
                        display: 'inline-block',
                        fontSize: '12px',
                        color: '#666',
                        marginTop: '-36px',
                      }}>
                      暂无数据
                    </span>
                  }>
                  <Button type="primary" onClick={this.handleAddCluster}>
                    添加集群
                  </Button>
                </Empty>
              </div>
            </div>
          </Spin>
        )}
        {modalShow && (
          <InfoModal info={info} handleCancel={this.handleModalShow} />
        )}
      </div>
    );
  }
}

const InfoModal = (props) => {
  const { info = {}, handleCancel } = props;
  return (
    <Modal
      title="集群详情"
      visible={true}
      footer={
        <Button type="primary" onClick={() => handleCancel(false)}>
          关闭
        </Button>
      }
      onCancel={() => handleCancel(false)}>
      <Descriptions className="c-cluster__description" column={1} bordered>
        <Descriptions.Item label="集群名称">{info.name}</Descriptions.Item>
        <Descriptions.Item label="集群模式">
          {clusterTypeMap[info.type][info.mode]}
        </Descriptions.Item>
        <Descriptions.Item label="集群描述">
          {info.desc || '--'}
        </Descriptions.Item>
        <Descriptions.Item label="集群标签">
          {info.tags
            ? info.tags.split(',').map((tag) => (
                <Tag
                  className="c-text-ellipsis c-cluster__tag"
                  style={{ maxWidth: 325 }}
                  key={tag}>
                  {tag}
                </Tag>
              ))
            : '--'}
        </Descriptions.Item>
        {info.type === 'kubernetes' && (
          <Descriptions.Item label="kubernetes版本">
            {info.version || '--'}
          </Descriptions.Item>
        )}
        <Descriptions.Item label="节点数">{info.nodes}</Descriptions.Item>
        <Descriptions.Item label="创建人">{info.create_user}</Descriptions.Item>
        <Descriptions.Item label="创建时间">
          {moment(info.create_time).format('YYYY-MM-DD HH:mm:ss')}
        </Descriptions.Item>
        <Descriptions.Item label="最近修改人">
          {info.update_user}
        </Descriptions.Item>
        <Descriptions.Item label="最近修改时间">
          {moment(info.update_time).format('YYYY-MM-DD HH:mm:ss')}
        </Descriptions.Item>
      </Descriptions>
    </Modal>
  );
};
