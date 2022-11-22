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

package view

type NamespaceEventRsp struct {
	Size int `json:"size"`
	Events []NamespaceEvent `json:"events"`
}

type NamespaceEvent struct {
	Id  int 	`json:"id"`
	Time string `json:"time"`
	Type string `json:"type"`
	Reason string `json:"reason"`
	Resource string `json:"resource"`
	Message string `json:"message"`
}
