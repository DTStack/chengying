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

package node

import (
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/xke-service"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
	clustergenerator "dtstack.com/dtstack/easymatrix/matrix/k8s/cluster"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
)

const (
	NODE_ROLE_CONTROL = "controlplane"
	NODE_ROLE_ETCD    = "etcd"
	NODE_ROLE_WORKER  = "worker"
)

var (
	AllNodeRoles  = []string{NODE_ROLE_CONTROL, NODE_ROLE_ETCD, NODE_ROLE_WORKER}
	RoleToRkeRole = map[string]string{
		"Etcd":    NODE_ROLE_ETCD,
		"Worker":  NODE_ROLE_WORKER,
		"Control": NODE_ROLE_CONTROL,
	}
)

const (
	MaxNodeCache = 64
)

var NodeManager *nodeManager

func init() {
	NodeManager = &nodeManager{nodeQ: make(chan Node, MaxNodeCache)}
	go NodeManager.nodeHandler()
}

type nodeManager struct {
	nodeQ chan Node
}

func (this *nodeManager) AddNode(clusterId, aid int, sid, name, roles string) {
	this.nodeQ <- Node{
		ClusterId: clusterId, Sid: sid, Roles: roles, Aid: aid, Name: name,
	}
}

func (this *nodeManager) nodeHandler() {
	for {
		node, ok := <-this.nodeQ

		if !ok {
			return
		}
		//if !this.IsClusterReady(node) {
		//	model.DeployHostList.UpdateStatus(node.Aid, host.K8SNodeInitializeFail, host.K8S_SUCCESS_DOCKCER_INIT)
		//	model.DeployHostList.UpdateSteps(node.Aid, host.K8SNodeInitializeFail)
		//	continue
		//}
		err := this.freshCluster(node)
		if err != nil {
			log.Errorf("freshCluster err: %v", err.Error())
		}
	}
}

func (this *nodeManager) IsClusterReady(node Node) bool {
	list, err := model.DeployClusterHostRel.GetClusterHostRelList(node.ClusterId)

	if err != nil {
		log.Errorf("%v", err.Error())
		return false
	}
	rolesMap := make(map[string]int)

	for _, rel := range list {
		roles := strings.Split(rel.Roles, ",")
		for _, role := range roles {
			if _, ok := RoleToRkeRole[role]; !ok {
				log.Errorf("role is not legal: %v", role)
				continue
			}
			rkeRole := RoleToRkeRole[role]
			if _, ok := rolesMap[rkeRole]; ok {
				rolesMap[rkeRole] = rolesMap[rkeRole] + 1
			} else {
				rolesMap[rkeRole] = 1
			}
		}
	}
	count := 0
	log.Infof("cluster roles: %v", rolesMap)
	for _, role := range AllNodeRoles {
		if _, ok := rolesMap[role]; ok {
			count += 1
		}
	}
	if count == len(AllNodeRoles) {
		log.Infof("\t\tcluster is ready to install, clusterId: %v", node.ClusterId)
		return true
	}
	log.Infof("\t\tcluster is not ready to install, clusterId: %v", node.ClusterId)
	return false
}

func (this *nodeManager) freshCluster(node Node) error {
	log.Infof("-->Step freshCluster, cluser: %v", node)
	log.OutputInfof(node.Sid, "-->Step freshCluster, cluser: %v", node)
	query := "SELECT h.status, h.steps, h.ip, roles, yaml, deploy_cluster_list.version " +
		"FROM deploy_cluster_list " +
		"LEFT JOIN deploy_cluster_host_rel ON deploy_cluster_list.id = deploy_cluster_host_rel.clusterId " +
		"LEFT JOIN deploy_host as h ON deploy_cluster_host_rel.sid = h.sid " +
		"WHERE deploy_cluster_list.id = ? and deploy_cluster_list.is_deleted = 0 and h.status >= ? and h.isDeleted = 0"
	hostList := make([]HostNode, 0)
	if err := model.USE_MYSQL_DB().Select(&hostList, query, node.ClusterId, host.K8SDockerInitializeOk); err != nil {
		return err
	}
	rkeConfig, err := xke_service.BuildRKEConfigFromRaw(hostList[0].Yaml)
	if err != nil {
		return err
	}

	sshKeyPath := xke_service.DEFALT_CONFIG_RKE_SSH_KEY_PATH
	sshPort := xke_service.DEFALT_CONFIG_RKE_SSH_PORT
	sshUser := xke_service.DEFALT_CONFIG_RKE_SSH_USER
	for _, hn := range hostList {
		//remove later
		var rkeRoles []string
		for _, r := range strings.Split(hn.Roles, ",") {
			if _, ok := RoleToRkeRole[r]; ok {
				rkeRoles = append(rkeRoles, RoleToRkeRole[r])
			}
		}
		xke_service.AddNodeToConfig(rkeConfig, hn.Ip, sshPort, strings.Join(rkeRoles, ","), sshKeyPath, sshUser)
	}

	// 保留不带 rancher 后缀的 version 保存到数据库
	tmpConfig := rkeConfig
	tmpConfigYaml, err := yaml.Marshal(tmpConfig)
	if err != nil {
		return err
	}

	// example: v1.16.3 转为 v1.16.3-rancher1-1
	rkeConfig.Version, err = kube.DeployClusterK8sAvailable.GetRealVersion(hostList[0].Version)
	if err != nil {
		return fmt.Errorf("database err: %v", err)
	}

	xke, err := xke_service.NewXkeService()
	if err != nil {
		return err
	}
	log.Infof("create k8s with rke config: %v", rkeConfig)

	rkeConfigYaml, err := yaml.Marshal(rkeConfig)
	if err != nil {
		return err
	}

	err = xke.Create(rkeConfig.ClusterName, string(rkeConfigYaml), node.ClusterId)
	if err != nil {
		model.DeployHostList.UpdateStatus(node.Aid, host.K8SNodeInitializeFail, err.Error())
		return err
	}
	model.DeployHostList.UpdateStatus(node.Aid, host.K8SNodeInitializeOk, host.K8S_SUCCESS_NODE_INIT)
	model.DeployHostList.UpdateSteps(node.Aid, host.K8SNodeInitializeOk)

	model.DeployClusterList.UpdateYamlById(node.ClusterId, string(tmpConfigYaml))
	info := &clustergenerator.GeneratorInfo{
		Type:        constant.TYPE_SELF_BUILD,
		ClusterInfo: &modelkube.ClusterInfo{
			Id:            node.ClusterId,
			Name:          node.Name,
		},
	}
	selfBuild, err := clustergenerator.GetTemplateFile(info,false)
	if err != nil {
		return err
	}
	err = xke.DeployWithF(rkeConfig.ClusterName, string(selfBuild))
	if err != nil {
		model.DeployHostList.UpdateStatus(node.Aid, host.K8SNodeDeploymentFailed, err.Error())
		return err
	}
	log.Infof("<--Step freshCluster success, cluser: %v", node.ClusterId)
	log.OutputInfof(node.Sid, "<--Step freshCluster success, cluser: %v", node.ClusterId)
	return nil
}
