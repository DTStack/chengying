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
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/log"
)

var importClusterTemplateFiles = []constant.TemplateFile{
	constant.TPL_CLUSTER_RESOURCE,
	constant.TPL_CLUSTER_RESOURCE_V1BETA1,
}

type ImportClusterGenerator struct {

}

func (g *ImportClusterGenerator) Generate() (map[string][]byte,error){
	yamls := make(map[string][]byte,len(importClusterTemplateFiles))
	asset.ResetImportClusterTemplateWithLocalFile()
	for _,tplName := range importClusterTemplateFiles{
		bts,err := asset.Asset(tplName.FileName)
		if err != nil{
			log.Errorf("[import_cluster]: read cluster resource %s, error : %v",tplName.FileName,err)
			return nil ,err
		}
		yamls[tplName.FileName] = bts
	}
	return yamls,nil
}

func (g *ImportClusterGenerator) GetFileNames() []constant.TemplateFile{
	return importClusterTemplateFiles
}
