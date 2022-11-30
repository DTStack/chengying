// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package api

import (
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/api/impl"
	"github.com/kataras/iris/context"
)

var ClusterEasyMatrixAPIRoutes = apibase.Route{
	Path: "cluster",
	SubRoutes: []apibase.Route{
		{
			Path: "hosts",
			SubRoutes: []apibase.Route{{
				//http:://xxxx/api/v2/cluster/hosts/create
				Path: "create",
				POST: impl.CreateHostCluster,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "创建主机集群[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/hosts/update
				Path: "update",
				POST: impl.UpdateHostCluster,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "编辑主机集群[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/hosts/alert
				Path: "alert",
				GET:  impl.GetHostClusterAlert,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获取主机集群主机报警信息[EasyMatrix API]",
					},
				},
			}, {
				Path: "{cluster_id:int}",
				SubRoutes: []apibase.Route{{
					//http:://xxxx/api/v2/cluster/hosts/{cluster_id}/info
					Path: "info",
					GET:  impl.GetHostClusterInfo,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取主机集群信息[EasyMatrix API]",
						},
					},
				}, {
					//http:://xxxx/api/v2/cluster/hosts/{cluster_id}/delete
					Path: "delete",
					POST: impl.DeleteHostCluster,
					Docs: apibase.Docs{
						POST: &apibase.ApiDoc{
							Name: "删除主机集群[EasyMatrix API]",
						},
					},
				}, {
					//http:://xxxx/api/v2/cluster/hosts/{cluster_id}/overview
					Path: "overview",
					GET:  impl.GetHostClusterOverView,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取主机集群总览信息[EasyMatrix API]",
						},
					},
				}, {
					//http:://xxxx/api/v2/cluster/hosts/{cluster_id}/performance
					Path: "performance",
					GET:  impl.GetHostClusterPerformance,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取主机集群性能趋势[EasyMatrix API]",
						},
					},
				}},
			}, {
				//http:://xxxx/api/v2/cluster/hosts/hosts
				Path: "hosts",
				GET:  impl.GetHostClusterHostList,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获取主机集群主机列表[EasyMatrix API]",
					},
				},
			},
				{
					//http:://xxxx/api/v2/cluster/hosts/role
					Path: "role",
					POST: impl.EditRole,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "修改主机角色",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/role_list
					Path: "role_list",
					GET:  impl.RoleList,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取主机角色列表",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/role_rename
					Path: "role_rename",
					POST: impl.RoleRename,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "修改角色名",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/role_info
					Path: "role_info",
					GET:  impl.RoleInfo,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取集群下机器的角色信息",
						},
					},
				},

				{
					//http:://xxxx/api/v2/cluster/hosts/role_delete
					Path: "role_delete",
					POST: impl.RoleDelete,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取主机角色列表",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/role_add
					Path: "role_add",
					POST: impl.RoleAdd,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "角色添加",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/auto_orchestration
					Path: "auto_orchestration",
					POST: impl.AutoOrchestration,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "自动编排",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/auto_svcgroup
					Path: "auto_svcgroup",
					POST: impl.AutoSvcGroup,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取实时编排结果",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/auto_deploy
					Path: "auto_deploy",
					POST: impl.AutoDeploy,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "自动部署",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/auto_deploy_cancel
					Path: "auto_deploy_cancel",
					POST: impl.AutoDeployCancel,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "取消自动部署",
						},
					},
				},
				{
					//http:://xxxx/api/v2/cluster/hosts/orchestration_history
					Path: "orchestration_history",
					GET:  impl.OrchestrationHistory,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取历史安装记录",
						},
					},
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/hostgroups?type=&id=
			Path: "hostgroups",
			GET:  impl.GetHostGroups,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取集群主机分组[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/list?type=&sort-by=id&sort-dir=desc&limit=10&start=0
			Path: "list",
			GET:  impl.GetClusterList,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取集群列表[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/products
			Path: "products",
			GET:  impl.GetClusterProductList,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取集群产品列表[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/productsInfo
			Path: "productsInfo",
			GET:  impl.GetClusterProductsInfo,
			Middlewares: []context.Handler{
				apibase.CheckPermission3,
			},
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取集群产品详情[EasyMatrix API]",
				},
			},
		}, {
			//http://xxxx/api/v2/cluster/restartServices
			Path: "restartServices",
			GET:  impl.GetRestartServices,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取需要重启服务",
				},
			},
		}, {
			//http://xxxx/api/v2/cluster/currExecCount?clusterId=1
			Path: "currExecCount",
			GET:  impl.GetCurrentExecCount,
			Middlewares: []context.Handler{
				apibase.CheckPermission3,
			},
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取集群当前正在运行的命令数",
				},
			},
		}, {
			//http://xxxx/api/v2/cluster/orderList?clusterId=1&operationType=1&object=DTBase&status=success&startTime=2021-05-16 13:13:12&endTime=2021-05-16 13:13:12&page=1&pageSize=10
			Path: "orderList",
			GET:  impl.OrderList,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "查看所有命令",
				},
			},
		},
		{
			//http://xxxx/api/v2/cluster/OrderDetail
			Path: "orderDetail",
			GET:  impl.OrderDetail,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "查看命令详情",
				},
			},
		},
		{
			Path: "seqReport",
			//http://xxxx/api/v2/cluster/seqReport
			POST: impl.SeqReport,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "汇报 seq",
				},
			},
		},
		{
			Path: "isShowLog",
			//http://xxxx/api/v2/cluster/isShowLog?seq=123
			GET: impl.IsShowLog,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "判断该 seq 是否进行记录内容与日志",
				},
			},
		},
		{
			//http://xxxx/api/v2/cluster/shellStatusReport
			Path: "shellStatusReport",
			POST: impl.ShellStatusReport,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "汇报脚本状态",
				},
			},
		},
		{
			//http://xxxx/api/v2/cluster/showShellLog?execId=222
			Path: "showShellLog",
			GET:  impl.ShowShellLog,
			Middlewares: []context.Handler{
				apibase.CheckPermission3,
			},
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "查看脚本日志",
				},
			},
		},

		{
			//http://xxxx/api/v2/cluster/downLoadShellLog?execId=222
			Path: "downLoadShellLog",
			GET:  impl.DownLoadShellLog,
			Middlewares: []context.Handler{
				apibase.CheckPermission3,
			},
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "下载脚本日志",
				},
			},
		},

		{
			//http://xxxx/api/v2/cluster/downLoadShellContent?execId=222
			Path: "downLoadShellContent",
			GET:  impl.DownLoadShellContent,
			Middlewares: []context.Handler{
				apibase.CheckPermission3,
			},
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "下载脚本内容",
				},
			},
		},
		{
			//http://xxxx/api/v2/cluster/previewShellContent?execId=222
			Path: "previewShellContent",
			GET:  impl.PreviewShellContent,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "预览脚本内容",
				},
			},
		},

		{
			//http://xxxx/api/v2/cluster/listObjectValue?clusterId=1
			Path: "listObjectValue",
			GET:  impl.ListObjectValue,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "list object value",
				},
			},
		},
	},
}
