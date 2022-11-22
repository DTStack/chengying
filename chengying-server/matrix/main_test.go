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
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFiles(t *testing.T) {
	t.SkipNow()
	dir := filepath.Join("/tmp", "test")

	files, err := util.ListFiles(dir)

	if err != nil {
		fmt.Println(err.Error())
	}
	var sub []string
	for _, f := range files {
		sub = append(sub, strings.Replace(f, dir, "", -1))
	}
	fmt.Println(strings.Join(sub, "\n"))
}

func TestFiles(t *testing.T) {
	t.SkipNow()
	var te []int

	if te == nil {
		fmt.Println("sp is sb")
		fmt.Println(te)
	}

	te = append(te, 1, 2, 3)

	fmt.Println(te)
}

func TestMainfile(t *testing.T) {
	//t.SkipNow()
	fmt.Println(base.VERSION)
	fmt.Println("Copyright (c) 2017 DTStack Inc.")
	base.ConfigureProductVersion(base.VERSION)

	app := cli.NewApp()
	app.Version = base.VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config,c",
			Usage: "config path",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug info",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		log.SetDebug(ctx.Bool("debug"))
		if err := ParseConfig(ctx.String("config")); err != nil {
			return err
		}
		return base.Run()
	}

	args := []string{"./matrix", "--config", "example-config.yml", "--debug"}

	err := app.Run(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit with failure: %v\n", err)
	}
}

func Test12(t *testing.T) {
	t.SkipNow()
	str := ` setSchemaFieldServiceAddr
sds    
`
	t.Log(str)
	t.Log(strings.Split(strings.TrimSpace(str), "\n"))

}

//func TestExec(t *testing.T) {
//
//	fmt.Println(base.VERSION)
//	fmt.Println("Copyright (c) 2017 DTStack Inc.")
//	base.ConfigureProductVersion(base.VERSION)
//
//	cfgPath := "/Users/duanjiaxing/GolandProjects/easymatrix/matrix/example-config.yml"
//	configContent, err := yaml.NewConfigWithFile(cfgPath, ucfg.PathSep("."))
//
//	if err != nil {
//		panic(err)
//	}
//
//	config := Config{}
//	if err := configContent.Unpack(&config); err != nil {
//		panic(err)
//	}
//
//	clog := &config.Log
//	if err := log.ConfigureLogger(clog.Dir, clog.MaxSize, clog.MaxBackups, clog.MaxAge); err != nil {
//		panic(err)
//	} else {
//		fmt.Printf("Saving logs at %s\n", clog.Dir)
//	}
//	//the other componet use the common package log
//	commonlog.Config("EM")
//
//	db := &config.MysqlDb
//	if err := model.ConfigureMysqlDatabase(db.Host, db.Port, db.User, db.Password, db.DbName); err != nil {
//		panic(err)
//
//	}
//	if err := modelkube.Build(); err != nil {
//		panic(err)
//
//	}
//	if err := union.Build(); err != nil {
//		panic(err)
//
//	}
//	if err := resource.InitResource(); err != nil {
//		panic(err)
//
//	}
//	if err := resetDeployStatus(); err != nil {
//		panic(err)
//
//	}
//
//	agentConfig := &config.Agent
//	apiconf := &config.Api
//	host.InitAgentInstall(agentConfig.Host, apiconf.StaticUrl, agentConfig.InstallPath)
//	agent.InitAgentClient(agentConfig.Host)
//	deployConfig := &config.Deploy
//	base.ConfigureDeployInstallPath(deployConfig.InstallPath)
//	log.SetDebug(true)
//
//	cmd := fmt.Sprintf("#!/bin/sh\nset -x\nfind %s -regex  \"^.*-.*-.*~$\"", base.INSTALL_CURRRENT_PATH)
//	agent.InitAgentClient("hwcld-03:8889")
//
//	execId := uuid.NewV4().String()
//	operationId := uuid.NewV4().String()
//
//	model.OperationList.Insert(model.OperationInfo{
//		OperationId:     operationId,
//		OperationType:   enums.OperationType.HostInit.Code,
//		OperationStatus: enums.ExecStatusType.Running.Code,
//		ObjectType:      enums.OperationObjType.Host.Code,
//		ObjectValue:     "51bc8884-dc76-47f4-946f-68e5609887b8",
//	})
//	err = model.ExecShellList.InsertExecShellInfo(operationId, execId, "", "", "51bc8884-dc76-47f4-946f-68e5609887b8", enums.ShellType.Exec.Code)
//	if err != nil {
//		panic(err)
//	}
//	content, err := agent.AgentClient.ToExecCmd("51bc8884-dc76-47f4-946f-68e5609887b8", "", strings.TrimSpace(cmd), execId)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(content)
//}
