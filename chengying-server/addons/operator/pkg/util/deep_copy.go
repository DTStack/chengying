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

import "strconv"

// copy desired's value to existing.
// if the field of desired is not exist in existing, add the filed and value
// if the field of desired is exist in existing, change the value
// if a field exist in exsiting but not in desired, not changed
func DeepCopy(desired, existing interface{}) (interface{}, error) {

	switch desired := desired.(type) {
	case map[string]interface{}:
		if existing == nil {
			return desired, nil
		}
		existing, ok := existing.(map[string]interface{})
		if !ok {
			return nil, &InvalidError{
				Fields: []string{},
				ErrMsg: "structure is different with k8s object",
			}
		}
		for k, v := range desired {
			ev := existing[k]
			rv, err := DeepCopy(v, ev)
			if err != nil {
				invalidError := err.(*InvalidError)
				invalidError.Fields = append(invalidError.Fields, k)
				return nil, err
			}
			existing[k] = rv
		}
		return existing, nil
	case []interface{}:
		if existing == nil {
			return desired, nil
		}
		existing, ok := existing.([]interface{})
		if !ok {
			return nil, &InvalidError{
				Fields: []string{},
				ErrMsg: "structure is different with k8s object",
			}
		}
		for i, v := range desired {
			if len(existing) >= i {
				ev := existing[i]
				rv, err := DeepCopy(v, ev)
				if err != nil {
					invalidError := err.(*InvalidError)
					invalidError.Fields = append(invalidError.Fields, strconv.Itoa(i))
					return nil, err
				}
				existing[i] = rv
			} else {
				existing = append(existing, v)
			}
		}
		return existing, nil
	case string, int64, bool, float64, nil:
		return desired, nil
	default:
		return nil, &InvalidError{
			Fields: []string{},
			ErrMsg: "the type is not support",
		}
	}
}
