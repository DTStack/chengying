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

package docker

import (
	"bufio"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"fmt"
	"github.com/heroku/docker-registry-client/registry"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func Login(username, address, password string) error {
	log.Debugf("docker login ...")
	login := exec.Command("docker", "login", "-u", username, address, "-p", password)
	login.Stdout = os.Stdout
	login.Stderr = os.Stderr
	if err := login.Run(); err != nil {
		return err
	}
	return nil
}

func OutputStdLog(std io.Reader, deployUUID string) {
	//实时循环读取输出流中的一行内容
	reader := bufio.NewReader(std)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		log.OutputInfof(deployUUID, "%v", line)
	}
}

func OutputStdLogWithRet(std io.Reader, deployUUID string, output *[]string) {
	//实时循环读取输出流中的一行内容
	reader := bufio.NewReader(std)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		*output = append(*output, line)
		log.OutputInfof(deployUUID, "%v", line)
	}
}

/*
b2d5eeeaba3a: Loading layer [==================================================>]   5.88MB/5.88MB
b5d14f4aebad: Loading layer [==================================================>]  18.24MB/18.24MB
35e4637a9d6c: Loading layer [==================================================>]  3.072kB/3.072kB
7ff80c2c03d5: Loading layer [==================================================>]  4.096kB/4.096kB
6e8309ec6bfd: Loading layer [==================================================>]  3.584kB/3.584kB
a7c11ed26fd5: Loading layer [==================================================>]  7.168kB/7.168kB
7f598054d646: Loading layer [==================================================>]  2.048kB/2.048kB
53062ac1a12b: Loading layer [==================================================>]  3.072kB/3.072kB
92c327b8726e: Loading layer [==================================================>]  1.672MB/1.672MB
0a91289117b4: Loading layer [==================================================>]  3.072kB/3.072kB
20fa73f7ad0f: Loading layer [==================================================>]  5.489MB/5.489MB
464db9ec72d8: Loading layer [==================================================>]  3.433MB/3.433MB
3f666b083f1c: Loading layer [==================================================>]    130kB/130kB
5e7458eed436: Loading layer [==================================================>]  3.584kB/3.584kB
b860ae2f768f: Loading layer [==================================================>]  4.608kB/4.608kB
88e637af7547: Loading layer [==================================================>]  4.608kB/4.608kB
Loaded image: easymanager/manage-front:2.11.5-rel
*/
func Load(file, deployUuid string) (string, error) {
	log.Debugf("docker load ...")
	load := exec.Command("docker", "load", "-i", file)
	stdout, err := load.StdoutPipe()
	output := []string{}
	var image string
	if err == nil {
		go OutputStdLogWithRet(stdout, deployUuid, &output)
	}
	stderr, err := load.StderrPipe()
	if err == nil {
		go OutputStdLog(stderr, deployUuid)
	}
	if err := load.Run(); err != nil {
		return "", err
	}

	for _, line := range output {
		if strings.Contains(line, "Loaded image:") {
			images := strings.Split(line, ": ")
			if len(images) == 2 {
				image = images[1]
			}
			break
		}
	}

	if image != "" {
		image = strings.Replace(image, "\n", "", -1)
	}

	return image, nil
}

func Tag(new, old string) error {
	tag := exec.Command("docker", "tag", old, new)
	tag.Stdout = os.Stdout
	tag.Stderr = os.Stderr
	if err := tag.Run(); err != nil {
		return err
	}
	return nil
}

func Push(image, deployUuid string) error {
	push := exec.Command("docker", "push", image)
	stdout, err := push.StdoutPipe()
	if err == nil {
		go OutputStdLog(stdout, deployUuid)
	}
	stderr, err := push.StderrPipe()
	if err == nil {
		go OutputStdLog(stderr, deployUuid)
	}
	if err := push.Run(); err != nil {
		return err
	}
	return nil
}

func newTransport(transport http.RoundTripper, registryURL, username, password string) *registry.Registry {
	transport = registry.WrapTransport(transport, registryURL, username, password)
	registry := &registry.Registry{
		URL: registryURL,
		Client: &http.Client{
			Transport: transport,
		},
		Logf: registry.Log,
	}
	return registry
}

func NewRegClient(registryURL, username, password string) (*registry.Registry, error) {
	transport := http.DefaultTransport
	url := fmt.Sprintf("http://%s", registryURL)
	registry := newTransport(transport, url, username, password)
	if err := registry.Ping(); err != nil {
		url = fmt.Sprintf("https://%s", registryURL)
		registry = newTransport(transport, url, username, password)
		if err := registry.Ping(); err != nil {
			return nil, err
		}
	}
	return registry, nil
}
