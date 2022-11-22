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
	"context"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	molev1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/mole/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var GVK = schema.GroupVersionKind{
	Group:   "operator.dtstack.com",
	Version: "v1",
	Kind:    "Mole",
}
func New() *molev1.Mole{
	return &molev1.Mole{}
}
func Ping(client kube.Client, namespace string) error {
	ping := &molev1.Mole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dtstack-dryru",
			Namespace: namespace,
		},
		Spec: molev1.MoleSpec{
			Product: molev1.SchemaConfig{
				Service: map[string]molev1.ServiceConfig{
					"ping": molev1.ServiceConfig{},
				},
			},
		},
	}
	if _, err := client.Get(context.Background(), ping); err != nil {
		return err
	}

	if err := client.DryRun(base.Create, ping); err != nil {
		return err
	}
	return nil
}
