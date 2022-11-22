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

package events

import (
	"fmt"
)

var (
	DEFAULT_RESOURCES_POD_KIND        = "Pod"
	DEFAULT_RESOURCES_SERVICE_KIND    = "Service"
	DEFAULT_RESOURCES_DEPLOYMENT_KIND = "Deployment"
	DEFAULT_RESOURCES_INGRESSES_KIND  = "Ingress"
	DEFAULT_RESOURCES_EVENT_KIND      = "Event"
)

type Event struct {
	Namespace   string      `json:"namespace"`
	Resource    string      `json:"resource"`
	Key         string      `json:"key"`
	Operation   string      `json:"operation"`
	Object      interface{} `json:"object"`
	Workspaceid int         `json:"workspaceid"`
}

type EventResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (e *Event) Info() string {
	return fmt.Sprintf("%s-%s-%s-%s", e.Namespace, e.Resource, e.Key, e.Operation)
}
