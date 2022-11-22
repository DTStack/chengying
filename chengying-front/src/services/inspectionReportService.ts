import * as http from '@/utils/http';
export default {
  // 获取节点状态
  getNodeStatus(params?: any) {
    return http.get('/api/v2/inspect/host/status', params);
  },
  // 获取应用状态
  getAppStatus(params?: any) {
    return http.get('/api/v2/inspect/service/status', params);
  },
  // 获取告警历史
  getHistoryAlarm(params?: any) {
    return http.get('/api/v2/inspect/alert/history', params);
  },
  // 获取图表配置列表
  getChartConfigList(params?: any) {
    return http.get('/api/v2/inspect/graph/config', params);
  },
  // 获取图表数据
  getChartData(params?: any) {
    return http.get('/api/v2/inspect/graph/data', params);
  },
  // 生成巡检报告
  generatorReport(params?: any) {
    return http.post('/api/v2/inspect/generate', params);
  },
  // 查询生成进度
  getReportProgress(params?: any) {
    return http.get('/api/v2/inspect/progress', params);
  },
  downloadFile(params?: any) {
    return http.get('/api/v2/inspect/download', params);
  },
};
