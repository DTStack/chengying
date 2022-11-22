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

var ServerManageApis = apibase.Route{
	Path: "server",
	SubRoutes: []apibase.Route{{
		Path: "dashboard",
		SubRoutes: []apibase.Route{{
			//http://xxxx/api/v1/server/dashboard/url
			Path: "url",
			GET:  impl.RetDashboardUrl,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "提供给出 dashboard的 URL列表API",
					Query: apibase.ApiParams{
						"type": apibase.ApiParam{"string", "'cluster' or 'services'", "", true},
						"id":   apibase.ApiParam{"string", "服务器组ID", "", true},
					},
					Returns: []apibase.ApiReturnGroup{{
						Fields: apibase.ResultFields{
							"$[*].url": apibase.ApiReturn{"string", "dashboard仪表盘链接URL"},
						},
					}},
				},
			},
		}},
	}},
}
