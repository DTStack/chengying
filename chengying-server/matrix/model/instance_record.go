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
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

const (
	INCLUDE_SEP = ","
)

type DeployInstanceUpdateRecordInfo struct {
	ID                 int               `db:"id"`
	UpdateUUID         uuid.UUID         `db:"update_uuid"`
	InstanceId         int               `db:"instance_id"`
	Sid                string            `db:"sid"`
	Ip                 string            `db:"ip"`
	ProductName        string            `db:"product_name"`
	ProductNameDisplay string            `db:"product_name_display"`
	ProductVersion     string            `db:"product_version"`
	Group              string            `db:"group"`
	ServiceName        string            `db:"service_name"`
	ServiceNameDisplay string            `db:"service_name_display"`
	ServiceVersion     string            `db:"service_version"`
	Status             string            `db:"status"`
	StatusMessage      string            `db:"status_message"`
	Progress           uint              `db:"progress"`
	UpdateDate         dbhelper.NullTime `db:"update_time"`
	CreateDate         dbhelper.NullTime `db:"create_time"`
}

type deployInstanceUpdateRecord struct {
	dbhelper.DbTable
}

var DeployInstanceUpdateRecord = &deployInstanceUpdateRecord{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_INSTANCE_UPDATE_RECORD},
}

type DeployInstanceRecordInfo struct {
	ID                 int               `db:"id"`
	DeployUUID         uuid.UUID         `db:"deploy_uuid"`
	InstanceId         int               `db:"instance_id"`
	Sid                string            `db:"sid"`
	Ip                 string            `db:"ip"`
	ProductName        string            `db:"product_name"`
	ProductNameDisplay string            `db:"product_name_display"`
	ProductVersion     string            `db:"product_version"`
	Group              string            `db:"group"`
	ServiceName        string            `db:"service_name"`
	ServiceNameDisplay string            `db:"service_name_display"`
	ServiceVersion     string            `db:"service_version"`
	Status             string            `db:"status"`
	StatusMessage      string            `db:"status_message"`
	Progress           uint              `db:"progress"`
	UpdateDate         dbhelper.NullTime `db:"update_time"`
	CreateDate         dbhelper.NullTime `db:"create_time"`
}

type deployInstanceRecord struct {
	dbhelper.DbTable
}

var DeployInstanceRecord = &deployInstanceRecord{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_INSTANCE_RECORD},
}

func (l *deployInstanceRecord) GetInstanceInfoByWhere(cause dbhelper.WhereCause) (error, *DeployInstanceRecordInfo) {
	info := DeployInstanceRecordInfo{}
	err := l.GetWhere(nil, cause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployInstanceRecord) CreateOrUpdate(d *DeployInstanceRecordInfo) (error, int64, string) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("deploy_uuid", d.DeployUUID)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("instance_id", d.InstanceId)

	err, info := l.GetInstanceInfoByWhere(whereCause)
	var  id int64
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"deploy_uuid":          d.DeployUUID,
			"instance_id":          d.InstanceId,
			"sid":                  d.Sid,
			"ip":                   d.Ip,
			"product_name":         d.ProductName,
			"product_name_display": d.ProductNameDisplay,
			"product_version":      d.ProductVersion,
			"group":                d.Group,
			"service_name":         d.ServiceName,
			"service_name_display": d.ServiceNameDisplay,
			"service_version":      d.ServiceVersion,
			"status":               d.Status,
			"status_message":       d.StatusMessage,
			"progress":             d.Progress,
			"update_time":          time.Now(),
			"create_time":          time.Now(),
		})
		if err != nil {
			return err, -1, d.DeployUUID.String()
		}
		id, _ = ret.LastInsertId()
	} else {
		id  = int64(info.ID)
		err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"sid":                  d.Sid,
			"ip":                   d.Ip,
			"product_name":         d.ProductName,
			"product_name_display": d.ProductNameDisplay,
			"product_version":      d.ProductVersion,
			"group":                d.Group,
			"service_name":         d.ServiceName,
			"service_name_display": d.ServiceNameDisplay,
			"service_version":      d.ServiceVersion,
			"status":               d.Status,
			"status_message":       d.StatusMessage,
			"progress":             d.Progress,
			"update_time":          time.Now(),
		}, false)
		if err != nil {
			return err, int64(info.ID), d.DeployUUID.String()
		}
	}
	return nil, id, d.DeployUUID.String()
}

