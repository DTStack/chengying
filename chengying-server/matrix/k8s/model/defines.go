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
	"fmt"
	"strings"
)

var (
	SECRET_VERSION  = "v1"
	SECRET_KIND     = "Secret"
	SECRET_GROUP    = ""
	SECRET_RESOURCE = "secrets"

	CONFIGMAP_KIND     = "ConfigMap"
	CONFIGMAP_GROUP    = ""
	CONFIGMAP_VERSION  = "v1"
	CONFIGMAP_RESOURCE = "configmaps"

	MOLE_VERSION  = "v1"
	MOLE_KIND     = "Mole"
	MOLE_GROUP    = "operator.dtstack.com"
	MOLE_RESOURCE = "moles"

	NAMESPACE_VERSION  = "v1"
	NAMESPACE_GROUP    = ""
	NAMESPACE_KIND     = "Namespace"
	NAMESPACELIST_KIND = "NamespaceList"

	NAMESPACE_PREFIX = "dtstack-"

	SupportResources = map[string]struct{}{
		"cpu":    {},
		"memory": {},
	}
)

func BuildResourceName(resourceType, parentProductName, productName, serviceName string) string {
	return fmt.Sprintf("%v-%v-%v-%v", resourceType, ConvertDNSRuleName(parentProductName), ConvertDNSRuleName(productName), ConvertDNSRuleName(serviceName))
}
func ConvertDNSRuleName(s string) string {
	s = strings.Replace(s, "_", "", -1)
	s = strings.ToLower(s)
	return s
}

func BuildResourceNameWithNamespace(resourceType, parentProductName, productName, serviceName, namespace string) string {
	return fmt.Sprintf("%v-%v-%v-%v.%v", resourceType, ConvertDNSRuleName(parentProductName), ConvertDNSRuleName(productName), ConvertDNSRuleName(serviceName), namespace)
}

func BuildWorkloadServiceName(productName, serviceName, partName, stepName, namespace string) string {
	return fmt.Sprintf("%v-%v-%v-%v.%v", ConvertDNSRuleName(productName), ConvertDNSRuleName(serviceName), ConvertDNSRuleName(partName), ConvertDNSRuleName(stepName), namespace)
}

func BuildWorkloadPodName(productName, serviceName, partName, stsnumber string) string {
	return fmt.Sprintf("%v-%v-%v-%v", ConvertDNSRuleName(productName), ConvertDNSRuleName(serviceName), ConvertDNSRuleName(partName), stsnumber)
}
