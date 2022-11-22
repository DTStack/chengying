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

package impl

import (
	"database/sql"
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/handler"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/matrix/instance"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/service"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"dtstack.com/dtstack/easymatrix/schema"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
	"io"
	"strconv"
	"strings"
	"time"
)

type switchDetail struct {
	SwitchType    string `db:"switch_type" json:"switch_type"`
	Status        string `db:"status" json:"status"`
	StatusMessage string `db:"status_message" json:"status_message"`
	Progress      int    `db:"progress" json:"progress"`
}

func GetSwitchDetail(ctx context.Context) apibase.Result {

	recordId, err := ctx.URLParamInt("record_id")
	if err != nil {
		log.Errorf("[SwitchDetail] Record id: %v is illegal, err: %v", recordId, err)
		return err
	}
	detail, err := model.SwitchRecord.GetRecordById(recordId)
	if err != nil {
		log.Errorf("[SwitchDetail] Query database error: %v", err)
	}
	return switchDetail{
		SwitchType:    detail.Type,
		Status:        detail.Status,
		StatusMessage: detail.StatusMessage,
		Progress:      detail.Progress,
	}
}

func CheckSwitchRecord(ctx context.Context) apibase.Result {
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	productName := ctx.Params().Get("product_name")
	serviceName := ctx.Params().Get("service_name")
	switchName := ctx.URLParam("name")
	currentRecord, err := model.SwitchRecord.GetCurrentSwitchRecord(clusterId, productName, serviceName, switchName)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("%v", err)
		return err
	}
	if currentRecord == nil {
		return map[string]interface{}{"record_id": 0}
	}
	return map[string]interface{}{"record_id": currentRecord.Id}
}

func OperateSwitch(ctx context.Context) (rlt apibase.Result) {
	info := model.SwitchRecordInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	switchName := ctx.FormValue("name")
	// type: on开启，off关闭
	switchType := ctx.FormValue("type")
	productVersion := ctx.FormValue("product_version")
	productInfo, err := model.DeployProductList.GetByProductNameAndVersion(info.ProductName, productVersion)
	if err != nil {
		log.Errorf("[Operate Switch] Get product(product_name:%v, product_version:%v) info error: %v", info.ProductName, productVersion, err)
		return err
	}
	// change value of instance.switch.kerberos.is_on and config.kerberos_on
	product := productInfo.Schema
	sc, err := schema.Unmarshal(product)
	if err != nil {
		log.Errorf("[Operate Switch] Unmarshal schema error: %v", err)
		return err
	}
	switcher, err := NewSwitcher(clusterId, info.ProductName, info.ServiceName, switchName, switchType, sc)
	if err != nil {
		log.Errorf("[Switcher] New switcher error: %v", err)
	}
	err = switcher.NewSwitcherRecord()
	if err != nil {
		log.Errorf("[Switcher] New switcher record error: %v", err)
		err := switcher.AddProgressLog("FAIL", err.Error(), 0)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
	}
	switcher.AddProgressLog("RUNNING", currentTime()+" 开关操作初始化成功", 10)
	go doSwitch(switcher, switchType)
	return map[string]interface{}{"record_id": switcher.GetId()}
}

func doSwitch(switcher Switcher, switchType string) (rlt interface{}) {
	// step1: modify field
	rlt = switcher.ModifyConfigField()
	if _, ok := rlt.(error); ok {
		log.Errorf("[Switcher] Modify field error: %v", rlt)
		switcher.AddProgressLog("FAIL", currentTime()+" 更改数据库field字段失败", 10)
		return rlt
	}
	switcher.AddProgressLog("RUNNING", currentTime()+" 更改数据库field字段成功", 30)
	// step2: update config
	rlt = switcher.UpdateConfig()
	if rlt != nil {
		log.Errorf("[Switcher] Update config error: %v", rlt)
		switcher.AddProgressLog("FAIL", currentTime()+fmt.Sprintf(" 配置文件下发失败: %v", rlt), 30)
		return rlt
	}
	switcher.AddProgressLog("RUNNING", currentTime()+" 配置文件下发成功", 50)
	// step3: execute on/off operation
	if switchType == base.SwitchOn {
		rlt = switcher.SwitchOn()
	} else {
		rlt = switcher.SwitchOff()
	}
	if rlt != nil {
		log.Errorf("[Switcher] Execute script error: %v", rlt)
		switcher.AddProgressLog("FAIL", currentTime()+fmt.Sprintf(" 执行脚本失败: %v", rlt), 50)
		return rlt
	}
	switcher.ModifyInstanceField()
	switcher.AddProgressLog("SUCCESS", currentTime()+" 开关操作成功", 100)
	return nil
}

