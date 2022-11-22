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

//黑夜给了我黑色的眼睛，专治各种疑难buf(g)；
package instance

import (
	"bytes"
	"database/sql"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/asset"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/event"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/schema"
)

const (
	STATUS_CHAN_BUFFER    = 16
	CMD_SEP               = " "
	INSTALL_PARAM_SEP     = ","
	INSTALL_CURRRENT_PATH = "/opt/dtstack/"
)

var InstallScriptTemplate = &template.Template{}

var PatchUpdateScriptTemplae = &template.Template{}

type instance struct {
	id          int64
	clusterId   int
	pid         int
	ip          string
	operationId string
	sid         string
	name        string
	schema      *schema.SchemaConfig
	agentId     string
	/*
		2021.01.22.em4.1.3
		instance mode
		0 install
		1 upgrade
	*/
	mode int

	statusCh chan event.Event
}

type Instancer interface {
	Install(onlyAgent bool) error
	UnInstall(onlyAgent bool) error
	UpdateConfig() error
	Start() error
	Stop(stopAgentOptionsTypeArr ...int) error
	ExecScript(script, timeout string, execId string) error
	SetPid(pid int) error
	SetMode(mode int)
	ID() int
	Clear()
	PatchUpdate(productname, servicename, path, downloadfile string) error

	GetInstanceInfo() (error, *model.DeployInstanceInfo)
	GetStatusChan() <-chan event.Event
}

func NewCommonInstancer(clusterid int, pid int, ip string, sid string, serviceName string, agentid string, operationId string) Instancer {
	newPatchUpdateInstance := &instance{
		clusterId:   clusterid,
		pid:         pid,
		ip:          ip,
		sid:         sid,
		name:        serviceName,
		agentId:     agentid,
		operationId: operationId,
	}
	return newPatchUpdateInstance
}

func NewInstancer(pid int, ip, serviceName string, schema *schema.SchemaConfig, operationId string) (Instancer, error) {
	if schema.Service[serviceName].Instance == nil || schema.Service[serviceName].Instance.UseCloud {
		err := fmt.Errorf("NewInstancer err: schema instance invalid")
		return nil, err
	}
	err, hostInfo := model.DeployHostList.GetHostInfoByIp(ip)
	if err != nil {
		return nil, err
	}
	hostRel, err := model.DeployClusterHostRel.GetClusterHostRelBySid(hostInfo.SidecarId)
	if err != nil {
		return nil, err
	}
	newInstance := &instance{
		clusterId:   hostRel.ClusterId,
		pid:         pid,
		operationId: operationId,
		ip:          ip,
		sid:         hostInfo.SidecarId,
		name:        serviceName,
		schema:      schema,
		statusCh:    make(chan event.Event, STATUS_CHAN_BUFFER),
	}
	schemaByte, err := json.Marshal(schema.Service[serviceName])
	if err != nil {
		return nil, err
	}
	healthState := model.INSTANCE_HEALTH_NOTSET
	prometheusPort, _ := strconv.Atoi(schema.Service[serviceName].Instance.PrometheusPort)
	if schema.Service[serviceName].Instance.HealthCheck != nil {
		healthState = model.INSTANCE_HEALTH_WAITING
	}
	err, id, agentId := model.DeployInstanceList.NewInstanceRecord(hostRel.ClusterId, pid, prometheusPort, healthState, ip, hostInfo.SidecarId,
		schema.Service[serviceName].Group, serviceName, schema.Service[serviceName].ServiceDisplay, schema.Service[serviceName].Version, schema.Service[serviceName].Instance.HARoleCmd, schemaByte)
	if err != nil {
		return nil, err
	}
	newInstance.id = id
	newInstance.agentId = agentId

	event.GetEventManager().AddObserver(newInstance)

	return newInstance, nil
}

func (this *instance) SetMode(mode int) {
	this.mode = mode
}

func (this *instance) handleError(typ, msg, id string, execId string) error {
	if execId != "" {
		execShellInfo, err := model.ExecShellList.GetByExecId(execId)
		if err != nil {
			log.Debugf("ExecShellList GetByExecId error: %v", err)
		}
		if err == nil && execShellInfo.ExecStatus == enums.ExecStatusType.Running.Code {
			now := time.Now()
			duration := now.Sub(execShellInfo.CreateTime.Time).Seconds()
			err = model.ExecShellList.UpdateStatusByExecId(execId, enums.ExecStatusType.Failed.Code, dbhelper.NullTime{Time: now, Valid: true}, sql.NullFloat64{Float64: duration, Valid: true})
			if err != nil {
				log.Errorf("ExecShellList UpdateStatusBySeq error: %v", err)
			}
		}

		operationInfo, err := model.OperationList.GetByOperationId(execShellInfo.OperationId)
		if err != nil {
			return err
		}

		now := time.Now()
		duration := now.Sub(operationInfo.CreateTime.Time).Seconds()
		model.OperationList.UpdateStatusByOperationId(operationInfo.OperationId, enums.ExecStatusType.Failed.Code, dbhelper.NullTime{Time: now, Valid: true}, sql.NullFloat64{Float64: duration, Valid: true})
	}

	this.updateStatus(typ, msg)
	log.Errorf("%v, id: %v", msg, id)
	return errors.New(msg)
}

