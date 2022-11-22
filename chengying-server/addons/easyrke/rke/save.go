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

package rke

import (
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func save(img interface{}) error {
	if _, ok := img.(string); !ok {
		return fmt.Errorf("error format: %v", img)
	}
	image := img.(string)
	log.Infof("pull image: %v", image)
	fmt.Println("pull image:", image)
	pull := exec.Command("docker", "pull", image)
	pull.Stdout = os.Stdout
	pull.Stderr = os.Stderr
	if err := pull.Run(); err != nil {
		return err
	}
	log.Infof("save image: %v", image)
	fmt.Println("save image:", image)
	save := exec.Command("docker", "save", "-o", "images/"+strings.Replace(image, "/", "_", -1)+".tar", image)
	save.Stdout = os.Stdout
	save.Stderr = os.Stderr
	if err := save.Run(); err != nil {
		return err
	}
	log.Infof("save image: %v success", image)
	fmt.Println("save image:", image, "success")
	return nil
}

func SaveImage(k8sVersion string) {
	images := LoadK8sRKESystemImages()

	if _, ok := images[k8sVersion]; !ok {
		fmt.Errorf("no target k8s  version %v, support version list: %v", k8sVersion, reflect.ValueOf(images).MapKeys())
		log.Errorf("no target k8s  version %v, support version list: %v", k8sVersion, reflect.ValueOf(images).MapKeys())
		return
	}
	target := images[k8sVersion]
	ims := reflect.ValueOf(&target).Elem()
	for i := 0; i < ims.NumField(); i++ {
		image := ims.Field(i).Interface()
		fmt.Println("get image:", image)
		log.Infof("get image: %v", image)
		if err := save(image); err != nil {
			fmt.Errorf("err %v", err.Error())
			log.Errorf("err %v", err.Error())
		}
	}
	log.Infof("Finished %v", k8sVersion)
}
