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

package impl

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/matrix/instance"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/model/upgrade"
	"dtstack.com/dtstack/easymatrix/schema"
	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
)

var (
	BackupSqlFile = "/data/dtstack/backup/%s:%s:%s/%d/%s.sql"
	DatabaseTxt   = "database.txt"
	MysqlDb       = "mysql_db"
	MysqlBackupDb = "mysql_backup_db"
	MysqlDumpFile = "mysql_dump_absolute_path"
)

type BackupResp struct {
	Status     string `json:"status"`
	ExecId     string `json:"exec_id"`
	BackupName string `json:"backup_name"`
	BackupSqls string `json:"backup_sqls"`
}

func BackupDatabase(ctx context.Context) apibase.Result {
	var backupParam = struct {
		ClusterId     int    `json:"cluster_id"`
		SourceVersion string `json:"source_version"`
		TargetVersion string `json:"target_version"`
	}{}
	if err := ctx.ReadJSON(&backupParam); err != nil {
		log.Errorf("BackupDatabase read param error: %v", err)
		return err
	}
	productName := ctx.Params().Get("product_name")
	var productVersion = backupParam.TargetVersion
	var isRollback bool
	if CompareVersion(backupParam.SourceVersion, backupParam.TargetVersion) > 0 {
		productVersion = backupParam.SourceVersion
		isRollback = true
	}
	productInfo, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		log.Errorf("BackupDatabase-GetProduct productName: %s, productVersion: %s db error: %v", productName,
			backupParam.TargetVersion, err)
		return err
	}
	sc, err := schema.Unmarshal(productInfo.Schema)
	if err != nil {
		log.Errorf("BackupDatabase-UnmarshalSchema for productName: %s, productVersion: %s error: %v", productName,
			backupParam.TargetVersion, err)
		return err
	}
	tx := model.USE_MYSQL_DB().MustBegin()
	now := time.Now()
	var backupScript, svc, errExecId string
	var backupDbs, backupSqls []string
	for svcName, svcConfig := range sc.Service {
		if sqlValid(svcName) && svcConfig.Instance != nil && svcConfig.Instance.Backup != "" {
			svc = svcName
			if err = inheritBaseService(backupParam.ClusterId, sc, model.USE_MYSQL_DB()); err != nil {
				log.Errorf("BackupDatabase-InheritBaseService error: %v", err)
				return err
			}
			if err = setSchemaFieldServiceAddr(backupParam.ClusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
				log.Debugf("BackupDatabase-SetSchemaFieldServiceAddr err: %v", err)
				return err
			}
			if err = sc.ParseVariable(); err != nil {
				log.Errorf("BackupDatabase-ParseVariable error: %v", err)
				return err
			}
			if v, ok := svcConfig.Config[MysqlDb]; ok {
				backupDbs = append(backupDbs, v.(schema.VisualConfig).String())
			}
			if !isRollback {
				dependsDatabaseBytes, err := ioutil.ReadFile(filepath.Join(base.WebRoot, productName, productVersion,
					svcName, DatabaseTxt))
				if err != nil && !os.IsNotExist(err) {
					log.Errorf("BackupDatabase-ReadDependDatabase error: %v", err)
					return err
				}
				dependDatabases := strings.Split(string(dependsDatabaseBytes), ",")
				if len(dependDatabases) > 0 {
					backupDbs = append(backupDbs, dependDatabases...)
				}
			}
			if len(backupDbs) > 0 {
				deployInstances, _ := model.DeployInstanceList.GetInstanceBelongServiceWithNamespace(productName, svc, backupParam.ClusterId, "")
				if len(deployInstances) == 0 {
					log.Errorf("BackupDatabase-QueryInstance, cluster: %d, product: %s, service: %s,  error: %v",
						backupParam.ClusterId, productName, svc, fmt.Errorf("未部署SQL组件"))
					return fmt.Errorf("未部署SQL组件")
				}
				for _, dbName := range backupDbs {
					if dbName != "" && dbName != "\n" {
						backupSql := fmt.Sprintf(BackupSqlFile, productName, backupParam.SourceVersion, backupParam.TargetVersion,
							now.Unix(), dbName)
						svcConfig.Config[MysqlBackupDb] = schema.VisualConfig{
							Default: &dbName,
							Desc:    "internal",
							Type:    "internal",
							Value:   &dbName,
						}
						svcConfig.Config[MysqlDumpFile] = schema.VisualConfig{
							Default: &backupSql,
							Desc:    "internal",
							Type:    "internal",
							Value:   &backupSql,
						}
						backupFileName := filepath.Join(base.WebRoot, sc.ProductName, sc.ProductVersion, svcName, svcConfig.Instance.Backup)
						tpl, err := template.ParseFiles(backupFileName)
						if err != nil {
							log.Errorf("BackupDatabase-ParseTemplate error: %v", err)
							return err
						}
						buf := &bytes.Buffer{}
						if err = tpl.Option("missingkey=error").Execute(buf, svcConfig.Config); err != nil {
							log.Errorf("BackupDatabase-ExecuteTemplate error: %v", err)
							return err
						}
						backupScript = buf.String()
						latestBackupInfo, err := upgrade.BackupHistory.GetLatestRecord(backupParam.ClusterId, dbName)
						if err != nil && err != sql.ErrNoRows {
							log.Errorf("BackupDatabase-QueryBackupHistory db: %s, error: %v", dbName, err)
							continue
						}
						if isRollback {
							backupSqls = append(backupSqls, backupSql)
						} else {
							if latestBackupInfo != nil && time.Since(latestBackupInfo.CreateTime) <= 1*time.Hour {
								backupSqls = append(backupSqls, latestBackupInfo.BackupSql)
								continue
							} else {
								backupSqls = append(backupSqls, backupSql)
							}
						}
						//生成 operationid 并且落库
						operationId := uuid.NewV4().String()
						err = model.OperationList.Insert(model.OperationInfo{
							ClusterId:       backupParam.ClusterId,
							OperationId:     operationId,
							OperationType:   enums.OperationType.Backup.Code,
							OperationStatus: enums.ExecStatusType.Running.Code,
							ObjectType:      enums.OperationObjType.Host.Code,
							ObjectValue:     deployInstances[0].Ip,
						})
						execId := uuid.NewV4().String()
						err = model.ExecShellList.InsertExecShellInfo(backupParam.ClusterId, operationId, execId, productInfo.ProductName, svc, deployInstances[0].Sid, enums.ShellType.Exec.Code)
						if err != nil {
							log.Errorf("ExecShellList InsertExecShellInfo error: %v")
						}
						param := &agent.ExecScriptParams{
							ExecScript: backupScript,
							Parameter:  "",
							Timeout:    "1h",
						}
						err = agentExec(deployInstances[0].Sid, execId, param)
						if err != nil {
							log.Errorf("BackupDatabase-exec script error: %v", err)
							errExecId = execId
							break
						}
						_, err = upgrade.BackupHistory.InsertRecord(backupParam.ClusterId, dbName, backupSql, productName, tx)
						if err != nil {
							log.Errorf("BackupDatabase-InsertBackup dbName: %v, error: %v", dbName, err)
							errExecId = execId
							break
						}
					}
				}
			}
			break
		}
	}

	result := &BackupResp{}
	if errExecId != "" {
		tx.Rollback()
		result.Status = "fail"
		result.ExecId = errExecId
	} else {
		if err := tx.Commit(); err != nil {
			result.Status = "fail"
			result.ExecId = errExecId
		} else {
			result.Status = "success"
			result.BackupName = now.Format(base.TsLayout)
			result.BackupSqls = strings.Join(backupSqls, ",")
		}
	}
	return result
}

