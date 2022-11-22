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
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"dtstack.com/dtstack/easymatrix/schema"
	"strings"
)

type deployInstanceList struct {
	dbhelper.DbTable
}

var DeployInstanceList = &deployInstanceList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_INSTANCE_LIST},
}

const (
	INSTANCE_STATUS_INSTALLING             = "installing"
	INSTANCE_STATUS_INSTALLED              = "installed"
	INSTANCE_STATUS_INSTALL_FAIL           = "install fail"
	INSTANCE_STATUS_UPGRADE_FAIL           = "upgrade fail"
	INSTANCE_STATUS_UNINSTALLING           = "uninstalling"
	INSTANCE_STATUS_UNINSTALLED            = "uninstalled"
	INSTANCE_STATUS_UNINSTALL_FAIL         = "uninstall fail"
	INSTANCE_STATUS_INSTALLING_CANCELLED   = "installing cancelled"
	INSTANCE_STATUS_RUNNING                = "running"
	INSTANCE_STATUS_RUN_FAIL               = "run fail"
	INSTANCE_STATUS_HEALTH_CHECKED         = "health-checked"
	INSTANCE_STATUS_HEALTH_CHECK_FAIL      = "health-check fail"
	INSTANCE_STATUS_HEALTH_CHECK_CANCELLED = "health-check cancelled"
	INSTANCE_STATUS_HEALTH_CHECK_WAITING   = "health-check waiting"
	INSTANCE_STATUS_STOPPED                = "stopped"
	INSTANCE_STATUS_STOPPING               = "stopping"
	INSTANCE_STATUS_STOP_FAIL              = "stop fail"
	INSTANCE_STATUS_UPDATE_CONFIG_FAIL     = "update-config fail"
	INSTANCE_ROLLBACK_FAIL                 = "rollback fail"

	INSTANCE_HEALTH_OK      = 1
	INSTANCE_HEALTH_BAD     = 0
	INSTANCE_HEALTH_NOTSET  = -1
	INSTANCE_HEALTH_WAITING = -2

	INSTANCE_EVENT_INSTALL       = "install"
	INSTANCE_EVENT_UNINSTALL     = "uninstall"
	INSTANCE_EVENT_START         = "start"
	INSTANCE_EVENT_STOP          = "stop"
	INSTANCE_EVENT_EXEC          = "exec"
	INSTANCE_EVENT_ERROR         = "error"
	INSTANCE_EVENT_CONFIG_UPDATE = "config update"
	INSTANCE_EVENT_UNKNOWN       = "unkown"

	INSTANCE_NORMAL   = "NORMAL"
	INSTANCE_ABNORMAL = "ABNORMAL"
)

var (
	OUT_OF_EVENTREPORT_STATUS_LIST = []string{INSTANCE_STATUS_INSTALLING, INSTANCE_STATUS_INSTALLED, INSTANCE_STATUS_INSTALL_FAIL, INSTANCE_STATUS_UNINSTALLING, INSTANCE_STATUS_UNINSTALLED, INSTANCE_STATUS_UNINSTALL_FAIL}

	OUT_OF_START_STATUS_LIST = []string{INSTANCE_STATUS_INSTALLING, INSTANCE_STATUS_INSTALLED, INSTANCE_STATUS_INSTALL_FAIL, INSTANCE_STATUS_UNINSTALLING, INSTANCE_STATUS_UNINSTALLED, INSTANCE_STATUS_UNINSTALL_FAIL, INSTANCE_STATUS_RUNNING}
	OUT_OF_STOP_STATUS_LIST  = []string{INSTANCE_STATUS_INSTALLING, INSTANCE_STATUS_INSTALLED, INSTANCE_STATUS_INSTALL_FAIL, INSTANCE_STATUS_UNINSTALLING, INSTANCE_STATUS_UNINSTALLED, INSTANCE_STATUS_UNINSTALL_FAIL, INSTANCE_STATUS_STOPPED}
)

type InstanceAndProductInfo struct {
	DeployInstanceInfo
	ProductName        string `db:"product_name"`
	ProductNameDisplay string `db:"product_name_display"`
	ProductVersion     string `db:"product_version"`
}

