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
	"crypto/tls"
	"easyagent/internal/proto"
	"easyagent/internal/sidecar/base"
	"easyagent/internal/sidecar/config"
	"easyagent/internal/sidecar/controller/util"
	"easyagent/internal/sidecar/event"
	"easyagent/internal/sidecar/monitor"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/satori/go.uuid"
)

var (
	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
)

type Agenter interface {
	// it used for controller run
	Run()
	// it used for controller stop
	Kill()
	// only run goroutine to recv control command
	StartRecvCtl()
	// exit goroutine to stop recv control command
	StopRecvCtl()
	// excute control command serial
	ExecCtl(ctlResp *proto.ControlResponse) error
	// get agent config
	GetConfig() config.AgentConfig
}

type agent struct {
	ctlCh         chan *proto.ControlResponse
	stopRecvCh    chan struct{}
	killCh        chan struct{}
	stopSleepCh   chan struct{}
	flushAgentsCh chan struct{}

	cg      *util.Cgroup
	classid uint32

	wgSpv sync.WaitGroup

	binaryPath        string
	agentId           uuid.UUID
	configurationPath string
	parameter         []string
	workdir           string
	name              string

	wgHealth          sync.WaitGroup
	healthShell       string
	healthPeriod      time.Duration
	healthStartPeriod time.Duration
	healthTimeout     time.Duration
	healthRetries     uint64
	healthStopCh      chan struct{}
	healthCtx         context.Context
	healthCancel      context.CancelFunc

	startParmaMu sync.RWMutex
	cpuLimit     float32
	memLimit     uint64
	netLimit     uint64
	environment  map[string]string
	runUser      string

	ctxMu    sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	opCtx    context.Context
	opCancel context.CancelFunc

	isRunning atomic.Value
	enabled   atomic.Value

	// just for supervisor
	restartCount int
	startTime    time.Time
	stdBuff      *util.PrefixSuffixSaver
}

func NewAgent(cfg config.AgentConfig, flushAgentsCh chan struct{}) Agenter {
	ag := &agent{
		binaryPath:        cfg.BinaryPath,
		agentId:           cfg.AgentId.UUID,
		configurationPath: cfg.ConfigurationPath,
		parameter:         cfg.Parameter,
		workdir:           cfg.Workdir,
		name:              cfg.Name,
		healthShell:       cfg.HealthShell,
		healthPeriod:      cfg.HealthPeriod,
		healthStartPeriod: cfg.HealthStartPeriod,
		healthTimeout:     cfg.HealthTimeout,
		healthRetries:     cfg.HealthRetries,
		cpuLimit:          cfg.CpuLimit,
		memLimit:          cfg.MemLimit,
		netLimit:          cfg.NetLimit,
		environment:       cfg.Environment,
		runUser:           cfg.RunUser,

		ctlCh:         make(chan *proto.ControlResponse),
		stopRecvCh:    make(chan struct{}),
		killCh:        make(chan struct{}),
		flushAgentsCh: flushAgentsCh,
		stdBuff:       &util.PrefixSuffixSaver{N: 256 << 10},
	}
	ag.setStopSleepCh()
	ag.setContextCancel()
	ag.setEnabled(cfg.Enabled)
	ag.isRunning.Store(false)

	if ag.healthShell != "" {
		if ag.healthRetries == 0 {
			ag.healthRetries = 1
		}
		if ag.healthTimeout > ag.healthPeriod || ag.healthTimeout <= 0 {
			ag.healthTimeout = ag.healthPeriod
		}
	}

	return ag
}

// it used for controller run
func (ag *agent) Run() {
	ag.StartRecvCtl()

	if ag.getEnabled() {
		cpuLimit, memLimit, netLimit, env := ag.getStartParam()
		if err := ag.run(cpuLimit, memLimit, netLimit, env); err != nil {
			base.Errorf("%v", err)
		}
	}
}

// it used for controller kill
func (ag *agent) Kill() {
	close(ag.killCh)

	if !ag.isRunning.Load().(bool) {
		return
	}

	base.Infof("killing agent %v", ag.agentId)

	ag.cancelAgent()
	ag.wgSpv.Wait()

	if err := ag.delTcClass(); err != nil {
		base.Errorf("%v", err)
	}
	if err := ag.unInstallCgroup(); err != nil {
		base.Errorf("%v", err)
	}
}