func (this *instance) saveEvent(eventType string, event interface{}) {
	content, err := json.MarshalIndent(event, "", "\t")
	if err != nil {
		log.Errorf("%v", err)
	}
	model.DeployInstanceEvent.NewInstanceEvent(this.id, eventType, string(content[:]))
}

func (this *instance) PatchUpdate(productname, servicename, path string, downloadfile string) (ret error) {
	execId := uuid.NewV4().String()
	err := model.ExecShellList.InsertExecShellInfo(this.clusterId, this.operationId, execId, productname, this.name, this.sid, enums.ShellType.Exec.Code)
	if err != nil {
		log.Errorf("ExecShellList InsertExecShellInfo error: %v")
	}
	patchUpdateScript, err := this.getPatchUpdateScript(productname, servicename, path, downloadfile)
	if err != nil {
		msg := fmt.Sprintf("get patches update script err: %v", err)
		return this.handleError("patches update fail,script init fail", msg, this.sid, execId)
	}

	err = this.ExecScript(patchUpdateScript, EXEC_TIMEOUT, execId)
	if err != nil {
		msg := fmt.Sprintf("exec patches update script err: %v", err)
		return this.handleError("update fail", msg, this.sid, execId)
	}

	return nil
}

func (this *instance) Install(onlyAgent bool) (ret error) {
	installEvent := &agent.InstallEvent{
		InstanceEvent: agent.InstanceEvent{
			Type: model.INSTANCE_EVENT_INSTALL,
		},
	}
	defer func() {
		if ret != nil {
			installEvent.Message = ret.Error()
		} else {
			installEvent.Message = "install success"
		}
		this.saveEvent(installEvent.GetType(), installEvent)
	}()

	execId := uuid.NewV4().String()
	productInfo, err := model.DeployProductList.GetProductInfoById(this.pid)
	if err != nil {
		return err
	}
	err = model.ExecShellList.InsertExecShellInfo(this.clusterId, this.operationId, execId, productInfo.ProductName, this.name, this.sid, enums.ShellType.Install.Code)
	if err != nil {
		log.Errorf("ExecShellList InsertExecShellInfo err: %v", err)
	}

	if this.getInstanceSchema().Instance == nil || this.getInstanceSchema().Instance.UseCloud {
		return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, "schema instance invalid", this.sid, execId)
	}
	// 构建agent端部署参数(安装脚本、用户、应用目录、下载脚本、agentid等)
	param := &agent.InstallParms{
		Timeout: fmt.Sprintf("%dm", cache.SysConfig.GlobalConfig.ServiceInstallTimeoutLimit),
	}
	installEvent.InstallParam = param
	param.Name = this.name
	param.CollectorId = this.sid
	param.WorkDir = this.getInstanceHomeDir() + this.name
	param.RunUser = this.getInstanceSchema().Instance.RunUser
	if len(this.getInstanceSchema().Instance.ConfigPaths) > 0 {
		param.ConfigurationPath = this.getInstanceSchema().Instance.ConfigPaths[0]
	}
	cmds := strings.Split(this.getInstanceSchema().Instance.Cmd, CMD_SEP)
	if len(cmds) > 0 {
		param.BinaryPath = cmds[0]
	} else {
		return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, "schema cmd is null", this.sid, execId)
	}
	if len(cmds) > 1 {
		param.Parameter = strings.Join(cmds[1:], INSTALL_PARAM_SEP)
	}

	if this.getInstanceSchema().Instance.HealthCheck != nil {
		param.HealthShell = this.getInstanceSchema().Instance.HealthCheck.Shell
		param.HealthPeriod = this.getInstanceSchema().Instance.HealthCheck.Period
		param.HealthStartPeriod = this.getInstanceSchema().Instance.HealthCheck.StartPeriod
		param.HealthTimeout = this.getInstanceSchema().Instance.HealthCheck.Timeout
		param.HealthRetries = this.getInstanceSchema().Instance.HealthCheck.Retries
	}
	// 生成install_agentx.sh脚本下载安装包进行安装
	if !onlyAgent && !this.isInstanceEmptyCar() {
		param.InstallScript, err = this.getInstanceInstallScript(param.BinaryPath, param.RunUser, this.getInstanceSchema().Instance.DataDir)
		if err != nil {
			msg := fmt.Sprintf("get agent install script err: %v", err)
			return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
		}
	} else {
		log.Debugf("Install %v with onlyAgent flag", this.name)
		param.InstallScript = LINUX_EXEC_HEADER
	}

	// 执行组件包进行安装
	err, resp := agent.AgentClient.AgentInstall(param, execId)
	if err != nil {
		msg := fmt.Sprintf("exec agent install err: %v", err)
		return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
	}
	if resp.Data == nil {
		return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, "exec agent install resp is null", this.sid, execId)
	}
	agentId, ok := resp.Data.(map[string]interface{})["agent_id"]
	if !ok {
		return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, "agent_id is null", this.sid, execId)
	}
	err = this.updateAgentId(agentId.(string))
	if err != nil {
		msg := fmt.Sprintf("update agent_id fail, err: %v", err)
		return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
	}
	respResult, _ := resp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	installEvent.InstallResp = respResult
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["message"]
		if ok {
			msg := fmt.Sprintf("exec agent install err: %v", message.(string))
			return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
		}
		return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, "agent install unkown error", this.sid, execId)
	}

	if !this.isInstanceEmptyCar() && len(this.getInstanceSchema().Instance.ConfigPaths) > 0 {
		baseDir := filepath.Join(base.WebRoot, this.schema.ProductName, this.schema.ProductVersion)
		cfgContents, err := this.schema.ParseServiceConfigFiles(baseDir, this.name)
		if err != nil {
			msg := fmt.Sprintf("parse service config err: %v, service: %v", err, this.name)
			return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
		}
		for index, content := range cfgContents {
			configParam := &agent.ConfigParams{ConfigContent: string(content[:]), ConfigPath: this.getInstanceSchema().Instance.ConfigPaths[index],
				WorkDir: this.getInstanceHomeDir() + this.name}
			installEvent.ConfigParam = append(installEvent.ConfigParam, configParam.ConfigPath)
			err, resp = agent.AgentClient.AgentConfigUpdate(this.sid, this.agentId, configParam, execId)
			if err != nil {
				msg := fmt.Sprintf("exec update config err: %v, config path: %v", err, configParam.ConfigPath)
				return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
			}
			respResult, _ := resp.Data.(map[string]interface{})["result"]
			failed, ok := respResult.(map[string]interface{})["failed"]
			installEvent.ConfigResp = append(installEvent.ConfigResp, respResult)
			if ok && failed.(bool) == true {
				message, ok := respResult.(map[string]interface{})["message"]
				if ok {
					msg := fmt.Sprintf("exec update config err: %v, config path: %v", message.(string), configParam.ConfigPath)
					return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
				}
				msg := fmt.Sprintf("agent config update unkown error, config path:%v", configParam.ConfigPath)
				return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
			}
		}
	}
	if !onlyAgent && this.mode != 1 && this.mode != 3 &&
		!this.isInstanceEmptyCar() &&
		this.getInstanceSchema().Instance.PostDeploy != "" &&
		this.mode != 1 {
		execId := uuid.NewV4().String()
		err := model.ExecShellList.InsertExecShellInfo(this.clusterId, this.operationId, execId, productInfo.ProductName, this.name, this.sid, enums.ShellType.Exec.Code)
		if err != nil {
			log.Errorf("ExecShellList InsertExecShellInfo error: %v")
		}
		err = this.ExecScript(LINUX_EXEC_HEADER+this.getInstanceSchema().Instance.PostDeploy, EXEC_TIMEOUT, execId)
		if err != nil {
			msg := fmt.Sprintf("exec post deploy err: %v", err)
			return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
		}
		respResult, _ := resp.Data.(map[string]interface{})["result"]
		failed, ok := respResult.(map[string]interface{})["failed"]
		installEvent.PostDeployResp = respResult
		if ok && failed.(bool) == true {
			message, ok := respResult.(map[string]interface{})["response"]
			if ok {
				msg := fmt.Sprintf("exec post deploy err: %v", message)
				return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, msg, this.sid, execId)
			}
			return this.handleError(model.INSTANCE_STATUS_INSTALL_FAIL, "post deploy unkown error", this.sid, execId)
		}
	}
	//surpport upgrade
	if this.mode == 1 || this.mode == 3 &&
		!onlyAgent &&
		!this.isInstanceEmptyCar() &&
		this.getInstanceSchema().Instance.PostUpGrade != "" {
		execId := uuid.NewV4().String()
		err := model.ExecShellList.InsertExecShellInfo(this.clusterId, this.operationId, execId, productInfo.ProductName, this.name, this.sid, enums.ShellType.Exec.Code)
		if err != nil {
			log.Errorf("ExecShellList InsertExecShellInfo error: %v")
		}
		err = this.ExecScript(LINUX_EXEC_HEADER+this.getInstanceSchema().Instance.PostUpGrade, EXEC_TIMEOUT, execId)
		if err != nil {
			msg := fmt.Sprintf("exec post upgrade err: %v", err)
			return this.handleError(model.INSTANCE_STATUS_UPGRADE_FAIL, msg, this.sid, execId)
		}
		respResult, _ := resp.Data.(map[string]interface{})["result"]
		failed, ok := respResult.(map[string]interface{})["failed"]
		installEvent.PostDeployResp = respResult
		if ok && failed.(bool) == true {
			message, ok := respResult.(map[string]interface{})["response"]
			if ok {
				msg := fmt.Sprintf("exec post upgrade err: %v", message)
				return this.handleError(model.INSTANCE_STATUS_UPGRADE_FAIL, msg, this.sid, execId)
			}
			return this.handleError(model.INSTANCE_STATUS_UPGRADE_FAIL, "post upgrade unkown error", this.sid, execId)
		}
	}
	this.updateStatus(model.INSTANCE_STATUS_INSTALLED, "")

	return nil
}

