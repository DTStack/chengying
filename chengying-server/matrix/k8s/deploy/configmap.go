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
	"context"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/workload/support"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/util"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/model"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/configmap"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"dtstack.com/dtstack/easymatrix/schema"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	kschema "k8s.io/apimachinery/pkg/runtime/schema"
	"path/filepath"
	"strings"
)

func ApplyConfigMaps(cache kube.ClientCache,sc *schema.SchemaConfig, clusterId int, namespace string) error {

	baseDir := filepath.Join(base.WebRoot, sc.ProductName, sc.ProductVersion)
	for svcName, svc := range sc.Service {
		cfgFiles := map[string]string{}
		if svc.Instance != nil && len(svc.Instance.ConfigPaths) > 0 {
			err := sc.ParseServiceVariable(svcName)
			if err != nil {
				return err
			}
			contents, err := sc.ParseServiceConfigFiles(baseDir, svcName)
			if err != nil {
				return err
			}
			for index, content := range contents {
				cfgFile := svc.Instance.ConfigPaths[index]
				cfgFile = strings.Replace(cfgFile, "/", "_", -1)
				cfgFiles[cfgFile] = string(content[:])
			}
			configMap := model.MakeConfigMap(namespace, sc.ParentProductName, sc.ProductName, svcName, cfgFiles)
			if cache != nil{
				 err := cache.GetClient(namespace).Apply(context.Background(),configMap)
				 if err != nil{
				 	return err
				 }
			}else{
				gvr := &kschema.GroupVersionResource{
					Group:    model.CONFIGMAP_GROUP,
					Version:  model.CONFIGMAP_VERSION,
					Resource: model.CONFIGMAP_RESOURCE,
				}
				configBytes, err := json.Marshal(configMap)
				if err != nil {
					return err
				}
				configDynamic := NewDynamic(configBytes, gvr, model.CONFIGMAP_KIND)
				err = ApplyDynamicResource(configDynamic, clusterId)
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func GetConfigMaps(cache kube.ClientCache,sc *schema.SchemaConfig, clusterId int, namespace, serviceName string) (interface{}, error) {
	var conf *corev1.ConfigMap
	if sc.DeployType == "workload"{
		svc := sc.Service[serviceName]
		workloadVersion := strings.Split(svc.Workload,"@")
		wlTyp := workloadVersion[0]
		wlversion := ""
		if len(workloadVersion) ==2 {
			wlversion = workloadVersion[1]
		}
		def, err := modelkube.WorkloadDefinition.Get(wlTyp,wlversion)
		if err != nil || def == nil{
			return nil, err
		}
		parts,err := modelkube.WorkloadPart.Select(def.Id)
		if err != nil || parts == nil{
			return nil, err
		}
		steps, err := modelkube.WorkloadStep.SelectType(parts[0].Id,support.CreateTypeConfigmap)
		if err != nil || steps == nil{
			return nil, err
		}
		conf = configmap.New()
		conf.Name = util.BuildStepName(util.BuildBaseName(util.BuildWorkloadName(sc.ProductName,serviceName),parts[0].Name),steps[0].Name)
		conf.Namespace = namespace
	}else{
		conf = model.MakeConfigMap(namespace, sc.ParentProductName, sc.ProductName, serviceName, nil)
	}
	if cache != nil{
		exist,err := cache.GetClient(namespace).Get(context.Background(),conf)
		if err != nil{
			return nil,err
		}
		if !exist{
			return nil,nil
		}
		return conf,nil
	}
	gvr := &kschema.GroupVersionResource{
		Group:    model.CONFIGMAP_GROUP,
		Version:  model.CONFIGMAP_VERSION,
		Resource: model.CONFIGMAP_RESOURCE,
	}
	configBytes, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}
	configDynamic := NewDynamic(configBytes, gvr, model.CONFIGMAP_KIND)
	return GetDynamicResource(configDynamic, clusterId)
}
