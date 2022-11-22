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

package model

import (
	"strings"
	"time"
)

type Agent struct {
	Id             string     `db:"id"`
	SidecarId      string     `db:"sidecar_id"`
	Type           int        `db:"type"`
	Name           string     `db:"name"`
	Version        string     `db:"version"`
	IsUninstalled  int        `db:"is_uninstalled"`
	DeployDate     *time.Time `db:"deploy_date"`
	AutoDeployment int        `db:"auto_deployment"`
	LastUpdateDate *time.Time `db:"last_update_date"`
	AutoUpdated    int        `db:"auto_updated"`
}

type ExecScriptResponse struct {
	Seqno    uint32 `json:"seqno,omitempty"`
	Failed   bool   `json:"failed,omitempty"`
	Response string `json:"response,omitempty"`
	AgentId  string `json:"agentId"`
}

func GetAgentByServices(names []string) ([]Agent, error) {
	for i := range names {
		names[i] = "\"" + names[i] + "\""
	}
	query := "SELECT * FROM agent_list WHERE name IN (" + strings.Join(names, ",") + ")"
	agents := make([]Agent, 0)

	err := USE_MYSQL_DB().Select(&agents, query)
	return agents, err
}
