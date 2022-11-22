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
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"strconv"
	"time"
)

type deployAddonsList struct {
	dbhelper.DbTable
}

var DeployAddonsList = &deployAddonsList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_ADDONS_LIST},
}

const (
	ADDON_STATUS_INSTALLING           = "installing"
	ADDON_STATUS_INSTALLED            = "installed"
	ADDON_STATUS_INSTALL_FAIL         = "install fail"
	ADDON_STATUS_UNINSTALLING         = "uninstalling"
	ADDON_STATUS_UNINSTALLED          = "uninstalled"
	ADDON_STATUS_UNINSTALL_FAIL       = "uninstall fail"
	ADDON_STATUS_INSTALLING_CANCELLED = "installing cancelled"
	ADDON_STATUS_RUNNING              = "running"
	ADDON_STATUS_RUN_FAIL             = "run fail"
	ADDON_STATUS_STOPPED              = "stopped"
	ADDON_STATUS_STOPPING             = "stopping"
	ADDON_STATUS_STOP_FAIL            = "stop fail"
	ADDON_STATUS_UPDATE_CONFIG_FAIL   = "update-config fail"
)

type DeployAddonsInfo struct {
	Id            int       `db:"id"`
	Aid           int       `db:"aid"`
	SidecarId     string    `db:"sid" json:"sid"`
	AgentId       string    `db:"agentId"`
	Config        string    `db:"config"`
	AddonType     string    `db:"addon_type"`
	AddonVersion  string    `db:"addon_version"`
	Status        int       `db:"status"`
	StatusMessage string    `db:"status_message"`
	IsDeleted     int       `db:"isDeleted" json:"-"`
	UpdateDate    base.Time `db:"updated" json:"updated"`
	CreateDate    base.Time `db:"created" json:"created"`
}

var _getAddonListFields = utils.GetTagValues(DeployAddonsInfo{}, "db")

func (l *deployAddonsList) NewDeployAddonRecord(aid int64, sid, addonType, addonVersion string, config string) (error, int64, string) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("aid", aid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("sid", sid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("addon_type", addonType)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("addon_version", addonVersion)

	addonId := int64(-1)
	agentId := ""

	err, info := l.GetDeployAddonInfoByWhere(whereCause)
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"aid":           aid,
			"sid":           sid,
			"addon_type":    addonType,
			"addon_version": addonVersion,
			"config":        config,
			"update_time":   time.Now(),
			"create_time":   time.Now(),
		})
		if err != nil {
			return err, addonId, agentId
		}
		addonId, _ = ret.LastInsertId()
	} else {
		err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"aid":           aid,
			"sid":           sid,
			"addon_type":    addonType,
			"addon_version": addonVersion,
			"config":        config,
			"update_time":   time.Now(),
		}, false)
		if err != nil {
			return err, addonId, agentId
		}
		addonId = int64(info.Id)
		agentId = info.AgentId
	}
	return nil, addonId, agentId
}

func (l *deployAddonsList) GetDeployAddonInfoByWhere(cause dbhelper.WhereCause) (error, *DeployAddonsInfo) {
	info := DeployAddonsInfo{}
	err := l.GetWhere(nil, cause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployAddonsList) GetAddonList() ([]DeployAddonsInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.NotEqual("agent_id", "")
	rows, _, err := l.SelectWhere(_getAddonListFields, whereCause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployAddonsInfo{}
	for rows.Next() {
		info := DeployAddonsInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployAddonsList) GetAddonInfoById(id int) (error, *DeployAddonsInfo) {
	whereCause := dbhelper.WhereCause{}
	info := DeployAddonsInfo{}
	err := l.GetWhere(nil, whereCause.Equal("id", id), &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployAddonsList) GetAddonInfoByAgentId(agentId string) (error, *DeployAddonsInfo) {
	whereCause := dbhelper.WhereCause{}
	info := DeployAddonsInfo{}
	err := l.GetWhere(nil, whereCause.Equal("agent_id", agentId), &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployAddonsList) UpdateAddonStatusByAgentId(agentId string, status, statusMsg string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("agent_id", agentId)

	var err error
	err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"status":         status,
		"status_message": statusMsg,
		"update_time":    time.Now(),
	}, false)
	return err
}

func (l *deployAddonsList) UpdateAddonStatusByAgentPerformance(agentId string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("agent_id", agentId)
	whereCause = whereCause.And()
	whereCause = whereCause.NotEqual("status", ADDON_STATUS_RUNNING)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"status":         ADDON_STATUS_RUNNING,
		"status_message": "",
		"update_time":    time.Now(),
	}, false)
	if err != sql.ErrNoRows {
		log.Errorf("[addon] UpdateAddonStatusByAgentPerformance err: %v, agentId: %v", err.Error(), agentId)
	}
	return err
}

func (l *deployAddonsList) UpdateAddonStatusById(id int, status, statusMsg string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", id)

	var err error
	err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"status":         status,
		"status_message": statusMsg,
		"update_time":    time.Now(),
	}, false)

	return err
}

func (l *deployAddonsList) UpdateAddonAgentId(id int, agentId string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", id)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"agent_id":    agentId,
		"update_time": time.Now(),
	}, false)
	return err
}

//func (l *deployAddonsList) GetAddonConfig(id int, configfile string) (interface{}, error) {
//	whereCause := dbhelper.WhereCause{}
//	info := DeployAddonsInfo{}
//	err := l.GetWhere(nil, whereCause.Equal("id", id), &info)
//
//	if err != nil {
//		log.Errorf("[GetAddonConfig] get addon info error:%v, addon id: %v", err.Error(), id)
//		return nil, fmt.Errorf("error:%v, addon id: %v", err.Error(), id)
//	}
//	if info.AgentId != "" && configfile != "" {
//		params := &agent.ExecScriptParams{}
//		params.ExecScript = "#!/bin/sh\ncat " + util.ShellQuote(configfile)
//		params.AgentId = info.AgentId
//		params.Timeout = "10s"
//
//		err, respBody := agent.AgentClient.AgentExec(info.SidecarId, params,"")
//		if err != nil {
//			log.Errorf("[GetAddonConfig] response err: %v, sid: %v", err, info.SidecarId)
//			return nil, err
//		}
//
//		result, exists := respBody.Data.(map[string]interface{})["result"]
//		if !exists {
//			log.Errorf("[GetAddonConfig] server response without result: %v", respBody)
//			return nil, fmt.Errorf("without result: %v", respBody)
//		}
//		failed, exists := result.(map[string]interface{})["failed"]
//		if exists && failed.(bool) == true {
//			return respBody, fmt.Errorf(result.(map[string]interface{})["response"].(string))
//		}
//		return result.(map[string]interface{})["response"].(string), nil
//	}
//	return nil, fmt.Errorf("config file or agent+id is empty , id: %v", id)
//}

func (l *deployAddonsList) DeleteByAgentId(agentId string) error {
	query := "DELETE from " + TBL_DEPLOY_ADDONS_LIST + " "
	query += "WHERE agent_id='" + agentId + "'"
	_, err := l.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (l *deployAddonsList) DeleteById(id int) error {
	query := "DELETE from " + TBL_DEPLOY_ADDONS_LIST + " "
	query += "WHERE id='" + strconv.Itoa(id) + "'"
	_, err := l.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}
