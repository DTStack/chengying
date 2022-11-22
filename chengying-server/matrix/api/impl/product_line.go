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
	"database/sql"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/dag"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/context"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	jsonExt = ".json"
)

var productLineLock sync.Mutex

type ProductLineInfoResult struct {
	ID                 int             `json:"id"`
	ProductLineName    string          `json:"product_line_name"`
	ProductLineVersion string          `json:"product_line_version"`
	CreateTime         string          `json:"create_time"`
	UpdateTime         string          `json:"update_time"`
	DeployProcess      []DeployProcess `json:"deploy_process"`
}

type DeployProcess struct {
	model.ProductSerial
	IsExist bool `json:"is_exist"`
}

// ProductLineInfo
// @Description  	GET Product Line Info
// @Summary      	获取所有产品线信息接口
// @Tags         	product_line
// @Accept          application/json
// @Produce 		application/json
// @Param           product_type query  string  false  "产品类型"
// @Success         200  {object} string "{"msg":"ok","code":0,"data":{"list":"","count":""}}"
// @Router          /api/v2/product_line [get]
func ProductLineInfo(ctx context.Context) apibase.Result {
	log.Debugf("[ProductLine->ProductLineInfo] ProductLineInfo from EasyMatrix API ")

	type res struct {
		List  []ProductLineInfoResult `json:"list"`
		Count int                     `json:"count"`
	}
	result := make([]ProductLineInfoResult, 0)
	productType := ctx.URLParam("product_type")
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.DeployProductLineInfo{})
	list, total := model.DeployProductLineList.GetProductLineList(pagination)
	for _, info := range list {
		m := ProductLineInfoResult{}
		m.ID = info.ID
		m.ProductLineName = info.ProductLineName
		m.ProductLineVersion = info.ProductLineVersion
		if info.CreateTime.Valid {
			m.CreateTime = info.CreateTime.Time.Format(base.TsLayout)
		}
		if info.UpdateTime.Valid {
			m.UpdateTime = info.UpdateTime.Time.Format(base.TsLayout)
		}
		//DeployProcess
		deploySerial, err := GetProductLineDeploySerial(info)
		if err != nil {
			log.Errorf("[ProductLine->ProductLineInfo] GetProductLineDeploySerial error: %v", err)
			return err
		}
		for _, p := range deploySerial {
			dep := DeployProcess{}
			dep.ID = p.ID
			dep.ProductName = p.ProductName
			dep.Dependee = p.Dependee
			productInfoList, _ := model.DeployProductList.GetProductList(p.ProductName, productType, nil, nil)
			if len(productInfoList) != 0 {
				dep.IsExist = true
			}
			m.DeployProcess = append(m.DeployProcess, dep)
		}
		result = append(result, m)
	}

	return res{
		List:  result,
		Count: total,
	}
}

