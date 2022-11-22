import * as React from 'react';

interface iHostDashProps {
  searchStr: string;
  cur_parent_cluster: any;
}

/**
 * 主机仪表盘模块
 * @param {*} props
 * @return {*}
 */
const HostDash: React.FC<iHostDashProps> = (props) => {
  const { cur_parent_cluster, searchStr } = props;

  const getRequestParams = () => {
    const url: string = searchStr; // 获取url中"?"符后的字串
    const theRequest = {};
    console.log(theRequest);
    if (url.indexOf('?') !== -1) {
      const str = url.substr(1);
      const strs = str.split('&');
      for (let i = 0; i < strs.length; i++) {
        // tslint:disable-line
        theRequest[strs[i].split('=')[0]] = unescape(strs[i].split('=')[1]);
      }
    }
    return theRequest;
  };

  const host: string = getRequestParams()?.host;
  const dashUrl: string =
    window.location.protocol +
    '//' +
    window.location.hostname +
    ':' +
    window.APPCONFIG.GRAFANA_PORT +
    '/grafana/d/Ne_roaViz/host-overview?orgId=1&var-node=' +
    host +
    '&var-cluster=' +
    cur_parent_cluster.name +
    '&var-job=node_exporter&var-port=9100&theme=light&no-feedback=true';

  return (
    <iframe
      className="style-iframe"
      style={{ height: document.body.clientHeight - 88 - 10 }}
      src={dashUrl}
      frameBorder="0"></iframe>
  );
};
export default HostDash;
