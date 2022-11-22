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

package upgrade

import (
	"database/sql"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"time"
)

const (
	SMOOTH_UPGRADE_MODE   = "smooth"
	STANDARD_UPGRADE_MODE = "standard"
)

type upgradeHistory struct {
	dbhelper.DbTable
}

var UpgradeHistory = &upgradeHistory{
	dbhelper.DbTable{
		GetDB:     model.USE_MYSQL_DB,
		TableName: model.TBL_UPGRADE_HISTORY,
	},
}

type HistoryInfo struct {
	Id                int          `db:"id" json:"id"`
	ClusterId         int          `db:"cluster_id" json:"cluster_id"`
	ProductName       string       `db:"product_name" json:"product_name"`
	SourceVersion     string       `db:"source_version" json:"source_version"`
	TargetVersion     string       `db:"target_version" json:"target_version"`
	BackupName        string       `db:"backup_name" json:"backup_name"`
	SourceServiceIp   []byte       `db:"source_service_ip" json:"source_service_ip"`
	SourceConfig      []byte       `db:"source_config" json:"source_config"`
	SourceMultiConfig []byte       `db:"source_multi_config" json:"source_multi_config"`
	CreateTime        sql.NullTime `db:"create_time" json:"create_time"`
	Type              int          `db:"type" json:"type"`
	BackupSql         string       `db:"backup_sql" json:"backup_sql"`
	UpgradeMode       string       `db:"upgrade_mode" json:"upgrade_mode"`
	IsDeleted         bool         `db:"is_deleted" json:"is_deleted"`
}

func (u *upgradeHistory) InsertRecord(clusterId, upgradeType int, upgradeMode, productName, sourceVersion, targetVersion, backupName, backupSql string,
	serviceIp, sourceConfig, sourceMultiConfig []byte) (int64, error) {
	var history HistoryInfo
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("cluster_id", clusterId).And().
		Equal("product_name", productName).And().
		Equal("source_version", sourceVersion).And().
		Equal("target_version", targetVersion).And().
		Equal("upgrade_mode", upgradeMode).And().
		Equal("is_deleted", 0)
	if upgradeMode != SMOOTH_UPGRADE_MODE {
		whereCause = whereCause.And().Equal("backup_name", backupName)
	}
	err := u.GetWhere(nil, whereCause, &history)
	if err != nil && err == sql.ErrNoRows {
		r, err := u.InsertWhere(dbhelper.UpdateFields{
			"cluster_id":          clusterId,
			"product_name":        productName,
			"source_version":      sourceVersion,
			"target_version":      targetVersion,
			"backup_name":         backupName,
			"source_service_ip":   serviceIp,
			"source_config":       sourceConfig,
			"source_multi_config": sourceMultiConfig,
			"create_time":         time.Now(),
			"type":                upgradeType,
			"backup_sql":          backupSql,
			"upgrade_mode":        upgradeMode,
			"is_deleted":          0,
		})
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	} else if err == nil {
		return int64(history.Id), nil
	} else {
		return 0, err
	}
}

func (u *upgradeHistory) GetByClsAndProductNameAndSourceVersion(clusterId int, productName, sourceVersion, upgradeMode string) ([]HistoryInfo, error) {
	whereClause := dbhelper.MakeWhereCause().Equal("cluster_id", clusterId).And().
		Equal("product_name", productName).And().
		Equal("is_deleted", 0)
	if sourceVersion != "" {
		whereClause = whereClause.And().Equal("source_version", sourceVersion)
	}
	if upgradeMode != "" {
		whereClause = whereClause.And().Equal("upgrade_mode", upgradeMode)
	}
	rows, _, err := u.SelectWhere(nil, whereClause, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("rows close error: %v", err)
			return
		}
	}()
	var infoList []HistoryInfo
	for rows.Next() {
		row := HistoryInfo{}
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, row)
	}
	return infoList, nil
}

func (u *upgradeHistory) GetOne(clusterId int, productName, sourceVersion, backupName string) (*HistoryInfo, error) {
	whereClause := dbhelper.MakeWhereCause().Equal("cluster_id", clusterId).And().
		Equal("product_name", productName).And().
		Equal("source_version", sourceVersion).And().
		Equal("backup_name", backupName).And().
		Equal("is_deleted", 0)
	var info HistoryInfo
	err := u.GetWhere(nil, whereClause, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (u *upgradeHistory) GetTargetVersionInfo(clusterId int, productName, sourceVersion, productType, upgradeMode string) ([]model.DeployProductListInfoWithNamespace, error) {
	list := make([]model.DeployProductListInfoWithNamespace, 0)
	var values []interface{}
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_version, p.product_type from deploy_product_list as p " +
		"LEFT JOIN deploy_upgrade_history as u ON p.product_version = u.target_version AND p.product_name = u.product_name " +
		"WHERE u.cluster_id = ? AND u.product_name = ? AND u.source_version = ? AND u.upgrade_mode = ? AND u.is_deleted = 0"
	values = append(values, clusterId, productName, sourceVersion, upgradeMode)
	if productType != "" {
		query += " AND p.product_type = ?"
		values = append(values, productType)
	}
	err := model.USE_MYSQL_DB().Select(&list, query, values...)
	if err != nil {
		return list, err
	}
	return list, nil
}
