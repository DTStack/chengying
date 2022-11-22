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
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"fmt"
)

type deployClusterSmoothUpgradeProductRel struct {
	dbhelper.DbTable
}

var DeployClusterSmoothUpgradeProductRel = &deployClusterSmoothUpgradeProductRel{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CLUSTER_SMOOTH_UPGRADE_PRODUCT_REL},
}

type ClusterSmoothUpgradeProductRel struct {
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

func (l *deployClusterSmoothUpgradeProductRel) GetCurrentProductByProductNameClusterId(productName string, clusterId int) (*DeployProductListInfo, error) {
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.product, p.schema, p.product_type, p.create_time, " +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_smooth_upgrade_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ? AND p.product_name = ? AND r.is_deleted=0 "
	product := &DeployProductListInfo{}
	err := USE_MYSQL_DB().Get(product, query, clusterId, productName)
	return product, err
}

func (l *deployClusterSmoothUpgradeProductRel) CheckProductReadyForDeploy(productName string) error {
	var status []string
	query := "SELECT p.status FROM deploy_product_list as p " +
		"LEFT JOIN deploy_cluster_smooth_upgrade_product_rel ON p.id=deploy_cluster_smooth_upgrade_product_rel.pid " +
		"WHERE p.product_name=? AND deploy_cluster_smooth_upgrade_product_rel.is_deleted=0 "

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

func (l *deployClusterSmoothUpgradeProductRel) GetCurrentProductByProductNameClusterIdNamespace(productName string, clusterId int, namespace string) (*DeployProductListInfo, error) {
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.product, p.schema, p.product_type, p.create_time, " +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_smooth_upgrade_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ? AND p.product_name = ? AND r.is_deleted=0 AND r.namespace=? "
	product := &DeployProductListInfo{}
	err := USE_MYSQL_DB().Get(product, query, clusterId, productName, namespace)
	return product, err
}

func (l *deployClusterSmoothUpgradeProductRel) GetProductByPid(pid int) ([]DeployProductListInfo, error) {
	query := "SELECT deploy_product_list.* from deploy_cluster_smooth_upgrade_product_rel " +
		"LEFT JOIN deploy_product_list ON deploy_product_list.id = deploy_cluster_smooth_upgrade_product_rel.pid " +
		"WHERE pid = ? AND deploy_cluster_smooth_upgrade_product_rel.is_deleted=0 "
	productList := make([]DeployProductListInfo, 0)
	err := USE_MYSQL_DB().Select(&productList, query, pid)
	return productList, err
}

func (l *deployClusterSmoothUpgradeProductRel) GetSmoothUpgradeProductRelByClusterIdAndPid(clusterId, pid int) (*ClusterSmoothUpgradeProductRel, error) {
	info := ClusterSmoothUpgradeProductRel{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("clusterId", clusterId).And().
		Equal("pid", pid).And().
		Equal("is_deleted", 0)
	if err := l.GetWhere(nil, whereCause, &info); err != nil {
		return &info, err
	}
	return &info, nil
}

func (l *deployClusterSmoothUpgradeProductRel) GetByPidAndClusterIdNamespace(pid, clusterId int, namespace string) (ClusterSmoothUpgradeProductRel, error) {
	info := ClusterSmoothUpgradeProductRel{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("pid", pid).And().
		Equal("clusterId", clusterId).And().
		Equal("is_deleted", 0).And().
		Equal("namespace", namespace), &info)
	return info, err
}
