import * as React from 'react';
import {
  Table,
  Button,
  Card,
  DatePicker,
  message,
  Badge,
  Modal,
  Icon,
} from 'antd';
import { ServiceProp } from '@/pages/serviceStatus';
// import { RangePickerPresetRange, RangePickerValue } from 'antd/lib/date-picker/interface';
import { CardTabListType } from 'antd/lib/card';
import ServiceAlert from './serviceAlert';
import HostAlert from '../../hostsStatus/hostAlert';
import HealthCheck from './healthCheck';
import ServiceExpand from './serviceExpand';
import NormalTable from './normalTable';
import * as Cookie from 'js-cookie';
import HostInfoConfig from './hostInfoConfig';
import { servicePageService } from '@/services';
import { InstanceReplicaParams } from '@/services/ServicePageService';
import isEqual from 'lodash/isEqual';
import utils from '@/utils/utils';

const { confirm } = Modal;

const HA_INTERVAL: number = 10000; // 服务实例轮询间隔10s
interface IProps extends ServiceProp {
  cur_service_info: any;
  cur_product_info: any;
  dashId: string;
  getServiceGroup: () => void; // 获取服务信息
  products: any; // 滚动重启
  setCurrentService: Function;
}
interface IState {
  rangeDate: any;
  configVisible: boolean;
  config_paths: any[];
  current_instance_id: string;
  curent_config_path: string;
  hasServiceAlert: boolean;
  hasHealthAlert: boolean;
  hasHostAlert: boolean;
  use_cloud: boolean;
  haRoleClass: string; // 是否隐藏ha_role列
  expandVisible: boolean;
  key: string;
  currentHostDetail: any;
  currentHosts: any[];
}
export default class RuntimeStatus extends React.PureComponent<IProps, IState> {
  private runtimeHost: any = null;
  private runtimeAlert: any = null;

  state: IState = {
    rangeDate: [],
    configVisible: false,
    config_paths: [],
    current_instance_id: '', // 当前实例id
    curent_config_path: '',
    hasServiceAlert: false,
    hasHealthAlert: false,
    hasHostAlert: false,
    use_cloud: false,
    haRoleClass: 'hide',
    expandVisible: false,
    key: 'tab1',
    currentHostDetail: {},
    currentHosts: [],
  };

  componentDidUpdate(prevProps: IProps, prevState: IState) {
    const { cur_service_info, dashId } = this.props;
    const { currentHosts } = this.state;
    // 切换服务，重置数据
    if (
      !isEqual(prevProps.dashId, dashId) ||
      !isEqual(prevProps.cur_service_info, cur_service_info)
    ) {
      this.setState(
        {
          hasServiceAlert: false,
          hasHealthAlert: false,
          rangeDate: [],
          key: 'tab1',
          currentHostDetail:
            JSON.parse(sessionStorage.getItem('host_service')) || {},
          use_cloud: true,
        },
        () => {
          this.getServiceInfo();
          this.runtimeAlert = setTimeout(() => {
            sessionStorage.removeItem('host_service');
          }, 100);
        }
      );
    }
    if (!isEqual(currentHosts, prevState.currentHosts) && currentHosts.length) {
      this.getServiceAlert();
      this.getHealthCheck();
      this.getHostAlert();
    }
  }

  componentWillUnmount() {
    clearTimeout(this.runtimeAlert);
    clearTimeout(this.runtimeHost);
  }

  /**
   * 获取host
   */
  getServiceInfo = () => {
    const { cur_service_info, cur_product_info } = this.props;
    const { getHostsList } = this.props.actions;
    getHostsList(
      {
        product_name: cur_product_info?.product_name,
        service_name: cur_service_info?.service_name,
        namespace: utils.k8sNamespace,
      },
      (use_cloud, list: any) => {
        this.setState(
          {
            use_cloud,
            haRoleClass: list[0] && list[0].ha_role ? '' : 'hide',
            currentHosts: list,
          },
          () => {
            clearTimeout(this.runtimeHost);
            this.runtimeHost = setTimeout(() => {
              this.getServiceInfo();
            }, HA_INTERVAL);
          }
        );
      }
    );
  };

  getTableRow = (record) => {
    this.setState({
      currentHostDetail: record,
    });
  };

  // ip地址的render
  ipRender = (text: string, record: any) => {
    try {
      return JSON.parse(record.schema).Instance.HomePage ? (
        <a
          target="_blank"
          rel="noopener noreferrer"
          href={JSON.parse(record.schema).Instance.HomePage.replace(
            'link_ip',
            text
          )}
          title={JSON.parse(record.schema).Instance.HomePage.replace(
            'link_ip',
            text
          )}>
          {text}
        </a>
      ) : (
        text
      );
    } catch (error) {
      return text;
    }
  };

