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
	getPartSql = "select * from workload_part where workload_id = :workload_id order by id asc"
	getPartSts *sqlx.NamedStmt
	WorkloadPart = &worloadPart{
		PrepareFunc: prepareWorkloadPart,
	}
)

func prepareWorkloadPart() error{
	var err error
	getPartSts,err = model.USE_MYSQL_DB().PrepareNamed(getPartSql)
	if err != nil{
		log.Errorf("[kube workload_part]: init sql: %s , error %v",getPartSql,err)
		return err
	}
	return nil
}
type worloadPart struct {
	PrepareFunc
}
type WorkloaPartSchema struct {
	Id		int 	`db:"id"`
	Name 	string 	`db:"name"`
	Type 	string 	`db:"type"`
	Parameters	string 	`db:"parameters"`
	WorkloadID	int 	`db:"workload_id"`
}

func (w *worloadPart) Select(workloadId int) ([]WorkloaPartSchema,error){
	arg := &WorkloaPartSchema{
		WorkloadID: workloadId,
	}
	rows,err := getPartSts.Queryx(arg)
	if err != nil{
		log.Errorf("[kube workload_part] select sql %s, workloadid %d, error %v",getPartSql,workloadId,err)
		return nil,err
	}
	result := []WorkloaPartSchema{}
	for rows.Next(){
		tbsc := WorkloaPartSchema{}
		if err = rows.StructScan(&tbsc);err != nil{
			log.Errorf("[kube workload_part] sql %s result failed struct scan to workload_part error %v",getPartSql,err)
			return nil,err
		}
		result = append(result,tbsc)
	}
	if len(result) == 0{
		return nil, nil
	}
	return result,nil
}

func (w *worloadPart) InsertOrUpdate(db *sqlx.Tx,name, parttype, parameters string, workloadid int64)(workloadPartId int64, err error){

	var workloadpart WorkloaPartSchema
	var workloadpart_id int64

	query := "select id,name,type,workload_id from workload_part where workload_id=?"
	wkpart_insert := "insert into workload_part (name,type,parameters,workload_id) values(?,?,?,?)"
	wkpart_update := "update workload_part set parameters=?,name=?,type=? where id=? and workload_id=?"
	err=db.Get(&workloadpart,query,workloadid)
	if err == sql.ErrNoRows {
		relid,err:=db.Exec(wkpart_insert,name,parttype,parameters,workloadid)
		if err !=nil {
			log.Errorf("Database insert workloadpart error:%v, sql:%v\n",err,wkpart_insert)
			return 0,err
		}
		workloadpart_id,_ = relid.LastInsertId()
	}else {
		_,err:=db.Exec(wkpart_update,parameters,name,parttype,workloadpart.Id,workloadid)
		if err !=nil {
			log.Errorf("Database update workloadpart error:%v, sql:%v\n",err,wkpart_update)
			return 0,err
		}
		workloadpart_id = int64(workloadpart.Id)
	}
	return workloadpart_id,nil
}