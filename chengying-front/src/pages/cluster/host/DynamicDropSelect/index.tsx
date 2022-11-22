// TODO: 代码写的很乱，待优化,先cover掉功能
import * as React from 'react';
import { Select, Input, Icon, message, Modal } from 'antd';
import { clusterHostService } from '@/services';
import classnames from 'classnames';
import './style.scss';
const { Option } = Select;

interface IOption {
  value: string | number;
  key: string | number;
  label: string | number;
  roleType: 1 | 2;
}

interface IPropsDynamicDropSelect {
  cur_parent_cluster: any;
  container: any;
  value: any[];
  onChange: Function;
  updateRoleList: Function;
}

interface IStateDynamicDropSelect {
  isOpen: boolean;
  status: {
    flagInputFocus: boolean;
    flagSelectFocus: boolean;
  };
  isADD: boolean;
  isEDIT: boolean;
  editID: string;
  validateStatus: boolean;
  input: string;
  options: IOption[];
}

export default class DynamicDropSelect extends React.Component<
  IPropsDynamicDropSelect,
  IStateDynamicDropSelect
> {
  state = {
    isOpen: false,
    value: [],
    input: '',
    status: {
      flagInputFocus: false,
      flagSelectFocus: false,
    },
    options: [],
    isADD: false,
    isEDIT: false,
    editID: '',
    validateStatus: true,
  };

  getRoleListByCluster = async (callback?) => {
    const { cur_parent_cluster } = this.props;
    const res = await clusterHostService.getRoleList({
      cluster_id: cur_parent_cluster.id,
    });
    const { code, data } = res.data;
    if (code === 0) {
      const roleOptions = data.map((role) => ({
        key: '' + role.role_id,
        value: '' + role.role_id,
        label: role.role_name,
        roleType: role.role_type,
      }));
      this.setState(
        {
          options: roleOptions,
          isADD: false,
          input: '',
        },
        () => {
          if (typeof callback === 'function') {
            callback(roleOptions);
          }
        }
      );
    }
  };

  addRole = async (roleName: string) => {
    const { cur_parent_cluster } = this.props;
    const res = await clusterHostService.addRole({
      cluster_id: cur_parent_cluster.id,
      role_name: roleName,
    });
    const { code, msg } = res.data;
    if (code === 0) {
      message.success('添加成功');
      this.getRoleListByCluster();
    } else {
      message.error(msg);
    }
  };

  deleteRole = async (id: number) => {
    const { cur_parent_cluster } = this.props;
    const res = await clusterHostService.deleteRole({
      cluster_id: cur_parent_cluster.id,
      role_id: id,
    });
    const { code } = res.data;
    if (code === 0) {
      message.success('删除成功');
      this.getRoleListByCluster((roles) => {
        const { value } = this.props;
        const refreshValue = this.refreshSelectedValue(value, roles);
        this.props.updateRoleList(refreshValue);
      });
    } else {
      message.error('删除失败');
    }
  };

  componentDidMount() {
    this.getRoleListByCluster();
  }

  handleSelectStatusChange = (flag: boolean) => {
    const { status } = this.state;
    this.setState({
      status: {
        ...status,
        flagSelectFocus: flag,
      },
    });
  };

  handleInputStatusChange = (flag: boolean) => {
    const { status } = this.state;
    this.setState({
      status: {
        ...status,
        flagInputFocus: flag,
      },
    });
  };

  hideSelect = () => {
    const { status, isADD, isEDIT } = this.state;
    if (isADD || isEDIT) {
      message.error('请先保存');
      return false;
    }
    this.setState({
      status: {
        ...status,
        flagInputFocus: false,
        flagSelectFocus: false,
      },
    });
    return true;
  };

  handleOptionAdd = () => {
    const { isADD, isEDIT } = this.state;
    if (isADD === true || isEDIT) {
      message.error('请先保存');
      return;
    }
    this.setState({
      isADD: true,
    });
  };

  handleOptionSave = () => {
    const { input } = this.state;
    const validateStatus = this.validator(input);
    if (!validateStatus.status) {
      message.error(validateStatus.message);
      this.setState({
        validateStatus: validateStatus.status,
      });
      return;
    }
    this.addRole(input);
  };

  handleOptionDelete = (item?) => {
    const { isADD, isEDIT } = this.state;
    if (item && (isADD || isEDIT)) {
      message.error('请先保存');
      return;
    }

    if (!item) {
      this.setState({
        isADD: false,
        input: '',
        validateStatus: true,
      });
    } else {
      Modal.confirm({
        title: '确认要删除该角色？',
        content: '删除后该角色将不可用！',
        icon: (
          <Icon
            type="close-circle"
            theme="filled"
            style={{ color: '#FF5F5C' }}
          />
        ),
        okType: 'danger',
        onOk: () => {
          if (this.state.isADD === true) {
            message.error('请先保存');
            return;
          }
          this.deleteRole(parseInt(item.value));
        },
        okText: '删除',
        zIndex: 1000000,
      });
    }
  };

  renameRoleName = (id, roleName, callback?) => {
    const validateStatus = this.validator(roleName);
    if (!validateStatus.status) {
      message.error(validateStatus.message);
      this.setState({
        validateStatus: validateStatus.status,
      });
      return;
    }
    const { cur_parent_cluster } = this.props;
    clusterHostService
      .modifyRole({
        cluster_id: cur_parent_cluster.id,
        role_id: parseInt(id),
        new_name: roleName,
      })
      .then((res) => {
        const { code, msg } = res.data;
        if (code === 0) {
          if (typeof callback === 'function') {
            callback();
          }
          message.success('保存成功');
          this.getRoleListByCluster();
          this.setState({
            isEDIT: false,
            input: '',
          });
        } else {
          message.error(msg);
        }
      });
  };

  validator = (value) => {
    const validateStatus = {
      status: true,
      message: '',
    };

    if (value.trim() === '') {
      validateStatus.status = false;
      validateStatus.message = '名称不能为空';
    } else if (!/^[0-9a-zA-Z_]{1,}$/.test(value.trim())) {
      validateStatus.status = false;
      validateStatus.message = '仅支持英文、数字、下划线';
    }

    this.setState({
      validateStatus: validateStatus.status,
    });
    return validateStatus;
  };

  // 数据清理，避免出现想后端传递不存在的角色
  refreshSelectedValue = (value, map) => {
    const filter = value.filter(
      (id) => map.findIndex((option) => option.value === `${id}`) > -1
    );
    return filter;
  };

  handleOptionEdit = (id, roleName) => {
    const { isADD, isEDIT } = this.state;
    if (isADD || isEDIT) {
      message.error('请先保存');
      return;
    }
    this.setState({
      isEDIT: true,
      editID: id,
      input: roleName,
    });
  };

  renderDropItem = (optionItem: IOption) => {
    const { editID, input, validateStatus } = this.state;
    const { value } = this.props;
    const container = (chidlren, isSelected = false) => {
      return (
        <div
          onClick={(e) => message.error('请先保存')}
          className="ant-select-dropdown-menu-item ds-option"
          key={optionItem.key}
          style={{
            width: '100%',
            height: '32px',
            padding: '0px 32px 0px 12px',
            lineHeight: '32px',
            background: isSelected ? 'rgb(242, 249, 255)' : undefined,
          }}>
          {chidlren}
        </div>
      );
    };
    if (optionItem.key === editID) {
      return container(
        <div onClick={(e) => e.stopPropagation()}>
          <Input
            onChange={(e) => {
              this.setState({
                input: e.currentTarget.value,
              });
            }}
            className={classnames({
              rolenameinput: true,
              rolenamenull: !validateStatus,
            })}
            onClick={(e) => e.stopPropagation()}
            style={{ height: '24px', width: '260px' }}
            value={input}
          />
          <Icon
            className="ml-5"
            onClick={() => {
              this.renameRoleName(editID, this.state.input);
            }}
            style={{
              fontSize: '12px',
              color: '#999999',
              cursor: 'pointer',
              lineHeight: '28px',
              marginRight: '2px',
              marginLeft: '6px',
            }}
            type="check"
          />
          <Icon
            type="close"
            style={{
              cursor: 'pointer',
              color: '#999999',
              fontSize: '12px',
              lineHeight: '28px',
              marginLeft: '2px',
            }}
            onClick={() => {
              this.setState({
                isEDIT: false,
                input: '',
              });
            }}
          />
        </div>
      );
    } else {
      const isSelected =
        value.findIndex((key) => `${key}` === optionItem.key) > -1;
      return container(
        <div>
          <div style={{ display: 'inline-block', width: '260px' }}>
            {optionItem.label}
          </div>
          {optionItem.roleType !== 1 ? (
            <>
              <Icon
                type="delete"
                className="fr emicon"
                style={{
                  lineHeight: '34px',
                  color: '#999999',
                  fontSize: '12px',
                  marginLeft: '10px',
                  marginRight: '-5px',
                }}
                onClick={(e) => {
                  this.handleOptionDelete(optionItem);
                  e.stopPropagation();
                }}
              />
              <Icon
                type="edit"
                onClick={(e) => {
                  this.handleOptionEdit(optionItem.key, optionItem.label);
                  e.stopPropagation();
                }}
                className="fr emicon"
                style={{
                  lineHeight: '34px',
                  color: '#999999',
                  fontSize: '12px',
                  marginRight: '-6px',
                }}
              />
            </>
          ) : null}
          {isSelected ? (
            <Icon
              type="check"
              style={{
                cursor: 'pointer',
                color: '#12BC6A',
                fontSize: '12px',
                lineHeight: '28px',
                marginLeft: '37px',
              }}
              onClick={() => {
                this.setState({
                  isEDIT: false,
                  input: '',
                });
              }}
            />
          ) : null}
        </div>,
        isSelected
      );
    }
  };

  getDropRender = (menu, options: IOption[]) => {
    const { isADD, input, validateStatus, isEDIT } = this.state;
    const base = (
      <div>
        {isADD ? (
          <div
            className="edit-role-item"
            style={{ padding: '4px 20px 4px 12px' }}>
            <Input
              style={{ width: '260px', height: '24px' }}
              className={classnames({
                rolenameinput: true,
                rolenamenull: !validateStatus,
              })}
              value={input}
              onChange={(e) => {
                this.setState({
                  input: e.currentTarget.value,
                });
              }}
              onBlur={() => this.handleInputStatusChange(false)}
              onFocus={() => this.handleInputStatusChange(true)}
            />
            <Icon
              className="ml-5"
              onClick={this.handleOptionSave}
              style={{
                color: '#999999',
                lineHeight: '22px',
                verticalAlign: 'middle',
                marginLeft: '6px',
              }}
              type="check"
            />
            <Icon
              type="close"
              // className="emicon"
              style={{
                cursor: 'pointer',
                color: '#999999',
                fontSize: '12px',
                lineHeight: '22px',
                marginLeft: '3px',
                verticalAlign: 'middle',
              }}
              onClick={() => this.handleOptionDelete()}
            />
          </div>
        ) : null}
        <div
          style={{ width: '100%', height: '1px', background: '#E8E8E8' }}></div>
        <div>
          <div className="drop-render">
            <div className="circle-plus">
              <span className="emicon emicon-plus" />
            </div>
            <a
              style={{
                flex: 'none',
                padding: '8px',
                display: 'inline-block',
                cursor: 'pointer',
              }}
              onClick={this.handleOptionAdd}>
              自定义类型
            </a>
          </div>
        </div>
      </div>
    );

    if (isEDIT) {
      return (
        <div className="dynamic-drop-select">
          <div style={{ padding: '4px 0' }}>
            {options.map(this.renderDropItem)}
          </div>
          {base}
        </div>
      );
    } else {
      return (
        <div className="dynamic-drop-select">
          {menu}
          {base}
        </div>
      );
    }
  };

  render() {
    const { status, options } = this.state;
    const isOpen = status.flagInputFocus || status.flagSelectFocus;

    return (
      <span onClick={(e) => e.stopPropagation()}>
        <Select
          className="ds-select"
          mode="multiple"
          open={isOpen}
          value={this.props.value.map((item) => `${item}`)}
          suffixIcon={<Icon type="delete" />}
          placeholder="请选择角色"
          optionLabelProp="label"
          showSearch={false}
          onChange={(value) => {
            this.props.onChange(
              value.map((roleId) => {
                return parseInt(roleId);
              })
            );
          }}
          optionFilterProp="label"
          onFocus={() => {
            this.handleSelectStatusChange(true);
          }}
          // getPopupContainer={() => this.props.container}
          dropdownRender={(menu) => this.getDropRender(menu, options)}>
          {options.map((item) => (
            <Option
              className="ds-option"
              label={item.label}
              key={item.value}
              value={item.value}>
              <span onClick={(e) => e.stopPropagation()}>
                {item.label}
                {/* TODO: 编辑操作 */}
              </span>
              {item.roleType !== 1 ? (
                <>
                  <Icon
                    type="delete"
                    className="fr emicon"
                    style={{
                      marginTop: '4px',
                      color: '#999999',
                      fontSize: '12px',
                      marginLeft: '10px',
                      marginRight: '-5px',
                    }}
                    onClick={(e) => {
                      this.handleOptionDelete(item);
                      e.stopPropagation();
                    }}
                  />
                  <Icon
                    type="edit"
                    onClick={(e) => {
                      this.handleOptionEdit(item.key, item.label);
                      e.stopPropagation();
                    }}
                    className="fr emicon"
                    style={{
                      marginTop: '4px',
                      color: '#999999',
                      fontSize: '12px',
                      marginRight: '-6px',
                    }}
                  />
                </>
              ) : null}
            </Option>
          ))}
        </Select>
      </span>
    );
  }
}
