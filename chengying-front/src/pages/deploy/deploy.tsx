import * as React from 'react';
import {
  Button,
  Steps,
  Collapse,
  Table,
  Modal,
  Progress,
  Alert,
  Tooltip,
  Switch,
} from 'antd';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as DeployAction from '@/actions/deployAction';
import * as HeaderAction from '@/actions/headerAction';
import { AppStoreTypes } from '@/stores';
import Logtail from '@/components/logtail';
import { EditPanel } from './edit.comp';
import { IPEditPanel } from './ip.comp';
import utils from '@/utils/utils';
import './style.scss';

const Step = Steps.Step;
const Panel = Collapse.Panel;

// const mapStateToProps = (state: AppStoreTypes) => ({
//   deploy: state.deployStore,
//   topnav: state.HeaderStore
// })

const mapStateToProps = (state: AppStoreTypes) => {
  console.log(state);
  return {
    deploy: state.DeployStore,
    topnav: state.HeaderStore,
  };
};

const mapDispatchToProps = (dispatch: any) => ({
  actions: bindActionCreators(Object.assign({}, DeployAction), dispatch),
  headerAction: bindActionCreators(Object.assign({}, HeaderAction), dispatch),
});

interface Prop {
  deploy: any;
  topnav: any;
  actions: any;
  headerAction: any;
}

interface State {
  step: number;
  showCancel: boolean;
  product_name: string;
  product_version: string;
  log_modal_visible: boolean;
  logpaths: any[];
  log_service_id: number;
}

