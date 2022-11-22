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
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/wangqi811/gomonkey/v2"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

/*
 @Author: zhijian
 @Date: 2021/3/25 20:49
 @Description:
*/

// http://hwcld-03:88/api/v2/product/backup?clusterId=1
func TestGetBackupPackage(t *testing.T) {

	Convey("TestGetBackupPackage", t, func() {
		patches := ApplySingletonMethod(reflect.TypeOf(model.DeployClusterHostRel), "GetClusterHostRelList",
			func(_ *interface{}, _ int) ([]model.ClusterHostRel, error) {
				return []model.ClusterHostRel{
						{
							Id:         1,
							Sid:        "970d3196-1b56-440c-9a49-90e7a1edadef",
							ClusterId:  1,
							Roles:      "",
							UpdateTime: dbhelper.NullTime{},
							CreateTime: dbhelper.NullTime{},
							IsDeleted:  0,
						},
						{
							Id:         2,
							Sid:        "abd68c8c-5493-4cc3-9d00-687a8ef8b363",
							ClusterId:  1,
							Roles:      "",
							UpdateTime: dbhelper.NullTime{},
							CreateTime: dbhelper.NullTime{},
							IsDeleted:  0,
						},
						{
							Id:         3,
							Sid:        "8fd07e06-1e95-43c9-834b-22c536b6c600",
							ClusterId:  1,
							Roles:      "",
							UpdateTime: dbhelper.NullTime{},
							CreateTime: dbhelper.NullTime{},
							IsDeleted:  0,
						},
						{
							Id:         4,
							Sid:        "26eaaa4c-a58a-44f9-b7fa-20542011d0ba",
							ClusterId:  1,
							Roles:      "",
							UpdateTime: dbhelper.NullTime{},
							CreateTime: dbhelper.NullTime{},
							IsDeleted:  0,
						},
					},
					nil
			})
		defer patches.Reset()
		hostInfoBySidOutputs := []OutputCell{
			{Values: Params{nil, &model.HostInfo{Ip: "192.168.1.1"}}}, // 模拟函数的第1次输出
			{Values: Params{nil, &model.HostInfo{Ip: "192.168.1.2"}}}, // 模拟函数的第2次输出
			{Values: Params{nil, &model.HostInfo{Ip: "192.168.1.3"}}}, // 模拟函数的第3次输出
			{Values: Params{nil, &model.HostInfo{Ip: "192.168.1.4"}}}, // 模拟函数的第4次输出
			{Values: Params{nil, &model.HostInfo{Ip: "192.168.1.5"}}}, // 模拟函数的第5次输出

		}
		patches.ApplyMethodSeq(reflect.TypeOf(model.DeployHostList), "GetHostInfoBySid", hostInfoBySidOutputs)

		toExecCmd := []OutputCell{
			{Values: Params{`/opt/dtstack/DTBase/mysql_1610970097_4.0.2~`, nil}}, // 模拟函数的第1次输出
			{Values: Params{`/opt/dtstack/easymanager/mysql_1610970097_4.0.2~
/opt/dtstack/easymanager/mysql_1610970097_4.0.3~
/opt/dtstack/DTUic/dtuic_1610970097_4.0.0~`, nil}}, // 模拟函数的第2次输出
			{Values: Params{`/opt/dtstack/easymanager/mysql_1610970097_4.0.2~
/opt/dtstack/DTStream/StreamFront_1610970097_4.0.2~
/opt/dtstack/easymanager/mysql_1610970097_4.0.4~
/opt/dtstack/DTUic/dtuic_1610970097_4.0.0~`, nil}}, // 模拟函数的第3次输出
			{Values: Params{`/opt/dtstack/easymanager/mysql_1610970097_4.0.2~
/opt/dtstack/DTStream/StreamFront_1610970097_4.0.2~
/opt/dtstack/easymanager/mysql_1610970097_4.0.3~
/opt/dtstack/DTUic/dtuic_1610970097_4.0.0~`, nil}}, // 模拟函数的第4次输出
		}
		// /opt/dtstack/easymanager/mysql_1610970097_4.0.2~
		patches.ApplyMethodSeq(reflect.TypeOf(agent.AgentClient), "ToExecCmd", toExecCmd)
		patches.ApplyFunc(log.Infof, func(format string, args ...interface{}) {})
		app := iris.New()
		ctx := context.NewContext(app)
		request := httptest.NewRequest("GET", "http://hwcld-03:88/api/v2/product/backup?clusterId=1", nil)
		responseRecorder := httptest.NewRecorder()
		ctx = app.ContextPool.Acquire(responseRecorder, request)
		result := GetBackupPackage(ctx)
		t.Logf("result: %s", result)
		So(result, ShouldNotBeNil)
	})
}

//http://hwcld-03:88/api/v2/product/clean
func TestCleanBackupPackage(t *testing.T) {

	Convey("TestCleanBackupPackage", t, func() {
		patches := ApplySingletonMethod(reflect.TypeOf(agent.AgentClient), "ToExecCmd", func(_ *interface{}, _, _, _ string) (content string, err error) {
			return content, nil
		})
		defer patches.Reset()
		reader := strings.NewReader(`[
        {
            "product": "DTBase",
            "version": "4.0.3",
            "service": [
                {
                    "name": "mysql",
                    "host_info": [
                        {
                            "sid": "007605dd-45b2-4a52-94ee-a86b5ea723df",
                            "ip": "192.168.0.55"
                        },
                        {
                            "sid": "007605dd-45b2-4a52-94ee-a86b5ea723df",
                            "ip": "192.168.0.88"
                        }
                    ]
                }
            ]
        },
        {
            "product": "DTBase",
            "version": "4.0.2",
            "service": [
                {
                    "name": "mysql",
                    "host_info": [
                        {
                            "sid": "007605dd-45b2-4a52-94ee-a86b5ea723df",
                            "ip": "192.168.0.55"
                        }
                    ]
                },
                {
                    "name": "test",
                    "host_info": [
                        {
                            "sid": "4960dfc6-90bf-47db-9aa1-c783cac1149c",
                            "ip": "192.168.0.88"
                        }
                    ]
                }
            ]
        },
        {
            "product": "easymanager",
            "version": "4.0.4",
            "service": [
                {
                    "name": "mysql",
                    "host_info": [
                        {
                            "sid": "4960dfc6-90bf-47db-9aa1-c783cac1149c",
                            "ip": "192.168.0.88"
                        }
                    ]
                },
                {
                    "name": "dtuic",
                    "host_info": [
                        {
                            "sid": "4960dfc6-90bf-47db-9aa1-c783cac1149c",
                            "ip": "192.168.0.88"
                        }
                    ]
                }
            ]
        },
        {
            "product": "easymanager",
            "version": "4.0.2",
            "service": [
                {
                    "name": "mysql",
                    "host_info": [
                        {
                            "sid": "4960dfc6-90bf-47db-9aa1-c783cac1149c",
                            "ip": "192.168.0.88"
                        }
                    ]
                }
            ]
        }
    ]`)
		patches.ApplyFunc(log.Infof, func(format string, args ...interface{}) {})
		app := iris.New()
		ctx := context.NewContext(app)
		request := httptest.NewRequest("GET", "http://localhost", reader)
		responseRecorder := httptest.NewRecorder()
		ctx = app.ContextPool.Acquire(responseRecorder, request)
		result := CleanBackupPackage(ctx)
		So(result, ShouldBeNil)
	})

}
