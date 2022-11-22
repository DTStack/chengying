import React, { useState, useEffect } from 'react';
import { Modal, Table, Badge, Pagination, message } from 'antd';
import { scriptManager } from '@/services';
import './style.scss';
import { RESULT_STATUS, RESULT_FILTER_HISTORY } from '../const';
import { EllipsisText } from 'dt-react-component';

interface IProps {
  title: string;
  visible: boolean;
  close: () => void;
  id: number;
}

const renderStatus = (type: any) => {
  switch (type) {
    case RESULT_STATUS.NORMAL:
      return (
        <span>
          <Badge color="#12BC6A" /> 正常
        </span>
      );
    case RESULT_STATUS.UNNORMAL:
      return (
        <span>
          <Badge color="#FF5F5C" /> 异常
        </span>
      );
    case RESULT_STATUS.RUN:
      return (
        <span>
          <Badge color="#3F87FF" /> 运行中
        </span>
      );
    case RESULT_STATUS.UNRUN:
      return (
        <span>
          <Badge color="#BFBFBF" /> 未运行
        </span>
      );
  }
};

const TaskHistory: React.FC<IProps> = (props) => {
  const { title, visible, close, id } = props;
  const [data, setData] = useState([]);
  const [start, setStart] = useState(0);
  const [total, setTotal] = useState(0);
  const [execStatus, setExecStatus] = useState('');
  const [expandedKeys, setExpandedKeys] = useState([]);
  const columns: any = () => {
    return [
      {
        title: '编号',
        dataIndex: 'count',
        key: 'count',
        width: 80,
        fixed: 'left',
      },
      {
        title: '主机IP',
        dataIndex: 'ip',
        width: 130,
        key: 'ip',
        fixed: 'left',
        render: (text: string) => {
          return <span>{text || '--'}</span>;
        },
      },
      {
        title: '执行时间',
        dataIndex: 'end_time',
        key: 'end_time',
        width: 150,
        render: (text: string) => {
          return <span>{text || '--'}</span>;
        },
      },
      {
        title: '执行方式',
        dataIndex: 'exec_type',
        key: 'exec_type',
        width: 100,
        render: (text: number) => {
          if (text == 0) {
            return <span>定时执行</span>;
          } else {
            return <span>手动执行</span>;
          }
        },
      },
      {
        title: '执行结果',
        dataIndex: 'exec_status',
        key: 'exec_status',
        filters: RESULT_FILTER_HISTORY,
        filterMultiple: false,
        render: renderStatus,
        width: 100,
      },
      {
        title: '详情',
        dataIndex: 'exec_result',
        key: 'exec_result',
        render: (text: string) => {
          return (
            <div>
              {text?.length > 20 ? (
                <EllipsisText value={text} maxWidth={260} />
              ) : (
                <span>{text || '--'}</span>
              )}
            </div>
          );
        },
      },
    ];
  };

  const getLog = (
    pstart?: number,
    isFilter?: boolean,
    pexecStatus?: any,
    isOpen?: boolean
  ) => {
    const param = {
      id,
      limit: 10,
      start: pstart ?? start,
      execStatus: isFilter ? pexecStatus : execStatus,
    };
    scriptManager.getTaskLog(param).then((res: any) => {
      let arr = [];
      setData([]);
      const { data } = res;
      if (data.code == 0) {
        data.data.list.map((item: any, index: number) => {
          item.count = data.data.count - start - index;
          item.key = item.count;
          if (isOpen) {
            arr.push(item.key);
          }
          item.children?.map((citem: any, cIndex: number) => {
            citem.key = `${item.key}-${cIndex}`;
            return citem;
          });
          return item;
        });
        setExpandedKeys(arr);
        setData(data.data.list);
        setTotal(data.data.count);
      } else {
        message.error(data.data.msg);
      }
    });
  };
  // 表格筛选
  const handleChangeTable = (pagination, filters) => {
    let data = filters.exec_status;
    setExecStatus(data.length > 0 ? data[0] : '');
    if (data.length == 0 && execStatus) {
      setStart(0);
      getLog(0, true, '');
    }
  };

  const doExpand = (expanded, record: any) => {
    if (expanded) {
      setExpandedKeys([...expandedKeys, record.key]);
    } else {
      let arr = expandedKeys.filter((v) => {
        return v !== record.key;
      });
      setExpandedKeys(arr);
    }
  };

  useEffect(() => {
    if (!visible) {
      return;
    }
    // 调用接口
    if (execStatus) {
      getLog(start, true, execStatus, true);
      return;
    }
    getLog();
  }, [id, visible, start]);

  useEffect(() => {
    if (execStatus) {
      getLog(0, true, execStatus, true);
      return;
    }
  }, [execStatus]);

  // 改变页码数
  const onChangePage = (page: number) => {
    setStart((page - 1) * 10);
  };

  return (
    <Modal
      title={title}
      visible={visible}
      onCancel={close}
      width="900px"
      footer={null}>
      <Table
        onExpand={doExpand}
        expandedRowKeys={expandedKeys}
        pagination={false}
        columns={columns()}
        onChange={handleChangeTable}
        dataSource={data}
        scroll={{ x: 900 }}
      />
      <div
        className="paginationBox"
        style={{ padding: '13px 20px 13px 0', textAlign: 'right' }}>
        <Pagination
          current={start / 10 + 1}
          size="small"
          total={total}
          onChange={onChangePage}
          showTotal={(total) => (
            <span>
              共<span style={{ color: '#3F87FF' }}>{total}</span>
              条数据，每页显示
              <span style={{ color: '#3F87FF' }}>10</span>条
            </span>
          )}
        />
      </div>
    </Modal>
  );
};
export default TaskHistory;
