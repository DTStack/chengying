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
	"easyagent/internal/proto"
	"golang.org/x/net/context"
)

type fakeRpcServer struct{}

func (rpc *fakeRpcServer) RegisterSidecar(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return &proto.RegisterResponse{}, nil
}

func (rpc *fakeRpcServer) ReadyForControl(request *proto.ControlRequest, stream proto.EasyAgentService_ReadyForControlServer) error {
	return nil
}

func (rpc *fakeRpcServer) ReportEvent(ctx context.Context, event *proto.Event) (*proto.EmptyResponse, error) {
	return &proto.EmptyResponse{}, nil
}
