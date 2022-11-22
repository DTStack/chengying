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
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/util"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/view/response"
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"github.com/kataras/iris/context"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strconv"
	"strings"
)

func GetAllocated(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetAllocated: %v", ctx.Request().RequestURI)
	c := clientCache.GetClient("")
	nodes := &corev1.NodeList{}
	if err := c.Lists(ctx, nodes); err != nil {
		return err
	}
	allocated := response.AllocatedResponse{}
	allocated.Nodes = len(nodes.Items)
	var memSize, memUsed, cpuSize, cpuUsed int64

	for _, node := range nodes.Items {
		if client.GetNodeConditionStatus(node, "Ready") != "True" {
			allocated.ErrorNodes++
			continue
		}
		pods, err := client.GetNodePods(ctx, c, &node)
		if err != nil {
			return err
		}
		resources := client.GetNodeAllocatedResources(&node, pods)
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
	c := clientCache.GetClient("")
	nodes := &corev1.NodeList{}
	if err := c.Lists(ctx, nodes); err != nil {
		return err
	}
	log.Debugf("EasyKube -> GetTop5: node list %v", nodes.Items)
	//获取allocated数据
	resources := make([]client.NodeAllocatedResources, len(nodes.Items))
	for i, node := range nodes.Items {
		if client.GetNodeConditionStatus(node, "Ready") != "True" {
			continue
		}
		pods, err := client.GetNodePods(ctx, c, &node)
		if err != nil {
			return err
		}
		resources[i] = client.GetNodeAllocatedResources(&node, pods)
	}
	// top5排序
	count := 0
	top5 := response.Top5Response{}
	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].CPURequestsFraction > resources[j].CPURequestsFraction
	})
	for _, v := range resources {
		if count >= 5 {
			break
		}
		top5.CpuTop5 = append(top5.CpuTop5, response.Top5Attribute{
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
		top5.MemTop5 = append(top5.MemTop5, response.Top5Attribute{
			Ip:    v.LocalIp,
			Usage: v.MemoryRequestsFraction,
		})
		count++
	}

	return top5
}

func GetWorkLoad(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetWorkLoad: %v", ctx.Request().RequestURI)

	workload, err := client.GetWorkLoad(ctx, clientCache.GetClient(""))
	if err != nil {
		return err
	}
	return workload
}

func GetAllocatedPodList(ctx context.Context) apibase.Result {
	log.Debugf("EasyKube -> GetAllocatedPod: %v", ctx.Request().RequestURI)

	c := clientCache.GetClient("")
	nodes := &corev1.NodeList{}
	if err := c.Lists(ctx, nodes); err != nil {
		return err
	}
	log.Debugf("EasyKube -> GetAllocatedPodList: node list %v", nodes.Items)
	podList := response.PodListResponse{}
	for _, node := range nodes.Items {
		if client.GetNodeConditionStatus(node, "Ready") != "True" {
			continue
		}
		pods, _ := client.GetNodePods(ctx, c, &node)
		resources := client.GetNodeAllocatedResources(&node, pods)
		podList.List = append(podList.List, response.NodePod{
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
	c := clientCache.GetClient("")
	nodes := &corev1.NodeList{}
	if err := c.Lists(ctx, nodes); err != nil {
		return err
	}
	roleSet := make(map[string]int)
	component := response.ComponentResponse{}
	if len(nodes.Items) > 0 {
		component.List = append(component.List, response.Component{
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
					component.List = append(component.List, response.Component{
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
				component.List = append(component.List, response.Component{
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

	c := clientCache.GetClient("")
	node, err := client.GetNodeByIp(ctx, c, hostIp)
	if err != nil {
		return err
	}

	extraInfo := response.ExtraInfoResponse{}
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
