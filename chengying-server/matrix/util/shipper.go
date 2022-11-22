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

package util

var apiShipper = &ApiShipperPwd{}

type ApiShipperPwd struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Group    string `json:"group"`
}

type PwdConnectParams struct {
	ApiShipperPwd
	ClusterId   int    `json:"cluster_id"`
	ClusterType string `json:"cluster_type"`
	Role        string `json:"role"`
}

type PwdInstallParams struct {
	ApiShipperPwd
	ClusterId   int    `json:"cluster_id"`
	ClusterType string `json:"cluster_type"`
	Role        string `json:"role"`
	Cmd         string `json:"cmd"`
}

type ApiShipperPk struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	User  string `json:"user"`
	Pk    string `json:"pk"`
	Group string `json:"group"`
}

type PkConnectParams struct {
	ApiShipperPk
	ClusterId   int    `json:"cluster_id"`
	ClusterType string `json:"cluster_type"`
	Role        string `json:"role"`
}

type PkInstallParams struct {
	ApiShipperPk
	ClusterId   int    `json:"cluster_id"`
	ClusterType string `json:"cluster_type"`
	Role        string `json:"role"`
	Cmd         string `json:"cmd"`
}

type ApiShipperCheck struct {
	Aid int `json:"aid"`
}

type ApiShipperCheckByIp struct {
	Ip string `json:"ip"`
}
