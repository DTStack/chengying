import * as React from 'react';
import { Card, Row, Col, Tooltip, Table, Progress, Badge, message } from 'antd';
import { CardTabListType } from 'antd/lib/card';
import HostAlert from './hostAlert';
import utils from '@/utils/utils';
import isEqual from 'lodash/isEqual';
import { servicePageService } from '@/services';
import * as Cookies from 'js-cookie';

interface Prop {
  detailData: any;
  runningServices: any[];
  cur_parent_cluster: any;
  getServicesList: Function;
  history?: any;
  selectedHost: any;
}

interface Stat {
  hasHostAlert: boolean;
  key: string;
}

class HostDetail extends React.Component<Prop, Stat> {
  state: Stat = {
    key: 'tab1',
    hasHostAlert: false,
  };

  componentDidUpdate(prevProps, prevState) {
    if (
      !isEqual(this.props.detailData, prevProps.detailData) &&
      this.props?.detailData?.ip
    ) {
      this.props.getServicesList({
        ip: this.props?.detailData?.ip,
        pid_list: this.props?.detailData?.pid_list,
      });
      this.getHostAlert();
    }
  }

  computeProgress = (a: string, b: string): number => {
    const usedCount: number = utils.formatGBUnit(a);
    const totalCount: number = utils.formatGBUnit(b);
    const usedPct: number = usedCount / totalCount;
    // console.log('计算结果');
    // console.log(usedCount, totalCount);
    // console.log(usedPct);
    return usedPct;
  };

  goToServer = (record, e) => {
    const { detailData = {} } = this.props;
    const path = `/opscenter/service?component=${record.product_name}&service_group=${record.group}&service=${e}`;
    utils.setNaviKey('menu_ops_center', 'sub_menu_service');
    Cookies.set('em_product_name', record.product_name);
    sessionStorage.setItem(
      'host_service',
      JSON.stringify({ ip: detailData.ip })
    );
    this.props.history.push(path);
  };

  onTabChange = (key) => {
    this.setState({ key });
  };

  getHostAlert = () => {
    const { detailData } = this.props;
    servicePageService
      .getHostAlert({
        ip: detailData?.ip,
      })
      .then((res: any) => {
        let data = res.data;
        if (data.code === 0 && data?.data) {
          const hasHostAlert = data.data.data?.some(
            (o) => o.state === 'alerting'
          );
          this.setState({
            hasHostAlert,
          });
        } else {
          message.error(data.msg);
        }
      });
  };

