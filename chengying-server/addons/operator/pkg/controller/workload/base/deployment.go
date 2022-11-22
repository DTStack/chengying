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
	"context"
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/internal"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workload/support"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Deployment struct {
	MaxUnavailable intstr.IntOrString
	MaxSurge       intstr.IntOrString
	Replicas       int32
}

func (deploy Deployment) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}
}

func (d *Deployment) Status(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, partName string) {
	name := util.BuildBaseName(owner.Name, partName)
	namespace := owner.Namespace
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	deployment := &appsv1.Deployment{}
	if err := info.Get(context.TODO(), key, deployment); err != nil {
		//ignore the error
		return
	}
	partStatus := workloadv1beta1.WorkloadPending
	if deployment.Status.UnavailableReplicas == 0 {
		partStatus = workloadv1beta1.WorkloadRunning
	}
	if owner.Status.PartStatus == nil {
		owner.Status.PartStatus = map[string]workloadv1beta1.WorkLoadPartStatus{}
	}
	owner.Status.PartStatus[partName] = workloadv1beta1.WorkLoadPartStatus{
		Status: partStatus,
	}
}

func (d *Deployment) DefaultObject(partName, ownerName, namespace string) runtime.Object {
	name := util.BuildBaseName(ownerName, partName)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &d.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":        name,
					"control-by": "deployment",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        name,
						"control-by": "deployment",
					},
				},
				Spec: corev1.PodSpec{},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: "RollingUpdate",
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &d.MaxUnavailable,
					MaxSurge:       &d.MaxSurge,
				},
			},
		},
	}
}

func (d *Deployment) GetMutateFunction(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, part workloadv1beta1.WorkLoadPart) internal.MutateFunction {
	return func(object runtime.Object) (bool, error) {
		// sync deployment's metadata and deployment's spec
		deploy := object.(*appsv1.Deployment)
		existing, err := util.ToMap(deploy.Spec)
		if err != nil {
			log.Error(err, "deploy to map fail")
			return false, err
		}
		synced, err := util.DeepCopy(part.BaseWorkLoad.Parameters.Object, existing)
		if err != nil {
			log.Error(err, "copy desired deploy params to existing fail")
			return false, err
		}
		if err = util.ToObject(synced, deploy); err != nil {
			log.Error(err, "synced data to deploy fail")
			return false, err
		}
		spec := &deploy.Spec.Template.Spec
		ifReStart := false
		for _, step := range part.Steps {
			if step.Action == workloadv1beta1.CreateOrUpdateAction {
				continue
			}
			if step.Type == support.BoundTypeContainer {
				c := &corev1.Container{}
				if err := util.ToObject(step.Object.Object, c); err != nil {
					log.Error(err, "convert container date to corev1.container fail")
					return false, err
				}
				c.Name = step.Name
				if err = ensureContainer(spec, c); err != nil {
					return false, err
				}
				continue
			}
			if step.Type == support.BoundTypeInitContainer {
				c := &corev1.Container{}
				if err := util.ToObject(step.Object.Object, c); err != nil {
					log.Error(err, "convert container date to corev1.container fail")
					return false, err
				}
				c.Name = step.Name
				if err = ensureInitContainer(spec, c); err != nil {
					return false, err
				}
				continue
			}
			if step.Type == support.BoundTypeVolume {
				v := &corev1.Volume{}
				if err := util.ToObject(step.Object.Object, v); err != nil {
					log.Error(err, "convert volume data to corev1.Volume fail")
					return false, err
				}

				v.Name = step.Name
				if err = ensureVolume(spec, v); err != nil {
					return false, err
				}
				ifReStart, err = ifConfigMapChanged(v, owner, &part, info.Client)
				if err != nil {
					return false, err
				}
				continue
			}
		}
		return ifReStart, nil
	}
}
