import * as React from 'react';
import 'public/fonts/emfont/iconfont.css';
const IcoLoad = require('public/imgs/icon_loading.svg');

export default () => {
  return (
    <div className="loading-box">
      <img className="anticon-spin" src={IcoLoad.default} width={16} />
    </div>
  );
};