type DeployInstanceInfo struct {
	ID                 int               `db:"id"`
	ClusterId          int               `db:"cluster_id"`
	Namespace          string            `db:"namespace"`
	AgentId            string            `db:"agent_id"`
	Sid                string            `db:"sid"`
	Pid                int               `db:"pid"`
	Ip                 string            `db:"ip"`
	Group              string            `db:"group"`
	PrometheusPort     int               `db:"prometheus_port"`
	ServiceName        string            `db:"service_name"`
	ServiceNameDisplay string            `db:"service_name_display"`
	ServiceVersion     string            `db:"service_version"`
	Schema             []byte            `db:"schema"`
	HaRoleCmd          string            `db:"ha_role_cmd"`
	HealthState        int               `db:"health_state"`
	Status             string            `db:"status"`
	StatusMessage      string            `db:"status_message"`
	HeartDate          dbhelper.NullTime `db:"heart_time"`
	UpdateDate         dbhelper.NullTime `db:"update_time"`
	CreateDate         dbhelper.NullTime `db:"create_time"`
}

type DeployInstanceRecordByDeployIdInfo struct {
	DeployInstanceRecordInfo
	Schema sql.NullString `db:"schema"`
}

type DeployInstanceUpdateRecordByUpdateIdInfo struct {
	DeployInstanceUpdateRecordInfo
	Schema sql.NullString `db:"schema"`
}

func (l *deployInstanceList) NewInstanceRecord(clusterId, pid, prometheusPort, healthState int, ip, sid, groupName, serviceName, serviceDisplay, serviceVersion string, haRoleCmd string, schema []byte) (error, int64, string) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("pid", pid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("ip", ip)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("sid", sid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", serviceName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", clusterId)

	instanceId := int64(-1)
	agentId := ""

	if serviceDisplay == "" {
		serviceDisplay = serviceName
	}

	err, info := l.GetInstanceInfoByWhere(whereCause)
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"cluster_id":           clusterId,
			"pid":                  pid,
			"ip":                   ip,
			"sid":                  sid,
			"group":                groupName,
			"prometheus_port":      prometheusPort,
			"service_name":         serviceName,
			"service_name_display": serviceDisplay,
			"service_version":      serviceVersion,
			"ha_role_cmd":          haRoleCmd,
			"schema":               schema,
			"health_state":         healthState,
			"update_time":          time.Now(),
			"create_time":          time.Now(),
		})
		if err != nil {
			return err, instanceId, agentId
		}
		instanceId, _ = ret.LastInsertId()
	} else {
		err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"pid":                  pid,
			"ip":                   ip,
			"sid":                  sid,
			"group":                groupName,
			"prometheus_port":      prometheusPort,
			"service_name":         serviceName,
			"service_name_display": serviceDisplay,
			"service_version":      serviceVersion,
			"ha_role_cmd":          haRoleCmd,
			"schema":               schema,
			"update_time":          time.Now(),
		}, false)
		if err != nil {
			return err, instanceId, agentId
		}
		instanceId = int64(info.ID)
		agentId = info.AgentId
	}
	return nil, instanceId, agentId
}

func (l *deployInstanceList) NewPodInstanceRecord(clusterId, pid, prometheusPort, healthState int, namespace, ip, sid, groupName, serviceName, agentId, serviceVersion, status string, message string, schema []byte) (error, int64, string) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("cluster_id", clusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", namespace)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("pid", pid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("sid", sid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("agent_id", agentId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("service_name", serviceName)

	instanceId := int64(-1)

	err, info := l.GetInstanceInfoByWhere(whereCause)
	if err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"cluster_id":           clusterId,
			"namespace":            namespace,
			"pid":                  pid,
			"agent_id":             agentId,
			"ip":                   ip,
			"sid":                  sid,
			"group":                groupName,
			"prometheus_port":      prometheusPort,
			"service_name":         serviceName,
			"service_name_display": serviceName,
			"service_version":      serviceVersion,
			"schema":               schema,
			"health_state":         healthState,
			"status":               status,
			"status_message":       message,
			"update_time":          time.Now(),
			"create_time":          time.Now(),
		})
		if err != nil {
			return err, instanceId, agentId
		}
		instanceId, _ = ret.LastInsertId()
	} else {
		err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"pid":                  pid,
			"ip":                   ip,
			"sid":                  sid,
			"group":                groupName,
			"prometheus_port":      prometheusPort,
			"service_name":         serviceName,
			"service_name_display": serviceName,
			"service_version":      serviceVersion,
			"schema":               schema,
			"health_state":         healthState,
			"status":               status,
			"status_message":       message,
			"update_time":          time.Now(),
		}, false)
		if err != nil {
			return err, instanceId, agentId
		}
		instanceId = int64(info.ID)
		agentId = info.AgentId
	}
	return nil, instanceId, agentId
}