func (ag *agent) cancelAgent() {
	ag.ctxMu.RLock()
	defer ag.ctxMu.RUnlock()

	ag.cancel()
}

func (ag *agent) getCtxAgent() context.Context {
	ag.ctxMu.RLock()
	defer ag.ctxMu.RUnlock()

	return ag.ctx
}

// only run goroutine to recv control command
func (ag *agent) StartRecvCtl() {
	wait := make(chan struct{}, 0)

	go func() {
		defer base.Debugf("agent %v stop receive ctlResp", ag.agentId)

		base.Debugf("agent %v start receive ctlResp...", ag.agentId)

		ready := make(chan struct{}, 1)
		for {
			select {
			case ctlResp := <-ag.ctlCh:
				switch ctl := ctlResp.Options.(type) {
				case *proto.ControlResponse_InstallAgentOptions_:
					ag.StartOperationCtx()
					ag.install(ctl, ctlResp.Seqno)
				case *proto.ControlResponse_UninstallAgentOptions_:
					ag.StartOperationCtx()
					ag.uninstall(ctl, ctlResp.Seqno)
				case *proto.ControlResponse_UpdateAgentOptions_:
					ag.StartOperationCtx()
					ag.update(ctl, ctlResp.Seqno)
				case *proto.ControlResponse_StartAgentOptions_:
					ag.start(ctl, ctlResp.Seqno)
				case *proto.ControlResponse_StopAgentOptions_:
					ag.stop(ctl, ctlResp.Seqno)
				case *proto.ControlResponse_UpdateAgentConfigOptions_:
					ag.StartOperationCtx()
					ag.updateConfig(ctl, ctlResp.Seqno)
				case *proto.ControlResponse_ExecScriptOptions_:
					ag.StartOperationCtx()
					ag.execScript(ctl, ctlResp.Seqno)
				case *proto.ControlResponse_ExecRestOptions_:
					ag.StartOperationCtx()
					ag.execRest(ctl, ctlResp.Seqno)
				default:
					base.Errorf("unknow control command type")
				}
			case <-ag.stopRecvCh:
				return
			case <-ag.killCh:
				return
			case ready <- struct{}{}:
				close(wait)
			}
		}
	}()

	// we should make sure ctlCh is receiving
	<-wait
}

// exit goroutine to stop recv control command
func (ag *agent) StopRecvCtl() {
	close(ag.stopRecvCh)
}

// excute control command serial
func (ag *agent) ExecCtl(ctlResp *proto.ControlResponse) (err error) {
	// won't be dropped
	switch ctlResp.Options.(type) {
	case *proto.ControlResponse_CancelOptions_:
		ag.CancelOperation()
		ev := &proto.Event_OperationProgress{
			Seqno: ctlResp.Seqno,
		}
		ev.Failed = true
		event.ReportEvent(ev)
		return
	}

	select {
	case ag.ctlCh <- ctlResp:
	case <-ag.stopRecvCh:
		err = fmt.Errorf("agent is stopping recevice, droped")
	case <-ag.killCh:
		err = fmt.Errorf("sidecar is be killing, droped")
	default:
		err = fmt.Errorf("agent is running task, droped")
		ev := &proto.Event_OperationProgress{
			Seqno:   ctlResp.Seqno,
			AgentId: ag.agentId.Bytes(),
			Failed:  true,
			Message: err.Error(),
		}
		event.ReportEvent(ev)
	}

	if err != nil {
		base.Errorf("%v", err)
	}
	return
}

// get agent enabled
func (ag *agent) getEnabled() bool {
	return ag.enabled.Load().(bool)
}

// set agent enabled
func (ag *agent) setEnabled(enabled bool) {
	ag.enabled.Store(enabled)
}

func (ag *agent) getStartParam() (cpuLimit float32, memLimit uint64, netLimit uint64, env map[string]string) {
	ag.startParmaMu.RLock()
	cpuLimit, memLimit, netLimit, env = ag.cpuLimit, ag.memLimit, ag.netLimit, ag.environment
	ag.startParmaMu.RUnlock()

	return
}

