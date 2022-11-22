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
	"database/sql"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"fmt"
	"time"
)

type clusterImageStore struct {
	dbhelper.DbTable
}

var ClusterImageStore = &clusterImageStore{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_CLUSTER_IMAGE_STORE},
}

type ImageStore struct {
	Id         int               `db:"id" json:"id"`
	ClusterId  int               `db:"clusterId" json:"clusterId"`
	IsDefault  int               `db:"is_default" json:"is_default"`
	Name       string            `db:"name" json:"name"`
	Alias      string            `db:"alias" json:"alias"`
	Address    string            `db:"address" json:"address"`
	Username   string            `db:"username" json:"username"`
	Password   string            `db:"password" json:"password"`
	Email      string            `db:"email" json:"email"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

func (l *clusterImageStore) InsertImageStore(store ImageStore) (int, error) {
	info := ImageStore{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("name", store.Name).And().
		Equal("clusterId", store.ClusterId).And().
		Equal("is_deleted", 0), &info)
	if err != nil && err == sql.ErrNoRows {
		isDefault := 0
		_, err = l.GetDefaultStoreByClusterId(store.ClusterId)
		if err == sql.ErrNoRows {
			isDefault = 1
		}
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"clusterId":  store.ClusterId,
			"name":       store.Name,
			"alias":      store.Alias,
			"address":    store.Address,
			"username":   store.Username,
			"password":   store.Password,
			"email":      store.Email,
			"is_default": isDefault,
		})
		if err != nil {
			apibase.ThrowDBModelError(err)
			return -1, err
		}
		id, _ := ret.LastInsertId()
		return int(id), err
	} else if err == nil {
		return info.Id, fmt.Errorf("镜像仓库:%v 已存在", store.Name)
	} else {
		return -1, err
	}
}

func (l *clusterImageStore) DeleteImageStoreById(id int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)

	return err
}

func (l *clusterImageStore) UpdateImageStoreById(store ImageStore) error {
	info := ImageStore{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().
		Equal("name", store.Name).And().
		Equal("clusterId", store.ClusterId).And().
		Equal("is_deleted", 0), &info)
	if err == nil && info.Id != store.Id {
		return fmt.Errorf("镜像仓库:%v 已存在", store.Name)
	}
	err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", store.Id), dbhelper.UpdateFields{
		"clusterId":   store.ClusterId,
		"name":        store.Name,
		"alias":       store.Alias,
		"address":     store.Address,
		"username":    store.Username,
		"password":    store.Password,
		"email":       store.Email,
		"update_time": time.Now(),
	}, false)

	return err
}

func (l *clusterImageStore) GetImageStoreInfoByClusterId(clusterId int) ([]ImageStore, error) {
	info := make([]ImageStore, 0)
	err := ClusterImageStore.GetDB().Select(&info, "select * from "+ClusterImageStore.TableName+
		" where clusterId=? and is_deleted=0", clusterId)
	return info, err
}

func (l *clusterImageStore) GetImageStoreInfoById(id int) (ImageStore, error) {
	info := ImageStore{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("id", id).And().Equal("is_deleted", 0), &info)
	return info, err
}

func (l *clusterImageStore) SetDefaultById(id, clusterId int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("clusterId", clusterId), dbhelper.UpdateFields{
		"is_default":  0,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		return err
	}
	err = l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"is_default":  1,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *clusterImageStore) GetDefaultStoreByClusterId(clusterId int) (ImageStore, error) {
	info := ImageStore{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("clusterId", clusterId).And().
		Equal("is_deleted", 0).And().Equal("is_default", 1), &info)
	return info, err
}

func (l *clusterImageStore) GetStoreByClusterIdAndNamespace(clusterId int, namespace string) (ImageStore, error) {
	info := ImageStore{}
	status := "valid"
	query := "select imagestore.alias,imagestore.address,imagestore.username,imagestore.password from "+ ClusterImageStore.TableName+
		" as imagestore LEFT JOIN deploy_namespace_list as list  on imagestore.id = list.registry_id where list.cluster_id=? and namespace=? and list.is_deleted=0 and list.status=?"
	err := ClusterImageStore.GetDB().Get(&info,query,clusterId,namespace,status)

	return info,err
}