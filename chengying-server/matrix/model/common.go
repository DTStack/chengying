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
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	UICDB   *sqlx.DB
	MYSQLDB *sqlx.DB
)

const (
	TBL_DEPLOY_HOST                               = "deploy_host"
	TBL_DEPLOY_NODE                               = "deploy_node"
	TBL_DEPLOY_INSTANCE_LIST                      = "deploy_instance_list"
	TBL_DEPLOY_INSTANCE_RECORD                    = "deploy_instance_record"
	TBL_DEPLOY_INSTANCE_EVENT                     = "deploy_instance_event"
	TBL_DEPLOY_PRODUCT_LIST                       = "deploy_product_list"
	TBL_DEPLOY_PRODUCT_HISTORY                    = "deploy_product_history"
	TBL_DEPLOY_PRODUCT_UPDATE_HISTORY             = "deploy_product_update_history"
	TBL_DEPLOY_SCHEMA_FIELD_MODIFY                = "deploy_schema_field_modify"
	TBL_DEPLOY_SERVICE_IP_LIST                    = "deploy_service_ip_list"
	TBL_DEPLOY_UNCHECKED_SERVICE                  = "deploy_unchecked_service"
	TBL_USER_LIST                                 = "user_list"
	TBL_ROLE_LIST                                 = "role_list"
	TBL_DEPLOY_USER_ROLE                          = "user_role"
	TBL_USER_CLUSTER_RIGHT                        = "user_cluster_right"
	TBL_ADDONS_LIST                               = "addons_list"
	TBL_DEPLOY_ADDONS_LIST                        = "deploy_addons_list"
	TBL_EVENT_LIST                                = "deploy_instance_runtime_event"
	TBL_DEPLOY_STRATEGY_LIST                      = "deploy_strategy_list"
	TBL_DEPLOY_STRATEGY_RESOURCE_LIST             = "deploy_strategy_resource_list"
	TBL_DEPLOY_STRATEGY_ASSIGN_LIST               = "deploy_strategy_assign_list"
	TBL_DEPLOY_CLUSTER_IMAGE_STORE                = "deploy_cluster_image_store"
	TBL_DEPLOY_CLUSTER_LIST                       = "deploy_cluster_list"
	TBL_DEPLOY_CLUSTER_HOST_REL                   = "deploy_cluster_host_rel"
	TBL_PRODUCT_BACKUP_CONFIG                     = "product_backup_config"
	TBL_DEPLOY_CLUSTER_K8S_AVAILABLE              = "deploy_cluster_k8s_available"
	TBL_DEPLOY_CLUSTER_K8S_ONLY                   = "deploy_cluster_k8s_only"
	TBL_DEPLOY_CLUSTER_PRODUCT_REL                = "deploy_cluster_product_rel"
	TBL_DEPLOY_CLUSTER_SMOOTH_UPGRADE_PRODUCT_REL = "deploy_cluster_smooth_upgrade_product_rel"
	TBL_DEPLOY_CLUSTER_KUBE_POD_LIST              = "deploy_cluster_kube_pod_list"
	TBL_DEPLOY_CLUSTER_KUBE_SERVICE_LIST          = "deploy_cluster_kube_service_list"
	TBL_DEPLOY_KUBE_BASE_PRODUCT_LIST             = "deploy_kube_base_product_list"
	TBL_DEPLOY_KUBE_PRODUCT_LOCK                  = "deploy_kube_product_lock"
	TBL_SAFETY_AUDIT_LIST                         = "safety_audit_list"
	TBL_DEPLOY_INSTANCE_UPDATE_RECORD             = "deploy_instance_update_record"
	TBL_SCHEMA_MULTI_FIELD                        = "deploy_schema_multi_fields"
	TBL_INSPECT_REPORT_TEMPLATE                   = "inspect_report_template"
	TBL_NOTIFY_EVENT                              = "deploy_notify_event"
	HOST_ROLE                                     = "host_role"
	DEPLOY_PRODUCT_SELECT_HISTORY                 = "deploy_product_select_history"
	DEPLOY_UUID                                   = "deploy_uuid"
	TBL_SWITCH_RECORD                             = "deploy_switch_record"
	TBL_UPLOAD_RECORD                             = "deploy_upload_record"
	TBL_INSPECT_REPORT                            = "deploy_inspect_report"
	OPERATION_LIST                                = "operation_list"
	EXECSHELL_LIST                                = "exec_shell_list"
	TBL_AUTO_TEST                                 = "smoke_testing"
	TBL_SERVICE_HEALTH_CHECK                      = "service_health_check"
	TBL_UPGRADE_HISTORY                           = "deploy_upgrade_history"
	TBL_BACKUP_HISTORY                            = "deploy_backup_history"
	TBL_DEPLOY_SMOOTH_UPGRADE_LIST                = "deploy_smooth_upgrade_list"
	TBL_DEPLOY_MYSQL_IP_LIST                      = "deploy_mysql_ip_list"
	TBL_TASK_LIST                                 = "task_list"
	TBL_TASK_HOST                                 = "task_host"
	TBL_TASK_LOG                                  = "task_log"
	TBL_DEPLOY_PRODUCT_LINE_LIST                  = "deploy_product_line_list"
	TBL_DEPLOY_SERVICE_RELATIONS_LIST             = "deploy_service_relations_list"
)

func USE_MYSQL_DB() *sqlx.DB {
	return MYSQLDB
}

func USE_UIC_DB() *sqlx.DB {
	return UICDB
}

func connectDatabase(host, user, password, dbname string, port int) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=Local&parseTime=true", user, password, host, port, dbname))
}

func ConfigureMysqlDatabase(host string, port int, user, password, dbname string) error {
	var err error
	MYSQLDB, err = connectDatabase(host, user, password, dbname, port)
	trycount := 0
	for {
		if trycount >= 3 {
			break
		}
		err := MYSQLDB.Ping()
		if err != nil {
			trycount++
			continue
		}
		break
	}
	cache.Db = MYSQLDB
	return err
}
