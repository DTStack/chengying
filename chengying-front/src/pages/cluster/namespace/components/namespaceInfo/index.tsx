import * as React from 'react';
import { Tabs } from 'antd';
import OverviewPage from './overview/index';
import EventView from './eventView/index';

const TabPane = Tabs.TabPane;

interface IProps {
  namespace: string;
  history?: any;
}
interface IState {
  activeKey: string;
}
export default class NamespaceInfo extends React.PureComponent<IProps, IState> {
  state: IState = {
    activeKey: 'overview',
  };

  render() {
    const { activeKey } = this.state;
    const { namespace, history } = this.props;
    const componentProps = {
      namespace,
      history,
    };
    return (
      <Tabs
        activeKey={activeKey}
        onChange={(activeKey: string) => this.setState({ activeKey })}>
        <TabPane key="overview" tab="服务概览">
          {activeKey === 'overview' && <OverviewPage {...componentProps} />}
        </TabPane>
        <TabPane key="eventview" tab="事件查看">
          {activeKey === 'eventview' && <EventView {...componentProps} />}
        </TabPane>
      </Tabs>
    );
  }
}
