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
	"dtstack.com/dtstack/easymatrix/matrix/grafana"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/handler"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/discover"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/matrix/event"
	"dtstack.com/dtstack/easymatrix/matrix/harole"
	"dtstack.com/dtstack/easymatrix/matrix/instance"
	kdeploy "dtstack.com/dtstack/easymatrix/matrix/k8s/deploy"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	kmodel "dtstack.com/dtstack/easymatrix/matrix/k8s/model"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"dtstack.com/dtstack/easymatrix/schema"
	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
)

const (
	LINUX_SYSTEM_SLASH  = "/"
	LINUX_SYSTEM_LINES  = "\n"
	LINUX_FILE_TYPE_TXT = "txt"
	LINUX_FILE_TYPE_ZIP = "zip"
	EASYFILER_PORT      = ":7788"
	EASYFILER_TMP_ROOT  = "/tmp/"
)

const (
	ERROR_NO_ROW_IN_RESULT = "no rows in result set"
)
const (
	TIME_LAYOUT     = "2006-01-02 15:04:05"
	DATE_LAYOUT     = "2006-01-02"
	HOURS_PER_DAY   = 24
	MINUTES_PER_DAY = 86400
)

var (
	LINUX_SYSTEM_ZIP_EXTENTION = []string{".zip", ".tgz", ".tar.gz", ".tar"}
	LOG_MORE_ACTION_LATEST     = "latest"
	LOG_MORE_ACTION_DOWN       = "down"
	LOG_MORE_ACTION_UP         = "up"
	LOG_MORE_PREVIEW_COUNT     = 300
	LOG_MORE_PREVIEW_RATE      = 20
)

const (
	TAR_SUFFIX = ""
)

var (
	sMu        sync.RWMutex
	sessionMap = map[string]*preSession{}
)

type preSession struct {
	id            string
	instanceId    string
	log           string
	logWcl        int64
	anchor        int64
	lastAction    string
	startTimeNano int64
	start         int64
	end           int64
	total         int64
}

func checkInstanceStatus(info *model.DeployInstanceInfo, statusList []string) {
	paramErrs := apibase.NewApiParameterErrors()

	for _, status := range statusList {
		if status == info.Status {
			paramErrs.AppendError("checkInstanceStatus", "instance status failed %s", status)
		}
	}
	paramErrs.CheckAndThrowApiParameterErrors()
}

func checkInstanceAgentId(ctx context.Context) (string, *model.DeployInstanceInfo) {
	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	if agentIdStr == "" {
		paramErrs.AppendError("checkInstanceAgentId", "instance agent id is empty")
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	_, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("checkInstanceAgentId", "instance agent_id not uuid format")
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	err, info := model.DeployInstanceList.GetInstanceInfoByAgentId(agentIdStr)

	if err != nil {
		paramErrs.AppendError("checkInstanceAgentId", err.Error())
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	return agentIdStr, info
}

func Start(ctx context.Context) apibase.Result {
	agentId, info := checkInstanceAgentId(ctx)

	checkInstanceStatus(info, model.OUT_OF_START_STATUS_LIST)

	sc := &schema.ServiceConfig{}
	err := json.Unmarshal(info.Schema, sc)
	if err != nil {
		log.Errorf("[Start] json.Unmarshal service schema error: %v, agentId: %v", err, agentId)
		return fmt.Errorf("[Start] json.Unmarshal service schema error: %v", err)
	}

	env := map[string]string{}

	if sc.Instance != nil && sc.Instance.Environment != nil {
		for key, value := range sc.Instance.Environment {
			env[key] = *value
		}
	}
	param := &agent.StartParams{
		AgentId:     agentId,
		Environment: env,
	}
	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	err, info = model.DeployInstanceList.GetInstanceInfoByAgentId(agentId)
	if err != nil {
		return err
	}
	err = model.OperationList.Insert(model.OperationInfo{
		ClusterId:       info.ClusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.SvcStart.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Svc.Code,
		ObjectValue:     info.ServiceName,
	})
	if err != nil {
		log.Errorf("OperationList Insert err:%v", err)
	}
	productInfo, err := model.DeployProductList.GetProductInfoById(info.Pid)
	if err != nil {
		log.Errorf("DeployProductList GetProductInfoById err:%v", err)
	}
	execId := uuid.NewV4().String()
	err = model.ExecShellList.InsertExecShellInfo(info.ClusterId, operationId, execId, productInfo.ProductName, info.ServiceName, info.Sid, enums.ShellType.Start.Code)
	if err != nil {
		log.Errorf("ExecShellList InsertExecShellInfo err:%v", err)
	}
	err, agentServerResp := agent.AgentClient.AgentStartWithParam(param, param.AgentId, execId)
	if err != nil {
		log.Errorf("[Start] AgentStartWithParam err: %v, agentId: %v, resp: %v", err, agentId, agentServerResp)
		model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_RUN_FAIL, err.Error())
		return err
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["message"]
		msg := "start unkown error"
		if ok {
			msg = fmt.Sprintf("start err: %v", message.(string))
		}
		if !strings.Contains(msg, agent.IS_ALREADY_RUNNING) {
			log.Errorf("%v", msg)
			model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_RUN_FAIL, msg)
			return errors.New(msg)
		}
	}
	if info.HealthState != model.INSTANCE_HEALTH_NOTSET {
		err = model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_RUNNING, "", model.INSTANCE_HEALTH_WAITING)
	} else {
		err = model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_RUNNING, "")
	}
	if err != nil {
		log.Errorf("[Start] UpdateInstanceStatusByAgentId err: %v", err)
	}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(info.ClusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		productInfo, err := model.DeployProductList.GetProductInfoById(info.Pid)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "集群运维", "服务启动", "集群名称："+clusterInfo.Name+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion+
			", 服务组："+info.Group+", 服务名称："+info.ServiceName+info.ServiceVersion+", 服务实例："+info.Ip); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}

		// 检查重启服务列表中是否存在，若存在则删除掉对应ip的记录
		if err := model.NotifyEvent.DeleteNotifyEvent(info.ClusterId, 0, productInfo.ProductName, info.ServiceName, info.Ip, false); err != nil {
			log.Errorf("deleted notify event error: %v", err)
		}

	}()
	return err
}

func Stop(ctx context.Context) apibase.Result {
	agentId, info := checkInstanceAgentId(ctx)

	checkInstanceStatus(info, model.OUT_OF_STOP_STATUS_LIST)

	err, agentServerResp := agent.AgentClient.AgentStop(agentId, agent.AGENT_STOP_UNRECOVER, "")

	if err != nil {
		log.Errorf("[Stop] AgentStop agentId: %v, err: %v, resp: %v", agentId, err, agentServerResp)
		model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_STOP_FAIL, err.Error())
		return err
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	if ok && failed.(bool) == true {
		message, ok := respResult.(map[string]interface{})["message"]
		msg := "stop unkown error"
		if ok {
			msg = fmt.Sprintf("stop err: %v", message.(string))
		}
		log.Errorf("%v", msg)
		model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_STOP_FAIL, msg)
		return errors.New(msg)
	}
	if info.HealthState != model.INSTANCE_HEALTH_NOTSET {
		err = model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_STOPPED, "", model.INSTANCE_HEALTH_WAITING)
	} else {
		err = model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_STOPPED, "")
	}
	if err != nil {
		log.Errorf("[Stop] UpdateInstanceStatus err: %v, agentId: %v", err, agentId)
	}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(info.ClusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		productInfo, err := model.DeployProductList.GetProductInfoById(info.Pid)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "集群运维", "服务停止", "集群名称："+clusterInfo.Name+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion+
			", 服务组："+info.Group+", 服务名称："+info.ServiceName+info.ServiceVersion+", 服务实例："+info.Ip); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	return err
}

func ServiceList(ctx context.Context) apibase.Result {
	log.Debugf("[Instance->ServiceList] get service list by product name for EasyMatrix API ")
	productName := ctx.Params().Get("product_name")

	if productName == "" {
		log.Errorf("[Instance->ServiceList]product name is null")
		return fmt.Errorf("product name is null")
	}

	serviceList := []string{}
	query := "SELECT DISTINCT IL.service_name FROM " +
		model.DeployInstanceList.TableName + " AS IL LEFT JOIN " + model.DeployProductList.TableName + " AS PL ON IL.pid = PL.id WHERE PL.product_name=? ORDER BY service_name"

	log.Debugf("%v", query)
	if err := model.USE_MYSQL_DB().Select(&serviceList, query, productName); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return map[string]interface{}{
		"services": serviceList,
		"count":    len(serviceList),
	}
}

