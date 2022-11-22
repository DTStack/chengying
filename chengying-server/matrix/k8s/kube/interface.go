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

package kube

import (
	"context"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	"k8s.io/apimachinery/pkg/runtime"
)

type Client interface {
	Apply(ctx context.Context, object runtime.Object) error
	Create(ctx context.Context, object runtime.Object) error
	Update(ctx context.Context, object runtime.Object) error
	Delete(ctx context.Context, object runtime.Object) error
	Get(ctx context.Context, object runtime.Object) (bool,error)
	List(ctx context.Context, object runtime.Object, namespace string) error
	Status(ctx context.Context, object runtime.Object) error
    DryRun(action base.DryRunAction,object runtime.Object) error
}

type ClientCache interface {
	Connect(connectStr,workspace string) error
	GetClient(workspace string) Client
	DeleteClient(workspace string)
	Copy() ClientCache
	//need to update?
	//UpdateClient(connectStr,workspace string) error
}

type ImportType string

func (i ImportType) String() string{
	return string(i)
}
var (
	IMPORT_KUBECONFIG              ImportType = "kubeconfig"
	IMPORT_AGENT                   ImportType = "agent"
)
