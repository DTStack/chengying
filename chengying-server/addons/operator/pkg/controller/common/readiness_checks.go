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

package common

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
)

const (
	ConditionStatusSuccess = "True"
)

func IsIngressReady(ingress *v1beta1.Ingress) bool {
	if ingress == nil {
		return false
	}

	return len(ingress.Status.LoadBalancer.Ingress) > 0
}

func IsDeploymentReady(deployment *appsv1.Deployment) (bool, error) {
	if deployment == nil {
		return false, nil
	}
	// A deployment has an array of conditions
	for _, condition := range deployment.Status.Conditions {
		// One failure condition exists, if this exists, return the Reason
		if condition.Type == appsv1.DeploymentReplicaFailure {
			return false, errors.New(condition.Reason)
			// A successful deployment will have the progressing condition type as true
		} else if condition.Type == appsv1.DeploymentProgressing && condition.Status != ConditionStatusSuccess {
			return false, nil
		}
	}

	return deployment.Status.ReadyReplicas == deployment.Status.Replicas, nil
}
