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
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"dtstack.com/dtstack/easymatrix/matrix/grafana"
	"dtstack.com/dtstack/easymatrix/matrix/harole"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"encoding/json"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/kataras/iris/context"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"io/ioutil"
	"net/url"

	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var gofpdfDir string

func init() {
	gofpdfDir = "./"
}

func FontDir() string {
	return filepath.Join(gofpdfDir, "font")
}

func FontFile(fileStr string) string {
	return filepath.Join(FontDir(), fileStr)
}

func PdfDir(id int) string {
	return filepath.Join(base.WebRoot, "pdf", strconv.Itoa(id))
}

func PdfFile(id int, fileStr string) string {
	return filepath.Join(PdfDir(id), fileStr)
}

func ImageDir(id int) string {
	return filepath.Join(base.WebRoot, "img", strconv.Itoa(id))
}

func ImageFile(id int, fileStr string) string {
	return filepath.Join(ImageDir(id), fileStr)
}

func Filename(id int, baseStr string) string {
	return PdfFile(id, baseStr+".pdf")
}

type resultLists struct {
	ServiceName string `json:"service_name"`
	Ip          string `json:"ip"`
	Status      string `json:"status"`
	HARole      string `json:"ha_role"`
}

func GetServiceStatus(ctx context.Context) apibase.Result {
	// 获取当前所在父产品名称
	parentProduct, err := GetCurrentParentProduct(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	// 获取当前集群id
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	productList, err := model.DeployProductList.GetDeploySonProductName(parentProduct, clusterId)
	if err != nil {
		log.Errorf("get son product name error: %v", err)
		return err
	}

	result := map[string][]resultLists{}
	for _, productName := range productList {
		instanceList, err := model.DeployInstanceList.FindByProductNameAndClusterId(productName, clusterId)
		if err != nil {
			log.Errorf("%v", err)
		}
		var resultInfoList []resultLists
		for _, instance := range instanceList {
			resultInfo := resultLists{
				ServiceName: instance.ServiceName,
				Ip:          instance.Ip,
			}
			status := model.INSTANCE_NORMAL
			if instance.Status != model.INSTANCE_STATUS_RUNNING {
				status = model.INSTANCE_ABNORMAL
			} else if instance.HealthState != model.INSTANCE_HEALTH_OK && instance.HealthState != model.INSTANCE_HEALTH_NOTSET {
				status = model.INSTANCE_ABNORMAL
			}
			resultInfo.Status = status
			roleData := harole.RoleData(instance.Pid, instance.ServiceName)
			if roleData != nil {
				haRole, ok := roleData[instance.AgentId]
				if !ok {
					haRole = "-"
				}
				resultInfo.HARole = strings.Replace(haRole, "\r", "", 1)
			}
			resultInfoList = append(resultInfoList, resultInfo)
		}
		result[productName] = resultInfoList
	}

	return result

}

type alertInfo struct {
	AlertName      string `json:"alert_name"`
	State          string `json:"state"`
	DashboardName  string `json:"dashboard_name"`
	DashboardTitle string `json:"dashboard_title"`
	Time           string `json:"time"`
}

func GetAlertHistory(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	from, err := ctx.URLParamInt("from")
	if err != nil {
		paramErrs.AppendError("$", "param from is empty")
	}
	to, err := ctx.URLParamInt("to")
	if err != nil {
		paramErrs.AppendError("$", "param to is empty")
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	alertList := FormatAlertList(from, to)
	return map[string]interface{}{
		"count": len(alertList),
		"data":  alertList,
	}
}

func FormatAlertList(from, to int) []alertInfo {
	params := map[string]string{}
	params["from"] = strconv.Itoa(from)
	params["to"] = strconv.Itoa(to)
	params["type"] = "alert"
	params["limit"] = "1000"
	err, resp := grafana.GetAnnotations(params)
	if err != nil {
		log.Errorf("get annotations from grafana error: %v", err)
		return make([]alertInfo, 0)
	}

	var alertList []alertInfo

	for _, info := range resp {
		if info.NewState != "ok" && info.NewState != "paused" && info.NewState != "pending" {
			param := map[string]string{
				"dashboardId": strconv.Itoa(info.DashboardId),
				"panelId":     strconv.Itoa(info.PanelId),
			}
			err, alertRule := grafana.GrafanaAlertsSearch(param)
			if err != nil || len(alertRule) == 0 {
				log.Errorf("get alert rule: %v error: %v", info.AlertId, err)
				continue
			}
			panelTitle, dashboardName := RetrievePanelTitle(alertRule[0].DashboardUid, alertRule[0].PanelId)
			alert := alertInfo{
				AlertName:      info.AlertName,
				State:          info.NewState,
				DashboardName:  dashboardName,
				DashboardTitle: panelTitle,
				Time:           time.Unix(info.Time/1000, 0).Format("2006-01-02 15:04:05"),
			}
			alertList = append(alertList, alert)
		}
	}
	return alertList
}

type ReportHostInfo struct {
	Ip         string `json:"ip"`
	Cpu        string `json:"cpu"`
	Mem        string `json:"mem"`
	SystemDisk string `json:"system_disk"`
	DataDisk   string `json:"data_disk"`
}

func GetHostStatus(ctx context.Context) apibase.Result {
	// 获取当前所在父产品名称
	parentProduct, err := GetCurrentParentProduct(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	// 获取当前集群Id
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	// 获取当前集群下所有接入的主机
	query := "select deploy_host.ip from deploy_cluster_host_rel " +
		"left join deploy_host on deploy_cluster_host_rel.sid=deploy_host.sid " +
		"left join deploy_instance_list on deploy_host.sid=deploy_instance_list.sid " +
		"left join deploy_product_list on deploy_instance_list.pid=deploy_product_list.id " +
		"left join sidecar_list on sidecar_list.id=deploy_host.sid where deploy_host.sid!='' " +
		"and deploy_host.isDeleted=0 and deploy_product_list.parent_product_name=? and deploy_cluster_host_rel.clusterId=? group by deploy_host.sid"
	var ipList []string
	if err := model.USE_MYSQL_DB().Select(&ipList, query, parentProduct, clusterId); err != nil {
		log.Errorf("get hosts error: %v", err)
		return err
	}

	hostStatusMap := map[string]ReportHostInfo{}
	// 初始化主机正常状态
	for _, ip := range ipList {
		hostStatusMap[ip] = ReportHostInfo{
			Ip:         ip,
			Cpu:        model.INSTANCE_NORMAL,
			Mem:        model.INSTANCE_NORMAL,
			SystemDisk: model.INSTANCE_NORMAL,
			DataDisk:   model.INSTANCE_NORMAL,
		}
	}

	// 获取Host Overview仪表盘信息
	err, dashboardResp := grafana.GetDashboardByUid("Ne_roaViz")
	if err != nil {
		log.Errorf("get host overview dashboard error: %v", err)
	}
	// 获取cpu告警信息
	err, cpuAlerts := grafana.AlertRuleTest(&grafana.AlertRuleTestParam{
		Dashboard: dashboardResp.Dashboard,
		PanelId:   38,
	})
	if err != nil {
		log.Errorf("test cpu alerts error: %v", err)
	}
	reg := regexp.MustCompile(`\w+-(?P<instance>(\d+\.)+\d+):\d+\s?(?P<mnt>/\w*)?`)
	if cpuAlerts.State != "ok" {
		for _, match := range cpuAlerts.Matches {
			metric := reg.FindStringSubmatch(match["metric"].(string))
			if metric != nil {
				if host, ok := hostStatusMap[metric[1]]; ok {
					host.Cpu = model.INSTANCE_ABNORMAL
					hostStatusMap[metric[1]] = host
				}
			}
		}
	}

	// 获取内存告警信息
	err, memoryAlerts := grafana.AlertRuleTest(&grafana.AlertRuleTestParam{
		Dashboard: dashboardResp.Dashboard,
		PanelId:   50,
	})
	if err != nil {
		log.Errorf("test memory alerts error: %v", err)
	}
	if memoryAlerts.State != "ok" {
		for _, match := range memoryAlerts.Matches {
			metric := reg.FindStringSubmatch(match["metric"].(string))
			if metric != nil {
				if host, ok := hostStatusMap[metric[1]]; ok {
					host.Mem = model.INSTANCE_ABNORMAL
					hostStatusMap[metric[1]] = host
				}
			}
		}

	}

	// 获取磁盘告警信息
	err, diskAlerts := grafana.AlertRuleTest(&grafana.AlertRuleTestParam{
		Dashboard: dashboardResp.Dashboard,
		PanelId:   44,
	})
	if err != nil {
		log.Errorf("test disk alerts error: %v", err)
	}
	if diskAlerts.State != "ok" {
		for _, match := range diskAlerts.Matches {
			metric := reg.FindStringSubmatch(match["metric"].(string))
			if metric != nil {
				if host, ok := hostStatusMap[metric[1]]; ok {
					if metric[3] == "/" {
						host.SystemDisk = model.INSTANCE_ABNORMAL
					}
					if metric[3] == "/data" {
						host.DataDisk = model.INSTANCE_ABNORMAL
					}
					hostStatusMap[metric[1]] = host
				}
			}
		}
	}

	var resultList = make([]ReportHostInfo, 0)
	for _, v := range hostStatusMap {
		resultList = append(resultList, v)
	}

	return resultList
}

func GetGraphConfig(ctx context.Context) apibase.Result {
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	configList, err := model.InspectReportTemplate.GetTemplateConfig(clusterId)
	if err != nil {
		log.Errorf("%v", err)
	}
	return configList
}

type ChartsInfo struct {
	X []float64                `json:"x"`
	Y []map[string]interface{} `json:"y"`
}

type targetInfo struct {
	Expr         string `json:"expr"`
	LegendFormat string `json:"legend_format"`
}

func GetGraphData(ctx context.Context) apibase.Result {
	paramsErr := apibase.NewApiParameterErrors()
	from, err := ctx.URLParamInt("from")
	if err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param from is empty"))
	}
	to, err := ctx.URLParamInt("to")
	if err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param to is empty"))
	}
	targets := ctx.URLParam("targets")
	var targetsList []targetInfo
	if err := json.Unmarshal([]byte(targets), &targetsList); err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param targets format error:%v", err))
	}
	//unit := ctx.URLParam("unit")
	decimal, err := ctx.URLParamInt("decimal")
	if err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param decimal is empty"))
	}
	paramsErr.CheckAndThrowApiParameterErrors()

	var chartInfo ChartsInfo
	// 取第一个target
	if len(targetsList) > 0 {
		var x = make([]float64, 0)
		var y = make([]map[string]interface{}, 0)
		var xEdited = false
		for _, target := range targetsList {
			legendFormat := target.LegendFormat
			err, queryResponse := grafana.GrafanaQuery(target.Expr, from/1000, to/1000, (to/1000-from/1000)/600)
			if err != nil {
				log.Errorf("grafana query error: %v", err)
				return nil
			}
			var maxValue = 0

			resultList := queryResponse.Data.Result
			sort.Sort(resultList)
			if len(resultList) > 0 {
				for index, result := range resultList {
					metric := result.Metric
					if len(metric) == 0 {
						continue
					}
					values := result.Values
					item := map[string]interface{}{}
					item["title"] = formatLegend(legendFormat, metric)
					var data []interface{}
					for _, value := range values {
						if index == 0 && !xEdited {
							x = append(x, value[0].(float64))
						}
						value, _ := FormatFloatCeil(value[1].(string), decimal)
						data = append(data, value)
					}
					if maxValue <= len(values) {
						maxValue = len(values)
					} else {
						zeroSlice := make([]interface{}, 0)
						for i := 0; i < maxValue-len(values); i++ {
							zeroSlice = append(zeroSlice, float64(0))
						}
						data = append(zeroSlice, data...)
					}
					item["data"] = data
					y = append(y, item)

				}
			}
			xEdited = true
		}
		chartInfo.X = x
		chartInfo.Y = y
	}

	return chartInfo
}

