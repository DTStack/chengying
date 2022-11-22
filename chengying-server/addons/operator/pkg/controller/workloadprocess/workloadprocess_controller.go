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

package workloadprocess

import (
	"context"
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workloadprocess/reconciler"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const controllerName = "workloadprocess-controller"

var log = logf.Log.WithName(controllerName)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconcile(mgr))
}

func newReconcile(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWorkloadProcess{
		Client:   mgr.GetClient(),
		Schema:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor(controllerName),
	}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	ctl, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		log.Error(err, "create workloadprocess controller fail")
		return err
	}
	if err = ctl.Watch(
		&source.Kind{
			Type: &workloadv1beta1.WorkloadProcess{},
		},
		&handler.EnqueueRequestForObject{},
		predicate.ResourceVersionChangedPredicate{},
	); err != nil {
		log.Error(err, "watch resource workloadprocess fail")
		return err
	}

	if err = ctl.Watch(
		&source.Kind{
			Type: &workloadv1beta1.WorkLoad{},
		},
		&handler.EnqueueRequestForOwner{
			OwnerType:    &workloadv1beta1.WorkloadProcess{},
			IsController: true,
		},
		predicate.ResourceVersionChangedPredicate{},
	); err != nil {
		log.Error(err, "watch controllered resource workload fail")
		return err
	}
	return nil
}

type ReconcileWorkloadProcess struct {
	client.Client
	Schema   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *ReconcileWorkloadProcess) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	wp := &workloadv1beta1.WorkloadProcess{}
	err := r.Get(context.TODO(), request.NamespacedName, wp)
	if err != nil {
		// it is for delete action, if the workload is deleted, the reconcile will be forgot.
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		//return error, it will be log 'Reconciler error' with controllername and request and error, and then requeue
		return reconcile.Result{}, err
	}

	if err = r.manageStatus(wp); err != nil {
		return reconcile.Result{}, err
	}
	// so that you can directly update workload.
	// only the deployuuid changed, would do the next thing
	if wp.Spec.DeployUUId == wp.Spec.LastDeployUUId {
		return reconcile.Result{}, nil
	}
	for svcName, svcWorkload := range wp.Spec.WorkLoads {
		rc := &reconciler.LastApplyReconciler{
			Name:     svcName,
			Version:  svcWorkload.Version,
			Group:    svcWorkload.Group,
			Owner:    wp,
			Workload: &svcWorkload.WorkLoad,
			Schema:   r.Schema,
		}
		if err = rc.Reconcile(context.TODO(), r, r.Recorder); err != nil {
			wp.Status = workloadv1beta1.ProcessStatus{Phase: workloadv1beta1.ProcessFail}
			r.Status().Update(context.TODO(), wp)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileWorkloadProcess) manageStatus(process *workloadv1beta1.WorkloadProcess) error {
	deployedCount := len(process.Spec.WorkLoads)
	readyCount := 0
	for serviceName, svcWokrload := range process.Spec.WorkLoads {
		workload := svcWokrload.WorkLoad
		name := workload.Name
		if len(name) == 0 {
			name = util.BuildWorkloadName(process.Name, serviceName)
		}
		namespace := process.Namespace
		key := types.NamespacedName{Name: name, Namespace: namespace}
		err := r.Get(context.TODO(), key, &workload)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			log.Error(err, "manage status fail")
			return err
		}
		if workload.Status.Phase == workloadv1beta1.WorkloadRunning {
			readyCount++
		}
	}
	if deployedCount == readyCount {
		process.Status = workloadv1beta1.ProcessStatus{Phase: workloadv1beta1.ProcessFinish}
	} else {
		process.Status = workloadv1beta1.ProcessStatus{Phase: workloadv1beta1.ProcessPending}
	}
	r.Status().Update(context.TODO(), process)
	return nil
}
