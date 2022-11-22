/*
 * @Description: This is desc
 * @Author: wulin
 * @Date: 2021-05-10 11:26:40
 * @LastEditors: wulin
 * @LastEditTime: 2021-06-24 11:52:37
 */
import * as React from 'react';
import { Descriptions } from 'antd';
import moment from 'moment';
import { clusterTypeMap } from '@/constants/const';
import './style.scss';

interface IProps {
  overview: any;
  cur_parent_cluster: any;
}

export default class InfoView extends React.PureComponent<IProps, any> {
  render() {
    const { cur_parent_cluster = {}, overview = {} } = this.props;
    const { type = 'hosts', mode } = cur_parent_cluster;
    return (
      <div className="index-view-card">
        <p className="title">集群信息</p>
        <Descriptions className="c-cluster__description" column={1} bordered>
          <Descriptions.Item label="集群模式">
            {clusterTypeMap[type] ? clusterTypeMap[type][mode] : ''}
          </Descriptions.Item>
          {type === 'kubernetes' && (
            <Descriptions.Item label="kubernetes版本">
              {overview.version || '--'}
            </Descriptions.Item>
          )}
          <Descriptions.Item label="节点数">{overview.nodes}</Descriptions.Item>
          <Descriptions.Item label="CPU">
            {overview.cpu_size_display}
          </Descriptions.Item>
          <Descriptions.Item label="内存">
            {overview.mem_size_display}
          </Descriptions.Item>
          {type === 'kubernetes' ? (
            <Descriptions.Item label="PODS">
              {overview.pod_size_display}
            </Descriptions.Item>
          ) : (
            <Descriptions.Item label="磁盘">
              {overview.disk_size_display}
            </Descriptions.Item>
          )}
          <Descriptions.Item label="创建人">
            {overview.create_user}
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {moment(overview.create_time).format('YYYY-MM-DD HH:mm:ss')}
          </Descriptions.Item>
        </Descriptions>
      </div>
    );
  }
}