func GroupList(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		log.Errorf("product name is null")
		return fmt.Errorf("product name is null")
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	cloud := ctx.URLParam("cloud")
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)

	type serviceInfo struct {
		ServiceName string `db:"service_name"`
		Group       string `db:"group"`
		HealthState int    `db:"health_state"`
		Status      string `db:"status"`
	}
	type resultList struct {
		ServiceName        string `json:"service_name"`
		ServiceNameDisplay string `json:"service_name_display"`
		Alert              bool   `json:"alert"`
	}

	groupAndServices := map[string][]resultList{}
	serviceInfoList := []serviceInfo{}

	// Avoid deploying the same product package with multiple namespaces
	query := "SELECT IL.service_name, IL.group, IL.health_state, IL.status FROM " +
		model.DeployInstanceList.TableName + " AS IL LEFT JOIN " + model.DeployProductList.TableName + " AS PL ON IL.pid = PL.id WHERE PL.product_name=? AND IL.cluster_id=? AND IL.namespace=? ORDER BY service_name"
	if err := model.USE_MYSQL_DB().Select(&serviceInfoList, query, productName, clusterId, namespace); err != nil {
		log.Errorf("%v", err)
		return fmt.Errorf("Database query error %v", err)
	}

	// Avoid deploying the same product package with multiple namespaces
	//err, info := model.DeployProductList.GetCurrentProductInfoByName(productName)
	var info *model.DeployProductListInfo

	info, err = model.DeployClusterProductRel.GetCurrentProductByProductNameClusterIdNamespace(productName, clusterId, namespace)
	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"groups": groupAndServices,
			"count":  len(groupAndServices),
		}
	}
	if err != nil {
		return err
	}

	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		return err
	}

	serviceDisplayMap := map[string]string{}
	for name, svc := range sc.Service {
		if svc.ServiceDisplay != "" {
			serviceDisplayMap[name] = svc.ServiceDisplay
		}
	}

	var lastServiceName string
	svcMap := make(map[string]struct{}, 0)
	for _, info := range serviceInfoList {
		var alertState string
		var execFail bool
		if _, ok := svcMap[info.ServiceName]; !ok {
			alertState, err = getInstanceAlertState(productName, info.ServiceName, "")
			if err != nil {
				log.Errorf("%v", err)
			}
			infoList, err := model.HealthCheck.GetInfoByClusterIdAndProductNameAndServiceName(clusterId, productName, info.ServiceName, "")
			if err != nil {
				log.Errorf("%v", err)
			}
			for _, info := range infoList {
				if info.ExecStatus == enums.ExecStatusType.Failed.Code {
					execFail = true
					break
				}
			}
		}
		svcMap[info.ServiceName] = struct{}{}
		r := groupAndServices[info.Group]
		if info.ServiceName != lastServiceName {
			serviceDisplay, ok := serviceDisplayMap[info.ServiceName]
			if !ok {
				serviceDisplay = info.ServiceName
			}
			r = append(r, resultList{ServiceName: info.ServiceName, ServiceNameDisplay: serviceDisplay})
		}
		if info.Status != model.INSTANCE_STATUS_RUNNING {
			r[len(r)-1].Alert = true
		} else if info.HealthState != model.INSTANCE_HEALTH_OK && info.HealthState != model.INSTANCE_HEALTH_NOTSET {
			r[len(r)-1].Alert = true
		} else if alertState == "alert" || execFail == true {
			r[len(r)-1].Alert = true
		}
		groupAndServices[info.Group] = r
		lastServiceName = info.ServiceName
	}
	if cloud == "true" {
		// 解析schema中的服务组件的ServiceAddr
		if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), namespace); err != nil {
			log.Errorf("setSchemaFieldServiceAddr err: %v", err)
			return err
		}
		for name, svc := range sc.Service {
			if svc.Instance != nil && svc.Instance.UseCloud {
				r := groupAndServices[svc.Group]
				serviceDisplay, ok := serviceDisplayMap[name]
				if !ok {
					serviceDisplay = name
				}
				r = append(r, resultList{ServiceName: name, ServiceNameDisplay: serviceDisplay})
				groupAndServices[svc.Group] = r
			}
		}
	}
	return map[string]interface{}{
		"groups": groupAndServices,
		"count":  len(groupAndServices),
	}
}

func getInstanceAlertState(productName, serviceName, ip string) (string, error) {
	params := make([]map[string]string, 0)
	params = append(params, map[string]string{"tag": productName})
	params = append(params, map[string]string{"tag": serviceName})
	err, response := grafana.GetDashboard(params)
	if err != nil || len(response) == 0 {
		return "unknown", err
	}
	for _, resp := range response {
		alertList := ServiceAlertList(strconv.Itoa(resp.Id), ip)
		for _, alert := range alertList {
			if alert.State != "ok" && alert.State != "paused" && alert.State != "pending" {
				return "alert", nil
			}
		}
	}
	return "ok", nil
}

func HealthReport(ctx context.Context) apibase.Result {
	agentId, info := checkInstanceAgentId(ctx)
	params := &agent.HealthCheck{}
	if err := ctx.ReadJSON(&params); err != nil {
		log.Errorf("[HealthReport] err: %v", err)
		return fmt.Errorf("[HealthReport] err: %v", err)
	}
	_, instanceInfo := model.DeployInstanceList.GetInstanceInfoByAgentId(agentId)
	if instanceInfo.HealthState != -1 {
		checkInstanceStatus(info, append(model.OUT_OF_EVENTREPORT_STATUS_LIST, model.INSTANCE_STATUS_STOPPED))
		err := model.DeployInstanceList.UpdateInstanceHealthCheck(agentId, !params.Failed)
		if err != nil {
			return fmt.Errorf("[HealthReport] UpdateInstanceHealthCheck err: %v, agentId: %v", err, agentId)
		}
	}
	var status int
	if params.Failed {
		status = enums.ExecStatusType.Failed.Code
	} else {
		status = enums.ExecStatusType.Success.Code
	}
	execShellInfo, err := model.ExecShellList.GetBySeq(int(params.Seqno))
	//if err != nil {
	//	log.Debugf("ExecShellList GetBySeq error: %v", int(params.Seqno))
	//}
	if err == nil && execShellInfo.ExecStatus == enums.ExecStatusType.Running.Code {
		now := time.Now()
		duration := now.Sub(execShellInfo.CreateTime.Time).Seconds()
		err = model.ExecShellList.UpdateStatusBySeq(int(params.Seqno), status, dbhelper.NullTime{Time: now, Valid: true}, sql.NullFloat64{Float64: duration, Valid: true})
		if err != nil {
			log.Errorf("ExecShellList UpdateStatusBySeq error: %v", err)
		}
		err := UpdateOperationStatusBySeq(int(params.Seqno))
		if err != nil {
			log.Errorf("UpdateOperationStatusBySeq error: %v", err)
		}
	}

	ev := &event.Event{
		AgentId: agentId,
		Type:    event.REPORT_EVENT_HEALTH_CHECK,
		Data:    params,
	}
	event.GetEventManager().EventReciever(ev)
	return nil
}

func ErrorReport(ctx context.Context) apibase.Result {
	agentId, info := checkInstanceAgentId(ctx)

	checkInstanceStatus(info, append(model.OUT_OF_EVENTREPORT_STATUS_LIST, model.INSTANCE_STATUS_STOPPED))

	params := &agent.AgentError{}
	if err := ctx.ReadJSON(&params); err != nil {
		log.Errorf("[ErrorReport] err: %v", err)
		return fmt.Errorf("[ErrorReport] err: %v", err)
	}

	execShellInfo, err := model.ExecShellList.GetBySeq(int(params.Seqno))
	if err != nil {
		log.Debugf("ExecShellList GetBySeq error: %v", int(params.Seqno))
	}
	if err == nil && execShellInfo.ExecStatus == enums.ExecStatusType.Running.Code {
		now := time.Now()
		duration := now.Sub(execShellInfo.CreateTime.Time).Seconds()
		err = model.ExecShellList.UpdateStatusBySeq(int(params.Seqno), enums.ExecStatusType.Failed.Code, dbhelper.NullTime{Time: now, Valid: true}, sql.NullFloat64{Float64: duration, Valid: true})
		if err != nil {
			log.Errorf("ExecShellList UpdateStatusBySeq error: %v", err)
		}
		err := UpdateOperationStatusBySeq(int(params.Seqno))
		if err != nil {
			log.Errorf("UpdateOperationStatusBySeq error: %v", err)
		}
	}

	if strings.Contains(params.ErrStr, agent.IS_ALREADY_RUNNING) {
		return nil
	}
	if info.HealthState != model.INSTANCE_HEALTH_NOTSET {
		err = model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_RUN_FAIL,
			params.LastUpdateDate+":"+params.ErrStr, model.INSTANCE_HEALTH_WAITING)
	} else {
		err = model.DeployInstanceList.UpdateInstanceStatusByAgentId(agentId, model.INSTANCE_STATUS_RUN_FAIL,
			params.LastUpdateDate+":"+params.ErrStr)
	}
	if err != nil {
		return fmt.Errorf("[ErrorReport] UpdateInstanceStatus err: %v, agentId: %v", err, agentId)
	}
	ev := &event.Event{
		AgentId: agentId,
		Type:    event.REPORT_EVENT_INSTANCE_ERROR,
		Data:    params,
	}
	event.GetEventManager().EventReciever(ev)

	return nil
}

func HostResourcesReport(ctx context.Context) apibase.Result {
	return nil
}

