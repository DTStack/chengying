import * as React from 'react';
import { Form, Radio } from 'antd';
import Account from './account.comp';
import Shell from './shell.comp';
import { formLayout } from './constant';
import './style.scss';

const FormItem = Form.Item;
const RadioGroup = Radio.Group;

interface Prop {
  afterInstall: () => void;
  //   refreshHost?: () => void
  //   refreshGroup?: () => void
  refresh?: (params) => void;
  clusterInfo?: any;
  onCancel?: () => void;
}

interface State {
  radioType: string;
}

class CreateNewHost extends React.Component<Prop, State> {
  state: State = {
    radioType: 'account',
  };

  handleRadioChange = (e) => {
    console.log(e);
    this.setState({ radioType: e.target.value });
  };

  render() {
    const { radioType } = this.state;
    const { afterInstall, refresh, clusterInfo, onCancel } = this.props;
    const addHostProps = {
      afterInstall,
      refresh,
      clusterInfo,
      onCancel,
    };
    return (
      <div className="create-host-container">
        <Form>
          <FormItem label="接入方式" {...formLayout}>
            <RadioGroup value={radioType} onChange={this.handleRadioChange}>
              <Radio value="account">账号接入</Radio>
              <Radio value="shell">命令行接入</Radio>
            </RadioGroup>
          </FormItem>
        </Form>
        {radioType === 'account' ? (
          <Account {...addHostProps} />
        ) : (
          <Shell {...addHostProps} />
        )}
      </div>
    );
  }
}

export default CreateNewHost;
