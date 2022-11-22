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

package kube

import (
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"github.com/jmoiron/sqlx"
)

type deployClusterK8sAvailable struct {
	PrepareFunc
}
var (
	getAvailableSql = "select * from deploy_cluster_k8s_available where mode = :mode and is_deleted = 0"
	getAvailableVersionSql = "select version from deploy_cluster_k8s_available where version like :version and is_deleted=0"
	getAvailableSts *sqlx.NamedStmt
	getAvailableVersionSts *sqlx.NamedStmt
	DeployClusterK8sAvailable = &deployClusterK8sAvailable{
		PrepareFunc: prepareDeployClusterK8sAvailable,
	}
)

func prepareDeployClusterK8sAvailable() error{
	var err error
	getAvailableSts,err = model.USE_MYSQL_DB().PrepareNamed(getAvailableSql)
	if err != nil{
		log.Errorf("[kube cluster_k8s_availbale]: init sql: %s , error %v",getAvailableSql,err)
		return err
	}
	getAvailableVersionSts,err = model.USE_MYSQL_DB().PrepareNamed(getAvailableVersionSql)
	if err != nil{
		log.Errorf("[kube cluster_k8s_availbale]: init sql: %s, error %v",getAvailableVersionSql,err)
	}
	return nil
}

type ClusterK8sAvailable struct {
	Id         int               `db:"id"`
	Mode       int               `db:"mode"`
	Version    string            `db:"version"`
	Properties string            `db:"properties"`
	UpdateTime dbhelper.NullTime `db:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time"`
	IsDeleted  int               `db:"is_deleted"`
}

func (l *deployClusterK8sAvailable) GetClusterK8sAvailableByMode(mode int) ([]ClusterK8sAvailable, error) {
	info := []ClusterK8sAvailable{}
	arg := &ClusterK8sAvailable{
		Mode: mode,
	}
	rows,err := getAvailableSts.Queryx(arg)
	if err != nil{
		log.Errorf("[kube cluster_k8s_availbale] getk8sAvailable by mode sql: %s, value: %d, error: %v",getAvailableSql,mode,err)
		return nil,err
	}
	for rows.Next(){
		tbsc := ClusterK8sAvailable{}
		if err = rows.StructScan(&tbsc); err != nil{
			log.Errorf("[kube cluster_k8s_availbale]: struct scan to ClusterK8sAvailable error :%v",err)
			return nil, err
		}
		info = append(info, tbsc)
	}
	return info, nil
}

func (l *deployClusterK8sAvailable) GetRealVersion(version string) (string, error) {
	if len(version) == 0{
		return "",nil
	}
	version = version + "%"
	available := &ClusterK8sAvailable{
		Version: version,
	}
	if err := getAvailableVersionSts.Get(&version,available);err != nil{
		log.Errorf("[kube cluster_k8s_available] getRealVersion sql: %s, value: %s, error: %v",getAvailableVersionSql,version,err)
		return "",err
	}

	return version, nil
}
