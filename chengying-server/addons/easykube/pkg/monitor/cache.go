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
	monitorevents "dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
	"sync"
)

var EventCache = &Cache{
	Events: make([]monitorevents.Eventer, 0, 1024),
}

type Cache struct {
	Events []monitorevents.Eventer
	mu     sync.Mutex
}

func (c *Cache) Push(evnet monitorevents.Eventer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Events = append(c.Events, evnet)
}

func (c *Cache) Get() []monitorevents.Eventer {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.Events) == 0 {
		return nil
	}
	events := c.Events
	c.Events = make([]monitorevents.Eventer, 0, 1024)
	return events
}
