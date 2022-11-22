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

package workload

import (
	"context"
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workload/base"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workload/support"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const controllerName = "workload-controller"

var log = logf.Log.WithName(controllerName)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconcile(mgr))
}

func newReconcile(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWorkload{
		CtrlInfo: &base.CtrlInfo{
			Client:   mgr.GetClient(),
			Recorder: mgr.GetEventRecorderFor(controllerName),
			Schema:   mgr.GetScheme(),
		},
	}
}

// watch resource workload and the support resource
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	ctl, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		log.Error(err, "create workload controller fail")
		return err
	}
	if err = ctl.Watch(
		&source.Kind{
			Type: &workloadv1beta1.WorkLoad{},
		},
		&handler.EnqueueRequestForObject{},
		predicate.ResourceVersionChangedPredicate{},
	); err != nil {
		log.Error(err, "watch resource workload fail")
		return err
	}

	objs, err := support.GetTypes()
	if err != nil {
		return err
	}
	for k, obj := range objs {
		if err := ctl.Watch(
			&source.Kind{
				Type: obj,
			},
			&handler.EnqueueRequestForOwner{
				OwnerType:    &workloadv1beta1.WorkLoad{},
				IsController: true,
			},
			predicate.ResourceVersionChangedPredicate{},
		); err != nil {
			log.Error(err, "watch support type fail", "type", k)
			return err
		}
	}
	return nil
}

type ReconcileWorkload struct {
	*base.CtrlInfo
}

func (r *ReconcileWorkload) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	wl := &workloadv1beta1.WorkLoad{}
	err := r.Get(context.TODO(), request.NamespacedName, wl)
	if err != nil {
		// it is for delete action, if the workload is deleted, the reconcile will be forgot.
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		//return error, it will be log 'Reconciler error' with controllername and request and error, and then requeue
		return reconcile.Result{}, err
	}

	for _, part := range wl.Spec.WorkLoadParts {

		rc, err := base.GetObjectListReconciler(r.CtrlInfo, wl, part)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err = rc.Reconcile(context.TODO(), r.Client, r.Recorder); err != nil {
			r.manageFail(wl, &part, err)
			return reconcile.Result{}, err
		}
		// the baseworkload judge the status it self
		base.Status(r.CtrlInfo, wl, part)
	}
	r.manageStauts(wl)
	return reconcile.Result{}, nil
}

func (r *ReconcileWorkload) manageFail(wl *workloadv1beta1.WorkLoad, part *workloadv1beta1.WorkLoadPart, err error) {
	wl.Status.Phase = workloadv1beta1.WorkloadFail
	if wl.Status.PartStatus == nil {
		wl.Status.PartStatus = map[string]workloadv1beta1.WorkLoadPartStatus{}
	}
	wl.Status.PartStatus[part.BaseWorkLoad.Name] = workloadv1beta1.WorkLoadPartStatus{
		Status:  workloadv1beta1.WorkloadFail,
		Message: err.Error(),
	}
	r.Status().Update(context.TODO(), wl)
}

func (r *ReconcileWorkload) manageStauts(wl *workloadv1beta1.WorkLoad) {
	workloadStatus := workloadv1beta1.WorkloadPending
	if len(wl.Status.PartStatus) == len(wl.Spec.WorkLoadParts) {
		ifPending := false
		for _, ps := range wl.Status.PartStatus {
			if ps.Status == workloadv1beta1.WorkloadPending {
				ifPending = true
			}
		}
		if !ifPending {
			workloadStatus = workloadv1beta1.WorkloadRunning
		}
	}
	wl.Status.Phase = workloadStatus
	r.Status().Update(context.TODO(), wl)
}
