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
	insertNamespaceClientSql = "insert into deploy_namespace_client (yaml, namespace_id, file_name) values (:yaml, :namespace_id, :file_name)"
	updateNamespaceClientSql = "update deploy_namespace_client set yaml = :yaml, file_name = :file_name where namespace_id = :namespace_id"
	selectNamespaceClientSql = "select * from deploy_namespace_client where namespace_id = :namespace_id"
	insertNamespaceClientSts *sqlx.NamedStmt
	updateNamespaceClientSts *sqlx.NamedStmt
	selectNamespaceClientSts *sqlx.NamedStmt
	DeployNamespaceClient = &deployNamespaceClient{
		PrepareFunc: prepareDeployNamespaceClinet,
	}
)
type deployNamespaceClient struct {
	PrepareFunc
}

type DeployNamespaceClientSchema struct {
	Id  int `db:"id"`
	Yaml string `db:"yaml"`
	NamespaceId int `db:"namespace_id"`
	Filename    string 	`db:"file_name"`
}

func prepareDeployNamespaceClinet() error{
	var err error
	insertNamespaceClientSts,err = model.USE_MYSQL_DB().PrepareNamed(insertNamespaceClientSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_client]: init sql: %s , error %v",insertNamespaceClientSql,err)
		return err
	}
	updateNamespaceClientSts,err = model.USE_MYSQL_DB().PrepareNamed(updateNamespaceClientSql)
	if err !=nil{
		log.Errorf("[kube deploy_namespace_client]: init sql: %s , error %v",updateNamespaceClientSql,err)
		return err
	}
	selectNamespaceClientSts,err = model.USE_MYSQL_DB().PrepareNamed(selectNamespaceClientSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_client]: init sql: %s , error %v",selectNamespaceClientSql,err)
		return err
	}
	return nil
}

func (c *deployNamespaceClient)Insert(tbsc *DeployNamespaceClientSchema) error{
	_,err := insertNamespaceClientSts.Exec(tbsc)
	if err != nil{
		log.Errorf("[deploy_namespace_client]: insert sql exec error %v",err)
		return err
	}
	return nil
}

func (c *deployNamespaceClient)Update(tbsc *DeployNamespaceClientSchema) error{
	_,err := updateNamespaceClientSts.Exec(tbsc)
	if err != nil{
		log.Errorf("[deploy_namespace_client]: update sql exec error %v",err)
		return err
	}
	return nil
}

func (c *deployNamespaceClient)Get(namespaceId int) (*DeployNamespaceClientSchema,error){
	tbsc := &DeployNamespaceClientSchema{
		NamespaceId: namespaceId,
	}
	err := selectNamespaceClientSts.Get(tbsc,tbsc)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil,nil
		}
		return nil,err
	}
	return tbsc,nil
}
