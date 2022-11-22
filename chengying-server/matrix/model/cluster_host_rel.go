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
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"fmt"
	"time"
)

type deployClusterHostRel struct {
	dbhelper.DbTable
}

var DeployClusterHostRel = &deployClusterHostRel{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CLUSTER_HOST_REL},
}

type ClusterHostRel struct {
	Id         int               `db:"id" json:"id"`
	Sid        string            `db:"sid" json:"sid"`
	ClusterId  int               `db:"clusterId" json:"clusterId"`
	Roles      string            `db:"roles" json:"roles"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

var _getRelListFields = utils.GetTagValues(ClusterHostRel{}, "db")

func (l *deployClusterHostRel) InsertClusterHostRel(clusterId int, sid, roles string) (int, error) {
	rel := ClusterHostRel{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("sid", sid).And().
		NotEqual("clusterId", clusterId).And().
		Equal("is_deleted", 0), &rel)

	if err == nil {
		log.Errorf("该主机已经接入了其他集群 sid: %s", sid)
		return 0, fmt.Errorf("该主机已经接入了其他集群")
	}

	err = l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("sid", sid).And().
		Equal("clusterId", clusterId).And().
		Equal("is_deleted", 0), &rel)
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"sid":         sid,
			"clusterId":   clusterId,
			"roles":       roles,
			"update_time": time.Now(),
			"create_time": time.Now(),
		})
		if err != nil {
			apibase.ThrowDBModelError(err)
			return -1, err
		}
		id, _ := ret.LastInsertId()
		return int(id), err
	} else {
		err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("sid", sid).And().Equal("clusterId", clusterId), dbhelper.UpdateFields{
			"roles":       roles,
			"update_time": time.Now(),
		}, false)
		return rel.Id, err
	}
}

func (l *deployClusterHostRel) GetClusterHostRelList(clusterId int) ([]ClusterHostRel, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("is_deleted", 0)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("clusterId", clusterId)
	rows, _, err := l.SelectWhere(_getRelListFields, whereCause, nil)
	list := []ClusterHostRel{}
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		info := ClusterHostRel{}
		err = rows.StructScan(&info)
		if err != nil {
			return list, err
		}
		list = append(list, info)
	}
	return list, nil

}

func (l *deployClusterHostRel) GetClusterHostRelBySid(sid string) (ClusterHostRel, error) {
	rel := ClusterHostRel{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("sid", sid).And().
		Equal("is_deleted", 0), &rel)

	if err != nil {
		return rel, err
	}
	return rel, nil
}

func (l *deployClusterHostRel) GetClusterHostRelByIp(ip string) (ClusterHostRel, error) {
	query := fmt.Sprintf("select a.* from deploy_cluster_host_rel as a"+
		" inner join (select sid from deploy_host where status != %d and ip = '%s') as b"+
		" on a.sid = b.sid where a.is_deleted = 0", host.SidecarOffline, ip)
	rel := ClusterHostRel{}
	if err := l.GetDB().Get(&rel, query); err != nil {
		return rel, err
	}
	log.Debugf("%+v", rel)
	return rel, nil
}

func (l *deployClusterHostRel) GetClusterHostMap() map[string]int {
	rows, _ := l.GetDB().Queryx("SELECT b.ip,a.clusterId FROM deploy_cluster_host_rel as a inner join deploy_host as b on a.sid = b.sid")
	defer rows.Close()
	item := make(map[string]int)
	for rows.Next() {
		d := struct {
			Ip        string `db:"ip"`
			ClusterId int    `db:"clusterId"`
		}{}
		_ = rows.StructScan(&d)
		item[d.Ip] = d.ClusterId
	}
	return item
}

func (l *deployClusterHostRel) DeleteWithSid(sid string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("sid", sid), dbhelper.UpdateFields{
		"is_deleted": 1,
	}, false)
	return err
}

func (l *deployClusterHostRel) UpdateRolesWithSid(sid, roles string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("sid", sid), dbhelper.UpdateFields{
		"roles":       roles,
		"update_time": time.Now(),
	}, false)

	return err
}

type InspectClusterHostRel struct {
	IP      string  `db:"ip" json:"ip"`
	DirSize float64 `json:"dir_size"`
	Sid     string  `db:"sid" json:"sid"`
}

func (l *deployClusterHostRel) GetInspectClusterHostRelList(clusterId int) ([]InspectClusterHostRel, error) {
	query := fmt.Sprintf("select dh.ip,dchr.sid from %s dchr "+
		"left join %s dh on dh.sid=dchr.sid and dchr.is_deleted=0 "+
		"where clusterId = ? and dh.isDeleted=0", l.TableName, DeployHostList.TableName)
	res := make([]InspectClusterHostRel, 0)
	if err := l.GetDB().Select(&res, query, clusterId); err != nil {
		return res, err
	}
	return res, nil
}
