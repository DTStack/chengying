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
	"fmt"
	"strings"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
)

type mysqlIpList struct {
	dbhelper.DbTable
}

var DeployMysqlIpList = &mysqlIpList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_MYSQL_IP_LIST},
}

type MysqlSetIpParam struct {
	Ip string `json:"ip"`
}

type DeployMysqlIpInfo struct {
	ID          int               `db:"id"`
	ClusterId   int               `db:"cluster_id"`
	NameSpace   string            `db:"namespace"`
	ProductName string            `db:"product_name"`
	MysqlIpList string            `db:"mysql_ip_list"`
	UpdateDate  dbhelper.NullTime `db:"update_time"`
	CreateDate  dbhelper.NullTime `db:"create_time"`
}

func (l *mysqlIpList) GetMysqlIpListByName(productName string, clusterId int, namespace string) (error, *DeployMysqlIpInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("product_name", productName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", clusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", namespace)
	info := DeployMysqlIpInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *mysqlIpList) SetMysqlIp(productName, ipList string, clusterId int, namespace string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("product_name", productName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", clusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", namespace)

	err, _ := l.GetMysqlIpListByName(productName, clusterId, namespace)
	//if err == nil {
	//	//	err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
	//	//		"mysql_ip_list": ipList,
	//	//		"update_time":   time.Now(),
	//	//	}, false)
	//	//} else {
	if err == sql.ErrNoRows {
		_, err = l.InsertWhere(dbhelper.UpdateFields{
			"product_name":  productName,
			"mysql_ip_list": ipList,
			"update_time":   time.Now(),
			"create_time":   time.Now(),
			"cluster_id":    clusterId,
			"namespace":     namespace,
		})
	}
	if err != nil {
		log.Errorf("[SetMysqlIp] err: %v", err)
		return err
	}
	return nil
}

func (l *mysqlIpList) Delete(namespace, productName string, clusterId int) error {
	query := "DELETE from " + TBL_DEPLOY_MYSQL_IP_LIST + " "
	query += "WHERE cluster_id=? and namespace=? and product_name=?"
	_, err := l.GetDB().Exec(query, clusterId, namespace, productName)
	if err != nil {
		return err
	}
	return nil
}

func (l *mysqlIpList) GetMysqlIpList(clusterId int, productName string) ([]string, error) {
	var ipList string
	query := fmt.Sprintf("SELECT mysql_ip_list FROM %s WHERE cluster_id=? and product_name=?", l.TableName)
	if err := l.GetDB().Get(&ipList, query, clusterId, productName); err != nil {
		log.Errorf("%v", err)
		if err == sql.ErrNoRows {
			var result = make([]string, 0)
			return result, nil
		} else {
			return nil, err
		}
	}
	return strings.Split(ipList, ","), nil
}
