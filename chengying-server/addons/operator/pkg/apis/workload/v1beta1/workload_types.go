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

package v1beta1

import (
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type WorkLoadPart struct {
	BaseWorkLoad BaseWorkLoad       `json:"baseworkload"`
	Steps        []WorkLoadPartStep `json:"steps,omitempty"`
}

type BaseWorkLoad struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Parameters Object `json:"parameters,omitempty"`
}

type WorkLoadPartStep struct {
	Name   string          `json:"name"`
	Type   string          `json:"type"`
	Action ComponentAction `json:"action"`
	Object Object          `json:"object"`
}

//type BaseWorkLoadSpecValue interface{}
//
//type K8sObjectSpec interface{}

type ComponentAction string

const (
	BoundAction          ComponentAction = "bound"
	CreateOrUpdateAction ComponentAction = "createorupdate"
	WorkloadRunning                      = "running"
	WorkloadPending                      = "pending"
	WorkloadFail                         = "fail"
	WorkloadCreated                      = "created"
)

// workload is a namespaced resource, that consists of a series of k8s object,
// it is used to expand a series of functions of native workload
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=workloads,scope=Namespaced
type WorkLoad struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WorkLoadSpec   `json:"spec,omitempty"`
	Status            WorkLoadStatus `json:"status,omitempty"`
}

type WorkLoadSpec struct {
	WorkLoadParts []WorkLoadPart `json:"workloadpatrs,omitempty"`
}

type WorkLoadStatus struct {
	Phase      string                        `json:"phase"`
	PartStatus map[string]WorkLoadPartStatus `json:"part_status,omitempty"`
}

type WorkLoadPartStatus struct {
	Status  string `json:"status"`
	Message string `json:"message",omitempty`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WorkLoadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkLoad `json:"items"`
}
type Object unstructured.Unstructured

func (u *Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Object)
}

func (u *Object) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		fmt.Println("0-------------", err.Error())
		return err
	}
	u.Object = m
	return nil
}

func (o *Object) DeepCopyInto(out *Object) {
	clone := o.DeepCopy()
	*out = *clone
	return
}

func (o *Object) DeepCopy() *Object {
	if o.Object == nil {
		return nil
	}
	return &Object{
		Object: runtime.DeepCopyJSON(o.Object),
	}
}
