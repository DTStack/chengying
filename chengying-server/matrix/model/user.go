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
	"fmt"
	"strings"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/log"
)

type userList struct {
	dbhelper.DbTable
}

var UserList = &userList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_USER_LIST},
}

type userRoleList struct {
	dbhelper.DbTable
}

var UserRoleList = &userRoleList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_USER_ROLE},
}

type roleList struct {
	dbhelper.DbTable
}

var RoleList = &roleList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_ROLE_LIST},
}

type clusterRightList struct {
	dbhelper.DbTable
}

var ClusterRightList = &clusterRightList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_USER_CLUSTER_RIGHT},
}

const (
	ROLE_ADMIN_ID        = 1
	ROLE_OPS_ID          = 2
	ROLE_READER_ID       = 3
	PASSWORD_NOT_CHANGED = 0
	PASSWORD_CHANGED     = 1
)

type UserInfo struct {
	ID                  int               `db:"id" json:"id"`
	UserName            string            `db:"username" json:"username"`
	PassWord            string            `db:"password" json:"password"`
	Company             string            `db:"company" json:"company"`
	FullName            string            `db:"full_name" json:"full_name"`
	Email               string            `db:"email" json:"email"`
	Phone               string            `db:"phone" json:"phone"`
	Status              int               `db:"status" json:"status"`
	ResetPasswordStatus int               `db:"reset_password_status" json:"reset_password_status"`
	UpdateTime          dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime          dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted           int               `db:"is_deleted" json:"is_deleted"`
}

type RoleInfo struct {
	ID         int               `db:"id" json:"id"`
	RoleName   string            `db:"role_name" json:"role_name"`
	RoleValue  string            `db:"role_value" json:"role_value"`
	RoleDesc   string            `db:"role_desc" json:"role_desc"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"-"`
	IsDeleted  int               `db:"is_deleted" json:"-"`
}

