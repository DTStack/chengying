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
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	"easyagent/internal/server/rpc"
	. "easyagent/internal/server/tracy"
	"github.com/kataras/iris/context"
	"github.com/satori/go.uuid"
)

func getRegisterServerAddr(sidecarId uuid.UUID) (string, int, error) {

	info, err := model.SidecarList.GetSidecarInfo(sidecarId)

	if err != nil {
		return "", -1, err
	}
	if len(info.ServerHost) == 0 {
		return "", -1, fmt.Errorf("ServerHost is empty")
	}
	return info.ServerHost, info.ServerPort, nil
}

func checkSidConnected(target string, sidecarId uuid.UUID) (bool, string, int, error) {
	serverHost, serverPort, err := getRegisterServerAddr(sidecarId)

	if err != nil {
		return false, serverHost, serverPort, err
	}
	serverAddr := net.JoinHostPort(serverHost, strconv.Itoa(serverPort))

	if !strings.Contains(target, serverAddr) {
		return false, serverHost, serverPort, nil
	}
	return true, serverHost, serverPort, nil
}

func checkForSLB(ctx context.Context, sidecarId uuid.UUID) bool {

	if rpc.SidecarClient.IsClientExist(sidecarId) {
		return false
	}
	serverHost, serverPort, err := getRegisterServerAddr(sidecarId)

	if ctx.Request().Host == net.JoinHostPort(serverHost, strconv.Itoa(serverPort)) {
		return false
	}
	log.Debugf("serverHost %v, serverPort%v, err%v", serverHost, serverPort, err)
	ControlProgressLog("[AGENT-CONTROL] serverHost %v, serverPort%v, err%v", serverHost, serverPort, err)

	if err != nil {
		log.Errorf("checkForSLB err: %v", err)
		return false
	}
	redirectUrl := "http://" + serverHost + ":" + strconv.Itoa(serverPort) + ctx.Request().RequestURI

	log.Debugf("redirect:%v", redirectUrl)
	ControlProgressLog("[AGENT-CONTROL] redirect:%v", redirectUrl)

	ctx.Redirect(redirectUrl, http.StatusTemporaryRedirect)

	return true
}
