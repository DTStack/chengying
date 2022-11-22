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

package tracy

import (
	"context"
	"fmt"
	"time"

	"easyagent/internal/server/log"
	"easyagent/internal/server/publisher"
)

const (
	LogTsLayout = "2006-01-02 15:04:05.000000"
)

const (
	AGENT_INSTALL_PROGRESS_LOG = "agent-install"
	AGENT_CONTROL_PROGRESS_LOG = "agent-control"
)

func TracyOutput2Path(path string, toOutput bool, format string, args ...interface{}) error {

	if toOutput {
		format = time.Now().Format(LogTsLayout) + " " + format
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		index := "dtlog-1-easyagent-nodelete-" + time.Now().Format("2006.01.02") + "_000001.alias"
		body := struct {
			Msg        string `json:"msg"`
			LastUpdate string `json:"last_update_date"`
		}{fmt.Sprintf(format, args...), time.Now().Format(LogTsLayout)}
		if err := publisher.Publish.OutputJson(ctx, "", index, "dt_agent_install_progress",
			body, []byte{}); err != nil {
		}
		cancel()
	}
	return log.Output2Path(path, format, args...)
}

func InstallProgressLog(format string, args ...interface{}) {
	TracyOutput2Path(AGENT_INSTALL_PROGRESS_LOG, true, format, args...)
}

func ControlProgressLog(format string, args ...interface{}) {
	TracyOutput2Path(AGENT_CONTROL_PROGRESS_LOG, true, format, args...)
}
