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
	molev1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/mole/v1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/model"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceState struct {
	Name           string
	MoleConfig     *corev1.ConfigMap
	MoleIngress    *v1beta1.Ingress
	MoleService    *corev1.Service
	MoleDeployment *appsv1.Deployment
	MoleJob        *batchv1.Job
}

func NewServiceState(name string) *ServiceState {
	return &ServiceState{
		Name: name,
	}
}

func (i *ServiceState) Read(ctx context.Context, cr *molev1.Mole, client client.Client) error {
	//if job, no deployment,service,ingress
	if cr.Spec.Product.Service[i.Name].IsJob {
		return i.readMoleJob(ctx, cr, client)
	}
	err := i.readMoleDeployment(ctx, cr, client, i.Name)
	if err != nil {
		return err
	}
	err = i.readMoleService(ctx, cr, client, i.Name)
	if err != nil {
		return err
	}
	if cr.Spec.Product.Service[i.Name].IsDeployIngress {
		err = i.readMoleIngress(ctx, cr, client)
	}

	return err
}

func (i *ServiceState) readMoleJob(ctx context.Context, cr *molev1.Mole, reader client.Reader) error {
	currentState := &batchv1.Job{}
	selector := model.MoleJobSelector(cr, i.Name)
	err := reader.Get(ctx, selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	i.MoleJob = currentState.DeepCopy()
	return nil
}

func (i *ServiceState) readMoleService(ctx context.Context, cr *molev1.Mole, client client.Client, name string) error {
	//currentState := model.MoleService(cr, name)
	currentState := &corev1.Service{}
	selector := model.MoleServiceSelector(cr, name)
	err := client.Get(ctx, selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	i.MoleService = currentState.DeepCopy()

	return nil
}

func (i *ServiceState) readMoleIngress(ctx context.Context, cr *molev1.Mole, client client.Client) error {
	//currentState := model.MoleIngress(cr, i.Name)
	currentState := &v1beta1.Ingress{}
	selector := model.MoleIngressSelector(cr, i.Name)
	err := client.Get(ctx, selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	i.MoleIngress = currentState.DeepCopy()
	return nil
}

func (i *ServiceState) readMoleDeployment(ctx context.Context, cr *molev1.Mole, client client.Client, name string) error {
	//currentState := model.MoleDeployment(cr, name)
	currentState := &appsv1.Deployment{}
	selector := model.MoleDeploymentSelector(cr, name)
	err := client.Get(ctx, selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	i.MoleDeployment = currentState.DeepCopy()
	return nil
}
