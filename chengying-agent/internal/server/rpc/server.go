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

package rpc

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	dbhelper "easyagent/go-common/db-helper"
	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	"easyagent/internal/server/publisher"
	. "easyagent/internal/server/tracy"

	"github.com/prometheus/client_golang/prometheus"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

var (
	ErrUnregitered        = errors.New("unregistered sidecar")
	ErrInstalled          = errors.New("uninstalled sidecar")
	ErrEndReadyForControl = errors.New("end ReadyForControl")
	ErrAlreadyConnected   = errors.New("already connected")
	ErrUnsupportEventType = errors.New("unsupport event type")

	registerSidecarTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "register_sidecar_total",
		Help: "Total Number of RegisterSidecar",
	})
	readyForControlTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ready_for_control_total",
		Help: "Total Number of ReadyForControl",
	})
	reportEventTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "report_event_total",
		Help: "Total Number of ReportEvent",
	})
	rpcServerErrorTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rpc_server_error_total",
		Help: "Total Number of rpc server errors",
	})
	agentErrorTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "agent_error_total",
		Help: "Total Number of agent error",
	})

	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "cpu usage",
		},
		[]string{"agent_id"},
	)
	memory = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory",
			Help: "memory",
		},
		[]string{"agent_id"},
	)
	bytesSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_sent",
			Help: "bytes sent",
		},
		[]string{"agent_id"},
	)

	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
)

const TsLayout = "2006-01-02 15:04:05.000000"

type sidecar struct {
	sid    uuid.UUID
	ctlCh  chan *proto.ControlResponse
	stopCh chan struct{}
}

type rpcService struct {
	sync.RWMutex
	sidecarMap map[uuid.UUID]sidecar
	apiHost    string
	apiPort    int
}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(registerSidecarTotal)
	prometheus.MustRegister(readyForControlTotal)
	prometheus.MustRegister(reportEventTotal)
	prometheus.MustRegister(rpcServerErrorTotal)
	prometheus.MustRegister(agentErrorTotal)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memory)
	prometheus.MustRegister(bytesSent)
}

