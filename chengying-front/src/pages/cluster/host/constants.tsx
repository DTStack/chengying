import * as React from 'react';
import { Badge, Tag, Tooltip, Icon, message } from 'antd';
import ProgressBar from '@/components/progressBar';
import './style.scss';

export const columnGenerator = ({
  statusInfo,
  isKubernetes,
  isKubernetesCustom,
  authorityList,
  progressProp,
  hostGroups,
  handleRoleSet,
}) => {
  const roleAuth = authorityList.sub_menu_role_manage;
  return [
    {
      title: '节点',
      dataIndex: 'hostname',
      render: (text, record) => {
        return (
          <div>
            <p>{record.hostname}</p>
            <p className="text-gray">{record.ip}</p>
          </div>
        );
      },
    },
    {
      title: 'agent状态',
      dataIndex: 'is_running',
      filters: [
        { text: '运行中', value: 'true' },
        { text: '已停止', value: 'false' },
      ],
      filterMultiple: true,
      render: (text) => {
        return (
          <Badge
            color={text ? '#12BC6A' : '#FF5F5C'}
            text={text ? '运行中' : '已停止'}
          />
        );
      },
    },
    {
      title: (
        <span>
          主机初始化状态
          <Tooltip title={statusInfo.title}>
            <Icon type="info-circle" />
          </Tooltip>
        </span>
      ),
      dataIndex: 'status',
      width: '15%',
      filters: statusInfo.filters,
      filterMultiple: true,
      render: (text, record: any) => {
        const finalStatus = statusInfo.finalStatus;
        return (
          <React.Fragment>
            <Badge
              color={progressProp(record.status).color}
              text={progressProp(record.status).title || record.errorMsg}
            />
            {record.status !== finalStatus && (
              <Tooltip title={record.errorMsg}>
                &nbsp;
                <Icon type="info-circle" />
              </Tooltip>
            )}
          </React.Fragment>
        );
      },
    },
    isKubernetes
      ? {
          title: '角色',
          dataIndex: 'roles',
          width: '18%',
          filters: [
            { text: 'Control', value: 'Control' },
            { text: 'Worker', value: 'Worker' },
            { text: 'Etcd', value: 'Etcd' },
            { text: '全部', value: 'all' },
          ],
          filterMultiple: true,
          render: (text, record: any) => {
            const roles = [];
            for (const key in record.roles) {
              record.roles[key] && roles.push(key);
            }
            return (
              <React.Fragment>
                {roles.map((item) => (
                  <Tag className="c-cluster__tag" key={item}>
                    {item}
                  </Tag>
                ))}
              </React.Fragment>
            );
          },
        }
      : undefined,
    {
      title: '主机分组',
      dataIndex: 'group',
      filters: hostGroups,
      filterMultiple: true,
    },
    {
      title: 'CPU',
      dataIndex: 'cpu_usage_pct',
      sorter: true,
      render: (text, record: any) => {
        return (
          <ProgressBar
            now={record.cpu_core_used_display}
            total={record.cpu_core_size_display}
          />
        );
      },
    },
    isKubernetes
      ? {
          title: 'POD',
          dataIndex: 'pod_usage_pct',
          sorter: true,
          render: (text, record: any) => {
            return (
              <ProgressBar
                now={record.pod_used_display}
                total={record.pod_size_display}
              />
            );
          },
        }
      : {
          title: '磁盘',
          dataIndex: 'disk_usage_pct',
          sorter: true,
          render: (text, record: any) => {
            return (
              <ProgressBar
                pStyle={{ width: 140, marginLeft: '-28px' }}
                now={record.disk_used_display}
                total={record.disk_size_display}
              />
            );
          },
        },
    !isKubernetesCustom
      ? {
          title: '角色',
          dataIndex: 'role_list_display',
          ellipsis: true,
          render: (roleList = []) => {
            return roleList.length === 0 ? (
              '--'
            ) : (
              <Tooltip title={roleList.map((role) => role.role_name).join(',')}>
                <span>{roleList.map((role) => role.role_name).join(',')}</span>
              </Tooltip>
            );
          },
        }
      : undefined,
    !isKubernetesCustom
      ? {
          title: '操作',
          render: (value, record) => {
            record.role_list_display = record.role_list_display || [];
            const roleList = record.role_list_display.map(
              (item) => item.role_id
            );
            return (
              <span
                className="link-btn"
                onClick={() => {
                  if (roleAuth) {
                    handleRoleSet(record.sid, roleList);
                  } else {
                    message.error('权限不足，请联系管理员！');
                  }
                }}>
                编辑角色
              </span>
            );
          },
        }
      : undefined,
  ].filter((item) => item !== undefined);
};