@(connect(mapStateToProps, mapDispatchToProps) as any)
class DeployPage extends React.Component<Prop, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      step: 0,
      showCancel: true,
      product_name:
        (utils.getParamsFromUrl(window.location.href) as any).product_name ||
        'dtlog',
      product_version:
        (utils.getParamsFromUrl(window.location.href) as any).product_version ||
        '',
      log_modal_visible: false,
      logpaths: [],
      log_service_id: -1,
    };
    this.deployListInterval = null;
  }

  private deployListInterval: any = null;
  componentDidMount() {
    this.updateProductConfig({
      product_name: this.state.product_name,
      product_version: this.state.product_version,
    });
  }

  // 获取产品部署配置schema
  updateProductConfig(params: any) {
    this.props.actions.updateProductConfig(params);
  }

  // 选择切换版本
  // handleMenuClick() { }
  // 修改schema配置
  handleServiceSchemaUpdate(params: any) {
    this.props.actions.modifyServiceConfig({
      product_name: this.state.product_name,
      service_name: params.pindex,
      product_version: this.state.product_version,
      field_path: params.path,
      field: params.value,
    });
  }

  handleServiceConfigReset(params: any) {
    this.props.actions.resetServiceConfig({
      product_name: this.state.product_name,
      service_name: params.pindex,
      product_version: this.state.product_version,
      field_path: params.field_path,
    });
  }

  addIpToConfig(pindex: any) {
    this.props.actions.addIpToConfig(pindex);
  }

  handleSetServiceIp(pindex: any, ips: any) {
    this.props.actions.handleSetServiceIp({
      product_name: this.state.product_name,
      service_name: pindex,
      product_version: this.state.product_version,
      ip: ips,
    });
  }

  // 展开获取ip
  handleCollapseChange = (index: any) => {
    if (index !== undefined) {
      const { deploy } = this.props;
      const products = [];
      for (const p in deploy.product.Service) {
        products.push({
          name: p,
          data: deploy.product.Service[p],
        });
      }
      console.log(products);
      // actions.getServiceIpList({
      //     product_name: product_name,
      //     service_name: products[i - 1].name
      // });
    }
  };

  // 开始部署
  // tslint:disable
  handleDoDeploy() {
    const self = this;
    const { headerAction } = this.props;
    this.props.actions.startProductDeploy(
      {
        product_name: this.state.product_name,
        product_version: this.props.deploy.product_version,
      },
      function () {
        headerAction.getProductList();
        self.setState({
          step: 1,
        });
        self.props.actions.getDeployList(
          {
            product_name: self.state.product_name,
            product_version: self.state.product_version,
            deploy_uuid: self.props.deploy.deploy.deploy_uuid,
          },
          function () {
            clearInterval(self.deployListInterval);
          }
        );
        self.deployListInterval = setInterval(function () {
          self.props.actions.getDeployList(
            {
              product_name: self.state.product_name,
              product_version: self.state.product_version,
              deploy_uuid: self.props.deploy.deploy.deploy_uuid,
            },
            function () {
              clearInterval(self.deployListInterval);
            }
          );
        }, 3000);
      }
    );
  }

  // 返回第一步
  handleReturnToStepOne() {
    this.props.actions.returnToStepOne();
    clearInterval(this.deployListInterval);
  }

  // 显示部署日志
  handleShowDeployLog(record: any) {
    const schema = JSON.parse(record.schema);
    this.setState({
      log_modal_visible: true,
      logpaths: schema.Instance.Logs,
      log_service_id: record.instance_id,
    });
  }

  // 关闭部署日志
  handleLogModalCancel() {
    this.setState({
      log_modal_visible: false,
    });
  }

  // 返回产品列表
  handleGotoProducts() {
    // window.location.href = '/version';
    window.location.href = '/products';
    // this.props.router.replace("/products");
  }

  componentWillUnmount() {
    clearInterval(this.deployListInterval);
    this.props.actions.resetDeployStatus();
  }

  // 切换 use cloud
  handleSwitchUseCloud(product: any) {
    const { deploy, actions } = this.props;
    const service = Object.assign({}, deploy.product.Service);
    for (const s in service) {
      if (s === product.name) {
        service[s].Instance.UseCloud = !service[s].Instance.UseCloud;
        service[s].Instance.Ips = '';
      }
    }
    actions.handleSetServiceIp({
      product_name: this.state.product_name,
      service_name: product.name,
      product_version: this.state.product_version,
      ip: '',
    });
    actions.switchUseCloudByProduct(service);
  }

  handleUpdateIpByInput(pindex: any, ips: any) {
    this.props.actions.updateIpsByService(pindex, ips);
  }

  handleCancelDeploy = () => {
    this.props.actions.cancelDeploy({
      product_name: this.state.product_name,
      product_version: this.state.product_version,
    });
  };

  render() {
    const { deploy, actions } = this.props;
    // let menu = (
    //     <Menu onClick={this.handleMenuClick.bind(this)}>
    //         {deploy.versionList.map((v: any, i: any) => {
    //             return (
    //                 <Menu.Item key={i}>{v}</Menu.Item>
    //             )
    //         })}
    //     </Menu>
    // )
    const products = [];
    for (const p in deploy.product.Service) {
      products.push({
        name: p,
        data: deploy.product.Service[p],
      });
    }

    const columns = [
      {
        title: '执行时间',
        dataIndex: 'update_time',
        key: 'update_time',
      },
      {
        title: '组件',
        dataIndex: 'service_name',
        key: 'service_name',
      },
      {
        title: 'IP',
        dataIndex: 'ip',
        key: 'ip',
      },
      {
        title: '产品版本',
        dataIndex: 'product_version',
        key: 'product_version',
      },
      {
        title: '启动及初始化',
        dataIndex: 'progress',
        key: 'progress',
        render: (text: any, record: any) => {
          let s_status = 'active';
          switch (record.status) {
            case 'install fail':
            case 'run fail':
            case 'health-check fail':
            case 'health-check cancelled':
              s_status = 'exception';
              break;
          }
          return (
            <div style={{ width: '170px' }}>
              <Progress
                percent={record.progress}
                size="small"
                status={s_status as any}
              />
            </div>
          );
        },
      },
      {
        title: '启动状态',
        dataIndex: 'status',
        key: 'status',
        render: (text: any, record: any) => {
          let service_status = {};
          switch (record.status) {
            case 'install fail':
            case 'run fail':
            case 'health-check fail':
            case 'health-check cancelled':
              service_status = {
                color: '#FF5F5C',
              };
              break;
            case 'installed':
            case 'health-checked':
              service_status = {
                color: '#12BC6A',
              };
              break;
          }
          return (
            <Tooltip title={record.status_message}>
              <span style={service_status}>{text}</span>
            </Tooltip>
          );
        },
      },
      {
        title: '查看详情',
        key: 'detail',
        render: (text: any, record: any) => {
          const schema = JSON.parse(record.schema);
          // let logBtnCls = {}
          if (schema && schema.Instance && schema.Instance.Logs) {
            return (
              <span>
                <a
                  style={
                    schema.Instance.Logs.length
                      ? { display: 'inline' }
                      : { display: 'none' }
                  }
                  onClick={this.handleShowDeployLog.bind(this, record)}>
                  部署详情
                </a>
              </span>
            );
          } else {
            return <span></span>;
          }
        },
      },
    ];
    const alertSucCls = {
      margin: '10px 0',
      display: 'none',
    };
    const alertFailCls = {
      margin: '10px 0',
      display: 'none',
    };
    switch (deploy.complete) {
      case 'deploy fail':
        alertSucCls.display = 'none';
        alertFailCls.display = 'block';
        break;
      case 'deployed':
        alertSucCls.display = 'block';
        alertFailCls.display = 'none';
        break;
    }

    return (
      <div className="deploy-page">
        {/* <div className="option-bar clearfix" style={deploy.deploy.deploy_status ? { display: 'none' } : { display: 'block' }}>
                </div> */}
        <div className="config-list-wrapper">
          <Steps current={this.state.step} style={{ marginBottom: '40px' }}>
            <Step title="组件编排" description="" />
            <Step title="执行部署" description="" />
          </Steps>
          <div
            style={
              deploy.deploy.deploy_status
                ? { display: 'none' }
                : { display: 'block' }
            }>
            <Collapse accordion onChange={this.handleCollapseChange}>
              {products.map((p, index) => {
                let title = p.name + '   版本:' + p.data.Version;
                if (p.data.BaseParsed) {
                  title =
                    p.name +
                    '   版本:' +
                    p.data.Version +
                    '   继承组件:' +
                    p.data.BaseProduct;
                }
                // if (p.data.ServiceAddr && p.data.ServiceAddr.IP.length) {
                //     title += '    ip地址:' + p.data.ServiceAddr.IP.join(',');
                // }
                // console.log(p);
                let useCloudHtml: any = Object;
                if (p.data.Instance) {
                  useCloudHtml = (
                    <div className="edit-panel">
                      <span
                        style={{
                          fontSize: '14px',
                          color: '#333',
                          fontWeight: 500,
                          width: '8%',
                          display: 'inline-block',
                        }}>
                        使用云主机：
                      </span>
                      <Switch
                        size="small"
                        checked={p.data.Instance.UseCloud}
                        onChange={this.handleSwitchUseCloud.bind(this, p)}
                      />
                    </div>
                  );
                }
                return (
                  <Panel header={title} key={(index + 1).toString()}>
                    {useCloudHtml}
                    <IPEditPanel
                      data={p.data.ServiceAddr || {}}
                      pindex={p.name}
                      pname={this.state.product_name}
                      sname={p.name}
                      instance={p.data.Instance}
                      updateip={this.handleUpdateIpByInput.bind(this)}
                      setip={this.handleSetServiceIp.bind(this)}
                    />
                    <EditPanel
                      data={p.data.Instance}
                      type="Instance"
                      pindex={p.name}
                      modifyServiceConfig={actions.modifyServiceInput}
                      resetServiceConfig={this.handleServiceConfigReset.bind(
                        this
                      )}
                      onBlur={this.handleServiceSchemaUpdate.bind(this)}
                    />
                    <EditPanel
                      data={p.data.Config}
                      dpd={p.data.DependsOn || []}
                      type="Config"
                      pindex={p.name}
                      modifyServiceConfig={actions.modifyServiceInput}
                      resetServiceConfig={this.handleServiceConfigReset.bind(
                        this
                      )}
                      onBlur={this.handleServiceSchemaUpdate.bind(this)}
                    />
                  </Panel>
                );
              })}
            </Collapse>
            <div className="deploy-btn-box">
              <Button
                type="primary"
                style={
                  deploy.status === 'deploying'
                    ? { display: 'none' }
                    : { display: 'inline-block' }
                }
                onClick={this.handleDoDeploy.bind(this)}>
                开始部署
              </Button>
              <Button
                disabled
                style={
                  deploy.status === 'deploying'
                    ? { display: 'inline-block' }
                    : { display: 'none' }
                }>
                正在部署
              </Button>
            </div>
          </div>
          <div
            style={
              deploy.deploy.deploy_status
                ? { display: 'block' }
                : { display: 'none' }
            }>
            <Table
              columns={columns}
              dataSource={deploy.deploy_list}
              pagination={false}></Table>
            <Alert style={alertSucCls} message="部署完成" type="success" />
            <Alert style={alertFailCls} message="部署失败" type="error" />
            <div className="option-bar">
              <Button onClick={this.handleCancelDeploy}>取消部署</Button>
              <Button
                onClick={this.handleReturnToStepOne.bind(this)}
                style={
                  deploy.status === 'installing'
                    ? { display: 'none', margin: '0 10px' }
                    : { display: 'inline-block', margin: '0 10px' }
                }>
                上一步
              </Button>
              <Button onClick={this.handleGotoProducts.bind(this)}>
                返回版本管理
              </Button>
            </div>
            <Modal
              title="执行日志"
              footer={null}
              width={900}
              visible={this.state.log_modal_visible}
              onCancel={this.handleLogModalCancel.bind(this)}>
              <Logtail
                logs={this.state.logpaths}
                serviceid={this.state.log_service_id}
                isreset={!this.state.log_modal_visible}
              />
            </Modal>
          </div>
        </div>
      </div>
    );
  }
}

export default DeployPage;