func DoExtentionOperation(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	serviceName := ctx.Params().Get("service_name")
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	operationType := ctx.URLParam("type")
	value := ctx.URLParam("value")
	if operationType == "download" {
		productInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
		if err != nil {
			log.Errorf("[Download] Get product by name error: %v", err)
			return nil
		}
		list, err := model.DeployInstanceList.GetInstanceListByPidServiceName(productInfo.ID, clusterId, serviceName)
		if err != nil {
			return nil
		}
		if len(list) == 0 {
			log.Errorf("[Download] No instance in service: %s", serviceName)
			return fmt.Errorf("服务%s不存在实例", serviceName)
		}
		info := list[0]
		if !strings.HasPrefix(value, LINUX_SYSTEM_SLASH) {
			pwdCmd := "#!/bin/sh\n echo `pwd`"
			content, err := agent.AgentClient.ToExecCmd(info.Sid, info.AgentId, pwdCmd, "")
			if err != nil {
				return err
			}
			log.Debugf("pwd response: %v", content)
			value = strings.Replace(content, LINUX_SYSTEM_LINES, "", -1) + LINUX_SYSTEM_SLASH + value
			log.Debugf("download file: %v", value)
		}
		ip := strings.Split(info.Ip, "/")[0]
		fileSlice := strings.Split(value, LINUX_SYSTEM_SLASH)
		target := ip + EASYFILER_PORT
		var data = make(chan []byte, 1)
		ctx.Header("Content-Disposition", "attachment;filename="+fileSlice[len(fileSlice)-1]+TAR_SUFFIX)
		cancel := make(chan string, 1)
		go func() {
			if err := handler.DownloadWithoutStorage(target, value, data, cancel); err == io.EOF {
				log.Infof("Download file %v succeed.", value)
				return
			}
			log.Errorf("Download file %v fail", value)
			return
		}()
	LOOP:
		for {
			select {
			case ca := <-cancel:
				log.Errorf("cancel error: %v", ca)
				return errors.New(ca)
			case res := <-data:
				if string(res) == "done" {
					break LOOP
				}
				if string(res) == "error" {
					ctx.Write([]byte("download again?"))
					break LOOP
				}
				ctx.Write(res)
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
	return apibase.EmptyResult{}
}

type Switcher interface {
	NewSwitcherRecord() error
	GetId() int
	GetSwitcherRecord() (*model.SwitchRecordInfo, error)
	ModifyInstanceField() interface{}
	ModifyConfigField() interface{}
	UpdateConfig() error
	SwitchOn() error
	SwitchOff() error
	InstanceSwitchPath() string
	ConfigSwitchPath(configName string) string
	AddProgressLog(status, message string, progress int) error
}

type switcher struct {
	id          int
	clusterId   int
	productName string
	serviceName string
	switchName  string
	switchType  string
	ip          []string
	schema      *schema.SchemaConfig
}

func NewSwitcher(clusterId int, productName, serviceName, switchName, switchType string, schema *schema.SchemaConfig) (Switcher, error) {
	productInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	if err != nil {
		log.Errorf("[Switcher] Get product by name error: %v", err)
		return nil, err
	}
	list, err := model.DeployInstanceList.GetInstanceListByPidServiceName(productInfo.ID, clusterId, serviceName)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("[Switcher] Instance not exist in service " + serviceName)
	}
	switcher := &switcher{
		clusterId:   clusterId,
		productName: productName,
		serviceName: serviceName,
		switchName:  switchName,
		switchType:  switchType,
		schema:      schema,
		ip:          []string{},
	}
	for _, instanceInfo := range list {
		switcher.ip = append(switcher.ip, instanceInfo.Ip)
	}
	return switcher, nil
}

func (s *switcher) NewSwitcherRecord() error {
	recordId, err := model.SwitchRecord.NewSwitchRecord(s.switchName, s.switchType, s.productName, s.serviceName, "RUNNING", currentTime()+" 开始运行", s.clusterId, 0)
	if err != nil {
		log.Errorf("[Switcher] New switch record error: %v", err)
		return err
	}
	s.id = int(recordId)
	return nil
}

func (s *switcher) GetSwitcherRecord() (*model.SwitchRecordInfo, error) {
	switcherRecord, err := model.SwitchRecord.GetRecordById(s.id)
	if err != nil {
		log.Errorf("[Switcher] Get switcher record by id error: %v", s.id)
		return nil, err
	}
	return switcherRecord, nil
}

func (s *switcher) GetId() int {
	return s.id
}

func (s *switcher) ModifyInstanceField() (rlt interface{}) {
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	switches := s.schema.Service[s.serviceName].Instance.Switch
	if _, exist := switches[s.switchName]; exist {
		// modify instance switch field
		rlt = saveSchemaModifyField(s.clusterId, s.productName, s.serviceName, s.InstanceSwitchPath(), strconv.FormatBool(s.switchType == base.SwitchOn), tx)
		if rlt != nil {
			return errors.New(rlt.(string))
		}
	} else {
		log.Errorf("[Switcher] Switch %v not exists in schema", s.switchName)
		return errors.New("switch not exist in schema")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return rlt
}

func (s *switcher) ModifyConfigField() (rlt interface{}) {
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		defer func() {
			if _, ok := rlt.(error); ok {
				tx.Rollback()
			}
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
	}()
	switches := s.schema.Service[s.serviceName].Instance.Switch
	if sn, exist := switches[s.switchName]; exist {
		// modify config switch field
		rlt = saveSchemaModifyField(s.clusterId, s.productName, s.serviceName, s.ConfigSwitchPath(sn.Config), decideConfig(s.switchType), tx)
		if rlt != nil {
			return errors.New(rlt.(string))
		}
	} else {
		log.Errorf("[Switcher] Switch %v not exists in schema", s.switchName)
		return errors.New("switch not exist in schema")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return rlt
}

func (s *switcher) UpdateConfig() error {
	productInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(s.productName, s.clusterId)
	if err != nil {
		log.Errorf("[Switcher] Get product by name error: %v", err)
		return err
	}
	servicer, err := service.NewServicer(productInfo.ID, s.clusterId, s.serviceName, "")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	return servicer.RollingConfigUpdate(s.clusterId)
}

func (s *switcher) parseVariable() (*schema.SchemaConfig, error) {
	var err error
	// select one sidecar random to execute script
	var newSchema *schema.SchemaConfig
	if newSchema, err = schema.Clone(s.schema); err != nil {
		log.Errorf("[Switcher] Clone schema error")
		return nil, err
	}
	if err = inheritBaseService(s.clusterId, newSchema, model.USE_MYSQL_DB()); err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	if err = setSchemaFieldServiceAddr(s.clusterId, newSchema, model.USE_MYSQL_DB(), ""); err != nil {
		log.Debugf("[Product->ServiceGroup] setSchemaFieldServiceAddr err: %v", err)
		return nil, err
	}
	var node *model.ServiceIpNode
	node, err = model.GetServiceIpNode(s.clusterId, newSchema.ProductName, s.serviceName, s.ip[0])
	if err != nil {
		log.Errorf("%v", err)
	}
	idToIndex, err := getIdToIndex(s.clusterId, s.productName, s.serviceName, s.ip)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	if err = newSchema.SetServiceNodeIP(s.serviceName, util.FoundIpIdx(s.ip, s.ip[0]), node.NodeId, idToIndex); err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	//兼容存在uncheck部署模式
	if err = newSchema.SetEmptyServiceAddr(); err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	if err := newSchema.ParseServiceVariable(s.serviceName); err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return newSchema, nil
}

func (s *switcher) SwitchOn() error {
	newSchema, err := s.parseVariable()
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	if switchConfig, ok := newSchema.Service[s.serviceName].Instance.Switch[s.switchName]; ok {
		//生成 operationid 并且落库
		operationId := uuid.NewV4().String()
		err = model.OperationList.Insert(model.OperationInfo{
			ClusterId:       s.clusterId,
			OperationId:     operationId,
			OperationType:   enums.OperationType.OpenKerberos.Code,
			OperationStatus: enums.ExecStatusType.Running.Code,
			ObjectType:      enums.OperationObjType.Svc.Code,
			ObjectValue:     s.serviceName,
		})
		if err != nil {
			log.Errorf("OperationList Insert err:%v", err)
		}
		// execute on script
		if switchConfig.OnScript != "" {
			err, hostInfo := model.DeployHostList.GetHostInfoByIp(s.ip[0])
			if err != nil {
				return err
			}
			productInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(s.productName, s.clusterId)
			if err != nil {
				log.Errorf("[Switcher] Get product by name error: %v", err)
				return err
			}

			_, instanceInfo := model.DeployInstanceList.GetInstanceInfoByWhere(dbhelper.MakeWhereCause().Equal("cluster_id", s.clusterId).And().Equal("pid", productInfo.ID).And().Equal("service_name", s.serviceName).And().Equal("sid", hostInfo.SidecarId))
			_, err = s.ExecScript(hostInfo.SidecarId, instanceInfo.AgentId, instance.LINUX_EXEC_HEADER+switchConfig.OnScript, instance.EXEC_TIMEOUT, operationId)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
		}
		s.AddProgressLog("RUNNING", currentTime()+" 脚本执行成功", 80)
		// execute post on script
		if switchConfig.PostOnScript != nil {
			if switchConfig.PostOnScript.Type == "restart" {
				value := switchConfig.PostOnScript.Value
				if value != "" {
					serviceList := strings.Split(value, ",")
					productInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(s.productName, s.clusterId)
					if err != nil {
						log.Errorf("[Switcher] Get product by name error: %v", err)
						return err
					}
					for _, serviceName := range serviceList {
						servicer, err := service.NewServicer(productInfo.ID, s.clusterId, serviceName, operationId)
						if err != nil {
							log.Errorf("%v", err)
							return err
						}
						servicer.RollingRestart()
					}
				}
				s.AddProgressLog("SUCCESS", currentTime()+" 服务重启完成", 100)
			}
		}
	}
	return nil
}

func (s *switcher) SwitchOff() error {
	newSchema, err := s.parseVariable()
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	if switchConfig, ok := newSchema.Service[s.serviceName].Instance.Switch[s.switchName]; ok {
		//生成 operationid 并且落库
		operationId := uuid.NewV4().String()
		err = model.OperationList.Insert(model.OperationInfo{
			ClusterId:       s.clusterId,
			OperationId:     operationId,
			OperationType:   enums.OperationType.CloseKerberos.Code,
			OperationStatus: enums.ExecStatusType.Running.Code,
			ObjectType:      enums.OperationObjType.Svc.Code,
			ObjectValue:     s.serviceName,
		})
		if err != nil {
			log.Errorf("OperationList Insert err:%v", err)
		}
		// execute off script
		if switchConfig.OffScript != "" {
			err, hostInfo := model.DeployHostList.GetHostInfoByIp(s.ip[0])
			if err != nil {
				return err
			}
			productInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(s.productName, s.clusterId)
			if err != nil {
				log.Errorf("[Switcher] Get product by name error: %v", err)
				return err
			}

			_, instanceInfo := model.DeployInstanceList.GetInstanceInfoByWhere(dbhelper.MakeWhereCause().Equal("cluster_id", s.clusterId).And().Equal("pid", productInfo.ID).And().Equal("service_name", s.serviceName).And().Equal("sid", hostInfo.SidecarId))
			_, err = s.ExecScript(hostInfo.SidecarId, instanceInfo.AgentId, instance.LINUX_EXEC_HEADER+switchConfig.OffScript, instance.EXEC_TIMEOUT, operationId)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
		}
		s.AddProgressLog("RUNNING", currentTime()+" 脚本执行成功", 80)

		// execute post off script
		if switchConfig.PostOnScript != nil {
			if switchConfig.PostOffScript.Type == "restart" {
				value := switchConfig.PostOffScript.Value
				if value != "" {
					serviceList := strings.Split(value, ",")
					productInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(s.productName, s.clusterId)
					if err != nil {
						log.Errorf("[Switcher] Get product by name error: %v", err)
						return err
					}
					for _, serviceName := range serviceList {
						servicer, err := service.NewServicer(productInfo.ID, s.clusterId, serviceName, operationId)
						if err != nil {
							log.Errorf("%v", err)
							return err
						}
						servicer.RollingRestart()
					}
				}
				s.AddProgressLog("SUCCESS", currentTime()+" 服务重启完成", 100)
			}
		}
	}
	return nil
}

func (s *switcher) AddProgressLog(status, message string, progress int) error {
	record, err := model.SwitchRecord.GetRecordById(s.id)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	message = record.StatusMessage + "\n" + message
	_, err = model.SwitchRecord.NewSwitchRecord(s.switchName, s.switchType, s.productName, s.serviceName, status, message, s.clusterId, progress)
	if err != nil {
		return err
	}
	return nil
}

func (s *switcher) ExecScript(sid, agentId, script, timeout string, operationId string) (string, error) {
	execEvent := &agent.ExecEvent{
		InstanceEvent: agent.InstanceEvent{
			Type: model.INSTANCE_EVENT_EXEC,
		},
	}
	if len(sid) == 0 {
		return "", fmt.Errorf("ExecScript err: sid is null")
	}
	param := &agent.ExecScriptParams{
		ExecScript: script,
		Parameter:  "",
		Timeout:    timeout,
		AgentId:    agentId,
	}
	execEvent.ExecScriptParam = param

	execId := uuid.NewV4().String()
	err := model.ExecShellList.InsertExecShellInfo(s.clusterId, operationId, execId, s.productName, s.serviceName, sid, enums.ShellType.Exec.Code)
	if err != nil {
		log.Errorf("ExecShellList InsertExecShellInfo err:%v", err)
	}
	err, agentServerResp := agent.AgentClient.AgentExec(sid, param, execId)
	if err != nil {
		log.Errorf("ExecScript sid: %v, err: %v, resp: %v", sid, err, agentServerResp)
		return "", err
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	execEvent.ExecResp = respResult
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["response"]
		if ok && message != nil {
			log.Errorf("%v", message.(string))
			return "", fmt.Errorf(message.(string))
		}
		log.Errorf("ExecScript error unknow")
		return "", fmt.Errorf("ExecScript error unknow")
	}
	return "", nil
}

func (s *switcher) InstanceSwitchPath() string {
	return "Instance.Switch." + s.switchName + ".IsOn"
}

func (s *switcher) ConfigSwitchPath(configName string) string {
	return "Config." + configName + ".Value"
}

func saveSchemaModifyField(clusterId int, productName, serviceName, fieldPath, field string, tx *sqlx.Tx) interface{} {
	schemaModifyField := model.SchemaFieldModifyInfo{
		ClutserId:   clusterId,
		ProductName: productName,
		ServiceName: serviceName,
		FieldPath:   fieldPath,
		Field:       field,
	}
	return modifyField(&schemaModifyField, tx)
}

func decideConfig(switchType string) string {
	if switchType == base.SwitchOn {
		return "1"
	}
	return "0"
}

func currentTime() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}
