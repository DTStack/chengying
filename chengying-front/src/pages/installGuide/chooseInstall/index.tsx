import * as React from 'react';
import { Alert, Row, Col, message, Spin } from 'antd';
import { connect } from 'react-redux';
import { Dispatch, bindActionCreators } from 'redux';
import * as installGuideAction from '@/actions/installGuideAction';
import { installGuideService } from '@/services';
import ChoiceCard from '@/components/choiceCard';
import './style.scss';

const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, installGuideAction), dispatch),
});

interface IProps {
  actions: installGuideAction.InstallGuideActionTypes;
  history: any;
  location: any;
}

interface IState {
  clusterType: string;
  showAlert: boolean;
  loading: boolean;
  lists: {
    hostsList: any[];
    kubernetesList: any[];
  };
}

@(connect(undefined, mapDispatchToProps) as any)
export default class ChooseInstall extends React.PureComponent<IProps, IState> {
  state: IState = {
    clusterType: 'hosts',
    showAlert: false,
    loading: true,
    lists: {
      hostsList: [],
      kubernetesList: [],
    },
  };

  componentDidMount() {
    this.checkCluster();
  }

  // 是否有可用集群
  checkCluster = async () => {
    const response = await installGuideService.getInstallClusterList({
      'sort-by': 'id',
      'sort-dir': 'desc',
      limit: 0,
      start: 0,
    });
    const res = response.data;
    const hostsList = [];
    const kubernetesList = [];
    if (res.code === 0) {
      res.data.counts > 0 &&
        res.data.clusters.forEach((item) => {
          if (item.status === 'Running' || item.status === 'Error') {
            if (item.type === 'hosts') {
              hostsList.push(item);
            } else {
              kubernetesList.push(item);
            }
          }
        });
      this.setState({
        lists: {
          hostsList,
          kubernetesList,
        },
      });
    } else {
      message.error(res.msg);
    }
    this.setState({ loading: false });
  };

  // 点击
  handleTypeClick = (type: string) => {
    const list = this.state.lists[`${type}List`];
    if (list.length) {
      let p = `/deploycenter/appmanage/installs?type=${type}`;
      if (this.props.location.search) {
        p += '&' + this.props.location.search.split('?')[1];
      }
      this.props.history.push(p);
    }
    this.setState({
      showAlert: true,
      clusterType: type,
    });
  };

  render() {
    const { clusterType, showAlert, loading } = this.state;
    return (
      <div className="choose-install-page">
        <p className="text-title-bold mb-20">部署方式</p>
        {showAlert && (
          <Alert
            className="mb-20"
            type="warning"
            showIcon
            message={
              <span>
                未检测到可用{clusterType === 'hosts' ? '主机' : 'Kubernetes'}
                集群，请前往”<a href="/clustermanage/create">集群管理</a>
                “进行集群创建或检查
              </span>
            }
          />
        )}
        {!loading ? (
          <Row gutter={20}>
            {/* <Col span={8}>
              <ChoiceCard
                className="cluster-box-style"
                title="Kubernetes部署模式"
                content="应用部署在K8S集群上，适用于容器部署模式"
                imgSrc={require('public/imgs/install_kubernetes.png')}
                handleTypeClick={this.handleTypeClick.bind(this, 'kubernetes')}
              />
            </Col> */}
            <Col span={8}>
              <ChoiceCard
                className="cluster-box-style"
                title="物理/虚拟机部署模式"
                content="应用部署在物理/虚拟机上，适用于传统部署模式"
                imgSrc={require('public/imgs/install_host.png')}
                handleTypeClick={this.handleTypeClick.bind(this, 'hosts')}
              />
            </Col>
          </Row>
        ) : (
          <Spin spinning={loading}>
            <div style={{ height: '100vh' }}></div>
          </Spin>
        )}
      </div>
    );
  }
}
