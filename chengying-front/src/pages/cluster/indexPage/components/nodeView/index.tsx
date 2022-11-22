import * as React from 'react';
import { cloneDeep } from 'lodash';
import { pieOption } from '@/pages/cluster/indexPage/constant';
import Chart from '@/components/chart';

interface IProps {
  overview: any;
}
export default class NodeView extends React.PureComponent<IProps, any> {
  setOption = () => {
    const { overview = {} } = this.props;
    const { nodes = 0, error_nodes = 0 } = overview;
    const option = cloneDeep(pieOption);
    const data = [
      {
        name: '正常',
        value: nodes - error_nodes,
      },
      {
        name: '异常',
        value: error_nodes,
      },
    ];

    // 重置
    if (nodes === 0) {
      option.color = ['#C9D0D4'];
    }
    option.title[0].text = `{value|${nodes}}\n{name|总量 (个)}`;
    option.legend.formatter = (name: string) => {
      const obj: any = data.find((item) => item.name === name) || {};
      return `{name|${name}}{value|${obj.value}}`;
    };
    option.series[0].data = data;
    return option;
  };

  render() {
    const option = this.setOption();
    return (
      <div className="index-view-card">
        <p className="title">节点</p>
        <Chart
          style={{ marginLeft: '-20%' }}
          option={option}
          width="100%"
          height="100%"
        />
      </div>
    );
  }
}
