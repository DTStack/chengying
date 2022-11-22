/*
 * @Description: 脚本回显详情
 * @Author: wulin
 * @Date: 2021-06-01 14:00:41
 * @LastEditors: wulin
 * @LastEditTime: 2021-06-25 09:55:47
 */
import * as React from 'react';
import { connect } from 'react-redux'
import { Table, Icon, Spin, message } from 'antd';
import CommandPopbox from './popbox';
import * as Cookie from 'js-cookie';
import '../style.scss';
import Utils from '@/utils/utils';
import TimerInc from '@/components/timerInc';
import { AppStoreTypes } from '@/stores';

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});

type IRecordProps = {
  records: any;
  history: any;
  states: any;
};

// const objectTypeMap = new Map<number, string>([
//   [1, '产品包'],
//   [2, '服务'],
//   [3, '主机'],
// ])

type ITableColumn = {
  operationName?: string;
  name?: string;
  startTime: string;
  duration: number;
  execId?: string;
  shellType?: number;
  objectType?: number;
  objectValue?: string;
  hostIp?: string;
  desc?: string;
  status?: number;
};

export type IPopboxState = {
  visible: boolean;
  title: string;
  type: string;
  execId: string;
};

const EchoTableDetail: React.FC<IRecordProps> = (props) => {
  const { records, history } = props;
  const [expandedRow, setExpand] = React.useState<string[]>([]);
  // const [data, getCommandData] = React.useState<Array<ITableColumn>>(record.children);
  const [visibleInfo, changePopboxInfo] = React.useState<IPopboxState>({
    visible: false,
    title: '',
    type: '',
    execId: '',
  });

  /**
   * 展开收起
   * @param expanded
   * @param record
   */
  const handleExpand = (expanded, record): void => {
    let rowKeys = [...expandedRow];
    const currentRow = tableKey(record);
    if (expanded && !rowKeys.includes(currentRow)) {
      rowKeys.push(currentRow);
    } else {
      rowKeys = expandedRow.filter((o) => o != currentRow);
    }
    setExpand(rowKeys);
  };

  /**
   * 弹窗
   * @param record
   * @param key
   */
  const handleEvent = (record?: ITableColumn, key?: string): void => {
    changePopboxInfo((prevValue) => ({
      visible: !prevValue.visible,
      title: key === 'echo' ? '脚本查看' : '日志查看',
      type: key,
      execId: record && record.execId,
    }));
  };

  /**
   * 自定义icon
   * @param props
   * @returns
   */
  const custemIcon = (props): React.ReactNode => {
    const { expanded, onExpand, record } = props;
    return record && record.children ? (
      <a className="custem-icon" onClick={(e) => onExpand(record, e)}>
        <Icon type={!expanded ? 'right' : 'down'} />
      </a>
    ) : (
      <span className="no-icon"></span>
    );
  };

  /**
   * 转码
   * @param text
   * @returns
   */
  const transcode = (text: string): string => {
    return encodeURIComponent(text);
  };

  /**
   * key值重复问题
   * @param record
   * @returns
   */
  const tableKey = (record): string =>
    `${transcode(record.name)}_${record.objectValue}_${record.execId}`;

  const handleGoToResult = (record: any): any => {
    message.destroy();
    if (!record.isExist) return message.warning('该组件已从该集群卸载！');
    // 1 产品包、2 服务、3 主机
    let path: string = '';
    switch (record.objectType) {
      case 1:
        path = `/deploycenter/cluster/detail/deployed?from=${location.pathname}`;
        Utils.setNaviKey('menu_deploy_center', 'deployed');
        break;
      case 2:
        path = `/opscenter/service?from=${location.pathname}`;
        record = { ...record, productName: records.productName };
        Cookie.set('em_product_id', '');
        Cookie.set('em_product_name', '');
        Cookie.set('em_current_parent_product', records.parentProductName);
        Utils.setNaviKey('menu_ops_center', 'sub_menu_service');
        break;
      default:
        path = `/opscenter/hoststatus?from=${location.pathname}`;
        Cookie.set('em_product_id', '');
        Cookie.set('em_product_name', records.productName);
        Cookie.set('em_current_parent_product', records.parentProductName);
        Utils.setNaviKey('menu_ops_center', 'sub_menu_host');
        break;
    }
    sessionStorage.setItem('service_object', JSON.stringify(record));
    history.push(path);
  };

  return (
    <>
      <Table
        className="echo-table dt-table-fixed-base"
        style={{ height: 'calc(100vh - 380px)', boxShadow: 'none' }}
        scroll={{ y: true }}
        rowKey={tableKey}
        showHeader={false}
        expandIcon={custemIcon}
        expandedRowKeys={expandedRow}
        dataSource={records.children}
        onExpand={handleExpand}
        indentSize={25}
        pagination={false}>
        <Table.Column
          title="名称"
          key="name"
          render={(text: string, record: ITableColumn) => {
            return (
              <div className="td-name">
                {record.status === 1 && (
                  <Spin size="small" style={{ margin: '0 10px' }} />
                )}
                {record.status === 2 && (
                  <Icon
                    style={{ color: '#12BC6A' }}
                    type="check-circle"
                    theme="filled"
                  />
                )}
                {record.status === 3 && (
                  <Icon
                    style={{ color: 'red' }}
                    type="close-circle"
                    theme="filled"
                  />
                )}
                {record.desc}
                {record.execId && (
                  <span className="opt-tool">
                    <a onClick={(e) => handleEvent(record, 'log')}>日志查看</a>
                    <a onClick={(e) => handleEvent(record, 'echo')}>脚本查看</a>
                  </span>
                )}
                {/* <div className='desc'>{ record.shellType === 1 ? record.desc : ''}</div> */}
              </div>
            );
          }}
        />
        <Table.Column
          title="类型"
          width="20%"
          render={(text: string, record: ITableColumn) => {
            return (
              <a onClick={(e) => handleGoToResult(record)}>
                {record.objectType != 3 ? record.objectValue : record.hostIp}
              </a>
            );
          }}
        />
        <Table.Column title="开始时间" width="20%" dataIndex="startTime" />
        <Table.Column
          title="持续时间"
          width="12%"
          dataIndex="duration"
          render={(text: string, record: any) => {
            return record.status === 1 ? (
              <TimerInc initial={parseFloat(record.status)} />
            ) : (
              text + 's'
            );
          }}
        />
      </Table>

      {visibleInfo.visible && (
        <CommandPopbox {...visibleInfo} onColse={handleEvent} />
      )}
    </>
  );
};

export default connect(mapStateToProps)(EchoTableDetail);