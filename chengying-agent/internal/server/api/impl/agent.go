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
	apibase "easyagent/go-common/api-base"
	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	"easyagent/internal/server/shipper"

	"fmt"
	"time"

	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
)

func QueryAgentList(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentInfo(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentName(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentSid(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	return map[string]interface{}{
		"sidecar_id": collectorId.String(),
	}
}

func SetAgentName(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentStatus(ctx context.Context) apibase.Result {

	return nil
}

func StopAgent(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	stopAgentOptionsType, err := ctx.URLParamInt("stop_agent_options_type")
	if err != nil {
		stopAgentOptionsType = int(proto.STOP_UNRECOVER)
	}

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.STOP_AGENT)], collectorId, agentId)

	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Stop agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err = shipper.GetApiShipper().StopShipper(uint32(seq), collectorId, agentId, stopAgentOptionsType)

	if err != nil {
		apibase.ThrowRpcHandleError(err)
	}

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
	}
}

func GetAgentStoppingProgress(ctx context.Context) apibase.Result {

	return nil
}

func StartAgent(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.START_AGENT)], collectorId, agentId)

	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Start agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err = shipper.GetApiShipper().StartShipper(uint32(seq), collectorId, agentId)

	if err != nil {
		apibase.ThrowRpcHandleError(err)
	}

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
	}
}

func StartAgentWithParam(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}
	params := &shipper.StartParams{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	paramErrs.CheckAndThrowApiParameterErrors()

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.START_AGENT)], collectorId, agentId)

	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Start agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err = shipper.GetApiShipper().StartShipperWithParam(uint32(seq), collectorId, agentId, *params)

	if err != nil {
		apibase.ThrowRpcHandleError(err)
	}

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
	}
}

func GetAgentStartingProgress(ctx context.Context) apibase.Result {

	return nil
}

func RestartAgent(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentRestartingProgress(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentConfig(ctx context.Context) apibase.Result {

	return nil
}

func UpdateAgentConfig(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}

	params := &shipper.ConfigParams{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	configContent := params.ConfigContent

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.UPDATE_AGENT_CONFIG)], collectorId, agentId)

	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Update agent config collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err = shipper.GetApiShipper().UpdateAgentConfigShipper(uint32(seq), collectorId, agentId, configContent)

	if err != nil {
		apibase.ThrowRpcHandleError(err)
	}

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
	}
}

func GetAgentConfigUpdatingProgress(ctx context.Context) apibase.Result {

	return nil
}

func ReinstallAgent(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentReinstallProgress(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentInstallProgress(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentUpdateInstallProgress(ctx context.Context) apibase.Result {

	return nil
}

func InstallAgent(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	params := &shipper.InstallParms{
		HealthPeriod:  "1m",
		HealthRetries: 1,
		Timeout:       "15m",
	}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	collectorId, err := uuid.FromString(params.CollectorId)

	if err != nil {
		paramErrs.AppendError("collectorId", "collectorId not uuid format")
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	if params.Name == "" {
		params.Name = "UNTITLED"
	}

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	agentId := model.AgentList.NewAgentRecord(collectorId, 0, params.Name, "")
	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.INSTALL_AGENT)], collectorId, agentId)

	if err != nil {
		apibase.ThrowDBModelError()
	}

	log.Debugf("Install agent collectorid %v, agentid %v, seq %d", collectorId, agentId, seq)

	err = shipper.GetApiShipper().InstallShipper(uint32(seq), collectorId, agentId, params)

	if err != nil {
		apibase.ThrowRpcHandleError(err)
	}

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
	}
}

func UpdateAgent(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}

	params := &shipper.ShellParams{Timeout: "15m"}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.UPDATE_AGENT)], collectorId, agentId)

	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("Update agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err = shipper.GetApiShipper().UpdateShipper(uint32(seq), collectorId, agentId, params)

	if err != nil {
		apibase.ThrowRpcHandleError(err)
	}

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
	}
}

func UninstallAgent(ctx context.Context) apibase.Result {

	paramErrs := apibase.NewApiParameterErrors()

	agentIdStr := ctx.Params().Get("agent_id")

	agentId, err := uuid.FromString(agentIdStr)

	if err != nil {
		paramErrs.AppendError("agent_id", "agent_id not uuid format")
	}

	params := &shipper.ShellParams{Timeout: "15m"}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}

	collectorId := model.AgentList.GetAgentSidecarId(agentId)

	if checkForSLB(ctx, collectorId) {
		return nil
	}

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.UNINSTALL_AGENT)], collectorId, agentId)

	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	log.Debugf("UnInstall agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	err = shipper.GetApiShipper().UnInstallShipper(uint32(seq), collectorId, agentId, params)

	if err != nil {
		apibase.ThrowRpcHandleError(err)
	}

	return map[string]interface{}{
		"agent_id":      agentId.String(),
		"operation_seq": seq,
	}
}

func GetAgentUninstallProgress(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentProgress(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentAutoupdateConfig(ctx context.Context) apibase.Result {

	return nil
}

func SetAgentAutoupdateConfig(ctx context.Context) apibase.Result {

	return nil
}

func GetAgentEventHistory(ctx context.Context) apibase.Result {

	return nil
}

func CancelOperation(ctx context.Context) apibase.Result {
	params := &shipper.CancelParams{}
	if err := ctx.ReadJSON(&params); err != nil {
		return fmt.Errorf("param error: %v", err)
	}

	results := map[string]interface{}{}
	var seqs []int64
	for sidecarId, agentIds := range params.Agents {
		for _, agentId := range agentIds {
			sidecarUuid, _ := uuid.FromString(sidecarId)
			agentUuid, _ := uuid.FromString(agentId)

			seq, err := model.AgentOperation.NewOperationRecord(
				proto.ControlResponse_ControlCmd_name[int32(proto.CANCEL_AGENTS)], sidecarUuid, agentUuid)

			err, result := shipper.GetApiShipper().CancelShipperSync(uint32(seq), sidecarUuid, agentUuid, 5*time.Second)
			if err != nil {
				results[agentId] = err.Error()
			} else if result == nil {
				results[agentId] = "nil result"
			} else {
				results[agentId] = "ok"
			}
			seqs = append(seqs, seq)
		}
	}

	return map[string]interface{}{
		"agents":        params.Agents,
		"operation_seq": seqs,
		"result":        results,
	}
}
