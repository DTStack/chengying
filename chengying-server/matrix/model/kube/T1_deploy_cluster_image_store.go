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
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"github.com/jmoiron/sqlx"
)
var (
	getImageStoreByIdSql = "select * from deploy_cluster_image_store where id = :id and is_deleted = 0"
	getImageStoreByIdSts *sqlx.NamedStmt
	getImageStoreByClusterIdSql = "select * from deploy_cluster_image_store where clusterid = :clusterId and is_deleted = 0"
	getImageStoreByClusterIdSts *sqlx.NamedStmt
	DeployClusterImageStore = &deployClusterImageStore{
		PrepareFunc: prepareDeployClusterImageStore,
	}
)
type deployClusterImageStore struct {
	PrepareFunc
}

func prepareDeployClusterImageStore() error{
	var err error
	getImageStoreByIdSts,err = model.USE_MYSQL_DB().PrepareNamed(getImageStoreByIdSql)
	if err != nil{
		log.Errorf("[kube deploy_cluster_image_store]: init sql: %s , error %v",getImageStoreByIdSql,err)
		return err
	}
	getImageStoreByClusterIdSts,err = model.USE_MYSQL_DB().PrepareNamed(getImageStoreByClusterIdSql)
	if err != nil {
		log.Errorf("[kube deploy_cluster_image_store]: init sql: %s , error %v",getImageStoreByClusterIdSql,err)
		return err
	}
	return nil
}

type DeployClusterImageStoreSchema struct {
	Id         int               `db:"id"`
	ClusterId  int               `db:"clusterId"`
	IsDefault  int               `db:"is_default"`
	Name       string            `db:"name"`
	Alias      string            `db:"alias"`
	Address    string            `db:"address"`
	Username   string            `db:"username"`
	Password   string            `db:"password"`
	Email      string            `db:"email"`
	UpdateTime sql.NullTime      `db:"update_time"`
	CreateTime sql.NullTime      `db:"create_time"`
	IsDeleted  int               `db:"is_deleted"`
}

func (s *deployClusterImageStore)GetById(id int) (*DeployClusterImageStoreSchema,error){
	sc := &DeployClusterImageStoreSchema{
		Id:         id,
	}
	if err := getImageStoreByIdSts.Get(sc,sc);err!=nil{
		if err == sql.ErrNoRows{
			return nil,nil
		}
		log.Errorf("[kube deploy_cluster_image_store]: get imagestore %s by id %d error: %v",getImageStoreByIdSql,id,err)
		return nil,err
	}
	return sc,nil
}

func (s *deployClusterImageStore)GetByClusterId(cid int) (*[]DeployClusterImageStoreSchema,error){

	arg := &DeployClusterImageStoreSchema{
		ClusterId: cid,
	}

	rows, err := getImageStoreByClusterIdSts.Queryx(arg)
	if err != nil {
		log.Errorf("[kube deploy_cluster_image_store]: init sql: %s, value %+v , error %v",getImageStoreByClusterIdSql,*arg,err)
		return nil,err
	}
	result := []DeployClusterImageStoreSchema{}
	for rows.Next(){
		imageStore:= DeployClusterImageStoreSchema{}
		if err = rows.StructScan(&imageStore);err != nil{
			log.Errorf("[kube deploy_cluster_image_store]: init sql: %s, error %v",getImageStoreByClusterIdSql,err)
			return nil,err
		}
		result = append(result,imageStore)
	}

	return &result,nil
}