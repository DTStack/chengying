import asyncComponent from '@/components/asyncComponent';

/* ----------- layout布局 ----------- */
const ClusterLayout = asyncComponent(
  import('@/layouts/clusterLayout'),
  'ClusterLayout'
);
const RootLayout = asyncComponent(import('@/layouts/rootLayout'), 'RootLayout');
const Login = asyncComponent(import('@/pages/login'), 'Login');

/* ----------- 其他 ----------- */
const HostPage = asyncComponent(import('@/pages/host'), 'HostPage');

const NoPermission = asyncComponent(
  import('@/components/noPermission'),
  'NoPermission'
);

/* ---------- 一级菜单 ---------- */
/* ------ 集群管理 ------ */
const ClusterManagerPage = asyncComponent(
  import('@/pages/clusterManager/list'),
  'ClusterManagerPage'
);
const ClusterManagerCreate = asyncComponent(
  import('@/pages/clusterManager/createCluster'),
  'ClusterManagerCreate'
);
const ClusterManagerEdit = asyncComponent(
  import('@/pages/clusterManager/editCluster'),
  'ClusterManagerEdit'
);
/* ------ 应用管理 ------ */
const DeployedApp = asyncComponent(
  import('@/pages/components/container'),
  'DeployedApp'
);
const DeployHistory = asyncComponent(
  import('@/pages/components/container'),
  'DeployHistory'
);
const PatchHistory = asyncComponent(
  import('@/pages/components/container'),
  'PatchHistory'
);
const ProductPackage = asyncComponent(
  import('@/pages/appManager/productPackage/productPackage'),
  'ProductPackage'
);
/* ------ 监控 ------ */
const DashboardPage = asyncComponent(
  import('@/pages/dashboard/dashboard'),
  'DashboardPage'
);
const DashboardDetail = asyncComponent(
  import('@/pages/dashboard/detail'),
  'DashboardDetail'
);
const AlertManager = asyncComponent(
  import('@/pages/alertManager'),
  'AlertManager'
);
const AddAlertChannelPanel = asyncComponent(
  import('@/pages/alertChannel/addAlertChannel/addAlertChannel'),
  'AddAlertChannelPanel'
);

const Steps = asyncComponent(import('@/pages/installGuide/steps'), 'Steps');

/* ------ 系统配置 ------ */
const Security = asyncComponent(
  import('@/pages/systemConfig/security'),
  'Security'
);

const GlobalConfig = asyncComponent(
  import('@/pages/systemConfig/globalConfig'),
  'GlobalConfig'
);

/* ------ 脚本回显 ------ */
const CommandEchoIndex = asyncComponent(
  import('@/pages/command'),
  'CommandEchoIndex'
);
const CommandEchoDetail = asyncComponent(
  import('@/pages/command/details'),
  'CommandEchoDetail'
);

/* ---------- 用户管理 ---------- */
const Members = asyncComponent(import('@/pages/userCenter/members'), 'Members');
const RoleManage = asyncComponent(
  import('@/pages/userCenter/roleManage'),
  'RoleManage'
);
const RoleInfo = asyncComponent(
  import('@/pages/userCenter/roleManage/roleInfo'),
  'RoleInfo'
);
const SelfInfo = asyncComponent(
  import('@/pages/userCenter/selfInfo'),
  'SelfInfo'
);
const SecurityAudit = asyncComponent(
  import('@/pages/userCenter/securityAudit'),
  'SecurityAudit'
);

/* ---------- 脚本管理 ---------- */
const BackupConfig = asyncComponent(
  import('@/pages/platformManager/backupConfig'),
  'BackupConfig'
);
const ScriptManager = asyncComponent(
  import('@/pages/platformManager/scriptManager'),
  'ScriptManager'
);
const ClusterInspection = asyncComponent(
  import('@/pages/platformManager/clusterInspection'),
  'ClusterInspection'
);


/* ----------- 二级菜单 ----------- */
/* ------ 产品 ------ */
const ServicePage = asyncComponent(
  import('@/pages/serviceStatus'),
  'ServicePage'
);
const HostDetailPage = asyncComponent(
  import('@/pages/hostDetail/hostDetail'),
  'HostDetailPage'
);
const AddHost = asyncComponent(import('@/pages/hostAdd/hostAdd'), 'AddHost');
const Deploy = asyncComponent(import('@/pages/deploy/deploy'), 'Deploy');
const HostStatus = asyncComponent(import('@/pages/hostsStatus'), 'HostStatus');
const IndexPage = asyncComponent(import('@/pages/indexPage'), 'IndexPage');

/* -- 诊断 -- */
const Log = asyncComponent(import('@/pages/diagnosis/log'), 'Log');
const Event = asyncComponent(import('@/pages/event/event'), 'Event');
const Config = asyncComponent(
  import('@/pages/diagnosis/config/config'),
  'Config'
);
const InspectionReport = asyncComponent(
  import('@/pages/diagnosis/inspectionReport/inspectionReport'),
  'InspectionReport'
);
const Backup = asyncComponent(
  import('@/pages/diagnosis/backup/backup'),
  'Backup'
);

