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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type Job struct {
}

func (job Job) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "batch",
		Version: "v1",
		Kind:    "Job",
	}
}
func (job *Job) Status(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, partName string) {
	name := util.BuildBaseName(owner.Name, partName)
	namespace := owner.Namespace
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	j := &batchv1.Job{}
	if err := info.Get(context.TODO(), key, j); err != nil {
		return
	}
	partStatus := workloadv1beta1.WorkloadPending
	if *j.Spec.Completions == j.Status.Succeeded {
		partStatus = workloadv1beta1.WorkloadRunning
	}
	if owner.Status.PartStatus == nil {
		owner.Status.PartStatus = map[string]workloadv1beta1.WorkLoadPartStatus{}
	}
	owner.Status.PartStatus[partName] = workloadv1beta1.WorkLoadPartStatus{
		Status: partStatus,
	}
}

func (job *Job) DefaultObject(partName, ownerName, namespace string) runtime.Object {
	name := util.BuildBaseName(ownerName, partName)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        name,
						"control-by": "job",
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: "Never",
				},
			},
		},
	}
}

func (job *Job) GetMutateFunction(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, part workloadv1beta1.WorkLoadPart) internal.MutateFunction {
	return func(object runtime.Object) (bool, error) {
		job := object.(*batchv1.Job)
		existing, err := util.ToMap(job)
		if err != nil {
			log.Error(err, "job to map fail")
			return false, err
		}
		synced, err := util.DeepCopy(part.BaseWorkLoad.Parameters.Object, existing)
		if err != nil {
			log.Error(err, "copy desired job params to existing fail")
			return false, err
		}
		if err = util.ToObject(synced, job); err != nil {
			log.Error(err, "synced data to job fail")
			return false, err
		}
		spec := &job.Spec.Template.Spec
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
		}
		return false, nil
	}
}