func (this *instance) handleUnInstallScript() error {
	info, err := model.DeployProductList.GetProductInfoById(this.pid)
	if err != nil {
		return nil
	}
	psc := &schema.SchemaConfig{}
	err = json.Unmarshal(info.Product, psc)
	if err != nil {
		log.Debugf("pid %d json.Unmarshal %v", this.pid, err)
		return nil
	}
	if len(psc.Service[this.name].Instance.UnInstall) == 0 {
		return nil
	}
	newSchema, err := schema.Clone(psc)
	if err != nil {
		return err
	}
	// 单scheme
	var infoList []model.SchemaFieldModifyInfo
	query := "SELECT field_path, field FROM " + model.DeploySchemaFieldModify.TableName + " WHERE product_name=? AND cluster_id=? AND service_name = ? "
	if err := model.USE_MYSQL_DB().Select(&infoList, query, this.schema.ProductName, this.clusterId, this.name); err != nil {
		return err
	}
	for _, modify := range infoList {
		if _, err := newSchema.SetField(this.name+"."+modify.FieldPath, modify.Field); err != nil {
			return err
		}
	}

	// 多scheme
	multiFields, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndIp(this.clusterId, newSchema.ProductName, this.name, this.ip)
	if err != nil {
		log.Debugf("SchemaMultiField.GetByProductNameAndServiceNameAndIp pid %v err %v", this.pid, err)
		return err
	}

	for _, multiField := range multiFields {
		if _, err := newSchema.SetField(this.name+"."+multiField.FieldPath, multiField.Field); err != nil {
			return err
		}
	}
	serviceConfig := newSchema.Service[this.name]
	baseDir := filepath.Join(base.WebRoot, this.schema.ProductName, this.schema.ProductVersion, this.name)
	unInstallScript, err1 := serviceConfig.ParseUnInstallScript(baseDir)
	unInstallParameter, err2 := serviceConfig.ParseUnInstallParamter()
	if err1 != nil && err2 != nil {
		return fmt.Errorf("serviceConfig.ParseUnInstallScript %v, ParseUnInstallParamter %v", err1, err2)
	}
	if unInstallScript != "" {
		param := &agent.ExecScriptParams{
			ExecScript: unInstallScript,
			Parameter:  strings.Join(unInstallParameter, ","),
			Timeout:    "30m",
			AgentId:    this.agentId,
		}
		err, unInstallResp := agent.AgentClient.AgentExec(this.sid, param, "")
		if err != nil {
			return fmt.Errorf("AgentClient.AgentExec %v", err)
		}
		result, ok := unInstallResp.Data.(map[string]interface{})["result"]
		if !ok {
			return fmt.Errorf("unInstallResp.Data %v", unInstallResp)
		}
		failed, ok := result.(map[string]interface{})["failed"]
		if ok && failed.(bool) {
			return fmt.Errorf("unInstallResp failed %v", failed)
		}
	}
	return nil
}