func (rpc *rpcService) RegisterSidecar(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	registerSidecarTotal.Inc()

	sid, _ := uuid.FromBytes(request.Id)
	InstallProgressLog("[INSTALL] RegisterSidecar sid :%v...", sid)

	p, _ := peer.FromContext(ctx)
	clientIP, _, _ := net.SplitHostPort(p.Addr.String())

	sidecarInfo, err := model.SidecarList.GetSidecarInfo(sid)
	if err != nil {
		rpcServerErrorTotal.Inc()
		log.Errorf("register sidecar %v(%v) error: %v", sid, clientIP, err)
		InstallProgressLog("[INSTALL] RegisterSidecar register sidecar error: %v, sid: %v", err, sid)
		return nil, err
	}
	if sidecarInfo.Status == model.SIDECAR_NOT_INSTALLED {
		rpcServerErrorTotal.Inc()
		log.Errorf("register sidecar %v(%v) error: %v", sid, clientIP, ErrInstalled)
		InstallProgressLog("[INSTALL] RegisterSidecar SIDECAR_NOT_INSTALLED register sidecar error: %v, sid %v", err, sid)
		return nil, ErrInstalled
	}

	log.Debugf("sidecar %v(%v) registerd success", sid, clientIP)
	InstallProgressLog("[INSTALL] RegisterSidecar sidecar registerd success, sid: %v", sid)

	updateFields := dbhelper.UpdateFields{
		"last_update_date": time.Now(),
		"os_type":          request.OsType,
		"os_platform":      request.OsPlatform,
		"os_version":       request.OsVersion,
		"cpu_serial":       request.CpuSerial,
		"cpu_cores":        request.CpuCores,
		"mem_size":         request.MemSize,
		"swap_size":        request.SwapSize,
		"host":             request.Host,
		"server_host":      rpc.apiHost,
		"server_port":      rpc.apiPort,
		"ssh_host":         clientIP,
		"local_ip":         request.LocalIp,
	}
	if err := model.SidecarList.UpdateSidecar(sid, updateFields); err != nil {
		rpcServerErrorTotal.Inc()
		log.Errorf("RegisterSidecar UpdateSidecar %v(%v) error: %v", sid, clientIP, err)
		InstallProgressLog("[INSTALL] RegisterSidecar UpdateSidecar error: %v, sid: %v", err, sid)
		return nil, err
	}

	for _, callback := range request.CallBack {
		go func(callback string) {
			v, err := url.ParseQuery(callback)
			if err != nil {
				log.Errorf("ParseQuery %v error: %v", callback, err)
				InstallProgressLog("[INSTALL] RegisterSidecar ParseQuery %v error: %v", callback, err)
				return
			}
			callbackUrl := v.Get("CallBack")
			if callbackUrl == "" {
				log.Errorf("CallBack is empty!")
				InstallProgressLog("[INSTALL] RegisterSidecar CallBack is empty, sid: %v", sid)
				return
			}
			v.Del("CallBack")

			q, err := url.QueryUnescape(callbackUrl)
			if err != nil {
				log.Errorf("QueryUnescape %v error: %v", callbackUrl, err)
				InstallProgressLog("[INSTALL] RegisterSidecar SshHost or CallBack is empty, sid: %v", sid)
				return
			}
			b, err := base64.StdEncoding.DecodeString(q)
			if err != nil {
				InstallProgressLog("[INSTALL] RegisterSidecar DecodeString %v error: %v, sid: %v", callbackUrl, err, sid)
				log.Errorf("DecodeString %v error: %v", callbackUrl, err)
				return
			}
			callbackUrl = string(b)
			log.Debugf("ready to send %v", callbackUrl)
			InstallProgressLog("[INSTALL] RegisterSidecar ready to send %v, sid: %v", callbackUrl, sid)
			req, err := http.NewRequest("GET", callbackUrl, nil)
			if err != nil {
				log.Errorf("NewRequest %v error: %v", callbackUrl, err)
				InstallProgressLog("[INSTALL] RegisterSidecar NewRequest %v error: %v, sid", callbackUrl, err, sid)
				return
			}
			for hk := range v {
				req.Header.Set(hk, v.Get(hk))
			}
			req.Header.Set("SID", sid.String())
			req.Header.Set("HostName", request.Host)
			req.Header.Set("IP", request.LocalIp)

			for !rpc.IsClientExist(sid) {
				log.Debugf("%v waiting ReadyForControl rpc...", sid)
				InstallProgressLog("[INSTALL] %v waiting ReadyForControl rpc...", sid)
				time.Sleep(1 * time.Second)
			}

			headers, _ := json.Marshal(req.Header)
			log.Debugf("Now Request Header: %s", headers)
			InstallProgressLog("[INSTALL] RegisterSidecar callback headers%s, sid: %v", headers, sid)

			resp, err := client.Do(req)
			if err != nil {
				InstallProgressLog("[INSTALL] RegisterSidecar send request error: %v, sid: %v", err, sid)
				log.Errorf("send request error: %v", err)
				return
			}
			respContent, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("read respose nody error: %v", err)
				InstallProgressLog("[INSTALL] RegisterSidecar read respose nody error: %v", err)
			}
			resp.Body.Close()

			InstallProgressLog("[INSTALL] RegisterSidecarsend get resp code: %v, sid: %v, response: %v", resp.StatusCode, sid, string(respContent))

			log.Debugf("get resp code: %v, response content: %v", resp.StatusCode, string(respContent))
		}(callback)
	}

	return &proto.RegisterResponse{}, nil
}

func (rpc *rpcService) ReadyForControl(request *proto.ControlRequest, stream proto.EasyAgentService_ReadyForControlServer) error {
	readyForControlTotal.Inc()

	sid, _ := uuid.FromBytes(request.Id)

	p, _ := peer.FromContext(stream.Context())
	clientIP, _, _ := net.SplitHostPort(p.Addr.String())

	sidecarInfo, err := model.SidecarList.GetSidecarInfo(sid)
	if err != nil {
		rpcServerErrorTotal.Inc()
		log.Errorf("ReadyForControl sidecar %v(%v) error: %v", sid, clientIP, err)
		return err
	}
	if sidecarInfo.Status == model.SIDECAR_NOT_INSTALLED {
		rpcServerErrorTotal.Inc()
		log.Errorf("ReadyForControl sidecar %v(%v) error: %v", sid, clientIP, ErrInstalled)
		return ErrInstalled
	}

	rpc.Lock()
	sc, ok := rpc.sidecarMap[sid]
	if !ok {
		sc = sidecar{
			sid:    sid,
			ctlCh:  make(chan *proto.ControlResponse),
			stopCh: make(chan struct{}),
		}
		rpc.sidecarMap[sid] = sc
	} else {
		rpc.Unlock()
		rpcServerErrorTotal.Inc()
		log.Errorf("sidecar %v(%v) already connected", sid, clientIP)
		return ErrAlreadyConnected
	}
	rpc.Unlock()

	defer func() {
		rpc.Lock()
		delete(rpc.sidecarMap, sid)
		rpc.Unlock()

		close(sc.stopCh)
	}()

	for {
		select {
		case <-stream.Context().Done():
			err := stream.Context().Err()
			rpcServerErrorTotal.Inc()
			log.Errorf("sidecar %v(%v) stream error: %v", sid, clientIP, err)
			return err
		case ctl, ok := <-sc.ctlCh:
			if !ok {
				rpcServerErrorTotal.Inc()
				log.Infof("sidecar %v(%v) closed", sid, clientIP)
				return ErrEndReadyForControl
			}
			if err := stream.Send(ctl); err != nil {
				rpcServerErrorTotal.Inc()
				log.Errorf("sidecar %v(%v) stream send error: %v", sid, clientIP, err)
				return err
			}
			log.Debugf("sidecar %v(%v) stream send success: %v", sid, clientIP, ctl)
		}
	}
}

