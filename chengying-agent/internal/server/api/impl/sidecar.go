/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package impl

import (
	"encoding/base64"
	"fmt"

	apibase "easyagent/go-common/api-base"
	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	"easyagent/internal/server/shipper"
	. "easyagent/internal/server/tracy"

	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
)

const TsLayout = "2006-01-02 15:04:05"

func QuerySidecarList(ctx context.Context) apibase.Result {
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx)
	list, count := model.SidecarList.GetSidecarList(pagination)
	sidecarList := []map[string]interface{}{}
	for _, s := range list {
		id, err := uuid.FromString(s.ID)
		if err != nil {
			continue
		}
		m := map[string]interface{}{}
		m["id"] = id.String()
		m["name"] = s.Name
		m["os"] = s.OsType
		m["version"] = s.Version
		m["auto-deploy"] = s.AutoDeploy
		m["auto-update"] = s.AutoUpdate
		sidecarList = append(sidecarList, m)
	}
	return map[string]interface{}{
		"list":  sidecarList,
		"total": count,
	}
}

func NewSidecar(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	params := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	if params.Version == "" {
		paramErrs.AppendError("version", "Missing version")
	}

	if params.Name == "" {
		params.Name = "UNTITLED"
	}

	id, err := model.SidecarList.NewSidecarRecord(params.Name, params.Version)
	if err != nil {
		paramErrs.AppendError("NewSidecarRecord", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	return id
}

func GetSidecarInformation(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	id, err := uuid.FromString(ctx.Params().Get("sidecar_id"))
	if err != nil {
		paramErrs.AppendError("sid format err", err.Error())
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	info, err := model.SidecarList.GetSidecarInfo(id)
	if err != nil {
		return fmt.Errorf("No such sidecar information which id = %v: %v", id.String(), err)
	}
	returnFields := map[string]interface{}{
		"host":             info.Host,
		"local_ip":         info.LocalIp,
		"cpu_usage":        info.CpuUsage,
		"cpu_cores":        info.CpuCores,
		"cpu_serial":       info.CpuSerial,
		"os_platform":      info.OsPlatform,
		"os_version":       info.OsVersion,
		"last_update_date": info.UpdateDate.Time.Format(TsLayout),
	}
	return returnFields
}

func GetSidecarAgents(ctx context.Context) apibase.Result {
	id, err := uuid.FromString(ctx.Params().Get("sidecar_id"))
	if err != nil {
		paramErrs := apibase.NewApiParameterErrors()
		paramErrs.AppendError("sidecar_id", err)
		paramErrs.CheckAndThrowApiParameterErrors()
	}

	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx)
	agents, count := model.AgentList.GetAgentsBySidecarId(pagination, id)

	list := []map[string]interface{}{}

	for _, s := range agents {
		id, err := uuid.FromString(s.ID)

		if err != nil {
			continue
		}

		m := map[string]interface{}{}
		m["agent_id"] = id.String()
		m["deploy_date"] = s.DeployDate.Time
		list = append(list, m)
	}

	return map[string]interface{}{
		"list":  list,
		"total": count,
	}

}

func InstallSidecar(ctx context.Context) apibase.Result {
	return nil
}

func GetSidecarStatus(ctx context.Context) apibase.Result {
	return nil
}

func ChangeSidecarName(ctx context.Context) apibase.Result {
	return nil
}

func ControlSidecar(ctx context.Context) apibase.Result {
	return nil
}

func ExecScript(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	params := &shipper.ExecScriptParams{Timeout: "15m"}

	if err := ctx.ReadJSON(&params); err != nil {
		InstallProgressLog("[INSTALL] ExecScript ReadJSON err: %v", err)
		paramErrs.AppendError("$", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	SidecarIdStr := ctx.Params().Get("sidecar_id")
	SidecarId, err := uuid.FromString(SidecarIdStr)
	if err != nil {
		InstallProgressLog("[INSTALL] ExecScriptsidecar_id not uuid format: %v", SidecarIdStr)
		paramErrs.AppendError("sidecar_id", "sidecar_id not uuid format!")
	}

	if checkForSLB(ctx, SidecarId) {
		return nil
	}

	agentId := uuid.Nil
	if params.AgentId != "" {
		agentId, err = uuid.FromString(params.AgentId)
		if err != nil {
			return fmt.Errorf("[AGENT-CONTROL]ExecScriptSync param agentId is not uuid format: %v", params)
		}
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.EXEC_SCRIPT)], SidecarId, agentId)

	if err != nil {
		InstallProgressLog("[INSTALL] ExecScriptsidecar_NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError(err)
	}

	log.Debugf("exec script sidecar_id %v, agentid %v, seq %v", SidecarId.String(), agentId, seq)
	InstallProgressLog("[INSTALL] ExecScriptsidecar_sidecar_id %v, agentid %v, params %v", SidecarId.String(), agentId, *params)

	err = shipper.GetApiShipper().ExecScriptShipper(uint32(seq), SidecarId, agentId, params)

	if err != nil {
		InstallProgressLog("[INSTALL] ExecScriptsidecar_sidecar_id %v, ExecScriptShipper err: %v", SidecarId.String(), err)
		apibase.ThrowRpcHandleError(err)
	}
	return map[string]interface{}{
		"sidecar_id":    SidecarId.String(),
		"operation_seq": seq,
	}
}

func ExecScriptSync(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	params := &shipper.ExecScriptParams{Timeout: "15m"}

	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSync read json err: %v", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	ControlProgressLog("[AGENT-CONTROL] ExecScriptSync params: %v", *params)

	SidecarIdStr := ctx.Params().Get("sidecar_id")
	SidecarId, err := uuid.FromString(SidecarIdStr)
	if err != nil {
		paramErrs.AppendError("sidecar_id", "sidecar_id not uuid format!")
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSync sidecar_id not uuid format: %v", SidecarIdStr)
	}

	if checkForSLB(ctx, SidecarId) {
		return nil
	}
	agentId := uuid.Nil
	if params.AgentId != "" {
		agentId, err = uuid.FromString(params.AgentId)
		if err != nil {
			return fmt.Errorf("[AGENT-CONTROL]ExecScriptSync param agentId is not uuid format: %v", params)
		}
	}

	seq, err := model.AgentOperation.NewOperationSeqno(SidecarId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSync err: %v", err)
		apibase.ThrowDBModelError(err)
	}

	log.Debugf("exec script sidecar_id %v, agentid %v, seq %v", SidecarId.String(), agentId, seq)
	ControlProgressLog("[AGENT-CONTROL] ExecScriptSync exec script sidecar_id %v, agentid %v, seq %v", SidecarId.String(), agentId, seq)

	err, result := shipper.GetApiShipper().ExecScriptShipperSync(uint32(seq), SidecarId, agentId, params)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSync ExecScriptShipperSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}
	ControlProgressLog("[AGENT-CONTROL] ExecScriptSync ExecScriptShipperSync result: %v", result)

	return map[string]interface{}{
		"sidecar_id":    SidecarId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func ExecScriptSyncBase64(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	params := &shipper.ExecScriptParams{Timeout: "15m"}

	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSyncBase64 read json err: %v", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	ControlProgressLog("[AGENT-CONTROL] ExecScriptSyncBase64 params: %v", *params)

	SidecarIdStr := ctx.Params().Get("sidecar_id")
	SidecarId, err := uuid.FromString(SidecarIdStr)
	if err != nil {
		paramErrs.AppendError("sidecar_id", "sidecar_id not uuid format!")
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSyncBase64 sidecar_id not uuid format: %v", SidecarIdStr)
	}

	if checkForSLB(ctx, SidecarId) {
		return nil
	}

	agentId := uuid.Nil
	if params.AgentId != "" {
		agentId, err = uuid.FromString(params.AgentId)
		if err != nil {
			return fmt.Errorf("[AGENT-CONTROL]ExecScriptSync param agentId is not uuid format: %v", params)
		}
	}

	seq, err := model.AgentOperation.NewOperationSeqno(SidecarId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSyncBase64 err: %v", err)
		apibase.ThrowDBModelError(err)
	}

	log.Debugf("exec script sidecar_id %v, agentid %v, seq %v", SidecarId.String(), agentId, seq)
	ControlProgressLog("[AGENT-CONTROL] ExecScriptSyncBase64 exec script sidecar_id %v, agentid %v, seq %v", SidecarId.String(), agentId, seq)

	err, result := shipper.GetApiShipper().ExecScriptShipperSync(uint32(seq), SidecarId, agentId, params)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] ExecScriptSyncBase64 ExecScriptShipperSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}

	ControlProgressLog("[AGENT-CONTROL] ExecScriptSyncBase64 ExecScriptShipperSync result: %v", result)

	if result == nil {
		return fmt.Errorf("ExecScriptSyncBase64 result is nil: %s", SidecarId.String())
	}
	resultExec := result.(*proto.Event_ExecScriptResponse)
	resultExec.Response = base64.StdEncoding.EncodeToString([]byte(resultExec.Response))

	return map[string]interface{}{
		"sidecar_id":    SidecarId.String(),
		"operation_seq": seq,
		"result":        resultExec,
	}
}

func ExecRestSync(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	params := &shipper.ExecRestParams{Timeout: "1m"}

	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		ControlProgressLog("[AGENT-CONTROL] ExecRestSync read json err: %v", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	ControlProgressLog("[AGENT-CONTROL] ExecRestSync params: %v", *params)

	SidecarIdStr := ctx.Params().Get("sidecar_id")
	SidecarId, err := uuid.FromString(SidecarIdStr)
	if err != nil {
		paramErrs.AppendError("sidecar_id", "sidecar_id not uuid format!")
		ControlProgressLog("[AGENT-CONTROL] ExecRestSync sidecar_id not uuid format: %v", SidecarIdStr)
	}

	if checkForSLB(ctx, SidecarId) {
		return nil
	}

	agentId := uuid.Nil
	if params.AgentId != "" {
		agentId, err = uuid.FromString(params.AgentId)
		if err != nil {
			return fmt.Errorf("[AGENT-CONTROL]ExecRestSync param agentId is not uuid format: %v", params)
		}
	}

	seq, err := model.AgentOperation.NewOperationSeqno(SidecarId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] ExecRestSync err: %v", err)
		apibase.ThrowDBModelError(err)
	}

	log.Debugf("exec script sidecar_id %v, agentid %v, seq %v", SidecarId.String(), agentId, seq)
	ControlProgressLog("[AGENT-CONTROL] ExecRestSync exec script sidecar_id %v, agentid %v, seq %v", SidecarId.String(), agentId, seq)

	err, result := shipper.GetApiShipper().ExecRestShipperSync(uint32(seq), SidecarId, agentId, params)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] ExecRestSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}

	ControlProgressLog("[AGENT-CONTROL] ExecRestSync result: %v", result)

	if result == nil {
		return fmt.Errorf("ExecRestSync result is nil: %s", SidecarId.String())
	}
	return map[string]interface{}{
		"sidecar_id":    SidecarId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func GetExecScriptProgress(ctx context.Context) apibase.Result {
	return nil
}
