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
	HealthRetries     int    `json:"healthRetries"`
	WorkDir           string `json:"workDir"`
	RunUser           string `json:"runUser"`
	Timeout           string `json:"timeout,omitempty"`
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
}

type ConfigParams struct {
	ConfigContent string `json:"config_content"`
	ConfigPath    string `json:"config_path"`
	WorkDir       string `json:"work_dir"`
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
	Body    []byte `json:"body"`
	Timeout string `json:"timeout"`
	AgentId string `json:"agentId"`
}

type EasyagentServerResponse struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type HealthCheck struct {
	Sid     string `json:"sid"`
	AgentId string `json:"agentId"`
	Failed  bool   `json:"failed"`
	Seqno   uint32 `json:"seqno,omitempty"`
}

type AgentError struct {
	Sid            string `json:"sid"`
	AgentId        string `json:"agentId"`
	ErrStr         string `json:"errstr"`
	Seqno          uint32 `json:"seqno,omitempty"`
	LastUpdateDate string `json:"last_update_date"`
}

type AgentPerformance struct {
	Sid            string  `json:"sid"`
	AgentId        string  `json:"agentId"`
	CpuUsage       float32 `json:"cpuUsage"`
	Memory         uint64  `json:"memory"`
	Cmd            string  `json:"cmd"`
	BytesSent      uint64  `json:"bytesSent"`
	BytesRecv      uint64  `json:"bytesRecv"`
	LastUpdateDate string  `json:"last_update_date"`
}

const (
	IS_ALREADY_RUNNING          = "is already running"
	RETRY_AGENT_IS_RUNNING_TASK = "agent is running task, droped"
)

const (
	AGENT_STOP_UNRECOVER = 0
	AGENT_STOP_RECOVER   = 1
)