  // 配置
  handleTableConfig = (text, record) => {
    const { authorityList } = this.props;
    if (!authorityList.service_view) {
      return '--';
    }
    if ('schema' in record && record.schema) {
      const { Instance } = JSON.parse(record.schema);
      return (
        <span>
          <a
            style={{
              display: Instance.ConfigPaths ? 'inline' : 'none',
              paddingRight: '5px',
            }}
            onClick={() => this.handleOpenHostConfig(record)}>
            配置
          </a>
          <a
            style={{ display: Instance.Logs ? 'inline' : 'none' }}
            onClick={() => this.handleToLog(record)}>
            日志
          </a>
        </span>
      );
    } else {
      return '--';
    }
  };

  // fakeServiceColumn 初始化
  initFakeServiceColumns = () => {
    const { authorityList } = this.props;
    const CAN_SERVER_VIEW = authorityList.service_view;
    const fakeServiceColumn = [
      {
        title: 'IP地址',
        dataIndex: 'ip',
        key: 'ip',
        render: this.ipRender,
        width: '30%',
      },
      {
        title: '组件版本',
        dataIndex: 'product_version',
        key: 'product_version',
        width: '30%',
      },
      {
        title: '更新时间',
        dataIndex: 'update_time',
        key: 'update_time',
        width: '30%',
      },
      {
        title: '查看',
        dataIndex: 'config',
        key: 'config',
        render: (text: any, record: any) => {
          if (!CAN_SERVER_VIEW) {
            return '--';
          }
          if ('schema' in record && record.schema) {
            const schema = JSON.parse(record.schema);
            const style = schema.Instance.ConfigPaths
              ? { display: 'inline', padding: '0 5px' }
              : { display: 'none' };
            return (
              <a
                style={style}
                onClick={() => this.handleOpenHostConfig(record)}>
                配置
              </a>
            );
          } else {
            return '--';
          }
        },
      },
    ];
    return fakeServiceColumn;
  };

  // columnUserCloud 初始化
  initColumnUserCloud = () => {
    const columnUserCloud = [
      {
        title: 'IP地址',
        dataIndex: 'ip',
        key: 'ip',
        render: this.ipRender,
        width: '30%',
      },
      {
        title: '产品版本',
        dataIndex: 'product_version',
        key: 'product_version',
        width: '70%',
      },
    ];
    return columnUserCloud;
  };

  /**
   * 实例配置显示/关闭
   */
  handleOpenHostConfig = (record: any) => {
    const schema = JSON.parse(record.schema);
    this.props.actions.getHostConfig(
      {
        instance_id: record.id,
        configfile: schema.Instance.ConfigPaths[0],
      },
      () => {
        this.setState({
          configVisible: true,
          config_paths: schema.Instance.ConfigPaths,
          current_instance_id: record.id,
          curent_config_path: schema.Instance.ConfigPaths[0],
        });
      }
    );
  };

  // 实例配置关闭
  handleCloseHostConfig = () => {
    this.setState({
      configVisible: false,
    });
  };

  // 实例配置切换
  handleConfigPathClick = (path: any) => {
    const { current_instance_id } = this.state;
    this.props.actions.getHostConfig(
      {
        instance_id: current_instance_id,
        configfile: path,
      },
      () => {
        this.setState({
          curent_config_path: path,
        });
      }
    );
  };

  // 跳转到日志页面
  handleToLog = (item: any) => {
    const schema = JSON.parse(item.schema);
    utils.setNaviKey('menu_ops_center', 'sub_menu_diagnose_log');
    this.props.history.push('/opscenter/diagnosis/log', {
      product_name: item.product_name,
      host_ip: item.ip,
      service_name: item.service_name,
      logpaths: schema.Instance.Logs,
      log_service_id: item.id,
    });
  };

  // 获取告警列表
  getServiceAlert = () => {
    const {
      cur_service_info: { service_name },
      ServiceStore: {
        cur_product: { product_name },
      },
      dashId,
    } = this.props;
    const {
      currentHostDetail: { ip },
      currentHosts,
    } = this.state;
    if (!dashId) {
      return;
    }
    servicePageService
      .getAlertsHistory({
        dashboardId: dashId,
        ip: ip || currentHosts.map((o) => o.ip).join(','),
        product_name,
        service_name,
      })
      .then((res: any) => {
        const rst = res.data;
        if (rst.code === 0 && rst?.data) {
          const hasServiceAlert = rst.data.data?.some(
            (o) => !['ok', 'paused', 'pending'].includes(o.state)
          );
          this.setState({
            hasServiceAlert,
          });
        }
      });
  };

