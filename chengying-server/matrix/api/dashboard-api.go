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
)

var DashboardOperationEasyMatrixRoutes = apibase.Route{
	Path: "dashboard",
	SubRoutes: []apibase.Route{
		{
			Path: "import",
			POST: impl.ImportDashboard,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "导入grafana仪表盘接口",
				},
			},
		},
		{
			Path: "export",
			GET: impl.ExportDashboard,
			Docs: apibase.Docs{
				POST: &apibase.ApiDoc{
					Name: "导出grafana仪表盘接口",
				},
			},
		}, {
			Path: "alerts",
			GET:  impl.GetDashboardAlerts,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取告警规则列表",
				},
			},
			SubRoutes: []apibase.Route{
				{
					Path: "pause",
					POST: impl.DashboardAlertsPause,
					Docs: apibase.Docs{
						POST: &apibase.ApiDoc{
							Name: "告警规则停止开启",
						},
					},
				},
			},
		},
	},
}