func (this *instance) UnInstall(onlyAgent bool) (ret error) {
	unInstallEvent := &agent.UnInstallEvent{
		InstanceEvent: agent.InstanceEvent{
			Type: model.INSTANCE_EVENT_UNINSTALL,
		},
	}
	defer func() {
		if ret != nil {
			unInstallEvent.Message = ret.Error()
		} else {
			unInstallEvent.Message = "uninstall success"
		}
		this.saveEvent(unInstallEvent.GetType(), unInstallEvent)
	}()
	if this.agentId == "" {
		err := model.DeployInstanceList.DeleteByInstanceId(int(this.id))
		if err != nil {
			log.Errorf("DeleteByInstanceId err:%v", err.Error())
		}
		return nil
	}
	if this.getInstanceSchema().Instance == nil || this.getInstanceSchema().Instance.UseCloud {
		return this.handleError(model.INSTANCE_STATUS_UNINSTALL_FAIL, "schema instance invalid", this.agentId, "")
	}
	var err error
	this.updateStatus(model.INSTANCE_STATUS_UNINSTALLING, "")

	if err := this.handleUnInstallScript(); err != nil {
		log.Debugf("%s handle UnInstall Script %v", this.ip, err)
		return this.handleError(model.INSTANCE_STATUS_UNINSTALL_FAIL, "handle UnInstall Script", this.agentId, "")
	}
	param := &agent.ShellParams{}
	unInstallEvent.UnInstallParam = param
	if !onlyAgent && !this.isInstanceEmptyCar() {
		param.ShellScript = this.getUninstallScript()
	} else {
		log.Debugf("UnInstall %v with onlyAgent flag", this.name)
		param.ShellScript = LINUX_EXEC_HEADER
	}
	err, agentServerResp := agent.AgentClient.AgentUninstall(this.agentId, param, "")
	log.Infof("AgentClient.AgentUninstall Response: %s", agentServerResp)
	if err != nil {
		msg := fmt.Sprintf("exec agent uninstall err: %v", err.Error())
		return this.handleError(model.INSTANCE_STATUS_UNINSTALL_FAIL, msg, this.agentId, "")
	} else {
		respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
		failed, ok := respResult.(map[string]interface{})["failed"]
		unInstallEvent.UnInstallResp = respResult
		if ok && failed.(bool) == true {
			message, ok := respResult.(map[string]interface{})["message"]
			if !ok {
				message = ""
			}
			if !strings.Contains(message.(string), "can't find agent, uninstall fail") {
				return this.handleError(model.INSTANCE_STATUS_UNINSTALL_FAIL, message.(string), this.agentId, "")
			}
		}
		if !onlyAgent && !this.isInstanceEmptyCar() && this.getInstanceSchema().Instance.PostUndeploy != "" {
			err = this.ExecScriptGlobal(LINUX_EXEC_HEADER+this.getInstanceSchema().Instance.PostUndeploy, EXEC_TIMEOUT)
			if err != nil {
				msg := fmt.Sprintf("exec post undeploy err: %v", err.Error())
				return this.handleError(model.INSTANCE_STATUS_UNINSTALL_FAIL, msg, this.agentId, "")
			}
		}
		this.updateStatus(model.INSTANCE_STATUS_UNINSTALLED, "")
	}
	err = model.DeployInstanceList.DeleteByagentId(this.agentId)
	if err != nil && err != sql.ErrNoRows {
		msg := fmt.Sprintf("delete instance err: %v", err)
		log.Errorf("%v, id: %v", msg, this.agentId)
		return errors.New(msg)
	}
	return nil
}

