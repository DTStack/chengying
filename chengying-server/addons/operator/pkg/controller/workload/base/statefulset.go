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
)

type StatefulSet struct {
	Replicas int32
}

func (sts StatefulSet) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "StatefulSet",
	}
}

func (sts *StatefulSet) Status(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, partName string) {
	name := util.BuildBaseName(owner.Name, partName)
	namespace := owner.Namespace
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	statefulset := &appsv1.StatefulSet{}
	if err := info.Get(context.TODO(), key, statefulset); err != nil {
		//ignore the error
		return
	}
	partStatus := workloadv1beta1.WorkloadPending
	if statefulset.Status.Replicas > 0 &&
		statefulset.Status.CurrentReplicas == statefulset.Status.Replicas &&
		statefulset.Status.ReadyReplicas == statefulset.Status.Replicas {
		partStatus = workloadv1beta1.WorkloadRunning
	}
	if owner.Status.PartStatus == nil {
		owner.Status.PartStatus = map[string]workloadv1beta1.WorkLoadPartStatus{}
	}
	owner.Status.PartStatus[partName] = workloadv1beta1.WorkLoadPartStatus{
		Status: partStatus,
	}
}

func (sts *StatefulSet) DefaultObject(partName, ownerName, namespace string) runtime.Object {
	name := util.BuildBaseName(ownerName, partName)
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &sts.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":        name,
					"control-by": "statefulset",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        name,
						"control-by": "statefulset",
					},
				},
				Spec: corev1.PodSpec{},
			},
		},
	}
}

func (sts *StatefulSet) GetMutateFunction(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, part workloadv1beta1.WorkLoadPart) internal.MutateFunction {
	return func(object runtime.Object) (bool, error) {
		sts := object.(*appsv1.StatefulSet)
		existing, err := util.ToMap(sts)
		if err != nil {
			log.Error(err, "sts to map fail")
			return false, err
		}
		synced, err := util.DeepCopy(part.BaseWorkLoad.Parameters.Object, existing)
		if err != nil {
			log.Error(err, "copy desired sts params to existing fail")
			return false, err
		}
		if err = util.ToObject(synced, sts); err != nil {
			log.Error(err, "synced data to sts fail")
			return false, err
		}
		spec := &sts.Spec.Template.Spec

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
			if step.Type == support.BoundTypePvc {
				if sts.Spec.VolumeClaimTemplates == nil {
					pvc := &corev1.PersistentVolumeClaim{}
					if err := util.ToObject(step.Object.Object, pvc); err != nil {
						log.Error(err, "convert pvc data to corev1.PersistentVolumeClaim fail")
						return false, err
					}
					pvc.Name = step.Name
					ensurePvc(&sts.Spec, pvc)
				}
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

// statefulset's PersistentVolumeClaim can't not modify
func ensurePvc(spec *appsv1.StatefulSetSpec, desired *corev1.PersistentVolumeClaim) {
	if len(spec.VolumeClaimTemplates) == 0 {
		spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
	}
	exist := false
	for _, v := range spec.VolumeClaimTemplates {
		if v.Name == desired.Name {
			exist = true
			break
		}
	}
	if !exist {
		spec.VolumeClaimTemplates = append(spec.VolumeClaimTemplates, *desired)
	}
}
