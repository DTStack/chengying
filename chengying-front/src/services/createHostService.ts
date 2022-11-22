import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { createHost } = apis;

export default {
  testConnection_pwd(params: any[]) {
    return http[createHost.testConnection_pwd.method](
      createHost.testConnection_pwd.url,
      params
    );
  },
  testConnection_pk(params: any[]) {
    return http[createHost.testConnection_pk.method](
      createHost.testConnection_pk.url,
      params
    );
  },
  installHostByPk(params: any) {
    return http[createHost.pkInstallUrl.method](
      createHost.pkInstallUrl.url,
      params
    );
  },
  installHostByPwd(params: any) {
    return http[createHost.pwdInstallUrl.method](
      createHost.pwdInstallUrl.url,
      params
    );
  },
  /* -- 合并 -- */
  // 测试连通性
  testConnection(params: any[], type: number) {
    return type === 1
      ? this.testConnection_pk(params)
      : this.testConnection_pwd(params);
  },
  // 安装主机
  installHost(params: any, type: number) {
    return type === 1
      ? this.installHostByPk(params)
      : this.installHostByPwd(params);
  },
};
