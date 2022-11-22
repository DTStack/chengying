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
	molev1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/mole/v1"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller/model"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/workload"
	"dtstack.com/dtstack/easymatrix/schema"
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apires "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

//type StatusPhase string
//
//type Instance struct {
//	ConfigPaths    []string                 `json:"config_path,omitempty"`
//	Logs           []string                 `json:"logs,omitempty"`
//	DataDir        []string                 `json:"data_dir,omitempty"`
//	Environment    map[string]string        `json:"environment,omitempty"`
//	Cmd            string                   `json:"cmd,omitempty,omitempty"`
//	PrometheusPort string                   `json:"prometheus_port,omitempty"`
//	Ingress        *MoleIngress             `json:"ingress,omitempty"`
//	Service        *MoleService             `json:"service,omitempty"`
//	Deployment     *MoleDeployment          `json:"deployment,omitempty"`
//	Resources      corev1.ResourceRequirements `json:"resources,omitempty"`
//	PostDeploy     string                   `json:"post_deploy,omitempty"`
//}
//
//type ServiceConfig struct {
//	ServiceDisplay  string   `json:"service_display,omitempty"`
//	IsDeployIngress bool     `json:"is_deploy_ingress,omitempty"`
//	IsJob           bool     `json:"is_job,omitempty"`
//	Version         string   `json:"version,omitempty"`
//	Instance        Instance `json:"instance,omitempty"`
//	Group           string   `json:"group,omitempty"`
//	DependsOn       []string `json:"depends_on,omitempty"`
//	BaseProduct     string   `json:"base_product"`
//	BaseService     string   `json:"base_service,omitempty"`
//	BaseParsed      bool     `json:"base_parsed,omitempty"`
//	BaseAttribute   string   `json:"base_attribute,omitempty"`
//}
//
//type SchemaConfig struct {
//	Pid                int                      `json:"pid,omitempty"`
//	ClusterId          int                      `json:"cluster_id,omitempty"`
//	DeployUUid         string                   `json:"deploy_uuid,omitempty"`
//	ParentProductName  string                   `json:"parent_product_name,omitempty"`
//	ProductName        string                   `json:"product_name,omitempty"`
//	ProductNameDisplay string                   `json:"product_name_display,omitempty"`
//	ProductVersion     string                   `json:"product_version,omitempty"`
//	ImagePullSecret    string                   `json:"imagePullSecret,omitempty"`
//	Service            map[string]ServiceConfig `json:"service"`
//}
//
//type MoleIngress struct {
//	Annotations   map[string]string `json:"annotations,omitempty"`
//	Labels        map[string]string `json:"labels,omitempty"`
//	Hostname      string            `json:"hostname,omitempty"`
//	Path          string            `json:"path,omitempty"`
//	Enabled       bool              `json:"enabled,omitempty"`
//	TLSEnabled    bool              `json:"tlsEnabled,omitempty"`
//	TLSSecretName string            `json:"tlsSecretName,omitempty"`
//	TargetPort    string            `json:"targetPort,omitempty"`
//}
//
//type MoleService struct {
//	Annotations map[string]string `json:"annotations,omitempty"`
//	Labels      map[string]string `json:"labels,omitempty"`
//	Type        corev1.ServiceType    `json:"type,omitempty"`
//	Ports       []corev1.ServicePort  `json:"ports,omitempty"`
//}
//
//type MoleDeployment struct {
//	Annotations                   map[string]string      `json:"annotations,omitempty"`
//	Labels                        map[string]string      `json:"labels,omitempty"`
//	Replicas                      int32                  `json:"replicas,omitempty"`
//	Image                         string                 `json:"image,omitempty"`
//	Ports                         []int                  `json:"ports,omitempty"`
//	Containers                    []SideContainer        `json:"containers,omitempty"`
//	NodeSelector                  map[string]string      `json:"nodeSelector,omitempty"`
//	Tolerations                   []corev1.Toleration        `json:"tolerations,omitempty"`
//	Affinity                      *corev1.Affinity           `json:"affinity,omitempty"`
//	SecurityContext               *corev1.PodSecurityContext `json:"securityContext,omitempty"`
//	TerminationGracePeriodSeconds int64                  `json:"terminationGracePeriodSeconds,omitempty"`
//}
//
//type SideContainer struct {
//	Image string `json:"image,omitempty"`
//	Name  string `json:"name,omitempty"`
//}
//
//type MoleSpec struct {
//	Product SchemaConfig `json:"product,omitempty"`
//}
//
//type MoleStatus struct {
//	Phase   StatusPhase `json:"phase"`
//	Message string      `json:"message"`
//}
//
//type Mole struct {
//	metav1.TypeMeta   `json:",inline"`
//	metav1.ObjectMeta `json:"metadata,omitempty"`
//
//	Spec   MoleSpec   `json:"spec,omitempty"`
//	Status MoleStatus `json:"status,omitempty"`
//}

func NewMole(productName, namespace string) *molev1.Mole {
	return &molev1.Mole{
		TypeMeta: metav1.TypeMeta{
			Kind:       MOLE_KIND,
			APIVersion: MOLE_GROUP + "/" + MOLE_VERSION,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      BuildMoleName(productName),
			Namespace: namespace,
		},
	}

}