func (ag *agent) setStartParam(cpuLimit float32, memLimit uint64, netLimit uint64, env map[string]string) {
	ag.startParmaMu.Lock()
	ag.cpuLimit, ag.memLimit, ag.netLimit, ag.environment = cpuLimit, memLimit, netLimit, env
	ag.startParmaMu.Unlock()
}

// get agent config
func (ag *agent) GetConfig() config.AgentConfig {
	cpuLimit, memLimit, netLimit, env := ag.getStartParam()
	return config.AgentConfig{
		BinaryPath:        ag.binaryPath,
		AgentId:           config.WrapperUUID{UUID: ag.agentId},
		ConfigurationPath: ag.configurationPath,
		Parameter:         ag.parameter,
		Workdir:           ag.workdir,
		Name:              ag.name,
		HealthShell:       ag.healthShell,
		HealthPeriod:      ag.healthPeriod,
		HealthStartPeriod: ag.healthStartPeriod,
		HealthTimeout:     ag.healthTimeout,
		HealthRetries:     ag.healthRetries,
		Enabled:           ag.getEnabled(),
		CpuLimit:          cpuLimit,
		MemLimit:          memLimit,
		NetLimit:          netLimit,
		Environment:       env,
		RunUser:           ag.runUser,
	}
}

func (ag *agent) isKilled() bool {
	select {
	case <-ag.killCh:
		return true
	default:
		return false
	}
}

func (ag *agent) supervisor(cmd *util.Cmd) {
	reqStop := false
	oldEnv := cmd.Env
	oldAttr := cmd.SysProcAttr

	defer func() {
		ag.isRunning.Store(false)
		ag.wgSpv.Done()

		err := fmt.Errorf("stop supervisor: %v", ag.agentId)
		base.Debugf("%v", err)
		if !reqStop || ag.isKilled() {
			event.ReportEvent(&proto.Event_AgentError{AgentId: ag.agentId.Bytes(), Errstr: "agent run error(unexpected):" + err.Error()})
		}
	}()

	base.Debugf("start supervisor: %v", ag.agentId)

	ag.restartCount = 0
	ag.startTime = time.Now()

	for {
		monitor.StartMonitAgent(ag.agentId, cmd.Process.Pid, ag.classid)
		ag.startHealthCheck()

		err := cmd.Wait()
		base.Infof("agent %v exit(%v),std error: %v", ag.agentId, err, string(ag.stdBuff.Bytes()))
		ag.stopHealthCheck()
		monitor.StopMonitAgent(ag.agentId)
		if err == context.Canceled {
			// request stop
			reqStop = true
			return
		}

		ag.sleep(2 * time.Second)

		for {
			// after 60 seconds we can reset the restart counter
			if time.Since(ag.startTime) > 60*time.Second {
				ag.restartCount = 0
			}
			// don't continue to restart after 3 tries, exit the supervisor and wait for a configuration update
			// or manual restart
			if ag.restartCount >= 3 {
				return
			}

			binaryPath := ag.binaryPath
			if !filepath.IsAbs(ag.binaryPath) {
				path := filepath.Join(ag.workdir, ag.binaryPath)
				if _, err := os.Stat(path); os.IsExist(err) {
					binaryPath = path
				}
			}
			cmd = util.CommandContext(ag.getCtxAgent(), ag.runUser, true, ag.cg, binaryPath, ag.parameter...)
			if ag.workdir == "" {
				if cmd.Dir = filepath.Dir(ag.binaryPath); cmd.Dir == "." {
					cmd.Dir = os.TempDir()
				}
			} else {
				cmd.Dir = ag.workdir
			}
			cmd.Env = oldEnv
			cmd.SysProcAttr = oldAttr
			ag.restartCount++
			ag.startTime = time.Now()
			if err = cmd.Start(); err == nil {
				// start success
				break
			} else if err == context.Canceled {
				// request stop
				reqStop = true
				return
			}
			// start fail
			base.Errorf("run agent %v error: %v", ag.agentId, err)
			event.ReportEvent(&proto.Event_AgentError{AgentId: ag.agentId.Bytes(), Errstr: "run agent error: " + err.Error()})
			ag.sleep(2 * time.Second)
		}
		base.Infof("run agent %v success(pid: %d)", ag.agentId, cmd.Process.Pid)
	}
}

