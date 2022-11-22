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

package driver

import (
	"bytes"
	"context"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/util"
	kutil "dtstack.com/dtstack/easymatrix/matrix/k8s/util"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	mutil "dtstack.com/dtstack/easymatrix/matrix/util"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func RkeCreate(clusterName, config string, clusterId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	configf, err := createRkeConfig(clusterName, config)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	rkeBin, err := getRunTimeBin("rke")
	if err != nil {
		return err
	}
	param := []string{"up", "--config", configf}
	cmd := exec.CommandContext(ctx, rkeBin, param...)

	//logf, err := ioutil.TempFile(k8s.ClusterStoreDir, cluster+"*.log")
	logf, err := os.OpenFile(kutil.BuildClusterLogName(clusterName, clusterId), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		return err
	}
	defer logf.Close()
	cmd.Stdout = logf
	cmd.Stderr = logf

	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func createRkeConfig(clusterName, config string) (string, error) {
	f, err := util.NewFile(constant.ClusterStoreDir, clusterName+constant.TEMPLATE_SUFFIX)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err = f.WriteString(config); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func getClusterKubeconfigF(clusterName string) (string, error) {
	config := constant.ClusterStoreDir + "kube_config_" + clusterName + constant.TEMPLATE_SUFFIX
	if !mutil.IsPathExist(config) {
		return "", fmt.Errorf("%v not exist!", config)
	}
	return config, nil
}

func getRunTimeBin(name string) (string, error) {
	bin := constant.RuntimeBinDir + name
	if !mutil.IsPathExist(bin) {
		return "", fmt.Errorf("%v not exist!", bin)
	}
	return bin, nil
}

func DeployWithKubeCtl(clusterName, yaml string) error {
	buf := bytes.NewBufferString(yaml)
	kubeConfig, err := getClusterKubeconfigF(clusterName)
	if err != nil {
		return err
	}
	kubeCtl, err := getRunTimeBin("kubectl")
	if err != nil {
		return err
	}
	cmd := exec.Command(kubeCtl, "--kubeconfig", kubeConfig, "apply", "-f", "-")
	cmd.Stdin = buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DeployWithKubeCtlWithFile(clusterName, file string) error {
	kubeConfig, err := getClusterKubeconfigF(clusterName)
	fmt.Println("-----------------------kubeconfig------------------------",kubeConfig)
	fmt.Println("---------------------------file-----------------------------",file)
	if err != nil {
		return err
	}
	kubeCtl, err := getRunTimeBin("kubectl")
	if err != nil {
		return err
	}
	cmd := exec.Command(kubeCtl, "--kubeconfig", kubeConfig, "apply", "-f", file)
	return cmd.Run()
}
