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
	"easyagent/go-common/utils"
	"sync"
	"time"

	apibase "easyagent/go-common/api-base"
	dbhelper "easyagent/go-common/db-helper"
	"easyagent/internal/server/log"

	uuid "github.com/satori/go.uuid"
)

var (
	seqnoMu      = sync.Mutex{}
	lastSeqno    = int64(0)
	seqno        = int64(0)
	lastInsertId = int64(0)
)

const (
	PEROID_SEQNO = 4096
	MAX_SEQNO    = 65535
)

type agentOperation struct {
	dbhelper.DbTable
}

type operationProgress struct {
	dbhelper.DbTable
}

var AgentOperation = &agentOperation{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_OP_HISTORY},
}

var OperationProgress = &operationProgress{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_PROGRESS_HISTORY},
}

type Operation struct {
	SEQ         int               `db:"seq"`
	OpName      string            `db:"op_name"`
	OpTime      dbhelper.NullTime `db:"op_time"`
	Target      string            `db:"target"`
	SendStatus  int               `db:"send_status"`
	OpResult    int               `db:"auto_deployment"`
	OpResultMsg sql.NullString    `db:"op_result_msg"`
	FinishTime  dbhelper.NullTime `db:"finish_time"`
}

type Progress struct {
	SEQ       int               `db:"seq"`
	Ts        dbhelper.NullTime `db:"ts"`
	Progress  float32           `db:"progress"`
	SidecarId string            `db:"sidecar_id"`
	AgentId   string            `db:"agent_id"`
	Msg       string            `db:"msg"`
}

var _getAgentOperateFields = utils.GetTagValues(Operation{}, "db")
var _getProgressFields = utils.GetTagValues(Progress{}, "db")

func (l *agentOperation) NewOperationSeqno(target uuid.UUID, agentId uuid.UUID) (int64, error) {
	seqnoMu.Lock()
	defer seqnoMu.Unlock()
	log.Infof("NewOperationSeqno: %v, target: %v, agentId: %v", lastSeqno, target, agentId)

	//ugly fixup
	//这里单纯生成全局唯一序列号，不入库，减少agentOperation 数据库增长和binlog 日志增长
	//self auto 1-65535
	seqno = seqno + 1
	if seqno > MAX_SEQNO {
		seqno = 0
	}
	//跳过数据库里的seq记录
	//step over db record seq
	peroid := lastInsertId + PEROID_SEQNO

	seq := seqno + peroid
	lastSeqno = seq

	return seq, nil
}

func (l *agentOperation) NewOperationRecord(opName string, target uuid.UUID, agentId uuid.UUID) (int64, error) {
	seqnoMu.Lock()
	defer seqnoMu.Unlock()
	ret, err := l.InsertWhere(dbhelper.UpdateFields{
		"op_name": opName,
		"op_time": time.Now(),
		"target":  target.String(),
	})
	if err != nil {
		apibase.ThrowDBModelError(err)
		return -1, err
	}
	seq, _ := ret.LastInsertId()
	OperationProgress.NewProgressRecord(seq, target, agentId)
	lastInsertId = seq

	return seq, nil
}

func (l *agentOperation) GetOperation(seq int) *Operation {
	whereCause := dbhelper.WhereCause{}
	opera := Operation{}
	err := l.GetWhere(nil, whereCause.Equal("seq", seq), &opera)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	return &opera
}

func (l *agentOperation) Exists(seq int) bool {
	whereCause := dbhelper.WhereCause{}
	opera := Operation{}
	err := l.GetWhere(nil, whereCause.Equal("seq", seq), &opera)
	if err != nil {
		return false
	}
	return true
}

func (l *agentOperation) GetAgentOperations(pagination *apibase.Pagination) ([]Operation, int) {
	rows, totalRecords, err := l.SelectWhere(_getAgentOperateFields, nil, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	list := []Operation{}
	for rows.Next() {
		info := Operation{}
		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return list, totalRecords
}

func (l *agentOperation) UpdateAgentOperation(seq uint32, updateFields dbhelper.UpdateFields) error {
	if !l.Exists(int(seq)) {
		return nil
	}
	return l.UpdateWhere(dbhelper.MakeWhereCause().Equal("seq", seq), updateFields, false)
}

////
//
//Progress
//
func (l *operationProgress) NewProgressRecord(seq int64, sidecarId uuid.UUID, agentId uuid.UUID) {
	if _, err := l.InsertWhere(dbhelper.UpdateFields{
		"seq":        seq,
		"ts":         time.Now(),
		"progress":   0,
		"sidecar_id": sidecarId.String(),
		"agent_id":   agentId.String(),
	}); err != nil {
		apibase.ThrowDBModelError(err)
	}
}

func (l *operationProgress) GetOperation(seq int) *Progress {
	whereCause := dbhelper.WhereCause{}
	opera := Progress{}
	err := l.GetWhere(nil, whereCause.Equal("seq", seq), &opera)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	return &opera
}

func (l *operationProgress) UpdateProgress(seq int, progress float32, msg string) {
	l.UpdateWhere(dbhelper.MakeWhereCause().Equal("seq", seq), dbhelper.UpdateFields{
		"ts":       time.Now(),
		"progress": progress,
		"msg":      msg,
	}, true)
}
