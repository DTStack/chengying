import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { scriptManager } = apis;

export default {
    // 获取脚本列表
    getList(params: any) {
        return http[scriptManager.getList.method](
            scriptManager.getList.url,
            params
        );
    },
    getTaskContent(params: any) {
        return http[scriptManager.getTaskContent(params).method](
          scriptManager.getTaskContent(params).url
        );
    },
    getTaskStatus(params: any) {
        return http[scriptManager.getTaskStatus.method](
            scriptManager.getTaskStatus.url,
            params
        );
    },
    getSettingStatus(params: any) {
        const { host_id, spec } = params;
        return http[scriptManager.getSettingStatus(params).method](
          scriptManager.getSettingStatus(params).url,
          {
            host_id,
            spec
          }
        );
    },
    checkCronParse(params: any) {
        const { spec, next } = params;
        return http[scriptManager.checkCronParse(params).method](
          scriptManager.checkCronParse(params).url,
          {
            spec, 
            next 
          }
        );
    },
    taskRun(params: any) {
        const { host_id } = params;
        return http[scriptManager.taskRun(params).method](
          scriptManager.taskRun(params).url,
          {
            host_id
          }
        );
    },
    getTaskLog(params: any) {
        const { limit, start, execStatus } = params;
        return http[scriptManager.getTaskLog(params).method](
          scriptManager.getTaskLog(params).url,
          {
            'sort-by': 'id',
            'sort-dir': 'desc',
            'exec-status': execStatus,
            limit,
            start
          }
        );
    },
    deleteTask(params: any) {
        return http[scriptManager.deleteTask(params).method](
          scriptManager.deleteTask(params).url
        );
    },
    // 编辑脚本
    editScript(params: any) {
      const { describe, exec_timeout, log_retention } = params;
      return http[scriptManager.editScript(params).method](
        scriptManager.editScript(params).url,
        {
          describe,
          exec_timeout,
          log_retention
        }
      );
    }
}