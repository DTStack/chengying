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

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/satori/go.uuid"
	y2 "gopkg.in/yaml.v2"
)

type logConfig struct {
	Dir        string `config:"dir" validate:"required"`
	MaxSize    int    `config:"max-logger-size"`
	MaxBackups int    `config:"max-logger-backups"`
	MaxAge     int    `config:"days-to-keep"`
}

type rpcConfig struct {
	Server        string `config:"server" validate:"required"`
	Port          int    `config:"port" validate:"required"`
	CertFile      string `config:"cert"`
	Tls           bool   `config:"tls"`
	TlsSkipVerify bool   `config:"tls-skip-verify"`
}

type easyAgentConfig struct {
	Uuid            WrapperUUID   `config:"uuid" validate:"required"`
	Network         []string      `config:"network"`
	MonitorInterval time.Duration `config:"monitor-interval" validate:"required,min=1s"`
}

type Config struct {
	Log       logConfig       `config:"log" validate:"required"`
	Rpc       rpcConfig       `config:"rpc" validate:"required"`
	EasyAgent easyAgentConfig `config:"easyagent" validate:"required"`
	CallBack  []string        `config:"callback"`
}

type AgentConfig struct {
	BinaryPath        string            `yaml:"binary-path"`
	AgentId           WrapperUUID       `yaml:"agentId"`
	ConfigurationPath string            `yaml:"configuration-path"`
	Parameter         []string          `yaml:"parameter,flow"`
	Workdir           string            `yaml:"workdir,omitempty"`
	Name              string            `yaml:"name"`
	HealthShell       string            `yaml:"health-shell,omitempty"`
	HealthPeriod      time.Duration     `yaml:"health-period,omitempty"`
	HealthStartPeriod time.Duration     `yaml:"health-start-period,omitempty"`
	HealthTimeout     time.Duration     `yaml:"health-timeout,omitempty"`
	HealthRetries     uint64            `yaml:"health-retries,omitempty"`
	Enabled           bool              `yaml:"enabled"`
	CpuLimit          float32           `yaml:"cpu-limit,omitempty"`
	MemLimit          uint64            `yaml:"mem-limit,omitempty"`
	NetLimit          uint64            `yaml:"net-limit,omitempty"`
	Environment       map[string]string `yaml:"environment,omitempty"`
	RunUser           string            `yaml:"run-user,omitempty"`
}

func ParseConfig(configFile string) (*Config, error) {
	configContent, err := yaml.NewConfigWithFile(configFile, ucfg.PathSep("."))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found!")
		}
		return nil, err
	}

	cfg := Config{}
	if err = configContent.Unpack(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseAgents(configFile string) ([]AgentConfig, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var agents []AgentConfig
	err = y2.Unmarshal(content, &agents)
	if err != nil {
		return nil, err
	}

	return agents, checkAgents(agents)
}

func WriteAgents(agents []AgentConfig, configFile string) error {
	content, err := y2.Marshal(agents)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, content, 0600)
}

func checkAgents(agents []AgentConfig) error {
	m := make(map[WrapperUUID]struct{}, len(agents))
	for _, agent := range agents {
		if agent.AgentId.UUID == uuid.Nil {
			return fmt.Errorf("an agent agentId is empty")
		}
		if _, ok := m[agent.AgentId]; ok {
			return fmt.Errorf("agent agentId deplicate: %s", agent.AgentId)
		}
		m[agent.AgentId] = struct{}{}

		if agent.HealthShell != "" && agent.HealthPeriod < time.Second {
			return fmt.Errorf("agent %v health-period less than 1 sec", agent.AgentId)
		}
	}

	return nil
}
