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
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"fmt"
)

/*
 @Author: zhijian
 @Date: 2021/6/2 15:19
 @Description:
*/

type execShellList struct {
	dbhelper.DbTable
}
type ExecShellInfo struct {
	Id          int               `db:"id" json:"-"`
	ClusterId   int               `db:"cluster_id" json:"clusterId"`
	ExecId      string            `db:"exec_id" json:"execId"`
	OperationId string            `db:"operation_id" json:"operationId"`
	ProductName string            `db:"product_name" json:"name"`
	ServiceName string            `db:"service_name" json:"serviceName"`
	Seq         sql.NullInt32     `db:"seq" json:"seq"`
	ExecStatus  int               `db:"exec_status" json:"status"`
	Sid         string            `db:"sid" json:"sid"`
	ShellType   int               `db:"shell_type" json:"shellType"`
	CreateTime  dbhelper.NullTime `db:"create_time" json:"-"`
	EndTime     dbhelper.NullTime `db:"end_time" json:"-"`
	Duration    sql.NullFloat64   `db:"duration" json:"-"`
	UpdateTime  dbhelper.NullTime `db:"update_time" json:"-"`
	HostIp      string            `db:"host_ip" json:"hostIp"`

	ShowCreateTime string  `json:"startTime"`
	ShowEndTime    string  `json:"endTime"`
	ShowDuration   float64 `json:"duration"`
	ShellDesc      string  `json:"desc"`
}

var ExecShellList = &execShellList{
	dbhelper.DbTable{USE_MYSQL_DB, EXECSHELL_LIST},
}

func (ex *execShellList) InsertExecShellInfo(clusterId int, operationId, execId, productName, svcName, sid string, shellType int) error {
	if operationId == "" {
		return nil
	}
	err, info := DeployHostList.GetHostInfoBySid(sid)
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("INSERT INTO %s (cluster_id,exec_id,operation_id, product_name,service_name,exec_status,sid,host_ip,shell_type) VALUES (?,?,?,?,?,?,?,?,?)", EXECSHELL_LIST)
	_, err = ex.GetDB().Exec(sql, clusterId, execId, operationId, productName, svcName, enums.ExecStatusType.Running.Code, sid, info.Ip, shellType)
	if err != nil {
		return err
	}
	return nil
}

func (ex *execShellList) IsExist(seq int) (bool, error) {
	//存在则是 1  不存在 返回 0
	sql := fmt.Sprintf("select IFNULL((select 1 from %s where seq = ? limit 1),0)", EXECSHELL_LIST)
	var ret int
	err := ex.GetDB().Get(&ret, sql, seq)
	if ret == 1 {
		return true, nil
	} else {
		return false, err
	}
}

func (ex *execShellList) UpdateSeqByExecId(execId string, seq int) error {
	sql := fmt.Sprintf("update %s set seq =? where exec_id = ?", EXECSHELL_LIST)
	_, err := ex.GetDB().Exec(sql, seq, execId)
	if err != nil {
		return err
	}
	return nil
}
func (ex *execShellList) UpdateStatusBySeq(seq int, status int, endTime dbhelper.NullTime, duration sql.NullFloat64) error {
	sql := fmt.Sprintf("update %s set exec_status =? ,end_time=? , duration=? where seq =?", EXECSHELL_LIST)
	_, err := ex.GetDB().Exec(sql, status, endTime, duration, seq)
	if err != nil {
		return err
	}
	return nil
}

func (ex *execShellList) UpdateStatusByExecId(execId string, status int, endTime dbhelper.NullTime, duration sql.NullFloat64) error {
	sql := fmt.Sprintf("update %s set exec_status =? ,end_time=? , duration=? where exec_id =?", EXECSHELL_LIST)
	_, err := ex.GetDB().Exec(sql, status, endTime, duration, execId)
	if err != nil {
		return err
	}
	return nil
}

func (ex *execShellList) GetBySeq(seq int) (*ExecShellInfo, error) {
	sql := fmt.Sprintf("select * from %s where seq =?", EXECSHELL_LIST)
	var execShellInfo ExecShellInfo
	err := ex.GetDB().Get(&execShellInfo, sql, seq)
	if err != nil {
		return nil, err
	}
	return &execShellInfo, nil
}

func (ex *execShellList) SelectShellGroupBySeq(seq int) ([]ExecShellInfo, error) {
	sql := fmt.Sprintf("select * from %s where operation_id = (select operation_id from %s where seq = ?)", EXECSHELL_LIST, EXECSHELL_LIST)
	var shellGroup []ExecShellInfo
	err := ex.GetDB().Select(&shellGroup, sql, seq)
	if err != nil {
		return nil, err
	}
	return shellGroup, nil
}

func (ex *execShellList) SelectShellGroupByOperationId(operationId string) ([]ExecShellInfo, error) {
	sql := fmt.Sprintf("select * from %s where operation_id = ?", EXECSHELL_LIST)
	var shellGroup []ExecShellInfo
	err := ex.GetDB().Select(&shellGroup, sql, operationId)
	if err != nil {
		return nil, err
	}
	return shellGroup, nil
}

func (ex *execShellList) GetByOperationId(operationId string) (*ExecShellInfo, error) {
	sql := fmt.Sprintf("select * from %s where operation_id = ? ", EXECSHELL_LIST)
	var one ExecShellInfo
	err := ex.GetDB().Get(&one, sql, operationId)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

func (ex *execShellList) GetByExecId(execId string) (*ExecShellInfo, error) {
	sql := fmt.Sprintf("select * from %s where exec_id = ?", EXECSHELL_LIST)
	var execShellInfo ExecShellInfo
	err := ex.GetDB().Get(&execShellInfo, sql, execId)
	if err != nil {
		return nil, err
	}
	return &execShellInfo, nil
}

func (ex *execShellList) GetCountByOperationId(operationId string) (*int, error) {
	sql := fmt.Sprintf("select count(1) from %s where operation_id = ?", EXECSHELL_LIST)
	var count int
	err := ex.GetDB().Get(&count, sql, operationId)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (ex *execShellList) DeleteBySid(sid string) error {
	sql := fmt.Sprintf("delete from %s where sid = ? ", EXECSHELL_LIST)
	_, err := ex.GetDB().Exec(sql, sid)
	if err != nil {
		return err
	}
	return nil
}
