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
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"fmt"
	"strconv"
	"time"
)

type deployClusterProductRel struct {
	dbhelper.DbTable
}

var DeployClusterProductRel = &deployClusterProductRel{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CLUSTER_PRODUCT_REL},
}

type ClusterProductRel struct {
	Id            int               `db:"id" json:"id"`
	Pid           int               `db:"pid" json:"pid"`
	ClusterId     int               `db:"clusterId" json:"clusterId"`
	Namespace     string            `db:"namespace" json:"namespace"`
	ProductParsed []byte            `db:"product_parsed" json:"product_parsed"`
	Status        string            `db:"status" json:"status"`
	DeployUUID    string            `db:"deploy_uuid" json:"deploy_uuid"`
	AlertRecover  int               `db:"alert_recover" json:"alert_recover"`
	UserId        int               `db:"user_id" json:"user_id"`
	IsDeleted     int               `db:"is_deleted" json:"is_deleted"`
	UpdateTime    dbhelper.NullTime `db:"update_time" json:"update_time"`
	DeployTime    dbhelper.NullTime `db:"deploy_time"`
	CreateTime    dbhelper.NullTime `db:"create_time" json:"create_time"`
}

func (l *deployClusterProductRel) GetProductListByClusterIdAndStatus(clusterId int, deployStatus []string) ([]DeployProductListInfo, error) {
	query := "SELECT deploy_product_list.* from deploy_cluster_product_rel " +
		"LEFT JOIN deploy_product_list ON deploy_product_list.id = deploy_cluster_product_rel.pid " +
		"WHERE clusterId = ? ANDis_deleted=0 "
	if len(deployStatus) > 0 {
		query += "AND deploy_cluster_product_rel.status IN ("
		for i, v := range deployStatus {
			if i > 0 {
				query += "," + v
			} else {
				query += v
			}
		}
		query += ")"
	}
	productList := make([]DeployProductListInfo, 0)
	err := USE_MYSQL_DB().Select(&productList, query, clusterId)
	return productList, err
}

func (l *deployClusterProductRel) GetProductByPid(pid int) ([]DeployProductListInfo, error) {
	query := "SELECT deploy_product_list.* from deploy_cluster_product_rel " +
		"LEFT JOIN deploy_product_list ON deploy_product_list.id = deploy_cluster_product_rel.pid " +
		"WHERE pid = ? AND deploy_cluster_product_rel.is_deleted=0 "
	productList := make([]DeployProductListInfo, 0)
	err := USE_MYSQL_DB().Select(&productList, query, pid)
	return productList, err
}

func (l *deployClusterProductRel) GetProductListByClusterId(clusterId int, status string) ([]DeployProductListInfo, error) {
	query := "SELECT deploy_product_list.*,deploy_cluster_product_rel.namespace from deploy_cluster_product_rel " +
		"LEFT JOIN deploy_product_list ON deploy_product_list.id = deploy_cluster_product_rel.pid " +
		"WHERE clusterId = ? AND deploy_cluster_product_rel.is_deleted=0 AND deploy_cluster_product_rel.status = ?"
	productList := make([]DeployProductListInfo, 0)
	err := USE_MYSQL_DB().Select(&productList, query, clusterId, status)
	return productList, err
}

func (l *deployClusterProductRel) GetParentProductNameListByClusterIdNamespace(clusterId int, namespace string) ([]string, error) {
	query := "SELECT distinct deploy_product_list.parent_product_name from deploy_cluster_product_rel " +
		"LEFT JOIN deploy_product_list ON deploy_product_list.id = deploy_cluster_product_rel.pid " +
		"WHERE clusterId = ? AND deploy_cluster_product_rel.is_deleted=0 AND deploy_cluster_product_rel.namespace=? "
	parentProductNameList := make([]string, 0)
	err := USE_MYSQL_DB().Select(&parentProductNameList, query, clusterId, namespace)
	return parentProductNameList, err
}

func (l *deployClusterProductRel) GetCurrentProductByProductNameClusterId(productName string, clusterId int) (*DeployProductListInfo, error) {
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.product, p.schema, p.product_type, p.create_time, " +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ? AND p.product_name = ? AND r.is_deleted=0 "
	product := &DeployProductListInfo{}
	err := USE_MYSQL_DB().Get(product, query, clusterId, productName)
	return product, err
}

