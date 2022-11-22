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
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
	"strings"
)

func getAffinities(cr *molev1.Mole, name string) *corev1.Affinity {
	var affinity = corev1.Affinity{}
	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.Affinity != nil {
		affinity = *cr.Spec.Product.Service[name].Instance.Deployment.Affinity
	}
	return &affinity
}

func getSecurityContext(cr *molev1.Mole, name string) *corev1.PodSecurityContext {
	var securityContext = corev1.PodSecurityContext{}
	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.SecurityContext != nil {
		securityContext = *cr.Spec.Product.Service[name].Instance.Deployment.SecurityContext
	}
	return &securityContext
}

func getReplicas(cr *molev1.Mole, name string) *int32 {
	var replicas int32 = 1
	if cr.Spec.Product.Service[name].Instance.Deployment == nil {
		return &replicas
	}
	if cr.Spec.Product.Service[name].Instance.Deployment.Replicas <= 0 {
		return &replicas
	} else {
		return &cr.Spec.Product.Service[name].Instance.Deployment.Replicas
	}
}

func getRollingUpdateStrategy() *appsv1.RollingUpdateDeployment {
	var maxUnavailable = intstr.FromInt(0)
	var maxSurge = intstr.FromString("25%")
	return &appsv1.RollingUpdateDeployment{
		MaxUnavailable: &maxUnavailable,
		MaxSurge:       &maxSurge,
	}
}

func getPodAnnotations(cr *molev1.Mole, existing map[string]string, name string) map[string]string {
	var annotations = map[string]string{}
	// Add fixed annotations
	annotations["prometheus.io/scrape"] = "true"
	//annotations["prometheus.io/port"] = fmt.Sprintf("%v", GetMolePort(cr, name))
	annotations = MergeAnnotations(annotations, existing)

	if cr.Spec.Product.Service[name].Instance.Deployment != nil {
		annotations = MergeAnnotations(cr.Spec.Product.Service[name].Instance.Deployment.Annotations, annotations)
	}
	return annotations
}

func getPodLabels(cr *molev1.Mole, name string) map[string]string {
	var labels = map[string]string{}

	labels["app"] = BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name)
	labels["pid"] = strconv.Itoa(cr.Spec.Product.Pid)
	labels["deploy_uuid"] = cr.Spec.Product.DeployUUid
	labels["cluster_id"] = strconv.Itoa(cr.Spec.Product.ClusterId)
	labels["product_name"] = cr.Spec.Product.ProductName
	labels["product_version"] = cr.Spec.Product.ProductVersion
	labels["parent_product_name"] = cr.Spec.Product.ParentProductName
	labels["service_name"] = name
	labels["service_version"] = cr.Spec.Product.Service[name].Version
	labels["group"] = cr.Spec.Product.Service[name].Group
	labels["com"] = MoleCom

	return labels
}
func getDeploymentLabels(cr *molev1.Mole, name string) map[string]string {
	var labels = map[string]string{}

	labels["pid"] = strconv.Itoa(cr.Spec.Product.Pid)
	labels["deploy_uuid"] = cr.Spec.Product.DeployUUid
	labels["cluster_id"] = strconv.Itoa(cr.Spec.Product.ClusterId)
	labels["product_name"] = cr.Spec.Product.ProductName
	labels["product_version"] = cr.Spec.Product.ProductVersion
	labels["parent_product_name"] = cr.Spec.Product.ParentProductName
	labels["service_name"] = name
	labels["service_version"] = cr.Spec.Product.Service[name].Version
	labels["group"] = cr.Spec.Product.Service[name].Group
	labels["com"] = MoleCom

	return labels
}

func getNodeSelectors(cr *molev1.Mole, name string) map[string]string {
	var nodeSelector = map[string]string{}

	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.NodeSelector != nil {
		nodeSelector = cr.Spec.Product.Service[name].Instance.Deployment.NodeSelector
	}
	return nodeSelector

}

func getTerminationGracePeriod(cr *molev1.Mole, name string) *int64 {
	var tcp int64 = 30
	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.TerminationGracePeriodSeconds != 0 {
		tcp = cr.Spec.Product.Service[name].Instance.Deployment.TerminationGracePeriodSeconds
	}
	return &tcp

}

func getTolerations(cr *molev1.Mole, name string) []corev1.Toleration {
	tolerations := []corev1.Toleration{}

	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.Tolerations != nil {
		for _, val := range cr.Spec.Product.Service[name].Instance.Deployment.Tolerations {
			tolerations = append(tolerations, val)
		}
	}
	return tolerations
}