func PerformanceReport(ctx context.Context) apibase.Result {
	agentIdStr := ctx.Params().Get("agent_id")

	if agentIdStr == "00000000-0000-0000-0000-000000000000" || strings.Contains(agentIdStr, "00000000") {
		//filter os metric
		params := &agent.AgentPerformance{}
		if err := ctx.ReadJSON(&params); err != nil {
			log.Errorf("[PerformanceReport] err: %v", err)
		} else {
			if params.Sid != "" {
				err := model.DeployHostList.UpdateUpdatedWithSid(params.Sid)
				if err != nil {
					model.DeployNodeList.UpdateUpdatedWithSid(params.Sid)
				}
			}
		}
		return nil
	}

	agentId, info := checkInstanceAgentId(ctx)
	checkInstanceStatus(info, model.OUT_OF_EVENTREPORT_STATUS_LIST)
	if err := model.DeployInstanceList.UpdateInstanceStatusByAgentPerformance(agentId); err != sql.ErrNoRows {
		return fmt.Errorf("[PerformanceReport] UpdateInstanceStatus err: %v, agentId: %v", err, agentId)
	}

	return nil
}

func InstanceUpdateRecord(ctx context.Context) apibase.Result {
	// by deploy_uuid get deploy instance record info
	paramErrs := apibase.NewApiParameterErrors()

	updateUUID := ctx.Params().Get("update_uuid")
	if updateUUID == "" {
		paramErrs.AppendError("$", fmt.Errorf("[Instance->InstanceUpdateRecord] update_uuid is empty"))
	}
	var status []string
	if ctx.URLParam("status") != "" {
		status = strings.Split(ctx.URLParam("status"), ",")
	}
	serviceName := ctx.URLParam("serviceName")

	paramErrs.CheckAndThrowApiParameterErrors()

	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.DeployInstanceUpdateRecordInfo{})
	info, count, complete := model.DeployInstanceUpdateRecord.GetDeployInstanceUpdateRecordByUpdateId(pagination, updateUUID, status, serviceName)

	list := []map[string]interface{}{}
	for _, s := range info {
		r := map[string]interface{}{}
		r["id"] = s.ID
		r["update_uuid"] = s.UpdateUUID
		r["instance_id"] = s.InstanceId
		r["sid"] = s.Sid
		r["ip"] = s.Ip
		if s.Schema.Valid {
			r["schema"] = s.Schema.String
		} else {
			r["schema"] = "[]"
		}
		r["product_name"] = s.ProductName
		r["product_version"] = s.ProductVersion
		r["group"] = s.Group
		r["service_name"] = s.ServiceName
		r["service_version"] = s.ServiceVersion
		r["status"] = s.Status
		r["status_message"] = s.StatusMessage
		r["progress"] = s.Progress

		if s.UpdateDate.Valid == true {
			r["update_time"] = s.UpdateDate.Time.Format(base.TsLayout)
		} else {
			r["update_time"] = ""
		}

		if s.CreateDate.Valid == true {
			r["create_time"] = s.CreateDate.Time.Format(base.TsLayout)
		} else {
			r["create_time"] = ""
		}

		list = append(list, r)
	}

	return map[string]interface{}{
		"list":     list,
		"count":    count,
		"complete": complete,
	}
}

const (
	autoDeployType   = "auto"
	manualDeployType = "manual"
)

func InstanceRecord(ctx context.Context) apibase.Result {
	// by deploy_uuid get deploy instance record info
	paramErrs := apibase.NewApiParameterErrors()

	deployUUID := ctx.Params().Get("deploy_uuid")
	if deployUUID == "" {
		paramErrs.AppendError("$", fmt.Errorf("[Instance->InstanceRecord] deploy_uuid is empty"))
	}
	var status []string
	if ctx.URLParam("status") != "" {
		status = strings.Split(ctx.URLParam("status"), ",")
	}
	serviceName := ctx.URLParam("serviceName")

	paramErrs.CheckAndThrowApiParameterErrors()

	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.DeployInstanceRecordInfo{})
	isManualDeployUUID := false
	uuidInfo, err := model.DeployUUID.GetInfoByUUID(deployUUID)
	//如果是没找到，那么应该是老的部署历史数据
	if errors.Is(err, sql.ErrNoRows) {
		isManualDeployUUID = true
		goto list
	}
	// db err
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	//查询到了 并且是自动部署
	if uuidInfo.UuidType == model.AutoDeployUUIDType || uuidInfo.UuidType == model.AutoDeployChildrenUUIDType {
		isManualDeployUUID = false
	} else {
		//手动部署
		isManualDeployUUID = true
	}
list:
	//如果是uuid 是自动部署或者自动部署的子产品 uuid 则所有uuid 部署完毕视为部署成功，list 接口查看所有uuid的部署信息
	var allComplete string
	var count int
	var allInfos []model.DeployInstanceRecordByDeployIdInfo
	//返回给前端部署类型，前端依靠不同部署类型调用想用类型的停止接口
	var deployType string
	if !isManualDeployUUID {
		deployType = autoDeployType
		var prodUuidInfos []model.DeployUUIDInfo
		//如果是uuid 为自动部署返回的 uuid
		if uuidInfo.UuidType == model.AutoDeployUUIDType {
			prodUuidInfos, err = model.DeployUUID.GetUUIDListByParentUUID(uuidInfo.UUID)
		} else {
			//如果是uuid 为自动产品中产生的uuid
			prodUuidInfos, err = model.DeployUUID.GetUUIDListByParentUUID(uuidInfo.ParentUUID)
		}
		//如果没查询到 正在准备状态
		if errors.Is(err, sql.ErrNoRows) {
			allComplete = model.PRODUCT_STATUS_PENDING
			return map[string]interface{}{
				"deploy_type": deployType,
				"list":        nil,
				"count":       0,
				"complete":    allComplete,
			}
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		allComplete = model.PRODUCT_STATUS_DEPLOYED

		//查询每一个uuid 的信息，并且每个都部署成功才视为最后成功
		for _, uuidInfo := range prodUuidInfos {

			//这里的分页查询为不限制
			productInstanceRecordInfo := getAutoDeployInstanceRecordByDeployId(&apibase.Pagination{}, uuidInfo, status, serviceName)
			allInfos = append(allInfos, productInstanceRecordInfo...)
			count += len(productInstanceRecordInfo)
			deployResInfo, err := model.DeployProductHistory.GetDeployHistoryByDeployUUID(uuidInfo.UUID)

			//如果有还有产品包没来得及部署
			if errors.Is(err, sql.ErrNoRows) {
				allComplete = model.PRODUCT_STATUS_DEPLOYING
				continue
			}

			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			//如果某个产品包部署状态为部署失败那么意味着这整个自动部署流程失败了
			if deployResInfo.Status == model.PRODUCT_STATUS_DEPLOY_FAIL {
				allComplete = model.PRODUCT_STATUS_DEPLOY_FAIL
				continue
			}
			//如果某个产品包部署状态为不是失败也不是部署完毕 那么意味着整个自动部署流程还在部署中
			if deployResInfo.Status != model.PRODUCT_STATUS_DEPLOYED {
				allComplete = model.PRODUCT_STATUS_DEPLOYING
			}
			//	如果每个产品包以上逻辑都没有进入，那么 allComplete 不会改变为部署成功
		}

	} else {
		deployType = manualDeployType
		//如果是手动部署 uuid 走老的逻辑
		allInfos, count, allComplete = model.DeployInstanceRecord.GetDeployInstanceRecordByDeployId(pagination, deployUUID, status, serviceName)
	}

	list := []map[string]interface{}{}
	for _, s := range allInfos {
		r := map[string]interface{}{}
		r["id"] = s.ID
		r["deploy_uuid"] = s.DeployUUID
		r["instance_id"] = s.InstanceId
		r["sid"] = s.Sid
		r["ip"] = s.Ip
		if s.Schema.Valid {
			r["schema"] = s.Schema.String
		} else {
			r["schema"] = "[]"
		}
		r["product_name"] = s.ProductName
		r["product_version"] = s.ProductVersion
		r["group"] = s.Group
		r["service_name"] = s.ServiceName
		r["service_version"] = s.ServiceVersion
		r["status"] = s.Status
		r["status_message"] = s.StatusMessage
		r["progress"] = s.Progress

		if s.UpdateDate.Valid == true {
			r["update_time"] = s.UpdateDate.Time.Format(base.TsLayout)
		} else {
			r["update_time"] = ""
		}

		if s.CreateDate.Valid == true {
			r["create_time"] = s.CreateDate.Time.Format(base.TsLayout)
		} else {
			r["create_time"] = ""
		}

		list = append(list, r)
	}

	//如果部署不是手动部署,分页处理
	if !isManualDeployUUID {
		count = len(list)
		//如果是记录条数小于本次请求区间上限
		if count < pagination.Start+pagination.Limit {
			list = list[pagination.Start:count]
		} else {
			//如果没有超过上限
			list = list[pagination.Start : pagination.Start+pagination.Limit]
		}
	}

	return map[string]interface{}{
		"deploy_type": deployType,
		"list":        list,
		"count":       count,
		"complete":    allComplete,
	}
}