func (ag *agent) install(ctl *proto.ControlResponse_InstallAgentOptions_, seqno uint32) {
	ev := &proto.Event_OperationProgress{
		Seqno:   seqno,
		AgentId: ctl.InstallAgentOptions.AgentId,
	}

	script, err := util.CreateTempScript(ctl.InstallAgentOptions.InstallScript, "install-script-")
	if err != nil {
		ev.Failed = true
		base.Errorf("%v", err)
		ev.Message = err.Error()
		event.ReportEvent(ev)
		return
	}
	defer os.Remove(script)

	ctx, cancel := context.WithTimeout(ag.opCtx, ctl.InstallAgentOptions.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, script, ctl.InstallAgentOptions.InstallParameter...)
	stdBuf := &util.PrefixSuffixSaver{N: 64 << 10}
	cmd.Stdout = stdBuf
	cmd.Stderr = stdBuf
	cmd.Dir = filepath.Dir(script)
	if err := cmd.Run(); err != nil {
		ev.Failed = true
		if _, ok := err.(*exec.ExitError); !ok {
			base.Errorf("%v", err)
			ev.Message = err.Error()
		} else {
			ev.Message = string(stdBuf.Bytes())
		}
	} else {
		ev.Message = string(stdBuf.Bytes())
	}

	event.ReportEvent(ev)
}

func (ag *agent) uninstall(ctl *proto.ControlResponse_UninstallAgentOptions_, seqno uint32) {
	ev := &proto.Event_OperationProgress{
		Seqno:   seqno,
		AgentId: ctl.UninstallAgentOptions.AgentId,
	}

	script, err := util.CreateTempScript(ctl.UninstallAgentOptions.UninstallScript, "uninstall-script-")
	if err != nil {
		ev.Failed = true
		base.Errorf("%v", err)
		ev.Message = err.Error()
		event.ReportEvent(ev)
		return
	}
	defer os.Remove(script)

	ctx, cancel := context.WithTimeout(ag.opCtx, ctl.UninstallAgentOptions.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, script, ctl.UninstallAgentOptions.Parameter...)
	stdBuf := &util.PrefixSuffixSaver{N: 64 << 10}
	cmd.Stdout = stdBuf
	cmd.Stderr = stdBuf
	if ag.workdir == "" {
		if cmd.Dir = filepath.Dir(ag.binaryPath); cmd.Dir == "." {
			cmd.Dir = os.TempDir()
		}
	} else {
		cmd.Dir = ag.workdir
	}
	if _, err = os.Stat(cmd.Dir); os.IsNotExist(err) {
		cmd.Dir = filepath.Dir(script)
	}
	if err := cmd.Run(); err != nil {
		ev.Failed = true
		if _, ok := err.(*exec.ExitError); !ok {
			base.Errorf("%v", err)
			ev.Message = err.Error()
		} else {
			ev.Message = string(stdBuf.Bytes())
		}
	} else {
		ev.Message = string(stdBuf.Bytes())
	}

	event.ReportEvent(ev)
}

func (ag *agent) update(ctl *proto.ControlResponse_UpdateAgentOptions_, seqno uint32) {
	ev := &proto.Event_OperationProgress{
		Seqno:   seqno,
		AgentId: ctl.UpdateAgentOptions.AgentId,
	}

	script, err := util.CreateTempScript(ctl.UpdateAgentOptions.UpdateScript, "update-script-")
	if err != nil {
		ev.Failed = true
		base.Errorf("%v", err)
		ev.Message = err.Error()
		event.ReportEvent(ev)
		return
	}
	defer os.Remove(script)

	ctx, cancel := context.WithTimeout(ag.opCtx, ctl.UpdateAgentOptions.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, script, ctl.UpdateAgentOptions.Parameter...)
	stdBuf := &util.PrefixSuffixSaver{N: 64 << 10}
	cmd.Stdout = stdBuf
	cmd.Stderr = stdBuf
	if ag.workdir == "" {
		if cmd.Dir = filepath.Dir(ag.binaryPath); cmd.Dir == "." {
			cmd.Dir = os.TempDir()
		}
	} else {
		cmd.Dir = ag.workdir
	}
	if err := cmd.Run(); err != nil {
		ev.Failed = true
		if _, ok := err.(*exec.ExitError); !ok {
			base.Errorf("%v", err)
			ev.Message = err.Error()
		} else {
			ev.Message = string(stdBuf.Bytes())
		}
	} else {
		ev.Message = string(stdBuf.Bytes())
	}
	event.ReportEvent(ev)
}

