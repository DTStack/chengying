import * as React from 'react';
import { Button } from 'antd';
import './style.scss';
import utils from '@/utils/utils';

declare var APP: any;

interface IState {
  clientWidth: number;
  clientHeight: number;
  currentWidth: number;
  currentHeight: number;
}

class InstallGuide extends React.Component<any, IState> {
  state: IState = {
    clientWidth: 0,
    clientHeight: 0,
    currentWidth: 0,
    currentHeight: 0,
  };

  componentDidMount() {
    this.resize();
    window.addEventListener('resize', this.resize);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resize);
  }

  resize = () => {
    const clientWidth = document.body.clientWidth;
    const clientHeight = document.body.clientHeight;
    const scale = 1440 / 864;

    this.setState({
      clientHeight,
      clientWidth,
      currentHeight: clientHeight - 48,
      currentWidth: scale * clientHeight,
    });
  };

  handleStartDeploy = () => {
    let path = '/deploycenter/appmanage/installs';
    utils.setNaviKey('menu_deploy_center', 'sub_menu_product_deploy');
    this.props.history.push(path);
  };

  render() {
    const { clientWidth, currentWidth, currentHeight } = this.state;
    const right = Math.round((clientWidth - currentWidth) / 2);
    const titleFontSize = Math.round(36 / (1200 / currentWidth));
    const descWidth = Math.round(506 / (1200 / currentWidth));
    return (
      <div
        className="install-guide-container"
        style={{ height: currentHeight }}>
        <div className="install-guide-bg bg-st">
          <img
            src={require('public/imgs/guide_bg1.png')}
            height={currentHeight}
            width={clientWidth - right}
          />
        </div>
        <div className="install-guide-bg bg-nd" style={{ right }}>
          <img
            src={require('public/imgs/guide_bg2.png')}
            height={currentHeight}
          />
        </div>
        <div
          className="install-guide-content"
          style={{
            width: currentWidth,
            height: currentHeight,
          }}>
          <div
            style={{
              marginTop: `-${Math.round(120 / (1200 / currentWidth))}px`,
            }}>
            <p
              className="title mb-20"
              style={{
                fontSize:
                  titleFontSize < 24
                    ? 24
                    : titleFontSize > 36
                    ? 36
                    : titleFontSize,
              }}>
              欢迎使用 ChengYing 产品部署工具
            </p>
            <p className="desc" style={{ width: descWidth }}>
              ChengYing
              是一款全自动化全生命周期的运维管家，支持Kubernetes容器与物理/虚拟机两种方式进行多集群部署，提供产品部署、扩缩节点、日志诊断、集群运维、监控告警等一站式服务，轻松实现大数据产品全流程部署与监控。
            </p>
            <Button
              type="primary"
              shape="round"
              onClick={this.handleStartDeploy}>
              开始向导
            </Button>
          </div>
        </div>
        {/* <Route path={`${this.props.match.path}/install`} component={Steps}/> */}
        <div
          className="version"
          style={{ position: 'fixed', bottom: 10, color: '#fff' }}>
          ChengYing@V{APP.VERSION}
        </div>
      </div>
    );
  }
}

export default InstallGuide;
