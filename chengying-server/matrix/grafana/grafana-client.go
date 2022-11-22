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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

var GrafanaURL *url.URL

func InitDashboard() {
	err, dashboardList := GetAllDashboards()
	if err != nil {
		log.Errorf("Dashboard init error: %v", err)
	}
	for _, dashboard := range dashboardList {
		err, dashboardResp := GetDashboardByUid(dashboard.Uid)
		if err != nil {
			log.Errorf("Dashboard: %v get error: %v", dashboard.Uid, err)
		}
		saveParam := &DashboardImportParam{
			Dashboard: dashboardResp.Dashboard,
			Overwrite: false,
			Message:   "",
		}
		err = SaveDashboard(saveParam)
		if err != nil {
			log.Errorf("Save dashboard %v error: %v", dashboardResp.Dashboard.Uid, err)
		}
	}

}

func InitGrafanaClient(grafanaUrl string) (err error) {
	GrafanaURL, err = url.Parse(grafanaUrl)
	if err != nil {
		log.Errorf("[InitGrafanaClient] init err: %v", err)
	}
	return err
}

func restCore(method, uri string, params map[string]string, body interface{}, resp interface{}) error {
	c := util.NewClient(util.DefaultClient)
	c.BaseURL = GrafanaURL

	r, err := c.NewRequest(method, uri, params, body, "")
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")

	_, err = c.Do(r, resp)

	return err
}

func GrafanaAlertsSearch(params map[string]string) (error, AlertRules) {
	//log.Debugf("[GrafanaClient] GrafanaSearch with params:%v ", params)
	resp := new([]GrafanaAlertSearchResponse)
	err := restCore("GET", GRAFANA_API_ALERTS, params, nil, resp)
	return err, *resp
}

func GrafanaAlertsPause(alertId string, pause bool) (error, GrafanaAlertPausedResponse) {
	log.Debugf("[GrafanaClient] GrafanaAlertsPause with alertId:%v ", alertId)
	resp := new(GrafanaAlertPausedResponse)
	body := map[string]bool{"paused": pause}
	err := restCore("POST", fmt.Sprintf(GRAFANA_API_ALERTS_PAUSE, alertId), nil, body, resp)
	return err, *resp
}

func GrafanaImportDashboard(param *DashboardImportParam) (error, DashboardImportResponse) {
	log.Debugf("[GrafanaClient] GrafanaDashboardImport")
	resp := new(DashboardImportResponse)
	err := restCore("POST", fmt.Sprintf(GRAFANA_API_DASHBOARD_IMPORT), nil, param, resp)
	return err, *resp
}

func GetDashboardByUid(uid string) (error, GetDashboardResponse) {
	resp := new(GetDashboardResponse)
	err := restCore("GET", fmt.Sprintf(GRAFANA_API_DASHBOARD_GET, uid), nil, nil, resp)
	return err, *resp
}

func ExportDashboardByUid(uid string) (error, DashboardJsonExport) {
	resp := new(DashboardJsonExport)
	err := restCore("GET", fmt.Sprintf(GRAFANA_API_DASHBOARD_GET, uid), nil, nil, resp)
	return err, *resp
}

func SaveDashboard(param *DashboardImportParam) error {
	resp := new(SaveDashboardResponse)
	err := restCore("POST", GRAFANA_API_DASHBOARD_SAVE, nil, param, resp)
	return err
}

func GetAllDashboards() (error, []DashboardListResponse) {
	resp := new([]DashboardListResponse)
	err := restCore("GET", GRAFANA_API_DASHBOARD_LIST, nil, nil, resp)
	return err, *resp
}

func GetAnnotations(params map[string]string) (error, []AnnotationResponse) {
	resp := new([]AnnotationResponse)
	err := restCore("GET", GRAFANA_API_ANNOTAION_LIST, params, nil, resp)
	return err, *resp
}

func AlertRuleTest(param *AlertRuleTestParam) (error, AlertRuleTestResponse) {
	resp := new(AlertRuleTestResponse)
	err := restCore("POST", GRAFANA_API_ALERTS_TEST, nil, param, resp)
	return err, *resp
}

func GrafanaQuery(query string, start, end, step int) (error, GrafanaQueryResponse) {
	resp := &GrafanaQueryResponse{}
	params := map[string]string{
		"query": query,
		"start": strconv.Itoa(start),
		"end":   strconv.Itoa(end),
		"step":  strconv.Itoa(step),
	}
	err := restCore("GET", GRAFANA_API_DATASOURCE_QUERY, params, nil, resp)
	return err, *resp
}

func GetDashboard(params []map[string]string) (error, []DashboardListResponse) {
	resp := new([]DashboardListResponse)
	c := util.NewClient(util.DefaultClient)
	c.BaseURL = GrafanaURL
	rel, err := url.Parse(GRAFANA_API_DASHBOARD_LIST)
	q := c.BaseURL.Query()
	if len(params) != 0 {
		for _, item := range params {
			for k, v := range item {
				q.Add(k, v)
				rel.RawQuery = q.Encode()
			}
		}
	}
	rel.Path = path.Join(c.BaseURL.Path, rel.Path)
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err, *resp
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = c.Do(req, resp)
	return err, *resp
}
