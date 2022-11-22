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

package workload

import (
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/util"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	matrixSchema "dtstack.com/dtstack/easymatrix/schema"
	"encoding/json"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"path/filepath"
	"strconv"
	"strings"
)

var GVK = schema.GroupVersionKind{
	Group:   "dtstack.com",
	Version: "v1beta1",
	Kind:    "WorkLoad",
}

type Builder struct {
	Def         *modelkube.WorkloadDefinitionSchema
	Parts       []modelkube.WorkloaPartSchema
	Steps       map[int][]modelkube.WorloadStepSchema
	Schema      *matrixSchema.SchemaConfig
	ProductName string
	ServiceName string
	Namespace   string
	Store       *model.ImageStore
}

func New() *workloadv1beta1.WorkLoad {
	return &workloadv1beta1.WorkLoad{}
}

func (g *Builder) Build() (*workloadv1beta1.WorkLoad, error) {
	parts, err := g.buildParts()
	if err != nil {
		return nil, err
	}
	workload := &workloadv1beta1.WorkLoad{
		Spec: workloadv1beta1.WorkLoadSpec{
			WorkLoadParts: parts,
		},
	}
	if err = g.FieldSet(workload); err != nil {
		return nil, err
	}
	if workload.Annotations == nil {
		workload.Annotations = map[string]string{}
	}
	workload.Annotations["type"] = g.Def.Name
	workload.Annotations["version"] = g.Def.Version
	return workload, nil

}

func (g *Builder) buildParts() ([]workloadv1beta1.WorkLoadPart, error) {

	parts := make([]workloadv1beta1.WorkLoadPart, 0, len(g.Parts))

	for _, part := range g.Parts {
		steps, err := g.buildSteps(part.Id, part.Name)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{}
		if err = json.Unmarshal([]byte(part.Parameters), &m); err != nil {
			log.Errorf("[workload] unmashal part parameters %s to map[string]interface{} fail", part.Parameters)
			return nil, err
		}
		g.parse(m, part.Name, g.WorkloadName())

		//处理往容器中添加自定义host映射
		podSpec := map[string]interface{}{}

		hostAlias := g.Schema.Service[g.ServiceName].Instance.HostAlias
		if hostAlias != nil {
			hostAliases := []map[string]interface{}{}

			//["ip:hostname","ip:hostname"]
			ipAndhostname := strings.Split(*hostAlias, ",")

			for _, alias := range ipAndhostname {
				tmpalias := map[string]interface{}{}
				ipHost := strings.Split(alias, ":")

				tmpalias["ip"] = ipHost[0]
				tmpalias["hostnames"] = []string{ipHost[1]}
				hostAliases = append(hostAliases, tmpalias)
			}

			imagePullSecret := map[string]interface{}{
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"imagePullSecrets": []map[string]string{
								{"name": g.Store.Alias},
							},
							"hostAliases": hostAliases,
						},
					},
				},
			}

			podSpec = imagePullSecret
		} else {
			//set imagepullsecret
			imagePullSecret := map[string]interface{}{
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"imagePullSecrets": []map[string]string{
								{"name": g.Store.Alias},
							},
						},
					},
				},
			}

			podSpec = imagePullSecret
		}

		util.DeepCopy(podSpec, m)
		p := workloadv1beta1.WorkLoadPart{
			BaseWorkLoad: workloadv1beta1.BaseWorkLoad{
				Type:       part.Type,
				Name:       part.Name,
				Parameters: workloadv1beta1.Object{Object: m},
			},
			Steps: steps,
		}
		parts = append(parts, p)
	}
	return parts, nil
}

func (g *Builder) buildSteps(partId int, partName string) ([]workloadv1beta1.WorkLoadPartStep, error) {
	steps := g.Steps[partId]
	workloadSteps := make([]workloadv1beta1.WorkLoadPartStep, 0, len(steps))
	for _, step := range steps {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(step.Object), &m); err != nil {
			log.Errorf("[workload] unmashal step Object %s to map[string]interface{} fail", step.Object)
			return nil, err
		}
		//parse @xxx
		m = g.parse(m, partName, g.WorkloadName()).(map[string]interface{})
		st := workloadv1beta1.WorkLoadPartStep{
			Name:   step.Name,
			Type:   step.Type,
			Action: workloadv1beta1.ComponentAction(step.Action),
			Object: workloadv1beta1.Object{Object: m},
		}
		workloadSteps = append(workloadSteps, st)
	}
	return workloadSteps, nil
}

