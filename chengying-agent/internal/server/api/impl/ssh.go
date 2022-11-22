/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package impl

import (
	apibase "easyagent/go-common/api-base"
	"easyagent/internal/server/log"
	sshs "easyagent/internal/server/ssh"
	. "easyagent/internal/server/tracy"

	"github.com/kataras/iris/context"
)

func CheckByUserPwd(ctx context.Context) apibase.Result {
	InstallProgressLog("[INSTALL] CheckByUserPwd ...%v", "")
	paramErrs := apibase.NewApiParameterErrors()

	params := struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	}{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		InstallProgressLog("[INSTALL] CheckByUserPwd ReadParams err: %v", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	param := &sshs.SshParam{
		Host: params.Host,
		User: params.User,
		Pass: params.Password,
		Port: params.Port,
		Mode: 1,
	}
	InstallProgressLog("[INSTALL] CheckByUserPwd params:%v", *param)
	_, err := sshs.SSHManager.CheckConnection(param)
	if err != nil {
		InstallProgressLog("[INSTALL] CheckByUserPwd CheckConnection err: %v", err)
		apibase.ThrowSshHandleError(err)
	}
	InstallProgressLog("[INSTALL] CheckByUserPwd: %v", "SUCCESS")
	return map[string]interface{}{}
}

func CheckByUserPk(ctx context.Context) apibase.Result {
	InstallProgressLog("[INSTALL] CheckByUserPk ...")
	paramErrs := apibase.NewApiParameterErrors()

	params := struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		User string `json:"user"`
		Pk   string `json:"pk"`
	}{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
		InstallProgressLog("[INSTALL] CheckByUserPk ReadParams err: %v", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	param := &sshs.SshParam{
		Host: params.Host,
		User: params.User,
		Pk:   params.Pk,
		Port: params.Port,
		Mode: 2,
	}
	InstallProgressLog("[INSTALL] CheckByUserPk params:%v", *param)
	_, err := sshs.SSHManager.CheckConnection(param)
	if err != nil {
		InstallProgressLog("[INSTALL] CheckByUserPwd CheckConnection err: %v", err)
		apibase.ThrowSshHandleError(err)
	}
	InstallProgressLog("[INSTALL] CheckByUserPk: %v", "SUCCESS")
	return map[string]interface{}{}
}

func RunWithUserPwd(ctx context.Context) apibase.Result {
	InstallProgressLog("[INSTALL] RunWithUserPwd ...")
	paramErrs := apibase.NewApiParameterErrors()

	params := struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Cmd      string `json:"cmd"`
	}{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	param := &sshs.SshParam{
		Host: params.Host,
		User: params.User,
		Pass: params.Password,
		Port: params.Port,
		Mode: 1,
		Cmd:  params.Cmd,
	}
	InstallProgressLog("[INSTALL] RunWithUserPwd params:%v", *param)
	log.Debugf("RunWithUserPwd:%v", param)
	result, err := sshs.SSHManager.RunWithSSH(param, true)
	if err != nil {
		InstallProgressLog("[INSTALL] RunWithUserPwd err:%v, result: %v", err, result)
		apibase.ThrowSshHandleError(err.Error() + ": " + result)
	}
	log.Debugf("RunWithUserPwd:%v, result:%v", param, result)
	InstallProgressLog("[INSTALL] RunWithUserPwd result:%v", result)
	return map[string]interface{}{
		"result": result,
	}
}

func RunWithUserPk(ctx context.Context) apibase.Result {
	InstallProgressLog("[INSTALL] RunWithUserPk ...%v", "")
	paramErrs := apibase.NewApiParameterErrors()

	params := struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		User string `json:"user"`
		Pk   string `json:"pk"`
		Cmd  string `json:"cmd"`
	}{}
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	param := &sshs.SshParam{
		Host: params.Host,
		User: params.User,
		Pk:   params.Pk,
		Port: params.Port,
		Mode: 2,
		Cmd:  params.Cmd,
	}
	InstallProgressLog("[INSTALL] RunWithUserPk params:%v", *param)
	log.Debugf("RunWithUserPk:%v", *param)
	result, err := sshs.SSHManager.RunWithSSH(param, true)
	if err != nil {
		InstallProgressLog("[INSTALL] RunWithUserPk err:%v, result: %v", err, result)
		apibase.ThrowSshHandleError(err.Error() + ": " + result)
	}
	log.Debugf("RunWithUserPk:%v, result:%v", *param, result)
	InstallProgressLog("[INSTALL] RunWithUserPk result:%v", result)
	return map[string]interface{}{
		"result": result,
	}
}
