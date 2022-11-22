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

package api

import (
	"dtstack.com/dtstack/easymatrix/addons/oldkube/pkg/api/impl"
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
)

var ClientGoAPIRoutes = apibase.Route{
	Path: "clientgo",
	SubRoutes: []apibase.Route{{
		//http:://xxxx/api/v1/clientgo/allocated
		Path: "allocated",
		GET:  impl.GetAllocated,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取集群mem,cpu,pod等资源总使用情况[EasyKube API]",
			},
		},
	}, {
		//http:://xxxx/api/v1/clientgo/top5
		Path: "top5",
		GET:  impl.GetTop5,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取cpu、mem资源使用率前五高的机器[EasyKube API]",
			},
		},
	}, {
		//http:://xxxx/api/v1/clientgo/workload
		Path: "workload",
		GET:  impl.GetWorkLoad,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取资源工作负载[EasyKube API]",
			},
		},
	}, {
		//http:://xxxx/api/v1/clientgo/allocatedPodList
		Path: "allocatedPodList",
		GET:  impl.GetAllocatedPodList,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取集群每个节点的pod资源总使用情况[EasyKube API]",
			},
		},
	}, {
		//http:://xxxx/api/v1/clientgo/componentStatus
		Path: "componentStatus",
		GET:  impl.GetComponentStatus,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取集群的组件状态[EasyKube API]",
			},
		},
	}, {
		//http:://xxxx/api/v1/clientgo/extraInfo
		Path: "extraInfo",
		GET:  impl.GetExtraInfo,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取节点的role信息，k8s版本信息[EasyKube API]",
			},
		},
	}, {
		Path: "dynamic",
		SubRoutes: []apibase.Route{{
			//http:://xxxx/api/v1/clientgo/dynamic/apply
			Path: "apply",
			POST: impl.ApplyDynamicResource,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "在集群内部动态创建更新k8s资源[EasyKube API]",
				},
			},
		}, {
			//http:://xxxx/api/v1/clientgo/dynamic/delete
			Path: "delete",
			POST: impl.DeleteDynamicResource,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "在集群内部删除k8s资源[EasyKube API]",
				},
			},
		}, {
			//http:://xxxx/api/v1/clientgo/dynamic/get
			Path: "get",
			POST: impl.GetDynamicResource,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "在集群内部获取k8s资源[EasyKube API]",
				},
			},
		}},
	}, {
		Path: "namespace",
		SubRoutes: []apibase.Route{{
			//http:://xxxx/api/v1/clientgo/namespace/list
			Path: "list",
			GET:  impl.GetNamespaceList,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "获取集群namespace list[EasyKube API]",
				},
			},
		}, {
			//http:://xxxx/api/v1/clientgo/namespace/create
			Path: "create",
			POST: impl.CreateNamespace,
			Docs: apibase.Docs{
				GET: &apibase.ApiDoc{
					Name: "创建集群namespace[EasyKube API]",
				},
			},
		}},
	}},
}
