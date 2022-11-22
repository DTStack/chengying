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

type InstanceEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (this *InstanceEvent) GetType() string {
	return this.Type
}

func (this *InstanceEvent) GetMessage() string {
	return this.Message
}

type InstallEvent struct {
	InstanceEvent
	InstallSchema  interface{}   `json:"installSchema"`
	InstallParam   *InstallParms `json:"installParam"`
	InstallResp    interface{}   `json:"installResp"`
	ConfigParam    []interface{} `json:"configParam"`
	ConfigResp     []interface{} `json:"configResp"`
	PostDeployResp interface{}   `json:"configUpdateResp"`
}

type UnInstallEvent struct {
	InstanceEvent
	UnInstallParam *ShellParams `json:"unInstallParam"`
	UnInstallResp  interface{}  `json:"unInstallResp"`
}

type ConfigEvent struct {
	InstanceEvent
	ConfigSchema interface{}   `json:"configSchema"`
	ConfigPath   []string      `json:"configPath"`
	ConfigResp   []interface{} `json:"configResp"`
}

type StartEvent struct {
	InstanceEvent
	StartParam *StartParams `json:"startParam"`
	StartResp  interface{}  `json:"startResp"`
}

type StopEvent struct {
	InstanceEvent
	StopResp interface{} `json:"startResp"`
}

type ExecEvent struct {
	InstanceEvent
	ExecScriptParam *ExecScriptParams `json:"execScriptParam"`
	ExecResp        interface{}       `json:"execResp"`
}

type ErrorEvent struct {
	InstanceEvent
	ErrorResp interface{} `json:"errorResp"`
}

type UnknownEvent struct {
	InstanceEvent
}
