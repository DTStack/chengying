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

package health

import (
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	ReloadInterval time.Duration
}

type HealthCheck struct {
	config  Config
	Infos   []*model.HealthCheckInfo
	Ticker  *time.Ticker
	IdTickM map[int]*time.Ticker
	//context  context.Context
	mtx sync.RWMutex
}

func NewHealthCheck(cfg Config) *HealthCheck {
	//ctx, _ := context.WithCancel(context.Background())
	healthCheck := HealthCheck{
		config: cfg,
		//context:  ctx,
		IdTickM: make(map[int]*time.Ticker, 0),
	}
	return &healthCheck
}

func (h *HealthCheck) Run() {
	sql := fmt.Sprintf("select id,cluster_id,product_name,pid,service_name,agent_id,"+
		"sid,script_name,script_name_display,auto_exec,period,retries from %s", model.TBL_SERVICE_HEALTH_CHECK)
	if err := model.HealthCheck.GetDB().Select(&h.Infos, sql); err != nil {
		log.Errorf("err: %v", err.Error())
		return
	}
	if len(h.Infos) != 0 {
		for _, info := range h.Infos {
			h.StartOneHealthCheck(info)
		}
	}
	idleDelay := time.NewTimer(h.config.ReloadInterval)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "FATAL!!! %v\n", r)
		}
		idleDelay.Stop()
	}()
LOOP:
	for {
		idleDelay.Reset(h.config.ReloadInterval)
		select {
		case sig := <-sigs:
			fmt.Printf("Quit according to signal '%s'\n", sig.String())
			break LOOP
		case <-idleDelay.C:
			h.Update()
		}
	}
	return
}

func (h *HealthCheck) Stop(id int) {
	h.IdTickM[id].Stop()
}

func (h *HealthCheck) Update() {
	var newInfos []*model.HealthCheckInfo
	sql := fmt.Sprintf("select id,cluster_id,product_name,pid,service_name,agent_id,"+
		"sid,script_name,script_name_display,auto_exec,period,retries from %s", model.TBL_SERVICE_HEALTH_CHECK)
	if err := model.HealthCheck.GetDB().Select(&newInfos, sql); err != nil {
		log.Errorf("err: %v", err.Error())
		return
	}
	oldInfosMap := make(map[int]*model.HealthCheckInfo, 0)
	for _, oldInfo := range h.Infos {
		oldInfosMap[oldInfo.ID] = oldInfo
	}
	newInfosMap := make(map[int]*model.HealthCheckInfo, 0)
	for _, newInfo := range newInfos {
		newInfosMap[newInfo.ID] = newInfo
	}
	for _, newInfo := range newInfos {
		if _, ok := oldInfosMap[newInfo.ID]; ok {
			if newInfo.ID == oldInfosMap[newInfo.ID].ID && newInfo.AutoExec != oldInfosMap[newInfo.ID].AutoExec {
				if oldInfosMap[newInfo.ID].AutoExec == true && newInfo.AutoExec == false {
					h.Stop(oldInfosMap[newInfo.ID].ID)
				}
				if oldInfosMap[newInfo.ID].AutoExec == false && newInfo.AutoExec == true {
					h.StartOneHealthCheck(newInfo)
				}
			}
		} else {
			if _, ok := h.IdTickM[newInfo.ID]; !ok {
				h.StartOneHealthCheck(newInfo)
			}

		}
	}
	for _, oldInfo := range h.Infos {
		if _, ok := newInfosMap[oldInfo.ID]; !ok {
			h.Stop(oldInfo.ID)
		}
	}
	h.Infos = newInfos
}

func (h *HealthCheck) StartOneHealthCheck(info *model.HealthCheckInfo) {
	if info.AutoExec == true {
		h.mtx.Lock()
		defer h.mtx.Unlock()
		if info.Retries <= 0 {
			info.Retries = 1
		}
		var period time.Duration
		if !regexp.MustCompile(`\${.*?}`).MatchString(info.Period) {
			var err error
			if period, err = time.ParseDuration(info.Period); err != nil {
				log.Errorf("service %s's period is invalid", info.ServiceName)
				return
			} else if period < time.Second {
				log.Errorf("service %s's period less than 1s", info.ServiceName)
				return
			}
		}
		h.IdTickM[info.ID] = time.NewTicker(period)
		go func(i *model.HealthCheckInfo) {
			var count int
			for {
				<-h.IdTickM[i.ID].C
				cmd := fmt.Sprintf("#!/bin/sh\n %s", i.ScriptName)
				_, err := agent.AgentClient.ToExecCmdWithTimeout(i.Sid, i.AgentId, cmd, "60s", "", "")
				if err != nil {
					count++
					if count >= i.Retries {
						sql := "UPDATE " + model.HealthCheck.TableName + " SET `exec_status`=?, `error_message`=?, `start_time`=NOW() WHERE `id`=?"
						if _, err = model.DeployInstanceList.GetDB().Exec(sql, enums.ExecStatusType.Failed.Code, err.Error(), i.ID); err != nil {
							log.Errorf("%v", err)
							return
						}
					}
				} else {
					count = 0
					sql := "UPDATE " + model.HealthCheck.TableName + " SET `exec_status`=?, `error_message`=?, `start_time`=NOW() WHERE `id`=?"
					if _, err = model.DeployInstanceList.GetDB().Exec(sql, enums.ExecStatusType.Success.Code, "", i.ID); err != nil {
						log.Errorf("%v", err)
						return
					}
				}
			}
		}(info)
	}
}
