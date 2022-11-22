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
	SubRoutes: []apibase.Route{{
		Path: "kubernetes",
		SubRoutes: []apibase.Route{{
			Path: "imageStore",
			SubRoutes: []apibase.Route{{
				//http:://xxxx/api/v2/cluster/kubernetes/imageStore/create
				Path: "create",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				POST: impl.CreateImageStore,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "创建k8s镜像仓库[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/kubernetes/imageStore/update
				Path: "update",
				POST: impl.UpdateImageStore,
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "编辑k8s镜像仓库[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/kubernetes/imageStore/setDefault
				Path: "setDefault",
				POST: impl.SetDefaultImageStore,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "设置k8s集群默认镜像仓库[EasyMatrix API]",
					},
				},
			}, {
				Path: "{cluster_id:int}",
				SubRoutes: []apibase.Route{{
					//http:://xxxx/api/v2/cluster/kubernetes/imageStore/{cluster_id}/clusterInfo
					Path: "clusterInfo",
					GET:  impl.GetImageStoreInfoByClusterId,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "通过clusterId获取k8s镜像仓库[EasyMatrix API]",
						},
					},
				}, {
					//http:://xxxx/api/v2/cluster/kubernetes/imageStore/{cluster_id}/checkDefault
					Path: "checkDefault",
					GET:  impl.CheckDefaultImageStore,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "检查集群是否设置默认仓库[EasyMatrix API]",
						},
					},
				}},
			}, {
				Path: "{store_id:int}",
				SubRoutes: []apibase.Route{{
					//http:://xxxx/api/v2/cluster/kubernetes/imageStore/{store_id}/info
					Path: "info",
					GET:  impl.GetImageStoreInfoById,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "通过仓库Id获取k8s镜像仓库[EasyMatrix API]",
						},
					},
				}},
			}, {
				//http:://xxxx/api/v2/cluster/kubernetes/imageStore/delete
				Path: "delete",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				POST: impl.DeleteImageStore,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "通过仓库Id删除k8s镜像仓库[EasyMatrix API]",
					},
				},
			}},
		}, {
			//http:://xxxx/api/v2/cluster/kubernetes/rketemplate?version=V1.14.3&network_plugin=flannel
			Path: "rketemplate",
			GET:  impl.GetRkeTemplate,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "获取rke配置模版[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/kubernetes/create
			Path: "create",
			Middlewares: []context.Handler{
				apibase.CheckPermission3,
			},
			POST: impl.CreateK8sCluster,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "创建k8s集群[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/kubernetes/update
			Path: "update",
			Middlewares: []context.Handler{
				apibase.CheckPermission3,
			},
			POST: impl.UpdateK8sCluster,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "编辑k8s集群[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/kubernetes/available
			Path: "available",
			GET:  impl.GetK8sAvailable,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取k8s支持列表[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/kubernetes/hosts
			Path: "hosts",
			GET:  impl.GetK8sClusterHostList,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取k8s集群主机列表[EasyMatrix API]",
				},
			},
		}, {
			//http:://xxxx/api/v2/cluster/kubernetes/installCmd
			Path: "installCmd",
			GET:  impl.GetK8sClusterImportCmd,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取导入命令[EasyMatrix API]",
				},
			},
		}, {
			Path: "{cluster_id:int}",
			SubRoutes: []apibase.Route{{
				//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/installLog
				Path: "installLog",
				GET:  impl.GetK8sClusterInstallLog,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获取k8s集群安装日志[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/info
				Path: "info",
				GET:  impl.GetK8sClusterInfo,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获取k8s集群信息[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/delete
				Path: "delete",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				POST: impl.DeleteK8sCluster,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "删除k8s集群[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/overview
				Path: "overview",
				GET:  impl.GetK8sClusterOverView,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获取k8s集群总览信息[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/performance
				Path: "performance",
				GET:  impl.GetK8sClusterPerformance,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获取K8S集群性能趋势[EasyMatrix API]",
					},
				},
			}, {
				Path: "namespace",
				SubRoutes: []apibase.Route{{
					//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/namespace/list
					Path: "list",
					GET:  impl.GetK8sClusterNameSpaceList,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取k8s集群namespace list[EasyKube API]",
						},
					},
				}, {
					//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/namespace/create
					Path: "create",
					POST: impl.CreateK8sClusterNamespace,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "创建k8s集群namespace[EasyKube API]",
						},
					},
				}, {
					Path: "{namespace_name:string}",
					SubRoutes: []apibase.Route{{
						Path: "product",
						SubRoutes: []apibase.Route{{
							Path: "{pid:int}",
							SubRoutes: []apibase.Route{{
								//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/namespace/{namespace_name}/product/{pid}/depends
								Path: "depends",
								GET:  impl.GetK8sClusterProductDepends,
								Docs: apibase.Docs{
									GET: &apibase.ApiDoc{
										Name: "获得k8s集群指定namespace下产品依赖[EasyKube API]",
									},
								},
							}, {
								//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/namespace/{namespace_name}/product/{pid}/deploy
								Path: "deploy",
								POST: impl.DeployK8sProduct,
								Docs: apibase.Docs{
									GET: &apibase.ApiDoc{
										Name: "指定k8s集群" +
											"ce部署产品[EasyKube API]",
									},
								},
							}, {
								Middlewares: []context.Handler{
									apibase.CheckPermission1,
								},
								//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/namespace/{namespace_name}/product/{pid}/stop
								Path: "stop",
								POST: impl.StopUndeployingK8sProduct,
								Docs: apibase.Docs{
									GET: &apibase.ApiDoc{
										Name: "停止卸载k8s产品包[EasyKube API]",
									},
								},
							}, {
								//http:://xxxx/api/v2/cluster/kubernetes/{cluster_id}/namespace/{namespace_name}/product/{pid}/installLog
								Path: "installLog",
								GET:  impl.GetProductInstallLog,
								Docs: apibase.Docs{
									GET: &apibase.ApiDoc{
										Name: "获取产品部署日志[EasyKube API]",
									},
								},
							}},
						}},
					}},
				}},
			}},
		}, {
			//http:://xxxx/api/v2/cluster/kubernetes/listwatch
			Path: "listwatch",
			SubRoutes: []apibase.Route{
				{
					Path: "events",
					POST: impl.K8sProductListWatch,
					Docs: apibase.Docs{
						POST: &apibase.ApiDoc{
							Name: "提交k8s产品部署组件事件信息[EasyMatrix API]",
						},
					},
				},
			},
		},
		}}, {
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
