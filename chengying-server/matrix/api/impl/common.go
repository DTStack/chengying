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
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"encoding/json"
	"fmt"
	enhanceyaml "github.com/ghodss/yaml"
	"github.com/kataras/iris/context"
	"io/ioutil"
	"strconv"
)

func File2text(ctx context.Context) apibase.Result {
	f, _, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return string(content)
}

func GetCurrentClusterFromParam(ctx context.Context) (int, error) {
	id := ctx.FormValue("clusterId")
	return strconv.Atoi(id)
}

func GetCurrentClusterId(ctx context.Context) (int, error) {
	id := ctx.GetCookie(COOKIE_CURRENT_CLUSTER_ID)
	if id == "" {
		return -1, fmt.Errorf("未找到集群信息: %v", COOKIE_CURRENT_CLUSTER_ID)
	}
	cid, err := strconv.Atoi(id)
	if err != nil {
		return -1, fmt.Errorf("未找到集群信息: %v", err.Error())
	}
	return cid, nil
}

func GetCurrentParentProduct(ctx context.Context) (string, error) {
	parentProduct := ctx.GetCookie(COOKIE_PARENT_PRODUCT_NAME)
	if parentProduct == "" {
		return "", fmt.Errorf("未找到父产品信息: %v", COOKIE_PARENT_PRODUCT_NAME)
	}
	return parentProduct, nil
}

func GetMetaLableFromKubeObj(obj interface{}, label string) string {
	return ""
}

func GetSafetyAuditModule(ctx context.Context) apibase.Result {
	log.Debugf("[GetSafetyAuditModule] GetSafetyAuditModule from EasyMatrix API ")

	list, err := model.SafetyAuditList.GetSafetyAuditModuleList()
	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"count": 0,
			"list":  nil,
		}
	} else if err != nil {
		return err
	}

	return map[string]interface{}{
		"count": len(list),
		"list":  list,
	}
}

func GetSafetyAuditOperation(ctx context.Context) apibase.Result {
	log.Debugf("[GetSafetyAuditOperation] GetSafetyAuditOperation from EasyMatrix API ")

	module := ctx.URLParam("module")
	if module == "" {
		return fmt.Errorf("请先选择模块\n")
	}

	list, err := model.SafetyAuditList.GetSafetyAuditOperationList(module)
	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"count": 0,
			"list":  nil,
		}
	} else if err != nil {
		return err
	}

	return map[string]interface{}{
		"count": len(list),
		"list":  list,
	}
}

func GetSafetyAuditList(ctx context.Context) apibase.Result {
	log.Debugf("[GetSafetyAuditList] GetSafetyAuditList from EasyMatrix API ")

	operator := sqlTransfer(ctx.URLParam("operator"))
	ip := sqlTransfer(ctx.URLParam("ip"))
	operation := ctx.URLParam("operation")
	module := ctx.URLParam("module")
	content := sqlTransfer(ctx.URLParam("content"))
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.SafetyAuditInfo{})
	from := ctx.URLParam("from")
	to := ctx.URLParam("to")

	list, count, err := model.SafetyAuditList.SelectSafetyAuditListByWhere(pagination, module, operation, operator, ip, content, from, to)
	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"count": 0,
			"list":  nil,
		}
	} else if err != nil {
		return err
	}

	return map[string]interface{}{
		"count": count,
		"list":  list,
	}

}

func addSafetyAuditRecord(ctx context.Context, module, operation, content string) error {
	return model.SafetyAuditList.InsertSafetyAuditRecord(ctx.GetCookie("em_username"), module, operation, ctx.RemoteAddr(), content)
}

