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
	"database/sql"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"github.com/jmoiron/sqlx"

	//apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
)

var (
	getClusterSql string = "select * from deploy_cluster_list where name = :name and type = 'kubernetes' and is_deleted = :is_deleted"
	insertClusterSql string = "insert into deploy_cluster_list (name,`desc`,tags,mode,configs,yaml,version,type) values (:name, :desc, :tags, :mode, :configs, :yaml, :version, 'kubernetes')"
	getClusterByIdSql string = "select * from deploy_cluster_list where id = :id and is_deleted = 0"
	getClusterSts *sqlx.NamedStmt
	insertClusterSts *sqlx.NamedStmt
	getClusterByIdSts *sqlx.NamedStmt
	DeployClusterList = &deployClusterList{
		PrepareFunc: prepareDeployClusterList,
	}
)
type deployClusterList struct {
	PrepareFunc
}

func prepareDeployClusterList() error{
	var err error
	getClusterSts,err = model.USE_MYSQL_DB().PrepareNamed(getClusterSql)
	if err != nil{
		log.Errorf("[kube cluster_list] init sql: %s , error %v",getClusterSql,err)
		return err
	}
	insertClusterSts,err = model.USE_MYSQL_DB().PrepareNamed(insertClusterSql)
	if err != nil{
		log.Errorf("[kube cluster_list] init sql: %s , error %v",insertClusterSql,err)
		return err
	}
	getClusterByIdSts,err = model.USE_MYSQL_DB().PrepareNamed(getClusterByIdSql)
	if err != nil{
		log.Errorf("[kube cluster_list] init sql: %s, error %v",getClusterByIdSql,err)
		return err
	}
	return nil
}


type ClusterInfo struct {
	Id         int               `db:"id" json:"id"`
	Name       string            `db:"name" json:"name"`
	Type       string            `db:"type" json:"type"`
	Mode       int               `db:"mode" json:"mode"`
	Version    string            `db:"version" json:"version"`
	Desc       string            `db:"desc" json:"desc"`
	Tags       string            `db:"tags" json:"tags"`
	Configs    sql.NullString    `db:"configs" json:"configs"`
	Yaml       sql.NullString    `db:"yaml" json:"yaml"`
	Status     int               `db:"status" json:"status"`
	ErrorMsg   string            `db:"errorMsg" json:"errorMsg"`
	CreateUser string            `db:"create_user" json:"create_user"`
	UpdateUser string            `db:"update_user" json:"update_user"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

func (l *deployClusterList) InsertK8sCluster(cluster *ClusterInfo) (int, error) {
	err := getClusterSts.Get(cluster,cluster)
	if err != nil && err == sql.ErrNoRows {
		result,err := insertClusterSts.Exec(cluster)
		if err != nil {
			log.Errorf("[kube cluster_list] exec sql: %s, value: %+v, error: %v:",insertClusterSts.QueryString,*cluster,err)
			return -1, err
		}
		id, err := result.LastInsertId()
		if err != nil{
			log.Errorf("[kube cluster_list] get lastInsertId error :%v",err)
		}
		return int(id), err
	} else if err == nil {
		return cluster.Id, fmt.Errorf("k8s集群:%v 已存在", cluster.Name)
	} else {
		log.Errorf("[kube cluster_list] query sql: %s, value: %+v, error: %v",getClusterSts.QueryString,*cluster,err)
		return -1, err
	}
}

func (l *deployClusterList) GetClusterById(id int) (*ClusterInfo,error){
	info := ClusterInfo{}
	info.Id = id
	err := getClusterByIdSts.Get(&info,info)
	if err == nil{
		return &info,nil
	}
	if err == sql.ErrNoRows{
		return nil, nil
	}
	return nil, fmt.Errorf("[kube cluster_list] query sql: %s, id: %d, error: %v",getClusterByIdSql,id,err)
}
