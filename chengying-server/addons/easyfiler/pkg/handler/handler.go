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

package handler

import (
	"context"
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/proto"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"os"
	"strings"
	"time"
)

const (
	EASYFILER_TMP_ROOT = "/tmp/"
	LINUX_SYSTEM_SLASH = "/"
)

func Preview(target, path, action string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dail host")
	}
	defer conn.Close()
	c := proto.NewTransferServiceClient(conn)
	reqData := &proto.PreviewRequest{
		Path:   path,
		Action: action,
	}
	cli, err := c.Preview(context.Background(), reqData)
	if err != nil {
		return nil, err
	}
	var lines = make([]string, 1)
	for {
		res, err := cli.Recv()
		if err != nil {
			break
		}
		lines = append(lines, res.LogRow)
	}
	return lines, nil
}

func List(target, suffix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	c := proto.NewTransferServiceClient(conn)
	reqData := &proto.ListRequest{
		Suffix: suffix,
	}
	cli, err := c.List(context.Background(), reqData)
	if err != nil {
		return nil, err
	}
	var lists = make([]string, 1)
	for {
		res, err := cli.Recv()
		if err != nil {
			break
		}
		lists = append(lists, res.Name)
	}
	return lists, nil
}

func Download(target, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
		return err
	}
	defer conn.Close()

	c := proto.NewTransferServiceClient(conn)

	reqData := &proto.DownloadRequest{
		Path: path,
	}
	dir, filename := spiteDirAndFilename(path)
	cli, err := c.Download(context.Background(), reqData)
	if err != nil {
		return err
	}

	absoluteDir := EASYFILER_TMP_ROOT + target + LINUX_SYSTEM_SLASH + dir
	if !isPathExist(absoluteDir) {
		if err := os.MkdirAll(absoluteDir, 0666); err != nil {
			fmt.Printf("failed to create file, err:%v\n", err)
			return err
		}
	}
	if isPathExist(absoluteDir + filename) {
		if err := os.Remove(absoluteDir + filename); err != nil {
			return err
		}
	}
	file, err := os.Create(absoluteDir + filename + ".tar")
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		res, err := cli.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(len(res.Data))
		n, _ := file.Seek(0, os.SEEK_END)
		file.WriteAt(res.Data, n)
		if err != nil {
			fmt.Printf("failed to write file, err:%v\n", err)
			return err
		}
	}
	return nil
}

func DownloadWithoutStorage(target, path string, data chan<- []byte, close chan string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	defer func() {
	}()
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	defer conn.Close()
	if err != nil {
		close <- err.Error()
		return err
	}
	c := proto.NewTransferServiceClient(conn)
	reqData := &proto.DownloadRequest{
		Path: path,
	}
	cli, err := c.Download(context.Background(), reqData)
	if err != nil {
		close <- err.Error()
		return err
	}
	res, err := cli.Recv()
	if err != nil {
		close <- err.Error()
		return err
	}
LOOP:
	for {
		select {
		case data <- res.Data:
			res, err = cli.Recv()
			if err == io.EOF {
				data <- []byte("done")
				break LOOP
			}
			if err != nil && err != io.EOF {
				data <- []byte("error")
				close <- err.Error()
				break LOOP
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
	return err
}

func Upload(target, path string) error {
	return nil
}

func isPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func spiteDirAndFilename(path string) (dir, filename string) {
	pathSlice := strings.SplitAfter(path, LINUX_SYSTEM_SLASH)
	for i := 0; i < len(pathSlice)-1; i++ {
		dir += pathSlice[i]
	}
	filename = pathSlice[len(pathSlice)-1]
	return
}
