import API from '@/constants/apis';
import * as http from '@/utils/http';


export default {
  getSecurity(params?: any) {
    return http[API.security.getSecurity.method](
      API.security.getSecurity.url,
      params
    );
  },
  setSecurity(params: any) {
    return http[API.security.setSecurity.method](
      API.security.setSecurity.url,
      params
    )
  }
};