// it wrapper run
func (ag *agent) start(ctl *proto.ControlResponse_StartAgentOptions_, seqno uint32) {
	ev := &proto.Event_OperationProgress{
		Seqno:   seqno,
		AgentId: ctl.StartAgentOptions.AgentId,
	}
	if err := ag.run(
		ctl.StartAgentOptions.CpuLimit,
		ctl.StartAgentOptions.MemLimit,
		ctl.StartAgentOptions.NetLimit,
		ctl.StartAgentOptions.Environment,
	); err != nil {
		ev.Failed = true
		base.Errorf("%v", err)
		ev.Message = err.Error()
	} else {
		ag.setStartParam(
			ctl.StartAgentOptions.CpuLimit,
			ctl.StartAgentOptions.MemLimit,
			ctl.StartAgentOptions.NetLimit,
			ctl.StartAgentOptions.Environment,
		)
		ag.setEnabled(true)
		ag.flushAgents()
	}
	event.ReportEvent(ev)
}

// it start agent and run supervisor goroutine
func (ag *agent) run(cpuLimit float32, memLimit, netLimit uint64, environment map[string]string) error {
	if ag.isRunning.Load().(bool) {
		return fmt.Errorf("agent %v is already running", ag.agentId)
	}

	if err := ag.installCgroup(); err != nil {
		base.Errorf("agent %v install cgroup error: %v", ag.agentId, err)
	}
	if err := ag.updateCgroup(cpuLimit, memLimit); err != nil {
		base.Errorf("agent %v update cgroup error: %v", ag.agentId, err)
	}
	if err := ag.setTcClassRate(netLimit); err != nil {
		base.Errorf("agent %v set tc class error: %v", ag.agentId, err)
	}

	binaryPath := ag.binaryPath
	if !filepath.IsAbs(ag.binaryPath) {
		path := filepath.Join(ag.workdir, ag.binaryPath)
		if _, err := os.Stat(path); os.IsExist(err) {
			binaryPath = path
		}
	}
	cmd := util.CommandContext(ag.getCtxAgent(), ag.runUser, true, ag.cg, binaryPath, ag.parameter...)
	cmd.Stdout = ag.stdBuff
	cmd.Stderr = ag.stdBuff
	if ag.workdir == "" {
		if cmd.Dir = filepath.Dir(ag.binaryPath); cmd.Dir == "." {
			cmd.Dir = os.TempDir()
		}
	} else {
		cmd.Dir = ag.workdir
	}
	if ag.runUser != "" {
		u, err := user.Lookup(ag.runUser)
		if err != nil {
			return fmt.Errorf("agent %v run with user %v error: %v", ag.agentId, ag.runUser, err)
		}
		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}
		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			return fmt.Errorf("agent %v run with user %v error: %v", ag.agentId, ag.runUser, err)
		}
		gid, err := strconv.Atoi(u.Gid)
		if err != nil {
			return fmt.Errorf("agent %v run with user %v error: %v", ag.agentId, ag.runUser, err)
		}
		base.Infof("run agent with user id : %v, gid: %v", uid, gid)
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}
	cmd.Env = util.Env(environment)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("run agent %v error: %v", ag.agentId, err)
	}
	base.Infof("run agent %v success(pid: %d),with user: %v", ag.agentId, cmd.Process.Pid, ag.runUser)

	ag.isRunning.Store(true)

	// supervisor goroutine
	ag.wgSpv.Add(1)
	go ag.supervisor(cmd)

	return nil
}

func (ag *agent) stop(ctl *proto.ControlResponse_StopAgentOptions_, seqno uint32) {
	ev := &proto.Event_OperationProgress{
		Seqno:   seqno,
		AgentId: ctl.StopAgentOptions.AgentId,
	}
	if ag.isRunning.Load().(bool) {
		close(ag.stopSleepCh) // quit sleep in supervisor
		ag.cancelAgent()
		ag.wgSpv.Wait()

		if err := ag.delTcClass(); err != nil {
			base.Errorf("%v", err)
		}
		if err := ag.unInstallCgroup(); err != nil {
			base.Errorf("%v", err)
		}

		// reset stopSleepCh
		ag.setStopSleepCh()
		// reset context
		ag.setContextCancel()

		if ctl.StopAgentOptions.StopAgentOptionsType == proto.STOP_UNRECOVER {
			ag.setEnabled(false)
		} else if ctl.StopAgentOptions.StopAgentOptionsType == proto.STOP_RECOVER {
			ag.setEnabled(true)
		}

		ag.flushAgents()
	}

	event.ReportEvent(ev)
}

