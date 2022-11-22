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

package client_go

import (
	"context"
	"dtstack.com/dtstack/easymatrix/addons/oldkube/pkg/util"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type DynamicData struct {
	Data     string `json:"data"`
	Group    string `json:"group"`
	Resource string `json:"resource"`
	Version  string `json:"version"`
}

func GetDynamicClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("[EasyKube] init config err:%v", err)
		return nil, err
	}
	clientset, err := dynamic.NewForConfig(config)

	return clientset, err
}

func ApplyDynamicResource(ctx context.Context, client dynamic.Interface, body DynamicData) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    body.Group,
		Version:  body.Version,
		Resource: body.Resource,
	}

	decodeData := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(body.Data), &decodeData)
	if err != nil {
		return nil, fmt.Errorf("yaml unmarshal err:%v", err)
	}

	obj := util.MapConvert(decodeData) // convert map[interface{}]interface{} to map[string]interface{}
	u := unstructured.Unstructured{obj}

	resource, err := client.Resource(gvr).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})
	if err == nil {
		if resource.Object["spec"] != nil {
			resource.Object["spec"] = u.Object["spec"]
		} else {
			resource.Object["data"] = u.Object["data"]
		}

		return client.Resource(gvr).Namespace(u.GetNamespace()).Update(ctx, resource, metav1.UpdateOptions{})
	}
	return client.Resource(gvr).Namespace(u.GetNamespace()).Create(ctx, &u, metav1.CreateOptions{})
}

func DeleteDynamicResource(ctx context.Context, client dynamic.Interface, body DynamicData) error {
	gvr := schema.GroupVersionResource{
		Group:    body.Group,
		Version:  body.Version,
		Resource: body.Resource,
	}

	decodeData := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(body.Data), &decodeData)
	if err != nil {
		return fmt.Errorf("yaml unmarshal err:%v", err)
	}

	obj := util.MapConvert(decodeData) // convert map[interface{}]interface{} to map[string]interface{}
	u := unstructured.Unstructured{obj}

	_, err = client.Resource(gvr).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})

	if err != nil { // if resource not exist , no need delete
		log.Errorf(err.Error())
		return nil
	}
	return client.Resource(gvr).Namespace(u.GetNamespace()).Delete(ctx, u.GetName(), metav1.DeleteOptions{})
}

func GetDynamicResource(ctx context.Context, client dynamic.Interface, body DynamicData) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    body.Group,
		Version:  body.Version,
		Resource: body.Resource,
	}

	decodeData := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(body.Data), &decodeData)
	if err != nil {
		return nil, fmt.Errorf("yaml unmarshal err:%v", err)
	}

	obj := util.MapConvert(decodeData) // convert map[interface{}]interface{} to map[string]interface{}
	u := unstructured.Unstructured{obj}

	return client.Resource(gvr).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})
}
