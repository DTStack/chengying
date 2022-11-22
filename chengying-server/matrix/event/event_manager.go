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

package event

const (
	maxEventCache = 128
)

var eventManager *eventManage

func init() {
	eventManager = &eventManage{eventQueue: make(chan *Event, maxEventCache)}
	go eventManager.EventHandler()
}

type eventManage struct {
	Observable
	eventQueue chan *Event
}

type EventManager interface {
	EventReciever(event *Event)
	EventHandler()
	AddObserver(observer Observer)
	RemoveObserver(observer Observer)
	Notify(event *Event)
}

func (this *eventManage) AddObserver(observer Observer) {
	this.RegistObserver(observer)
}

func (this *eventManage) RemoveObserver(observer Observer) {
	this.UnregistObserver(observer)
}

func (this *eventManage) Notify(event *Event) {
	this.NofityObservers(*event)
}

func (this *eventManage) EventHandler() {
	for {
		ev, ok := <-this.eventQueue

		if !ok {
			return
		}
		this.Notify(ev)
	}
}

func (this *eventManage) EventReciever(event *Event) {
	// log.Debugf("[EventReciever] receive: %+v", event.Data)
	this.eventQueue <- event
}

func GetEventManager() EventManager {
	return eventManager
}