var _getInstanceListFields = utils.GetTagValues(DeployInstanceInfo{}, "db")

func (l *deployInstanceList) GetInstanceListByPidServiceName(pid, clusterId int, name string) ([]DeployInstanceInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("pid", pid)
	if name != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("service_name", name)
	}
	if clusterId > 0 {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("cluster_id", clusterId)
	}
	rows, _, err := l.SelectWhere(_getInstanceListFields, whereCause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []DeployInstanceInfo{}
	for rows.Next() {
		info := DeployInstanceInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployInstanceList) GetInstanceListByPidGroup(pid, clusterId int, name string) ([]DeployInstanceInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("pid", pid)
	if name != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("group", name)
	}
	if clusterId > 0 {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("cluster_id", clusterId)
	}
	rows, _, err := l.SelectWhere(_getInstanceListFields, whereCause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployInstanceInfo{}
	for rows.Next() {
		info := DeployInstanceInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployInstanceList) GetInstanceListByWhere(cause dbhelper.WhereCause) ([]DeployInstanceInfo, error) {
	rows, _, err := l.SelectWhere(_getInstanceListFields, cause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployInstanceInfo{}
	for rows.Next() {
		info := DeployInstanceInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployInstanceList) GetInstanceList() ([]DeployInstanceInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.NotEqual("agent_id", "")
	rows, _, err := l.SelectWhere(_getInstanceListFields, whereCause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployInstanceInfo{}
	for rows.Next() {
		info := DeployInstanceInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployInstanceList) GetInstanceListByClusterId(clusterId, pid int) ([]DeployInstanceInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.NotEqual("agent_id", "").And().Equal("cluster_id", clusterId).And().Equal("pid", pid)
	rows, _, err := l.SelectWhere(_getInstanceListFields, whereCause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployInstanceInfo{}
	for rows.Next() {
		info := DeployInstanceInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployInstanceList) GetInstanceListByClusterIdNamespace(clusterId, pid int, namespace string) ([]DeployInstanceInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.NotEqual("agent_id", "").And().Equal("cluster_id", clusterId).And().Equal("namespace", namespace).And().Equal("pid", pid)
	rows, _, err := l.SelectWhere(_getInstanceListFields, whereCause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployInstanceInfo{}
	for rows.Next() {
		info := DeployInstanceInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *deployInstanceList) GetInstanceInfoById(id int) (error, *DeployInstanceInfo) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", id)
	info := DeployInstanceInfo{}
	err := l.GetWhere(nil, whereCause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployInstanceList) GetInstanceInfoByAgentId(agentId string) (error, *DeployInstanceInfo) {
	whereCause := dbhelper.WhereCause{}
	info := DeployInstanceInfo{}
	err := l.GetWhere(nil, whereCause.Equal("agent_id", agentId), &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployInstanceList) GetInstanceInfoByWhere(cause dbhelper.WhereCause) (error, *DeployInstanceInfo) {
	info := DeployInstanceInfo{}
	err := l.GetWhere(nil, cause, &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *deployInstanceList) UpdateInstanceHealthCheck(agentId string, isHealth bool) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("agent_id", agentId)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"health_state": isHealth,
		"heart_time":   time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[instance] UpdateInstanceHealthCheck err: %v, agentId: %v", err.Error(), agentId)
	}
	return err
}

func (l *deployInstanceList) UpdateInstanceStatusByAgentId(agentId string, status, statusMsg string, healthState ...int) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("agent_id", agentId)

	var err error
	if healthState != nil {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"health_state":   healthState[0],
			"status":         status,
			"status_message": statusMsg,
			"update_time":    time.Now(),
		}, false)
	} else {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"status":         status,
			"status_message": statusMsg,
			"update_time":    time.Now(),
		}, false)
	}
	return err
}

func (l *deployInstanceList) UpdateInstanceStatusByAgentPerformance(agentId string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("agent_id", agentId)
	whereCause = whereCause.And()
	whereCause = whereCause.NotEqual("status", INSTANCE_STATUS_RUNNING)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"status":         INSTANCE_STATUS_RUNNING,
		"status_message": "",
		"update_time":    time.Now(),
	}, false)
	if err != sql.ErrNoRows {
		log.Errorf("[instance] UpdateInstanceStatusByAgentPerformance err: %v, agentId: %v", err.Error(), agentId)
	}
	return err
}

func (l *deployInstanceList) GetInstanceServiceConfig(instanceId int, configfile string) (interface{}, error) {
	whereCause := dbhelper.WhereCause{}
	info := DeployInstanceInfo{}
	err := l.GetWhere(nil, whereCause.Equal("id", instanceId), &info)

	if err != nil {
		log.Errorf("[GetInstanceServiceConfig] get service schema info error:%v, instance id: %v", err.Error(), instanceId)
		return nil, fmt.Errorf("error:%v, instance id: %v", err.Error(), instanceId)
	}

	serviceSchema := schema.ServiceConfig{}
	err = json.Unmarshal(info.Schema, &serviceSchema)
	if err != nil {
		log.Errorf("[GetInstanceServiceConfig] json.Unmarshal service schema error:%v, instance id: %v", err.Error(), instanceId)
		return nil, fmt.Errorf("schema error:%v, instance id: %v", err.Error(), instanceId)
	}

	if configfile == "" && len(serviceSchema.Instance.ConfigPaths) > 0 {
		configfile = serviceSchema.Instance.ConfigPaths[0]
	}
	if serviceSchema.Instance != nil && configfile != "" {
		params := &agent.ExecScriptParams{}
		params.ExecScript = "#!/bin/sh\ncat " + util.ShellQuote(configfile)
		params.AgentId = info.AgentId
		params.Timeout = "10s"

		err, respBody := agent.AgentClient.AgentExec(info.Sid, params, "")
		if err != nil {
			log.Errorf("[GetInstanceServiceConfig] response err: %v, sid: %v", err, info.Sid)
			return nil, err
		}

		result, exists := respBody.Data.(map[string]interface{})["result"]
		if !exists {
			log.Errorf("[GetInstanceServiceConfig] server response without result: %v", respBody)
			return nil, fmt.Errorf("without result: %v", respBody)
		}
		failed, exists := result.(map[string]interface{})["failed"]
		responseMsg, ok := result.(map[string]interface{})["response"]
		if exists && failed.(bool) == true && ok {
			return respBody, fmt.Errorf(responseMsg.(string))
		}
		if ok {
			return responseMsg.(string), nil
		}
		return "", nil
	}
	return nil, fmt.Errorf("service schema is empty ,instance id: %v", instanceId)
}

func (l *deployInstanceList) GetInstanceBelongService(productName string, serviceName string, clusterId int) ([]InstanceAndProductInfo, int) {
	// 需要表关联查询 deploy_instance_list join deploy_product_list
	query := "SELECT IL.*,PL.product_name, PL.product_name_display, PL.product_version FROM " +
		DeployInstanceList.TableName + " AS IL LEFT JOIN " + DeployProductList.TableName + " AS PL ON IL.pid = PL.id " +
		"WHERE PL.product_name =? AND IL.service_name =? AND IL.cluster_id =?"

	rows, err := DeployInstanceList.GetDB().Queryx(query, productName, serviceName, clusterId)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()

	res := []InstanceAndProductInfo{}
	for rows.Next() {
		info := InstanceAndProductInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}

		res = append(res, info)
	}
	return res, len(res)
}

func (l *deployInstanceList) DeleteByagentId(agentId string) error {
	query := "DELETE from " + TBL_DEPLOY_INSTANCE_LIST + " "
	query += "WHERE agent_id='" + agentId + "'"
	_, err := l.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (l *deployInstanceList) DeleteByInstanceId(instanceId int) error {
	query := "DELETE from " + TBL_DEPLOY_INSTANCE_LIST + " "
	query += "WHERE id='" + strconv.Itoa(instanceId) + "'"
	_, err := l.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (l *deployInstanceList) DeleteByClusterIdPidNamespace(pid, clusterId int, namespace string) error {
	query := "DELETE from " + TBL_DEPLOY_INSTANCE_LIST + " "
	query += "WHERE pid='" + strconv.Itoa(pid) + "' and cluster_id='" + strconv.Itoa(clusterId) + "' and namespace='" + namespace + "'"
	_, err := l.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (l *deployInstanceList) GetDeployInstanceLogById(id string, logFile string, isMatch string, tailNum int) (int, interface{}, error) {
	var isDir int = 0
	var result interface{}
	info := DeployInstanceInfo{}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", id)
	err := l.GetWhere(nil, whereCause, &info)

	if err == sql.ErrNoRows {
		log.Errorf("[GetDeployInstanceLogById] empty instance info by id %v, err %v", id, err.Error())
		return isDir, result, err
	}

	if err != nil {
		log.Errorf("[GetDeployInstanceLogById] error get instance info by id %v, err %v", id, err.Error())
		return isDir, result, err
	}
	sid := info.Sid
	if logFile != "" && sid != "" {
		if isMatch == "true" {
			params := &agent.ExecScriptParams{}
			params.ExecScript = "#!/bin/sh\nif [ -d " + util.ShellQuote(logFile) + " ];then\n echo 'dir'\nelse\nif [ -f " + util.ShellQuote(logFile) + " ];then\n" +
				"echo 'file'\nelse\necho 'null'\nfi\nfi\n "
			params.AgentId = info.AgentId
			params.Timeout = "10s"

			err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
			if err != nil {
				log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
				return isDir, result, err
			}

			res := resqBody.Data.(map[string]interface{})["result"]
			failed, exist := res.(map[string]interface{})["failed"]

			if exist && failed.(bool) == true {
				log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
				err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
				return isDir, result, err
			}

			resonseRes, exists := res.(map[string]interface{})["response"]
			if exists {
				if resonseRes == "dir" || resonseRes == "dir\n" {
					isDir = 1
					params := &agent.ExecScriptParams{}
					params.ExecScript = "#!/bin/sh\nls -m " + util.ShellQuote(logFile)
					params.AgentId = info.AgentId
					params.Timeout = "10s"

					err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
					if err != nil {
						log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
						return isDir, result, err
					}

					res := resqBody.Data.(map[string]interface{})["result"]
					failed, exist := res.(map[string]interface{})["failed"]

					if exist && failed.(bool) == true {
						log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
						err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
						return isDir, result, err
					}
					return isDir, res.(map[string]interface{})["response"], nil
				} else {
					params := &agent.ExecScriptParams{}
					params.ExecScript = "#!/bin/sh\ntail -n " + strconv.Itoa(tailNum) + " " + util.ShellQuote(logFile)
					params.AgentId = info.AgentId
					params.Timeout = "10s"

					err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
					if err != nil {
						log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
						return isDir, result, err
					}

					res := resqBody.Data.(map[string]interface{})["result"]
					failed, exist := res.(map[string]interface{})["failed"]

					if exist && failed.(bool) == true {
						log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
						err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
						return isDir, result, err
					}
					return isDir, res.(map[string]interface{})["response"], nil
				}
			} else {
				if resonseRes == "file" || resonseRes == "file\n" {
					params := &agent.ExecScriptParams{}
					params.ExecScript = "#!/bin/sh\ntail -n " + strconv.Itoa(tailNum) + " " + util.ShellQuote(logFile)
					params.AgentId = info.AgentId
					params.Timeout = "10s"

					err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
					if err != nil {
						log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
						return isDir, result, err
					}

					res := resqBody.Data.(map[string]interface{})["result"]
					failed, exist := res.(map[string]interface{})["failed"]

					if exist && failed.(bool) == true {
						log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
						err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
						return isDir, result, err
					}
					return isDir, res.(map[string]interface{})["response"], nil
				}
				return isDir, result, fmt.Errorf("unknow error")
			}
		} else {
			params := &agent.ExecScriptParams{}
			params.ExecScript = "#!/bin/sh\nls -m " + logFile
			params.AgentId = info.AgentId
			params.Timeout = "10s"

			err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
			if err != nil {
				log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
				return isDir, result, err
			}

			res := resqBody.Data.(map[string]interface{})["result"]
			failed, exist := res.(map[string]interface{})["failed"]

			if exist && failed.(bool) == true {
				log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
				err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
				return isDir, result, err
			}
			return isDir, res.(map[string]interface{})["response"], nil
		}
	}

	return isDir, result, nil
}

func (l *deployInstanceList) GetDeployPodInstanceLog(info *DeployInstanceInfo, logFile string, isMatch string, tailNum int) (int, interface{}, error) {
	var isDir int = 0
	var result interface{}
	product, err := DeployProductList.GetProductInfoById(info.Pid)
	if err != nil {
		log.Errorf("%v", err)
		return isDir, result, err
	}
	err, pod := DeployKubePodList.GetPodInfoByAgentId(info.AgentId)
	if err != nil {
		return isDir, result, err
	}
	sid := info.Sid
	var agentId string
	var baseDir string
	///tmp/dtstack/ 为业务测约定约定写入hostpath的路径
	baseDir = "/host/tmp/dtstack/" + product.ProductName + "/" + info.ServiceName + "/" + pod.PodName + "/"
	if !strings.Contains(logFile, baseDir) {
		logFile = baseDir + logFile
	}
	if logFile != "" && sid != "" {
		if isMatch == "true" {
			params := &agent.ExecScriptParams{}
			params.ExecScript = "#!/bin/sh\nif [ -d " + util.ShellQuote(logFile) + " ];then\n echo 'dir'\nelse\nif [ -f " + util.ShellQuote(logFile) + " ];then\n" +
				"echo 'file'\nelse\necho 'null'\nfi\nfi\n "
			params.AgentId = agentId
			params.Timeout = "10s"

			err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
			if err != nil {
				log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
				return isDir, result, err
			}

			res := resqBody.Data.(map[string]interface{})["result"]
			failed, exist := res.(map[string]interface{})["failed"]

			if exist && failed.(bool) == true {
				log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
				err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
				return isDir, result, err
			}

			resonseRes, exists := res.(map[string]interface{})["response"]
			if exists {
				if resonseRes == "dir" || resonseRes == "dir\n" {
					isDir = 1
					params := &agent.ExecScriptParams{}
					params.ExecScript = "#!/bin/sh\nls -1 " + util.ShellQuote(logFile)
					params.AgentId = agentId
					params.Timeout = "10s"

					err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
					if err != nil {
						log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
						return isDir, result, err
					}

					res := resqBody.Data.(map[string]interface{})["result"]
					failed, exist := res.(map[string]interface{})["failed"]

					if exist && failed.(bool) == true {
						log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
						err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
						return isDir, result, err
					}
					result := strings.Split(res.(map[string]interface{})["response"].(string), "\n")
					return isDir, strings.Join(result, ","), nil
				} else {
					params := &agent.ExecScriptParams{}
					params.ExecScript = "#!/bin/sh\ntail -n " + strconv.Itoa(tailNum) + " " + util.ShellQuote(logFile)
					params.AgentId = agentId
					params.Timeout = "10s"

					err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
					if err != nil {
						log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
						return isDir, result, err
					}

					res := resqBody.Data.(map[string]interface{})["result"]
					failed, exist := res.(map[string]interface{})["failed"]

					if exist && failed.(bool) == true {
						log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
						err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
						return isDir, result, err
					}
					return isDir, res.(map[string]interface{})["response"], nil
				}
			} else {
				if resonseRes == "file" || resonseRes == "file\n" {
					params := &agent.ExecScriptParams{}
					params.ExecScript = "#!/bin/sh\ntail -n " + strconv.Itoa(tailNum) + " " + util.ShellQuote(logFile)
					params.AgentId = agentId
					params.Timeout = "10s"

					err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
					if err != nil {
						log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
						return isDir, result, err
					}

					res := resqBody.Data.(map[string]interface{})["result"]
					failed, exist := res.(map[string]interface{})["failed"]

					if exist && failed.(bool) == true {
						log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
						err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
						return isDir, result, err
					}
					return isDir, res.(map[string]interface{})["response"], nil
				}
				return isDir, result, fmt.Errorf("unknow error")
			}
		} else {
			params := &agent.ExecScriptParams{}
			params.ExecScript = "#!/bin/sh\nls -1 " + logFile
			params.AgentId = agentId
			params.Timeout = "10s"

			err, resqBody := agent.AgentClient.AgentExec(sid, params, "")
			if err != nil {
				log.Errorf("[GetDeployInstanceLogById] get log info err %v", err.Error())
				return isDir, result, err
			}

			res := resqBody.Data.(map[string]interface{})["result"]
			failed, exist := res.(map[string]interface{})["failed"]

			if exist && failed.(bool) == true {
				log.Errorf("[GetDeployInstanceLogById] server response failed %v", resqBody)
				err = fmt.Errorf("%v", res.(map[string]interface{})["response"])
				return isDir, result, err
			}

			result := strings.Split(strings.TrimSpace(res.(map[string]interface{})["response"].(string)), "\n")
			return isDir, strings.Join(result, ","), nil
		}
	}

	return isDir, result, nil
}

func (l *deployInstanceList) UpdateInstanceStatusById(id int, status, statusMsg string, healthState ...int) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", id)

	var err error
	if healthState != nil {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"health_state":   healthState[0],
			"status":         status,
			"status_message": statusMsg,
			"update_time":    time.Now(),
		}, false)
	} else {
		err = l.UpdateWhere(whereCause, dbhelper.UpdateFields{
			"status":         status,
			"status_message": statusMsg,
			"update_time":    time.Now(),
		}, false)
	}

	return err
}

func (l *deployInstanceList) UpdateInstanceAgentId(id int, agentId string) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", id)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"agent_id":    agentId,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *deployInstanceList) UpdateInstancePid(id int, pid int) error {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", id)
	err := l.UpdateWhere(whereCause, dbhelper.UpdateFields{
		"pid":         pid,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *deployInstanceList) GetServiceNameByClusterIdAndPid(clusterId, pid int, namespace string) (map[string]bool, error) {
	names := make([]string, 0)
	err := DeployInstanceList.GetDB().Select(&names, "select distinct service_name from "+DeployInstanceList.TableName+" where cluster_id=? and pid = ? and namespace = ?", clusterId, pid, namespace)
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}
	return nameSet, err
}

func (l *deployInstanceList) GetServiceNameSetByClusterIdAndPid(clusterId, pid int) (map[string]bool, error) {
	serviceNames := make([]string, 0)
	err := DeployInstanceList.GetDB().Select(&serviceNames, "select service_name from "+DeployInstanceList.TableName+
		" where cluster_id=? and pid=?", clusterId, pid)
	serviceNameSet := make(map[string]bool)
	for _, v := range serviceNames {
		serviceNameSet[v] = true
	}
	return serviceNameSet, err
}

func (l *deployInstanceList) DeleteBySid(sid string) error {
	const deleteSqlStr = "DELETE FROM deploy_instance_list where sid = ?"
	_, err := l.GetDB().Exec(deleteSqlStr, sid)
	if err != nil {
		return err
	}
	return nil
}

func (l *deployInstanceList) DeleteByIp(ip string) error {
	const deleteSqlStr = "DELETE FROM deploy_instance_list where ip = ?"
	_, err := l.GetDB().Exec(deleteSqlStr, ip)
	if err != nil {
		return err
	}
	return nil
}

type InstanceInfo struct {
	ServiceName string `db:"service_name"`
	Ip          string `db:"ip"`
	Status      string `db:"status"`
	HealthState int    `db:"health_state"`
	Pid         int    `db:"pid"`
	AgentId     string `db:"agent_id"`
}

func (l *deployInstanceList) FindByProductNameAndClusterId(productName string, clusterId int) ([]InstanceInfo, error) {
	var instanceList []InstanceInfo
	query := "SELECT IL.service_name, IL.status, IL.health_state, IL.ip, IL.pid, IL.agent_id FROM " + DeployInstanceList.TableName +
		" AS IL LEFT JOIN " + DeployProductList.TableName + " AS PL ON IL.pid=PL.id WHERE PL.product_name=? AND " +
		"IL.cluster_id=? ORDER BY service_name"
	if err := USE_MYSQL_DB().Select(&instanceList, query, productName, clusterId); err != nil {
		log.Errorf("get product: %v instance error: %v", productName, err)
		return nil, err
	}
	return instanceList, nil
}

func (l *deployInstanceList) GetInstanceBelongServiceWithNamespace(productName string, serviceName string, clusterId int, namespace string) ([]InstanceAndProductInfo, int) {
	// 需要表关联查询 deploy_instance_list join deploy_product_list
	query := "SELECT IL.*,PL.product_name, PL.product_name_display, PL.product_version FROM " +
		DeployInstanceList.TableName + " AS IL LEFT JOIN " + DeployProductList.TableName + " AS PL ON IL.pid = PL.id " +
		"WHERE PL.product_name =? AND IL.service_name =? AND IL.cluster_id =? AND IL.namespace=?"

	rows, err := DeployInstanceList.GetDB().Queryx(query, productName, serviceName, clusterId, namespace)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()

	res := []InstanceAndProductInfo{}
	for rows.Next() {
		info := InstanceAndProductInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}

		res = append(res, info)
	}
	return res, len(res)
}

