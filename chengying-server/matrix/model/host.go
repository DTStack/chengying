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
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"fmt"
	"time"
)

type deployHostList struct {
	dbhelper.DbTable
}

var DeployHostList = &deployHostList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_HOST},
}

type HostInfo struct {
	ID         int            `db:"id" json:"id"`
	SidecarId  string         `db:"sid" json:"sid"`
	HostName   string         `db:"hostname" json:"hostname"`
	Ip         string         `db:"ip" json:"ip"`
	Status     int            `db:"status" json:"status"`
	Steps      int            `db:"steps" json:"steps"`
	ErrorMsg   string         `db:"errorMsg" json:"errorMsg"`
	IsDeleted  int            `db:"isDeleted" json:"-"`
	UpdateDate base.Time      `db:"updated" json:"updated"`
	CreateDate base.Time      `db:"created" json:"created"`
	Group      string         `db:"group" json:"group"`
	RoleList   sql.NullString `db:"role_list" json:"-"`
}

type HostRunningInfo struct {
	MemSize   int64          `db:"mem_size"`
	MemUsage  int64          `db:"mem_usage"`
	DiskUsage sql.NullString `db:"disk_usage"`
	CpuCores  int            `db:"cpu_cores"`
	CpuUsage  float64        `db:"cpu_usage"`
	Load1     float64        `db:"load1"`
	Status    int            `db:"status"`
	LocalIp   string         `db:"local_ip"`
	Id        int            `db:"id"`
	Updated   base.Time      `db:"updated"`
}

func (l *deployHostList) AutoCreateAid(host, group string) (error, int) {
	whereCause := dbhelper.WhereCause{}
	info := HostInfo{}
	err := l.GetWhere(nil, whereCause.Equal("ip", host), &info)
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"ip":      host,
			"updated": time.Now(),
			"created": time.Now(),
			"group":   group,
		})
		if err != nil {
			return err, -1
		}
		seq, _ := ret.LastInsertId()
		return nil, int(seq)
	} else if err == nil {
		if info.Group != group {
			err = l.UpdateWhere(whereCause.Equal("ip", host), dbhelper.UpdateFields{
				"updated": time.Now(),
				"group":   group,
			}, false)
		}
		return err, info.ID
	} else {
		return err, -1
	}
}

func (l *deployHostList) InsertHostRecord(sidecarId, hostName, ip string) (error, int) {
	whereCause := dbhelper.WhereCause{}
	info := HostInfo{}
	err := l.GetWhere(nil, whereCause.Equal("ip", ip), &info)
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
		err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("ip", ip), dbhelper.UpdateFields{
			"sid":       sidecarId,
			"hostname":  hostName,
			"ip":        ip,
			"updated":   time.Now(),
			"isDeleted": 0,
		}, false)
		return err, info.ID
	}
}

func (l *deployHostList) GetHostInfoById(aid string) (error, *HostInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", aid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	info := HostInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployHostList) GetHostInfoBySid(sid string) (error, *HostInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("sid", sid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	info := HostInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployHostList) UpdateRoleBySid(sid, roleList string) error {
	updateSql := fmt.Sprintf("UPDATE %s SET role_list = ? WHERE sid = ?", l.TableName)
	_, err := l.GetDB().Exec(updateSql, roleList, sid)
	if err != nil {
		return err
	}
	return nil
}

func (l *deployHostList) GetHostListBySid(sid string) (error, []HostInfo, int) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("sid", sid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	rows, totalRecords, err := l.SelectWhere(_getHostListFields, whereCause, nil)
	if err != nil {
		return err, nil, 0
	}
	defer rows.Close()
	list := []HostInfo{}
	for rows.Next() {
		info := HostInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return err, nil, 0
		}
		list = append(list, info)
	}
	return nil, list, totalRecords
}

func (l *deployHostList) GetHostInfoByIp(ip string) (error, *HostInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("ip", ip)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	info := HostInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployHostList) GetHostInfoByIpAndStatus(ip string, statusLow int, statusHigh int) (error, *HostInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("ip", ip)
	whereCause = whereCause.And()
	whereCause = whereCause.GreaterAndEqualThan("status", statusHigh)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	info := HostInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployHostList) GetHostInfoBySidAndStatus(sid string, status int) (error, *HostInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("sid", sid)
	whereCause = whereCause.And()
	whereCause = whereCause.GreaterThan("status", status)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("isDeleted", 0)
	info := HostInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

var _getHostListFields = utils.GetTagValues(HostInfo{}, "db")

func (l *deployHostList) GetHostList(pagination *apibase.Pagination) ([]HostInfo, int) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("isDeleted", 0)
	rows, totalRecords, err := l.SelectWhere(_getHostListFields, whereCause, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()

	list := []HostInfo{}
	for rows.Next() {
		info := HostInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return list, totalRecords
}

func (l *deployHostList) GetGroupList() (groups []string) {
	if err := l.GetDB().Select(&groups, "SELECT DISTINCT `group` FROM "+l.TableName); err != nil {
		apibase.ThrowDBModelError(err)
	}
	return
}

func (l *deployHostList) UpdateStatus(aid int, status int, msg string) {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", aid), dbhelper.UpdateFields{
		"status":   status,
		"errorMsg": msg,
		"updated":  time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateStatus err: %v", err)
	}
}

func (l *deployHostList) UpdateSteps(aid int, step int) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", aid)
	whereCause = whereCause.And()
	whereCause = whereCause.LittleThan("steps", step)

	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"steps": step,
	}, false)
	if err != sql.ErrNoRows {
		log.Errorf("[deployHostList] UpdateSteps err: %v", err)
	}
}

