import * as React from 'react';
import { Row, Col } from 'antd';
import ChoiceCard from '@/components/choiceCard';
import './style.scss';

interface IProps {
  history: any;
}

export default class CreateCluster extends React.PureComponent<IProps, any> {
  // 创建跳转
  handleTypeClick = (type, mode) => {
    console.log(this.props);
    this.props.history.push(
      `/deploycenter/cluster/create/add?id=-1&type=${type}&mode=${mode}`
    );
  };

  render() {
    return (
      <div className="cluster-container">
        <div className="cluster-type">
          <p className="text-title-bold">主机集群</p>
          <Row gutter={20}>
            <Col span={8}>
              <ChoiceCard
                className="cluster-box-style"
                title="创建主机集群"
                content="创建基于物理机/虚拟机的集群，直接将应用部署物理机/虚拟机上"
                imgSrc={require('public/imgs/cluster_hosts.png')}
                handleTypeClick={this.handleTypeClick.bind(this, 'hosts', 0)}
              />
            </Col>
          </Row>
        </div>
        <div className="cluster-type" style={{ display: 'none' }}>
          <p className="text-title-bold">Kubernetes集群</p>
          <Row gutter={20}>
            <Col span={8}>
              <ChoiceCard
                className="cluster-box-style"
                title="导入已有Kubernetes集群"
                content="导入现有的Kubernetes集群，Kubernetes集群的管理和配置继续由提供商提供"
                imgSrc={require('public/imgs/kubernetes_import.png')}
                handleTypeClick={this.handleTypeClick.bind(
                  this,
                  'kubernetes',
                  1
                )}
              />
            </Col>
            <Col span={8}>
              <ChoiceCard
                className="cluster-box-style"
                title="自建Kubernetes集群"
                content="从现有的物理机或虚拟机中创建一个新的Kubernetes集群"
                imgSrc={require('public/imgs/kubernetes_self.png')}
                handleTypeClick={this.handleTypeClick.bind(
                  this,
                  'kubernetes',
                  0
                )}
              />
            </Col>
          </Row>
        </div>
      </div>
    );
  }
}
