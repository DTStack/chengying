import * as React from 'react';
const alarmStateMap = {
  ok: '#12BC6A',
  pending: '#F9DB13',
  no_data: '#FFA941',
  alerting: '#FF5F5C',
  paused: '#BFBFBF',
};
export const nodeStatusColumns = [
  {
    key: 'ip',
    dataIndex: 'ip',
    title: '节点',
  },
  {
    key: 'cpu',
    dataIndex: 'cpu',
    title: 'cpu',
    render: (value) => {
      const bg = value === 'NORMAL' ? '#12BC6A' : '#FF5F5C';
      const text = value === 'NORMAL' ? '正常' : '异常';
      return (
        <div className="circle-box">
          <div className="circle" style={{ background: bg }}></div>
          <div className="circle-text">{text}</div>
        </div>
      );
    },
  },
  {
    key: 'mem',
    dataIndex: 'mem',
    title: '内存',
    render: (value) => {
      const bg = value === 'NORMAL' ? '#12BC6A' : '#FF5F5C';
      const text = value === 'NORMAL' ? '正常' : '异常';
      return (
        <div className="circle-box">
          <div className="circle" style={{ background: bg }}></div>
          <div className="circle-text">{text}</div>
        </div>
      );
    },
  },
  {
    key: 'system_disk',
    dataIndex: 'system_disk',
    title: '系统盘',
    render: (value) => {
      const bg = value === 'NORMAL' ? '#12BC6A' : '#FF5F5C';
      const text = value === 'NORMAL' ? '正常' : '异常';
      return (
        <div className="circle-box">
          <div className="circle" style={{ background: bg }}></div>
          <div className="circle-text">{text}</div>
        </div>
      );
    },
  },
  {
    key: 'data_disk',
    dataIndex: 'data_disk',
    title: '数据盘',
    render: (value) => {
      const bg = value === 'NORMAL' ? '#12BC6A' : '#FF5F5C';
      const text = value === 'NORMAL' ? '正常' : '异常';
      return (
        <div className="circle-box">
          <div className="circle" style={{ background: bg }}></div>
          <div className="circle-text">{text}</div>
        </div>
      );
    },
  },
];

export const dtBaseColumns = [
  {
    key: 'service_name',
    dataIndex: 'service_name',
    title: '服务',
  },
  {
    key: 'ip',
    dataIndex: 'ip',
    title: '节点',
  },
  {
    key: 'ha_role',
    dataIndex: 'ha_role',
    title: '角色',
    render: (value) => {
      if (value) {
        return value;
      } else {
        return '--';
      }
    },
  },
  {
    key: 'status',
    dataIndex: 'status',
    title: '状态',
    render: (value) => {
      const bg = value === 'NORMAL' ? '#12BC6A' : '#FF5F5C';
      const text = value === 'NORMAL' ? '正常' : '异常';
      return (
        <div className="circle-box">
          <div className="circle" style={{ background: bg }}></div>
          <div className="circle-text">{text}</div>
        </div>
      );
    },
  },
];

export const dtBathColumns = [
  {
    key: 'service_name',
    dataIndex: 'service_name',
    title: '服务',
  },
  {
    key: 'ip',
    dataIndex: 'ip',
    title: '节点',
  },
  {
    key: 'status',
    dataIndex: 'status',
    title: '状态',
    render: (value) => {
      const bg = value === 'NORMAL' ? '#12BC6A' : '#FF5F5C';
      const text = value === 'NORMAL' ? '正常' : '异常';
      return (
        <div className="circle-box">
          <div className="circle" style={{ background: bg }}></div>
          <div className="circle-text">{text}</div>
        </div>
      );
    },
  },
];

export const alarmColumns = [
  {
    key: 'alert_name',
    dataIndex: 'alert_name',
    title: '告警名称',
  },
  {
    key: 'state',
    dataIndex: 'state',
    title: '状态',
    render: (value) => {
      return (
        <div className="circle-box">
          <div
            className="circle"
            style={{ background: alarmStateMap[value] }}></div>
          <div className="circle-text">{value}</div>
        </div>
      );
    },
  },
  {
    key: 'dashboard_name',
    dataIndex: 'dashboard_name',
    title: '仪表盘名称(组件)',
  },
  {
    key: 'dashboard_title',
    dataIndex: 'dashboard_title',
    title: '仪表盘标题',
  },
  {
    key: 'time',
    dataIndex: 'time',
    title: '告警时间',
  },
];
