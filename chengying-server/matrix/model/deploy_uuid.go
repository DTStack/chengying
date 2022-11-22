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
)

/*
 @Author: zhijian
 @Date: 2021/4/15 22:15
 @Description:
*/

type deployUUID struct {
	dbhelper.DbTable
}
type DeployUUIDInfo struct {
	Id         int               `db:"id" json:"id"`
	UUID       string            `db:"uuid" json:"uuid"`
	ParentUUID string            `db:"parent_uuid" json:"parent_uuid"`
	UuidType   int               `db:"type" json:"type"`
	Pid        string            `db:"pid" json:"pid"` //  type 为自动部署的时候  pid 为本次自动部署的所有 pidlist  格式 1,2,3
	CreateTime dbhelper.NullTime `db:"create_time"`
}

const (
	// ManualDeployUUIDType 手动
	ManualDeployUUIDType = 1
	// AutoDeployUUIDType 自动部署
	AutoDeployUUIDType = 2
	// AutoDeployChildrenUUIDType 自动部署中子产品部署时的 UUID
	AutoDeployChildrenUUIDType = 3
)

var DeployUUID = &deployUUID{
	dbhelper.DbTable{USE_MYSQL_DB, DEPLOY_UUID},
}

func (d *deployUUID) GetInfoByUUID(UUID string) (*DeployUUIDInfo, error) {
	var deployUUID DeployUUIDInfo
	sql := fmt.Sprintf("select * from %s where uuid = ?", DEPLOY_UUID)
	err := d.Get(&deployUUID, sql, UUID)
	if err != nil {
		return nil, err
	}
	return &deployUUID, nil
}

func (d *deployUUID) InsertOne(UUID, parentUUID string, UUIDType, pid int) error {
	sql := fmt.Sprintf("insert into %s ( uuid, type, parent_uuid,pid) values (?,?,?,?);", DEPLOY_UUID)
	_, err := d.GetDB().Exec(sql, UUID, UUIDType, parentUUID, pid)
	if err != nil {
		return err
	}
	return nil
}
func (d *deployUUID) GetUUIDListByParentUUID(parentUUID string) ([]DeployUUIDInfo, error) {
	var UUIDInfoList []DeployUUIDInfo
	sql := fmt.Sprintf("select * from %s where parent_uuid = ?", DEPLOY_UUID)
	err := d.GetDB().Select(&UUIDInfoList, sql, parentUUID)
	if err != nil {
		return nil, err
	}
	return UUIDInfoList, nil
}

func (d *deployUUID) SetPidByUUID(UUID, pid string) error {
	sql := fmt.Sprintf("update %s set pid = ? where uuid = ?", DEPLOY_UUID)
	_, err := d.GetDB().Exec(sql, pid, UUID)
	if err != nil {
		return err
	}
	return nil
}
