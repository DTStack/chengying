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
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/dag"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"fmt"
	"time"
)

type deployProductLineList struct {
	dbhelper.DbTable
}

var DeployProductLineList = &deployProductLineList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_PRODUCT_LINE_LIST},
}

type DeployProductLineInfo struct {
	ID                 int               `db:"id" json:"id"`
	ProductLineName    string            `db:"product_line_name" json:"product_line_name"`
	ProductLineVersion string            `db:"product_line_version" json:"product_line_version"`
	ProductSerial      []byte            `db:"product_serial" json:"product_serial"`
	CreateTime         dbhelper.NullTime `db:"create_time" json:"create_time"`
	UpdateTime         dbhelper.NullTime `db:"update_time" json:"update_time"`
	IsDeleted          int               `db:"is_deleted" json:"is_deleted"`
}

type ProductLineTemplate struct {
	ProductLineName    string          `json:"product_line_name" validation:"required"`
	ProductLineVersion string          `json:"product_line_version" validation:"required"`
	ProductSerial      []ProductSerial `json:"product_serial" validation:"required"`
}

type ProductSerial struct {
	ID          dag.Node `json:"id" validation:"required"`
	ProductName string   `json:"product_name" validation:"required"`
	Dependee    dag.Node `json:"dependee"`
}

func (l *deployProductLineList) InsertProductLineIfNotExist(productLineName, productLineVersion string, productSerial []byte) (error, int) {
	info := DeployProductLineInfo{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("product_line_name", productLineName).And().
		Equal("product_line_version", productLineVersion).And().
		Equal("is_deleted", 0), &info)
	if err != nil && err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"product_line_name":    productLineName,
			"product_line_version": productLineVersion,
			"product_serial":       productSerial,
		})
		if err != nil {
			apibase.ThrowDBModelError(err)
			return err, -1
		}
		seq, _ := ret.LastInsertId()
		return nil, int(seq)
	} else if err == nil {
		aid := info.ID
		return fmt.Errorf("%v(%v) 已存在", productLineName, productLineVersion), aid
	} else {
		return err, 0
	}
}

func (l *deployProductLineList) GetProductLineList(pagination *apibase.Pagination) ([]DeployProductLineInfo, int) {
	fields := []string{"id", "product_line_name", "product_line_version", "product_serial", "create_time", "update_time"}
	whereCause := dbhelper.MakeWhereCause().GreaterThan("id", "0")
	whereCause = whereCause.And().Equal("is_deleted", 0)
	rows, total, err := l.SelectWhere(fields, whereCause, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()
	list := make([]DeployProductLineInfo, 0)
	for rows.Next() {
		info := DeployProductLineInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return list, total
}

func (l *deployProductLineList) GetProductLineListByNameAndVersion(productLineName, productLineVersion string) (*DeployProductLineInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("product_line_name", productLineName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_line_version", productLineVersion)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("is_deleted", 0)
	info := DeployProductLineInfo{}
	if err := l.GetWhere(nil, whereCause, &info); err != nil {
		return &info, err
	}
	return &info, nil
}

func (l *deployProductLineList) DeleteProductLineById(id int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	return err
}
