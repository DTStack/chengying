import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { event } = apis;

export default {
  getEventType(params: any) {
    return http[event.getEventType.method](event.getEventType.url, params);
  },
  getEventCount(params: any) {
    return http[event.getEventCount.method](event.getEventCount.url, params);
  },
  getEventEcharts(params: any) {
    return http[event.getEventEcharts.method](
      event.getEventEcharts.url,
      params
    );
  },
  getEventList(params: any) {
    return http[event.getEventList.method](event.getEventList.url, params);
  },
  getEventProductRank(params: any, config?: any) {
    return http[event.getEventProductRank(params).method](
      event.getEventProductRank(params).url,
      config
    );
  },
};
