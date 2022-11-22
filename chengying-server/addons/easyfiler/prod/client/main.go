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
	"fmt"
	"github.com/urfave/cli"
	"os"
)

const (
	VERSION = "1.0.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "Easyfiler"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "e.g: --config ./config.yml",
		},
	}
	app.Commands = []cli.Command{
		listCommand(),
		downloadCommand(),
		uploadCommand(),
		previewCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("exit with failure: %v\n", err)
		os.Exit(1)
	}

}
