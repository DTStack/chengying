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
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	UICDB   *sqlx.DB
	MYSQLDB *sqlx.DB
)

const (
	TBL_PRODUCT_LIST     = "product_list"
	TBL_SIDECAR_LIST     = "sidecar_list"
	TBL_AGENT_LIST       = "agent_list"
	TBL_OP_HISTORY       = "operation_history"
	TBL_UPDATE_HISTORY   = "update_history"
	TBL_PROGRESS_HISTORY = "progress_history"
	TBL_RESOURCE_USAGES  = "resource_usages"
	TBL_DEPLOY_CALLBACK  = "deploy_callback"
	TBL_DASHBOARD_LIST   = "dashboard_list"

	TBL_TRIGGER_LIST  = "trigger_list"
	TBL_STRATEGY_LIST = "strategy_list"
)

func USE_MYSQL_DB() *sqlx.DB {
	return MYSQLDB
}

func USE_UIC_DB() *sqlx.DB {
	return UICDB
}

func connectDatabase(host, user, password, dbname string, port int) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true", user, password, host, port, dbname))
}

func ConfigureUicDatabase(host string, port int, user, password, dbname string) error {
	var err error
	UICDB, err = connectDatabase(host, user, password, dbname, port)
	return err
}

func ConfigureMysqlDatabase(host string, port int, user, password, dbname string) error {
	var err error
	MYSQLDB, err = connectDatabase(host, user, password, dbname, port)
	return err
}
