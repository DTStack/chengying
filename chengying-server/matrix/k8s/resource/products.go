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
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/deployment"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/endpoints"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/job"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/workload"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"dtstack.com/dtstack/easymatrix/matrix/model/kube/union"
	"dtstack.com/dtstack/easymatrix/schema"
	batchv1 "k8s.io/api/batch/v1"
	"strconv"
	"strings"
)

func GetParentProductList(namespace,clusterid string) (*view.NamespacedProductsRsp,error){
	products,err := getAll(namespace,clusterid)
	if err != nil{
		return nil, err
	}
	for k,_ := range products{
		products[k] = nil
	}
	return &view.NamespacedProductsRsp{
		Size:     len(products),
		Products: products,
	},nil
}
//get product name below parentproduct
func GetProductList(namespace,clusterid,parentproductName string) (*view.NamespacedProductsRsp,error){
	products,err := getByParentProduct(namespace,clusterid,parentproductName)
	if err != nil{
		return nil, err
	}
	for _,pp := range products{
		for pn,_ := range pp{
			pp[pn] = nil
		}
	}
	return &view.NamespacedProductsRsp{
		Size: len(products),
		Products: products,
	},nil
}

func GetServiceList(namespace,clusterid,parentproductName,productName string) (*view.NamespacedProductsRsp,error){
	products,err := getByProduct(namespace,clusterid,parentproductName,productName)
	if err != nil{
		return nil, err
	}
	for _,pp := range products{
		for _,p := range pp{
			for sn,_ := range p{
				p[sn] = nil
			}
		}
	}
	return &view.NamespacedProductsRsp{
		Size: len(products),
		Products: products,
	},nil
}

func GetService(ctx context.Context,namespace,clusterid, parentproductName, productName, serviceName string) (*view.NamespacedProductsRsp,error){
	products,err := getByService(namespace,clusterid, parentproductName, productName, serviceName)
	if err != nil{
		return nil,err
	}
	cid,err := strconv.Atoi(clusterid)
	if err != nil{
		log.Errorf("[resource products]: convert clusterid %s to int error %v",clusterid,err)
	}
	tbsc,err := modelkube.DeployNamespaceList.Get(namespace,cid)
	if err != nil{
		return nil,err
	}
	if tbsc == nil{
		return nil,nil
	}
	cache,err := kube.ClusterNsClientCache.GetClusterNsClient(clusterid).GetClientCache(kube.ImportType(tbsc.Type))
	if err != nil{
		return nil,err
	}
	setServiceStatus(ctx,cache,namespace,&products)
	return &view.NamespacedProductsRsp{
		Size: len(products),
		Products: products,
	},nil
}

