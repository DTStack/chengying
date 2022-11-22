// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package rpc_server

import (
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/proto"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	Port    string
	Root    string
	WithDB  bool
	Rate    int
	Isziped bool
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Errorf("failed to listen, err: %v", err)
		return err
	}
	srv := grpc.NewServer()
	proto.RegisterTransferServiceServer(srv, &FileTransferService{Root: s.Root, WithDB: s.WithDB, Rate: s.Rate, IsZiped: s.Isziped})
	reflection.Register(srv)
	if err := srv.Serve(listener); err != nil {
		log.Errorf("failed to serve, err: %v", err)
		return err
	}
	return nil
}