//获取自动部署的记录
func getAutoDeployInstanceRecordByDeployId(pagination *apibase.Pagination, uuidInfo model.DeployUUIDInfo, status []string, serviceName string) []model.DeployInstanceRecordByDeployIdInfo {

	deployUUID := uuidInfo.UUID
	whereCause := dbhelper.MakeWhereCause()
	whereCause = whereCause.Equal("deploy_uuid", deployUUID)
	var values []interface{}
	if serviceName != "" {
		for _, v := range strings.Split(serviceName, ",") {
			values = append(values, v)
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("IR.service_name", values...)
	}
	if len(status) > 0 {
		s := make([]interface{}, 0, len(status))
		for _, s_ := range status {
			s = append(s, interface{}(s_))
		}
		whereCause = whereCause.And().Included("IR.status", s...)
	}
	where, value := whereCause.SQL()
	query := "SELECT IR.*, IL.schema FROM " +
		model.DeployInstanceRecord.TableName + " AS IR LEFT JOIN " +
		model.DeployInstanceList.TableName + " AS IL ON IR.instance_id = IL.id " + where + " " + pagination.AsQuery()
	queryCount := "SELECT COUNT(*) FROM " + model.DeployInstanceRecord.TableName + " AS IR " + where

	var list []model.DeployInstanceRecordByDeployIdInfo
	var count int
	if err := model.DeployInstanceRecord.GetDB().Get(&count, queryCount, value...); err != nil {
		log.Errorf("queryCount: %v, value: %v, err: %v", queryCount, value, err)
		apibase.ThrowDBModelError(err)
	}

	if count > 0 {
		rows, err := model.USE_MYSQL_DB().Queryx(query, value...)
		if err != nil {
			log.Errorf("query: %v, value: %v, err: %v", query, value, err)
			apibase.ThrowDBModelError(err)
		}

		defer rows.Close()

		for rows.Next() {
			info := model.DeployInstanceRecordByDeployIdInfo{}
			if err := rows.StructScan(&info); err != nil {
				apibase.ThrowDBModelError(err)
			}
			list = append(list, info)
		}
	}
	return list
}

func InstanceServiceLog(ctx context.Context) apibase.Result {
	// by instance id get deploy instance install and some log
	log.Debugf("[Instance->InstanceServiceLog] return deploy log info from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()

	id := ctx.Params().Get("id")

	logFile := ctx.FormValue("logfile")
	isMatch := ctx.FormValue("is_match")
	tailNumStr := ctx.FormValue("tail_num")
	tailNum, err := strconv.Atoi(tailNumStr)
	if err != nil {
		tailNum = 30
	}
	if id == "" {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	instanceId, err := strconv.Atoi(id)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	err, info := model.DeployInstanceList.GetInstanceInfoById(instanceId)
	if err != nil {
		return err
	}
	cluster, err := model.DeployClusterList.GetClusterInfoById(info.ClusterId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	var isDir int
	var result interface{}
	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_HOSTS {
		isDir, result, err = model.DeployInstanceList.GetDeployInstanceLogById(id, logFile, isMatch, tailNum)
	} else {
		isDir, result, err = model.DeployInstanceList.GetDeployPodInstanceLog(info, logFile, isMatch, tailNum)
	}
	if err != nil {
		log.Errorf(err.Error())
	}
	return map[string]interface{}{
		"logfile": logFile,
		"is_dir":  isDir,
		"result":  result,
	}
}

func InstanceBelongService(ctx context.Context) apibase.Result {
	// by product name and service name  get deploy instance  info
	log.Debugf("[Instance->InstanceBelongService] return deploy instance info by product name and service name for EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()

	productName := ctx.Params().Get("product_name")

	serviceName := ctx.Params().Get("service_name")
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)

	if productName == "" || serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name or service_name is empty"))
	}
	clusterIdStr := ctx.URLParam("clusterId")
	paramErrs.CheckAndThrowApiParameterErrors()
	var clusterId int
	var err error
	if clusterIdStr == "" {
		clusterId, err = GetCurrentClusterId(ctx)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
	} else {
		clusterId, err = strconv.Atoi(clusterIdStr)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
	}
	list := []map[string]interface{}{}

	var pr *model.DeployProductListInfo
	pr, err = model.DeployClusterProductRel.GetCurrentProductByProductNameClusterIdNamespace(productName, clusterId, namespace)
	//when watch the service in the front, uninstall the service will occur sql no errors cyclic
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			log.Errorf("%v", err)
			return err
		}
	}

	//pr, err = model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	////when watch the service in the front, uninstall the service will occur sql no errors cyclic
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		return nil
	//	} else {
	//		log.Errorf("%v", err)
	//		return err
	//	}
	//}

	sc, err := schema.Unmarshal(pr.Product)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), namespace); err != nil {
		log.Errorf("[Instance->InstanceBelongService] setSchemaFieldServiceAddr err: %v", err)
		return err
	}

	//正常部署的组件，非使用外部主机
	info := make([]model.InstanceAndProductInfo, 0)
	var count int

	info, count = model.DeployInstanceList.GetInstanceBelongServiceWithNamespace(productName, serviceName, clusterId, namespace)
	//info, count = model.DeployInstanceList.GetInstanceBelongService(productName, serviceName, clusterId)

	for _, s := range info {
		r := map[string]interface{}{}
		r["id"] = s.ID
		r["agent_id"] = s.AgentId
		r["sid"] = s.Sid
		r["pid"] = s.Pid
		r["ip"] = s.Ip
		r["product_name"] = s.ProductName
		r["product_version"] = s.ProductVersion
		r["group"] = s.Group
		r["prometheus_port"] = s.PrometheusPort
		r["health_state"] = s.HealthState
		r["service_name"] = s.ServiceName
		r["service_version"] = s.ServiceVersion
		r["status"] = s.Status
		r["status_message"] = s.StatusMessage
		r["schema"] = string(s.Schema)

		roleData := harole.RoleData(s.Pid, serviceName)
		if roleData != nil {
			haRole, ok := roleData[s.AgentId]
			if !ok {
				haRole = "-"
			}
			r["ha_role"] = haRole
		}

		if s.HeartDate.Valid == true {
			r["heart_time"] = s.HeartDate.Time.Format(base.TsLayout)
		} else {
			r["heart_time"] = ""
		}

		if s.UpdateDate.Valid == true {
			r["update_time"] = s.UpdateDate.Time.Format(base.TsLayout)
		} else {
			r["update_time"] = ""
		}

		if s.CreateDate.Valid == true {
			r["create_time"] = s.CreateDate.Time.Format(base.TsLayout)
		} else {
			r["create_time"] = ""
		}
		// k8s模式下才有该字段
		if s.Namespace != "" {
			r["namespace"] = s.Namespace
		}

		list = append(list, r)
	}
	//兼容使用外部主机请服务
	for name, svc := range sc.Service {
		if name == serviceName && svc.Instance != nil && svc.Instance.UseCloud {
			//已经部署过的组件，误设置了使用外部主机
			_, count := model.DeployInstanceList.GetInstanceBelongService(productName, serviceName, clusterId)
			if count > 0 {
				continue
			}
			err, ipInfo := model.DeployServiceIpList.GetServiceIpListByName(productName, serviceName, clusterId, "")
			if err != nil {
				return err
			}
			ipList := strings.Split(ipInfo.IpList, IP_LIST_SEP)
			for _, ip := range ipList {
				r := map[string]interface{}{}
				r["ip"] = ip
				r["product_name"] = sc.ProductName
				r["product_version"] = sc.ProductVersion
				r["group"] = svc.Group
				r["service_name"] = serviceName
				r["service_name"] = svc.Version
				list = append(list, r)
			}
			return map[string]interface{}{
				"use_cloud": true,
				"list":      list,
				"count":     len(list),
			}
		}
	}

	return map[string]interface{}{
		"list":  list,
		"count": count,
	}
}

