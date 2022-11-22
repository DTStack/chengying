import * as React from 'react';
import { Row, Col, Select, Button } from 'antd';
import * as Cookie from 'js-cookie';
const Option = Select.Option;

/**
 *
 * @param props
 * @returns
 */
const ConfigToolBar: React.FC<any> = (props) => {
  const {
    cur_service,
    configFile,
    handleFileChange,
    isRestart,
    handleRestartServiceInTurn,
    couldSaveConfig,
    handleSaveServiceConfig,
    distributeServiceConfig,
    handleAddParamConfigShow,
    handleResetServiceConfig,
    isKubernetes,
  } = props;
  return (
    <Row className="c-paramConfig-toolbar">
      <Col span={8} style={{ display: 'flex', alignItems: 'center' }}>
        配置文件：
        <Select
          placeholder="请选择配置文件"
          value={cur_service?.Instance?.ConfigPaths ? configFile : undefined}
          style={{ width: '264px', marginRight: '-50px' }}
          onChange={handleFileChange}>
          {cur_service?.Instance?.ConfigPaths && (
            <Option value={'all'}>全部</Option>
          )}
          {(cur_service?.Instance?.ConfigPaths || []).map((file: string) => (
            <Option key={file} value={file}>
              {file}
            </Option>
          ))}
        </Select>
      </Col>
      <Button
        disabled={!couldSaveConfig}
        type="primary"
        style={{ marginRight: 10 }}
        icon="save"
        onClick={handleSaveServiceConfig}>
        保存
      </Button>
      <Button
        disabled={!couldSaveConfig}
        type="default"
        icon="redo"
        className={couldSaveConfig ? 'c-paramConfig-toolbar__btn' : null}
        onClick={handleResetServiceConfig}
        style={{ marginRight: 10 }}>
        重置
      </Button>
      <Button
        disabled={
          !!couldSaveConfig ||
          !cur_service?.Instance?.ConfigPaths ||
          !cur_service?.Instance?.ConfigPaths.length
        }
        type="primary"
        icon="plus"
        style={{ marginRight: 10 }}
        onClick={() => handleAddParamConfigShow(true)}>
        添加参数
      </Button>
      {!isKubernetes && Cookie.get('em_current_cluster_id') && (
        <>
          <Button
            type="primary"
            ghost
            // className="c-paramConfig-toolbar__btn"
            onClick={distributeServiceConfig}
            style={{ marginRight: 10 }}
            disabled={!cur_service?.Instance?.ConfigPaths?.length || couldSaveConfig}
            icon="setting">
            配置下发
          </Button>
          <Button
            type="default"
            className={isRestart ? '' : 'c-paramConfig-toolbar__btn'}
            disabled={isRestart}
            loading={isRestart}
            onClick={handleRestartServiceInTurn}>
            {!isRestart && (
              <i
                className="emicon emicon-poweroff"
                style={{
                  lineHeight: '12px',
                  verticalAlign: '-2px',
                  marginRight: 5,
                }}
              />
            )}
            滚动重启
          </Button>
        </>
      )}
    </Row>
  );
};

export default ConfigToolBar;
