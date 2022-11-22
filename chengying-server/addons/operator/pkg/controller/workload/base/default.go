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
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetReconciler(info *CtrlInfo, owner *workloadv1beta1.WorkLoad, partName string, step workloadv1beta1.WorkLoadPartStep) (internal.Reconciler, error) {
	gvk, err := support.GetGvk(step.Type)
	if err != nil {
		return nil, err
	}
	obj, err := info.Schema.New(gvk)

	if err != nil {
		log.Error(err, "check the support type, the gvk is not support by k8s", "gvk", gvk.String())
		return nil, err
	}
	name := util.BuildStepName(util.BuildBaseName(owner.Name, partName), step.Name)
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		log.Error(err, "convert owner to metav1.Object fail", "type", gvk)
		return nil, err
	}
	objMeta.SetName(name)
	objMeta.SetNamespace(owner.Namespace)

	return &internal.ObjectReconciler{
		Name:   gvk.Kind,
		Object: obj,
		Owner:  owner,
		Schema: info.Schema,
		MutateFn: func(object runtime.Object) (bool, error) {

			existingMap, err := util.ToMap(obj)
			if err != nil {
				log.Error(err, "object to map fail", "type", fmt.Sprintf("%T", obj))
				return false, err
			}
			syncd, err := util.DeepCopy(step.Object.Object, existingMap)
			if err != nil {
				log.Error(err, "copy desired to existing fail")
				return false, err
			}
			if err = util.ToObject(syncd, obj); err != nil {
				bts, _ := json.Marshal(syncd)
				log.Error(err, "to runtime.Object fail", "json", string(bts), "object-type", fmt.Sprintf("%T", obj))
				return false, err
			}

			return false, nil
		},
	}, nil
}

func ifConfigMapChanged(v *corev1.Volume, owner *workloadv1beta1.WorkLoad, part *workloadv1beta1.WorkLoadPart, c client.Client) (bool, error) {
	if v.ConfigMap != nil {
		cn := v.ConfigMap.Name
		existing := &corev1.ConfigMap{}
		// the reader client is cache reader
		// the error occures when talk with apiserver ignore.
		if err := c.Get(context.TODO(), client.ObjectKey{Name: cn, Namespace: owner.Namespace}, existing); err != nil {
			return false, nil
		}
		//compare if the volume refs configmap in the step is different with the configmap in the k8s
		for _, step := range part.Steps {
			if step.Type == support.CreateTypeConfigmap {
				configmapName := util.BuildStepName(util.BuildBaseName(owner.Name, part.BaseWorkLoad.Name), step.Name)
				if configmapName == v.ConfigMap.Name {
					existingMap, err := util.ToMap(existing)
					if err != nil {
						log.Error(err, "existing configmap to map fail")
						return false, err
					}
					syncd, err := util.DeepCopy(step.Object.Object, existingMap)
					if err != nil {
						log.Error(err, "copy desired configmap to existing fail")
						return false, err
					}
					obj := existing.DeepCopy()

					util.ToObject(syncd, existing)

					return !reflect.DeepEqual(obj, existing), nil
				}
			}
		}
	}
	return false, nil
}

func ensureVolume(spec *corev1.PodSpec, desired *corev1.Volume) error {
	if len(spec.Volumes) == 0 {
		spec.Volumes = []corev1.Volume{}
	}
	exist := false
	for i, v := range spec.Volumes {
		if v.Name == desired.Name {
			exist = true
			vMap, err := util.ToMap(v)
			if err != nil {
				log.Error(err, "ensureVolume existing to map fail")
				return err
			}
			dMap, err := util.ToMap(desired)
			if err != nil {
				log.Error(err, "ensureVolume desired to map fail")
				return err
			}
			syncd, err := util.DeepCopy(dMap, vMap)
			if err != nil {
				log.Error(err, "ensureInitContainer deepcopy fail")
				return err
			}
			util.ToObject(syncd, &v)
			spec.Volumes[i] = v
		}
	}
	if !exist {
		spec.Volumes = append(spec.Volumes, *desired)
	}
	return nil
}

func ensureInitContainer(spec *corev1.PodSpec, desired *corev1.Container) error {
	if len(spec.InitContainers) == 0 {
		spec.InitContainers = []corev1.Container{}
	}
	exist := false
	for i, c := range spec.InitContainers {
		if c.Name == desired.Name {
			exist = true
			cMap, err := util.ToMap(c)
			if err != nil {
				log.Error(err, "ensureInitContainer existing to map fail")
				return err
			}
			dMap, err := util.ToMap(desired)
			if err != nil {
				log.Error(err, "ensureInitContainer desired to map fail")
				return err
			}
			syncd, err := util.DeepCopy(dMap, cMap)
			if err != nil {
				log.Error(err, "ensureInitContainer deepcopy fail")
				return err
			}
			util.ToObject(syncd, &c)
			spec.InitContainers[i] = c
		}
	}
	if !exist {
		spec.InitContainers = append(spec.InitContainers, *desired)
	}
	return nil
}

func ensureContainer(spec *corev1.PodSpec, desired *corev1.Container) error {
	if len(spec.Containers) == 0 {
		spec.Containers = []corev1.Container{}
	}
	exist := false
	for i, c := range spec.Containers {
		if c.Name == desired.Name {
			exist = true
			cMap, err := util.ToMap(c)
			if err != nil {
				log.Error(err, "ensureContainer existing to map fail")
				return err
			}
			dMap, err := util.ToMap(desired)
			if err != nil {
				log.Error(err, "ensureContainer desired to map fail")
				return err
			}
			syncd, err := util.DeepCopy(dMap, cMap)
			if err != nil {
				log.Error(err, "ensureContainer deepcopy fail")
				return err
			}
			util.ToObject(syncd, &c)
			spec.Containers[i] = c
		}
	}
	if !exist {
		spec.Containers = append(spec.Containers, *desired)
	}
	return nil
}
