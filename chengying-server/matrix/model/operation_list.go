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
	"database/sql"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"fmt"
	"github.com/jmoiron/sqlx"
)

/*
 @Author: zhijian
 @Date: 2021/6/2 14:46
 @Description:
*/
type operationList struct {
	dbhelper.DbTable
}
type OperationInfo struct {
	Id                int               `db:"id" json:"-"`
	ClusterId         int               `db:"cluster_id" json:"clusterId"`
	OperationId       string            `db:"operation_id" json:"operationId"`
	OperationType     int               `db:"operation_type" json:"operationType"`
	OperationStatus   int               `db:"operation_status" json:"operationStatus"`
	ObjectType        int               `db:"object_type" json:"objectType"`
	ObjectValue       string            `db:"object_value" json:"objectValue"`
	CreateTime        dbhelper.NullTime `db:"create_time" json:"-"`
	EndTime           dbhelper.NullTime `db:"end_time"  json:"-"`
	Duration          sql.NullFloat64   `db:"duration" json:"-"`
	UpdateTime        dbhelper.NullTime `db:"update_time" json:"-"`
	OperationName     string            `json:"operationName"`
	ShowCreateTime    string            `json:"startTime"`
	ShowEndTime       string            `json:"endTime"`
	ShowDuration      float64           `json:"duration"`
	Group             string            `json:"group"`
	ProductName       string            `json:"productName"`
	IsExist           bool              `json:"isExist"`
	ParentProductName string            `json:"parentProductName"`
}

var OperationList = &operationList{
	dbhelper.DbTable{USE_MYSQL_DB, OPERATION_LIST},
}

func (opl *operationList) GetRunningCount(clusterId int) (*int, error) {
	sql := fmt.Sprintf("select count(1) from %s where operation_status = ? and cluster_id= ?", OPERATION_LIST)
	var count int
	err := opl.GetDB().Get(&count, sql, enums.ExecStatusType.Running.Code, clusterId)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (opl *operationList) Insert(one OperationInfo) error {
	sql := fmt.Sprintf("INSERT INTO %s (cluster_id,operation_id,operation_type, operation_status,object_type,object_value) VALUES (?,?,?,?,?,?)", OPERATION_LIST)
	_, err := opl.GetDB().Exec(sql, one.ClusterId, one.OperationId, one.OperationType, one.OperationStatus, one.ObjectType, one.ObjectValue)
	if err != nil {
		return err
	}
	return nil
}

func (opl *operationList) InsertWithTx(tx *sqlx.Tx, one OperationInfo) error {
	sql := fmt.Sprintf("INSERT INTO %s (cluster_id,operation_id,operation_type, operation_status,object_type,object_value) VALUES (?,?,?,?,?,?)", OPERATION_LIST)
	_, err := tx.Exec(sql, one.ClusterId, one.OperationId, one.OperationType, one.OperationStatus, one.ObjectType, one.ObjectValue)
	if err != nil {
		return err
	}
	return nil
}

func (opl *operationList) GetByOperationId(operationId string) (*OperationInfo, error) {
	sql := fmt.Sprintf("select * from %s where operation_id = ?", OPERATION_LIST)
	var one OperationInfo
	err := opl.GetDB().Get(&one, sql, operationId)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

func (opl *operationList) UpdateStatusByOperationId(operationId string, status int, endTime dbhelper.NullTime, duration sql.NullFloat64) error {
	sql := fmt.Sprintf("update %s set operation_status  = ? ,end_time= ? , duration= ?  where operation_id = ? ", OPERATION_LIST)
	_, err := opl.GetDB().Exec(sql, status, endTime, duration, operationId)
	if err != nil {
		return err
	}
	return nil
}
func (opl *operationList) GetByOperationTypeAndObjectValue(operationType int, objectValue string) (*OperationInfo, error) {
	sql := fmt.Sprintf("select * from %s where operation_type = ? and object_value = ?  order by create_time desc limit 1", OPERATION_LIST)
	var one OperationInfo
	err := opl.GetDB().Get(&one, sql, operationType, objectValue)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

func (opl *operationList) ListObjectValue(clusterId int) ([]string, error) {
	sql := fmt.Sprintf("select distinct object_value from %s where cluster_id = ? ", OPERATION_LIST)
	var objectValueList []string
	err := opl.GetDB().Select(&objectValueList, sql, clusterId)
	if err != nil {
		return nil, err
	}
	return objectValueList, nil
}

func (opl *operationList) DeleteBySid(sid string) error {
	sql := fmt.Sprintf("delete  from %s where object_value = ? ", OPERATION_LIST)
	_, err := opl.GetDB().Exec(sql, sid)
	if err != nil {
		return err
	}
	return nil
}