  render() {
    const { detailData = {}, cur_parent_cluster, history } = this.props;
    const memPre = detailData?.mem_used_display
      ? this.computeProgress(
          detailData.mem_used_display,
          detailData.mem_size_display
        )
      : 0;
    const diskPre = detailData?.disk_used_display
      ? this.computeProgress(
          detailData.disk_used_display,
          detailData.disk_size_display
        )
      : 0;
    const filePre = detailData?.file_used_display
      ? this.computeProgress(
          detailData.file_used_display,
          detailData.file_size_display
        )
      : 0;

    const columns = [
      {
        title: '组件',
        dataIndex: 'product_name_display',
        width: '33.3%',
      },
      {
        title: '服务组',
        dataIndex: 'group',
        width: '33.3%',
      },
      {
        title: '服务',
        dataIndex: 'service',
        render: (text: any, record: any) => {
          return (
            <div className="services-box">
              {record.service_name_list.split(',').map((e: any, i: number) => {
                return (
                  <a
                    style={{ marginRight: 10 }}
                    id={i + ''}
                    key={e}
                    onClick={() => this.goToServer(record, e)}>
                    {e}
                  </a>
                );
              })}
            </div>
          );
        },
      },
    ];

    const tabList: CardTabListType[] = [
      {
        key: 'tab1',
        tab: <Badge>运行服务</Badge>,
      },
      {
        key: 'tab2',
        tab: <Badge dot={this.state.hasHostAlert}>主机告警</Badge>,
      },
    ];

    const contentList = {
      tab1: (
        <Table
          rowKey={(record: any, index: number) => record.product_name + record.group + index}
          size="middle"
          className="dt-pagination-lower box-shadow-style"
          columns={columns}
          dataSource={this.props.runningServices}
          pagination={false}
        />
      ),
      tab2: <HostAlert ip={detailData.ip} history={history} />,
    };

    return (
      <div className="host-detail-page">
        <p className="text-title-bold">详细信息</p>
        <Card
          className="detail-card box-shadow-style mb-20"
          style={{ cursor: 'default' }}
          bordered={false}>
          <Row>
            <Col span={8}>
              <span>主机IP：{detailData.ip}</span>
            </Col>
            <Col span={8}>
              <span>主机名称：{detailData.hostname}</span>
            </Col>
            <Col span={8}>
              <span>主机创建时间：{detailData.created}</span>
            </Col>
          </Row>
          <Row>
            <Col span={8}>
              <span>
                agent状态：
                <span
                  className={
                    detailData.is_running ? 'c-status-normal' : 'c-status-stop'
                  }
                  style={{
                    color: detailData.is_running ? '#3f87ff' : '#FF5F5C',
                  }}>
                  {detailData.is_running ? '运行中' : '关闭'}
                </span>
              </span>
            </Col>
            <Col span={8} className="c-text-ellipsis">
              <Tooltip placement="topLeft" title={detailData.errorMsg}>
                <span>初始化状态：{detailData.errorMsg}</span>
              </Tooltip>
            </Col>
            <Col span={8}>
              <span>最近心跳时间: {detailData.updated}</span>
            </Col>
          </Row>
          <Row>
            <Col span={24}>
              <span>部署组件：{detailData.product_name_display_list}</span>
            </Col>
          </Row>
          <Row>
            <Col span={8}>
              <span className="progress-box">
                <span>
                  物理内存：{detailData.mem_used_display} /{' '}
                  {detailData.mem_size_display}
                </span>
                <Progress percent={memPre * 100} showInfo={false} />
              </span>
            </Col>
            {cur_parent_cluster.type === 'kubernetes' ? (
              <React.Fragment>
                <Col span={8}>
                  <span className="progress-box">
                    CPU使用情况：{detailData.cpu_core_used_display} /{' '}
                    {detailData.cpu_core_size_display}
                    <Progress
                      percent={detailData.cpu_usage_pct}
                      showInfo={false}
                    />
                  </span>
                </Col>
                <Col span={8}>
                  <span className="progress-box">
                    <span>
                      POD：{detailData.pod_used_display} /{' '}
                      {detailData.pod_size_display}
                    </span>
                    <Progress
                      percent={detailData.pod_usage_pct}
                      showInfo={false}
                    />
                  </span>
                </Col>
              </React.Fragment>
            ) : (
              <React.Fragment>
                <Col span={8}>
                  <span className="progress-box">
                    <span>
                      磁盘使用情况：{detailData.disk_used_display} /{' '}
                      {detailData.disk_size_display}
                    </span>
                    <Progress percent={diskPre * 100} showInfo={false} />
                  </span>
                </Col>
                <Col span={8}>
                  <span className="progress-box">
                    <span>
                      文件系统空间：{detailData.file_used_display} /{' '}
                      {detailData.file_size_display}
                    </span>
                    <Progress percent={filePre * 100} showInfo={false} />
                  </span>
                </Col>
              </React.Fragment>
            )}
          </Row>
        </Card>
        <p className="text-title-bold">运行状况</p>
        <Card
          className="box-shadow-style status-card"
          bordered={false}
          tabList={tabList}
          activeTabKey={this.state.key}
          onTabChange={this.onTabChange}>
          {contentList[this.state.key]}
        </Card>
      </div>
    );
  }
}

export default HostDetail;
