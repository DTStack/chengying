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

package resource

import (
	"context"
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/configmap"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/deployment"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/mole"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/resourcequota"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/secret"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/service"
	"strings"
)

func NamespacePing(ctx context.Context,clusterid string,vo *view.NamespacePingReq) error{
	cache,err := kube.ClusterNsClientCache.GetClusterNsClient(clusterid).GetClientCache(kube.IMPORT_AGENT)
	if err != nil{
		return err
	}
	ip := vo.Ip
	port := vo.Port
	if !strings.HasPrefix(ip,"http"){
		ip = "http://"+ip
	}
	host := ip + ":" +port
	ns := vo.Namespace
	err = cache.Connect(host,ns)
	if err != nil{
		return err
	}
	client := cache.GetClient(ns)
	return ping(client,ns)
}

func ping(client kube.Client, ns string) error{
	var err error
	if err = service.Ping(client,ns); err != nil{
		return err
	}
	if err = secret.Ping(client,ns); err != nil{
		return err
	}
	if err = resourcequota.Ping(client,ns); err != nil{
		return err
	}
	if err = mole.Ping(client,ns); err != nil{
		return err
	}
	if err = deployment.Ping(client,ns); err != nil{
		return err
	}
	if err = configmap.Ping(client,ns); err != nil{
		return err
	}
	return nil
}
