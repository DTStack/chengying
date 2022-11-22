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
	operatorapis "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"encoding/json"
	"fmt"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	clientgoschema "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	sigclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var (
	Schema       = runtime.NewScheme()
	LAST_APPLIED = "dtstack/last-applied-configuration"
	Create       = DryRunAction("create")
	Update       = DryRunAction("update")
	Delete       = DryRunAction("delete")
	Patch        = DryRunAction("patch")
	All          = DryRunAction("all")
)

func init() {
	clientgoschema.AddToScheme(Schema)
	crdv1.AddToScheme(Schema)
	operatorapis.AddToScheme(Schema)
}

type ClientCache struct {
	mapper          meta.RESTMapper
	workspaceClinet map[string]*Client
	mu              sync.RWMutex
}

type Client struct {
	workspace string
	c         sigclient.Client
}

type DryRunAction string

// note: the workspace input should be understood as a connection representing different permissions
//       it is often equal with runtime object's namespace
//       but it is not strictly required that the workspace input must be equal with the runtiem object's namespace
//       eg: when useincluster config build client,we often use "" represent workspace
func (c *ClientCache) Connect(kubeconfig, workspace string) error {
	config, err := buildRestConfig(kubeconfig)
	if err != nil {
		return err
	}
	mapper, err := c.getMapper(config)
	if err != nil {
		return err
	}
	client, err := sigclient.New(config, sigclient.Options{
		Scheme: Schema,
		Mapper: mapper,
	})
	if err != nil {
		log.Errorf("[kube_client]: create sigclient error : %v", err)
		return err
	}
	//save a new client
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.workspaceClinet == nil {
		c.workspaceClinet = make(map[string]*Client)
	}
	c.workspaceClinet[workspace] = &Client{workspace, client}
	return nil
}

func (c *ClientCache) GetClient(workspace string) *Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.workspaceClinet == nil {
		return nil
	}
	return c.workspaceClinet[workspace]
}

func (c *ClientCache) DeleteClient(workspace string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.workspaceClinet == nil {
		return
	}
	delete(c.workspaceClinet, workspace)
}

func (c *Client) Status(ctx context.Context, object runtime.Object) error {
	obj := object.DeepCopyObject()
	exist, err := c.Get(ctx, obj)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}
	//original data
	bts, err := json.Marshal(object)
	if err != nil {
		log.Errorf("[kube_client]: mashal object %T error %v", object, err)
		return err
	}
	//the deepcopyed one, either the object input will be modified with the data in k8s
	if err = c.c.Status().Patch(ctx, object, sigclient.RawPatch(types.MergePatchType, bts)); err != nil {
		log.Errorf("[kube_client]: status patch fail, error: %v", err)
		return err
	}
	return nil
}

//apply: create if object is exists, patch if it doesn't exist
func (c *Client) Apply(ctx context.Context, object runtime.Object) error {
	obj := object.DeepCopyObject()
	exist, err := c.Get(ctx, obj)
	if err != nil {
		return err
	}
	if !exist {
		if err = c.Create(ctx, object); err != nil {
			return err
		}
		return nil
	}
	//if err = c.Update(ctx,object); err != nil{
	//	return err
	//}
	if err = c.Patch(ctx, object); err != nil {
		return err
	}
	return nil
}

func (c *Client) Create(ctx context.Context, object runtime.Object) error {
	metaObj, err := metaAccess(object)
	if err != nil {
		return err
	}
	if err = c.c.Create(ctx, object); err != nil {
		log.Errorf("[kube_client]: use client in workspace <%s> to create object %T in namespace <%s>,error: %v",
			c.workspace, object, metaObj.GetNamespace(), err)
		return err
	}
	return nil
}

func (c *Client) Update(ctx context.Context, object runtime.Object) error {
	metaObj, err := metaAccess(object)
	if err != nil {
		return err
	}
	if err := c.c.Update(ctx, object); err != nil {
		log.Errorf("[kube_client]: use client in workspace <%s> to update object %T in namespace <%s>,error : %v",
			c.workspace, object, metaObj.GetNamespace(), err)
		return err
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, object runtime.Object) error {
	exist, err := c.Get(ctx, object)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}
	//delete in the background
	deleteOption := sigclient.PropagationPolicy(metav1.DeletionPropagation(metav1.DeletePropagationBackground))
	if err := c.c.Delete(ctx, object, deleteOption); err != nil {
		log.Errorf("[kube_client]: use client in workspace <%s> to delete object %T,error :%v", c.workspace, object, err)
		return err
	}
	return nil
}

