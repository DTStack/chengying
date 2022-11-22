/*
 * @Description: 脚本回显详情
 * @Author: wulin
 * @Date: 2021-06-01 14:00:41
 * @LastEditors: wulin
 * @LastEditTime: 2021-06-17 11:00:16
 */
import * as React from 'react';
import { PageHeader, Spin, Breadcrumb, message } from 'antd';
import * as Cookie from 'js-cookie';
import EchoTableDetail from './detail';
import { echoService } from '@/services';
import { statusMap, IRecordColumns } from '../index';

import '../style.scss';

type IDetailProps = {
  history: any;
  location: any;
  match: any;
};

type IColProps = {
  name: string;
  contentType: string;
  hostIp?: string;
  desc?: string;
  status?: string;
};

type IDetailTableProps = IRecordColumns & {
  children?: Array<IRecordColumns & IColProps>;
};

const CommandEchoDetail: React.FC<IDetailProps> = (props) => {
  const { location, history } = props;
  const { state = {} } = location;
  const [commandDetail, getCommandData] = React.useState<IDetailTableProps>({});
  let timeEvent;

  React.useEffect(() => {
    getEchoOrderDetail();
    return () => clearTimeout(timeEvent);
  }, []);

  /**
   * 获取列表数据
   */
  const getEchoOrderDetail = () => {
    echoService
      .getEchoOrderDetail({
        clusterId: Cookie.get('em_current_cluster_id'),
        operationId: state.operationId,
        // operationId: '7a179aaf-d675-4dac-a29c-89d01c15e634',
      })
      .then((res: any) => {
        const ret = res.data;
        if (ret.code === 0) {
          getCommandData(ret.data);
          timeEvent = setTimeout(() => {
            getEchoOrderDetail();
          }, 5000);
        } else {
          message.error(ret.msg);
        }
      });
  };

  const content: React.ReactNode = (
    <>
      <div className="echo-status clearfix mb-12">
        <span className="mr-48">
          状态：{(statusMap[(commandDetail || {}).status] || {}).label || '--'}
        </span>
        <span className="mr-48">
          对象：{(commandDetail || {}).productName || '--'}
        </span>
        <span className="mr-48">
          开始时间：{(commandDetail || {}).startTime || '--'}
        </span>
        <span className="mr-48">
          持续时间：{(commandDetail || {}).duration || '--'}s
        </span>
      </div>
      <EchoTableDetail
        records={commandDetail}
        history={history}
        states={state}
      />
    </>
  );

  return (
    <div className="scriptecho-page">
      <Breadcrumb>
        <Breadcrumb.Item
          onClick={(e): React.ReactEventHandler =>
            history.push(
              `/deploycenter/cluster/detail/echoList?from=${location.pathname}`
            )
          }>
          集群命令
        </Breadcrumb.Item>
        <Breadcrumb.Item>{(commandDetail || {}).name || '--'}</Breadcrumb.Item>
      </Breadcrumb>
      <PageHeader
        title={
          <span>
            {(commandDetail || {}).status === 1 && <Spin size="small" />}&nbsp;
            {(commandDetail || {}).name}
          </span>
        }
        footer={content}
      />
    </div>
  );
};

export default CommandEchoDetail;