func (this *instance) UpdateConfig() (ret error) {
	configEvent := &agent.ConfigEvent{
		InstanceEvent: agent.InstanceEvent{
			Type: model.INSTANCE_EVENT_CONFIG_UPDATE,
		},
	}
	defer func() {
		if ret != nil {
			configEvent.Message = ret.Error()
		} else {
			configEvent.Message = "config success"
		}
		this.saveEvent(configEvent.GetType(), configEvent)
	}()
	if this.agentId == "" {
		return fmt.Errorf("agent id is null")
	}
	if this.getInstanceSchema().Instance == nil || this.getInstanceSchema().Instance.UseCloud {
		return this.handleError(model.INSTANCE_STATUS_RUN_FAIL, "schema instance invalid", this.agentId, "")
	}

	if !this.isInstanceEmptyCar() && len(this.getInstanceSchema().Instance.ConfigPaths) > 0 {
		baseDir := filepath.Join(base.WebRoot, this.schema.ProductName, this.schema.ProductVersion)
		configEvent.ConfigSchema = this.schema
		cfgContents, err := this.schema.ParseServiceConfigFiles(baseDir, this.name)
		if err != nil {
			msg := fmt.Sprintf("parse service config err: %v, service: %v", err, this.name)
			return errors.New(msg)
		}
		for index, content := range cfgContents {
			configParam := &agent.ConfigParams{ConfigContent: string(content[:]), ConfigPath: this.getInstanceSchema().Instance.ConfigPaths[index],
				WorkDir: this.getInstanceHomeDir() + this.name}
			configEvent.ConfigPath = append(configEvent.ConfigPath, configParam.ConfigPath)
			err, resp := agent.AgentClient.AgentConfigUpdate(this.sid, this.agentId, configParam, this.operationId)
			if err != nil {
				msg := fmt.Sprintf("exec update config err: %v, config path: %v", err, configParam.ConfigPath)
				return errors.New(msg)
			}
			respResult, _ := resp.Data.(map[string]interface{})["result"]
			failed, ok := respResult.(map[string]interface{})["failed"]
			configEvent.ConfigResp = append(configEvent.ConfigResp, respResult)
			if ok && failed.(bool) == true {
				message, ok := respResult.(map[string]interface{})["message"]
				if ok {
					msg := fmt.Sprintf("exec update config err: %v, config path: %v", message.(string), configParam.ConfigPath)
					return errors.New(msg)
				}
				msg := fmt.Sprintf("agent config update unkown error, config path:%v", configParam.ConfigPath)
				return errors.New(msg)
			}
		}
	}
	return nil
}