func (g *Builder) FieldSet(workload *workloadv1beta1.WorkLoad) error {
	instance := g.Schema.Service[g.ServiceName].Instance
	if instance == nil {
		return nil
	}
	instanceMap, err := util.ToMap(instance)
	if err != nil {
		log.Errorf("[workload] ServiceConfig to map fail, error %v", err)
		return err
	}
	workloadMap, err := util.ToMap(workload)
	if err != nil {
		log.Errorf("[workload] WorkLoad to map fail,error %v", err)
		return err
	}
	//deal schema's config_path
	if err = g.Schema.ParseServiceVariable(g.ServiceName); err != nil {
		log.Errorf("parse servie %s fail, error: %v", g.ServiceName, err)
		return err
	}
	if len(instance.ConfigPaths) > 0 {
		configmap := map[string]string{}
		baseDir := filepath.Join(base.WebRoot, g.Schema.ProductName, g.Schema.ProductVersion)
		contents, err := g.Schema.ParseServiceConfigFiles(baseDir, g.ServiceName)
		if err != nil {
			log.Errorf("parse %s config variables to config file fail,error:%v", g.ServiceName, err)
			return err
		}
		for index, content := range contents {
			fileName := instance.ConfigPaths[index]

			fileName = fileName[strings.LastIndex(fileName, "/")+1:]
			configmap[fileName] = string(content[:])
		}
		instanceMap["ConfigPaths"] = configmap
	}

	params := g.Def.Params.ToStruct()
	for _, param := range params {
		v, err := FieldGet(instanceMap, param.Key)
		if err != nil {
			return err
		}
		if v == nil {
			continue
		}
		if err = FieldSet(workloadMap, param.Ref, v); err != nil {
			return err
		}
	}
	if err = util.ToObject(workloadMap, workload); err != nil {
		return err
	}
	return nil
}

func (g *Builder) WorkloadName() string {
	return util.BuildWorkloadName(g.ProductName, g.ServiceName)
}

func FieldGet(data interface{}, key string) (interface{}, error) {
	v := data
	for _, field := range strings.Split(key, ".") {
		switch v.(type) {
		case map[string]interface{}:
			m := v.(map[string]interface{})
			v = m[field]
		case []interface{}:
			sl := v.([]interface{})
			index, err := strconv.Atoi(field)
			if err != nil {
				log.Errorf("[resource/workload]: field get convert field to int err, field %s, error %v", field, err)
				return nil, err
			}
			v = sl[index]
		}
	}
	return v, nil
}

func FieldSet(data interface{}, key string, value interface{}) error {
	v := data
	fields := strings.Split(key, ".")
	l := len(fields)

	for i, field := range strings.Split(key, ".") {
		switch v.(type) {
		case map[string]interface{}:
			m := v.(map[string]interface{})
			if i == l-1 {
				m[field] = value
				continue
			}
			v = m[field]
		case []interface{}:
			sl := v.([]interface{})
			index, err := strconv.Atoi(field)
			if err != nil {
				log.Errorf("[resource/workload]: field get convert field to int err, field %s, error %v", field, err)
				return err
			}
			if i == l-1 {
				sl[index] = value
				continue
			}
			v = sl[index]
		}
	}
	return nil
}

func (g *Builder) parse(obj interface{}, partName, workloadName string) interface{} {
	switch obj.(type) {
	case map[string]interface{}:
		m := obj.(map[string]interface{})
		for k, v := range m {
			if k == "image" {
				m[k] = g.Store.Address + "/" + v.(string)
				continue
			}
			m[k] = g.parse(v, partName, workloadName)
		}
		return m
	case []interface{}:
		sl := obj.([]interface{})
		for i, v := range sl {
			sl[i] = g.parse(v, partName, workloadName)
		}
		return sl
	case string:
		s := obj.(string)
		if strings.Contains(s, "@") {
			if s[1:] == partName {
				return util.BuildBaseName(workloadName, partName)
			}
			return util.BuildStepName(util.BuildBaseName(workloadName, partName), s[1:])
		}
		return s
	default:
		return obj
	}
}
