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

var TaskOperationEasyMatrixAPIRoutes = apibase.Route{
	//http://xxxx/api/v2/task
	Path: "task",
	GET:  impl.TaskList,
	Docs: apibase.Docs{
		GET: &apibase.ApiDoc{
			Name: "获取所有任务信息接口[EasyMatrix API]",
		},
	},
	SubRoutes: []apibase.Route{{
		//http://xxxx/api/v2/task/upload
		Path: "upload",
		Middlewares: []context.Handler{
			apibase.CheckPermission3,
		},
		POST: impl.TaskUpload,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "上传脚本接口[EasyMatrix API]",
			},
		},
	}, {
		//http://xxxx/api/v2/task/status
		Path: "status",
		Middlewares: []context.Handler{
			apibase.CheckPermission3,
		},
		POST: impl.ModifyStatus,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "修改任务定时状态接口[EasyMatrix API]",
			},
		},
	}, {
		//http://xxxx/api/v2/task/{id:int(primary key)}
		Path: "{id:int}",
		Middlewares: []context.Handler{
			apibase.CheckPermission3,
		},
		DELETE: impl.TaskDelete,
		Docs: apibase.Docs{
			DELETE: &apibase.ApiDoc{
				Name: "删除任务接口[EasyMatrix API]",
			},
		},
	}, {
		//http://xxxx/api/v2/task/{id:int(primary key)}
		Path: "{id:int}",
		Middlewares: []context.Handler{
			apibase.CheckPermission3,
		},
		POST: impl.TaskUpdate,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "修改定时设置接口[EasyMatrix API]",
			},
		},
		SubRoutes: []apibase.Route{
			{
				//http://xxxx/api/v2/task/{id:int(primary key)}/content
				Path: "content",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				GET: impl.TaskFileContent,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "获取脚本内容接口[EasyMatrix API]",
					},
				},
			}, {
				//http://xxxx/api/v2/task/{id:int(primary key)}/edit
				Path: "edit",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				POST: impl.TaskEdit,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "编辑脚本接口[EasyMatrix API]",
					},
				},
			}, {
				//http://xxxx/api/v2/task/{id:int(primary key)}/cronParse
				Path: "cronParse",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				POST: impl.ParseSpec,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "校验cron表达式合法性接口[EasyMatrix API]",
					},
				},
			}, {
				//http://xxxx/api/v2/task/{id:int(primary key)}/run
				Path: "run",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				POST: impl.TaskRun,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "校验cron表达式合法性接口[EasyMatrix API]",
					},
				},
			}, {
				//http://xxxx/api/v2/task/{id:int(primary key)}/log
				Path: "log",
				Middlewares: []context.Handler{
					apibase.CheckPermission3,
				},
				GET: impl.TaskLogs,
				Docs: apibase.Docs{
					POST: &apibase.ApiDoc{
						Name: "获取执行历史接口[EasyMatrix API]",
					},
				},
			},
		}},
	},
}
