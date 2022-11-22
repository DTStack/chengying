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

type upgradeType struct {
	Upgrade  EnumValueType
	Rollback EnumValueType
}

var UpgradeType = upgradeType{
	Upgrade: EnumValueType{
		Code: 0,
		Desc: "升级",
	},
	Rollback: EnumValueType{
		Code: 1,
		Desc: "回滚",
	},
}
