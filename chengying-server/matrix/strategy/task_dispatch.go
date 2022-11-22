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

import (
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"strings"
	"sync"
	"time"
)

type SafeStrategyTaskChanMap struct {
	sync.RWMutex
	h  map[string]chan TaskInfo
	s  map[string]chan TaskInfo
	hl map[string]time.Time
	sl map[string]time.Time
}

var (
	StrategyTaskMap = &SafeStrategyTaskChanMap{
		h:  make(map[string]chan TaskInfo),
		s:  make(map[string]chan TaskInfo),
		hl: make(map[string]time.Time),
		sl: make(map[string]time.Time),
	}
)

func (this *SafeStrategyTaskChanMap) Add(task TaskInfo) {
	this.Lock()
	defer this.Unlock()
	switch task.TaskType {
	case STR_TYPE_HOST:
		if _, ok := this.h[task.SidecarId]; !ok {
			this.h[task.SidecarId] = make(chan TaskInfo, MAX_HOST_TASK_NUM)
		}
		this.h[task.SidecarId] <- task
	case STR_TYPE_SERVICE:
		if _, ok := this.s[task.AgentId]; !ok {
			this.s[task.AgentId] = make(chan TaskInfo, MAX_SERVICE_TASK_NUM)
		}
		this.s[task.AgentId] <- task
	default:
		log.Errorf("no such task type: %v", task.TaskType)
	}
}

func (this *SafeStrategyTaskChanMap) WaitForTask() {
	for {
		for key, _ := range this.h {
			if _, ok := this.hl[key]; !ok {
				log.Infof("open host channel, sid: %v", key)
				this.hl[key] = time.Now()
				go this.waitForHostTask(key)
			}
		}
		for key, _ := range this.s {
			if _, ok := this.sl[key]; !ok {
				log.Infof("open service channel, agentid: %v", key)
				this.sl[key] = time.Now()
				go this.waitForServiceTask(key)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (this *SafeStrategyTaskChanMap) waitForHostTask(sid string) {
	log.Infof("start Host task listen:  %v", sid)
	for {
		select {
		case task := <-this.h[sid]:
			log.Infof("%v", task)
			log.Infof("exec host task, sid: %v. task name: %v, %v", sid, task.StrategyName, task.Host)
			content, err := agent.AgentClient.ToExecCmdWithTimeout(task.SidecarId, "", task.ExecScript, task.Timeout, task.Parameter, "")
			if err != nil {
				log.Errorf("%v", err.Error())
			}
			log.Infof("sid: %v, result: %v", sid, content)
			if len(content) > 0 {
				for _, line := range strings.Split(content, "\n") {
					if len(line) == 0 {
						continue
					}
					model.EventList.NewEvent("", "", "", task.Host, task.StrategyName, line)
				}
			}
		}
	}
}

func (this *SafeStrategyTaskChanMap) waitForServiceTask(agentId string) {
	log.Infof("start service task listen:  %v", agentId)
	for {
		select {
		case task := <-this.s[agentId]:
			log.Infof("exec service task, agentId: %v. task name: %v, %v-%v", agentId, task.StrategyName, task.Host, task.ServiceName)
			content, err := agent.AgentClient.ToExecCmdWithTimeout(task.SidecarId, agentId, task.ExecScript, task.Timeout, task.Parameter, "")
			if err != nil {
				log.Errorf("%v", err.Error())
				continue
			}
			log.Infof("agentId: %v, result: %v", agentId, content)
			if len(content) > 0 {
				for _, line := range strings.Split(content, "\n") {
					if len(line) == 0 {
						continue
					}
					model.EventList.NewEvent(task.ParentProductName, task.ProductName, task.ServiceName, task.Host, task.StrategyName, line)
				}
			}
		}
	}
}
