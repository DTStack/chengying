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

package agent

const (
	AGENT_INSTALL_URI           = "/api/v1/agent/install"
	AGENT_INSTALL_SYNC_URI      = "/api/v1/agent/installSync"
	AGENT_UNINSTALL_URI         = "/api/v1/agent/%s/uninstall"
	AGENT_UNINSTALL_SYNC_URI    = "/api/v1/agent/%s/uninstallSync"
	AGENT_CANCEL_SYNC_URI       = "/api/v1/agent/cancelOperation"
	AGENT_START_URI             = "/api/v1/agent/%s/start"
	AGENT_START_SYNC_URI        = "/api/v1/agent/%s/startSync"
	AGENT_START_SYNC_PARAMS_URI = "/api/v1/agent/%s/startSyncWithParam"
	AGENT_STOP_URI              = "/api/v1/agent/%s/stop"
	AGENT_STOP_SYNC_URI         = "/api/v1/agent/%s/stopSync"
	AGENT_EXEC_SYNC_URI         = "/api/v1/sidecar/%s/execscriptSync"
	AGENT_REST_SYNC_URI         = "/api/v1/sidecar/%s/execrestSync"
	AGENT_EXEC_OFTEN_SYNC_URI   = "/api/v1/sidecar/%s/execscriptOftenSync"
	AGENT_CONFIG_SYNC_URI       = "/api/v1/agent/%s/configSync"
)
