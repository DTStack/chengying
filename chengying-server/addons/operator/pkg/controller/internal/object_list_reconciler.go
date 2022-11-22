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
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

type ObjectListReconciler struct {
	Reconcilers map[int]Reconciler
}

func (list *ObjectListReconciler) Reconcile(ctx context.Context, c client.Client, recorder record.EventRecorder) error {
	reconcilers := list.sort()
	for _, r := range reconcilers {
		if err := r.Reconcile(ctx, c, recorder); err != nil {
			return err
		}
	}
	return nil
}

func (list *ObjectListReconciler) Append(r Reconciler, index int) {
	if list.Reconcilers == nil {
		list.Reconcilers = map[int]Reconciler{}
	}
	list.Reconcilers[index] = r
}

func (list *ObjectListReconciler) sort() []Reconciler {
	keys := make([]int, 0, len(list.Reconcilers))
	for k := range list.Reconcilers {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	reconcilers := make([]Reconciler, 0, len(list.Reconcilers))
	for _, k := range keys {
		reconcilers = append(reconcilers, list.Reconcilers[k])
	}
	return reconcilers
}
