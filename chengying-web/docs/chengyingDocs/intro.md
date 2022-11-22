---
id: intro
title: 什么是Chengying
sidebar_position: 1
---

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

## 简介

承影是一个全自动化全生命周期的运维管家，提供大数据产品的一站式部署、运维、监控服务，其可实现产品部署、产品升级、版本回滚、扩缩节点、日志诊断、集群监控、实时告警等功能，致力于最大化节省运维成本，降低线上故障率与运维难度，为客户提供安全稳定的产品部署与监控。

## 系统架构

![q1](/img/intro/img.png)

产品部署之前需添加主机、及上传安装包，主机和安装包通过nginx进行服务转发，matinx是主程序服务，其中包含schema安装包文件解析，orchestrate组件依赖关系处理，instance exec执行引擎3部分。除此之外，promethus进行集群监控，以及通过grafana进行仪表盘展现，MySQL存储Easy Manager相关数据。Easy Manager所有的组件通过easyagent server服务与外部实体连接、调度。

![q1](/img/intro/img_1.png)

Easyagent Server通过标准的GPRC协议与sidecar进行双向通讯，通过GRPC全双工通讯，客户端网络环境不需要开监听端口，所有控制指令进行实时传输。Sidecar进行不同agent的服务管理，如kafaka、ES等服务，可实现对agent服务的cgroup分组以及tc控制，控制主机资源使用、上报监控资源使用情况，以及进行健康检查。同时Easyagent Server抽象出七大REST接口，安装、启动、停止、更新、配置修改、卸载、执行等与上层应用进行交互，使agent类别和功能可轻松无限扩展。

![q1](/img/intro/img_2.png)

## 主要特性

* 自动化部署

EasyManager通过规范化的步骤和参数约定制作出产品安装包，发布包中的Schema文件中配置了安装包中所有的服务，包含各服务的配置参数、健康检查参数，服务之间的依赖关系等。产品部署时可根据Schema中的相关配置实现一键全自动化部署。

* 界面化集群运维

Hadoop集群、大数据平台在日常运维中涉及到的节点扩容缩容、组件停止启动、服务滚动重启、服务参数修改、版本升级与回滚等多种运维操作，通过逻辑化、流程化的产品界面展现，方便运维人员操作和监控，提高运维效率。

* 仪表盘集群监控

通过集成开源的promethus和grafana，实现对集群、服务、节点的核心参数监控，并通过灵活形象的仪表盘进行数据展现。包含CPU占用率，RAM使用率、磁盘空间、I/O读写速率等核心参数进行监控，实时掌握集群、服务、节点的运行状态，降低运维故障率。同时，支持用户自建仪表盘及监控项，实现自定义监控项。

* 实时告警

支持实时监控集群中各组件服务的运行指标，如CPU、内存、磁盘、读写IO等，并支持短信、钉钉、邮件告警通道配置，集成多种第三方消息插件。当集群服务出现异常时，可触发告警条件，系统将及时通知接收人。

* 强扩展性

通过自研的Easyagent Server抽象出七大REST接口，安装、启动、停止、更新、配置修改、卸载、执行等与上层应用进行交互，可使agent类别和功能可轻松无限扩展。

* 安全稳定

数据安全、产品安全是大数据产品需要重点考虑的问题。Easy Manager在产品设计中过滤掉rm、drop等命令行，防止对数据库的误操作，通过更加安全的方式执行相关命令。同时提供服务的滚动重启、产品的断电重启，解决运维时服务不停止运行的场景并节省运维时间。


## 名词解释



产品	"产品"指用EM部署的大数据产品，如数栈（DTinsight）、云日志（EasyLog）等，作为EM最高级别的组织单位。
组件	"组件"一般和"组件安装包"等同，指产品下包含多个组件，亦产品下包含多个组件安装包，如数栈（DTinsight）产品下包含Hadoop、DTBase、DTCommon、DTBatch、DTStream等多个组件，代表不同的应用。
服务组	"服务组"值组件下的服务分类，如Hadoop组件下包含HDFS、Spark、Yarn、Flink、Default等服务组，进行服务区分。
服务	"服务"指服务组下的具体服务，如HDFS服务组下包含hdfs_datanode、hdfs_journalnode、hdfs_namenode、hdfs_zkfc等组件；Spark服务组下包含hivemetastore、spark_historyserver、thriftserver等组件。
主机分组	"主机分组"指对主机进行人为分组，当主机数量过多时可进行分组，方便管理。
主机	"主机"指服务器，包含物理机、虚拟机，指产品部署时需要的硬件资源，通常以主机IP或主机名称进行区分。