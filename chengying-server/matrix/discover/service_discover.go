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

package discover

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
	"text/template"

	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
)

var serviceLock = sync.Mutex{}

const (
	SERVICE_SD_FILE = "/prometheus/service_sd_file.yml"

	SERVICE_SD_TPL = `{{range $_ := .}}- targets: ['{{.Ip}}:{{.PrometheusPort}}']
  labels:
    product_name: {{.ProductName}}
    product_version: {{.ProductVersion}}
    cluster_name: {{.ClusterName}}
    service_name: {{.ServiceName}}
    service_version: {{.ServiceVersion}}
    group: {{.Group}}
{{end}}`
)

func FlushServiceDiscover() {
	serviceLock.Lock()
	defer serviceLock.Unlock()

	var instanceInfo []struct {
		ClusterName string `db:"cluster_name"`
		model.InstanceAndProductInfo
	}

	query := fmt.Sprintf("SELECT %s.*, %s.name AS cluster_name, product_name, product_name_display, product_version FROM %s LEFT JOIN %s ON pid=%s.id LEFT JOIN %s ON cluster_id=%s.id WHERE prometheus_port>0 AND %s.`status` IN (?, ?, ?, ?, ?) AND product_name is not NULL ",
		model.DeployInstanceList.TableName,
		model.DeployClusterList.TableName,
		model.DeployInstanceList.TableName,
		model.DeployProductList.TableName,
		model.DeployProductList.TableName,
		model.DeployClusterList.TableName,
		model.DeployClusterList.TableName,
		model.DeployInstanceList.TableName,
	)
	if err := model.USE_MYSQL_DB().Select(&instanceInfo, query,
		model.INSTANCE_STATUS_RUNNING,
		model.INSTANCE_STATUS_RUN_FAIL,
		model.INSTANCE_STATUS_STOPPED,
		model.INSTANCE_STATUS_STOPPING,
		model.INSTANCE_STATUS_STOP_FAIL); err != nil {
		log.Errorf("%v", err)
		return
	}

	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("service_discover").Option("missingkey=error").Parse(SERVICE_SD_TPL))
	if err := tpl.Execute(buf, instanceInfo); err != nil {
		log.Errorf("%v", err)
		return
	}

	log.Infof("Flush Service Discover: %v", string(buf.Bytes()))

	if err := ioutil.WriteFile(SERVICE_SD_FILE, buf.Bytes(), 0755); err != nil {
		log.Errorf("%v", err)
		return
	}
}
