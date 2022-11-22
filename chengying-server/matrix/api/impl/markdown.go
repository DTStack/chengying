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
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/markdown"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"github.com/kataras/iris/context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DeployInfoMarkdownPage struct {
	PageTemplate      string
	DeployArchData    []DeployArchInfoStruct
	DeployVersionData []model.DeployVersionStruct
	DeployListData    []DeployListStruct
	IPStr             string
	IPLen             int
}

// DeployInfoGenerate 	godoc
// @Summary      	生成部署信息文档
// @Description  	生成部署信息文档
// @Tags         	common
// @Produce      	json
// @Param			em_version query string true "EM 版本号"
// @Success      	200		{object}  string	"{"msg":"ok","code":0,"data":{}}"
// @Router       	/api/v2/common/deployInfo/generate [post]
func DeployInfoGenerate(ctx context.Context) apibase.Result {
	var reqParams struct {
		EMVersion string `json:"em_version"`
	}
	err := ctx.ReadJSON(&reqParams)
	if err != nil {
		return err
	}

	//获取部署架构信息
	deployArchInfo, err := GetDeployArchInfo()
	if err != nil {
		return err
	}
	//获取部署版本信息
	deployVersionInfo, err := GetDeployVersionInfo(reqParams.EMVersion)
	if err != nil {
		return err
	}
	//获取部署清单信息
	deployListInfo, err := GetDeployListInfo(reqParams.EMVersion)
	ipMap := make(map[string]string, 0)
	for _, v := range deployArchInfo {
		ipMap[v.HostIP] = v.HostName
	}
	deployList, ipHead, ipLen := DealDeployListInfo(deployListInfo, ipMap)
	deployInfoMDData := DeployInfoMarkdownPage{
		DeployArchData:    deployArchInfo,
		DeployVersionData: deployVersionInfo,
		DeployListData:    deployList,
		IPStr:             ipHead,
		IPLen:             ipLen,
	}
	err = deployInfoMDData.RenderAndGenerateDeployInfo()
	return err

}

type DeployArchInfoStruct struct {
	HostIP          string
	HostName        string
	CPU             string
	MemSizeDisplay  string `json:"mem_size_display"`
	DiskSizeDisplay string `json:"disk_size_display"`
	FileSizeDisplay string `json:"file_size_display"`
	OS              string `json:"os_display"`
	Pwd             string `json:"pwd"`
}

// DeployInfoDownload 	godoc
// @Summary      	部署信息文档下载
// @Description  	部署信息文档下载
// @Tags         	common
// @Produce      	json
// @Success      	200		{object} string	"{"msg":"ok","code":0,"data":{}}"
// @Router       	/api/v2/common/deployInfo/download [get]
func DeployInfoDownload(ctx context.Context) apibase.Result {
	defer func() {
		if err := os.RemoveAll(filepath.Join(base.MDDir, "集群部署信息.md")); err != nil {
			log.Errorf("Remove md file error: %v", err)
		}
	}()
	filePath := filepath.Join(base.MDDir, "集群部署信息.md")
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("[DeployInfoDownload] open report file error")
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	fileName := url.QueryEscape("集群部署信息")
	ctx.Header("Content-Disposition", "attachment;filename=\""+fileName+".md\"")
	ctx.Write(content)
	return apibase.EmptyResult{}
}

func GetDeployArchInfo() ([]DeployArchInfoStruct, error) {
	var result []DeployArchInfoStruct
	hostsList, err := model.DeployClusterProductRel.GetDeployArchInfo()
	if err != nil {
		return nil, err
	}
	for i, list := range hostsList {
		hostsList[i].MemSizeDisplay, _ = MultiSizeConvert(list.MemSize, 0)
		hostsList[i].CpuCoreSizeDisplay = strconv.Itoa(list.CpuCores) + "core"
		if list.DiskUsage.Valid {
			hostsList[i].DiskSizeDisplay, _, hostsList[i].FileSizeDisplay, _ = diskUsageConvert(list.DiskUsage.String)
		}
		result = append(result, DeployArchInfoStruct{
			HostIP:          hostsList[i].IP,
			HostName:        hostsList[i].HostName,
			CPU:             hostsList[i].CpuCoreSizeDisplay,
			MemSizeDisplay:  hostsList[i].MemSizeDisplay,
			DiskSizeDisplay: hostsList[i].DiskSizeDisplay,
			FileSizeDisplay: hostsList[i].FileSizeDisplay,
			OS:              hostsList[i].OSDisplay,
			Pwd:             "",
		})
	}
	return result, nil
}

func GetDeployVersionInfo(EMVersion string) ([]model.DeployVersionStruct, error) {
	result, err := model.DeployClusterProductRel.GetDeployVersionInfo()
	if err != nil {
		return nil, err
	}
	result = append([]model.DeployVersionStruct{{
		ProductName:    "EasyManager",
		ProductVersion: EMVersion}}, result...)
	return result, nil
}

func GetDeployListInfo(EMVersion string) ([]model.DeployListStruct, error) {
	result, err := model.DeployInstanceList.GetDeployListInfo()
	if err != nil {
		return nil, err
	}
	result = append([]model.DeployListStruct{{
		ProductName:    "EasyManager",
		ProductVersion: EMVersion,
		ServiceName:    "",
		ServiceVersion: "",
		IPs:            ""}}, result...)
	return result, nil
}

type DeployListStruct struct {
	ProductName    string
	ProductVersion string
	ServiceName    string
	ServiceVersion string
	DeployIPs      map[string]bool
	IPStr          string
}

