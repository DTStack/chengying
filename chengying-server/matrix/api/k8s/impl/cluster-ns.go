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
	"fmt"
	"github.com/kataras/iris/context"
	"strconv"
)

func NamespaceSave(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	requestInfo := &view.NamespaceSaveReq{}
	err := ctx.ReadJSON(requestInfo)
	if err != nil{
		return fmt.Errorf("read json to NamespaceSaveReq error:%v",err)
	}

	clusterid := ctx.GetCookie(view.ClusterId)
	userName := ctx.GetCookie(view.User)
	return resource.Save(ctx,clusterid,userName,requestInfo)
}

func NamespaceStatus(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	ns := ctx.Params().Get("namespace")
	status := ctx.URLParam("status")
	desc := ctx.URLParam("desc")
	typ := ctx.URLParam("type")
	clusterid := ctx.GetCookie(view.ClusterId)
	resp,err := resource.GetNamespaceStatus(ctx,ns,clusterid,status,desc,typ)
	if err != nil{
		return err
	}
	return resp
}

func NameSpaceListStatus(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	clusterid := ctx.GetCookie(view.ClusterId)
	status := ctx.URLParam("status")
	desc := ctx.URLParam("desc")
	typ := ctx.URLParam("type")
	resp,err := resource.GetNamespaceListStatus(ctx,clusterid,status,desc,typ)
	if err != nil{
		return err
	}
	return resp
}

func AgentGenerate(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	req := &view.AgentGenerateReq{}
	if err := ctx.ReadJSON(req); err != nil{
		return fmt.Errorf("read json to AgentGenerateReq error :%v",err)
	}
	clusterid := ctx.GetCookie(view.ClusterId)
	var cid int
	if clusterid == ""{
		return fmt.Errorf("AgentGenerate Get cluster id error from cookie")
	}
	cid, _= strconv.Atoi(clusterid)
	rsp,err := resource.AgentGenerate(req,cid)
	if err != nil{
		return err
	}
	return rsp
}

func NamespaceGet(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	ns := ctx.Params().Get("namespace")
	clusterid := ctx.GetCookie(view.ClusterId)
	resp,err := resource.NamespaceGet(ns,clusterid)
	if err != nil{
		return err
	}
	return resp
}

func NamespaceDeleteConfirm(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	ns := ctx.Params().Get("namespace")
	clusterid := ctx.GetCookie(view.ClusterId)
	if err := resource.NamespaceDeleteConfirm(ns,clusterid);err != nil{
		return &view.NamespaceDeleteConfirmRsp{Status: false}
	}
	return &view.NamespaceDeleteConfirmRsp{Status: true}
}

func NamespaceDelete(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	ns := ctx.Params().Get("namespace")
	clusterid := ctx.GetCookie(view.ClusterId)
	return resource.NamespaceDelete(ctx,ns,clusterid)
}

func NamespacePing(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	clusterid := ctx.GetCookie(view.ClusterId)
	req := &view.NamespacePingReq{}
	if err := ctx.ReadJSON(req); err != nil{
		return fmt.Errorf("read json to NamespacepingReq error :%v",err)
	}
	return resource.NamespacePing(ctx,clusterid,req)
}

func NamespaceEvent(ctx context.Context) apibase.Result{
	log.Debugf("[cluster-ns]: %s",ctx.Request().RequestURI)
	clusterid := ctx.GetCookie(view.ClusterId)
	namespace := ctx.Params().Get("namespace")
	limit := ctx.URLParam("limit")
	start := ctx.URLParam("start")
	si,err := strconv.Atoi(start)
	if err != nil{
		return fmt.Errorf("the start %s is not valid",start)
	}
	li,err := strconv.Atoi(limit)
	if err != nil{
		return fmt.Errorf("the limit %s is not valid",limit)
	}
	result,err := resource.NamsapceEvent(clusterid,namespace,li,si)
	if err != nil{
		return err
	}
	return result

}
