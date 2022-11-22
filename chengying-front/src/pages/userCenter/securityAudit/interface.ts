interface AuditReqParams {
  from: number;
  to: number;
  module: string;
  operation: string;
  content: string;
  operator: string;
  ip: string;
}
export type AuditReqParamsType = AuditReqParams;
