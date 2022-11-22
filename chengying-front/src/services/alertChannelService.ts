import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { alertChannel } = apis;

export default {
  getAlertNotifications() {
    return http[alertChannel.getAlertNotifications.method](
      alertChannel.getAlertNotifications.url
    );
  },
  delAlertNotification(notificationId: any) {
    return http[alertChannel.delAlertNotification.method](
      alertChannel.delAlertNotification.url + notificationId
    );
  },
  grafanaAlertChannelTest(params: any) {
    return http[alertChannel.grafanaAlertChannelTest.method](
      alertChannel.grafanaAlertChannelTest.url,
      params
    );
  },
  grafanaAlertChannelSave(params: any) {
    return http[alertChannel.grafanaAlertChannelSave.method](
      alertChannel.grafanaAlertChannelSave.url,
      params
    );
  },
  dtstackAlertChannelSave(params: any) {
    return http[alertChannel.dtstackAlertChannelSave.method](
      alertChannel.dtstackAlertChannelSave.url,
      params
    );
  },
  dtstackAlertChannelList(params: any) {
    return http[alertChannel.dtstackAlertChannelList.method](
      alertChannel.dtstackAlertChannelList.url,
      params
    );
  },
  getDtstackAlertDetail(params: any) {
    return http[alertChannel.getDtstackAlertDetail.method](
      alertChannel.getDtstackAlertDetail.url,
      params
    );
  },
  getGrafanaAlertDetail(params: any) {
    return http[alertChannel.getGrafanaAlertDetail(params).method](
      alertChannel.getGrafanaAlertDetail(params).url
    );
  },
  dtstackAlertChannelDel(params: any) {
    return http[alertChannel.dtstackAlertChannelDel.method](
      alertChannel.dtstackAlertChannelDel.url,
      params
    );
  },

  grafanaAlertChannelUpdate(params: any) {
    return http[alertChannel.grafanaAlertChannelUpdate(params).method](
      alertChannel.grafanaAlertChannelUpdate(params).url,
      params
    );
  },
};