func GetProductLineDeploySerial(productLineInfo model.DeployProductLineInfo) ([]model.ProductSerial, error) {
	serials := make([]model.ProductSerial, 0)
	if err := json.Unmarshal(productLineInfo.ProductSerial, &serials); err != nil {
		log.Errorf("json unmarshal error: %v", err)
		return nil, err
	}
	nodes := make([]dag.Node, 0)
	edges := make([]dag.Edge, 0)
	serialMap := map[dag.Node]model.ProductSerial{}
	for _, serial := range serials {
		serialMap[serial.ID] = serial
		nodes = append(nodes, serial.ID)
	}
	for _, serial := range serialMap {
		if _, ok := serialMap[serial.Dependee]; ok {
			edges = append(edges, dag.Edge{
				Depender: serial.ID,
				Dependee: serial.Dependee,
			})
		}
	}
	graph := &dag.Graph{
		Nodes: nodes,
		Edges: edges,
	}
	deploySerial := make([]model.ProductSerial, 0)
	if err := dag.Execute(graph, func(node dag.Node) error {
		deploySerial = append(deploySerial, model.ProductSerial{
			ID:          serialMap[node].ID,
			ProductName: serialMap[node].ProductName,
			Dependee:    serialMap[node].Dependee,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return deploySerial, nil
}

type UploadProductLineResult struct {
	ID                 int    `json:"id"`
	ProductLineName    string `json:"product_line_name"`
	ProductLineVersion string `json:"product_line_version"`
	ProductSerial      string `json:"product_serial"`
	CreateTime         string `json:"create_time"`
}

// UploadProductLine
// @Description  	Upload Product Line
// @Summary      	上传产品线接口
// @Tags         	product_line
// @Accept          application/json
// @Produce 		application/json
// @Param           file body string true "-F 'file=@DTBatch产品线.json'"
// @Success         200  {object} string "{"msg": "ok","code": 0,"data": {"id": "","product_line_name": "","product_line_version": "","product_serial": "","create_time": ""}}"
// @Router          /api/v2/product_line/upload [post]
func UploadProductLine(ctx context.Context) apibase.Result {
	log.Debugf("[ProductLine->UploadProductLine] UploadProductLine from EasyMatrix API ")

	productLineLock.Lock()
	defer productLineLock.Unlock()
	file, head, err := ctx.FormFile("file")
	if errors.Is(err, http.ErrMissingFile) {
		return fmt.Errorf("请上传产品线")
	} else if err != nil {
		return err
	}
	defer file.Close()

	//限制文件后缀
	fileExt := path.Ext(head.Filename)
	if fileExt != jsonExt {
		return fmt.Errorf("仅支持 %v 格式文件", jsonExt)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("[ProductLine->UploadProductLine] read json file err: %v", err)
		return err
	}
	//解析文件
	var template model.ProductLineTemplate
	if err = json.Unmarshal(content, &template); err != nil {
		log.Errorf("[ProductLine->UploadProductLine] json %s unmarshal error: %v", head.Filename, err)
		return fmt.Errorf("解析文件内容错误, err: %v", err)
	}
	if err := unmarshalValidation(&template); err != nil {
		return fmt.Errorf("解析文件内容错误, err: %v", err)
	}
	productSerialBytes, _ := json.Marshal(template.ProductSerial)
	err, id := model.DeployProductLineList.InsertProductLineIfNotExist(template.ProductLineName, template.ProductLineVersion, productSerialBytes)
	if err != nil {
		log.Errorf("[ProductLine->UploadProductLine] insert product line err: %v", err)
		return fmt.Errorf("上传失败，%v", err)
	}
	result := UploadProductLineResult{
		ID:                 id,
		ProductLineName:    template.ProductLineName,
		ProductLineVersion: template.ProductLineVersion,
		ProductSerial:      string(productSerialBytes),
		CreateTime:         time.Now().Format(base.TsLayout),
	}

	return result
}

func unmarshalValidation(v interface{}) error {
	var err error
	fields := reflect.ValueOf(v).Elem()
	for i := 0; i < fields.NumField(); i++ {
		if fields.Type().Field(i).Type.Kind() == reflect.Slice {
			serials := fields.Field(i).Interface().([]model.ProductSerial)
			serialMap := map[dag.Node]model.ProductSerial{}
			for _, serial := range serials {
				if serial.ID <= 0 {
					err = fmt.Errorf("the 'id' field is invalid")
					return err
				}
				serialMap[serial.ID] = serial
			}
			for _, serial := range serialMap {
				if serial.Dependee != 0 {
					if _, ok := serialMap[serial.Dependee]; !ok {
						err = fmt.Errorf("the 'dependee' field is invalid")
						return err
					}
				}
				err = unmarshalValidation(&serial)
				if err != nil {
					return err
				}
			}
		}
		labelRequired := fields.Type().Field(i).Tag.Get("validation")
		labelValue := fields.Type().Field(i).Tag.Get("json")
		if strings.Contains(labelRequired, "required") && fields.Field(i).IsZero() {
			err = fmt.Errorf("the '%s' field is required", labelValue)
			return err
		}
	}
	return err
}

// DeleteProductLine
// @Description  	Delete Product Line
// @Summary      	删除产品线接口
// @Tags         	product_line
// @Accept          application/json
// @Produce 		application/json
// @Param           id query  int  false  "产品线id"
// @Success         200  {object} string "{"msg": "ok","code": 0,"data": null}}"
// @Router          /api/v2/product_line/{id:int(primary key)} [delete]
func DeleteProductLine(ctx context.Context) apibase.Result {
	log.Debugf("[ProductLine->DeleteProductLine] DeleteProductLine from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	idStr := ctx.Params().Get("id")
	if idStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	if err := model.DeployProductLineList.DeleteProductLineById(id); err != nil {
		log.Errorf("[ProductLine->DeleteProductLine] delete product line err: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("id: %v, is not exist", id)
		}
		return err
	}

	return nil
}

// ProductListOfProductLine
// @Description  	Get Product List
// @Summary      	产品包列表接口
// @Tags         	product_line
// @Accept          application/json
// @Produce 		application/json
// @Param           product_line_name query  string  false  "产品线名称"
// @Param          product_line_version query  string  false  "产品线版本"
// @Param           product_type query  string  false  "产品类型"
// @Param           deploy_status query  string  false  "部署状态"
// @Success         200  {object} string "{"msg":"ok","code":0,"data":{"list":"","count":""}}"
// @Router          /api/v2/product_line/product_list [get]
func ProductListOfProductLine(ctx context.Context) apibase.Result {
	log.Debugf("[ProductLine->ProductListOfProductLine] ProductListOfProductLine from EasyMatrix API ")

	var deployStatus []string
	productLineName := ctx.URLParam("product_line_name")
	productLineVersion := ctx.URLParam("product_line_version")
	productType := ctx.URLParam("product_type")
	if status := ctx.URLParam("deploy_status"); status != "" {
		deployStatus = strings.Split(status, ",")
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	type ProductChildren struct {
		ID                 int    `json:"id"`
		ParentProductName  string `json:"parent_product_name"`
		ProductName        string `json:"product_name"`
		ProductNameDisplay string `json:"product_name_display"`
		ProductVersion     string `json:"product_version"`
		Status             string `json:"status"`
		DeployUUID         string `json:"deploy_uuid"`
		DeployTime         string `json:"deploy_time"`
		CreateTime         string `json:"create_time"`
		ProductType        int    `json:"product_type"`
		Namespace          string `json:"namespace"`
		IsDefault          bool   `json:"is_default"`
	}
	type DeployProductList struct {
		ProductName string            `json:"product_name"`
		Status      string            `json:"status"`
		Children    []ProductChildren `json:"children"`
	}
	resultList := make([]DeployProductList, 0)
	productMap := map[string][]model.DeployProductListInfo{}
	if productLineName != "" && productLineVersion != "" {
		info, err := model.DeployProductLineList.GetProductLineListByNameAndVersion(productLineName, productLineVersion)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[ProductLine->ProductListOfProductLine] get product line err: %v", err)
			return err
		}
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("产品线 `%v(%v)` 不存在", productLineName, productLineVersion)
		}
		serials := make([]model.ProductSerial, 0)
		if err := json.Unmarshal(info.ProductSerial, &serials); err != nil {
			log.Errorf("[ProductLine->ProductListOfProductLine] json unmarshal error: %v", err)
			return err
		}
		temp := map[string]struct{}{}
		for _, serial := range serials {
			if _, ok := temp[serial.ProductName]; !ok {
				temp[serial.ProductName] = struct{}{}
				productInfoList, _ := model.DeployProductList.GetProductList(serial.ProductName, productType, nil, nil)
				for _, productInfo := range productInfoList {
					productMap[productInfo.ProductName] = append(productMap[productInfo.ProductName], productInfo)
				}
			}
		}
	} else if productLineName == "" && productLineVersion == "" {
		productInfoList, _ := model.DeployProductList.GetProductList("", productType, nil, nil)
		for _, productInfo := range productInfoList {
			productMap[productInfo.ProductName] = append(productMap[productInfo.ProductName], productInfo)
		}
	} else {
		return fmt.Errorf("product_line_name or product_line_version is empty")
	}
	for productName, productInfoList := range productMap {
		m := DeployProductList{}
		childrenList := make([]ProductChildren, 0)
		m.ProductName = productName
		bFound := false
		for _, productInfo := range productInfoList {
			children := ProductChildren{}
			children.ID = productInfo.ID
			children.ParentProductName = productInfo.ParentProductName
			children.ProductName = productInfo.ProductName
			children.ProductNameDisplay = productInfo.ProductNameDisplay
			children.ProductVersion = productInfo.ProductVersion
			children.ProductType = productInfo.ProductType
			if productInfo.CreateTime.Valid {
				children.CreateTime = productInfo.CreateTime.Time.Format(base.TsLayout)
			}
			installed, err := model.DeployClusterProductRel.GetByPidAndClusterId(productInfo.ID, clusterId)
			if err == nil {
				children.Status = installed.Status
				children.DeployUUID = installed.DeployUUID
				children.Namespace = installed.Namespace
				children.IsDefault = true
				bFound = true
				m.Status = installed.Status
				if installed.DeployTime.Valid {
					children.DeployTime = installed.DeployTime.Time.Format(base.TsLayout)
				}
			} else {
				children.Status = model.PRODUCT_STATUS_UNDEPLOYED
			}
			childrenList = append(childrenList, children)
		}
		m.Children = childrenList
		if len(m.Children) > 0 {
			// 默认按照版本号大小由前至后排列
			sort.SliceStable(m.Children, func(i, j int) bool {
				return m.Children[i].ProductVersion > m.Children[j].ProductVersion
			})
			//全部没有部署，默认选中最新版本，状态为未部署
			if !bFound {
				m.Children[0].IsDefault = true
				m.Status = model.PRODUCT_STATUS_UNDEPLOYED
			}
			//部署状态筛选
			if len(deployStatus) == 0 {
				resultList = append(resultList, m)
			} else {
				for _, status := range deployStatus {
					if m.Status == status {
						resultList = append(resultList, m)
					}
				}
			}
		}
	}
	// 默认以组件名排序
	sort.SliceStable(resultList, func(i, j int) bool {
		return strings.Compare(resultList[i].ProductName, resultList[j].ProductName) == -1
	})

	return map[string]interface{}{
		"list":  resultList,
		"count": len(resultList),
	}
}
