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

package host

const (
	InitStatus               = -1
	InstallSidecarOk         = 1
	InstallScriptWrapperOk   = 2
	InstallScriptWrapperFail = -2
	InitInitializeShOk       = 3
	InitInitializeShFail     = -3
	SidecarOffline           = -4
	K8SDockerInitializeOk    = 5
	K8SDockerInitializeFail  = -5
	K8SNodeInitializeOk      = 6
	K8SNodeInitializeFail    = -6
	K8SNodeDeploymentOk      = 7
	K8SNodeDeploymentFailed  = -7
)

const (
	SCRIPT_WRAPPER_SH             string = "script_wrapper"
	ENVIRONMENT_INIT_SH           string = "environment_init"
	DOCKER_ENVIRONMENT_INIT_SH    string = "docker_environment_init"
	INSTALL_AGENTX_SH             string = "install_agentx"
	INSTALL_KUBERNETES_SH         string = "install_kubernetes"
	INSTALL_KUBERNETES_V1BETA1_SH string = "install_kubernetes_v1beta1"
	KUBERNETES_MODE               string = "1"
)

const (
	SUCCESS_SIDECAR_INSTALL          string = "管控安装成功"
	ERROR_SIDECAR_INSTALL            string = "管控安装失败"
	SUCCESS_SCRIPT_WRAPPER_INSTALL   string = "script wrapper安装成功, 启动成功"
	ERROR_SCRIPT_WRAPPER_START       string = "script wrapper安装成功, 启动失败"
	ERROR_SCRIPT_WRAPPER_INSTALL     string = "script wrapper安装失败"
	SUCCESS_HOST_INIT                string = "主机初始化成功"
	ERROR_HOST_INIT                  string = "主机初始化失败"
	ERROR_HOST_OFFLINE               string = "主机下线"
	K8S_SUCCESS_DOCKCER_INIT         string = "K8S DOCKER初始化成功"
	K8S_ERROR_DOCKCER_INIT           string = "K8S DOCKER初始化失败"
	K8S_SUCCESS_NODE_INIT            string = "K8S NODE初始化成功"
	K8S_ERROR_NODE_INIT              string = "K8S NODE初始化失败"
	K8S_SUCCESS_NODE_DEPLOYMENT_INIT string = "K8S NODE部署成功"
	K8S_ERROR_NODE_DEPLOYMENT_INIT   string = "K8S NODE部署失败"
)
