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
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/resourcequota"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	corev1 "k8s.io/api/core/v1"
	"math"
	"strconv"
)

func GetNamespaceListStatus (ctx context.Context, clusterid,status,desc,typ string)([]view.NamespaceStatusRsp,error){

	return GetNamespaceStatus(ctx,"",clusterid,status,desc,typ)
}

//the namespace is support like. so get the namespaced client by namesapce, can't use the namespace input
func GetNamespaceStatus(ctx context.Context, namespace, clusterId, status, descStr, typ string)([]view.NamespaceStatusRsp,error){
	tbscs,err := modelkube.DeployNamespaceList.Select(clusterId,status,descStr,namespace,typ)
	if err != nil{
		return nil, err
	}
	if tbscs == nil{
		return nil, nil
	}
	return getStatus(ctx,clusterId,tbscs)
}

func getStatus(ctx context.Context, clusterId string, namespacelist []modelkube.DeployNamespaceListSchema) ([]view.NamespaceStatusRsp,error){
	resp := make([]view.NamespaceStatusRsp,0,len(namespacelist))
	for _,tbsc := range namespacelist{
		var quota *corev1.ResourceQuota
		if tbsc.Status == constant.NAMESPACE_VALID{
			clientCache,err := kube.ClusterNsClientCache.GetClusterNsClient(clusterId).GetClientCache(kube.ImportType(tbsc.Type))
			if err != nil{
				return nil, err
			}
			quota,err = resourcequota.Get(ctx,clientCache.GetClient(tbsc.Namespace),tbsc.Namespace)
			if err != nil{
				return nil,err
			}
		}

		cpuUsed  := "--"
		cpuTotal := "--"
		memUsed  := "--"
		memTotal := "--"
		var cpupercent float64
		var mempercent float64

		if quota != nil{

			r := quota.Status.Hard["limits.cpu"]
			//this value's unit is m
			cpuTotalfloat := float64(r.MilliValue())
			cpuTotal = formatFloat(cpuTotalfloat/1000,2)+"core"

			r = quota.Status.Hard["limits.memory"]
			//this value's unit is byte
			memTotalfloat := float64(r.Value())
			memTotal = formatFloat(memTotalfloat/1024/1024/1024,2)+"GB"

			r = quota.Status.Used["limits.cpu"]
			cpuUsedfloadt := float64(r.MilliValue())
			cpuUsed = formatFloat(cpuUsedfloadt/1000,2)+"core"

			r = quota.Status.Used["limits.memory"]
			memUsedfloat := float64(r.Value())
			memUsed = formatFloat(memUsedfloat/1024/1024/1024,2)+"GB"

			cpupercent = (cpuUsedfloadt*100)/cpuTotalfloat
			mempercent = (memUsedfloat*100)/memTotalfloat

		}
		result := view.NamespaceStatusRsp{
			Id:         tbsc.Id,
			Namespace:  tbsc.Namespace,
			Status:     tbsc.Status,
			CpuUsed:    cpuUsed,
			CpuTotal:   cpuTotal,
			CpuPercent: cpupercent,
			MemUsed:    memUsed,
			MemTotal:   memTotal,
			MemPercent: mempercent,
			User:       tbsc.User,
			UpdateTime: tbsc.UpdateTime.Format(constant.DATE_FORMAT),
			Type:       tbsc.Type,

		}
		resp = append(resp,result)
	}
	return resp,nil
}

func formatFloat(num float64, decimal int) string{
	d := float64(1)
	if decimal > 0 {
		d = math.Pow10(decimal)
	}
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}
