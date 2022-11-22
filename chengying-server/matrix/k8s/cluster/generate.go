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

package cluster

import (
	"dtstack.com/dtstack/easymatrix/matrix/asset"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/util"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"encoding/base64"
	"fmt"
	"strconv"
	"text/template"
)

type GeneratorInfo struct {
	Type constant.GenerateType
	HostIp string
	*modelkube.ClusterInfo
}


type Generator interface {
	Generate() (map[string][]byte,error)
	GetFileNames() []constant.TemplateFile
}
// import cluster and self_build_cluster is different
func getGenerater(info *GeneratorInfo) (Generator,error){
	//import cluster
	switch info.Type {
	case constant.TYPE_IMPORT_CLUSTER:
		return &ImportClusterGenerator{},nil
	case constant.TYPE_SELF_BUILD:
		return &SelfBuilClusterGenerator{
			ClusterId:  strconv.Itoa(info.Id),
			CallBackIp: info.HostIp,
			Aid:        -1,
		},nil
	default:
		return nil ,fmt.Errorf("[generate]: unknown mode %d",info.Mode)
	}
}

func GenerateTemplate(info *GeneratorInfo) error{
	gen,err := getGenerater(info)
	if err != nil{
		return err
	}
	yamls,err := gen.Generate()
	if err != nil{
		return err
	}
	for filename ,yaml := range yamls{
		token := generateToken(strconv.Itoa(info.Id),filename)
		f,err := util.NewFile(constant.CLUSTER_TEMPLATE_DIR,token+constant.TEMPLATE_SUFFIX)
		if err != nil{
			return err
		}
		defer f.Close()
		_,err = f.Write(yaml)
		if err != nil{
			log.Errorf("write yaml to file %s fail,error : %v",constant.CLUSTER_TEMPLATE_DIR+token+constant.TEMPLATE_SUFFIX,err)
			return err
		}
	}
	return nil
}

func GetTemplateUrl(info *GeneratorInfo, old bool) ([]byte,error){
	gen,err := getGenerater(info)
	if err != nil{
		return nil,err
	}
	files := gen.GetFileNames()
	for _, file := range files{
		if file.Old == old{
			token := generateToken(strconv.Itoa(info.Id),file.FileName)
			return []byte(host.AgentInstall.StaticHost + constant.TEMPLATES_FILE_SERVER_PREFIX + token + constant.TEMPLATE_SUFFIX), nil
		}
	}
	return nil,fmt.Errorf("the template url of old %v is not found %+v %+v",old,*info,*info.ClusterInfo)
}

func GetTemplateFile(info *GeneratorInfo, old bool) ([]byte, error){
	gen,err := getGenerater(info)
	if err != nil{
		return nil, err
	}
	files := gen.GetFileNames()
	for _, file := range files{
		if file.Old == old{
			token := generateToken(strconv.Itoa(info.Id),file.FileName)
			return []byte(constant.CLUSTER_TEMPLATE_DIR+token+constant.TEMPLATE_SUFFIX),nil
		}
	}
	return nil, fmt.Errorf("the template file for old %v is not found %+v %+v",old,*info,*info.ClusterInfo)
}

func ReadAndParseTemplate(tplFile string) (*template.Template,error){
	bts,err := asset.Asset(tplFile)
	if err != nil{
		log.Errorf("[generate]: read template %s error : %v",tplFile,err)
		return nil,err
	}
	tpl,err := template.New(tplFile).Parse(string(bts))
	if err != nil{
		log.Errorf("[generate]: parse template %s error: %v",tplFile,err)
		return nil,err
	}
	return tpl,nil
}

func generateToken(clusterId,filename string) string {
	return base64.StdEncoding.EncodeToString([]byte(clusterId + filename))
}