type DeployListStruct struct {
	ProductName    string `db:"product_name"`
	ProductVersion string `db:"product_version"`
	ServiceName    string `db:"service_name"`
	ServiceVersion string `db:"service_version"`
	IPs            string `db:"ips"`
}

func (l *deployInstanceList) GetDeployListInfo() ([]DeployListStruct, error) {
	deployList := make([]DeployListStruct, 0)
	query := fmt.Sprintf("select dpl.product_name,dpl.product_version,dil.service_name,dil.service_version," +
		"group_concat(dil.ip order by dil.ip separator '/') as ips " +
		"from " + DeployInstanceList.TableName + " dil " +
		"inner join " + DeployProductList.TableName + " dpl on dil.pid=dpl.id " +
		"inner join  " + TBL_DEPLOY_HOST + " dh on dh.sid=dil.sid  " +
		"group by dil.pid,dil.service_name,dil.service_version;")
	if err := USE_MYSQL_DB().Select(&deployList, query); err != nil {
		return nil, fmt.Errorf("[GetDeployListInfo] Database err: %v", err)
	}
	return deployList, nil
}

func (l *deployInstanceList) GetInspectServiceInfoById(id int) (int, int, error) {
	var UnRunningNum, UnHealthyNum int
	query := fmt.Sprintf("SELECT COUNT(*) AS un_running_count "+
		"FROM %s dil "+
		"WHERE dil.status=? AND cluster_id=? group by service_name", DeployInstanceList.TableName)
	if err := USE_MYSQL_DB().Get(&UnRunningNum, query, INSTANCE_STATUS_RUN_FAIL, id); err != nil && err != sql.ErrNoRows {
		log.Errorf("GetInspectServiceInfoById query: %v, values %d, err: %v", query, id, err)
		return UnRunningNum, UnHealthyNum, err
	}
	query = fmt.Sprintf("SELECT COUNT(*) AS un_healthy_count FROM %s dil "+
		"WHERE  dil.health_state=? AND cluster_id=? ", DeployInstanceList.TableName)
	if err := USE_MYSQL_DB().Get(&UnHealthyNum, query, INSTANCE_STATUS_HEALTH_CHECK_FAIL, id); err != nil && err != sql.ErrNoRows {
		log.Errorf("GetInspectServiceInfoById query: %v, values %d, err: %v", query, id, err)
		return UnRunningNum, UnHealthyNum, err
	}
	return UnRunningNum, UnHealthyNum, nil
}

