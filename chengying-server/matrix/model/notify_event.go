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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"time"

	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
)

type notifyEvent struct {
	dbhelper.DbTable
}

var NotifyEvent = &notifyEvent{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_NOTIFY_EVENT},
}

type NotifyEventInfo struct {
	Id                int               `db:"id"`
	ClusterId         int               `db:"cluster_id"`
	Type              int               `db:"type"`
	ProductName       string            `db:"product_name"`
	ServiceName       string            `db:"service_name"`
	DependProductName string            `db:"depend_product_name"`
	DependServiceName string            `db:"depend_service_name"`
	Host              string            `db:"host"`
	IsDeleted         int               `db:"is_deleted"`
	CreateTime        dbhelper.NullTime `db:"create_time"`
	UpdateTime        dbhelper.NullTime `db:"update_time"`
	ProductStopped    int               `db:"product_stopped"`
}

func (e *notifyEvent) InsertNotifyEvent(clusterId, eventType int, productName, serviceName, dependProductName,
	dependServiceName, host string) error {

	err := e.UpdateWhere(dbhelper.MakeWhereCause().Equal("cluster_id", clusterId).
		And().Equal("type", eventType).And().Equal("product_name", productName).
		And().Equal("service_name", serviceName).And().Equal("depend_product_name", dependProductName).
		And().Equal("depend_service_name", dependServiceName).And().Equal("host", host).And().Equal("is_deleted", 0),
		dbhelper.UpdateFields{
			"cluster_id":          clusterId,
			"type":                eventType,
			"product_name":        productName,
			"service_name":        serviceName,
			"depend_product_name": dependProductName,
			"depend_service_name": dependServiceName,
			"host":                host,
			"update_time":         time.Now(),
		}, true)
	if err != nil {
		return err
	}
	return nil
}

func (e *notifyEvent) DeleteNotifyEvent(clusterId, eventType int, productName, serviceName, host string, isProduct bool) error {
	whereClause := dbhelper.MakeWhereCause().Equal("cluster_id", clusterId).And().Equal("type", eventType).
		And().Equal("product_name", productName).And().Equal("is_deleted", 0)
	if serviceName != "" {
		whereClause = whereClause.And().Equal("service_name", serviceName)
	}
	if host != "" {
		whereClause = whereClause.And().Equal("host", host)
	}
	if isProduct {
		whereClause = whereClause.And().Equal("product_stopped", 1)
	}
	return e.UpdateWhere(whereClause, dbhelper.UpdateFields{
		"is_deleted": 1,
		"update_time": time.Now(),
	}, false)

}

func (e *notifyEvent) UpdateProductStopped(clusterId, eventType int, productName string) error {
	whereClause := dbhelper.MakeWhereCause().Equal("cluster_id", clusterId).And().Equal("type", eventType).
		And().Equal("product_name", productName).And().Equal("is_deleted", 0)
	return e.UpdateWhere(whereClause, dbhelper.UpdateFields{
		"product_stopped": 1,
		//"update_time":     time.Now(),
	}, false)
}

func (e *notifyEvent) GetServiceLastStartTime(clusterId int, productName, serviceName string) (*time.Time,error) {
	event := NotifyEventInfo{}
	query := "SELECT * FROM " + e.TableName + " where cluster_id = ? and product_name = ? and service_name = ? and is_deleted = 1 order by id desc limit 1"
	if err := e.GetDB().Get(&event, query, clusterId, productName, serviceName); err != nil {
		log.Debugf("[notifyEvent.GetServiceLastStartTime] %s", err)
		return nil, err
	}
	return &event.CreateTime.Time,nil
}
