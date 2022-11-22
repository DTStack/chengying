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

package base

import (
	"fmt"
)

const (
	NORMAL_EXIT = iota
	NETWORK_FAILURE
)

type SystemFailure struct {
	ExitCode int
	Reason   string
}

func SystemExitWithFailure(exitCode int, r interface{}, args ...interface{}) {
	if exitCode == NORMAL_EXIT {
		_SYSTEM_FAIL <- SystemFailure{NORMAL_EXIT, ""}
	} else if reason, ok := r.(string); ok {
		_SYSTEM_FAIL <- SystemFailure{exitCode, fmt.Sprintf(reason, args...)}
	} else if err, ok := r.(error); ok {
		_SYSTEM_FAIL <- SystemFailure{exitCode, err.Error()}
	} else if code, ok := r.(int); ok {
		_SYSTEM_FAIL <- SystemFailure{exitCode, fmt.Sprintf("Code: %d", code)}
	} else {
		_SYSTEM_FAIL <- SystemFailure{exitCode, "Unknown reason"}
	}
}