func (r *deployInstanceRecord) GetDeployInstanceRecordByDeployId(pagination *apibase.Pagination, deployUUID string, status []string, serviceName string) ([]DeployInstanceRecordByDeployIdInfo, int, string) {
	whereCause := dbhelper.MakeWhereCause()
	whereCause = whereCause.Equal("deploy_uuid", deployUUID)
	var values []interface{}
	if serviceName != "" {
		for _, v := range strings.Split(serviceName, ",") {
			values = append(values, v)
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("IR.service_name", values...)
	}
	if len(status) > 0 {
		s := make([]interface{}, 0, len(status))
		for _, s_ := range status {
			s = append(s, interface{}(s_))
		}
		whereCause = whereCause.And().Included("IR.status", s...)
	}
	where, value := whereCause.SQL()
	query := "SELECT IR.*, IL.schema FROM " +
		DeployInstanceRecord.TableName + " AS IR LEFT JOIN " +
		DeployInstanceList.TableName + " AS IL ON IR.instance_id = IL.id " + where + " " + pagination.AsQuery()
	queryCount := "SELECT COUNT(*) FROM " + DeployInstanceRecord.TableName + " AS IR " + where

	var list []DeployInstanceRecordByDeployIdInfo
	var count int
	if err := DeployInstanceRecord.GetDB().Get(&count, queryCount, value...); err != nil {
		log.Errorf("queryCount: %v, value: %v, err: %v", queryCount, value, err)
		apibase.ThrowDBModelError(err)
	}

	if count > 0 {
		rows, err := USE_MYSQL_DB().Queryx(query, value...)
		if err != nil {
			log.Errorf("query: %v, value: %v, err: %v", query, value, err)
			apibase.ThrowDBModelError(err)
		}

		defer rows.Close()

		for rows.Next() {
			info := DeployInstanceRecordByDeployIdInfo{}
			if err := rows.StructScan(&info); err != nil {
				apibase.ThrowDBModelError(err)
			}
			list = append(list, info)
		}
	}

	var complete string
	if err := DeployProductHistory.GetWhere([]string{"status"}, dbhelper.MakeWhereCause().Equal("deploy_uuid", deployUUID), &complete); err != nil {
		apibase.ThrowDBModelError(err)
	}

	return list, count, complete
}

func (r *deployInstanceRecord) GetDeployInstanceRecordById(id int) (info DeployInstanceRecordInfo, err error) {
	err = r.GetWhere(nil, dbhelper.MakeWhereCause().Equal("id", id), &info)
	return
}

func (r *deployInstanceRecord) GetDeployInstanceRecordByInstanceId(instanceId int) (info DeployInstanceRecordInfo, err error) {
	err = r.GetWhere(nil, dbhelper.MakeWhereCause().Equal("instance_id", instanceId), &info)
	return
}

func (r *deployInstanceRecord) GetDeployInstanceRecordByInstanceIdAndStatus(instanceId int, status string) (info DeployInstanceRecordInfo, err error) {
	err = r.GetWhere(nil, dbhelper.MakeWhereCause().Equal("instance_id", instanceId).And().Equal("status", status), &info)
	return
}

func (r *deployInstanceRecord) UpdateDeployInstanceRecord(id int, fields dbhelper.UpdateFields) error {
	return r.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), fields, false)
}

func (r *deployInstanceUpdateRecord) GetDeployInstanceUpdateRecordByUpdateId(pagination *apibase.Pagination, updateUUID string, status []string, serviceName string) ([]DeployInstanceUpdateRecordByUpdateIdInfo, int, string) {
	whereCause := dbhelper.MakeWhereCause()
	whereCause = whereCause.Equal("update_uuid", updateUUID)
	var values []interface{}
	if serviceName != "" {
		for _, v := range strings.Split(serviceName, ",") {
			values = append(values, v)
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("IR.service_name", values...)
	}
	if len(status) > 0 {
		s := make([]interface{}, 0, len(status))
		for _, s_ := range status {
			s = append(s, interface{}(s_))
		}
		whereCause = whereCause.And().Included("IR.status", s...)
	}
	where, value := whereCause.SQL()
	query := "SELECT IR.*, IL.schema FROM " +
		DeployInstanceUpdateRecord.TableName + " AS IR LEFT JOIN " +
		DeployInstanceList.TableName + " AS IL ON IR.instance_id = IL.id " + where + " " + pagination.AsQuery()
	queryCount := "SELECT COUNT(*) FROM " + DeployInstanceUpdateRecord.TableName + " AS IR " + where

	var list []DeployInstanceUpdateRecordByUpdateIdInfo
	var count int
	if err := DeployInstanceUpdateRecord.GetDB().Get(&count, queryCount, value...); err != nil {
		log.Errorf("queryCount: %v, value: %v, err: %v", queryCount, value, err)
		apibase.ThrowDBModelError(err)
	}

	if count > 0 {
		rows, err := USE_MYSQL_DB().Queryx(query, value...)
		if err != nil {
			log.Errorf("query: %v, value: %v, err: %v", query, value, err)
			apibase.ThrowDBModelError(err)
		}

		defer rows.Close()

		for rows.Next() {
			info := DeployInstanceUpdateRecordByUpdateIdInfo{}
			if err := rows.StructScan(&info); err != nil {
				apibase.ThrowDBModelError(err)
			}
			list = append(list, info)
		}
	}

	var complete string
	if err := DeployProductUpdateHistory.GetWhere([]string{"status"}, dbhelper.MakeWhereCause().Equal("update_uuid", updateUUID), &complete); err != nil {
		apibase.ThrowDBModelError(err)
	}

	return list, count, complete
}
