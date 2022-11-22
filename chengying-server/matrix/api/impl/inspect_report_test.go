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

package impl

import (
	"dtstack.com/dtstack/easymatrix/matrix/harole"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"github.com/kataras/iris"
	iriscontext "github.com/kataras/iris/context"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/wangqi811/gomonkey/v2"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	parentProductName = "DTinsight"
)

func TestGetServiceStatus(t *testing.T) {
	Convey("TestGetServiceStatus", t, func() {
		patches := ApplyFunc(GetCurrentParentProduct, func(ctx iriscontext.Context) (string, error) {
			fmt.Println("in mock function")
			return parentProductName, nil
		})
		defer patches.Reset()
		patches.ApplyFunc(GetCurrentClusterId, func(ctx iriscontext.Context) (int, error) {
			fmt.Println("in mock function")
			return 1, nil
		})
		patches.ApplySingletonMethod(reflect.TypeOf(model.DeployProductList), "GetDeploySonProductName",
			func(_ *interface{}, _ string, _ int) ([]string, error) {
				return []string{"DTBase"}, nil
			})
		patches.ApplySingletonMethod(reflect.TypeOf(model.DeployInstanceList), "FindByProductNameAndClusterId", func(_ *interface{}, _ string, _ int) ([]model.InstanceInfo, error) {
			dest := []model.InstanceInfo{
				{
					ServiceName: "kafka",
					Status:      "running",
					HealthState: 1,
					Ip:          "172.16.100.184",
					Pid:         57,
					AgentId:     "b82bf2ee-8c86-4779-a009-c0608f8e77e1",
				},
				{
					ServiceName: "mysql",
					Status:      "running",
					HealthState: 1,
					Ip:          "172.16.101.142",
					Pid:         57,
					AgentId:     "41aa76c4-8291-4e51-94ab-d70a4831f16c",
				},
			}
			return dest, nil
		})
		patches.ApplyFunc(harole.RoleData, func(_ int, _ string) map[string]string {
			return map[string]string{}
		})
		app := iris.New()
		ctx := iriscontext.NewContext(app)
		request := httptest.NewRequest("GET", "http://localhost", nil)
		responseRecorder := httptest.NewRecorder()
		ctx = app.ContextPool.Acquire(responseRecorder, request)
		result := GetServiceStatus(ctx)
		t.Logf("result: %s", result)
		So(result, ShouldNotBeNil)
	})
}
