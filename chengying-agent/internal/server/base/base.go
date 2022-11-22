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

package base

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	apibase "easyagent/go-common/api-base"
	"easyagent/internal/proto"
	slog "easyagent/internal/server/log"
	"easyagent/internal/server/rpc"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/kataras/iris"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	VERSION         = "dev-snapshot"
	RPCSERVICE      proto.EasyAgentServiceServer
	_RPCSERVER      *grpc.Server
	_SYSTEM_FAIL    = make(chan SystemFailure)
	RPC_SERVER_PORT int
	RPC_USE_TLS     bool
	API_HOST        string = "localhost"
	API_PORT        int    = 8889
)

func ConfigureProductVersion(v string) {
	VERSION = v
}

func ConfigureRpcService(port int, certFile, keyFile string, apiHost string, apiPort int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	RPC_USE_TLS = certFile != ""

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	}
	if RPC_USE_TLS {
		cert, err := tls.X509KeyPair([]byte(certFile), []byte(keyFile))
		if err != nil {
			return err
		}
		creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		opts = append(opts, grpc.Creds(creds))
	}

	if apiHost == "" {
		ifs, err := net.Interfaces()
		if err != nil {
			return err
		}
	IFS:
		for _, v := range ifs {
			if (v.Flags&net.FlagUp) > 0 &&
				(v.Flags&net.FlagLoopback) == 0 &&
				!strings.Contains(v.Name, "docker") &&
				!strings.Contains(v.Name, "veth") &&
				!strings.Contains(v.Name, "flannel") &&
				!strings.Contains(v.Name, "cni") &&
				!strings.Contains(v.Name, ":") {
				addrs, err := v.Addrs()
				if err != nil {
					continue
				}
				for _, addr := range addrs {
					ip := addr.(*net.IPNet).IP
					if ip.IsGlobalUnicast() && ip.To4() != nil {
						apiHost = ip.String()
						slog.Infof("found api host %v(%v) for rpc server", apiHost, v.Name)
						break IFS
					}
				}
			}
		}
	}
	if apiHost == "" {
		return errors.New("not found api host interface ip")
	}

	_RPCSERVER = grpc.NewServer(opts...)
	RPCSERVICE = rpc.NewRpcService(apiHost, apiPort)
	proto.RegisterEasyAgentServiceServer(_RPCSERVER, RPCSERVICE)
	grpc_prometheus.Register(_RPCSERVER)

	go func() {
		if err := _RPCSERVER.Serve(lis); err != nil {
			SystemExitWithFailure(NETWORK_FAILURE, "RPC server failure: %v", err)
		}
	}()
	RPC_SERVER_PORT = port
	return nil
}

func ConfigureApiServer(host string, port int, root *apibase.Route, restrictSchema bool) error {
	API_HOST = host
	API_PORT = port
	app := iris.New()
	apibase.RegisterUUIDStringMacro(app)
	app.AttachLogger(&lumberjack.Logger{
		Filename:   filepath.Join(slog.LOGDIR, "api.log"),
		MaxSize:    slog.LOGGER_MAX_SIZE,
		MaxBackups: slog.LOGGER_MAX_BKS,
		MaxAge:     slog.LOGGER_MAX_AGE,
	})
	app.Get("metrics", iris.ToHandler(promhttp.Handler())) // must add first
	app.StaticServe("easyagent", "/easyagent")
	if err := apibase.InitSchema(app, root, restrictSchema); err != nil {
		return err
	}

	go func() {
		err := app.Run(iris.Addr(net.JoinHostPort(host, strconv.Itoa(port))))
		if err != nil {
			SystemExitWithFailure(NETWORK_FAILURE, "API server failure: %v", err)
		}
	}()
	return nil
}
