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

var ServiceOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "service",
	SubRoutes: []apibase.Route{
		{
			Path: "{pid:int}",
			SubRoutes: []apibase.Route{
				{
					Path: "{service_name:string}",
					SubRoutes: []apibase.Route{{
						Path: "start",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						POST: impl.ServiceStart,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "service start[EasyMatrix API]",
							},
						},
					}, {
						Path: "stop",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						POST: impl.ServiceStop,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "service stop[EasyMatrix API]",
							},
						},
					}, {
						Path: "rolling_restart",
						Middlewares: []context.Handler{
							apibase.CheckPermission3,
						},
						POST: impl.ServiceRollingRestart,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "service rolling restart[EasyMatrix API]",
							},
						},
					}, {
						Path: "config_update",
						Middlewares: []context.Handler{
							apibase.CheckPermission1,
						},
						POST: impl.ServiceRollingConfigUpdate,
						Docs: apibase.Docs{
							GET: &apibase.ApiDoc{
								Name: "service config update[EasyMatrix API]",
							},
						},
					}},
				},
			},
		}, {
			Path: "license",
			POST: impl.License,
		},
	},
}
