/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package impl

import (
	apibase "easyagent/go-common/api-base"
	"easyagent/internal/server/model"

	"fmt"
	"strconv"

	"github.com/kataras/iris/context"
)

func QueryDeploymentProducts(ctx context.Context) apibase.Result {
	errs := apibase.NewApiParameterErrors()
	pagination := apibase.GetPaginationFromQueryParameters(errs, ctx)

	withType := ctx.URLParam("type")
	var (
		prods []model.ProductListInfo
		count int
	)
	switch withType {
	case "sidecar", "easyagent", strconv.Itoa(model.PROD_SIDECAR):
		prods, count = model.ProductList.QueryProductList(model.PROD_SIDECAR, pagination)
	case "easydb-agent", "easydb", strconv.Itoa(model.PROD_EASYDB_AGENT):
		prods, count = model.ProductList.QueryProductList(model.PROD_EASYDB_AGENT, pagination)
	case "easylog-agent", "easylog", strconv.Itoa(model.PROD_EASYLOG_AGENT):
		prods, count = model.ProductList.QueryProductList(model.PROD_EASYLOG_AGENT, pagination)
	case "":
		prods, count = model.ProductList.QueryProductList(0, pagination)
	default:
		errs.AppendError("type", "Illegal query type name: '%s'", withType)
	}
	errs.CheckAndThrowApiParameterErrors()

	products := []map[string]interface{}{}
	for _, p := range prods {
		products = append(products, map[string]interface{}{
			"id":      p.ID,
			"type":    model.ProdTypeString(p.Type),
			"name":    p.Name,
			"version": p.Version,
		})
	}
	return map[string]interface{}{
		"products": products,
		"total":    count,
	}
}

func GetDeploymentProductInfo(ctx context.Context) apibase.Result {
	prodId, err := ctx.Params().GetInt("prod_id")
	if err != nil {
		e := apibase.NewApiParameterErrors()
		e.AppendError("prod_id", err)
		e.CheckAndThrowApiParameterErrors()
	}

	p := model.ProductList.GetProductInfo(prodId)
	if p != nil {
		return map[string]interface{}{
			"id":                p.ID,
			"type":              p.Type,
			"name":              p.Name,
			"version":           p.Version,
			"description":       p.Description.String,
			"checksum":          p.CheckSum,
			"release_date":      p.ReleaseDate.Time,
			"deploy_times":      p.DeployTimes,
			"deploy_fail_times": p.DeployFailedTimes,
		}
	} else {
		return fmt.Errorf("No such product which id = %v", prodId)
	}
}

func NewDeploymentProductInfo(ctx context.Context) apibase.Result {

	return nil
}

func RemoveDeploymentProduct(ctx context.Context) apibase.Result {

	return nil
}

func UpdateDeploymentProductInfo(ctx context.Context) apibase.Result {

	return nil
}
