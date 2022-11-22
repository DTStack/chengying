import * as React from 'react';
import Resize from '@/components/resize';

// 引入 ECharts 主模块
const echarts = require('echarts/lib/echarts');
// 引入柱状图
require('echarts/lib/chart/line');
// 引入提示框和标题组件
require('echarts/lib/component/legend');
require('echarts/lib/component/tooltip');
require('echarts/lib/component/title');

interface Props {
  echartsTime: any[];
  echartsData: any[];
  currentEventType: string;
}
class Echarts extends React.Component<Props, any> {
  constructor(props: any) {
    super(props);
  }

  componentDidMount() {
    this.initLineChart();
  }

  componentDidUpdate(nextProps) {
    if (nextProps.echartsTime != this.props.echartsTime) {
      this.initLineChart();
    }
  }

  initLineChart() {
    const { echartsTime, echartsData } = this.props;
    const myChart = echarts.init(document.getElementById('echarts'));
    const option = {
      color: '#20aefc',
      tooltip: {
        trigger: 'axis',
      },
      grid: {
        left: '45px', // 距离左边的距离
        right: '20px', // 距离右边的距离
        bottom: '30px', // 距离下边的距离
        top: '42px', // 距离上边的距离
      },
      calculable: true,
      xAxis: [
        {
          type: 'category',
          data: echartsTime,
          axisLine: {
            lineStyle: {
              color: '#CECECE',
            },
          },
          axisTick: {
            // 坐标轴刻度相关设置
            length: '0', // 长度设置为0
          },
          axisLabel: {
            show: true,
            margin: 12,
            color: '#7E7E7E',
          },
          boundaryGap: true,
        },
      ],
      yAxis: [
        {
          type: 'value',
          name: '数量           ',
          axisTick: {
            // 坐标轴的刻度
            show: false,
          },
          axisLine: {
            show: false,
          },
          axisLabel: {
            show: true,
            color: '#7E7E7E',
          },
          splitLine: {
            // 坐标轴分割线。默认数值轴显示，类目轴不显示
            show: true,
            lineStyle: {
              color: '#CECECE',
              width: 1,
              type: 'dashed',
            },
          },
        },
      ],
      series: [
        // 系列列表
        {
          name: this.props.currentEventType + '次数',
          type: 'line',
          symbol: 'none',
          data: echartsData,
          lineStyle: {
            // 线条的样式
            normal: {
              color: '#20aefc', // 折线的颜色
            },
          },
          smooth: 0.3, // 是否平滑处理，如果是Number类型(取值范围为0到1)，表示平滑程度，越小越接近折线段，反之
        },
      ],
    };
    // 绘制图表
    myChart.setOption(option);
    this.setState({ lineChart: myChart });
  }

  onResize = () => {
    this.state.lineChart.resize();
  };

  render() {
    return (
      <Resize onResize={this.onResize}>
        <div id="echarts" style={{ width: '100%', height: '100%' }}></div>
      </Resize>
    );
  }
}
export default Echarts;
