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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"github.com/satori/go.uuid"
)

type deploySchemaFieldModify struct {
	dbhelper.DbTable
}

var DeploySchemaFieldModify = &deploySchemaFieldModify{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_SCHEMA_FIELD_MODIFY},
}

const (
	PRODUCT_STATUS_UNDEPLOYED    = "undeployed"
	PRODUCT_STATUS_PENDING       = "pending"
	PRODUCT_STATUS_DEPLOYING     = "deploying"
	PRODUCT_STATUS_DEPLOYED      = "deployed"
	PRODUCT_STATUS_DEPLOY_FAIL   = "deploy fail"
	PRODUCT_STATUS_UNDEPLOYING   = "undeploying"
	PRODUCT_STATUS_UNDEPLOY_FAIL = "undeploy fail"
	PRODUCT_HEALTH_OK            = 1
	PRODUCT_HEALTH_BAD           = 0
)

type SchemaFieldModifyInfo struct {
	ID          int       `db:"id"`
	ClutserId   int       `db:"cluster_id"`
	ProductName string    `db:"product_name"`
	ServiceName string    `db:"service_name"`
	FieldPath   string    `db:"field_path"`
	Field       string    `db:"field"`
	UpdateDate  time.Time `db:"update_time"`
	CreateDate  time.Time `db:"create_time"`
	Namespace   string    `db:"namespace"`
}

func (s SchemaFieldModifyInfo) Value() (driver.Value, error) {
	return []interface{}{s.ClutserId, s.ProductName, s.ServiceName, s.FieldPath, s.Field, s.CreateDate, s.UpdateDate,
		s.Namespace}, nil
}

type deployProdcutList struct {
	dbhelper.DbTable
}

var DeployProductList = &deployProdcutList{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_PRODUCT_LIST},
}

type DeployProductListInfo struct {
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
	DeployTime         dbhelper.NullTime `db:"deploy_time"`
	CreateTime         dbhelper.NullTime `db:"create_time"`
	UserId             int               `db:"user_id"`
	Schema             []byte            `db:"schema"`
	ProductType        int               `db:"product_type"`
	Namespace          string            `db:"namespace"`
}

type DeployProductListInfoWithNamespace struct {
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
	DeployTime         dbhelper.NullTime `db:"deploy_time"`
	CreateTime         dbhelper.NullTime `db:"create_time"`
	UserId             int               `db:"user_id"`
	Schema             []byte            `db:"schema"`
	ProductType        int               `db:"product_type"`
	Namespace          string            `db:"namespace"`
}

type deployProductHistory struct {
	dbhelper.DbTable
}

var DeployProductHistory = &deployProductHistory{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_PRODUCT_HISTORY},
}

type DeployProductHistoryInfo struct {
	ID                 int               `db:"id"`
	ClusterId          int               `db:"cluster_id"`
	Namespace          string            `json:"namespace"`
	DeployUUID         uuid.UUID         `db:"deploy_uuid"`
	ProductName        string            `db:"product_name"`
	ProductNameDisplay string            `db:"product_name_display"`
	ParentProductName  string            `db:"parent_product_name"`
	ProductVersion     string            `db:"product_version"`
	ProductType        int               `db:"product_type"`
	Status             string            `db:"status"`
	CreateTime         dbhelper.NullTime `db:"create_time"`
	DeployStartTime    dbhelper.NullTime `db:"deploy_start_time"`
	DeployEndTime      dbhelper.NullTime `db:"deploy_end_time"`
	UserId             int               `db:"user_id"`
}

type deployProductUpdateHistory struct {
	dbhelper.DbTable
}

var DeployProductUpdateHistory = &deployProductUpdateHistory{
	dbhelper.DbTable{USE_MYSQL_DB, TBL_DEPLOY_PRODUCT_UPDATE_HISTORY},
}

