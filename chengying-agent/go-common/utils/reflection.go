/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"fmt"
	"reflect"
)

func getAsInterfaceValueAndType(data interface{}) (reflect.Type, reflect.Value, bool) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Struct {
		data = v.Interface()
	} else if t.Kind() == reflect.Ptr {
		data = v.Elem().Interface()
	} else if t.Kind() != reflect.Interface {
		return reflect.TypeOf(nil), reflect.ValueOf(nil), false
	}
	v = reflect.ValueOf(data)
	t = reflect.TypeOf(data)
	return t, v, true
}

func GetTagValues(data interface{}, tag string) []string {
	list := []string{}
	t, _, ok := getAsInterfaceValueAndType(data)
	if !ok {
		return list
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(tag)
		if tag != "" {
			list = append(list, tag)
		}
	}
	return list
}

func MapTagFromStruct(data interface{}, tag string) (map[string]interface{}, error) {
	t, v, ok := getAsInterfaceValueAndType(data)
	if !ok {
		return nil, fmt.Errorf("Not a struct")
	}
	result := map[string]interface{}{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(tag)
		if tag != "" {
			result[tag] = v.Field(i).Interface()
		}
	}
	return result, nil
}
