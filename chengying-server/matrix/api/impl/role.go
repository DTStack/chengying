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

package impl

import (
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/asset"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"github.com/kataras/iris/context"
	"gopkg.in/yaml.v2"
	"strconv"
)

type Permission struct {
	Title      string       `yaml:"title" json:"title"`
	Code       string       `yaml:"code" json:"code"`
	Permission int          `yaml:"permission" json:"permission"`
	Children   []Permission `yaml:"children" json:"children"`
	Selected   bool         `json:"selected"`
}

type PermissionList struct {
	Permissions []Permission `yaml:"permissions"`
}

func GetRolePermissions(ctx context.Context) apibase.Result {
	log.Debugf("GetRolePermissions: %v", ctx.Request().RequestURI)
	id, err := ctx.Params().GetInt("role_id")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	err, info := model.RoleList.GetInfoByRoleId(id)
	if err != nil {
		log.Errorf("Database err: %v", err)
		return err
	}
	roleValue, _ := strconv.Atoi(info.RoleValue)

	data, err := asset.Asset("templates/role-permissions.yml")
	if err != nil {
		log.Errorf("get role-permissions.yml error %v", err)
		return err
	}

	permissionTree := &PermissionList{}
	err = yaml.Unmarshal(data, permissionTree)
	if err != nil {
		log.Errorf("yaml unmarshal error %v", err)
		return err
	}

	for i := 0; i < len(permissionTree.Permissions); i++ {
		roleFilter(&permissionTree.Permissions[i], roleValue)
	}

	return map[string]interface{}{
		"role_name":   info.RoleName,
		"description": info.RoleDesc,
		"permissions": permissionTree,
	}
}

func roleFilter(p *Permission, roleValue int) {
	p.Selected = roleValue&p.Permission > 0
	for i := 0; i < len(p.Children); i++ {
		roleFilter(&p.Children[i], roleValue)
	}
}

func GetRolePermissionCodes(ctx context.Context) apibase.Result {
	log.Debugf("GetRolePermissionCodes: %v", ctx.Request().RequestURI)

	userPermission, err := apibase.GetTokenUserPermission(ctx)
	if err != nil {
		log.Errorf(err.Error())
	}

	data, err := asset.Asset("templates/role-permissions.yml")
	if err != nil {
		log.Errorf("get role-permissions.yml error %v", err)
		return err
	}

	permissionTree := &PermissionList{}
	err = yaml.Unmarshal(data, permissionTree)
	if err != nil {
		log.Errorf("yaml unmarshal error %v", err)
		return err
	}

	list := make([]string, 0)
	for i := 0; i < len(permissionTree.Permissions); i++ {
		roleCodeFilter(permissionTree.Permissions[i], userPermission, &list)
	}

	return list
}

func roleCodeFilter(p Permission, roleValue int, list *[]string) {
	if roleValue&p.Permission > 0 {
		*list = append(*list, p.Code)
	}
	for i := 0; i < len(p.Children); i++ {
		roleCodeFilter(p.Children[i], roleValue, list)
	}
}

func GetRoleList(ctx context.Context) apibase.Result {
	log.Debugf("GetRoleList: %v", ctx.Request().RequestURI)
	list, err := model.RoleList.GetList()
	if err != nil {
		return fmt.Errorf("get role list database err %v", err)
	}
	res := make([]map[string]interface{}, 0)
	for _, role := range list {
		res = append(res, map[string]interface{}{
			"id":          role.ID,
			"name":        role.RoleName,
			"desc":        role.RoleDesc,
			"update_time": role.UpdateTime.Time,
		})
	}
	return res
}