func (this *instance) Start() (ret error) {
	startEvent := &agent.StartEvent{
		InstanceEvent: agent.InstanceEvent{
			Type: model.INSTANCE_EVENT_START,
		},
	}
	defer func() {
		if ret != nil {
			startEvent.Message = ret.Error()
		} else {
			startEvent.Message = "start success"
		}
		this.saveEvent(startEvent.GetType(), startEvent)
	}()
	err, info := model.DeployInstanceList.GetInstanceInfoByAgentId(this.agentId)
	if err != nil {
		log.Errorf("DeployInstanceList GetInstanceInfoByAgentId error: %v")

	}
	productInfo, err := model.DeployProductList.GetProductInfoById(info.Pid)
	if err != nil {
		log.Errorf("DeployProductList GetProductInfoById error: %v")
	}
	execId := uuid.NewV4().String()
	err = model.ExecShellList.InsertExecShellInfo(this.clusterId, this.operationId, execId, productInfo.ProductName, info.ServiceName, info.Sid, enums.ShellType.Start.Code)
	if err != nil {
		log.Errorf("ExecShellList InsertExecShellInfo error: %v")
	}
	if this.agentId == "" {
		return fmt.Errorf("agent id is null")
	}
	if this.getInstanceSchema().Instance == nil || this.getInstanceSchema().Instance.UseCloud {
		return this.handleError(model.INSTANCE_STATUS_RUN_FAIL, "schema instance invalid", this.agentId, execId)
	}

	env := map[string]string{}
	if this.getInstanceSchema().Instance.Environment != nil {
		for key, value := range this.getInstanceSchema().Instance.Environment {
			env[key] = *value
		}
	}
	param := &agent.StartParams{
		AgentId:     this.agentId,
		Environment: env,
	}
	startEvent.StartParam = param
	err, agentServerResp := agent.AgentClient.AgentStartWithParam(param, this.agentId, execId)
	if err != nil {
		msg := fmt.Sprintf("exec start err: %v", err.Error())
		return this.handleError(model.INSTANCE_STATUS_RUN_FAIL, msg, this.agentId, execId)
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	startEvent.StartResp = respResult
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["message"]
		msg := "unkown error"
		if ok {
			msg = fmt.Sprintf("exec start err: %v", message)
		}
		if !strings.Contains(msg, agent.IS_ALREADY_RUNNING) {
			log.Errorf("%v", msg)
			return this.handleError(model.INSTANCE_STATUS_RUN_FAIL, msg, this.agentId, execId)
		}
	}
	if this.getInstanceSchema().Instance.HealthCheck != nil {
		this.updateStatus(model.INSTANCE_STATUS_RUNNING, "", model.INSTANCE_HEALTH_WAITING)
	} else {
		this.updateStatus(model.INSTANCE_STATUS_RUNNING, "")
	}
	return nil
}

func (this *instance) Stop(stopAgentOptionsTypeArr ...int) (ret error) {
	stopEvent := &agent.StopEvent{
		InstanceEvent: agent.InstanceEvent{
			Type: model.INSTANCE_EVENT_STOP,
		},
	}
	defer func() {
		if ret != nil {
			stopEvent.Message = ret.Error()
		} else {
			stopEvent.Message = "stop success"
		}
		this.saveEvent(stopEvent.GetType(), stopEvent)
	}()
	if this.agentId == "" {
		return fmt.Errorf("[Instancer] Stop err: agent id is null")
	}
	if this.getInstanceSchema().Instance == nil || this.getInstanceSchema().Instance.UseCloud {
		return this.handleError(model.INSTANCE_STATUS_STOP_FAIL, "schema instance invalid", this.agentId, "")
	}

	stopAgentOptionsType := agent.AGENT_STOP_UNRECOVER
	if stopAgentOptionsTypeArr != nil {
		stopAgentOptionsType = stopAgentOptionsTypeArr[0]
	}

	err, agentServerResp := agent.AgentClient.AgentStop(this.agentId, stopAgentOptionsType, "")
	if err != nil {
		msg := fmt.Sprintf("exec stop err: %v", err.Error())
		return this.handleError(model.INSTANCE_STATUS_STOP_FAIL, msg, this.agentId, "")
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	stopEvent.StopResp = respResult
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["message"]
		if ok {
			msg := fmt.Sprintf("exec stop err: %v", message)
			return this.handleError(model.INSTANCE_STATUS_STOP_FAIL, msg, this.agentId, "")
		}
		msg := "unkown error"
		return this.handleError(model.INSTANCE_STATUS_STOP_FAIL, msg, this.agentId, "")
	}

	if this.getInstanceSchema().Instance.HealthCheck != nil {
		this.updateStatus(model.INSTANCE_STATUS_STOPPED, "", model.INSTANCE_HEALTH_WAITING)
	} else {
		this.updateStatus(model.INSTANCE_STATUS_STOPPED, "")
	}
	return nil
}

