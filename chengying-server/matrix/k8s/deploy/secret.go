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
	sqlModel "dtstack.com/dtstack/easymatrix/matrix/model"
	"encoding/json"
	kschema "k8s.io/apimachinery/pkg/runtime/schema"
)

func ApplyImageSecret(cache kube.ClientCache,clusterId int, namespace string, store sqlModel.ImageStore) error {
	secret := k8sModel.NewDockerConfigSecret(namespace, store)
	gvr := &kschema.GroupVersionResource{
		Group:    k8sModel.SECRET_GROUP,
		Version:  k8sModel.SECRET_VERSION,
		Resource: k8sModel.SECRET_RESOURCE,
	}
	secretBytes, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	secretDynamic := NewDynamic(secretBytes, gvr, k8sModel.SECRET_KIND)
	if cache == nil{
		return ApplyDynamicResource(secretDynamic, clusterId)
	}
	return cache.GetClient(namespace).Apply(context.Background(),secret)
}
