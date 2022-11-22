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

package reconciler

import (
	"context"
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/internal"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/util"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
)

const LAST_WITHOUT_ANNOTATION = "last_without_annotation"

var log = logf.Log.WithName("last_apply_reconciler")

type LastApplyReconciler struct {
	Name     string
	Version  string
	Group    string
	Owner    *workloadv1beta1.WorkloadProcess
	Workload *workloadv1beta1.WorkLoad
	Schema   *runtime.Scheme
}

func (r *LastApplyReconciler) Reconcile(ctx context.Context, c client.Client, recorder record.EventRecorder) error {

	r.Workload.Namespace = r.Owner.Namespace
	if len(r.Workload.Name) == 0 {
		r.Workload.Name = util.BuildWorkloadName(r.Owner.Name, r.Name)
	}
	key := types.NamespacedName{
		Namespace: r.Workload.Namespace,
		Name:      r.Workload.Name,
	}
	operationResult, err := r.createOrUpdate(ctx, c)
	if err != nil {
		recorder.Event(
			r.Owner,
			internal.EventTypeWarning,
			fmt.Sprintf("%sSyncFailed", r.Name),
			fmt.Sprintf("%s %s sync failed, error: %v", r.ObjectType(), key, err),
		)
		return err
	}
	if operationResult == internal.OperationResultNone {
		return nil
	}
	recorder.Event(
		r.Owner,
		internal.EventTypeNormal,
		fmt.Sprintf("%sSyncSuccess", r.Name),
		fmt.Sprintf("%s %s %s sync success", r.ObjectType(), key, operationResult),
	)
	return nil
}

// the error occures when talk with api-server ignore
func (r *LastApplyReconciler) createOrUpdate(ctx context.Context, c client.Client) (string, error) {
	existing := &workloadv1beta1.WorkLoad{}
	key, err := client.ObjectKeyFromObject(r.Workload)
	if err != nil {
		log.Error(err, "get namespacename from Object fail", "kind", r.ObjectType())
		return internal.OperationResultNone, err
	}
	//create object
	if err = c.Get(ctx, key, existing); err != nil {
		if !errors.IsNotFound(err) {
			return internal.OperationResultNone, nil
		}

		if err = r.Mutate(existing); err != nil {
			fmt.Println("------", err.Error())
			return internal.OperationResultNone, err
		}

		if err = c.Create(ctx, r.Workload); err != nil {
			return internal.OperationResultNone, err
		}
		return internal.OperationResultCreate, nil
	}

	lastApply := existing.Annotations[LAST_WITHOUT_ANNOTATION]
	nowApply, err := json.Marshal(r.Workload)
	if err != nil {
		log.Error(err, "unmashal workload fail")
		return internal.OperationResultNone, err
	}

	//the product deploy need deployuuid and the uuid is used in the pod, but it is changed every service in product when deploy
	//even if the service has no changes before.
	//so compare if changed without the annotation deployuuid
	if lastApply == string(nowApply) {
		return internal.OperationResultNone, nil
	}
	if err = r.Mutate(existing); err != nil {
		return internal.OperationResultNone, err
	}

	//update
	if err = c.Update(ctx, existing); err != nil {
		return internal.OperationResultNone, err
	}

	return internal.OperationResultUpdate, nil
}

func (r *LastApplyReconciler) ObjectType() string {
	return fmt.Sprintf("%T", r.Workload)
}

// add the lastapply first
// then add the deployuuid and the other product info.
// final add the ownerReferences
func (r *LastApplyReconciler) Mutate(existing *workloadv1beta1.WorkLoad) error {
	lastApply, err := json.Marshal(r.Workload)
	if err != nil {
		log.Error(err, "json marshal fail", "kind", r.ObjectType())
		return err
	}
	if r.Workload.Annotations == nil {
		r.Workload.Annotations = map[string]string{}
	}
	r.Workload.Annotations[LAST_WITHOUT_ANNOTATION] = string(lastApply)

	for _, part := range r.Workload.Spec.WorkLoadParts {
		params := part.BaseWorkLoad.Parameters.Object

		productInfo := map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"deploy_uuid":         r.Owner.Spec.DeployUUId,
							"parent_product_name": r.Owner.Spec.ParentProductName,
							"pid":                 strconv.Itoa(r.Owner.Spec.ProductId),
							"product_name":        r.Owner.Spec.ProductName,
							"product_version":     r.Owner.Spec.ProductVersion,
							"cluster_id":          strconv.Itoa(r.Owner.Spec.ClusterId),
							"com":                 "dtstack.com",
							"service_name":        r.Name,
							"service_version":     r.Version,
							"group":               r.Group,
						},
					},
				},
			},
		}
		util.DeepCopy(productInfo, params)
	}
	if err = controllerutil.SetControllerReference(r.Owner, r.Workload, r.Schema); err != nil {
		log.Error(err, "workload set owner workloadprocess fail")
		return err
	}

	existingMap, err := util.ToMap(existing)
	if err != nil {
		log.Error(err, "exist workload to map fail")
		return err
	}
	desired, err := util.ToMap(r.Workload)
	if err != nil {
		log.Error(err, "desired workload to map fail")
		return err
	}
	syncd, err := util.DeepCopy(desired, existingMap)
	if err != nil {
		log.Error(err, "copy desired to existing workload fail")
		return err
	}
	if err = util.ToObject(syncd, existing); err != nil {
		bts, _ := json.Marshal(syncd)
		log.Error(err, "to workload fail", "json", string(bts))
		return err
	}
	return nil
}

func newMapIfNil(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	return m
}
