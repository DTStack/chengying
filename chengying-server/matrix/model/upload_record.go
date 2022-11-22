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
	"fmt"
	"time"
)

type uploadRecord struct {
	dbhelper.DbTable
}

var UploadRecord = &uploadRecord{
	dbhelper.DbTable{
		USE_MYSQL_DB, TBL_UPLOAD_RECORD,
	},
}

type UploadRecordInfo struct {
	Id         int               `db:"id"`
	UploadType int               `db:"upload_type"`
	Name       string            `db:"name"`
	Progress   float64           `db:"progress"`
	Status     string            `db:"status"`
	CreateTime dbhelper.NullTime `db:"create_time"`
	UpdateTime dbhelper.NullTime `db:"update_time"`
	IsDeleted  int               `db:"is_deleted"`
}

func (u *uploadRecord) AddUploadRecord(name, status string, uploadType int) (int64, error) {
	result, err := u.InsertWhere(dbhelper.UpdateFields{
		"upload_type": uploadType,
		"name":        name,
		"progress":    0,
		"status":      status,
		"create_time": time.Now(),
		"update_time": time.Now(),
		"is_deleted":  0,
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (u *uploadRecord) UpdateProgress(id int, progress float64, status string) error {
	whereClause := dbhelper.MakeWhereCause().Equal("id", id)
	err := u.UpdateWhere(whereClause, dbhelper.UpdateFields{
		"progress":    progress,
		"status":      status,
		"update_time": time.Now(),
	}, false)
	return err
}

func (u *uploadRecord) GetUploadingProduct(runningStatus, successStatus string) ([]UploadRecordInfo, error) {
	query := fmt.Sprintf("select * from %s where status in (?,?) and is_deleted=0", u.TableName)
	uploadingProductList := make([]UploadRecordInfo, 0)
	if err := u.GetDB().Select(&uploadingProductList, query, runningStatus, successStatus); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return uploadingProductList, nil
}

func (u *uploadRecord) DeleteUploadingProduct(uploadId int) error {
	query := fmt.Sprintf("update %s set is_deleted=1 where id=? and is_deleted=0", u.TableName)
	_, err := u.GetDB().Exec(query, uploadId)
	return err
}

func (u *uploadRecord) GetCancelProductById(uploadId int) (UploadRecordInfo, error) {
	whereClause := dbhelper.MakeWhereCause().Equal("id", uploadId).And().Equal("is_deleted", 1)
	var info UploadRecordInfo
	err := u.GetWhere(nil, whereClause, &info)
	return info, err
}