  getHealthCheck = () => {
    const {
      cur_service_info: { service_name },
      ServiceStore: {
        cur_product: { product_name },
      },
    } = this.props;
    const {
      currentHostDetail: { ip },
    } = this.state;
    servicePageService
      .getHealthCheck({
        product_name,
        service_name,
        ip,
      })
      .then((res: any) => {
        let data = res.data;
        if (data.code === 0 && data?.data) {
          const hasHealthAlert = data.data.list?.some(
            (o) => o.exec_status === 3
          );
          this.setState({
            hasHealthAlert,
          });
        } else {
          message.error(data.msg);
        }
      });
  };

  getHostAlert = () => {
    const {
      currentHostDetail: { ip },
      currentHosts,
    } = this.state;
    if (!currentHosts.length) {
      return;
    }
    servicePageService
      .getHostAlert({
        ip: ip || currentHosts.map((o) => o.ip).join(','),
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

  // 扩缩容模态框显示
  handleExOrCap = () => {
    const { expandVisible } = this.state;
    const { authorityList } = this.props;
    if (utils.noAuthorityToDO(authorityList, 'service_replica')) {
      return;
    }
    this.setState({ expandVisible: !expandVisible });
  };

  // 扩缩容
  handlePodChange = (replica: number) => {
    const { cur_service_info, ServiceStore, getServiceGroup } = this.props;
    const { cur_product } = ServiceStore;
    const conf: InstanceReplicaParams = {
      replica,
    };
    utils.k8sNamespace && (conf.namespace = utils.k8sNamespace);
    servicePageService
      .instanceReplica(
        {
          product_name: cur_product.product_name,
          service_name: cur_service_info.service_name,
        },
        conf
      )
      .then((response: any) => {
        const res = response.data;
        const { code, msg } = res;
        if (code === 0) {
          getServiceGroup();
        } else {
          message.error(msg);
        }
        this.handleExOrCap();
      });
  };

  // 滚动重启服务
  handleRestartServiceInTurn = () => {
    const { authorityList, ServiceStore, cur_service_info, products, actions } =
      this.props;
    if (utils.noAuthorityToDO(authorityList, 'service_roll_restart')) {
      return;
    }
    const { cur_product } = ServiceStore;
    let pid = -1;
    for (const p of products) {
      if (p.product_name === cur_product.product_name) {
        pid = p.id;
      }
    }
    confirm({
      title: '确定要滚动重启该服务下的所有主机吗？',
      content:
        '重启后最新的运行配置将生效。重启过程中会出现部分主机上服务暂时停止的情况，但服务整体将正常运行。',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      className: 'rollstart-comfirm-dialog',
      onOk: () => {
        actions.setServiceRollRestartState(cur_service_info.service_name, true);
        return new Promise((resolve, reject) => {
          setTimeout(resolve, 500);
          servicePageService
            .setServiceRollRestart({
              pid: pid,
              service_name: cur_service_info.service_name,
            })
            .then((res: any) => {
              // debugger;
              const { code, msg } = res.data;
              code === 0
                ? message.success('滚动重启完成')
                : message.success(msg);
              actions.setServiceRollRestartState(
                cur_service_info.service_name,
                false
              );
            })
            .catch(() => {
              message.success('重启完成');
              actions.setServiceRollRestartState(
                cur_service_info.service_name,
                false
              );
            });
        });
      },
      onCancel() {
        actions.setServiceRollRestartState(
          cur_service_info.service_name,
          false
        );
      },
    });
  };

  onTabChange = (key) => {
    this.setState({ key });
  };

  render() {
    const {
      cur_service_info: { service_name, IsJob },
      cur_service_info,
      HeaderStore: { cur_parent_cluster },
      ServiceStore: { cur_product, config },
      ServiceStore,
      history,
      dashId,
    } = this.props;
    const {
      configVisible,
      expandVisible,
      curent_config_path,
      config_paths,
      hasHealthAlert,
      hasServiceAlert,
      hasHostAlert,
      use_cloud,
      haRoleClass,
      currentHosts,
      rangeDate,
      currentHostDetail,
    } = this.state;

    // 当前服务完整信息
    const cur_service_infos = cur_product?.product?.Service[service_name];

    // 获取sHosts, isRestart
    const isRestart: boolean = cur_service_infos?.isRestart || false;

    // 历史遗留的奇怪的判断逻辑
    // 分了三种来分别渲染表格
    const PSEUDO_NO_CLOUD: boolean =
      cur_service_infos?.Instance?.Pseudo && !use_cloud;
    const NORMAL_USE_CLOUD: boolean = !PSEUDO_NO_CLOUD && use_cloud;
    const NORMAL_NO_CLOUD: boolean = !PSEUDO_NO_CLOUD && !use_cloud;

    let tableColumns = [];
    if (PSEUDO_NO_CLOUD) {
      tableColumns = this.initFakeServiceColumns();
    }
    if (NORMAL_USE_CLOUD) {
      tableColumns = this.initColumnUserCloud();
    }

    // 是否可以扩缩容 - k8s导入，且有isJob属性
    const isKubernetes: boolean = cur_parent_cluster?.type === 'kubernetes';

    const CAN_REPLICA: boolean =
      isKubernetes && cur_parent_cluster?.mode === 1 && !IsJob;

    const tabList: CardTabListType[] = [
      {
        key: 'tab1',
        tab: <Badge dot={hasServiceAlert}>服务告警</Badge>,
      },
      {
        key: 'tab2',
        tab: <Badge dot={hasHostAlert}>主机告警</Badge>,
      },
      {
        key: 'tab3',
        tab: <Badge dot={hasHealthAlert}>健康检查</Badge>,
      },
    ];

    const contentList = {
      tab1: (
        <ServiceAlert
          dashId={dashId}
          cur_service_info={cur_service_info}
          ServiceStore={ServiceStore}
          currentHostDetail={currentHostDetail}
          currentHosts={currentHosts}
          history={history}
        />
      ),
      tab2: (
        <HostAlert
          ip={currentHostDetail?.ip || currentHosts.map((o) => o.ip).join(',')}
          history={history}
        />
      ),
      tab3: (
        <HealthCheck
          curProduct={cur_product}
          curService={cur_service_info}
          currentHostDetail={currentHostDetail}
        />
      ),
    };

    return (
      <>
        <div className="card-wrapper has-card-extra max-content__scroll">
          <div className="mb-12 clearfix">
            <p className="text-title-bold fl-l">运行主机</p>
            {NORMAL_NO_CLOUD && !isKubernetes && Cookie.get('em_current_cluster_id') && (
              <Button
                className="fl-r"
                type="primary"
                disabled={isRestart}
                loading={isRestart}
                onClick={this.handleRestartServiceInTurn}>
                滚动重启
              </Button>
            )}
            {CAN_REPLICA && (
              <Button
                className="fl-r"
                type="primary"
                onClick={this.handleExOrCap}>
                服务扩缩容
              </Button>
            )}
          </div>
          <div className="service-list" style={{ marginBottom: 20 }}>
            {NORMAL_NO_CLOUD ? (
              <NormalTable
                {...this.props}
                sHosts={currentHosts}
                haRoleClass={haRoleClass}
                isRestart={isRestart}
                ipRender={this.ipRender}
                handleOpenHostConfig={this.handleOpenHostConfig}
                handleClickRow={this.getTableRow}
                selectedRow={currentHostDetail?.ip}
              />
            ) : (
              <Table
                rowKey="id"
                scroll={{ y: 49 * 10, x: false }}
                size="middle"
                className="dt-pagination-lower box-shadow-style"
                pagination={false}
                columns={tableColumns}
                dataSource={currentHosts}
              />
            )}
          </div>
          {NORMAL_NO_CLOUD && (
            <>
              <div className="mb-12 clearfix">
                <p className="text-title-bold fl-l">
                  运行状况&nbsp;<span>{currentHostDetail?.ip}</span>
                </p>
              </div>
              <Card
                className="box-shadow-style status-card"
                bordered={false}
                tabList={tabList}
                activeTabKey={this.state.key}
                onTabChange={this.onTabChange}>
                {contentList[this.state.key]}
              </Card>
            </>
          )}
        </div>
        {configVisible && (
          <HostInfoConfig
            configVisible={configVisible}
            curent_config_path={curent_config_path}
            data={config}
            handleCloseHostConfig={this.handleCloseHostConfig}
            handleConfigPathClick={this.handleConfigPathClick}
            config_paths={config_paths}
          />
        )}
        {/* -- k8s扩缩容 -- */}
        {expandVisible && (
          <ServiceExpand
            visible={expandVisible}
            defaultValue={currentHosts.length}
            onOk={this.handlePodChange}
            onCancel={this.handleExOrCap}
          />
        )}
      </>
    );
  }
}