func setServiceStatus(ctx context.Context,cache kube.ClientCache,namespace string, products *view.Products ) error{
	for ppn,pp := range *products{
		for pn,p := range pp{
			for sn,_ := range p{
				svc := p[sn]
				// workload type installl
				if len(svc.WorkloadType) != 0{
					wl := workload.New()
					wl.Name = util.BuildWorkloadName(pn.String(),sn.String())
					wl.Namespace = namespace
					_,err := cache.GetClient(namespace).Get(ctx,wl)
					if err != nil{
						return err
					}
					endPointsAdress := map[string]struct{}{}
					replica := 0
					for _, part := range wl.Spec.WorkLoadParts{
						typ := part.BaseWorkLoad.Type
						if typ == support.WorkloadTypeDeployment || typ == support.WorkloadTypeStatefulset{
							partReplica,err := workload.FieldGet(part.BaseWorkLoad.Parameters.Object,"spec.replicas")
							if err != nil{
								return err
							}
							replica = replica + int(partReplica.(float64))
							for _,step := range part.Steps{
								if step.Type != support.CreateTypeService{
									continue
								}
								svcName := util.BuildStepName(util.BuildBaseName(wl.Name,part.BaseWorkLoad.Name),step.Name)
								endpoints := endpoints.New()
								endpoints.Namespace = namespace
								endpoints.Name = svcName
								exist,err := cache.GetClient(namespace).Get(ctx,endpoints)
								if err != nil{
									return err
								}
								if !exist{
									continue
								}
								for _, adress := range endpoints.Subsets[0].Addresses{
									endPointsAdress[adress.TargetRef.Name] = struct{}{}
								}
							}
						}
					}
					ensureHealth(len(endPointsAdress),replica,svc)
					continue
				}
				// the job type of mole install
				if svc.IsJob{
					jobName := molemodel.BuildResourceName(molemodel.MoleJobName,ppn.String(),pn.String(),sn.String())
					job := job.New()
					job.Name = jobName
					job.Namespace = namespace
					exist,err := cache.GetClient(namespace).Get(ctx,job)
					if err != nil{
						return err
					}
					if !exist{
						job = nil
					}
					ensureHealthByJob(job,svc)
					continue
				}
				// the deployment of mole install
				svcName := molemodel.BuildResourceName(molemodel.MoleServiceName,ppn.String(),pn.String(),sn.String())
				endpoints := endpoints.New()
				endpoints.Namespace = namespace
				endpoints.Name = svcName
				exist,err := cache.GetClient(namespace).Get(ctx,endpoints)
				if err != nil{
					return err
				}
				succeed := 0
				if exist{
					if len(endpoints.Subsets) >0 && len(endpoints.Subsets[0].Addresses) > 0{
						succeed = len(endpoints.Subsets[0].Addresses)
					}
				}

				deployName := molemodel.BuildResourceName(molemodel.MoleDeploymentName,ppn.String(),pn.String(),sn.String())
				deploy := deployment.New()
				deploy.Namespace = namespace
				deploy.Name = deployName
				exist,err = cache.GetClient(namespace).Get(ctx,deploy)
				if err != nil{
					return err
				}
				replica := 0
				if exist{
					replica = int(*deploy.Spec.Replicas)
				}
				ensureHealth(succeed,replica,svc)
			}
		}
	}
	return nil
}

func ensureHealth(succeed int, desired int, svc *view.Service){
	healthState := "unhealthy"
	healthStateCount := 0
	serviceStatus := "unavailable"
	serviecStatusCount := 0
	if succeed >0 {
		serviceStatus = "available"
		serviecStatusCount = succeed
	}
	if serviceStatus == "available"{
		if desired == serviecStatusCount{
			healthState = "healthy"
			healthStateCount = serviecStatusCount
		}else{
			healthState = "unhealthy"
			healthStateCount = desired - serviecStatusCount
		}
	}else{
		healthState = "unhealthy"
		healthStateCount = desired
		serviceStatus = "unavailable"
		serviecStatusCount = desired
	}
	svc.HealthState = healthState
	svc.HealthStateCount = healthStateCount
	svc.ServiceStatus = serviceStatus
	svc.ServiecStatusCount = serviecStatusCount
}

func ensureHealthByJob(job *batchv1.Job, svc *view.Service){

	svc.ServiceStatus = "available"
	svc.ServiecStatusCount = 0
	svc.HealthState = "healthy"
	svc.HealthStateCount = 0

	if job == nil{
		svc.ServiceStatus = "unavailable"
		svc.HealthState = "unhealthy"
	}
	desired := job.Spec.Completions
	succeed := job.Status.Succeeded
	if desired != nil{
		if *desired == succeed{
			svc.ServiceStatus = "available"
			svc.ServiecStatusCount = int(succeed)
			svc.HealthState = "healthy"
			svc.HealthStateCount = int(succeed)
		}else{
			svc.ServiceStatus = "unavailable"
			svc.ServiecStatusCount =int(*desired)
			svc.HealthState = "unhealthy"
			svc.HealthStateCount = int(*desired - succeed)

		}
	}else{
		if succeed >0 {
			svc.ServiceStatus = "available"
			svc.ServiecStatusCount = int(succeed)
			svc.HealthState = "healthy"
			svc.HealthStateCount = int(succeed)
		}
	}
}

