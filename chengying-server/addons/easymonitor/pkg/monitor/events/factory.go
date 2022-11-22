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

package events

import (
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"encoding/json"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	OPERATION_CREATE_OR_UPDATE         = "createOrUpdate"
	OPERATION_DELETE                   = "delete"
	RESOURCE_LABLE_CLUSTER_ID          = "cluster_id"
	RESOURCE_LABLE_DEPLOY_UUID         = "deploy_uuid"
	RESOURCE_LABLE_PARENT_PRODUCT_NAME = "parent_product_name"
	RESOURCE_LABLE_PRODUCT_NAME        = "product_name"
	RESOURCE_LABLE_PRODUCT_VERSION     = "product_version"
	RESOURCE_LABLE_SERVICE_NAME        = "service_name"
	RESOURCE_LABLE_SERVICE_VERSION     = "service_version"
	RESOURCE_LABLE_SERVICE_GROUP       = "group"
	RESOURCE_LABLE_PRODUCT_ID          = "pid"
)

func getMetaInfo(object interface{}) (string, string, string, string) {
	switch object.(type) {
	case *unstructured.Unstructured:
		objectR := object.(*unstructured.Unstructured)
		ns := objectR.GetNamespace()
		apiVersion := objectR.GetAPIVersion()
		kind := objectR.GroupVersionKind().Kind
		name := objectR.GetName()
		return ns, apiVersion, kind, name
	case *appv1.Deployment:
		deploy := object.(*appv1.Deployment)
		return deploy.Namespace, deploy.GetSelfLink(), DEFAULT_RESOURCES_DEPLOYMENT_KIND, deploy.Name
	case *v1.Pod:
		pod := object.(*v1.Pod)
		return pod.Namespace, pod.GetSelfLink(), DEFAULT_RESOURCES_POD_KIND, pod.Name
	case *v1.Service:
		service := object.(*v1.Service)
		return service.Namespace, service.GetSelfLink(), DEFAULT_RESOURCES_SERVICE_KIND, service.Name
	case *v1.Event:
		event := object.(*v1.Event)
		return event.Namespace, event.GetSelfLink(), DEFAULT_RESOURCES_EVENT_KIND, event.Name
	default:
		log.Errorf("unkown object type: %v", object)
		return "", "", "", ""
	}
}

func logEvent(operation string, object interface{}) {
	content, err := json.Marshal(object)
	if err != nil {
		log.Errorf("%v", err.Error())
	}
	ns, api, kind, name := getMetaInfo(object)
	log.Infof("New event, ns: %v, kind: %v, api: %v, name: %v, type: %v, object: %v",
		ns,
		kind,
		api,
		name,
		operation,
		string(content[:]))
}

func NewEvent(key, operation string, object interface{}) Eventer {
	//logEvent(operation, object)
	event := &Event{}
	ns, _, kind, _ := getMetaInfo(object)
	event.Resource = kind
	event.Namespace = ns
	event.Key = key
	event.Operation = operation
	event.Object = object
	return event
}
