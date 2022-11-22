/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"database/sql"
	apibase "easyagent/go-common/api-base"
	dbhelper "easyagent/go-common/db-helper"
	"easyagent/go-common/utils"

	uuid "github.com/satori/go.uuid"
)

type sidecarList struct {
	dbhelper.DbTable
}

type deployCallback struct {
	dbhelper.DbTable
}

var SidecarList = &sidecarList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_SIDECAR_LIST},
}

var DeployCallback = &deployCallback{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CALLBACK},
}

const (
	SIDECAR_UNKNOWN_STATUS = -1
	SIDECAR_NOT_INSTALLED  = iota
	SIDECAR_INSTALLING
	SIDECAR_INSTALL_FAILURE
	SIDECAR_INSTALLED
	SIDECAR_UPDATING
	SIDECAR_UPDATE_FAILURE
)

type SimpleSidecarInfo struct {
	ID         string            `db:"id"`
	Name       string            `db:"name"`
	Version    string            `db:"version"`
	Disabled   bool              `db:"disabled"`
	Status     int               `db:"status"`
	OsType     string            `db:"os_type"`
	DeployDate dbhelper.NullTime `db:"deploy_date"`
	AutoDeploy bool              `db:"auto_deployment"`
	UpdateDate dbhelper.NullTime `db:"last_update_date"`
	AutoUpdate bool              `db:"auto_updated"`
}

type SidecarInfo struct {
	SimpleSidecarInfo
	EcsFlag      bool           `db:"is_ecs"`
	OsPlatform   string         `db:"os_platform"`
	OsVersion    string         `db:"os_version"`
	Host         string         `db:"host"`
	LocalIp      string         `db:"local_ip"`
	CpuSerial    string         `db:"cpu_serial"`
	CpuCores     uint           `db:"cpu_cores"`
	MemSize      uint64         `db:"mem_size"`
	SwapSize     uint64         `db:"swap_size"`
	ServerHost   string         `db:"server_host"`
	ServerPort   int            `db:"server_port"`
	SshHost      string         `db:"ssh_host"`
	SshPort      int            `db:"ssh_port"`
	SshUser      string         `db:"ssh_user"`
	SshPwd       string         `db:"ssh_password"`
	CpuUsage     float32        `db:"cpu_usage"`
	MemUsage     int64          `db:"mem_usage"`
	SwapUsage    int64          `db:"swap_usage"`
	Load1        float32        `db:"load1"`
	Uptime       float64        `db:"uptime"`
	DiskUsage    sql.NullString `db:"disk_usage"`
	DiskUsagePct float32        `db:"disk_usage_pct"`
	NetUsage     sql.NullString `db:"net_usage"`
}

type DeployCallbackInfo struct {
	ID          int64  `db:"id"`
	Time        int64  `db:"time"`
	ClientID    string `db:"client_id"`
	InstallType string `db:"install_type"`
	InstallRes  string `db:"install_res"`
	MSG         []byte `db:"msg"`
	RequestUrl  string `db:"request_url"`
	IP          string `db:"ip"`
}

func (l *deployCallback) CreateDeployCallback(st *DeployCallbackInfo) (*DeployCallbackInfo, error) {
	if _, err := l.InsertWhere(dbhelper.UpdateFields{
		"time":         st.Time,
		"client_id":    st.ClientID,
		"install_type": st.InstallType,
		"install_res":  st.InstallRes,
		"msg":          st.MSG,
		"request_url":  st.RequestUrl,
		"ip":           st.IP,
	}); err != nil {
		return nil, err
	}
	return st, nil
}

func (l *sidecarList) CreateSidecarByDeploy(id uuid.UUID) (uuid.UUID, error) {
	if _, err := l.InsertWhere(dbhelper.UpdateFields{
		"id":       id,
		"status":   SIDECAR_INSTALLED,
		"disabled": false,
	}); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (l *sidecarList) NewSidecarRecord(name string, version string) (uuid.UUID, error) {
	id := uuid.NewV4()
	if _, err := l.InsertWhere(dbhelper.UpdateFields{
		"id":       id.String(),
		"name":     name,
		"status":   SIDECAR_NOT_INSTALLED,
		"disabled": false,
	}); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (l *sidecarList) GetSidecarInfo(id uuid.UUID) (*SidecarInfo, error) {
	whereCause := dbhelper.WhereCause{}
	info := SidecarInfo{}
	err := l.GetWhere(nil, whereCause.Equal("id", id), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

var _getSidecarListFields = utils.GetTagValues(SimpleSidecarInfo{}, "db")

func (l *sidecarList) GetSidecarList(pagination *apibase.Pagination) ([]SimpleSidecarInfo, int) {
	rows, totalRecords, err := l.SelectWhere(_getSidecarListFields, nil, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}

	list := []SimpleSidecarInfo{}
	for rows.Next() {
		info := SimpleSidecarInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return list, totalRecords
}

func (l *sidecarList) UpdateSidecar(id uuid.UUID, updateFields dbhelper.UpdateFields) error {
	return l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), updateFields, false)
}
