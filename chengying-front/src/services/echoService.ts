import apis from '@/constants/apis';
import * as http from '@/utils/http';

const { scriptEcho } = apis;

export default {
  getEchoCount(params: any) {
    return http[scriptEcho.getEchoCount(params).method](
      scriptEcho.getEchoCount(params).url,
      params
    );
  },
  getEchoOrder(params: any) {
    return http[scriptEcho.getEchoOrder(params).method](
      scriptEcho.getEchoOrder(params).url,
      params
    );
  },
  getEchoSearchList(params: any) {
    return http[scriptEcho.getEchoSearchList(params).method](
      scriptEcho.getEchoSearchList(params).url,
      params
    );
  },
  getEchoOrderDetail(params: any) {
    return http[scriptEcho.getEchoOrderDetail(params).method](
      scriptEcho.getEchoOrderDetail(params).url,
      params
    );
  },
  // 显示shell日志回显
  showShellLog(params: any) {
    return http[scriptEcho.showShellLog.method](
      scriptEcho.showShellLog.url,
      params
    );
  },
  showShellContent(params: any) {
    return http[scriptEcho.showShellContent.method](
      scriptEcho.showShellContent.url,
      params
    );
  },
};