func getVolumes(cr *molev1.Mole, name string) []corev1.Volume {
	var volumes []corev1.Volume
	// Volume to mount the config file from a configMap
	volumes = append(volumes, corev1.Volume{
		Name: BuildResourceName(MoleConfigVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: BuildResourceName(MoleConfigName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				},
				DefaultMode: &VolumeConfigMapMode,
			},
		},
	})

	if len(cr.Spec.Product.Service[name].Instance.Logs) != 0 {
		logConfigMapVolume := corev1.Volume{
			Name: "log-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: Configmap,
					},
					DefaultMode: &VolumeConfigMapMode,
				},
			},
		}
		commonLogVolume := corev1.Volume{
			Name: LogVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
		return append(volumes, logConfigMapVolume, commonLogVolume)
	}
	deploy := cr.Spec.Product.Service[name].Instance.Deployment
	if deploy != nil && len(deploy.Containers) != 0 {
		commonLogVolume := corev1.Volume{
			Name: LogVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
		return append(volumes, commonLogVolume)
	}

	return volumes

}

func getVolumeMounts(cr *molev1.Mole, name string) []corev1.VolumeMount {
	var log = logf.Log.WithName("volumeMount")
	var mounts []corev1.VolumeMount
	for _, configPath := range cr.Spec.Product.Service[name].Instance.ConfigPaths {
		subPath := strings.Replace(configPath, "/", "_", -1)
		mounts = append(mounts, corev1.VolumeMount{
			Name:      BuildResourceName(MoleConfigVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			SubPath:   subPath,
			MountPath: fmt.Sprintf("opt/dtstack/%v/%v/%v", cr.Spec.Product.ProductName, name, configPath),
		})
	}

	if len(cr.Spec.Product.Service[name].Instance.Logs) == 0 {
		log.Info(fmt.Sprintf("service: %s no logs path", name))
		return mounts
	}

	log.Info(fmt.Sprintf("logpath name: %s paths: %+v\n", name, cr.Spec.Product.Service[name].Instance.Logs))
	for _, logPath := range cr.Spec.Product.Service[name].Instance.Logs {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      LogVolume,
			MountPath: fmt.Sprintf("/opt/dtstack/%v/%v/%v", cr.Spec.Product.ProductName, name, logPath[:strings.LastIndex(logPath, "/")]),
			SubPath:   LogSubpath,
		})
	}
	log.Info(fmt.Sprintf("mounts %+v", mounts))
	return mounts
}

// it is not a good way to deal plugin
// but mole will be Deprecated in the future
// the molecontainer and enviroment is not used before, use these two parameters to realize,
// although it is not a way
func getInitPlugin(cr *molev1.Mole, name string) []corev1.Container {

	if len(cr.Spec.Product.Service[name].Instance.Deployment.Containers) == 0 {
		return nil
	}
	containers := make([]corev1.Container, 0, len(cr.Spec.Product.Service[name].Instance.Deployment.Containers))

	env := cr.Spec.Product.Service[name].Instance.Environment
	for _, c := range cr.Spec.Product.Service[name].Instance.Deployment.Containers {

		sourceAndTarget := strings.Split(env[c.Name], ":")
		source := sourceAndTarget[0]
		target := sourceAndTarget[1]
		subPath := target[strings.LastIndex(target, "/")+1:]
		mountPath := "/tmp/" + subPath
		init := corev1.Container{
			Name:            strings.ToLower(strings.Replace(c.Name, "_", "-", -1)),
			Image:           c.Image,
			ImagePullPolicy: "Always",
			Env: []corev1.EnvVar{
				{
					Name:  "PLUGIN_PATH",
					Value: source,
				},
			},
			Command: []string{
				"/bin/sh",
			},
			Args: []string{
				"-c",
				"cd " + mountPath + " && rm -rf ./* && mv ${PLUGIN_PATH}/* " + mountPath + "/",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: mountPath,
					Name:      LogVolume,
					SubPath:   subPath,
				},
			},

			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: apiresource.MustParse("50Mi"),
					corev1.ResourceCPU:    apiresource.MustParse("50m"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: apiresource.MustParse("0Mi"),
					corev1.ResourceCPU:    apiresource.MustParse("0m"),
				},
			},
		}
		containers = append(containers, init)
	}
	return containers
}

