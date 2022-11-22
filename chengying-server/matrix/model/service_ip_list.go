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
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
)

type serviceIpList struct {
	dbhelper.DbTable
}

var DeployServiceIpList = &serviceIpList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_SERVICE_IP_LIST},
}

type ServiceSetIpParam struct {
	Ip string `json:"ip"`
}

type DeployServiceIpInfo struct {
	ID          int               `db:"id"`
	ClusterId   int               `db:"cluster_id"`
	NameSpace   string            `db:"namespace"`
	ProductName string            `db:"product_name"`
	ServiceName string            `db:"service_name"`
	IpList      string            `db:"ip_list"`
	UpdateDate  dbhelper.NullTime `db:"update_time"`
	CreateDate  dbhelper.NullTime `db:"create_time"`
}

func (d DeployServiceIpInfo) Value() (driver.Value, error) {
	return []interface{}{d.ClusterId, d.NameSpace, d.ProductName, d.ServiceName, d.IpList, d.CreateDate, d.UpdateDate}, nil
}

func (l *serviceIpList) GetServiceIpInfoById(id int) (error, *DeployServiceIpInfo) {
	whereCause := dbhelper.WhereCause{}
	info := DeployServiceIpInfo{}
	err := l.GetWhere(nil, whereCause.Equal("id", id), &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *serviceIpList) GetServiceIpListByPNameAndClusterId(productName string, clusterId int) ([]DeployServiceIpInfo, error) {
	sql := fmt.Sprintf("select * from %s where cluster_id = ? and product_name = ?", TBL_DEPLOY_SERVICE_IP_LIST)
	deployServiceIpInfos := make([]DeployServiceIpInfo, 0)
	err := l.GetDB().Select(&deployServiceIpInfos, sql, clusterId, productName)
	if err != nil {
		return nil, err
	}
	return deployServiceIpInfos, nil
}

func (l *serviceIpList) GetServiceIpListByName(productName, serviceName string, clusterId int, namespace string) (error, *DeployServiceIpInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("product_name", productName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", serviceName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", clusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", namespace)
	info := DeployServiceIpInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *serviceIpList) GetPodServiceIpListByName(namespace, productName, serviceName string, clusterId int) (error, *DeployServiceIpInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("namespace", namespace)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_name", productName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", serviceName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", clusterId)
	info := DeployServiceIpInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *serviceIpList) SetServiceIp(productName, serviceName, ipList string, clusterId int, namespace string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("product_name", productName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", serviceName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", clusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", namespace)

	err, _ := l.GetServiceIpListByName(productName, serviceName, clusterId, namespace)
	if err == nil {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"ip_list":     ipList,
			"update_time": time.Now(),
		}, false)
	} else {
		_, err = l.InsertWhere(dbhelper.UpdateFields{
			"product_name": productName,
			"service_name": serviceName,
			"ip_list":      ipList,
			"update_time":  time.Now(),
			"create_time":  time.Now(),
			"cluster_id":   clusterId,
			"namespace":    namespace,
		})
	}
	if err != nil {
		log.Errorf("[SetServiceIp] err: %v", err)
		return err
	}
	return nil
}

func (l *serviceIpList) SetPodServiceIp(namespace, productName, serviceName, ipList string, clusterId int) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("namespace", namespace)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_name", productName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", serviceName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", clusterId)

	err, _ := l.GetPodServiceIpListByName(namespace, productName, serviceName, clusterId)
	if err == nil {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"ip_list":     ipList,
			"update_time": time.Now(),
		}, false)
	} else {
		log.Errorf("%v", err.Error())
		_, err = l.InsertWhere(dbhelper.UpdateFields{
			"product_name": productName,
			"service_name": serviceName,
			"ip_list":      ipList,
			"update_time":  time.Now(),
			"create_time":  time.Now(),
			"cluster_id":   clusterId,
			"namespace":    namespace,
		})
	}
	if err != nil {
		log.Errorf("[SetServiceIp] err: %v", err)
		return err
	}
	return nil
}

func (l *serviceIpList) Delete(namespace, productName, serviceName string, clusterId int) error {
	query := "DELETE from " + TBL_DEPLOY_SERVICE_IP_LIST + " "
	query += "WHERE cluster_id=? and namespace=? and product_name=? and service_name=?"
	_, err := l.GetDB().Exec(query, clusterId, namespace, productName, serviceName)
	if err != nil {
		return err
	}
	return nil
}

func (l *serviceIpList) DeleteByClusterIdNamespaceProduct(namespace, productName string, clusterId int) error {
	query := "DELETE from " + TBL_DEPLOY_SERVICE_IP_LIST + " "
	query += "WHERE cluster_id=? and namespace=? and product_name=?"
	_, err := l.GetDB().Exec(query, clusterId, namespace, productName)
	if err != nil {
		return err
	}
	return nil
}

func (l *serviceIpList) HostOffByIp(ip string) error {
	const selectSqlStr = "SELECT id, ip_list FROM deploy_service_ip_list"
	const deleteSqlStr = "DELETE FROM deploy_service_ip_list WHERE ip_list = ?"
	const updateSqlStr = "UPDATE deploy_service_ip_list SET ip_list = ? WHERE id = ?"

	_, err := l.GetDB().Exec(deleteSqlStr, ip)
	if err != nil {
		return err
	}

	type idAndIpList struct {
		Id     int    `db:"id"`
		IpList string `db:"ip_list"`
	}
	connectList := []idAndIpList{}
	err = l.GetDB().Select(&connectList, selectSqlStr)
	if err != nil {
		return err
	}
	for _, connect := range connectList {
		ipL := strings.Split(connect.IpList, ",")
		for idx, currentIp := range ipL {
			if currentIp == ip {
				ipL = append(ipL[:idx], ipL[idx+1:]...)
				newIpL := strings.Replace(strings.Trim(fmt.Sprint(ipL), "[]"), " ", ",", -1)
				_, err := l.GetDB().Exec(updateSqlStr, newIpL, connect.Id)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func (l *serviceIpList) GetServiceIpList(clusterId int, productName, serviceName string) ([]string, error) {
	var ipList string
	query := fmt.Sprintf("SELECT ip_list FROM %s WHERE cluster_id=? and product_name=? and service_name=?", l.TableName)
	if err := l.GetDB().Get(&ipList, query, clusterId, productName, serviceName); err != nil {
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

func (l *serviceIpList) CountServiceIpByClusterId(clusterId int) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(1) FROM %s WHERE cluster_id=? ", l.TableName)
	err := l.GetDB().Get(&count, query, clusterId)
	if err != nil && err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (l *serviceIpList) BatchInsertServiceIp(infos []interface{}, tx *sqlx.Tx) error {
	argsList := make([]string, 0)
	for range infos {
		argsList = append(argsList, "(?)")
	}
	query, args, err := sqlx.In(fmt.Sprintf("INSERT INTO %s (cluster_id, namespace, product_name, service_name, ip_list, "+
		"create_time, update_time) VALUES %s", l.TableName, strings.Join(argsList, ",")), infos...)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, args...)
	return err
}
