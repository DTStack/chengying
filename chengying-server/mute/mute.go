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

package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	input := flag.String("n", "operator.dtstack.com_moles_crd.yaml", "The name of the file to be modified")
	output := flag.String("o", "operator.dtstack.com_moles_crd.yaml", "The name of the modified file")
	flag.Parse()

	obj := make(map[interface{}]interface{})
	yamlFile, _ := ioutil.ReadFile(*input)
	_ = yaml.Unmarshal(yamlFile, &obj)
	in := RemoveDescription(obj)
	out, _ := yaml.Marshal(in)

	f, _ := os.Create(*output)
	defer f.Close()
	_, _ = f.Write(out)
}

func RemoveDescription(obj interface{}) interface{} {
	switch valueType := obj.(type) {
	case map[interface{}]interface{}:
		result := make(map[interface{}]interface{})
		for k, v := range valueType {
			if k == "description" {
				continue
			}
			result[k] = RemoveDescription(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(valueType))
		for i, v := range valueType {
			result[i] = RemoveDescription(v)
		}
		return result
	default:
		return valueType
	}
}
