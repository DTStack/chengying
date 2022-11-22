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

package kube

import (
	"database/sql"
)

type DeployProductListSchema struct {
	ID                 int               `db:"id"`
	ParentProductName  string            `db:"parent_product_name"`
	ProductName        string            `db:"product_name"`
	ProductNameDisplay string            `db:"product_name_display"`
	ProductVersion     string            `db:"product_version"`
	Product            []byte            `db:"product"`
	ProductParsed      []byte            `db:"product_parsed"`
	IsCurrentVersion   int               `db:"is_current_version"`
	Status             string            `db:"status"`
	DeployUUID         string            `db:"deploy_uuid"`
	AlertRecover       int               `db:"alert_recover"`
	DeployTime         sql.NullTime		 `db:"deploy_time"`
	CreateTime         sql.NullTime		 `db:"create_time"`
	UserId             int               `db:"user_id"`
	Schema             []byte            `db:"schema"`
	ProductType        int               `db:"product_type"`
}
