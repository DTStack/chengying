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
	"database/sql"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/schema"
	"fmt"
	"github.com/kataras/iris"
	iriscontext "github.com/kataras/iris/context"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/wangqi811/gomonkey/v2"
	"net/http/httptest"
	"reflect"
	"testing"
)

const DefaultSchema = "{\"ParentProductName\":\"DTinsight\",\"ProductName\":\"DTBase\",\"ProductNameDisplay\":\"\",\"ProductVersion\":\"2.1.7_haiyang\",\"Service\":{\"clean_history\":{\"ServiceDisplay\":\"\",\"Version\":\"1.0.0\",\"Instance\":{\"ConfigPaths\":[\"backup_clean.sh\"],\"RunUser\":\"\",\"Cmd\":\"./waiting.sh\",\"PostDeploy\":\"./post_deploy.sh\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"\"},\"Group\":\"default\",\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"},\"kafka\":{\"ServiceDisplay\":\"streamtest\",\"Version\":\"1.1.1_1\",\"Instance\":{\"ConfigPaths\":[\"config/server.properties\"],\"Logs\":[\"logs/*.log\"],\"HealthCheck\":{\"Shell\":\"./health.sh ${@kafka} 9092\",\"Period\":\"20s\",\"StartPeriod\":\"\",\"Timeout\":\"\",\"Retries\":1},\"RunUser\":\"\",\"Cmd\":\"./start_kafka.sh\",\"PostDeploy\":\"./post_deploy.sh\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"9525\"},\"Group\":\"default\",\"DependsOn\":[\"zookeeper\"],\"Config\":{\"kafka_data\":{\"Default\":\"/data/kafka/logs\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"/data/kafka/logs\"},\"kafka_ip\":{\"Default\":\"${@kafka}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${@kafka}\"},\"zk_ip\":{\"Default\":\"${@zookeeper}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${@zookeeper}\"}},\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"},\"kdcserver\":{\"ServiceDisplay\":\"\",\"Version\":\"\",\"Group\":\"default\",\"BaseProduct\":\"Kerberos\",\"BaseProductVersion\":\"\",\"BaseService\":\"kdcserver\",\"BaseParsed\":false,\"BaseAtrribute\":\"optional\"},\"kerberos_pkg\":{\"ServiceDisplay\":\"\",\"Version\":\"\",\"Group\":\"default\",\"Config\":{\"default_realm\":{\"Default\":\"DTSTACK.COM\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"DTSTACK.COM\"}},\"BaseProduct\":\"Kerberos\",\"BaseProductVersion\":\"\",\"BaseService\":\"kerberos_pkg\",\"BaseParsed\":false,\"BaseAtrribute\":\"optional\"},\"mysql\":{\"ServiceDisplay\":\"\",\"Version\":\"5.7.33-1\",\"Instance\":{\"ConfigPaths\":[\"my.cnf\"],\"Logs\":[\"/data/my3306/log/alert.log\"],\"HealthCheck\":{\"Shell\":\"./health.sh ${@mysql} ${user} ${password}\",\"Period\":\"20s\",\"StartPeriod\":\"\",\"Timeout\":\"\",\"Retries\":2},\"RunUser\":\"\",\"Cmd\":\"./bin/mysql.sh\",\"HARoleCmd\":\"./show_mysql_state.sh ${@mysql} ${user} ${password} 2\\u003e/dev/null\",\"PostDeploy\":\"./post_deploy.sh ${user} ${password} ${repl_user} ${repl_password} ${check_num}\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"9104\",\"MaxReplica\":\"1\",\"UpdateRecreate\":true},\"Group\":\"mysql\",\"Config\":{\"check_num\":{\"Default\":\"3\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"3\"},\"db\":{\"Default\":\"metastore\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"metastore\"},\"encryption_mysql_password\":{\"Default\":\"XdVyOxV50S7MA8qMfYJJFwzHpmz0qTeiObElsfOYVcQsqjUGZEZyRW+lI76vZ/BIuHb66ac3YmAvFWJD8NOPkcWcJGdMJJK+XgVs9iIVe8NP+7hmXLR7Ah5mgI5lYqbIKQ8D82TiA5D5WuU203wUqdUARGt8s8gwLXGimoIOJjg=\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"XdVyOxV50S7MA8qMfYJJFwzHpmz0qTeiObElsfOYVcQsqjUGZEZyRW+lI76vZ/BIuHb66ac3YmAvFWJD8NOPkcWcJGdMJJK+XgVs9iIVe8NP+7hmXLR7Ah5mgI5lYqbIKQ8D82TiA5D5WuU203wUqdUARGt8s8gwLXGimoIOJjg=\"},\"password\":{\"Default\":\"DT@Stack#123\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"DT@Stack#123\"},\"port\":{\"Default\":\"3306\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"3306\"},\"repl_password\":{\"Default\":\"drpEco@123\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"drpEco@123\"},\"repl_user\":{\"Default\":\"repl\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"repl\"},\"tengine_ip\":{\"Default\":\"127.0.0.1\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"127.0.0.1\"},\"user\":{\"Default\":\"drpeco\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"drpeco\"}},\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"},\"mysql_slave\":{\"ServiceDisplay\":\"\",\"Version\":\"5.7.33-1\",\"Instance\":{\"ConfigPaths\":[\"crond\",\"my.cnf\"],\"Logs\":[\"/data/my3306/log/alert.log\"],\"HealthCheck\":{\"Shell\":\"./health.sh ${@mysql_slave} ${mysql_user} ${mysql_password}\",\"Period\":\"20s\",\"StartPeriod\":\"\",\"Timeout\":\"\",\"Retries\":2},\"RunUser\":\"\",\"Cmd\":\"./bin/mysql.sh\",\"HARoleCmd\":\"./show_mysql_state.sh ${@mysql_slave} ${mysql_user} ${mysql_password} 2\\u003e/dev/null\",\"PostDeploy\":\"./post_deploy.sh ${@mysql} ${mysql.repl_user} ${mysql.repl_password} ${mysql.user} ${mysql.password} ${mysql.check_num}\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"9104\",\"MaxReplica\":\"1\",\"UpdateRecreate\":true},\"Group\":\"mysql\",\"DependsOn\":[\"mysql\"],\"Config\":{\"mysql_host\":{\"Default\":\"${@mysql}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${@mysql}\"},\"mysql_password\":{\"Default\":\"${mysql.password}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${mysql.password}\"},\"mysql_user\":{\"Default\":\"${mysql.user}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${mysql.user}\"},\"node_list\":{\"Default\":\"node001,node002,node003\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"node001,node002,node003\"}},\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"},\"prometheus\":{\"ServiceDisplay\":\"\",\"Version\":\"2.23.0-5\",\"Instance\":{\"ConfigPaths\":[\"prometheus.yml\"],\"Logs\":[\"logs/*.log\"],\"HealthCheck\":{\"Shell\":\"./health.sh ${@prometheus} 9090\",\"Period\":\"20s\",\"StartPeriod\":\"\",\"Timeout\":\"\",\"Retries\":3},\"RunUser\":\"\",\"Cmd\":\"./start_prometheus.sh\",\"PostDeploy\":\"\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"\",\"MaxReplica\":\"1\"},\"Group\":\"default\",\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"},\"pushgateway\":{\"ServiceDisplay\":\"\",\"Version\":\"0.4.0-4\",\"Instance\":{\"HealthCheck\":{\"Shell\":\"./health.sh ${@pushgateway} 9091\",\"Period\":\"20s\",\"StartPeriod\":\"\",\"Timeout\":\"\",\"Retries\":3},\"RunUser\":\"\",\"Cmd\":\"./start_pushgateway.sh\",\"PostDeploy\":\"\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"\",\"MaxReplica\":\"1\"},\"Group\":\"default\",\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"},\"redis\":{\"ServiceDisplay\":\"\",\"Version\":\"3.2.13-8\",\"Instance\":{\"ConfigPaths\":[\"conf/redis.conf\",\"conf/sentinel.conf\",\"bin/show_redis_role.sh\",\"bin/start_redis.sh\"],\"Logs\":[\"logs/*.log\"],\"HealthCheck\":{\"Shell\":\"./bin/health.sh ${@redis} 16379\",\"Period\":\"20s\",\"StartPeriod\":\"\",\"Timeout\":\"\",\"Retries\":2},\"RunUser\":\"\",\"Cmd\":\"bin/start_redis.sh\",\"HARoleCmd\":\"./bin/show_redis_role.sh 2\\u003e/dev/null\",\"PostDeploy\":\"./post_deploy.sh\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"9121\"},\"Group\":\"default\",\"Config\":{\"db\":{\"Default\":\"1\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"1\"},\"encryption_password\":{\"Default\":\"XdVyOxV50S7MA8qMfYJJFwzHpmz0qTeiObElsfOYVcQsqjUGZEZyRW+lI76vZ/BIuHb66ac3YmAvFWJD8NOPkcWcJGdMJJK+XgVs9iIVe8NP+7hmXLR7Ah5mgI5lYqbIKQ8D82TiA5D5WuU203wUqdUARGt8s8gwLXGimoIOJjg=\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"XdVyOxV50S7MA8qMfYJJFwzHpmz0qTeiObElsfOYVcQsqjUGZEZyRW+lI76vZ/BIuHb66ac3YmAvFWJD8NOPkcWcJGdMJJK+XgVs9iIVe8NP+7hmXLR7Ah5mgI5lYqbIKQ8D82TiA5D5WuU203wUqdUARGt8s8gwLXGimoIOJjg=\"},\"host\":{\"Default\":\"${@redis}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${@redis}\"},\"is_redis_standalone\":{\"Default\":\"false\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"false\"},\"password\":{\"Default\":\"DT@Stack#123\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"DT@Stack#123\"},\"rd_port\":{\"Default\":\"16379\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"16379\"},\"redis_port\":{\"Default\":\"26379\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"26379\"}},\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"},\"zookeeper\":{\"ServiceDisplay\":\"\",\"Version\":\"3.4.14-7_3\",\"Instance\":{\"ConfigPaths\":[\"conf/zoo.cfg\",\"conf/java.env\",\"conf/jaas.conf\",\"bin/on.sh\"],\"Logs\":[\"logs/zookeeper.out\"],\"HealthCheck\":{\"Shell\":\"bin/healthcheck.sh ${@zookeeper} 2181\",\"Period\":\"20s\",\"StartPeriod\":\"\",\"Timeout\":\"\",\"Retries\":3},\"RunUser\":\"\",\"Cmd\":\"./bin/zkServer.sh\",\"HARoleCmd\":\"bin/show_zookeeper_role.sh 2\\u003e/dev/null\",\"PostDeploy\":\"bin/post_deploy.sh ${@zookeeper}\",\"PostUpGrade\":\"\",\"PostUndeploy\":\"\",\"PrometheusPort\":\"9505\",\"Switch\":{\"Kerberos\":{\"Config\":\"kerberos_on\",\"Desc\":\"Kerberos switch\",\"IsOn\":false,\"OnScript\":\"bin/on.sh\",\"OffScript\":\"\",\"PostOnScript\":{\"Type\":\"restart\",\"Value\":\"zookeeper\",\"Desc\":\"\"},\"PostOffScript\":{\"Type\":\"restart\",\"Value\":\"zookeeper\",\"Desc\":\"\"},\"Extention\":{\"Type\":\"download\",\"Value\":\"/data/kerberos/zookeeper.keytab\",\"Desc\":\"\"}}}},\"Group\":\"default\",\"Config\":{\"kdcserver_ip\":{\"Default\":\"${@kdcserver}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${@kdcserver}\"},\"kerberos_on\":{\"Default\":\"0\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"0\"},\"realms\":{\"Default\":\"${kerberos_pkg.default_realm}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${kerberos_pkg.default_realm}\"},\"zk_ip\":{\"Default\":\"${@zookeeper}\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"${@zookeeper}\"},\"zk_port\":{\"Default\":\"2181\",\"Desc\":\"internal\",\"Type\":\"internal\",\"Value\":\"2181\"}},\"BaseProduct\":\"\",\"BaseProductVersion\":\"\",\"BaseService\":\"\",\"BaseParsed\":false,\"BaseAtrribute\":\"\"}}}"