func (this *instance) ExecScript(script, timeout string, execId string) (ret error) {
	execEvent := &agent.ExecEvent{
		InstanceEvent: agent.InstanceEvent{
			Type: model.INSTANCE_EVENT_EXEC,
		},
	}
	defer func() {
		if ret != nil {
			execEvent.Message = ret.Error()
		} else {
			execEvent.Message = "stop success"
		}
		this.saveEvent(execEvent.GetType(), execEvent)
	}()
	if len(this.sid) == 0 {
		return fmt.Errorf("ExecScript err: sid is null")
	}
	param := &agent.ExecScriptParams{
		ExecScript: script,
		Parameter:  "",
		Timeout:    timeout,
		AgentId:    this.agentId,
	}
	execEvent.ExecScriptParam = param

	err, agentServerResp := agent.AgentClient.AgentExec(this.sid, param, execId)
	if err != nil {
		log.Errorf("ExecScript sid: %v, err: %v, resp: %v", this.sid, err, agentServerResp)
		return err
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	execEvent.ExecResp = respResult
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["response"]
		if ok && message != nil {
			log.Errorf("ExecScript err: %v", message.(string))
			return fmt.Errorf(message.(string))
		}
		log.Errorf("ExecScript error unknow")
		return fmt.Errorf("ExecScript error unknow")
	}
	return nil
}

func (this *instance) ExecScriptGlobal(script, timeout string) error {
	if len(this.sid) == 0 {
		return fmt.Errorf("ExecScript err: sid is null")
	}
	param := &agent.ExecScriptParams{
		ExecScript: script,
		Parameter:  "",
		Timeout:    timeout,
	}

	execId := uuid.NewV4().String()
	productInfo, err := model.DeployProductList.GetProductInfoById(this.pid)
	if err != nil {
		log.Errorf("DeployProductList GetProductInfoById error: %v")
	}
	err = model.ExecShellList.InsertExecShellInfo(this.clusterId, this.operationId, execId, productInfo.ProductName, this.name, this.sid, enums.ShellType.Exec.Code)
	if err != nil {
		log.Errorf("ExecShellList InsertExecShellInfo error: %v")
	}
	err, agentServerResp := agent.AgentClient.AgentExec(this.sid, param, execId)
	if err != nil {
		log.Errorf("ExecScript sid: %v, err: %v, resp: %v", this.sid, err, agentServerResp)
		return err
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["response"]
		if ok && message != nil {
			log.Errorf("ExecScript err: %v", message.(string))
			return fmt.Errorf(message.(string))
		}
		log.Errorf("ExecScript error unknow")
		return fmt.Errorf("ExecScript error unknow")
	}
	return nil
}

//update instance pid
func (this *instance) SetPid(pid int) error {
	var err error
	if err = model.DeployInstanceList.UpdateInstancePid(int(this.id), pid); err == nil {
		this.pid = pid
	}
	return err
}

func (this *instance) updateStatus(status, statusMsg string, healthState ...int) error {
	return model.DeployInstanceList.UpdateInstanceStatusById(int(this.id), status, statusMsg, healthState...)
}

func (this *instance) updateAgentId(agentId string) error {
	if len(agentId) == 0 {
		return fmt.Errorf("agent_id is null")
	}
	this.agentId = agentId
	return model.DeployInstanceList.UpdateInstanceAgentId(int(this.id), agentId)
}

func (this *instance) GetStatusChan() <-chan event.Event {
	return this.statusCh
}

func (this *instance) GetInstanceInfo() (error, *model.DeployInstanceInfo) {
	return model.DeployInstanceList.GetInstanceInfoById(int(this.id))
}

func (this *instance) ID() int {
	return int(this.id)
}

func (this *instance) Notify(ev event.Event) {
	if this.agentId != ev.AgentId {
		return
	}
	switch ev.Type {
	case event.REPORT_EVENT_HEALTH_CHECK, event.REPORT_EVENT_HEALTH_CHECK_CANCEL, event.REPORT_EVENT_INSTANCE_ERROR:
		select {
		case this.statusCh <- ev:
		default:
		}
	default:
		break
	}
}

func (this *instance) Clear() {
	log.Debugf("[Instancer] Clear")
	event.GetEventManager().RemoveObserver(this)
}

// 初始化生成补丁更新脚本模板到内存中
func (this *instance) initPatchUpdateScript() {
	asset.ResetPatchUpdateScriptWithLocalFile()

	typ := "templates/patchupdate.sh"

	// 返回上面脚本内容的数据流
	tmpScriptData, err := asset.Asset(typ)
	if err != nil {
		panic(err)
	}

	PatchUpdateScriptTemplae = template.Must(template.New(typ).Parse(string(tmpScriptData)))
}

func (this *instance) getPatchUpdateScript(productName string, serviceName string, path string, downloadfile string) (string, error) {
	this.initPatchUpdateScript()
	script := &bytes.Buffer{}

	// 补丁下载路径
	// 如：http://xxx/easyagent/DTBase/2.1.5/redis/patches_package/abc.jar
	uri := base.API_STATIC_URL + strings.TrimPrefix(downloadfile, ".")
	urlEncode, err := url.Parse(uri)
	if err != nil {
		log.Errorf("getPatchesDownloadUrl parse err: %v", err)
		return "", err
	}
	downloadUrl := urlEncode.String()

	// 补丁包更新目标，如:path="dtuic/tools/arthas/arthas-spy.jar" --> /opt/dtstack/DTUic/dtuic/tools/arthas/arthas-spy.jar
	patches_path := base.INSTALL_CURRRENT_PATH + productName + "/" + path

	//服务组件目录(用于备份原文件)，如：/opt/dtstack/DTBase/redis
	app_dir := base.INSTALL_CURRRENT_PATH + productName + "/" + serviceName

	//填充脚本中变量
	err = PatchUpdateScriptTemplae.Execute(script, map[string]interface{}{
		"PATCHES_PATH":       patches_path,
		"AGENT_DIR":          app_dir,
		"AGENT_DOWNLOAD_URL": downloadUrl,
	})
	if err != nil {
		return "", err
	}
	return script.String(), nil
}

