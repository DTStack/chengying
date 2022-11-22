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
	hd "dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/handler"
	"fmt"
	"github.com/urfave/cli"
)

func listCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "get file tree under the easyfiler root",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "target, t",
				Usage: "e.g: --target 127.0.0.1:7788",
			},
		},
		Action: func(ctx *cli.Context) error {
			target := ctx.String("target")
			_ = ctx.String("path")
			lists, err := hd.List(target, "txt")
			if err == nil {
				for i := range lists {
					fmt.Println(lists[i])
				}
			}
			return err
		},
	}
}

func downloadCommand() cli.Command {
	return cli.Command{
		Name:  "download",
		Usage: "download file from easyfiler server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "target, t",
				Usage: "e.g: --target 127.0.0.1:7788",
			},
			cli.StringFlag{
				Name:  "path, p",
				Usage: "e.g: --path easyfiler/easyfiler.log",
			},
		},
		Action: func(ctx *cli.Context) error {
			target := ctx.String("target")
			path := ctx.String("path")
			return hd.Download(target, path)
		},
	}
}

func uploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "upload file to easyfiler server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "target, t",
				Usage: "e.g: --target 127.0.0.1:7788",
			},
			cli.StringFlag{
				Name:  "file, f",
				Usage: "e.g: --file ./config.yml",
			},
		},
		Action: func(ctx *cli.Context) error {
			target := ctx.String("target")
			file := ctx.String("file")
			return hd.Upload(target, file)
		},
	}
}

func previewCommand() cli.Command {
	return cli.Command{
		Name:  "preview",
		Usage: "preview txt file on easyfiler server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "target, t",
				Usage: "e.g: --target 127.0.0.1",
			},
			cli.StringFlag{
				Name:  "path,p",
				Usage: "e.g --path easyfiler/easyfiler.log",
			},
		},
		Action: func(ctx *cli.Context) error {
			target := ctx.String("target")
			path := ctx.String("path")

			lines, err := hd.Preview(target, path, "latest")
			if err == nil {
				for i := range lines {
					fmt.Println(lines[i])
				}
			}
			return err
		},
	}
}
