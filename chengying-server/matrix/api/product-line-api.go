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
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/api/impl"
)

var ProductLineOperationEasyMatrixAPIRoutes = apibase.Route{
	//http://xxxx/api/v2/product_line
	Path: "product_line",
	GET:  impl.ProductLineInfo,
	Docs: apibase.Docs{
		GET: &apibase.ApiDoc{
			Name: "获取所有产品线信息接口[EasyMatrix API]",
		},
	},
	SubRoutes: []apibase.Route{{
		//http://xxxx/api/v2/product_line/upload
		Path: "upload",
		POST: impl.UploadProductLine,
		Docs: apibase.Docs{
			POST: &apibase.ApiDoc{
				Name: "上传产品线接口[EasyMatrix API]",
			},
		},
	}, {
		//http://xxxx/api/v2/product_line/{id:int(primary key)}
		Path:   "{id:int}",
		DELETE: impl.DeleteProductLine,
		Docs: apibase.Docs{
			DELETE: &apibase.ApiDoc{
				Name: "删除产品线接口[EasyMatrix API]",
			},
		},
	}, {
		//http://xxxx/api/v2/product_line/product_list
		Path: "product_list",
		GET:  impl.ProductListOfProductLine,
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "产品包列表接口[EasyMatrix API]",
			},
		},
	}},
}
