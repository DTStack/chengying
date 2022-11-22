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

package base

import (
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/internal"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workload/support"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BaseWorkload interface {
	// modify the object for the reconcile later
	GetMutateFunction(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, part workloadv1beta1.WorkLoadPart) internal.MutateFunction
	// assign default value
	DefaultObject(name, owner, namespace string) runtime.Object

	Status(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, partName string)
}

type CtrlInfo struct {
	client.Client
	Recorder record.EventRecorder
	Schema   *runtime.Scheme
}

//judge the status of baseworkload
func Status(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, part workloadv1beta1.WorkLoadPart) {
	baseWorkload, _ := GetBaseWorkload(part.BaseWorkLoad.Type)
	baseWorkload.Status(info, owner, part.BaseWorkLoad.Name)
}

func GetObjectListReconciler(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, part workloadv1beta1.WorkLoadPart) (internal.Reconciler, error) {

	baseWorkload, err := GetBaseWorkload(part.BaseWorkLoad.Type)
	if err != nil {
		return nil, err
	}

	os := &internal.ObjectListReconciler{}

	var workloadReconciler *internal.ObjectReconciler

	workloadIndex := 0
	for i, step := range part.Steps {
		if step.Action == workloadv1beta1.CreateOrUpdateAction {
			reconciler, err := GetReconciler(info, owner, part.BaseWorkLoad.Name, step)
			if err != nil {
				return nil, err
			}
			os.Append(reconciler, i)
		} else {
			if workloadReconciler == nil {
				workloadReconciler = &internal.ObjectReconciler{
					Name:     part.BaseWorkLoad.Type,
					Object:   baseWorkload.DefaultObject(part.BaseWorkLoad.Name, owner.Name, owner.Namespace),
					MutateFn: baseWorkload.GetMutateFunction(info, owner, part),
					Owner:    owner,
					Schema:   info.Schema,
				}
			}
			workloadIndex = i
		}
	}
	//workload's reconcile after the bound is finish
	os.Append(workloadReconciler, workloadIndex)

	return os, nil
}

var BaseWorkloadList = map[schema.GroupVersionKind]BaseWorkload{
	Deployment{}.GroupVersionKind(): &Deployment{
		MaxUnavailable: intstr.FromInt(0),
		MaxSurge:       intstr.FromString("25%"),
		Replicas:       1,
	},
	StatefulSet{}.GroupVersionKind(): &StatefulSet{
		Replicas: 1,
	},
	DaemonSet{}.GroupVersionKind(): &DaemonSet{},
	Job{}.GroupVersionKind():       &Job{},
}

func GetBaseWorkload(typ string) (BaseWorkload, error) {
	gvk, err := support.GetGvk(typ)
	if err != nil {
		return nil, err
	}
	bw, exist := BaseWorkloadList[gvk]
	if !exist {
		return nil, fmt.Errorf("baseworkload %s is not support", typ)
	}
	return bw, nil
}
