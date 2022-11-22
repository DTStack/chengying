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

package xke_service

import (
	"fmt"
	rutil "github.com/rancher/kontainer-engine/drivers/util"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"gopkg.in/yaml.v2"
	"strings"
)

const (
	SPLIT_SEP_ROLES                = ","
	DEFALT_CONFIG_RKE_SSH_KEY_PATH = "~/.ssh/id_rsa"
	DEFALT_CONFIG_RKE_SSH_PORT     = "22"
	DEFALT_CONFIG_RKE_SSH_USER     = "docker"
	DEFALT_CONFIG_RKE_ROLES        = "controlplane,etcd,worker"
	DEFALT_CONFIG_RKE_VERSION      = "v1.16.3"
)

var ignoreDockerVersion = true

var rkeConfigTemplate = &v3.RancherKubernetesEngineConfig{
	Version: DEFALT_CONFIG_RKE_VERSION,
	Authorization: v3.AuthzConfig{
		Mode: "none",
	},
	ClusterName: "dtstack",
	Network: v3.NetworkConfig{
		Plugin: "flannel",
	},
	IgnoreDockerVersion: &ignoreDockerVersion,
	Ingress: v3.IngressConfig{
		Provider: "nginx",
	},
	Monitoring: v3.MonitoringConfig{
		Provider: "metrics-server",
	},
	Restore: v3.RestoreConfig{
		Restore: false,
	},
}

func BuildRKEConfigFromRaw(raw string) (*v3.RancherKubernetesEngineConfig, error) {
	rkeConfig, err := rutil.ConvertToRkeConfig(raw)
	if err != nil {
		return &rkeConfig, err
	}
	return &rkeConfig, nil
}

func GetDefaultRKEconfig(version, clusterName, networkPlugin string) v3.RancherKubernetesEngineConfig {
	tmpConfig := *rkeConfigTemplate
	tmpConfig.ClusterName = clusterName
	tmpConfig.Network.Plugin = networkPlugin
	tmpConfig.Version = version

	return tmpConfig
}

func GetDefaultRKEconfigRaw(version, clusterName, networkPlugin string) (string, error) {
	tmpConfig := *rkeConfigTemplate
	tmpConfig.ClusterName = clusterName
	tmpConfig.Network.Plugin = networkPlugin
	tmpConfig.Version = version

	config, err := yaml.Marshal(tmpConfig)
	if err != nil {
		return "", err
	}
	return string(config), nil
}

func AddNodeToRawConfig(raw string, nodeIp, port, roles, sshKeyPath, user string) error {
	config, err := BuildRKEConfigFromRaw(raw)
	if err != nil {
		return err
	}
	for _, node := range config.Nodes {
		if node.Address == nodeIp {
			return fmt.Errorf("error add, node exist %v", nodeIp)
		}
	}
	config.Nodes = append(config.Nodes, v3.RKEConfigNode{
		Address:    nodeIp,
		User:       user,
		SSHKeyPath: sshKeyPath,
		Port:       port,
		Role:       strings.Split(roles, SPLIT_SEP_ROLES),
	})
	return nil
}

func AddNodeToConfig(config *v3.RancherKubernetesEngineConfig, nodeIp, port, roles, sshKeyPath, user string) error {
	for _, node := range config.Nodes {
		if node.Address == nodeIp {
			return fmt.Errorf("error add, node exist %v", nodeIp)
		}
	}
	config.Nodes = append(config.Nodes, v3.RKEConfigNode{
		Address:    nodeIp,
		User:       user,
		SSHKeyPath: sshKeyPath,
		Port:       port,
		Role:       strings.Split(roles, SPLIT_SEP_ROLES),
	})
	return nil
}

func AddDefaultNodeToConfig(config *v3.RancherKubernetesEngineConfig, nodeIp string) error {
	for _, node := range config.Nodes {
		if node.Address == nodeIp {
			return fmt.Errorf("error add, node exist %v", nodeIp)
		}
	}
	config.Nodes = append(config.Nodes, v3.RKEConfigNode{
		Address:    nodeIp,
		User:       DEFALT_CONFIG_RKE_SSH_USER,
		SSHKeyPath: DEFALT_CONFIG_RKE_SSH_KEY_PATH,
		Port:       DEFALT_CONFIG_RKE_SSH_PORT,
		Role:       strings.Split(DEFALT_CONFIG_RKE_ROLES, SPLIT_SEP_ROLES),
	})
	return nil
}

func RemoveNodeFromConfig(config *v3.RancherKubernetesEngineConfig, nodeIp string) {
	for index, node := range config.Nodes {
		if node.Address == nodeIp {
			config.Nodes = append(config.Nodes[:index], config.Nodes[index+1:]...)
			break
		}
	}
}
