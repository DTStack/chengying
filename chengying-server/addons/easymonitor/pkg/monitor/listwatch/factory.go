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

package listwatch

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type listWatchFactory struct {
	listWatch cache.ListerWatcher
}

// WithCustomObject create listwatch with custom kuber client & object & namespace & fieldSelector.
func WithClientResourceNsLabelSelector(client kubernetes.Clientset, resource, namespace string, labelsSelector labels.Selector) cache.ListerWatcher {
	optionsModifier := func(options *v1.ListOptions) {
		options.LabelSelector = labelsSelector.String()
	}
	listWatcher := cache.NewFilteredListWatchFromClient(client.CoreV1().RESTClient(), resource, namespace, optionsModifier)
	return listWatcher
}

// WithCustomObject create listwatch with custom kuber client & object & namespace & fieldSelector.
func WithConfigResourceNslabelSelector(getter rest.Interface, master, kubeconfig, resource, namespace string, labelsSelector labels.Selector) (error, cache.ListerWatcher) {
	optionsModifier := func(options *v1.ListOptions) {
		options.LabelSelector = labelsSelector.String()
	}
	listWatcher := cache.NewFilteredListWatchFromClient(getter, resource, namespace, optionsModifier)
	return nil, listWatcher
}
