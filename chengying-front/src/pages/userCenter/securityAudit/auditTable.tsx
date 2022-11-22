import * as React from 'react';
import { Table, message } from 'antd';
import { isEqual } from 'lodash';
import { SecurityAuditService } from '@/services';
import { AuditReqParamsType } from './interface';

interface IProps {
  reqParams: AuditReqParamsType;
}

interface IState {
  pageParams: {
    start: number;
    limit: number;
    'sort-by': string;
    'sort-dir': string;
  };
  auditList: {
    list: any[];
    count: number;
  };
  loading: boolean;
}

export default class AuditTable extends React.PureComponent<IProps, IState> {
  state: IState = {
    pageParams: {
      start: 0,
      limit: 10,
      'sort-by': 'create_time',
      'sort-dir': 'desc',
    },
    auditList: {
      list: [],
      count: 0,
    },
    loading: false,
  };

  componentDidMount() {
    this.getAuditList();
  }

  componentDidUpdate(nextProps: IProps) {
    if (!isEqual(nextProps.reqParams, this.props.reqParams)) {
      this.setState(
        {
          pageParams: {
            ...this.state.pageParams,
            start: 0,
          },
        },
        this.getAuditList
      );
    }
  }

  // 获取表格数据
  getAuditList = () => {
    const { pageParams } = this.state;
    const { reqParams } = this.props;
    this.setState({ loading: true });
    SecurityAuditService.getSafetyAudit({
      ...reqParams,
      ...pageParams,
    }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({
          auditList: {
            count: res.data ? res.data.count : 0,
            list: res.data ? res.data.list : [],
          },
        });
      } else {
        message.error(res.msg);
      }
      this.setState({ loading: false });
    });
  };

  // 表格分页
  handleTableChange = (pagination: any) => {
    const pageParams = {
      ...this.state.pageParams,
      start: (pagination.current - 1) * pagination.pageSize,
    };
    this.setState({ pageParams }, this.getAuditList);
  };

  initColumns = () => {
    const columns = [
      {
        title: '时间',
        dataIndex: 'create_time',
        key: 'create_time',
        render: (text: string) => text || '--',
      },
      {
        title: '操作人',
        dataIndex: 'operator',
        key: 'operator',
        render: (text: string) => text || '--',
      },
      {
        title: '操作模块',
        dataIndex: 'module',
        key: 'module',
      },
      {
        title: '动作',
        dataIndex: 'operation',
        key: 'operation',
      },
      {
        title: '详细内容',
        dataIndex: 'content',
        key: 'content',
        width: '25%',
        render: (text: string) => text || '--',
      },
      {
        title: '来源IP',
        dataIndex: 'ip',
        key: 'ip',
        render: (text: string) => text || '--',
      },
    ];
    return columns;
  };

  render() {
    const { auditList, pageParams, loading } = this.state;
    const columns = this.initColumns();
    const pagination = {
      size: 'small',
      pageSize: pageParams.limit,
      total: auditList.count,
      current: pageParams.start / pageParams.limit + 1,
      showTotal: (total) => (
        <span>
          共<span style={{ color: '#3F87FF' }}>{total}</span>条数据，每页显示
          {pageParams.limit}条
        </span>
      ),
    };
    return (
      <Table
        rowKey="id"
        className="dt-table-fixed-base"
        style={{ height: 'calc(100vh - 135px)' }}
        scroll={{ y: true }}
        loading={loading}
        columns={columns}
        dataSource={auditList.list}
        pagination={pagination}
        onChange={this.handleTableChange}
      />
    );
  }
}
