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
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	uuid "github.com/satori/go.uuid"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/group"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"github.com/kataras/iris/context"
)

func GroupStart(ctx context.Context) apibase.Result {
	log.Debugf("[GroupStart] GroupStart from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	groupName := ctx.Params().Get("group_name")

	if groupName == "" {
		log.Errorf("%v", "group_name is empty")
		return fmt.Errorf("group_name is empty")
	}

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	pInfo, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		return err
	}

	err = model.OperationList.Insert(model.OperationInfo{
		ClusterId:       clusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.ProductStart.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Product.Code,
		ObjectValue:     pInfo.ProductName,
	})
	if err != nil {
		log.Errorf("OperationList Insert err:%v", err)
	}
	grouper, err := group.NewGrouper(pid, clusterId, groupName, operationId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	grouper.Start()

	return grouper.GetResult()
}

func GroupStop(ctx context.Context) apibase.Result {
	log.Debugf("[GroupStop] GroupStop from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	groupName := ctx.Params().Get("group_name")

	if groupName == "" {
		log.Errorf("%v", "group_name is empty")
		return fmt.Errorf("group_name is empty")
	}

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	grouper, err := group.NewGrouper(pid, clusterId, groupName, "")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	grouper.Stop()

	return grouper.GetResult()
}
