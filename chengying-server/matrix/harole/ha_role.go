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

package harole

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
)

var (
	roleDataMutex sync.RWMutex
	roleData      = map[int]map[string]map[string]string{}
	insProcessMap = map[string]*insProcess{}

	//instance 失效统计
	expireCount = map[string]uint32{}
	expireTimes = 1

	end = "end"
)

type roleInfo struct {
	productId   int
	serviceName string
	agentId     string
	roleName    string
}

type insProcess struct {
	productId   int
	serviceName string
	agentId     string
	sid         string
	haRoleCmd   string
	quitCount   int
	quit        chan interface{}
}

func StartRoleRunner() error {

	roleChan := make(chan roleInfo, 10)
	go waitForChannel(roleChan, expireCount)

	runnerWhereCause := roleRunnerWhereCase()
	for {
		//获取符合要求的实例
		instances, err := model.DeployInstanceList.GetInstanceListByWhere(runnerWhereCause)
		if err != nil {
			return err
		}
		insAgentIds := []string{}
		if len(instances) > 0 {
			for _, ins := range instances {
				if p, ok := insProcessMap[ins.AgentId]; ok {
					if p.haRoleCmd != ins.HaRoleCmd {
						p.haRoleCmd = ins.HaRoleCmd
					}
					p.quitCount = 0
				} else {
					p = &insProcess{
						productId:   ins.Pid,
						serviceName: ins.ServiceName,
						agentId:     ins.AgentId,
						sid:         ins.Sid,
						haRoleCmd:   ins.HaRoleCmd,
						quit:        make(chan interface{}, 1),
					}
					insProcessMap[ins.AgentId] = p
					go p.dealHaRoleCmd(roleChan)
				}
				insAgentIds = append(insAgentIds, ins.AgentId)
			}
		}

		//检测 instances 中是否存在并递增 quitCount
		for pk, pv := range insProcessMap {
			exist := checkInsProcessForCount(pv.agentId, insAgentIds)
			if !exist {
				pv.quitCount = pv.quitCount + 1
			}
			if pv.quitCount >= expireTimes {
				delete(insProcessMap, pk)
				select {
				case pv.quit <- "":
				default:
				}
			}
		}

		time.Sleep(10 * time.Second)
	}

	return nil
}

func checkInsProcessForCount(agentId string, insAgentIds []string) (exist bool) {
	for _, insAgentId := range insAgentIds {
		if agentId == insAgentId {
			exist = true
			break
		}
	}
	return
}

func (p *insProcess) dealHaRoleCmd(roleChan chan roleInfo) {
	ticker := time.NewTicker(180 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case _, ok := <-p.quit:
			if ok {
				roleChan <- roleInfo{
					productId:   p.productId,
					serviceName: p.serviceName,
					agentId:     p.agentId,
					roleName:    end,
				}
				log.Infof("quit dealHaRoleCmd insProcess.productId: %v,insProcess.serviceName: %v,insProcess.agentId: %v", p.productId, p.serviceName, p.agentId)
			}
			return
		}
		roleName := runHaRoleCmd(p.haRoleCmd, p.agentId, p.sid)
		//将执行结果输入通道
		roleChan <- roleInfo{
			productId:   p.productId,
			serviceName: p.serviceName,
			agentId:     p.agentId,
			roleName:    roleName,
		}
	}
}

func runHaRoleCmd(haRoleCmd string, agentId string, sId string) (roleName string) {
	params := &agent.ExecScriptParams{}
	params.ExecScript = "#!/bin/sh\n" + haRoleCmd
	params.AgentId = agentId
	params.Timeout = "10s"

	err, instance := model.DeployInstanceList.GetInstanceInfoByAgentId(agentId)
	if err != nil || instance.Status != model.INSTANCE_STATUS_RUNNING {
		return "-"
	}

	//执行ha_role_cmd，并且返回结果
	err, resp := agent.AgentClient.AgentExecOften(sId, params, "")
	if err != nil {
		log.Errorf("[runHaRoleCmd] sid: %v, err: %v, resp: %v", sId, err, resp)
		return
	}

	data, _ := json.Marshal(resp.Data.(map[string]interface{})["result"])
	execResp := model.ExecScriptResponse{}
	json.Unmarshal(data, &execResp)

	if execResp.Failed == true {
		log.Errorf("[runHaRoleCmd] server response failed %v", execResp)
		err = fmt.Errorf("%v", execResp.Response)
		return
	}

	roleName = strings.TrimRight(execResp.Response, "\n")
	return
}

func waitForChannel(roleChan chan roleInfo, expireCount map[string]uint32) error {
	for {
		select {
		case role := <-roleChan:
			if role.roleName == end {
				delRoleData(role)
				delete(expireCount, role.agentId)
			} else if role.roleName == "" {
				if _, ok := getRoleData(role); ok {
					count := expireCount[role.agentId]
					expireCount[role.agentId] = atomic.AddUint32(&count, 1)

					if atomic.LoadUint32(&count) >= 3 {
						delRoleData(role)
						delete(expireCount, role.agentId)
					}
				}
			} else {
				setRoleData(role)
				//重置
				_, ok := expireCount[role.agentId]
				if ok {
					count := expireCount[role.agentId]
					atomic.StoreUint32(&count, 0)
				}
			}
		}
	}
}

func RoleData(pid int, serviceName string) map[string]string {
	roleDataMutex.RLock()
	defer roleDataMutex.RUnlock()
	productMap, ok := roleData[pid]
	if ok {
		serviceMap, ok := productMap[serviceName]
		if ok {
			return serviceMap
		}
	}
	return nil
}

//return (value,bool)
//1. the roleName
//2. true(exist) or false(not exist)
func getRoleData(role roleInfo) (string, bool) {
	roleDataMutex.RLock()
	defer roleDataMutex.RUnlock()
	productMap, ok := roleData[role.productId]
	if ok {
		serviceMap, ok := productMap[role.serviceName]
		if ok {
			roleName, ok := serviceMap[role.agentId]
			if ok {
				return roleName, ok
			}
		}
	}
	return "", false
}

//return true(set success) or false(set failed)
func setRoleData(role roleInfo) {
	roleDataMutex.Lock()
	defer roleDataMutex.Unlock()
	productMap, ok := roleData[role.productId]
	if !ok {
		productMap = map[string]map[string]string{}
		roleData[role.productId] = productMap
	}
	serviceMap, ok := productMap[role.serviceName]
	if !ok {
		serviceMap = map[string]string{}
		productMap[role.serviceName] = serviceMap
	}
	serviceMap[role.agentId] = role.roleName
}

func delRoleData(role roleInfo) {
	roleDataMutex.Lock()
	defer roleDataMutex.Unlock()
	productMap, ok := roleData[role.productId]
	if ok {
		serviceMap, ok := productMap[role.serviceName]
		if ok {
			delete(serviceMap, role.agentId)
		}
	}
}

func roleRunnerWhereCase() (whereCause dbhelper.WhereCause) {
	whereCause = whereCause.NotEqual("ha_role_cmd", "")
	whereCause = whereCause.And()
	whereCause = whereCause.Included("health_state", model.INSTANCE_HEALTH_OK, model.INSTANCE_HEALTH_NOTSET)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("status", model.INSTANCE_STATUS_RUNNING)
	return
}
