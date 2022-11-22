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

var InstanceRecordOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "instance_record",
	SubRoutes: []apibase.Route{
		{
			Path: "{id:int}",
			SubRoutes: []apibase.Route{
				{
					Middlewares: []context.Handler{
						apibase.CheckPermission1,
					},
					Path: "force_stop",
					POST: impl.ForceStop,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "force stop instance specified by record id[EasyMatrix API]",
						},
					},
				}, {
					Middlewares: []context.Handler{
						apibase.CheckPermission1,
					},
					Path: "force_uninstall",
					POST: impl.ForceUninstall,
					Docs: apibase.Docs{
						GET: &apibase.ApiDoc{
							Name: "force uninstall instance specified by record id[EasyMatrix API]",
						},
					},
				},
			},
		},
	},
}