func agentExec(sid, execId string, param *agent.ExecScriptParams) error {
	err, agentServerResp := agent.AgentClient.AgentExec(sid, param, execId)
	if err != nil {
		log.Errorf("ExecScript sid: %v, err: %v, resp: %v", sid, err, agentServerResp)
		return err
	}
	respResult, _ := agentServerResp.Data.(map[string]interface{})["result"]
	failed, ok := respResult.(map[string]interface{})["failed"]
	if ok && failed.(bool) {
		message, ok := respResult.(map[string]interface{})["response"]
		if ok && message != nil {
			log.Errorf("ExecScript err: %v", message.(string))
			return fmt.Errorf(message.(string))
		}
		log.Errorf("ExecScript error unknow")
		return fmt.Errorf("unknown error")
	}
	return nil
}

func sqlValid(svcName string) bool {
	return strings.HasSuffix(svcName, "Sql") || strings.HasSuffix(svcName, "sql")
}

type serviceIpInfo struct {
	ServiceName string `json:"service_name"`
	IpList      string `json:"ip_list"`
}

type fieldModifyInfo struct {
	ServiceName string `json:"service_name"`
	FieldPath   string `json:"field_path"`
	Field       string `json:"field"`
}

type multiFieldModifyInfo struct {
	ServiceName string `json:"service_name"`
	FieldPath   string `json:"field_path"`
	Field       string `json:"field"`
	Hosts       string `json:"hosts"`
}

