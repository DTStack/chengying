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

// NEED REMOVE LATER!

package impl

import (
	"bytes"
	"errors"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	"easyagent/internal/server/shipper"
	"github.com/satori/go.uuid"
)

func startScriptWrapper(collectorId, agentId uuid.UUID) error {

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.START_AGENT)], collectorId, agentId)

	if err != nil {
		log.Errorf("[startScriptWrapper]startScriptWrapper error: %v", err)
		return err
	}

	err, result := shipper.GetApiShipper().StartShipperSync(uint32(seq), collectorId, agentId)

	if err != nil {
		log.Errorf("[startScriptWrapper]startScriptWrapper error: %v", err)
		return err
	}

	if result != nil && result.(*proto.Event_OperationProgress).Failed {
		log.Errorf("[startScriptWrapper]startScriptWrapper error: %v", result.(*proto.Event_OperationProgress).Message)
		return errors.New(result.(*proto.Event_OperationProgress).Message)
	}

	log.Debugf("Start agent collectorid %v, agentid %v, seq %v", collectorId.String(), agentId.String(), seq)

	return nil
}

func uninstallScriptWrapper(collectorId, agentId uuid.UUID) error {
	params := &shipper.ShellParams{Timeout: "15m"}
	params.Parameter = ""
	params.ShellScript = "#!/bin/sh\nrm -rf /opt/dtstack/easyagent/script_wrapper/"

	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.UNINSTALL_AGENT)], collectorId, agentId)

	if err != nil {
		log.Errorf("[uninstallScriptWrapper]uninstallScriptWrapper error: %v", err)
		return err
	}
	err, result := shipper.GetApiShipper().UnInstallShipperSync(uint32(seq), collectorId, agentId, params)
	if err != nil {
		log.Errorf("[uninstallScriptWrapper]uninstallScriptWrapper error: %v", err.Error())
		return err
	}
	if result != nil && result.(*proto.Event_OperationProgress).Failed {
		log.Errorf("[uninstallScriptWrapper]uninstallScriptWrapper result error: %v", result.(*proto.Event_OperationProgress).Message)
		return errors.New(result.(*proto.Event_OperationProgress).Message)
	}

	log.Debugf("[uninstallScriptWrapper]uninstallScriptWrapper agent collectorid %v, agentid %v, seq %d", collectorId, agentId, seq)
	return nil
}

func LoopInstallScriptWrapper(sidecarId, server string) {
	var err error
	agentId := uuid.NewV4()
	tryCount := 1
	for {
		err = installScriptWrapper(sidecarId, agentId, server)

		if err == nil {
			break
		}
		if tryCount > 5 {
			break
		} else {
			log.Errorf("[ERROR]LoopInstallScriptWrapper")
		}
		tryCount = tryCount + 1
		log.Debugf(">>>LoopInstallScriptWrapper")
		time.Sleep(3 * time.Second)
	}
}

func installScriptWrapper(sidecarId string, agentId uuid.UUID, server string) error {
	log.Debugf("[InstallScriptWrapper]Install agent collectorid %v", sidecarId)

	typ := "wrapper"
	script := &bytes.Buffer{}
	scriptTemplates[typ].Execute(script, map[string]interface{}{
		"easyagent_server": server,
	})

	params := &shipper.InstallParms{Timeout: "15m"}
	params.Name = "script-wrapper"
	params.BinaryPath = "/opt/dtstack/easyagent/script_wrapper/script-wrapper"
	params.ConfigurationPath = "/opt/dtstack/easyagent/script_wrapper/script-wrapper.yml"
	params.Parameter = "-c,/opt/dtstack/easyagent/script_wrapper/script-wrapper.yml,--debug"
	params.CollectorId = sidecarId
	params.InstallScript = script.String()
	params.InstallParameter = ""

	collectorId, err := uuid.FromString(params.CollectorId)

	if err != nil {
		log.Errorf("[InstallScriptWrapper]installScriptWrapper error: %v", params.CollectorId)
		return err
	}
	if params.Name == "" {
		params.Name = "UNTITLED"
	}
	seq, err := model.AgentOperation.NewOperationRecord(
		proto.ControlResponse_ControlCmd_name[int32(proto.INSTALL_AGENT)], collectorId, agentId)

	if err != nil {
		log.Errorf("[InstallScriptWrapper]installScriptWrapper error: %v", err.Error())
		return err
	}

	err, result := shipper.GetApiShipper().InstallShipperSync(uint32(seq), collectorId, agentId, params)

	if err != nil {
		log.Errorf("[InstallScriptWrapper]installScriptWrapper error: %v", err.Error())

		uninstallScriptWrapper(collectorId, agentId)
		return err
	}

	if result != nil && result.(*proto.Event_OperationProgress).Failed {
		log.Errorf("[uninstallScriptWrapper]uninstallScriptWrapper result error: %v", result.(*proto.Event_OperationProgress).Message)

		uninstallScriptWrapper(collectorId, agentId)
		return errors.New(result.(*proto.Event_OperationProgress).Message)
	}

	model.AgentList.InsertAgentRecord(collectorId, agentId, 0, "script-wrapper", "")
	log.Debugf("[InstallScriptWrapper]Install agent collectorid %v, agentid %v, seq %d", collectorId, agentId, seq)

	err = startScriptWrapper(collectorId, agentId)
	if err != nil {
		log.Errorf("[InstallScriptWrapper]startScriptWrapper error: %v", err.Error())
	}
	return nil
}