func InstanceServiceConfig(ctx context.Context) apibase.Result {
	// by product name and service name  get deploy instance  info
	log.Debugf("[Instance->InstanceServiceConfig] return instance config by instance primary key id for EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()

	instanceIdStr := ctx.Params().Get("id")
	configfile := ctx.FormValue("configfile")
	instanceId, err := strconv.Atoi(instanceIdStr)
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)

	if instanceIdStr == "" || err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	err, info := model.DeployInstanceList.GetInstanceInfoById(instanceId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	cluster, err := model.DeployClusterList.GetClusterInfoById(info.ClusterId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	// 区分主机模式和k8s模式在deploy_cluster_product_rel表中查询
	var productRel model.ClusterProductRel

	productRel, err = model.DeployClusterProductRel.GetByPidAndClusterIdNamespace(info.Pid, info.ClusterId, namespace)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("[InstanceServiceConfig] get namespace by pid and clusterid error: %v", err)
		return fmt.Errorf("[InstanceServiceConfig] get namespace by pid and clusterid error: %v", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		smoothUpgradeProductRel, err := model.DeployClusterSmoothUpgradeProductRel.GetByPidAndClusterIdNamespace(info.Pid, info.ClusterId, namespace)
		if err != nil {
			log.Errorf("[InstanceServiceConfig] get namespace by pid and clusterid error: %v", err)
			return fmt.Errorf("[InstanceServiceConfig] get namespace by pid and clusterid error: %v", err)
		}
		productRel = model.ClusterProductRel{
			Id:            smoothUpgradeProductRel.Id,
			Pid:           smoothUpgradeProductRel.Pid,
			ClusterId:     smoothUpgradeProductRel.ClusterId,
			Namespace:     smoothUpgradeProductRel.Namespace,
			ProductParsed: smoothUpgradeProductRel.ProductParsed,
			Status:        smoothUpgradeProductRel.Status,
			DeployUUID:    smoothUpgradeProductRel.DeployUUID,
			AlertRecover:  smoothUpgradeProductRel.AlertRecover,
			UserId:        smoothUpgradeProductRel.UserId,
			IsDeleted:     smoothUpgradeProductRel.IsDeleted,
			UpdateTime:    smoothUpgradeProductRel.UpdateTime,
			DeployTime:    smoothUpgradeProductRel.DeployTime,
			CreateTime:    smoothUpgradeProductRel.CreateTime,
		}
	}

	//productRel, err = model.DeployClusterProductRel.GetByPidAndClusterId(info.Pid, info.ClusterId)
	//if err != nil {
	//	log.Errorf("%v", err)
	//	return err
	//}

	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_HOSTS {
		res, err := model.DeployInstanceList.GetInstanceServiceConfig(instanceId, configfile)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		return map[string]interface{}{"result": res}
	}

	sc, err := schema.Unmarshal(productRel.ProductParsed)
	if err != nil {
		return err
	}
	nsTbsc, err := modelkube.DeployNamespaceList.Get(productRel.Namespace, info.ClusterId)
	if err != nil {
		return err
	}
	cache, err := kube.ClusterNsClientCache.GetClusterNsClient(strconv.Itoa(info.ClusterId)).GetClientCache(kube.ImportType(nsTbsc.Type))
	if err != nil {
		return err
	}
	resp, err := kdeploy.GetConfigMaps(cache, sc, cluster.Id, productRel.Namespace, info.ServiceName)
	if err != nil {
		return err
	}
	configMap, err := kmodel.ConvertConfigMap(resp)
	if err != nil {
		return err
	}
	targetFile := ""
	// Files with the same name are not supported temporarily
	if sc.DeployType == "workload" {
		targetFile = configfile[strings.LastIndex(configfile, "/")+1:]
	} else {
		targetFile = strings.Replace(configfile, "/", "_", -1)
	}

	return map[string]interface{}{"result": configMap.Data[targetFile]}

}

func InstanceEvent(ctx context.Context) apibase.Result {
	// by product name and service name  get deploy instance  info
	log.Debugf("[Instance->InstanceEvent] return instance event by instance primary key id and event id for EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()

	instanceIdStr := ctx.Params().Get("id")
	eventIdStr := ctx.URLParam("eventId")

	instanceId, err := strconv.Atoi(instanceIdStr)
	if instanceIdStr == "" || err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	eventId, err := strconv.Atoi(eventIdStr)
	if eventIdStr == "" || err != nil {
		paramErrs.AppendError("$", fmt.Errorf("event id is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	event, err := model.DeployInstanceEvent.GetEventInfoByInstanceIdEventId(int64(instanceId), int64(eventId))

	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	return event.Content
}

func InstancerControllTest(ctx context.Context) apibase.Result {
	// by product name and service name  get deploy instance  info
	log.Debugf("[Instance->InstancerControllTest]")
	//runtime.GC()

	params := instance.NewInstancerParam{}

	if err := ctx.ReadJSON(&params); err != nil {
		log.Errorf("[InstancerControllTest] err: %v", err)
		return fmt.Errorf("[InstancerControllTest] err: %v", err)
	}

	log.Debugf("[InstancerControllTest] params: %v", params)

	sc := &schema.SchemaConfig{}
	err := json.Unmarshal([]byte(params.Schema), sc)

	if err != nil {
		log.Errorf("[InstancerControllTest] err: %v", err)
		return fmt.Errorf("[InstancerControllTest] err: %v", err)
	}

	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	//model.OperationList.Insert(model.OperationInfo{
	//	OperationId:     operationId,
	//	OperationType:   enums.OperationType.ProductDeploy.Code,
	//	OperationStatus: enums.ExecStatusType.Running.Code,
	//	ObjectType:      enums.OperationObjType.Svc.Code,
	//	ObjectValue:     params.ServiceName,
	//})
	inst, err := instance.NewInstancer(params.Pid, params.Ip, params.ServiceName, sc, operationId)

	if err != nil {
		log.Errorf("[InstancerControllTest] err: %v", err)
		return fmt.Errorf("[InstancerControllTest] err: %v", err)
	}

	err = inst.Install(false)
	if err != nil {
		log.Errorf("[InstancerControllTest] err: %v", err)
		return fmt.Errorf("[InstancerControllTest] err: %v", err)
	}

	err = inst.Start()
	if err != nil {
		log.Errorf("[InstancerControllTest] err: %v", err)
		return fmt.Errorf("[InstancerControllTest] err: %v", err)
	}

	//err = inst.Stop()
	//if err != nil {
	//	log.Errorf("[InstancerControllTest] err: %v", err.Error())
	//	return fmt.Errorf("[InstancerControllTest] err: %v", err.Error())
	//}
	//
	//err = inst.UnInstall(false)
	//if err != nil {
	//	log.Errorf("[InstancerControllTest] err: %v", err.Error())
	//	return fmt.Errorf("[InstancerControllTest] err: %v", err.Error())
	//}

	return nil
}

func ListLogFiles(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("ListLogFiles: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[ListLogFiles] err: %v", rlt)
		}
	}()
	paramErrs := apibase.NewApiParameterErrors()

	instanceIdStr := ctx.Params().Get("instance_id")

	ftype := ctx.URLParam("type")

	instanceId, err := strconv.Atoi(instanceIdStr)

	if instanceIdStr == "" || err != nil {
		paramErrs.AppendError("$", fmt.Errorf("instance_id is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	var lists []string

	err, info := model.DeployInstanceList.GetInstanceInfoById(instanceId)
	//兼容组件部署失败的情况
	if err != nil && strings.Contains(err.Error(), ERROR_NO_ROW_IN_RESULT) {
		return map[string]interface{}{
			"count": len(lists),
			"list":  lists,
		}
	}
	cluster, err := model.DeployClusterList.GetClusterInfoById(info.ClusterId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	product, err := model.DeployProductList.GetProductInfoById(info.Pid)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	sc := &schema.ServiceConfig{}

	if err = json.Unmarshal([]byte(info.Schema), sc); err != nil {
		return err
	}
	if sc.Instance == nil {
		return fmt.Errorf("service %v have no log", info.ServiceName)
	}
	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_HOSTS {
		listCmd := "#!/bin/sh\nfind -L %s -maxdepth 5 -size -4096M -type f -mtime -7 -print"
		for _, dir := range sc.Instance.Logs {
			cmd := listCmd
			dir = dir + "*"
			cmd = fmt.Sprintf(cmd, dir)
			log.Debugf("listCmd: %v", cmd)
			logs, err := agent.AgentClient.ToExecCmd(info.Sid, info.AgentId, cmd, "")
			log.Debugf("response: %v", logs)
			if err != nil && strings.Contains(err.Error(), "No such file or directory") {
				continue
			}
			if err != nil {
				return err
			}
			lists = append(lists, strings.Split(logs, LINUX_SYSTEM_LINES)...)
		}
	} else {
		err, pod := model.DeployKubePodList.GetPodInfoByAgentId(info.AgentId)
		if err != nil {
			return err
		}
		///tmp/dtstack/ 为业务测约定约定写入hostpath的路径
		baseDir := "/host/tmp/dtstack/" + product.ProductName + "/" + info.ServiceName + "/" + pod.PodName + "/"
		listCmd := "#!/bin/sh\nfind -L %s -maxdepth 5 -size -1048576k -type f -mtime -7 -print"
		for _, dir := range sc.Instance.Logs {
			cmd := listCmd
			dir = baseDir + dir + "*"
			cmd = fmt.Sprintf(cmd, dir)
			log.Debugf("listCmd: %v", cmd)
			logs, err := agent.AgentClient.ToExecCmd(info.Sid, "", cmd, "")
			log.Debugf("response: %v", logs)
			if err != nil && strings.Contains(err.Error(), "No such file or directory") {
				continue
			}
			if err != nil {
				return err
			}
			lists = append(lists, strings.Split(logs, LINUX_SYSTEM_LINES)...)
		}
	}
	//filter by txt or zip
	var txtList, zipList []string
	for _, line := range lists {
		if len(line) == 0 || line == "" {
			continue
		}
		isZip := false
		for _, et := range LINUX_SYSTEM_ZIP_EXTENTION {
			if strings.Contains(line, et) {
				isZip = true
				break
			}
		}
		if isZip {
			zipList = append(zipList, line)
		} else {
			txtList = append(txtList, line)
		}
	}
	if ftype == LINUX_FILE_TYPE_ZIP {
		lists = zipList
	} else {
		lists = txtList
	}
	sort.Strings(lists)
	return map[string]interface{}{
		"count": len(lists),
		"list":  lists,
	}
}

func PreviewLogFile(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("PreviewLogFile: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[PreviewLogFile] err: %v", rlt)
		}
	}()
	paramErrs := apibase.NewApiParameterErrors()

	instanceIdStr := ctx.Params().Get("instance_id")
	logFile := ctx.URLParam("logfile")
	start := ctx.URLParam("start")
	end := ctx.URLParam("end")

	instanceId, err := strconv.Atoi(instanceIdStr)
	if instanceIdStr == "" || logFile == "" || err != nil {
		paramErrs.AppendError("$", fmt.Errorf("instance_id or logfile is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	err, info := model.DeployInstanceList.GetInstanceInfoById(instanceId)
	if err != nil {
		return err
	}
	cluster, err := model.DeployClusterList.GetClusterInfoById(info.ClusterId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	var lists []string
	var session *preSession
	var agentId string
	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_HOSTS {
		agentId = info.AgentId
	} else {
		agentId = ""
	}

	requestedFile := filepath.Clean(filepath.Join("/", logFile))
	rel, err := filepath.Rel("/", requestedFile)
	if err != nil {
		log.Errorf("failed to get the relative path, err: %v", err)
		return err
	}
	pwdCmd := "#!/bin/sh\n echo `pwd`"
	res, err := agent.AgentClient.ToExecCmd(info.Sid, agentId, pwdCmd, "")
	if err != nil {
		return err
	}
	absolutePath := strings.Replace(res, LINUX_SYSTEM_LINES, "", -1)
	logFilePath := filepath.Join(absolutePath, rel)

	prevewCmd := fmt.Sprintf("#!/bin/sh\n cat %s | wc -l", logFilePath)
	total, err := agent.AgentClient.ToExecCmd(info.Sid, agentId, prevewCmd, "")
	if err != nil {
		return err
	}
	total = strings.Replace(total, LINUX_SYSTEM_LINES, "", -1)
	log.Debugf("total: %v", total)
	totalInt, err := strconv.ParseInt(total, 10, 64)
	if err != nil {
		return err
	}
	var startInt, endInt int64
	if start == "" {
		startInt = totalInt - 300 + 1
		endInt = totalInt
	} else {
		startInt, err = strconv.ParseInt(start, 10, 64)
		if err != nil {
			return err
		}
		endInt, err = strconv.ParseInt(end, 10, 64)
		if err != nil {
			return err
		}
	}
	if startInt < 1 {
		startInt = 1
	}
	if endInt < 1 {
		endInt = 1
	}

	log.Debugf("session: %v", session)
	previewCmd := fmt.Sprintf("#!/bin/sh\n sed -n '%d,%dp' %s ", startInt, endInt, logFile)
	log.Debugf("preview cmd: %v", previewCmd)
	content, err := agent.AgentClient.ToExecCmd(info.Sid, agentId, previewCmd, "")
	if err != nil {
		return err
	}
	log.Debugf("preview response: %v", content)
	lists = strings.Split(content, LINUX_SYSTEM_LINES)
	if len(lists) > 0 && lists[len(lists)-1] == "" {
		return map[string]interface{}{
			"count": len(lists[:len(lists)-1]),
			"list":  lists[:len(lists)-1],
			"total": totalInt,
		}
	}
	return map[string]interface{}{
		"count": len(lists),
		"list":  lists,
		"total": totalInt,
	}
}

func processPreviewSession(info *model.DeployInstanceInfo, sessionMd5, action, agentId string) (error, *preSession) {
	sMu.Lock()
	defer sMu.Unlock()
	session := sessionMap[sessionMd5]
	ancher := session.anchor
	switch action {
	case LOG_MORE_ACTION_UP:
		if ancher-int64(LOG_MORE_PREVIEW_RATE) > int64(LOG_MORE_PREVIEW_COUNT) {
			session.anchor = ancher - int64(LOG_MORE_PREVIEW_RATE)
		} else {
			session.anchor = int64(LOG_MORE_PREVIEW_COUNT)
		}
	case LOG_MORE_ACTION_DOWN:
		session.anchor = ancher + int64(LOG_MORE_PREVIEW_RATE)
	case LOG_MORE_ACTION_LATEST:
		prevewCmd := fmt.Sprintf("#!/bin/sh\n cat %s | wc -l", sessionMap[sessionMd5].log)
		content, err := agent.AgentClient.ToExecCmd(info.Sid, agentId, prevewCmd, "")
		if err != nil {
			return err, nil
		}
		content = strings.Replace(content, LINUX_SYSTEM_LINES, "", -1)
		count, err := strconv.ParseInt(content, 10, 64)
		if err != nil {
			return err, nil
		}
		session.anchor = count
	default:
		return fmt.Errorf("wrong action"), nil
	}
	session.lastAction = action
	sessionMap[sessionMd5] = session
	return nil, session
}
func initPreviewSession(info *model.DeployInstanceInfo, log, action, agentId string) (error, *preSession) {
	sMu.Lock()
	defer sMu.Unlock()
	nano := strconv.FormatInt(time.Now().UnixNano(), 10)
	//uuid 19_redis_lgs/redis.log__1399999999
	uuid := strconv.Itoa(info.ID) + "_" + info.ServiceName + "_" + log + "_" + nano
	session := &preSession{id: util.Md5(uuid)}
	if _, ok := sessionMap[uuid]; ok {
		sessionMap[uuid] = nil
	}
	prevewCmd := fmt.Sprintf("#!/bin/sh\n cat %s | wc -l", log)
	content, err := agent.AgentClient.ToExecCmd(info.Sid, agentId, prevewCmd, "")
	if err != nil {
		return err, nil
	}
	content = strings.Replace(content, LINUX_SYSTEM_LINES, "", -1)
	count, err := strconv.ParseInt(content, 10, 64)
	if err != nil {
		return err, nil
	}
	session.instanceId = strconv.Itoa(info.ID)
	session.log = log
	session.logWcl = count
	session.anchor = count
	session.lastAction = action
	session.startTimeNano = time.Now().UnixNano()
	sessionMap[session.id] = session
	return nil, session
}

func DownloadLogFile(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("DownloadLogFile: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[DownloadLogFile] err: %v", rlt)
		}
	}()
	paramErrs := apibase.NewApiParameterErrors()

	logFile := ctx.URLParam("logfile")
	instanceIdStr := ctx.Params().Get("instance_id")
	instanceId, err := strconv.Atoi(instanceIdStr)
	if instanceIdStr == "" || logFile == "" || err != nil {
		paramErrs.AppendError("$", fmt.Errorf("instance_id or logfile is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	err, info := model.DeployInstanceList.GetInstanceInfoById(instanceId)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(logFile, LINUX_SYSTEM_SLASH) {
		pwdCmd := "#!/bin/sh\n echo `pwd`"
		content, err := agent.AgentClient.ToExecCmd(info.Sid, info.AgentId, pwdCmd, "")
		if err != nil {
			return err
		}
		log.Debugf("pwd response: %v", content)
		logFile = strings.Replace(content, LINUX_SYSTEM_LINES, "", -1) + LINUX_SYSTEM_SLASH + logFile
		log.Debugf("logfile: %v", logFile)
	}

	ip := strings.Split(info.Ip, "/")[0]
	target := ip + EASYFILER_PORT

	var data = make(chan []byte, 1)
	ctx.Header("Content-Disposition", "attachment;filename="+logFile+TAR_SUFFIX)
	cancel := make(chan string, 1)
	go func() {
		if err := handler.DownloadWithoutStorage(target, logFile, data, cancel); err == io.EOF {
			log.Infof("DownloadLogFile: %v succeed", ctx.Request().RequestURI)
			return
		}
		log.Errorf("DownloadLogFile: %v failed", ctx.Request().RequestURI)
		return
	}()

LOOP:
	for {
		select {
		case ca := <-cancel:
			log.Errorf("cancel err %v", ca)
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

	return apibase.EmptyResult{}
}

func _DownloadLogFile(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("DownloadLogFile: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[DownloadLogFile] err: %v", rlt)
		}
	}()

	paramErrs := apibase.NewApiParameterErrors()

	logFile := ctx.URLParam("logfile")
	instanceIdStr := ctx.Params().Get("instance_id")
	instanceId, err := strconv.Atoi(instanceIdStr)
	if instanceIdStr == "" || logFile == "" || err != nil {
		paramErrs.AppendError("$", fmt.Errorf("instance_id or logfile is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	err, info := model.DeployInstanceList.GetInstanceInfoById(instanceId)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(logFile, LINUX_SYSTEM_SLASH) {
		pwdCmd := "#!/bin/sh\n echo `pwd`"
		content, err := agent.AgentClient.ToExecCmd(info.Sid, info.AgentId, pwdCmd, "")
		if err != nil {
			return err
		}
		log.Debugf("pwd response: %v", content)
		logFile = strings.Replace(content, LINUX_SYSTEM_LINES, "", -1) + LINUX_SYSTEM_SLASH + logFile
		log.Debugf("logfile: %v", logFile)
	}
	target := info.Ip + EASYFILER_PORT
	log.Infof("DownloadLogFile: %v", target)

	if err := handler.Download(target, logFile); err != nil {
		return err
	}

	if !util.IsPathExist(EASYFILER_TMP_ROOT + target + LINUX_SYSTEM_SLASH + logFile) {
		return fmt.Errorf("attempt to download a non-existent file")
	}

	ctx.SendFile(EASYFILER_TMP_ROOT+target+LINUX_SYSTEM_SLASH+logFile, strings.SplitAfter(logFile, LINUX_SYSTEM_SLASH)[len(strings.SplitAfter(logFile, LINUX_SYSTEM_SLASH))-1])
	return nil
}

func EventTypeList(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("EventTypeList: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[EventTypeList] err: %v", rlt)
		}
	}()

	var eventTypes []string
	query := `SELECT DISTINCT deploy_instance_runtime_event.event_type FROM deploy_instance_runtime_event WHERE deploy_instance_runtime_event.isDeleted=0`
	if err := model.USE_MYSQL_DB().Select(&eventTypes, query); err == sql.ErrNoRows {
		return map[string]interface{}{
			"count": 0,
			"list":  nil,
		}
	} else if err != nil {
		log.Errorf("event type list query: %v, err: %v", query, err)
		apibase.ThrowDBModelError(err)
	}

	return map[string]interface{}{
		"count": len(eventTypes),
		"list":  eventTypes,
	}
}

func EventList(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("EventList: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[EventList] err: %v", rlt)
		}
	}()
	paramErrs := apibase.NewApiParameterErrors()

	parentProductName := ctx.URLParam("parentProductName")
	eventType := ctx.URLParam("eventType")
	productNames := ctx.URLParam("productNames")
	serviceNames := ctx.URLParam("serviceNames")
	hosts := ctx.URLParam("hosts")
	keyWord := ctx.URLParam("keyWord")

	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.EventInfo{})

	if parentProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("parent_product_name is empty"))
	}
	if eventType == "" {
		paramErrs.AppendError("$", fmt.Errorf("event_type is empty"))
	}

	from, err := ctx.URLParamInt64("from")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("invalid start time"))
	}
	startTime := time.Unix(from, 0).Format(TIME_LAYOUT)
	to, err := ctx.URLParamInt64("to")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("invalid end time"))
	}
	endTime := time.Unix(to, 0).Format(TIME_LAYOUT)

	paramErrs.CheckAndThrowApiParameterErrors()

	var whereCause string
	var values []interface{}
	var count int
	_eventUrlParamsParse(ctx, true, &whereCause, &values)
	query := `SELECT COUNT(*) FROM deploy_instance_runtime_event `
	query += whereCause
	if keyWord != "" {
		query += ` AND deploy_instance_runtime_event.content Like ?`
		values = append(values, "%"+keyWord+"%")
	}
	if err := model.USE_MYSQL_DB().Get(&count, query, values...); err != nil {
		return err
	}

	eventInfoList, err := model.EventList.SelectEventListByWhere(pagination, eventType, parentProductName, productNames, serviceNames, hosts, startTime, endTime, keyWord)
	if err != nil {
		return err
	}

	var eventContentList []string
	for _, l := range eventInfoList {
		eventContentList = append(eventContentList, l.Content)
	}

	return map[string]interface{}{
		"count": count,
		"list":  eventContentList,
	}
}

type eventRankResponse struct {
	Name  string `db:"name" json:"name"`
	Times int    `db:"times" json:"times"`
	Rank  int    `json:"rank"`
}

func _eventUrlParamsParseIn(param, name string, whereCause *string, values *[]interface{}) {
	if param != "" {
		*whereCause += ` AND ` + name + ` IN (`
		for i, v := range strings.Split(param, ",") {
			if i > 0 {
				*whereCause += `,`
			}
			*whereCause += `?`
			*values = append(*values, v)
		}
		*whereCause += `)`
	}
}

func _eventUrlParamsParse(ctx context.Context, defaultTimeQuantum bool, whereCause *string, values *[]interface{}) {
	parentProductName := ctx.URLParam("parentProductName")
	eventType := ctx.URLParam("eventType")
	productNames := ctx.URLParam("productNames")
	serviceNames := ctx.URLParam("serviceNames")
	hosts := ctx.URLParam("hosts")

	paramErrs := apibase.NewApiParameterErrors()
	if eventType == "" {
		paramErrs.AppendError("$", fmt.Errorf("param $(eventType) is required"))
	}
	if parentProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("param $(parentProductName) is required"))
	}

	from, err := ctx.URLParamInt64("from")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("invalid start time"))
	}
	startTime := time.Unix(from, 0).Format(TIME_LAYOUT)
	to, err := ctx.URLParamInt64("to")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("invalid end time"))
	}
	endTime := time.Unix(to, 0).Format(TIME_LAYOUT)
	paramErrs.CheckAndThrowApiParameterErrors()

	if to <= from {
		paramErrs.AppendError("$", fmt.Errorf("invalid params: $(from) must be less than param $(to)"))
	}

	*whereCause += ` WHERE deploy_instance_runtime_event.event_type=? AND deploy_instance_runtime_event.parent_product_name=?`
	*values = append(*values, eventType, parentProductName)

	if defaultTimeQuantum {
		*whereCause += ` AND deploy_instance_runtime_event.create_time BETWEEN ? AND ?`
		*values = append(*values, startTime, endTime)
	}
	_eventUrlParamsParseIn(productNames, "deploy_instance_runtime_event.product_name", whereCause, values)
	_eventUrlParamsParseIn(serviceNames, "deploy_instance_runtime_event.service_name", whereCause, values)
	_eventUrlParamsParseIn(hosts, "deploy_instance_runtime_event.host", whereCause, values)

	*whereCause += ` AND deploy_instance_runtime_event.isDeleted=0`
}