// 初始化生成install_agent.sh模板
func (this *instance) initAagentInstallSh() {
	asset.ResetInstallAgentXShWithLocalFile()

	typ := "templates/install.agentx.sh"

	tmpldata, err := asset.Asset(typ)
	if err != nil {
		panic(err)
	}
	InstallScriptTemplate = template.Must(template.New(typ).Parse(string(tmpldata)))
}

func (this *instance) getInstanceInstallScript(appBin, runUser string, dataDir []string) (string, error) {
	installDir := this.getInstanceHomeDir()
	this.initAagentInstallSh()
	script := &bytes.Buffer{}

	downloadUrl, err := this.getInstanceDownloadUrl()
	if err != nil {
		return "", err
	}
	err = InstallScriptTemplate.Execute(script, map[string]interface{}{
		"AGENT_ZIP":          this.name + ".zip",
		"AGENT_DIR":          installDir + this.name,
		"AGENT_BIN":          installDir + this.name + "/" + appBin,
		"AGENT_DOWNLOAD_URL": downloadUrl,
		"RUN_USER":           runUser,
		"DATA_DIR":           strings.Join(dataDir, " "),
	})
	if err != nil {
		return "", err
	}
	return script.String(), nil
}

func (this *instance) getUninstallScript() string {
	var (
		path    = model.ClusterBackupConfig.GetPathByClusterId(this.clusterId)
		dstDir  = filepath.Join(path, this.schema.ProductName) + "/"
		dstName = this.getInstanceBackupName()
	)
	script := fmt.Sprintf(`%s 
mkdir -p %s 
mv "%s" "%s"
dirs=%s

declare -i count=0
declare -i keep=3

current=""
for dir in ${dirs}
do
    product=${dir%%%%-*}
    if [ "${current}" != "${product}" ]
    then
        current=${product}
        count=0
    fi

    if test ${count} -ge ${keep}
    then
		%s
    fi
    count=$((${count}+1))
done
`,
		LINUX_EXEC_HEADER,
		dstDir,
		this.getInstanceHomeDir()+this.name,
		dstName,
		//mysql_1610970097_4.0.2~
		"`ls -r "+dstDir+" | grep ^.*-.*-.*~$`",
		"rm -rf "+dstDir+"${dir}",
	)

	return script
}

func (this *instance) isInstanceEmptyCar() bool {

	if this.getInstanceSchema().Instance.EmptyCar {
		return true
	}
	return false
}

func (this *instance) getInstanceSchema() schema.ServiceConfig {
	return this.schema.Service[this.name]
}

func (this *instance) getInstanceHomeDir() string {
	installPath := this.getInstanceSchema().Instance.InstallPath
	if len(installPath) > 0 {
		return installPath
	}
	return base.INSTALL_CURRRENT_PATH + this.schema.ProductName + "/"
}

const backupSeparator = "-"

func (this *instance) getInstanceBackupName() string {
	// mysql-1610970097-4.0.2~
	configPath := model.ClusterBackupConfig.GetPathByClusterId(this.clusterId)
	return path.Join(configPath, this.schema.ProductName) + "/" + this.name + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + this.schema.ProductVersion + "~"
}

func (this *instance) getInstanceDownloadUrl() (string, error) {
	uri := base.API_STATIC_URL + "/easyagent/" + this.schema.ProductName + "/" + this.schema.ProductVersion + "/" + this.name + ".zip"
	urlEncode, err := url.Parse(uri)
	if err != nil {
		log.Errorf("getInstanceDownloadUrl parse err: %v", err)
		return uri, err
	}
	return urlEncode.String(), nil
}

func (this *instance) parseAgentServerResponse(resp *agent.EasyagentServerResponse, method, status, field string) error {
	respResult, ok := resp.Data.(map[string]interface{})["result"]
	if !ok {
		err := fmt.Errorf("%s err: %v", method, "response data.result is null")
		log.Errorf(err.Error())
		this.updateStatus(status, err.Error())
		return err
	}
	failed, ok := respResult.(map[string]interface{})["failed"]
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})[field]
		if ok {
			err := fmt.Errorf("%s err: %v", method, message.(string))
			log.Errorf(err.Error())
			this.updateStatus(status, err.Error())
			return fmt.Errorf(message.(string))
		}
	}
	return nil
}
