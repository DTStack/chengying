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

package strategy

const (
	STR_TYPE_HOST     = 0
	STR_TYPE_SERVICE  = 1
	CRON_TYPE_MINUTES = 0
	CRON_TYPE_HOURS   = 1
	CRON_TYPE_DAYS    = 2
)

const (
	MAX_HOST_TASK_NUM    = 10
	MAX_SERVICE_TASK_NUM = 1
)

type TaskInfo struct {
	TaskType          int
	StrategyId        int
	StrategyName      string
	SidecarId         string
	AgentId           string
	ExecScript        string
	Timeout           string
	Parameter         string
	Host              string
	ParentProductName string
	ProductName       string
	ServiceName       string
}
