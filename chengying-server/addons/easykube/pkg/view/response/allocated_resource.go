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

package response

type AllocatedResponse struct {
	Nodes          int    `json:"nodes"`
	ErrorNodes     int    `json:"error_nodes"`
	MemSizeDisplay string `json:"mem_size_display"`
	MemUsedDisplay string `json:"mem_used_display"`
	CpuSizeDisplay string `json:"cpu_size_display"`
	CpuUsedDisplay string `json:"cpu_used_display"`
	PodSizeDisplay int64  `json:"pod_size_display"`
	PodUsedDisplay int    `json:"pod_used_display"`
}