func (l *deployInstanceList) GetNameNodeConfigByIdAndServiceName(id int) (string, string, error) {
	whereCause := dbhelper.WhereCause{}
	info := DeployInstanceInfo{}
	err := l.GetWhere(nil, whereCause.Equal("cluster_id", id).And().
		Equal("service_name", "hadoop_pkg"), &info)
	serviceSchema := schema.ServiceConfig{}
	err = json.Unmarshal(info.Schema, &serviceSchema)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("[GetNameNodeConfigByIdAndServiceName] json.Unmarshal service schema error:%v, cluster id: %v", err.Error(), id)
		return "", "", err
	}
	var namenodeOpts, datanodeOpts string
	for k, v := range serviceSchema.Config {
		if k == "namenode_opts" {
			for key, value := range v.(map[string]interface{}) {
				if key == "Value" {
					namenodeOpts = value.(string)
				}
			}

		}
		if k == "datanode_opts" {
			for key, value := range v.(map[string]interface{}) {
				if key == "Value" {
					datanodeOpts = value.(string)
				}
			}

		}

	}
	namenodeOpts = strings.Split(strings.Split(namenodeOpts, " ")[0], "-Xmx")[1]
	datanodeOpts = strings.Split(strings.Split(datanodeOpts, " ")[0], "-Xmx")[1]
	return namenodeOpts, datanodeOpts, nil
}

type InspectServiceList struct {
	ProductName string `db:"product_name"`
	ServiceName string `db:"service_name"`
}

func (l *deployInstanceList) GetServerListNotHadoopById(id int) ([]InspectServiceList, error) {
	result := make([]InspectServiceList, 0)
	query := fmt.Sprintf("select dpl.product_name,dil.service_name "+
		"from %s dil left join %s dpl on dil.pid=dpl.id "+
		"where dil.status = 'running' and dpl.product_name <> 'Hadoop' and dil.cluster_id = ? "+
		"group by dpl.product_name, dil.service_name", l.TableName, TBL_DEPLOY_PRODUCT_LIST)
	if err := USE_MYSQL_DB().Select(&result, query, id); err != nil {
		log.Errorf("GetServerListNotHadoopById query: %v, cluster_id %d, err: %v", query, id, err)
		return nil, err
	}
	return result, nil
}
