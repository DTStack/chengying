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
	"context"
	"errors"

	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/satori/go.uuid"
)

var (
	SidecarClient SidecarClienter

	sendControlTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "send_control_total",
		Help: "Total Number of SendControl(include SendControlSync)",
	})
	sendControlSyncTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "send_control_sync_total",
		Help: "Total Number of SendControlSync",
	})
	sendControlErrorTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "send_control_error_total",
		Help: "Total Number of SendControl errors",
	})

	ErrNotFound = errors.New("not found sid")
	ErrRpcStop  = errors.New("rpc stoped")
)

type SidecarClienter interface {
	SendControl(ctx context.Context, sid uuid.UUID, ctlResp *proto.ControlResponse) error
	SendControlSync(ctx context.Context, sid uuid.UUID, ctlResp *proto.ControlResponse) (op interface{}, err error)
	IsClientExist(sid uuid.UUID) bool
	Close(sid uuid.UUID) error
}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(sendControlTotal)
	prometheus.MustRegister(sendControlSyncTotal)
	prometheus.MustRegister(sendControlErrorTotal)
}

func (rpc *rpcService) SendControl(ctx context.Context, sid uuid.UUID, ctlResp *proto.ControlResponse) error {
	sendControlTotal.Inc()

	rpc.RLock()
	sc, ok := rpc.sidecarMap[sid]
	rpc.RUnlock()
	if !ok {
		sendControlErrorTotal.Inc()
		log.Errorf("sidecar %v not found", sid)
		return ErrNotFound
	}

	var err error

	select {
	case sc.ctlCh <- ctlResp:
	case <-sc.stopCh:
		err = ErrRpcStop
		sendControlErrorTotal.Inc()
	case <-ctx.Done():
		err = ctx.Err()
		sendControlErrorTotal.Inc()
	}

	return err
}

func (rpc *rpcService) SendControlSync(ctx context.Context, sid uuid.UUID, ctlResp *proto.ControlResponse) (op interface{}, err error) {
	sendControlSyncTotal.Inc()

	ch := waitSeqno(ctlResp.Seqno)
	err = rpc.SendControl(ctx, sid, ctlResp)
	if err != nil {
		stopSeqno(ctlResp.Seqno, nil)
		return
	}

	select {
	case op = <-ch:
	case <-ctx.Done():
		err = ctx.Err()
		sendControlErrorTotal.Inc()
	}

	return
}

func (rpc *rpcService) Close(sid uuid.UUID) error {
	rpc.RLock()
	sc, ok := rpc.sidecarMap[sid]
	rpc.RUnlock()
	if !ok {
		log.Errorf("sidecar %v not found", sid)
		return ErrNotFound
	}

	close(sc.stopCh)
	return nil
}
