import * as React from 'react';
import { Modal, Select, Tooltip } from 'antd';
const { Option } = Select;

export interface HostInfoConfigProps {
  configVisible: boolean;
  curent_config_path: string;
  data: string;
  handleCloseHostConfig: () => void;
  handleConfigPathClick: (param: string) => void;
  config_paths: Array<any>;
}

const HostInfoConfig: React.FC<HostInfoConfigProps> = ({
  configVisible,
  curent_config_path,
  data,
  handleCloseHostConfig,
  handleConfigPathClick,
  config_paths,
}) => (
  <Modal
    title="配置信息"
    width={900}
    footer={null}
    visible={configVisible}
    onCancel={handleCloseHostConfig}>
    <div style={{ marginBottom: '10px' }}>
      <Select
        style={{ minWidth: 160 }}
        value={curent_config_path}
        onChange={handleConfigPathClick}>
        {config_paths.map((l, i) => (
          <Option key={`${i}`} value={l}>
            <Tooltip title={l}>{l}</Tooltip>
          </Option>
        ))}
      </Select>
    </div>
    <pre
      style={{
        backgroundColor: '#F5F5F5',
        border: '1px dashed #DDD',
        maxHeight: 500,
        overflowY: 'auto',
      }}
      className="service-config-pre">
      {data}
    </pre>
  </Modal>
);

export default HostInfoConfig;
