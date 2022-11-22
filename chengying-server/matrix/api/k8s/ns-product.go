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
)

var NSProductAPIRoutes = apibase.Route{
	Path: "product/manage",
	SubRoutes: []apibase.Route{
		{
			Path: 	"{namespace:string}",
			GET: 	impl.GetParentProducts,
			Docs: 	apibase.Docs{
				Desc: 	"get deployed parent product name in namespace",
			},
			SubRoutes: 	[]apibase.Route{
				{
					Path:	"{parent_product_name:string}",
					GET: 	impl.GetProducts,
					Docs: 	apibase.Docs{
						Desc: 	"get deployed product named in parent product",
					},
					SubRoutes: []apibase.Route{
						{
							Path: 		"{product_name:string}",
							GET: 		impl.GetServiceList,
							SubRoutes: 	[]apibase.Route{
								{
									Path: 		"{service_name:string}",
									GET: 		impl.GetService,
								},
							},
						},
					},
				},
			},
		},
	},
}
