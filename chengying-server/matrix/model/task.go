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
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status int

const (
	NOT_RUNNING Status = 0 //未运行
	RUNNING     Status = 1 //运行中
	FINISH      Status = 2 //正常
	FAILURE     Status = 3 //异常

)

const (
	TASK_FILES_DIR      = "task-files"
	TASK_STATUS_DISABLE = 0
	TASK_STATUS_ENABLE  = 1
	TASK_AUTO_RUN       = 0
	TASK_MANUAL_RUN     = 1
)

type taskList struct {
	dbhelper.DbTable
}

var TaskList = &taskList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_TASK_LIST},
}

type taskHostList struct {
	dbhelper.DbTable
}

var TaskHostList = &taskHostList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_TASK_HOST},
}

type taskLogList struct {
	dbhelper.DbTable
}

var TaskLogList = &taskLogList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_TASK_LOG},
}

type TaskInfo struct {
	ID           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Describe     string `db:"describe" json:"describe"`
	Spec         string `db:"spec" json:"spec"`
	ExecTimeout  int    `db:"exec_timeout" json:"exec_timeout"`
	LogRetention int    `db:"log_retention" json:"log_retention"`
	Status       int    `db:"status" json:"status"`
	Hosts        []ResHostInfo
	ExecType     int
	CreateTime   dbhelper.NullTime `db:"create_time" json:"create_time"`
	UpdateTime   dbhelper.NullTime `db:"update_time" json:"update_time"`
	IsDeleted    int               `db:"is_deleted" json:"is_deleted"`
}

type TaskHost struct {
	ID     int `db:"id" json:"id"`
	TaskId int `db:"task_id" json:"task_id"`
	HostId int `db:"host_id" json:"host_id"`
}

type ResHostInfo struct {
	Ip        string `db:"ip" json:"ip"`
	Sid       string `db:"sid" json:"sid"`
	HostName  string `db:"hostname" json:"hostname"`
	ClusterId int    `db:"cluster_id" json:"cluster_id"`
}

type TaskLog struct {
	ID          int               `db:"id" json:"id"`
	TaskId      int               `db:"task_id" json:"task_id"`
	Name        string            `db:"name" json:"name"`
	Spec        string            `db:"spec" json:"spec"`
	Ip          string            `db:"ip" json:"ip"`
	OperationId string            `db:"operation_id" json:"operationId"`
	Command     string            `db:"command" json:"command"`
	ExecType    int               `db:"exec_type" json:"exec_type"`
	ExecStatus  Status            `db:"exec_status" json:"exec_status"`
	ExecResult  string            `db:"exec_result" json:"exec_result"`
	StartTime   dbhelper.NullTime `db:"start_time" json:"start_time"`
	EndTime     dbhelper.NullTime `db:"end_time" json:"end_time"`
}

func (l *taskList) InsertTaskIfNotExist(name, describe, spec string, execTimeout, logRetention int) (error, int) {
	info := TaskInfo{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("name", name).And().Equal("is_deleted", 0), &info)
	if err != nil && err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"name":          name,
			"describe":      describe,
			"spec":          spec,
			"exec_timeout":  execTimeout,
			"log_retention": logRetention,
		})
		if err != nil {
			apibase.ThrowDBModelError(err)
			return err, -1
		}
		seq, _ := ret.LastInsertId()
		return nil, int(seq)
	} else if err == nil {
		aid := info.ID
		return fmt.Errorf("%v 已存在", name), aid
	} else {
		return err, 0
	}
}

func (l *taskHostList) InsertTaskHost(taskId int, hosts []TaskHost) error {
	if err := l.DeleteTaskHostByTaskId(taskId); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("delete task host err: %v", err)
		return err
	}
	valueStrings := make([]string, 0, len(hosts))
	valueArgs := make([]interface{}, 0, len(hosts)*2)
	for _, host := range hosts {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, host.TaskId)
		valueArgs = append(valueArgs, host.HostId)
	}
	stmt := fmt.Sprintf("INSERT INTO %s (task_id, host_id) VALUES %s", TBL_TASK_HOST, strings.Join(valueStrings, ","))
	_, err := l.Exec(stmt, valueArgs...)
	return err
}

func (l *taskLogList) InsertTaskLog(log TaskLog) (int64, error) {
	ret, err := l.InsertWhere(dbhelper.UpdateFields{
		"task_id":      log.TaskId,
		"name":         log.Name,
		"spec":         log.Spec,
		"ip":           log.Ip,
		"operation_id": log.OperationId,
		"command":      log.Command,
		"exec_type":    log.ExecType,
		"exec_status":  log.ExecStatus,
		"exec_result":  log.ExecResult,
		"start_time":   log.StartTime.Time,
	})
	if err != nil {
		apibase.ThrowDBModelError(err)
		return -1, err
	}
	id, _ := ret.LastInsertId()
	return id, err
}

func (l *taskList) GetTaskInfoList() ([]TaskInfo, error) {
	list := make([]TaskInfo, 0)
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	rows, _, err := l.SelectWhere(_getTaskInfoFields, whereCause, nil)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		info := TaskInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return list, err
		}
		list = append(list, info)
	}
	return list, nil
}

