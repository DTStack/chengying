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
	"archive/zip"
	"bytes"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/grafana"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/context"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func ImportDashboard(ctx context.Context) apibase.Result {
	dashboardImportParam := &grafana.DashboardImportParam{}
	err := ctx.ReadJSON(dashboardImportParam)
	if err != nil {
		return fmt.Errorf("import dashboard read json %T err: %v", dashboardImportParam, err)
	}

	// call grafana import api
	err, resp := grafana.GrafanaImportDashboard(dashboardImportParam)
	if err != nil {
		return fmt.Errorf("import grafana dashboard err:%v", err)
	}
	log.Debugf("import dashboard, resp: %v", resp)

	// get dashboard id from grafana api via uid
	err, dashboardResp := grafana.GetDashboardByUid(dashboardImportParam.Dashboard.Uid)
	if err != nil {
		return fmt.Errorf("get dashboard by uid error: %v", err)
	}
	dashboardId := dashboardResp.Dashboard.Id
	// call grafana api to save dashboard for saving alerts automatically
	saveParam := &*dashboardImportParam
	saveParam.Overwrite = false
	saveParam.Dashboard.Id = dashboardId
	saveParam.Dashboard.Version = dashboardResp.Dashboard.Version
	saveParam.Message = ""
	err = grafana.SaveDashboard(saveParam)
	if err != nil {
		return fmt.Errorf("save dashboard error: %v", err)
	}
	return true
}

func GetDashboardAlerts(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	page, err := ctx.URLParamInt("page")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("start is empty"))
	}
	size, err := ctx.URLParamInt("size")
	if err != nil {
		size = 20
	}
	query := ctx.URLParam("query")
	paramErrs.CheckAndThrowApiParameterErrors()
	var alertList []grafana.GrafanaAlertSearchResponse
	urlParams := ctx.URLParams()
	err, alerts := grafana.GrafanaAlertsSearch(urlParams)
	if err != nil {
		log.Errorf("grafana search alerts error: %v", err)
		return err
	}
	// 为告警规则添加panel标题字段
	sort.Sort(alerts)
	for _, alertRule := range alerts {
		//搜索"告警名称"&"仪表盘名称"，不区分大小写
		if query != "" {
			query = strings.ToLower(query)
			newName := strings.ToLower(alertRule.Name)
			newDashboardSlug := strings.ToLower(alertRule.DashboardSlug)
			if strings.Contains(newName, query) || strings.Contains(newDashboardSlug, query) {
				alertList = append(alertList, alertRule)
			}
		} else {
			alertList = append(alertList, alertRule)
		}
	}
	var interfaceSlice = make([]interface{}, len(alertList))
	for index, data := range alertList {
		interfaceSlice[index] = data
	}
	pageData := util.Paginate(interfaceSlice, page, size)
	var resultList = make([]grafana.GrafanaAlertSearchResponse, 0)
	for _, page := range pageData {
		rule := page.(grafana.GrafanaAlertSearchResponse)
		panelTitle, _ := RetrievePanelTitle(rule.DashboardUid, rule.PanelId)
		rule.PanelTitle = panelTitle
		resultList = append(resultList, rule)
	}
	return map[string]interface{}{
		"data":  resultList,
		"total": len(alertList),
	}
}

func RetrievePanelTitle(dashboardUid string, alertPanelId int) (string, string) {
	err, dashboard := grafana.GetDashboardByUid(dashboardUid)
	if err != nil {
		log.Errorf("get dashboard by uid: %v error: %v", dashboardUid, err)
		return "", ""
	}

	panels := dashboard.Dashboard.Panels
	panelKind := reflect.TypeOf(panels).Kind()
	if panelKind == reflect.Slice || panelKind == reflect.Array {
		panelSlice := panels.([]interface{})
		for _, elem := range panelSlice {
			dat := elem.(map[string]interface{})
			panelId := dat["id"]
			if subPanels, ok := dat["panels"]; ok {
				subPanelSlice := subPanels.([]interface{})
				for _, subPanel := range subPanelSlice {
					panelMap := subPanel.(map[string]interface{})
					subPanelId := panelMap["id"]
					if subPanelId.(float64) == float64(alertPanelId) {
						return panelMap["title"].(string), dashboard.Dashboard.Title
					}
				}
			} else {
				if panelId.(float64) == float64(alertPanelId) {
					return dat["title"].(string), dashboard.Dashboard.Title
				}
			}
		}
	} else {
		log.Errorf("dashboard: %v, panel not array format", dashboardUid)
	}
	return "", ""
}

func DashboardAlertsPause(ctx context.Context) apibase.Result {
	var param struct {
		AlertId string `json:"alertId"`
		Paused  bool   `json:"paused"`
	}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}

	resultList := make([]grafana.GrafanaAlertPausedResponse, 0)
	alertIdList := strings.Split(param.AlertId, ",")
	for _, alertId := range alertIdList {
		err, resp := grafana.GrafanaAlertsPause(alertId, param.Paused)
		if err != nil {
			log.Errorf("pause grafana alerts: %v", err)
			return err
		}
		resultList = append(resultList, resp)
	}
	return map[string]interface{}{
		"data":  resultList,
		"total": len(resultList),
	}
}

func ExportDashboard(ctx context.Context) apibase.Result {
	uids := strings.Split(ctx.URLParam("dashboardId"), ",")
	err, dashs := grafana.GetAllDashboards()
	if err != nil {
		return fmt.Errorf("grafana server err %v", err)
	}
	checkUidExist := func(id int) string {
		for _, dash := range dashs {
			if dash.Id == id {
				return dash.Uid
			}
		}
		return ""
	}
	var (
		buf = new(bytes.Buffer)
		w   = zip.NewWriter(buf)
	)
	for _, uid := range uids {
		id, err := strconv.Atoi(uid)
		if err != nil {
			return fmt.Errorf("%s format not int", uid)
		}
		uid := checkUidExist(id)
		if uid == "" {
			return fmt.Errorf("dashbroad %d not exist", id)
		}
		err, response := grafana.ExportDashboardByUid(uid)
		if err != nil {
			return fmt.Errorf("get dashbroad %d file err %v", id, err)
		}
		f, err := w.Create(response.Dashboard.Title + ".json")
		if err != nil {
			return fmt.Errorf("get dashbroad %d file err %v", id, err)
		}
		data, _ := json.Marshal(response.Dashboard)
		_, err = f.Write(data)
	}
	err = w.Close()
	if err != nil {
		log.Errorf("zip err %v", err)
	}
	ctx.Header("Content-Disposition", "attachment; filename=dashbroad.zip")
	ctx.Header("Content-Type", "application/zip")
	_, err = buf.WriteTo(ctx.ResponseWriter())
	if err != nil {
		return fmt.Errorf("download zip err %v", err)
	}
	return nil
}
