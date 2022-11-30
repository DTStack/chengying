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

package main

import (
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	commonlog "dtstack.com/dtstack/easymatrix/go-common/log"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/api"
	"dtstack.com/dtstack/easymatrix/matrix/api/impl"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"dtstack.com/dtstack/easymatrix/matrix/discover"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt"
	"dtstack.com/dtstack/easymatrix/matrix/grafana"
	"dtstack.com/dtstack/easymatrix/matrix/harole"
	"dtstack.com/dtstack/easymatrix/matrix/health"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/monitor"
	"dtstack.com/dtstack/easymatrix/matrix/task"
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"net"
	"os"
	"strconv"
	"time"
)

type DatabaseConfig struct {
	User     string `config:"user" validate:"required"`
	Password string `config:"password" validate:"required"`
	Host     string `config:"host" validate:"required"`
	Port     int    `config:"port"`
	DbName   string `config:"dbname" validate:"required"`
}

type LogConfig struct {
	Dir        string `config:"dir" validate:"required"`
	MaxSize    int    `config:"max-logger-size"`
	MaxBackups int    `config:"max-logger-backups"`
	MaxAge     int    `config:"days-to-keep"`
}

type ApiConfig struct {
	Host      string `config:"host"`
	Port      int    `config:"port" validate:"required"`
	Restrict  bool   `config:"restrict-api-check"`
	StaticUrl string `config:"static-url"`
}

type AgentConfig struct {
	Host        string `config:"host"`
	InstallPath string `config:"install-path"`
}

type DeployConfig struct {
	InstallPath string `config:"install-path"`
}

type VerifyConfig struct {
	Identity bool `config:"identity" validate:"required"`
	Code     bool `config:"code"`
}

type PrometheusConfig struct {
	NodeExporterPort int `config:"node-exporter-port"`
}

type Config struct {
	MysqlDb    DatabaseConfig `config:"mysqldb" validate:"required"`
	Log        LogConfig      `config:"log" validate:"required"`
	Api        ApiConfig      `config:"api" validate:"required"`
	Agent      AgentConfig    `config:"agent" validate:"required"`
	Deploy     DeployConfig   `config:"deploy" validate:"required"`
	GrafanaUrl string         `config:"grafana-url" validate:"required"`
	Verify     VerifyConfig
	Prometheus PrometheusConfig `config:"prometheus"`
	Swagger    bool             `config:"swagger"`
}

func resetDeployStatus() error {
	var err error
	var query string

	db := model.USE_MYSQL_DB()

	query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET `status`=?  WHERE `status`=?"
	if _, err = db.Exec(query, model.PRODUCT_STATUS_DEPLOY_FAIL, model.PRODUCT_STATUS_DEPLOYING); err != nil {
		return err
	}
	query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET `status`=?  WHERE `status`=?"
	if _, err = db.Exec(query, model.PRODUCT_STATUS_UNDEPLOY_FAIL, model.PRODUCT_STATUS_UNDEPLOYING); err != nil {
		return err
	}
	query = "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?  WHERE `status`=?"
	if _, err = db.Exec(query, model.PRODUCT_STATUS_DEPLOY_FAIL, model.PRODUCT_STATUS_DEPLOYING); err != nil {
		return err
	}
	query = "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?  WHERE `status`=?"
	if _, err = db.Exec(query, model.PRODUCT_STATUS_UNDEPLOY_FAIL, model.PRODUCT_STATUS_UNDEPLOYING); err != nil {
		return err
	}
	query = "UPDATE " + model.DeployInstanceList.TableName + " SET `status`=?, status_message='reset installing', update_time=NOW() WHERE `status`=?"
	if _, err = db.Exec(query, model.INSTANCE_STATUS_INSTALL_FAIL, model.INSTANCE_STATUS_INSTALLING); err != nil {
		return err
	}
	query = "UPDATE " + model.DeployInstanceList.TableName + " SET `status`=?, status_message='reset uninstalling', update_time=NOW() WHERE `status`=?"
	if _, err = db.Exec(query, model.INSTANCE_STATUS_UNINSTALL_FAIL, model.INSTANCE_STATUS_UNINSTALLING); err != nil {
		return err
	}

	return nil
}

func ParseConfig(configFile string) error {
	configContent, err := yaml.NewConfigWithFile(configFile, ucfg.PathSep("."))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found!")
		}
		return err
	}

	config := Config{}
	if err := configContent.Unpack(&config); err != nil {
		return err
	}

	clog := &config.Log
	if err := log.ConfigureLogger(clog.Dir, clog.MaxSize, clog.MaxBackups, clog.MaxAge); err != nil {
		return err
	} else {
		fmt.Printf("Saving logs at %s\n", clog.Dir)
	}
	//the other componet use the common package log
	commonlog.Config("EM")

	db := &config.MysqlDb
	if err := model.ConfigureMysqlDatabase(db.Host, db.Port, db.User, db.Password, db.DbName); err != nil {
		return err
	}
	if err := resetDeployStatus(); err != nil {
		return err
	}

	agentConfig := &config.Agent
	apiconf := &config.Api
	host.InitAgentInstall(agentConfig.Host, apiconf.StaticUrl, agentConfig.InstallPath)
	agent.InitAgentClient(agentConfig.Host)

	deployConfig := &config.Deploy
	base.ConfigureDeployInstallPath(deployConfig.InstallPath)

	grafana.InitGrafanaClient(config.GrafanaUrl)
	grafana.InitDashboard()
	apibase.InitSchemaVerify(config.Verify.Identity)
	impl.VerifyCode = config.Verify.Code
	impl.VerifyIdentity = config.Verify.Identity

	prometheusConfig := &config.Prometheus
	discover.NodeExporterPort = prometheusConfig.NodeExporterPort

	// 系统配置初始化
	cache.InitSysConfig()
	encrypt.InitPlatformEncrypt()

	if err := base.ConfigureApiServer(apiconf.Host, apiconf.Port, apiconf.StaticUrl, &api.ApiV2Schema, apiconf.Restrict, config.Swagger); err != nil {
		return err
	} else {
		fmt.Printf("Running API service at %v\n", net.JoinHostPort(apiconf.Host, strconv.Itoa(apiconf.Port)))
	}

	//启动角色状态线程
	go harole.StartRoleRunner()
	go grafana.StartMonitorAlertRecover()

	//启动策略
	//go strategy.SyncStrategies()
	//go strategy.ScheduleCache.Refresh()
	//go strategy.StrategyTaskMap.WaitForTask()

	//启动监控
	go monitor.StartMonitot()

	//初始化定时任务
	task.ServiceTask.Initialize()

	//启动多脚本健康检查
	c := health.Config{ReloadInterval: time.Minute * 5}
	healthCheck := health.NewHealthCheck(c)
	healthCheck.Run()

	return nil
}
