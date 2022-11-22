import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { clusterHost } = apis;

export default {
  // 主机删除，沿用之前接口
  deleteHost(params: any) {
    return http[clusterHost.deleteHost.method](
      clusterHost.deleteHost.url,
      params
    );
  },
  // 获取集群下角色清单
  getRoleList(params: any) {
    return http[clusterHost.getRoleList.method](
      clusterHost.getRoleList.url,
      params
    );
  },
  // 删除集群下角色
  deleteRole(params: any) {
    return http[clusterHost.deleteRole.method](
      clusterHost.deleteRole.url,
      params
    );
  },
  // 添加角色
  addRole(params: any) {
    return http[clusterHost.addRole.method](clusterHost.addRole.url, params);
  },
  modifyRole(params: any) {
    return http[clusterHost.modifyRole.method](
      clusterHost.modifyRole.url,
      params
    );
  },
  // 编辑主机角色
  bindHostRoles(params: any) {
    return http[clusterHost.bindHostRoles.method](
      clusterHost.bindHostRoles.url,
      params
    );
  },
};
