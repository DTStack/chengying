import * as React from 'react';
import { Row, Col, Button } from 'antd';
import InputBox from './inputBox';
import { isEqual, cloneDeep } from 'lodash';
import { connect } from 'react-redux';
import { Dispatch, bindActionCreators } from 'redux';
import { AppStoreTypes } from '@/stores';
import * as installGuideAction from '@/actions/installGuideAction';
interface Prop {
  instanceData: any;
  saveParamValue: Function;
  setAllConfig: Function;
  saveAllConfig: Function;
  resetParamFieldvalue: Function;
  handleParamBlur: Function;
  isCloud?: boolean;
  deployState?: string;
  actions?: any;
}
interface State {
  localState: any;
  configEditState: 'edit' | 'normal';
  beforeEditConfig: any;
  editFiled: string[];
}
const mapStateToProps = (state: AppStoreTypes) => ({
  deployState: state.InstallGuideStore.deployState,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(
    Object.assign({}, { editDeployState: installGuideAction.editDeployState }),
    dispatch
  ),
});
@(connect(mapStateToProps, mapDispatchToProps) as any)
class InstanceConfig extends React.Component<Prop, State> {
  state: State = {
    localState: {},
    configEditState: 'normal',
    beforeEditConfig: {},
    editFiled: [],
  };

  componentDidMount() {
    this.initLocalState();
  }

  componentDidUpdate(prevProps: Prop) {
    if (!isEqual(prevProps.instanceData, this.props.instanceData)) {
      this.initLocalState();
    }
  }

  initLocalState = () => {
    const { instanceData, isCloud } = this.props;
    console.log(instanceData);

    const instanceDataList = Object.keys(instanceData);
    const localState = {};
    if (JSON.stringify(instanceData) === '{}' || isCloud) {
      return;
    }
    instanceDataList.forEach((o: any, i: number) => {
      if (instanceData[o] instanceof Array) {
        // 数组
        instanceData[o].forEach((q: any, j: number) => {
          if (!localState[`${o}.${j}`]) {
            localState[`${o}.${j}`] = `${q}`;
          }
        });
      } else if (
        typeof instanceData[o] === 'string' ||
        typeof instanceData[o] === 'number'
      ) {
        // 字符串 或 number
        if (!localState[`${o}`]) {
          localState[`${o}`] = instanceData[o];
        }
      } else {
        // json
        Object.keys(instanceData[o]).forEach((m: any, j: number) => {
          if (!localState[`${o}.${m}`]) {
            localState[`${o}.${m}`] = instanceData[o][m];
          }
        });
      }
    });
    this.setState({ localState });
  };

  handleInputChange = (newValue: any, oldvalue: string, name: string) => {
    if (newValue.target.value !== oldvalue) {
      const key = `Instance.${name}`;
      const editFiled = [...this.state.editFiled];
      if (editFiled.indexOf(key) === -1) {
        editFiled.push(key);
        this.setState({
          editFiled,
        });
      }
      this.props.saveParamValue(newValue.target.value, key, 'Instance');
    }
  };

  handleBlur = (e: any, f: string) => {
    const { localState } = this.state;
    if (localState[f].toString() !== e.target.value) {
      this.props.handleParamBlur(e.target.value, `Instance.${f}`);
      this.setState({
        localState: Object.assign({}, localState, {
          [f]: e.target.value,
        }),
      });
    }
  };

