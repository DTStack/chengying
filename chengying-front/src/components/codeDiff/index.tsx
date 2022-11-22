import React, { useRef, useEffect, useState } from 'react';
import { Layout, Modal, Spin } from 'antd';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import 'monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution.js';
import './style.scss';

const { Sider, Content } = Layout;

export interface ConfDiffProps {
  visible: boolean;
  title: string;
  data: string[];
  extra?: boolean | React.ReactNode;
  handleCancle: () => void;
  handleSubmit?: () => void;
}
let diffEditor = null;

const ConfDiff: React.FC<ConfDiffProps> = (props) => {
  const { visible, title, data, extra, handleCancle, handleSubmit } = props;
  const [loading, changeStatusLoading] = useState<boolean>(false);
  const modalWidth: number = 939 + (extra ? 125 : 0);
  const diffRef = useRef(null);
  const defaultOptions = {
    language: 'yaml',
    folding: true,
    readOnly: true,
  };

  useEffect(() => {
    if (diffRef.current) {
      changeStatusLoading(true);
      const [source, target] = data;
      diffEditor = monaco.editor.createDiffEditor(diffRef.current, {
        ...defaultOptions,
      });
      changeStatusLoading(false);
      const originalModel = monaco.editor.createModel(source, 'text/plain');
      const modifiedModel = monaco.editor.createModel(target, 'text/plain');
      diffEditor.setModel({
        original: originalModel,
        modified: modifiedModel,
      });
    }
  }, [diffRef.current]);

  useEffect(() => {
    if (diffEditor) {
      const [source, target] = data;
      const { original, modified } = diffEditor.getModel();
      if (original?.getValue() !== source) {
        original?.setValue(source || '');
      }
      if (modified?.getValue() !== target) {
        modified?.setValue(target || '');
      }
    }
  }, [diffEditor, data]);

  return (
    <Modal
      className="code-diff__wrapper"
      title={title}
      width={modalWidth}
      visible={visible}
      onCancel={handleCancle}
      onOk={handleSubmit}>
      <Layout>
        {extra && <Sider width={125}>{extra}</Sider>}
        <Content>
          <div className="ace-title">
            <div className="title">当前版本</div>
            <div className="title">最新修改</div>
          </div>
          <Spin spinning={loading} tip="正在初始化...">
            <div ref={diffRef} style={{ height: 500 }}></div>
          </Spin>
        </Content>
      </Layout>
    </Modal>
  );
};

export default ConfDiff;