func (l *deployHostList) UpdateStatusWithSid(sid string, status int, msg string) {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("sid", sid), dbhelper.UpdateFields{
		"status":   status,
		"errorMsg": msg,
		"updated":  time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateStatusWithSid err: %v", err)
	}
}

func (l *deployHostList) UpdateUpdatedWithSid(sid string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("sid", sid), dbhelper.UpdateFields{
		"updated": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateUpdatedWithSid err: %v", err)
	}
	return err
}

func (l *deployHostList) UpdateStatusWithAid(aid int, status int, msg string) {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", aid), dbhelper.UpdateFields{
		"status":   status,
		"errorMsg": msg,
		"updated":  time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateStatusWithAid err: %v", err)
	}
}

func (l *deployHostList) UpdateWithAid(aid int, sid, hostName, ip string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", aid), dbhelper.UpdateFields{
		"sid":       sid,
		"hostname":  hostName,
		"ip":        ip,
		"updated":   time.Now(),
		"isDeleted": 0,
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateWithAid err: %v", err)
		return err
	}
	return nil
}

func (l *deployHostList) UpdateGroupWithAid(aid int, group string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", aid), dbhelper.UpdateFields{
		"group":   group,
		"updated": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateGroupWithAid err: %v", err)
		return err
	}
	return nil
}

func (l *deployHostList) UpdateGroup(old, new string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("group", old), dbhelper.UpdateFields{
		"group":   new,
		"updated": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateGroup err: %v", err)
		return err
	}
	return nil
}

func (l *deployHostList) DeleteWithAid(aid int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", aid), dbhelper.UpdateFields{
		"status":    host.SidecarOffline,
		"isDeleted": 1,
		"role_list": "",
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] DeleteHost err: %v", err)
		return err
	}
	return nil
}

//根据集群 id 获取集群下所有未删除的主机信息
func (l *deployHostList) GetHostListByClusterId(clusterId int) ([]HostInfo, error) {
	var list []HostInfo
	sql := "SELECT deploy_host.id as id," +
		"deploy_host.sid as sid," +
		"deploy_host.hostname as hostname," +
		"deploy_host.ip as ip," +
		"deploy_host.status as status," +
		"deploy_host.steps as steps," +
		"deploy_host.errorMsg as errorMsg," +
		"deploy_host.isDeleted as isDeleted," +
		"deploy_host.updated as updated," +
		"deploy_host.created as created," +
		"deploy_host.`group` as `group`," +
		"deploy_host.role_list as role_list" +
		" from deploy_host left join deploy_cluster_host_rel host_rel on deploy_host.sid = host_rel.sid where clusterId = ? and isDeleted = 0"
	err := l.GetDB().Select(&list, sql, clusterId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (l *deployHostList) GetRunHostListByClusterId(id int) (hosts []HostInfo) {
	sql := "SELECT d.* FROM deploy_host AS d INNER JOIN dtagent.deploy_cluster_host_rel c " +
		"ON c.sid = d.sid AND c.clusterId = ? " +
		"WHERE d.status = 3 AND d.steps = 3 AND d.isDeleted = 0"
	if err := l.GetDB().Select(&hosts, sql, id); err != nil {
		log.Errorf("[deployHostList] GetRunHostListByClusterId err: %v", err)
		return nil
	}
	return
}

type HostInspectInfo struct {
	Ip          string `db:"ip" json:"ip"`
	IsRunning   bool   `db:"is_running" json:"status"`
	ServiceList string `db:"service_list" json:"service_list"`
}

func (l *deployHostList) GetInspectNodeInfoByClusterId(clusterId int) ([]HostInspectInfo, error) {
	var list []HostInspectInfo
	query := fmt.Sprintf("select deploy_host.ip as ip," +
		"IF(TIMESTAMPDIFF(MINUTE,deploy_host.updated,NOW())<3,true,false) as is_running, " +
		"IFNULL(GROUP_CONCAT(DISTINCT(concat(deploy_product_list.product_name,'-',deploy_instance_list.service_name))),'') as service_list " +
		"from deploy_host " +
		"left join deploy_cluster_host_rel host_rel on deploy_host.sid = host_rel.sid " +
		"left join deploy_instance_list on deploy_host.sid = deploy_instance_list.sid " +
		"left join deploy_product_list on deploy_instance_list.pid = deploy_product_list.id " +
		"where clusterId = ? and isDeleted = 0 " +
		"group by deploy_host.ip")
	err := l.GetDB().Select(&list, query, clusterId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (l *deployHostList) GetHostRunningInfoListByClusterId(clusterId int) ([]HostRunningInfo, error) {
	query := "SELECT deploy_host.id,mem_size,mem_usage,disk_usage,cpu_cores,cpu_usage,load1,local_ip,deploy_host.status,deploy_host.updated " +
		"FROM deploy_cluster_host_rel " +
		"LEFT JOIN deploy_host ON deploy_cluster_host_rel.sid = deploy_host.sid " +
		"LEFT JOIN sidecar_list ON sidecar_list.id = deploy_host.sid " +
		"WHERE deploy_cluster_host_rel.clusterId = ? and deploy_cluster_host_rel.is_deleted = 0"
	list := make([]HostRunningInfo, 0)
	if err := l.GetDB().Select(&list, query, clusterId); err != nil {
		return nil, err
	}
	return list, nil
}
