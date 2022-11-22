import * as React from 'react';
import { Modal, Button } from 'antd';

interface IAutoTestProps {
  visible: boolean;
  autoTestInfo: any;
  autoTestMsg: string;
  handleClose: () => void;
  handleTest: () => void;
}

const AutoTest: React.FC<IAutoTestProps> = (props) => {
  const { visible, autoTestInfo, autoTestMsg, handleTest, handleClose } = props;
  return (
    <Modal
      title="自动化测试"
      footer={null}
      visible={visible}
      onCancel={handleClose}>
      <div className="auto-smoke">
        <div>
          点击按钮开始自动化测试
          <Button
            type="primary"
            style={{ marginLeft: 20 }}
            disabled={autoTestInfo.exec_status === 1}
            onClick={handleTest}>
            开始测试
          </Button>
          {autoTestInfo.exec_status === 1 && (
            <Button type="link" loading>
              运行中
            </Button>
          )}
        </div>
        <p>最近一次运行时间：{autoTestInfo.end_time || '--'}</p>
        <p>
          查看报告：
          {!autoTestInfo?.exec_status && '--'}
          {autoTestInfo.exec_status === 3 && (
            <span className="error">运行失败</span>
          )}
          {autoTestInfo.report_url && autoTestInfo.exec_status === 2 && (
            <a href={autoTestInfo.report_url} target="_blank">
              测试报告
            </a>
          )}
        </p>
        {autoTestInfo.exec_status === 3 && autoTestMsg && (
          <div className="exec-error-log">
            错误日志：
            <code>{autoTestMsg}</code>
          </div>
        )}
      </div>
    </Modal>
  );
};

export default AutoTest;