var _getTaskInfoFields = utils.GetTagValues(TaskInfo{}, "db")

func (l *taskList) GetTaskInfoListByName(name string) ([]TaskInfo, int) {
	list := make([]TaskInfo, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE is_deleted = 0", TBL_TASK_LIST)
	if name != "" {
		query += " AND ( `name` LIKE '%" + name + "%' OR `describe` LIKE '%" + name + "%' )"
	}
	if err := USE_MYSQL_DB().Select(&list, query); err != nil {
		apibase.ThrowDBModelError(err)
	}
	return list, len(list)
}

func (l *taskHostList) GetTaskHostInfoByTaskId(taskId int) (error, []ResHostInfo) {
	list := make([]ResHostInfo, 0)
	query := "SELECT ip,hostname,host_rel.sid,host_rel.clusterId AS cluster_id FROM deploy_host " +
		"LEFT JOIN task_host ON deploy_host.id = task_host.host_id " +
		"LEFT JOIN deploy_cluster_host_rel host_rel ON deploy_host.sid = host_rel.sid " +
		"WHERE isDeleted=0 AND task_host.task_id=? AND host_rel.is_deleted=0"
	err := l.GetDB().Select(&list, query, taskId)
	if err != nil {
		return err, list
	}
	return nil, list
}

func (l *taskHostList) GetTaskHostByHostIds(hostIds []string) (error, []ResHostInfo) {
	list := make([]ResHostInfo, 0)
	query := "SELECT ip,hostname,host_rel.sid,host_rel.clusterId AS cluster_id FROM deploy_host " +
		"LEFT JOIN deploy_cluster_host_rel host_rel ON deploy_host.sid = host_rel.sid " +
		"WHERE isDeleted=0 AND host_rel.is_deleted=0 "
	if len(hostIds) > 0 {
		query += " AND deploy_host.id IN ("
		for i, v := range hostIds {
			if i > 0 {
				query += "," + v
			} else {
				query += v
			}
		}
		query += ")"
	}
	err := l.GetDB().Select(&list, query)
	if err != nil {
		return err, list
	}
	return nil, list
}

var _getTaskLogFields = utils.GetTagValues(TaskLog{}, "db")

func (l *taskLogList) GetTaskLogByTaskId(taskId int) ([]TaskLog, int) {
	list := make([]TaskLog, 0)
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("task_id", taskId)
	rows, totalRecords, err := l.SelectWhere(_getTaskLogFields, whereCause, nil)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()
	for rows.Next() {
		info := TaskLog{}
		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return list, totalRecords
}

func (l *taskList) GetTaskInfoByTaskId(taskId int) (TaskInfo, error) {
	info := TaskInfo{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", taskId)
	whereCause = whereCause.And().Equal("is_deleted", 0)
	return info, l.GetWhere(nil, whereCause, &info)
}

func (l *taskList) UpdateStatusByTaskIds(taskIds []string, status int) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	if len(taskIds) > 0 {
		ids := make([]interface{}, 0, len(taskIds))
		for _, id := range taskIds {
			ids = append(ids, interface{}(id))
		}
		whereCause = whereCause.And().Included("id", ids...)
	}
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"status":      status,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *taskList) DeleteTaskByTaskId(taskId int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", taskId), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *taskHostList) DeleteTaskHostByTaskId(taskId int) error {
	query := fmt.Sprintf("delete from %s where task_id = ?", TBL_TASK_HOST)
	_, err := l.GetDB().Exec(query, taskId)
	if err != nil {
		return err
	}
	return nil
}

func (l *taskList) UpdateSpecByTaskId(taskId int, spec string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	whereCause = whereCause.And().Equal("id", taskId)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"spec":        spec,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *taskLogList) UpdateTaskLogById(logId int64, status Status, result string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", logId)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"exec_status": status,
		"exec_result": result,
		"end_time":    time.Now(),
	}, false)
	return err
}

func (l *taskLogList) GetOperationIdByTaskId(taskId int, execStatus string, pagination *apibase.Pagination) ([]TaskLog, int) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("task_id", taskId)
	if execStatus != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("exec_status", execStatus)
	}
	where, value := whereCause.SQL()
	query := "SELECT operation_id FROM " + TBL_TASK_LOG + " " + where + " " + "GROUP BY operation_id " + pagination.AsQuery()
	queryCount := "SELECT COUNT(*) FROM " + "(SELECT operation_id FROM task_log " + where + " GROUP BY operation_id) AS t"
	var count int
	list := make([]TaskLog, 0)
	if err := l.GetDB().Get(&count, queryCount, value...); err != nil {
		apibase.ThrowDBModelError(err)
	}
	if count > 0 {
		rows, err := USE_MYSQL_DB().Queryx(query, value...)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		defer rows.Close()
		for rows.Next() {
			var info TaskLog
			if err := rows.StructScan(&info); err != nil {
				apibase.ThrowDBModelError(err)
			}
			list = append(list, info)
		}
	}
	return list, count
}
