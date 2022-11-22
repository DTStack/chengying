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
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/api/impl"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
)

var ResourceApi = apibase.Route{
	Path: "kube/resource",
	SubRoutes: []apibase.Route{{
		//http:://xxxx/api/v1/kube/resource/get
		Path: "get",
		POST: impl.Get,
		Docs: apibase.Docs{
			Desc: "get k8s resource",
		},
	}, {
		//http:://xxxx/api/v1/kube/resource/delete
		Path: "delete",
		POST: impl.Delete,
		Docs: apibase.Docs{
			Desc: "delete k8s resource",
		},
	}, {
		//http:://xxxx/api/v1/kube/resource/apply
		Path: "apply",
		POST: impl.Apply,
		Docs: apibase.Docs{
			Desc: "apply k8s resource",
		},
	}, {
		//http:://xxxx/api/v1/kube/resource/list
		Path: "list",
		POST: impl.List,
		Docs: apibase.Docs{
			Desc: "list k8s resource",
		},
	}, {
		//http:://xxxx/api/v1/kube/resource/create
		Path: "create",
		POST: impl.Create,
		Docs: apibase.Docs{
			Desc: "create k8s resource",
		},
	}, {
		Path: "dryrun",
		POST: impl.DryRun,
		Docs: apibase.Docs{
			Desc: "dry run",
		},
	}, {
		Path: "update",
		POST: impl.Update,
		Docs: apibase.Docs{
			Desc: "update resource",
		},
	}, {
		Path: "events",
		GET:  impl.Events,
		Docs: apibase.Docs{
			Desc: "get listwathc events",
		},
	}, {
		Path: "status",
		POST: impl.Status,
		Docs: apibase.Docs{
			Desc: "update status",
		},
	}},
}
