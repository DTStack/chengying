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

package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/sidecar/base"
	"easyagent/internal/sidecar/client"
	"easyagent/internal/sidecar/config"
	"easyagent/internal/sidecar/controller/agent"
	"easyagent/internal/sidecar/controller/util"
	"easyagent/internal/sidecar/event"
	"github.com/satori/go.uuid"
)

var (
	httpclient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
)

type Controller struct {
	client     client.EaClienter
	agentsFile string

	flushAgentsCh chan struct{}

	maxExec int32

	sync.RWMutex
	agents map[uuid.UUID]agent.Agenter
}

func NewController(client client.EaClienter, agents []config.AgentConfig, agentsFile string) *Controller {
	c := &Controller{
		client:        client,
		agentsFile:    agentsFile,
		flushAgentsCh: make(chan struct{}),
		maxExec:       10,
		agents:        make(map[uuid.UUID]agent.Agenter, len(agents)),
	}
	for _, ag := range agents {
		c.agents[ag.AgentId.UUID] = agent.NewAgent(ag, c.flushAgentsCh)
	}
	return c
}

func (c *Controller) AddAgent(agentId uuid.UUID, agent agent.Agenter) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.agents[agentId]; !ok {
		c.agents[agentId] = agent
	}
}

func (c *Controller) getAgent(agentId uuid.UUID) (agent.Agenter, bool) {
	c.RLock()
	defer c.RUnlock()

	ag, ok := c.agents[agentId]
	return ag, ok
}

func (c *Controller) delAgent(agentId uuid.UUID) {
	c.Lock()
	defer c.Unlock()

	delete(c.agents, agentId)
}

