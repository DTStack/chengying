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

package schema

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"text/template"
)

func TestParseSchemaConfig(t *testing.T) {
	cfg, err := ParseSchemaConfigFile(`schema_test.yml`)
	if err != nil {
		t.Fatal(err)
	}

	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if cfg, err = Unmarshal(jsonCfg); err != nil {
		t.Fatal(err)
	}

	if err := cfg.SetServiceAddr("es", []string{"127.0.0.1", "127.0.0.2"}, []string{"local1", "local2"}); err != nil {
		t.Fatal(err)
	}
	if err := cfg.SetServiceAddr("dtlog", []string{"172.0.0.8", "172.0.0.9"}, []string{"local8", "local9"}); err != nil {
		t.Fatal(err)
	}
	newSchema, err := Clone(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(newSchema, cfg) {
		t.Fatal("not DeepEqual")
	}
	t.Logf("%p --- %p", newSchema.Service["es"].ServiceAddr, cfg.Service["es"].ServiceAddr)
	newSchema.Service["es"].ServiceAddr.IP[0] = "my god"
	t.Logf("%#v --- %#v", newSchema.Service["es"].ServiceAddr, cfg.Service["es"].ServiceAddr)
	if cfg.Service["es"].ServiceAddr.IP[0] == "my god" {
		t.Fatal("newSchema and cfg affect etch other")
	}
	//if err := cfg.SetBaseService("dtuic", ConfigMap{"jjww": VisualConfig{Type: "internal", Desc: "internal", Default: "abc", Value: "abc"}}); err != nil {
	//	t.Fatal(err)
	//}
	if err = cfg.ParseVariable(); err != nil {
		t.Fatal(err)
	}
	jsonCfg, err = json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("cfg: %s", jsonCfg)
	//if cfg, err = Unmarshal(jsonCfg); err != nil {
	//	t.Fatal(err)
	//}

	cfgM, err := cfg.ParseConfigFiles("")
	if err != nil {
		t.Fatal(err)
	}
	for svcName, cfg := range cfgM {
		t.Logf("%v config file: %s", svcName, cfg)
	}

	if err := cfg.SetServiceNodeIP("dtlog", 1, 1, nil); err != nil {
		t.Fatal(err)
	}
	if err := cfg.ParseServiceVariable("dtlog"); err != nil {
		t.Fatal(err)
	}
	jsonCfg, err = json.MarshalIndent(cfg.Service["dtlog"], "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dtlog node1 cfg: %s", jsonCfg)
	jsonCfgs, err := cfg.ParseServiceConfigFiles("", "dtlog")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dtlog node1 config file: %s", jsonCfgs[0])
}

func TestSchemaConfig_SetField(t *testing.T) {
	cfg, err := ParseSchemaConfigFile(`schema_test.yml`)
	if err != nil {
		t.Fatal(err)
	}

	old, err := cfg.SetField("", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField(".", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.not_found", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.Instance.Cmd.some", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.Instance.ConfigPaths.some", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.Config", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.Config.log_port.Desc", "test")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.Instance.Logs.10", "test")
	if err == nil {
		t.Fatal(err)
	}
	old, err = cfg.SetField("dtlog.Config.log_port.Value", "test")
	if err != nil {
		t.Fatal(err)
	}
	old, err = cfg.SetField("dtlog.Instance.Environment.MYSQL_ADDRESS", "test")
	if err != nil {
		t.Fatal(err)
	}
	old, err = cfg.SetField("dtlog.Instance.Logs.0", "test")
	if err != nil {
		t.Fatal(err)
	}
	old, err = cfg.SetField("dtlog.Instance.UseCloud", "abc")
	if err == nil {
		t.Fatal("expect error, but its not")
	}
	old, err = cfg.SetField("dtlog.Instance.UseCloud", "true")
	if err != nil {
		t.Fatal(err)
	}
	jsonCfg, err := json.MarshalIndent(cfg.Service["dtlog"], "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("old field: %v, %s", old, jsonCfg)
}

func TestVisualConfig_String(t *testing.T) {
	cfg := `arg: -c {{.ips}}    {{.LastSegIP .ips}}      {{range $i, $v := .ips.Value.IP}}{{if $i}},{{end}}http://{{$v}}:{{$.port}}|{{$.AddOne $i}}{{end}}
http://{{.Join .ips ":" .port ","}}:{{.port}}    http://{{.Joinx .ips "," ":" .port}}     {{.NodeCount .ips}}     {{if eq (.NodeID .ips) 2}}NodeID=2{{end}}
JoinHost: http://{{.JoinHost "ips" ":" .port ","}}:{{.port}}
JoinxHost: http://{{.JoinxHost "ips" "," ":" .port}}
Hostname: {{.Hostname .ips}}
NodeIndex: {{.NodeIndex "ips"}}
{{range $v := .IPList .ips}}{{$v}},{{end}}
{{range $v := .HostList "ips"}}{{$v}},{{end}}
GetIpByNodeID: {{.GetIpByNodeID .ips (.NodeID .ips)}}
`
	port := "8080"
	config := ConfigMap{
		"ips": VisualConfig{
			Value: &ServiceAddrStruct{
				map[uint]uint{1: 0, 2: 1, 3: 2},
				[]string{"local1", "local2", "local3"},
				[]string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}, 1, 2},
		},
		"port": VisualConfig{Value: &port},
	}
	tpl := template.Must(template.New("test").Option("missingkey=error").Parse(cfg))
	err := tpl.Execute(os.Stdout, config)
	if err != nil {
		t.Fatal(err)
	}
}
