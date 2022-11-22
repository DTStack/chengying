export const option = {
  title: {
    text: 'ECharts 入门示例',
  },
  tooltip: {},
  xAxis: {
    data: ['衬衫', '羊毛衫', '雪纺衫', '裤子', '高跟鞋', '袜子'],
  },
  yAxis: {},
  series: [
    {
      name: '销量',
      type: 'bar',
      data: [5, 20, 36, 10, 10, 20],
    },
  ],
};

// 普通饼图
export const pieOption = {
  title: [
    {
      text: '{value|5}\n{name|总量 (个)}',
      // textAlign: "center",
      top: 'center',
      left: 'center',
      padding: [0, 10],
      textStyle: {
        rich: {
          value: {
            fontSize: 48,
            fontWeight: 600,
            color: '#333333',
            padding: [10, 0],
          },
          name: {
            fontSize: 12,
            color: '#666666',
          },
        },
      },
    },
  ],
  // tooltip: {
  //   trigger: 'item'
  // },
  legend: {
    icon: 'circle',
    orient: 'vertical',
    top: 'center',
    right: '0%',
    itemGap: 50,
    itemHeight: 9,
    itemWidth: 9,
    textStyle: {
      fontSize: 20,
      rich: {
        name: {
          width: 60,
          color: '#666',
        },
        value: {
          color: '#333',
          fontWeight: 600,
        },
      },
    },
    formatter: (name) => {
      return `{name|${name}}{value|1132}`;
    },
    data: ['正常', '异常'],
  },
  color: ['#12BC6A', '#FF5F5C'],
  series: [
    {
      type: 'pie',
      radius: ['55%', '70%'],
      // center: ['30%', '50%'],
      hoverAnimation: false,
      label: {
        show: false,
      },
      data: [
        {
          name: '正常',
          value: 5,
        },
        {
          name: '异常',
          value: 1,
        },
      ],
    },
  ],
};

// 柱狀圖
export const barOption = {
  tooltip: {
    trigger: 'axis',
  },
  grid: {
    top: '26px',
    left: '10px',
    right: '4px',
    bottom: '18px',
    containLabel: true,
  },
  dataset: {
    dimensions: ['date', 'value'],
    source: [],
  },
  xAxis: [
    {
      type: 'category',
      axisLine: {
        lineStyle: {
          color: '#e8e8e8',
        },
      },
      axisLabel: {
        color: 'rgba(0,0,0,0.45)',
      },
    },
  ],
  yAxis: [
    {
      type: 'value',
      axisLine: {
        show: false,
      },
      axisLabel: {
        color: 'rgba(0,0,0,0.45)',
        borderWidth: 0,
        formatter: '{value}%',
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
  ],
  color: ['#3F87FF'],
  series: [
    {
      name: 'CPU usage',
      type: 'bar',
      barWidth: '20',
    },
  ],
};
