import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { config } = apis;

export default {
  getConfigAlertGroups(params: any) {
    return http[config.getConfigAlertGroups.method](
      config.getConfigAlertGroups.url,
      params
    );
  },
  getConfigAlertaction(params: any) {
    return http[config.getConfigAlertaction.method](
      config.getConfigAlertaction.url,
      params
    );
  },
};
