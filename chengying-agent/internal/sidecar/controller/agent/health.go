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

package agent

import (
	"context"
	"os"
	"runtime"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/sidecar/base"
	"easyagent/internal/sidecar/controller/util"
	"easyagent/internal/sidecar/event"
)

func (ag *agent) startHealthCheck() {
	if ag.healthShell == "" {
		return
	}
	if ag.healthStopCh != nil {
		// health check already started
		return
	}

	ag.wgHealth.Add(1)
	go func() {
		ag.healthStopCh = make(chan struct{})
		ag.healthCtx, ag.healthCancel = context.WithCancel(context.Background())

		base.Infof("agent %v start health check: %v, %v", ag.agentId, ag.healthShell, ag.healthPeriod)
		defer func() {
			base.Infof("agent %v stop health check", ag.agentId)
			ag.wgHealth.Done()
		}()

		//sleep(ag.healthCtx, ag.healthStartPeriod)
		//base.Infof("agent %v health check waited %v", ag.agentId, ag.healthStartPeriod)

		ticker := time.NewTicker(ag.healthPeriod)
		defer ticker.Stop()
		healthFailedCount := uint64(0)
		healthStart := time.Now()
		for {
			select {
			case <-ag.healthStopCh:
				return
			case <-ticker.C:
				ev := &proto.Event_AgentHealthCheck{AgentId: ag.agentId.Bytes()}
				if ag.runHealthCheck() {
					healthFailedCount = 0
					event.ReportEvent(ev)
				} else {
					healthFailedCount++
					if healthFailedCount >= ag.healthRetries && time.Now().Sub(healthStart) > ag.healthStartPeriod {
						ev.Failed = true
						event.ReportEvent(ev)
						base.Infof("agent %v health check failed %d times", ag.agentId, healthFailedCount)
					}
				}
			}
		}
	}()
}

func (ag *agent) stopHealthCheck() {
	if ag.healthShell == "" {
		return
	}
	if ag.healthStopCh != nil {
		close(ag.healthStopCh)
		ag.healthCancel()
		ag.wgHealth.Wait()

		// reset healthStopCh for start
		ag.healthStopCh = nil
	}
}

func (ag *agent) runHealthCheck() bool {
	ctx, cancel := context.WithTimeout(ag.healthCtx, ag.healthTimeout)
	defer cancel()

	var cmd *util.Cmd
	if runtime.GOOS == "windows" {
		cmd = util.CommandContext(ctx, "", false, nil, "cmd.exe", "/c", ag.healthShell)
	} else {
		cmd = util.CommandContext(ctx, "", false, nil, "sh", "-c", ag.healthShell)
	}

	if ag.workdir != "" {
		cmd.Dir = ag.workdir
	} else {
		cmd.Dir = os.TempDir()
	}
	err := cmd.Run()

	if err == nil {
		return true
	} else {
		base.Infof("agent %v health check failed: %s", ag.agentId, err.Error())
		return false
	}
}

func sleep(ctx context.Context, duration time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	<-ctx.Done()
}
