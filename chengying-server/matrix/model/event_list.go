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
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"strings"
	"time"
)

type eventList struct {
	dbhelper.DbTable
}

var EventList = &eventList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_EVENT_LIST},
}

type EventInfo struct {
	Id                int       `db:"id" json:"-"`
	EventType         string    `db:"event_type" json:"event_type"`
	ProductName       string    `db:"product_name" json:"product_name"`
	ParentProductName string    `db:"parent_product_name" json:"parent_product_name"`
	ServiceName       string    `db:"service_name" json:"service_name"`
	Host              string    `db:"host" json:"host"`
	Content           string    `db:"content" json:"content"`
	IsDeleted         int       `db:"isDeleted" json:"-"`
	UpdateDate        base.Time `db:"update_time" json:"update_time"`
	CreateDate        base.Time `db:"create_time" json:"create_time"`
}

var _getEventListFields = utils.GetTagValues(EventInfo{}, "db")

func (e *eventList) NewEvent(parentProductName, productName, serviceName, host, eventType string, content string) error {
	_, err := e.InsertWhere(dbhelper.UpdateFields{
		"event_type":          eventType,
		"parent_product_name": parentProductName,
		"product_name":        productName,
		"service_name":        serviceName,
		"host":                host,
		"content":             content,
		"update_time":         time.Now(),
		"create_time":         time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func _eventListWhereCauseIn(param, name string, whereCause *dbhelper.WhereCause) {
	var values []interface{}

	if param != "" {
		for _, v := range strings.Split(param, ",") {
			values = append(values, v)
		}
		*whereCause = whereCause.And()
		*whereCause = whereCause.Included(name, values...)
	}
}

func _eventListWhereCause(eventType, parentProductName, productNames, serviceNames, hosts, startTime, endTime, keyWord string) dbhelper.WhereCause {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("deploy_instance_runtime_event.isDeleted", 0)
	if eventType != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("deploy_instance_runtime_event.event_type", eventType)
	}
	if parentProductName != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("deploy_instance_runtime_event.parent_product_name", parentProductName)
	}

	_eventListWhereCauseIn(productNames, "deploy_instance_runtime_event.product_name", &whereCause)
	_eventListWhereCauseIn(serviceNames, "deploy_instance_runtime_event.service_name", &whereCause)
	_eventListWhereCauseIn(hosts, "deploy_instance_runtime_event.host", &whereCause)

	whereCause = whereCause.And()
	whereCause = whereCause.Like("deploy_instance_runtime_event.content", "%"+keyWord+"%")

	whereCause = whereCause.And()
	whereCause = whereCause.Between("create_time", startTime, endTime)

	return whereCause
}
func (e *eventList) SelectEventListByWhere(pagination *apibase.Pagination, eventType, parentProductName, productNames, serviceNames, hosts, startTime, endTime, keyWord string) ([]EventInfo, error) {
	whereCause := _eventListWhereCause(eventType, parentProductName, productNames, serviceNames, hosts, startTime, endTime, keyWord)

	rows, _, err := e.SelectWhere(_getEventListFields, whereCause, pagination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	eventInfoList := []EventInfo{}
	for rows.Next() {
		info := EventInfo{}
		if err = rows.StructScan(&info); err != nil {
			return nil, err
		}
		eventInfoList = append(eventInfoList, info)
	}
	return eventInfoList, nil
}

func (e *eventList) GetEventListByWhere(eventType, parentProductName, productNames, serviceNames, hosts, startTime, endTime string) (*eventList, error) {
	whereCause := _eventListWhereCause(eventType, parentProductName, productNames, serviceNames, hosts, startTime, endTime, "")
	info := eventList{}
	err := e.GetWhere(nil, whereCause, &info)
	if err != nil {
		return nil, err
	}
	return &info, err
}