type UserRole struct {
	ID         int               `db:"id" json:"id"`
	RoleId     int               `db:"role_id" json:"role_id"`
	UserId     int               `db:"user_id" json:"user_id"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

type ClusterRight struct {
	ID         int               `db:"id" json:"id"`
	UserId     int               `db:"user_id" json:"user_id"`
	ClusterId  int               `db:"cluster_id" json:"cluster_id"`
	UpdateTime dbhelper.NullTime `db:"update_time" json:"update_time"`
	CreateTime dbhelper.NullTime `db:"create_time" json:"create_time"`
	IsDeleted  int               `db:"is_deleted" json:"is_deleted"`
}

type ResClusterInfo struct {
	ClusterId   int    `db:"cluster_id" json:"cluster_id"`
	ClusterName string `db:"cluster_name" json:"cluster_name"`
	ClusterType string `db:"cluster_type" json:"cluster_type"`
}

type ResInfo struct {
	UserInfo
	ClusterList      []ResClusterInfo `db:"cluster_list" json:"cluster_list"`
	RoleId           int              `db:"role_id" json:"role_id"`
	RoleName         string           `db:"role_name" json:"role_name"`
	UpdateTimeFormat string
}

func (l *userList) InsertUserIfNotExist(userName, password, company, fullName, email, phone string) (error, int) {
	info := UserInfo{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("username", userName).And().Equal("is_deleted", 0), &info)
	if err != nil && err == sql.ErrNoRows {
		ret, err := l.InsertWhere(dbhelper.UpdateFields{
			"username":  userName,
			"password":  password,
			"company":   company,
			"full_name": fullName,
			"email":     email,
			"phone":     phone,
		})
		if err != nil {
			apibase.ThrowDBModelError(err)
			return err, -1
		}
		seq, _ := ret.LastInsertId()
		return nil, int(seq)
	} else if err == nil {
		aid := info.ID
		return fmt.Errorf("用户名:%v 已存在", userName), int(aid)
	} else {
		return err, 0
	}
}

func (l *userRoleList) InsertUserRole(userId, roleId int) (error, int) {
	ret, err := l.InsertWhere(dbhelper.UpdateFields{
		"role_id": roleId,
		"user_id": userId,
	})
	if err != nil {
		apibase.ThrowDBModelError(err)
		return err, -1
	}
	seq, _ := ret.LastInsertId()
	return nil, int(seq)
}

func (l *userList) UpdateInfoByUserId(company, fullName, email, phone string, userId int) bool {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", userId), dbhelper.UpdateFields{
		"company":     company,
		"full_name":   fullName,
		"email":       email,
		"phone":       phone,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[UpdatePwdByUserId] UpdatePwdByUserId err: %v", err)
		return false
	}
	return true
}

func (l *userList) UpdateStatusByUserId(status int, userId int) bool {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", userId), dbhelper.UpdateFields{
		"status":      status,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[UpdateStatusByUserId] UpdateStatusByUserId err: %v", err)
		return false
	}
	return true
}

func (l *userList) UpdatePwdByUserId(password string, userId int) bool {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", userId), dbhelper.UpdateFields{
		"password":              password,
		"reset_password_status": PASSWORD_CHANGED,
		"update_time":           time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[UpdatePwdByUserId] UpdatePwdByUserId err: %v", err)
		return false
	}
	return true
}

func (l *userRoleList) DeleteByUserId(userId int) bool {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("user_id", userId), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateStatus err: %v", err)
		return false
	}
	return true
}

func (l *userList) DeleteByUserId(userId int) bool {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", userId), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[deployHostList] UpdateStatus err: %v", err)
		return false
	}
	return true
}

func (l *userList) GetInfoByUserId(userId int) (error, ResInfo) {
	info := ResInfo{}
	query := `SELECT user_list.*, role_list.id AS role_id, role_list.role_name FROM user_list
	LEFT JOIN user_role ON user_list.id=user_role.user_id
	LEFT JOIN role_list ON user_role.role_id = role_list.id
	WHERE user_list.id=? and user_list.is_deleted=0 and user_role.is_deleted=0 and role_list.is_deleted=0`
	err := l.GetDB().Get(&info, query, userId)
	if err != nil {
		return err, info
	}
	return nil, info
}

var _getUserRoleFields = utils.GetTagValues(UserRole{}, "db")

func (l *userRoleList) ListUserRoleByUserId(userId int) (error, []UserRole) {
	rows, _, err := l.SelectWhere(_getUserRoleFields, dbhelper.MakeWhereCause().Equal("user_id", userId).And().Equal("is_deleted", 0), nil)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()
	list := []UserRole{}
	for rows.Next() {
		info := UserRole{}

		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return nil, list
}

func (l *userList) GetInfoByUserName(username string) (error, *UserInfo) {
	info := UserInfo{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("username", username).And().Equal("is_deleted", 0), &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

var _getUserInfoFields = utils.GetTagValues(UserInfo{}, "db")

func (l *userList) UserInfoList(username, status, roleId string, pagination *apibase.Pagination) ([]ResInfo, interface{}) {
	fields := []string{}
	for _, field := range _getUserInfoFields {
		if field == "password" {
			continue
		}
		fields = append(fields, field)
	}

	fromQuery := `FROM user_list
LEFT JOIN user_role ON user_list.id = user_role.user_id 
LEFT JOIN role_list ON user_role.role_id = role_list.id
WHERE user_list.id != -1 and user_list.is_deleted = 0 and user_role.is_deleted = 0 and role_list.is_deleted = 0`
	if status != "" {
		fromQuery += " and user_list.status in (" + status + ")"
	}
	if roleId != "" {
		fromQuery += " and role_list.id in (" + roleId + ")"
	}
	if username != "" {
		fromQuery += " and ( user_list.username like '%" + username + "%' or user_list.full_name like '%" + username + "%' )"
	}

	var err error
	var total int
	var list []ResInfo
	query := "SELECT COUNT(1) " + fromQuery
	stmt, err := USE_MYSQL_DB().Prepare(query)
	if err != nil {
		msg := fmt.Errorf("prepare sql failed, %v", err)
		apibase.ThrowDBModelError(msg)
	}
	defer stmt.Close()
	stmt.Exec()
	if err = USE_MYSQL_DB().Get(&total, query); err != nil {
		apibase.ThrowDBModelError(err)
	}

	if total > 0 {
		query := "SELECT user_list." + strings.Join(fields, ",user_list.") + " ,role_list.role_name, role_list.id AS role_id  " +
			fromQuery + " " + pagination.AsQuery()
		log.Debugf("UserInfoList query: %v", query)
		stmt, err := USE_MYSQL_DB().Prepare(query)
		if err != nil {
			msg := fmt.Errorf("prepare sql failed, %v", err)
			apibase.ThrowDBModelError(msg)
		}
		defer stmt.Close()
		stmt.Exec()
		if err = USE_MYSQL_DB().Select(&list, query); err != nil {
			apibase.ThrowDBModelError(err)
		}
	}

	return list, total
}

func (l *roleList) GetInfoByRoleId(roleId int) (error, *RoleInfo) {
	info := RoleInfo{}
	err := l.GetWhere(nil, dbhelper.MakeWhereCause().Equal("id", roleId).And().Equal("is_deleted", 0), &info)
	if err != nil {
		return err, &info
	}
	return nil, &info
}

func (l *roleList) GetList() ([]RoleInfo, error) {
	list := make([]RoleInfo, 0)
	err := l.GetDB().Select(&list, "select * from "+RoleList.TableName+
		" where id!=-1 and is_deleted=0")
	return list, err
}

func (l *clusterRightList) InsertUserClusterRight(userId int, clusterId int) (error, int) {
	ret, err := l.InsertWhere(dbhelper.UpdateFields{
		"user_id":    userId,
		"cluster_id": clusterId,
	})
	if err != nil {
		apibase.ThrowDBModelError(err)
		return err, -1
	}
	seq, _ := ret.LastInsertId()
	return nil, int(seq)
}

func (l *clusterRightList) UpdateUserClusterRightByUserId(userId int, clusterId string) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("user_id", userId), dbhelper.UpdateFields{
		"cluster_id":  clusterId,
		"update_time": time.Now(),
	}, false)
	return err
}

func (l *clusterRightList) DeleteByUserId(userId int) bool {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("user_id", userId), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[DeleteByUserId] DeleteByUserId err: %v", err)
		return false
	}
	return true
}

func (l *clusterRightList) DeleteByClusterId(clusterId int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("cluster_id", clusterId), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[DeleteByClusterId] DeleteByClusterId err: %v", err)
		return err
	}
	return err
}

func (l *clusterRightList) DeleteById(id int) error {
	err := l.UpdateWhere(dbhelper.MakeWhereCause().Equal("id", id), dbhelper.UpdateFields{
		"is_deleted":  1,
		"update_time": time.Now(),
	}, false)
	if err != nil {
		log.Errorf("[DeleteById] DeleteById err: %v", err)
		return err
	}
	return err
}

func (l *clusterRightList) GetUserClusterRightByUserId(userId int) ([]ClusterRight, error) {
	list := make([]ClusterRight, 0)
	query := "select * from user_cluster_right " +
		"WHERE user_id = ? " +
		"and is_deleted = 0 "
	err := USE_MYSQL_DB().Select(&list, query, userId)
	return list, err
}
