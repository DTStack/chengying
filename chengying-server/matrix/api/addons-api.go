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

var AddonsEasyMatrixAPIRoutes = apibase.Route{
	Path: "addons",
	SubRoutes: []apibase.Route{{
		//http://xxxx/api/v2/addons/upload
		Path: "upload",
		GET:  impl.AddonUpload,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "插件上传[EasyMatrix API]",
			},
		},
	}, {
		//http://xxxx/api/v2/addons/list
		Path: "list",
		GET:  impl.AddonList,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取支持的addons列表[EasyMatrix API]",
			},
		},
	}, {
		//http://xxxx/api/v2/addons/add
		Path: "add",
		GET:  impl.AddonAdd,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "增加插件[EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "0 True Or !=0 False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		//http://xxxx/api/v2/addons/delete
		Path: "delete",
		GET:  impl.AddonDelete,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "删除插件[EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "0 True Or !=0 False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		//http://xxxx/api/v2/addons/install
		Path: "install",
		GET:  impl.AddonInstall,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "安装插件[EasyMatrix API]",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.code": apibase.ApiReturn{"int", "0 True Or !=0 False"},
						"$.msg":  apibase.ApiReturn{"string", "message"},
						"$.data": apibase.ApiReturn{"string", "data"},
					},
				}},
			}},
	}, {
		Path: "{agent_id:string}",
		SubRoutes: []apibase.Route{
			{
				//http://xxxx/api/v2/addons/{agent_id}/stop
				Path: "stop",
				GET:  impl.AddonStop,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "停止插件[EasyMatrix API]",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.code": apibase.ApiReturn{"int", "0 True Or !=0 False"},
								"$.msg":  apibase.ApiReturn{"string", "message"},
								"$.data": apibase.ApiReturn{"string", "data"},
							},
						}},
					}},
			}, {
				//http://xxxx/api/v2/addons/{agent_id}/start
				Path: "start",
				POST: impl.AddonStart,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "启动插件[EasyMatrix API]",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.code": apibase.ApiReturn{"int", "0 True Or !=0 False"},
								"$.msg":  apibase.ApiReturn{"string", "message"},
								"$.data": apibase.ApiReturn{"string", "data"},
							},
						}},
					}},
			}, {
				//http://xxxx/api/v2/addons/{agent_id}/config
				Path: "config",
				POST: impl.AddonConfig,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "插件配置更新[EasyMatrix API]",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.code": apibase.ApiReturn{"int", "0 True Or !=0 False"},
								"$.msg":  apibase.ApiReturn{"string", "message"},
								"$.data": apibase.ApiReturn{"string", "data"},
							},
						}},
					}},
			}, {
				//http://xxxx/api/v2/addons/{agent_id}/uninstall
				Path: "uninstall",
				GET:  impl.AddonUnInstall,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "卸载插件[EasyMatrix API]",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.code": apibase.ApiReturn{"int", "0 True Or !=0 False"},
								"$.msg":  apibase.ApiReturn{"string", "message"},
								"$.data": apibase.ApiReturn{"string", "data"},
							},
						}},
					}},
			}, {
				//http://xxxx/api/v2/addons/{agent_id}/upgrade
				Path: "upgrade",
				POST: impl.AddonUpgrade,
				Docs: apibase.Docs{
					GET: &apibase.ApiDoc{
						Name: "插件升级 [EasyMatrix API]",
						Returns: []apibase.ApiReturnGroup{{
							Fields: apibase.ResultFields{
								"$.code": apibase.ApiReturn{"int", "True Or False"},
								"$.msg":  apibase.ApiReturn{"string", "message"},
								"$.data": apibase.ApiReturn{"string", "data"},
							},
						}},
					}},
			},
		},
	},
	},
}
