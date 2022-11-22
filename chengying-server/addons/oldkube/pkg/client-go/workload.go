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

package client_go

import (
	"context"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	batch2 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetWorkLoad(client kubernetes.Interface) (WorkLoadResponse, error) {
	workload := WorkLoadResponse{}

	// CornJobs WorkLoad
	cronJobs, err := client.BatchV1beta1().CronJobs("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return workload, err
	}
	workload.CornJobs.Capacity = len(cronJobs.Items)
	workload.CornJobs.Load = getCornJobsLoadNum(cronJobs)

	// Jobs WorkLoad
	jobs, err := client.BatchV1().Jobs("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return workload, err
	}
	workload.Jobs.Capacity = len(jobs.Items)
	workload.Jobs.Load = getJobsLoadNum(jobs)

	// Pods WorkLoad
	pods, err := client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return workload, err
	}
	workload.Pods.Capacity = len(pods.Items)
	workload.Pods.Load = getPodsLoadNum(pods)

	// DaemonSets WorkLoad
	daemonSet, err := client.AppsV1().DaemonSets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return workload, err
	}
	workload.DaemonSets.Capacity = len(daemonSet.Items)
	workload.DaemonSets.Load = getDaemonSetsLoadNum(daemonSet)

	// Deployment WorkLoad
	deployment, err := client.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return workload, err
	}

	workload.Deployments.Capacity = len(deployment.Items)
	workload.Deployments.Load = getDeploymentsLoadNum(deployment)

	// ReplicaSet WorkLoad
	replicaSets, err := client.AppsV1().ReplicaSets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return workload, err
	}
	workload.ReplicaSets.Capacity = len(replicaSets.Items)
	workload.ReplicaSets.Load = getReplicaSetsLoadNum(replicaSets)

	// ReplicationController WorkLoad
	replicationControllers, err := client.CoreV1().ReplicationControllers("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return workload, err
	}
	workload.ReplicationControllers.Capacity = len(replicationControllers.Items)
	workload.ReplicationControllers.Load = getReplicationControllersLoadNum(replicationControllers)

	return workload, nil

}

// 获取 CornJob load 数量
func getCornJobsLoadNum(list *batch2.CronJobList) int {
	num := 0
	for _, cronJob := range list.Items {
		if cronJob.Spec.Suspend != nil && !(*cronJob.Spec.Suspend) {
			num++
		}
	}
	return num
}

// 获取 Job load 数量
func getJobsLoadNum(list *batch.JobList) int {
	num := 0
	for _, job := range list.Items {
		jobStatus := false
		for _, condition := range job.Status.Conditions {
			if condition.Type == batch.JobComplete && condition.Status == v1.ConditionTrue {
				jobStatus = true
				break
			}
		}
		if jobStatus {
			num++
		}
	}
	return num
}

// 获取 Pod load 数量
func getPodsLoadNum(list *v1.PodList) int {
	num := 0
	for _, pod := range list.Items {
		podStatus := false
		for _, condition := range pod.Status.Conditions {
			if condition.Type == v1.PodReady && condition.Status == v1.ConditionTrue {
				podStatus = true
				break
			}
		}
		if podStatus {
			num++
		}
	}
	return num
}

// 获取 DaemonSet Load 数量
func getDaemonSetsLoadNum(list *apps.DaemonSetList) int {
	num := 0
	for _, daemonSet := range list.Items {
		if daemonSet.Status.NumberReady == daemonSet.Status.CurrentNumberScheduled {
			num++
		}
	}
	return num
}

// 获取 Deployment Load 数量
func getDeploymentsLoadNum(list *apps.DeploymentList) int {
	num := 0
	for _, deployment := range list.Items {
		if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
			num++
		}
	}
	return num
}

// 获取 ReplicaSets Load 数量
func getReplicaSetsLoadNum(list *apps.ReplicaSetList) int {
	num := 0
	for _, replicaSet := range list.Items {
		if replicaSet.Status.Replicas == replicaSet.Status.ReadyReplicas {
			num++
		}
	}
	return num
}

// 获取 ReplicationControllerList Load 数量
func getReplicationControllersLoadNum(list *v1.ReplicationControllerList) int {
	num := 0
	for _, replicationController := range list.Items {
		if replicationController.Status.Replicas == replicationController.Status.ReadyReplicas {
			num++
		}
	}
	return num
}
