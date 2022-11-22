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

func MergeAnnotations(requested map[string]string, existing map[string]string) map[string]string {
	//if existing == nil {
	//	return requested
	//}
	//
	//for k, v := range requested {
	//	existing[k] = v
	//}
	return existing
}

func BuildResourceName(resourceType, parentProductName, productName, serviceName string) string {
	return fmt.Sprintf("%v-%v-%v-%v", resourceType, ConvertDNSRuleName(parentProductName), ConvertDNSRuleName(productName), ConvertDNSRuleName(serviceName))
}

func BuildResourceLabel(parentProductName, productName, serviceName string) string {
	return fmt.Sprintf("%v-%v-%v", ConvertDNSRuleName(parentProductName), ConvertDNSRuleName(productName), ConvertDNSRuleName(serviceName))
}

func BuildPortName(serviceName string, index int) string {
	return fmt.Sprintf("%v-%v", ConvertDNSRuleName(serviceName), index)
}

func ConvertDNSRuleName(s string) string {
	s = strings.Replace(s, "_", "", -1)
	s = strings.ToLower(s)
	return s
}
