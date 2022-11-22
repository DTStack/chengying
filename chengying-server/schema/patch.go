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

package schema

const (
	PATCH_FILE = "patch"
)

type DiffService struct {
	ServiceName  string
	NewFiles     []string
	DiffFiles    []string
	DeletedFiles []string
}

type Patch struct {
	ProductName       string
	NewProductVersion string
	OldProductVersion string
	NewServices       []string
	DeletedServices   []string
	DiffServices      []*DiffService
}

func (p Patch) IsDeletedService(name string) bool {
	for _, deleteSvcName := range p.DeletedServices {
		if deleteSvcName == name {
			return true
		}
	}
	return false
}

func (p Patch) IsDiffService(name string) bool {
	for _, diffSvc := range p.DiffServices {
		if diffSvc.ServiceName == name {
			return true
		}
	}
	return false
}

func (p Patch) IsChangedService(name string) bool {
	return p.IsDeletedService(name) || p.IsDiffService(name)
}
