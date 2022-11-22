import * as React from 'react';
import { Icon } from 'antd';

const NoPermission: React.FC<any> = () => {
  const style: React.CSSProperties = {
    lineHeight: '200px',
    textAlign: 'center',
    fontSize: '20px',
    fontWeight: 500,
  };
  return (
    <p style={style}>
      <Icon className="mr-8" type="frown-o" />
      您尚无访问权限，或该页面不存在
    </p>
  );
};
export default NoPermission;
