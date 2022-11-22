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

package secret

import (
	"context"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"encoding/base64"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)
var GVK = schema.GroupVersionKind{
	Group:   "",
	Version: "v1",
	Kind:    "Secret",
}

func Convert(obj runtime.Object) *corev1.Secret{
	return obj.(*corev1.Secret)
}
func GetDockerConfigJson(tbsc *[]modelkube.DeployClusterImageStoreSchema,namespace string, registryName string) (*corev1.Secret,error){
	dockerConfigJson,err := getDockerConfigAuth(tbsc)
	if err != nil{
		return nil,err
	}
	// 生成镜像仓库认证信息的secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: registryName,
			Namespace:namespace,
		},
		Data:       dockerConfigJson,
		Type:       corev1.SecretTypeDockerConfigJson,
	}
	return secret,nil
}

func getDockerConfigAuth(tbsc *[]modelkube.DeployClusterImageStoreSchema) (map[string][]byte,error) {
	// 创建多个镜像仓库的认证口令信息
	Dockerconfig := make(map[string]interface{},0)
	multiauths := make(map[string]interface{},0)
	for _,registry := range *tbsc{
		src := []byte(registry.Username + ":" + registry.Password)
		dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
		base64.StdEncoding.Encode(dst, src)
		multiauths[registry.Address]=map[string]interface{}{
		    "username": registry.Username,
			"password": registry.Password,
			"auth":     string(dst),
		}

	}
	Dockerconfig["auths"] = multiauths
	//src := []byte(tbsc.Username + ":" + tbsc.Password)
	//dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	//base64.StdEncoding.Encode(dst, src)
	//m := map[string]interface{}{
	//
	//	"auths": map[string]interface{}{
	//		tbsc.Address: map[string]interface{}{
	//			"username": tbsc.Username,
	//			"password": tbsc.Password,
	//			"auth":     string(dst),
	//		},
	//	},
	//}
	bts,err := json.Marshal(Dockerconfig)
	if err != nil{
		log.Errorf("[secret]: json mashal dockerconfigjson err %v",err)
		return nil,err
	}
	return map[string][]byte{corev1.DockerConfigJsonKey:bts},nil
}

func Ping(client kube.Client, namespace string) error{
	ping := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name: "dtstack-dryru",
		},
	}
	if _,err := client.Get(context.Background(),ping);err != nil{
		return err
	}

	if err := client.DryRun(base.Create,ping);err != nil{
		return err
	}
	return nil
}
