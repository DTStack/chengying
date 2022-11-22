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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"fmt"
	"strconv"
	"time"
)

type healthCheck struct {
	dbhelper.DbTable
}

var HealthCheck = &healthCheck{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_SERVICE_HEALTH_CHECK},
}

type HealthCheckInfo struct {
	ID                int               `db:"id"`
	ClusterId         int               `db:"cluster_id"`
	ProductName       string            `db:"product_name"`
	Pid               int               `db:"pid"`
	ServiceName       string            `db:"service_name"`
	AgentId           string            `db:"agent_id"`
	Sid               string            `db:"sid"`
	Ip                string            `db:"ip"`
	ScriptName        string            `db:"script_name"`
	ScriptNameDisplay string            `db:"script_name_display"`
	AutoExec          bool              `db:"auto_exec"`
	Period            string            `db:"period"`
	Retries           int               `db:"retries"`
	ExecStatus        int               `db:"exec_status"`
	ErrorMessage      string            `db:"error_message"`
	StartTime         dbhelper.NullTime `db:"start_time"`
	EndTime           dbhelper.NullTime `db:"end_time"`
	CreateTime        dbhelper.NullTime `db:"create_time"`
	UpdateTime        dbhelper.NullTime `db:"update_time"`
}

func (h *healthCheck) GetInfoByClusterIdAndProductNameAndServiceName(clusterId int, productName, serviceName, hostIp string) ([]HealthCheckInfo, error) {
	s := "select * from " + TBL_SERVICE_HEALTH_CHECK + " where cluster_id = " + "'" + strconv.Itoa(clusterId) + "'" +
		" and product_name = " + "'" + productName + "'" + " and service_name = " + "'" + serviceName + "'"
	if len(hostIp) != 0 {
		s = s + " and ip = '" + hostIp + "'"
	}
	s = s + " order by start_time desc"
	rows, err := USE_MYSQL_DB().Queryx(s)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Errorf("%v", err)
		return nil, err
	}
	list := make([]HealthCheckInfo, 0)
	for rows.Next() {
		info := HealthCheckInfo{}
		if err := rows.StructScan(&info); err != nil {
			log.Errorf("%v", err)
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (h *healthCheck) GetInfoById(id int) (HealthCheckInfo, error) {
	info := HealthCheckInfo{}
	err := h.GetWhere(nil, dbhelper.MakeWhereCause().Equal("id", id), &info)
	return info, err
}

func (h *healthCheck) UpdateAutoexecById(id int, autoexec bool) error {
	err := h.UpdateWhere(dbhelper.MakeWhereCause().
		Equal("id", id), dbhelper.UpdateFields{
		"auto_exec":   autoexec,
		"update_time": time.Now(),
	}, false)
	return err
}

func (h *healthCheck) UpdateHealthCheckStatus(id, execStatus int, errorMessage string, startTime, endTime dbhelper.NullTime) error {
	err := h.UpdateWhere(dbhelper.MakeWhereCause().
		Equal("id", id), dbhelper.UpdateFields{
		"exec_status":   execStatus,
		"error_message": errorMessage,
		"start_time":    startTime,
		"end_time":      endTime,
	}, false)
	return err
}

func (h *healthCheck) DeleteByIp(ip string) error {
	deleteSqlStr := fmt.Sprintf("DELETE FROM %s WHERE ip = ? ", h.TableName)
	_, err := h.GetDB().Exec(deleteSqlStr, ip)
	if err != nil {
		return err
	}
	return nil
}