func (c *Client) Get(ctx context.Context, object runtime.Object) (bool, error) {
	metaObj, err := meta.Accessor(object)
	if err != nil {
		log.Errorf("[kube_client]: get action convert object %T into object meta error :%v", object, err)
		return false, err
	}
	name := metaObj.GetName()
	if len(name) == 0 {
		return false, fmt.Errorf("[kube_client]: the object's name is missing")
	}
	namespace := metaObj.GetNamespace()
	if err := c.c.Get(ctx, sigclient.ObjectKey{namespace, name}, object); err != nil {
		if err, b := err.(*errors.StatusError); b {
			if err.Status().Reason == metav1.StatusReasonNotFound {
				return false, nil
			}
		}
		log.Errorf("[kube_client]: use client in workspace <%s> to get object %T <%s> in namespace <%s>, error : %v",
			c.workspace, object, name, namespace, err)
		return false, err
	}
	return true, nil
}
func (c *Client) List(ctx context.Context, object runtime.Object, namespace string) error {
	opt := sigclient.InNamespace(namespace)
	return c.Lists(ctx, object, opt)
}

func (c *Client) Lists(ctx context.Context, object runtime.Object, opts ...sigclient.ListOption) error {

	if err := c.c.List(ctx, object, opts...); err != nil {
		log.Errorf("[kube_client]: use client in workspace <%s> to list object %T,error : %v",
			c.workspace, object, err)
		return err
	}
	return nil
}

//create patch data in the cient. server-side apply is not a necessary feature
//use application/strategic-merge-patch+json strateg to patch and create patch data
func (c *Client) Patch(ctx context.Context, object runtime.Object) error {
	metaObj, err := metaAccess(object)
	if err != nil {
		return err
	}
	//modified, err := json.Marshal(object)
	//if err != nil {
	//	log.Errorf("[kube_client]: json marshal mofified object %T , err: %v", object, err)
	//	return err
	//}
	//exist, err := c.Get(ctx, object)
	//if err != nil {
	//	return err
	//}
	//if !exist {
	//	return fmt.Errorf("[kube_client]: can't do patch action,the object <%T> %s in namespace<%s> is not exist, please create",
	//		object, metaObj.GetName(), metaObj.GetNamespace())
	//}
	//get last deploy original object
	//metaObj, err = meta.Accessor(object)
	//if err != nil {
	//	log.Errorf("[kube_client]: patch action convert get object %T into object meta error :%v", object, err)
	//	return err
	//}
	//For compatibility with previous versions, directly update object.
	//so that next time, it can do patch action
	//if metaObj.GetAnnotations() == nil || len(metaObj.GetAnnotations()[LAST_APPLIED]) == 0 {
	//	json.Unmarshal(modified, object)
	//	if err = c.Update(ctx, object); err != nil {
	//		return err
	//	}
	//}
	//// get last deploy orignal object
	//ori := metaObj.GetAnnotations()[LAST_APPLIED]
	//
	//patchdata, err := strategicpatch.CreateTwoWayMergePatch([]byte(ori), modified, object)
	//if err != nil {
	//	log.Errorf("[kube_client]: create mergepatch error: %v \n ori: \n %s \n modified: \n %s",
	//		err, ori, string(modified))
	//	return err
	//}
	bts, err := json.Marshal(object)
	if err != nil {
		log.Errorf("[kube_client]: mashal object %T error %v", object, err)
		return err
	}
	if err = c.c.Patch(ctx, object, sigclient.RawPatch(types.MergePatchType, bts)); err != nil {
		log.Errorf("[kube_client]: use client in workspace <%s> to patch object %T in namespace <%s>, error: %v ",
			c.workspace, object, metaObj.GetNamespace(), err)
		return err
	}
	return nil
}

func (c *Client) DryRun(action DryRunAction, object runtime.Object) error {
	ctx := context.Background()
	switch action {
	case Create:
		return c.dryRunCreate(ctx, object)
	case Update:
		return c.dryRunUpdate(ctx, object)
	case Delete:
		return c.dryRunDelete(ctx, object)
	case Patch:
		return c.dryRunPatch(ctx, object)
	case All:
		if err := c.dryRunCreate(ctx, object); err != nil {
			return err
		}
		if err := c.dryRunUpdate(ctx, object); err != nil {
			return err
		}
		if err := c.dryRunPatch(ctx, object); err != nil {
			return err
		}
		return c.dryRunDelete(ctx, object)
	}
	return nil
}

