import apis from '@/constants/apis';
import * as http from '@/utils/http';

const { securityAudit } = apis;

export default {
  // 获取安全审计列表
  getSafetyAudit(params: any) {
    return http[securityAudit.getSafetyAudit.method](
      securityAudit.getSafetyAudit.url,
      params
    );
  },
  // 获取审计操作模块
  getAuditModule() {
    return http[securityAudit.getAuditModule.method](
      securityAudit.getAuditModule.url
    );
  },
  // 获取审计操作模块
  getAuditOperation(params: any) {
    return http[securityAudit.getAuditOperation.method](
      securityAudit.getAuditOperation.url,
      params
    );
  },
};
