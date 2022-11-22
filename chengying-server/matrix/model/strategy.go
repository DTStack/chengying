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
	"dtstack.com/dtstack/easymatrix/go-common/utils"
)

type strategyList struct {
	dbhelper.DbTable
}

var StrategyList = &strategyList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_STRATEGY_LIST},
}

type StrategyInfo struct {
	ID           int               `db:"id"`
	Name         string            `db:"name"`
	Desc         string            `db:"desc"`
	Property     int               `db:"property"`
	StrategyType int               `db:"strategy_type"`
	DeployStatus int               `db:"deploy_status"`
	ExeStatus    int               `db:"exe_status"`
	ErrorMsg     string            `db:"error_message"`
	StartDate    dbhelper.NullTime `db:"start_date"`
	EndDate      dbhelper.NullTime `db:"end_date"`
	StartTime    dbhelper.NullTime `db:"start_time"`
	EndTime      dbhelper.NullTime `db:"end_time"`
	CronPeriod   int               `db:"cron_period"`
	CronInterval int               `db:"cron_interval"`
	CronTime     dbhelper.NullTime `db:"cron_time"`
	Params       string            `db:"params"`
	TimeOut      int               `db:"time_out"`
	IsDeleted    int               `db:"is_deleted"`
	GmtCreate    dbhelper.NullTime `db:"gmt_create"`
	GmtModify    dbhelper.NullTime `db:"gmt_modified"`
}

var _getStrategyListFields = utils.GetTagValues(StrategyInfo{}, "db")

func (l *strategyList) GetStrategyList() (error, []*StrategyInfo) {
	list := []*StrategyInfo{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	rows, _, err := l.SelectWhere(_getStrategyListFields, whereCause, nil)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		info := &StrategyInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return err, nil
		}
		list = append(list, info)
	}
	return nil, list
}

func (l *strategyList) GetDeployedStrategyList() (error, []*StrategyInfo) {
	list := []*StrategyInfo{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("deploy_status", 1)
	rows, _, err := l.SelectWhere(_getStrategyListFields, whereCause, nil)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		info := &StrategyInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return err, nil
		}
		list = append(list, info)
	}
	return nil, list
}
