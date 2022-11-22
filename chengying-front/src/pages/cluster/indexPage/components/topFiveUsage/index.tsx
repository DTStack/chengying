import * as React from 'react';
import { Row, Col, Empty } from 'antd';
import { cloneDeep } from 'lodash';
import Chart from '@/components/chart';
import { barOption } from '@/pages/cluster/indexPage/constant';
import './style.scss';

interface IProps {
  overview: any;
}
export default class TopFiveUsage extends React.PureComponent<IProps, any> {
  setOption = (data: any = [], type: string) => {
    const option: any = cloneDeep(barOption);
    option.dataset = {
      dimensions: ['ip', 'usage'],
      source: data,
    };
    option.series[0].name = `${type} Usage`;
    option.tooltip = {
      trigger: 'axis',
      formatter: (params) => {
        const value = params[0].value;
        return (
          '<div style="font-size:12px">' +
          '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:10px;height:10px;background-color:#3f87ff;"></span>' +
          '<span>' +
          value.ip +
          '</span></br>' +
          '<span>' +
          params[0].seriesName +
          ': ' +
          Number(parseFloat(value.usage).toFixed(2)) +
          '%' +
          '</span>' +
          '</div>'
        );
      },
    };
    return option;
  };

  render() {
    const { metrics = {} } = this.props.overview;
    return (
      <div className="top-five-usage-page">
        <div className="index-view-card mb-20">
          <p className="title">节点使用率TOP5</p>
        </div>
        <Row gutter={20}>
          <Col span={12}>
            <div className="index-view-card" style={{ height: 379 }}>
              <p className="title">CPU</p>
              {Array.isArray(metrics.cpu_top5) && metrics.cpu_top5.length ? (
                <Chart
                  option={this.setOption(metrics.cpu_top5, 'CPU')}
                  width="100%"
                  height="100%"
                />
              ) : (
                <Empty className="c-ant_empty-center" />
              )}
            </div>
          </Col>
          <Col span={12}>
            <div className="index-view-card" style={{ height: 379 }}>
              <p className="title">內存</p>
              {Array.isArray(metrics.mem_top5) && metrics.mem_top5.length ? (
                <Chart
                  option={this.setOption(metrics.mem_top5, 'Memory')}
                  width="100%"
                  height="100%"
                />
              ) : (
                <Empty className="c-ant_empty-center" />
              )}
            </div>
          </Col>
        </Row>
      </div>
    );
  }
}