type CurrentProductResult struct {
	ServiceIpList       []serviceIpInfo        `json:"service_ip_list"`
	FieldModifyList     []fieldModifyInfo      `json:"field_modify_list"`
	FieldMultifieldList []multiFieldModifyInfo `json:"field_multifield_list"`
}

func ProductCurrentInfo(ctx context.Context) apibase.Result {
	paramErr := apibase.ApiParameterErrors{}
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		paramErr.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	paramErr.CheckAndThrowApiParameterErrors()
	param := struct {
		ClusterId int `json:"cluster_id"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("ProductCurrentInfo-parse param error: %v", err)
		return err
	}
	serviceIpList, err := model.DeployServiceIpList.GetServiceIpListByPNameAndClusterId(productName, param.ClusterId)
	if err != nil {
		log.Errorf("ProductCurrentInfo-query service ip error: %v", err)
		return fmt.Errorf("查询数据错误")
	}
	var result CurrentProductResult
	for _, serviceIp := range serviceIpList {
		result.ServiceIpList = append(result.ServiceIpList, serviceIpInfo{
			ServiceName: serviceIp.ServiceName,
			IpList:      serviceIp.IpList,
		})
	}
	fieldModifyList, err := model.DeploySchemaFieldModify.GetByProductNameClusterId(productName, param.ClusterId)
	if err != nil {
		log.Errorf("ProductCurrentInfo-query field modify error: %v", err)
		return fmt.Errorf("查询数据错误")
	}
	for _, fieldModify := range fieldModifyList {
		result.FieldModifyList = append(result.FieldModifyList, fieldModifyInfo{
			ServiceName: fieldModify.ServiceName,
			FieldPath:   fieldModify.FieldPath,
			Field:       fieldModify.Field,
		})
	}
	multiFieldModifyList, err := model.SchemaMultiField.GetListByProductName(param.ClusterId, productName)
	if err != nil {
		log.Errorf("ProductCurrentInfo-query multi field modify error: %v", err)
		return fmt.Errorf("查询数据错误")
	}
	for _, multiFiled := range multiFieldModifyList {
		result.FieldMultifieldList = append(result.FieldMultifieldList, multiFieldModifyInfo{
			ServiceName: multiFiled.ServiceName,
			FieldPath:   multiFiled.FieldPath,
			Field:       multiFiled.Field,
			Hosts:       multiFiled.Hosts,
		})
	}
	return result
}

type SaveUpgradeResult struct {
	UpgradeId  int    `json:"upgrade_id"`
	CreateTime string `json:"create_time"`
}

func SaveUpgrade(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		log.Errorf("SaveUpgrade-product_name is empty")
		return fmt.Errorf("product_name为空")
	}
	param := struct {
		ClusterId       int                    `json:"cluster_id"`
		SourceVersion   string                 `json:"source_version"`
		TargetVersion   string                 `json:"target_version"`
		BackupName      string                 `json:"backup_name"`
		BackupSqls      string                 `json:"backup_sqls"`
		UpgradeMode     string                 `json:"upgrade_mode"`
		ServiceIpList   []serviceIpInfo        `json:"service_ip_list"`
		FieldModifyList []fieldModifyInfo      `json:"field_modify_list"`
		FieldMultiList  []multiFieldModifyInfo `json:"field_multifield_list"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("SaveUpgrade-parse param error: %v", err)
		return err
	}
	result := SaveUpgradeResult{}
	serviceIpBytes, _ := json.Marshal(param.ServiceIpList)
	fieldModifyBytes, _ := json.Marshal(param.FieldModifyList)
	multiFieldBytes, _ := json.Marshal(param.FieldMultiList)
	id, err := upgrade.UpgradeHistory.InsertRecord(param.ClusterId, enums.UpgradeType.Upgrade.Code, param.UpgradeMode, productName, param.SourceVersion,
		param.TargetVersion, param.BackupName, param.BackupSqls, serviceIpBytes, fieldModifyBytes, multiFieldBytes)
	if err != nil {
		log.Errorf("SaveUpgrade-InsertRecord error: %v", err)
		return err
	}
	result.UpgradeId = int(id)
	result.CreateTime = time.Now().Format(base.TsLayout)
	return result
}

