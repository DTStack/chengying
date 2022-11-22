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
	"fmt"
	"time"

	"database/sql"

	apibase "easyagent/go-common/api-base"
	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	"easyagent/internal/server/shipper"
	. "easyagent/internal/server/tracy"

	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
)

func InstallAgentSync(ctx context.Context) apibase.Result {
	InstallProgressLog("[INSTALL] InstallAgentSync ...%v", "")
	paramErrs := apibase.NewApiParameterErrors()
	params := &shipper.InstallParms{
		HealthPeriod:  "1m",
		HealthRetries: 1,
		Timeout:       "15m",
	}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		InstallProgressLog("[INSTALL] InstallAgentSync param err: %v", err)
	}
	collectorId, err := uuid.FromString(params.CollectorId)

	if err != nil {
		paramErrs.AppendError("collectorId", "collectorId not uuid format")
		InstallProgressLog("[INSTALL] InstallAgentSync param err: collectorId not uuid format")
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	if params.Name == "" {
		params.Name = "UNTITLED"
	}

	agentId := uuid.NewV4()
	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.INSTALL_AGENT)], collectorId, agentId)

	if err != nil {
		InstallProgressLog("[INSTALL] InstallAgentSync NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError()
	}

	log.Debugf("Install agent sync collectorid %v, agentid %v, seq %d", collectorId, agentId, seq)
	InstallProgressLog("[INSTALL] Install agent sync collectorid %v, agentid %v, seq %d, params: %v", collectorId, agentId, seq, *params)

	err, result := shipper.GetApiShipper().InstallShipperSync(uint32(seq), collectorId, agentId, params)

	if err != nil {
		log.Errorf("[InstallAgentSync]installSync error: %v", err.Error())
		InstallProgressLog("[INSTALL] InstallAgentSync InstallShipperSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}
	if result == nil {
		InstallProgressLog("[INSTALL] InstallAgentSync failed with messsage: %v", "result is nil")
		return fmt.Errorf("InstallAgentSync result is nil: %s", collectorId.String())
	}
	model.AgentList.InsertAgentRecord(collectorId, agentId, 0, params.Name, "")

	InstallProgressLog("[INSTALL] InstallAgentSync success with result: %v", result)

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func UninstallAgentSync(ctx context.Context) apibase.Result {
	InstallProgressLog("[INSTALL] UninstallAgentSync ...%v", "")
	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
		InstallProgressLog("[INSTALL] UninstallAgentSync param err: agent_id not uuid format")
	}

	params := &shipper.ShellParams{Timeout: "15m"}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		InstallProgressLog("[INSTALL] UninstallAgentSync param err: %v", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	InstallProgressLog("[INSTALL] UninstallAgentSync params: %v", *params)
	collectorId, err := model.AgentList.CheckAgentId(agentId)

	if err == sql.ErrNoRows {
		InstallProgressLog("[INSTALL] UninstallAgentSync result, agent already uninstall: %v", agentId.String())
		return map[string]interface{}{
			"agent_id":      agentId.String(),
			"operation_seq": -1,
			"result":        map[string]interface{}{"message": "agent already uninstall"},
		}
	}

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.UNINSTALL_AGENT)], collectorId, agentId)

	if err != nil {
		InstallProgressLog("[INSTALL] UninstallAgentSync NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("UnInstall agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	InstallProgressLog("[INSTALL] UnInstall agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	err, result := shipper.GetApiShipper().UnInstallShipperSync(uint32(seq), collectorId, agentId, params)

	if err != nil {
		InstallProgressLog("[INSTALL] UninstallAgentSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}

	if result == nil {
		InstallProgressLog("[INSTALL] UninstallAgentSync failed with messsage: %v", "result is nil")
		return fmt.Errorf("UninstallAgentSync result is nil: %s", collectorId.String())
	}
	resultUnInstall := result.(*proto.Event_OperationProgress)

	if resultUnInstall.Failed {
		InstallProgressLog("[INSTALL] UninstallAgentSync failed with messsage: %v", resultUnInstall.Message)
	} else {
		model.AgentList.DeleteByagentId(agentId.String())
		InstallProgressLog("[INSTALL] UninstallAgentSync success with result: %v", result)
	}
	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func StopAgentSync(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] StopAgentSync agent_id not uuid format: %v", agentIdStr)
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}

	stopAgentOptionsType, err := ctx.URLParamInt("stop_agent_options_type")
	if err != nil {
		stopAgentOptionsType = int(proto.STOP_UNRECOVER)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.STOP_AGENT)], collectorId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] StopAgentSync NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Stop agent sync collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	ControlProgressLog("[AGENT-CONTROL] Stop agent sync collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err, result := shipper.GetApiShipper().StopShipperSync(uint32(seq), collectorId, agentId, stopAgentOptionsType)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] StopShipperSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}
	ControlProgressLog("[AGENT-CONTROL] StopShipperSync result: %v", result)

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func StartAgentSync(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
		ControlProgressLog("[AGENT-CONTROL] StartAgentSync agent_id not uuid format: %v", agentIdStr)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.START_AGENT)], collectorId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] StartAgentSync NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Start agent sync  collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	ControlProgressLog("[AGENT-CONTROL] start agent sync collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err, result := shipper.GetApiShipper().StartShipperSync(uint32(seq), collectorId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] StartAgentSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}
	log.Debugf("Start result: %v", result)
	ControlProgressLog("[AGENT-CONTROL] StartAgentSync result: %v", result)

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func StartAgentSyncWithParam(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
		ControlProgressLog("[AGENT-CONTROL] StartAgentSyncWithParam agent_id not uuid format: %v", agentIdStr)
	}

	params := &shipper.StartParams{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.START_AGENT)], collectorId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] StartAgentSyncWithParam NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Start agent sync with param collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	ControlProgressLog("[AGENT-CONTROL] start agent sync with param collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err, result := shipper.GetApiShipper().StartShipperSyncWithParam(uint32(seq), collectorId, agentId, *params)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] StartAgentSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}
	log.Debugf("Start result: %v", result)
	ControlProgressLog("[AGENT-CONTROL] StartAgentSyncWithParam result: %v", result)

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func RestartAgentSync(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)
	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
		ControlProgressLog("[AGENT-CONTROL] RestartAgentSync agent_id not uuid format: %v", agentIdStr)
	}

	stopAgentOptionsType, err := ctx.URLParamInt("stop_agent_options_type")
	if err != nil {
		stopAgentOptionsType = int(proto.STOP_UNRECOVER)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.STOP_AGENT)], collectorId, agentId)
	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] RestartAgentSync NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("restart agent stop sync  collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	ControlProgressLog("[AGENT-CONTROL] restart agent stop sync  collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err, result := shipper.GetApiShipper().StopShipperSync(uint32(seq), collectorId, agentId, stopAgentOptionsType)
	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] RestartAgentSync  StopShipperSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}

	log.Debugf("restart stop result: %v", result)
	ControlProgressLog("[AGENT-CONTROL] RestartAgentSync stop result: %v", result)

	time.Sleep(2 * time.Second)

	seq, err = model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.START_AGENT)], collectorId, agentId)
	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] RestartAgentSync NewOperationRecord err: %v", err)
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("restart agent start sync  collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	ControlProgressLog("[AGENT-CONTROL] restart agent start sync  collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err, result = shipper.GetApiShipper().StartShipperSync(uint32(seq), collectorId, agentId)
	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] RestartAgentSync  StartShipperSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}

	log.Debugf("restart start result: %v", result)
	ControlProgressLog("[AGENT-CONTROL] RestartAgentSync start result: %v", result)

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}

func UpdateAgentConfigSync(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
		ControlProgressLog("[AGENT-CONTROL] UpdateAgentConfigSync agent_id not uuid format: %v", agentIdStr)
	}

	params := &shipper.ConfigParams{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		ControlProgressLog("[AGENT-CONTROL] UpdateAgentConfigSync readjson err: %v", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	ControlProgressLog("[AGENT-CONTROL] UpdateAgentConfigSync params: %v", *params)

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.UPDATE_AGENT_CONFIG)], collectorId, agentId)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] UpdateAgentConfigSync NewOperationRecord: %v", err)
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Update agent config sync collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)
	ControlProgressLog("[AGENT-CONTROL] UpdateAgentConfigSync collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err, result := shipper.GetApiShipper().UpdateAgentConfigShipperSync(uint32(seq), collectorId, agentId, *params)

	if err != nil {
		ControlProgressLog("[AGENT-CONTROL] UpdateAgentConfigSync UpdateAgentConfigShipperSync err: %v", err)
		apibase.ThrowRpcHandleError(err)
	}
	log.Debugf("Update agent config result: %v", result)
	ControlProgressLog("[AGENT-CONTROL] UpdateAgentConfigSync Update agent config result: %v", result)

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
		"result":        result,
	}
}
