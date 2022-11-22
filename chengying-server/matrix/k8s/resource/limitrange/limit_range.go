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

package limitrange

import (
	"context"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var GVK = schema.GroupVersionKind{
	Group:   "",
	Version: "v1",
	Kind:    "LimitRange",
}

func Get(ctx context.Context,client kube.Client, namespace string) (*corev1.LimitRange,error){
	if client == nil{
		return nil,fmt.Errorf("the namespace client is not exist")
	}
	limitrange := &corev1.LimitRangeList{}
	if err := client.List(ctx,limitrange,namespace);err != nil{
		return nil,err
	}
	if len(limitrange.Items) == 0{
		return nil,nil
	}
	limirange := limitrange.Items[0]
	return &limirange,nil
}
