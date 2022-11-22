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

package common

import (
	"context"
	stdErr "errors"
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type ActionRunner interface {
	RunAll(desiredState DesiredServiceState) error
	create(obj runtime.Object) error
	update(obj runtime.Object) error
	delete(obj runtime.Object) error
	ingressReady(obj runtime.Object) error
	deploymentReady(obj runtime.Object) error
}

type ServiceAction interface {
	Run(runner ActionRunner) (string, error)
}

// The desired cluster state is defined by a list of actions that have to be run to
// get from the current state to the desired state
type DesiredServiceState []ServiceAction

func (d *DesiredServiceState) AddAction(action ServiceAction) DesiredServiceState {
	if action != nil {
		*d = append(*d, action)
	}
	return *d
}

func (d *DesiredServiceState) AddActions(actions []ServiceAction) DesiredServiceState {
	for _, action := range actions {
		d.AddAction(action)
	}
	return *d
}

type ServiceActionRunner struct {
	scheme *runtime.Scheme
	client client.Client
	ctx    context.Context
	log    logr.Logger
	cr     runtime.Object
}

func NewServiceActionRunner(ctx context.Context, client client.Client, scheme *runtime.Scheme, cr runtime.Object) ActionRunner {
	return &ServiceActionRunner{
		scheme: scheme,
		client: client,
		log:    logf.Log.WithName("action-runner"),
		ctx:    ctx,
		cr:     cr,
	}
}

func (i *ServiceActionRunner) RunAll(desiredState DesiredServiceState) error {
	for index, action := range desiredState {

		msg, err := action.Run(i)
		if err != nil {
			i.log.Info(fmt.Sprintf("(%5d) %10s %s", index, "FAILED", msg))
			i.log.Info(fmt.Sprintf("err:%v", err))
			return err
		}
		i.log.Info(fmt.Sprintf("(%5d) %10s %s", index, "SUCCESS", msg))
	}

	return nil
}

func (i *ServiceActionRunner) create(obj runtime.Object) error {
	err := controllerutil.SetControllerReference(i.cr.(metav1.Object), obj.(metav1.Object), i.scheme)
	if err != nil {
		return err
	}
	return i.client.Create(i.ctx, obj)
}

func (i *ServiceActionRunner) update(obj runtime.Object) error {
	err := controllerutil.SetControllerReference(i.cr.(metav1.Object), obj.(metav1.Object), i.scheme)
	if err != nil {
		return err
	}
	err = i.client.Update(i.ctx, obj)
	if err != nil {
		// Update conflicts can happen frequently when kubernetes updates the resource
		// in the background
		if errors.IsConflict(err) {
			fmt.Println("conflict----------------------------", err.Error())
			return nil
		}
		return err
	}
	return nil
}

func (i *ServiceActionRunner) delete(obj runtime.Object) error {
	return i.client.Delete(i.ctx, obj)
}

func (i *ServiceActionRunner) ingressReady(obj runtime.Object) error {
	ready := IsIngressReady(obj.(*v1beta1.Ingress))
	if !ready {
		return stdErr.New("ingress not ready")
	}
	return nil
}

func (i *ServiceActionRunner) deploymentReady(obj runtime.Object) error {
	ready, err := IsDeploymentReady(obj.(*appsv1.Deployment))
	if err != nil {
		return err
	}

	if !ready {
		return stdErr.New("deployment not ready")
	}
	return nil
}

// An action to create generic kubernetes resources
// (resources that don't require special treatment)
type GenericCreateAction struct {
	Ref runtime.Object
	Msg string
}

// An action to update generic kubernetes resources
// (resources that don't require special treatment)
type GenericUpdateAction struct {
	Ref runtime.Object
	Msg string
}

type LogAction struct {
	Msg string
}

type IngressReadyAction struct {
	Ref runtime.Object
	Msg string
}

type DeploymentReadyAction struct {
	Ref runtime.Object
	Msg string
}

// An action to delete generic kubernetes resources
// (resources that don't require special treatment)
type GenericDeleteAction struct {
	Ref runtime.Object
	Msg string
}

func (i GenericCreateAction) Run(runner ActionRunner) (string, error) {
	return i.Msg, runner.create(i.Ref)
}

func (i GenericUpdateAction) Run(runner ActionRunner) (string, error) {

	return i.Msg, runner.update(i.Ref)
}

func (i GenericDeleteAction) Run(runner ActionRunner) (string, error) {
	return i.Msg, runner.delete(i.Ref)
}

func (i LogAction) Run(runner ActionRunner) (string, error) {
	return i.Msg, nil
}

func (i IngressReadyAction) Run(runner ActionRunner) (string, error) {
	return i.Msg, runner.ingressReady(i.Ref)
}

func (i DeploymentReadyAction) Run(runner ActionRunner) (string, error) {
	return i.Msg, runner.deploymentReady(i.Ref)
}
