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

package monitor

import (
	easymonitor "dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg"
	monitorevents "dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
)

func StartMonitor(namespace string, ch chan struct{}) error {
	transmitor := &AgentModeTransmitor{}
	monitorevents.Transmitor = transmitor
	err := easymonitor.StartMonitorController("", "", namespace, ch)
	if err != nil {
		return err
	}
	return nil
}

type AgentModeTransmitor struct {
}

func (a *AgentModeTransmitor) Push(event monitorevents.Eventer) {
	e := event.(*monitorevents.Event)
	EventCache.Push(e)
}

func (a *AgentModeTransmitor) Process() {

}