func (ag *agent) setContextCancel() {
	ag.ctxMu.Lock()
	defer ag.ctxMu.Unlock()

	ag.ctx, ag.cancel = context.WithCancel(context.Background())
}

func (ag *agent) CancelOperation() {
	if ag.opCancel != nil {
		ag.opCancel()
	}
}

func (ag *agent) StartOperationCtx() {
	ag.opCtx, ag.opCancel = context.WithCancel(context.Background())
}

func (ag *agent) flushAgents() {
	ag.flushAgentsCh <- struct{}{}
}

func (ag *agent) updateConfig(ctl *proto.ControlResponse_UpdateAgentConfigOptions_, seqno uint32) {
	ev := &proto.Event_OperationProgress{
		Seqno:   seqno,
		AgentId: ctl.UpdateAgentConfigOptions.AgentId,
	}
	configurationPath := ag.configurationPath
	if ctl.UpdateAgentConfigOptions.ConfigPath != "" {
		configurationPath = ctl.UpdateAgentConfigOptions.ConfigPath
	}
	if !filepath.IsAbs(configurationPath) {
		configurationPath = filepath.Join(ag.workdir, configurationPath)
	}
	base.Infof("updateConfig path %v", configurationPath)
	if err := ioutil.WriteFile(configurationPath, []byte(ctl.UpdateAgentConfigOptions.ConfigContent), 0600); err != nil {
		ev.Failed = true
		ev.Message = err.Error()
	}
	event.ReportEvent(ev)
}
func (ag *agent) execRest(ctl *proto.ControlResponse_ExecRestOptions_, seqno uint32) {
	ev := &proto.Event_ExecRestResponse{
		Seqno:   seqno,
		AgentId: ctl.ExecRestOptions.AgentId,
	}
	//ctx, cancel := context.WithTimeout(ag.opCtx, ctl.ExecRestOptions.Timeout)
	//defer cancel()
	// TODO add host:port to config
	url := "http:/127.0.0.1:8899/api/v1/" + ctl.ExecRestOptions.Path + ctl.ExecRestOptions.Query
	req, err := http.NewRequest(ctl.ExecRestOptions.Method, url, nil)
	if err != nil {
		base.Errorf("NewRequest %v error: %v", url, err)
		ev.Response = []byte(err.Error())
		event.ReportEvent(ev)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		base.Errorf("send request error: %v", err)
		ev.Response = []byte(err.Error())
		event.ReportEvent(ev)
		return
	}
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		base.Errorf("read respose nody error: %v", err)
		ev.Response = []byte(err.Error())
		event.ReportEvent(ev)
		return
	}
	ev.Response = respContent
	resp.Body.Close()
	event.ReportEvent(ev)
}

func (ag *agent) execScript(ctl *proto.ControlResponse_ExecScriptOptions_, seqno uint32) {
	ev := &proto.Event_ExecScriptResponse{
		Seqno:   seqno,
		AgentId: ctl.ExecScriptOptions.AgentId,
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

	ctx, cancel := context.WithTimeout(ag.opCtx, ctl.ExecScriptOptions.Timeout)
	defer cancel()
	cmd := util.CommandContext(ctx, "", false, nil, script, ctl.ExecScriptOptions.Parameter...)
	stdBuf := &util.PrefixSuffixSaver{N: 128 << 10}
	cmd.Stdout = stdBuf
	cmd.Stderr = stdBuf
	if ag.workdir == "" {
		if cmd.Dir = filepath.Dir(ag.binaryPath); cmd.Dir == "." {
			cmd.Dir = os.TempDir()
		}
	} else {
		cmd.Dir = ag.workdir
	}
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

// sleep d duration, except stop signal
func (ag *agent) sleep(d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-ag.killCh:
	case <-ag.stopSleepCh:
	}
}

func (ag *agent) setStopSleepCh() {
	ag.stopSleepCh = make(chan struct{})
}
