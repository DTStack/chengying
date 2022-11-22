import * as React from 'react';
import { message, Breadcrumb, Form, Tree } from 'antd';
import { userCenterService } from '@/services';
import { Link } from 'react-router-dom';
import { formItemCenterLayout } from '@/constants/formLayout';
import './style.scss';
import utils from '@/utils/utils';

const FormItem = Form.Item;

interface IProps {
  location: any;
}

interface PermissionItemType {
  tilte: string;
  key: string;
  code: string;
  permission: number;
  children: PermissionItemType[];
  selected: boolean;
}

interface RoleInfoType {
  roleName: string;
  description: string;
  permissions: PermissionItemType[];
}

const RoleInfo: React.FC<IProps> = (props) => {
  const [authorityKeys, setAuthorityKeys] = React.useState<any[]>([]);
  const [roleInfo, setRoleInfo] = React.useState<RoleInfoType>({
    roleName: '',
    description: '',
    permissions: [],
  });

  // 获取权限树
  const getAuthorityTree = (id: number) => {
    userCenterService.getAuthorityTree({ roleId: id }).then((res: any) => {
      res = res.data;
      if (res.code === 0 && Object.keys(res.data).length) {
        const { role_name, description, permissions } = res.data;
        console.log(role_name, description, permissions);
        const list = permissions.Permissions;
        const keys = getAuthorityKeys(list, []);
        setRoleInfo({
          roleName: role_name,
          description: description,
          permissions: list,
        });
        setAuthorityKeys(keys);
      } else {
        message.error(res.msg);
      }
    });
  };

  // 获取当前id
  const getRoleId = () => {
    const search = props.location.search;
    const urlParams: any = utils.getParamsFromUrl(search);
    return urlParams.id;
  };

  // 勾选项
  const getAuthorityKeys = (tree: any[], authorityKeys: string[]) => {
    tree.forEach((item: PermissionItemType) => {
      const { selected, code, children } = item;
      item.key = item.code;
      selected && authorityKeys.push(code);
      if (children?.length) {
        getAuthorityKeys(children, authorityKeys);
      }
    });
    return authorityKeys;
  };

  React.useEffect(() => {
    const id = getRoleId();
    getAuthorityTree(id);
  }, []);

  console.log('roleInfo', roleInfo);

  return (
    <div className="role-manage-page">
      <Breadcrumb>
        <Breadcrumb.Item>
          <Link to="/usercenter/role">角色管理</Link>
        </Breadcrumb.Item>
        <Breadcrumb.Item>查看角色</Breadcrumb.Item>
      </Breadcrumb>
      <div className="role-info-content box-shadow-style">
        <Form>
          <FormItem label="角色名称" {...formItemCenterLayout}>
            {roleInfo.roleName}
          </FormItem>
          <FormItem label="角色描述" {...formItemCenterLayout}>
            {roleInfo.description}
          </FormItem>
          <FormItem label="功能权限" {...formItemCenterLayout}>
            <Tree
              checkable
              disabled
              checkStrictly={true}
              checkedKeys={authorityKeys}
              treeData={roleInfo.permissions}
              style={{
                maxHeight: 'calc(100vh - 235px)',
                overflowY: 'auto',
              }}
            />
          </FormItem>
        </Form>
      </div>
    </div>
  );
};
export default RoleInfo;
