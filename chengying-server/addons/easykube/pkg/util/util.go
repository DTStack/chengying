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

import (
	"fmt"
)

func MultiSizeConvert(size1, size2 int64) (string, string) {
	sizeUnits := [...]string{"B", "KB", "MB", "GB", "TB"}
	f1 := float32(size1)
	f2 := float32(size2)
	for _, v := range sizeUnits {
		if f1 < 1024 && f2 < 1024 {
			return fmt.Sprintf("%.2f"+v, f1), fmt.Sprintf("%.2f"+v, f2)
		} else {
			f1 = f1 / 1024
			f2 = f2 / 1024
		}
	}
	return fmt.Sprintf("%.2f"+sizeUnits[len(sizeUnits)-1], f1), fmt.Sprintf("%.2f"+sizeUnits[len(sizeUnits)-1], f1)
}

func MapConvert(m map[interface{}]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range m {
		switch valueType := v.(type) {
		case map[interface{}]interface{}:
			result[fmt.Sprint(k)] = MapConvert(valueType)
		default:
			result[fmt.Sprint(k)] = v
		}
	}
	return result
}
