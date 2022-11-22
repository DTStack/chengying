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

package client

import (
	"errors"
	"testing"

	"easyagent/internal/proto"
	"easyagent/internal/sidecar/base"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var defaultFakeEasyAgentServiceClient *fakeEasyAgentServiceClient
var defaultFakeEasyAgentServiceReadyForControlClient *fakeEasyAgentServiceReadyForControlClient

func init() {
	base.ConfigureLogger("/tmp", 0, 0, 0)
}

func defaultFakeReset() {
	defaultFakeEasyAgentServiceClient = &fakeEasyAgentServiceClient{}
	defaultFakeEasyAgentServiceReadyForControlClient = &fakeEasyAgentServiceReadyForControlClient{}
}

type fakeEasyAgentServiceClient struct {
	errRegisterSidecar error
	errReadyForControl error
	errReportEvent     error

	errReadyForControlCount int
}

func (c *fakeEasyAgentServiceClient) RegisterSidecar(ctx context.Context, in *proto.RegisterRequest, opts ...grpc.CallOption) (*proto.RegisterResponse, error) {
	out := new(proto.RegisterResponse)
	return out, c.errRegisterSidecar
}

func (c *fakeEasyAgentServiceClient) ReadyForControl(ctx context.Context, in *proto.ControlRequest, opts ...grpc.CallOption) (proto.EasyAgentService_ReadyForControlClient, error) {
	err := c.errReadyForControl
	if defaultFakeEasyAgentServiceClient.errReadyForControlCount <= 0 {
		err = nil
	}
	defaultFakeEasyAgentServiceClient.errReadyForControlCount--
	return defaultFakeEasyAgentServiceReadyForControlClient, err
}

func (c *fakeEasyAgentServiceClient) ReportEvent(ctx context.Context, in *proto.Event, opts ...grpc.CallOption) (*proto.EmptyResponse, error) {
	out := new(proto.EmptyResponse)
	return out, c.errReportEvent
}

type fakeEasyAgentServiceReadyForControlClient struct {
	err      error
	errCount int
	grpc.ClientStream
}

func (x *fakeEasyAgentServiceReadyForControlClient) Recv() (*proto.ControlResponse, error) {
	m := new(proto.ControlResponse)
	err := x.err
	if defaultFakeEasyAgentServiceReadyForControlClient.errCount <= 0 {
		err = nil
	}
	defaultFakeEasyAgentServiceReadyForControlClient.errCount--
	return m, err
}

func (x *fakeEasyAgentServiceReadyForControlClient) CloseSend() error {
	return nil
}

func NewFakeEasyAgentClient() (EaClienter, error) {
	eaClient := &easyAgentClient{
		registerCh: make(chan struct{}),
		client:     defaultFakeEasyAgentServiceClient,
	}
	eaClient.registerOk.Store(false)
	return eaClient, nil
}

func TestRegisterSidecarError(t *testing.T) {
	defaultFakeReset()
	defaultFakeEasyAgentServiceClient.errRegisterSidecar = errors.New("fake test always return error")
	c, err := NewFakeEasyAgentClient()
	if err != nil {
		t.Fatal(err)
	}

	err = c.RegisterSidecar(&proto.RegisterRequest{})
	if err == nil {
		t.Fatal("expect RegisterSidecar error, but its nil")
	}

	if err = c.ReportEvent(&proto.Event{}); err != NotRegister {
		t.Fatalf("expect ReportEvent return NotRegister, but its %v", err)
	}

	registerCh := c.(*easyAgentClient).registerCh
	select {
	case <-registerCh:
		t.Fatal("expect pending read from registerCh, but its not")
	default:
		t.Log("ok, registerCh is pending")
	}
	close(registerCh)

	// shoud not pending
	c.GetControlResponse()
}

func TestRegisterSidecarSuccess(t *testing.T) {
	defaultFakeReset()
	c, err := NewFakeEasyAgentClient()
	if err != nil {
		t.Fatal(err)
	}

	err = c.RegisterSidecar(&proto.RegisterRequest{})
	if err != nil {
		t.Fatalf("RegisterSidecar error: %v", err)
	}

	if err = c.ReportEvent(&proto.Event{}); err != nil {
		t.Fatalf("ReportEvent error: %v", err)
	}

	_, ok := <-c.(*easyAgentClient).registerCh
	if ok {
		t.Fatal("expect registerChan is closed, but its not")
	}

	// shoud not pending
	c.GetControlResponse()
}

func TestGetControlResponse(t *testing.T) {
	defaultFakeReset()
	defaultFakeEasyAgentServiceClient.errReadyForControl = errors.New("fake test ReadyForControl return error")
	defaultFakeEasyAgentServiceClient.errReadyForControlCount = 1
	defaultFakeEasyAgentServiceReadyForControlClient.err = errors.New("fake test Recv return error")
	defaultFakeEasyAgentServiceReadyForControlClient.errCount = 1
	c, err := NewFakeEasyAgentClient()
	if err != nil {
		t.Fatal(err)
	}

	err = c.RegisterSidecar(&proto.RegisterRequest{})
	if err != nil {
		t.Fatalf("RegisterSidecar error: %v", err)
	}

	c.GetControlResponse()
	if defaultFakeEasyAgentServiceClient.errReadyForControlCount != -2 {
		t.Fatalf("err to ReadyForControl should loop, but its not")
	}
	if defaultFakeEasyAgentServiceReadyForControlClient.errCount != -1 {
		t.Fatalf("err to Recv should loop, but its not")
	}
}
