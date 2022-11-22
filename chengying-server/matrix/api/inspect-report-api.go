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

var InspectReportEasyMatrixRoutes = apibase.Route{
	Path: "inspect",
	SubRoutes: []apibase.Route{
		{
			Path: "service/status",
			GET:  impl.GetServiceStatus,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取当前集群各产品下服务健康情况",
				},
			},
		}, {
			Path: "alert/history",
			GET:  impl.GetAlertHistory,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取指定时间段内告警历史",
				},
			},
		}, {
			Path: "host/status",
			GET:  impl.GetHostStatus,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取当前集群节点CPU、内存、磁盘健康信息",
				},
			},
		}, {
			Path: "graph",
			SubRoutes: []apibase.Route{
				{
					Path: "config",
					GET:  impl.GetGraphConfig,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取图表配置列表",
						},
					},
				},
				{
					Path: "data",
					GET:  impl.GetGraphData,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "获取图表数据",
						},
					},
				},
			},
		}, {
			Path: "generate",
			POST: impl.StartGenerateReport,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "生成巡检报告",
				},
			},
		}, {
			Path: "progress",
			GET:  impl.GetReportProgress,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "查看生成pdf报告进度",
				},
			},
		}, {
			Path: "download",
			GET:  impl.Download,
			Middlewares: []context.Handler{
				apibase.CheckPermission1,
			},
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "下载巡检报告",
				},
			},
		},
	},
}
