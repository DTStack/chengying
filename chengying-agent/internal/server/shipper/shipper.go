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
	"strings"
	"time"

	apibase "easyagent/go-common/api-base"
	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"easyagent/internal/server/rpc"

	uuid "github.com/satori/go.uuid"
)

var apiShipper = &ApiShipper{}

type ApiShipper struct {
}

type InstallParms struct {
	CollectorId       string `json:"collectorId"`
	ConfigurationPath string `json:"configurationPath"`
	BinaryPath        string `json:"binaryPath"`
	Name              string `json:"name"`
	Parameter         string `json:"parameter"`
	InstallScript     string `json:"installScript"`
	InstallParameter  string `json:"installParameter"`
	HealthShell       string `json:"healthShell"`
	HealthPeriod      string `json:"healthPeriod"`
	HealthStartPeriod string `json:"healthStartPeriod"`
	HealthTimeout     string `json:"healthTimeout"`
	HealthRetries     uint64 `json:"healthRetries"`
	WorkDir           string `json:"workDir"`
	Timeout           string `json:"timeout"`
	Runuser           string `json:"runUser"`
}

type StartParams struct {
	AgentId     string            `json:"agentId"`
	CpuLimit    float32           `json:"cpuLimit"`
	MemLimit    uint64            `json:"memLimit"`
	NetLimit    uint64            `json:"netLimit"`
	Environment map[string]string `json:"environment"`
}

type CancelParams struct {
	Agents map[string][]string `json:"agents"`
}

type ShellParams struct {
	Parameter   string `json:"parameter"`
	ShellScript string `json:"shellScript"`
	Timeout     string `json:"timeout"`
}

type ConfigParams struct {
	ConfigContent string `json:"config_content"`
	ConfigPath    string `json:"config_path"`
}

type ExecScriptParams struct {
	ExecScript string `json:"execScript"`
	Parameter  string `json:"parameter"`
	Timeout    string `json:"timeout"`
	AgentId    string `json:"agentId"`
}

type ExecRestParams struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Query   string `json:"query"`
	Body    []byte `json:"body,omitempty"`
	Timeout string `json:"timeout"`
	AgentId string `json:"agentId"`
}

func GetApiShipper() *ApiShipper {
	return apiShipper
}

func (shipper *ApiShipper) shipperCore(collectorId uuid.UUID, response *proto.ControlResponse) error {
	log.Infof("Shipper RPC response:%v", response)
	if err := rpc.SidecarClient.SendControl(context.Background(), collectorId, response); err != nil {
		log.Errorf("SendControl error: %v", err)
		return err
	}
	return nil
}

func (shipper *ApiShipper) InstallShipper(seq uint32, collectorId, agentId uuid.UUID, parms *InstallParms) error {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("installAgent timeout", "wrong timeout format")
		return paramErrs
	}
	var healthPeriod, healthStartPeriod, healthTimeout time.Duration
	if parms.HealthShell != "" {
		if healthPeriod, err = time.ParseDuration(parms.HealthPeriod); err != nil {
			paramErrs.AppendError("installAgent healthPeriod", "wrong healthPeriod format")
			return paramErrs
		}
		if parms.HealthStartPeriod != "" {
			if healthStartPeriod, err = time.ParseDuration(parms.HealthStartPeriod); err != nil {
				paramErrs.AppendError("installAgent healthStartPeriod", "wrong healthStartPeriod format")
				return paramErrs
			}
		}
		if parms.HealthTimeout != "" {
			if healthTimeout, err = time.ParseDuration(parms.HealthTimeout); err != nil {
				paramErrs.AppendError("installAgent healthTimeout", "wrong healthTimeout format")
				return paramErrs
			}
		}

		if healthPeriod < time.Second {
			paramErrs.AppendError("installAgent healthPeriod", "healthPeriod less than 1 sec")
			return paramErrs
		}
		if healthTimeout <= 0 || healthTimeout > healthPeriod {
			healthTimeout = healthPeriod
		}
		if parms.HealthRetries <= 0 {
			parms.HealthRetries = 1
		}
	}
	param := []string{}
	if len(parameters) > 0 {
		param = strings.Split(parameters, ",")
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
				InstallParameter:  []string{},
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

	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) UnInstallShipper(seq uint32, collectorId, agentId uuid.UUID, parms *ShellParams) error {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("uninstallAgent timeout", "wrong timeout format")
		return paramErrs
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
	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) StartShipper(seq uint32, collectorId, agentId uuid.UUID) error {

	response := &proto.ControlResponse{
		Cmd:   proto.START_AGENT,
		Seqno: seq,
		Options: &proto.ControlResponse_StartAgentOptions_{
			StartAgentOptions: &proto.ControlResponse_StartAgentOptions{
				AgentId:     agentId.Bytes(),
				CpuLimit:    0,
				MemLimit:    0,
				NetLimit:    0,
				Environment: make(map[string]string),
			},
		},
	}
	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) StartShipperWithParam(seq uint32, collectorId, agentId uuid.UUID, params StartParams) error {

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
	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) StopShipper(seq uint32, collectorId, agentId uuid.UUID, stopAgentOptionsType int) error {

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
	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) UpdateShipper(seq uint32, collectorId, agentId uuid.UUID, parms *ShellParams) error {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("updateAgent timeout", "wrong timeout format")
		return paramErrs
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
	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) UpdateAgentConfigShipper(seq uint32, collectorId, agentId uuid.UUID, configContent string) error {
	response := &proto.ControlResponse{
		Cmd:   proto.UPDATE_AGENT_CONFIG,
		Seqno: seq,
		Options: &proto.ControlResponse_UpdateAgentConfigOptions_{
			UpdateAgentConfigOptions: &proto.ControlResponse_UpdateAgentConfigOptions{
				AgentId:       agentId.Bytes(),
				ConfigContent: configContent,
			},
		},
	}
	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) ExecScriptShipper(seq uint32, collectorId, agentId uuid.UUID, parms *ExecScriptParams) error {
	paramErrs := apibase.NewApiParameterErrors()
	parameters := parms.Parameter
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("execScript timeout", "wrong timeout format")
		return paramErrs
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
	return shipper.shipperCore(collectorId, response)
}

func (shipper *ApiShipper) ExecRestShipper(seq uint32, collectorId, agentId uuid.UUID, parms *ExecRestParams) error {
	paramErrs := apibase.NewApiParameterErrors()
	timeout, err := time.ParseDuration(parms.Timeout)
	if err != nil {
		paramErrs.AppendError("execRest timeout", "wrong timeout format")
		return paramErrs
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
	return shipper.shipperCore(collectorId, response)
}
