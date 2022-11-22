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

package host

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"text/template"
	"time"

	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/asset"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"github.com/satori/go.uuid"
)

type agentInstall struct {
	AgentHost            string
	StaticHost           string
	AgentInstallCallback string
	AgentInstallPath     string
}

var (
	AgentInstall = &agentInstall{
		AgentHost:        "127.0.0.1:8889",
		StaticHost:       "http://127.0.0.1:8864",
		AgentInstallPath: "/opt/dtstack/easymanager/easyagent",
	}
)

var ScriptTemplates = map[string]*template.Template{}

func init() {
	initTemplates()
}

func initTemplates() {
	for _, typ := range []string{SCRIPT_WRAPPER_SH, ENVIRONMENT_INIT_SH, INSTALL_AGENTX_SH, DOCKER_ENVIRONMENT_INIT_SH} {
		switch typ {
		case SCRIPT_WRAPPER_SH:
			tmpldata, err := asset.Asset("templates/install.script.wrapper.sh")
			if err != nil {
				panic(err)
			}
			ScriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		case ENVIRONMENT_INIT_SH:
			tmpldata, err := asset.Asset("templates/environment.init.sh")
			if err != nil {
				panic(err)
			}
			ScriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		case INSTALL_AGENTX_SH:
			tmpldata, err := asset.Asset("templates/install.agentx.sh")
			if err != nil {
				panic(err)
			}
			ScriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		case DOCKER_ENVIRONMENT_INIT_SH:
			tmpldata, err := asset.Asset("templates/docker.environment.init.sh")
			if err != nil {
				panic(err)
			}
			ScriptTemplates[typ] = template.Must(template.New(typ).Parse(string(tmpldata)))
		}
	}
}

func InitAgentInstall(agentHost, staticHost, installPath string) {
	AgentInstall.AgentHost = agentHost
	AgentInstall.StaticHost = staticHost
	AgentInstall.AgentInstallPath = installPath
}

func (this *agentInstall) getCallbackURL(aid int, localAddr string) (error, string) {
	return nil, this.getCallbackUrlBase64(fmt.Sprintf("http://%s/api/v2/agent/install/callback?aid=%d", localAddr, aid))
}

func (this *agentInstall) getCallbackUrlBase64(url string) string {
	return base64.StdEncoding.EncodeToString([]byte(url))
}

//curl 'http://172.16.10.74:8889/api/v1/deploy/sidecar/install/shell?TargetPath=~/dtstack/easyagent&CallBack=aHR0cDovL2xvY2FsaG9zdDo4ODU0L29wZW4vYXBpL3YyL2FnZW50L2luc3RhbGw%2FdGVuYW50SWQ9MSZhaWQ9LTE%3D' | sh
func (this *agentInstall) GetAgentInstallCmd(aid int, localAddr, ctype, clusterId, roles string) (error, string) {
	cmd := "curl -s 'http://"

	cmd = cmd + this.AgentHost + "/api/v1/deploy/sidecar/install/shell?"
	cmd = cmd + "TargetPath=" + this.AgentInstallPath
	err, callback := this.getCallbackURL(aid, localAddr)
	if err != nil {
		return err, ""
	}
	cmd = cmd + "&CallBack=" + callback
	cmd = cmd + "&Type=" + ctype
	cmd = cmd + "&ClusterId=" + clusterId
	cmd = cmd + "&Roles=" + roles
	cmd = cmd + "' | sh"

	log.Debugf("GetAgentInstallCmd: %v", cmd)

	return nil, cmd
}

func (this *agentInstall) GetAgentCallBack(aid int, localAddr string) (string, error) {
	err, callback := this.getCallbackURL(aid, localAddr)
	if err != nil {
		return "", err
	}
	return callback, err
}

func (this *agentInstall) resetTemplates() {
	err := asset.ResettemplatesWithLocalFile()
	if err != nil {
		log.Errorf("InitScriptInstallSh error: %v", err)
		return
	}
	initTemplates()
}

func (this *agentInstall) LoopInstallScriptWrapper(sidecarId, server, targetPath, debug string) (error, string) {
	this.resetTemplates()
	var err error
	if len(server) == 0 {
		server = this.AgentHost
	}
	tryCount := 1
	agentId := ""
	for {
		log.Debugf("[InstallScriptWrapper]LoopInstallScriptWrapper: %v, sid: %v", tryCount, sidecarId)
		err, agentId = this.installScriptWrapper(sidecarId, server, targetPath, debug)

		if err == nil {
			break
		} else if agentId != "" {
			this.uninstallScriptWrapper(agentId)
		}
		if tryCount > 3 {
			break
		} else {
			log.Errorf("[InstallScriptWrapper]LoopInstallScriptWrapper: %v, sid: %v", err, sidecarId)
		}
		tryCount = tryCount + 1
		time.Sleep(3 * time.Second)
	}
	return err, agentId
}

