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
	"easyagent/internal/server/api/impl"
)

var SidecarManagementRoutes = apibase.Route{
	Path: "sidecar",
	GET:  impl.QuerySidecarList,
	Docs: apibase.Docs{
		GET: &apibase.ApiDoc{
			Name:  "获取Sidecart列表",
			Query: apibase.ApiParams{},
			Returns: []apibase.ApiReturnGroup{{
				Fields: apibase.ResultFields{
					"$[*].id":          apibase.ApiReturn{"string", "Sidecar ID"},
					"$[*].name":        apibase.ApiReturn{"string", "名称／备注"},
					"$[*].os":          apibase.ApiReturn{"string", "操作系统"},
					"$[*].version":     apibase.ApiReturn{"string", "版本"},
					"$[*].auto-deploy": apibase.ApiReturn{"bool", "是否自动部署"},
					"$[*].auto-update": apibase.ApiReturn{"bool", "是否自动升级"},
				},
			}},
		},
	},
	SubRoutes: []apibase.Route{{
		Path: "new",
		POST: impl.NewSidecar,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Body: apibase.ApiParams{
					"$.name":    apibase.ApiParam{"string", "Sidecar部署名称", "", false},
					"$.version": apibase.ApiParam{"string", "部署的版本号", "", true},
				},
			},
		},
	}, {
		Path: "{sidecar_id:string}",
		GET:  impl.GetSidecarInformation,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取客户端信息",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.id":           apibase.ApiReturn{"string", "客户端ID"},
						"$.name":         apibase.ApiReturn{"string", "客户端备注名"},
						"$.capabilities": apibase.ApiReturn{"array", "客户端支持的功能"},
						"$.state":        apibase.ApiReturn{"string", "客户端状态"},
					},
				}},
			},
		},
		SubRoutes: []apibase.Route{{
			//http://xxxx/api/v1/sidecar/{sidecar_id}/execscript
			Path: "execscript",
			POST: impl.ExecScript,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.operation_seq": apibase.ApiReturn{"int", "操作序列号"},
							"$.agent_id":      apibase.ApiReturn{"string", "agent唯一标识"},
						},
					}},
				},
			},
		}, {
			//http://xxxx/api/v1/sidecar/{sidecar_id}/execscriptSync
			Path: "execscriptSync",
			POST: impl.ExecScriptSync,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.operation_seq": apibase.ApiReturn{"int", "操作序列号"},
							"$.agent_id":      apibase.ApiReturn{"string", "agent唯一标识"},
							"$.result":        apibase.ApiReturn{"string", "exec result"},
						},
					}},
				},
			},
		}, {
			//http://xxxx/api/v1/sidecar/{sidecar_id}/execscriptSyncBase64
			Path: "execscriptSyncBase64",
			POST: impl.ExecScriptSyncBase64,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.operation_seq": apibase.ApiReturn{"int", "操作序列号"},
							"$.agent_id":      apibase.ApiReturn{"string", "agent唯一标识"},
							"$.result":        apibase.ApiReturn{"string", "exec result"},
						},
					}},
				},
			},
		}, {
			//http://xxxx/api/v1/sidecar/{sidecar_id}/execrestSync
			Path: "execrestSync",
			POST: impl.ExecRestSync,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.operation_seq": apibase.ApiReturn{"int", "操作序列号"},
							"$.agent_id":      apibase.ApiReturn{"string", "agent唯一标识"},
							"$.result":        apibase.ApiReturn{"string", "exec result"},
						},
					}},
				},
			},
		}, {
			//http://xxxx/api/v1/sidecar/{sidecar_id}/exescript/progress
			Path: "exescript/progress",
			POST: impl.GetExecScriptProgress,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "",
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
			},
		}, {
			//http://xxxx/api/v1/sidecar/{sidecar_id}/agents
			Path: "agents",
			GET:  impl.GetSidecarAgents,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取采集器(collector,sidecar)下具体agent列表",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.agent_id":    apibase.ApiReturn{"string", "Agent ID"},
							"$.deploy_date": apibase.ApiReturn{"date", "部署日期"},
						},
					}},
				},
			},
		}, {
			Path: "install",
			POST: impl.InstallSidecar,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "安装客户端",
					Body: apibase.ApiParams{
						"$.host":       apibase.ApiParam{"string", "目标主机地址", "", true},
						"$.ssh.port":   apibase.ApiParam{"int", "SSH端口号", "22", false},
						"$.ssh.user":   apibase.ApiParam{"string", "SSH登陆帐号", "root", true},
						"$.ssh.passwd": apibase.ApiParam{"string", "SSH登陆密码", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.id": apibase.ApiReturn{"int", "Agent ID"},
						},
					}},
				},
			},
		}, {
			Path: "status",
			GET:  impl.GetSidecarStatus,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取客户端状态",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.status": apibase.ApiReturn{"string", "客户端状态"},
						},
					}},
				},
			},
		}, {
			Path: "rename",
			POST: impl.ChangeSidecarName,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "修改客户端备注名",
					Body: apibase.ApiParams{
						"$.name": apibase.ApiParam{"string", "新的备注名", "", true},
					},
				},
			},
		}, {
			Path: "control",
			POST: impl.ControlSidecar,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "控制客户端",
					Body: apibase.ApiParams{
						"$.cmd": apibase.ApiParam{"string", "控制命令", "", true},
					},
				},
			},
		}},
	}},
}
