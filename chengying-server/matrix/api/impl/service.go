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

//黑夜给了我黑色的眼睛，专治各种buf(g)；
package impl

import (
	"database/sql"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/schema"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"

	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"

	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/service"
	"github.com/kataras/iris/context"
)

func ServiceStart(ctx context.Context) apibase.Result {
	log.Debugf("[ServiceStart] ServiceStart from EasyMatrix API ")

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	serviceName := ctx.Params().Get("service_name")

	if serviceName == "" {
		log.Errorf("%v", "service_name is empty")
		return fmt.Errorf("service_name is empty")
	}
	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	err = model.OperationList.Insert(model.OperationInfo{
		ClusterId:       clusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.SvcStart.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Svc.Code,
		ObjectValue:     serviceName,
	})
	if err != nil {
		log.Errorf("OperationList Insert err:%v", err)
	}
	servicer, err := service.NewServicer(pid, clusterId, serviceName, operationId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		productInfo, err := model.DeployProductList.GetProductInfoById(pid)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "集群运维", "服务启动", "集群名称："+clusterInfo.Name+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion+
			", 服务名称："+serviceName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	return servicer.Start()
}

func ServiceStop(ctx context.Context) apibase.Result {
	log.Debugf("[ServiceStop] ServiceStop from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	serviceName := ctx.Params().Get("service_name")

	if serviceName == "" {
		log.Errorf("%v", "service_name is empty")
		return fmt.Errorf("service_name is empty")
	}

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	servicer, err := service.NewServicer(pid, clusterId, serviceName, "")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		productInfo, err := model.DeployProductList.GetProductInfoById(pid)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		_, deployInstanceInfo := model.DeployInstanceList.GetInstanceInfoByWhere(dbhelper.MakeWhereCause().Equal("pid", pid).And().Equal("cluster_id", clusterId).And().Equal("service_name", serviceName))
		if err := addSafetyAuditRecord(ctx, "集群运维", "服务停止", "集群名称："+clusterInfo.Name+", 服务组："+deployInstanceInfo.Group+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion+
			", 服务名称："+serviceName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	return servicer.Stop()
}

func ServiceRollingRestart(ctx context.Context) apibase.Result {
	log.Debugf("[ServiceRollingRestart] ServiceRollingRestart from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	serviceName := ctx.Params().Get("service_name")

	if serviceName == "" {
		log.Errorf("%v", "service_name is empty")
		return fmt.Errorf("service_name is empty")
	}

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	err = model.OperationList.Insert(model.OperationInfo{
		ClusterId:       clusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.SvcRollingRestart.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Svc.Code,
		ObjectValue:     serviceName,
	})
	if err != nil {
		log.Errorf("OperationList Insert err:%v", err)
	}
	servicer, err := service.NewServicer(pid, clusterId, serviceName, operationId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	productInfo, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		log.Errorf("%v\n", err)
	}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		_, deployInstanceInfo := model.DeployInstanceList.GetInstanceInfoByWhere(dbhelper.MakeWhereCause().Equal("pid", pid).And().Equal("cluster_id", clusterId).And().Equal("service_name", serviceName))
		if err := addSafetyAuditRecord(ctx, "集群运维", "服务滚动重启", "集群名称："+clusterInfo.Name+", 服务组："+deployInstanceInfo.Group+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion+
			", 服务名称："+serviceName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	// 查看待重启服务列表中是否有该服务，如果存在，滚动重启后清除记录
	if err := model.NotifyEvent.DeleteNotifyEvent(clusterId, 0, productInfo.ProductName, serviceName, "", false); err != nil {
		log.Errorf("delete notify event error: %v", err)
	}
	return servicer.RollingRestart()
}

func fieldHasChange(clusterId int, productName, serviceName string) bool {
	var (
		lastSendTime, err1 = model.NotifyEvent.GetServiceLastStartTime(clusterId, productName, serviceName)
		modifyTime, err2   = model.DeploySchemaFieldModify.GetServiceModifyTime(clusterId, productName, serviceName)
		multiTime, err3    = model.SchemaMultiField.GetServiceModifyTime(clusterId, productName, serviceName)
	)
	// 情况1: lastSendTime 没数据，配置没下发过
	// modifyTime 没数据，没修改过 && multiTime 没数据，没修改过 --->不需要重启
	// modifyTime 有数据，需要重启
	// multiTime  有数据，需要重启

	// 情况2: lastSendTime 有数据，配置下发过
	// modifyTime 没数据，没修改过 && multiTime 没数据，没修改过 --->不需要重启
	// modifyTime 无数据 && multiTime 有数据，时间早 ---> 不需要重启
	// modifyTime 有数据，时间早 && multiTime 无数据 ---> 不需要重启
	// modifyTime 有数据，时间早 && multiTime 有数据，时间早 ---> 不需要重启
	if err1 == sql.ErrNoRows {
		if err2 == sql.ErrNoRows && err3 == sql.ErrNoRows {
			log.Debugf("[skip add restart service]: db no row to restart")
			return false
		}
	} else {
		if err2 == sql.ErrNoRows && err3 == sql.ErrNoRows {
			log.Debugf("[skip add restart service]: field no row to restart")
			return false
		}
		if err2 != sql.ErrNoRows && err3 == sql.ErrNoRows && modifyTime.Before(*lastSendTime) {
			log.Debugf("[skip add restart service]: multiTime time old")
			return false
		}
		if err2 == sql.ErrNoRows && err3 != sql.ErrNoRows && multiTime.Before(*lastSendTime) {
			log.Debugf("[skip add restart service]: modifyTime time old")
			return false
		}
		if err2 != sql.ErrNoRows && err3 != sql.ErrNoRows && multiTime.Before(*lastSendTime) && modifyTime.Before(*lastSendTime) {
			log.Debugf("[skip add restart service]: multiTime and modifyTime time old")
			return false
		}
	}
	return true
}

func addRestartService(clusterId, pid int, serviceName string, productInfo *model.DeployProductListInfo) {
	if !fieldHasChange(clusterId, productInfo.ProductName, serviceName) {
		log.Debugf("[addRestartService.fieldHasChange]")
		return
	}

	// 获取当前集群已经部署的所有产品信息
	clusterProductRelList, err := model.DeployClusterProductRel.GetListByClusterIdAndStatus(clusterId, []string{model.PRODUCT_STATUS_DEPLOYED}, "")
	if err != nil {
		log.Errorf("get deployed products error: %v", err)
	}
	type productServiceInfo struct {
		ProductId     int
		ProductName   string
		ServiceName   string
		ServiceConfig schema.ServiceConfig
	}
	var productServiceList []productServiceInfo
	for _, clusterProductRel := range clusterProductRelList {
		sc, err := schema.Unmarshal(clusterProductRel.ProductParsed)
		if err != nil {
			log.Errorf("%v", err)
		}
		for svc, config := range sc.Service {
			// 重启服务需要添加当前服务
			if clusterProductRel.ClusterId == clusterId && clusterProductRel.Pid == pid && svc == serviceName {
				productServiceList = append(productServiceList, productServiceInfo{
					ProductName:   productInfo.ProductName,
					ServiceName:   serviceName,
					ProductId:     pid,
					ServiceConfig: config,
				})
			}
			// 添加依赖服务
			if len(config.RelyOn) > 0 {
				for _, relyOn := range config.RelyOn {
					infoSlice := strings.Split(relyOn, ".")
					if len(infoSlice) == 2 && productInfo.ProductName == infoSlice[0] && serviceName == infoSlice[1] {
						productServiceList = append(productServiceList, productServiceInfo{
							ProductId:     clusterProductRel.Pid,
							ProductName:   sc.ProductName,
							ServiceName:   svc,
							ServiceConfig: config,
						})
					}
				}
			}
		}
	}
	for _, productServiceInfo := range productServiceList {
		var ipList []string
		query := "select distinct(dh.ip) from deploy_host dh left join deploy_instance_list il on il.sid=dh.sid " +
			"where il.pid=? and il.service_name=? and il.status='running' and il.cluster_id=?"
		if err := model.USE_MYSQL_DB().Select(&ipList, query, productServiceInfo.ProductId, productServiceInfo.ServiceName, clusterId); err != nil {
			log.Errorf("%v", err)
			continue
		}
		if productServiceInfo.ServiceConfig.Instance != nil && !productServiceInfo.ServiceConfig.Instance.Pseudo {
			for _, ip := range ipList {
				err := model.NotifyEvent.InsertNotifyEvent(clusterId, 0, productServiceInfo.ProductName, productServiceInfo.ServiceName,
					productInfo.ProductName, serviceName, ip)
				if err != nil {
					log.Errorf("%v", err)
				}
			}
		}
	}
}

func ServiceRollingConfigUpdate(ctx context.Context) apibase.Result {
	log.Debugf("[ServiceConfigUpdate] ServiceConfogUpdate from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	serviceName := ctx.Params().Get("service_name")

	if serviceName == "" {
		log.Errorf("%v", "service_name is empty")
		return fmt.Errorf("service_name is empty")
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	servicer, err := service.NewServicer(pid, clusterId, serviceName, "")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		productInfo, err := model.DeployProductList.GetProductInfoById(pid)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		instance, _ := model.DeployInstanceList.GetInstanceListByClusterId(clusterId, pid)
		if err := addSafetyAuditRecord(ctx, "集群运维", "配置下发", "集群名称："+clusterInfo.Name+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion+
			", 服务组："+instance[0].Group+", 服务名称："+serviceName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
		addRestartService(clusterId, pid, serviceName, productInfo)
	}()
	return servicer.RollingConfigUpdate(clusterId)
}

type serviceLicenseConfig struct {
	NodeNum   int    `json:"node_num"`
	Exclusive string `json:"exclusive,omitempty"` // 以","隔离，服务的排他性，不允许和exclusive中的服务一起部署
}

type ProductLicenseConfig struct {
	Seq      int                             `json:"seq"` //部署顺序
	Services map[string]serviceLicenseConfig `json:"services"`
	Version  string
	DoDeploy bool
}

/**
请求格式：
{
	"dtinsight": {
		"dtapp": {
			"seq": 1,
			"services": {
				"api": {
					"node_num": 2,
					"exclusive": "hdfs_namenode,yarn_resourcemanager"
				},
				"console": {
					"node_num": 2
				},
				"gateway": {
					"node_num": 2
				}
			}
		},
		"hadoop": {
			"seq": 2,
			"services": {
				"hdfs_datanode": {
					"node_num": 3,
					"exclusive": "hdfs_namenode"
				},
				"yarn_nodemanager": {
					"node_num": 3,
					"exclusive": "yarn_resourcemanager"
				}
			}
		}
	}
}
*/
func License(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("[License] License from EasyMatrix API ")
	paramErrs := apibase.NewApiParameterErrors()
	license := map[string]map[string]ProductLicenseConfig{}
	if err := ctx.ReadJSON(&license); err != nil {
		paramErrs.AppendError("$", fmt.Errorf("params is error"))
		log.Debugf("[License] err:%v", err.Error())
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	log.Debugf("[License] params:%v", license)

	// ----- license 自动化部署实现 -----
	//tx := model.USE_MYSQL_DB().MustBegin()
	//defer func() {
	//	if _, ok := rlt.(error); ok {
	//		tx.Rollback()
	//		log.Errorf("[License] error(rollback): %v", rlt)
	//	}
	//	if r := recover(); r != nil {
	//		tx.Rollback()
	//		log.Errorf("[License] panic(rollback): %v", r)
	//		rlt = r
	//	}
	//}()

	//err := autoDeployProduct(license)
	err = setMaxReplica(license, clusterId)
	if err != nil {
		return err
	}

	//for _, products := range license {
	//	for deployProductName := range products {
	//		seq := math.MaxInt32
	//		for productName, productLicense := range products {
	//			if productLicense.Seq < seq {
	//				seq = productLicense.Seq
	//				deployProductName = productName
	//			}
	//		}
	//		res := DealDeploy(deployProductName, products[deployProductName].Version, -1)
	//		if err, ok := res.(error); ok {
	//			return err
	//		}
	//		delete(products, deployProductName)
	//	}
	//}
	//
	//if err := tx.Commit(); err != nil {
	//	return err
	//}

	return nil
}

func setMaxReplica(license map[string]map[string]ProductLicenseConfig, clusterId int) (err error) {
	for _, products := range license {
		for productName, productLicense := range products {
			for serviceName, serviceLicense := range productLicense.Services {
				modifyField := &model.SchemaFieldModifyInfo{
					ClutserId:   clusterId,
					FieldPath:   "instance.max_replica",
					Field:       strconv.Itoa(serviceLicense.NodeNum),
					ProductName: productName,
					ServiceName: serviceName,
				}
				if err, ok := CommonModifySchemaField(modifyField).(error); ok {
					return err
				}
			}
		}
	}
	return
}

//func autoDeployProduct(license map[string]map[string]ProductLicenseConfig, tx *sqlx.Tx, clusterId int) (err error) {
//	for parentProductName, products := range license {
//		for productName, productLicense := range products {
//			for serviceName, serviceLicense := range productLicense.Services {
//				err, ipInfo := model.DeployServiceIpList.GetServiceIpListByName(productName, serviceName, clusterId)
//				if err != nil {
//					return err
//				}
//				ipList := strings.Split(ipInfo.IpList, IP_LIST_SEP)
//				if len(ipList) == serviceLicense.NodeNum {
//					continue
//				}
//				modifyField := &model.SchemaFieldModifyInfo{
//					FieldPath:   "instance.max_replica",
//					Field:       strconv.Itoa(serviceLicense.NodeNum),
//					ProductName: productName,
//					ServiceName: serviceName,
//				}
//				if err, ok := CommonModifySchemaField(modifyField).(error); ok {
//					return err
//				}
//
//				ipListNews, err := getLicenseIpList(ipList, serviceLicense.NodeNum, serviceLicense.Exclusive)
//				if err != nil {
//					return err
//				}
//				updateSql := "UPDATE " + model.DeployServiceIpList.TableName + " set ip_list= ?, update_time=now() where product_name=? and service_name=?"
//				if _, err := tx.Exec(updateSql, strings.Join(ipListNews, IP_LIST_SEP), productName, serviceName); err != nil {
//					log.Errorf("[SetIP] SetServiceIp err: %v", err)
//					return err
//				}
//				err, productInfo := model.DeployProductList.GetCurrentProductInfoByName(productName)
//				if err != nil {
//					return err
//				}
//				productLicense.Version = productInfo.ProductVersion
//				productLicense.DoDeploy = true
//				license[parentProductName][productName] = productLicense
//			}
//			if !productLicense.DoDeploy {
//				delete(products, productName)
//			}
//		}
//	}
//	return
//}

func getLicenseIpList(ipList []string, nodeNum int, exclusive string) ([]string, error) {
	if nodeNum < len(ipList) {
		return ipList[:nodeNum], nil
	} else {
		var ipListNews []string

		var values []interface{}
		var whereCause string
		for i, ip := range ipList {
			if i > 0 {
				whereCause += ","
			}
			whereCause += "?"
			values = append(values, ip)
		}
		if exclusive == "" {
			if err := model.USE_MYSQL_DB().Select(&ipListNews, fmt.Sprintf("select ip from deploy_host where steps = 3 and ip NOT IN (%s)", whereCause), values...); err != nil {
				apibase.ThrowDBModelError(err)
			}
		} else {
			var exclusWhereCause string
			exclus := strings.Split(exclusive, IP_LIST_SEP)
			for i, exclu := range exclus {
				if i > 0 {
					exclusWhereCause += ","
				}
				exclusWhereCause += "?"
				values = append(values, exclu)
			}
			if err := model.USE_MYSQL_DB().Select(&ipListNews, fmt.Sprintf("select ip from deploy_host where steps = 3 and ip NOT IN (%s) and ip not in (select distinct ip from deploy_instance_list where service_name IN (%s))", whereCause, exclusWhereCause), values...); err != nil {
				apibase.ThrowDBModelError(err)
			}
		}
		increaseNode := nodeNum - len(ipList)
		if len(ipListNews) < increaseNode {
			return ipList, errors.New("node_num overflow limit node")
		}
		ipList = append(ipList, ipListNews[:increaseNode]...)
		return ipList, nil
	}
}
