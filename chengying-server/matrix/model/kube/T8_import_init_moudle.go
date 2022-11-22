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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"github.com/jmoiron/sqlx"
)

var (
	getMoudleSql = "select * from import_init_moudle where is_deleted = 0"
	getMoudleSts *sqlx.Stmt
	ImportInitMoudle = &importInitMoudle{
		PrepareFunc: prepareImportInitMoudle,
	}
)
type importInitMoudle struct {
	PrepareFunc
}

func prepareImportInitMoudle() error{
	var err error
	getMoudleSts,err = model.USE_MYSQL_DB().Preparex(getMoudleSql)
	if err != nil{
		log.Errorf("[kube import_init_moudle]: init sql: %s , error %v",getMoudleSql,err)
		return err
	}
	return nil
}
type ImportInitMoudleSchema struct {
	Id				int 	`db:"id"`
	ServiceAccount 	string 	`db:"service_account"`
	Role 			string 	`db:"role"`
	RoleBinding 	string 	`db:"role_binding"`
	Operator 		string 	`db:"operator"`
	LogConfig       string  `db:"log_config"`
	IsDeleted 		int 	`db:"is_deleted"`
}

func (i *importInitMoudle) GetInitMoudle() (*ImportInitMoudleSchema,error){
	sc := &ImportInitMoudleSchema{}
	if err := getMoudleSts.Get(sc);err!=nil{
		log.Errorf("[kube import_init_moudle]: get init moudle error %v",err)
		return nil ,err
	}
	return sc,nil
}
