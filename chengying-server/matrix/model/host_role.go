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

package model

import (
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

/*
 @Author: zhijian
 @Date: 2021/4/12 14:18
 @Description:
*/

const (
	DefaultRoleType = 1
	CustomRoleType  = 2
)

type hostRole struct {
	dbhelper.DbTable
}
type HostRoleInfo struct {
	Id        int    `db:"id" json:"role_id"`
	ClusterId int    `db:"cluster_id" json:"-"`
	RoleName  string `db:"role_name" json:"role_name"`
	RoleType  int    `db:"role_type" json:"role_type"`
}

var HostRole = &hostRole{
	dbhelper.DbTable{USE_MYSQL_DB, HOST_ROLE},
}

func (h *hostRole) GetDefaultRoleNameMap() map[string]struct{} {
	sql := fmt.Sprintf("select role_name from %s where role_type = ?", HOST_ROLE)
	var roleNameList []string
	roleNameMap := make(map[string]struct{})
	h.GetDB().Select(&roleNameList, sql, DefaultRoleType)

	for _, name := range roleNameList {
		roleNameMap[name] = struct{}{}
	}
	return roleNameMap
}

func (h *hostRole) GetRoleListByRoleIdList(roleIdList []int) ([]HostRoleInfo, error) {
	var hostRoleInfos []HostRoleInfo

	for _, id := range roleIdList {
		roleInfo, err := h.GetRoleInfoById(id)
		if err != nil {
			return nil, err
		}
		hostRoleInfos = append(hostRoleInfos, roleInfo)
	}

	sort.Slice(hostRoleInfos, func(i, j int) bool {
		return hostRoleInfos[i].RoleName < hostRoleInfos[j].RoleName
	})

	return hostRoleInfos, nil
}

//根据 roleid(1,2,3) 获取 rolename 切片
func (h *hostRole) GetRoleNameListStrByIdList(idStrList string) ([]string, error) {
	if strings.TrimSpace(idStrList) == "" {
		return nil, nil
	}
	idMap := make(map[int]struct{})
	for _, idStr := range strings.Split(idStrList, ",") {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		idMap[id] = struct{}{}
	}
	var roleName []string
	for id, _ := range idMap {
		roleInfo, err := h.GetRoleInfoById(id)
		if err != nil {
			return nil, err
		}
		roleName = append(roleName, roleInfo.RoleName)
	}
	//必须排序，否则 schema 可能不同
	sort.Strings(roleName)
	return roleName, nil
}

func (h *hostRole) GetRoleInfoById(roleId int) (HostRoleInfo, error) {
	sql := fmt.Sprintf("select * from %s where id = ?", HOST_ROLE)
	hostRoleInfo := HostRoleInfo{}
	err := h.GetDB().Get(&hostRoleInfo, sql, roleId)
	if err != nil {
		return hostRoleInfo, err
	}
	return hostRoleInfo, nil
}

func (h *hostRole) AddRole(clusterId int, roleName string) error {

	sql := fmt.Sprintf("INSERT INTO %s (cluster_id,role_name, role_type) VALUES (? ,?, %d);", HOST_ROLE, CustomRoleType)
	_, err := h.GetDB().Exec(sql, clusterId, roleName)
	if err != nil {
		return err
	}
	return nil
}

func (h *hostRole) DeleteRole(roleId int) error {

	sql := fmt.Sprintf("delete from %s where id = ?", HOST_ROLE)
	_, err := h.GetDB().Exec(sql, roleId)
	if err != nil {
		return err
	}
	return nil
}

func (h *hostRole) InitNewCluster(clusterId int) error {
	defaultRole := []string{"web", "manager", "worker"}
	for _, roleName := range defaultRole {
		sql := fmt.Sprintf("INSERT INTO %s ( cluster_id, role_name, role_type) VALUES ( ?, '%s', %d)", HOST_ROLE, roleName, DefaultRoleType)
		_, err := h.GetDB().Exec(sql, clusterId)
		if err != nil {
			return err
		}
	}
	return nil
}
func (h *hostRole) GetRoleListByClusterId(clusterId int) ([]HostRoleInfo, error) {
	sql := fmt.Sprintf("select * from %s where cluster_id = ? order by id asc", HOST_ROLE)
	var hostRoleList []HostRoleInfo
	err := h.GetDB().Select(&hostRoleList, sql, clusterId)
	if err != nil {
		return hostRoleList, err
	}
	return hostRoleList, nil
}

func (h *hostRole) GetRoleInfoByRoleNameAndClusterId(clusterId int, roleName string) (*HostRoleInfo, error) {
	sql := fmt.Sprintf("select * from %s where cluster_id = ? and role_name= ? ", HOST_ROLE)
	var hostRoleInfo HostRoleInfo
	err := h.GetDB().Get(&hostRoleInfo, sql, clusterId, roleName)
	if err != nil {
		return nil, err
	}
	return &hostRoleInfo, nil
}

func (h *hostRole) ReNameByRoleId(roleId int, roleName string) error {
	sql := fmt.Sprintf("update %s set role_name = ? where id = ? ", HOST_ROLE)
	_, err := h.GetDB().Exec(sql, roleName, roleId)
	if err != nil {
		return err
	}
	return nil
}

func (h *hostRole) DeleteRoleByClusterIdAndRoleId(roleId int) error {
	sql := fmt.Sprintf("delete from %s where id = ? ", HOST_ROLE)
	_, err := h.GetDB().Exec(sql, roleId)
	if err != nil {
		return err
	}
	return nil
}
