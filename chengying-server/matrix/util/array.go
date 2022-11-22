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

package util

import "strings"

func ArrayUniqueStr(arr []string) []string {
	var list []string

	keys := make(map[string]bool)
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ArrayUniqueNotNullStr(arr []string) []string {
	var list []string

	keys := make(map[string]bool)
	for _, entry := range arr {
		entry = strings.Replace(entry, " ", "", -1)
		if len(entry) == 0 {
			continue
		}
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func RemoveDuplicateStrElement(elements []string) []string {
	result := make([]string, 0, len(elements))
	temp := map[string]struct{}{}
	for _, item := range elements {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

//求交集
func IntersectionString(slice1 []string, slice2 []string) []string {
	m := make(map[string]struct{})
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v] = struct{}{}
	}
	for _, v := range slice2 {
		if _, ok := m[v]; ok {
			nn = append(nn, v)
		}
	}
	return nn
}

//求差集 slice1 - slice2并集
func DifferenceString(slice1, slice2 []string) []string {
	m := make(map[string]struct{})
	nn := make([]string, 0)
	inter := IntersectionString(slice1, slice2)
	for _, v := range inter {
		m[v] = struct{}{}
	}

	for _, value := range slice1 {
		if _, ok := m[value]; !ok {
			nn = append(nn, value)
		}
	}
	return nn
}

func IndexOfString(slice []string, target string) int {
	for idx, s := range slice {
		if s == target {
			return idx
		}
	}
	return -1
}

// compare no ordered slice
func DiffNoOrderStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	if len(diff) == 0 {
		return true
	}
	return false
}
