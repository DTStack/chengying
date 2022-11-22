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
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"os"
)

const LOG_PREFIX = "easyrke"

type LogConfig struct {
	Dir        string `config:"dir" validate:"required"`
	MaxSize    int    `config:"max-logger-size"`
	MaxBackups int    `config:"max-logger-backups"`
	MaxAge     int    `config:"days-to-keep"`
}

type Config struct {
	Log LogConfig `config:"log" validate:"required"`
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
	clog := config.Log
	if err := log.ConfigureLogger(LOG_PREFIX, clog.Dir, clog.MaxSize, clog.MaxBackups, clog.MaxAge); err != nil {
		return err
	} else {
		fmt.Printf("Saving logs at %s\n", clog.Dir)
	}
	return nil
}