func TestGetSwitchDetail(t *testing.T) {
	Convey("TestGetSwitchDetail", t, func() {
		patches := ApplyFunc(log.Errorf, func(format string, args ...interface{}) {})
		defer patches.Reset()

		Convey("TestGetSwitchDetailSuccess", func() {
			patches.ApplySingletonMethod(reflect.TypeOf(model.SwitchRecord), "GetRecordById",
				func(_ *interface{}, _ int) (*model.SwitchRecordInfo, error) {
					ret := &model.SwitchRecordInfo{
						Id:            1,
						ClusterId:     1,
						Name:          "Kerberos",
						ProductName:   "DTBase",
						ServiceName:   "zookeeper",
						Status:        "SUCCESS",
						Progress:      100,
						StatusMessage: "SUCCESS",
						Type:          "on",
					}
					return ret, nil
				})
			app := iris.New()
			ctx := iriscontext.NewContext(app)
			request := httptest.NewRequest("GET", "http://localhost?record_id=1", nil)
			responseRecorder := httptest.NewRecorder()
			ctx = app.ContextPool.Acquire(responseRecorder, request)
			result := GetSwitchDetail(ctx)
			r := result.(switchDetail)
			So(responseRecorder.Code, ShouldEqual, 200)
			So(r.Progress, ShouldEqual, 100)
			So(r.SwitchType, ShouldEqual, "on")
		})
	})
}

