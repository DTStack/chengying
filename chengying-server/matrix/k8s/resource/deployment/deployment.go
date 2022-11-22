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

package deployment

import (
	"context"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"encoding/json"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var GVK = schema.GroupVersionKind{
	Group:   "apps",
	Version: "v1",
	Kind:    "Deployment",
}

func New() *appsv1.Deployment{
	return &appsv1.Deployment{}
}

func Convert(obj runtime.Object) *appsv1.Deployment{
	return obj.(*appsv1.Deployment)
}

func Get(ctx context.Context,client kube.Client,namespace,name string) (*appsv1.Deployment,error){
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: namespace,
		},
	}
	exist,err := client.Get(ctx,deploy)
	if err != nil{
		return nil,err
	}
	if !exist{
		return nil,nil
	}
	return deploy,nil
}


func ToObject(bts []byte)(*appsv1.Deployment,error){
	r,err := base.Schema.New(GVK)
	if err != nil{
		log.Errorf("[deployment]: new object error: %v",err)
		return nil,err
	}
	err = json.Unmarshal(bts,r)
	if err!= nil{
		log.Errorf("[deployment]: json %s unmarshal error: %v",string(bts),err)
		return nil,err
	}
	deploy := r.(*appsv1.Deployment)
	return deploy,nil
}

func Ping(client kube.Client, namespace string) error{
	ping := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name: "dtstack-dryru",
		},
		Spec:       appsv1.DeploymentSpec{
			Selector:                &metav1.LabelSelector{
				MatchLabels: 	map[string]string{
						"app":			"dtstack-dryru",
				},
			},
			Template:                corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
						Name:          		"dtstack-dryru",
						Labels: 			map[string]string{
							"app":			"dtstack-dryru",
						},
				},
				Spec:       corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: 					"dtstack-dryru",
							Image:                  "dtstack-dryru",
						},
					},
				},
			},
		},
	}
	if _,err := client.Get(context.Background(),ping);err != nil{
		return err
	}

	if err := client.DryRun(base.Create,ping);err != nil{
		return err
	}
	return nil
}