/* ------ 集群 ------ */
const ClusterIndexPage = asyncComponent(
  import('@/pages/cluster/indexPage/indexPage'),
  'ClusterIndexPage'
);
const ClusterHost = asyncComponent(
  import('@/pages/cluster/host/host'),
  'ClusterHost'
);
const ClusterImageStore = asyncComponent(
  import('@/pages/cluster/imageStore/imageStore'),
  'ClusterImageStore'
);
const ClusterNamespace = asyncComponent(
  import('@/pages/cluster/namespace/namespace'),
  'ClusterNamespace'
);

export interface RouterConfItemType {
  path: string;
  redirect?: string;
  component?: any;
  layout?: any;
  children?: RouterConfItemType[];
}

export const RouterConf: RouterConfItemType[] = [
  /* ---------- 登录 ---------- */
  {
    path: '/login',
    component: Login,
  },
  /* ---------- 主页面 ---------- */
  {
    path: '/',
    layout: RootLayout,
    children: [
      /* ------ 运维中心 ------ */
      {
        path: '/opscenter',
        children: [
          { path: '/overview', component: IndexPage },
          { path: '/service', component: ServicePage },
          { path: '/hostAdd', component: AddHost },
          { path: '/hostDetail', component: HostDetailPage },
          { path: '/deploy', component: Deploy },
          { path: '/hoststatus', component: HostStatus },
          { path: '/diagnosis/log', component: Log },
          { path: '/diagnosis/event', component: Event },
          { path: '/diagnosis/config', component: Config },
          { path: '/diagnosis/inspectionReport', component: InspectionReport },
          { path: '/diagnosis/backup', component: Backup },
        ],
      },
      /* ------ 部署中心 ------ */
      {
        path: '/deploycenter',
        children: [
          { path: '/cluster/list', component: ClusterManagerPage },
          { path: '/cluster/create', component: ClusterManagerCreate },
          { path: '/cluster/create/:type', component: ClusterManagerEdit },
          { path: '/appmanage/products', component: ProductPackage },
          { path: '/appmanage/installs', component: Steps },
          { path: '/monitoring/dashboard', component: DashboardPage },
          { path: '/monitoring/dashdetail', component: DashboardDetail },
          { path: '/monitoring/alert', component: AlertManager },
          { path: '/monitoring/addAlert', component: AddAlertChannelPanel },
        ],
      },
      /* ------ 平台管理 ------ */
      {
        path: '/platform',
        children: [
          { path: '/backup', component: BackupConfig },
          { path: '/scriptManager', component: ScriptManager },
          { path: '/clusterInspection', component: ClusterInspection },
        ]
      },
      {
        path: '/deploycenter',
        children: [
          { path: '/cluster/list', component: ClusterManagerPage },
          { path: '/cluster/create', component: ClusterManagerCreate },
          { path: '/cluster/create/:type', component: ClusterManagerEdit },
          { path: '/appmanage/products', component: ProductPackage },
          { path: '/appmanage/installs', component: Steps },
          { path: '/monitoring/dashboard', component: DashboardPage },
          { path: '/monitoring/dashdetail', component: DashboardDetail },
          { path: '/monitoring/alert', component: AlertManager },
          { path: '/monitoring/addAlert', component: AddAlertChannelPanel },
        ],
      },

      /* -------- 集群详情 -------- */
      {
        path: '/deploycenter/cluster/detail',
        layout: ClusterLayout,
        children: [
          { path: '/index', component: ClusterIndexPage },
          { path: '/host', component: ClusterHost },
          { path: '/imagestore', component: ClusterImageStore },
          { path: '/namespace', component: ClusterNamespace },
          { path: '/deployed', component: DeployedApp },
          { path: '/history', component: DeployHistory },
          { path: '/patchHistory', component: PatchHistory },
          { path: '/echoList', component: CommandEchoIndex },
          { path: '/echoDetail', component: CommandEchoDetail },
        ],
      },
      {
        path: '/systemconfig',
        children: [
          { path: '/security', component: Security },
          { path: '/globalConfig', component: GlobalConfig },
        ],
      },
      /* ------ 用户中心 ------ */
      {
        path: '/usercenter',
        children: [
          { path: '/members', component: Members },
          { path: '/role', component: RoleManage },
          { path: '/role/view', component: RoleInfo },
          { path: '/selfinfo', component: SelfInfo },
          { path: '/audit', component: SecurityAudit },
        ],
      },
      {
        path: '/host',
        component: HostPage,
      },
      {
        path: '/nopermission',
        component: NoPermission,
      },
      {
        path: '*',
        component: NoPermission,
      },
    ],
  },
];
