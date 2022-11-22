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

var RoleOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "role",
	SubRoutes: []apibase.Route{{
		Path: "{role_id:int}",
		SubRoutes: []apibase.Route{{
			//http:://xxxx/api/v2/role/{role_id}/permissions
			Path: "permissions",
			GET:  impl.GetRolePermissions,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "通过role_id获取角色的所有权限点",
				},
			},
		}},
	}, {
		Path: "list",
		GET:  impl.GetRoleList,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取角色列表",
			},
		},
	}, {
		//http:://xxxx/api/v2/role/codes
		Path: "codes",
		GET:  impl.GetRolePermissionCodes,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "通过role_id获取角色的所有权限code，方便前端比较",
			},
		},
	}},
}
