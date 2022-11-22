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

package pkg

import (
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor"
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/crd"
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/listwatch"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

const LOG_PREFIX = "easymonitor"

var (
	DTSTACK_NAMESPACE            = "dtstack-system"
	DEFAULT_NAMESPACE            = "default"
	DTSTACK_LABEL_COM            = "dtstack.com"
	DEFAULT_RESOURCES_POD        = "pods"
	DEFAULT_RESOURCES_SERVICE    = "services"
	DEFAULT_RESOURCES_DEPLOYMENT = "deployments"
	DEFAULT_RESOURCES_INGRESSES  = "ingresses"
	DEFAULT_RESOURCES_EVENTS     = "events"
)

type MatrixConfig struct {
	Host string `config:"host"`
}

type KubeConfig struct {
	Namespace string `config:"namespace"`
}

type Config struct {
	Matrix MatrixConfig `config:"matrix" validate:"required"`
	Kube   KubeConfig   `config:"kube" validate:"required"`
}

func ParseConfig(configFile string, stopCh chan struct{}) error {
	configContent, err := yaml.NewConfigWithFile(configFile, ucfg.PathSep("."))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found!")
		}
		return err
	}
	config := Config{}
	if err := configContent.Unpack(&config); err != nil {
		return err
	}
	log.Config(LOG_PREFIX)
	matrix := config.Matrix
	events.SelfBuildModeTransmitor(matrix.Host, stopCh)
	if err != nil {
		return fmt.Errorf("init matrix config error!")
	}
	kube := config.Kube
	DTSTACK_NAMESPACE = kube.Namespace
	if namespace := os.Getenv("WATCH_NAMESPACE"); namespace != "" {
		DTSTACK_NAMESPACE = namespace
	}
	if len(DTSTACK_NAMESPACE) == 0 {
		return fmt.Errorf("config namespace error, namespace is nil!")
	}
	return nil
}

func StartDefaultNsMonitorController(master, kubeconfig string, stopCh <-chan struct{}) error {
	return StartMonitorController(master, kubeconfig, DTSTACK_NAMESPACE, stopCh)
}

func StartMonitorController(master, kubeconfig string, namespace string, stopCh <-chan struct{}) error {
	return startMonitorController(master, kubeconfig, namespace, stopCh, nil)
}

func StartMonitorControllerWithTransmitor(master, kubeconfig string, namespace string, stopCh <-chan struct{}, transmitor events.TransmitorInterface) error {
	return startMonitorController(master, kubeconfig, namespace, stopCh, transmitor)
}

func startMonitorController(master, kubeconfig string, namespace string, stopCh <-chan struct{}, transmitor events.TransmitorInterface) error {
	log.Infof("Start Resource Monitor...")
	labelSelector := labels.SelectorFromSet(labels.Set(map[string]string{
		"com": DTSTACK_LABEL_COM,
	}))
	var config *rest.Config
	var err error
	if len(kubeconfig) == 0 {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Errorf("[easymonitor]: create inclusterconfig error :%v", err)
			return err
		}
	} else {
		apiconfig, err := clientcmd.Load([]byte(kubeconfig))
		if err != nil {
			log.Errorf("[easymonitor]: load kubeconfig error %v", err.Error())
			return err
		}
		config, err = clientcmd.NewNonInteractiveClientConfig(*apiconfig, "", &clientcmd.ConfigOverrides{}, nil).ClientConfig()
		if err != nil {
			log.Errorf("[easymonitor]: create rest config from kubeconfig error :%v", err.Error())
			return err
		}
	}
	// creates the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	//podMonitor
	err, podListWatch := listwatch.WithConfigResourceNslabelSelector(client.CoreV1().RESTClient(), master, kubeconfig, DEFAULT_RESOURCES_POD, namespace, labelSelector)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	//serviceMonitor
	err, serviceListWatch := listwatch.WithConfigResourceNslabelSelector(client.CoreV1().RESTClient(), master, kubeconfig, DEFAULT_RESOURCES_SERVICE, namespace, labelSelector)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	//deploymentMonitor
	err, deploymentListWatch := listwatch.WithConfigResourceNslabelSelector(client.AppsV1().RESTClient(), master, kubeconfig, DEFAULT_RESOURCES_DEPLOYMENT, namespace, labelSelector)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	//event
	err, eventListWatch := listwatch.WithConfigResourceNslabelSelector(client.CoreV1().RESTClient(), master, kubeconfig, DEFAULT_RESOURCES_EVENTS, namespace, labels.Nothing())
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	// ingressMonitor
	//err, ingressListWatch := listwatch.WithConfigResourceNslabelSelector(client.CoreV1().RESTClient(), master, kubeconfig, DEFAULT_RESOURCES_INGRESSES, DTSTACK_NAMESPACE, labelSelector)
	//if err != nil {
	//	log.Errorf("%v", err.Error())
	//	return err
	//}
	if transmitor == nil {
		transmitor = events.Transmitor
	}
	sController := monitor.NewStaticController(transmitor)
	dController, err := monitor.NewDynamicController(config, namespace, transmitor)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	//add static resource listwatch
	sController.Add(podListWatch, &v1.Pod{})
	sController.Add(serviceListWatch, &v1.Service{})
	sController.Add(deploymentListWatch, &appv1.Deployment{})
	sController.Add(eventListWatch, &v1.Event{})
	//sController.Add(ingressListWatch, &v1beta1.Ingress{})

	//add dynamic resource listwatch

	dController.Add(crd.MoleGvrk)
	dController.Add(crd.WorkloadProcessGvrk)

	go sController.Run(1, stopCh)
	go dController.Run(1, stopCh)

	log.Infof("Start Resource Monitor Ok！")

	return nil
}
