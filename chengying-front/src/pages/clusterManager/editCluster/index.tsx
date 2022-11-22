import * as React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { Dispatch, bindActionCreators } from 'redux';
import * as EditClusterAction from '@/actions/editClusterAction';
import * as HostAction from '@/actions/hostAction';
import { Icon, Button } from 'antd';
import utils from '@/utils/utils';
import StepOne from './stepOne';
import StepFinal from './stepFinal';
import './style.scss';
import { clusterManagerService } from '@/services';

export const linkMap = {
  hosts: ['主机集群'],
  kubernetes: ['自建Kubernetes集群', '导入已有Kubernetes集群'],
};

interface IProps {
  location: any;
  match: any;
  history: any;
  actions: EditClusterAction.EditClusterActionTypes &
    HostAction.HostActionTypes;
  clusterInfo: any;
}
interface IState {
  step: number;
  isEdit: boolean;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  clusterInfo: state.editClusterStore.clusterInfo,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(
    Object.assign({}, EditClusterAction, HostAction),
    dispatch
  ),
});

@(connect(mapStateToProps, mapDispatchToProps) as any)
export default class EditCluster extends React.PureComponent<IProps, IState> {
  constructor(props: IProps) {
    super(props);
    this.state = {
      step: 1,
      isEdit: false,
    };
    this.props.actions.resetClusterInfo();
  }

  stepOneForm: any = null;

  componentDidMount() {
    const {
      location,
      match: { params },
    } = this.props;
    const { id, type, mode, step } = utils.getParamsFromUrl(
      location.search || ''
    ) as any;
    const isEdit = params.type === 'edit';
    this.props.actions.setClusterInfo({ id: +id, type, mode: +mode });
    if (isEdit) {
      this.fetchClusterInfo(id, type);
    }
    this.setState({
      isEdit,
      step: step || 1,
    });
  }

  componentWillUnmount() {
    this.props.actions.resetHostList();
  }

  // 获取集群信息
  fetchClusterInfo = (id: number, type: string) => {
    clusterManagerService
      .getClusterInfo({ cluster_id: id }, type)
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.props.actions.setClusterInfo({ ...res.data });
        }
      });
  };

  // 提交保存基礎信息
  handleNextStep = () => {
    const { location, clusterInfo } = this.props;
    const { type, mode, yaml } = clusterInfo;
    const { isEdit } = this.state;
    this.stepOneForm.props.form.validateFields(async (err, values) => {
      if (!err) {
        values.tags = values.tags.join(',');
        const callback = () => {
          const step = this.state.step + 1;
          this.setState({
            step,
          });
          window.history.pushState(
            {},
            '0',
            location.pathname +
              `?id=${this.props.clusterInfo.id}&type=${type}&mode=${mode}&step=${step}`
          );
        };

        // k8s自建的是否可以下一步判断
        let result = true;
        if (type === 'kubernetes' && mode === 0) {
          if (!yaml) {
            await this.props.actions.getTemplate(
              Object.assign({}, this.props.clusterInfo, values)
            );
          }
          result = this.stepOneForm.validateYaml();
        }

        if (result) {
          this.props.actions.clusterSumbmitOperate(
            Object.assign({}, this.props.clusterInfo, values),
            isEdit,
            callback
          );
        }
      }
    });
  };

  handleSubmit = () => {
    const link =
      this.props.location.search.split('from=')[1] ||
      '/deploycenter/cluster/list';
    utils.setNaviKey('sub_menu_cluster_edit', 'sub_menu_cluster_list');
    this.props.history.push(link);
  };

  // 返回
  handleCancel = () => {
    const { isEdit } = this.state;
    const link =
      this.props.location.search.split('from=')[1] ||
      `/deploycenter/cluster${isEdit ? '/list' : '/create'}`;
    this.props.history.push(link);
  };

  render() {
    const { step, isEdit } = this.state;
    const { clusterInfo } = this.props;
    const stepProps = {
      ...this.props,
      action: this.props.actions,
      isEdit,
      handleCancel: this.handleCancel,
    };
    return (
      <div className="cluster-container">
        <Link
          className="nav-link text-title-bold"
          to={
            this.props.location.search.split('from=')[1] ||
            `/deploycenter/cluster${isEdit ? '/list' : '/create'}`
          }>
          <Icon type="left" style={{ marginRight: 5 }} />
          {isEdit ? '编辑' : '添加'}集群（
          {linkMap[clusterInfo.type][clusterInfo.mode]}）
        </Link>
        <div className="edit-cluster-container">
          {step === 1 ? (
            <StepOne
              {...stepProps}
              wrappedComponentRef={(form) => (this.stepOneForm = form)}
            />
          ) : (
            <StepFinal {...stepProps} />
          )}
          <div className="edit-cluster-bottom">
            {step === 1 ? (
              <React.Fragment>
                <Button
                  type="default"
                  className="mr-10"
                  onClick={this.handleCancel}>
                  取消
                </Button>
                <Button type="primary" onClick={this.handleNextStep}>
                  下一步
                </Button>
              </React.Fragment>
            ) : (
              <Button type="primary" onClick={this.handleSubmit}>
                完成
              </Button>
            )}
          </div>
        </div>
      </div>
    );
  }
}
