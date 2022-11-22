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

package client_go

import (
	"context"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Namespace struct {
	Name string `json:"namespace"`
}

type NamespaceListResponse struct {
	Namespaces []Namespace `json:"namespaces"`
}

type ContentResponse struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

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

type NodePod struct {
	LocalIp     string  `json:"local_ip"`
	PodUsed     int     `json:"pod_used"`
	PodSize     int64   `json:"pod_size"`
	PodUsagePct float64 `json:"pod_usage_pct"`
}

type PodListResponse struct {
	List []NodePod `json:"list"`
}

type ExtraInfoResponse struct {
	Roles   string `json:"roles"`
	Version string `json:"version"`
}

type Component struct {
	Role    string   `json:"role"`
	Status  int      `json:"status"`
	Message []string `json:"errors"`
}

type ComponentResponse struct {
	List []Component `json:"list"`
}

type WorkLoad struct {
	Load     int `json:"load"`
	Capacity int `json:"capacity"`
}

type WorkLoadResponse struct {
	CornJobs               WorkLoad `json:"CornJobs"`
	Jobs                   WorkLoad `json:"Jobs"`
	Pods                   WorkLoad `json:"Pods"`
	DaemonSets             WorkLoad `json:"DaemonSets"`
	Deployments            WorkLoad `json:"Deployments"`
	ReplicaSets            WorkLoad `json:"ReplicaSets"`
	ReplicationControllers WorkLoad `json:"ReplicationControllers"`
}

type Top5Attribute struct {
	Ip    string  `json:"ip"`
	Usage float64 `json:"usage"`
}

type Top5Response struct {
	CpuTop5 []Top5Attribute `json:"cpu_top5"`
	MemTop5 []Top5Attribute `json:"mem_top5"`
}

func GetClientSet() (kubernetes.Interface, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("[EasyKube] init config err:%v", err)
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)

	return clientset, err
}

func GetNodeList(client kubernetes.Interface) (*v1.NodeList, error) {
	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})

	return nodes, err
}

func GetNodeByName(client kubernetes.Interface, hostName string) (*v1.Node, error) {
	node, err := client.CoreV1().Nodes().Get(context.Background(), hostName, metav1.GetOptions{})
	return node, err
}

func GetNodeByIp(client kubernetes.Interface, hostIp string) (*v1.Node, error) {
	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == "InternalIP" && address.Address == hostIp {
				return &node, err
			}
		}
	}
	return nil, fmt.Errorf("cannot match ip %v from cluster", hostIp)
}

func GetNodePods(client kubernetes.Interface, node v1.Node) (*v1.PodList, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name +
		",status.phase!=" + string(v1.PodSucceeded) +
		",status.phase!=" + string(v1.PodFailed))

	if err != nil {
		return nil, err
	}
	return client.CoreV1().Pods(v1.NamespaceAll).List(context.Background(), metav1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})
}

func GetNodeConditionStatus(node v1.Node, conditionType v1.NodeConditionType) v1.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return v1.ConditionUnknown
}

func GetNodeConditionMessage(node v1.Node, conditionType v1.NodeConditionType) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Message
		}
	}
	return ""
}

func GetNamespaceList(client kubernetes.Interface) (*v1.NamespaceList, error) {
	return client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
}

func CreateNamespace(client kubernetes.Interface, name string) (*v1.Namespace, error) {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return client.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
}
