import * as React from 'react';
import { Row, Col, message } from 'antd';
import { connect } from 'react-redux';
import { HeaderStateTypes } from '@/stores/headerReducer';
import { AppStoreTypes } from '@/stores';
import ResourceUsage from './components/resourceUsage';
import NodeView from './components/nodeView';
import './style.scss';
import ComponentStatus from './components/componentStatus';
import InfoView from './components/infoView';
import WorkLoad from './components/workLoad';
import TopFiveUsage from './components/topFiveUsage';
import clusterIndexService from '@/services/clusterIndexService';

interface IProps {
  HeaderStore: HeaderStateTypes;
  location: any;
}
interface IState {
  overview: any;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
});

@(connect(mapStateToProps, undefined) as any)
export default class IndexPage extends React.PureComponent<IProps, IState> {
  state: IState = {
    overview: {},
  };

  componentDidMount() {
    this.getClusterOverview();
  }

  componentDidUpdate(prevProps: IProps) {
    if (
      prevProps.HeaderStore.cur_parent_cluster.id !==
      this.props.HeaderStore.cur_parent_cluster.id
    ) {
      this.getClusterOverview();
    }
  }

  // 获取总览
  getClusterOverview = () => {
    const { cur_parent_cluster } = this.props.HeaderStore;
    clusterIndexService
      .getClusterOverview(
        {
          cluster_id: cur_parent_cluster.id,
        },
        cur_parent_cluster.type
      )
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.setState({ overview: res.data });
        } else {
          message.error(res.msg);
        }
      });
  };

  render() {
    const { overview = {} } = this.state;
    const { cur_parent_cluster } = this.props.HeaderStore;
    const isKubernetes = cur_parent_cluster.type === 'kubernetes';
    return (
      <div className="cluster-page-body cluster-overflow">
        <ResourceUsage
          overview={overview}
          cur_parent_cluster={cur_parent_cluster}
        />
        {isKubernetes ? (
          <Row className="mt-20" gutter={20}>
            <Col span={8}>
              <NodeView overview={overview} />
            </Col>
            <Col span={8}>
              <ComponentStatus overview={overview} />
            </Col>
            <Col span={8}>
              <InfoView
                overview={overview}
                cur_parent_cluster={cur_parent_cluster}
              />
            </Col>
          </Row>
        ) : (
          <Row className="mt-20" gutter={20}>
            <Col span={12}>
              <NodeView overview={overview} />
            </Col>
            <Col span={12}>
              <InfoView
                overview={overview}
                cur_parent_cluster={cur_parent_cluster}
              />
            </Col>
          </Row>
        )}
        {isKubernetes && (
          <div className="mt-20">
            <WorkLoad overview={overview} />
          </div>
        )}
        <div className="mt-20">
          <TopFiveUsage overview={overview} />
        </div>
      </div>
    );
  }
}
