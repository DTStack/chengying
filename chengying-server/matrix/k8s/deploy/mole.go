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
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	k8sModel "dtstack.com/dtstack/easymatrix/matrix/k8s/model"
	"dtstack.com/dtstack/easymatrix/schema"
	"encoding/json"
	kschema "k8s.io/apimachinery/pkg/runtime/schema"
)

func ApplyMole(cache kube.ClientCache,sc *schema.SchemaConfig, uncheckedServices []string, clusterId, pid int, namespace, deployUUID, secret string) error {
	moleCR := k8sModel.NewMole(sc.ProductName, namespace)
	if err := k8sModel.FillSchema(moleCR,sc, uncheckedServices, clusterId, pid, deployUUID, secret);err != nil{
		return err
	}
	if cache != nil{
		// the mole'status can directly update, if not set, status will be set to "" when update
		// it will occurs the product deploy success, but the instance's process almost %30
		existing := moleCR.DeepCopy()
		exist,err := cache.GetClient(namespace).Get(context.TODO(),existing)
		if err != nil{
			return err
		}
		if exist{
			moleCR.Status = existing.Status
		}
		return cache.GetClient(namespace).Apply(context.TODO(),moleCR)
	}
	gvr := &kschema.GroupVersionResource{
		Group:    k8sModel.MOLE_GROUP,
		Version:  k8sModel.MOLE_VERSION,
		Resource: k8sModel.MOLE_RESOURCE,
	}
	moleBytes, err := json.Marshal(moleCR)
	if err != nil {
		return err
	}
	moleDynamic := NewDynamic(moleBytes, gvr, k8sModel.MOLE_KIND)
	err = ApplyDynamicResource(moleDynamic, clusterId)
	return err
}

func DeleteMole(cache kube.ClientCache, productName, namespace string, clusterId int) error {
	moleCR := k8sModel.NewMole(productName, namespace)
	if cache != nil{
		return cache.GetClient(namespace).Delete(context.Background(),moleCR)
	}
	gvr := &kschema.GroupVersionResource{
		Group:    k8sModel.MOLE_GROUP,
		Version:  k8sModel.MOLE_VERSION,
		Resource: k8sModel.MOLE_RESOURCE,
	}
	moleBytes, err := json.Marshal(moleCR)
	if err != nil {
		return err
	}
	moleDynamic := NewDynamic(moleBytes, gvr, k8sModel.MOLE_KIND)
	err = DeleteDynamicResource(moleDynamic, clusterId)
	return err
}

func GetMole(cache kube.ClientCache,sc *schema.SchemaConfig, clusterId int, namespace string) (interface{}, error) {
	moleCR := k8sModel.NewMole(sc.ProductName, namespace)
	if cache != nil{
		exist,err := cache.GetClient(namespace).Get(context.Background(),moleCR)
		if err != nil{
			return nil, err
		}
		if !exist{
			return nil,nil
		}
		return moleCR,nil
	}
	gvr := &kschema.GroupVersionResource{
		Group:    k8sModel.MOLE_GROUP,
		Version:  k8sModel.MOLE_VERSION,
		Resource: k8sModel.MOLE_RESOURCE,
	}
	moleBytes, err := json.Marshal(moleCR)
	if err != nil {
		return nil, err
	}
	moleDynamic := NewDynamic(moleBytes, gvr, k8sModel.MOLE_KIND)
	return GetDynamicResource(moleDynamic, clusterId)
}
