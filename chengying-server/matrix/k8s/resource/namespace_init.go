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

package resource

import (
	"bytes"
	"context"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/model"
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/configmap"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/deployment"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/role"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/rolebinding"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/secret"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/serviceaccount"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"fmt"
	"html/template"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"strconv"
	"strings"
)

func Save(ctx context.Context, clusterID string, user string, vo *view.NamespaceSaveReq) error {
	status := constant.NAMESPACE_NOT_CONNECT
	var cache kube.ClientCache
	var lastCache kube.ClientCache
	var err error
	//init clientcache
	if vo.Type == kube.IMPORT_KUBECONFIG.String() {
		cache, err = kube.ClusterNsClientCache.GetClusterNsClient(clusterID).GetClientCache(kube.IMPORT_KUBECONFIG)
		if err != nil {
			return err
		}
		if cache != nil {
			lastCache = cache.Copy()
		}
		if err = kubeconfigClientInit(ctx, cache, vo); err != nil {
			return err
		}
	} else if vo.Type == kube.IMPORT_AGENT.String() {
		if len(vo.Ip) != 0 && len(vo.Port) != 0 {
			cache, err = kube.ClusterNsClientCache.GetClusterNsClient(clusterID).GetClientCache(kube.IMPORT_AGENT)
			if err != nil {
				return err
			}
			if cache != nil {
				lastCache = cache.Copy()
			}
			if err = agentClientInit(ctx, cache, vo); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("[namespace_init]: unknow import type %s", vo.Type)
	}
	cid, err := strconv.Atoi(clusterID)
	if err != nil {
		return fmt.Errorf("[namespace_init]: clusterid %s is illegal", clusterID)
	}
	// apply init resource
	if cache != nil {
		namespace := vo.Namespace
		registryid := vo.RegistryId
		nsTbsc, err := modelkube.DeployNamespaceList.Get(vo.Namespace, cid)
		if err != nil {
			return err
		}
		//when input the ip and port, it will in
		if vo.Type == kube.IMPORT_AGENT.String() {
			if nsTbsc == nil {
				return fmt.Errorf("something is wrong, the namespace should be exist")
			}
			registryid = nsTbsc.RegistryId
		} else {
			if err := ping(cache.GetClient(namespace), namespace); err != nil {
				return err
			}
		}
		//either new or update namespace, apply the optional and required object
		//namespaceInit create mole operator deployment
		err = namespaceInit(ctx, cache, namespace, registryid, cid)
		if err != nil {
			if lastCache != nil {
				kube.ClusterNsClientCache.GetClusterNsClient(clusterID).PutNsClient(lastCache)
			}
			return err
		}
		status = constant.NAMESPACE_VALID
	}

	//operate mysql
	//add the ip and port is specail
	if len(vo.Ip) != 0 && len(vo.Port) != 0 {
		nsTbsc, err := modelkube.DeployNamespaceList.Get(vo.Namespace, cid)
		if err != nil {
			return err
		}
		if nsTbsc == nil {
			return fmt.Errorf("update ip and port something is wrong, the namespace should be exist")
		}
		nsTbsc.Ip = vo.Ip
		nsTbsc.Port = vo.Port
		nsTbsc.Status = status
		vo = &nsTbsc.NamespaceSaveReq
		err = modelkube.DeployNamespaceList.Update(nsTbsc)
		if err != nil {
			return err
		}
	} else {
		//new namespace,insert namespace record into deploy_namespace_list and deploy_namespace_client
		if vo.Id == 0 {
			tbsc := &modelkube.DeployNamespaceListSchema{
				ClusterId:        cid,
				NamespaceSaveReq: *vo,
				IsDeleted:        0,
				User:             user,
				Status:           status,
			}
			id, err := modelkube.DeployNamespaceList.Insert(tbsc)
			if err != nil {
				return err
			}
			vo.Id = id
			nsclientSc := &modelkube.DeployNamespaceClientSchema{
				Yaml:        vo.Yaml,
				NamespaceId: int(id),
				Filename:    vo.FileName,
			}
			if err = modelkube.DeployNamespaceClient.Insert(nsclientSc); err != nil {
				return err
			}
		} else {
			//update namespace
			nsTbsc, err := modelkube.DeployNamespaceList.Get(vo.Namespace, cid)
			if err != nil {
				return err
			}
			if nsTbsc == nil {
				return fmt.Errorf("update namespace something is wrong, the namespace should be exist")
			}
			vo.Ip = nsTbsc.Ip
			vo.Port = nsTbsc.Port
			nsTbsc.Status = status
			nsTbsc.User = user
			nsTbsc.NamespaceSaveReq = *vo
			err = modelkube.DeployNamespaceList.Update(nsTbsc)
			if err != nil {
				return err
			}
			nsclientSc := &modelkube.DeployNamespaceClientSchema{
				Yaml:        vo.Yaml,
				NamespaceId: nsTbsc.Id,
				Filename:    vo.FileName,
			}
			if err = modelkube.DeployNamespaceClient.Update(nsclientSc); err != nil {
				return err
			}
		}
	}
	//if the namespaced client is built, do some goroutine job(start easymonitor for list and watch pod、deployment、service、event resource)
	if cache != nil {
		if err = StartGoroutines(ctx, clusterID, vo); err != nil {
			return err
		}
	}
	return nil
}

func agentClientInit(ctx context.Context, cache kube.ClientCache, vo *view.NamespaceSaveReq) error {
	ip := vo.Ip
	port := vo.Port
	if !strings.HasPrefix(ip, "http") {
		ip = "http://" + ip
	}
	host := ip + ":" + port
	cache.Connect(host, vo.Namespace)
	return nil

}

//use kubeconfig to build client with k8s
func kubeconfigClientInit(ctx context.Context, cache kube.ClientCache, vo *view.NamespaceSaveReq) error {
	err := cache.Connect(vo.Yaml, vo.Namespace)
	if err != nil {
		return err
	}
	return nil
}

//install dockerconfigjson secret, operator deployment required
//sa,role,rolebinding optional
func namespaceInit(ctx context.Context, cache kube.ClientCache, namespace string, registryId, clusterid int) error {
	logType := "promtail"
	logServerAddress := "loki:3000"
	logImage := "promtail:latest"

	_, exist := os.LookupEnv("LOG_SERVER_ADDRESS")
	if exist {
		logServerAddress = os.Getenv("LOG_SERVER_ADDRESS")
	}

	_, exist = os.LookupEnv("LOG_IMAGE")
	if exist {
		logImage = os.Getenv("LOG_IMAGE")
		if !strings.Contains(logImage, "promtail") {
			logType = "filebeat"
		}
	}

	logSwitch, err := strconv.ParseBool(os.Getenv(model.ENV_LOG_SWITCH))
	if err != nil {
		logSwitch = false
	}

	c := cache.GetClient(namespace)
	// get mole operator resource template(deployment、role、rolebinding、sa)from import_init_moudle table
	moudle, err := modelkube.ImportInitMoudle.GetInitMoudle()
	if err != nil {
		return err
	}
	// get mole operator k8s resource object(role、rolebinding、sa、configmap) from template
	objects, err := generateOptional(moudle, namespace)
	if err != nil {
		return err
	}
	//if saName = "",the cluster can't create sa,role and rolebinding, that means the namespace has default sa for use.
	saName := ""
	for k, obj := range objects {

		if err := c.Apply(ctx, obj); err != nil {
			saName = ""
			break
		}
		if k == serviceaccount.GVK.Kind {
			sa := serviceaccount.Convert(obj)
			saName = sa.Name
		}
	}
	// get mole operator k8s resource object(secret、deployment) from template
	tbsc, err := modelkube.DeployClusterImageStore.GetById(registryId)
	if err != nil {
		return err
	}
	if tbsc == nil {
		return fmt.Errorf("the registry of %d is nil in the database", registryId)
	}

	objs, err := generateRequired(moudle, namespace, tbsc, clusterid)
	if err != nil {
		return err
	}

	for k, v := range objs {
		if len(saName) != 0 && k == deployment.GVK.Kind {
			deploy := deployment.Convert(v)
			deploy.Spec.Template.Spec.ServiceAccountName = saName

			currentEnvs := deploy.Spec.Template.Spec.Containers[0].Env
			switchEnv := corev1.EnvVar{
				Name:  model.ENV_LOG_SWITCH,
				Value: strconv.FormatBool(logSwitch),
			}
			logEnvVars := []corev1.EnvVar{}
			if logSwitch {
				logEnvVars = []corev1.EnvVar{
					{
						Name:  "LOG_IMAGE",
						Value: fmt.Sprintf("%s/%s", tbsc.Address, logImage),
					},
					{
						Name:  "LOG_TYPE",
						Value: logType,
					},
					{
						Name:  "LOG_SERVER_ADDRESS",
						Value: logServerAddress,
					},
				}
			}
			currentEnvs = append(currentEnvs, switchEnv)
			currentEnvs = append(currentEnvs, logEnvVars...)
			deploy.Spec.Template.Spec.Containers[0].Env = currentEnvs
		}
		if err := c.Apply(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

func getDefaultWithEnv(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return v
}

func createMoleImage(registryUrl string) string {
	image := getDefaultWithEnv("MOLE_IMAGE", "easymatrix/operator:apadp1.23.3")
	if !strings.HasSuffix(registryUrl, "/") {
		registryUrl = registryUrl + "/"
	}
	if strings.HasPrefix(image, "/") {
		image = strings.TrimPrefix(image, "/")
	}
	return registryUrl + image
}

//secret of type dockerconfigjson and mole operator is required
func generateRequired(sc *modelkube.ImportInitMoudleSchema, namespace string, tbsc *modelkube.DeployClusterImageStoreSchema, clusterId int) (map[string]runtime.Object, error) {
	// dockerconfigjson secret
	objects := map[string]runtime.Object{}
	//tbsc, err := modelkube.DeployClusterImageStore.GetById(registryId)
	clusterImageStore, err := modelkube.DeployClusterImageStore.GetByClusterId(clusterId)
	if err != nil {
		return nil, err
	}
	//if tbsc == nil {
	//	return nil, fmt.Errorf("the registry of %d is nil in the database", registryId)
	//}
	se, err := secret.GetDockerConfigJson(clusterImageStore, namespace, tbsc.Alias)
	if err != nil {
		return nil, err
	}
	objects[secret.GVK.Kind] = se

	// mole operator
	var buf bytes.Buffer
	if len(sc.Operator) == 0 {
		return nil, fmt.Errorf("[namespace_init]: operator is not exist, please check")
	}
	data := map[string]string{
		"NAME_SPACE":  namespace,
		"MOLE_IMAGE":  createMoleImage(tbsc.Address),
		"SECRET_NAME": se.Name,
	}
	if err := generateTemplate("operator", sc.Operator, data, &buf); err != nil {
		return nil, err
	}
	deploy, err := deployment.ToObject(buf.Bytes())
	if err != nil {
		return nil, err
	}
	objects[deployment.GVK.Kind] = deploy
	return objects, nil
}

//servicaaccount,role and rolebinding is optional.
func generateOptional(sc *modelkube.ImportInitMoudleSchema, namespace string) (map[string]runtime.Object, error) {

	var buf bytes.Buffer
	objects := map[string]runtime.Object{}
	data := map[string]string{"NAME_SPACE": namespace}
	// serviceaccount
	if len(sc.ServiceAccount) != 0 {
		if err := generateTemplate("sa", sc.ServiceAccount, data, &buf); err != nil {
			return nil, err
		}
		sa, err := serviceaccount.ToObject(buf.Bytes())
		if err != nil {
			return nil, err
		}
		buf.Reset()
		objects[serviceaccount.GVK.Kind] = sa
	}
	//role
	if len(sc.Role) != 0 {
		if err := generateTemplate("role", sc.Role, data, &buf); err != nil {
			return nil, err
		}
		r, err := role.ToObject(buf.Bytes())
		if err != nil {
			return nil, err
		}
		buf.Reset()
		objects[role.GVK.Kind] = r
	}
	//rolebinding
	if len(sc.RoleBinding) != 0 {
		if err := generateTemplate("rolebinding", sc.RoleBinding, data, &buf); err != nil {
			return nil, err
		}
		rb, err := rolebinding.ToObject(buf.Bytes())
		if err != nil {
			return nil, err
		}
		buf.Reset()
		objects[rolebinding.GVK.Kind] = rb
	}
	// configmap for log
	if len(sc.LogConfig) != 0 {
		if err := generateTemplate("logconfig", sc.LogConfig, data, &buf); err != nil {
			return nil, err
		}
		config, err := configmap.ToObject(buf.Bytes())
		if err != nil {
			return nil, err
		}
		objects[configmap.GVK.Kind] = config
	}
	return objects, nil
}

func generateTemplate(tplname string, tpl string, data interface{}, buf *bytes.Buffer) error {
	t, err := template.New(tplname).Parse(tpl)
	if err != nil {
		log.Errorf("[namespace_init]: parse %s template error: %v", tplname, err)
		return err
	}
	err = t.Execute(buf, data)
	if err != nil {
		log.Errorf("[namespace_init] generate %s template error: %v", tplname, err)
		return err
	}
	return nil
}
