import * as React from 'react';
import { message, Table } from 'antd';
import { ClusterNamespaceService } from '@/services';
import { PaginationConfig } from 'antd/lib/table';
import './style.scss';

interface IProps {
  namespace: string;
}
const EventView: React.FC<IProps> = (props) => {
  const [reqParams, setReqParams] = React.useState<any>({
    start: 0,
    limit: 10,
  });
  const [eventList, setEventLists] = React.useState<any>({
    list: [],
    count: 0,
  });
  const [tableLoading, setTableLoading] = React.useState<boolean>(false);
  const columns = [
    {
      title: '发生时间',
      dataIndex: 'time',
      key: 'time',
      width: '20%',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: '15%',
    },
    {
      title: '原因',
      dataIndex: 'reason',
      key: 'reason',
      width: '15%',
    },
    {
      title: '对象',
      dataIndex: 'resource',
      key: 'resource',
      width: '20%',
    },
    {
      title: '信息',
      dataIndex: 'message',
      key: 'message',
      width: '20%',
      render: (text: string) =>
        text ? <span style={{ wordBreak: 'break-all' }}>{text}</span> : '--',
    },
  ];
  const pagination = {
    size: 'small',
    pageSize: reqParams.limit,
    total: eventList.count,
    current: reqParams.start / reqParams.limit + 1,
    showTotal: (total) => (
      <span>
        共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示
        {reqParams.limit}条
      </span>
    ),
  };

  React.useEffect(() => {
    getEventList();
  }, [reqParams]);

  // 获取事件列表
  function getEventList() {
    // setTableLoading(true);
    ClusterNamespaceService.getEventLists({
      ...reqParams,
      namespace: props.namespace,
    }).then((response: any) => {
      const res = response.data;
      const { code, data, msg } = res;
      if (code === 0) {
        setEventLists({
          list: data.events || [],
          count: data.size || 0,
        });
      } else {
        message.error(msg);
      }
      setTableLoading(false);
    });
  }

  // 表格分页
  function handleTableChange(pagination: PaginationConfig) {
    const { current, pageSize } = pagination;
    setReqParams({
      start: (current - 1) * pageSize,
      limit: pageSize,
    });
  }

  return (
    <Table
      size="middle"
      rowKey="id"
      className="dt-pagination-lower dt-table-border dt-table-last-row-noborder"
      scroll={{ y: 'calc(100vh - 90px - 32px - 89px - 88px)' }}
      loading={tableLoading}
      columns={columns}
      dataSource={eventList.list}
      pagination={pagination}
      onChange={handleTableChange}
    />
  );
};
export default EventView;
