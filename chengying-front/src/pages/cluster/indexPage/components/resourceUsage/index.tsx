import * as React from 'react';
import { Button, Icon, message, Empty } from 'antd';
import classnames from 'classnames';
import { cloneDeep } from 'lodash';
import { pieOption, lineOption, timeRange } from './constant';
import Chart from '@/components/chart';
import './style.scss';
import clusterIndexService from '@/services/clusterIndexService';
import moment from 'moment';

interface IProps {
  overview: any;
  cur_parent_cluster: any;
}
interface IState {
  metric: string;
  performanceOption: any;
  timePicker: boolean;
  selectTime: any;
}

export default class ResourceUsage extends React.PureComponent<IProps, IState> {
  state: IState = {
    metric: 'cpu',
    performanceOption: {},
    timePicker: false,
    selectTime: {
      value: '15m',
      text: 'Last 15 minutes',
      times: [moment().subtract('15', 'minutes'), moment()],
    },
  };

  componentDidMount() {
    this.getPerformance();
  }

  componentDidUpdate(prevProps: IProps) {
    if (prevProps.cur_parent_cluster.id !== this.props.cur_parent_cluster.id) {
      this.getPerformance();
    }
  }

  handleClick = (metric: string) => {
    this.setState(
      {
        metric,
        performanceOption: metric !== 'pod' ? this.state.performanceOption : {},
      },
      () => {
        metric !== 'pod' && this.getPerformance();
      }
    );
  };

  // 获取使用率趋势图
  getPerformance = () => {
    const { cur_parent_cluster } = this.props;
    const {
      metric,
      selectTime: { times },
    } = this.state;
    const metricMap = {
      cpu: 'CPU',
      memory: '内存',
      disk: '磁盘',
      pod: '容器量',
    };
    clusterIndexService
      .getClusterPerformance(
        cur_parent_cluster.id,
        {
          metric: metric,
          from: Math.floor(times[0].valueOf() / 1000),
          to: Math.floor(times[1].valueOf() / 1000),
        },
        cur_parent_cluster.type
      )
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          if (res.data.lists.length) {
            const performanceOption = cloneDeep(lineOption);
            performanceOption.yAxis.name = `${metricMap[metric]}使用情况（%）`;
            performanceOption.dataset = {
              dimensions: ['date', 'value'],
              source: res.data.lists,
            };
            this.setState({
              performanceOption,
            });
          }
        } else {
          message.error(res.msg);
        }
      });
  };

  // 时间控件展开
  handleTimePickerShow = () => {
    this.setState({
      timePicker: !this.state.timePicker,
    });
  };

  // 时间选择
  handleTimeChange = (item: any) => {
    this.setState(
      {
        selectTime: item,
        timePicker: false,
      },
      () => {
        this.getPerformance();
      }
    );
  };

  render() {
    const { overview, cur_parent_cluster } = this.props;
    const { metric, performanceOption } = this.state;
    const cardProps = {
      metric,
      handleClick: this.handleClick,
    };
    return (
      <div className="resource-usage-page">
        <div className="pg-l">
          <ResourceCard
            {...cardProps}
            type="cpu"
            title="CPU"
            usage={overview.cpu_used_display}
            size={overview.cpu_size_display}
          />
          <ResourceCard
            {...cardProps}
            type="memory"
            title="內存"
            usage={overview.mem_used_display}
            size={overview.mem_size_display}
          />
          {cur_parent_cluster.type === 'hosts' ? (
            <ResourceCard
              {...cardProps}
              type="disk"
              title="磁盘"
              usage={overview.disk_used_display}
              size={overview.disk_size_display}
            />
          ) : (
            <ResourceCard
              {...cardProps}
              type="pod"
              title="容量"
              usage={overview.pod_used_display}
              size={overview.pod_size_display}
            />
          )}
        </div>
        <div className="pg-r">
          <div className="component-timer-absolute">
            <Button type="default" onClick={this.handleTimePickerShow}>
              <Icon type="clock-circle" /> {this.state.selectTime.text}
            </Button>
            {this.state.timePicker && (
              <div className="timer-picker box-shadow-style">
                {timeRange.map((list, index) => (
                  <ul key={index}>
                    {list.map((item) => (
                      <li
                        key={item.value}
                        className={classnames({
                          active: this.state.selectTime.value === item.value,
                        })}
                        onClick={this.handleTimeChange.bind(this, item)}>
                        {item.text}
                      </li>
                    ))}
                  </ul>
                ))}
              </div>
            )}
          </div>
          {Object.keys(performanceOption).length ? (
            <Chart option={performanceOption} height="100%" />
          ) : (
            <Empty className="c-ant_empty-center" />
          )}
        </div>
      </div>
    );
  }
}

interface ResourceCardProps {
  handleClick: Function;
  metric: string;
  type: string;
  title: string;
  usage: string;
  size: string;
}
const ResourceCard = (props: ResourceCardProps) => {
  function formatNum(value: string) {
    const data = parseFloat(value);
    if (data === undefined || isNaN(data)) {
      return 0;
    }
    return data;
  }

  function setOption(usage: number, size: number) {
    const option = cloneDeep(pieOption);
    const value = size === 0 ? 0 : Number(((usage / size) * 100).toFixed(2));
    option.series[1].data = [
      {
        ...option.series[1].data[0],
        itemStyle: {
          color:
            value >= window.APPCONFIG.ALARM_WARN
              ? '#FF5F5C'
              : value >= window.APPCONFIG.ALARM_HEALTH
              ? '#FFB310'
              : '#12BC6A',
        },
        value,
      },
      {
        ...option.series[1].data[1],
        value: 150 - value,
      },
    ];
    return option;
  }

  const { type, title } = props;
  const usage = formatNum(props.usage);
  const size = formatNum(props.size);
  const option = setOption(usage, size);
  return (
    <div
      className={classnames('resource-usage-card', {
        'card-focus': props.metric === type,
      })}
      onClick={() => props.handleClick(type)}>
      <div className="card-content">
        <Chart option={option} width={110} height={110} />
        <div className="card-info">
          <p>
            {title}
            {props.size && `（${props.size.replace(/[\d.]+/g, '')}）`}
          </p>
          <p className="mt-10">
            <span className="info-usage">{usage}</span> / {size}
          </p>
        </div>
      </div>
    </div>
  );
};