func FormatGraphData(from, to, decimal int, targets string) ChartsInfo {
	var targetsList []targetInfo
	if err := json.Unmarshal([]byte(targets), &targetsList); err != nil {
		log.Errorf("param targets format error:%v", err)
	}
	return FormatCharInfo(targetsList, from, to, decimal)
}

func FormatCharInfo(targetsList []targetInfo, from, to, decimal int) ChartsInfo {
	var chartInfo ChartsInfo
	// 取第一个target
	if len(targetsList) > 0 {
		var x = make([]float64, 0)
		var y = make([]map[string]interface{}, 0)
		var xEdited = false
		for _, target := range targetsList {
			legendFormat := target.LegendFormat
			err, queryResponse := grafana.GrafanaQuery(target.Expr, from/1000, to/1000, (to/1000-from/1000)/600)
			if err != nil {
				log.Errorf("grafana query error: %v", err)
			}
			var maxValue = 0

			resultList := queryResponse.Data.Result
			sort.Sort(resultList)
			if len(resultList) > 0 {
				for index, result := range resultList {
					metric := result.Metric
					if len(metric) == 0 {
						continue
					}
					values := result.Values
					item := map[string]interface{}{}
					item["title"] = formatLegend(legendFormat, metric)
					var data []interface{}
					for _, value := range values {
						if index == 0 && !xEdited {
							x = append(x, value[0].(float64))
						}
						value, _ := FormatFloatCeil(value[1].(string), decimal)
						data = append(data, value)
					}
					if maxValue <= len(values) {
						maxValue = len(values)
					} else {
						zeroSlice := make([]interface{}, 0)
						for i := 0; i < maxValue-len(values); i++ {
							zeroSlice = append(zeroSlice, float64(0))
						}
						data = append(zeroSlice, data...)
					}
					item["data"] = data
					y = append(y, item)

				}
			}
			xEdited = true
		}
		chartInfo.X = x
		chartInfo.Y = y
	}

	return chartInfo
}

func formatLegend(legendFormat string, metric map[string]interface{}) string {
	reg := regexp.MustCompile(`\{{(?P<field>\w+)}}`)
	for {
		matches := reg.FindStringSubmatch(legendFormat)
		if matches != nil {
			if value, ok := metric[matches[1]]; ok {
				legendFormat = strings.Replace(legendFormat, matches[0], value.(string), 1)
			}
		} else {
			break
		}
	}
	return legendFormat
}

