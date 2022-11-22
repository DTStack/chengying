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

package model

import (
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/base"
)

type addonsList struct {
	dbhelper.DbTable
}

var AddonsList = &addonsList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_ADDONS_LIST},
}

type AddonsInfo struct {
	Id         int       `db:"id"`
	Type       string    `db:"type"`
	Desc       string    `db:"desc"`
	Os         string    `db:"os"`
	Version    string    `db:"version"`
	Schema     string    `db:"schema"`
	IsDeleted  int       `db:"isDeleted" json:"-"`
	UpdateDate base.Time `db:"updated" json:"updated"`
	CreateDate base.Time `db:"created" json:"created"`
}

func (l *addonsList) GetAddonInfoById(aid string) (error, *AddonsInfo) {
	whereCause := dbhelper.WhereCause{}
	info := AddonsInfo{}
	err := l.GetWhere(nil, whereCause.Equal("id", aid), &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}
