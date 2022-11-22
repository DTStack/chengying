import * as React from 'react';
// import { bindActionCreators, Dispatch } from 'redux';
// import { connect } from 'react-redux';
// import MainLayout from '@/layouts/mainLayout';
import { Layout } from 'antd';

interface HostDetailProps {
  match: any;
}
export default class HostDetailPage extends React.Component<HostDetailProps> {
  constructor(props: object) {
    super(props as HostDetailProps);
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
      this.props.match.params.host +
      '&var-job=node_exporter&var-port=9100&theme=light';
    console.log(dashUrl);
    return (
      <div>
        <Layout className="c-sub-layout">
          <iframe
            className="style-iframe"
            src={dashUrl}
            frameBorder="0"></iframe>
        </Layout>
      </div>
    );
  }
}
