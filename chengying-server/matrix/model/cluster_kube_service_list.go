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
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"k8s.io/api/core/v1"
	"time"
)

type deployKubeServiceList struct {
	dbhelper.DbTable
}

var DeployKubeServiceList = &deployKubeServiceList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CLUSTER_KUBE_SERVICE_LIST},
}

type ClusterKubeService struct {
	Id             int               `db:"id" json:"id"`
	Pid            int               `db:"pid" json:"pid"`
	ClusterId      int               `db:"clusterId" json:"clusterId"`
	NameSpace      string            `db:"namespace" json:"namespace"`
	ProductName    string            `db:"product_name" json:"product_name"`
	ProductVersion string            `db:"product_version" json:"product_version"`
	ServiceName    string            `db:"service_name" json:"service_name"`
	ServiceVersion string            `db:"service_version" json:"service_version"`
	ClusterIp      string            `db:"cluster_ip" json:"cluster_ip"`
	Type           string            `db:"type" json:"type"`
	IsDeleted      int               `db:"is_deleted" json:"is_deleted"`
	UpdateTime     dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime     dbhelper.NullTime `db:"create_time" json:"create_time"`
}

func (l *deployKubeServiceList) UpdateOrCreate(service *v1.Service, labels *ClusterKubeLabels) (error, int) {
	info := ClusterKubeService{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("clusterId", labels.ClusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", service.GetNamespace())
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_name", labels.ProductName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", labels.ServiceName)

	err := l.GetWhere(nil, whereCause, &info)
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"pid":             labels.Pid,
			"clusterId":       labels.ClusterId,
			"namespace":       service.GetNamespace(),
			"product_name":    labels.ProductName,
			"product_version": labels.ProductVersion,
			"service_name":    labels.ServiceName,
			"service_version": labels.ServiceVersion,
			"cluster_ip":      service.Spec.ClusterIP,
			"type":            service.Spec.Type,
			"updated":         time.Now(),
			"created":         time.Now(),
		})
		if err != nil {
			return err, -1
		}
		seq, _ := ret.LastInsertId()
		return nil, int(seq)
	} else if err == nil {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"pid":             labels.Pid,
			"product_version": labels.ProductVersion,
			"service_version": labels.ServiceVersion,
			"cluster_ip":      service.Spec.ClusterIP,
			"type":            service.Spec.Type,
			"updated":         time.Now(),
		}, false)
		return err, info.Id
	} else {
		return err, -1
	}
	return nil, -1
}

func (l *deployKubeServiceList) Delete(service *v1.Service, labels *ClusterKubeLabels) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("clusterId", labels.ClusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", service.GetNamespace())
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_name", labels.ProductName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", labels.ServiceName)
	whereCause = whereCause.And()

	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"updated":   time.Now(),
		"isDeleted": 1,
	}, false)
	if err != nil {
		log.Errorf("[deployKubeServiceList] Delete err: %v", err)
		return err
	}
	return nil
}
