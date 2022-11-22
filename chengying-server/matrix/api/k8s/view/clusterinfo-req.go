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

import "fmt"

type ClusterInfoReq struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Mode       int    `json:"mode"`
	Version    string `json:"version"`
	Desc       string `json:"desc"`
	Tags       string `json:"tags"`
	Configs    string `json:"configs"`
	Yaml       string `json:"yaml"`
	Status     int    `json:"status"`
	ErrorMsg   string `json:"errorMsg"`
	CreateUser string `json:"create_user"`
	NetworkPlugin NetWorkPlugin `json:"network_plugin"`
}


type NetWorkPlugin string

func (n NetWorkPlugin) String() string{
	return fmt.Sprintf("{\"network_plugin\":\"%s\"}",string(n))
}
