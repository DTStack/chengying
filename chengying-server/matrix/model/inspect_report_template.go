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
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
)

type inspectReportTemplate struct {
	dbhelper.DbTable
}

var InspectReportTemplate = &inspectReportTemplate{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_INSPECT_REPORT_TEMPLATE},
}

type InspectReportTemplateInfo struct {
	BaseTemplateConfig
	Id         int               `db:"id" json:"id"`
	ClusterId  int               `db:"cluster_id" json:"cluster_id"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

type BaseTemplateConfig struct {
	Type    int    `db:"type" json:"type"`
	Module  string `db:"module" json:"module"`
	Metric  string `db:"metric" json:"metric"`
	Targets string `db:"targets" json:"targets"`
	Unit    string `db:"unit" json:"unit"`
	Decimal int    `db:"decimal" json:"decimal"`
}

func (i *inspectReportTemplate) GetTemplateConfig(clusterId int) ([]BaseTemplateConfig, error) {
	var configList []BaseTemplateConfig
	query := "select `type`, `module`, `metric`, targets, unit, `decimal` from " + TBL_INSPECT_REPORT_TEMPLATE + " where is_deleted=0"
	if err := USE_MYSQL_DB().Select(&configList, query); err != nil {
		log.Errorf("get report template config error: %v", err)
		return nil, err
	}
	return configList, nil
}

func (i *inspectReportTemplate) GetPlatformTemplateConfig() ([]BaseTemplateConfig, error) {
	var configList []BaseTemplateConfig
	query := "select `type`, `module`, `metric`, targets, unit, `decimal` " +
		"from " + TBL_INSPECT_REPORT_TEMPLATE + " where is_deleted=0 and type > 1"
	if err := USE_MYSQL_DB().Select(&configList, query); err != nil {
		log.Errorf("get report template config error: %v", err)
		return nil, err
	}
	return configList, nil
}
