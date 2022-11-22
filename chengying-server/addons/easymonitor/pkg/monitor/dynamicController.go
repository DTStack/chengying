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

package monitor

import (
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/crd"
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"strings"
	"sync"
	"time"
)

func newUnstructured(apiVersion, kind string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
		},
	}
}

func capitalize(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

type DynamicController struct {
	mutex                  sync.Mutex
	dc                     dynamic.Interface
	gvrs                   []*schema.GroupVersionResource
	informers              dynamicinformer.DynamicSharedInformerFactory
	informerUpdateObjectCh chan *unstructured.Unstructured
	informerDeleteObjectCh chan *unstructured.Unstructured
	transmitor             events.TransmitorInterface
}

func NewNewDynamicSharedInformerFactoryNs(client dynamic.Interface, defaultResync time.Duration, ns string) dynamicinformer.DynamicSharedInformerFactory {
	return dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, defaultResync, ns, nil)
}

func NewDynamicController(config *rest.Config, namespace string, trs events.TransmitorInterface) (*DynamicController, error) {

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Errorf("%v", err.Error())
		return nil, err
	}
	return &DynamicController{
		informers:              NewNewDynamicSharedInformerFactoryNs(client, 0, namespace),
		dc:                     client,
		informerUpdateObjectCh: make(chan *unstructured.Unstructured, 1),
		informerDeleteObjectCh: make(chan *unstructured.Unstructured, 1),
		transmitor:             trs,
	}, nil
}

func (dc *DynamicController) Add(gvrk crd.GroupVersionResourceKind) error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	gvr := gvrk.GroupVersionResource()
	dc.gvrs = append(dc.gvrs, &gvr)
	dc.informers.ForResource(gvr).Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			o := obj.(*unstructured.Unstructured)
			o.SetGroupVersionKind(gvrk.GroupVersionKind())
			dc.informerUpdateObjectCh <- o
		},
		UpdateFunc: func(old, updated interface{}) {
			oldObj := old.(*unstructured.Unstructured)
			updateObj := updated.(*unstructured.Unstructured)
			if oldObj.Object["status"] == nil || updateObj.Object["status"] == nil {
				return
			}
			oldStatus := oldObj.Object["status"].(map[string]interface{})["phase"]
			newStatus := updateObj.Object["status"].(map[string]interface{})["phase"]
			log.Infof("old %v new %v", oldStatus, newStatus)
			if oldStatus == newStatus {
				return
			}
			updateObj.SetGroupVersionKind(gvrk.GroupVersionKind())
			dc.informerUpdateObjectCh <- updateObj
		},
		DeleteFunc: func(obj interface{}) {
			o := obj.(*unstructured.Unstructured)
			o.SetGroupVersionKind(gvrk.GroupVersionKind())
			dc.informerDeleteObjectCh <- o
		},
	})
	return nil
}

func (dc *DynamicController) Run(threadiness int, stopCh <-chan struct{}) {
	// Let the workers stop when we are done
	log.Infof("Starting dynamic rc monitor controller...")
	dc.informers.Start(stopCh)
	synced := dc.informers.WaitForCacheSync(stopCh)
	for _, gvr := range dc.gvrs {
		if !synced[*gvr] {
			log.Errorf("%v not synced!", *gvr)
			return
		}
	}
	go func() {
		for {
			select {
			case objFromInformer := <-dc.informerUpdateObjectCh:
				dc.syncUpdate(objFromInformer.GetSelfLink(), objFromInformer)
			case objFromInformer := <-dc.informerDeleteObjectCh:
				dc.syncDelete(objFromInformer.GetSelfLink(), objFromInformer)
			case <-stopCh:
				log.Infof("receive stopped, stop informer")
				return
			}
		}
	}()
	<-stopCh
	log.Infof("Stopping  dynamic rc controller...")
}

func (dc *DynamicController) syncUpdate(key string, obj interface{}) error {
	dc.transmitor.Push(events.NewEvent(key, events.OPERATION_CREATE_OR_UPDATE, obj))
	return nil
}

func (dc *DynamicController) syncDelete(key string, obj interface{}) error {
	dc.transmitor.Push(events.NewEvent(key, events.OPERATION_DELETE, obj))
	return nil
}
