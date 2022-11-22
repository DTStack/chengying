import * as React from 'react';
import { Modal, Icon } from 'antd';
export const alertModal = (
  runtimeState: string,
  deployState: string,
  callback?: Function
) => {
  if (runtimeState === 'edit' || deployState === 'edit') {
    const arr = [];
    if (runtimeState === 'edit') {
      arr.push('运行配置');
    }
    if (deployState === 'edit') {
      arr.push('部署配置');
    }
    Modal.warning({
      title: '配置信息未保存！',
      content: `${arr.join(',')}未保存，请先点击保存！`,
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okText: '确定',
      onCancel: () => {
        callback && callback(false);
      },
      onOk: () => {
        callback && callback(false);
      },
    });
    return false;
  }
  callback && callback(true);
  return true;
};
