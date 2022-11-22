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
	"fmt"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"github.com/kataras/iris/context"
)

func ForceStop(ctx context.Context) apibase.Result {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		return err
	}

	info, err := model.DeployInstanceRecord.GetDeployInstanceRecordById(id)
	if err != nil {
		return err
	}

	if info.Status != model.INSTANCE_STATUS_STOP_FAIL {
		return fmt.Errorf("instance status %v can't force stop", info.Status)
	}

	err, instance := model.DeployInstanceList.GetInstanceInfoById(info.InstanceId)
	if err != nil {
		return err
	}
	if instance.HealthState == model.INSTANCE_HEALTH_NOTSET {
		err = model.DeployInstanceList.UpdateInstanceStatusById(info.InstanceId, model.INSTANCE_STATUS_STOPPED, "force stop")
	} else {
		err = model.DeployInstanceList.UpdateInstanceStatusById(info.InstanceId, model.INSTANCE_STATUS_STOPPED, "force stop", model.INSTANCE_HEALTH_WAITING)
	}
	if err != nil {
		return err
	}

	return model.DeployInstanceRecord.UpdateDeployInstanceRecord(id, dbhelper.UpdateFields{
		"status":         model.INSTANCE_STATUS_STOPPED,
		"status_message": "force stop",
		"progress":       30,
		"update_time":    time.Now(),
	})
}

func ForceUninstall(ctx context.Context) apibase.Result {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		return err
	}

	info, err := model.DeployInstanceRecord.GetDeployInstanceRecordById(id)
	if err != nil {
		return err
	}

	if info.Status != model.INSTANCE_STATUS_UNINSTALL_FAIL {
		return fmt.Errorf("instance status %v can't force uninstall", info.Status)
	}

	if err = model.DeployInstanceList.DeleteByInstanceId(info.InstanceId); err != nil {
		return err
	}

	return model.DeployInstanceRecord.UpdateDeployInstanceRecord(id, dbhelper.UpdateFields{
		"status":         model.INSTANCE_STATUS_UNINSTALLED,
		"status_message": "force uninstall",
		"progress":       100,
		"update_time":    time.Now(),
	})
}
