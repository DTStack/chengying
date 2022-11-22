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

package model

import (
	molev1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/mole/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	sigsclient "sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

func MoleJob(cr *molev1.Mole, name string) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: BuildResourceName(MoleJobName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Labels: map[string]string{
				"deploy_uuid":         cr.Spec.Product.DeployUUid,
				"cluster_id":          strconv.Itoa(cr.Spec.Product.ClusterId),
				"product_name":        cr.Spec.Product.ProductName,
				"product_version":     cr.Spec.Product.ProductVersion,
				"parent_product_name": cr.Spec.Product.ParentProductName,
				"service_name":        name,
				"service_version":     cr.Spec.Product.Service[name].Version,
				"group":               cr.Spec.Product.Service[name].Group,
				"com":                 MoleCom,
			},
			Namespace: cr.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        BuildResourceName(MolePodName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
					Labels:      getPodLabels(cr, name),
					Annotations: getPodAnnotations(cr, nil, name),
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: getImagePullSecrets(cr),
					Containers: []corev1.Container{
						{
							Name:            ConvertDNSRuleName(name),
							Image:           cr.Spec.Product.Service[name].Instance.Deployment.Image,
							Command:         strings.Split(cr.Spec.Product.Service[name].Instance.PostDeploy, " "),
							Resources:       getResources(cr, name),
							ImagePullPolicy: "Always",
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
}

func MoleJobSelector(cr *molev1.Mole, name string) sigsclient.ObjectKey {
	return sigsclient.ObjectKey{
		Namespace: cr.Namespace,
		Name:      BuildResourceName(MoleJobName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	}
}
