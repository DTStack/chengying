import * as React from 'react';
import { Card, Input, Button, Icon } from 'antd';
import classnames from 'classnames';
import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/mode-yaml';
import 'ace-builds/src-noconflict/theme-kuroir';
import './style.scss';

const { TextArea } = Input;

export interface YamlEditorProps {
  yaml: string;
  readOnly?: boolean;
  onChange?: (value: string, event?: any) => void;
}
const YamlEditor: React.FC<YamlEditorProps> = (props) => {
  const { yaml, readOnly } = props;
  return (
    <Card
      title="YAML"
      size="small"
      className="c-yaml-editor__ant-card"
      extra={
        <div style={{ display: 'flex' }}>
          <TextArea
            id="J_CMDCode"
            value={yaml}
            style={{
              position: 'absolute',
              opacity: 0,
              zIndex: -10,
            }}
          />
          <Button
            id="J_CopyBtn"
            type="link"
            style={{ padding: 0 }}
            data-clipboard-action="cut"
            data-clipboard-target="#J_CMDCode">
            <Icon type="copy" /> 复制到剪贴板
          </Button>
        </div>
      }>
      <AceEditor
        className={classnames({ yaml_editor_readonly: readOnly })}
        mode="yaml"
        theme="kuroir"
        value={yaml}
        readOnly={readOnly || false}
        width="100%"
        height="100%"
        onChange={props.onChange}
        name="ADVANCED_MODE_DIV"
      />
    </Card>
  );
};
export default YamlEditor;
