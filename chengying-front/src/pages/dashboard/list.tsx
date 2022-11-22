import * as React from 'react';
import { Checkbox, Button, Icon, Modal, Tag, Select, message } from 'antd';
import { dashboardService } from '@/services';
import TagColors from './serviceColorMap';
import * as _ from 'lodash';
const Option = Select.Option;

// const mapStateToProps = (state: AppStoreTypes) => ({
//   list: state.DashBoardStore.dashboards
// });

interface ListProps {
  list: any;
  canEdit: boolean;
  onChange: () => void;
}
interface ListState {
  dashes: any[];
  sFolder: any[];
  sDash: any[];
  moveModalVisible: boolean;
  mFolder: any;
}

// @(connect(mapStateToProps) as any)
export default class DashListComp extends React.Component<
  ListProps,
  ListState
> {
  constructor(props: ListProps) {
    super(props);
  }

  convertDashList = (list: any[]) => {
    list.forEach((folder: any) => {
      if (folder.id === 0) {
        folder.isOpen = true;
      } else {
        folder.isOpen = false;
      }
      folder.checked = false;
      folder.list.forEach((dash: any) => {
        dash.checked = false;
      });
    });
    return list;
  };

  state: ListState = {
    dashes: [],
    sFolder: [],
    sDash: [],
    moveModalVisible: false,
    mFolder: {
      id: '',
      uid: '',
      title: '',
    },
  };

  componentDidMount() {
    this.setState({
      dashes: this.convertDashList(this.props.list),
    });
  }

  UNSAFE_componentWillReceiveProps(nextProps: ListProps) {
    if (nextProps.list.length > 0) {
      this.setState({
        dashes: this.convertDashList(nextProps.list),
      });
    }
  }

  buildDashList = (list: any[]) => {
    return list.map((dash: any) => (
      <li
        className="dash-list-sub-row clearfix"
        key={dash.id}
        onClick={this.handleDashClick.bind(this, dash)}>
        <Checkbox
          checked={dash.checked}
          onChange={this.handleSelectDash.bind(this, dash.id)}
          style={{ marginLeft: 48 }}
        />
        <Icon type="dashboard" style={{ margin: '0px 10px 0 16px' }} />
        {dash.title}
        {dash.tags.length
          ? dash.tags.map((tag: any, i: number) => {
              return (
                <Tag
                  style={{
                    float: 'right',
                    border: TagColors[tag]
                      ? `1px solid rgb(${TagColors[tag]})`
                      : `1px solid rgb(${TagColors.default})`,
                    backgroundColor: TagColors[tag]
                      ? `rgba(${TagColors[tag]},0.10)`
                      : `rgb(${TagColors.default},0.10)`,
                    color: TagColors[tag]
                      ? `rgb(${TagColors[tag]})`
                      : `rgb(${TagColors.default})`,
                  }}
                  key={i}>
                  {tag}
                </Tag>
              );
            })
          : ''}
      </li>
    ));
  };

  /**
   * dashboard文件夹展开/关闭操作
   * @param folder(obj) 被点击的文件夹
   * @memberof DashListComp
   */
  handleFolderClick = (folder: any, className?: string) => {
    // antd3.9.0之后icon使用了svg，i标签层级就变成了i>svg>path，点击之后就会冒泡到svg上拿不到className，改成传参方式判断
    const { dashes } = this.state;

    if (
      className.indexOf('dash-list-row') > -1 ||
      className.indexOf('anticon') > -1
    ) {
      for (const f of dashes) {
        if (f.id === folder.id) {
          f.isOpen = !folder.isOpen;
        } else {
          f.isOpen = false;
        }
      }
      this.setState({
        dashes,
      });
    }
  };

  /**
   * 点击打开dashboard
   * @param dash(obj) 被点击的dash对象
   * @memberof DashListComp
   */
  handleDashClick = (dash: any, e: any) => {
    e.stopPropagation();
    if (e.target.getAttribute('class').indexOf('dash-list-sub-row') > -1) {
      window.location.href =
        '/deploycenter/monitoring/dashdetail?url=' +
        encodeURIComponent(dash.url);
    }
  };

  /**
   * 全选/不选 folder&dashboard
   * @memberof DashListComp
   */
  handleSelectAll = (e: any) => {
    const { dashes } = this.state;
    const sf = [];
    const sd = [];
    if (e.target.checked) {
      for (const f of dashes) {
        f.checked = true;
        sf.push(f.uid);
        for (const d of f.list) {
          d.checked = true;
          sd.push(d);
        }
      }
    } else {
      for (const f of dashes) {
        f.checked = false;
        for (const d of f.list) {
          d.checked = false;
        }
      }
    }
    // debugger;
    this.setState({
      dashes: dashes,
      sFolder: sf,
      sDash: sd,
    });
  };

  /**
   * 选择文件夹checkbox
   * @memberof DashListComp
   */
  handleSelectFolder = (folderUid: number, e: any) => {
    // debugger;
    e.stopPropagation();
    let { dashes, sFolder, sDash } = this.state;
    console.log('QQ');
    console.log(dashes);
    for (const f of dashes) {
      console.log(f);
      if (f.uid === folderUid) {
        f.checked = e.target.checked;
        for (const d of f.list) {
          d.checked = e.target.checked;
        }
        if (e.target.checked) {
          sFolder.push(f.uid);
          sDash = _.concat(sDash, f.list);
        } else {
          // sFolder = sFolder.splice(sFolder.indexOf(f.id)-1,1)
          sFolder = _.difference(sFolder, [folderUid]);
          sDash = _.differenceBy(sDash, f.list, 'id');
        }
      }
    }
    this.setState(
      {
        dashes,
        sFolder,
        sDash,
      },
      () => {
        console.log(this.state.dashes);
        console.log(this.state.sFolder);
        console.log(this.state.sDash);
      }
    );
  };

  /**
   * 选择仪表盘checkbox
   * @memberof DashListComp
   */
  handleSelectDash = (dashId: number, e: any) => {
    e.stopPropagation();
    let { dashes, sDash } = this.state;
    for (const f of dashes) {
      for (const d of f.list) {
        if (d.id === dashId) {
          d.checked = e.target.checked;
          if (e.target.checked) {
            sDash.push(d);
          } else {
            sDash = _.differenceBy(sDash, [d], 'id');
          }
        }
      }
    }
    this.setState({
      dashes,
      sDash,
    });
  };

  /**
   * 导出仪表盘
   */
  handleDashExport = () => {
    const { canEdit } = this.props;
    if (!canEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    const { sDash } = this.state;
    const dashIds = sDash.map((o) => o.id);
    const url = `/api/v2/dashboard/export?dashboardId=${dashIds.join(',')}`;
    window.open(url);
  };

  /**
   * 删除文件夹或者仪表盘
   * @memberof DashListComp
   */
  showDeleteModal = () => {
    const { canEdit } = this.props;
    if (!canEdit) {
      message.error('权限不足，请联系管理员！');
      return;
    }
    Modal.confirm({
      title: '确定你要删除选中的dashboard吗？',
      icon: <Icon type="exclamation-circle" theme="filled" />,
      okType: 'danger',
      onOk: () => {
        this.handleConfirmDel();
      },
    });
  };

  /**
   * 确认删除仪表盘
   * @memberof DashListComp
   */
  handleConfirmDel = () => {
    const self = this;
    const { sDash, sFolder } = this.state;
    let dash = [];
    for (const f in sFolder) {
      for (const b of sDash) {
        if (b.folderUid === sFolder[f]) {
          dash.push(b);
        }
      }
      dashboardService.delFolderByUid(sFolder[f]).then((rst: any) => {
        // rst = rst.data;
        if (sFolder.length - 1 === parseInt(f)) {
          self.setState({
            sFolder: [],
            // delModalVisible: false
          });
          if (sDash.length === 0) {
            self.props.onChange();
          }
        }
      });
    }
    dash = _.differenceBy(sDash, dash, 'id');
    for (const d in dash) {
      dashboardService.delDashByUid(dash[d].uid).then((res: any) => {
        // res = res.data;
        if (dash.length - 1 === parseInt(d)) {
          self.setState({
            sFolder: [],
            sDash: [],
            // delModalVisible: false
          });
          self.props.onChange();
        }
      });
    }
  };

  /**
   * 打开/关闭移动操作弹层
   *
   * @memberof DashListComp
   */
  showMoveModal = () => {
    this.setState({
      moveModalVisible: true,
    });
  };

  hideMoveModal = () => {
    this.setState({
      moveModalVisible: false,
      mFolder: {
        id: '',
        uid: '',
        title: '',
      },
    });
  };

  handleChangeMoveFolder = (folderId: any) => {
    let folder = null;
    const { dashes } = this.state;
    for (const f of dashes) {
      if (f.id === folderId) {
        folder = f;
      }
    }
    this.setState({
      mFolder: {
        id: folder.id,
        title: folder.title,
        uid: folder.uid,
      },
    });
  };

  handleConfirmMove = () => {
    const self = this;
    const { sDash, mFolder } = this.state;
    for (const d in sDash) {
      dashboardService
        .createDashboard({
          dashboard: sDash[d],
          folderId: mFolder.id,
          overwrite: true,
        })
        .then((res: any) => {
          if (sDash.length - 1 === parseInt(d)) {
            self.setState({
              sFolder: [],
              sDash: [],
              moveModalVisible: false,
            });
            self.props.onChange();
          }
        });
    }
  };

  render() {
    const { dashes, sFolder, sDash, moveModalVisible, mFolder } = this.state;
    return (
      <div className="dash-list-wrapper box-shadow-style">
        <div className="dash-list-header">
          <Checkbox onChange={this.handleSelectAll} />
          {sFolder.length || sDash.length ? (
            <>
              <Button
                style={{ marginRight: 10 }}
                onClick={this.handleDashExport}>
                导出
              </Button>
              <Button onClick={this.showDeleteModal}>删除</Button>
            </>
          ) : (
            ''
          )}
        </div>
        <div className="dash-list-body">
          <ul className="folder-list">
            {dashes.length
              ? dashes.map((folder) => {
                  return (
                    <li
                      key={folder.id}
                      onClick={(e) =>
                        this.handleFolderClick(
                          folder,
                          folder.isOpen
                            ? 'dash-list-row active clearfix'
                            : 'dash-list-row clearfix'
                        )
                      }
                      className={
                        folder.isOpen
                          ? 'dash-list-row active clearfix'
                          : 'dash-list-row clearfix'
                      }>
                      <Checkbox
                        checked={folder.checked}
                        onChange={this.handleSelectFolder.bind(
                          this,
                          folder.uid
                        )}
                        style={{ marginLeft: 16 }}
                      />
                      <Icon type="folder-open" />
                      <Icon type="folder" />
                      {folder.title}
                      <Icon type="down" style={{ marginRight: 16 }} />
                      <Icon type="right" style={{ marginRight: 16 }} />
                      <ul
                        className="dash-list"
                        style={
                          folder.isOpen
                            ? {
                                display: 'block',
                                maxHeight: 'calc(100vh - 230px)',
                                overflowY: 'auto',
                              }
                            : { display: 'none' }
                        }>
                        {this.buildDashList(folder.list)}
                      </ul>
                    </li>
                  );
                })
              : ''}
          </ul>
        </div>
        <Modal
          title="移动"
          visible={moveModalVisible}
          onOk={this.handleConfirmMove}
          onCancel={this.hideMoveModal}
          okText="确认"
          cancelText="取消">
          <p style={{ lineHeight: '40px' }}>将选择的dashboard移动到：</p>
          <p>
            <span>文件夹：</span>
            <Select
              style={{ width: 160 }}
              onChange={this.handleChangeMoveFolder}
              value={mFolder.id}>
              {dashes.map((folder) => {
                return (
                  <Option key={folder.id} value={folder.id}>
                    {folder.title}
                  </Option>
                );
              })}
            </Select>
          </p>
        </Modal>
      </div>
    );
  }
}