func RollbackVersions(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		log.Errorf("RollbackVersions-product_name is empty")
		return fmt.Errorf("product_name为空")
	}
	param := struct {
		ClusterId      int    `json:"cluster_id"`
		ProductVersion string `json:"product_version"`
		UpgradeMode    string `json:"upgrade_mode"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("RollbackVersions-parse param error: %v", err)
		return err
	}
	upgradeInfoList, err := upgrade.UpgradeHistory.GetByClsAndProductNameAndSourceVersion(param.ClusterId, productName, "", param.UpgradeMode)
	if err != nil {
		log.Errorf("RollbackVersions-query db error: %v", err)
		return err
	}
	var versionList []string
	for _, upgradeInfo := range upgradeInfoList {
		// 判断回滚的目标版本产品包是否被删除，如果被删除则不显示在回滚版本列表中
		_, err := model.DeployProductList.GetByProductNameAndVersion(productName, upgradeInfo.SourceVersion)
		if err != nil && err == sql.ErrNoRows {
			continue
		}
		if CompareVersion(upgradeInfo.SourceVersion, param.ProductVersion) < 0 &&
			!contains(versionList, upgradeInfo.SourceVersion) && upgradeInfo.SourceVersion != param.ProductVersion {
			versionList = append(versionList, upgradeInfo.SourceVersion)
		}
	}
	return versionList
}

func CanRollback(clusterId int, productName, productVersion string) bool {
	_, err := model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("query db error: %v", err)
		return false
	}
	if err == nil {
		return false
	}
	upgradeInfoList, err := upgrade.UpgradeHistory.GetByClsAndProductNameAndSourceVersion(clusterId, productName, "", "")
	if err != nil {
		log.Errorf("query db error: %v", err)
		return false
	}
	for _, upgradeInfo := range upgradeInfoList {
		// 判断回滚的目标版本产品包是否被删除，如果被删除则不能回滚
		_, err := model.DeployProductList.GetByProductNameAndVersion(productName, upgradeInfo.SourceVersion)
		if err != nil && err == sql.ErrNoRows {
			continue
		}
		if CompareVersion(upgradeInfo.SourceVersion, productVersion) < 0 && upgradeInfo.SourceVersion != productVersion {
			return true
		}
	}
	return false
}

func BackupTimes(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		log.Errorf("BackupTimes-product_name is empty")
		return fmt.Errorf("product_name为空")
	}
	param := struct {
		ClusterId     int    `json:"cluster_id"`
		TargetVersion string `json:"target_version"`
		UpgradeMode   string `json:"upgrade_mode"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("BackupTimes-parse param error: %v", err)
		return err
	}
	upgradeInfoList, err := upgrade.UpgradeHistory.GetByClsAndProductNameAndSourceVersion(param.ClusterId, productName,
		param.TargetVersion, param.UpgradeMode)
	if err != nil {
		log.Errorf("BackupTimes-query db error: %v", err)
		return err
	}
	var timeList []string
	for _, upgradeInfo := range upgradeInfoList {
		if upgradeInfo.BackupName != "" {
			timeList = append(timeList, upgradeInfo.BackupName)
		}
	}
	sort.Slice(timeList, func(i, j int) bool {
		timeI, err := time.Parse(base.TsLayout, timeList[i])
		if err != nil {
			return false
		}
		timeJ, err := time.Parse(base.TsLayout, timeList[j])
		if err != nil {
			return false
		}
		return timeI.After(timeJ)
	})
	return timeList
}

