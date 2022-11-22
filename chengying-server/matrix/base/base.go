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
	"dtstack.com/dtstack/easymatrix/matrix/swagger"
	"net"
	"path/filepath"
	"strconv"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	slog "dtstack.com/dtstack/easymatrix/matrix/log"
	"github.com/kataras/iris"
	"github.com/natefinch/lumberjack"
)

var (
	VERSION               = "EasyMatrix-2.1.3"
	_SYSTEM_FAIL          = make(chan SystemFailure)
	API_HOST              = "localhost"
	API_PORT              = 8864
	API_STATIC_URL        = "http://localhost:8864"
	HTTP_PROTOCOL         = "http://"
	INSTALL_CURRRENT_PATH = "/opt/dtstack/"
)

func ConfigureProductVersion(v string) {
	VERSION = v
}

func ConfigureDeployInstallPath(path string) {
	INSTALL_CURRRENT_PATH = path
}

func ConfigureApiServer(host string, port int, staticUrl string, root *apibase.Route, restrictSchema, swaggerSwitch bool) error {
	API_HOST = host
	API_PORT = port
	API_STATIC_URL = staticUrl
	app := iris.New()
	apibase.RegisterUUIDStringMacro(app)

	app.AttachLogger(&lumberjack.Logger{
		Filename:   filepath.Join(slog.LOGDIR, "matrix-api.log"),
		MaxSize:    slog.LOGGER_MAX_SIZE,
		MaxBackups: slog.LOGGER_MAX_BKS,
		MaxAge:     slog.LOGGER_MAX_AGE,
	})

	// 创建静态web url
	app.StaticServe(WebRoot, "/easyagent")

	if err := apibase.InitSchema(app, root, restrictSchema); err != nil {
		return err
	}

	// 判断是否启动swagger
	if swaggerSwitch {
		swagger.InitializeSwagger(app)
	}

	go func() {
		err := app.Run(iris.Addr(net.JoinHostPort(host, strconv.Itoa(port))), iris.WithoutBodyConsumptionOnUnmarshal) //二次消费body
		if err != nil {
			SystemExitWithFailure(NETWORK_FAILURE, "API server failure: %v", err)
		}
	}()
	return nil
}
