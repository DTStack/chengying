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

import (
	client "dtstack.com/dtstack/easymatrix/addons/oldkube/pkg/client-go"
	"dtstack.com/dtstack/easymatrix/addons/oldkube/pkg/util"
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"github.com/kataras/iris/context"
	"sort"
	"strconv"
	"strings"
)

func GetAllocated(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetAllocated: %v", ctx.Request().RequestURI)

	clientset, err := client.GetClientSet()
	if err != nil {
		log.Errorf("[EasyKube] -> GetAllocated：create clientSet err:%v", err)
		return fmt.Errorf("create clientset err:%v", err)
	}

	// 通过client-go获取k8s节点，并计算cpu、mem等数据
	nodes, err := client.GetNodeList(clientset)
	log.Debugf("EasyKube -> GetAllocated: node list %v", nodes.Items)
	if err != nil {
		return fmt.Errorf("GetNodeList err:%v", err)
	}
	allocated := client.AllocatedResponse{}
	allocated.Nodes = len(nodes.Items)
	var memSize, memUsed, cpuSize, cpuUsed int64

	for _, node := range nodes.Items {
		if client.GetNodeConditionStatus(node, "Ready") != "True" {
			allocated.ErrorNodes++
			continue
		}
		pods, _ := client.GetNodePods(clientset, node)
		resources, err := client.GetNodeAllocatedResources(node, pods)
		if err != nil {
			return fmt.Errorf("GetNodeAllocatedResources err:%v", err)
		}
		memSize += resources.MemoryCapacity
		memUsed += resources.MemoryRequests
		cpuSize += resources.CPUCapacity
		cpuUsed += resources.CPURequests
		allocated.PodSizeDisplay += resources.PodCapacity
		allocated.PodUsedDisplay += resources.AllocatedPods
	}
	allocated.MemSizeDisplay, allocated.MemUsedDisplay = util.MultiSizeConvert(memSize, memUsed)
	allocated.CpuSizeDisplay = strconv.FormatFloat(float64(cpuSize)/1000, 'f', 2, 64) + "core"
	allocated.CpuUsedDisplay = strconv.FormatFloat(float64(cpuUsed)/1000, 'f', 2, 64) + "core"
	return allocated
}

func GetTop5(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetTop5: %v", ctx.Request().RequestURI)

	clientset, err := client.GetClientSet()
	if err != nil {
		log.Errorf("[EasyKube] -> GetTop5：create clientSet err:%v", err)
		return fmt.Errorf("create clientset err:%v", err)
	}
	nodes, err := client.GetNodeList(clientset)
	log.Debugf("EasyKube -> GetTop5: node list %v", nodes.Items)
	if err != nil {
		return fmt.Errorf("GetNodeList err:%v", err)
	}
	//获取allocated数据
	resources := make([]client.NodeAllocatedResources, len(nodes.Items))
	for i, node := range nodes.Items {
		if client.GetNodeConditionStatus(node, "Ready") != "True" {
			continue
		}
		pods, _ := client.GetNodePods(clientset, node)
		resources[i], err = client.GetNodeAllocatedResources(node, pods)
		if err != nil {
			return fmt.Errorf("GetNodeAllocatedResources err:%v", err)
		}
	}
	// top5排序
	count := 0
	top5 := client.Top5Response{}
	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].CPURequestsFraction > resources[j].CPURequestsFraction
	})
	for _, v := range resources {
		if count >= 5 {
			break
		}
		top5.CpuTop5 = append(top5.CpuTop5, client.Top5Attribute{
			Ip:    v.LocalIp,
			Usage: v.CPURequestsFraction,
		})
		count++
	}

	count = 0
	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].MemoryRequestsFraction > resources[j].MemoryRequestsFraction
	})
	for _, v := range resources {
		if count >= 5 {
			break
		}
		top5.MemTop5 = append(top5.MemTop5, client.Top5Attribute{
			Ip:    v.LocalIp,
			Usage: v.MemoryRequestsFraction,
		})
		count++
	}

	return top5
}

func GetWorkLoad(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetWorkLoad: %v", ctx.Request().RequestURI)

	clientset, err := client.GetClientSet()
	if err != nil {
		log.Errorf("[EasyKube] -> GetWorkLoad：create clientSet err:%v", err)
		return fmt.Errorf("create clientset err:%v", err)
	}
	workload, err := client.GetWorkLoad(clientset)
	if err != nil {
		return fmt.Errorf("GetWorkLoad err:%v", err)
	}
	return workload
}

func GetAllocatedPodList(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetAllocatedPod: %v", ctx.Request().RequestURI)

	clientset, err := client.GetClientSet()
	if err != nil {
		log.Errorf("[EasyKube] -> GetAllocatedPodList：create clientSet err:%v", err)
		return fmt.Errorf("create clientset err:%v", err)
	}
	nodes, err := client.GetNodeList(clientset)
	log.Debugf("EasyKube -> GetAllocatedPodList: node list %v", nodes.Items)
	if err != nil {
		return fmt.Errorf("GetNodeList err:%v", err)
	}
	podList := client.PodListResponse{}
	for _, node := range nodes.Items {
		if client.GetNodeConditionStatus(node, "Ready") != "True" {
			continue
		}
		pods, _ := client.GetNodePods(clientset, node)
		resources, err := client.GetNodeAllocatedResources(node, pods)
		if err != nil {
			return fmt.Errorf("GetNodeAllocatedResources err:%v", err)
		}
		podList.List = append(podList.List, client.NodePod{
			LocalIp:     resources.LocalIp,
			PodUsed:     resources.AllocatedPods,
			PodSize:     resources.PodCapacity,
			PodUsagePct: resources.PodFraction,
		})

	}
	return podList
}

