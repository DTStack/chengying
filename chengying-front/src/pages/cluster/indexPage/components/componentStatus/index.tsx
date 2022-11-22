import * as React from 'react';
import { Icon, Tooltip, Empty } from 'antd';
import './style.scss';

interface IProps {
  overview: any;
}

export default class ComponentStatus extends React.PureComponent<IProps, any> {
  render() {
    const { overview = {} } = this.props;
    return (
      <div className="index-view-card">
        <p className="title">组件状态</p>
        <div style={{ padding: 38 }}>
          {Array.isArray(overview.component) && overview.component.length ? (
            overview.component.map((item: any) => (
              <div
                key={item.role}
                className={
                  'component-status-item ' +
                  (item.status === 0 ? 'item-success' : 'item-error')
                }>
                <div className="item-icon">
                  <Icon
                    type={item.status === 0 ? 'check' : 'exclamation'}
                    style={{ color: '#fff' }}
                  />
                </div>
                <div className="item-content">{item.role}</div>
                {item.status === 1 && (
                  <Tooltip title={item.errors.join(',')}>
                    <Icon
                      className="tooltip-error"
                      type="exclamation-circle"
                      style={{ color: '#999', alignSelf: 'center' }}
                    />
                  </Tooltip>
                )}
              </div>
            ))
          ) : (
            <Empty className="c-ant_empty-center" />
          )}
        </div>
      </div>
    );
  }
}
