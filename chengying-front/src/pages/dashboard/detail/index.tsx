import * as React from 'react';
import { Layout } from 'antd';
const { Content } = Layout;

interface DetailProp {
  location: any;
  history?: any;
}

export default class DashDetailPage extends React.Component<DetailProp> {
  state = {
    dashurl: '/',
  };

  componentDidMount() {
    this.frame.current.style.height = document.body.offsetHeight - 52 + 'px';
    // postmessage监听iframe仪表盘回退事件
    window.addEventListener('message', (e) => {
      if (e.data === 'rollback' && location.href.indexOf('service') < 0) {
        // window.history.back()
        this.props.history.push('/deploycenter/monitoring/dashboard');
        // window.history.pushState({},'/dashboard/board');
      }
    });
  }

  componentWillUnmount() {
    // 注销postmessage监听
    window.removeEventListener('message', () => {}); // tslint:disable-line
  }

  private frame = React.createRef<HTMLIFrameElement>();
  render() {
    const { location } = this.props;
    console.log('仪表盘详情', location);
    const param: any = {};
    location.search
      .slice(1)
      .split('&')
      .map((o: string) => {
        param[o.split('=')[0]] = o.split('=')[1];
      });
    let listUrl = param.panelId
      ? window.location.protocol +
        '//' +
        window.location.hostname +
        ':' +
        window.APPCONFIG.GRAFANA_PORT +
        decodeURIComponent(param.url) +
        '?theme=light&tab=alert&panelId=' +
        param.panelId
      : window.location.protocol +
        '//' +
        window.location.hostname +
        ':' +
        window.APPCONFIG.GRAFANA_PORT +
        decodeURIComponent(param.url) +
        '?theme=light';
    if (param['var-cluster']) {
      listUrl += '&var-cluster=' + param['var-cluster'];
    }
    return (
      <Layout>
        <Layout>
          <Content>
            <iframe
              ref={this.frame}
              className="style-iframe"
              src={listUrl}
              frameBorder="0"></iframe>
          </Content>
        </Layout>
      </Layout>
    );
  }
}