// WorkloadDefinaInit
// @Description  	workload加载接口
// @Summary      	workload加载接口
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/workloadinit [get]
func WorkloadDefinaInit(ctx context.Context) (rlt apibase.Result) {
	var msg string
	// 获取目录下workload yaml文件到切片中
	files := []string{}

	workloadfile_path := base.WebRoot + "/workload-definition"
	wkfiles, err := ioutil.ReadDir(workloadfile_path)
	if err != nil {
		log.Errorf("open file error: %v", err)
		return err
	}
	for _, f := range wkfiles {
		if f.IsDir() {
			continue
		}
		files = append(files, workloadfile_path+"/"+f.Name())
	}

	// 遍历workload yaml文件进行解析入库
	if len(files) != 0 {
		for _, file := range files {
			conf := new(modelkube.WorkloadDefinitionYaml)
			tx := model.USE_MYSQL_DB().MustBegin()

			yamlFile, err := ioutil.ReadFile(file)
			if err != nil {
				log.Errorf("[WorkloadDefinaInit] read workload yaml %v err:%v", file, err)
				return fmt.Errorf("[WorkloadDefinaInit] read workload yaml %v err:%v", file, err)
			}

			js, err := enhanceyaml.YAMLToJSON(yamlFile)
			if err != nil {
				log.Errorf("[WorkloadDefinaInit] %v to json err: %v\n", file, err)
				return fmt.Errorf("[WorkloadDefinaInit] %v to json err: %v\n", file, err)
			}

			err = json.Unmarshal(js, conf)
			if err != nil {
				log.Errorf("[WorkloadDefinaInit] unmarshal json error %v:", err)
				return fmt.Errorf("[WorkloadDefinaInit] unmarshal json error %v:", err)
			}
			// workload definition表用到的参数
			wkdef_version := conf.ApiVersion
			wkdef_name := conf.Metadata.Name
			wkdef_params, err := json.Marshal(conf.Spec.Params)
			if err != nil {
				log.Errorf("[WorkloadDefinaInit] marshal workloaddefinition_params json error %v:", err)
				return fmt.Errorf("[WorkloadDefinaInit] marshal workloaddefinition_params json error %v:", err)
			}
			// workload part表用到的参数
			wkpart_name := conf.Spec.WorkloadPatrs[0].Baseworkload.Name
			wkpart_type := conf.Spec.WorkloadPatrs[0].Baseworkload.Type
			wkpart_parameters, err := json.Marshal(conf.Spec.WorkloadPatrs[0].Baseworkload.Parameters)
			if err != nil {
				log.Errorf("[WorkloadDefinaInit] marshal workloadpart_parameters json error %v:", err)
				return fmt.Errorf("[WorkloadDefinaInit] marshal workloadpart_parameters json error %v:", err)
			}
			// workload definition表处理
			workload_id, err := modelkube.WorkloadDefinition.InsertOrUpdate(tx, wkdef_name, wkdef_version, string(wkdef_params))
			if err != nil {
				log.Errorf("[WorkloadDefinaInit] insert or update workload definifion error:%v", err)
				return fmt.Errorf("[WorkloadDefinaInit] insert or update workload definifion error:%v", err)
			}
			// workload part表处理
			workloadPartId, err := modelkube.WorkloadPart.InsertOrUpdate(tx, wkpart_name, wkpart_type, string(wkpart_parameters), workload_id)
			if err != nil {
				log.Errorf("[WorkloadDefinaInit] insert or update workload part error:%v", err)
				return fmt.Errorf("[WorkloadDefinaInit] insert or update workload part error:%v", err)
			}

			// workload step表处理
			for _, step := range conf.Spec.WorkloadPatrs[0].Steps {
				wkstep_name := step.Name
				wkstep_type := step.Type
				wkstep_action := step.Action
				wkstep_object, err := json.Marshal(step.Object)
				if err != nil {
					log.Errorf("[WorkloadDefinaInit] marshal workload_step object json error %v:", err)
					return err
				}
				err = modelkube.WorkloadStep.InsertOrUpdate(tx, wkstep_name, wkstep_type, wkstep_action, string(wkstep_object), workloadPartId)
				if err != nil {
					log.Errorf("[WorkloadDefinaInit] insert or update workload step error:%v", err)
					return fmt.Errorf("[WorkloadDefinaInit] insert or update workload step error:%v", err)
				}
			}
			if err := tx.Commit(); err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		msg = workloadfile_path + "目录下不存在workload定义文件！"
	}

	return msg
}
