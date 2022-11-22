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

package base

import (
	"net"
	"path/filepath"
	"strconv"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	slog "dtstack.com/dtstack/easymatrix/go-common/log"
	"github.com/kataras/iris"
	"github.com/natefinch/lumberjack"
)

var (
	VERSION       = "Easykube-0.0.1"
	_SYSTEM_FAIL  = make(chan SystemFailure)
	API_HOST      = "localhost"
	API_PORT      = 8899
	HTTP_PROTOCOL = "http://"
)

func ConfigureProductVersion(v string) {
	VERSION = v
}

func ConfigureApiServer(host string, port int, root *apibase.Route, restrictSchema bool, stop chan struct{}) error {
	API_HOST = host
	API_PORT = port
	app := iris.New()
	apibase.RegisterUUIDStringMacro(app)

	app.AttachLogger(&lumberjack.Logger{
		Filename:   filepath.Join(slog.LOGDIR, "kube-api.log"),
		MaxSize:    slog.LOGGER_MAX_SIZE,
		MaxBackups: slog.LOGGER_MAX_BKS,
		MaxAge:     slog.LOGGER_MAX_AGE,
	})

	if err := apibase.InitSchema(app, root, restrictSchema); err != nil {
		return err
	}

	go func() {
		err := app.Run(iris.Addr(net.JoinHostPort(host, strconv.Itoa(port))), iris.WithoutBodyConsumptionOnUnmarshal) //二次消费body
		if err != nil {
			close(stop)
			SystemExitWithFailure(NETWORK_FAILURE, "API server failure: %v", err)
		}
	}()
	return nil
}
