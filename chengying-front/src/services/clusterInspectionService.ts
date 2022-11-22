import apis from '@/constants/apis';

import * as http from '@/utils/http';
const { clusterInspectionApi } = apis;

export default {
    // 获取集群巡检相关的
    getClusterInspectionStatisSetData(params?: any) {
        return http[clusterInspectionApi.getClusterInspectionStatisSet.method](
            clusterInspectionApi.getClusterInspectionStatisSet.url,
            params
        )
    },
    getClusterInspectionBaseInfoData(params?: any) {
        return http[clusterInspectionApi.getClusterInspectionBaseInfo.method](
            clusterInspectionApi.getClusterInspectionBaseInfo.url,
            params
        )
    },
    getClusterInspectionTableData(params?: any) {
        return http[clusterInspectionApi.getClusterInspectionTable.method](
            clusterInspectionApi.getClusterInspectionTable.url,
            params
        )
    },
    // 获取应用服务运行接口
    getApplicationService(params?: any) {
        return http[clusterInspectionApi.getClusterrApplicationServerInfo.method](
            clusterInspectionApi.getClusterrApplicationServerInfo.url,
            params
        )
    },
    getApplicationGraphData(params?: any) {
        return http[clusterInspectionApi.getClusterApplicationServerTableInfo.method](
            clusterInspectionApi.getClusterApplicationServerTableInfo.url,
            params
        )
    },
    // 大数据服务运行
    getClusterBigDataServer(params?: any) {
        return http[clusterInspectionApi.getClusterBigDataInfo.method](
            clusterInspectionApi.getClusterBigDataInfo.url,
            params
        )
    }, 
    getSettingData(params?: any) {
        return http[clusterInspectionApi.getSettingData.method](
            clusterInspectionApi.getSettingData.url,
            params
        )
    }
};