func EventRank(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("EventRank: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[EventRank] err: %v", rlt)
		}
	}()

	groupBy := ctx.Params().Get("product_or_service")

	if groupBy != "product" && groupBy != "service" {
		return fmt.Errorf("invalid param {product_or_service}")
	}

	var whereCause string
	var values []interface{}
	_eventUrlParamsParse(ctx, true, &whereCause, &values)

	var eventRankResList []eventRankResponse
	query := `SELECT deploy_instance_runtime_event.%s_name AS name,COUNT(*) AS times FROM deploy_instance_runtime_event  
				` + whereCause + ` GROUP BY deploy_instance_runtime_event.%s_name ORDER BY times DESC LIMIT 0, 5`
	query = fmt.Sprintf(query, groupBy, groupBy)
	if err := model.USE_MYSQL_DB().Select(&eventRankResList, query, values...); err == sql.ErrNoRows {
		return ""
	} else if err != nil {
		return err
	}
	for i := range eventRankResList {
		eventRankResList[i].Rank = i + 1
	}
	if len(eventRankResList) > 5 {
		return eventRankResList[0:5]
	}
	return eventRankResList
}

type eventStatisticsResponse struct {
	Count        int `db:"count" json:"count"`
	ProductCount int `db:"product_count" json:"product_count"`
	ServiceCount int `db:"service_count" json:"service_count"`
}

