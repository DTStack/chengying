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

package enums

import (
	"fmt"
	"reflect"
)

/*
 @Author: zhijian
 @Date: 2021/5/28 10:03
 @Description: shell 操作类型枚举
*/

type shellType struct {
	Install EnumValueType
	Start   EnumValueType
	Exec    EnumValueType
}

func (c shellType) List() (enumValues []EnumValueType) {
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		enumValues = append(enumValues, field.Interface().(EnumValueType))
	}
	return enumValues
}

func (c shellType) GetByCode(code int) (*EnumValueType, error) {
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		enum := v.Field(i).Interface().(EnumValueType)
		if enum.Code == code {
			return &enum, nil
		}
	}
	return nil, fmt.Errorf("not found by code %d", code)
}

var ShellType = shellType{
	Install: EnumValueType{
		Code: 1,
		Desc: "服务安装",
	},
	Start: EnumValueType{
		Code: 2,
		Desc: "服务启动",
	},
	Exec: EnumValueType{
		Code: 3,
		Desc: "执行脚本",
	},
}
