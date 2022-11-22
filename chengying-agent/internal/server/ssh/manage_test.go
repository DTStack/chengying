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
	"testing"

	"easyagent/internal/server/log"
	"fmt"
)

func init() {
	log.ConfigureLogger("/tmp", 0, 0, 0)
}

func TestManage_RunWithSSH(t *testing.T) {
	param := &SshParam{
		Host: "172.16.10.108",
		User: "dtstack",
		Pass: "abc123",
		Port: 22,
		Mode: 1,
		Cmd:  "sh /opt/dtstack/easymanager/easyagent/easyagent.sh restart",
	}
	result, err := SSHManager.RunWithSSH(param, true)
	if err != nil {
		t.Errorf("RunWithSSH err:%v, %v", err, result)
	}
	fmt.Println(result)
}

func TestManage_RunWithSSHS(t *testing.T) {
	param := &SshParam{
		Host: "172.16.10.108",
		User: "dtstack",
		Pass: "abc123",
		Port: 22,
		Mode: 1,
		Cmd:  "sudo systemctl status network",
	}
	params := []*SshParam{param}
	SSHManager.RunWithSSHS(params, true)
}
