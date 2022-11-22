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

package k8s

import (
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/impl"
	"github.com/kataras/iris/context"
)

var ClusterResourceAPIRoutes = apibase.Route{
	Path:        "cluster/manage",
	SubRoutes:   []apibase.Route{
		{
			Path:        "namespaces",
			GET:         impl.NameSpaceListStatus,
			Docs:        apibase.Docs{
				Desc:   	"query namespaces situation",
			},
			SubRoutes:   []apibase.Route{
				{
					Path: "{namespace:string}",
					GET: impl.NamespaceStatus,
					Docs: apibase.Docs{
						Desc: "query namespace situation",
					},
				},
			},
		},
		{
			Path:        "namespace",
			SubRoutes:   []apibase.Route{
				{
					Path: 		"save",
					POST: 		impl.NamespaceSave,
					Middlewares: []context.Handler{
						apibase.CheckPermission3,
					},
					Docs: 		apibase.Docs{
						Desc: 		"save data about the client in different namespace that represents different permissions",
					},
				},
				{
					Path: 		"agent/generate",
					POST:       impl.AgentGenerate,
					Docs:       apibase.Docs{
						Desc:       "generate agent mode import yaml",
					},
				},
				{
					Path: 		"ping",
					POST: 		impl.NamespacePing,
					Docs: 		apibase.Docs{
						Desc: 		"test the agent mode if connect",
					},
				},
			},
		},
		{
			Path:		"{namespace:string}",
			SubRoutes: 	[]apibase.Route{
				{
					Path: 		"get",
					GET: 		impl.NamespaceGet,
					Docs: 		apibase.Docs{
						Desc: 		"get namespace client info",
					},
				},
				{
					Path: 		"events",
					GET: 		impl.NamespaceEvent,
					Docs: 		apibase.Docs{
						Desc: 		"get namespace event",
					},
				},
				{
					Path: 		"delete",
					POST: 		impl.NamespaceDelete,
					Middlewares: []context.Handler{
						apibase.CheckPermission3,
					},
					Docs: 		apibase.Docs{
						Desc: 		"delete namespace client info",
					},
					SubRoutes: 	[]apibase.Route{
						{
							Path: 	"confirm",
							GET: 	impl.NamespaceDeleteConfirm,
							Middlewares: []context.Handler{
								apibase.CheckPermission3,
							},
							Docs: 	apibase.Docs{
								Desc: "Confirm whether it can be deleted",
							},
						},
					},
				},
			},
		},
	},
}
