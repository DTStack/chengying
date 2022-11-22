import apis from '@/constants/apis';
import * as http from '@/utils/http';
import {
  ServiceTree,
  ServiceFile,
  DistributeServiceConfig,
} from '@/model/apis';

const { product } = apis;

class ProductService {
  /**
   * 全部产品包-查看-配置-选择一个文件查看内容
   */
  public static getServiceTree(params: ServiceTree) {
    const req = product.getServiceTree(params);
    return http[req.method](req.url);
  }

  /**
   * 全部产品包-查看-配置-选择一个服务后返回文件列表
   */
  public static getServiceFile(params: ServiceFile) {
    const req = product.getServiceFile(params);
    return http[req.method](req.url, { file: params.file });
  }

  /**
   * 服务配置下发
   */
  public static distributeServiceConfig(params: DistributeServiceConfig) {
    const req = product.distributeServiceConfig(params);
    return http[req.method](req.url);
  }
}

export default ProductService;
