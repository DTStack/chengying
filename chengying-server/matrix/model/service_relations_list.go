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
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"time"
)

const (
	RELATIONS_TYPE_CONFLICT string = "conflict"
	RELATIONS_TYPE_RELYON   string = "relyOn"
)

type deployServiceRelationsList struct {
	dbhelper.DbTable
}

var DeployServiceRelationsList = &deployServiceRelationsList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_SERVICE_RELATIONS_LIST},
}

type DeployServiceRelationsInfo struct {
	Id                int          `db:"id"`
	RelationsType     string       `db:"relations_type"`
	SourceProductName string       `db:"source_product_name"`
	SourceServiceName string       `db:"source_service_name"`
	TargetProductName string       `db:"target_product_name"`
	TargetServiceName string       `db:"target_service_name"`
	CreateTime        sql.NullTime `db:"create_time"`
	UpdateTime        sql.NullTime `db:"update_time"`
	IsDeleted         int          `db:"is_deleted"`
}

func (l *deployServiceRelationsList) InsertDeployServiceRelationsRecord(info DeployServiceRelationsInfo) (int64, error) {
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("is_deleted", 0).And().
		Equal("relations_type", info.RelationsType).And().
		Equal("source_product_name", info.SourceProductName).And().
		Equal("source_service_name", info.SourceServiceName).And().
		Equal("target_product_name", info.TargetProductName).And().
		Equal("target_service_name", info.TargetServiceName), &info)
	if err != nil && err == sql.ErrNoRows {
		r, err := l.InsertWhere(dbhelper.UpdateFields{
			"relations_type":      info.RelationsType,
			"source_product_name": info.SourceProductName,
			"source_service_name": info.SourceServiceName,
			"target_product_name": info.TargetProductName,
			"target_service_name": info.TargetServiceName,
			"create_time":         time.Now(),
			"update_time":         time.Now(),
			"is_deleted":          0,
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

var _getServiceRelationsListFields = utils.GetTagValues(DeployServiceRelationsInfo{}, "db")

func (l *deployServiceRelationsList) GetServiceRelationsList() (error, []DeployServiceRelationsInfo) {
	list := make([]DeployServiceRelationsInfo, 0)
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	rows, _, err := l.SelectWhere(_getServiceRelationsListFields, whereCause, nil)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		info := DeployServiceRelationsInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return err, nil
		}
		list = append(list, info)
	}
	return nil, list
}