func (this *agentInstall) installScriptWrapper(sidecarId, server, targetPath, debug string) (error, string) {
	log.Debugf("[InstallScriptWrapper]Install with collectorid %v", sidecarId)

	typ := SCRIPT_WRAPPER_SH
	script := &bytes.Buffer{}
	ScriptTemplates[typ].Execute(script, map[string]interface{}{
		"EASYAGENT_SERVER": server,
		"INSTALL_PATH":     targetPath,
	})
	params := &agent.InstallParms{}
	params.Name = "script-wrapper"
	params.BinaryPath = targetPath + "/script_wrapper/script-wrapper"
	params.ConfigurationPath = targetPath + "/script_wrapper/script-wrapper.yml"
	if len(debug) > 0 {
		params.Parameter = "-c," + targetPath + "/script_wrapper/script-wrapper.yml,--" + debug
	} else {
		params.Parameter = "-c," + targetPath + "/script_wrapper/script-wrapper.yml"
	}
	params.CollectorId = sidecarId
	params.InstallScript = script.String()
	params.InstallParameter = ""

	collectorId, err := uuid.FromString(params.CollectorId)
	if err != nil {
		log.Errorf("[InstallScriptWrapper]installScriptWrapper error: %v, sid: %v", params.CollectorId, collectorId)
		return err, ""
	}
	err, respBody := agent.AgentClient.AgentInstall(params, "")

	if err != nil {
		log.Errorf("[InstallScriptWrapper] script_wrapper install err: %v, sid: %v", err, sidecarId)
		return err, ""
	}
	if respBody.Data == nil {
		log.Errorf("[InstallScriptWrapper] script_wrapper install err: %v, sid: %v", respBody.Msg, sidecarId)
		return fmt.Errorf("script_wrapper install err: %v, sid: %v", respBody.Msg, sidecarId), ""
	}
	agentId, exists := respBody.Data.(map[string]interface{})["agent_id"]
	if !exists {
		agentId = ""
	}
	return nil, agentId.(string)
}

func (this *agentInstall) uninstallScriptWrapper(agentId string) error {
	log.Debugf("[unInstallScriptWrapper]unInstall with agentId %v", agentId)
	var err error
	params := &agent.ShellParams{}
	params.Parameter = ""
	params.ShellScript = "#!/bin/sh\nrm -rf " + util.ShellQuote(this.AgentInstallPath)

	err, _ = agent.AgentClient.AgentUninstall(agentId, params, "")
	if err != nil {
		log.Errorf("[InstallScriptWrapper] script_wrapper uninstall err: %v, agentId: %v", err, agentId)
		return err
	}
	return nil
}

func (this *agentInstall) StartScriptWrapper(agentId string) error {
	var err error

	err, _ = agent.AgentClient.AgentStart(agentId, "")

	if err != nil {
		log.Errorf("[InstallScriptWrapper] script_wrapper start err: %v, agentId: %v", err, agentId)
		return err
	}
	return nil
}

func (this *agentInstall) EnvironmentInit(sid string, execId string) error {
	this.resetTemplates()
	log.Debugf("[EnvironmentInit] sid:%v ", sid)
	var err error
	u, err := url.Parse(this.StaticHost)
	if err != nil {
		return err
	}
	typ := ENVIRONMENT_INIT_SH
	script := &bytes.Buffer{}
	err = ScriptTemplates[typ].Execute(script, map[string]interface{}{
		"STATIC_HOST": this.StaticHost,
		"NTP_SERVER":  strings.Split(u.Host, ":")[0],
	})
	if err != nil {
		return err
	}
	return this.execScript(sid, script.String(), execId)
}

func (this *agentInstall) DockerEnvironmentInit(sid string) error {
	this.resetTemplates()
	log.Infof("-->Step DockerEnvironmentInit sid:%v ", sid)
	log.OutputInfof(sid, "-->Step DockerEnvironmentInit sid:%v ", sid)
	var err error
	typ := DOCKER_ENVIRONMENT_INIT_SH
	script := &bytes.Buffer{}
	u, err := url.Parse(this.StaticHost)
	if err != nil {
		return err
	}
	err = ScriptTemplates[typ].Execute(script, map[string]interface{}{
		"STATIC_HOST": this.StaticHost,
		"MATRIX_IP":   strings.Split(u.Host, ":")[0],
	})
	if err != nil {
		return err
	}
	err = this.execScript(sid, script.String(), "")
	if err != nil {
		return err
	}
	log.Infof("<--Step DockerEnvironmentInit success sid:%v ", sid)
	log.OutputInfof(sid, "<--Step DockerEnvironmentInit success sid:%v ", sid)
	return nil
}

func (this *agentInstall) execScript(sid, script string, execId string) error {
	var err error
	params := &agent.ExecScriptParams{}
	params.Parameter = ""
	params.ExecScript = script
	params.Timeout = "30m"

	err, respBody := agent.AgentClient.AgentExec(sid, params, execId)
	if err != nil {
		return err
	}
	result, exists := respBody.Data.(map[string]interface{})["result"]
	if !exists {
		return fmt.Errorf("no result")
	}
	failed, exists := result.(map[string]interface{})["failed"]
	if exists && failed.(bool) == true {
		return fmt.Errorf("%v", result.(map[string]interface{})["response"])
	}
	return nil
}
