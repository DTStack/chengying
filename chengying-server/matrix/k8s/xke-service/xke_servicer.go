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

package xke_service

import "dtstack.com/dtstack/easymatrix/matrix/k8s/xke-service/driver"

type XkeServicer interface {
	Create(name, config string, id int) error
	Deploy(name, yaml string) error
	DeployWithF(name, file string) error
}

type xkeService struct {
}

func NewXkeService() (XkeServicer, error) {
	newService := &xkeService{}
	return newService, nil
}

func (this *xkeService) Create(clusterName, config string, clusterId int) error {
	return driver.RkeCreate(clusterName, config, clusterId)
}

func (this *xkeService) Deploy(clusterName, yaml string) error {
	return driver.DeployWithKubeCtl(clusterName, yaml)
}

func (this *xkeService) DeployWithF(clusterName, file string) error {
	return driver.DeployWithKubeCtlWithFile(clusterName, file)
}
