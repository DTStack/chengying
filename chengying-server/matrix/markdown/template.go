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

package markdown

const DeployInfoPage = `
# 一、登录方式



# 二、部署架构
{deployArchTable}


# 三、部署版本
{deployVersionTable}

# 四、访问地址
|WEB访问           |URL              |账号密码        |
|--------|--------------------|---------------------|
|数栈地址 |
|EM地址  |

# 五、部署清单
{deployListTable}

`
const DeployArchTable = `
|主机IP      |主机名称       |CPU     |内存    |数据磁盘    |系统盘    |操作系统    |账号/密码    |
|-----------|--------------|--------|-------|-----------|---------|-----------|------------|
{params}`

const DeployArchParam = `|{hostIP}   |{hostName}   |{CPU}      |{MEM}    |{dataDisk}  |{systemDisk}  |{operatorSys}  |   |
`

const DeployVersionTable = `
|交付产品      |交付版本       |部署人     |部署日期    |
|-------------|-------------|----------|-----------|
{params}`

const DeployVersionParam = `|{product}   |{version}   |     	 |  		 |
`

const DeployListTable = `
|组件      |组件版本       |服务        |服务版本       |{ip/host}    |
|---------|--------------|-----------|--------------|--------------|
{params}`

const DeployListParam = `|{product}   |{version}   |{service}      |{service_version}  	 |{ips/hosts}	|
`

const DeploySelectFlag = "✅"
