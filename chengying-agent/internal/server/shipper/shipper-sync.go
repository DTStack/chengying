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

package shipper

import (
	"context"
	apibase "easyagent/go-common/api-base"
	"strings"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"easyagent/internal/server/rpc"

	uuid "github.com/satori/go.uuid"
)

const (
	defaultTimeout = 1 * time.Minute
)

func (shipper *ApiShipper) shipperCoreSync(collectorId uuid.UUID, response *proto.ControlResponse, timeout time.Duration) (error, interface{}) {
	log.Infof("ShipperSync RPC response:%v", response)
	ctx, cancel := context.WithTimeout(context.Background(), timeout+5*time.Second)
	defer cancel()
	result, err := rpc.SidecarClient.SendControlSync(ctx, collectorId, response)
	if err != nil {
		log.Errorf("SendControlSync err %v", err)
	}
	return err, result
}

func (shipper *ApiShipper) InstallShipperSync(seq uint32, collectorId, agentId uuid.UUID, parms *InstallParms) (error, interface{}) {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("installAgent timeout", "wrong timeout format")
		return paramErrs, nil
	}
	var healthPeriod, healthStartPeriod, healthTimeout time.Duration
	if parms.HealthShell != "" {
		if healthPeriod, err = time.ParseDuration(parms.HealthPeriod); err != nil {
			paramErrs.AppendError("installAgent healthPeriod", "wrong healthPeriod format")
			return paramErrs, nil
		}
		if parms.HealthStartPeriod != "" {
			if healthStartPeriod, err = time.ParseDuration(parms.HealthStartPeriod); err != nil {
				paramErrs.AppendError("installAgent healthStartPeriod", "wrong healthPeriod format")
				return paramErrs, nil
			}
		}
		if parms.HealthTimeout != "" {
			if healthTimeout, err = time.ParseDuration(parms.HealthTimeout); err != nil {
				paramErrs.AppendError("installAgent healthTimeout", "wrong healthPeriod format")
				return paramErrs, nil
			}
		}

		if healthPeriod < time.Second {
			paramErrs.AppendError("installAgent healthPeriod", "healthPeriod less than 1 sec")
			return paramErrs, nil
		}
		if healthTimeout <= 0 || healthTimeout > healthPeriod {
			healthTimeout = healthPeriod
		}
		if parms.HealthRetries <= 0 {
			parms.HealthRetries = 1
		}
	}
	installParameters := parms.InstallParameter

	param := []string{}
	installParam := []string{}

	if len(parameters) > 0 {
		param = strings.Split(parameters, ",")
	}
	if len(installParameters) > 0 {
		installParam = strings.Split(installParameters, ",")
	}

	response := &proto.ControlResponse{
		Cmd:   proto.INSTALL_AGENT,
		Seqno: seq,
		Options: &proto.ControlResponse_InstallAgentOptions_{
			InstallAgentOptions: &proto.ControlResponse_InstallAgentOptions{
				AgentId:           agentId.Bytes(),
				ConfigurationPath: parms.ConfigurationPath,
				BinaryPath:        parms.BinaryPath,
				Parameter:         param,
				Name:              parms.Name,
				InstallParameter:  installParam,
				InstallScript:     parms.InstallScript,
				HealthShell:       parms.HealthShell,
				HealthPeriod:      healthPeriod,
				HealthStartPeriod: healthStartPeriod,
				HealthTimeout:     healthTimeout,
				HealthRetries:     parms.HealthRetries,
				Timeout:           timeout,
				Workdir:           parms.WorkDir,
				RunUser:           parms.Runuser,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, timeout)
}

func (shipper *ApiShipper) CancelShipperSync(seq uint32, collectorId, agentId uuid.UUID, timeout time.Duration) (error, interface{}) {
	response := &proto.ControlResponse{
		Cmd:   proto.CANCEL_AGENTS,
		Seqno: seq,
		Options: &proto.ControlResponse_CancelOptions_{
			CancelOptions: &proto.ControlResponse_CancelOptions{
				AgentId: agentId.Bytes(),
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, timeout)
}

func (shipper *ApiShipper) UnInstallShipperSync(seq uint32, collectorId, agentId uuid.UUID, parms *ShellParams) (error, interface{}) {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("uninstallAgent timeout", "wrong timeout format")
		return paramErrs, nil
	}
	param := []string{}

	if len(parameters) > 0 {
		param = strings.Split(parameters, ",")
	}
	response := &proto.ControlResponse{
		Cmd:   proto.UNINSTALL_AGENT,
		Seqno: seq,
		Options: &proto.ControlResponse_UninstallAgentOptions_{
			UninstallAgentOptions: &proto.ControlResponse_UninstallAgentOptions{
				AgentId:         agentId.Bytes(),
				Parameter:       param,
				UninstallScript: parms.ShellScript,
				Timeout:         timeout,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, timeout)
}

func (shipper *ApiShipper) StartShipperSync(seq uint32, collectorId, agentId uuid.UUID) (error, interface{}) {

	response := &proto.ControlResponse{
		Cmd:   proto.START_AGENT,
		Seqno: seq,
		Options: &proto.ControlResponse_StartAgentOptions_{
			StartAgentOptions: &proto.ControlResponse_StartAgentOptions{
				AgentId:  agentId.Bytes(),
				CpuLimit: 0,
				MemLimit: 0,
				NetLimit: 0,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, defaultTimeout)
}

func (shipper *ApiShipper) StartShipperSyncWithParam(seq uint32, collectorId, agentId uuid.UUID, params StartParams) (error, interface{}) {

	response := &proto.ControlResponse{
		Cmd:   proto.START_AGENT,
		Seqno: seq,
		Options: &proto.ControlResponse_StartAgentOptions_{
			StartAgentOptions: &proto.ControlResponse_StartAgentOptions{
				AgentId:     agentId.Bytes(),
				CpuLimit:    params.CpuLimit,
				MemLimit:    params.MemLimit,
				NetLimit:    params.NetLimit,
				Environment: params.Environment,
			},
		},
	}
	return shipper.shipperCoreSync(collectorId, response, defaultTimeout)
}

func (shipper *ApiShipper) StopShipperSync(seq uint32, collectorId, agentId uuid.UUID, stopAgentOptionsType int) (error, interface{}) {

	optionsTypeEnum := proto.STOP_UNRECOVER
	if stopAgentOptionsType == int(proto.STOP_RECOVER) {
		optionsTypeEnum = proto.STOP_RECOVER
	}

	response := &proto.ControlResponse{
		Cmd:   proto.STOP_AGENT,
		Seqno: seq,
		Options: &proto.ControlResponse_StopAgentOptions_{
			StopAgentOptions: &proto.ControlResponse_StopAgentOptions{
				AgentId:              agentId.Bytes(),
				StopAgentOptionsType: optionsTypeEnum,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, defaultTimeout)
}

func (shipper *ApiShipper) UpdateShipperSync(seq uint32, collectorId, agentId uuid.UUID, parms *ShellParams) (error, interface{}) {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("updateAgent timeout", "wrong timeout format")
		return paramErrs, nil
	}
	param := []string{}

	if len(parameters) > 0 {
		param = strings.Split(parameters, ",")
	}
	response := &proto.ControlResponse{
		Cmd:   proto.UPDATE_AGENT,
		Seqno: seq,
		Options: &proto.ControlResponse_UpdateAgentOptions_{
			UpdateAgentOptions: &proto.ControlResponse_UpdateAgentOptions{
				AgentId:      agentId.Bytes(),
				Parameter:    param,
				UpdateScript: parms.ShellScript,
				Timeout:      timeout,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, timeout)
}

func (shipper *ApiShipper) UpdateAgentConfigShipperSync(seq uint32, collectorId, agentId uuid.UUID, param ConfigParams) (error, interface{}) {
	response := &proto.ControlResponse{
		Cmd:   proto.UPDATE_AGENT_CONFIG,
		Seqno: seq,
		Options: &proto.ControlResponse_UpdateAgentConfigOptions_{
			UpdateAgentConfigOptions: &proto.ControlResponse_UpdateAgentConfigOptions{
				AgentId:       agentId.Bytes(),
				ConfigContent: param.ConfigContent,
				ConfigPath:    param.ConfigPath,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, defaultTimeout)
}

func (shipper *ApiShipper) ExecScriptShipperSync(seq uint32, collectorId, agentId uuid.UUID, parms *ExecScriptParams) (error, interface{}) {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("execScript timeout", "wrong timeout format")
		return paramErrs, nil
	}
	var execAgentId []byte
	if agentId != uuid.Nil {
		execAgentId = agentId.Bytes()
	}
	param := []string{}

	if len(parameters) > 0 {
		param = strings.Split(parameters, ",")
	}
	response := &proto.ControlResponse{
		Cmd:   proto.EXEC_SCRIPT,
		Seqno: seq,
		Options: &proto.ControlResponse_ExecScriptOptions_{
			ExecScriptOptions: &proto.ControlResponse_ExecScriptOptions{
				ExecScript: parms.ExecScript,
				Parameter:  param,
				Timeout:    timeout,
				AgentId:    execAgentId,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, timeout)
}

func (shipper *ApiShipper) ExecRestShipperSync(seq uint32, collectorId, agentId uuid.UUID, parms *ExecRestParams) (error, interface{}) {
	paramErrs := apibase.NewApiParameterErrors()
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("execRest timeout", "wrong timeout format")
		return paramErrs, nil
	}
	var execAgentId []byte
	if agentId != uuid.Nil {
		execAgentId = agentId.Bytes()
	}

	response := &proto.ControlResponse{
		Cmd:   proto.EXEC_REST,
		Seqno: seq,
		Options: &proto.ControlResponse_ExecRestOptions_{
			ExecRestOptions: &proto.ControlResponse_ExecRestOptions{
				Method:  parms.Method,
				Path:    parms.Path,
				Query:   parms.Query,
				Body:    parms.Body,
				Timeout: timeout,
				AgentId: execAgentId,
			},
		},
	}

	return shipper.shipperCoreSync(collectorId, response, timeout)
}