func GetComponentStatus(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetComponentStatus: %v", ctx.Request().RequestURI)
	clientset, err := client.GetClientSet()
	if err != nil {
		log.Errorf("[EasyKube] -> GetWorkLoad：create clientSet err:%v", err)
		return fmt.Errorf("create clientset err:%v", err)
	}
	nodes, err := client.GetNodeList(clientset)
	if err != nil {
		return fmt.Errorf("GetNodeList err:%v", err)
	}
	roleSet := make(map[string]int)
	component := client.ComponentResponse{}
	if len(nodes.Items) > 0 {
		component.List = append(component.List, client.Component{
			Role:    "nodes",
			Status:  0,
			Message: nil,
		})
	}
	count := 1

	for _, node := range nodes.Items {
		// 先处理 error node 情况
		isError := client.GetNodeConditionStatus(node, "Ready") != "True"
		if isError {
			message := client.GetNodeConditionMessage(node, "Ready")
			component.List[0].Message = append(component.List[0].Message, node.Name+message)
		}
		isWorker := true
		// 处理其他 role
		for label := range node.Labels {
			if strings.HasPrefix(label, "node-role.kubernetes.io") {
				isWorker = false
				role := strings.Split(label, "/")[1]
				if _, ok := roleSet[role]; !ok {
					roleSet[role] = count
					component.List = append(component.List, client.Component{
						Role:    role,
						Status:  0,
						Message: nil,
					})
					count += 1
				}
				if isError {
					component.List[roleSet[role]].Status = 1
					component.List[roleSet[role]].Message = append(component.List[roleSet[role]].Message, node.Name)
				}
			}

		}
		// 特判没有node-role标签为worker节点
		if isWorker {
			if _, ok := roleSet["worker"]; !ok {
				roleSet["worker"] = count
				component.List = append(component.List, client.Component{
					Role:    "worker",
					Status:  0,
					Message: nil,
				})
				count += 1
			}
			if isError {
				component.List[roleSet["worker"]].Status = 1
				component.List[roleSet["worker"]].Message = append(component.List[roleSet["worker"]].Message, node.Name)
			}
		}
	}
	return component
}

func GetExtraInfo(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetExtraInfo: %v", ctx.Request().RequestURI)
	hostIp := ctx.URLParam("hostIp")
	if hostIp == "" {
		return fmt.Errorf("hostIp is empty")
	}

	clientset, err := client.GetClientSet()
	if err != nil {
		return fmt.Errorf("create clientset err:%v", err)
	}
	node, err := client.GetNodeByIp(clientset, hostIp)
	if err != nil {
		return fmt.Errorf("GetNodeList err:%v", err)
	}

	extraInfo := client.ExtraInfoResponse{}

	isWorker := true
	for label := range node.Labels {
		if strings.HasPrefix(label, "node-role.kubernetes.io") {
			isWorker = false
			if extraInfo.Roles == "" {
				extraInfo.Roles += strings.Split(label, "/")[1]
			} else {
				extraInfo.Roles += "," + strings.Split(label, "/")[1]
			}
		}
	}
	if isWorker {
		extraInfo.Roles = "Worker"
	}
	extraInfo.Version = node.Status.NodeInfo.KubeletVersion
	return extraInfo
}

func GetNamespaceList(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetNamespaceList: %v", ctx.Request().RequestURI)
	clientset, err := client.GetClientSet()
	if err != nil {
		return fmt.Errorf("create clientset err:%v", err)
	}
	namespaces, err := client.GetNamespaceList(clientset)
	if err != nil {
		return fmt.Errorf("list namespace err:%v", err)
	}
	list := client.NamespaceListResponse{}
	for _, namespace := range namespaces.Items {
		list.Namespaces = append(list.Namespaces, client.Namespace{Name: namespace.Name})
	}
	return list
}

func CreateNamespace(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> CreateNamespace: %v", ctx.Request().RequestURI)
	param := client.Namespace{}

	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	log.Debugf("param:%+v", param)
	clientset, err := client.GetClientSet()
	if err != nil {
		return fmt.Errorf("create clientset err:%v", err)
	}
	namespace, err := client.CreateNamespace(clientset, param.Name)
	if err != nil {
		return fmt.Errorf("create namespace err:%v", err)
	}
	log.Debugf("namespace:%+v", namespace)
	return client.Namespace{Name: namespace.Name}
}

func ApplyDynamicResource(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> CreateMole: %v", ctx.Request().RequestURI)

	param := client.DynamicData{}

	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}

	clientset, err := client.GetDynamicClient()
	if err != nil {
		return fmt.Errorf("create clientset err:%v", err)
	}
	resp, err := client.ApplyDynamicResource(ctx, clientset, param)
	if err != nil {
		return fmt.Errorf("apply %v err:%v", param.Resource, err)
	}
	return resp
}

func DeleteDynamicResource(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> DeleteDynamicResource: %v", ctx.Request().RequestURI)

	param := client.DynamicData{}

	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}

	clientset, err := client.GetDynamicClient()
	if err != nil {
		return fmt.Errorf("create clientset err:%v", err)
	}
	err = client.DeleteDynamicResource(ctx, clientset, param)
	if err != nil {
		return fmt.Errorf("delete %v err:%v", param.Resource, err)
	}
	return nil
}

func GetDynamicResource(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetDynamicResource: %v", ctx.Request().RequestURI)

	param := client.DynamicData{}

	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	clientset, err := client.GetDynamicClient()
	if err != nil {
		return fmt.Errorf("create clientset err:%v", err)
	}
	resp, err := client.GetDynamicResource(ctx, clientset, param)
	if err != nil {
		return fmt.Errorf("get %v err:%v", param.Resource, err)
	}

	return resp.Object
}