func (l *deployClusterProductRel) GetProductsByParentProductNameClusterId(parentProductName string, clusterId int) ([]DeployProductListInfo, error) {
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.product, p.schema, p.product_type, p.create_time, " +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ? AND p.parent_product_name = ? AND r.is_deleted=0 "
	productList := make([]DeployProductListInfo, 0)
	err := USE_MYSQL_DB().Select(&productList, query, clusterId, parentProductName)
	return productList, err
}

func (l *deployClusterProductRel) CheckProductReadyForDeploy(productName string) error {
	var status []string
	query := "SELECT p.status FROM deploy_product_list as p " +
		"LEFT JOIN deploy_cluster_product_rel ON p.id=deploy_cluster_product_rel.pid " +
		"WHERE p.product_name=? AND deploy_cluster_product_rel.is_deleted=0 "

	err := USE_MYSQL_DB().Select(&status, query, productName)

	if err != nil {
		return err
	}
	for _, s := range status {
		if s == PRODUCT_STATUS_DEPLOYING || s == PRODUCT_STATUS_UNDEPLOYING {
			return fmt.Errorf("product %v is deploying or undeploying", productName)
		}
	}
	return nil
}

func (l *deployClusterProductRel) GetListByClusterIdAndStatus(clusterId int, deployStatus []string, namespace string) ([]ClusterProductRel, error) {
	list := make([]ClusterProductRel, 0)
	query := "SELECT * FROM " + DeployClusterProductRel.TableName + " WHERE is_deleted=0 AND clusterId=? "
	if len(deployStatus) > 0 {
		query += "AND deploy_cluster_product_rel.status IN ("
		for i, v := range deployStatus {
			if i > 0 {
				query += ",'" + v + "'"
			} else {
				query += "'" + v + "'"
			}
		}
		query += ")"
	}
	//distinguish k8s multiple namespaces
	if namespace != "" {
		query += "AND namespace=" + "'" + namespace + "'"
	}
	err := DeployClusterList.GetDB().Select(&list, query, clusterId)
	return list, err
}

func (l *deployClusterProductRel) GetByPidAndClusterId(pid, clusterId int) (ClusterProductRel, error) {
	info := ClusterProductRel{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("pid", pid).And().
		Equal("clusterId", clusterId).And().
		Equal("is_deleted", 0), &info)
	return info, err
}

func (l *deployClusterProductRel) GetByPidAndClusterIdNamespacce(pid, clusterId int, namespace string) (ClusterProductRel, error) {
	info := ClusterProductRel{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("pid", pid).And().
		Equal("clusterId", clusterId).And().
		Equal("namespace", namespace).And().
		Equal("is_deleted", 0), &info)
	return info, err
}

