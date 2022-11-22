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

package client

import (
	"context"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/view/response"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	sigclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetWorkLoad(ctx context.Context, c *base.Client) (*response.WorkLoadResponse, error) {
	// CornJobs WorkLoad
	cronJobs := &batchv1beta1.CronJobList{}
	opt := sigclient.InNamespace(metav1.NamespaceAll)
	err := c.Lists(ctx, cronJobs, opt)
	if err != nil {
		return nil, err
	}
	workload := &response.WorkLoadResponse{}
	workload.CornJobs.Capacity = len(cronJobs.Items)
	workload.CornJobs.Load = getCornJobsLoadNum(cronJobs)

	// Jobs WorkLoad
	jobs := &batchv1.JobList{}
	err = c.Lists(ctx, jobs, opt)
	if err != nil {
		return nil, err
	}
	workload.Jobs.Capacity = len(jobs.Items)
	workload.Jobs.Load = getJobsLoadNum(jobs)

	// Pods WorkLoad
	pods := &corev1.PodList{}
	err = c.Lists(ctx, pods, opt)
	if err != nil {
		return nil, err
	}
	workload.Pods.Capacity = len(pods.Items)
	workload.Pods.Load = getPodsLoadNum(pods)

	// DaemonSets WorkLoad
	daemonSets := &appsv1.DaemonSetList{}
	err = c.Lists(ctx, daemonSets, opt)
	if err != nil {
		return nil, err
	}
	workload.DaemonSets.Capacity = len(daemonSets.Items)
	workload.DaemonSets.Load = getDaemonSetsLoadNum(daemonSets)

	// Deployment WorkLoad
	deployments := &appsv1.DeploymentList{}
	err = c.Lists(ctx, deployments, opt)
	if err != nil {
		return nil, err
	}

	workload.Deployments.Capacity = len(deployments.Items)
	workload.Deployments.Load = getDeploymentsLoadNum(deployments)

	// ReplicaSet WorkLoad
	replicaSets := &appsv1.ReplicaSetList{}
	err = c.Lists(ctx, replicaSets, opt)
	if err != nil {
		return nil, err
	}
	workload.ReplicaSets.Capacity = len(replicaSets.Items)
	workload.ReplicaSets.Load = getReplicaSetsLoadNum(replicaSets)

	// ReplicationController WorkLoad
	replicationControllers := &corev1.ReplicationControllerList{}
	err = c.Lists(ctx, replicationControllers, opt)
	if err != nil {
		return nil, err
	}
	workload.ReplicationControllers.Capacity = len(replicationControllers.Items)
	workload.ReplicationControllers.Load = getReplicationControllersLoadNum(replicationControllers)

	return workload, nil

}

// 获取 CornJob load 数量
func getCornJobsLoadNum(list *batchv1beta1.CronJobList) int {
	num := 0
	for _, cronJob := range list.Items {
		if cronJob.Spec.Suspend != nil && !(*cronJob.Spec.Suspend) {
			num++
		}
	}
	return num
}

// 获取 Job load 数量
func getJobsLoadNum(list *batchv1.JobList) int {
	num := 0
	for _, job := range list.Items {
		jobStatus := false
		for _, condition := range job.Status.Conditions {
			if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
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
func getPodsLoadNum(list *corev1.PodList) int {
	num := 0
	for _, pod := range list.Items {
		podStatus := false
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
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
func getDaemonSetsLoadNum(list *appsv1.DaemonSetList) int {
	num := 0
	for _, daemonSet := range list.Items {
		if daemonSet.Status.NumberReady == daemonSet.Status.CurrentNumberScheduled {
			num++
		}
	}
	return num
}

// 获取 Deployment Load 数量
func getDeploymentsLoadNum(list *appsv1.DeploymentList) int {
	num := 0
	for _, deployment := range list.Items {
		if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
			num++
		}
	}
	return num
}

// 获取 ReplicaSets Load 数量
func getReplicaSetsLoadNum(list *appsv1.ReplicaSetList) int {
	num := 0
	for _, replicaSet := range list.Items {
		if replicaSet.Status.Replicas == replicaSet.Status.ReadyReplicas {
			num++
		}
	}
	return num
}

// 获取 ReplicationControllerList Load 数量
func getReplicationControllersLoadNum(list *corev1.ReplicationControllerList) int {
	num := 0
	for _, replicationController := range list.Items {
		if replicationController.Status.Replicas == replicationController.Status.ReadyReplicas {
			num++
		}
	}
	return num
}