func FormatFloatCeil(num string, decimal int) (float64, error) {
	value, _ := strconv.ParseFloat(num, 64)
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Ceil(value*d)/d, 'f', -1, 64)
	return strconv.ParseFloat(res, 64)
}

type GenerateParam struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func StartGenerateReport(ctx context.Context) apibase.Result {
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		return err
	}
	var param GenerateParam
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("[StartGenerateReport] invalid params")
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("found panic: %v", err)
		}
	}()
	fromDate := time.Unix(int64(param.From/1000), 0).Format("2006-01-02")
	toDate := time.Unix(int64(param.To)/1000, 0).Format("2006-01-02")
	reportName := fmt.Sprintf("运维周报（%s至%s）", fromDate, toDate)
	reportId, err := model.InspectReport.NewInspectReport(reportName, "RUNNING", clusterId)
	if err != nil {
		log.Errorf("[StartGenerateReport] new db record error: %v", err)
		return err
	}
	go generateReport(ctx, param, int(reportId), clusterId)
	return map[string]interface{}{
		"report_id": reportId,
	}
}

func GetReportProgress(ctx context.Context) apibase.Result {
	reportId, err := ctx.URLParamInt("id")
	if err != nil {
		log.Errorf("Get report by id param error: %v", err)
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("found panic: %v", err)
		}
	}()
	reportInfo, err := model.InspectReport.GetById(reportId)
	if err != nil {
		log.Errorf("Get report by id db error: %v", err)
		return err
	}
	return reportInfo
}

