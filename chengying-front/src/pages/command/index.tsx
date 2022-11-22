/*
 * @Description: 命令列表
 * @Author: wulin
 * @Date: 2021-06-01 14:00:41
 * @LastEditors: wulin
 * @LastEditTime: 2021-06-25 13:57:31
 */
import * as React from 'react';
import { Table, Badge, Spin, message } from 'antd';
import * as Cookie from 'js-cookie';
import { PaginationProps } from 'antd/es/pagination/Pagination';
import { echoService } from '@/services';
import EchoSearch from './echoSearch';
import TimerInc from '@/components/timerInc';
import Utils from '@/utils/utils';
import './style.scss';

type IEchoProps = {
  history: any;
  location: any;
  match: any;
};

type IReqParams = {
  clusterId: string;
  operationType: number;
  objectValue: string;
  status: string;
  startTime: string;
  endTime: string;
  page: number;
  pageSize: number;
};

export type IRecordColumns = { [key: string]: string | number } & {
  isExist?: boolean;
};

export const statusMap = {
  '2': {
    label: '正常',
    value: 'success',
  },
  '3': {
    label: '失败',
    value: 'error',
  },
  '1': {
    label: '进行中',
    value: 'processing',
  },
};

export const operationMaps: IRecordColumns = {
  '1': '产品包部署',
  '2': '补丁包升级',
  '3': '产品包启动',
  '4': '服务启动',
  '5': '服务滚动重启',
  '6': '主机初始化',
  '7': 'Kerberos开启',
  '8': 'Kerberos关闭',
};

const CommandEchoList: React.FC<IEchoProps> = (props) => {
  const { history, location } = props;
  const [reqParams, changeReq] = React.useState<IReqParams>({
    clusterId: Cookie.get('em_current_cluster_id'),
    operationType: undefined,
    objectValue: '',
    status: '',
    startTime: '',
    endTime: '',
    page: 1,
    pageSize: 10,
  });
  const [echoData, getEchoData] = React.useState<Array<IRecordColumns>>([]);
  const [total, setTotal] = React.useState<number>(0);
  const [polling, changePolling] = React.useState<boolean>(false);
  const [loading, changeStatusLoading] = React.useState<boolean>(false);

  // 分页配置
  const pagination: PaginationProps = {
    current: reqParams.page,
    pageSize: reqParams.pageSize,
    total: total,
    size: 'small',
    showTotal: (total) =>
      `共 ${total} 条数据，每页显示 ${reqParams.pageSize} 条`,
  };
  let timeEvent = React.useRef(null);

  React.useEffect(() => {
    getEchoOrderList();
    return () => {
      changePolling(true);
      clearTimeout(timeEvent.current);
    };
  }, [reqParams]);

  /**
   * 获取列表数据
   */
  const getEchoOrderList = () => {
    const params = {
      ...reqParams,
      startTime: reqParams.startTime || undefined,
      endTime: reqParams.endTime || undefined,
      status: reqParams.status || undefined,
      objectValue: reqParams.objectValue || undefined,
    };
    changeStatusLoading(true);
    echoService.getEchoOrder(params).then((res: any) => {
      const ret = res.data;
      if (ret.code === 0) {
        getEchoData(ret.data.list);
        setTotal(ret.data.count);
        changeStatusLoading(false);
        if (polling) return null;
        timeEvent.current = setTimeout(() => {
          getEchoOrderList();
        }, 5000);
      } else {
        message.error(ret.msg);
      }
    });
  };

  /**
   * table的change函数
   * @param pagination
   * @param filters
   * @param sorter
   */
  const handleTableChange = (
    pagination: { current: any },
    filters: any,
    sorter: any
  ): void => {
    changeReq((prevValue) => ({ ...prevValue, page: pagination.current }));
  };

  const handleGoToDetail = (record: IRecordColumns): void => {
    history.push(
      `/deploycenter/cluster/detail/echoDetail?from=${location.pathname}`,
      record
    );
  };

  const handleGoToResult = (record: any): any => {
    message.destroy();
    if (!record.isExist) return message.warning('该组件已从该集群卸载！');
    // 1 产品包、2 服务、3 主机
    let path: string = '';
    switch (record.objectType) {
      case 1:
        path = `/deploycenter/cluster/detail/deployed`;
        Utils.setNaviKey('menu_deploy_center', 'sub_menu_cluster_list');
        break;
      case 2:
        path = `/opscenter/service`;
        Cookie.set('em_product_id', '');
        Cookie.set('em_product_name', record.productName);
        Cookie.set('em_current_parent_product', record.parentProductName);
        Utils.setNaviKey('menu_ops_center', 'sub_menu_service');
        break;
      default:
        path = `/opscenter/hoststatus`;
        // 兼容hostIp
        record = { ...record, hostIp: record.objectValue };
        Cookie.set('em_product_id', '');
        Cookie.set('em_product_name', record.productName);
        Cookie.set('em_current_parent_product', record.parentProductName);
        Utils.setNaviKey('menu_ops_center', 'sub_menu_host');
        break;
    }
    sessionStorage.setItem('service_object', JSON.stringify(record));
    history.push(path);
  };

  return (
    <div className="scriptecho-page">
      <EchoSearch handleEvent={changeReq} />
      <Table
        rowKey="operationId"
        className="dt-table-fixed-base"
        style={{ height: 'calc(100vh - 230px)' }}
        dataSource={echoData}
        onChange={handleTableChange}
        pagination={pagination}
        loading={loading}
        scroll={{ y: true }}>
        <Table.Column
          title=""
          width={30}
          key="operationId"
          render={(v, record: any): React.ReactNode => {
            return (
              <div>
                {record.operationStatus === 1 ? (
                  <Spin size="small" />
                ) : (
                  <span>&nbsp;</span>
                )}
              </div>
            );
          }}
        />
        <Table.Column
          title="对象"
          width={'20%'}
          dataIndex="objectValue"
          render={(text: string, record: IRecordColumns): React.ReactNode => {
            return <a onClick={(e) => handleGoToResult(record)}>{text}</a>;
          }}
        />
        <Table.Column
          title="操作"
          dataIndex="operationName"
          render={(text: string, record: IRecordColumns): React.ReactNode => {
            return <a onClick={(e) => handleGoToDetail(record)}>{text}</a>;
          }}
        />
        <Table.Column
          title="状态"
          dataIndex="operationStatus"
          render={(text: string): React.ReactNode => {
            const currentStatus = statusMap[text];
            return (
              <span>
                <Badge status={currentStatus.value} />
                &nbsp;{currentStatus.label}
              </span>
            );
          }}
        />
        <Table.Column title="开始时间" dataIndex="startTime" />
        <Table.Column
          title="持续时间"
          width={'12%'}
          dataIndex="duration"
          render={(time, record: any) => {
            return record.operationStatus === 1 ? (
              <TimerInc initial={parseFloat(time)} />
            ) : (
              time + 's'
            );
          }}
        />
      </Table>
    </div>
  );
};

export default CommandEchoList;
