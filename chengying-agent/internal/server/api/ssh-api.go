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

var SshhApiRoutes = apibase.Route{
	Path: "ssh",
	SubRoutes: []apibase.Route{{
		Path: "checkByUserPwd",
		POST: impl.CheckByUserPwd,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "ssh连通性(用户名密码)",
				Body: apibase.ApiParams{
					"$.host":     apibase.ApiParam{"string", "主机域名orIP", "", true},
					"$.port":     apibase.ApiParam{"int", "端口", "", true},
					"$.user":     apibase.ApiParam{"string", "用户名", "", true},
					"$.password": apibase.ApiParam{"string", "登录密码", "", true},
				},
			},
		},
	}, {
		Path: "runWithUserPwd",
		POST: impl.RunWithUserPwd,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "ssh安装(用户名密码)",
				Body: apibase.ApiParams{
					"$.host":     apibase.ApiParam{"string", "主机域名orIP", "", true},
					"$.port":     apibase.ApiParam{"int", "端口", "", true},
					"$.user":     apibase.ApiParam{"string", "用户名", "", true},
					"$.password": apibase.ApiParam{"string", "登录密码", "", true},
					"$.cmd":      apibase.ApiParam{"string", "一键安装脚本", "", true},
				},
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.result": apibase.ApiReturn{"string", "执行结果"},
					},
				}},
			},
		},
	}, {
		Path: "checkByUserPk",
		POST: impl.CheckByUserPk,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "ssh安装(用户名密钥)",
				Body: apibase.ApiParams{
					"$.host": apibase.ApiParam{"string", "主机域名orIP", "", true},
					"$.port": apibase.ApiParam{"int", "端口", "", true},
					"$.user": apibase.ApiParam{"string", "用户名", "", true},
					"$.pk":   apibase.ApiParam{"string", "秘钥内容", "", true},
				},
			},
		},
	}, {
		Path: "runWithUserPk",
		POST: impl.RunWithUserPk,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "ssh安装(用户名密钥)",
				Body: apibase.ApiParams{
					"$.host": apibase.ApiParam{"string", "主机域名orIP", "", true},
					"$.port": apibase.ApiParam{"int", "端口", "", true},
					"$.user": apibase.ApiParam{"string", "用户名", "", true},
					"$.pk":   apibase.ApiParam{"string", "秘钥内容", "", true},
					"$.cmd":  apibase.ApiParam{"string", "一键安装脚本", "", true},
				},
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.result": apibase.ApiReturn{"string", "执行结果"},
					},
				}},
			},
		},
	}},
}
