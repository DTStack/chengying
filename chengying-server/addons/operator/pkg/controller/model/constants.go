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

package model

import corev1 "k8s.io/api/core/v1"

const (
	MoleServiceName              = "service"
	MoleConfigName               = "configmap"
	MoleIngressName              = "ingress"
	MoleJobName                  = "job"
	MoleDeploymentName           = "deployment"
	MolePodName                  = "pod"
	MoleHealthEndpoint           = "/api/health"
	MoleConfigVolumeName         = "volume"
	MoleLogsVolumeName           = "log"
	MoleMountPath                = "/mount"
	LogPath                      = "/tmp/dtstack/"
	MoleServiceAccountName       = "dtstack"
	MoleCom                      = "dtstack.com"
	DefaultMemoryRequest         = "100Mi"
	DefaultMemoryLimit           = "1Gi"
	DefaultLogSidecarMemoryLimit = "500Mi"
	DefaultCpuLimit              = "500m"
	DefaultCpuRequest            = "0"
	EnvHostAlias                 = "HostAlias"
)

var SupportResource = map[corev1.ResourceName]struct{}{
	corev1.ResourceCPU:    {},
	corev1.ResourceMemory: {},
}

var VolumeConfigMapMode int32 = 0755
