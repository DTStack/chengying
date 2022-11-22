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
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

var (
	insertNamespaceSql = "insert into deploy_namespace_list (namespace,type,registry_id,ip,port,cluster_id,is_deleted,status,user) values (:namespace,:type,:registry_id,:ip,:port,:cluster_id,:is_deleted,:status,:user)"
	selectNamespaceSql = "select * from deploy_namespace_list where namespace = :namespace and cluster_id = :cluster_id and is_deleted=0"
	updateNamespaceStatusSql = "update deploy_namespace_list set status = :status where namespace = :namespace and cluster_id = :cluster_id and is_deleted=0"
	updateNamespaceSql = "update deploy_namespace_list set namespace=:namespace,type=:type,registry_id=:registry_id,ip=:ip,port=:port,user =:user,`status`=:status,update_time=:update_time where id=:id"
	deleteNamespaceSql = "update deploy_namespace_list set is_deleted=1, update_time=:update_time where namespace = :namespace and cluster_id = :cluster_id and is_deleted=0"
	selectLikedNamespaceDescSql = "select * from deploy_namespace_list where cluster_id = :cluster_id and namespace like :namespace and is_deleted=0 and `status` =:status order by update_time desc"
	selectLikedNamespaceAscSql = "select * from deploy_namespace_list where cluster_id = :cluster_id and namespace like :namespace and is_deleted=0 and `status` =:status order by update_time asc"
	selectAllSql = "select * from deploy_namespace_list where is_deleted=0"
	updateIpAndPortSql = "update deploy_namespace_list set ip=:ip,port=:port where namespace=:namespace and cluster_id=:cluster_id and is_deleted=0"
	//selectNamespacesByClusterSql = "select * from deploy_namespace_list where cluster_id = :cluster_id and is_deleted=0 order by update_time desc"
	updateNamespaceStatusSts *sqlx.NamedStmt
	insertNamespaceSts *sqlx.NamedStmt
	selectNamespaceSts *sqlx.NamedStmt
	updateNamespaceSts *sqlx.NamedStmt
	deleteNamespaceSts *sqlx.NamedStmt
	selectLikedNamespaceDescSts *sqlx.NamedStmt
	selectLIkedNamespaceAscSts *sqlx.NamedStmt
	updateIpAndPortSts *sqlx.NamedStmt
	selectAllSts *sqlx.Stmt
	//selectNamespacesByClusterSts *sqlx.NamedStmt
	DeployNamespaceList = &deployNamespaceList{
		PrepareFunc: prepareDeployNamespaceList,
	}
)
type deployNamespaceList struct {
	PrepareFunc
}
type DeployNamespaceListSchema struct {
	ClusterId   int 	`db:"cluster_id"`
	Status      string 	`db:"status"`
	IsDeleted   int 	`db:"is_deleted"`
	User        string  `db:"user"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
	view.NamespaceSaveReq
}

func prepareDeployNamespaceList() error{
	var err error
	insertNamespaceSts,err = model.USE_MYSQL_DB().PrepareNamed(insertNamespaceSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql: %s , error %v",insertNamespaceSql,err)
		return err
	}
	selectNamespaceSts,err = model.USE_MYSQL_DB().PrepareNamed(selectNamespaceSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql: %s, error %v",selectNamespaceSql,err)
		return err
	}
	updateNamespaceStatusSts,err = model.USE_MYSQL_DB().PrepareNamed(updateNamespaceStatusSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql: %s, error %v",updateNamespaceStatusSql,err)
		return err
	}
	updateNamespaceSts,err = model.USE_MYSQL_DB().PrepareNamed(updateNamespaceSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql: %s, error %v",updateNamespaceSql,err)
		return err
	}
	deleteNamespaceSts,err = model.USE_MYSQL_DB().PrepareNamed(deleteNamespaceSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql: %s, error %v",deleteNamespaceSql,err)
		return err
	}
	selectLikedNamespaceDescSts,err = model.USE_MYSQL_DB().PrepareNamed(selectLikedNamespaceDescSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql %s, error %v",selectLikedNamespaceDescSql,err)
		return err
	}
	selectLIkedNamespaceAscSts,err = model.USE_MYSQL_DB().PrepareNamed(selectLikedNamespaceAscSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql %s, error %v",selectLikedNamespaceAscSql,err)
		return err
	}
	selectAllSts,err = model.USE_MYSQL_DB().Preparex(selectAllSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql %s, error %v",selectAllSql,err)
		return err
	}
	updateIpAndPortSts,err = model.USE_MYSQL_DB().PrepareNamed(updateIpAndPortSql)
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: init sql %s, error %v",updateIpAndPortSql,err)
		return err
	}
	//selectNamespacesByClusterSts,err = model.USE_MYSQL_DB().PrepareNamed(selectNamespacesByClusterSql)
	//if err != nil{
	//	log.Errorf("[kube deploy_namespace_list]: init sql: %s, error %v",updateNamespaceSql,err)
	//	return err
	//}
	return nil
}

func (l *deployNamespaceList) Insert(s *DeployNamespaceListSchema) (int,error){
	err := selectNamespaceSts.Get(s,s)
	if err != nil && err == sql.ErrNoRows{
		rs,err := insertNamespaceSts.Exec(s)
		if err != nil{
			log.Errorf("[deploy_namespace_list]: insert sql exec error: %v",err)
			return -1,err
		}
		id,err := rs.LastInsertId()
		if err!= nil{
			log.Errorf("[deploy_namespace_list] get last insert id error :%v",err)
			return -1,err
		}
		return int(id),nil
	}else if err == nil{
		return -1,fmt.Errorf("namespace %s is exist",s.Namespace)
	}else {
		log.Errorf("[deploy_namespace_list]: query sql: %s, value: %+v, error: %v",selectNamespaceSql,*s,err)
		return -1,err
	}
}

func (l *deployNamespaceList) SelectAll() ([]DeployNamespaceListSchema,error){
	rows,err := selectAllSts.Queryx()
	if err != nil{
		log.Errorf("[deploy_namespace_list]: select all sql: %s, error: %v",selectAllSql,err)
		return nil,err
	}
	tbscs := []DeployNamespaceListSchema{}
	for rows.Next(){
		tbsc := DeployNamespaceListSchema{}
		if err := rows.StructScan(&tbsc);err != nil{
			log.Errorf("[deploy_namespace_list]: sql %s result struct scan to deploy_namespace_list error %v",selectAllSql,err)
			return nil,err
		}
		tbscs = append(tbscs,tbsc)
	}
	return tbscs,nil
}

func (l *deployNamespaceList) Update(s *DeployNamespaceListSchema) error{
	s.UpdateTime = time.Now()
	if _,err := updateNamespaceSts.Exec(s);err != nil{
		log.Errorf("[deploy_namespace_list]: update sql: %s, error: %v",updateNamespaceSql,err)
		return err
	}
	return nil
}

func (l *deployNamespaceList) UpdateStatus(s *DeployNamespaceListSchema) error{
	_,err := updateNamespaceStatusSts.Exec(s)
	if err != nil{
		log.Errorf("[deploy_namespace_list]: update sql: %s, error: %v",updateNamespaceStatusSql,err)
		return err
	}
	return nil
}

func (l *deployNamespaceList) GetLike(namespace string, clusterid int, status string, desc bool)([]DeployNamespaceListSchema,error){
	namespace  = "%"+namespace+"%"
	arg := &DeployNamespaceListSchema{
		ClusterId: clusterid,
		Status: status,
		NamespaceSaveReq:view.NamespaceSaveReq{
			Namespace: namespace,
		},
	}
	var rows *sqlx.Rows
	var err error
	if desc{
		rows,err = selectLikedNamespaceDescSts.Queryx(arg)
	}else{
		rows,err = selectLikedNamespaceDescSts.Queryx(arg)
	}
	if err != nil{
		log.Errorf("[kube deploy_namespace_list]: query liked namespace sql %s, value %+v, error %v",selectNamespaceSql,*arg,err)
		return nil,err
	}
	result := []DeployNamespaceListSchema{}
	for rows.Next(){
		tbsc := DeployNamespaceListSchema{}
		if err = rows.StructScan(&tbsc);err != nil{
			log.Errorf("[kube deploy_namespace_list]: sql %s result struct scan to deplyonamespacelist error %v",selectNamespaceSql,err)
			return nil,err
		}
		result = append(result,tbsc)
	}
	if len(result) == 0{
		return nil,nil
	}
	return result,nil
}
func (l *deployNamespaceList) Get(namespace string, clusterid int) (*DeployNamespaceListSchema,error){
	tbsc := &DeployNamespaceListSchema{
		ClusterId:        clusterid,
		NamespaceSaveReq: view.NamespaceSaveReq{
			Namespace: namespace,
		},
	}
	err := selectNamespaceSts.Get(tbsc,tbsc)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil,nil
		}
		log.Errorf("[deploy_namespace_list]: get namespace info, sql %s, value %+v, error %v",selectNamespaceSql,*tbsc,err)
		return nil ,err
	}
	return tbsc,nil
}


func (l *deployNamespaceList) Select(clsuterid,status,desc,namespace,typ string) ([]DeployNamespaceListSchema,error){
	s := "select * from deploy_namespace_list where is_deleted = 0 and cluster_id = "+clsuterid
	if len(status) != 0{
		statusList := strings.Split(status,",")
		statusQuery := ""
		for i,item := range statusList{
			if i == len(statusList) -1 {
				statusQuery = statusQuery + " `status` = '"+item+"'"
			}else{
				statusQuery = statusQuery + " `status` = '"+item+"' or"
			}
		}
		s = s + " and ("+statusQuery+")"
	}
	if len(typ) !=0{
		typList := strings.Split(typ,",")
		typQuery := ""
		for i,item := range typList{
			if i == len(typList) - 1{
				typQuery = typQuery + " type = '"+item+"'"
			}else{
				typQuery = typQuery + " type = '"+item+"' or"
			}
		}
		s = s + " and ("+typQuery+")"
	}
	if len(namespace) !=0 {
		s = s + " and namespace like '%"+namespace+"%'"
	}

	if len(desc) != 0 && desc == "false"{
		s = s + " order by update_time asc"
	}else{
		s = s + " order by update_time desc"
	}
	rows,err := model.USE_MYSQL_DB().Queryx(s)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		log.Errorf("[deploy_namespace_list]: get namespaces info, sql %s, error %v",s,err)
		return nil,err
	}
	tbscs := []DeployNamespaceListSchema{}
	for rows.Next(){
		tbsc := DeployNamespaceListSchema{}
		if err := rows.StructScan(&tbsc); err != nil{
			log.Errorf("[deploy_namespace_list]: struct scan to DeployNamespaceListSchema error :%v",err)
			return nil, err
		}
		tbscs = append(tbscs,tbsc)
	}
	return tbscs,nil
}

func (l *deployNamespaceList) Delete(namespace string, clusterid int) error{
	tbsc := &DeployNamespaceListSchema{
		ClusterId:        clusterid,
		UpdateTime: 	  time.Now(),
		NamespaceSaveReq: view.NamespaceSaveReq{
			Namespace: namespace,
		},
	}
	_,err := deleteNamespaceSts.Exec(tbsc)
	if err != nil{
		log.Errorf("[deploy_namespace_list]: sql %s, value %+v, error %v",deleteNamespaceSql,tbsc,err)
		return err
	}
	return nil
}
