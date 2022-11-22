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

package grafana

import (
	"fmt"
	"strconv"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
)

var (
	exipreData = map[int]int64{}
)

func StopAlert(pid int) error {
	info, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		return err
	}

	params := map[string]string{"dashboardTag": info.ProductName}
	err, alerts := GrafanaAlertsSearch(params)
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		err, pausedResp := GrafanaAlertsPause(strconv.Itoa(alert.Id), true)
		if err != nil {
			log.Errorf("%v, resp:%v", err, pausedResp)
		}
		fmt.Println(pausedResp)
	}

	db := model.USE_MYSQL_DB()
	query := "UPDATE " + model.DeployProductList.TableName + " SET `alert_recover`=1  WHERE `id`=?"
	if _, err = db.Exec(query, pid); err != nil {
		return err
	}
	return nil
}

func StartAlert(pid int) error {
	info, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		return err
	}

	params := map[string]string{"dashboardTag": info.ProductName}
	err, alerts := GrafanaAlertsSearch(params)
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		err, pausedResp := GrafanaAlertsPause(strconv.Itoa(alert.Id), false)
		if err != nil {
			log.Errorf("%v, resp:%v", err, pausedResp)
		}
	}

	db := model.USE_MYSQL_DB()
	query := "UPDATE " + model.DeployProductList.TableName + " SET `alert_recover`=0  WHERE `id`=?"
	if _, err = db.Exec(query, pid); err != nil {
		return err
	}

	delete(exipreData, pid)

	return nil
}

func StartMonitorAlertRecover() error {
	for {
		//获取符合要求的实例
		products, err := model.DeployProductList.GetProductListByWhere(dbhelper.MakeWhereCause().Equal("alert_recover", 1).And().Equal("status", model.PRODUCT_STATUS_DEPLOYED))
		if err != nil {
			return err
		}
		if len(products) > 0 {
			for _, product := range products {
				if expire(product.ID) {
					StartAlert(product.ID)
					continue
				}
				instances, err := model.DeployInstanceList.GetInstanceListByWhere(
					dbhelper.MakeWhereCause().Equal("pid", product.ID).And().NotEqual("status", model.INSTANCE_STATUS_RUNNING).
						Or().
						Equal("pid", product.ID).And().NotEqual("health_state", model.INSTANCE_HEALTH_OK).And().NotEqual("health_state", model.INSTANCE_HEALTH_NOTSET))
				if err != nil {
					return err
				}
				if len(instances) <= 0 {
					StartAlert(product.ID)
				}
			}
		}

		time.Sleep(10 * time.Second)
	}

	return nil
}

func expire(pid int) bool {
	recordTime, ok := exipreData[pid]
	if ok && time.Now().Unix()-recordTime >= 60*20 {
		return true
	}
	return false
}

func Register(sid string) error {
	pids := []int{}
	err := model.DeployInstanceList.GetDB().Select(&pids, "select distinct pid from "+model.DeployInstanceList.TableName+" where sid=?", sid)
	if err != nil {
		return err
	}
	for _, pid := range pids {
		_, ok := exipreData[pid]
		if !ok {
			exipreData[pid] = time.Now().Unix()
		}
	}
	return nil
}
