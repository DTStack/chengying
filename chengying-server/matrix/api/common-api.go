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

var CommonOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "common",
	SubRoutes: []apibase.Route{
		{
			Path: "file2text",
			POST: impl.File2text,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "upload file convert to text[EasyMatrix API]",
				},
			},
		}, {
			Path: "safetyAudit",
			SubRoutes: []apibase.Route{
				{
					Path: "list",
					GET:  impl.GetSafetyAuditList,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "get safety audit list[EasyMatrix API]",
						},
					},
				}, {
					Path: "module",
					GET:  impl.GetSafetyAuditModule,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "get safety audit module[EasyMatrix API]",
						},
					},
				}, {
					Path: "operation",
					GET:  impl.GetSafetyAuditOperation,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "get safety audit operation[EasyMatrix API]",
						},
					},
				},
			},
		}, {
			Path: "deployInfo",
			SubRoutes: []apibase.Route{
				{
					//api/v2/common/deployInfo/generate
					Path: "generate",
					POST: impl.DeployInfoGenerate,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "部署信息markdown文件生成",
						},
					},
				}, {
					//api/v2/common/deployInfo/download
					Path: "download",
					GET:  impl.DeployInfoDownload,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "部署信息markdown文件下载",
						},
					},
				},
			},
		},
	},
}
