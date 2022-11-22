import { message } from 'antd';
import { Service } from '@/services';
import { HeaderActions } from '@/constants/actionTypes';
import * as Http from '@/utils/http';
import * as Cookies from 'js-cookie';
import { COOKIES } from '@/constants/const';
import utils from '@/utils/utils';

export const getParentProductList = (currentParentProduct?: any) => {
  return (dispatch: any) => {
    Service.getParentProductList({ limit: 0 }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const data: any = {
          name: res.data[0],
          list: res.dsta,
        };
        let matchIndex = -1;
        if (currentParentProduct) {
          for (const i in res.data) {
            if (res.data[i] === currentParentProduct) {
              matchIndex = parseInt(i, 10);
            }
          }
          data.name = matchIndex > -1 ? res.data[matchIndex] : res.data[0];
        }
        Cookies.set('em_current_parent_product', data.name);
        dispatch({
          type: HeaderActions.GET_PARENT_PROD_LIST,
          payload: data,
        });
      }
    });
  };
};

export interface ClusterItem {
  clusterId?: number;
  clusterName?: string;
  clusterType?: string;
  mode: 0 | 1;
  subdomain: {
    [key: string]: string[];
  };
}

// 获取集群-产品树
export const getClusterProductList = () => {
  function getProductName(
    currentParentName: string,
    subdomain: any,
    mode: 0 | 1
  ): string {
    if (mode === 0) {
      Cookies.set(COOKIES.NAMESPACE, '');
      return currentParentName || subdomain[0];
    } else {
      let isCurrentParentNameFound = false;
      Object.keys(subdomain).forEach((namespace) => {
        (subdomain[namespace] || []).forEach((p) => {
          // cookie存在 em_current_parent_product的情况，且namespace不存在情况
          if (
            !isCurrentParentNameFound &&
            p === currentParentName &&
            !Cookies.get(COOKIES.NAMESPACE)
          ) {
            isCurrentParentNameFound = true;
            Cookies.set(COOKIES.NAMESPACE, namespace);
          }
        });
      });
      // cookie中em_current_parent_product和namespace不存在情况，设置默认cookie namespace
      if (!isCurrentParentNameFound && !Cookies.get(COOKIES.NAMESPACE)) {
        const namespace = Object.keys(subdomain)[0];
        Cookies.set(COOKIES.NAMESPACE, namespace);
        return subdomain[Object.keys(subdomain)[0]][0];
      } else {
        return currentParentName;
      }
    }
  }
  return (dispatch: any) => {
    Service.getClusterProductList({ limit: 0 }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        const currentParentProduct =
          Cookies.get('em_current_parent_product') || 'DTinsight';
        const currentParentClusterId = Cookies.get('em_current_cluster_id');
        let cluster: any = {};
        if (currentParentClusterId) {
          const foundCluster = res.data.find(
            (item) => +item.clusterId === +currentParentClusterId
          );
          cluster = foundCluster || res.data[0];
          // 兼容默认为k8s状态
          if (!foundCluster && cluster?.clusterType === 'kubernetes') {
            const loading = message.loading('正在跳转默认集群...', 0);
            setTimeout(loading, 2500);
            utils.setNaviKey('menu_ops_center', 'sub_menu_service');
            window.location.href = `/opscenter/service`;
          }
        } else {
          cluster = res.data[0];
        }
        const productName = getProductName(
          currentParentProduct,
          cluster?.subdomain || {},
          cluster?.mode
        );
        if (cluster?.clusterId) {
          Cookies.set('em_current_cluster_id', cluster?.clusterId || -1);
          Cookies.set('em_current_cluster_type', cluster?.clusterType);
          Cookies.set('em_current_parent_product', productName);
          if (cluster?.clusterType === 'kubernetes') {
            const namespace = Object.keys(cluster?.subdomain);
            Cookies.set('em_current_k8s_namespace', namespace[0]);
          }
        }
        dispatch({
          type: HeaderActions.GET_PARENT_PROD_LIST,
          payload: Object.assign(
            {},
            {
              name: productName || '选择产品',
              list: res.data,
              cluster: {
                id: cluster?.clusterId || -1,
                name: cluster?.clusterName || '选择集群',
                type: cluster?.clusterType || 'hosts',
                mode: cluster?.mode || 0,
              },
            }
          ),
        });
      } else {
        console.log(res.msg);
      }
    });
  };
};

export const getProductList = () => {
  return (dispatch: any) => {
    // 获取产品下的组件
    Http.get('/api/v2/product', { limit: 0 }).then((data: any) => {
      data = data.data;
      const list = [];
      if (data.code === 0) {
        for (const p of data.data.list) {
          if (p.is_current_version === 1) {
            list.push(p);
          }
        }
        dispatch({
          type: HeaderActions.SET_PRODUCT_LIST,
          payload: list,
        });
      } else {
        message.error(data.msg);
      }
    });
  };
};

export const setCurrentProduct = (params: any) => {
  return {
    type: HeaderActions.SET_CUR_PRODUCT,
    payload: params,
  };
};

export const setCurrentParentProduct = (pn: any) => {
  return {
    type: HeaderActions.SET_CUR_PARENT_PROD,
    payload: pn,
  };
};

export const getServiceGroupList = () => {
  return (dispatch: any) => {
    Http.get('/api/container/getCluster', {}).then((data: any) => {
      data = data.data;
      const result = data.result.data || [];
      for (let i = 0, len = result.length; i < len; i++) {
        const item = result[i].services;
        const newServices = [];
        for (let j = 0, len2 = item.length; j < len2; j++) {
          const subitem = item[j];
          const name = subitem.name.toLowerCase();
          if (name !== 'rabbitmq' && name !== 'mysql') {
            newServices.push(subitem);
          }
        }
        result[i].services = newServices;
      }
      dispatch({
        type: HeaderActions.GET_S_G_LIST,
        payload: result,
      });
    });
  };
};
export const getInstanceList = () => {
  return (dispatch: any) => {
    Http.get('/api/container/queryContainer', { pageSize: 200 }).then(
      (data: any) => {
        data = data.data;
        dispatch({
          type: HeaderActions.GET_INSTANCE_LIST,
          payload: data.result.data,
        });
      }
    );
  };
};

export const getClusterList = (name?: string) => {
  return (dispath: any) => {
    Service.getClusterList({
      type: undefined,
      'sort-by': 'id',
      'sort-dir': 'desc',
      limit: 0,
      start: 0,
    }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        dispath({
          type: HeaderActions.GET_PARENT_CLUSTER_LIST,
          payload: res.data.clusters,
        });
        // set cluster
        const defaultCluster = res.data.clusters[0];
        const id = Cookies.get('em_current_cluster_id');
        const cluster = id
          ? res.data.clusters.find((item) => item.id === +id) || defaultCluster
          : defaultCluster;
        dispath({
          type: HeaderActions.SET_CUR_PARENT_CLUSTER,
          payload: cluster,
        });
        Cookies.set('em_current_cluster_id', cluster.id);
        Cookies.set('em_current_cluster_type', cluster.type);
      } else {
        message.error(res.msg);
      }
    });
  };
};

export const setCurrentParentCluster = (payload: any) => {
  return {
    type: HeaderActions.SET_CUR_PARENT_CLUSTER,
    payload,
  };
};

export interface HeaderActionTypes {
  getClusterProductList: Function;
  getParentProductList: Function;
  getProductList: Function;
  setCurrentProduct: Function;
  setCurrentParentProduct: Function;
  getServiceGroupList: Function;
  getInstanceList: Function;
  getClusterList: Function;
  setCurrentParentCluster: Function;
}
