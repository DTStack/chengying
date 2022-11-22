import * as React from 'react';
import moment from 'moment';
import { connect } from 'react-redux';
import { Layout, DatePicker, Button, Table, Icon, message } from 'antd';
import MoreChart from './components/moreChart';
import {
  nodeStatusColumns,
  dtBaseColumns,
  dtBathColumns,
  alarmColumns,
} from './constant';
import { inspectionReportService } from '@/services';
import { AppStoreTypes } from '@/stores';
import './style.scss';
const { Content } = Layout;
const { RangePicker } = DatePicker;

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList
});

const idGenerator = () => {
  let _id = 0;
  return (): string => {
    return (++_id).toString();
  };
};
const id = idGenerator();

enum EnumProgressStatus {
  SUCCESS = 'SUCCESS',
  RUNNING = 'RUNNING',
}

function inspectionReport(props) {
  const end = moment();
  const start = moment().subtract(6, 'days');
  const [time, setTime] = React.useState<any[]>([start, end]);
  const [nodeStatusData, setNodeStatusData] = React.useState([]);
  const [dtAppData, setDtAppData] = React.useState({});
  const [alarmData, setAlarmData] = React.useState([]);
  const [chartConfig, setChartConfig] = React.useState([]);
  const [loadingDownload, setLoadingDownload] = React.useState({
    pedding: false,
    process: 0,
  });
  const refPoll = React.useRef(null);
  const disabledDate = (current) => {
    // if (!current || !end) return false;
    const range = 14 * 24 * 60 * 60 * 1000;
    return (
      (current && current.valueOf() > Date.now()) ||
      (Date.now() - range > current && current.valueOf())
    );
  };
  const timeChange = (date) => {
    setTime(date);
  };

  const handleProgress = (reportId: number) => {
    // 根据reportId查询进度
    refPoll.current = setInterval(() => {
      inspectionReportService
        .getReportProgress({
          id: reportId,
        })
        .then((res) => {
          const { code, msg, data } = res.data;
          if (code !== 0) return message.error(msg);
          // 如果进度为100,状态为成功
          if (data.status === EnumProgressStatus.SUCCESS) {
            clearInterval(refPoll.current);
            // 请求加锁，避免二次请求
            setTimeout(() => {
              setLoadingDownload({
                pedding: false,
                process: 100,
              });
              window.open(
                `/api/v2/inspect/download?id=${reportId}&file_path=${data.file_path}`,
                '_self'
              );
            }, 200);
          }
          setLoadingDownload({
            pedding: true,
            process: data.progress,
          });
          return undefined;
        })
        .catch((error) => {
          message.error(error.message);
          setLoadingDownload({
            pedding: false,
            process: 0,
          });
        });
    }, 500);
  };

  const handleDownload = () => {
    const { authorityList } = props;
    if (!authorityList.sub_menu_diagnose_inspect_report) {
      message.error('权限不足，请联系管理员！')
      return
    }
    // 发请求启动后端下载进程
    setLoadingDownload({
      pedding: true,
      process: 0,
    });
    inspectionReportService
      .generatorReport({
        from: moment(`${time[0].format('YYYY-MM-DD')} 00:00:00`).valueOf(),
        to: moment(`${time[1].format('YYYY-MM-DD')} 23:59:59`).valueOf(),
      })
      .then((res) => {
        const { data, msg, code } = res.data;
        if (code !== 0) return message.error(msg);
        // 根据请求返回的report_id, 启动轮询查询进度
        handleProgress(data.report_id);
        return undefined;
      })
      .catch((error) => {
        setLoadingDownload({
          pedding: false,
          process: 0,
        });
        message.error(error.message);
      });
  };

  const getNodeStatus = async () => {
    const res: any = await inspectionReportService.getNodeStatus();
    const { code, msg, data } = res.data;
    if (code === 0) {
      setNodeStatusData(data);
    } else {
      message.destroy();
      message.error(msg);
    }
  };
  const getAppStatus = async () => {
    const res: any = await inspectionReportService.getAppStatus();
    const { code, msg, data } = res.data;
    if (code === 0) {
      setDtAppData(data);
    } else {
      message.destroy();
      message.error(msg);
    }
  };
  const getHistoryAlarm = async () => {
    const res: any = await inspectionReportService.getHistoryAlarm({
      from: moment(`${time[0].format('YYYY-MM-DD')} 00:00:00`).valueOf(),
      to: moment(`${time[1].format('YYYY-MM-DD')} 23:59:59`).valueOf(),
    });
    const { code, msg, data } = res.data;
    if (code === 0) {
      setAlarmData(data.data);
    } else {
      message.error(msg);
    }
  };
  const getChartConfigList = async () => {
    const res: any = await inspectionReportService.getChartConfigList();
    const { code, msg, data } = res.data;
    if (code === 0) {
      const moduleArr = [];
      const chartConfigArr = [];
      data.reduce((pre, cur, index, arr) => {
        const obj = {
          noData: false,
          ...cur,
          moduleName: moduleArr.includes(cur.module) ? '' : cur.module,
          moduleType:
            (pre.module === 'System' || cur.module === 'System') &&
              pre.module !== cur.module
              ? cur.module
              : '',
        };
        moduleArr.push(cur.module);
        chartConfigArr.push(obj);
        return cur;
      }, {});
      setChartConfig(chartConfigArr);
    } else {
      message.destroy();
      message.error(msg);
    }
  };
  const delNoDataChart = (index: any) => {
    const arr = [...chartConfig];
    arr[index].noData = true;
    setChartConfig(arr);
  };
  React.useEffect(() => {
    getNodeStatus();
    getAppStatus();
  }, []);
  React.useEffect(() => {
    getHistoryAlarm();
    getChartConfigList();
  }, [time]);

  const timeChunk = React.useMemo(() => {
    return (
      <div className="body-time">
        <span>
          报告内容周期：
          {`${time[0].format('YYYY-MM-DD')} 至 ${time[1].format('YYYY-MM-DD')}`}
        </span>
        <span>报告生成时间：{moment().format('YYYY-MM-DD HH:mm:ss')}</span>
      </div>
    );
  }, []);

  return (
    <Layout>
      <Content className="inspection-report">
        <div className="inspection-report-header">
          <div className="inspection-report-time">
            <span style={{ color: 'red' }}>*</span>
            <span style={{ margin: '0 6px' }}>时间:</span>
            <RangePicker
              className="inspection-report-RangePicker"
              suffixIcon={<span className="emicon emicon-calendar" />}
              value={time}
              disabledDate={disabledDate}
              separator="至"
              allowClear={false}
              onChange={timeChange}
            />
          </div>
          {loadingDownload.pedding ? (
            <Button className="inspection-repoty-btn" type="primary">
              <span className="margin-right-8">报告生成中</span>
              {loadingDownload.process}%
            </Button>
          ) : (
            <Button
              className="inspection-report-btn"
              type="primary"
              onClick={handleDownload}
              icon="download">
              立即下载
            </Button>
          )}
        </div>
        <div className="inspection-report-overflow">
          <div className="inspection-report-body" id="inspection-report-body">
            <div className="body-title">巡检报告</div>
            {timeChunk}
            <div className="body-l2-title">集群状态汇总</div>
            <div className="body-tips">
              <div className="body-tips-left">
                <Icon
                  type="info-circle"
                  theme="filled"
                  style={{ color: '#999' }}
                />
              </div>
              <div className="body-tips-right">
                <span>1.报告中涉及的状态为报告下载时间点的状态。</span>
                <span>
                  2.状态为“正常”表示节点或应用当前的健康状态为健康，监控指标没有告警。状态为“异常”表示节点或应用当前的健康状态为不健康，或者监控指标有告警
                </span>
              </div>
            </div>
            <div className="body-l3-title">节点状态</div>
            <Table
              rowKey={id}
              className="border-table"
              columns={nodeStatusColumns}
              dataSource={nodeStatusData}
              pagination={false}
            />
            <div className="body-l3-title">应用状态</div>
            {Object.keys(dtAppData).map((key) => {
              if (key === 'DTBase') {
                return (
                  <div key={key}>
                    <div className="body-l4-title">{key}</div>
                    <Table
                      rowKey={id}
                      className="border-table"
                      columns={dtBaseColumns}
                      dataSource={dtAppData[key]}
                      pagination={false}
                    />
                  </div>
                );
              } else {
                return (
                  <div key={key}>
                    <div className="body-l4-title">{key}</div>
                    <Table
                      rowKey={id}
                      className="border-table"
                      columns={dtBathColumns}
                      dataSource={dtAppData[key]}
                      pagination={false}
                    />
                  </div>
                );
              }
            })}
            <div className="body-l3-title">告警记录</div>
            <Table
              rowKey={id}
              className="border-table"
              columns={alarmColumns}
              dataSource={alarmData}
              pagination={false}
            />
            <div className="body-l2-title">集群状态详细内容</div>
            {chartConfig.map((item, index) => {
              if (item.noData) {
                return null;
              }
              return (
                <div key={index}>
                  {item.moduleType ? (
                    <div className="body-l3-title">
                      {item.moduleType === 'System' ? '节点状态' : '应用状态'}
                    </div>
                  ) : null}
                  {item.moduleName ? (
                    <div className="body-l4-title">{item.moduleName}</div>
                  ) : null}
                  <MoreChart
                    config={item}
                    time={time}
                    key={index}
                    index={index}
                    delNoDataChart={delNoDataChart}
                  />
                </div>
              );
            })}
          </div>
        </div>
      </Content>
    </Layout>
  );
}

export default connect(mapStateToProps)(inspectionReport);