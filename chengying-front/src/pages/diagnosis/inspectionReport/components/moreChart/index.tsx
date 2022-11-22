import * as React from 'react';
import Chart from '@/components/chart';
import { colorRange } from './constant';
import moment from 'moment';
import { message, Spin, Icon } from 'antd';
import { inspectionReportService } from '@/services';
import './style.scss';
interface Prop {
  config: any;
  time: any;
  index: number;
  delNoDataChart: Function;
}
export default function index(props: Prop) {
  const { config, time, index, delNoDataChart } = props;
  const [option, setOption] = React.useState<any>({});
  const [chartData, setChartData] = React.useState<any>({});
  const [dataLoading, setDataLoading] = React.useState(false);

  const getChartData = async (params: any) => {
    setDataLoading(true);
    const { targets, decimal, unit, metric } = params;
    const obj = {
      targets,
      decimal,
      unit,
      from: moment(`${time[0].format('YYYY-MM-DD')} 00:00:00`).valueOf(),
      to: moment(`${time[1].format('YYYY-MM-DD')} 23:59:59`).valueOf(),
    };
    const res: any = await inspectionReportService.getChartData(obj);
    const { code, msg, data } = res.data;
    if (data?.x?.length === 0) {
      delNoDataChart(index);
      return;
    }
    if (code === 0) {
      const obj = {
        ...data,
        title: metric,
      };
      setOption(handleOption(obj));
      setChartData(obj);
    } else {
      message.error(msg);
    }
    setDataLoading(false);
  };

  
  const handleOption = (data: any, clickIndex?: number) => {
    const option = {
      grid: {
        left: 20,
        right: 20,
        top: 56,
        bottom: 12,
        containLabel: true,
      },
      title: {
        text: data.title || '',
        left: '50%',
        top: 12,
        textStyle: {
          color: '#333333',
          fontFamily: 'PingFangSC-Regular, PingFang SC',
          fontSize: 14,
          fontWeight: 400,
        },
      },
      legend: {
        top: "87%",
        left: "6%",
        itemWidth: 23,
        itemHeight: 5
      },
      xAxis: {
        type: 'category',
        data: (data.x || []).map((item) =>
          moment(item * 1000).format('MM/DD HH:mm')
        ),
        axisLine: {
          lineStyle: {
            color: '#BFBFBF',
          },
        },
      },
      yAxis: {
        type: 'value',
        axisLine: {
          show: false,
          lineStyle: {
            color: '#BFBFBF',
          },
        },
        splitLine: {
          lineStyle: {
            type: 'dashed',
          },
        },
        axisLabel: {
          margin: 14,
        },
      },
      tooltip: {
        show: true,
        trigger: 'axis',
        textStyle: {
          fontSize: 12,
          lineHeight: 20,
          padding: 12,
        },
      },
      series: [],
    };
    const series = [];
    const len = colorRange.length;
    (data.y || []).forEach((item, index) => {
      series.push({
        data: item.data,
        type: 'line',
        smooth: true,
        lineStyle: {
          color: colorRange[index % len],
        },
      });
    });
    option.series = [...series];
    return { ...option };
  };
  React.useEffect(() => {
    getChartData(config);
  }, [config]);
  return (
    <div className="more-chart">
      {dataLoading ? (
        <div className="more-chart-mask">
          <Spin
            tip="数据加载中"
            indicator={<Icon type="loading" style={{ fontSize: 24 }} spin />}
          />
        </div>
      ) : null}
      <div className="more-chart-item">
        <Chart option={option} width="100%" height="100%"></Chart>
      </div>
      <div className="more-chart-legend">
        <div className="more-chart-legend-item" style={{ marginBottom: 12 }}>
          <span></span>
          <span style={{ color: '#3F87FF' }}>current</span>
        </div>
        {(chartData.y || []).map((item, index) => {
          const len = colorRange.length;
          const bg = colorRange[index % len];
          return (
            <div className="more-chart-legend-item" key={index}>
              <div className="legend-item-left">
                <div className="legend-icon" style={{ background: bg }}></div>
                <span className="legend-title">{item.title}</span>
              </div>
              <span>{item.data[0]}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
