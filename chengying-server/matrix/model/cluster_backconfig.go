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
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"time"
)

type clusterBackupConfig struct {
	dbhelper.DbTable
}

var ClusterBackupConfig = &clusterBackupConfig{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_PRODUCT_BACKUP_CONFIG},
}

type clusterBackupConfigModel struct {
	Id          int               `db:"id"  json:"-"`
	ClusterId   int               `db:"cluster_id" json:"clusterId"`
	ClusterName string            `json:"clusterName"`
	ConfigPath  string            `db:"config_path" json:"path"`
	CreateTime  dbhelper.NullTime `db:"create_time" json:"-"`
	UpdateTime  dbhelper.NullTime `db:"update_time" json:"-"`
}

var _getclsBCFields = utils.GetTagValues(clusterBackupConfigModel{}, "db")

func (c *clusterBackupConfig) GetClusterPathConfigID(clusterId int) (*clusterBackupConfigModel, error) {
	config := clusterBackupConfigModel{}
	if err := c.GetWhere(_getclsBCFields, dbhelper.MakeWhereCause().Equal("cluster_id", clusterId), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *clusterBackupConfig) GetPathByClusterId(clusterId int) string {
	var (
		path string
	)
	err := c.GetDB().Get(&path, "select config_path from "+c.TableName+" where cluster_id = ?", clusterId)
	if path == "" || err != nil {
		path = base.INSTALL_CURRRENT_PATH
	}
	return path
}

func (c *clusterBackupConfig) GetClusterPathConfigs(id int) (*clusterBackupConfigModel, error) {
	var (
		config clusterBackupConfigModel
	)
	err := ClusterBackupConfig.GetDB().Get(&config, "SELECT * FROM product_backup_config WHERE cluster_id = ? ", id)
	return &config, err
}

func (c *clusterBackupConfig) SaveClusterBackupConfig(clusterId int, path string) error {
	config, _ := c.GetClusterPathConfigID(clusterId)
	if config == nil {
		_, err := c.InsertWhere(dbhelper.UpdateFields{
			"cluster_id":  clusterId,
			"config_path": path,
			"update_time": time.Now(),
			"create_time": time.Now(),
		})
		return err
	}
	err := c.UpdateWhere(dbhelper.MakeWhereCause().Equal("cluster_id", clusterId), dbhelper.UpdateFields{
		"config_path": path,
		"update_time": time.Now(),
	}, false)
	return err
}
