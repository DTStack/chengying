### schema specification

eg. schema.yml
```yaml
parent_product_name: DTinsight
product_name: dtstack easy-stack
product_name_display: dtstack easy-stack（product_name的显示名，可选）
product_version: 2.1-xinwang

service:
  dtlog:
    service_display: dtlog（service的显示名，可选）
    version: 1.8.1
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

dtlog config file .eg dtlog.ini
```ini
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

### All keys must match same as program variable name limit.

## 顶层字段：
- parent_product_name: 产品名称
- product_name：产品名称，不能重复，不能修改[必填]
- product_name_display：产品显示名称 [选填]
- product_version: 产品版本号，对外唯一可见版本号（升级/回滚）[必填]
- product_type: 产品包类型，默认为传统包 [选填]
- service: 服务字典列表，对应服务包xxx.zip[非空]

## 服务字段：
xxx@other_product.yyy：xxx是本产品内的唯一服务名，
`@`表示继承于，other-product表示其它产品名称，yyy表示服务名称
xxx@other_product.yyy[#other_product_version]@xxxx：
继承服务扩展，xxxx为继承服务属性，如xxxx=option表示非部署依赖必须，为可选继承服务;xxxx=bridge表示跨集群桥接依赖
可选继承服务的服务IP缺省为[127.0.0.1]
[#other_product_version]可省，默认依赖同集群已部署服务，k8s部署模式下依赖传统集群服务的必须带



## 服务内关键字段：
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
- instance.image：实例镜像,kubernetes产品包必填 [选填]
- instance.replica：副本数量, kubernetes产品包必填 [可选]
- instance.home_page：实例home页面 [可选]
- instance.ha_role_cmd：实例HA角色探测命令 [可选]
- instance.pseudo：是否是伪服务实例，默认false [可选]
- instance.max_replica：最大实例副本数量 [可选]
- instance.start_after_install：当所有节点安装完成后再启动[可选]
- instance.update_recreate：先删除老实例，再启动新实例[可选]
- instance.config_path：配置文件路径[可选]
- instance.environment：环境变量列表 [必填]
- instance.post_deploy：部署完成后运行命令[可选]
- instance.post_upgrade：部署完成后执行特定升级命令，如执行增量sql[可选]
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

## 页面提示`可配置字段`（前端暂不支持）
- type：类型（filepath，port，ip，number，string），用来校验输入信息
- desc：描述信息
- default：默认值

## 普通`可配置字段`
除上面类型字段外的字段

## 变量引用
- `${[service.]field}`：引用service服务config块中的变量field。service可以省略，表示引用本服务块内的变量
- `${@service[.instance-name]}`：引用service服务中instance-name实例ip列表。instance-name可以省略，表示所有实例IP列表

## 服务配置文件
用golang模板引用本服务内的配置字段，详见[go template](https://golang.org/pkg/text/template/)
内置函数Join，JoinHost，Joinx，JoinxHost，NodeCount，NodeIndex，NodeID，LastSegIP，AddOne，
Hostname，IPList，HostList, GetIpByNodeID，GetHostByNodeID

## 虚拟服务
instance块不配置或instance.use_cloud为true

## 启动顺序
根据服务depends_on确定启动顺序，根据健康检查结果逐个启动所有服务。服务内实例从上到下启动。

## 服务守护
服务进程退出，会被再次拉起（1分钟内3次机会）

## 服务安装路径
/opt/dtstack/product_name/service_name

## 服务继承
IP地址继承：如果是虚拟服务，IP地址继承到本服务
version继承：如果是虚拟服务且version未设置，version继承到本服务
config块继承：字段变量实际值继承到本服务，继承字段会覆盖本服务字段
可选服务继承：可选服务继承仅提供缺省的IP地址

## 工作目录
服务安装路径

## 命令行文件路径
相对于服务tar包

## 监控系统
EasyManager服务端默认装有prometheus监控系统+grafana

## 产品包打包方法
mero pkg-dir [old-pkg]
pkg-dir是产品包目录：必须包含schema文件（名称必须是`schema.yml`），各个服务目录（目录名为服务名）
old-pkg是老的产品完整包
```
mero pkg-dir
Product package create success: "EasyLog_2.3.tar"

mero pkg-dir old-pkg
Product patch create success: "EasyLog_2.2_2.3.patch"
```