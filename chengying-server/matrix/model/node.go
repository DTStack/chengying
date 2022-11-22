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
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"strconv"
	"time"
)

type deployNodeList struct {
	dbhelper.DbTable
}

var DeployNodeList = &deployNodeList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_NODE},
}

type NodeInfo struct {
	ID         int       `db:"id" json:"id"`
	SidecarId  string    `db:"sid" json:"sid"`
	HostName   string    `db:"hostname" json:"hostname"`
	Ip         string    `db:"ip" json:"ip"`
	Status     int       `db:"status" json:"status"`
	Steps      int       `db:"steps" json:"steps"`
	ErrorMsg   string    `db:"errorMsg" json:"errorMsg"`
	IsDeleted  int       `db:"isDeleted" json:"-"`
	UpdateDate base.Time `db:"updated" json:"updated"`
	CreateDate base.Time `db:"created" json:"created"`
	Group      string    `db:"group" json:"group"`
}

func (l *deployNodeList) InsertNodeRecord(sidecarId, hostName, ip string) (error, int) {
	whereCause := dbhelper.WhereCause{}
	info := NodeInfo{}
	err := l.GetWhere(nil, whereCause.Equal("sid", sidecarId), &info)
	if err != nil {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"sid":      sidecarId,
			"hostname": hostName,
			"ip":       ip,
			"updated":  time.Now(),
			"created":  time.Now(),
		})
		if err != nil {
			return err, -1
		}
		seq, _ := ret.LastInsertId()
		return nil, int(seq)
	} else {
		err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("sid", sidecarId), dbhelper.UpdateFields{
			"sid":       sidecarId,
			"ip":        ip,
			"hostname":  hostName,
			"updated":   time.Now(),
			"isDeleted": 0,
		}, false)
		return err, info.ID
	}
}

func (l *deployNodeList) GetNodeInfoBySId(sid string) (error, *NodeInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("sid", sid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	info := NodeInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployNodeList) GetNodeInfoById(aid int) (error,*NodeInfo){
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id",aid).And().
							Equal("isDeleted",0)
	info := &NodeInfo{}
	err := l.GetWhere(nil,whereCause,info)
	if err != nil{
		log.Errorf("[deployNodeList] get nodeinfo by id error %v",err)
	}
	return err,info
}

func (l *deployNodeList) GetNodeInfoByNodeIp(ip string) (error, *NodeInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("ip", ip)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	info := NodeInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployNodeList) UpdateWithAid(aid int, sid, hostName, ip string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", aid), dbhelper.UpdateFields{
		"sid":       sid,
		"hostname":  hostName,
		"ip":        ip,
		"updated":   time.Now(),
		"isDeleted": 0,
	}, false)
	if err != nil {
		log.Errorf("[deployNodeList] UpdateWithAid err: %v", err)
		return err
	}
	return nil
}

func (l *deployNodeList) UpdateStatus(aid int, status int, msg string) {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", aid), dbhelper.UpdateFields{
		"status":   status,
		"errorMsg": msg,
		"updated":  time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployNodeList] UpdateStatus err: %v", err)
	}
}

func (l *deployNodeList) GetDeployNodeSidByClusterIdAndMode(clusterId, mode int) (string, error) {
	var sid, query string
	if strconv.Itoa(mode) == host.KUBERNETES_MODE {
		query = "SELECT deploy_node.sid FROM deploy_cluster_host_rel " +
			"LEFT JOIN deploy_node ON deploy_cluster_host_rel.sid = deploy_node.sid " +
			"LEFT JOIN sidecar_list ON sidecar_list.id = deploy_node.sid " +
			"WHERE clusterId = ? and deploy_cluster_host_rel.is_deleted = 0 and deploy_node.isDeleted = 0 and deploy_node.updated>=DATE_SUB(NOW(),INTERVAL 3 MINUTE) limit 1"
	} else {
		query = "SELECT deploy_node.sid FROM deploy_cluster_host_rel " +
			"LEFT JOIN deploy_host ON deploy_cluster_host_rel.sid = deploy_host.sid " +
			"LEFT JOIN deploy_node ON deploy_node.ip = deploy_host.ip " +
			"LEFT JOIN sidecar_list ON sidecar_list.id = deploy_node.sid " +
			"WHERE clusterId = ? and deploy_cluster_host_rel.is_deleted = 0 and deploy_node.isDeleted = 0 and deploy_host.updated>=DATE_SUB(NOW(),INTERVAL 3 MINUTE) limit 1"
	}
	err := USE_MYSQL_DB().Get(&sid, query, clusterId)
	log.Infof("[NODE_DEPLOY] sid:%v", sid)
	return sid, err
}

func (l *deployNodeList) DeleteWithIp(ip string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("ip", ip), dbhelper.UpdateFields{
		"isDeleted": 1,
	}, false)
	return err
}

func (l *deployNodeList) UpdateUpdatedWithSid(sid string) {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("sid", sid), dbhelper.UpdateFields{
		"updated": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployNodeList] UpdateUpdatedWithSid err: %v", err)
	}
}

func (l *deployNodeList) DeleteWithId(id int)error{
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id",id),dbhelper.UpdateFields{
		"isDeleted":1,
	},false)
	if err != nil && err != sql.ErrNoRows{
		log.Errorf("[deployNodeList] update isDeleted error : %v",err)
		return err
	}
	return nil
}