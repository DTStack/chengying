import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { dashboard, service } = apis;
export default {
  /** DashBoard service */

  createDashboard(params: any) {
    return http[dashboard.createDashboard.method](
      dashboard.createDashboard.url,
      params
    );
  },
  createDashFolder(params: any) {
    return http[dashboard.createDashFolder.method](
      dashboard.createDashFolder.url,
      params
    );
  },
  delFolderByUid(uid: any) {
    return http[dashboard.delFolderByUid.method](
      dashboard.delFolderByUid.url + uid
    );
  },

  delDashByUid(uid: any) {
    return http[dashboard.delDashByUid.method](
      dashboard.delDashByUid.url + uid
    );
  },

  getServiceDashInfo(params: any) {
    const url = service.getServiceDashInfo.url;
    const qArr = [];
    for (const i in params) {
      if (i === 'tags') {
        for (const t of params[i]) {
          qArr.push('tag=' + t);
        }
      } else {
        qArr.push(i + '=' + params[i]);
      }
    }
    const query = qArr.join('&');
    return http[service.getServiceDashInfo.method](`${url}?${query}`);
  },
  importDashboard(params: any) {
    return http[dashboard.importDashboard.method](
      dashboard.importDashboard.url,
      params
    );
  },

  exportDashByIds(params: any) {
    return http[dashboard.exportDash.method](
      dashboard.exportDash.url,
      params
    );
  },
};
