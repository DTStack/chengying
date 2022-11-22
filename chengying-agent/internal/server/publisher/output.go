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

package publisher

import (
	"context"
	"time"

	"github.com/elastic/go-ucfg"
)

type esConfig struct {
	Urls     []string `config:"hosts" validate:"required"`
	UserName string   `config:"username"`
	PassWord string   `config:"password"`
}

type kafkaConfig struct {
	Urls     []string      `config:"hosts" validate:"required"`
	UserName string        `config:"username"`
	PassWord string        `config:"password"`
	Timeout  time.Duration `config:"timeout"             validate:"min=1"`
}

type HttpConfig struct {
	ApiHost string `config:"host" validate:"required"`
	ApiUri  string `config:"uri" validate:"required"`
}

type influxdbConfig struct {
	Urls     []string `config:"hosts" validate:"required"`
	UserName string   `config:"username"`
	PassWord string   `config:"password"`
}

type fileConfig struct {
	Path string `config:"path" validate:"required"`
}

type TransferConfig struct {
	Concurrency   uint8         `config:"concurrency" validate:"min=1,max=32"`
	Timeout       time.Duration `config:"timeout" validate:"min=1s"`
	Server        string        `config:"server" validate:"required"`
	Port          int           `config:"port" validate:"required"`
	CertFile      string        `config:"cert"`
	Tls           bool          `config:"tls"`
	TlsSkipVerify bool          `config:"tls-skip-verify"`
}

type OutputConfig struct {
	EsConfig       esConfig       `config:"elasticsearch"`
	InfluxdbConfig influxdbConfig `config:"influxdb"`
	FileConfig     fileConfig     `config:"file"`
}

type Outputer interface {
	Name() string
	OutputJson(ctx context.Context, id, index, tpy string, js interface{}, key []byte) error
	Close()
}

type OutputCreater func(config map[string]*ucfg.Config) (Outputer, error)