func getContainers(cr *molev1.Mole, name string) []corev1.Container {
	var log = logf.Log.WithName("containers")

	container := corev1.Container{
		Name:            ConvertDNSRuleName(name),
		Image:           cr.Spec.Product.Service[name].Instance.Deployment.Image,
		WorkingDir:      "",
		Ports:           getContainerPorts(cr, name),
		VolumeMounts:    getVolumeMounts(cr, name),
		Resources:       getResources(cr, name),
		ImagePullPolicy: "Always",
	}

	env := cr.Spec.Product.Service[name].Instance.Environment

	for _, c := range cr.Spec.Product.Service[name].Instance.Deployment.Containers {
		sourceAndTarget := strings.Split(env[c.Name], ":")
		target := sourceAndTarget[1]
		subPath := target[strings.LastIndex(target, "/")+1:]
		vm := corev1.VolumeMount{
			Name:      LogVolume,
			MountPath: target,
			SubPath:   subPath,
		}
		container.VolumeMounts = append(container.VolumeMounts, vm)
	}

	containers := []corev1.Container{container}

	//for _, container := range cr.Spec.Product.Service[name].Instance.Deployment.Containers {
	//	containers = append(containers, corev1.Container{
	//		Name:            container.Name,
	//		Image:           container.Image,
	//		VolumeMounts:    getVolumeMounts(cr, name),
	//		ImagePullPolicy: "IfNotPresent",
	//	})
	//}
	logSwitch := os.Getenv(ENV_LOG_SWITCH)
	if logSwitch != "true" {
		return containers
	}
	// 如果该 service 没有 log 则不添加
	if len(cr.Spec.Product.Service[name].Instance.Logs) == 0 {
		log.Info(fmt.Sprintf("service: %s no logs path", name))
		return containers
	}

	var logPathEnv strings.Builder

	switch LogType {
	case "promtail":
		for _, path := range cr.Spec.Product.Service[name].Instance.Logs {
			logPathEnv.WriteString(fmt.Sprintf("%v/%v,", CommonLogPath, path[strings.LastIndex(path, "/")+1:]))
		}
	case "filebeat":
		for _, path := range cr.Spec.Product.Service[name].Instance.Logs {
			logPathEnv.WriteString(fmt.Sprintf("\"%v/%s\",", CommonLogPath, path[strings.LastIndex(path, "/")+1:]))
		}
	}
	log.Info(fmt.Sprintf("LOG_PATH ENV: " + strings.Trim(logPathEnv.String(), ",")))
	memLimit := DefaultLogSidecarMemoryLimit
	cpuLimit := "500m"
	_, exist := os.LookupEnv("LOG_MEM_LIMIT")
	if exist {
		memLimit = os.Getenv("LOG_MEM_LIMIT")
	}
	_, exist = os.LookupEnv("LOG_CPU_LIMIT")
	if exist {
		cpuLimit = os.Getenv("LOG_CPU_LIMIT")
	}

	resource := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: apiresource.MustParse(memLimit),
			corev1.ResourceCPU:    apiresource.MustParse(cpuLimit),
		},
	}
	_, memExist := os.LookupEnv("LOG_MEM_REQUEST")
	_, cpuExist := os.LookupEnv("LOG_CPU_REQUEST")
	if memExist && cpuExist {
		resource.Requests = corev1.ResourceList{
			corev1.ResourceMemory: apiresource.MustParse(os.Getenv("LOG_MEM_REQUEST")),
			corev1.ResourceCPU:    apiresource.MustParse(os.Getenv("LOG_CPU_REQUEST")),
		}
	}
	log.Info(fmt.Sprintf("logContainer resource %+v", resource))
	logContainer := corev1.Container{
		Name:      LogType,
		Image:     LogImage,
		Args:      LogArgs,
		Resources: resource,
		Env: []corev1.EnvVar{
			{
				Name: "HOSTNAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "spec.nodeName",
					},
				},
			}, {
				Name: "HOST_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "status.hostIP",
					},
				},
			}, {
				Name: "NAMESPACE",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "metadata.namespace",
					},
				},
			}, {
				Name: "SERVICE_ACCOUNT_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "spec.serviceAccountName",
					},
				},
			}, {
				Name: "POD_UID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "metadata.uid",
					},
				},
			}, {
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "metadata.name",
					},
				},
			}, {
				Name: "POD_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "status.podIP",
					},
				},
			}, {
				Name:  "JOB",
				Value: name,
			}, {
				Name:  "PRODUCT",
				Value: cr.Spec.Product.ProductName,
			}, {
				Name:  "LOG_PATH",
				Value: strings.Trim(logPathEnv.String(), ","),
			}, {
				Name:  "LOG_SERVER_ADDRESS",
				Value: LogServerAddress,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "log-config",
				ReadOnly:  true,
				MountPath: LogConfigPath,
			},
			{
				Name:      LogVolume,
				ReadOnly:  true,
				MountPath: CommonLogPath,
				SubPath:   LogSubpath,
			},
		},
		ImagePullPolicy: corev1.PullIfNotPresent,
	}
	if LogType == "promtail" {
		logContainer.Ports = append(logContainer.Ports, corev1.ContainerPort{
			Name:          "log-metrics",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: PromtailPort,
		})
	}
	containers = append(containers, logContainer)

	return containers
}

