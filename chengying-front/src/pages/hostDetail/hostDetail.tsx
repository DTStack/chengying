import * as React from 'react';
// import { bindActionCreators, Dispatch } from 'redux';
// import { connect } from 'react-redux';
// import MainLayout from '@/layouts/mainLayout';
import { Layout } from 'antd';

interface HostDetailProps {
  location: any;
}
interface HostDetailState {
  host: string;
}
export default class HostDetailPage extends React.Component<
  HostDetailProps,
  HostDetailState
> {
  constructor(props: object) {
    super(props as HostDetailProps);
  }

  public state: HostDetailState = {
    host: this.getRequestParams().host,
  };

  // tslint:disable
  componentDidMount() {
    // postmessage监听iframe仪表盘回退事件
    window.addEventListener('message', function (e) {
      if (e.data === 'rollback' && location.href.indexOf('service') < 0) {
        window.history.back();
      }
    });
  }

  componentWillUnmount() {
    // 注销postmessage监听
    window.removeEventListener('message', () => {});
  }

  getRequestParams(): any {
    const url = location.search; // 获取url中"?"符后的字串
    const theRequest = {};
    if (url.indexOf('?') !== -1) {
      const str = url.substr(1);
      const strs = str.split('&');
      for (let i = 0; i < strs.length; i++) {
        theRequest[strs[i].split('=')[0]] = unescape(strs[i].split('=')[1]);
      }
    }
    return theRequest;
  }

  render() {
    // let dashUrl = window.location.protocol + '//172.16.10.85:' + window.APPCONFIG.GRAFANA_PORT + '/grafana/d/Ne_roaViz/host-overview?orgId=1&var-node=' + this.props.match.params.host + '&var-job=node_exporter&var-port=9100&theme=light&inactive=true';
    const dashUrl =
      window.location.protocol +
      '//' +
      window.location.hostname +
      ':' +
      window.APPCONFIG.GRAFANA_PORT +
      '/grafana/d/Ne_roaViz/host-overview?orgId=1&var-node=' +
      this.state.host +
      '&var-job=node_exporter&var-port=9100&theme=light';
    console.log(dashUrl);
    return (
      <div>
        <Layout className="c-sub-layout">
          <iframe
            style={{ height: document.body.clientHeight - 88 }}
            className="style-iframe"
            src={dashUrl}
            frameBorder="0"></iframe>
        </Layout>
      </div>
    );
  }
}
