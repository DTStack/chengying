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
	"time"
)

type deployKubeProductLock struct {
	dbhelper.DbTable
}

var DeployKubeProductLock = &deployKubeProductLock{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_KUBE_PRODUCT_LOCK},
}

type KubeProductLock struct {
	Id         int               `db:"id" json:"id"`
	Pid        int               `db:"pid" json:"pid"`
	ClusterId  int               `db:"clusterId" json:"clusterId"`
	Namespace  string            `db:"namespace" json:"namespace"`
	IsDeploy   int               `db:"is_deploy" json:"is_deploy"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

func (l *deployKubeProductLock) GetByPidAndClusterIdAndNamespace(pid, clusterId int, namespace string) (KubeProductLock, error) {
	info := KubeProductLock{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("pid", pid).And().
		Equal("clusterId", clusterId).And().
		Equal("namespace", namespace).And().
		Equal("is_deleted", 0), &info)
	return info, err
}

func (l *deployKubeProductLock) InsertOrUpdateRecord(info KubeProductLock) error {
	tmp := KubeProductLock{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("pid", info.Pid).And().
		Equal("is_deleted", 0).And().Equal("clusterId", info.ClusterId).And().Equal("namespace", info.Namespace), &tmp)
	if err != nil && err == sql.ErrNoRows {
		_, err = l.InsertWhere(dbhelper.UpdateFields{
			"pid":       info.Pid,
			"clusterId": info.ClusterId,
			"namespace": info.Namespace,
			"is_deploy": info.IsDeploy,
		})
	} else if err == nil {
		err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("pid", info.Pid).And().Equal("is_deleted", 0).
			And().Equal("clusterId", info.ClusterId).And().Equal("namespace", info.Namespace), dbhelper.UpdateFields{
			"is_deploy":   info.IsDeploy,
			"update_time": time.Now(),
		}, false)
	} else {
		return err
	}
	return err
}
