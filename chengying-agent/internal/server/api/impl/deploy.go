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
	"bytes"
	"encoding/base64"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	apibase "easyagent/go-common/api-base"
	"easyagent/go-common/utils"
	"easyagent/internal/server/asset"
	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	. "easyagent/internal/server/tracy"

	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
)

var scriptTemplates = map[string]*template.Template{}

const EASYAGENT_INSTALL_FILE = "easyagent_install.sh"

func init() {
	for _, typ := range []string{"sidecar", "wrapper", "sidecar_for_win"} {
		switch typ {
		case "sidecar":
			tmpldata, err := asset.Asset("templates/install.sidecar.sh")
			if err != nil {
				panic(err)
			}
			scriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		case "wrapper":
			tmpldata, err := asset.Asset("templates/install.script.wrapper.sh")
			if err != nil {
				panic(err)
			}
			scriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		case "sidecar_for_win":
			tmpldata, err := asset.Asset("templates/install.sidecar.win.ps1")
			if err != nil {
				panic(err)
			}
			scriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		}
	}
}

func InitSidecarInstallSh() {
	err := asset.ResetInstallSidecarShWithLocalFile()

	if err != nil {
		log.Errorf("InitSidecarInstallSh error: %v", err)
		InstallProgressLog("[INSTALL] InitSidecarInstallSh error: %v", err)
		return
	}
	for _, typ := range []string{"sidecar", "sidecar_for_win"} {
		switch typ {
		case "sidecar":
			tmpldata, err := asset.Asset("templates/install.sidecar.sh")
			if err != nil {
				panic(err)
			}
			scriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		case "sidecar_for_win":
			tmpldata, err := asset.Asset("templates/install.sidecar.win.ps1")
			if err != nil {
				panic(err)
			}
			scriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		}
	}
	InstallProgressLog("[INSTALL] InitSidecarInstallSh success: %v", EASYAGENT_INSTALL_FILE)
}

func GetSidecarInstallShell(ctx context.Context) apibase.Result {
	InitSidecarInstallSh()
	paramErrs := apibase.NewApiParameterErrors()
	typ := "sidecar"
	sid := uuid.NewV4()
	now := time.Now().Unix()
	params := url.Values{}
	sidecatInstallPath := ""
	for key, value := range ctx.URLParams() {
		params.Set(key, value)
		if key == "TargetPath" && !strings.Contains(value, "~") {
			sidecatInstallPath = value
		}
	}
	v, err := url.ParseQuery(params.Encode())
	if err != nil {
		log.Errorf("ParseQuery %v error: %v", params.Encode(), err)
		InstallProgressLog("[INSTALL] GetSidecarInstallShell ParseQuery %v error: %v", params.Encode(), err)
		paramErrs.AppendError("GetSidecarInstallShell ParseQuery: %v", err)
	}
	if v.Get("CallBack") == "" {
		log.Errorf("CallBack is empty!")
		InstallProgressLog("[INSTALL] GetSidecarInstallShell CallBack is empty, sid: %v", sid)
		paramErrs.AppendError("GetSidecarInstallShell CallBack is empty %v", "")
	}
	mode := v.Get("Debug")
	if mode != "" {
		mode = "--debug"
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	if v.Get("TargetOs") == "windows" {
		typ = "sidecar_for_win"
	}
	rpcPort, err := strconv.Atoi(os.Getenv("RPC_PORT"))
	if err != nil {
		rpcPort = 8890
	}

	if t := scriptTemplates[typ]; t != nil {
		buf := &bytes.Buffer{}
		ip, port, _ := net.SplitHostPort(ctx.Host())
		t.Execute(buf, map[string]interface{}{
			"UUID":                 sid,
			"INSTALL_TYPE":         typ,
			"create_time":          now,
			"callback":             ctx.Host(),
			"SERVER_IP_ADDRESS":    ip,
			"SERVER_HOST_PORT":     port,
			"SIDECAR_INSTALL_PATH": sidecatInstallPath,
			"CALLBACK_TO_PROD":     params.Encode(),
			"DEBUG_MODE":           "--debug",
			"TARGET_OS":            v.Get("TargetOs"),
			"RPC_PORT":             rpcPort,
		})

		return buf.String()
	}
	return nil
}

func DeploySidecarMain(ctx context.Context) apibase.Result {
	return nil
}

func GetSidecarInstallTargz(ctx context.Context) apibase.Result {
	return nil
}

/*
 一键部署回调入库
*/
func GetSidecarInstallCallback(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	installType := ctx.FormValue("install_type")
	clientIDStr := ctx.FormValue("client_id")

	clientID, err := uuid.FromString(clientIDStr)
	log.Debugf("uuid: %v, ori_uuid : %v", clientID, clientIDStr)
	InstallProgressLog("[INSTALL] GetSidecarInstallCallback uuid: %v, ori_uuid : %v", clientID, clientIDStr)

	if err != nil {
		paramErrs.AppendError("client_id", "clientID not UUID format!")
		InstallProgressLog("[INSTALL] GetSidecarInstallCallback clientID not UUID format!")
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	installRes := ctx.FormValue("install_res")
	installTime, _ := strconv.ParseInt(ctx.FormValue("check_flg"), 10, 64)

	bMsg, err := base64.StdEncoding.DecodeString(ctx.FormValue("msg"))
	if err != nil {
		log.Errorf("decode msg error: %v", err)
	}

	callback := &model.DeployCallbackInfo{
		Time:        installTime,
		ClientID:    clientID.String(),
		InstallType: installType,
		InstallRes:  installRes,
		MSG:         bMsg,
		RequestUrl:  ctx.Request().URL.String(),
		IP:          ctx.RemoteAddr(),
	}
	res, err := model.DeployCallback.CreateDeployCallback(callback)
	if err != nil {
		paramErrs.AppendError("CreateDeployCallback", err)
		InstallProgressLog("[INSTALL] GetSidecarInstallCallback CreateDeployCallback err: %v", err)
		paramErrs.CheckAndThrowApiParameterErrors()
	}
	b, _ := utils.MapTagFromStruct(res, "db")

	if installRes != "success" {
		InstallProgressLog("[INSTALL] Sidecar Agent failed!")
		return b
	}

	_, err = model.SidecarList.CreateSidecarByDeploy(clientID)
	if err != nil {
		paramErrs.AppendError("createSidecarByDeploy", err)
		InstallProgressLog("[INSTALL] GetSidecarInstallCallback createSidecarByDeploy err: %v", err)
		paramErrs.CheckAndThrowApiParameterErrors()
	}

	InstallProgressLog("[INSTALL] GetSidecarInstallCallback success: %v", clientID)

	return b
}

func RetDashboardUrl(ctx context.Context) apibase.Result {
	id_ := ctx.FormValue("id")
	type_ := ctx.FormValue("type")

	if type_ == "cluster" {
		return map[string]interface{}{
			"url": model.DashboardList.RetUrlByClusterID(id_),
		}
	}
	if type_ == "services" {
		return map[string]interface{}{
			"url": model.DashboardList.RetUrlByServicesID(id_),
		}
	}
	return nil
}
