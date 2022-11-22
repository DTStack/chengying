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

package agent

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"errors"
)

var (
	AgentClient = &agentClient{}
)

const (
	REST_TRY_COUNT = 3
)

type agentClient struct {
	httpClient *http.Client
	baseUrl    *url.URL
}

func InitAgentClient(agentHost string) error {
	AgentClient.httpClient = util.DefaultClient

	var err error
	AgentClient.baseUrl, err = url.Parse("http://" + agentHost)
	if err != nil {
		log.Errorf("[InitAgentClient] init err: %v", err)
	}
	return err
}

func (this *agentClient) GetAgentHost() string {
	uri := ""
	if this.baseUrl != nil {
		uri = this.baseUrl.String()
	}
	return uri
}

func (this *agentClient) AgentRestCore(method, uri string, params map[string]string, body interface{}, execId string) (error, *EasyagentServerResponse) {
	c := util.NewClient(this.httpClient)
	c.BaseURL = this.baseUrl

	var err error
	var r *http.Request
	respBody := new(EasyagentServerResponse)
	tryCount := 1
	for {
		if tryCount >= REST_TRY_COUNT {
			break
		}
		log.Debugf("[AgentRestCore]LoopAgentRestCore: %v, request uri: %v", tryCount, uri)
		if r, err = c.NewRequest(method, uri, params, body, execId); err != nil {
			return err, respBody
		}
		log.Debugf("[AgentRestCore]LoopAgentRestCore: %v, response body: %v", tryCount, respBody)
		_, err = c.Do(r, respBody)

		if err != nil && tryCount < REST_TRY_COUNT {
			tryCount++
			//log.Errorf("AgentRestCore do request err: %v, try %v", err, tryCount)
			time.Sleep(1 * time.Second)
			continue
		}
		if respBody.Code != 0 {
			tryCount++
			err = fmt.Errorf("%v", respBody.Data)
			log.Errorf("AgentRestCore do request err: %v, try %v", err, tryCount)
			time.Sleep(1 * time.Second)
			continue
		}

		switch respBody.Data.(type) {
		case map[string]interface{}:
		default:
			err = fmt.Errorf("response data format not surpported: %v", respBody)
		}
		if err != nil {
			break
		}
		result, ok := respBody.Data.(map[string]interface{})["result"]
		if !ok {
			break
		}
		failed, ok := result.(map[string]interface{})["failed"]
		if !ok {
			break
		}
		if failed.(bool) {
			message, ok := result.(map[string]interface{})["message"]
			if ok && (message.(string) == RETRY_AGENT_IS_RUNNING_TASK) && tryCount < REST_TRY_COUNT {
				tryCount++
				log.Errorf("AgentRestCore response err: %v, try %v", message.(string), tryCount)
				time.Sleep(1 * time.Second)
				continue
			}
			response, ok := result.(map[string]interface{})["response"]
			if ok && (response.(string) == RETRY_AGENT_IS_RUNNING_TASK) && tryCount < REST_TRY_COUNT {
				tryCount++
				log.Errorf("AgentRestCore response err: %v, try %v", response.(string), tryCount)
				time.Sleep(1 * time.Second)
				continue
			}
			log.Errorf("AgentRestCore response err: %v, try %v", result, tryCount)
			break
		}
	}

	return err, respBody
}

func (this *agentClient) AgentInstall(param *InstallParms, execId string) (error, *EasyagentServerResponse) {
	log.Debugf("[AgentClient] AgentInstall with params:%v ", param)
	return this.AgentRestCore("POST", AGENT_INSTALL_SYNC_URI, nil, param, execId)
}

func (this *agentClient) AgentUninstall(agentId string, param *ShellParams, execId string) (error, *EasyagentServerResponse) {
	log.Debugf("[AgentClient] AgentUnInstall with params:%v ", param)

	return this.AgentRestCore("POST", fmt.Sprintf(AGENT_UNINSTALL_SYNC_URI, agentId), nil, param, execId)
}

func (this *agentClient) AgentCancel(param *CancelParams) (error, *EasyagentServerResponse) {
	log.Debugf("[AgentClient] AgentCancel with params:%v ", param)
	return this.AgentRestCore("POST", AGENT_CANCEL_SYNC_URI, nil, param, "")
}

func (this *agentClient) AgentStartWithParam(param *StartParams, agentId string, execId string) (error, *EasyagentServerResponse) {
	log.Debugf("[AgentClient] AgentStart with params:%v ", param)
	return this.AgentRestCore("POST", fmt.Sprintf(AGENT_START_SYNC_PARAMS_URI, agentId), nil, param, execId)
}

func (this *agentClient) AgentStart(agentId string, execId string) (error, *EasyagentServerResponse) {
	log.Debugf("[AgentClient] AgentStart with params:%v ", agentId)
	return this.AgentRestCore("GET", fmt.Sprintf(AGENT_START_SYNC_URI, agentId), nil, nil, execId)
}

func (this *agentClient) AgentStop(agentId string, stopAgentOptionsType int, execId string) (error, *EasyagentServerResponse) {
	log.Debugf("[AgentClient] AgentStop with params:%v ", agentId)
	params := map[string]string{"stop_agent_options_type": strconv.Itoa(stopAgentOptionsType)}
	return this.AgentRestCore("GET", fmt.Sprintf(AGENT_STOP_SYNC_URI, agentId), params, nil, execId)
}

