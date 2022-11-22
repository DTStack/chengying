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
	"testing"

	"gopkg.in/yaml.v2"
)

func TestParseConfig(t *testing.T) {
	cfg, err := ParseConfig(`C:\Users\guyan\go\src\dtstack.com\dtstack\easyagent\prod\sidecar\example-config.yml`)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("cfg: %#v", cfg)

	b, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("[%s]", b)
}

func TestParseAgents(t *testing.T) {
	agents, err := ParseAgents(`C:\Users\guyan\go\src\dtstack.com\dtstack\easyagent\prod\sidecar\agents-file.yml`)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cfg: %#v", agents)
}

func TestWriteAgents(t *testing.T) {
	agents, err := ParseAgents(`C:\Users\guyan\go\src\dtstack.com\dtstack\easyagent\prod\sidecar\agents-file.yml`)
	if err != nil {
		t.Fatal(err)
	}
	err = WriteAgents(agents, `C:\Users\guyan\go\src\dtstack.com\dtstack\easyagent\prod\sidecar\agents-file.yml`)
	if err != nil {
		t.Fatal(err)
	}
}