type RollbackResult struct {
	DeployUUID string `json:"deploy_uuid"`
}

func Rollback(ctx context.Context) (rlt apibase.Result) {
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		log.Errorf("Rollback-product_name is empty")
		return fmt.Errorf("product_name为空")
	}
	param := struct {
		ClusterId     int    `json:"cluster_id"`
		SourceVersion string `json:"source_version"`
		TargetVersion string `json:"target_version"`
		BackupName    string `json:"backup_name"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("Rollback-parse param error: %v", err)
		return err
	}

	upgradeHistoryInfo, err := upgrade.UpgradeHistory.GetOne(param.ClusterId, productName, param.TargetVersion, param.BackupName)
	if err != nil {
		log.Errorf("Rollback-query upgrade history info error: %v", err)
		return fmt.Errorf("备份库缺失")
	}
	var serviceIpList []serviceIpInfo
	if err := json.Unmarshal(upgradeHistoryInfo.SourceServiceIp, &serviceIpList); err != nil {
		log.Errorf("Rollback-unmarshal service ip list error: %v", err)
		return err
	}
	var fieldModifyList []fieldModifyInfo
	if err := json.Unmarshal(upgradeHistoryInfo.SourceConfig, &fieldModifyList); err != nil {
		log.Errorf("Rollback-unmarshal field modify list error: %v", err)
		return err
	}
	var fieldMultiList []multiFieldModifyInfo
	if err := json.Unmarshal(upgradeHistoryInfo.SourceMultiConfig, &fieldMultiList); err != nil {
		log.Errorf("Rollback-unmarshal multi field list error: %v", err)
		return err
	}
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			err := tx.Rollback()
			if err != nil {
				log.Errorf("Transactional rollback error: %v", err)
			} else {
				log.Infof("Transactional rollback success.")
			}
		}
		if r := recover(); r != nil {
			err := tx.Rollback()
			if err != nil {
				log.Errorf("Transaction rollback error: %v", err)
			} else {
				log.Infof("Transaction rollback success.")
			}
			rlt = r
		}
	}()
	// 回滚deploy_service_ip_list
	_, err = tx.Exec(fmt.Sprintf("delete from %s where cluster_id=? and product_name=?", model.DeployServiceIpList.TableName),
		param.ClusterId, productName)
	if err != nil {
		log.Errorf("delete deploy service ip list error: %v", err)
		return err
	}
	var newServiceIpList, newFieldModifyList, newMultiFieldList []interface{}
	for _, serviceIpInfo := range serviceIpList {
		newServiceIpList = append(newServiceIpList, model.DeployServiceIpInfo{
			ClusterId:   param.ClusterId,
			NameSpace:   "",
			ProductName: productName,
			ServiceName: serviceIpInfo.ServiceName,
			IpList:      serviceIpInfo.IpList,
			UpdateDate:  dbhelper.NullTime{Time: time.Now()},
			CreateDate:  dbhelper.NullTime{Time: time.Now()},
		})
	}
	err = model.DeployServiceIpList.BatchInsertServiceIp(newServiceIpList, tx)
	if err != nil {
		log.Errorf("batch insert service ip list error: %v", err)
		return err
	}
	// 回滚deploy_schema_field_modify
	_, err = tx.Exec(fmt.Sprintf("delete from %s where cluster_id=? and product_name=?", model.DeploySchemaFieldModify.TableName),
		param.ClusterId, productName)
	if err != nil {
		log.Errorf("delete schema field modify error: %v", err)
		return err
	}
	for _, fieldModify := range fieldModifyList {
		newFieldModifyList = append(newFieldModifyList, model.SchemaFieldModifyInfo{
			ClutserId:   param.ClusterId,
			ProductName: productName,
			ServiceName: fieldModify.ServiceName,
			FieldPath:   fieldModify.FieldPath,
			Field:       fieldModify.Field,
			UpdateDate:  time.Now(),
			CreateDate:  time.Now(),
			Namespace:   "",
		})
	}
	err = model.DeploySchemaFieldModify.BatchInsert(newFieldModifyList, tx)
	if err != nil {
		log.Errorf("batch insert schema field modify error: %v", err)
		return err
	}
	// 回滚deploy_schema_multi_field
	_, err = tx.Exec(fmt.Sprintf("delete from %s where cluster_id=? and product_name=?", model.SchemaMultiField.TableName), param.ClusterId, productName)
	if err != nil {
		log.Errorf("delete schema multi field error: %v", err)
		return err
	}
	for _, multiField := range fieldMultiList {
		newMultiFieldList = append(newMultiFieldList, model.SchemaMultiFieldInfo{
			ClusterId:   param.ClusterId,
			ProductName: productName,
			ServiceName: multiField.ServiceName,
			FieldPath:   multiField.FieldPath,
			Field:       multiField.Field,
			Hosts:       multiField.Hosts,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
			IsDeleted:   0,
		})
	}
	err = model.SchemaMultiField.BatchInsert(newMultiFieldList, tx)
	if err != nil {
		log.Errorf("batch insert schema multi field error: %v", err)
		return err
	}
	sourceProductInfo, err := model.DeployProductList.GetByProductNameAndVersion(productName, param.SourceVersion)
	if err != nil {
		log.Errorf("Rollback-GetProduct productName: %s, productVersion: %s db error: %v", productName,
			param.TargetVersion, err)
		return err
	}
	sc, err := schema.Unmarshal(sourceProductInfo.Schema)
	if err != nil {
		log.Errorf("Rollback-unmarshal product schema error: %v", err)
		return err
	}
	err = execRollback(productName, param.ClusterId, sc, upgradeHistoryInfo)
	if err != nil {
		log.Errorf("Rollback-exec rollback error: %v", err)
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Errorf("transaction commit error: %v", err)
		return err
	}
	uncheckedServiceInfo, err := model.DeployUncheckedService.GetUncheckedServicesByPidClusterId(sourceProductInfo.ID,
		param.ClusterId, "")
	if err != nil {
		log.Errorf("Rollback-GetUncheckedService error: %v", err)
		return err
	}
	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		return fmt.Errorf("get userId err: %v", err)
	}
	var unchecked []string
	if uncheckedServiceInfo.UncheckedServices != "" {
		unchecked = strings.Split(uncheckedServiceInfo.UncheckedServices, ",")
	}
	deployUUID := DealDeploy(productName, param.TargetVersion, param.SourceVersion, unchecked, userId, param.ClusterId, 2, false)
	return deployUUID
}

func execRollback(productName string, clusterId int, sc *schema.SchemaConfig, upgradeHistoryInfo *upgrade.HistoryInfo) error {
	var err error
	for svcName, svcConfig := range sc.Service {
		if sqlValid(svcName) && svcConfig.Instance != nil && svcConfig.Instance.Rollback != "" {
			if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
				log.Errorf("Rollback-InheritBaseService error: %v", err)
				return err
			}
			if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
				log.Debugf("Rollback-SetSchemaFieldServiceAddr err: %v", err)
				return err
			}
			if err = sc.ParseVariable(); err != nil {
				log.Errorf("Rollback-ParseVariable error: %v", err)
				return err
			}
			var dbName, sqlAbsolutePath string
			if v, ok := svcConfig.Config[MysqlDb]; ok {
				dbName = v.(schema.VisualConfig).String()
			}
			if upgradeHistoryInfo != nil && upgradeHistoryInfo.BackupSql != "" {
				backupSqlList := strings.Split(upgradeHistoryInfo.BackupSql, ",")
				for _, backupSql := range backupSqlList {
					if strings.Contains(backupSql, dbName) {
						sqlAbsolutePath = backupSql
						break
					}
				}
			}
			deployInstances, _ := model.DeployInstanceList.GetInstanceBelongServiceWithNamespace(productName, svcName,
				clusterId, "")
			if len(deployInstances) == 0 {
				log.Errorf("Rollback-QueryInstance, cluster: %d, product: %s, service: %s,  error: %v",
					clusterId, productName, svcName, fmt.Errorf("未部署SQL组件"))
				return fmt.Errorf("未部署SQL组件")
			}
			err = checkBackupFileExist(sqlAbsolutePath, deployInstances[0].Sid)
			if err != nil {
				log.Errorf("Rollback-checkBackupFile, file not exist: %s", sqlAbsolutePath)
				return fmt.Errorf("备份库缺失")
			}
			if dbName != "" {
				svcConfig.Config[MysqlDumpFile] = schema.VisualConfig{
					Default: &sqlAbsolutePath,
					Desc:    "internal",
					Type:    "internal",
					Value:   &sqlAbsolutePath,
				}
				rollbackFileName := filepath.Join(base.WebRoot, sc.ProductName, sc.ProductVersion, svcName, svcConfig.Instance.Rollback)
				tpl, err := template.ParseFiles(rollbackFileName)
				if err != nil {
					log.Errorf("Rollback-ParseTemplate error: %v", err)
					return err
				}
				buf := &bytes.Buffer{}
				if err = tpl.Option("missingkey=error").Execute(buf, svcConfig.Config); err != nil {
					log.Errorf("Rollback-ExecuteTemplate error: %v", err)
					return err
				}
				//生成 operationid 并且落库
				operationId := uuid.NewV4().String()
				err = model.OperationList.Insert(model.OperationInfo{
					ClusterId:       clusterId,
					OperationId:     operationId,
					OperationType:   enums.OperationType.Rollback.Code,
					OperationStatus: enums.ExecStatusType.Running.Code,
					ObjectType:      enums.OperationObjType.Host.Code,
					ObjectValue:     deployInstances[0].Ip,
				})
				execId := uuid.NewV4().String()
				err = model.ExecShellList.InsertExecShellInfo(clusterId, operationId, execId, productName, svcName, deployInstances[0].Sid, enums.ShellType.Exec.Code)
				if err != nil {
					log.Errorf("ExecShellList InsertExecShellInfo error: %v", err)
				}
				param := &agent.ExecScriptParams{
					ExecScript: buf.String(),
					Timeout:    "1h",
				}
				err = agentExec(deployInstances[0].Sid, execId, param)
				if err != nil {
					log.Errorf("Rollback-exec script error: %v", err)
					return err
				}
			}
			break
		}
	}
	return nil
}

func checkBackupFileExist(sqlFilePath, sid string) error {
	param := &agent.ExecScriptParams{
		ExecScript: fmt.Sprintf("%s ls -dp %s", instance.LINUX_EXEC_HEADER, sqlFilePath),
		Timeout:    "5m",
	}
	err := agentExec(sid, "", param)
	if err != nil {
		return err
	}
	return nil
}
