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

package apibase

// 1: Administrator, 2: Cluster Operator, 4: Cluster Reader

// permission7 = Administrator & Cluster Operator & Cluster Reader
// permission3 = Administrator & Cluster Operator
// permission1 = Administrator

import (
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

const (
	Administrator   = 1
	ClusterOperator = 2
	ClusterReader   = 4
)

func CheckPermission3(ctx context.Context) {
	rulePermission := Administrator | ClusterOperator
	userPermission, err := GetTokenUserPermission(ctx)
	if err != nil {
		log.Errorf(err.Error())
	}
	if rulePermission&userPermission == 0 {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	ctx.Next()
}

func CheckPermission1(ctx context.Context) {
	rulePermission := Administrator
	userPermission, err := GetTokenUserPermission(ctx)
	if err != nil {
		log.Errorf(err.Error())
	}
	if rulePermission&userPermission == 0 {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	ctx.Next()
}
