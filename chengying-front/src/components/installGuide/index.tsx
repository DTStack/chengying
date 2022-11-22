import * as React from 'react';
import { Modal, Button } from 'antd';
import utils from '@/utils/utils';
declare var APP: any;

export interface IProps {
  history?: any;
  visible: boolean;
  onClose: () => void;
}

const InstallGuideModal: React.FC<IProps> = (props) => {
  const { visible, onClose } = props;

  React.useEffect(() => {}, []);

  const jumpGuidePath = () => {
    let path = '/deploycenter/appmanage/installs';
    utils.setNaviKey('menu_deploy_center', 'sub_menu_product_deploy');
    props.history.push(path);
    onClose();
  };

  const jumpSelectHost = () => {
    let path = '/deploycenter/cluster/list';
    utils.setNaviKey('menu_deploy_center', 'sub_menu_cluster_manage');
    props.history.push(path);
    onClose();
  };
  const titleShow = () => {
    // return `欢迎使用ChengYing产品部署管家！(${packageInfo.version})`
    return (
      <span>
        欢迎使用ChengYing产品部署管家！
        <span style={{ color: '#BFBFBF' }}>（V{APP.VERSION}）</span>
      </span>
    );
  };

  return (
    <Modal
      className="installGuideModal"
      title={titleShow()}
      onCancel={onClose}
      visible={visible}
      width={1000}
      footer={
        <Button type="primary" onClick={jumpGuidePath}>
          知道了
        </Button>
      }>
      <div>
        <span>
          <img
            src={require('./css/deploymentPatterns1.png')}
            style={{ width: '976px', height: 'auto', display: 'block' }}
          />
          <img
            src={require('./css/deploymentPatterns2.png')}
            style={{ width: '976px', height: 'auto', display: 'block' }}
          />
        </span>
        <div style={{ marginLeft: '28px' }}>
          若需要对已部署组件进行 “升级” 或 “回滚” 等管理操作，请前往
          <span
            onClick={jumpSelectHost}
            style={{ color: '#3F87FF', cursor: 'pointer' }}>
            部署中心
          </span>
          选择集群 &rarr; 已部署组件。
        </div>
      </div>
    </Modal>
  );
};

export default InstallGuideModal;
