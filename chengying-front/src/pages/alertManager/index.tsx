import * as React from 'react';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { Tabs } from 'antd';
import AlertRule from '@/pages/alertRule';
import AlertChannel from '@/pages/alertChannel';
import util from '@/utils/utils';
const TabPane = Tabs.TabPane;

interface IState {
  activeKey: string;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});
@(connect(mapStateToProps, undefined) as any)
export default class AlertManager extends React.PureComponent<any, IState> {
  state: IState = {
    activeKey: 'rule',
  };

  static getDerivedStateFromProps(nextProps: any, prevState: IState) {
    const { authorityList } = nextProps;
    const { activeKey } = prevState;
    const codeMap = {
      rule: 'sub_menu_alarm_record',
      channel: 'sub_menu_alarm_channel',
    };
    if (
      Object.keys(authorityList).length &&
      !authorityList[codeMap[activeKey]]
    ) {
      return {
        activeKey: activeKey === 'rule' ? 'channel' : 'rule',
      };
    }
    return null;
  }

  componentDidMount() {
    const urlParams: any = util.getParamsFromUrl(this.props.location.search);
    if ('activeKey' in urlParams) {
      this.setState({ activeKey: urlParams.activeKey });
    }
  }

  // 设置activeKey
  handleChange = (activeKey: string) => {
    this.setState({ activeKey });
  };

  render() {
    const { authorityList } = this.props;
    return (
      <Tabs
        className="c-tabs-padding"
        style={{ padding: 20 }}
        activeKey={this.state.activeKey}
        onChange={this.handleChange}>
        {authorityList.sub_menu_alarm_record && (
          <TabPane tab="告警规则" key="rule">
            <AlertRule history={this.props.history} />
          </TabPane>
        )}
        {authorityList.sub_menu_alarm_channel && (
          <TabPane tab="告警通道" key="channel">
            <AlertChannel />
          </TabPane>
        )}
      </Tabs>
    );
  }
}