func TestCheckSwitchRecord(t *testing.T) {
	Convey("TestCheckSwitchRecord", t, func() {
		Convey("TestCheckSwitchRecordSuccess", func() {
			patches := ApplyFunc(GetCurrentClusterId, func(ctx iriscontext.Context) (int, error) {
				return 1, nil
			})
			defer patches.Reset()
			patches.ApplySingletonMethod(reflect.TypeOf(model.SwitchRecord), "GetCurrentSwitchRecord",
				func(_ *interface{}, _ int, _, _, _ string) (*model.SwitchRecordInfo, error) {
					ret := &model.SwitchRecordInfo{
						Id: 1,
					}
					return ret, nil
				})
			app := iris.New()
			ctx := iriscontext.NewContext(app)
			request := httptest.NewRequest("GET", "http://localhost?name=Kerberos", nil)
			response := httptest.NewRecorder()
			ctx = app.ContextPool.Acquire(response, request)
			ctx.Params().Set("product_name", "DTBase")
			ctx.Params().Set("service_name", "zookeeper")
			result := CheckSwitchRecord(ctx)
			r, ok := result.(map[string]interface{})
			So(ok, ShouldBeTrue)
			So(r["record_id"], ShouldEqual, 1)
		})

		Convey("TestCheckSwitchRecordFail", func() {
			patches := ApplyFunc(log.Errorf, func(format string, args ...interface{}) {})
			defer patches.Reset()
			patches.ApplyFunc(GetCurrentClusterId, func(ctx iriscontext.Context) (int, error) {
				return 1, nil
			})
			patches.ApplySingletonMethod(reflect.TypeOf(model.SwitchRecord), "GetCurrentSwitchRecord",
				func(_ *interface{}, _ int, _, _, _ string) (*model.SwitchRecordInfo, error) {
					return nil, fmt.Errorf("error")
				})
			app := iris.New()
			ctx := iriscontext.NewContext(app)
			request := httptest.NewRequest("GET", "http://localhost?name=Kerberos", nil)
			response := httptest.NewRecorder()
			ctx = app.ContextPool.Acquire(response, request)
			ctx.Params().Set("product_name", "DTBase")
			ctx.Params().Set("service_name", "zookeeper")
			result := CheckSwitchRecord(ctx)
			r, ok := result.(error)
			So(ok, ShouldBeTrue)
			So(r, ShouldNotEqual, sql.ErrNoRows)
		})

		Convey("TestCheckSwitchRecordFail2", func() {
			patches := ApplyFunc(log.Errorf, func(format string, args ...interface{}) {})
			defer patches.Reset()
			patches.ApplyFunc(GetCurrentClusterId, func(ctx iriscontext.Context) (int, error) {
				return 1, nil
			})
			patches.ApplySingletonMethod(reflect.TypeOf(model.SwitchRecord), "GetCurrentSwitchRecord",
				func(_ *interface{}, _ int, _, _, _ string) (*model.SwitchRecordInfo, error) {
					return nil, sql.ErrNoRows
				})
			app := iris.New()
			ctx := iriscontext.NewContext(app)
			request := httptest.NewRequest("GET", "http://localhost?name=Kerberos", nil)
			response := httptest.NewRecorder()
			ctx = app.ContextPool.Acquire(response, request)
			ctx.Params().Set("product_name", "DTBase")
			ctx.Params().Set("service_name", "zookeeper")
			result := CheckSwitchRecord(ctx)
			r, ok := result.(map[string]interface{})
			So(ok, ShouldBeTrue)
			So(r["record_id"], ShouldEqual, 0)
		})
	})
}