type DeployProductUpdateHistoryInfo struct {
	ID                 int               `db:"id"`
	ClusterId          int               `db:"cluster_id"`
	Namespace          string            `json:"namespace"`
	UpdateUUID         uuid.UUID         `db:"update_uuid"`
	ProductName        string            `db:"product_name"`
	ProductNameDisplay string            `db:"product_name_display"`
	ParentProductName  string            `db:"parent_product_name"`
	ProductVersion     string            `db:"product_version"`
	ProductType        int               `db:"product_type"`
	Status             string            `db:"status"`
	CreateTime         dbhelper.NullTime `db:"create_time"`
	UpdateStartTime    dbhelper.NullTime `db:"update_start_time"`
	UpdateEndTime      dbhelper.NullTime `db:"update_end_time"`
	UserId             int               `db:"user_id"`
	PackageName        string            `db:"package_name"`
	UpdateDir          string            `db:"update_dir"`
	BackupDir          string            `db:"backup_dir"`
	ProductId          int               `db:"product_id"`
}

var _getDeployProductListFields = utils.GetTagValues(DeployProductListInfo{}, "db")

func (d *deployProdcutList) GetProductListByWhere(cause dbhelper.WhereCause) ([]DeployProductListInfo, error) {
	rows, _, err := d.SelectWhere(_getDeployProductListFields, cause, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DeployProductListInfo{}
	for rows.Next() {
		info := DeployProductListInfo{}
		err = rows.StructScan(&info)
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func (d *deployProdcutList) GetDeployProductList(pagination *apibase.Pagination,
	parentProductName, productNames, productName, productVersionLike, productVersion, productType string) ([]DeployProductListInfo, int) {
	fields := []string{"id", "parent_product_name", "product_name", "product_name_display", "product_version", "product_type", "status", "deploy_time", "create_time"}
	whereCause := dbhelper.WhereCause{}
	var values []interface{}
	whereCause = whereCause.GreaterThan("id", "0")

	if productType != "" {
		pType, _ := strconv.Atoi(productType)
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("product_type", pType)
	}

	if parentProductName != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("parent_product_name", parentProductName)
	}

	if productNames != "" {
		for _, v := range strings.Split(productNames, ",") {
			values = append(values, v)
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("product_name", values...)
	}

	if productName != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("product_name", productName)
	}

	if productVersionLike != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Like("product_version", "%"+productVersionLike+"%")
	}

	if productVersion != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("product_version", productVersion)
	}

	rows, totalproducts, err := d.SelectWhere(fields, whereCause, pagination)

	if err != nil {
		apibase.ThrowDBModelError(err)
	}

	defer rows.Close()

	list := []DeployProductListInfo{}

	for rows.Next() {
		info := DeployProductListInfo{}

		//todo  product

		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}

	return list, totalproducts

}

func (d *deployProdcutList) GetDeployParentProductList() []string {
	list := []string{}

	err := DeployProductList.GetDB().Select(&list, "select distinct parent_product_name from "+DeployProductList.TableName)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	return list
}

func (d *deployProdcutList) GetDeployProductName(productType int) []string {
	list := []string{}

	err := DeployProductList.GetDB().Select(&list, "select distinct product_name from "+DeployProductList.TableName+" where product_type=?", productType)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	return list
}

