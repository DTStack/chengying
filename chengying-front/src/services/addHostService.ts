import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { addHost, host } = apis;

export default {
  pwdInstallUrl(params: any) {
    return http[addHost.pwdInstallUrl.method](
      addHost.pwdInstallUrl.url,
      params
    );
  },
  connectHost(params: any, type: string) {
    return type === 'pk'
      ? this.pkConnectUrl(params)
      : this.pwdConnectUrl(params);
  },
  pkConnectUrl(params: any) {
    return http[addHost.pkConnectUrl.method](addHost.pkConnectUrl.url, params);
  },
  pwdConnectUrl(params: any) {
    return http[addHost.pwdConnectUrl.method](
      addHost.pwdConnectUrl.url,
      params
    );
  },
  installHost(params: any, type: string) {
    return type === 'pk'
      ? this.pkInstallUrl(params)
      : this.pwdInstallUrl(params);
  },
  pkInstallUrl(params: any) {
    return http[addHost.pkInstallUrl.method](addHost.pkInstallUrl.url, params);
  },
  checkInstallUrl(params: any) {
    return http[addHost.checkInstallUrl.method](
      addHost.checkInstallUrl.url,
      params
    );
  },
  confirmMoveHost(params: any) {
    return http[host.confirmMoveHost.method](host.confirmMoveHost.url, params);
  },
  updateGroupName(params: any) {
    return http[host.updateGroupName.method](host.updateGroupName.url, params);
  },
  deleteHost(params: any) {
    return http[host.deleteHost.method](host.deleteHost.url, params);
  },
};
