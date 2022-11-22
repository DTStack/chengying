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
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type indexerInformer struct {
	informer cache.Controller
	indexer  cache.Indexer
}
type staticController struct {
	informers   []indexerInformer
	updateQ     workqueue.RateLimitingInterface
	deleteQ     workqueue.RateLimitingInterface
	deleteCache map[string]interface{}
	transmitor  events.TransmitorInterface
}

func NewStaticController(trs events.TransmitorInterface) *staticController {
	// create the workqueue
	uQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	dQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	return &staticController{
		informers:   make([]indexerInformer, 0),
		updateQ:     uQueue,
		deleteQ:     dQueue,
		deleteCache: make(map[string]interface{}),
		transmitor:  trs,
	}
}

func (mc *staticController) Add(lw cache.ListerWatcher, object runtime.Object) {
	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the object key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the obkect than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(lw, object, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				mc.updateQ.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				mc.updateQ.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta updateQ, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				mc.deleteQ.Add(key)
				mc.deleteCache[key] = obj
			}
		},
	}, cache.Indexers{})
	indexerInformer := indexerInformer{
		informer: informer,
		indexer:  indexer,
	}
	mc.informers = append(mc.informers, indexerInformer)
}

func (mc *staticController) Run(threadiness int, stopCh <-chan struct{}) {
	// Let the workers stop when we are done
	defer mc.updateQ.ShutDown()
	defer mc.deleteQ.ShutDown()

	for _, informer := range mc.informers {
		go informer.informer.Run(stopCh)

		// Wait for all involved caches to be synced, before processing items from the updateQ is started
		if !cache.WaitForCacheSync(stopCh, informer.informer.HasSynced) {
			log.Errorf("Timed out waiting for caches to sync")
			return
		}
	}
	log.Infof("Starting default rc monitor controller...")
	for i := 0; i < threadiness; i++ {

		go wait.Until(mc.runUpdateWorker, time.Second, stopCh)
		go wait.Until(mc.runDeleteWorker, time.Second, stopCh)
	}
	<-stopCh
	log.Infof("Stopping default rc monitor controller...")
}

func (mc *staticController) processUpdateNextItem() bool {
	// Wait until there is a new item in the working updateQ
	key, quit := mc.updateQ.Get()
	if quit {
		return false
	}
	log.Infof("process update event %v", key)
	defer mc.updateQ.Done(key)
	err := mc.syncUpdate(key.(string))
	if err == nil {
		mc.updateQ.Forget(key)
		return true
	}
	log.Errorf("syncUpdate err: %v", err.Error())
	//mc.updateQ.Add(key)
	mc.deleteQ.Forget(key)

	return true
}

func (mc *staticController) processDeleteNextItem() bool {
	// Wait until there is a new item in the working updateQ
	key, quit := mc.deleteQ.Get()
	if quit {
		return false
	}
	log.Infof("process delete event %v", key)
	defer mc.deleteQ.Done(key)

	err := mc.syncDelete(key.(string))
	if err == nil {
		mc.deleteQ.Forget(key)
		return true
	}
	log.Errorf("syncDelete err: %v", err.Error())
	//mc.deleteQ.Add(key)
	mc.deleteQ.Forget(key)

	return true
}

func (mc *staticController) runUpdateWorker() {
	for mc.processUpdateNextItem() {
	}
}

func (mc *staticController) runDeleteWorker() {
	for mc.processDeleteNextItem() {
	}
}

func (mc *staticController) getIndexer(key string) cache.Indexer {
	for _, indexerInformer := range mc.informers {
		_, exists, _ := indexerInformer.indexer.GetByKey(key)
		if exists {
			return indexerInformer.indexer
		}
	}
	return nil
}

func (mc *staticController) syncUpdate(key string) error {
	//index := strings.LastIndex(key,"/")
	//workspaceid := key[index+1:]
	//key = key[:index]
	indexer := mc.getIndexer(key)
	if indexer == nil {
		return fmt.Errorf("no cache indexer for key: %v", key)
	}
	obj, exists, err := indexer.GetByKey(key)
	if err != nil {
		return err
	}
	if !exists {
		// Dependent resources are cleaned up by K8s via OwnerReferences
		log.Infof("Dependent resources are cleaned up by K8s via OwnerReferences!")
		return nil
	}
	mc.transmitor.Push(events.NewEvent(key, events.OPERATION_CREATE_OR_UPDATE, obj))
	return nil
}

func (mc *staticController) syncDelete(key string) error {
	if _, ok := mc.deleteCache[key]; !ok {
		return fmt.Errorf("no delete cache for key: %v", key)
	}
	obj := mc.deleteCache[key]

	mc.transmitor.Push(events.NewEvent(key, events.OPERATION_DELETE, obj))
	delete(mc.deleteCache, key)
	return nil
}
