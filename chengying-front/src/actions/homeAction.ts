import { message } from 'antd';

export default {
  getServiceGroupList() {
    return (dispatch: any, getState: any, Api: any) => {
      Api.getProductList({ limit: 0 }).then((data: any) => {
        const list = [];
        if (data.code === 0) {
          for (const p of data.data.list) {
            if (p.is_current_version === 1) {
              list.push(p);
            }
          }
          dispatch({
            type: 'SET_PRODUCT_LIST',
            payload: list,
          });
        } else {
          message.error(data.msg);
        }
      });
    };
  },
  getConfig(params?: any) {
    return {
      type: 'SET_CUR_PRODUCT',
      payload: params,
    };
  },
};