  renderNode = (state) => {
    const { instanceData, isCloud } = this.props;
    const disabled = state === 'normal';
    if (JSON.stringify(instanceData) === '{}' || isCloud) {
      return <p>无</p>;
    } else {
      const instanceDataList = Object.keys(instanceData);
      const envName = 'environment';
      const envIdx = instanceDataList.findIndex((item) => item === envName);

      if (envIdx !== -1) {
        instanceDataList.splice(envIdx, 1);
        instanceDataList.push(envName);
      }
      return instanceDataList.map((o: any, i: number) => {
        if (instanceData[o] instanceof Array) {
          // 数组
          return (
            <Row style={{ marginBottom: 20 }} key={i}>
              <Col span={10}>
                <p className="param-input-label">{o}：</p>
              </Col>
              {instanceData[o].map((q: any, j: number) => {
                return (
                  <Row style={{ marginBottom: 10 }} key={j}>
                    <Col span={10} />
                    <Col span={14}>
                      <InputBox
                        inputDisabled={disabled}
                        tooltipOnClick={() =>
                          this.props.resetParamFieldvalue({
                            field_path: `Instance.${o}.${j}`,
                            type: '1',
                          })
                        }
                        defaultvalue={q}
                        // onBlur={(e) => {
                        //     this.handleBlur(e, `${o}.${j}`);
                        // }}
                        onChange={(e) =>
                          this.handleInputChange(
                            e,
                            instanceData[o],
                            `${o}.${j}`
                          )
                        }
                        disabled={o === 'MaxReplica'}
                        title={o}
                      />
                    </Col>
                  </Row>
                );
              })}
            </Row>
          );
        } else if (
          typeof instanceData[o] === 'string' ||
          typeof instanceData[o] === 'number'
        ) {
          // 字符串 或 number
          return (
            <Row style={{ marginBottom: 20 }} key={i}>
              <Col span={10}>
                <p className="param-input-label">{o}：</p>
              </Col>
              <Col span={14}>
                <InputBox
                  inputDisabled={disabled}
                  tooltipOnClick={() =>
                    this.props.resetParamFieldvalue({
                      field_path: `Instance.${o}`,
                      type: '1',
                    })
                  }
                  defaultvalue={instanceData[o]}
                  // onBlur={(e) => {
                  //     this.handleBlur(e, `${o}`);
                  // }}
                  onChange={(e) =>
                    this.handleInputChange(e, instanceData[o], `${o}`)
                  }
                  disabled={o === 'MaxReplica'}
                  title={o}
                />
              </Col>
            </Row>
          );
        } else {
          // json
          return (
            <Row style={{ marginBottom: 20 }} key={i}>
              {/* <Row>
                                <Col span={6}>
                                    <p style={{fontWeight: 'bold', marginRight: '25%', marginBottom: 10}}>{o}:</p>
                                </Col>
                                </Row> */}
              {Object.keys(instanceData[o]).map((m: any, j: number) => {
                return (
                  <Row style={{ marginBottom: 10 }} key={j}>
                    <Row>
                      <Col className="param-input-label" span={10}>
                        {o === envName && (
                          <i className="emicon emicon-env mr-8" />
                        )}
                        <span>
                          {o}.{m}：
                        </span>
                      </Col>
                      <Col span={14}>
                        <InputBox
                          inputDisabled={disabled}
                          tooltipOnClick={() =>
                            this.props.resetParamFieldvalue({
                              field_path: `Instance.${o}.${m}`,
                              type: '1',
                            })
                          }
                          // onBlur={(e) => {
                          //     this.handleBlur(e, `${o}.${m}`);
                          // }}
                          onChange={(e) =>
                            this.handleInputChange(e, m, `${o}.${m}`)
                          }
                          defaultvalue={instanceData[o][m]}
                          disabled={o === 'MaxReplica'}
                          title={o}
                        />
                      </Col>
                    </Row>
                  </Row>
                );
              })}
            </Row>
          );
        }
      });
    }
  };

  // 运行配置编辑
  runtimeEdit = () => {
    this.setState({
      configEditState: 'edit',
      beforeEditConfig: cloneDeep(this.props.instanceData),
    });
    this.props.actions.editDeployState('edit');
  };

  // 运行配置取消
  runtimeCancel = () => {
    this.props.setAllConfig('Instance', this.state.beforeEditConfig);
    const editObj = this.contrastAndReturnEditData(
      this.props.instanceData,
      this.state.beforeEditConfig
    );
    this.setState({
      configEditState: 'normal',
    });
    this.props.actions.editDeployState('normal');
    if (JSON.stringify(editObj) === '{}') return;
    this.props.saveAllConfig('Instance', editObj, this.state.editFiled);
  };

  // 运行配置保存
  runtimeSave = () => {
    this.setState({
      configEditState: 'normal',
    });
    const editObj = this.contrastAndReturnEditData(
      this.state.beforeEditConfig,
      this.props.instanceData
    );
    this.props.actions.editDeployState('normal');
    if (JSON.stringify(editObj) === '{}') return;
    this.props.saveAllConfig('Instance', editObj, this.state.editFiled, true);
  };

  // 对比新老数据返回修改的数据
  contrastAndReturnEditData = (oldData, newData) => {
    const obj: any = {};
    for (const filed in newData) {
      const newValue = newData[filed];
      const oldValue = oldData[filed];
      if (newValue != oldValue) {
        obj[filed] = newValue;
      }
    }
    return obj;
  };

  render() {
    const { configEditState } = this.state;
    const { instanceData } = this.props;
    const { isCloud } = this.props;
    return (
      <div className="config-param_COMP">
        {JSON.stringify(instanceData) === '{}' || isCloud ? null : (
          <div className="handle-btn">
            {configEditState === 'normal' ? (
              <div>
                <Button type="link" onClick={this.runtimeEdit}>
                  编辑
                </Button>
              </div>
            ) : (
              <div>
                <Button type="link" onClick={this.runtimeSave}>
                  保存
                </Button>
                <Button type="link" onClick={this.runtimeCancel}>
                  取消
                </Button>
              </div>
            )}
          </div>
        )}
        {this.renderNode(configEditState)}
      </div>
    );
  }
}

export default InstanceConfig;
