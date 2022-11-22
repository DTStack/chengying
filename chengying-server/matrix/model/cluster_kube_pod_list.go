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
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"k8s.io/api/core/v1"
	"time"
)

type deployKubePodList struct {
	dbhelper.DbTable
}

var DeployKubePodList = &deployKubePodList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CLUSTER_KUBE_POD_LIST},
}

type ClusterKubeLabels struct {
	DeployUuid     string `db:"deploy_uuid" json:"deploy_uuid"`
	Pid            int    `db:"pid" json:"pid"`
	ClusterId      int    `db:"clusterId" json:"cluster_id"`
	ProductName    string `db:"product_name" json:"product_name"`
	ProductVersion string `db:"product_version" json:"product_version"`
	ServiceName    string `db:"service_name" json:"service_name"`
	ServiceVersion string `db:"service_version" json:"service_version"`
	ServiceGroup   string `db:"group" json:"group"`
}

type ClusterKubePod struct {
	Id             int               `db:"id" json:"id"`
	Pid            int               `db:"pid" json:"pid"`
	ClusterId      int               `db:"clusterId" json:"clusterId"`
	NameSpace      string            `db:"namespace" json:"namespace"`
	ProductName    string            `db:"product_name" json:"product_name"`
	ProductVersion string            `db:"product_version" json:"product_version"`
	ServiceName    string            `db:"service_name" json:"service_name"`
	ServiceVersion string            `db:"service_version" json:"service_version"`
	HostIp         string            `db:"host_ip" json:"host_ip"`
	PodId          string            `db:"pod_id" json:"pod_id"`
	PodName        string            `db:"pod_name" json:"pod_name"`
	PodKey         string            `db:"pod_key" json:"pod_key"`
	SelfLink       string            `db:"self_link" json:"self_link"`
	PodIp          string            `db:"pod_ip" json:"pod_ip"`
	Phase          string            `db:"phase" json:"phase"`
	Message        string            `db:"message" json:"message"`
	IsDeleted      int               `db:"is_deleted" json:"is_deleted"`
	UpdateTime     dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime     dbhelper.NullTime `db:"create_time" json:"create_time"`
}

func (l *deployKubePodList) UpdateOrCreate(key string, pod *v1.Pod, labels *ClusterKubeLabels) (error, int) {
	info := ClusterKubePod{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("clusterId", labels.ClusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", pod.GetNamespace())
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_name", labels.ProductName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", labels.ServiceName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("pod_id", pod.GetUID())

	err := l.GetWhere(nil, whereCause, &info)
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"pid":             labels.Pid,
			"clusterId":       labels.ClusterId,
			"namespace":       pod.GetNamespace(),
			"product_name":    labels.ProductName,
			"product_version": labels.ProductVersion,
			"service_name":    labels.ServiceName,
			"service_version": labels.ServiceVersion,
			"host_ip":         pod.Status.HostIP,
			"pod_id":          pod.GetUID(),
			"pod_name":        pod.Name,
			"pod_key":         key,
			"self_link":       pod.GetSelfLink(),
			"pod_ip":          pod.Status.PodIP,
			"phase":           pod.Status.Phase,
			"message":         pod.Status.Message,
			"update_time":     time.Now(),
			"create_time":     time.Now(),
		})
		if err != nil {
			return err, -1
		}
		seq, _ := ret.LastInsertId()
		return nil, int(seq)
	} else if err == nil {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"pid":             labels.Pid,
			"product_version": labels.ProductVersion,
			"service_version": labels.ServiceVersion,
			"phase":           pod.Status.Phase,
			"message":         pod.Status.Message,
			"self_link":       pod.GetSelfLink(),
			"host_ip":         pod.Status.HostIP,
			"pod_ip":          pod.Status.PodIP,
			"update_time":     time.Now(),
		}, false)
		return err, info.Id
	} else {
		return err, -1
	}
	return nil, -1
}

func (l *deployKubePodList) Delete(pod *v1.Pod, labels *ClusterKubeLabels) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("clusterId", labels.ClusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", pod.GetNamespace())
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_name", labels.ProductName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", labels.ServiceName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("pod_id", pod.GetUID())

	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"update_time": time.Now(),
		"is_deleted":  1,
	}, false)
	if err != nil {
		log.Errorf("[deployKubePodList] Delete err: %v", err)
		return err
	}
	return nil
}

func (l *deployKubePodList) GetPodInfoByAgentId(agentId string) (error, *ClusterKubePod) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("pod_id", agentId)
	info := ClusterKubePod{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployKubePodList) GetByPodKey(key string) (*ClusterKubePod, error) {
	info := &ClusterKubePod{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("pod_key", key).And().Equal("is_deleted", 0), info)
	return info, err
}

func (l *deployKubePodList) DeleteByKey(key string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("pod_key", key)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"update_time": time.Now(),
		"is_deleted":  1,
	}, false)
	if err != nil {
		log.Errorf("[deployKubePodList] Delete err: %v", err)
		return err
	}
	return nil
}
