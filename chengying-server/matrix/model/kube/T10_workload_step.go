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

var(
	getStepSql = "select * from workload_step where workloadpart_id = :workloadpart_id order by id asc"
	getStepByTypeSql = "select * from workload_step where workloadpart_id = :workloadpart_id and type = :type"
	getStepSts *sqlx.NamedStmt
	getStepByTypeSts *sqlx.NamedStmt
	WorkloadStep = &workloadStep{
		PrepareFunc: prepareWorkload,
	}
)

func prepareWorkload()error{
	var err error
	getStepSts,err = model.USE_MYSQL_DB().PrepareNamed(getStepSql)
	if err != nil{
		log.Errorf("[kube workload_step]: init sql: %s , error %v",getStepSql,err)
		return err
	}
	getStepByTypeSts,err = model.USE_MYSQL_DB().PrepareNamed(getStepByTypeSql)
	if err != nil{
		log.Errorf("[[kube workload_step]: init sql: %s, error %v]",getStepByTypeSql,err)
	}
	return nil
}

type workloadStep struct {
	PrepareFunc
}
type WorloadStepSchema struct {
	Id 		int			`db:"id""`
	Name 	string 		`db:"name"`
	Type 	string 		`db:"type"`
	Action 	string 		`db:"action"`
	Object 	string 		`db:"object"`
	WorkloadPartId	int `db:"workloadpart_id"`
}

func (w *workloadStep) SelectType(workloadPartId int, typ string) ([]WorloadStepSchema,error){
	arg := &WorloadStepSchema{
		WorkloadPartId: workloadPartId,
		Type: typ,
	}
	rows,err := getStepByTypeSts.Queryx(arg)
	if err != nil{
		log.Errorf("[kube workload_step] select sql %s, workloadpartid %d, type %s, error %v",getStepByTypeSql,workloadPartId,typ,err)
		return nil,err
	}
	result := []WorloadStepSchema{}
	for rows.Next(){
		tbsc := WorloadStepSchema{}
		if err = rows.StructScan(&tbsc);err != nil{
			log.Errorf("[kube workload_step] sql %s result failed struct scan to workload_part error %v",getStepByTypeSql,err)
			return nil,err
		}
		result = append(result,tbsc)
	}
	if len(result) == 0{
		return nil, nil
	}
	return result,nil
}

func (w *workloadStep) Select(workloadPartId int) ([]WorloadStepSchema,error){
	arg := &WorloadStepSchema{
		WorkloadPartId: workloadPartId,
	}
	rows,err := getStepSts.Queryx(arg)
	if err != nil{
		log.Errorf("[kube workload_step] select sql %s, workloadpartid %d, error %v",getStepSql,workloadPartId,err)
		return nil,err
	}
	result := []WorloadStepSchema{}
	for rows.Next(){
		tbsc := WorloadStepSchema{}
		if err = rows.StructScan(&tbsc);err != nil{
			log.Errorf("[kube workload_step] sql %s result failed struct scan to workload_part error %v",getStepSql,err)
			return nil,err
		}
		result = append(result,tbsc)
	}
	if len(result) == 0{
		return nil, nil
	}
	return result,nil
}

func(w *workloadStep) InsertOrUpdate(db *sqlx.Tx ,name, parttype, action,object string, workloadPartid int64)error{

	var workloadstep WorloadStepSchema

	query := "select id,name,type,action,workloadpart_id from workload_step where name=? and workloadpart_id=?"
	wkstep_insert := "insert into workload_step (name,type,action,object,workloadpart_id) values(?,?,?,?,?)"
	wkstep_update := "update workload_step set object=?,type=?,action=? where id=? and workloadpart_id=?"
	err:=db.Get(&workloadstep,query,name,workloadPartid)
	if err == sql.ErrNoRows {
		_,err:=db.Exec(wkstep_insert,name,parttype,action,object,workloadPartid)
		if err !=nil {
			log.Errorf("Database insert workloadpart error:%v, sql:%v\n",err,wkstep_insert)
			return err
		}
	}else {
		_,err:=db.Exec(wkstep_update,object,parttype,action,workloadstep.Id,workloadPartid)
		if err !=nil {
			log.Errorf("Database update workloadpart error:%v, sql:%v\n",err,wkstep_update)
			return err
		}
	}
	return nil
}