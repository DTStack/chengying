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
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/model/mysql"
	serve "dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/rpc-server"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"os"
)

const LOG_PREFIX = "easyfiler"

type ServerConfig struct {
	Root    string `config:"root"`
	Port    string `config:"port"`
	WithDB  bool   `config:"withdb"`
	Rate    int    `config:"rate"`
	IsZiped bool   `config:"isziped"`
}

type DatabaseConfig struct {
	User     string `config:"user"`
	Password string `config:"password"`
	Host     string `config:"host"`
	Port     int    `config:"port"`
	DbName   string `config:"dbname"`
}

type LogConfig struct {
	Dir        string `config:"dir"`
	MaxSize    int    `config:"max-logger-size"`
	MaxBackups int    `config:"max-logger-backups"`
	MaxAge     int    `config:"days-to-keep"`
}

type Config struct {
	Mysql  DatabaseConfig `config:"mysql"`
	Log    LogConfig      `config:"log"`
	Server ServerConfig   `config:"server"`
}

func ParseConfig(filename string) (*serve.Server, error) {
	configContent, err := yaml.NewConfigWithFile(filename, ucfg.PathSep("."))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found!")
		}
		return nil, err
	}
	config := &Config{}
	if err := configContent.Unpack(config); err != nil {
		return nil, err
	}

	clog := config.Log
	if err := log.ConfigureLogger(LOG_PREFIX, clog.Dir, clog.MaxSize, clog.MaxBackups, clog.MaxAge); err != nil {
		return nil, err
	} else {
		fmt.Printf("Saving logs at %s\n", clog.Dir)
	}
	if config.Server.WithDB {
		db := config.Mysql
		if err := mysql.InitDB(mysql.DBconf{
			User:     db.User,
			Password: db.Password,
			Host:     db.Host,
			Port:     db.Port,
			DB:       db.DbName,
		}); err != nil {
			return nil, err
		}
	}

	server := config.Server
	srv := &serve.Server{
		Port:    server.Port,
		Root:    server.Root,
		WithDB:  server.WithDB,
		Rate:    server.Rate,
		Isziped: server.IsZiped,
	}
	return srv, nil
}
