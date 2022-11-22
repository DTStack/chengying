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
 @Description: 操作类型枚举
*/

type operationType struct {
	ProductDeploy     EnumValueType
	PatchUpdate       EnumValueType
	ProductStart      EnumValueType
	SvcStart          EnumValueType
	SvcRollingRestart EnumValueType
	HostInit          EnumValueType
	OpenKerberos      EnumValueType
	CloseKerberos     EnumValueType
	Backup            EnumValueType
	Rollback          EnumValueType
}

func (c operationType) List() (enumValues []EnumValueType) {
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		enumValues = append(enumValues, field.Interface().(EnumValueType))
	}
	return enumValues
}

func (c operationType) GetByCode(code int) (*EnumValueType, error) {
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

var OperationType = operationType{
	ProductDeploy: EnumValueType{
		Code: 1,
		Desc: "产品包部署",
	},
	PatchUpdate: EnumValueType{
		Code: 2,
		Desc: "补丁操作",
	},
	ProductStart: EnumValueType{
		Code: 3,
		Desc: "产品包启动",
	},

	SvcStart: EnumValueType{
		Code: 4,
		Desc: "服务启动",
	},

	SvcRollingRestart: EnumValueType{
		Code: 5,
		Desc: "服务滚动重启",
	},
	HostInit: EnumValueType{
		Code: 6,
		Desc: "主机初始化",
	},
	OpenKerberos: EnumValueType{
		Code: 7,
		Desc: "Kerberos开启",
	},
	CloseKerberos: EnumValueType{
		Code: 8,
		Desc: "Kerberos关闭",
	},
	Backup: EnumValueType{
		Code: 9,
		Desc: "备份数据库",
	},
	Rollback: EnumValueType{
		Code: 10,
		Desc: "回滚数据库",
	},
}
