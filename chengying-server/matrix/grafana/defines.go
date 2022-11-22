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

import "time"

type GrafanaAlertSearchResponse struct {
	Id            int    `json:"id"`
	DashboardId   int    `json:"dashboardId"`
	DashboardUid  string `json:"dashboardUid"`
	DashboardSlug string `json:"dashboardSlug"`
	PanelId       int    `json:"panelId"`
	Name          string `json:"name"`
	State         string `json:"state"`
	NewStateDate  string `json:"newStateDate"`
	EvalDate      string `json:"evalDate"`
	//EvalData       string `json:"evalData"`
	ExecutionError string `json:"executionError"`
	Url            string `json:"url"`
	PanelTitle     string `json:"panelTitle"`
	EvalData       struct {
		EvalMatches []EvalMatches `json:"evalMatches,omitempty"`
		NoData      bool          `json:"noData,omitempty"`
	} `json:"evalData"`
}

type EvalMatches struct {
	Tags   map[string]string `json:"tags"`
	Metric string            `json:"metric"`
	Value  float64           `json:"value"`
}

type GrafanaAlertPausedResponse struct {
	AlertId int    `json:"alertId"`
	State   string `json:"state"`
	Message string `json:"message"`
}

type DashboardImportParam struct {
	Dashboard DashboardInfo       `json:"dashboard"`
	Inputs    []map[string]string `json:"inputs"`
	Overwrite bool                `json:"overwrite"`
	Message   string              `json:"message"`
}

type DashboardImportResponse struct {
	DashboardId      int    `json:"dashboardId"`
	Description      string `json:"description"`
	Imported         bool   `json:"imported"`
	ImportedRevision int    `json:"importedRevision"`
	ImportedUri      string `json:"importedUri"`
	ImportedUrl      string `json:"importedUrl"`
	Path             string `json:"path"`
	PluginId         string `json:"pluginId"`
	Removed          bool   `json:"removed"`
	Revision         int    `json:"revision"`
	Slug             string `json:"slug"`
	Title            string `json:"title"`
}

type DashboardInfo struct {
	Annotations   interface{}       `json:"annotations"`
	Editable      bool              `json:"editable"`
	GnetId        interface{}       `json:"gnetId"`
	GraphToolTip  int               `json:"graphTooltip"`
	Id            int               `json:"id"`
	Iteration     int               `json:"iteration"`
	Links         interface{}       `json:"links"`
	Panels        interface{}       `json:"panels"`
	Refresh       bool              `json:"refresh"`
	SchemaVersion int               `json:"schemaVersion"`
	Style         string            `json:"style"`
	Tags          []string          `json:"tags"`
	Templating    interface{}       `json:"templating"`
	Time          map[string]string `json:"time"`
	TimePicker    interface{}       `json:"timepicker"`
	TimeZone      string            `json:"timezone"`
	Title         string            `json:"title"`
	Uid           string            `json:"uid"`
	Version       int               `json:"version"`
	Inputs        interface{}       `json:"__inputs"`
	Requires      interface{}       `json:"__requires"`
}

type GetDashboardResponse struct {
	Dashboard DashboardInfo          `json:"dashboard"`
	Meta      map[string]interface{} `json:"meta"`
}

type DashboardJsonExport struct {
	Dashboard struct {
		Annotations   interface{}       `json:"annotations"`
		Editable      bool              `json:"editable"`
		GnetId        interface{}       `json:"gnetId"`
		GraphToolTip  int               `json:"graphTooltip"`
		Iteration     int               `json:"iteration"`
		Links         interface{}       `json:"links"`
		Panels        interface{}       `json:"panels"`
		Refresh       bool              `json:"refresh"`
		SchemaVersion int               `json:"schemaVersion"`
		Style         string            `json:"style"`
		Tags          []string          `json:"tags"`
		Templating    interface{}       `json:"templating"`
		Time          map[string]string `json:"time"`
		TimePicker    interface{}       `json:"timepicker"`
		TimeZone      string            `json:"timezone"`
		Title         string            `json:"title"`
		Uid           string            `json:"uid"`
		Version       int               `json:"version"`
		Inputs        interface{}       `json:"__inputs"`
		Requires      interface{}       `json:"__requires"`
	} `json:"dashboard"`
}

type GetDashboardInfo struct {
	Id      int `json:"id"`
	Version int `json:"version"`
}

type SaveDashboardResponse struct {
	Id      int    `json:"id"`
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	Uid     string `json:"uid"`
	Url     string `json:"url"`
	Version int    `json:"version"`
}

type DashboardListResponse struct {
	Id  int    `json:"id"`
	Uid string `json:"uid"`
}

type AnnotationResponse struct {
	AlertId     int    `json:"alertId"`
	AlertName   string `json:"alertName"`
	DashboardId int    `json:"dashboardId"`
	NewState    string `json:"newState"`
	PanelId     int    `json:"panelId"`
	Time        int64  `json:"time"`
}

type AlertRuleResponse struct {
	DashboardUid string `json:"dashboardUId"`
	PanelId      int    `json:"panelId"`
}

type AlertRuleTestParam struct {
	Dashboard DashboardInfo `json:"dashboard"`
	PanelId   int           `json:"panelId"`
}

type AlertRuleTestResponse struct {
	ConditionEvals string                   `json:"conditionEvals"`
	Firing         bool                     `json:"firing"`
	Logs           []map[string]interface{} `json:"logs"`
	Matches        []map[string]interface{} `json:"matches"`
	State          string                   `json:"state"`
	TimeMS         string                   `json:"timeMs"`
}

type GrafanaQueryResponse struct {
	Data   DataInfo `json:"data"`
	Status string   `json:"status"`
}

type DataInfo struct {
	Result     ResultSlice `json:"result"`
	ResultType string      `json:"resultType"`
}

type ResultSlice []ResultInfo

func (r ResultSlice) Len() int {
	return len(r)
}

func (r ResultSlice) Less(i, j int) bool {
	return len(r[i].Values) > len(r[j].Values)
}

func (r ResultSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type ResultInfo struct {
	Metric map[string]interface{} `json:"metric"`
	Values [][]interface{}        `json:"values"`
}

type AlertRules []GrafanaAlertSearchResponse

func (a AlertRules) Len() int {
	return len(a)
}

func (a AlertRules) Less(i, j int) bool {
	iTime, _ := time.Parse("2006-01-02T15:04:05+08:00", a[i].NewStateDate)
	jTime, _ := time.Parse("2006-01-02T15:04:05+08:00", a[j].NewStateDate)
	return iTime.After(jTime)
}

func (a AlertRules) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
