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

package impl

import (
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"github.com/kataras/iris/context"
)

func GetParentProducts(ctx context.Context) apibase.Result{
	log.Debugf("[ns-product]: %s",ctx.Request().RequestURI)
	ns := ctx.Params().Get("namespace")
	clusterid := ctx.GetCookie(view.ClusterId)
	rsp,err := resource.GetParentProductList(ns,clusterid)
	if err != nil{
		return err
	}
	return rsp
}

func GetProducts(ctx context.Context) apibase.Result{
	log.Debugf("[ns-product]: %s",ctx.Request().RequestURI)
	parentProduct := ctx.Params().Get("parent_product_name")
	namespace := ctx.Params().Get("namespace")
	clusterid := ctx.GetCookie(view.ClusterId)
	rsp,err := resource.GetProductList(namespace,clusterid,parentProduct)
	if err != nil{
		return err
	}
	return rsp
}
//
func GetServiceList(ctx context.Context) apibase.Result{
	log.Debugf("[ns-product]: %s",ctx.Request().RequestURI)
	namespace := ctx.Params().Get("namespace")
	parentProduct := ctx.Params().Get("parent_product_name")
	productName := ctx.Params().Get("product_name")
	clusterid := ctx.GetCookie(view.ClusterId)
	rsp,err := resource.GetServiceList(namespace,clusterid,parentProduct,productName)
	if err != nil{
		return err
	}
	return rsp
}

func GetService(ctx context.Context) apibase.Result{
	log.Debugf("[ns-product]: %s",ctx.Request().RequestURI)
	namespace := ctx.Params().Get("namespace")
	parentProduct := ctx.Params().Get("parent_product_name")
	productName := ctx.Params().Get("product_name")
	servicename := ctx.Params().Get("service_name")
	clusterid := ctx.GetCookie(view.ClusterId)
	rsp,err := resource.GetService(ctx,namespace,clusterid,parentProduct,productName,servicename)
	if err != nil{
		return err
	}
	return rsp
}
