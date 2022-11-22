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
	molemodel "dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/model"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workload/support"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/util"
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/deployment"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/mole"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/resourcequota"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/workload"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"dtstack.com/dtstack/easymatrix/matrix/model/kube/union"
	"dtstack.com/dtstack/easymatrix/schema"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"
	"strings"
)

func InstanceReplica(ctx context.Context, clusterid string, vo *view.InstanceReplicaReq) error {

	cid, err := strconv.Atoi(clusterid)
	if err != nil {
		log.Errorf("[instance_replica]: convet %s to int error %v", clusterid, err)
		return err
	}
	tbscs, err := modelkube.DeployNamespaceList.GetLike(vo.Namespace, cid, constant.NAMESPACE_VALID, true)
	if err != nil {
		return err
	}
	if tbscs == nil {
		return fmt.Errorf("no valid namespaced client exist")
	}
	// now one cluster only can use one namspaceclient, until the em front modified
	tbsc := tbscs[0]
	namespace := tbsc.Namespace
	cache, err := kube.ClusterNsClientCache.GetClusterNsClient(clusterid).GetClientCache(kube.ImportType(tbsc.Type))
	if err != nil {
		return err
	}

	unionTbsc, err := union.UnionT4T7.GetProduct(cid, namespace, vo.ProductName)
	if err != nil {
		return err
	}
	if unionTbsc == nil {
		return fmt.Errorf("[instance_replica]: can not found the deployed product %s in cluster %d, in namespace %s", vo.ProductName, cid, namespace)
	}
	sc, err := schema.Unmarshal(unionTbsc.Product)
	if err != nil {
		log.Errorf("[instance_replica]: unmashal product fail, err: %v", err)
		return err
	}
	client := cache.GetClient(namespace)
	quota, err := resourcequota.Get(ctx, client, namespace)
	if err != nil {
		return err
	}
	svc := sc.Service[vo.ServiceName]
	var cpu_limit int64
	var mem_limit int64
	if len(svc.Workload) != 0 {
		if svc.Instance != nil && svc.Instance.NeedStorage != nil && *svc.Instance.NeedStorage {
			return fmt.Errorf("the service %s need storage, please check the storage is ready and redeploy", vo.ServiceName)
		}
		wl := workload.New()
		wl.Name = util.BuildWorkloadName(vo.ProductName, vo.ServiceName)
		wl.Namespace = namespace
		exist, err := cache.GetClient(namespace).Get(ctx, wl)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("the service %s in product %s in cluster %s not support replica", vo.ServiceName, vo.ProductName, clusterid)
		}
		// now, the service corresponds to one type instance, multiple part definitions are not supported
		if quota != nil {
			part := wl.Spec.WorkLoadParts[0]
			for _, step := range part.Steps {
				if step.Type == support.BoundTypeContainer || step.Type == support.BoundTypeInitContainer {
					limits, err := workload.FieldGet(step.Object.Object, "resources.limits")
					if err != nil {
						return err
					}
					m := limits.(map[string]interface{})
					quota, err := resource.ParseQuantity(m["memory"].(string))
					if err != nil {
						return err
					}
					mem_limit = mem_limit + quota.Value()

					quota, err = resource.ParseQuantity(m["cpu"].(string))
					if err != nil {
						return err
					}
					cpu_limit = cpu_limit + quota.Value()
				}
			}
			if err := judgeLimit(quota, cpu_limit, mem_limit, int64(vo.Replica)); err != nil {
				return err
			}
		}
		workload.FieldSet(wl.Spec.WorkLoadParts[0].BaseWorkLoad.Parameters.Object, "spec.replicas", vo.Replica)
		return cache.GetClient(namespace).Update(ctx, wl)

	} else {
		m := mole.New()
		m.Namespace = namespace
		m.Name = strings.ToLower(vo.ProductName)
		exist, err := client.Get(ctx, m)
		if err != nil {
			return err
		}
		if !exist {
			fmt.Errorf("not found the service %s in product %s in cluster %s", vo.ServiceName, vo.ProductName, clusterid)
		}
		if quota != nil {
			limits := m.Spec.Product.Service[vo.ServiceName].Instance.Resources.Limits
			resource_limit_cpu := limits["cpu"]
			resource_limit_mem := limits["memory"]
			//if there is no limit on the mole, use the default limit set by mole operator
			if limits == nil {
				name := molemodel.BuildResourceName(molemodel.MoleDeploymentName, m.Spec.Product.ParentProductName,
					m.Spec.Product.ProductName, vo.ServiceName)
				deploy, err := deployment.Get(ctx, client, namespace, name)
				if err != nil {
					return err
				}
				if deploy == nil {
					return fmt.Errorf("the service %s is not deployed", vo.ServiceName)
				}
				cs := deploy.Spec.Template.Spec.Containers
				for _, c := range cs {
					if c.Name == strings.ToLower(vo.ServiceName) {
						resource_limit_cpu = c.Resources.Limits["cpu"]
						resource_limit_mem = c.Resources.Limits["memory"]
					}
				}
			}
			cpu_limit = resource_limit_cpu.Value()
			mem_limit = resource_limit_mem.Value()
			if err := judgeLimit(quota, resource_limit_cpu.Value(), resource_limit_mem.Value(), int64(vo.Replica)); err != nil {
				return err
			}
		}
		m.Spec.Product.Service[vo.ServiceName].Instance.Deployment.Replicas = int32(vo.Replica)
		return cache.GetClient(namespace).Update(ctx, m)
	}
}

func judgeLimit(quota *corev1.ResourceQuota, limtcpu, limitmem, replica int64) error {
	r := quota.Status.Hard["limits.cpu"]
	cpuTotal := r.Value()

	r = quota.Status.Hard["limits.memory"]
	memTotal := r.Value()

	r = quota.Status.Used["limits.cpu"]
	cpuUsed := r.Value()

	r = quota.Status.Used["limits.memory"]
	memUsed := r.Value()

	if cpuUsed+replica*limtcpu > cpuTotal {
		return fmt.Errorf("The cpu resources are insufficient, please check")
	}

	if memUsed+replica*limitmem > memTotal {
		return fmt.Errorf("The memory resources are insufficient, please check")
	}
	return nil
}
