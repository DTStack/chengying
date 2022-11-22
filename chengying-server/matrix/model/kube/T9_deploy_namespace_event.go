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
	"time"
)

var (
	insertEventSql = "insert into deploy_namespace_event (type,reason,resource,message,namespace_id,time) values (:type, :reason, :resource, :message, :namespace_id,:time)"
	pageEventSql = "select * from deploy_namespace_event where namespace_id = ? order by time desc, id desc limit?,?"
	selectCountSql = "select count(*) from deploy_namespace_event where namespace_id = :namespace_id"
	cleanEventsSql = "delete from deploy_namespace_event where DATE_SUB(CURDATE(),INTERVAL ? DAY) >= DATE(time)"
	insertEventSts *sqlx.NamedStmt
	selectCountSts *sqlx.NamedStmt
	cleanEventsSts *sqlx.Stmt
	DeployNamespaceEvent = &deployNamespaceEvent{
		PrepareFunc: prepareDeloyNamespaceEvent,
	}
)

type deployNamespaceEvent struct {
	PrepareFunc
}

type DeployNamespaceEventSchema struct {
	Id			int 		`db:"id"`
	Type		string 		`db:"type"`
	Reason  	string  	`db:"reason"`
	Resource 	string 		`db:"resource"`
	Message 	string 		`db:"message"`
	NamespaceId   int 	`db:"namespace_id"`
	Time		time.Time	`db:"time"`
}

func prepareDeloyNamespaceEvent() error{
	var err error
	insertEventSts,err = model.USE_MYSQL_DB().PrepareNamed(insertEventSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_event]: init sql: %s , error %v",insertEventSql,err)
		return err
	}
	selectCountSts,err = model.USE_MYSQL_DB().PrepareNamed(selectCountSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_event]: init sql %s, error %v",selectCountSql,err)
		return err
	}
	cleanEventsSts,err = model.USE_MYSQL_DB().Preparex(cleanEventsSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_event]: init sql %s,err %v",deleteNamespaceSql,err)
	}
	return nil
}

func (e *deployNamespaceEvent) Insert(tbsc *DeployNamespaceEventSchema) error{
	_,err := insertEventSts.Exec(tbsc)
	if err != nil{
		log.Errorf("[deploy_namespace_event]: insert sql exec error %v",err)
		return err
	}
	return nil
}

func (e *deployNamespaceEvent) PageSelect(namespaceId, start, limit int)([]DeployNamespaceEventSchema,error){
	rows,err := model.USE_MYSQL_DB().Queryx(pageEventSql,namespaceId,start,limit)
	if err != nil{
		log.Errorf("[deploy_namespace_event]: pageselect sql %s, namespaceid %d, start %d, limit %d error %v",
			pageEventSql,namespaceId,start,limit,err)
		return nil,err
	}
	tbscs := []DeployNamespaceEventSchema{}
	for rows.Next(){
		tbsc := &DeployNamespaceEventSchema{}
		if err = rows.StructScan(tbsc);err != nil{
			log.Errorf("[deploy_namespace_event]: convert to DeployNamespaceEventSchema error %v",err)
			return nil,err
		}
		tbscs = append(tbscs,*tbsc)
	}
	return tbscs,nil
}

func (e *deployNamespaceEvent) SelectCount(namespaceid int) (int,error){
	var count int
	tbsc := &DeployNamespaceEventSchema{
		NamespaceId: namespaceid,
	}
	err := selectCountSts.Get(&count,tbsc)
	if err != nil{
		if err == sql.ErrNoRows{
			return 0,nil
		}
		log.Errorf("[deploy_namespace_event]: select count error: %v",err)
		return 0,err
	}
	return count,nil
}

func (e *deployNamespaceEvent) CleanHistory(day int) error{

	_,err := cleanEventsSts.Exec(day)
	if err != nil{
		log.Errorf("[deploy_namespace_event]: clean history error: %v",err)
		return err
	}
	return nil
}
