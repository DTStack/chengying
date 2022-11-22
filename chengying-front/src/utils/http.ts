// tslint:disable
import 'url-search-params-polyfill';
import axios from 'axios';
import { message } from 'antd';
require('es6-promise').polyfill();

// let count = 1;
// 封装好的get和post接口，调用方法情况action文件
const instance = axios.create({
  // baseURL: API_URL, //设置默认api路径
  timeout: 180000, // 设置超时时间
  headers: {
    // contentType: "application/json",
  },
  withCredentials: true,
});
// let blobInstance = axios.create({
//     timeout: 10000, //设置超时时间
//     headers: {
//         'X-Custom-Header': 'dtuic'
//     },
//     withCredentials:true,
//     responseType:'blob'
// });
// 拦截器，统一处理请求
instance.interceptors.request.use(
  (config) => {
    //  console.log('interceptor request:',config);
    return config;
  },
  (error) => {
    console.log('request error:', error);
    return Promise.reject(error);
  }
);
// 拦截器，统一处理未登录或者没有权限的情况
instance.interceptors.response.use(
  (response) => {
    if (response.data.code === 103 && window.APPCONFIG.userCenter) {
      window.location.href = '/login';
    }
    return response;
  },
  (error) => {
    console.log(error.response);
    if (error.response.status === 403) {
      message.destroy();
      message.error('权限不足,请联系管理员!');
      // window.location.href = '/nopermission'
    }
    return Promise.reject(error.response.data);
    // 返回接口返回的错误信息
  }
);

export const get = (url: string, param: any) => {
  return instance.get(`${url}`, { params: param });
};

/**
 * 后端需要莫名其妙的post格式
 * @param {String} url
 * @param {Obj} param
 */
export const postFormData = (url: string, param: any) => {
  // const formParams = new URLSearchParams();
  const formParams = new FormData();
  for (const p in param) {
    formParams.append(p, param[p]);
  }
  return instance.post(`${url}`, formParams);
};
/**
 * post json
 * @param {String} url
 * @param {Obj} param
 */
export const post = (url: string, param: any) => {
  return instance.post(`${url}`, param);
};

export const dele = (url: string, param: any) => {
  return instance.delete(url, param);
};
export const put = (url: string, param: any) => {
  return instance.put(url, param);
};

/**
 * 传统post方式
 * @param {String} url
 * @param {Obj} param
 */
export const postJsonData = (url: string, param: any) => {
  return instance.post(`${url}`, param);
};

export const getMultiData = (getFuncArr: any) => {
  return axios.all(getFuncArr);
};
