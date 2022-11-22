import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { alertRule } = apis;

export default {
  getAlertsByDashId(params: any) {
    return http[alertRule.getAlertsByDashId.method](
      alertRule.getAlertsByDashId.url,
      params
    );
  },
  getAlertRuleList(params: any) {
    return http[alertRule.getAlertRuleList.method](
      alertRule.getAlertRuleList.url,
      params
    );
  },
  switchGrafanaAlertPause(params: any) {
    return http[alertRule.switchGrafanaAlertPause.method](
      alertRule.switchGrafanaAlertPause.url,
      params
    );
  },
};
