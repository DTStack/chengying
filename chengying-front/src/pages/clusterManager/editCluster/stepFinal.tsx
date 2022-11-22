import * as React from 'react';
import Host from '@/pages/host';
import KubernetesCmd from './KubernetesCmd';

interface IProps {
  isEdit: boolean;
  clusterInfo: any;
  handleCancel: (event: React.MouseEvent<HTMLElement>) => void;
}

export default class StepFinal extends React.PureComponent<IProps, any> {
  render() {
    const { clusterInfo } = this.props;
    const hostProps = {
      ...this.props,
      clusterInfo,
    };
    return (
      <div data-testid="ec-step-final">
        {clusterInfo.type === 'kubernetes' && clusterInfo.mode === 1 ? (
          <KubernetesCmd clusterInfo={clusterInfo} />
        ) : (
          <React.Fragment>
            <p className="c-title__color">选择主机</p>
            <Host {...hostProps} />
          </React.Fragment>
        )}
      </div>
    );
  }
}