func getContainerPorts(cr *molev1.Mole, name string) []corev1.ContainerPort {
	//portName := BuildPortName(name, MoleHttpPortName)
	defaultPorts := make([]corev1.ContainerPort, 0)
	for index, port := range cr.Spec.Product.Service[name].Instance.Deployment.Ports {
		defaultPorts = append(defaultPorts, corev1.ContainerPort{
			Name:          BuildPortName("port", index),
			Protocol:      "TCP",
			ContainerPort: int32(port),
		})
	}
	if cr.Spec.Product.Service[name].Instance.PrometheusPort != "" {
		promPort, err := strconv.Atoi(cr.Spec.Product.Service[name].Instance.PrometheusPort)
		if err != nil {
			return defaultPorts
		}

		defaultPorts = append(defaultPorts, corev1.ContainerPort{
			Name:          "metrics",
			Protocol:      "TCP",
			ContainerPort: int32(promPort),
		})
	}

	return defaultPorts
}

func getHostAlias(cr *molev1.Mole, name string) []corev1.HostAlias {
	env := cr.Spec.Product.Service[name].Instance.Environment
	if env == nil {
		return nil
	}
	hostAlias, exist := env[EnvHostAlias]
	if !exist {
		return nil
	}
	alias := strings.Split(hostAlias, ",")
	results := make([]corev1.HostAlias, 0, len(alias))
	for _, alia := range alias {
		ipAndHosts := strings.Split(alia, ":")
		r := corev1.HostAlias{
			IP:        ipAndHosts[0],
			Hostnames: ipAndHosts[1:],
		}
		results = append(results, r)
	}
	return results
}

func getDeploymentSpec(cr *molev1.Mole, annotations map[string]string, name string) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas:        getReplicas(cr, name),
		MinReadySeconds: 10,

		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:        BuildResourceName(MolePodName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				Labels:      getPodLabels(cr, name),
				Annotations: getPodAnnotations(cr, annotations, name),
			},
			Spec: corev1.PodSpec{
				NodeSelector:     getNodeSelectors(cr, name),
				Tolerations:      getTolerations(cr, name),
				Affinity:         getAffinities(cr, name),
				SecurityContext:  getSecurityContext(cr, name),
				Volumes:          getVolumes(cr, name),
				InitContainers:   getInitPlugin(cr, name),
				Containers:       getContainers(cr, name),
				ImagePullSecrets: getImagePullSecrets(cr),
				HostAliases:      getHostAlias(cr, name),
				//ServiceAccountName: MoleServiceAccountName,
				//RestartPolicy:   corev1.RestartPolicyAlways,
				//TerminationGracePeriodSeconds: getTerminationGracePeriod(cr, name),
			},
		},
		Strategy: appsv1.DeploymentStrategy{
			Type:          "RollingUpdate",
			RollingUpdate: getRollingUpdateStrategy(),
		},
	}
}

func MoleDeployment(cr *molev1.Mole, name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BuildResourceName(MoleDeploymentName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Labels:    getDeploymentLabels(cr, name),
			Namespace: cr.Namespace,
		},
		Spec: getDeploymentSpec(cr, nil, name),
	}
}

func MoleDeploymentReconciled(cr *molev1.Mole, currentState *appsv1.Deployment, name string) *appsv1.Deployment {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = getDeploymentLabels(cr, name)
	reconciled.Spec = getDeploymentSpec(cr, currentState.Spec.Template.Annotations, name)
	return reconciled
}

func MoleDeploymentSelector(cr *molev1.Mole, name string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      BuildResourceName(MoleDeploymentName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	}
}

func getImagePullSecrets(cr *molev1.Mole) []corev1.LocalObjectReference {
	return []corev1.LocalObjectReference{
		{Name: cr.Spec.Product.ImagePullSecret},
	}
}

func getPodLifeCycle() *corev1.Lifecycle {
	return &corev1.Lifecycle{
		PostStart: &corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
					"mkdir -p /mount/${HOSTNAME}/logs && ln -s /mount/${HOSTNAME}/logs logs",
				},
			},
		},
	}
}
func getResources(cr *molev1.Mole, name string) corev1.ResourceRequirements {
	limits := corev1.ResourceList{
		corev1.ResourceMemory: apiresource.MustParse(DefaultMemoryLimit),
		corev1.ResourceCPU:    apiresource.MustParse(DefaultCpuLimit),
	}
	requests := corev1.ResourceList{
		corev1.ResourceMemory: apiresource.MustParse(DefaultMemoryRequest),
		corev1.ResourceCPU:    apiresource.MustParse(DefaultCpuRequest),
	}
	resources := cr.Spec.Product.Service[name].Instance.Resources
	for r, l := range resources.Limits {
		if _, support := SupportResource[r]; support {
			limits[r] = l
		}
	}
	for r, q := range resources.Requests {
		if _, support := SupportResource[r]; support {
			requests[r] = q
		}
	}

	return corev1.ResourceRequirements{
		Requests: requests,
		Limits:   limits,
	}
}