//func ensureHealthByEndpoints(endpoints *corev1.Endpoints, replica int, svc *view.Service){
//
//	healthState := "unhealthy"
//	healthStateCount := 0
//	serviceStatus := "unavailable"
//	serviecStatusCount := 0
//	if endpoints != nil{
//		if len(endpoints.Subsets) >0 && len(endpoints.Subsets[0].Addresses) > 0{
//			serviceStatus = "available"
//			serviecStatusCount = len(endpoints.Subsets[0].Addresses)
//		}
//	}
//	if serviceStatus == "available"{
//		if replica == serviecStatusCount{
//			healthState = "healthy"
//			healthStateCount = serviecStatusCount
//		}else{
//			healthState = "unhealthy"
//			healthStateCount = replica - serviecStatusCount
//		}
//	}else{
//		healthState = "unhealthy"
//		healthStateCount = replica
//		serviceStatus = "unavailable"
//		serviecStatusCount = replica
//	}
//	svc.HealthState = healthState
//	svc.HealthStateCount = healthStateCount
//	svc.ServiceStatus = serviceStatus
//	svc.ServiecStatusCount = serviecStatusCount
//}

func getByService(namespace,clusterid, parentproductName, productName, serviceName string) (view.Products,error){
	products,err := getByProduct(namespace,clusterid,parentproductName,productName)
	if err != nil{
		return nil, err
	}
	if serviceName == "all"{
		return products,nil
	}
	result := view.Products{}
	for ppn,pp := range products{
		parentproduct := view.ParentProduct{}
		result[ppn] = parentproduct
		for pn,p := range pp{
			prouct := view.Product{}
			parentproduct[pn] = prouct
			for sn,s := range p{
				if strings.Contains(strings.ToLower(string(sn)),strings.ToLower(serviceName)){
					prouct[sn] = s
				}
			}
		}
	}
	return result,nil
}

func getByProduct(namespace,clusterid, parentproductName, productName string) (view.Products,error){
	products,err := getByParentProduct(namespace,clusterid,parentproductName)
	if err != nil{
		return nil, err
	}
	if productName == "all"{
		return products,nil
	}
	reslut := view.Products{}
	for ppn,pp := range products{
		parentproduct := view.ParentProduct{}
		reslut[ppn] = parentproduct
		for pn,p := range pp{
			if string(pn) == productName{
				parentproduct[pn] = p
			}
		}
	}
	return reslut,nil
}

func getByParentProduct(namespace,clusterid, parentproductName string) (view.Products,error){
	products,err := getAll(namespace,clusterid)
	if err != nil{
		return nil, err
	}
	if parentproductName == "all"{
		return products,nil
	}
	parentProduct := products[view.ParentProductName(parentproductName)]
	result := view.Products{}
	result[view.ParentProductName(parentproductName)] = parentProduct
	return result,nil
}

func getAll(namespace,clusterid string) (view.Products,error){
	cid,err := strconv.Atoi(clusterid)
	if err != nil{
		log.Errorf("[products]: convert clusterid %s to int error %v",clusterid,err)
		return nil,err
	}
	tbscs ,err := union.UnionT4T7.SelectParentProduct(cid,namespace)
	if err != nil{
		return nil ,err
	}
	products := view.Products{}
	for _,tbsc := range tbscs{
		config,err := schema.Unmarshal(tbsc.Product)
		if err != nil{
			return nil,err
		}
		parent,exist := products[view.ParentProductName(config.ParentProductName)];
		if !exist{
			parent = view.ParentProduct{}
			products[view.ParentProductName(tbsc.ParentProductName)] = parent
		}
		product := view.Product{}
		for k,v := range config.Service{
			if v.Instance == nil{
				continue
			}
			service := view.Service{
				Version:            v.Version,
				Group:              v.Group,
				IsJob: 				v.IsJob,
			}
			if len(v.Workload) != 0{
				service.WorkloadType = v.Workload
			}
			product[view.ServiceName(k)] = &service
		}
		parent[view.ProudctName(config.ProductName)] = product
	}
	return products,nil
}
