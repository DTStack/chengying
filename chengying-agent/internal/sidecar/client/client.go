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
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"sync/atomic"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/sidecar/base"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	NotRegister = errors.New("not register sidecar")
)

type EaClienter interface {
	// RegisterSidecar must be first grpc call, otherwise call will pending
	RegisterSidecar(req *proto.RegisterRequest) error
	GetControlResponse() *proto.ControlResponse
	ReportEvent(event *proto.Event) error
	GetServerAddress() string
}

type easyAgentClient struct {
	registerOk atomic.Value
	registerCh chan struct{}

	uuid          proto.Uuid
	serverAddress string
	client        proto.EasyAgentServiceClient
	stream        proto.EasyAgentService_ReadyForControlClient
}

func (c *easyAgentClient) RegisterSidecar(req *proto.RegisterRequest) error {
	req.Id = c.uuid
	req.Systime = time.Now()
	ctx, cancal := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancal()

	_, err := c.client.RegisterSidecar(ctx, req)
	if err == nil {
		c.registerOk.Store(true)
		close(c.registerCh)
	}
	return err
}

func (c *easyAgentClient) GetControlResponse() *proto.ControlResponse {
	for {
		if c.stream == nil {
			req := &proto.ControlRequest{}
			req.Id = c.uuid
			req.Systime = time.Now()
			<-c.registerCh // waiting for RegisterSidecar

			stream, err := c.client.ReadyForControl(context.Background(), req)
			if err != nil {
				base.Errorf("GetControlResponse error: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}
			c.stream = stream
		}

		ctlResp, err := c.stream.Recv()
		if err != nil {
			base.Errorf("GetControlResponse Recv error: %v", err)
			c.stream.CloseSend()
			c.stream = nil // reset stream
			time.Sleep(3 * time.Second)
			continue
		}
		return ctlResp
	}
}

func (c *easyAgentClient) ReportEvent(event *proto.Event) error {
	if !c.registerOk.Load().(bool) {
		return NotRegister
	}

	event.Id = c.uuid
	event.Systime = time.Now()

	go func() {
		ctx, cancal := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancal()

		if _, err := c.client.ReportEvent(ctx, event); err != nil {
			base.Errorf("(%v) ReportEvent error: %v", event, err)
		}
	}()

	return nil
}

func (c *easyAgentClient) GetServerAddress() string {
	return c.serverAddress
}

// NewEasyAgentClient not block, it connect server when call grpc function
func NewEasyAgentClient(uuid uuid.UUID, serverAddr string, istls, tlsSkipVerify bool, certFile string) (EaClienter, error) {
	var opts []grpc.DialOption
	if istls {
		var cp *x509.CertPool
		if certFile != "" {
			b, err := ioutil.ReadFile(certFile)
			if err != nil {
				return nil, err
			}
			cp = x509.NewCertPool()
			if !cp.AppendCertsFromPEM(b) {
				return nil, errors.New("credentials: failed to append certificates")
			}
		}
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: tlsSkipVerify, RootCAs: cp})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	// dial is not block, it connect server when call grpc function
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	eaClient := &easyAgentClient{
		registerCh:    make(chan struct{}),
		uuid:          proto.Uuid(uuid.Bytes()),
		serverAddress: serverAddr,
		client:        proto.NewEasyAgentServiceClient(conn),
	}
	eaClient.registerOk.Store(false)
	return eaClient, nil
}
