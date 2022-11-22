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
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"fmt"
	"time"
)

type deployClusterList struct {
	dbhelper.DbTable
}

var DeployClusterStatus = map[int]string{
	0:  "Waiting",
	1:  "Pending",
	2:  "Running",
	-2: "Error",
}

var DeployClusterList = &deployClusterList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CLUSTER_LIST},
}

const (
	DEPLOY_CLUSTER_TYPE_HOSTS      = "hosts"
	DEPLOY_CLUSTER_TYPE_KUBERNETES = "kubernetes"
	DEPLOY_CLUSTER_STATUS_WAITING  = 0
	DEPLOY_CLUSTER_STATUS_PENDING  = 1
	DEPLOY_CLUSTER_STATUS_RUNNING  = 2
	DEPLOY_CLUSTER_STATUS_ERROR    = -2
	DEPLOY_CLUSTER_MODE_IMPORT     = 1
)

type ClusterInfo struct {
	Id         int               `db:"id" json:"id"`
	Name       string            `db:"name" json:"name"`
	Type       string            `db:"type" json:"type"`
	Mode       int               `db:"mode" json:"mode"`
	Version    string            `db:"version" json:"version"`
	Desc       string            `db:"desc" json:"desc"`
	Tags       string            `db:"tags" json:"tags"`
	Configs    sql.NullString    `db:"configs" json:"configs"`
	Yaml       sql.NullString    `db:"yaml" json:"yaml"`
	Status     int               `db:"status" json:"status"`
	ErrorMsg   string            `db:"errorMsg" json:"errorMsg"`
	CreateUser string            `db:"create_user" json:"create_user"`
	UpdateUser string            `db:"update_user" json:"update_user"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

type K8sConfigInfo struct {
	NetworkPlugin string `json:"network_plugin"`
}

type K8sCreateInfo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Mode       int    `json:"mode"`
	Version    string `json:"version"`
	Desc       string `json:"desc"`
	Tags       string `json:"tags"`
	Configs    string `json:"configs"`
	Yaml       string `json:"yaml"`
	Status     int    `json:"status"`
	ErrorMsg   string `json:"errorMsg"`
	CreateUser string `json:"create_user"`
}

func (l *deployClusterList) InsertHostCluster(cluster ClusterInfo, userName string) (int, error) {
	info := ClusterInfo{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("name", cluster.Name).And().
		Equal("type", "hosts").And().
		Equal("is_deleted", 0), &info)
	if err != nil && err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"name":        cluster.Name,
			"desc":        cluster.Desc,
			"tags":        cluster.Tags,
			"type":        "hosts",
			"create_user": userName,
			"update_user": userName,
		})
		if err != nil {
			apibase.ThrowDBModelError(err)
			return -1, err
		}
		id, _ := ret.LastInsertId()
		return int(id), err
	} else if err == nil {
		return info.Id, fmt.Errorf("主机集群:%v 已存在", cluster.Name)
	} else {
		return -1, err
	}
}

func (l *deployClusterList) DeleteHostClusterById(id int) error {
	info, err := l.GetClusterInfoById(id)
	if err != nil {
		return err
	}

	if info.Status != 0 {
		return fmt.Errorf("Running或者Error、Pending状态的集群不支持删除")
	}

	err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)

	return err
}

func (l *deployClusterList) UpdateHostCluster(cluster ClusterInfo, userName string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", cluster.Id), dbhelper.UpdateFields{
		"name":        cluster.Name,
		"desc":        cluster.Desc,
		"tags":        cluster.Tags,
		"update_time": time.Now(),
		"update_user": userName,
	}, false)

	return err
}

func (l *deployClusterList) UpdateClusterStatus(clusterId, status int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", clusterId), dbhelper.UpdateFields{
		"status":      status,
		"update_time": time.Now(),
	}, false)

	return err
}

func (l *deployClusterList) GetClusterInfoById(id int) (ClusterInfo, error) {
	info := ClusterInfo{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("id", id).And().Equal("is_deleted", 0), &info)
	return info, err
}

func (l *deployClusterList) GetClusterNameById(id int) (clusterName string) {
	l.GetDB().Get(&clusterName, "SELECT `name` FROM "+DeployClusterList.TableName+" where `id` = ?", id)
	return
}

func (l *deployClusterList) GetClusterList() ([]ClusterInfo, error) {
	info := make([]ClusterInfo, 0)
	err := DeployClusterList.GetDB().Select(&info, "select * from "+DeployClusterList.TableName+
		" where  is_deleted=0")
	return info, err
}

func (l *deployClusterList) SelectHostClusterList() ([]ClusterInfo, error) {
	var info []ClusterInfo
	err := DeployClusterList.GetDB().Select(&info, "select * from "+DeployClusterList.TableName+
		" where is_deleted=0 and type ='hosts'")
	return info, err
}

func (l *deployClusterList) GetDeployedClusterList() ([]ClusterInfo, error) {
	query := "SELECT distinct deploy_cluster_list.* from deploy_cluster_list " +
		"INNER JOIN deploy_cluster_product_rel ON deploy_cluster_list.id = deploy_cluster_product_rel.clusterId " +
		"WHERE deploy_cluster_product_rel.is_deleted=0 AND deploy_cluster_list.is_deleted=0"
	info := make([]ClusterInfo, 0)
	err := USE_MYSQL_DB().Select(&info, query)
	return info, err
}

func (l *deployClusterList) GetDeployedClusterListByUserId(userId int) ([]ClusterInfo, error) {
	query := "SELECT distinct deploy_cluster_list.* from deploy_cluster_list " +
		"INNER JOIN deploy_cluster_product_rel ON deploy_cluster_list.id = deploy_cluster_product_rel.clusterId " +
		"INNER JOIN user_cluster_right ON user_cluster_right.cluster_id = deploy_cluster_list.id " +
		"WHERE user_cluster_right.user_id=? AND user_cluster_right.is_deleted=0 " +
		"AND deploy_cluster_product_rel.is_deleted=0 AND deploy_cluster_list.is_deleted=0"
	info := make([]ClusterInfo, 0)
	err := USE_MYSQL_DB().Select(&info, query, userId)
	return info, err
}

func (l *deployClusterList) DeleteK8sClusterById(id int) error {
	info, err := l.GetClusterInfoById(id)
	if err != nil {
		return err
	}
	if info.Status != 0 {
		return fmt.Errorf("只可删除Waiting状态的k8s集群")
	}
	err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *deployClusterList) UpdateK8sCluster(cluster K8sCreateInfo) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", cluster.Id), dbhelper.UpdateFields{
		"name":        cluster.Name,
		"desc":        cluster.Desc,
		"tags":        cluster.Tags,
		"yaml":        cluster.Yaml,
		"configs":     cluster.Configs,
		"version":     cluster.Version,
		"update_time": time.Now(),
	}, false)

	return err
}

func (l *deployClusterList) UpdateVersionById(id int, version string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"version":     version,
		"update_time": time.Now(),
	}, false)

	return err
}

func (l *deployClusterList) UpdateYamlById(id int, yaml string) {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"yaml":        yaml,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("update yaml error:%v", err)
	}
}

func (l *deployClusterList) GetClusterInfoByHostId(hostid int) (*ClusterInfo, error) {
	query := "select cluster.* from deploy_cluster_list as cluster left join " +
		"(select host.id as hostid, rel.clusterId from deploy_host as host " +
		"left join deploy_cluster_host_rel as rel on host.sid = rel.sid where host.isDeleted = 0 and rel.is_deleted = 0) as hc" +
		"on cluster.id = hc.clusterId where hc.hostid = ?"
	info := &ClusterInfo{}
	if err := l.GetDB().Get(info, query, hostid); err != nil {
		log.Errorf("Get clusterinfo by hostid error: %v", err)
		return nil, err
	}
	return info, nil
}

func (l *deployClusterList) GetClusterListByUserId(userId int) ([]ClusterInfo, error) {
	list := make([]ClusterInfo, 0)
	query := "select deploy_cluster_list.* from deploy_cluster_list " +
		"left join user_cluster_right on user_cluster_right.cluster_id = deploy_cluster_list.id " +
		"WHERE user_cluster_right.user_id = ? " +
		"and user_cluster_right.is_deleted = 0 " +
		"and deploy_cluster_list.is_deleted = 0"
	err := USE_MYSQL_DB().Select(&list, query, userId)
	return list, err
}
