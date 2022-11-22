import * as React from 'react';
import { Result } from 'antd';

export default () => {
  return (
    <Result
      status="error"
      title="程序异常"
      subTitle="程序出现异常，请联系技术支持！"
    />
  );
};
