import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { clusterNamespace } = apis;

export default {
  // 获取命名空间列表
  getNamespaceList(params: any) {
    const request = { ...params };
    const param = { namespace: request.namespace };
    delete request.namespace;
    return http[clusterNamespace.getNamespaceList(param).method](
      clusterNamespace.getNamespaceList(param).url,
      request
    );
  },
  // 测试连通性
  pingConnect(params: any) {
    return http[clusterNamespace.pingConnect.method](
      clusterNamespace.pingConnect.url,
      params
    );
  },
  // 获取yaml文件
  getYamlFile(params: any) {
    return http[clusterNamespace.getYamlFile.method](
      clusterNamespace.getYamlFile.url,
      params
    );
  },
  // 保存namespace
  saveNamespace(params: any) {
    return http[clusterNamespace.saveNamespace.method](
      clusterNamespace.saveNamespace.url,
      params
    );
  },
  // 删除前进行校验，是否可以删除
  confirmDelete(params: any) {
    return http[clusterNamespace.confirmDelete(params).method](
      clusterNamespace.confirmDelete(params).url
    );
  },
  // 删除命名空间
  deleteNamespace(params: any) {
    return http[clusterNamespace.deleteNamespace(params).method](
      clusterNamespace.deleteNamespace(params).url
    );
  },
  // 获取命名空间信息
  getNamespaceInfo(params: any) {
    return http[clusterNamespace.getNamespaceInfo(params).method](
      clusterNamespace.getNamespaceInfo(params).url
    );
  },
  // 获取服务列表
  getServiceLists(params: any) {
    return http[clusterNamespace.getServiceLists(params).method](
      clusterNamespace.getServiceLists(params).url
    );
  },
  // 获取事件查看
  getEventLists(request: any) {
    const params = { namespace: request.namespace };
    return http[clusterNamespace.getEventLists(params).method](
      clusterNamespace.getEventLists(params).url,
      request
    );
  },
};
