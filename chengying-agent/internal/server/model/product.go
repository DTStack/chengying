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

package model

import (
	"database/sql"
	apibase "easyagent/go-common/api-base"
	dbhelper "easyagent/go-common/db-helper"
)

type productListTable struct {
	dbhelper.DbTable
}

var ProductList = productListTable{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_PRODUCT_LIST},
}

const (
	PROD_SIDECAR = iota + 1
	PROD_EASYDB_AGENT
	PROD_EASYLOG_AGENT
)

func ProdTypeString(t int) string {
	pt := ""
	switch t {
	case PROD_SIDECAR:
		pt = "sidecar"
	case PROD_EASYDB_AGENT:
		pt = "easydb-agent"
	case PROD_EASYLOG_AGENT:
		pt = "easylog-agent"
	default:
		pt = "unknown"
	}
	return pt
}

type ProductListInfo struct {
	ID          int            `db:"id"`
	Type        int            `db:"type"`
	Name        string         `db:"name"`
	Version     string         `db:"version"`
	Description sql.NullString `db:"desc"`
}

func (pl *productListTable) QueryProductList(prodType int, pagination *apibase.Pagination) ([]ProductListInfo, int) {
	where := dbhelper.MakeWhereCause()
	if prodType > 0 {
		where = where.Equal("type", prodType)
	}
	rows, count, err := pl.SelectWhere(dbhelper.GetDBColumnNamesFrom(ProductListInfo{}, "db"), where, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	products := []ProductListInfo{}
	for rows.Next() {
		d := ProductListInfo{}
		if err := rows.StructScan(&d); err != nil {
			apibase.ThrowDBModelError(err)
		}
		products = append(products, d)
	}
	return products, count
}

type ProductInfo struct {
	ID                int               `db:"id"`
	Type              int               `db:"type"`
	Name              string            `db:"name"`
	Version           string            `db:"version"`
	Description       sql.NullString    `db:"desc"`
	Url               string            `db:"url"`
	CheckSum          string            `db:"md5"`
	ReleaseDate       dbhelper.NullTime `db:"release_date"`
	DeployTimes       int               `db:"deploy_times"`
	DeployFailedTimes int               `db:"deploy_fail_times"`
}

func (pl *productListTable) GetProductInfo(id int) *ProductInfo {
	row := pl.SelectOneWhere(nil, dbhelper.MakeWhereCause().Equal("id", id))
	if row != nil {
		info := ProductInfo{}
		err := row.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		return &info
	} else {
		apibase.ThrowDBModelError("Cannot find product info where id = %d", id)
	}
	return nil
}