func FillSchema(m *molev1.Mole, sc *schema.SchemaConfig, uncheckedServices []string, clusterId, pid int, deployUUID, secret string) error{
	m.Spec.Product.Pid = pid
	m.Spec.Product.ClusterId = clusterId
	m.Spec.Product.DeployUUid = deployUUID
	m.Spec.Product.ProductName = sc.ProductName
	m.Spec.Product.ProductVersion = sc.ProductVersion
	m.Spec.Product.ParentProductName = sc.ParentProductName
	m.Spec.Product.ProductNameDisplay = sc.ProductNameDisplay
	m.Spec.Product.ImagePullSecret = secret

	uncheckedSet := make(map[string]bool)
	for _, serviceName := range uncheckedServices {
		uncheckedSet[serviceName] = true
	}
	a, _ := json.Marshal(sc.Service)
	log.Debugf("%+v", string(a))
	m.Spec.Product.Service = make(map[string]molev1.ServiceConfig)
	for name, config := range sc.Service {

		if uncheckedSet[name] || config.BaseProduct != "" {
			continue
		}
		if config.Instance == nil {
			continue
		}
		if config.Workload == workload.PLUGIN{
			continue
		}
		moleConfig := molev1.ServiceConfig{}
		moleConfig.DependsOn = config.DependsOn
		moleConfig.ServiceDisplay = config.ServiceDisplay
		moleConfig.Group = config.Group
		moleConfig.Version = config.Version
		moleConfig.IsDeployIngress = config.IsDeployIngress
		moleConfig.Instance.ConfigPaths = config.Instance.ConfigPaths
		moleConfig.IsJob = config.IsJob
		moleConfig.Instance.PostDeploy = config.Instance.PostDeploy
		moleConfig.Instance.Deployment = new(molev1.MoleDeployment)
		moleConfig.Instance.Deployment.Image = config.Instance.Image
		moleConfig.Instance.Deployment.Replicas = int32(config.Instance.Replica)
		moleConfig.Instance.Deployment.Ports = config.Instance.Ports
		moleConfig.Instance.Logs = config.Instance.Logs
		moleConfig.Instance.Logs = config.Instance.Logs
		moleConfig.Instance.PrometheusPort = config.Instance.PrometheusPort

		containers := make([]molev1.MoleContainer,0,len(config.Instance.PluginInit ))
		if config.Instance.PluginInit != nil{
			moleConfig.Instance.Environment = map[string]string{}
			for _, init := range config.Instance.PluginInit{
				bindPath := strings.Split(init,":")
				pluginName := bindPath[0]
				pluginSvc := sc.Service[pluginName]
				pluginPath := pluginSvc.Instance.PluginPath
				if len(pluginPath) == 0{
					pluginPath = "/plugin"
				}
				containers = append(containers,molev1.MoleContainer{
					Image: pluginSvc.Instance.Image,
					Name:  pluginName,
				})
				moleConfig.Instance.Environment[pluginName] = pluginPath + ":" + bindPath[1]
			}
		}
		moleConfig.Instance.Deployment.Containers = containers

		if config.Instance.HostAlias != nil && len(*config.Instance.HostAlias) != 0{
			if moleConfig.Instance.Environment == nil{
				moleConfig.Instance.Environment = map[string]string{}
			}
			moleConfig.Instance.Environment[model.EnvHostAlias] = *config.Instance.HostAlias
		}

		requests := corev1.ResourceList{}
		for r, q := range config.Instance.ResourceRequest {
			if _, support := SupportResources[r]; !support {
				log.Infof("request %s is not supported resource,wo now support 'cpu' and 'memory',please check", r)
				continue
			}
			quility,err := apires.ParseQuantity(*q)
			if err != nil{
				log.Errorf("service %s request %s value %s can not parse, err: %v",name,r,*q,err)
				return fmt.Errorf("service %s limit %s value %s can not parse, err: %v",name,r,*q,err)
			}
			requests[corev1.ResourceName(r)] = quility
		}
		limits := corev1.ResourceList{}
		for r, l := range config.Instance.ResourceLimit {
			if _, support := SupportResources[r]; !support {
				log.Infof("limit %s is not supported resource,wo now support 'cpu' and 'memory',please check", r)
				continue
			}
			quility,err := apires.ParseQuantity(*l)
			if err != nil{
				log.Errorf("service %s limit %s value %s can not parse, err: %v",name,r,*l,err)
				return fmt.Errorf("service %s limit %s value %s can not parse, err: %v",name,r,*l,err)
			}
			limits[corev1.ResourceName(r)] = quility
		}
		moleConfig.Instance.Resources.Limits = limits
		moleConfig.Instance.Resources.Requests = requests

		if config.Instance.Hostname != "" {
			moleConfig.Instance.Ingress.Hostname = config.Instance.Hostname
		}

		m.Spec.Product.Service[name] = moleConfig
	}
	return nil
}

func BuildMoleName(productName string) string {
	return ConvertDNSRuleName(productName)
}
