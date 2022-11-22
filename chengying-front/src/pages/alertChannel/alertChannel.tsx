import * as React from 'react';
import { Button, Table, Icon, Modal, message } from 'antd';
import { AppStoreTypes } from '@/stores';
import { connect } from 'react-redux';
import { alertChannelService } from '@/services';
import { Link } from 'react-router-dom';
import './style.scss';
const confirm = Modal.confirm;

interface AlertChannelProp {
  router?: any;
  match?: any;
  authorityList?: any;
}

interface AlertChannelState {
  channels: any[];
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});

@(connect(mapStateToProps, undefined) as any)
export default class AlertChannelPage extends React.Component<
  AlertChannelProp,
  AlertChannelState
> {
  constructor(props: any) {
    super(props);
    this.getAlertChannels();
  }

  state: AlertChannelState = {
    channels: [],
  };

  getAlertChannels = () => {
    const self = this;
    const channels: any = [];
    alertChannelService
      .dtstackAlertChannelList({ page: 1, size: 100 })
      .then((rst: any) => {
        if (rst.data.success) {
          alertChannelService.getAlertNotifications().then((res: any) => {
            // debugger;
            for (const c of rst.data.data.data) {
              for (const g of res.data) {
                if (c.alertGateName === g.name) {
                  channels.push({
                    ...c,
                    gid: g.id,
                    type: g.type,
                    isDefault: g.isDefault,
                  });
                }
              }
            }
            console.log('channels:', channels);
            // debugger;
            self.setState({
              channels: channels,
            });
          });
        } else {
          message.error('获取通道失败！');
        }
      });
  };

  jumpToAddChannel = () => {
    // this.props.router.replace("/dashboard/addAlertChannel");
    history.pushState({}, '/dashboard/addAlertChannel');
  };

  handleDelNotification = (notification: any) => {
    const self = this;
    confirm({
      title: '是否确认删除该告警通道？',
      content: '删除后，使用该告警通道的服务将不再发出告警。',
      autoFocusButton: 'cancel',
      okType: 'primary',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      className: 'del-comfirm-dialog',
      onOk() {
        alertChannelService
          .dtstackAlertChannelDel({ id: notification.id })
          .then((rst: any) => {
            if (rst.data.success) {
              alertChannelService
                .delAlertNotification(notification.gid)
                .then((res: any) => {
                  self.getAlertChannels();
                });
            } else {
              message.error(rst.message);
            }
          });
      },
    });
  };

  render() {
    const { channels } = this.state;
    const { authorityList } = this.props;
    const CAN_EDIT = authorityList.alarm_channel_edit;
    const columns = [
      {
        title: '通道名称',
        dataIndex: 'alertGateName',
        key: 'alertGateName',
      },
      {
        title: '告警类型',
        dataIndex: 'type',
        key: 'type',
      },
      {
        title: '操作',
        dataIndex: 'option',
        key: 'option',
        // width: 100,
        render: (text: any, record: any) => {
          const editUrl =
            '/deploycenter/monitoring/addAlert?gid=' +
            record.gid +
            '&tid=' +
            record.id;
          return (
            <React.Fragment>
              {CAN_EDIT ? (
                <React.Fragment>
                  <Link to={editUrl} style={{ marginRight: 8 }}>
                    <Icon type="edit" />
                  </Link>
                  <a onClick={this.handleDelNotification.bind(this, record)}>
                    <Icon type="delete" />
                  </a>
                </React.Fragment>
              ) : (
                '--'
              )}
            </React.Fragment>
          );
        },
      },
    ];
    return (
      <div
        className="alert-channel-page"
        style={{ minHeight: document.body.clientHeight - 140 }}>
        <div className="top-navbar mb-12 clearfix">
          {CAN_EDIT && (
            <Button
              type="primary"
              style={{ float: 'right' }}
              onClick={this.jumpToAddChannel}>
              <Link to="/deploycenter/monitoring/addAlert">新增告警通道</Link>
            </Button>
          )}
        </div>
        <div className="alert-channel-body">
          <Table
            rowKey="id"
            className="dt-table-fixed-base"
            style={{ height: 'calc(100vh - 48px - 40px - 89px)' }}
            scroll={{ y: true }}
            columns={columns}
            dataSource={channels}
            pagination={false}
          />
        </div>
      </div>
    );
  }
}
