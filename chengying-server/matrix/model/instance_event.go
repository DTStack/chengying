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
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"time"
)

type DeployInstanceEventInfo struct {
	ID         int               `db:"id"`
	InstanceId int               `db:"instance_id"`
	EventType  string            `db:"event_type"`
	Content    string            `db:"content"`
	UpdateDate dbhelper.NullTime `db:"update_time"`
	CreateDate dbhelper.NullTime `db:"create_time"`
}

type deployInstanceEvent struct {
	dbhelper.DbTable
}

var DeployInstanceEvent = &deployInstanceEvent{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_INSTANCE_EVENT},
}

var _getDeployInstanceEventFields = utils.GetTagValues(DeployInstanceEventInfo{}, "db")

func (l *deployInstanceEvent) NewInstanceEvent(instanceId int64, eventType string, content string) error {
	_, err := l.InsertWhere(dbhelper.UpdateFields{
		"instance_id": instanceId,
		"event_type":  eventType,
		"content":     content,
		"update_time": time.Now(),
		"create_time": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *deployInstanceEvent) GetDeployInstanceEventById(pagination *apibase.Pagination, instanceId int64, eventType string) ([]DeployInstanceEventInfo, error) {
	whereCause := dbhelper.MakeWhereCause()
	whereCause = whereCause.Equal("instance_id", instanceId)
	if eventType != "" {
		whereCause = whereCause.And().Equal("event_type", eventType)
	}
	if pagination == nil {
		pagination = &apibase.Pagination{SortBy: "id", SortDesc: true}
	}
	rows, _, err := r.SelectWhere(_getDeployInstanceEventFields, whereCause, pagination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployInstanceEventInfo{}
	for rows.Next() {
		info := DeployInstanceEventInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployInstanceEvent) GetEventInfoByInstanceIdEventId(instanceId int64, eventId int64) (*DeployInstanceEventInfo, error) {
	whereCause := dbhelper.WhereCause{}
	info := DeployInstanceEventInfo{}
	err := l.GetWhere(nil, whereCause.Equal("id", eventId), &info)
	if err != nil {
		return &info, err
	}
	return &info, nil
}
