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

package sshs

import (
	"errors"

	"easyagent/internal/server/log"
	"fmt"
)

var (
	ErrParamIsNull = errors.New("ssh connect params is nil")
)

var (
	MAX_SSH_CLIENT = 3
	CHANNEL_BUSY   = 1
)

var (
	SSHManager Manager = &manage{
		mga:  make(chan int, 1),
		stop: false,
	}
)

type Manager interface {
	CheckConnection(params *SshParam) (bool, error)
	RunWithSSH(params *SshParam, sync bool) (string, error)
	RunWithSSHS(params []*SshParam, sync bool)
}

type manage struct {
	mga  chan int
	stop bool
}

func (ma *manage) checkConnect(param *SshParam) (bool, error) {
	var err error
	if param == nil {
		return false, ErrParamIsNull
	}
	cli, err := CreateWithParam(param)

	if err != nil {
		return false, err
	}
	cli.Close()
	return true, nil
}

func (ma *manage) runWithSSH(param *SshParam, sync bool) (string, error) {
	var result string
	var err error
	if param == nil {
		return "", ErrParamIsNull
	}
	cli, err := CreateWithParam(param)

	if err != nil {
		return "", err
	}
	if sync {
		result, err = cli.RunSync(param.Cmd)
	} else {
		err = cli.Run(param.Cmd)
	}
	cli.Close()
	return result, err
}

func (ma *manage) CheckConnection(params *SshParam) (bool, error) {

	result, err := ma.checkConnect(params)

	if err != nil {
		return result, fmt.Errorf("ssh 连通性验证失败,请检查参数: %s-%d-%s", params.Host, params.Port, params.User)
	}
	return result, nil
}

func (ma *manage) RunWithSSH(params *SshParam, sync bool) (string, error) {

	//ma.mga <- CHANNEL_BUSY

	result, err := ma.runWithSSH(params, sync)

	//<-ma.mga
	return result, err
}

func (ma *manage) RunWithSSHS(params []*SshParam, sync bool) {
	for _, param := range params {
		if ma.stop {
			log.Debugf("Stop run ssh cmd:%v", param.Cmd)
		}
		result, err := ma.runWithSSH(param, sync)
		if err != nil {
			log.Errorf("RunWithSSHS err: %v, %v", err.Error(), result)
		}
	}
	if ma.stop {
		ma.stop = false
	}
}

func (ma *manage) StopAll() {
	ma.stop = true
}
