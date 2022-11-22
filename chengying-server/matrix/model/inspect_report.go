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
	"time"
)

type inspectReport struct {
	dbhelper.DbTable
}

var InspectReport = &inspectReport{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_INSPECT_REPORT},
}

type InspectReportInfo struct {
	Id         int               `db:"id" json:"id"`
	Name       string            `db:"name" json:"name"`
	Status     string            `db:"status" json:"status"`
	Progress   int               `db:"progress" json:"progress"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
	ClusterId  int               `db:"cluster_id" json:"cluster_id"`
	FilePath   string            `db:"file_path" json:"file_path"`
}

func (i *inspectReport) NewInspectReport(name, status string, clusterId int) (int64, error) {
	result, err := i.InsertWhere(dbhelper.UpdateFields{
		"name":        name,
		"status":      status,
		"progress":    0,
		"create_time": time.Now(),
		"update_time": time.Now(),
		"is_deleted":  0,
		"cluster_id":  clusterId,
		"file_path":   "",
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (i *inspectReport) UpdateProgress(id, progress int, path, status string) error {
	updateFields := dbhelper.UpdateFields{
		"progress":    progress,
		"update_time": time.Now(),
	}
	if path != "" {
		updateFields["file_path"] = path
	}
	if status != "" {
		updateFields["status"] = status
	}
	err := i.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), updateFields, false)
	return err
}

func (i *inspectReport) GetById(id int) (InspectReportInfo, error) {
	var info InspectReportInfo
	if err := i.GetWhere(nil, dbhelper.MakeWhereCause().Equal("id", id), &info); err != nil {
		return info, err
	}
	return info, nil
}
