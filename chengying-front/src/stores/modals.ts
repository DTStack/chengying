export interface ActionType {
  type: string;
  payload?: any;
}
export interface Products {
  create_time: string;
  deploy_time: string;
  id: string;
  is_current_version: number;
  product: object;
  product_name: string;
  product_version: string;
  status: string;
}
export interface HeaderStoreType {
  cur_product: {
    product_id?: number;
    product_name?: string;
  };
  products: Products[];
  cur_parent_product: string;
  parentProducts: any[];
  cur_parent_cluster: {
    id: number;
    name: string;
    type: string;
    mode: number;
  };
  parentClusters: any[];
}
export interface ServiceStore {
  products: any[]; // 组件列表
  cur_product: any; // 当前组件
  cur_service: any; // 当前服务
  services: any[]; // 当前组件服务组
  hosts: any[]; // 当前服务包含主机列表
  sgList: any[];
  config: string;
  configModify: any;
  use_cloud: boolean;
  configFile: string; // 运行配置 - 配置文件
  restartService?: {
    count: number;
    list: any[];
  }; // EM提醒信息
  redService?: {
    count: number;
    list: any[];
  };
}
export interface HostStore {
  filterSelectedItem: any[];
  filterData: any[];
  hostList: any[];
  clusterHostList: any[];
  pager: {
    total: number;
    pageSize: number;
  };
  selectRows: any[];
  searchValue: string;
  selectedHost: any; // 选中的主机
  selectedHostServices: any[];
  selectedIndex: number;
  hostGroupLists: any[];
  clusterHostGroupList: any[];
}
export interface DashBoardStore {
  dashboards: any[];
  tags: any[];
  folders: any[];
}

export interface InstallGuideStore {
  sqlErro: '',
  deployState: 'normal' | 'edit';
  runtimeState: 'normal' | 'edit';
  step: number;
  selectedProduct: any;
  productPackageList: any[];
  productServices: any[];
  clusterList: any[];
  namespaceList: any[];
  productServicesInfo: any;
  hostInstallToList: any[];
  smoothSelectService: any; // 运行平滑升级的配置项
  selectedService: any; // 选中的配置服务项
  serviceHostList: any[]; // 当前服务的可分配主机
  resourceState: any; // {selectedKeys: [],targetKeys: []}保存勾选ip状态，
  paramConfigState: any[]; // [{key:field,value: value}]保存param编辑状态做统一保存
  deployUUID: string; // 点击开始部署后保存后端返回的uuid用于查询列表，需要存到localstorage里，用户返回页面后查询一下是否为-1，不是的话跳转到最后一步，只有在部署完成和点击停止部署之后置为-1
  deployList: any[]; // 部署页面的list
  stopDeployBySelf: boolean; // 是否是用户手动停止部署
  deployFinished: boolean; // 部署是否完成
  start: number; // 分页处理(逻辑感觉有问题)
  count: number;
  complete: any;
  selectedServiceList: any;
  unSelectedServiceList: any[];
  installType: string;
  clusterId: number;
  namespace: string;
  baseClusterId: number;
  baseClusterInfo: {
    baseClusterList: any[];
    hasDepends: boolean;
    dependMessage: string;
  };
  oldHostInfo: any;
  selectProductLine: any; // 选中产品线
}

export interface AddHostStore {
  current: number;
  disabled: boolean;
  hostArr: any[];
  installMsg: any[];
  hostPKArr: any[];
  forms: any;
}
// 部署
export interface DeployStore {
  upgradeType: string;
  isFirstSmooth: boolean;
  forcedUpgrade: any[];
  versionList: any[];
  product_name: string;
  product_version: string;
  status: string;
  product: {
    ProductName: string;
    ProductVersion: string;
    Service: any[];
  };
  deploy: {
    deploy_uuid: string;
    deploy_status: boolean;
  };
  deploy_list: any[];
  complete: string;
  deployType: string;
}
// 卸载
export interface UnDeployStore {
  deploy_uuid: string;
  autoRefresh: boolean;
  complete: string;
  unDeployList: any[];
  start: number;
  count: number;
  unDeployLog: any;
}
export interface UserCenterStore {
  userName: string;
  userEmail: string; // 账户
  company: string; // 公司
  phone: string; // 手机号
  createTime: string; // 创建时间
  status: number; // 状态
  role: string; // 角色
  authorityList: any; // 权限表
  authorityRouter: any[]; // 权限路由
}

// 创建编辑主机集群
export interface EditClusterStore {
  clusterInfo: any;
}
