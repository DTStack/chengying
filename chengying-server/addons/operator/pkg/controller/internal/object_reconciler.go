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

package internal

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	OperationResultCreate  = "created"
	OperationResultUpdate  = "updated"
	OperationResultNone    = "unchanged"
	OperationResultRestart = "delete_to_restart"
	EventTypeWarning       = "Warning"
	EventTypeNormal        = "Normal"
)

var log = logf.Log.WithName("object-reconciler")

// the object's desired state if comes form the MutateFn
type ObjectReconciler struct {
	Name     string
	Object   runtime.Object
	MutateFn MutateFunction
	Owner    runtime.Object
	Schema   *runtime.Scheme
}

// does the actual reconcile
func (r *ObjectReconciler) Reconcile(ctx context.Context, c client.Client, recorder record.EventRecorder) error {
	key, err := client.ObjectKeyFromObject(r.Object)
	if err != nil {
		log.Error(err, "get namespaced from object fail", "kind", r.ObjectType())
		return err
	}
	operationResult, err := r.createOrUpdate(ctx, c, r.mutateFuncEnhance())
	if err != nil {
		recorder.Event(
			r.Owner,
			EventTypeWarning,
			fmt.Sprintf("%sSyncFailed", r.Name),
			fmt.Sprintf("%s %s sync failed, error: %v", r.ObjectType(), key, err),
		)
		return err
	}
	if operationResult == OperationResultNone {
		return nil
	}
	recorder.Event(
		r.Owner,
		EventTypeNormal,
		fmt.Sprintf("%sSyncSuccess", r.Name),
		fmt.Sprintf("%s %s %s sync success", r.ObjectType(), key, operationResult),
	)
	return nil
}

// accountding to the object existed in the k8s and the object given from the workload crd,
// determine whether the object need to create, restart or update
// the error when talk with apiserver ignore
func (r *ObjectReconciler) createOrUpdate(ctx context.Context, c client.Client, fn MutateFunction) (string, error) {
	obj := r.Object
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		log.Error(err, "get namespacename from Object fail", "kind", r.ObjectType())
		return OperationResultNone, err
	}
	//create object
	if err = c.Get(ctx, key, obj); err != nil {
		if !errors.IsNotFound(err) {
			return OperationResultNone, nil
		}
		//mutate object
		if _, err = fn(obj); err != nil {
			return OperationResultNone, err
		}
		if err = c.Create(ctx, obj); err != nil {
			return OperationResultNone, err
		}
		return OperationResultCreate, nil
	}

	existing := obj.DeepCopyObject()
	//mutate object, and determine whether need to restart
	restart, err := fn(obj)
	if err != nil {
		return OperationResultNone, err
	}
	//if the object is not changed, but figure out to restart, delete the object to restart
	if reflect.DeepEqual(obj, existing) {
		if !restart {
			return OperationResultNone, nil
		}
		if err = c.Delete(ctx, obj); err != nil {
			return OperationResultNone, err
		}
		return OperationResultRestart, nil
	}
	//update
	if err = c.Update(ctx, obj); err != nil {
		return OperationResultNone, err
	}

	return OperationResultUpdate, nil
}

func (r *ObjectReconciler) ObjectType() string {
	return fmt.Sprintf("%T", r.Object)
}

//set ownerrefrence
func (r *ObjectReconciler) mutateFuncEnhance() MutateFunction {
	return func(object runtime.Object) (bool, error) {
		restart, err := r.MutateFn(object)
		if err != nil {
			return false, err
		}
		if r.Owner != nil {
			desired, err := meta.Accessor(object)
			if err != nil {
				log.Error(err, "convert existing object to metav1.Object fail", "kind", r.ObjectType())
				return false, err
			}
			owner, err := meta.Accessor(r.Owner)
			if err != nil {
				log.Error(err, "convert owner to metav1.Object fail", "kind", fmt.Sprintf("%T", r.Owner))
				return false, err
			}
			err = controllerutil.SetControllerReference(owner, desired, r.Schema)
			if err != nil {
				log.Error(err, "set owner fail", "object-kind", r.ObjectType(), "owner-type", fmt.Sprintf("%T", r.Owner))
				return false, err
			}
			return restart, nil
		}
		return false, nil
	}
}