func EventStatistics(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("EventStatistics: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[EventStatistics] err: %v", rlt)
		}
	}()

	var whereCause string
	var values []interface{}
	_eventUrlParamsParse(ctx, true, &whereCause, &values)

	query := `SELECT COUNT(*) AS count,COUNT(DISTINCT deploy_instance_runtime_event.product_name) AS product_count,COUNT(DISTINCT deploy_instance_runtime_event.service_name) AS service_count 
			FROM deploy_instance_runtime_event ` + whereCause
	var eventStatistics eventStatisticsResponse
	if err := model.USE_MYSQL_DB().Get(&eventStatistics, query, values...); err == sql.ErrNoRows {
		return eventStatisticsResponse{
			Count:        0,
			ProductCount: 0,
			ServiceCount: 0,
		}
	} else if err != nil {
		return err
	}

	return eventStatistics
}

type coordinateResponse struct {
	Count int    `db:"count" json:"count"`
	Date  string `json:"date"`
}

func EventTimeCoordinate(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("EventStatistics: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[EventStatistics] err: %v", rlt)
		}
	}()

	paramErrs := apibase.NewApiParameterErrors()
	from, err := ctx.URLParamInt64("from")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("invalid start time"))
	}
	to, err := ctx.URLParamInt64("to")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("invalid end time"))
	}
	if to <= from {
		paramErrs.AppendError("$", fmt.Errorf("invalid params: $(from) must be less than param $(to)"))
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	var whereCause string
	var values []interface{}
	var coordinateResponseList []coordinateResponse

	_eventUrlParamsParse(ctx, false, &whereCause, &values)
	query := `SELECT COUNT(*) FROM deploy_instance_runtime_event ` + whereCause + ` AND deploy_instance_runtime_event.create_time BETWEEN ? AND ?`

	for i := from; i < to; i += MINUTES_PER_DAY {
		var count int
		if err := model.USE_MYSQL_DB().Get(&count, query, append(values, time.Unix(i, 0).Format(TIME_LAYOUT), time.Unix(i+MINUTES_PER_DAY, 0).Format(TIME_LAYOUT))...); err != nil {
			return err
		}
		coordinateResponseList = append(coordinateResponseList, coordinateResponse{
			Count: count,
			Date:  time.Unix(i, 0).Format(DATE_LAYOUT),
		})
	}

	return map[string]interface{}{
		"count": len(coordinateResponseList),
		"list":  coordinateResponseList,
	}
}

