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
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/monitor"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/view/request"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/context"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	sigclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Apply(ctx context.Context) apibase.Result {

	res := &request.Resource{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] apply read json error : %v", err)
	}
	obj, err := convertToObject(res)
	if err != nil {
		return err
	}
	if err = clientCache.GetClient("").Apply(ctx, obj); err != nil {
		return err
	}
	return nil
}

func Create(ctx context.Context) apibase.Result {
	res := &request.Resource{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] patch read json error : %v", err)
	}
	obj, err := convertToObject(res)
	if err != nil {
		return err
	}
	if err = clientCache.GetClient("").Create(ctx, obj); err != nil {
		return err
	}
	return nil
}

func Update(ctx context.Context) apibase.Result {
	res := &request.Resource{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] update read json error : %v", err)
	}
	obj, err := convertToObject(res)
	if err != nil {
		return err
	}
	if err = clientCache.GetClient("").Update(ctx, obj); err != nil {
		return err
	}
	return nil
}

func Delete(ctx context.Context) apibase.Result {
	res := &request.Resource{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] delete read json error : %v", err)
	}
	obj, err := convertToObject(res)
	if err != nil {
		return err
	}

	if err = clientCache.GetClient("").Delete(ctx, obj); err != nil {
		return err
	}
	return nil
}

func Get(ctx context.Context) apibase.Result {
	res := &request.Resource{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] get read json error : %v", err)
	}
	obj, err := convertToObject(res)
	if err != nil {
		return err
	}
	exist, err := clientCache.GetClient("").Get(ctx, obj)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}
	return obj
}

func Status(ctx context.Context) apibase.Result {
	res := &request.Resource{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] get read json error : %v", err)
	}
	obj, err := convertToObject(res)
	if err != nil {
		return err
	}
	if err = clientCache.GetClient("").Status(ctx, obj); err != nil {
		return err
	}
	return nil
}

func List(ctx context.Context) apibase.Result {
	res := &request.ResourceList{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] list read json error : %v", err)
	}
	gvk := schema.GroupVersionKind{
		Group:   res.Group,
		Version: res.Version,
		Kind:    res.Kind,
	}
	obj, err := base.Schema.New(gvk)
	if err != nil {
		return fmt.Errorf("[resource] can't get type by %s/%s,%s", res.Group, res.Version, res.Kind)
	}
	if !meta.IsListType(obj) {
		return fmt.Errorf("[resource] object %s/%s,%s is not a list type", res.Group, res.Version, res.Kind)
	}
	opt := sigclient.InNamespace(res.Namespace)
	if err = clientCache.GetClient("").Lists(ctx, obj, opt); err != nil {
		return err
	}
	return obj
}

func DryRun(ctx context.Context) apibase.Result {
	res := &request.Resource{}
	err := ctx.ReadJSON(res)
	if err != nil {
		return fmt.Errorf("[resource] dry run read json error : %v", err)
	}
	obj, err := convertToObject(res)
	if err != nil {
		return nil
	}
	return clientCache.GetClient("").DryRun(res.Action, obj)
}

func convertToObject(res *request.Resource) (runtime.Object, error) {
	gvk := schema.GroupVersionKind{
		Group:   res.Group,
		Version: res.Version,
		Kind:    res.Kind,
	}
	obj, err := base.Schema.New(gvk)
	if err != nil {
		return nil, fmt.Errorf("[resource] can't get type by %s/%s,%s", res.Group, res.Version, res.Kind)
	}
	if err = json.Unmarshal(res.Data, obj); err != nil {
		return nil, fmt.Errorf("[resource] unmarshal data to obj %T error : %v", obj, err)
	}
	return obj, nil
}

func Events(ctx context.Context) apibase.Result {
	return monitor.GetEvents()
}