func (rpc *rpcService) ReportEvent(ctx context.Context, event *proto.Event) (*proto.EmptyResponse, error) {
	reportEventTotal.Inc()

	sid, _ := uuid.FromBytes(event.Id)

	rpc.RLock()
	_, ok := rpc.sidecarMap[sid]
	rpc.RUnlock()
	if !ok {
		rpcServerErrorTotal.Inc()
		log.Errorf("sidecar %v ReportEvent error: %v", sid, ErrUnregitered)
		return nil, ErrUnregitered
	}

	log.Debugf("sidecar %v ReportEvent %v", sid, event)

	switch inst := event.Details.(type) {
	case *proto.Event_OpProgress:
		opResult := 0
		if inst.OpProgress.Failed {
			opResult = 1
		}
		updateFields := dbhelper.UpdateFields{
			"finish_time":   time.Now(),
			"op_result":     opResult,
			"op_return_msg": inst.OpProgress.Message,
		}
		if err := model.AgentOperation.UpdateAgentOperation(inst.OpProgress.Seqno, updateFields); err != nil {
			log.Errorf("seq %d UpdateAgentOperation error: %v", inst.OpProgress.Seqno, err)
		}
		stopSeqno(inst.OpProgress.Seqno, inst.OpProgress)
	case *proto.Event_AgentError_:
		agentErrorTotal.Inc()
		agentId, _ := uuid.FromBytes(inst.AgentError.AgentId)
		if agentId == uuid.Nil {
			log.Errorf("sidecar %v error: %v", sid, inst.AgentError.Errstr)
		} else {
			log.Errorf("agent %v error: %v", inst.AgentError.AgentId, inst.AgentError.Errstr)
		}
		index := "dtlog-1-easymanager-nodelete-" + time.Now().Format("2006.01.02") + "_000001.alias"
		inst.AgentError.Errstr = "agent异常退出:" + inst.AgentError.Errstr
		body := struct {
			Sid uuid.UUID `json:"sid"`
			*proto.Event_AgentError
			LastUpdate string `json:"last_update_date"`
		}{sid, inst.AgentError, time.Now().Format(TsLayout)}
		if err := publisher.Publish.OutputJson(ctx, agentId.String(), index, "dt_agent_error",
			body, []byte(sid.String())); err != nil {
			log.Errorf("sidecar %v OutputJson Event_AgentError error: %v", sid, err)
		}
	case *proto.Event_OsResourceUsages_:
		diskUsage, err := json.Marshal(inst.OsResourceUsages.DiskUsage)
		if err != nil {
			log.Errorf("sidecar %v Marshal DiskUsage error: %v", sid, err)
		}
		var diskSize, diskUsed, fileSize, fileUsed int64
		for _, disk := range inst.OsResourceUsages.DiskUsage {
			if disk.MountPoint != "/" {
				fileSize += int64(disk.TotalSpace)
				fileUsed += int64(disk.UsedSpace)
			} else {
				diskSize += int64(disk.TotalSpace)
				diskUsed += int64(disk.UsedSpace)
			}
		}
		// 计算disk_usage_pct
		diskUsagePct := float64(diskUsed) / float64(diskSize)
		netUsage, err := json.Marshal(inst.OsResourceUsages.NetStats)
		if err != nil {
			log.Errorf("sidecar %v Marshal NetStats error: %v", sid, err)
		}
		index := "dtlog-1-easymanager-nodelete-" + time.Now().Format("2006.01.02") + "_000001.alias"
		body := struct {
			Sid uuid.UUID `json:"sid"`
			*proto.Event_OsResourceUsages
			LastUpdate string `json:"last_update_date"`
		}{sid, inst.OsResourceUsages, time.Now().Format(TsLayout)}
		if err = publisher.Publish.OutputJson(ctx, sid.String(), index, "dt_agent_host_resource",
			body, []byte(sid.String())); err != nil {
			log.Errorf("sidecar %v OutputJson Event_OsResourceUsages error: %v", sid, err)
		}
		updateFields := dbhelper.UpdateFields{
			"last_update_date": time.Now(),
			"cpu_usage":        inst.OsResourceUsages.CpuUsage * 100,
			"mem_usage":        inst.OsResourceUsages.MemUsage,
			"swap_usage":       inst.OsResourceUsages.SwapUsage,
			"load1":            inst.OsResourceUsages.Load1,
			"uptime":           inst.OsResourceUsages.Uptime,
			"disk_usage":       diskUsage,
			"disk_usage_pct":   diskUsagePct,
			"net_usage":        netUsage,
		}
		if err = model.SidecarList.UpdateSidecar(sid, updateFields); err != nil {
			log.Errorf("sidecar %v update Event_OsResourceUsages error: %v", sid, err)
		}
	case *proto.Event_ProcResourceUsages:
		agentId, _ := uuid.FromBytes(inst.ProcResourceUsages.AgentId)
		cpuUsage.With(prometheus.Labels{"agent_id": agentId.String()}).Set(float64(inst.ProcResourceUsages.CpuUsage))
		memory.With(prometheus.Labels{"agent_id": agentId.String()}).Set(float64(inst.ProcResourceUsages.Memory))
		bytesSent.With(prometheus.Labels{"agent_id": agentId.String()}).Set(float64(inst.ProcResourceUsages.BytesSent))
		index := "dtlog-1-easymanager-nodelete-" + time.Now().Format("2006.01.02") + "_000001.alias"
		body := struct {
			Sid uuid.UUID `json:"sid"`
			*proto.Event_ProcessResourceUsages
			LastUpdate string `json:"last_update_date"`
		}{sid, inst.ProcResourceUsages, time.Now().Format(TsLayout)}
		if err := publisher.Publish.OutputJson(ctx, agentId.String(), index, "dt_agent_performance",
			body, nil); err != nil {
			log.Errorf("sidecar %v OutputJson Event_ProcessResourceUsages error: %v", sid, err)
		}
	case *proto.Event_ExecScriptResponse_:
		opResult := 0
		if inst.ExecScriptResponse.Failed {
			opResult = 1
		}
		updateFields := dbhelper.UpdateFields{
			"finish_time":   time.Now(),
			"op_result":     opResult,
			"op_return_msg": inst.ExecScriptResponse.Response,
		}
		if err := model.AgentOperation.UpdateAgentOperation(inst.ExecScriptResponse.Seqno, updateFields); err != nil {
			log.Errorf("seq %d UpdateAgentOperation error: %v", inst.ExecScriptResponse.Seqno, err)
		}
		stopSeqno(inst.ExecScriptResponse.Seqno, inst.ExecScriptResponse)
	case *proto.Event_ExecRestResponse_:
		stopSeqno(inst.ExecRestResponse.Seqno, inst.ExecRestResponse)
	case *proto.Event_AgentHealthCheck_:
		agentId, _ := uuid.FromBytes(inst.AgentHealthCheck.AgentId)
		index := "dtlog-1-easymanager-nodelete-" + time.Now().Format("2006.01.02") + "_000001.alias"
		body := struct {
			Sid uuid.UUID `json:"sid"`
			*proto.Event_AgentHealthCheck
			LastUpdate string `json:"last_update_date"`
		}{sid, inst.AgentHealthCheck, time.Now().Format(TsLayout)}
		if err := publisher.Publish.OutputJson(ctx, agentId.String(), index, "dt_agent_health_check",
			body, nil); err != nil {
			log.Errorf("sidecar %v OutputJson Event_ProcessResourceUsages error: %v", sid, err)
		}
	default:
		rpcServerErrorTotal.Inc()
		log.Errorf("unsupport event type: %v", proto.Event_EventType_name[int32(event.EventType)])
		return nil, ErrUnsupportEventType
	}

	return &proto.EmptyResponse{}, nil
}

func NewRpcService(apiHost string, apiPort int) *rpcService {
	rs := &rpcService{
		sidecarMap: make(map[uuid.UUID]sidecar),
		apiHost:    apiHost,
		apiPort:    apiPort,
	}

	SidecarClient = rs

	return rs
}

func (rpc *rpcService) IsClientExist(sid uuid.UUID) bool {
	rpc.RLock()
	_, ok := rpc.sidecarMap[sid]
	rpc.RUnlock()

	return ok
}
