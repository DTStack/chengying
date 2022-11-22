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

	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
)

type uncheckedService struct {
	dbhelper.DbTable
}

var DeployUncheckedService = &uncheckedService{
	dbhelper.DbTable{GetDB: USE_MYSQL_DB, TableName: TBL_DEPLOY_UNCHECKED_SERVICE},
}

type DeployUncheckedServiceInfo struct {
	ID                int               `db:"id"`
	ClusterId         int               `db:"cluster_id"`
	Pid               int               `db:"pid"`
	UncheckedServices string            `db:"unchecked_services"`
	UpdateDate        dbhelper.NullTime `db:"update_time"`
	CreateDate        dbhelper.NullTime `db:"create_time"`
	Namespace         string            `db:"namespace"`
}

func (us *uncheckedService) GetUncheckedServicesByPidClusterId(pid, clusterId int, namespace string) (info *DeployUncheckedServiceInfo, err error) {
	info = &DeployUncheckedServiceInfo{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("pid", pid)
	if clusterId > 0 {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("cluster_id", clusterId)
	}
	if namespace != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("namespace", namespace)
	}
	if err = us.GetWhere(nil, whereCause, info); err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (us *uncheckedService) SetUncheckedService(pid, clusterId int, uncheckedServices string) error {
	return us.UpdateWhere(dbhelper.MakeWhereCause(), dbhelper.UpdateFields{
		"pid":                pid,
		"cluster_id":         clusterId,
		"unchecked_services": uncheckedServices,
		"update_time":        time.Now(),
	}, true)
}
