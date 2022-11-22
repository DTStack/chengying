---
title: 产品包制作
---

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
# 产品包制作



> 首先，在制作产品包之前，先看看我们要准备哪些内容。

1.必须包含schema文件（名称必须是`schema.yml`）

2.各个服务目录（目录名为服务名）

3.[mero](https://dtstack-download.oss-cn-hangzhou.aliyuncs.com/chengying/mero)可执行文件



# schema配置介绍



> schema文件是一个标准的yaml文件



```yaml
parent_product_name: DTinsight
product_name: dtstack easy-stack
product_name_display: dtstack easy-stack（product_name的显示名，可选）
product_version: 2.1-xinwang
service:
  dtlog:
    service_display: dtlog（service的显示名，可选）
    version: 1.8.1
    orchestration:
      affinity:
        - web
      anti_affinity:
        - mysql
    instance:
      pseudo: true #伪服务
      ha_role_cmd: echo master #HA command
      #install path /opt/dtstack/dtstack easy-stack/dtlog
      cmd: bin/dtlog --port ${SLB.dtlog_front_port} --config ${config_path} --ca ${es.ca_file} --uic ${dtuic.jjww} #服务tar包内路径
      config_paths:
        - ${config_path}
      environment:
        MYSQL_ADDRESS: ${mysql.esip}:${es.bizdb}
        ES_ADDR: ${@es} #es service first ip
      healthcheck:
        shell: curl http://${dtlog_ip_port}/xxx/healthcheck
        #period: 30s #default 60s
        start_period: 10s #default 10s
        timeout: 10s #default 10s
        retries: 1 #default 1
      max_replica: 3
      start_after_install: false
      post_deploy: chown 0644 dtlog && zkcreate node xxx --ip ${@es}
      post_undeploy: rm -rf /var/data/dtlog
      logs:
      - log/abc.log
      - /var/log/dtlog/*.log
    group: default
    depends_on:
    - es
    - dtuic
    config:
      log_port: 8888
      config_path: dtlog.ini
      dtlog_ip_port: ${@dtlog}:8080 #self service's node ip
      es_ips: ${@es} #es service ip lists
      self_ip: ${@dtlog}
      
      
      
      
  es:
    version: 1.2
    instance:
      use_cloud: true  # same as instance == null
      cmd: start.sh
    config:
      ca_file:
        type: filepath
        desc: CA-File for EasyLog
        default: etc/es.ca.key
      bizdb: guyan${mysql.dsn}
      
      
      
  mysql:
    version: 5.7
    instance:
      cmd: start.sh
      max_replica: 1
      prometheus_port: ${prometheus_port}
      update_recreate: true
    config:
      username:
        type: string
        desc: Database Username for EasyLog
        default: root
      password:
        type: string
        desc: Database Password for EasyLog
        default: dtstack
      dsn: /dbname?encode=utf8
      prometheus_port: 9104
      esip: ${@es}
      
      
      
      
  SLB:
    version: 8.7
    #no instance, usually for cloud or other service/product
    config:
      dtlog_front_port:
        type: port
        desc: Uic port for EasyLog
        default: ${dtlog.log_port}
      dtlog_back_port: 8879
  dtuic@other_product.uic:
    #other-product uic's instance field not inherit
    version: 1.0
    #no instance, usually for cloud or other service/product
    depends_on:
    - es
    - SLB
    config:
      username:
        type: string
        desc: Username for UIC
        default: root
      password:
        type: string
        desc: Password for UIC
        default: dtstack
```

*配置项映射schema变量*

```
es-server: http://{{.Join .es_ips ","}}/aaa
arg: -c {{.config_path}} -self {{.dtlog_ip_port}} -host {{.Hostname .es_ips}}
node-count: {{.NodeCount .self_ip}}
node-index: {{.NodeIndex .self_ip}}
node-id: {{.NodeID .self_ip}} ##分配后不变
broker-id: {{.LastSegIP .self_ip}}  ##IP最后一段作为broker-id
add-one: {{.AddOne 100}}
JoinHost: http://{{.JoinHost .self_ip ","}}
JoinHost_: http://{{.JoinHost "self_ip" ","}} ##使用字符串，不推荐
JoinxHost: http://{{.JoinxHost .self_ip ","}}
Hostname: {{.Hostname .self_ip}}
IPList: {{range $v := .IPList .self_ip}}{{$v}},{{end}}
HostList: {{range $v := .HostList .self_ip}}{{$v}},{{end}}
HostList_: {{range $v := .HostList "self_ip"}}{{$v}},{{end}} ##使用字符串，不推荐
GetIpByNodeID: {{.GetIpByNodeID .NodeID .self_ip}}
```



# 配置介绍

## 顶层字段

- parent_product_name: 产品名称
- product_name：产品名称，不能重复，不能修改[必填]
- product_name_display：产品显示名称 [选填]
- product_version: 产品版本号，对外唯一可见版本号（升级/回滚）[必填]
- service: 服务字典列表，对应服务包xxx.zip[非空]

## 服务字段

xxx@other_product.yyy：

​        xxx是本产品内的唯一服务名,`@`表示继承于，other-product表示其它产品名称，yyy表示服务名称 xxx@other_product.yyy@xxxx： 继承服务扩展，xxxx为继承服务属性，如xxxx=option表示非部署依赖必须，为可选继承服务 可选继承服务的服务IP缺省为[127.0.0.1].

## 服务内关键字段

- service_display：服务的显示名 [选填]
- version：服务版本号（判断是否需要升级/回滚）[必填]
- depends_on：依赖服务列表 [可选]
- relatives: 部署血缘，需要部署在相同的节点上 [可选]
- instance：实例数组（从上到下顺序执行），没有表示`虚拟服务`，只提供配置信息给其它服务，如RDS，ELB eg. [可选]
- instance.empty_car：空车模式|空包模式，表示仅部署实例管控逻辑，不产生实际的文件部署，可用于现有系统被EM接管[可选]
- instance.install_path：指定实例的安装路径，缺省时使用默认安装路径[可选]
- instance.use_cloud：是否使用云上资源，跟`虚拟服务`意义相同[可选]
- instance.run_user：指定实例启动用户，没有表示跟父进程一致 [选填]
- instance.data_dir：数据目录配置，配置run_user时会按照run_user做目录chown run_user:run_user[选填]
- instance.cmd：实例启动命令行,必须hang住（空格间隔参数，非shell，忽略std/stderr输出） [必填]
- instance.home_page：实例home页面 [可选]
- instance.ha_role_cmd：实例HA角色探测命令 [可选]
- instance.pseudo：是否是伪服务实例，默认false [可选]
- instance.max_replica：最大实例/副本数量 [可选]  自动部署时作为该服务的副本数
- instance.start_after_install：当所有节点安装完成后再启动[可选]
- instance.update_recreate：先删除老实例，再启动新实例[可选]
- instance.config_path：配置文件路径[可选]
- instance.environment：环境变量列表 [必填]
- instance.post_deploy：部署完成后运行命令[可选]
- instance.post_upgrade：部署完成后执行特定命令，如升级时增量sql，4.1.3引入[可选]
- instance.post_undeploy：卸载部署完成后运行命令[可选]
- instance.prometheus_port：prometheus监听端口（如果有exporter，必须包含在产品包中且在cmd脚本中启动）[可选]
- instance.logs: 日志文件（支持*号匹配）[可选]
- instance.healthcheck：健康检查 [可选]
- instance.healthcheck.shell：return code=0表示成功[必填]
- instance.healthcheck.period：检查周期（>=1s），如"300s", "5m", "1.5h"[可选，默认60s]
- instance.healthcheck.start_period：开始检查时间，如"10s"[可选]
- instance.healthcheck.timeout：超时时间，如"10s"[可选]
- instance.healthcheck.retries：重试次数，如果都失败，才认为健康检查失败[可选，默认1]
- group: 组名称，默认default[可选]
- config: 可配置字段块，用来被本服务配置文件或所有服务instance,config字段引用
- orchestration.affinity：数组   自动编排角色亲和性  [可选]  但是没有该字段的话，该服务将无法参与自动编排
- orchestration.anti_affinity：数组   自动编排角色反亲和性  [可选]

## 变量引用

- `${[service.]field}`：引用service服务config块中的变量field。service可以省略，表示引用本服务块内的变量
- `${@service[.instance-name]}`：引用service服务中instance-name实例ip列表。instance-name可以省略，表示所有实例IP列表

## 服务配置文件

用golang模板引用本服务内的配置字段，详见[go template](https://golang.org/pkg/text/template/) 内置函数Join，JoinHost，Joinx，JoinxHost，NodeCount，NodeIndex，NodeID，LastSegIP，AddOne， Hostname，IPList，HostList, GetIpByNodeID，GetHostByNodeID

## 虚拟服务

instance块不配置或instance.use_cloud为true

## 启动顺序

根据服务depends_on确定启动顺序，根据健康检查结果逐个启动所有服务。服务内实例从上到下启动。

## 服务守护

服务进程退出，会被再次拉起（1分钟内3次机会）

## 服务安装路径

/opt/dtstack/product_name/service_name

## 服务继承

IP地址继承：如果是虚拟服务，IP地址继承到本服务 version继承：如果是虚拟服务且version未设置，version继承到本服务 config块继承：字段变量实际值继承到本服务，继承字段会覆盖本服务字段 可选服务继承：可选服务继承仅提供缺省的IP地址

## 工作目录

服务安装路径

## 命令行文件路径

相对于服务tar包

## 监控系统

EasyManager服务端默认装有prometheus监控系统+grafana

## 产品包打包方法

mero pkg-dir [old-pkg] pkg-dir是产品包目录：

必须包含schema文件（名称必须是`schema.yml`）

各个服务目录（目录名为服务名） old-pkg是老的产品完整包

```
mero pkg-dir
Product package create success: "product_name_product_version.tar"
```





# Jvm进程产品包的制作



> 在制作产品包之前，我们先回顾一下一个jar的启停维护方式。



**启动**

*java -jar jarFileName.jar或者使用nohup进行后台启动*

**维护**

*通常情况下，只能通过查看日志进行分析，或借助第三方工具。*



**对于传统模式下存在的问题**

* 手动启停
* 进程健康无法得到保障
* 无版本记录
* 等等





### 制作产品包前思考

1.这个包需要有几个服务

2.有哪些是伪服务

3.每个服务的启动是否有强依赖

4.有哪些是需要初始化



#### 示例思考

*假设我们已经部署好了Jvm所依赖的一些基础服务，所以我们只需要启动这个Jar。*

*我们需要部署启动一个Jvm进程，是一个有状态的进程，所以，他不是一个伪服务。*





### 制作产品包前准备

1.需要准备好每个服务的初始化脚本(如果需要)

2.需要准备好每个服务的启动脚本(伪服务不需要)

3.需要准备好每个服务的健康检查脚本(伪服务不需要)





#### 示例准备



**初始化脚本**

````bash
#!/bin/bash

mkdir -p {logs,run}
````



**启动脚本**

```bash
#!/bin/bash

source /etc/profile

exec $JAVA_HOME/bin/java -jar test.jar
```



**健康检查脚本**

````bash
#!/bin/bash

#usage ./health.sh ip port1 port2 port2
#parameter 1 is host
#parameters after 1 are list of ports to check
#all ports are ok echo 1, else echo 0 


if [ $# -lt 2 ] ; then
  echo "USAGE: $0 host [ports]" 
  exit
fi
 
address=$1
ret=0
shift
i=$@


for i in "$@"; do
  if command -v nc >/dev/null 2>&1; then
 # echo "exists nc"
   #  echo  $address $i
     nc -w 10  $address $i  < /dev/null >/dev/null 2>&1 
     
   #  echo status ======= $?
     if [ $? -eq 0 ] ; then
        ret=$(( $ret + 0 ))
     else
        ret=$(( $ret + 1 ))
     fi
  else
    ret=2
  fi
done


exit $ret

````



### 示例schema配置

```yaml
parent_product_name: DTinsight
product_name: Test_pkg
product_version: v1.0.1
 
service:

  Back:
    version: v1.0.1
    instance:
      cmd: ./bin/base.sh              ##服务的启动脚本
      post_deploy: ./bin/post_deploy.sh      ##服务在启动前的初始化操作
      environment:             ##环境变量导入，可以在脚本里直接引用
        JAVA_OPTS: ${java_opts}      
      config_paths:            ##服务的配置文件所在目录
      - conf/application.properties
      healthcheck:
        shell: ./bin/health.sh ${@Back} 8090     ##健康检查脚本配置
        period: 20s              ##健康检查的间隔时间
        retries: 3               ##健康检查的重试次数
      logs:
      - logs/*.log               ##配置服务的日志地址
      prometheus_port: 9515      ##配置服务启动的export的监听端口
    group: test          ##服务的分组，跟页面的服务展示有关系
    config:       ##配置映射,配置文件需要在config_paths中定义
      java_opts: "-Xms2048m -Xmx2048m"
      service_port: 8090
```



## mero打包

**打包准备**

````bash
[root@rel43-em mero]# ls -l
total 4
drwxr-xr-x 5 root root   40 May 30 00:09 Back
-rw-r--r-- 1 root root 1097 May 30 00:08 schema.yml
[root@rel43-em mero]# ls -R
.:
Back  schema.yml

./Back:
bin  conf  lib

./Back/bin:
base.sh  health.sh  post_deploy.sh

./Back/conf:
application.properties

./Back/lib:
test.jar
[root@rel43-em mero]# 
````



**执行mero**

````bas
[root@rel43-em mero]# ls -l
total 4
drwxr-xr-x 5 root root   40 May 30 00:09 Back
-rw-r--r-- 1 root root 1097 May 30 00:08 schema.yml
[root@rel43-em mero]# mero .
Product package create success: "Test_pkg_v1.0.1.tar"
[root@rel43-em mero]# ls -l
total 12
drwxr-xr-x 5 root root   40 May 30 00:09 Back
-rw-r--r-- 1 root root 1097 May 30 00:08 schema.yml
-rw-r--r-- 1 root root 5120 May 30 00:26 Test_pkg_v1.0.1.tar
[root@rel43-em mero]# 
````


至此，一个产品包制作完成