func (this *agentClient) AgentExec(sid string, params *ExecScriptParams, execId string) (error, *EasyagentServerResponse) {
	//log.Debugf("[AgentClient] AgentExec with params:%v, %v", sid, params)
	return this.AgentRestCore("POST", fmt.Sprintf(AGENT_EXEC_SYNC_URI, sid), nil, params, execId)
}

func (this *agentClient) AgentExecOften(sid string, params *ExecScriptParams, execId string) (error, *EasyagentServerResponse) {
	//log.Debugf("[AgentClient] AgentExec with params:%v, %v", sid, params)
	return this.AgentRestCore("POST", fmt.Sprintf(AGENT_EXEC_OFTEN_SYNC_URI, sid), nil, params, execId)
}

func (this *agentClient) AgentExecRest(sid string, params *ExecRestParams, execId string) (error, *EasyagentServerResponse) {
	//log.Debugf("[AgentClient] AgentExecRest with params:%v, %v", sid, params)
	return this.AgentRestCore("POST", fmt.Sprintf(AGENT_REST_SYNC_URI, sid), nil, params, execId)
}

func (this *agentClient) AgentConfigUpdate(sid, agentId string, params *ConfigParams, execId string) (error, *EasyagentServerResponse) {
	log.Debugf("[AgentClient] Agent config update with params:%v, %v", agentId, params)

	var configPath string
	// copy previous config file
	if !filepath.IsAbs(params.ConfigPath) {
		configPath = filepath.Join(params.WorkDir, params.ConfigPath)
	}
	copyCmd := fmt.Sprintf("#!/bin/sh\n cp -pr %s %s", configPath, configPath+"_"+strconv.FormatInt(time.Now().Unix(), 10))
	_, err := this.ToExecCmd(sid, agentId, copyCmd, execId)
	if err != nil {
		log.Errorf("[AgentClient] Agent config copy error: %v", err)
		return err, nil
	}
	// update config file
	return this.AgentRestCore("POST", fmt.Sprintf(AGENT_CONFIG_SYNC_URI, agentId), nil, params, execId)
}

func (this *agentClient) ToExecCmd(sid, agentId string, cmd string, execId string) (content string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	params := &ExecScriptParams{}
	params.AgentId = agentId
	params.Timeout = "60s"
	params.ExecScript = cmd
	err, respBody := this.AgentExec(sid, params, execId)
	if err != nil {
		return "", err
	}
	result, exists := respBody.Data.(map[string]interface{})["result"]
	if !exists {
		return "", fmt.Errorf("without result: %v", respBody)
	}
	failed, exists := result.(map[string]interface{})["failed"]
	if exists && failed.(bool) == true {
		if _, exists := result.(map[string]interface{})["response"]; exists {
			return "", fmt.Errorf("failed: %v", result.(map[string]interface{})["response"].(string))
		}
		if _, exists := result.(map[string]interface{})["message"]; exists {
			return "", fmt.Errorf("failed: %v", result.(map[string]interface{})["message"].(string))
		}
	}
	response, exists := result.(map[string]interface{})["response"]
	if exists {
		content = response.(string)
	}
	return content, nil
}

func (this *agentClient) ToExecCmdWithTimeout(sid, agentId string, cmd, timeout, parameter string, execId string) (content string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	params := &ExecScriptParams{}
	params.AgentId = agentId
	params.Timeout = timeout
	params.ExecScript = cmd
	params.Parameter = parameter
	err, respBody := this.AgentExec(sid, params, execId)
	if err != nil {
		return "", err
	}
	result, exists := respBody.Data.(map[string]interface{})["result"]
	if !exists {
		return "", fmt.Errorf("without result: %v", respBody)
	}
	failed, exists := result.(map[string]interface{})["failed"]
	if exists && failed.(bool) == true {
		if _, exists := result.(map[string]interface{})["response"]; exists {
			return "", fmt.Errorf("failed: %v", result.(map[string]interface{})["response"].(string))
		}
		if _, exists := result.(map[string]interface{})["message"]; exists {
			return "", fmt.Errorf("failed: %v", result.(map[string]interface{})["message"].(string))
		}
	}
	response, exists := result.(map[string]interface{})["response"]
	if exists {
		content = response.(string)
	}
	return content, nil
}

func (this *agentClient) ToExecRest(sid string, params *ExecRestParams, execId string) (content string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	err, respBody := this.AgentExecRest(sid, params, execId)
	if err != nil {
		return "", err
	}
	result, exists := respBody.Data.(map[string]interface{})["result"]
	if !exists {
		return "", fmt.Errorf("without result: %v", respBody)
	}
	failed, exists := result.(map[string]interface{})["failed"]
	if exists && failed.(bool) == true {
		if _, exists := result.(map[string]interface{})["response"]; exists {
			return "", fmt.Errorf("failed: %v", result.(map[string]interface{})["response"].(string))
		}
		if _, exists := result.(map[string]interface{})["message"]; exists {
			return "", fmt.Errorf("failed: %v", result.(map[string]interface{})["message"].(string))
		}
	}
	response, exists := result.(map[string]interface{})["response"]
	if exists {
		content = response.(string)
	}
	return content, nil
}
