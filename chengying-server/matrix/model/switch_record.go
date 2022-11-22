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
	"time"
)

type switchRecord struct {
	dbhelper.DbTable
}

type SwitchRecordI interface {
	NewSwitchRecord(name, switchType, productName, serviceName, status, StatusMessage string, clusterId, progress int) (int64, error)
	GetRecordById(id int) (*SwitchRecordInfo, error)
	GetCurrentSwitchRecord(clusterId int, productName, serviceName, switchName string) (*SwitchRecordInfo, error)
}

var SwitchRecord = &switchRecord{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_SWITCH_RECORD},
}

type SwitchRecordInfo struct {
	Id            int               `db:"id"`
	ClusterId     int               `db:"cluster_id"`
	Name          string            `db:"name"`
	Type          string            `db:"switch_type"`
	ProductName   string            `db:"product_name"`
	ServiceName   string            `db:"service_name"`
	Status        string            `db:"status"`
	StatusMessage string            `db:"status_message"`
	Progress      int               `db:"progress"`
	CreateTime    dbhelper.NullTime `db:"create_time"`
	UpdateTime    dbhelper.NullTime `db:"update_time" json:"update_time"`
	IsDeleted     int               `db:"is_deleted"`
}

func (s *switchRecord) NewSwitchRecord(name, switchType, productName, serviceName, status, StatusMessage string, clusterId, progress int) (int64, error) {
	whereClause := dbhelper.MakeWhereCause().Equal("name", name).
		And().Equal("cluster_id", clusterId).
		And().Equal("product_name", productName).
		And().Equal("service_name", serviceName).
		And().Equal("is_deleted", 0).
		And().NotEqual("status", "FAIL").
		And().NotEqual("status", "SUCCESS")
	switchRecord := SwitchRecordInfo{}
	err := s.GetWhere(nil, whereClause, &switchRecord)
	if err == sql.ErrNoRows {
		result, err := s.InsertWhere(dbhelper.UpdateFields{
			"cluster_id":     clusterId,
			"name":           name,
			"switch_type":    switchType,
			"product_name":   productName,
			"service_name":   serviceName,
			"status":         status,
			"status_message": StatusMessage,
			"progress":       progress,
			"create_time":    time.Now(),
			"update_time":    time.Now(),
			"is_deleted":     0,
		})
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	} else {
		err := s.UpdateWhere(whereClause, dbhelper.UpdateFields{
			"status":         status,
			"status_message": StatusMessage,
			"progress":       progress,
			"update_time":    time.Now(),
		}, false)
		if err != nil {
			return 0, err
		}
		return int64(switchRecord.Id), nil
	}
}

func (s *switchRecord) GetRecordById(id int) (*SwitchRecordInfo, error) {
	whereClause := dbhelper.MakeWhereCause().Equal("id", id).And().Equal("is_deleted", 0)
	info := SwitchRecordInfo{}
	err := s.GetWhere(nil, whereClause, &info)
	if err != nil {
		return nil, err
	}
	return &info, err
}

func (s *switchRecord) GetCurrentSwitchRecord(clusterId int, productName, serviceName, switchName string) (*SwitchRecordInfo, error) {
	whereClause := dbhelper.MakeWhereCause().Equal("cluster_id", clusterId).And().
		Equal("product_name", productName).And().
		Equal("service_name", serviceName).And().
		Equal("name", switchName).And().
		Equal("is_deleted", 0).And().
		Equal("status", "RUNNING")
	var info SwitchRecordInfo
	if err := s.GetWhere(nil, whereClause, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
