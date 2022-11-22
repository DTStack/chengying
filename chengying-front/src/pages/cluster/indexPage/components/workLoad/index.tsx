import * as React from 'react';
import { cloneDeep } from 'lodash';
import Chart from '@/components/chart';
import { pieOption } from '@/pages/cluster/indexPage/constant';
import './style.scss';

interface IProps {
  overview: any;
}

interface WorkLoadParams {
  capacity: number;
  load: number;
}
export default class WorkLoad extends React.PureComponent<IProps, any> {
  setOption = (data: WorkLoadParams) => {
    const option = cloneDeep(pieOption);
    const percent =
      data.capacity === 0
        ? 0
        : Number(((data.load / data.capacity) * 100).toFixed(2));
    option.title[0] = {
      ...option.title[0],
      text: `{value|${percent}}%\n{name|${data.load}/${data.capacity}}`,
      textStyle: {
        rich: {
          value: {
            fontSize: 20,
            fontWeight: 600,
            color: '#333333',
            padding: [5, 0],
          },
          name: {
            fontSize: 12,
            color: '#666666',
          },
        },
      },
    };
    option.legend = undefined;
    option.series[0] = {
      ...option.series[0],
      data: [
        {
          name: 'load',
          value: data.load,
        },
        {
          name: 'unload',
          value: data.capacity - data.load,
        },
      ],
    };
    if (data.load === 0) {
      option.color = ['#C9D0D4'];
    }
    return option;
  };

  render() {
    const { overview = {} } = this.props;
    return (
      <div className="index-view-card work-load-page">
        <p className="title">工作负载</p>
        <div className="work-load-content">
          {overview.workload &&
            Object.keys(overview.workload).map((item: string) => {
              const load = overview.workload[item];
              const option = this.setOption(load);
              return (
                <div key={item}>
                  <Chart option={option} width={150} height={150} />
                  <p>{item}</p>
                </div>
              );
            })}
        </div>
      </div>
    );
  }
}
