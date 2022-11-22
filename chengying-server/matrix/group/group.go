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

package group

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/service"
	"dtstack.com/dtstack/easymatrix/schema"
)

type Grouper interface {
	Start() error
	Stop(stopAgentOptionsTypeArr ...int) error

	GetResult() []ServiceResult
}

type ServiceResult struct {
	ServiceName string    `json:"service_name"`
	Status      bool      `json:"status"`
	BegineTime  base.Time `json:"begin_time"`
	EndTime     base.Time `json:"end_time"`
}

type group struct {
	pid        int
	clusterId  int
	name       string
	schema     *schema.SchemaConfig
	serviceMap map[string]struct{}

	svcMap      map[string]*sync.Once
	errRlt      int64 // 0 success, 1 fail
	operationId string
	rltMu       sync.Mutex
	result      []ServiceResult
}

func NewGrouper(pid, clusterId int, name string, operationId string) (Grouper, error) {

	info, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		return nil, err
	}
	schema, err := schema.Unmarshal(info.Product)
	if err != nil {
		return nil, err
	}

	list, err := model.DeployInstanceList.GetInstanceListByPidGroup(pid, clusterId, name)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("not found `%v` group", name)
	}
	newGrouper := &group{
		pid:         pid,
		name:        name,
		schema:      schema,
		clusterId:   clusterId,
		operationId: operationId,
		serviceMap:  map[string]struct{}{},
	}
	for _, instance := range list {
		newGrouper.serviceMap[instance.ServiceName] = struct{}{}
	}
	return newGrouper, nil
}

func (g *group) init() {
	g.svcMap = map[string]*sync.Once{}
	for svcName := range g.serviceMap {
		g.svcMap[svcName] = &sync.Once{}
	}
	g.errRlt = 0
	g.result = nil
}

func (g *group) GetResult() []ServiceResult {
	return g.result
}

func (g *group) Start() error {
	g.init()

	wg := sync.WaitGroup{}
	for svcName := range g.svcMap {
		wg.Add(1)
		go g.startService(svcName, &wg)
	}
	wg.Wait()

	if g.errRlt > 0 {
		return errors.New("some services start error")
	}
	return nil
}

func (g *group) Stop(stopAgentOptionsTypeArr ...int) error {
	g.init()

	stopAgentOptionsType := agent.AGENT_STOP_UNRECOVER
	if stopAgentOptionsTypeArr != nil {
		stopAgentOptionsType = stopAgentOptionsTypeArr[0]
	}

	wg := sync.WaitGroup{}
	for svcName := range g.svcMap {
		wg.Add(1)
		go g.stopService(svcName, stopAgentOptionsType, &wg)
	}
	wg.Wait()

	if g.errRlt > 0 {
		return errors.New("some services stop error")
	}
	return nil
}

// startService should be goroutine and recursive
func (g *group) startService(svcName string, wg *sync.WaitGroup) {
	defer wg.Done()

	once, exist := g.svcMap[svcName]
	if !exist {
		wgThis := sync.WaitGroup{}
		for _, dname := range g.schema.Service[svcName].DependsOn {
			wgThis.Add(1)
			go g.startService(dname, &wgThis)
		}
		wgThis.Wait()

		return
	}

	once.Do(func() {
		wgThis := sync.WaitGroup{}
		for _, dname := range g.schema.Service[svcName].DependsOn {
			wgThis.Add(1)
			go g.startService(dname, &wgThis)
		}
		wgThis.Wait()

		if atomic.LoadInt64(&g.errRlt) != 0 {
			return
		}

		log.Infof("startService %v ...", svcName)

		rlt := ServiceResult{ServiceName: svcName, Status: true, BegineTime: base.Time(time.Now())}

		servicer, err := service.NewServicer(g.pid, g.clusterId, svcName, g.operationId)
		if err == nil {
			if err = servicer.Start(); err == nil {
				log.Debugf("waiting service(%v) GetStartStatusChan...", svcName)
				status := <-servicer.GetStartStatusChan()
				log.Debugf("end service(%v) GetStartStatusChan...", svcName)

				if status == service.STATUS_FAILED {
					err = errors.New("service STATUS_FAILED")
				}
			}
		}
		if err != nil {
			atomic.StoreInt64(&g.errRlt, 1)
			rlt.Status = false
		}
		rlt.EndTime = base.Time(time.Now())

		g.rltMu.Lock()
		g.result = append(g.result, rlt)
		g.rltMu.Unlock()

		log.Infof("startService %v finish(%v)", svcName, err)
	})
}

// stopService should be goroutine and recursive
func (g *group) stopService(svcName string, stopAgentOptionsType int, wg *sync.WaitGroup) {
	defer wg.Done()

	once, exist := g.svcMap[svcName]
	if !exist {
		wgThis := sync.WaitGroup{}
		for _, beDname := range g.findBeDepends(svcName) {
			wgThis.Add(1)
			go g.stopService(beDname, stopAgentOptionsType, &wgThis)
		}
		wgThis.Wait()

		return
	}

	once.Do(func() {
		wgThis := sync.WaitGroup{}
		for _, beDname := range g.findBeDepends(svcName) {
			wgThis.Add(1)
			go g.stopService(beDname, stopAgentOptionsType, &wgThis)
		}
		wgThis.Wait()

		log.Infof("stopService %v ...", svcName)

		rlt := ServiceResult{ServiceName: svcName, Status: true, BegineTime: base.Time(time.Now())}

		servicer, err := service.NewServicer(g.pid, g.clusterId, svcName, g.operationId)
		if err == nil {
			err = servicer.Stop(stopAgentOptionsType)
		}
		if err != nil {
			atomic.StoreInt64(&g.errRlt, 1)
			rlt.Status = false
		}
		rlt.EndTime = base.Time(time.Now())

		g.rltMu.Lock()
		g.result = append(g.result, rlt)
		g.rltMu.Unlock()

		log.Infof("stopService %v finish(%v)", svcName, err)
	})
}

func (g *group) findBeDepends(name string) []string {
	list := make([]string, 0)
	for svcName, svc := range g.schema.Service {
		for _, dname := range svc.DependsOn {
			if dname == name {
				list = append(list, svcName)
			}
		}
	}
	return list
}
