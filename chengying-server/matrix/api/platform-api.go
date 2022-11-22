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

var PlatformOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "platform",
	SubRoutes: []apibase.Route{
		{
			Path: "inspect",
			SubRoutes: []apibase.Route{
				{
					Path: "baseInfo",
					SubRoutes: []apibase.Route{
						{
							Path: "name_node",
							GET:  impl.GetPlatformInspectNameNodeBaseInfo,
							Docs: apibase.Docs{
								GET: &apibase.ApiDoc{
									Name: "get inspect name_node baseInfo[EasyMatrix API]",
								},
							},
						}, {
							Path: "status",
							GET:  impl.GetPlatformInspectBaseInfoState,
							Docs: apibase.Docs{
								GET: &apibase.ApiDoc{
									Name: "get inspect baseInfo[EasyMatrix API]",
								},
							},
						},
					},
				}, {
					Path: "graph/config",
					GET:  impl.GetPlatformGraphConfig,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "get graph config [EasyMatrix API]",
						},
					},
				}, {
					Path: "statisticsConfig",
					GET:  impl.GetInspectConfig,
					SubRoutes: []apibase.Route{
						{
							Path: "update",
							POST: impl.ModifyInspectConfig,
							Docs: apibase.Docs{
								GET: &apibase.ApiDoc{
									Name: "update inspect config [EasyMatrix API]",
								},
							},
						},
					},
				}, {
					Path: "form/data",
					GET:  impl.GetPlatformInspectFormData,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "get form data by form title[EasyMatrix API]",
						},
					},
				},
			},
		},
	},
}