func (c *Controller) dispatch() {
	for {
		ctlResp := c.client.GetControlResponse()
		base.Debugf("recv control command: %v", ctlResp.Cmd)

		switch ctl := ctlResp.Options.(type) {
		case *proto.ControlResponse_InstallAgentOptions_:
			agentId, _ := uuid.FromBytes(ctl.InstallAgentOptions.AgentId)
			if agentId == uuid.Nil {
				base.Errorf("empty agentId!")
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.InstallAgentOptions.AgentId,
					Failed:  true,
					Message: "empty agentId!",
				})
				continue
			}
			if _, ok := c.getAgent(agentId); ok {
				err := fmt.Errorf("agent %v has already installed", agentId)
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.InstallAgentOptions.AgentId,
					Failed:  true,
					Message: err.Error(),
				})
				continue
			}
			if ctl.InstallAgentOptions.HealthShell != "" {
				if ctl.InstallAgentOptions.HealthPeriod < time.Second {
					err := fmt.Errorf("agent %v health-period less than 1 sec", agentId)
					base.Errorf("%v", err)
					event.ReportEvent(&proto.Event_OperationProgress{
						Seqno:   ctlResp.Seqno,
						AgentId: ctl.InstallAgentOptions.AgentId,
						Failed:  true,
						Message: err.Error(),
					})
					continue
				}
				if ctl.InstallAgentOptions.HealthRetries > 20 {
					err := fmt.Errorf("agent %v health-retries too big (> 20)", agentId)
					base.Errorf("%v", err)
					event.ReportEvent(&proto.Event_OperationProgress{
						Seqno:   ctlResp.Seqno,
						AgentId: ctl.InstallAgentOptions.AgentId,
						Failed:  true,
						Message: err.Error(),
					})
					continue
				} else if ctl.InstallAgentOptions.HealthRetries <= 0 {
					ctl.InstallAgentOptions.HealthRetries = 1
				}
			}

			ag := agent.NewAgent(config.AgentConfig{
				BinaryPath:        ctl.InstallAgentOptions.BinaryPath,
				AgentId:           config.WrapperUUID{UUID: agentId},
				ConfigurationPath: ctl.InstallAgentOptions.ConfigurationPath,
				Parameter:         ctl.InstallAgentOptions.Parameter,
				Workdir:           ctl.InstallAgentOptions.Workdir,
				Enabled:           false,
				Name:              ctl.InstallAgentOptions.Name,
				HealthShell:       ctl.InstallAgentOptions.HealthShell,
				HealthPeriod:      ctl.InstallAgentOptions.HealthPeriod,
				HealthStartPeriod: ctl.InstallAgentOptions.HealthStartPeriod,
				HealthTimeout:     ctl.InstallAgentOptions.HealthTimeout,
				HealthRetries:     ctl.InstallAgentOptions.HealthRetries,
				RunUser:           ctl.InstallAgentOptions.RunUser,
			}, c.flushAgentsCh)
			ag.StartRecvCtl()
			ag.ExecCtl(ctlResp)
			c.AddAgent(agentId, ag)

			c.flushAgents()
		case *proto.ControlResponse_UninstallAgentOptions_:
			agentId, _ := uuid.FromBytes(ctl.UninstallAgentOptions.AgentId)
			ag, ok := c.getAgent(agentId)
			if !ok {
				base.Errorf("uninstall %v, but can't find it", agentId)
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.UninstallAgentOptions.AgentId,
					Failed:  true,
					Message: "can't find agent, uninstall fail",
				})
				continue
			}
			if ag.ExecCtl(ctlResp) == nil {
				c.delAgent(agentId)
				ag.StopRecvCtl()
				c.flushAgents()
			}
		case *proto.ControlResponse_CancelOptions_:
			agentId, _ := uuid.FromBytes(ctl.CancelOptions.AgentId)
			a, _ := c.getAgent(agentId)
			if a != nil {
				a.ExecCtl(ctlResp)
			} else {
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.CancelOptions.AgentId,
					Failed:  true,
					Message: "can't find agent, cancel failed",
				})
				continue
			}
		case *proto.ControlResponse_UpdateAgentOptions_:
			agentId, _ := uuid.FromBytes(ctl.UpdateAgentOptions.AgentId)
			ag, ok := c.getAgent(agentId)
			if !ok {
				err := fmt.Errorf("update %v, but can't find it", agentId)
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.UpdateAgentOptions.AgentId,
					Failed:  true,
					Message: err.Error(),
				})
				continue
			}
			ag.ExecCtl(ctlResp)
		case *proto.ControlResponse_StartAgentOptions_:
			agentId, _ := uuid.FromBytes(ctl.StartAgentOptions.AgentId)
			ag, ok := c.getAgent(agentId)
			if !ok {
				err := fmt.Errorf("start %v, but can't find it", agentId)
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.StartAgentOptions.AgentId,
					Failed:  true,
					Message: err.Error(),
				})
				continue
			}
			ag.ExecCtl(ctlResp)
		case *proto.ControlResponse_StopAgentOptions_:
			agentId, _ := uuid.FromBytes(ctl.StopAgentOptions.AgentId)
			ag, ok := c.getAgent(agentId)
			if !ok {
				err := fmt.Errorf("stop %v, but can't find it", agentId)
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.StopAgentOptions.AgentId,
					Failed:  true,
					Message: err.Error(),
				})
				continue
			}
			ag.ExecCtl(ctlResp)
		case *proto.ControlResponse_UpdateAgentConfigOptions_:
			agentId, _ := uuid.FromBytes(ctl.UpdateAgentConfigOptions.AgentId)
			ag, ok := c.getAgent(agentId)
			if !ok {
				err := fmt.Errorf("update config %v, but can't find it", agentId)
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_OperationProgress{
					Seqno:   ctlResp.Seqno,
					AgentId: ctl.UpdateAgentConfigOptions.AgentId,
					Failed:  true,
					Message: err.Error(),
				})
				continue
			}
			ag.ExecCtl(ctlResp)
		case *proto.ControlResponse_ExecScriptOptions_:
			if ctl.ExecScriptOptions.Timeout <= 0 {
				err := fmt.Errorf("timeout must greater than zero")
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_ExecScriptResponse{
					Seqno:    ctlResp.Seqno,
					AgentId:  ctl.ExecScriptOptions.AgentId,
					Failed:   true,
					Response: err.Error(),
				})
				return
			}

			agentId, _ := uuid.FromBytes(ctl.ExecScriptOptions.AgentId)
			if agentId != uuid.Nil {
				ag, ok := c.getAgent(agentId)
				if !ok {
					err := fmt.Errorf("exec script %v, but can't find it", agentId)
					base.Errorf("%v", err)
					event.ReportEvent(&proto.Event_ExecScriptResponse{
						Seqno:    ctlResp.Seqno,
						AgentId:  ctl.ExecScriptOptions.AgentId,
						Failed:   true,
						Response: err.Error(),
					})
					continue
				}
				ag.ExecCtl(ctlResp)
			} else if maxExec := atomic.LoadInt32(&c.maxExec); maxExec > 0 {
				atomic.AddInt32(&c.maxExec, -1)
				go c.execScript(ctlResp)
			} else {
				err := fmt.Errorf("concurrency exec shell more than %v", maxExec)
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_ExecScriptResponse{
					Seqno:    ctlResp.Seqno,
					Failed:   true,
					Response: err.Error(),
				})
			}
		case *proto.ControlResponse_ExecRestOptions_:
			if ctl.ExecRestOptions.Timeout <= 0 {
				err := fmt.Errorf("timeout must greater than zero")
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_ExecRestResponse{
					Seqno:    ctlResp.Seqno,
					AgentId:  ctl.ExecRestOptions.AgentId,
					Failed:   true,
					Response: []byte(err.Error()),
				})
				return
			}
			if maxExec := atomic.LoadInt32(&c.maxExec); maxExec > 0 {
				atomic.AddInt32(&c.maxExec, -1)
				go c.execRest(ctlResp)
			} else {
				err := fmt.Errorf("concurrency exec rest more than %v", maxExec)
				base.Errorf("%v", err)
				event.ReportEvent(&proto.Event_ExecRestResponse{
					Seqno:    ctlResp.Seqno,
					Failed:   true,
					Response: []byte(err.Error()),
				})
			}
		default:
			base.Errorf("unsupport control command: %v", ctlResp.Cmd)
		}
	}
}

