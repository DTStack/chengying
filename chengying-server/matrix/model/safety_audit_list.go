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
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"strconv"
	"time"
)

type safetyAuditList struct {
	dbhelper.DbTable
}

var SafetyAuditList = &safetyAuditList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_SAFETY_AUDIT_LIST},
}

type SafetyAuditInfo struct {
	Id               int       `db:"id" json:"id"`
	Operator         string    `db:"operator" json:"operator"`
	Module           string    `db:"module" json:"module"`
	Operation        string    `db:"operation" json:"operation"`
	IP               string    `db:"ip" json:"ip"`
	Content          string    `db:"content" json:"content"`
	IsDeleted        int       `db:"is_deleted" json:"-"`
	UpdateTime       time.Time `db:"update_time" json:"-"`
	CreateTime       time.Time `db:"create_time" json:"-"`
	CreateTimeString string    `json:"create_time"`
}

var _getSafetyAuditListFields = utils.GetTagValues(SafetyAuditInfo{}, "db")

func (s *safetyAuditList) GetSafetyAuditOperationList(module string) ([]string, error) {
	query := `SELECT DISTINCT audit_item_list.operation FROM audit_item_list WHERE audit_item_list.module=? AND audit_item_list.is_deleted=0`
	var list []string
	err := USE_MYSQL_DB().Select(&list, query, module)
	return list, err
}

func (s *safetyAuditList) GetSafetyAuditModuleList() ([]string, error) {
	query := `SELECT DISTINCT audit_item_list.module FROM audit_item_list WHERE audit_item_list.is_deleted=0`
	var list []string
	err := USE_MYSQL_DB().Select(&list, query)
	return list, err
}

func (s *safetyAuditList) SelectSafetyAuditListByWhere(pagination *apibase.Pagination, module, operation, operator, ip, content, from, to string) ([]SafetyAuditInfo, int, error) {
	whereCause := dbhelper.WhereCause{}

	if module != "" {
		whereCause = whereCause.Equal("module", module).And()
	}
	if operation != "" {
		whereCause = whereCause.Equal("operation", operation).And()
	}
	if operator != "" {
		whereCause = whereCause.Like("operator", "%"+operator+"%").And()
	}
	if ip != "" {
		whereCause = whereCause.Like("ip", "%"+ip+"%").And()
	}
	if content != "" {
		whereCause = whereCause.Like("content", "%"+content+"%").And()
	}
	if from != "" {
		from, _ := strconv.ParseInt(from, 10, 0)
		to, _ := strconv.ParseInt(to, 10, 0)
		whereCause = whereCause.Between("create_time", time.Unix(from, 0).Format(base.TsLayout), time.Unix(to, 0).Format(base.TsLayout)).And()
	}
	whereCause = whereCause.Equal("is_deleted", 0)

	rows, count, err := s.SelectWhere(_getSafetyAuditListFields, whereCause, pagination)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	safetyAuditInfoList := []SafetyAuditInfo{}
	for rows.Next() {
		info := SafetyAuditInfo{}
		if err = rows.StructScan(&info); err != nil {
			return nil, 0, err
		}
		info.CreateTimeString = info.UpdateTime.Format(base.TsLayout)
		safetyAuditInfoList = append(safetyAuditInfoList, info)
	}
	return safetyAuditInfoList, count, nil
}

func (s *safetyAuditList) InsertSafetyAuditRecord(operator, module, operation, ip, content string) error {
	_, err := s.InsertWhere(dbhelper.UpdateFields{
		"operator":  operator,
		"operation": operation,
		"module":    module,
		"ip":        ip,
		"content":   content,
	})
	if err != nil {
		return err
	}
	return nil
}
