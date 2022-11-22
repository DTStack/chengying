import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { productLine } = apis;

export default {
    getProductLine(params?: any) {
        return http[productLine.getProductLine.method](
            productLine.getProductLine.url,
            {
                'sort-by': 'create_time',
                'sort-dir': 'desc',
                'limit': 0
            }
        );
    },
    uploadProductLine(params: any) {
        return http[productLine.uploadProductLine.method](
            productLine.uploadProductLine.url,
            params
        )
    },
    deleteProductLine(params: any) {
        return http[productLine.deleteProductLine(params).method](
            productLine.deleteProductLine(params).url
        )
    }
}