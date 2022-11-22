import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { globalConfig } = apis;

export default {
    getGlobalConfig () {
        return http[globalConfig.getGlobalConfig.method](
            globalConfig.getGlobalConfig.url
        );
    },
    setGlobalConfig (params: any) {
        return http[globalConfig.setGlobalConfig.method](
            globalConfig.setGlobalConfig.url,
            params
        );
    },
}
