import * as React from 'react';
import { Table, message } from 'antd';
import { userCenterService } from '@/services';
import { Link } from 'react-router-dom';
import moment from 'moment';
import './style.scss';

const RoleManage: React.FC<any> = (props) => {
  // 获取角色信息列表
  const getRoleList = () => {
    const [dataSource, setDataSouce] = React.useState<any[]>([]);
    const [loading, setLoading] = React.useState<boolean>(false);

    React.useEffect(() => {
      setLoading(true);
      userCenterService.getRoleList().then((res) => {
        res = res.data;
        if (res.code === 0) {
          setDataSouce(res.data);
        } else {
          message.error(res.msg);
        }
        setLoading(false);
      });
    }, []);

    return { loading, dataSource };
  };

  // 获取表格列
  const initColumns = () => {
    const columns = [
      {
        title: '角色名称',
        dataIndex: 'name',
        key: 'name',
      },
      {
        title: '角色描述',
        dataIndex: 'desc',
        key: 'desc',
      },
      {
        title: '最近修改时间',
        dataIndex: 'update_time',
        key: 'update_time',
        render: (text: string) => moment(text).format('YYYY-MM-DD HH:mm:ss'),
      },
      {
        title: '操作',
        dataIndex: 'action',
        key: 'action',
        render: (text, record: any) => (
          <Link to={`/usercenter/role/view?id=${record.id}`}>查看</Link>
        ),
      },
    ];
    return columns;
  };
  const columns = initColumns();
  const { dataSource, loading } = getRoleList();
  return (
    <div className="role-manage-page">
      <Table
        rowKey="id"
        className="dt-table-last-row-noborder box-shadow-style"
        columns={columns}
        dataSource={dataSource}
        loading={loading}
        pagination={false}
      />
    </div>
  );
};
export default RoleManage;
