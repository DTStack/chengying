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

package main

import (
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/api"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/base"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/monitor"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"net"
	"os"
	"strconv"
)

const LOG_PREFIX = "easykube"

var ENV_NAMESPACE = "WATCH_NAMESPACE"
var ENV_MONITOR_ENABLE = "MONITOR_ENABLE"
var ENV_OPERATOR_INIT = "OPERATOR_INIT"

type LogConfig struct {
	Dir        string `config:"dir" validate:"required"`
	MaxSize    int    `config:"max-logger-size"`
	MaxBackups int    `config:"max-logger-backups"`
	MaxAge     int    `config:"days-to-keep"`
}

type ApiConfig struct {
	Host     string `config:"host"`
	Port     int    `config:"port" validate:"required"`
	Restrict bool   `config:"restrict-api-check"`
}

type MonitorConfig struct {
	Enabled bool `config:"enabled"`
}

type Config struct {
	Log     LogConfig     `config:"log" validate:"required"`
	Api     ApiConfig     `config:"api" validate:"required"`
	Monitor MonitorConfig `config:"monitor" validate:"required"`
}

func ParseConfig(configFile string) error {
	configContent, err := yaml.NewConfigWithFile(configFile, ucfg.PathSep("."))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found!")
		}
		return err
	}

	config := Config{}
	if err := configContent.Unpack(&config); err != nil {
		return err
	}
	//clog := config.Log
	log.Config(LOG_PREFIX)
	//if err := log.ConfigureLogger(LOG_PREFIX, clog.Dir, clog.MaxSize, clog.MaxBackups, clog.MaxAge); err != nil {
	//	return err
	//} else {
	//	fmt.Printf("Saving logs at %s\n", clog.Dir)
	//}
	stop := make(chan struct{})
	enalbed := config.Monitor.Enabled
	if e := os.Getenv(ENV_MONITOR_ENABLE); len(e) != 0 {
		enalbed, err = strconv.ParseBool(e)
		if err != nil {
			return fmt.Errorf("env %s is not bool type", ENV_MONITOR_ENABLE)
		}
	}
	if enalbed {
		namespace := os.Getenv("WATCH_NAMESPACE")
		if len(namespace) == 0 {
			return fmt.Errorf("config namespace error, namespace is nil!")
		}
		err := monitor.StartMonitor(namespace, stop)
		if err != nil {
			return err
		}
	}

	apiconf := &config.Api

	if err := base.ConfigureApiServer(apiconf.Host, apiconf.Port, &api.ApiV2Schema, apiconf.Restrict, stop); err != nil {
		return err
	} else {
		fmt.Printf("Running API service at %v\n", net.JoinHostPort(apiconf.Host, strconv.Itoa(apiconf.Port)))
	}

	return nil
}
