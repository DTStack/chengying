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

package model

import (
	sqlModel "dtstack.com/dtstack/easymatrix/matrix/model"
	"encoding/base64"
	"encoding/json"
	v1 "k8s.io/api/core/v1"
	v2 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getDockerConfigAuth(store sqlModel.ImageStore) map[string]interface{} {
	src := []byte(store.Username + ":" + store.Password)
	maxLen := base64.StdEncoding.EncodedLen(len(src))
	dst := make([]byte, maxLen)
	base64.StdEncoding.Encode(dst, src)

	return map[string]interface{}{
		"auths": map[string]interface{}{
			store.Address: map[string]interface{}{
				"username": store.Username,
				"password": store.Password,
				"auth":     string(dst),
			},
		},
	}
}

func getDockerConfigData(store sqlModel.ImageStore) map[string][]byte {
	src, _ := json.Marshal(getDockerConfigAuth(store))

	return map[string][]byte{
		v1.DockerConfigJsonKey: src,
	}
}

func NewDockerConfigSecret(namespace string, store sqlModel.ImageStore) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: v2.ObjectMeta{
			Name:      store.Alias, //pull image secret name
			Namespace: namespace,
		},
		Data: getDockerConfigData(store),
		Type: v1.SecretTypeDockerConfigJson,
	}
}
