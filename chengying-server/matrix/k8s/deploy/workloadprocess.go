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

package deploy

import (
	"context"
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	workloadprocessrc "dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workloadprocess/reconciler"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/util"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/workload"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/workloadprocess"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"dtstack.com/dtstack/easymatrix/schema"
	"encoding/json"
	"fmt"
	"strings"
)

func ApplyWorkloadProcess(cache kube.ClientCache, sc *schema.SchemaConfig, uncheckedSvc []string,
	namespace string, deployuuid string, pid int, clusterID int, store *model.ImageStore) (bool, error) {

	process := workloadprocess.New()
	process.Name = strings.ToLower(sc.ProductName)
	process.Namespace = namespace
	process.Spec.ProductName = sc.ProductName
	process.Spec.ProductId = pid
	process.Spec.DeployUUId = deployuuid
	process.Spec.ParentProductName = sc.ParentProductName
	process.Spec.ProductVersion = sc.ProductVersion
	process.Spec.ClusterId = clusterID
	existing := process.DeepCopy()
	exist, err := cache.GetClient(namespace).Get(context.TODO(), existing)
	if err != nil {
		return false, err
	}
	if exist {
		process.Spec.LastDeployUUId = existing.Spec.DeployUUId
	}
	undeploy := make(map[string]bool, len(uncheckedSvc))
	for _, svcName := range uncheckedSvc {
		undeploy[svcName] = true
	}
	ifChanged := false

	// 处理workload类型产品包中的依赖组件为非workload类型的服务组件
	workloadService := map[string]schema.ServiceConfig{}
	for name, svc := range sc.Service {
		if svc.Workload != "" {
			workloadService[name] = svc
		}
	}

	for svcName, svc := range workloadService {
		if undeploy[svcName] {
			continue
		}
		//use the desired version, workload:version
		workloadVersion := strings.Split(svc.Workload, "@")
		wlTyp := workloadVersion[0]
		wlversion := ""
		if len(workloadVersion) == 2 {
			wlversion = workloadVersion[1]
		}
		wl, err := modelkube.WorkloadDefinition.Get(wlTyp, wlversion)
		if err != nil {
			return false, err
		}
		if wl == nil {
			return false, fmt.Errorf("the workload type %s is not support, please check the workload type", wlTyp)
		}
		parts, err := modelkube.WorkloadPart.Select(wl.Id)
		if err != nil {
			return false, err
		}
		if parts == nil {
			return false, fmt.Errorf("the part of workload type %s is nil, please check the workload type", wlTyp)
		}
		steps := make(map[int][]modelkube.WorloadStepSchema, len(parts))
		for _, part := range parts {
			stepTbsc, err := modelkube.WorkloadStep.Select(part.Id)
			if err != nil {
				return false, err
			}
			steps[part.Id] = stepTbsc
		}
		builder := &workload.Builder{
			Def:         wl,
			Parts:       parts,
			Steps:       steps,
			Schema:      sc,
			ProductName: sc.ProductName,
			ServiceName: svcName,
			Namespace:   namespace,
			Store:       store,
		}
		workload, err := builder.Build()
		if err != nil {
			return false, err
		}
		if process.Spec.WorkLoads == nil {
			process.Spec.WorkLoads = map[string]workloadv1beta1.ServiceWorkload{}
		}
		process.Spec.WorkLoads[svcName] = workloadv1beta1.ServiceWorkload{
			Version:  svc.Version,
			Group:    svc.Group,
			WorkLoad: *workload,
		}

		existing := workload.DeepCopy()
		existing.Namespace = namespace
		existing.Name = util.BuildWorkloadName(sc.ProductName, svcName)
		now, err := json.Marshal(existing)
		if err != nil {
			log.Errorf("[workloadprocess deploy]: marshal to json fail, err %v", err)
			return false, err
		}
		exist, err = cache.GetClient(namespace).Get(context.TODO(), existing)
		if err != nil {
			return false, err
		}
		if !exist {
			ifChanged = true
		}
		last := existing.Annotations[workloadprocessrc.LAST_WITHOUT_ANNOTATION]

		if string(now) != last {
			ifChanged = true
		}
	}
	if err = cache.GetClient(namespace).Apply(context.TODO(), process); err != nil {
		return false, err
	}
	return ifChanged, nil
}

func GetWorkloadProcess(cache kube.ClientCache, sc *schema.SchemaConfig, namespace string) (*workloadv1beta1.WorkloadProcess, error) {
	process := workloadprocess.New()
	process.Name = strings.ToLower(sc.ProductName)
	process.Namespace = namespace
	exist, err := cache.GetClient(namespace).Get(context.TODO(), process)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return process, nil
}

func DeleteWorkloadProcess(cache kube.ClientCache, productname, namespace string) error {
	process := workloadprocess.New()
	process.Namespace = namespace
	process.Name = strings.ToLower(productname)
	err := cache.GetClient(namespace).Delete(context.TODO(), process)
	if err != nil {
		return err
	}
	return nil
}
