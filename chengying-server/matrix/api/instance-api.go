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
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/api/impl"
	"github.com/kataras/iris/context"
)

var InstanceOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "instance",
	SubRoutes: []apibase.Route{
		{
			Path: "{agent_id:string}",
			SubRoutes: []apibase.Route{{
				//http://XXXXX/api/v2/instance/{agent_id:string}/start
				Path: "start",
				POST: impl.Start,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "实例启动[EasyMatrix API]",
					},
				},
			}, {
				//http://XXXXX/api/v2/instance/{agent_id:string}/stop
				Path: "stop",
				POST: impl.Stop,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "实例停止[EasyMatrix API]",
					},
				},
			}, {
				//http://XXXXX/api/v2/instance/{agent_id:string}/dt_agent_health_check
				Path: "dt_agent_health_check",
				POST: impl.HealthReport,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "instance health信息上报[EasyMatrix API]",
					},
				},
			}, {
				//http://XXXXX/api/v2/instance/{agent_id:string}/dt_agent_error
				Path: "dt_agent_error",
				POST: impl.ErrorReport,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "instance error信息上报[EasyMatrix API]",
					},
				},
			}, {
				//http://XXXXX/api/v2/instance/{agent_id:string}/dt_agent_host_resource
				Path: "dt_agent_host_resource",
				POST: impl.HostResourcesReport,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "host resources信息上报[EasyMatrix API]",
					},
				},
			}, {
				//http://XXXXX/api/v2/instance/{agent_id:string}/dt_agent_performance
				Path: "dt_agent_performance",
				POST: impl.PerformanceReport,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "agent performance信息上报[EasyMatrix API]",
					},
				},
			}},
		},
		{
			Path: "{deploy_uuid:string}",
			SubRoutes: []apibase.Route{
				{
					//http://XXXXX/api/v2/instance/{deploy_uuid:string}/list
					Path: "list",
					GET:  impl.InstanceRecord,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "实例部署记录列表，部署列表显示和进度查询[EasyMatrix API]",
						}},
				},
			},
		},
		{
			Path: "{update_uuid:string}",
			SubRoutes: []apibase.Route{
				{
					Middlewares: []context.Handler{
						apibase.CheckPermission1,
					},
					//http://XXXXX/api/v2/instance/{update_uuid:string}/listupdate
					Path: "listupdate",
					GET:  impl.InstanceUpdateRecord,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "实例部署记录列表，部署列表显示和进度查询[EasyMatrix API]",
						}},
				},
			},
		},
		{
			Path: "{id:string}",
			SubRoutes: []apibase.Route{
				{
					//http://XXXXX/api/v2/instance/{id:int(primary key)}/log
					Path: "log",
					GET:  impl.InstanceServiceLog,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "返回实例下日志文件或者目录[EasyMatrix API]",
						}},
				},
				{
					//http://XXXXX/api/v2/instance/{id:int(primary key)}/config
					Path: "config",
					GET:  impl.InstanceServiceConfig,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "返回实例相关配置文件内容[EasyMatrix API]",
						}},
				},
				{
					//http://XXXXX/api/v2/instance/{id:int(primary key)}/event
					Path: "event",
					GET:  impl.InstanceEvent,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "返回实例相关事件内容[EasyMatrix API]",
						}},
				},
			},
		},
		{
			Path: "product",
			SubRoutes: []apibase.Route{
				{
					//http://XXXXX/api/v2/instance/prodcut/{product_name:string}/service/{service_name:string}
					Path: "{product_name:string}",
					SubRoutes: []apibase.Route{
						{
							Path: "service",
							SubRoutes: []apibase.Route{
								{
									Path: "{service_name:string}",
									GET:  impl.InstanceBelongService,
									Docs: apibase.Docs{
										GET: &apibase.ApiDoc{
											Name: "返回服务列表[EasyMatrix API]",
										}},
									SubRoutes: []apibase.Route{
										{
											Path: "alert",
											GET:  impl.InstanceServiceAlert,
											Docs: apibase.Docs{
												GET: &apibase.ApiDoc{
													Name: "返回服务告警列表[EasyMatrix API]",
												}},
										},
										{
											Path: "healthCheck",
											GET:  impl.InstanceServiceHealthCheck,
											Docs: apibase.Docs{
												GET: &apibase.ApiDoc{
													Name: "返回健康检查列表[EasyMatrix API]",
												}},
											SubRoutes: []apibase.Route{
												{
													Path: "setAutoexecSwitch",
													POST: impl.SetHealthCheckAutoexecSwitch,
													Docs: apibase.Docs{
														GET: &apibase.ApiDoc{
															Name: "健康检查定时执行开关[EasyMatrix API]",
														}},
												},
												{
													Path: "manualExecution",
													POST: impl.InstanceServiceHealthCheckExec,
													Docs: apibase.Docs{
														GET: &apibase.ApiDoc{
															Name: "手动执行健康检查[EasyMatrix API]",
														}},
												},
											},
										},
									},
								},
							},
						},
						{
							Path: "service_list",
							GET:  impl.ServiceList,
							Docs: apibase.Docs{
								GET: &apibase.ApiDoc{
									Name: "根据产品名返回服务列表[EasyMatrix API]",
								}},
						},
						{
							Path: "group_list",
							GET:  impl.GroupList,
							Docs: apibase.Docs{
								GET: &apibase.ApiDoc{
									Name: "根据产品名返回group列表[EasyMatrix API]",
								}},
						},
					},
				},
			},
		},
		{
			Path: "test",
			SubRoutes: []apibase.Route{
				{
					//http://XXXXX/api/v2/instance/test/install
					Path: "install",
					POST: impl.InstancerControllTest,
					Docs: apibase.Docs{
						POST: &apibase.ApiDoc{
							Name: "实例部署测试[EasyMatrix API]",
						}},
				},
			},
		},
		{
			Path: "{instance_id:string}",
			SubRoutes: []apibase.Route{
				{
					//http://XXXXX/api/v2/instance/{instance_id:string}/logfiles?type=text|zip
					Path: "logfiles",
					GET:  impl.ListLogFiles,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "list {instance_id} $(dir) 下的文件列表[EasyMatrix API]",
							Returns: []apibase.ApiReturnGroup{{
								Fields: apibase.ResultFields{
									"$.code": apibase.ApiReturn{"int", "True Or False"},
									"$.msg":  apibase.ApiReturn{"string", "message"},
									"$.data": apibase.ApiReturn{"string", "data"},
								},
							}},
						},
					},
				},
				{
					//http://XXXXX/api/v2/instance/{instance_id:string}/logmore?logfile=$(path)&action=lastest
					Path: "logmore",
					GET:  impl.PreviewLogFile,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "预览 {instance_id} $(path) 指定文件[EasyMatrix API]",
							Returns: []apibase.ApiReturnGroup{{
								Fields: apibase.ResultFields{
									"$.code": apibase.ApiReturn{"int", "True Or False"},
									"$.msg":  apibase.ApiReturn{"string", "message"},
									"$.data": apibase.ApiReturn{"string", "data"},
								},
							}},
						},
					},
				},
				{
					//http://XXXXX/api/v2/instance/{instance_id:string}/logdown?logfile=$(path)
					Path: "logdown",
					GET:  impl.DownloadLogFile,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "下载 {instance_id} $(path) 指定文件[EasyMatrix API]",
						},
					},
				},
			},
		},
		{
			Path: "event",
			SubRoutes: []apibase.Route{
				{
					Path: "typeList",
					GET:  impl.EventTypeList,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取事件类型列表 [EasyMatrix API]",
						},
					},
				},
				{
					Path: "list",
					GET:  impl.EventList,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取事件列表 [EasyMatrix API]",
						},
					},
				},
				{
					Path: "coordinate",
					GET:  impl.EventTimeCoordinate,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "事件-时间折线图 [EasyMatrix API]",
						},
					},
				},
				{
					Path: "statistics",
					GET:  impl.EventStatistics,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "事件-时间折线图 [EasyMatrix API]",
						},
					},
				},
				{
					Path: "{product_or_service:string}",
					SubRoutes: []apibase.Route{{
						Path: "rank",
						GET:  impl.EventRank,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "(组件或服务)事件发生次数排行 [EasyMatrix API]",
							},
						},
					}},
				},
			},
		},
		{
			Path: "reload",
			POST: impl.DiscoveryReload,
		},
	},
}