func DiscoveryReload(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("DiscoveryReload: %v", ctx.Request().RequestURI)
	defer func() {
		if _, ok := rlt.(error); ok {
			log.Errorf("[DiscoveryReload] err: %v", rlt)
		}
	}()
	discovery := ctx.URLParam("discovery")

	switch discovery {
	case "service":
		discover.FlushServiceDiscover()
	case "node":
		discover.FlushNodeDiscover()
	default:
		return errors.New("invalid discovery")

	}
	return nil
}

func InstanceServiceAlert(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	dashboardId := ctx.URLParam("dashboardId")
	ips := ctx.URLParam("ip")
	if dashboardId == "" {
		paramErrs.AppendError("$", fmt.Errorf("param dashboardId is empty"))
	}
	if ips == "" {
		paramErrs.AppendError("$", fmt.Errorf("param ip is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	alertList := ServiceAlertList(dashboardId, ips)
	return map[string]interface{}{
		"count": len(alertList),
		"data":  alertList,
	}
}

type serviceAlertInfo struct {
	PanelTitle    string `json:"panel_title"`
	AlertName     string `json:"alert_name"`
	DashboardName string `json:"dashboard_name"`
	Url           string `json:"url"`
	State         string `json:"state"`
	Time          string `json:"time"`
}

func ServiceAlertList(dashboardId, ips string) []serviceAlertInfo {
	ipArr := strings.Split(ips, ",")
	ip := make(map[string]struct{}, 0)
	for _, k := range ipArr {
		ip[k] = struct{}{}
	}

	alertList := make([]serviceAlertInfo, 0)
	param := map[string]string{
		"dashboardId": dashboardId,
	}
	err, alerts := grafana.GrafanaAlertsSearch(param)
	if err != nil {
		log.Errorf("grafana search alerts error: %v", err)
		return alertList
	}

	for _, alert := range alerts {
		panelTitle, dashboardName := RetrievePanelTitle(alert.DashboardUid, alert.PanelId)
		//no_data, paused,alerting,ok, pending
		if alert.State == "ok" || alert.State == "paused" {
			alert.NewStateDate = ""
		}
		if ips == "" || alert.State != "alerting" || alert.EvalData.EvalMatches == nil {
			alert := serviceAlertInfo{
				PanelTitle:    panelTitle,
				State:         alert.State,
				AlertName:     alert.Name,
				DashboardName: dashboardName,
				Url:           alert.Url,
				Time:          alert.NewStateDate,
			}
			alertList = append(alertList, alert)
		} else if ips != "" && alert.EvalData.EvalMatches != nil {
			exist := false
			for _, match := range alert.EvalData.EvalMatches {
				if instance, ok := match.Tags["instance"]; ok {
					if _, oks := ip[strings.Split(instance, ":")[0]]; oks && !exist {
						alert := serviceAlertInfo{
							PanelTitle:    panelTitle,
							State:         alert.State,
							AlertName:     alert.Name,
							DashboardName: dashboardName,
							Url:           alert.Url,
							Time:          alert.NewStateDate,
						}
						alertList = append(alertList, alert)
						exist = true
					}
				} else {
					alert := serviceAlertInfo{
						PanelTitle:    panelTitle,
						State:         alert.State,
						AlertName:     alert.Name,
						DashboardName: dashboardName,
						Url:           alert.Url,
						Time:          alert.NewStateDate,
					}
					alertList = append(alertList, alert)
				}
			}
			if !exist {
				alert := serviceAlertInfo{
					PanelTitle:    panelTitle,
					State:         "ok",
					AlertName:     alert.Name,
					DashboardName: dashboardName,
					Url:           alert.Url,
					Time:          "",
				}
				alertList = append(alertList, alert)
			}
		}
	}
	sort.SliceStable(alertList, func(i, j int) bool {
		if alertList[i].State == alertList[j].State {
			return alertList[i].Time > alertList[j].Time
		}
		return alertList[i].State < alertList[j].State
	})
	return alertList
}

func InstanceServiceHealthCheck(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	serviceName := ctx.Params().Get("service_name")
	hostIp := ctx.URLParam("ip")
	if productName == "" || serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name or service_name is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	list := []map[string]interface{}{}
	infoList, err := model.HealthCheck.GetInfoByClusterIdAndProductNameAndServiceName(clusterId, productName, serviceName, hostIp)
	if err != nil {
		log.Errorf("err: %v", err.Error())
		apibase.ThrowDBModelError(err)
	} else if len(infoList) == 0 {
		return list
	}
	sort.SliceStable(infoList, func(i, j int) bool {
		return infoList[i].ExecStatus > infoList[j].ExecStatus
	})
	for _, s := range infoList {
		r := map[string]interface{}{}
		r["record_id"] = s.ID
		r["cluster_id"] = s.ClusterId
		r["product_name"] = s.ProductName
		r["service_name"] = s.ServiceName
		r["script_name"] = s.ScriptName
		r["script_name_display"] = s.ScriptNameDisplay
		r["exec_status"] = s.ExecStatus
		r["error_message"] = s.ErrorMessage
		r["auto_exec"] = s.AutoExec
		if s.StartTime.Valid == true {
			r["start_time"] = s.StartTime.Time.Format(base.TsLayout)
		} else {
			r["start_time"] = ""
		}
		if s.EndTime.Valid == true {
			r["end_time"] = s.EndTime.Time.Format(base.TsLayout)
		} else {
			r["end_time"] = ""
		}
		if s.CreateTime.Valid == true {
			r["create_time"] = s.CreateTime.Time.Format(base.TsLayout)
		} else {
			r["create_time"] = ""
		}
		if s.UpdateTime.Valid == true {
			r["update_time"] = s.UpdateTime.Time.Format(base.TsLayout)
		} else {
			r["update_time"] = ""
		}
		list = append(list, r)
	}
	return map[string]interface{}{
		"count": len(list),
		"list":  list,
	}
}

func SetHealthCheckAutoexecSwitch(ctx context.Context) apibase.Result {
	var param = struct {
		RecordId int  `json:"record_id"`
		AutoExec bool `json:"auto_exec"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if param.RecordId == 0 {
		return fmt.Errorf("record_id is empty")
	}

	err := model.HealthCheck.UpdateAutoexecById(param.RecordId, param.AutoExec)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	info, err := model.HealthCheck.GetInfoById(param.RecordId)
	if err != nil {
		return err
	}
	return map[string]interface{}{
		"record_id": info.ID,
		"auto_exec": info.AutoExec,
	}
}

func InstanceServiceHealthCheckExec(ctx context.Context) apibase.Result {
	param := struct {
		RecordId int `json:"record_id"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if param.RecordId == 0 {
		return fmt.Errorf("record_id is empty")
	}

	info, err := model.HealthCheck.GetInfoById(param.RecordId)
	if errors.Is(err, sql.ErrNoRows) {
		log.Errorf("%v", err)
		return fmt.Errorf("%s is not fond", info.ScriptNameDisplay)
	} else if err != nil {
		log.Errorf("err: %v", err.Error())
		return err
	}

	err, host := model.DeployHostList.GetHostInfoByIp(info.Ip)
	if err != nil {
		log.Errorf("get host info by ip error %v", err)
	}

	if host.IsDeleted == 0 && host.Status > 0 && time.Now().Sub(time.Time(host.UpdateDate)) < 3*time.Minute {
		var status int
		startTime := time.Now()
		status = enums.ExecStatusType.Running.Code
		if err = model.HealthCheck.UpdateHealthCheckStatus(info.ID, status, "", dbhelper.NullTime{Time: startTime, Valid: true}, dbhelper.NullTime{}); err != nil {
			log.Errorf("%v", err)
			return err
		}
		var cmd strings.Builder
		cmd.WriteString("#!/bin/sh\n")
		cmd.WriteString(fmt.Sprintf("/bin/sh %s", info.ScriptName))
		content, err := agent.AgentClient.ToExecCmdWithTimeout(info.Sid, info.AgentId, strings.TrimSpace(cmd.String()), "60s", "", "")
		if err != nil {
			content = err.Error()
			status = enums.ExecStatusType.Failed.Code
		} else {
			status = enums.ExecStatusType.Success.Code
		}
		endTime := time.Now()
		if err = model.HealthCheck.UpdateHealthCheckStatus(info.ID, status, content, dbhelper.NullTime{Time: startTime, Valid: true}, dbhelper.NullTime{Time: endTime, Valid: true}); err != nil {
			log.Errorf("%v", err)
			return err
		}
	} else {
		return fmt.Errorf("host %v is not running", host.Ip)
	}

	return map[string]interface{}{}
}
