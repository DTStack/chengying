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

package upgrade

import (
	"database/sql"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"time"
)

type smoothUpgrade struct {
	dbhelper.DbTable
}

var SmoothUpgrade = &smoothUpgrade{
	dbhelper.DbTable{
		GetDB:     model.USE_MYSQL_DB,
		TableName: model.TBL_DEPLOY_SMOOTH_UPGRADE_LIST,
	},
}

type SmoothUpgradeInfo struct {
	Id          int          `db:"id" json:"id"`
	ProductName string       `db:"product_name" json:"product_name"`
	ServiceName string       `db:"service_name" json:"service_name"`
	CreateTime  sql.NullTime `db:"create_time" json:"create_time"`
	IsDeleted   int          `db:"is_deleted" json:"is_deleted"`
}

func (s *smoothUpgrade) InsertSmoothUpgradeRecord(productName, serviceName string) (int64, error) {
	var info SmoothUpgradeInfo
	err := s.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("is_deleted", 0).And().
		Equal("product_name", productName).And().
		Equal("service_name", serviceName), &info)
	if err != nil && err == sql.ErrNoRows {
		r, err := s.InsertWhere(dbhelper.UpdateFields{
			"product_name": productName,
			"service_name": serviceName,
			"create_time":  time.Now(),
			"is_deleted":   0,
		})
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	} else if err == nil {
		return int64(info.Id), nil
	} else {
		return 0, err
	}
}

func (s *smoothUpgrade) GetByProductName(productName string) ([]SmoothUpgradeInfo, error) {
	whereClause := dbhelper.MakeWhereCause().
		Equal("product_name", productName).And().
		Equal("is_deleted", 0)
	rows, _, err := s.SelectWhere(nil, whereClause, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("rows close error: %v", err)
			return
		}
	}()
	var infoList []SmoothUpgradeInfo
	for rows.Next() {
		row := SmoothUpgradeInfo{}
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, row)
	}
	return infoList, nil
}
