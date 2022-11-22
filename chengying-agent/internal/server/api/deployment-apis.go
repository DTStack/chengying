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

var DeploymentManageApis = apibase.Route{
	Path: "deploy",
	SubRoutes: []apibase.Route{{
		Path: "sidecar",
		GET:  impl.DeploySidecarMain,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "sidecar部署接口",
			},
		},
		SubRoutes: []apibase.Route{{
			Path: "install/shell",
			GET:  impl.GetSidecarInstallShell,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "sidecar部署脚本下载接口",
				},
			},
		}, {
			Path: "install/download",
			GET:  impl.GetSidecarInstallTargz,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "sidecar部署tar包下载接口",
				},
			},
		}, {
			Path: "install/callback",
			GET:  impl.GetSidecarInstallCallback,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "sidecar部署完成回调接口",
					Query: apibase.ApiParams{
						"client_id":    apibase.ApiParam{"string", "客户端的uuid标识", "", true},
						"install_type": apibase.ApiParam{"string", "客户端的类型", "", true},
						"msg":          apibase.ApiParam{"string", "安装结论信息", "", true},
						"install_res":  apibase.ApiParam{"string", "安装结果标识[success,failed]", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$[*].sidecar_id": apibase.ApiReturn{"int", "部署Sidecar_ID"},
						},
					}},
				},
			},
		}},
	}, {
		Path: "product",
		GET:  impl.QueryDeploymentProducts,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取用于部署的产品列表",
				Query: apibase.ExtendPaginationQueryParamMap(apibase.ApiParams{
					"type": apibase.ApiParam{"string", "产品类型", "", false},
				}),
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$[*].id":      apibase.ApiReturn{"int", "部署ID"},
						"$[*].type":    apibase.ApiReturn{"string", "产品类型"},
						"$[*].name":    apibase.ApiReturn{"string", "产品名称"},
						"$[*].version": apibase.ApiReturn{"string", "版本"},
						"$.total":      apibase.ApiReturn{"int", "全部记录数"},
					},
				}},
			},
		},
		SubRoutes: []apibase.Route{{
			Path:   "{prod_id:int}",
			GET:    impl.GetDeploymentProductInfo,
			DELETE: impl.RemoveDeploymentProduct,
			POST:   impl.UpdateDeploymentProductInfo,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取用于部署的产品详细信息",
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$.id":                apibase.ApiReturn{"int", "部署ID"},
							"$.type":              apibase.ApiReturn{"string", "产品类型"},
							"$.name":              apibase.ApiReturn{"string", "产品名称"},
							"$.version":           apibase.ApiReturn{"string", "版本"},
							"$.description":       apibase.ApiReturn{"string", "描述／备注"},
							"$.checksum":          apibase.ApiReturn{"string", "MD5校验码"},
							"$.release_date":      apibase.ApiReturn{"date", "发布时间"},
							"$.deploy_times":      apibase.ApiReturn{"int", "成功部署的机器数"},
							"$.deploy_fail_times": apibase.ApiReturn{"int", "部署失败的机器数"},
						},
					}},
				},
				POST: &apibase.ApiDoc{
					Name: "更新产品的名称、版本、描述、URL路径等信息",
					Body: apibase.ApiParams{
						"$.name":        apibase.ApiParam{"string", "更新名称", "", false},
						"$.type":        apibase.ApiParam{"string", "更新类型，如果已经被部署过这不能更新类型", "", false},
						"$.version":     apibase.ApiParam{"string", "更新版本号,如果已经被部署过则不能更新版本", "", false},
						"$.url":         apibase.ApiParam{"string", "更新URL", "", false},
						"$.description": apibase.ApiParam{"string", "更新描述信息", "", false},
					},
				},
				DELETE: &apibase.ApiDoc{
					Name: "删除产品信息",
				},
			},
		}, {
			Path: "new",
			POST: impl.NewDeploymentProductInfo,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "添加新的产品",
					Body: apibase.ApiParams{
						"$.name":        apibase.ApiParam{"string", "名称", "", true},
						"$.type":        apibase.ApiParam{"string", "类型", "", true},
						"$.version":     apibase.ApiParam{"string", "版本号", "", true},
						"$.url":         apibase.ApiParam{"string", "URL", "", true},
						"$.description": apibase.ApiParam{"string", "描述信息", "", false},
					},
				},
			},
		}},
	}},
}
