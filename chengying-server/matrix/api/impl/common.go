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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
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
