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

package main

import (
	"fmt"
	"os"

	"easyagent/internal/server/api"
	"easyagent/internal/server/base"
	"easyagent/internal/server/log"
	"easyagent/internal/server/model"
	"easyagent/internal/server/publisher"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
)

type DatabaseConfig struct {
	User     string `config:"user" validate:"required"`
	Password string `config:"password" validate:"required"`
	Host     string `config:"host" validate:"required"`
	Port     int    `config:"port"`
	DbName   string `config:"dbname" validate:"required"`
}

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

type RpcConfig struct {
	Port     int    `config:"port" validate:"required"`
	CertFile string `config:"cert"`
	KeyFile  string `config:"key"`
}

type Config struct {
	MysqlDb DatabaseConfig          `config:"mysqldb" validate:"required"`
	Publish map[string]*ucfg.Config `config:"publish"`
	Log     LogConfig               `config:"log" validate:"required"`
	Api     ApiConfig               `config:"api" validate:"required"`
	Rpc     RpcConfig               `config:"rpc" validate:"required"`
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

	clog := &config.Log
	if err := log.ConfigureLogger(clog.Dir, clog.MaxSize, clog.MaxBackups, clog.MaxAge); err != nil {
		return err
	} else {
		fmt.Printf("Saving logs at %s\n", clog.Dir)
	}

	db := &config.MysqlDb
	if err := model.ConfigureMysqlDatabase(db.Host, db.Port, db.User, db.Password, db.DbName); err != nil {
		return err
	}

	publish := config.Publish
	if err := publisher.Publish.ConfigOutput(publish); err != nil {
		return err
	} else {
		fmt.Printf("config publish ok!\n")
	}

	rpc := &config.Rpc
	apiconf := &config.Api
	if err := base.ConfigureRpcService(rpc.Port, rpc.CertFile, rpc.KeyFile, apiconf.Host, apiconf.Port); err != nil {
		return err
	} else {
		fmt.Printf("Running RPC service at %d\n", rpc.Port)
	}

	if err := base.ConfigureApiServer(apiconf.Host, apiconf.Port, &api.ApiV1Schema, apiconf.Restrict); err != nil {
		return err
	} else {
		fmt.Printf("Running API service at %d\n", apiconf.Port)
	}

	return nil
}