func (l *deployClusterProductRel) UpdateStatus(clusterId, pid int, status string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().
		Equal("clusterId", clusterId).And().
		Equal("pid", pid).And().
		Equal("is_deleted", 0), dbhelper.UpdateFields{
		"status":      status,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *deployClusterProductRel) UpdateStatusWithNamespace(clusterId, pid int, namespace, status string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().
		Equal("clusterId", clusterId).And().
		Equal("pid", pid).And().
		Equal("namespace", namespace).And().
		Equal("is_deleted", 0), dbhelper.UpdateFields{
		"status":      status,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *deployClusterProductRel) GetCurrentProductByProductNameClusterIdNamespace(productName string, clusterId int, namespace string) (*DeployProductListInfo, error) {
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.product, p.schema, p.product_type, p.create_time, " +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ? AND p.product_name = ? AND r.is_deleted=0 AND r.namespace=? "
	product := &DeployProductListInfo{}
	err := USE_MYSQL_DB().Get(product, query, clusterId, productName, namespace)
	return product, err
}

func (l *deployClusterProductRel) GetByPidAndClusterIdNamespace(pid, clusterId int, namespace string) (ClusterProductRel, error) {
	info := ClusterProductRel{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("pid", pid).And().
		Equal("clusterId", clusterId).And().
		Equal("is_deleted", 0).And().
		Equal("namespace", namespace), &info)
	return info, err
}

func (d *deployClusterProductRel) GetDeployClusterProductList(
	parentProductName, productName, productVersionLike, productVersion, productType string, clusterId int, deployStatus, productNames []string, namespace string) ([]DeployProductListInfoWithNamespace, error) {

	query := "SELECT p.id,r.namespace,p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.product, p.schema, p.product_type, p.create_time, " +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ?  AND r.is_deleted=0 "

	//if len(deployStatus) > 0 {
	//	query += " AND r.status IN ("
	//	for i, v := range deployStatus {
	//		if i > 0 {
	//			query += ",'" + v + "'"
	//		} else {
	//			query += "'" + v + "'"
	//		}
	//	}
	//	query += ")"
	//}

	//distinguish k8s multiple namespaces
	if namespace != "" {
		query += " AND r.namespace=" + "'" + namespace + "'"
	}

	if productType != "" {
		pType, _ := strconv.Atoi(productType)
		query += fmt.Sprintf(" AND p.product_type= %d", pType)
	}

	if parentProductName != "" {
		query += " AND p.parent_product_name=" + "'" + parentProductName + "'"
	}

	if len(productNames) > 0 {
		query += " AND p.product_name IN ("
		for i, v := range productNames {
			if i > 0 {
				query += ",'" + v + "'"
			} else {
				query += "'" + v + "'"
			}
		}
		query += ")"
	}

	if productName != "" {
		query += " AND p.product_name=" + "'" + productName + "'"
	}

	//if productVersionLike != "" {
	//	query += " AND p.product_version like" + "'%" + productVersionLike + "%'"
	//}

	if productVersion != "" {
		query += " AND p.product_version=" + "'" + productVersion + "'"
	}

	list := make([]DeployProductListInfoWithNamespace, 0)

	stmt, err := USE_MYSQL_DB().Prepare(query)
	if err != nil {
		msg := fmt.Errorf("prepare sql failed, %v", err)
		apibase.ThrowDBModelError(msg)
	}
	defer stmt.Close()
	stmt.Exec(clusterId)

	rows, err := USE_MYSQL_DB().Queryx(query, clusterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		info := DeployProductListInfoWithNamespace{}

		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, err
}

type DeployArchHostInfo struct {
	MemSize            int64          `db:"mem_size"`
	DiskUsage          sql.NullString `db:"disk_usage"`
	CpuCores           int            `db:"cpu_cores"`
	DiskSizeDisplay    string
	CpuCoreSizeDisplay string
	FileSizeDisplay    string
	IP                 string `db:"ip"`
	HostName           string `db:"hostname"`
	OSDisplay          string `db:"os_display"`
	MemSizeDisplay     string
}

func (l *deployClusterProductRel) GetDeployArchInfo() ([]DeployArchHostInfo, error) {

	hostsList := make([]DeployArchHostInfo, 0)
	query := fmt.Sprintf("select dh.ip,dh.hostname,cpu_cores,mem_size,disk_usage,concat(os_platform,os_version) as os_display " +
		"from deploy_cluster_host_rel dchr " +
		"left join deploy_host dh on dchr.sid = dh.sid " +
		"left join sidecar_list sl on sl.id = dh.sid " +
		"where dchr.is_deleted = 0 order by dchr.id asc")
	if err := USE_MYSQL_DB().Select(&hostsList, query); err != nil {
		return nil, fmt.Errorf("[GetDeployArchInfo] Database err: %v", err)
	}
	return hostsList, nil
}

type DeployVersionStruct struct {
	Pid            int    `db:"pid"`
	ProductName    string `db:"product_name"`
	ProductVersion string `db:"product_version"`
}

func (l *deployClusterProductRel) GetDeployVersionInfo() ([]DeployVersionStruct, error) {
	ProductList := make([]DeployVersionStruct, 0)
	query := fmt.Sprintf("select r.pid,p.product_name, p.product_version " +
		"from deploy_cluster_product_rel as r " +
		"left join deploy_product_list as p on p.id = r.pid " +
		"where r.is_deleted=0 and r.status = ? " +
		"group by  p.product_name,p.product_version " +
		"order by p.product_name,p.product_version")
	if err := USE_MYSQL_DB().Select(&ProductList, query, "deployed"); err != nil {
		return nil, fmt.Errorf("[GetDeployVersionInfo] Database err: %v", err)
	}
	return ProductList, nil
}
