import * as echarts from 'echarts';
import moment from 'moment';
export const pieOption = {
  // 第一个图表
  series: [
    {
      type: 'pie',
      hoverAnimation: false, // 鼠标经过的特效
      radius: ['72%', '80%'],
      startAngle: 210,
      labelLine: {
        show: false,
      },
      label: {
        position: 'center',
      },
      data: [
        {
          value: window.APPCONFIG.ALARM_HEALTH,
          itemStyle: {
            color: '#12BC6A',
          },
        },
        {
          value: window.APPCONFIG.ALARM_WARN - window.APPCONFIG.ALARM_HEALTH,
          itemStyle: {
            color: '#FFB310',
          },
        },
        {
          value: 100 - window.APPCONFIG.ALARM_WARN,
          itemStyle: {
            color: '#FF5F5C',
          },
        },
        {
          value: 50,
          itemStyle: {
            color: 'rgba(0,0,0,0)',
            borderWidth: 0,
          },
        },
      ],
    },
    // 上层环形配置
    {
      type: 'pie',
      hoverAnimation: false, // 鼠标经过的特效
      radius: ['52%', '70%'],
      startAngle: 210,
      labelLine: {
        show: false,
      },
      label: {
        position: 'center',
      },
      data: [
        {
          value: 75,
          itemStyle: {
            color: '#FF5F5C',
          },
          label: {
            formatter: '{c}%',
            position: 'center',
            show: true,
            textStyle: {
              fontSize: '12',
              fontWeight: '600',
              color: '#333333',
            },
          },
        },
        {
          value: 75,
          itemStyle: {
            color: 'rgba(0,0,0,0)',
            borderWidth: 0,
          },
        },
      ],
    },
  ],
};

export const lineOption = {
  color: '#3f87ff',
  tooltip: {
    trigger: 'axis',
    formatter: (params) => {
      const value = params[0].value;
      return (
        '<div style="font-size:12px">' +
        '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:10px;height:10px;background-color:#3f87ff;"></span>' +
        '<span>' +
        value.date +
        '</span></br>' +
        '<span>' +
        params[0].seriesName +
        ': ' +
        Number(value.value.toFixed(2)) +
        '%' +
        '</span>' +
        '</div>'
      );
    },
  },
  grid: {
    top: '60px',
    left: '36px',
    right: '24px',
    bottom: '22px',
    containLabel: true,
  },
  dataset: {
    dimensions: ['date', 'value'],
    source: [],
  },
  xAxis: {
    type: 'category',
    min: (value) => value.min - 1,
    max: (value) => value.max + 1,
    axisLine: {
      lineStyle: {
        color: '#e8e8e8',
      },
    },
    axisLabel: {
      color: 'rgba(0,0,0,0.45)',
      margin: 12,
    },
  },
  yAxis: {
    name: 'CPU使用情况（%）',
    nameTextStyle: {
      fontSize: 14,
      padding: [0, 0, 12, 50],
    },
    axisLine: {
      show: false,
    },
    axisLabel: {
      color: 'rgba(0,0,0,0.45)',
      borderWidth: 0,
    },
    axisTick: {
      show: false,
    },
    splitLine: {
      lineStyle: {
        type: 'dashed',
      },
    },
  },
  series: [
    {
      name: '使用率',
      type: 'line',
      lineStyle: {
        color: '#3f87ff',
      },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(
          0,
          0,
          0,
          1,
          [
            { offset: 0, color: 'rgba(63, 135, 255, 1)' },
            { offset: 1, color: 'rgba(115, 188, 255, 0)' },
          ],
          false
        ),
      },
    },
  ],
};

export const timeRange = [
  [
    {
      value: '2d',
      text: 'Last 2 days',
      times: [moment().subtract('2', 'days'), moment()],
    },
    {
      value: '7d',
      text: 'Last 7 days',
      times: [moment().subtract('7', 'days'), moment()],
    },
  ],
  [
    {
      value: 'today',
      text: 'Today',
      times: [moment().startOf('day'), moment()],
    },
  ],
  [
    {
      value: '15m',
      text: 'Last 15 minutes',
      times: [moment().subtract('15', 'minutes'), moment()],
    },
    {
      value: '30m',
      text: 'Last 30 minutes',
      times: [moment().subtract('30', 'minutes'), moment()],
    },
    {
      value: '1h',
      text: 'Last 1 hour',
      times: [moment().subtract('1', 'hour'), moment()],
    },
    {
      value: '3h',
      text: 'Last 3 hours',
      times: [moment().subtract('3', 'hours'), moment()],
    },
    {
      value: '6h',
      text: 'Last 6 hours',
      times: [moment().subtract('6', 'hours'), moment()],
    },
    {
      value: '12h',
      text: 'Last 12 hours',
      times: [moment().subtract('12', 'hours'), moment()],
    },
    {
      value: '24h',
      text: 'Last 24 hours',
      times: [moment().subtract('24', 'hours'), moment()],
    },
  ],
];
