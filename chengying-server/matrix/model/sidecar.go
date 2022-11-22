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
	"time"
)

type Sidecar struct {
	Id             string         `db:"id"`
	Status         int            `db:"status"`
	Disabled       int            `db:"disabled"`
	Name           string         `db:"name"`
	Version        string         `db:"version"`
	Host           string         `db:"host"`
	OSType         string         `db:"os_type"`
	IsEcs          int            `db:"is_ecs"`
	OsPlatform     string         `db:"os_platform"`
	OsVersion      string         `db:"os_version"`
	CpuSerial      string         `db:"cpu_serial"`
	CpuCores       int            `db:"cpu_cores"`
	MemSize        int            `db:"mem_size"`
	SwapSize       int            `db:"swap_size"`
	DeployDate     *time.Time     `db:"deploy_date"`
	AutoDeployment int            `db:"auto_deployment"`
	LastUpdateDate *time.Time     `db:"last_update_date"`
	AutoUpdated    int            `db:"auto_updated"`
	ServerHost     string         `db:"server_host"`
	ServerPort     int            `db:"server_port"`
	SshHost        string         `db:"ssh_host"`
	SshUser        string         `db:"ssh_user"`
	SshPassword    string         `db:"ssh_password"`
	SshPort        string         `db:"ssh_port"`
	CpuUsage       float64        `db:"cpu_usage"`
	MemUsage       int64          `db:"mem_usage"`
	SwapUsage      int64          `db:"swap_usage"`
	Load1          float64        `db:"load1"`
	UpTime         float64        `db:"uptime"`
	DiskUsage      sql.NullString `db:"disk_usage"`
	NetUsage       sql.NullString `db:"net_usage"`
	LocalIp        string         `db:"local_ip"`
}

type DiskUsage struct {
	MountPoint string `json:"mountPoint"`
	UsedSpace  uint64 `json:"usedSpace"`
	TotalSpace uint64 `json:"totalSpace"`
}

type NetUsage struct {
	IfName    string   `json:"ifName"`
	IfIp      []string `json:"ifIp"`
	BytesSent uint64   `json:"bytesSent"`
	BytesRecv uint64   `json:"bytesRecv"`
}

type NetUsageDisplay struct {
	IfName    string   `json:"ifName"`
	IfIp      []string `json:"ifIp"`
	BytesSent string   `json:"bytesSent"`
	BytesRecv string   `json:"bytesRecv"`
}

func GetSidecars() ([]Sidecar, error) {
	query := "SELECT * FROM sidecar_list"
	sidecars := make([]Sidecar, 0)

	err := USE_MYSQL_DB().Select(&sidecars, query)
	return sidecars, err
}