func (d *deployProdcutList) GetProductList(productName, productType string, deployStatus []string, pagination *apibase.Pagination) ([]DeployProductListInfo, int) {
	fields := []string{"id", "parent_product_name", "product_name", "product_name_display", "product_version", "product_type", "status", "deploy_time", "create_time"}
	whereCause := dbhelper.MakeWhereCause().GreaterThan("id", "0")
	if productName != "" {
		whereCause = whereCause.And().Equal("product_name", productName)
	}
	if productType != "" {
		whereCause = whereCause.And().Equal("product_type", productType)
	}
	if len(deployStatus) > 0 {
		status := make([]interface{}, 0, len(deployStatus))
		for _, ds := range deployStatus {
			status = append(status, interface{}(ds))
		}
		whereCause = whereCause.And().Included("status", status...)
	}
	rows, total, err := d.SelectWhere(fields, whereCause, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()
	list := []DeployProductListInfo{}
	for rows.Next() {
		info := DeployProductListInfo{}

		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return list, total
}

func (d *deployProdcutList) GetProductListByNameAndType(productName, productType string, pagination *apibase.Pagination) ([]DeployProductListInfo, int) {
	fields := []string{"id", "parent_product_name", "product_name", "product_name_display", "product_version", "product_type", "product", "schema", "deploy_time", "create_time"}
	whereCause := dbhelper.MakeWhereCause().GreaterThan("id", "0")
	if productName != "" {
		whereCause = whereCause.And().Equal("product_name", productName)
	}
	if productType != "" {
		whereCause = whereCause.And().Equal("product_type", productType)
	}
	rows, total, err := d.SelectWhere(fields, whereCause, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}
	defer rows.Close()
	list := []DeployProductListInfo{}
	for rows.Next() {
		info := DeployProductListInfo{}

		err = rows.StructScan(&info)
		if err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}
	return list, total
}

func (d *deployProdcutList) GetProductInfoById(pid int) (*DeployProductListInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("id", pid)
	info := DeployProductListInfo{}

	return &info, d.GetWhere(nil, whereCause, &info)
}

var _getDeployProductHistoryFields = utils.GetTagValues(DeployProductHistoryInfo{}, "db")

func (h *deployProductHistory) GetDeployProductHistory(pagination *apibase.Pagination, parentProductName, productNames string, productType []int, deployStatus []string, productVersionLike, clusterId string) ([]DeployProductHistoryInfo, int) {
	whereCause := dbhelper.WhereCause{}
	var values []interface{}

	if parentProductName != "" {
		whereCause = whereCause.Equal("parent_product_name", parentProductName)
	}
	if productNames != "" {
		for _, v := range strings.Split(productNames, ",") {
			values = append(values, v)
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("product_name", values...)
	}
	if productVersionLike != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Like("product_version", "%"+productVersionLike+"%")
	}
	if clusterId != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("cluster_id", clusterId)
	}
	if len(deployStatus) > 0 {
		status := make([]interface{}, 0, len(deployStatus))
		for _, ds := range deployStatus {
			status = append(status, interface{}(ds))
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("status", status...)
	}
	if len(productType) > 0 {
		pType := make([]interface{}, 0, len(productType))
		for _, pt := range productType {
			pType = append(pType, interface{}(pt))
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("product_type", pType...)
	}

	rows, totalcounts, err := h.SelectWhere(_getDeployProductHistoryFields, whereCause, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}

	defer rows.Close()

	list := []DeployProductHistoryInfo{}
	for rows.Next() {
		info := DeployProductHistoryInfo{}
		if err := rows.StructScan(&info); err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}

	return list, totalcounts

}

var _getDeployProductUpdateHistoryFields = utils.GetTagValues(DeployProductUpdateHistoryInfo{}, "db")

func (h *deployProductUpdateHistory) GetDeployProductUpdateHistory(pagination *apibase.Pagination, parentProductName, productNames string, productType []int, updateStatus []string, productVersionLike, clusterId string) ([]DeployProductUpdateHistoryInfo, int) {
	whereCause := dbhelper.WhereCause{}
	var values []interface{}

	if parentProductName != "" {
		whereCause = whereCause.Equal("parent_product_name", parentProductName)
	}
	if productNames != "" {
		for _, v := range strings.Split(productNames, ",") {
			values = append(values, v)
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("product_name", values...)
	}
	if productVersionLike != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Like("product_version", "%"+productVersionLike+"%")
	}
	if clusterId != "" {
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("cluster_id", clusterId)
	}
	if len(updateStatus) > 0 {
		status := make([]interface{}, 0, len(updateStatus))
		for _, ds := range updateStatus {
			status = append(status, interface{}(ds))
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("status", status...)
	}
	if len(productType) > 0 {
		pType := make([]interface{}, 0, len(productType))
		for _, pt := range productType {
			pType = append(pType, interface{}(pt))
		}
		whereCause = whereCause.And()
		whereCause = whereCause.Included("product_type", pType...)
	}

	rows, totalcounts, err := h.SelectWhere(_getDeployProductUpdateHistoryFields, whereCause, pagination)
	if err != nil {
		apibase.ThrowDBModelError(err)
	}

	defer rows.Close()

	list := []DeployProductUpdateHistoryInfo{}
	for rows.Next() {
		info := DeployProductUpdateHistoryInfo{}
		if err := rows.StructScan(&info); err != nil {
			apibase.ThrowDBModelError(err)
		}
		list = append(list, info)
	}

	return list, totalcounts

}

func (d *deployProdcutList) GetByProductNameAndVersion(productName string, productVersion string) (*DeployProductListInfo, error) {
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("product_name", productName)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("product_version", productVersion)
	info := DeployProductListInfo{}
	return &info, d.GetWhere(nil, whereCause, &info)
}
func (d *deployProductHistory) GetDeployHistoryByDeployUUID(UUID string) (*DeployProductHistoryInfo, error) {
	var info DeployProductHistoryInfo
	sql := fmt.Sprintf("select * from %s where deploy_uuid = ?", d.TableName)
	err := d.GetDB().Get(&info, sql, UUID)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (d *deployProdcutList) GetDeploySonProductName(parentProductName string, clusterId int) ([]string, error) {
	list := []string{}
	query := fmt.Sprintf("SELECT product_name from %v LEFT JOIN %v ON %v.id = %v.pid WHERE parent_product_name=? AND clusterId=? AND %v.is_deleted=0",
		TBL_DEPLOY_PRODUCT_LIST,
		TBL_DEPLOY_CLUSTER_PRODUCT_REL,
		TBL_DEPLOY_PRODUCT_LIST,
		TBL_DEPLOY_CLUSTER_PRODUCT_REL,
		TBL_DEPLOY_CLUSTER_PRODUCT_REL)
	err := USE_MYSQL_DB().Select(&list, query, parentProductName, clusterId)
	return list, err
}

func (d *deployProdcutList) GetByProductNameClusterId(productName string, clusterId int) (*DeployProductListInfo, error) {
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.is_current_version, p.product, p.schema, p.product_type, p.create_time, " +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ? AND p.product_name = ? AND r.is_deleted=0"
	product := &DeployProductListInfo{}
	err := USE_MYSQL_DB().Get(product, query, clusterId, productName)
	return product, err
}

func (d *deployProdcutList) GetByProductNameAndParentProductNameClusterId(productName, parentProductName string, clusterId int) (*DeployProductListInfo, error) {
	query := "SELECT p.id, p.parent_product_name, p.product_name, p.product_name_display, p.product_version, p.is_current_version, p.product, p.schema, p.product_type, p.create_time," +
		"r.product_parsed, r.status, r.deploy_uuid, r.alert_recover, r.user_id, r.deploy_time " +
		"from deploy_cluster_product_rel as r " +
		"LEFT JOIN deploy_product_list as p ON p.id = r.pid " +
		"WHERE r.clusterId = ? AND p.product_name = ? AND p.parent_product_name = ? AND r.is_deleted=0"
	product := &DeployProductListInfo{}
	err := USE_MYSQL_DB().Get(product, query, clusterId, productName, parentProductName)
	return product, err
}

func (d *deploySchemaFieldModify) GetByProductNameClusterId(productName string, clusterId int) ([]SchemaFieldModifyInfo, error) {
	list := []SchemaFieldModifyInfo{}
	err := DeploySchemaFieldModify.GetDB().Select(&list, "select * from "+DeploySchemaFieldModify.TableName+
		" where product_name=? AND cluster_id=?", productName, clusterId)
	return list, err
}

func (d *deploySchemaFieldModify) GetByProductNameAndServiceAndClusterId(productName, serviceName string, clusterId int, namespace string) ([]SchemaFieldModifyInfo, error) {
	list := []SchemaFieldModifyInfo{}

	err := DeploySchemaFieldModify.GetDB().Select(&list, "select * from "+DeploySchemaFieldModify.TableName+
		" where product_name=? and service_name=? AND cluster_id=? AND namespace=?", productName, serviceName, clusterId, namespace)
	return list, err
}

func (d *deploySchemaFieldModify) GetServiceModifyTime(clusterId int, productName, serviceName string) (time.Time, error) {
	info := SchemaFieldModifyInfo{}
	query := "SELECT * FROM deploy_schema_field_modify where cluster_id = ? and product_name = ? and service_name = ? order by update_time desc limit 1"
	if err := d.GetDB().Get(&info, query, clusterId, productName, serviceName); err != nil {
		log.Errorf("[deploySchemaFieldModify.GetServiceModifyTime] %s", err)
		return time.Time{}, err
	}
	return info.UpdateDate, nil
}

func (d *deploySchemaFieldModify) GetFieldValue(clusterId int, productName, serviceName, fieldPath string) (string, error) {
	fieldValue := ""
	query := "SELECT field FROM deploy_schema_field_modify where cluster_id = ? and product_name = ? and service_name = ? and field_path = ?"
	if err := d.GetDB().Get(&fieldValue, query, clusterId, productName, serviceName, fieldPath); err != nil {
		return "", err
	}
	return fieldValue, nil
}

func (d *deploySchemaFieldModify) BatchInsert(modifyList []interface{}, tx *sqlx.Tx) error {
	if len(modifyList) > 0 {
		argsList := make([]string, 0)
		for range modifyList {
			argsList = append(argsList, "(?)")
		}
		query, args, err := sqlx.In(fmt.Sprintf("INSERT INTO %s (cluster_id, product_name, service_name, field_path, "+
			"field, update_time, create_time, namespace) VALUES %s", d.TableName, strings.Join(argsList, ",")), modifyList...)
		if err != nil {
			return err
		}
		_, err = tx.Exec(query, args...)
		return err
	}
	return nil
}

func (d *deployProdcutList) GetProductPidAndNameMap() (map[int]string, error) {
	list := make([]DeployProductListInfo, 0)
	err := DeployProductList.GetDB().Select(&list, "select * from "+DeployProductList.TableName)
	if err != nil {
		return nil, err
	}
	mp := make(map[int]string)
	for _, product := range list {
		mp[product.ID] = product.ProductName
	}
	return mp, nil

}

func (d *deployProdcutList) GetProductListInfo(
	parentProductName, productName, productVersionLike, productVersion, productType string, clusterId int, deployStatus, productNames []string, namespace string) ([]DeployProductListInfoWithNamespace, error) {

	query := "SELECT * from deploy_product_list " + "WHERE id > 0  "

	if productType != "" {
		pType, _ := strconv.Atoi(productType)
		query += fmt.Sprintf(" AND product_type= %d", pType)
	}

	if parentProductName != "" {
		query += " AND parent_product_name=" + "'" + parentProductName + "'"
	}

	if len(productNames) > 0 {
		query += " AND product_name IN ("
		for i, v := range productNames {
			if i > 0 {
				query += ",'" + v + "'"
			} else {
				query += "'" + v + "'"
			}
		}
		query += ")"
	}

	if productName != "" {
		query += " AND product_name=" + "'" + productName + "'"
	}

	if productVersionLike != "" {
		query += " AND product_version like" + "'%" + productVersionLike + "%'"
	}

	if productVersion != "" {
		query += " AND product_version=" + "'" + productVersion + "'"
	}

	list := make([]DeployProductListInfoWithNamespace, 0)
	err := USE_MYSQL_DB().Select(&list, query)
	if err != nil {
		return nil, fmt.Errorf("[ProductInfo] GetDeployProductLists query err: %s", err)
	}

	return list, nil
}
