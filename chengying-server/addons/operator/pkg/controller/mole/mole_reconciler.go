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

package mole

import (
	molev1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/mole/v1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/common"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/model"
)

type MoleReconciler struct {
	DsHash     string
	ConfigHash string
	PluginsEnv string
	Name       string
}

func NewMoleReconciler(name string) *MoleReconciler {
	return &MoleReconciler{
		DsHash:     "",
		ConfigHash: "",
		PluginsEnv: "",
		Name:       name,
	}
}

func (i *MoleReconciler) Reconcile(state *common.ServiceState, cr *molev1.Mole) common.DesiredServiceState {
	desired := common.DesiredServiceState{}
	if cr.Spec.Product.Service[i.Name].IsJob {
		return desired.AddAction(i.jobReconclie(state, cr))
	}
	desired = desired.AddAction(i.getMoleDeploymentDesiredState(state, cr))
	desired = desired.AddAction(i.getMoleServiceDesiredState(state, cr))
	if cr.Spec.Product.Service[i.Name].IsDeployIngress {
		desired = desired.AddAction(i.getMoleIngressDesiredState(state, cr))
	}
	return desired
}

func (i *MoleReconciler) jobReconclie(state *common.ServiceState, cr *molev1.Mole) common.ServiceAction {
	if state.MoleJob == nil {
		return common.GenericCreateAction{
			Ref: model.MoleJob(cr, i.Name),
			Msg: "create Mole Job",
		}
	}
	return nil
}

func (i *MoleReconciler) getMoleServiceDesiredState(state *common.ServiceState, cr *molev1.Mole) common.ServiceAction {
	if state.MoleService == nil {
		return common.GenericCreateAction{
			Ref: model.MoleService(cr, i.Name),
			Msg: "create Mole service",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.MoleServiceReconciled(cr, state.MoleService, i.Name),
		Msg: "update Mole service",
	}
}

func (i *MoleReconciler) getMoleIngressDesiredState(state *common.ServiceState, cr *molev1.Mole) common.ServiceAction {
	if state.MoleIngress == nil {
		return common.GenericCreateAction{
			Ref: model.MoleIngress(cr, i.Name),
			Msg: "create Mole ingress",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.MoleIngressReconciled(cr, state.MoleIngress, i.Name),
		Msg: "update Mole ingress",
	}
}

func (i *MoleReconciler) getMoleDeploymentDesiredState(state *common.ServiceState, cr *molev1.Mole) common.ServiceAction {
	if state.MoleDeployment == nil {
		return common.GenericCreateAction{
			Ref: model.MoleDeployment(cr, i.Name),
			Msg: "create Mole deployment",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.MoleDeploymentReconciled(cr, state.MoleDeployment, i.Name),
		Msg: "update Mole deployment",
	}
}