func DealDeployListInfo(deployListInfo []model.DeployListStruct, ipMap map[string]string) ([]DeployListStruct, string, int) {
	deployList := make([]DeployListStruct, 0)
	ipsStr := ""
	ipExist := make(map[string]bool, 0)
	var ipsList []string
	for _, v := range deployListInfo {
		tmpIPsArr := strings.Split(v.IPs, "/")
		for _, tmpV := range tmpIPsArr {
			if _, ok := ipExist[tmpV]; !ok && tmpV != "" {
				ipAndName := tmpV
				if tmpHostName, ok1 := ipMap[tmpV]; ok1 {
					ipAndName = ipAndName + "/" + tmpHostName
				}
				ipsList = append(ipsList, tmpV)
				ipExist[tmpV] = true
				if ipsStr != "" {
					ipsStr = ipsStr + "|" + ipAndName
				} else {
					ipsStr = ipAndName
				}
			}
		}
		deployIps := make(map[string]bool, 0)
		for _, ip := range tmpIPsArr {
			deployIps[ip] = true
		}
		deployList = append(deployList, DeployListStruct{
			ProductName:    v.ProductName,
			ProductVersion: v.ProductVersion,
			ServiceName:    v.ServiceName,
			ServiceVersion: v.ServiceVersion,
			DeployIPs:      deployIps,
			IPStr:          "",
		})
	}
	ipLen := len(ipsList)
	for di, d := range deployList {
		ipFlagStr := "|"
		for v1, ip := range ipsList {
			if d.DeployIPs[ip] == true {
				ipFlagStr = ipFlagStr + markdown.DeploySelectFlag + "|"
			} else {
				if v1 < ipLen {
					ipFlagStr = ipFlagStr + "|"
				}
			}
		}
		ipFlagStr = ipFlagStr[1 : len(ipFlagStr)-1]
		deployList[di].IPStr = ipFlagStr
	}
	return deployList, ipsStr, ipLen
}

func (d *DeployInfoMarkdownPage) RenderAndGenerateDeployInfo() error {
	pageTemplate := markdown.DeployInfoPage
	deployArchStr := ""
	if len(d.DeployArchData) != 0 {
		deployArchStr = d.RenderDeployArchTable()
	}
	pageTemplate = strings.Replace(pageTemplate, "{deployArchTable}", deployArchStr, 1)

	deployVersionStr := ""
	if len(d.DeployVersionData) != 0 {
		deployVersionStr = d.RenderDeployVersionTable()
	}
	pageTemplate = strings.Replace(pageTemplate, "{deployVersionTable}", deployVersionStr, 1)

	deployListStr := ""
	if len(d.DeployListData) != 0 {
		deployListStr = d.RenderDeployListTable()
	}
	pageTemplate = strings.Replace(pageTemplate, "{deployListTable}", deployListStr, 1)
	filePath := filepath.Join(base.MDDir, "集群部署信息.md")
	_ = os.Remove(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(pageTemplate))
	if err != nil {
		return fmt.Errorf("generate markdown file error :%v", err)
	}
	return nil
}

func (d *DeployInfoMarkdownPage) RenderDeployArchTable() string {
	template := markdown.DeployArchTable
	params := ""
	for _, v := range d.DeployArchData {
		tl := markdown.DeployArchParam
		tl = strings.Replace(tl, "{hostIP}", v.HostIP, 1)
		tl = strings.Replace(tl, "{hostName}", v.HostName, 1)
		tl = strings.Replace(tl, "{CPU}", v.CPU, 1)
		tl = strings.Replace(tl, "{MEM}", v.MemSizeDisplay, 1)
		tl = strings.Replace(tl, "{dataDisk}", v.DiskSizeDisplay, 1)
		tl = strings.Replace(tl, "{systemDisk}", v.FileSizeDisplay, 1)
		tl = strings.Replace(tl, "{operatorSys}", v.OS, 1)
		params = fmt.Sprintf("%s%s", params, tl)
	}
	template = strings.Replace(template, "{params}", params, 1)
	return template
}

func (d *DeployInfoMarkdownPage) RenderDeployVersionTable() string {
	template := markdown.DeployVersionTable
	params := ""
	for _, v := range d.DeployVersionData {
		tl := markdown.DeployVersionParam
		tl = strings.Replace(tl, "{product}", v.ProductName, 1)
		tl = strings.Replace(tl, "{version}", v.ProductVersion, 1)
		params = fmt.Sprintf("%s%s", params, tl)
	}
	template = strings.Replace(template, "{params}", params, 1)
	return template
}

func (d *DeployInfoMarkdownPage) RenderDeployListTable() string {
	template := markdown.DeployListTable
	if d.IPLen > 0 {
		template = strings.Replace(template, "{ip/host}", d.IPStr, 1)
		arStr := "" //markdown 动态列使用，动态扩充
		for i := 0; i < d.IPLen; i++ {
			arStr = arStr + "--------------|"
		}
		template = strings.Replace(template, "--------------|", arStr, 1)
	}
	params := ""
	for _, v := range d.DeployListData {
		tl := markdown.DeployListParam
		tl = strings.Replace(tl, "{product}", v.ProductName, 1)
		tl = strings.Replace(tl, "{version}", v.ProductVersion, 1)
		tl = strings.Replace(tl, "{service}", v.ServiceName, 1)
		tl = strings.Replace(tl, "{service_version}", v.ServiceVersion, 1)
		tl = strings.Replace(tl, "{ips/hosts}", v.IPStr, 1)
		params = fmt.Sprintf("%s%s", params, tl)
	}
	template = strings.Replace(template, "{params}", params, 1)
	return template
}
