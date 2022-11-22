import * as React from 'react';
import '../style.scss';
import ComponentContainer from '@/pages/components/container';

export default class DeployedApp extends React.Component<any, any> {
  render() {
    return (
      <div className="app-manager-container">
        <p className="title mb-12">应用列表</p>
        <ComponentContainer {...this.props} />
      </div>
    );
  }
}