func (c *Client) dryRunCreate(ctx context.Context, object runtime.Object) error {
	opt := sigclient.DryRunAll
	err := c.c.Create(ctx, object, opt)
	if err != nil {
		log.Errorf("[kube_client]: dryrun create fail error :%v", err)
		return err
	}
	return nil
}

func (c *Client) dryRunDelete(ctx context.Context, object runtime.Object) error {
	opt := sigclient.DryRunAll
	err := c.c.Delete(ctx, object, opt)
	if err != nil {
		log.Errorf("[kube_client]: dryrun delete fail error :%v", err)
		return err
	}
	return nil
}

func (c *Client) dryRunUpdate(ctx context.Context, object runtime.Object) error {
	opt := sigclient.DryRunAll
	err := c.c.Update(ctx, object, opt)
	if err != nil {
		log.Errorf("[kube_client]: dryrun update fail error :%v", err)
		return err
	}
	return nil
}

func (c *Client) dryRunPatch(ctx context.Context, obeject runtime.Object) error {
	opt := sigclient.DryRunAll
	err := c.c.Patch(ctx, obeject, sigclient.RawPatch(types.StrategicMergePatchType, []byte{}), opt)
	if err != nil {
		log.Errorf("[kube_client]: dryrun patch fail error :%v", err)
		return err
	}
	return nil
}

//func (c *ClientCache) isNamespaced(object runtime.Object) (bool,error){
//	gvk,err := apiutil.GVKForObject(object,myschema)
//	if err != nil{
//		log.Errorf("[kube_client]: get gvk of Object %T error : %v",object,err)
//		return false ,err
//	}
//	isnamespaced,exist := c.namespacedResource[gvk]
//	if exist{
//		return isnamespaced,nil
//	}
//	mapping,err :=c.mapper.RESTMapping(gvk.GroupKind(),gvk.Version)
//	if err != nil{
//		log.Errorf("[kube_client]: get mapping by gvk error : %v",err)
//		return false ,err
//	}
//	//it's not thread safe, but it doesn't matter
//	isnamespaced = mapping.Scope.Name() != meta.RESTScopeNameRoot
//	c.namespacedResource[gvk] = isnamespaced
//	return isnamespaced,nil
//}

// it is used to get resoucename and ensure if is namespaced later.
// a cluster only need a mapper
func (c *ClientCache) getMapper(config *rest.Config) (meta.RESTMapper, error) {
	if c.mapper == nil {
		m, err := restMapper(config)
		if err != nil {
			return nil, err
		}
		c.mapper = m
	}
	return c.mapper, nil
}
func metaAccess(object runtime.Object) (metav1.Object, error) {
	obj, err := meta.Accessor(object)
	if err != nil {
		log.Errorf("[kube_client] convert object %T into object meta, error : %v", object, err)
		return nil, err
	}
	return obj, nil
}

//add annotation "dtstack/last-applied-configuration" for action patch to compare
func addLastApplied(object runtime.Object) (metav1.Object, error) {
	obj, err := meta.Accessor(object)
	if err != nil {
		log.Errorf("[kube_client] convert object %T into object meta, error : %v", object, err)
		return nil, err
	}
	return obj, nil
	bts, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("[kube_client] json marshal object %T,error : %v", object, err)
		return nil, err
	}
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string, 1)
		obj.SetAnnotations(annotations)
	}
	annotations[LAST_APPLIED] = string(bts)
	return obj, nil
}

//if kubeconfig == "",use inclusterconfig
//if kubeconfig != "",use kubeconfig to build restconfig, and use currentcontext
func buildRestConfig(kubeconfig string) (*rest.Config, error) {
	if len(kubeconfig) == 0 {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Errorf("[kube_client]: create inclusterconfig error :%v", err)
			return nil, err
		}
		return config, nil
	}
	apiconfig, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		log.Errorf("[kube_client]: load kubeconfig error : %v", err)
		return nil, err
	}
	config, err := clientcmd.NewNonInteractiveClientConfig(*apiconfig, "", &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		log.Errorf("[kube_client]: create rest config from kubeconfig error :%v", err)
		return nil, err
	}
	return config, nil
}

func restMapper(config *rest.Config) (meta.RESTMapper, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Errorf("[kube_client]: create discovery client error : %v", err)
		return nil, err
	}
	grs, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		log.Errorf("[kube_client]: get apigroupresources by discoveryclient error : %v", err)
		return nil, err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(grs)
	return mapper, nil
}
