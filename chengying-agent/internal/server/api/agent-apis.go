/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	apibase "easyagent/go-common/api-base"
	"fmt"

	"easyagent/internal/server/api/impl"
)

var agentCommonDoc = apibase.ResultFields{
	"$[*].id": apibase.ApiReturn{"string", "Agent ID"},
}

func agentListInfo() apibase.ResultFields {
	n := apibase.ResultFields{}
	for k, v := range agentCommonDoc {
		n["$.info[*]."+k] = v
	}
	n["$.total"] = apibase.ApiReturn{"int", "记录总数"}
	return n
}

func extendAgentInfo(extended apibase.ResultFields) apibase.ResultFields {
	for k, v := range agentCommonDoc {
		k = "$." + k
		if _, existed := extended[k]; existed {
			continue
		}
		extended[k] = v
	}
	return extended
}

func progressReturnDoc(progName string) apibase.Docs {
	return apibase.Docs{
		GET: &apibase.ApiDoc{
			Name: fmt.Sprintf("返回%s的进度", progName),
			Returns: []apibase.ApiReturnGroup{{
				Fields: apibase.ResultFields{
					"$.operation_seq": apibase.ApiReturn{"int", "操作序列号"},
					"$.op_time":       apibase.ApiReturn{"date", "操作开始时间"},
					"$.op_name":       apibase.ApiReturn{"string", "操作名称"},
					"$.progress":      apibase.ApiReturn{"float", "当前进度百分比"},
					"$.op_result":     apibase.ApiReturn{"int", "操作执行状态"},
					"$.ret_msg":       apibase.ApiReturn{"string", "操作返回内容"},
					"$.msg":           apibase.ApiReturn{"string", "进度附带信息"},
					"$.finish_time":   apibase.ApiReturn{"date", "操作结束时间"},
					"$.ts":            apibase.ApiReturn{"date", "进度时间戳"},
					"$.collector_id":  apibase.ApiReturn{"string", "目标采集器标识"},
					"$.agent_id":      apibase.ApiReturn{"string", "目标agent标识"},
				},
			}},
		},
	}
}

func normalOperaReturnDoc(progName string) apibase.Docs {
	return apibase.Docs{
		GET: &apibase.ApiDoc{
			Returns: []apibase.ApiReturnGroup{{
				Fields: apibase.ResultFields{
					"$.operation_seq": apibase.ApiReturn{"int", "操作序列号"},
					"$.agent_id":      apibase.ApiReturn{"string", "agent唯一标识"},
				},
			}},
		},
	}
}

