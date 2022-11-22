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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"time"
)

var (
	Duration = 60
)

func SyncStrategies() {
	duration := time.Duration(Duration) * time.Second
	for {
		log.Infof("SyncStrategies ...")
		syncStrategies()
		time.Sleep(duration)
	}
}

func syncStrategies() {
	err, strategy := model.StrategyList.GetDeployedStrategyList()
	if err != nil {
		log.Errorf("err :%v", err.Error())
		return
	}
	err, strategyResources := model.StrategyResourceList.GetStrategyResourceList()
	if err != nil {
		log.Errorf("err :%v", err.Error())
		return
	}
	err, strategyAssigns := model.StrategyAssignList.GetStrategyAssignList()
	if err != nil {
		log.Errorf("err :%v", err.Error())
		return
	}
	rebuildStrategyMap(strategy, strategyResources, strategyAssigns)
}

func rebuildStrategyMap(strategys []*model.StrategyInfo,
	resources []*model.StrategyResourceInfo,
	assing []*model.StrategyAssignInfo) {
	s := make(map[int]*model.StrategyInfo)
	r := make(map[int]*model.StrategyResourceInfo)
	a := make(map[int][]*model.StrategyAssignInfo)

	for _, ss := range strategys {
		s[ss.ID] = ss
	}
	for _, rr := range resources {
		r[rr.StrategyId] = rr
	}
	for _, aa := range assing {
		a[aa.StrategyId] = append(a[aa.StrategyId], aa)
	}
	StrategyCacheMap.ReInit(s, r, a)
}
