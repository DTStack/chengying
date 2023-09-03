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

var ProductOperationEasyMatrixAPIRoutes = apibase.Route{
	//http://xxxx/api/v2/product
	Path: "product",
	GET:  impl.ProductInfo,
	Docs: apibase.Docs{
		GET: &apibase.ApiDoc{
			Name: "所有产品信息接口[EasyMatrix API]",
		},
	},
	SubRoutes: []apibase.Route{{
		Path: "upload",
		POST: impl.Upload,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "上传产品包接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "uploadAsync",
		POST: impl.UploadAsync,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "根据链接异步上传产品包接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "uploadSync",
		POST: impl.UploadSync,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "根据链接同步上传产品包接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "check_param",
		POST: impl.CheckAvailableLink,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "校验参数接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "in_progress",
		GET:  impl.GetUploadingProducts,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取正在上传的产品包列表接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "cancel_upload",
		POST: impl.CancelUpload,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "停止正在上传的产品包接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "backup",
		GET:  impl.GetBackupPackage,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "获取集群下所有备份包",
			},
		},
		SubRoutes: []apibase.Route{{
			Path: "setconfig",
			Middlewares: []context.Handler{
				apibase.CheckPermission1,
			},
			POST: impl.SetClusterBackupPATH,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "设置备份路径接口[EasyMatrix API]",
				},
			},
		}, {
			Path: "getconfig",
			Middlewares: []context.Handler{
				apibase.CheckPermission1,
			},
			GET: impl.GetClusterBackupPATH,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "查询备份路径接口[EasyMatrix API]",
				},
			},
		}},
	}, {
		Path: "clean",
		//Middlewares: []context.Handler{
		//	apibase.CheckPermission3,
		//},
		POST: impl.CleanBackupPackage,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "清理备份包",
			},
		},
	}, {
		Path: "patchpath",
		Middlewares: []context.Handler{
			apibase.CheckPermission3,
		},
		GET: impl.PatchPath,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "产品包组件目录获取接口[EasyMatrix API]",
				Query: apibase.ApiParams{
					"product_id": apibase.ApiParam{"string", "产品包名称", "", true},
				},
			},
		},
	}, {
		Path: "patchupdate",
		POST: impl.PatchUpdate,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "产品包组件补丁更新接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "patchupload",
		POST: impl.PatchUpload,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "产品包组件补丁包上传接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "parentProduct",
		GET:  impl.ParentProductInfo,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "父类产品列表接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "{pid:int}",
		SubRoutes: []apibase.Route{
			{
				Path: "start",
				GET:  impl.ProductStart,
			}, {
				Path: "stop",
				GET:  impl.ProductStop,
			}, {
				Path: "unchecked_services",
				GET:  impl.ProductUncheckedServices,
			},
		},
	}, {
		Path: "productList",
		GET:  impl.ProductList,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "产品包列表接口[EasyMatrix API]",
			},
		},
	}, {
		Path: "productName",
		GET:  impl.ProductName,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "产品组件包名接口[EasyMatrix API]",
			},
		},
	}, {
		//http:://xxxx/api/v2/product/{product_name}
		Path: "{product_name:string}",
		GET:  impl.ProductInfo,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "该产品版本信息接口[EasyMatrix API]",
			},
		},
		SubRoutes: []apibase.Route{
			{
				Path: "service",
				SubRoutes: []apibase.Route{{
					Path: "{service_name:string}",
					SubRoutes: []apibase.Route{{
						Path: "get_ip",
						GET:  impl.GetIP,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "设置服务IP接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "set_ip",
						POST: impl.SetIP,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "设置服务IP接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "modify_schema_field",
						POST: impl.ModifySchemaField,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "修改schema字段接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "modify_schema_field_devops",
						POST: impl.ModifySchemaFieldForDevOps,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "修改schema字段接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "modify_schema_field_batch",
						POST: impl.ModifySchemaFieldBatch,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "批量schema字段接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "reset_schema_field",
						Middlewares: []context.Handler{
							apibase.CheckPermission1,
						},
						POST: impl.ResetSchemaField,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "重置schema字段接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "reset_multi_schema_field",
						POST: impl.ResetSchemaMultiField,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "重置schema多字段接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "hosts",
						GET:  impl.AvailableHosts,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "get available hosts filter by productName and serviceName[EasyMatrix API]",
							},
						},
					}, {
						Path: "selected_hosts",
						//Middlewares: []context.Handler{
						//	apibase.CheckPermission1,
						//},
						GET: impl.SelectedHosts,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "get selected hosts by productName and serviceName[EasyMatrix API]",
							},
						},
					}, {
						Path: "modifyAll",
						// HIDE Temporary For DevOps
						//Middlewares: []context.Handler{
						//	apibase.CheckPermission3,
						//},
						POST: impl.ModifyAll,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "修改所有服务IP接口[EasyMatrix API]",
							},
						},
					}, {
						Path: "modifyMultiAll",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						POST: impl.ModifyAllSchemaMultiField,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "修改服务配置项配置多个值",
							},
						},
					}, {
						Path: "modifyMultiSingleField",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						POST: impl.ModifyMultiField,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "修改单个服务配置项配置多个值",
							},
						},
					}, {
						Path: "operate_switch",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						POST: impl.OperateSwitch,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "操作开关",
							},
						},
					}, {
						Path: "switch/record",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						GET: impl.CheckSwitchRecord,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获取开关操作记录",
							},
						},
					},
						{
							Path: "switch/detail",
							Middlewares: []context.Handler{
								apibase.CheckPermission3,
							},
							GET: impl.GetSwitchDetail,
							Docs: apibase.Docs{
								GET: &apibase.ApiDoc{
									Name: "获取开关操作详情",
								},
							},
						}, {
							Path: "extention_operation",
							Middlewares: []context.Handler{
								apibase.CheckPermission3,
							},
							GET: impl.DoExtentionOperation,
							Docs: apibase.Docs{
								GET: &apibase.ApiDoc{
									Name: "开关触发操作",
								},
							},
						}},
				}},
			},
			{
				//http:://xxxx/api/v2/product/{product_name}/current
				Path: "current",
				GET:  impl.Current,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获得当前运行的产品版本详细信息[EasyMatrix API]",
					},
				},
				//SubRoutes: []apibase.Route{
				//	{
				//		//http:://xxxx/api/v2/product/{product_name}/current/servicesStatus
				//		Path: "servicesStatus",
				//		GET:  impl.ServicesStatus,
				//		Docs: apibase.Docs{
				//			GET: &apibase.ApiDoc{
				//				Name: "获得产品服务状态[EasyMatrix API]",
				//			},
				//		},
				//	},
				//},
			},
			{
				//http:://xxxx/api/v2/product/{product_name}/history
				Path: "history",
				GET:  impl.History,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获得产品发布历史记录[EasyMatrix API]",
					},
				},
			},
			{
				//http:://xxxx/api/v2/product/{product_name}/updatehistory
				Path: "updatehistory",
				GET:  impl.UpdateHistory,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "获得产品发布后的组件更新历史记录[EasyMatrix API]",
					},
				},
			}, {
				//http:://xxxx/api/v2/product/{product_name}/serviceUpdate
				Path: "serviceUpdate",
				POST: impl.ServiceUpdate,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "批量设置产品服务信息[EasyMatrix API]",
						Body: apibase.ApiParams{
							"$.field_path": apibase.ApiParam{"string", "schema,如run_user", "", true},
							"$.field":      apibase.ApiParam{"interface", "value", "", true},
						},
					},
				},
			}, {
				Path: "autoTest",
				SubRoutes: []apibase.Route{
					{
						//http:://xxxx/api/v2/product/{product_name}/autoTest/start
						Path: "start",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						POST: impl.StartAutoTest,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "启动当前产品自动化测试",
							},
						},
					},
					{
						//http:://xxxx/api/v2/product/{product_name}/autoTest/history
						Path: "history",
						GET:  impl.AutoTestHistory,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获取当前产品最近的一次自动化测试记录",
							},
						},
					},
				},
			},
			{
				Path: "haRole",
				SubRoutes: []apibase.Route{
					{
						//http:://xxxx/api/v2/product/{product_name}/haRole/${service_name}
						Path: "{service_name:string}",
						GET:  impl.HaRole,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获得产品中ha的角色信息[EasyMatrix API]",
							},
						},
					},
				},
			},
			{
				Path: "version",
				SubRoutes: []apibase.Route{{
					Middlewares: []context.Handler{
						apibase.CheckPermission3,
					},
					//http:://xxxx/api/v2/product/{product_name}/version/{product_version}
					Path:   "{product_version:string}",
					DELETE: impl.ProductDelete,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "该产品该版本信息/删除接口[EasyMatrix API]",
						},
					},
				}, {
					Path: "{product_version:string}",
					GET:  impl.ProductInfo,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "该产品该版本信息/删除接口[EasyMatrix API]",
						},
					},
					SubRoutes: []apibase.Route{{
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/deployDevOps
						Path: "deployDevOps",
						POST: impl.DeployForDevOps,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "部署该产品该版本接口[EasyMatrix API]",
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/deploy
						Path: "deploy",
						POST: impl.Deploy,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "部署该产品该版本接口[EasyMatrix API]",
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/deployLogs
						Path: "deployLogs",
						GET:  impl.DeployLogs,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获取部署过程日志[EasyMatrix API]",
							},
						},
					}, {
						Middlewares: []context.Handler{
							apibase.CheckPermission1,
						},
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/undeploy
						Path: "undeploy",
						POST: impl.Undeploy,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "卸载该产品该版本接口[EasyMatrix API]",
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/cancel
						Path: "cancel",
						POST: impl.Cancel,
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service
						Path: "service",
						GET:  impl.Service,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获得当前产品的服务列表[EasyMatrix API]",
							},
						},
						SubRoutes: []apibase.Route{
							{
								//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceTree
								Path: "{service_name:string}",
								SubRoutes: []apibase.Route{
									{
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceTree
										Path: "serviceTree",
										GET:  impl.ServiceTree,
										Docs: apibase.Docs{
											GET: &apibase.ApiDoc{
												Name: "获得产品指定服务的包目录中文本文件列表，.sh/.conf/.properties/.sql[EasyMatrix API]",
											},
										},
									}, {
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceFile
										Path: "serviceFile",
										GET:  impl.ServiceFile,
										Docs: apibase.Docs{
											GET: &apibase.ApiDoc{
												Name: "获得产品指定服务的包目录中目标文件内容[EasyMatrix API]",
											},
										},
									}, {
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceGraphy
										Path: "serviceGraphy",
										POST: impl.ServiceGraphy,
										Docs: apibase.Docs{
											POST: &apibase.ApiDoc{
												Name: "服务自动编排[EasyMatrix API]",
											},
										},
									}, {
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceConfigFiles
										Path: "serviceConfigFiles",
										GET:  impl.ServiceConfigFiles,
										Docs: apibase.Docs{
											GET: &apibase.ApiDoc{
												Name: "获取服务所有配置模版文件[EasyMatrix API]",
											},
										},
									}, {
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceConfigFile?file=example-config.yml
										Path: "serviceConfigFile",
										GET:  impl.ServiceConfigFile,
										Docs: apibase.Docs{
											GET: &apibase.ApiDoc{
												Name: "获取服务所有配置模版文件[EasyMatrix API]",
												Query: apibase.ApiParams{
													"file": apibase.ApiParam{"string", "配置文件名", "", true},
												},
											},
										},
									}, {
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/configUpdate
										Path: "configUpdate",
										POST: impl.ConfigUpdate,
										Docs: apibase.Docs{
											POST: &apibase.ApiDoc{
												Name: "更新服务配置模版[EasyMatrix API]",
												Body: apibase.ApiParams{
													"$.file":    apibase.ApiParam{"string", "配置文件名", "", true},
													"$.content": apibase.ApiParam{"string", "配置文件内容", "", true},
													"$.values":  apibase.ApiParam{"map", "新增配置项目和默认值", "", true},
													"$.deleted": apibase.ApiParam{"string", "删除配置项", "", false},
												},
											},
										},
									}, {
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceConfig?file=redis.conf
										Path: "serviceConfig",
										GET:  impl.ServiceConfig,
										Docs: apibase.Docs{
											GET: &apibase.ApiDoc{
												Name: "获取服务配置信息[EasyMatrix API]",
												Query: apibase.ApiParams{
													"file": apibase.ApiParam{"string", "配置文件名", "", false},
												},
											},
										},
									}, {
										//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/service/{service_name}/serviceConfigDiff?file=redis.conf&ip=172.16.82.176
										Path: "serviceConfigDiff",
										GET:  impl.ServiceConfigDiff,
										Docs: apibase.Docs{
											GET: &apibase.ApiDoc{
												Name: "服务配置文件对比[EasyMatrix API]",
												Query: apibase.ApiParams{
													"file": apibase.ApiParam{"string", "配置文件名", "", false},
													"ip":   apibase.ApiParam{"string", "主机地址", "", false},
												},
											},
										},
									},
								},
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/serviceGroup
						Path: "serviceGroup",
						GET:  impl.ServiceGroup,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获得当前产品的服务组信息[EasyMatrix API]",
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/serviceGroupFile
						Path: "serviceGroupFile",
						GET:  impl.ServiceGroupFile,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获得当前产品服务关联文件的config[EasyMatrix API]",
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/serviceGraphy?unchecked_services=***
						Path: "serviceGraphy",
						POST: impl.ServicesGraphy,
						Docs: apibase.Docs{
							POST: &apibase.ApiDoc{
								Name: "产品自动编排[EasyMatrix API]",
								Query: apibase.ApiParams{
									"unchecked_services": apibase.ApiParam{"string", "未勾选的服务，逗号间隔", "", false},
								},
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/group_list
						Path: "group_list",
						GET:  impl.GetProductGroupList,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获得产品的组信息[EasyMatrix API]",
							},
						},
					}, {
						//http:://xxxx/api/v2/product/{product_name}/version/{product_version}/upgrade_candidate
						Path: "upgrade_candidate",
						GET:  impl.GetUpgradeCandidateList,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "获得产品的升级候选人信息[EasyMatrix API]",
							},
						},
					}},
				},
				},
			}, {
				//http:://xxxx/api/v2/product/{product_name}/checkMysqlAddr
				Path: "checkMysqlAddr",
				POST: impl.CheckMysqlAddr,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "MySQL地址校验[EasyMatrix API]",
					},
				},
			}, {
				//http://xxxx/api/v2/product/{product_name}/backupDb
				Path: "backupDb",
				POST: impl.BackupDatabase,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "备份数据库",
					},
				},
			}, {
				// http://xxxx/api/v2/product/{product_name}/currentInfo
				Path: "currentInfo",
				POST: impl.ProductCurrentInfo,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "获取老版本服务主机编排与配置信息",
					},
				},
			}, {
				// http://xxxx/api/v2/product/{product_name}/saveUpgrade
				Path: "saveUpgrade",
				POST: impl.SaveUpgrade,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "保存升级历史记录",
					},
				},
			}, {
				// http://xxxx/api/v2/product/{product_name}/rollbackVersions
				Path: "rollbackVersions",
				POST: impl.RollbackVersions,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "获取回滚版本列表",
					},
				},
			}, {
				// http://xxxx/api/v2/product/{product_name}/backupTimes
				Path: "backupTimes",
				POST: impl.BackupTimes,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "获取对应备份库记录时间",
					},
				},
			}, {
				// http://xxxx/api/v2/product/{product_name}/rollback
				Path: "rollback",
				POST: impl.Rollback,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "回滚操作",
					},
				},
			},
		},
	}, {
		//http:://xxx/api/v2/product/status
		Path: "status",
		GET:  impl.ProductStatus,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获得所有子产品服务状态接口[EasyMatrix API]",
			},
		},
	}, {
		//http://xxx/api/v2/product/anomalyService
		Path: "anomalyService",
		GET:  impl.GetProductAnomalyService,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取异常服务列表",
			},
		},
	}, {
		//http:://xxxx/api/v2/product/configAlterGroups
		Path: "configAlterGroups",
		GET:  impl.ConfigAlterGroups,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "返回配置修改的组信息[EasyMatrix API]",
			},
		},
	}, {
		//http:://xxxx/api/v2/product/configAlteration
		Path: "configAlteration",
		GET:  impl.ConfigAlteration,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "右侧服务组配置变更信息查询[EasyMatrix API]",
			},
		},
	}, {
		// NOT USED| DEPRECATED
		//http:://xxxx/api/v2/product/configAlterAll
		Path: "configAlterAll",
		GET:  impl.ConfigAlterAll,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "返回配置修改组的所有信息[EasyMatrix API]",
			},
		},
	}, {
		// FOR DEVOPS
		//http:://xxxx/api/v2/product/getServiceConfig?product_name=&clusterId=&pid=&configPath=
		Path: "getServiceConfig",
		GET:  impl.GetServiceConfig,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "返回服务配置信息[EasyMatrix API]",
			},
		},
	}, {
		//http:://xxxx/api/v2/product/product_name_list
		//GET Product NAME
		Path: "product_name_list",
		GET:  impl.GetProductNameList,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "返回组件名称列表接口[EasyMatrix API]",
			},
		},
	}, {
		//http:://xxxx/api/v2/product/deployCondition
		Path: "deployCondition",
		POST: impl.CheckDeployCondition,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "检查当前组件的部署条件[EasyMatrix API]",
			},
		},
	}},
}
