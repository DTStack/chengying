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
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
)

type strategyResourcceList struct {
	dbhelper.DbTable
}

var StrategyResourceList = &strategyResourcceList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_STRATEGY_RESOURCE_LIST},
}

type StrategyResourceInfo struct {
	ID         int               `db:"id"`
	StrategyId int               `db:"strategy_id"`
	Content    string            `db:"content"`
	IsDeleted  int               `db:"is_deleted"`
	GmtCreate  dbhelper.NullTime `db:"gmt_create"`
	GmtModify  dbhelper.NullTime `db:"gmt_modified"`
}

var _getStrategyResourceListFields = utils.GetTagValues(StrategyResourceInfo{}, "db")

func (l *strategyResourcceList) GetStrategyResourceList() (error, []*StrategyResourceInfo) {
	list := []*StrategyResourceInfo{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	rows, _, err := l.SelectWhere(_getStrategyResourceListFields, whereCause, nil)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		info := &StrategyResourceInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return err, nil
		}
		list = append(list, info)
	}
	return nil, list
}