func TestOperateSwitch(t *testing.T) {
	Convey("TestOperateSwitch", t, func() {
		Convey("TestOperateSwitchOn", func() {
			switchRecord := &switcher{
				id:          1,
				clusterId:   1,
				productName: "DTBase",
				serviceName: "zookeeper",
				switchName:  "Kerberos",
				switchType:  "on",
				schema:      nil,
				ip:          []string{"172.16.100.184", "172.16.8.85", "172.16.8.87"},
			}
			patches := ApplyFunc(GetCurrentClusterId, func(ctx iriscontext.Context) (int, error) {
				return 1, nil
			})
			defer patches.Reset()
			patches.ApplyFunc(log.Errorf, func(format string, args ...interface{}) {})
			patches.ApplySingletonMethod(reflect.TypeOf(model.DeployProductList), "GetByProductNameAndVersion",
				func(_ *interface{}, _, _ string) (*model.DeployProductListInfo, error) {
					return &model.DeployProductListInfo{
						Schema: []byte(DefaultSchema),
					}, nil
				})
			patches.ApplyFunc(NewSwitcher, func(clusterId int, productName, serviceName, switchName, switchType string, schema *schema.SchemaConfig) (Switcher, error) {
				return switchRecord, nil
			})
			patches.ApplySingletonMethod(reflect.TypeOf(model.DeployClusterProductRel), "GetCurrentProductByProductNameClusterId",
				func(_ *interface{}, _ string, _ int) (*model.DeployProductListInfo, error) {
					return &model.DeployProductListInfo{
						ID: 198,
					}, nil
				})
			patches.ApplySingletonMethod(reflect.TypeOf(model.DeployInstanceList), "GetInstanceListByPidServiceName",
				func(_ *interface{}, _, _ int, _ string) ([]model.DeployInstanceInfo, error) {
					ret := []model.DeployInstanceInfo{{
						Ip: "172.16.100.184",
					}, {
						Ip: "172.16.8.85",
					}, {
						Ip: "172.16.8.87",
					}}
					return ret, nil
				})
			patches.ApplySingletonMethod(reflect.TypeOf(model.SwitchRecord), "NewSwitchRecord",
				func(_ *interface{}, _, _, _, _, _, _ string, _, _ int) (int64, error) {
					return int64(1), nil
				})
			patches.ApplySingletonMethod(reflect.TypeOf(switchRecord), "AddProgressLog",
				func(_, _string, _ int) error {
					return nil
				})
			patches.ApplyFunc(doSwitch, func(_ Switcher, _ string) interface{} {
				return nil
			})
			app := iris.New()
			ctx := iriscontext.NewContext(app)
			request := httptest.NewRequest("GET", "http://localhost", nil)
			response := httptest.NewRecorder()
			ctx = app.ContextPool.Acquire(response, request)
			ctx.Params().Set("product_name", "DTBase")
			ctx.Params().Set("service_name", "zookeeper")
			result := OperateSwitch(ctx)
			r, ok := result.(map[string]interface{})
			So(ok, ShouldBeTrue)
			So(r["record_id"], ShouldEqual, 1)
		})
	})
}
