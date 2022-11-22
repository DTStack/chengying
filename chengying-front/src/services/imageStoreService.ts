import apis from '@/constants/apis';
import * as http from '@/utils/http';
const { imageStore } = apis;

export default {
  // 获取集群下镜像仓库列表
  getImageStoreList(params: { cluster_id: number }) {
    return http[imageStore.getImageStoreList(params).method](
      imageStore.getImageStoreList(params).url
    );
  },

  // 获取镜像仓库信息
  getImageStoreInfo(params: { store_id: number }) {
    return http[imageStore.getImageStoreInfo(params).method](
      imageStore.getImageStoreInfo(params).url,
      params
    );
  },

  // 创建镜像仓库
  createImageStore(params: any) {
    return http[imageStore.createImageStore.method](
      imageStore.createImageStore.url,
      params
    );
  },

  // 更新镜像仓库
  updateImageStore(params: any) {
    return http[imageStore.updateImageStore.method](
      imageStore.updateImageStore.url,
      params
    );
  },

  // 删除镜像仓库
  deleteImageStore(params: { id: any[] }) {
    return http[imageStore.deleteImageStore.method](
      imageStore.deleteImageStore.url,
      params
    );
  },

  // 设置默认仓库
  setDefaultStore(params: any) {
    return http[imageStore.setDefaultStore.method](
      imageStore.setDefaultStore.url,
      params
    );
  },
};
