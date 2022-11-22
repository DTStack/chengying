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

package impl

const (
	COOKIE_CURRENT_CLUSTER_ID    = "em_current_cluster_id"
	COOKIE_INSTALL_CLUSTER_ID    = "em_install_cluster_id"
	COOKIE_PARENT_PRODUCT_NAME   = "em_current_parent_product"
	COOKIE_CURRENT_K8S_NAMESPACE = "em_current_k8s_namespace"
)

type ConfigUpdateParam struct {
	File    string                 `json:"file"`
	Content string                 `json:"content"`
	Values  map[string]interface{} `json:"values"`
	Deleted string                 `json:"deleted"`
}

type AddonInstallParam struct {
	Sid         string                 `json:"sid"`
	AddonId     string                 `json:"addonId"`
	ConfigParam map[string]interface{} `json:"configParam"`
}
