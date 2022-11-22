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

package deploy

import (
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/docker"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/schema"
	"os"
	"path/filepath"
	"strings"
)

const (
	IMAGE_DIR    = "images"
	IMAGE_SUFFIX = ".tar"
)

func PushImages(store model.ImageStore, sc *schema.SchemaConfig, deployUUID string) error {
	log.Infof("starting push images to registry %v", store.Address)
	log.OutputInfof(deployUUID, "starting push images to registry %v", store.Address)
	err := docker.Login(store.Username, store.Address, store.Password)
	if err != nil {
		log.Errorf("docker login err: %v", err.Error())
		log.OutputInfof(deployUUID, "docker login error: %v", err.Error())
		return err
	}
	log.Infof("docker login success: %v", store.Address)
	log.OutputInfof(deployUUID, "docker login success: %v", store.Address)
	//regURL := strings.Split(store.Address, "/")[0]
	//hub, err := docker.NewRegClient(regURL, store.Username, store.Password)
	//if err != nil {
	//	log.Errorf("creat registry client err: %v", err.Error())
	//	log.OutputInfof(deployUUID, "creat registry client err: %v", err.Error())
	//}
	for name := range sc.Service {
		baseDir := filepath.Join(base.WebRoot, sc.ProductName, sc.ProductVersion, name, IMAGE_DIR)
		if sc.Service[name].Instance == nil || sc.Service[name].Instance.Image == "" {
			log.Infof("")
			continue
		}
		if sc.Service[name].BaseProduct != "" || sc.Service[name].BaseService != "" {
			continue
		}
		sourceImg := sc.Service[name].Instance.Image
		var SourceTag string
		var NewSchemaDefImg string
		var AttachImg string
		var newImgTag string
		log.Infof("push images with images dir: %v", baseDir)
		log.OutputInfof(deployUUID, "push images with images dir: %v", baseDir)
		err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			log.Infof("find file: %v", path)
			log.OutputInfof(deployUUID, "find file: %v", path)
			if err != nil {
				return err
			}
			if baseDir == path {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, IMAGE_SUFFIX) {
				log.Infof("not regular image file: %v", path)
				log.OutputInfof(deployUUID, "not regular image file: %v", path)
				return nil
			}
			//获取倒入镜像的名称，这里镜像分两种，一种是在schema中定义的镜像名称，另一种是产品包中附带的额外的镜像
			//schema中定义的镜像需要拼接上仓库地址重新赋值给schema文件，附带的镜像则不需要只需推送到仓库即可
			imageName, err := docker.Load(path, deployUUID)
			if err != nil {
				return err
			}
			if strings.ToLower(sourceImg) == strings.ToLower(imageName) {
				log.OutputInfof(deployUUID, "========= schema defined image are the same as load image")
				log.OutputInfof(deployUUID, "========= sourceImg is: %v", sourceImg)
				log.OutputInfof(deployUUID, "========= load images is: %v", imageName)
				SourceTag = sourceImg
				log.OutputInfof(deployUUID, "========= SourceTag is: %v", SourceTag)
				NewSchemaDefImg = store.Address + "/" + sourceImg
				newImgTag = NewSchemaDefImg
			} else {
				log.OutputInfof(deployUUID, "========= schema defined image are different from load image")
				log.OutputInfof(deployUUID, "========= sourceImg is: %v", sourceImg)
				log.OutputInfof(deployUUID, "========= load images is: %v", imageName)
				SourceTag = imageName
				log.OutputInfof(deployUUID, "========= SourceTag is: %v", SourceTag)
				AttachImg = store.Address + "/" + imageName
				newImgTag = AttachImg
			}
			log.Infof("load image %v success", path)
			log.OutputInfof(deployUUID, "load image %v success", path)

			log.OutputInfof(deployUUID, "tag %s to %s", SourceTag, newImgTag)
			err = docker.Tag(newImgTag, SourceTag)
			if err != nil {
				//the image can be changed in the front, should not affect the other image's push
				log.OutputInfof(deployUUID, "use the declare image %s in the front", sourceImg)
				return nil
			}
			log.Infof("tag %s to %s success", SourceTag, newImgTag)
			log.OutputInfof(deployUUID, "tag %s to %s success", SourceTag, newImgTag)
			//var exists bool
			//imgName := strings.SplitN(strings.Split(newImgTag, ":")[0], "/", 2)[1]
			//imgTag := strings.Split(newImgTag, ":")[1]
			//if hub != nil {
			//	tags, err := hub.Tags(imgName)
			//	if err != nil {
			//		log.Infof("search image tag err: %v", err)
			//	}
			//	for _, tag := range tags {
			//		if tag == imgTag {
			//			exists = true
			//			break
			//		}
			//	}
			//	if !exists {
			//		log.OutputInfof(deployUUID, "the image %v does not exist in the docker repository,start pushing the image ...", newImg)
			//		err = docker.Push(newImg, deployUUID)
			//		if err != nil {
			//			return err
			//		}
			//		log.Infof("push image %v success", path)
			//		log.OutputInfof(deployUUID, "push image %v success", path)
			//	} else {
			//		log.OutputInfof(deployUUID, "the image %v exists in the docker repository,skip pushing the image!", newImg)
			//	}
			//} else {
			//	err = docker.Push(newImg, deployUUID)
			//	if err != nil {
			//		return err
			//	}
			//	log.Infof("push image %v success", path)
			//	log.OutputInfof(deployUUID, "push image %v success", path)
			//}
			err = docker.Push(newImgTag, deployUUID)
			if err != nil {
				return err
			}

			log.Infof("push image %v success", newImgTag)
			log.OutputInfof(deployUUID, "push image %v success", newImgTag)
			sc.Service[name].Instance.Image = NewSchemaDefImg
			return nil
		})
		if err != nil {
			log.Errorf("push images err: %v", err.Error())
			log.OutputInfof(deployUUID, "push images error: %v", err.Error())
			return err
		}
	}
	return nil
}
