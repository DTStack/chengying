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
)

var SshOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "agent",
	SubRoutes: []apibase.Route{{
		Path: "install",
		GET:  impl.InstallInit,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "ssh 相关提供给操作台的接口[EasyMatrix API]",
			},
		},
		SubRoutes: []apibase.Route{{
			//http://xxxx/api/v2/agent/install/pwdconnect
			Path: "pwdconnect",
			POST: impl.CheckPwdConnect,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "ssh 连通性检查通过password密码[EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.host":     apibase.ApiParam{"string", "主机域名或者IP", "", true},
						"$.port":     apibase.ApiParam{"int", "端口", "", true},
						"$.user":     apibase.ApiParam{"string", "用户名", "", true},
						"$.password": apibase.ApiParam{"string", "登录密码", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/pkconnect
			Path: "pkconnect",
			POST: impl.CheckPkConnect,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "ssh 连通性检查通过pk秘钥文件[EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.host": apibase.ApiParam{"string", "主机域名或者IP", "", true},
						"$.port": apibase.ApiParam{"int", "端口", "", true},
						"$.user": apibase.ApiParam{"string", "用户名", "", true},
						"$.pk":   apibase.ApiParam{"string", "秘钥文件内容", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/pwdinstall
			Path: "pwdinstall",
			POST: impl.AgentInstallByPwd,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "agent 安装 通过ssh 密码建立连接 [EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.host":       apibase.ApiParam{"string", "主机域名或者IP", "", true},
						"$.port":       apibase.ApiParam{"int", "端口", "", true},
						"$.user":       apibase.ApiParam{"string", "用户名", "", true},
						"$.password":   apibase.ApiParam{"string", "密码", "", true},
						"$.group":      apibase.ApiParam{"string", "分组", "default", false},
						"$.type":       apibase.ApiParam{"string", "安装类型", "pwd", false},
						"$.cluster_id": apibase.ApiParam{"int", "集群id", "", true},
						"cluster_type": apibase.ApiParam{"string", "集群类型", "hosts", true},
						"$.role":       apibase.ApiParam{"string", "角色", "", false},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/installCmd
			Path: "installCmd",
			GET:  impl.AgentInstallCmd,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "agent 安装 sh [EasyMatrix API]",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/pkinstall
			Path: "pkinstall",
			POST: impl.AgentInstallByPk,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "agent 安装 通过ssh 秘钥文件建立连接 [EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.host":         apibase.ApiParam{"string", "主机域名或者IP", "", true},
						"$.port":         apibase.ApiParam{"int", "端口", "", true},
						"$.user":         apibase.ApiParam{"string", "用户名", "", true},
						"$.pk":           apibase.ApiParam{"string", "秘钥文件内容", "", true},
						"$.group":        apibase.ApiParam{"string", "分组", "default", false},
						"$.type":         apibase.ApiParam{"string", "安装类型", "pwd", false},
						"$.cluster_type": apibase.ApiParam{"string", "集群类型", "hosts", true},
						"$.cluster_id":   apibase.ApiParam{"int", "集群id", "", true},
						"$.role":         apibase.ApiParam{"string", "角色", "", false},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/checkinstall
			Path: "checkinstall",
			POST: impl.AgentInstallCheck,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "agent 安装 结果检查 [EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.aid": apibase.ApiParam{"int", "Agent ID", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/checkinstallall
			Path: "checkinstallall",
			GET:  impl.AgentInstallCheckAll,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "agent 安装 结果检查 [EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.start": apibase.ApiParam{"int", "start 起始翻页", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/checkinstallbyip
			Path: "checkinstallbyip",
			POST: impl.AgentInstallCheckByIp,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "agent 安装 结果检查 根据Ip地址查询 [EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.ip": apibase.ApiParam{"string", "IP", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/checkinstallbysid
			Path: "checkinstallbysid",
			GET:  impl.AgentInstallCheckBySid,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "agent 安装 结果检查 根据Sid查询 [EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.sid": apibase.ApiParam{"string", "sid", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.code": apibase.ApiReturn{"int", "True Or False"},
							"$.msg":  apibase.ApiReturn{"string", "message"},
							"$.data": apibase.ApiReturn{"string", "data"},
						},
					}},
				}},
		}, {
			//http://xxxx/api/v2/agent/install/callback?aid=-1
			Path: "callback",
			GET:  impl.AgentInstallCallBack,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "agent 安装 结果callback [EasyMatrix API]",
					Body: apibase.ApiParams{
						"$.aid": apibase.ApiParam{"int", "Agent ID", "", true},
					},
				},
			},
		}},
	}, {
		//http://xxxx/api/v2/agent/hosts
		Path: "hosts",
		GET:  impl.AgentHosts,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取所有可用的主机 [EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "True Or False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		//http://xxxx/api/v2/agent/hostgroups
		Path: "hostgroups",
		GET:  impl.AgentHostGroups,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取所有可用的主机组 [EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "True Or False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		//http://xxxx/api/v2/agent/hostgroup_rename
		Path: "hostgroup_rename",
		POST: impl.AgentHostGroupRename,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "主机组重命名 [EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "True Or False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		//http://xxxx/api/v2/agent/hostmove
		Path: "hostmove",
		POST: impl.AgentHostMove,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "主机移动至其它组 [EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "True Or False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		//http://xxxx/api/v2/agent/hostService
		Path: "hostService",
		GET:  impl.AgentHostService,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取主机下的服务列表 [EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "True Or False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		//http://xxxx/api/v2/agent/hostdelete
		Path: "hostdelete",
		POST: impl.AgentHostDelete,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "删除全部主机下主机选中主机 [EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "True Or False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			},
		},
	}},
}
