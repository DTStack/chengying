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
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"sync"
)

type safeStrategyCacheMap struct {
	sync.RWMutex
	s map[int]*model.StrategyInfo
	r map[int]*model.StrategyResourceInfo
	a map[int][]*model.StrategyAssignInfo
}

var (
	StrategyCacheMap = &safeStrategyCacheMap{
		s: make(map[int]*model.StrategyInfo),
		r: make(map[int]*model.StrategyResourceInfo),
		a: make(map[int][]*model.StrategyAssignInfo),
	}
)

func (this *safeStrategyCacheMap) ReInit(
	s map[int]*model.StrategyInfo,
	r map[int]*model.StrategyResourceInfo,
	a map[int][]*model.StrategyAssignInfo) {

	this.Lock()
	defer this.Unlock()
	this.s = s
	this.r = r
	this.a = a
}

func (this *safeStrategyCacheMap) GetStrategy() map[int]*model.StrategyInfo {
	this.RLock()
	defer this.RUnlock()
	return this.s
}

func (this *safeStrategyCacheMap) GetResource() map[int]*model.StrategyResourceInfo {
	this.RLock()
	defer this.RUnlock()
	return this.r
}

func (this *safeStrategyCacheMap) GetAssign() map[int][]*model.StrategyAssignInfo {
	this.RLock()
	defer this.RUnlock()
	return this.a
}
