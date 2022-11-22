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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type DaemonSet struct {
}

func (ds DaemonSet) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "DaemonSet",
	}
}

func (ds *DaemonSet) Status(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, partName string) {

}

func (ds *DaemonSet) DefaultObject(partName, ownerName, namespace string) runtime.Object {
	return nil
}

func (ds *DaemonSet) GetMutateFunction(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, part workloadv1beta1.WorkLoadPart) internal.MutateFunction {
	return nil
}