func (c *Controller) execScript(ctlResp *proto.ControlResponse) {
	defer atomic.AddInt32(&c.maxExec, 1)

	ctl, ok := ctlResp.Options.(*proto.ControlResponse_ExecScriptOptions_)
	if !ok {
		base.Errorf("it's not ExecShellOptions type")
		return
	}

	ev := &proto.Event_ExecScriptResponse{
		Seqno: ctlResp.Seqno,
	}

	script, err := util.CreateTempScript(ctl.ExecScriptOptions.ExecScript, "exec-script-")
	if err != nil {
		ev.Failed = true
		base.Errorf("%v", err)
		ev.Response = err.Error()
		event.ReportEvent(ev)
		return
	}
	defer os.Remove(script)

	ctx, cancel := context.WithTimeout(context.Background(), ctl.ExecScriptOptions.Timeout)
	defer cancel()
	cmd := util.CommandContext(ctx, "", false, nil, script, ctl.ExecScriptOptions.Parameter...)
	stdBuf := &util.PrefixSuffixSaver{N: 128 << 10}
	cmd.Stdout = stdBuf
	cmd.Stderr = stdBuf
	cmd.Dir = filepath.Dir(script)
	if err := cmd.Run(); err != nil {
		ev.Failed = true
		if _, ok := err.(*exec.ExitError); !ok {
			base.Errorf("%v", err)
			ev.Response = err.Error()
		} else {
			ev.Response = string(stdBuf.Bytes())
		}
	} else {
		ev.Response = string(stdBuf.Bytes())
	}
	event.ReportEvent(ev)
}

func (c *Controller) execRest(ctlResp *proto.ControlResponse) {
	defer atomic.AddInt32(&c.maxExec, 1)

	ctl, ok := ctlResp.Options.(*proto.ControlResponse_ExecRestOptions_)
	if !ok {
		base.Errorf("it's not ExecRestOptions type")
		return
	}

	ev := &proto.Event_ExecRestResponse{
		Seqno: ctlResp.Seqno,
	}
	//ctx, cancel := context.WithTimeout(ag.opCtx, ctl.ExecRestOptions.Timeout)
	//defer cancel()
	// TODO add host:port to config
	url := "http://127.0.0.1:8899/api/v1/" + ctl.ExecRestOptions.Path + ctl.ExecRestOptions.Query
	base.Infof("NewRequest: %v, method: %v", url, ctl.ExecRestOptions.Method)
	base.Infof("NewRequest: %v, body: %+v", url, ctl.ExecRestOptions.Body)
	req, err := http.NewRequest(ctl.ExecRestOptions.Method, url, bytes.NewReader(ctl.ExecRestOptions.Body))
	if err != nil {
		ev.Failed = true
		base.Errorf("NewRequest %v error: %v", url, err)
		ev.Response = []byte(err.Error())
		event.ReportEvent(ev)
		return
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		ev.Failed = true
		base.Errorf("send request error: %v", err)
		ev.Response = []byte(err.Error())
		event.ReportEvent(ev)
		return
	}
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ev.Failed = true
		base.Errorf("read respose nody error: %v", err)
		ev.Response = []byte(err.Error())
		event.ReportEvent(ev)
		return
	}
	base.Infof("response: %v", string(respContent))
	ev.Response = respContent
	resp.Body.Close()
	event.ReportEvent(ev)
}

func (c *Controller) Run() {
	// not lock for init
	for _, ag := range c.agents {
		ag.Run()
	}

	go c.dispatch()
	go c.flushAgentsLoop()
}

func (c *Controller) KillAgents() {
	agents := make([]agent.Agenter, 0)

	c.RLock()
	for _, ag := range c.agents {
		agents = append(agents, ag)
	}
	c.RUnlock()

	for _, ag := range agents {
		ag.Kill()
	}
}

func (c *Controller) flushAgents() {
	c.flushAgentsCh <- struct{}{}
}

func (c *Controller) flushAgentsLoop() {
	for range c.flushAgentsCh {
		agents := make([]config.AgentConfig, 0, 1)

		c.RLock()
		for _, ag := range c.agents {
			agents = append(agents, ag.GetConfig())
		}
		c.RUnlock()

		err := config.WriteAgents(agents, c.agentsFile)
		if err != nil {
			base.Errorf("flushAgents error: %v", err)
			event.ReportEvent(&proto.Event_AgentError{Errstr: err.Error()})
		}
	}
}
