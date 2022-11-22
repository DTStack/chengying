import * as React from 'react';
import { Collapse, Table, Tag, Tooltip } from 'antd';
import classnames from 'classnames';
import { uniq } from 'lodash';
import './style.scss';
const Panel = Collapse.Panel;

const SUCCESS_COLOR = '#12BC6A';
const WRAN_COLOR = '#FFB310';
const DANGER_COLOR = '#EF5350'; // 比通用红色偏深一点
const RED_COLOR = '#FF5F5C';

interface IProps {
  className?: string;
  style?: React.CSSProperties;
  paneKeyName?: string;
  activeKey: string[] | string;
  onChange: (key: string | string[]) => void;
  serviceList: any[];
  serviceGroup: string[];
  onTableRowClick: Function;
  serviceNameRender?: (value?: any, record?: any) => React.ReactNode; // serviceName render
}

const ServiceOverview: React.FC<IProps> = (props) => {
  const [colors, setColors] = React.useState<any>({});
  const {
    className,
    paneKeyName = 'product_name',
    activeKey,
    onChange,
    serviceList,
    serviceGroup,
  } = props;

  React.useEffect(() => {
    if (Array.isArray(serviceGroup) && serviceGroup.length) {
      uniqGroupAndSetColor(serviceGroup);
    }
  }, [serviceGroup]);

  // 确定色值
  function uniqGroupAndSetColor(group: string[]) {
    const color = [
      'blue',
      'purple',
      'green',
      'magenta',
      'orange',
      'cyan',
      'volcano',
      'lime',
      'gold',
      'geekblue',
      'red',
    ];
    const colorMap = {};
    uniq(group).forEach((groupName: string, index: number) => {
      const len = color.length;
      colorMap[groupName] = index < len ? color[index] : color[index - len];
    });
    setColors(colorMap);
  }

  const columns = [
    {
      dataIndex: 'service_name',
      key: 'service_name',
      width: '26%',
      render:
        props.serviceNameRender ||
        ((value: string) => (
          <Tooltip placement="right" title={value}>
            {value.length > 12 ? value.slice(0, 10) + '...' : value}
          </Tooltip>
        )),
    },
    {
      dataIndex: 'status',
      key: 'status',
      width: 70,
      render: (value: string, record: any) => (
        <div>
          {value === '正常' ? (
            <Tag className="table_avatar" color={SUCCESS_COLOR}>
              {record.status_count}
            </Tag>
          ) : (
            <Tag className="table_avatar" color={DANGER_COLOR}>
              {record.status_count}
            </Tag>
          )}
          {value}
        </div>
      ),
    },
    {
      dataIndex: 'health_state',
      key: 'health_state',
      width: 85,
      render: (value: number | string, record: any) => {
        const state =
          value > 0 || value === 'healthy'
            ? { text: '健康', color: SUCCESS_COLOR }
            : { text: '不健康', color: WRAN_COLOR };
        return (
          <div>
            <Tag className="table_avatar" color={state.color}>
              {record.health_state_count}
            </Tag>
            {state.text}
          </div>
        );
      },
    },
    {
      dataIndex: 'group',
      width: 68,
      key: 'group',
      render: (value: string) => (
        <Tooltip placement="right" title={value}>
          <Tag className="table_group" color={colors[value]}>
            <div className="table_group_div">{value}</div>
          </Tag>
        </Tooltip>
      ),
    },
  ];
  return (
    <Collapse
      className={classnames('c-overview__ant-collapse', className)}
      activeKey={activeKey}
      onChange={onChange}>
      {serviceList.map((service: any, index: number) => (
        <Panel
          key={service[paneKeyName]}
          header={
            <React.Fragment>
              <i
                className="emicon emicon-folder-anticon mr-8"
                style={{
                  color:
                    service.service_count &&
                    service.service_status !== 'healthy'
                      ? RED_COLOR
                      : SUCCESS_COLOR,
                }}
              />
              {service.product_name}
            </React.Fragment>
          }>
          {service.service_count ? (
            <Table
              rowKey="service_name"
              size="small"
              className="c-overview_ant-table"
              rowClassName={() => 'c-overview_row-trigger'}
              style={{ background: '#fff' }}
              columns={columns}
              dataSource={service.service_list}
              pagination={false}
              showHeader={false}
              onRow={(record) => ({
                onClick: (e) => props.onTableRowClick(record, service),
              })}
            />
          ) : (
            <div className="service-normal">
              <i
                className="emicon emicon-health_service"
                style={{ color: SUCCESS_COLOR }}
              />
              <p>服务全部正常</p>
            </div>
          )}
        </Panel>
      ))}
    </Collapse>
  );
};
export default ServiceOverview;
