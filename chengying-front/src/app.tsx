import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import configureStore from './configureStore';
// import './registerServiceWorker';
import zhCN from 'antd/es/locale/zh_CN';
import { ConfigProvider } from 'antd';
import '../public/styles/const.scss';
import 'ant-design-dtinsight-theme/theme/dt-theme/reset.less';
import 'ant-design-dtinsight-theme/theme/dt-theme/index.less';
import '../public/styles/base.scss';
import Routers from './router';
import moment from 'moment';
import 'moment/locale/zh-cn';
moment.locale('zh-cn');

const store = configureStore({});

ReactDOM.render(
  <ConfigProvider locale={zhCN}>
    <Provider store={store}>
      <Routers />
    </Provider>
  </ConfigProvider>,
  document.getElementById('root')
);