func Download(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	filePath := ctx.URLParam("file_path")
	if filePath == "" {
		paramErrs.AppendError("$", "缺少文件路径")
	}
	id, err := ctx.URLParamInt("id")
	if err != nil {
		paramErrs.AppendError("$", "缺少id")
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open report file error")
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	reportInfo, err := model.InspectReport.GetById(id)
	if err != nil {
		log.Errorf("get report info by id error: %v", err)
		return err
	}
	absolutePath := reportInfo.Name
	fileNames := absolutePath[strings.LastIndex(absolutePath, "/")+1:]
	fileName := url.QueryEscape(fileNames)
	ctx.Header("Content-Disposition", "attachment;filename=\""+fileName+".pdf\"")
	ctx.Write(content)
	return nil
}

func generateReport(ctx context.Context, param GenerateParam, id, clusterId int) error {
	defer func() {
		if err := os.RemoveAll(ImageDir(id)); err != nil {
			log.Errorf("Remove image dir of %d error: %v", id, err)
		}
	}()
	p := util.NewPdfGenerator(id)
	p.AddFont("Simhei", FontFile("simhei.ttf"))
	p.SetFont("Arial", "", 14)
	p.AddFooter()
	p.AddPage()
	// title
	p.AddText(p.PageWidth()*0.4, 18, 7, "Simhei", "巡检报告", 0, 0, 0)
	// 集群状态汇总
	p.AddLine(p.Left())
	p.AddText(p.Left(), 9, 5, "Simhei", "集群状态汇总", 0, 0, 0)
	p.AddText(p.Left(), 6, 2, "Simhei", "1.报告中设计的状态为报告下载时间点的状态", 112, 128, 144)
	p.AddText(p.Left(), 6, 2, "Simhei", "2.状态为“正常”表示节点或应用当前的健康状态为健康，监控指标没有告警。状态为“异常”表示节点或应用当前的健康状态为不健康，或者监控指标有告警", 112, 128, 144)
	p.Ln(3)

	hostStatus := func() {
		hostInfoListInterface := GetHostStatus(ctx)
		if hostInfoList, ok := hostInfoListInterface.([]ReportHostInfo); ok {
			p.AddText(p.Left(), 8, 7, "Simhei", "·节点状态", 0, 0, 0)
			headers := []string{"节点", "cpu", "内存", "系统盘", "数据盘"}
			datas := [][]string{}
			for _, hostInfo := range hostInfoList {
				data := []string{hostInfo.Ip, hostInfo.Cpu, hostInfo.Mem, hostInfo.SystemDisk, hostInfo.DataDisk}
				datas = append(datas, data)
			}
			p.AddTable(datas, headers)
		}
	}

	appStatus := func() {
		p.AddText(p.Left(), 8, 7, "Simhei", "·应用状态", 0, 0, 0)
		serviceStatusMapInterface := GetServiceStatus(ctx)
		if serviceStatusMap, ok := serviceStatusMapInterface.(map[string][]resultLists); ok {
			for k, v := range serviceStatusMap {
				p.AddText(p.Left(), 7, 6, "Simhei", k, 0, 0, 0)
				var headers []string
				datas := [][]string{}
				if k == "DTBase" {
					headers = []string{"服务", "节点", "角色", "状态"}
					for _, info := range v {
						data := []string{info.ServiceName, info.Ip, info.HARole, info.Status}
						datas = append(datas, data)
					}
				} else {
					headers = []string{"服务", "节点", "状态"}
					for _, info := range v {
						data := []string{info.ServiceName, info.Ip, info.Status}
						datas = append(datas, data)
					}
				}
				p.AddTable(datas, headers)
			}
		}
	}

	alertHistory := func() {
		p.AddText(p.Left(), 8, 7, "Simhei", "·告警记录", 0, 0, 0)
		alertHistoryList := FormatAlertList(param.From, param.To)
		headers := []string{"告警名称", "状态", "仪表盘名称(组件)", "仪表盘标题", "告警时间"}
		datas := [][]string{}
		for _, info := range alertHistoryList {
			data := []string{info.AlertName, info.State, info.DashboardName, info.DashboardTitle, info.Time}
			datas = append(datas, data)
		}
		p.AddTable(datas, headers)
	}

	hostStatus()
	appStatus()
	alertHistory()

	err := model.InspectReport.UpdateProgress(id, 30, "", "")

	var configMaps = map[string][]model.BaseTemplateConfig{}
	configList, err := model.InspectReportTemplate.GetTemplateConfig(clusterId)
	for _, config := range configList {
		module := strings.TrimSpace(config.Module)
		moduleConfigList, ok := configMaps[module]
		if !ok {
			moduleConfigList = []model.BaseTemplateConfig{}
		}
		moduleConfigList = append(moduleConfigList, config)
		configMaps[module] = moduleConfigList
	}

	// 集群状态详细内容
	p.AddLine(p.Left())
	p.AddText(p.Left(), 9, 5, "Simhei", "集群状态详细内容", 0, 0, 0)
	p.AddText(p.Left(), 8, 7, "Simhei", "·节点状态", 0, 0, 0)

	if err := os.MkdirAll(ImageDir(p.Id), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(PdfDir(p.Id), 0755); err != nil {
		return err
	}

	// 渲染节点监控折线图
	hostConfigs, ok := configMaps["System"]
	if !ok {
		log.Errorf("lack host config")
	} else {
		for _, hostConfig := range hostConfigs {
			addLineChart(param.From, param.To, hostConfig, p)
		}
	}

	err = model.InspectReport.UpdateProgress(id, 60, "", "")

	p.AddText(p.Left(), 8, 7, "Simhei", "·应用状态", 0, 0, 0)
	for k, v := range configMaps {
		if k != "System" {
			if existModuleData(param.From, param.To, v) {
				p.AddText(p.Left(), 8, 7, "Simhei", k, 0, 0, 0)
				for _, hostConfig := range v {
					addLineChart(param.From, param.To, hostConfig, p)
				}
			}
		}
	}

	fileStr := Filename(id, strconv.Itoa(id))
	err = p.OutputFileAndClose(fileStr)

	if err != nil {
		log.Errorf("generate pdf error: %v", err)
		err = model.InspectReport.UpdateProgress(id, 60, "", "FAIL")
	}

	err = model.InspectReport.UpdateProgress(id, 100, fileStr, "SUCCESS")
	return err
}

func existModuleData(from, to int, configList []model.BaseTemplateConfig) bool {
	for _, config := range configList {
		chartInfo := FormatGraphData(from, to, config.Decimal, config.Targets)
		if len(chartInfo.X) != 0 {
			return true
		}
	}
	return false
}

func addLineChart(from, to int, hostConfig model.BaseTemplateConfig, p *util.PdfGenerator) {
	chartsInfo := FormatGraphData(from, to, hostConfig.Decimal, hostConfig.Targets)
	series := make([]chart.Series, len(chartsInfo.Y))
	if len(chartsInfo.X) == 0 {
		return
	}
	for i := 0; i < len(chartsInfo.Y); i++ {
		xValues := make([]time.Time, len(chartsInfo.X))
		yValues := make([]float64, len(chartsInfo.X))
		title := chartsInfo.Y[i]["title"].(string)
		datas := chartsInfo.Y[i]["data"].([]interface{})
		for j := 0; j < len(chartsInfo.X); j++ {
			xValues[j] = time.Unix(int64(chartsInfo.X[j]), 0)
			if hostConfig.Unit == "byte" {
				yValues[j] = datas[j].(float64) / (1000 * 1000 * 1000)
			} else {
				yValues[j] = datas[j].(float64)
			}
		}
		series[i] = chart.TimeSeries{
			Name: title,
			Style: chart.Style{
				StrokeColor: drawing.Color{
					R: uint8(rand.Intn(256)),
					G: uint8(rand.Intn(256)),
					B: uint8(rand.Intn(256)),
					A: uint8(256 - 1),
				},
			},
			XValues: xValues,
			YValues: yValues,
		}
	}
	lineChartStyle := chart.Style{
		Padding: chart.Box{
			Top: 40,
		},
	}
	ya := defineYaxis(hostConfig.Metric)
	graph := chart.Chart{
		Title: hostConfig.Metric,
		TitleStyle: chart.Style{
			Font: GetChineseFont(),
		},
		Background: lineChartStyle,
		XAxis: chart.XAxis{
			Name: "时间",
			NameStyle: chart.Style{
				Font: GetChineseFont(),
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("2006-01-02 15:04:05"),
		},
		YAxis:  ya,
		Series: series,
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}
	imgName := ImageFile(p.Id, hostConfig.Module+"_"+strings.ReplaceAll(hostConfig.Metric, "/", "_")+".png")
	f, _ := os.Create(imgName)
	defer f.Close()
	fmt.Println(imgName)
	graph.Render(chart.PNG, f)
	p.AddLineChart("png", imgName)
}

func defineYaxis(metricName string) chart.YAxis {
	ya := chart.YAxis{
		Name: "值",
		NameStyle: chart.Style{
			Font: GetChineseFont(),
		},
	}
	var max float64
	if strings.Contains(metricName, "up") || strings.Contains(metricName, "Up") {
		max = 1
	} else if strings.Contains(metricName, "%") {
		max = 100
	} else {
		max = 10
	}
	ya.Range = &chart.ContinuousRange{
		Max: max,
	}
	return ya
}

func GetChineseFont() *truetype.Font {
	fontBytes, err := ioutil.ReadFile(FontFile("simhei.ttf"))
	if err != nil {
		return nil
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil
	}
	return font
}

const (
	MYSQL_IO_RUNNING           = "mysql_slave_status_slave_io_running{cluster_name='{cluster_name}'}"
	MYSQL_SQL_RUNNING          = "mysql_slave_status_slave_sql_running{cluster_name='{cluster_name}'}"
	DATA_NODE_LIVE_NUM         = "Hadoop_NameNode_NumLiveDataNodes{cluster_name='{cluster_name}',name='FSNamesystem',product_name='Hadoop',service_name='hdfs_namenode'}"
	DATA_NODE_DEAD_NUM         = "Hadoop_NameNode_NumDeadDataNodes{cluster_name='{cluster_name}',name='FSNamesystem',product_name='Hadoop',service_name='hdfs_namenode'}"
	NAME_NODE_TOTAL_FILE       = "Hadoop_NameNode_TotalFiles{cluster_name='{cluster_name}',product_name='Hadoop',service_name='hdfs_namenode'}"
	SERVICE_GC_FORM_TITLE      = "service_gc_statistics"
	DIR_USAGE_FORM_TITLE       = "dir_usage_statistics"
	ALARM_FROM_TITLE           = "alarm_indicators"
	HDFS_FILE_USAGE_FORM_TITLE = "hdfs_file_usage"
	CPU_USAGE_FROM_TITLE       = "cpu_usage"
	MEM_USAGE_FROM_TITLE       = "men_usage"
	DISK_USAGE_FROM_TITLE      = "disk_usage"
	INODE_USAGE_FROM_TITLE     = "inode_usage"
	SWAP_USAGE_FROM_TITLE      = "swap_usage"
	SERVICE_GC                 = "floor(delta(jvm_gc_collection_seconds_count{cluster_name='{cluster_name}',gc='{gc_type}'}[{time}m])) >= {freq}"
	FILE_USAGE                 = "100-node_filesystem_free{cluster_name = '{cluster_name}',mountpoint='/data'}/node_filesystem_size{cluster_name = '{cluster_name}',mountpoint='/data'} * 100 "
	CPU_USAGE                  = "100 - ( avg(irate(node_cpu{mode='idle',cluster_name='{cluster_name}'}[5m])) by (instance) * 100 ) >= {config}"
	INODE_USAGE                = "100*((node_filesystem_files{cluster_name='{cluster_name}',device=~'/dev/.*',mountpoint='/data'}-node_filesystem_files_free{cluster_name='{cluster_name}',device=~'/dev/.*',mountpoint='/data'})/node_filesystem_files{cluster_name='{cluster_name}',device=~'/dev/.*',mountpoint='/data'}) >= {config}"
	DISK_USAGE                 = "(100-100*((node_filesystem_avail{cluster_name='{cluster_name}',device=~'/dev/.*',mountpoint='/data'}/node_filesystem_size{cluster_name='{cluster_name}',device=~'/dev/.*',mountpoint='/data'}))) >= {config}"
	MEM_USAGE                  = "(1-(sum by(instance)(node_memory_MemFree+node_memory_Buffers+node_memory_Cached))/(sum by(instance)(node_memory_MemTotal)))*100 >= {config}"
	SWAP_USAGE                 = "((sum by(instance)(node_memory_SwapTotal))-(sum by(instance)(node_memory_SwapFree))-(sum by(instance)(node_memory_SwapCached)) )/(8*1024*1024*1024) >= 0"
)

type PlatformInspectCommonParam struct {
	ClusterId int `json:"cluster_id"`
	From      int `json:"from"`
	To        int `json:"to"`
}

type NodeInfo struct {
	Total         int `json:"total"`
	AgentErrorNum int `json:"agent_error_num"`
	AlertingNum   int `json:"alerting_num"`
}
type ServiceInfo struct {
	RunningFailNum       int `json:"running_fail_num"`
	HostDownNum          int `json:"host_down_num"`
	HealthyCheckErrorNum int `json:"healthy_check_error_num"`
	AlertingNum          int `json:"alerting_num"`
}
type Result struct {
	NodeStatus       NodeInfo    `json:"node_status"`
	ServiceStatus    ServiceInfo `json:"service_status"`
	MysqlSlaveStatus int         `json:"mysql_slave_status"`
	HaveNameNode     int         `json:"have_name_node"`
}

// GetPlatformInspectBaseInfoState 	godoc
// @Summary      	获取巡检报告基本信息
// @Description  	获取巡检报告基本信息
// @Tags         	Inspect
// @Produce      	json
// @Param			cluster_id query int true "Cluster ID"
// @Success      	200		{object} 	Result
// @Router       	/api/v2/platform/inspect/baseInfo/status [get]
func GetPlatformInspectBaseInfoState(ctx context.Context) apibase.Result {
	paramsErr := apibase.NewApiParameterErrors()
	clusterId, err := ctx.URLParamInt("cluster_id")
	if err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param cluster_id is empty"))
	}
	alertParam := make(map[string]string, 0)
	err, dashboardResp := grafana.GetDashboardByUid("Ne_roaViz")
	if err != nil {
		log.Errorf("get host overview dashboard error: %v", err)
	}
	allAlerts, err := GetAlertList(alertParam)
	alertParam = map[string]string{
		"dashboardId": strconv.Itoa(dashboardResp.Dashboard.Id),
	}
	hostAlerts, err := GetAlertList(alertParam)

	nodeStates, err := GetPlatformInspectBaseInfoByClusterId(clusterId, allAlerts, hostAlerts)
	result := Result{
		NodeStatus: NodeInfo{
			Total:         nodeStates.NodeTotal,
			AlertingNum:   nodeStates.AlertingNum,
			AgentErrorNum: nodeStates.AgentErrorNum,
		},
		ServiceStatus: ServiceInfo{
			AlertingNum:          nodeStates.ServiceAlertingNum,
			HostDownNum:          nodeStates.ServiceHostDownNum,
			HealthyCheckErrorNum: nodeStates.ServiceHealthyErrorNum,
			RunningFailNum:       nodeStates.ServiceRunErrorNum,
		},
		MysqlSlaveStatus: nodeStates.MysqlSlaveStatus,
		HaveNameNode:     nodeStates.HasNameNode,
	}
	return result
}

type Response struct {
	NameNodeMem      string `json:"name_node_mem"`
	DataNodeMem      string `json:"data_node_mem"`
	DataNodeLiveNums int    `json:"data_node_live_nums"`
	DataNodeDeadNums int    `json:"data_node_dead_nums"`
	HdfsFileNum      int    `json:"hdfs_file_num"`
}

// GetPlatformInspectNameNodeBaseInfo 	godoc
// @Summary      	大数据运行服务基础信息
// @Description  	大数据运行服务基础信息
// @Tags         	Inspect
// @Produce      	json
// @Param			cluster_id query int true "Cluster ID"
// @Success      	200		{object} 	Response
// @Router       	/api/v2/platform/inspect/baseInfo/name_node [get]
func GetPlatformInspectNameNodeBaseInfo(ctx context.Context) apibase.Result {
	var result Response
	paramsErr := apibase.NewApiParameterErrors()
	clusterId, err := ctx.URLParamInt("cluster_id")
	if err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param cluster_id is empty"))
	}
	namenodeMem, datanodeMem, err := model.DeployInstanceList.GetNameNodeConfigByIdAndServiceName(clusterId)
	if err != nil {
		return err
	}
	datanodeLiveExpr := DATA_NODE_LIVE_NUM
	datanodeDeadExpr := DATA_NODE_DEAD_NUM
	namenodeTotalFileExpr := NAME_NODE_TOTAL_FILE
	clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	datanodeLiveExpr = strings.Replace(datanodeLiveExpr, "{cluster_name}", clusterInfo.Name, -1)
	datanodeDeadExpr = strings.Replace(datanodeDeadExpr, "{cluster_name}", clusterInfo.Name, -1)
	namenodeTotalFileExpr = strings.Replace(namenodeTotalFileExpr, "{cluster_name}", clusterInfo.Name, -1)
	err, datanodeLiveResponse := grafana.GrafanaQuery(datanodeLiveExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("get data node live num error: %v", err)
		result.DataNodeLiveNums = 0
	}

	for _, r := range datanodeLiveResponse.Data.Result {
		values := r.Values
		v, _ := values[0][1].(string)
		result.DataNodeLiveNums, _ = strconv.Atoi(v)
		break
	}
	err, datanodeDeadResponse := grafana.GrafanaQuery(datanodeDeadExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("get data node dead num error: %v", err)
		result.DataNodeDeadNums = 0
	}

	for _, r := range datanodeDeadResponse.Data.Result {
		values := r.Values
		v, _ := values[0][1].(string)
		result.DataNodeDeadNums, _ = strconv.Atoi(v)
		break
	}
	err, namenodeTotalFileResponse := grafana.GrafanaQuery(namenodeTotalFileExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("get name node total file num error: %v", err)
		result.HdfsFileNum = 0
	}

	for _, r := range namenodeTotalFileResponse.Data.Result {
		values := r.Values
		v, _ := values[0][1].(string)
		result.HdfsFileNum, _ = strconv.Atoi(v)
		break
	}
	result.NameNodeMem = namenodeMem
	result.DataNodeMem = datanodeMem
	return result
}
func GetAlertList(param map[string]string) (grafana.AlertRules, error) {
	err, alerts := grafana.GrafanaAlertsSearch(param)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

type InspectNodeStates struct {
	NodeTotal              int `json:"total"`
	AgentErrorNum          int `json:"agent_error_num"`
	AlertingNum            int `json:"alerting_num"`
	ServiceRunErrorNum     int `json:"service_run_err_num"`
	ServiceHostDownNum     int `json:"service_host_down_num"`
	ServiceHealthyErrorNum int `json:"service_healthy_error_num"`
	ServiceAlertingNum     int `json:"service_alerting_num"`
	HasNameNode            int `json:"has_name_node"`
	MysqlSlaveStatus       int `json:"mysql_slave_status"`
}

func GetPlatformInspectBaseInfoByClusterId(id int, alerts, hostAlerts grafana.AlertRules) (InspectNodeStates, error) {
	var result InspectNodeStates
	clusterInfo, err := model.DeployClusterList.GetClusterInfoById(id)
	if err != nil {
		return result, fmt.Errorf("[GetPlatformInspectBaseInfoByClusterId]Get cluster info by cluster_id error: %v ", err)
	}

	hostsInfo, err := model.DeployHostList.GetInspectNodeInfoByClusterId(id)
	if err != nil {
		return result, fmt.Errorf("[GetPlatformInspectBaseInfoByClusterId]Get host info by cluster_id error: %v ", err)
	}
	result.MysqlSlaveStatus = GetMysqlSalveStatusByGrafana(clusterInfo.Name)
	unRunningNum, unHealthyNum, err := model.DeployInstanceList.GetInspectServiceInfoById(id)
	if err != nil {
		return result, fmt.Errorf("[GetPlatformInspectBaseInfoByClusterId]Get service info by cluster_id error: %v ", err)
	}
	result.ServiceRunErrorNum, result.ServiceHealthyErrorNum = unRunningNum, unHealthyNum
	result.NodeTotal = len(hostsInfo)
	hostAlerting := make(map[string]bool, 0)
	hostService := make(map[string][]string, 0)
	hostDownService := make(map[string]bool, 0)
	for _, host := range hostsInfo {

		serviceList := strings.Split(host.ServiceList, ",")
		hostService[host.Ip] = serviceList
		hostAlerting[host.Ip] = false
		if result.HasNameNode != 1 && strings.Index(host.ServiceList, "hdfs_namenode") != -1 {
			result.HasNameNode = 1
		}
		if host.IsRunning == false {
			result.AgentErrorNum = result.AgentErrorNum + 1
			for _, v := range serviceList {
				if _, ok := hostDownService[v]; !ok {
					hostDownService[v] = true
					result.ServiceHostDownNum++
				}
			}
		}
	}

	for _, alert := range hostAlerts {
		if alert.State != "ok" && alert.State != "paused" && alert.State != "pending" {
			for _, match := range alert.EvalData.EvalMatches {
				if instance, ok := match.Tags["instance"]; ok {
					ip := strings.Split(instance, ":")[0]
					if _, oks := hostAlerting[ip]; oks && hostAlerting[ip] == false {
						hostAlerting[ip] = true
						result.AlertingNum++
					}
				}
			}
		}
	}

	evalMatches := make([]grafana.EvalMatches, 0)
	for _, alert := range alerts {
		if alert.State != "ok" && alert.State != "paused" && alert.State != "pending" {
			for _, match := range alert.EvalData.EvalMatches {
				if name, ok := match.Tags["cluster_name"]; ok && name == clusterInfo.Name {
					evalMatches = append(evalMatches, match)
				}
			}
		}
	}
	serviceAlerting := make(map[string]bool, 0)
	for _, eval := range evalMatches {
		if instance, ok := eval.Tags["instance"]; ok {
			if _, okp := eval.Tags["product_name"]; !okp {
				continue
			}
			if _, okp := eval.Tags["service_name"]; !okp {
				continue
			}
			ip := strings.Split(instance, ":")[0]
			serviceStr := eval.Tags["product_name"] + "-" + eval.Tags["service_name"]
			for _, sv := range hostService[ip] {
				if sv == serviceStr {
					if _, sok := serviceAlerting[sv]; !sok {
						serviceAlerting[sv] = true
						result.ServiceAlertingNum++
					}
				}
			}
		}
	}
	return result, err
}

func GetMysqlSalveStatusByGrafana(clusterName string) int {
	ioExpr := MYSQL_IO_RUNNING
	sqlExpr := MYSQL_SQL_RUNNING
	ioExpr = strings.Replace(ioExpr, "{cluster_name}", clusterName, -1)
	sqlExpr = strings.Replace(sqlExpr, "{cluster_name}", clusterName, -1)
	err, ioResponse := grafana.GrafanaQuery(ioExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		return 0
	}
	err, sqlResponse := grafana.GrafanaQuery(sqlExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		return 0
	}
	resList := append(ioResponse.Data.Result, sqlResponse.Data.Result...)
	for _, result := range resList {
		values := result.Values
		for _, value := range values {
			v, _ := value[1].(string)
			if v == "0" {
				return 0
			}
		}
	}

	return 1
}

// GetPlatformGraphConfig 	godoc
// @Summary      	获取图表配置列表
// @Description  	获取图表配置列表
// @Tags         	Inspect
// @Produce      	json
// @Param			cluster_id query int true "Cluster ID"
// @Success      	200		{object} 	[]model.BaseTemplateConfig
// @Router       	/api/v2/platform/inspect/graph/config [get]
func GetPlatformGraphConfig(ctx context.Context) apibase.Result {
	paramsErr := apibase.NewApiParameterErrors()
	clusterId, err := ctx.URLParamInt("cluster_id")
	if err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param cluster_id is empty"))
	}
	clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	if err != nil {
		return fmt.Errorf("[GetPlatformGraphConfig]Get cluster info by cluster_id error: %v ", err)
	}
	name := clusterInfo.Name
	configList, err := model.InspectReportTemplate.GetPlatformTemplateConfig()
	if err != nil {
		log.Errorf("[GetPlatformGraphConfig] get config file %v ", err)
	}
	notHadoopServiceList, err := model.DeployInstanceList.GetServerListNotHadoopById(clusterId)
	if err != nil {
		return fmt.Errorf("[GetPlatformGraphConfig]Query application service list by cluster_id error: %v ", err)
	}
	serviceList := make([]model.InspectServiceList, 0)
	for _, v := range notHadoopServiceList {
		if !strings.Contains(strings.ToLower(v.ServiceName), "sql") &&
			!strings.Contains(strings.ToLower(v.ServiceName), "front") &&
			!strings.Contains(strings.ToLower(v.ServiceName), "kafka") &&
			!strings.Contains(strings.ToLower(v.ServiceName), "zookeeper") {
			serviceList = append(serviceList, v)
		}
	}

	for k, v := range configList {
		if v.Type == 2 && v.Metric == "Full GC Count (2minutes)" {
			configList = append(configList[0:k], configList[k+1:]...)
			for _, sv := range serviceList {
				temp := v
				temp.Module = sv.ServiceName
				temp.Targets = strings.Replace(temp.Targets, "{ProductName}", sv.ProductName, -1)
				temp.Targets = strings.Replace(temp.Targets, "{ServiceName}", sv.ServiceName, -1)
				configList = append(configList, temp)
			}
		}

	}

	sort.Slice(configList, func(i, j int) bool {
		if configList[i].Type == configList[j].Type {
			return strings.ToLower(configList[i].Module) < strings.ToLower(configList[j].Module)
		}
		return configList[i].Type < configList[j].Type
	})
	for k, v := range configList {
		configList[k].Targets = strings.Replace(v.Targets, "{cluster_name}", name, -1)
	}
	return configList
}

type FromStruct struct {
	FormHead  []string   `json:"form_head"`
	FormValue [][]string `json:"form_value"`
}

// GetPlatformInspectFormData 	godoc
// @Summary      	获取巡检报告表格信息
// @Description  	获取巡检报告表格信息
// @Tags         	Inspect
// @Produce      	json
// @Param			cluster_id query int 	true 	"Cluster ID"
// @Param			form_title query string true 	"form_title"
// @Success      	200		{object} 	FromStruct
// @Router       	/api/v2/platform/inspect/form/data [get]
func GetPlatformInspectFormData(ctx context.Context) apibase.Result {

	var result FromStruct
	inspectConfig := cache.SysConfig.InspectConfig
	paramsErr := apibase.NewApiParameterErrors()
	clusterId, err := ctx.URLParamInt("cluster_id")
	if err != nil {
		paramsErr.AppendError("$", fmt.Errorf("param cluster_id is empty"))
	}
	clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	if err != nil {
		return fmt.Errorf("[GetPlatformInspectFormData]Get cluster info by cluster_id error: %v ", err)
	}

	formTitle := ctx.URLParam("form_title")
	if formTitle == "" {
		paramsErr.AppendError("$", fmt.Errorf("param form_title is empty"))
	}
	switch formTitle {
	case SERVICE_GC_FORM_TITLE:
		gcHead := "GC次数（近{time}分钟大于等于{freq}次）"
		gcHead = strings.Replace(gcHead, "{time}", strconv.Itoa(inspectConfig.FullGCTime), 1)
		gcHead = strings.Replace(gcHead, "{freq}", strconv.Itoa(inspectConfig.FullGCFreq), 1)
		result.FormHead = []string{"服务", "所在节点", gcHead}
		result.FormValue = GetServiceGCFormValue(clusterInfo.Name, inspectConfig.FullGCTime, inspectConfig.FullGCFreq)
	case DIR_USAGE_FORM_TITLE:
		gcHead := "目录使用（大于等于{mem}G）"
		gcHead = strings.Replace(gcHead, "{mem}", strconv.Itoa(inspectConfig.DirSize), 1)
		result.FormHead = []string{"节点", gcHead}
		result.FormValue = GetDirUsageFormValue(clusterId, inspectConfig.DirSize)
	case ALARM_FROM_TITLE:
		result.FormHead = []string{"告警指标", "指标归属", "所在节点", "首次告警时间"}
		result.FormValue = GetAlertIndicatorsList(clusterId)
	case HDFS_FILE_USAGE_FORM_TITLE:
		result.FormHead = []string{"节点", "文件存储使用率"}
		result.FormValue = GetFileUsageFormValue(clusterInfo.Name)
	case CPU_USAGE_FROM_TITLE:
		gcHead := "CPU使用率（大于等于{cpu}%）"
		gcHead = strings.Replace(gcHead, "{cpu}", strconv.Itoa(inspectConfig.NodeCPUUsage), 1)
		result.FormHead = []string{"节点", gcHead}
		result.FormValue = GetCPUORMemUsageFormValue(clusterInfo.Name, inspectConfig.NodeCPUUsage, CPU_USAGE_FROM_TITLE)
	case MEM_USAGE_FROM_TITLE:
		gcHead := "内存使用率（大于等于{mem}%）"
		gcHead = strings.Replace(gcHead, "{mem}", strconv.Itoa(inspectConfig.NodeMEMUsage), 1)
		result.FormHead = []string{"节点", gcHead}
		result.FormValue = GetCPUORMemUsageFormValue(clusterInfo.Name, inspectConfig.NodeMEMUsage, MEM_USAGE_FROM_TITLE)
	case DISK_USAGE_FROM_TITLE:
		gcHead := "磁盘使用率（大于等于{disk}%）"
		gcHead = strings.Replace(gcHead, "{disk}", strconv.Itoa(inspectConfig.NodeDiskUsage), 1)
		result.FormHead = []string{"节点", gcHead}
		result.FormValue = GetCPUORMemUsageFormValue(clusterInfo.Name, inspectConfig.NodeDiskUsage, DISK_USAGE_FROM_TITLE)
	case INODE_USAGE_FROM_TITLE:
		gcHead := "inode使用率（大于等于{inode}%）"
		gcHead = strings.Replace(gcHead, "{inode}", strconv.Itoa(inspectConfig.NodeInodeUsage), 1)
		result.FormHead = []string{"节点", gcHead}
		result.FormValue = GetCPUORMemUsageFormValue(clusterInfo.Name, inspectConfig.NodeInodeUsage, INODE_USAGE_FROM_TITLE)
	case SWAP_USAGE_FROM_TITLE:
		result.FormHead = []string{"节点", "swap使用量"}
		result.FormValue = GetSwapUsageFormValue(clusterId)

	}
	return result
}
func GetServiceGCFormValue(name string, gctime, freq int) [][]string {
	serviceGC := make([][]string, 0)
	gcExpr := SERVICE_GC
	gcExpr = strings.Replace(gcExpr, "{cluster_name}", name, -1)
	gcExpr = strings.Replace(gcExpr, "{time}", strconv.Itoa(gctime), 1)
	gcExpr = strings.Replace(gcExpr, "{freq}", strconv.Itoa(freq), 1)
	gcExpr = strings.Replace(gcExpr, "{gc_type}", "G1 Old Generation", 1)
	err, gcResponse := grafana.GrafanaQuery(gcExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("[GetServiceGCFormValue] get G1 Old Generation GC err :%v", err)
		return serviceGC
	}
	resList := gcResponse.Data.Result
	gcExpr = strings.Replace(gcExpr, "G1 Old Generation", "PS MarkSweep", 1)
	err, gcResponse = grafana.GrafanaQuery(gcExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("[GetServiceGCFormValue] get PS MarkSweep GC err :%v", err)
		return serviceGC
	}
	resList = append(resList, gcResponse.Data.Result...)
	resultMap := make(map[string]int, 0)
	for _, v := range resList {
		ip, sv := "", ""
		if instance, ok := v.Metric["instance"]; ok {
			ip = strings.Split(instance.(string), ":")[0]
			if service, sok := v.Metric["service_name"]; sok {
				sv = service.(string)
			}
		}
		numStr, _ := v.Values[0][1].(string)
		num, _ := strconv.Atoi(numStr)
		resultMap[ip+"-"+sv] += num

	}
	type ServiceGC struct {
		IP      string `json:"ip"`
		Service string `json:"service"`
		Num     int    `json:"num"`
	}
	var resultSlice []ServiceGC
	for k, v := range resultMap {
		service, ip := strings.Split(k, "-")[1], strings.Split(k, "-")[0]
		resultSlice = append(resultSlice, ServiceGC{Service: service, IP: ip, Num: v})
	}
	sort.Slice(resultSlice, func(i, j int) bool {
		return resultSlice[i].Num > resultSlice[j].Num
	})
	for _, v := range resultSlice {
		serviceGC = append(serviceGC, []string{v.Service, v.IP, strconv.Itoa(v.Num)})
	}
	return serviceGC
}

func GetDirUsageFormValue(id int, config int) [][]string {
	formValue := make([][]string, 0)
	duCMD := "#!/bin/sh\n du -sm /opt/dtstack"
	hostList, err := model.DeployClusterHostRel.GetInspectClusterHostRelList(id)
	if err != nil {
		log.Errorf("[GetDirUsageFormValue] get host list err:%v", err)
		return formValue
	}

	for k, v := range hostList {
		content, err := agent.AgentClient.ToExecCmdWithTimeout(v.Sid, "", duCMD, "15s", "", "")
		if err != nil {
			log.Errorf("[GetDirUsageFormValue] exec cmd err:%v", err)
			return formValue
		}
		sizeStr := strings.Replace(content, "\t/opt/dtstack\n", "", -1)
		size, _ := strconv.ParseFloat(sizeStr, 64)
		hostList[k].DirSize = size
	}
	sort.Slice(hostList, func(i, j int) bool {
		return hostList[i].DirSize > hostList[j].DirSize
	})
	for _, v := range hostList {
		if v.DirSize/1024 < float64(config) {
			break
		}
		size := strconv.FormatFloat(v.DirSize/1024, 'f', 2, 64)
		size = size + "G"
		formValue = append(formValue, []string{v.IP, size})
	}
	return formValue
}

func GetAlertIndicatorsList(id int) [][]string {
	formValue := make([][]string, 0)
	alertParam := make(map[string]string, 0)
	hostList, err := model.DeployClusterHostRel.GetInspectClusterHostRelList(id)
	if err != nil {
		log.Errorf("[GetAlertIndexList]get host list error: %v", err)
		return formValue
	}
	IPMap := make(map[string]bool, 0)
	for _, v := range hostList {
		IPMap[v.IP] = true
	}
	err, dashboardResp := grafana.GetDashboardByUid("Ne_roaViz")
	if err != nil {
		log.Errorf("[GetAlertIndexList]get host overview dashboard error: %v", err)
		return formValue
	}
	allAlerts, err := GetAlertList(alertParam)

	alertParam = map[string]string{
		"dashboardId": strconv.Itoa(dashboardResp.Dashboard.Id),
	}
	hostAlerts, err := GetAlertList(alertParam)
	for _, alert := range hostAlerts {
		if alert.State != "ok" && alert.State != "paused" && alert.State != "pending" {
			timeStr, _ := time.Parse(time.RFC3339, alert.NewStateDate)
			for _, match := range alert.EvalData.EvalMatches {
				if instance, ok := match.Tags["instance"]; ok {
					ip := strings.Split(instance, ":")[0]
					if IPMap[ip] {
						formValue = append(formValue, []string{alert.Name, "节点", ip, timeStr.Format(base.TsLayout)})
					}

				}
			}
			if len(alert.EvalData.EvalMatches) == 0 {
				formValue = append(formValue, []string{alert.Name, "节点", "-", timeStr.Format(base.TsLayout)})
			}
		}
	}
	for _, alert := range allAlerts {

		if alert.State != "ok" && alert.State != "paused" && alert.State != "pending" && alert.DashboardUid != "Ne_roaViz" {
			timeStr, _ := time.Parse(time.RFC3339, alert.NewStateDate)
			for _, match := range alert.EvalData.EvalMatches {
				serviceName := ""
				if _, okp := match.Tags["service_name"]; !okp {
					serviceName = "-"
				} else {
					serviceName = match.Tags["service_name"]
				}
				if instance, ok := match.Tags["instance"]; ok {
					ip := strings.Split(instance, ":")[0]
					if IPMap[ip] {
						formValue = append(formValue, []string{alert.Name, serviceName, ip, timeStr.Format(base.TsLayout)})
					}
				}
			}

		}
	}
	return formValue
}

type InspectNodeUsage struct {
	IP    string
	Usage float64
}

func GetFileUsageFormValue(name string) [][]string {
	formValue := make([][]string, 0)
	fileUsageExpr := FILE_USAGE
	fileUsageExpr = strings.Replace(fileUsageExpr, "{cluster_name}", name, -1)
	err, fileResponse := grafana.GrafanaQuery(fileUsageExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("[GetFileUsageFormValue] get node file usage err :%v", err)
		return formValue
	}
	resList := fileResponse.Data.Result
	result := make([]InspectNodeUsage, 0)
	for _, v := range resList {
		ip := ""
		if instance, ok := v.Metric["instance"]; ok {
			ip = strings.Split(instance.(string), ":")[0]
			numStr, _ := v.Values[0][1].(string)
			usage, _ := strconv.ParseFloat(numStr, 64)
			result = append(result, InspectNodeUsage{IP: ip, Usage: usage})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Usage > result[j].Usage
	})
	for _, v := range result {
		usageStr := strconv.FormatFloat(v.Usage, 'f', 2, 64) + "%"
		formValue = append(formValue, []string{v.IP, usageStr})
	}
	return formValue
}

func GetCPUORMemUsageFormValue(name string, config int, formTitle string) [][]string {
	formValue := make([][]string, 0)
	UsageExpr := CPU_USAGE
	switch formTitle {
	case INODE_USAGE_FROM_TITLE:
		UsageExpr = INODE_USAGE
	case DISK_USAGE_FROM_TITLE:
		UsageExpr = DISK_USAGE
	case MEM_USAGE_FROM_TITLE:
		UsageExpr = MEM_USAGE
	}

	UsageExpr = strings.Replace(UsageExpr, "{cluster_name}", name, -1)
	UsageExpr = strings.Replace(UsageExpr, "{config}", strconv.Itoa(config), 1)
	err, Response := grafana.GrafanaQuery(UsageExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("[GetCPUORMemUsageFormValue] get node cpu usage err :%v", err)
		return formValue
	}
	result := make([]InspectNodeUsage, 0)
	resList := Response.Data.Result
	for _, v := range resList {
		ip := ""
		if instance, ok := v.Metric["instance"]; ok {
			ip = strings.Split(instance.(string), ":")[0]
			numStr, _ := v.Values[0][1].(string)
			usage, _ := strconv.ParseFloat(numStr, 64)
			result = append(result, InspectNodeUsage{IP: ip, Usage: usage})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Usage > result[j].Usage
	})
	for _, v := range result {
		usageStr := strconv.FormatFloat(v.Usage, 'f', 2, 64) + "%"
		formValue = append(formValue, []string{v.IP, usageStr})
	}
	return formValue
}

func GetSwapUsageFormValue(id int) [][]string {
	formValue := make([][]string, 0)
	UsageExpr := SWAP_USAGE
	//((sum by(instance)(node_memory_SwapTotal))-(sum by(instance)(node_memory_SwapFree))-(sum by(instance)(node_memory_SwapCached)) )/(8*1024*1024*1024)
	hostList, err := model.DeployClusterHostRel.GetInspectClusterHostRelList(id)
	IPMap := make(map[string]bool, 0)
	for _, v := range hostList {
		IPMap[v.IP] = true
	}

	err, Response := grafana.GrafanaQuery(UsageExpr, int(time.Now().Unix()), int(time.Now().Unix()), 1)
	if err != nil {
		log.Errorf("[GetSwapUsageFormValue] get node cpu usage err :%v", err)
		return formValue
	}
	result := make([]InspectNodeUsage, 0)
	resList := Response.Data.Result
	for _, v := range resList {
		ip := ""
		if instance, ok := v.Metric["instance"]; ok {
			ip = strings.Split(instance.(string), ":")[0]
			if IPMap[ip] {
				numStr, _ := v.Values[0][1].(string)
				usage, _ := strconv.ParseFloat(numStr, 64)
				result = append(result, InspectNodeUsage{IP: ip, Usage: usage})
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Usage > result[j].Usage
	})
	for _, v := range result {
		usageStr := strconv.FormatFloat(v.Usage, 'f', 4, 64) + "G"
		formValue = append(formValue, []string{v.IP, usageStr})
	}
	return formValue
}
