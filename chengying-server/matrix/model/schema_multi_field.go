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
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"

	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
)

type schemaMultiField struct {
	dbhelper.DbTable
}

var SchemaMultiField = &schemaMultiField{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_SCHEMA_MULTI_FIELD},
}

type SchemaMultiFieldInfo struct {
	Id          int       `db:"id"`
	ClusterId   int       `db:"cluster_id"`
	ProductName string    `db:"product_name"`
	ServiceName string    `db:"service_name"`
	FieldPath   string    `db:"field_path"`
	Field       string    `db:"field"`
	Hosts       string    `db:"hosts"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
	IsDeleted   int       `db:"is_deleted"`
}

func (m SchemaMultiFieldInfo) Value() (driver.Value, error) {
	return []interface{}{m.ClusterId, m.ProductName, m.ServiceName, m.FieldPath, m.Field, m.Hosts, m.CreateTime,
		m.UpdateTime, m.IsDeleted}, nil
}

type FieldConfig struct {
	Hosts string `json:"hosts"`
	Field string `json:"field"`
}

func (s *schemaMultiField) BatchInsert(multiFieldList []interface{}, tx *sqlx.Tx) error {
	if len(multiFieldList) > 0 {
		argsList := make([]string, 0)
		for range multiFieldList {
			argsList = append(argsList, "(?)")
		}
		query, args, err := sqlx.In(fmt.Sprintf("INSERT INTO %s (cluster_id,product_name,service_name,field_path,field,hosts,"+
			"create_time,update_time,is_deleted) VALUES %s", s.TableName, strings.Join(argsList, ",")), multiFieldList...)
		if err != nil {
			return err
		}
		_, err = tx.Exec(query, args...)
		return err
	}
	return nil
}

func (s *schemaMultiField) GetByProductNameAndServiceNameAndPath(clusterId int, productName, serviceName, fieldPath string) ([]string, error) {
	var list []string
	query := fmt.Sprintf("select cast(id as char) from %s where cluster_id=? and product_name=? and service_name=? and field_path=? and is_deleted=0",
		TBL_SCHEMA_MULTI_FIELD)
	err := USE_MYSQL_DB().Select(&list, query, clusterId, productName, serviceName, fieldPath)
	return list, err
}

func (s *schemaMultiField) GetListByFieldPath(clusterId int, productName, serviceName, fieldPath string) ([]SchemaMultiFieldInfo, error) {
	var list []SchemaMultiFieldInfo
	query := fmt.Sprintf("select cluster_id, product_name, service_name, field_path, field, hosts from %s where cluster_id = ? and product_name = ? and service_name = ? and field_path = ? and is_deleted=0", s.TableName)
	err := USE_MYSQL_DB().Select(&list, query, clusterId, productName, serviceName, fieldPath)
	return list, err
}

func (s *schemaMultiField) GetByProductNameAndServiceNameAndIp(clusterId int, productName, serviceName, ip string) ([]SchemaMultiFieldInfo, error) {
	var list []SchemaMultiFieldInfo
	query := fmt.Sprintf("select * from %s where cluster_id=? and product_name=? and service_name=? and hosts=? and is_deleted=0", TBL_SCHEMA_MULTI_FIELD)
	err := USE_MYSQL_DB().Select(&list, query, clusterId, productName, serviceName, ip)
	return list, err
}

func (s *schemaMultiField) GetByProductNameAndServiceNameAndClusterId(productName, serviceName string, clusterId int) ([]SchemaMultiFieldInfo, error) {
	var list []SchemaMultiFieldInfo
	query := fmt.Sprintf("select * from %s where product_name=? and service_name=? and cluster_id=? and is_deleted=0 order by id asc", TBL_SCHEMA_MULTI_FIELD)
	err := USE_MYSQL_DB().Select(&list, query, productName, serviceName, clusterId)
	return list, err
}

func (s *schemaMultiField) DeleteByIp(clusterId int, ip string) error {
	query := fmt.Sprintf("delete from %s where hosts=? and cluster_id=?", s.TableName)
	_, err := USE_MYSQL_DB().Exec(query, ip, clusterId)
	if err != nil {
		log.Errorf("delete multi field error: %v", err)
		return err
	}
	return nil
}

func (s *schemaMultiField) DeleteByProductNameAndServiceNameAndClusterId(productName, serviceName string, clusterId int, tx *sqlx.Tx) error {
	query := fmt.Sprintf("delete from %s where product_name=? and service_name=? and cluster_id=? and is_deleted=0", s.TableName)
	_, err := tx.Exec(query, productName, serviceName, clusterId)
	if err != nil {
		log.Errorf("delete multi field error: %v", err)
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func (s *schemaMultiField) DeleteByFieldPath(clusterId int, productName, serviceName, fieldPath string) error {
	query := fmt.Sprintf("delete from %s where product_name=? and service_name=? and cluster_id =? and field_path = ? and is_deleted=0", s.TableName)
	_, err := USE_MYSQL_DB().Exec(query, productName, serviceName, clusterId, fieldPath)
	if err != nil {
		log.Errorf("delete multi field by path error: %v", err)
		return err
	}
	return err
}

type FieldCount struct {
	FieldPath string `db:"field_path"`
	Field     string `db:"field"`
	Count     int    `db:"cnt"`
}

func (s *schemaMultiField) GetDistinctPathCount(clusterId int, productName, serviceName string) []FieldCount {
	query := fmt.Sprintf("select field_path, field, count(*) as cnt from %s where cluster_id =? and product_name=? and service_name=? and is_deleted = 0 group by field_path", s.TableName)
	var fieldCounts = make([]FieldCount, 0)
	err := USE_MYSQL_DB().Select(&fieldCounts, query, clusterId, productName, serviceName)
	if err != nil {
		log.Errorf("Validate multi schema field error:%s", err)
	}
	return fieldCounts
}

func (s *schemaMultiField) GetDistinctValueCount(clusterId int, productName, serviceName string) []FieldCount {
	query := fmt.Sprintf("select field_path, field, count(*) as cnt from %s where cluster_id =? and product_name=? and service_name=? and is_deleted = 0 group by field, field_path", s.TableName)
	var fieldCounts = make([]FieldCount, 0)
	err := USE_MYSQL_DB().Select(&fieldCounts, query, clusterId, productName, serviceName)
	if err != nil {
		log.Errorf("Validate multi schema field error:%s", err)
	}
	return fieldCounts
}

//func (s *schemaMultiField) Validate(clusterId int, productName, serviceName string, tx *sqlx.Tx) error {
//	query := fmt.Sprintf("select field_path, count(*) as cnt from %s where cluster_id =? and product_name=? and service_name=? and is_deleted = 0 group by field_path", s.TableName)
//	type fieldCount struct {
//		FieldPath string `db:"field_path"`
//		Count     int    `db:"cnt"`
//	}
//	var fieldCounts = make([]fieldCount, 0)
//	err := USE_MYSQL_DB().Select(&fieldCounts, query, clusterId, productName, serviceName)
//	if err != nil {
//		log.Errorf("Validate multi schema field error:%s", err)
//		return err
//	}
//	for _, fieldCountItem := range fieldCounts {
//		// 当剩余配置参数记录数小于2时，表示已经不存在配置富化，直接删除记录
//		if fieldCountItem.Count < 2 {
//			err = s.DeleteByFieldPath(clusterId, productName, serviceName, fieldCountItem.FieldPath, tx)
//			if err != nil {
//				log.Errorf("Validate delete multi schema field error:%s", err)
//			}
//		}
//	}
//	if err := tx.Commit(); err != nil {
//		log.Errorf("%v", err)
//		return err
//	}
//	return nil
//}

type ProductServiceTuple struct {
	ProductName string `db:"product_name"`
	ServiceName string `db:"service_name"`
}

func (s *schemaMultiField) GetProductServiceByIp(clusterId int, hosts string) ([]ProductServiceTuple, error) {
	query := fmt.Sprintf("select product_name, service_name from %s where cluster_id = ? and hosts = ? and is_deleted = 0 group by product_name, service_name", s.TableName)
	var tupleList []ProductServiceTuple
	err := USE_MYSQL_DB().Select(&tupleList, query, clusterId, hosts)
	if err != nil {
		log.Errorf("Get product service tuple error: %v", err)
		return nil, err
	}
	return tupleList, nil
}

func (s *schemaMultiField) UpdateField(info SchemaMultiFieldInfo) error {
	query := fmt.Sprintf("update %s set field = ? where cluster_id = ? and product_name = ? and service_name = ? and field_path = ? and is_deleted = 0 and hosts = ?", s.TableName)
	_, err := USE_MYSQL_DB().Exec(query, info.Field, info.ClusterId, info.ProductName, info.ServiceName, info.FieldPath, info.Hosts)
	return err
}

func (s *schemaMultiField) WhetherChangeField(clusterId int, productName, serviceName, filedPath string, fields []FieldConfig) bool {
	infos := []SchemaMultiFieldInfo{}
	query := fmt.Sprintf("select * from %s where cluster_id = ? and product_name = ? and service_name = ? and field_path = ?", s.TableName)
	err := s.GetDB().Select(&infos, query, clusterId, productName, serviceName, filedPath)
	if err != nil {
		log.Errorf("[schemaMultiFieldIsChangeField] %s", err)
		return false
	}
	// fields = [{"hosts":"172.16.82.232,172.16.82.233","field":"1"},{"hosts":"172.16.82.234","field":"5"}]
	fieldMap := map[string]string{}
	for _, info := range infos {
		fieldMap[info.Hosts] = info.Field
	}
	for _, field := range fields {
		ips := strings.Split(field.Hosts, ",")
		for _, ip := range ips {
			if v, ok := fieldMap[ip]; ok && v == field.Field {
				continue
			}
			return false
		}
	}
	return true
}

func (s *schemaMultiField) GetServiceModifyTime(clusterId int, productName, serviceName string) (time.Time, error) {
	info := SchemaMultiFieldInfo{}
	query := fmt.Sprintf("SELECT * FROM %s where cluster_id = ? and product_name = ? and service_name = ? order by update_time desc limit 1", s.TableName)
	if err := s.GetDB().Get(&info, query, clusterId, productName, serviceName); err != nil {
		return time.Time{}, err
	}
	return info.UpdateTime, nil
}

func (s *schemaMultiField) GetListByProductName(clusterId int, productName string) ([]SchemaMultiFieldInfo, error) {
	list := make([]SchemaMultiFieldInfo, 0)
	query := fmt.Sprintf("select cluster_id, product_name, service_name, field_path, field, hosts from %s where cluster_id = ? and product_name = ? and is_deleted=0", s.TableName)
	err := USE_MYSQL_DB().Select(&list, query, clusterId, productName)
	return list, err
}