var AgentListManagementRoutes = apibase.Route{
	Path: "agent",
	GET:  impl.QueryAgentList,
	Docs: apibase.Docs{
		GET: &apibase.ApiDoc{
			Name:  "获取Agent列表（混合方式）",
			Query: apibase.ExtendPaginationQueryParamMap(apibase.ApiParams{}),
			Returns: []apibase.ApiReturnGroup{{
				Fields: agentListInfo(),
			}},
		},
	},
	SubRoutes: []apibase.Route{{
		Path: "{agent_id:string}",
		GET:  impl.GetAgentInfo,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取与对应Sidecar关联的Agent",
				Returns: []apibase.ApiReturnGroup{{
					Desc:   "返回EasyDB Agent信息",
					Fields: extendAgentInfo(apibase.ResultFields{}),
				}, {
					Desc:   "返回easyLog Agent信息",
					Fields: extendAgentInfo(apibase.ResultFields{}),
				}},
			},
		},
		SubRoutes: []apibase.Route{{
			Path: "name",
			GET:  impl.GetAgentName,
			POST: impl.SetAgentName,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取Agent备注名",
				},
				POST: &apibase.ApiDoc{
					Name: "更改Agent备注名",
				},
			},
		}, {
			Path: "status",
			GET:  impl.GetAgentStatus,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取Agent当前状态",
					Returns: []apibase.ApiReturnGroup{{
						Desc: "返回Agent信息",
						Fields: apibase.ResultFields{
							"$.status": apibase.ApiReturn{"int", "Agent状态"},
						},
					}},
				},
			},
		}, {
			Path: "sid",
			GET:  impl.GetAgentSid,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取agent对应主机sid",
					Returns: []apibase.ApiReturnGroup{{
						Desc: "返回Agent对应主机sid",
						Fields: apibase.ResultFields{
							"$.sid": apibase.ApiReturn{"string", "Agent对应主机sid"},
						},
					}},
				},
			},
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/stop
			Path: "stop",
			GET:  impl.StopAgent,
			Docs: normalOperaReturnDoc("终止Agent运行"),
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/stop/progress
			Path: "stop/progress",
			GET:  impl.GetAgentStoppingProgress,
			Docs: progressReturnDoc("终止操作"),
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/start
			Path: "start",
			GET:  impl.StartAgent,
			Docs: normalOperaReturnDoc("启动Agent运行"),
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/startWithParam
			Path: "startWithParam",
			GET:  impl.StartAgentWithParam,
			Docs: normalOperaReturnDoc("启动Agent运行"),
		}, {
			Path: "start/progress",
			GET:  impl.GetAgentStartingProgress,
			Docs: progressReturnDoc("启动操作"),
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/restart
			Path: "restart",
			POST: impl.RestartAgent,
			Docs: normalOperaReturnDoc("重启Agent"),
		}, {
			Path: "restart/progress",
			GET:  impl.GetAgentRestartingProgress,
			Docs: progressReturnDoc("重启操作"),
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/config
			Path: "config",
			GET:  impl.GetAgentConfig,
			POST: impl.UpdateAgentConfig,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取Agent配置",
				},
				POST: &apibase.ApiDoc{
					Name: "更新Agent配置",
				},
			},
		}, {
			Path: "config/progress",
			GET:  impl.GetAgentConfigUpdatingProgress,
			Docs: progressReturnDoc("配置更新"),
		}, {
			Path: "reinstall",
			POST: impl.ReinstallAgent,
			Docs: normalOperaReturnDoc("重新安装Agent"),
		}, {
			Path: "reinstall/progress",
			GET:  impl.GetAgentReinstallProgress,
			Docs: progressReturnDoc("重新安装"),
		}, {
			Path: "install/progress",
			GET:  impl.GetAgentInstallProgress,
			Docs: progressReturnDoc("安装"),
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/update
			Path: "update",
			POST: impl.UpdateAgent,
			Docs: normalOperaReturnDoc("升级Agent"),
		}, {
			Path: "update/progress",
			GET:  impl.GetAgentUpdateInstallProgress,
			Docs: progressReturnDoc("升级安装"),
		}, {
			Path: "uninstall",
			POST: impl.UninstallAgent,
			Docs: normalOperaReturnDoc("卸载Agent"),
		}, {
			Path: "uninstall/progress",
			GET:  impl.GetAgentUninstallProgress,
			Docs: progressReturnDoc("卸载"),
		}, {
			Path: "progress",
			GET:  impl.GetAgentProgress,
			Docs: progressReturnDoc("获取进度"),
		}, {
			Path: "autoupdate",
			GET:  impl.GetAgentAutoupdateConfig,
			POST: impl.SetAgentAutoupdateConfig,
		}, {
			Path: "events",
			GET:  impl.GetAgentEventHistory,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name:  "获取Agent操作／事件记录",
					Query: apibase.ExtendPaginationQueryParamMap(apibase.ApiParams{}),
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$[*].name": apibase.ApiReturn{"string", "事件名称"},
							"$[*].ts":   apibase.ApiReturn{"date", "事件时间戳"},
						},
					}},
				},
			},
		}, {
			//http://xxxx/api/v1/agent/{agent_id}/uninstallSync
			Path: "uninstallSync",
			POST: impl.UninstallAgentSync,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "卸载agent",
					Body: apibase.ApiParams{
						"$.parameter":   apibase.ApiParam{"string", "卸载脚本参数", "", true},
						"$.shellScript": apibase.ApiParam{"string", "agent卸载脚本", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
							"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
							"$.result":        apibase.ApiReturn{"object", "卸载脚本执行结果"},
						},
					}},
				},
			},
		},
			{
				//http://xxxx/api/v1/agent/{agent_id}/startSync
				Path: "startSync",
				GET:  impl.StartAgentSync,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "start agent",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
								"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
								"$.result":        apibase.ApiReturn{"object", "执行结果"},
							},
						}},
					},
				},
			},
			{
				//http://xxxx/api/v1/agent/{agent_id}/startSyncWithParam
				Path: "startSyncWithParam",
				POST: impl.StartAgentSyncWithParam,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "start agent",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
								"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
								"$.result":        apibase.ApiReturn{"object", "执行结果"},
							},
						}},
					},
				},
			},
			{
				//http://xxxx/api/v1/agent/{agent_id}/stopSync
				Path: "stopSync",
				GET:  impl.StopAgentSync,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "stop agent",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
								"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
								"$.result":        apibase.ApiReturn{"object", "执行结果"},
							},
						}},
					},
				},
			},
			{
				//http://xxxx/api/v1/agent/{agent_id}/restartSync
				Path: "restartSync",
				GET:  impl.RestartAgentSync,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "restart agent",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
								"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
								"$.result":        apibase.ApiReturn{"object", "执行结果"},
							},
						}},
					},
				},
			},
			{
				//http://xxxx/api/v1/agent/{agent_id}/configSync
				Path: "configSync",
				POST: impl.UpdateAgentConfigSync,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "agent配置文件更新",
						Body: apibase.ApiParams{
							"$.config_content": apibase.ApiParam{"string", "", "配置文件内容", true},
						},
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
								"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
								"$.result":        apibase.ApiReturn{"object", "执行结果"},
							},
						}},
					},
				},
			}},
	}, {
		//http://xxxx/api/v1/agent/install
		Path: "install",
		POST: impl.InstallAgent,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "安装agent",
				Body: apibase.ApiParams{
					"$.collectorId":       apibase.ApiParam{"string", "目标采集器标识", "", true},
					"$.configurationPath": apibase.ApiParam{"string", "agent配置文件地址", "", true},
					"$.binaryPath":        apibase.ApiParam{"string", "agent可执行文件地址", "", true},
					"$.name":              apibase.ApiParam{"string", "agent名称", "", true},
					"$.parameter":         apibase.ApiParam{"string", "agent运行参数", "", true},
					"$.installScript":     apibase.ApiParam{"string", "agent安装脚本", "", true},
				},
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
						"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
					},
				}},
			},
		},
	}, {
		//http://xxxx/api/v1/agent/installSync
		Path: "installSync",
		POST: impl.InstallAgentSync,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "安装agent",
				Body: apibase.ApiParams{
					"$.collectorId":       apibase.ApiParam{"string", "目标采集器标识", "", true},
					"$.configurationPath": apibase.ApiParam{"string", "agent配置文件地址", "", true},
					"$.binaryPath":        apibase.ApiParam{"string", "agent可执行文件地址", "", true},
					"$.name":              apibase.ApiParam{"string", "agent名称", "", true},
					"$.parameter":         apibase.ApiParam{"string", "agent运行参数", "", true},
					"$.installScript":     apibase.ApiParam{"string", "agent安装脚本", "", true},
				},
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.agent_id":      apibase.ApiReturn{"string", "Agent ID"},
						"$.operation_seq": apibase.ApiReturn{"int", "操作序列号,可用于查询操作状态和进度"},
						"$.result":        apibase.ApiReturn{"object", "安装脚本执行结果"},
					},
				}},
			},
		},
	}, {
		//http://xxxx/api/v1/agent/cancelOperation
		Path: "cancelOperation",
		POST: impl.CancelOperation,
	}},
}
