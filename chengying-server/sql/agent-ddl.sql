CREATE DATABASE IF NOT EXISTS `dtagent` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
USE dtagent;

CREATE TABLE IF NOT EXISTS `agent_list` (
    `id` char(36) NOT NULL COMMENT 'Agent ID (UUID)',
    `sidecar_id` char(36) NOT NULL COMMENT 'sidecar id',
    `type` tinyint(1) NOT NULL COMMENT 'agent类型',
    `name` char(32) NULL DEFAULT '' COMMENT 'agent 名称',
    `version` char(32) NOT NULL DEFAULT '' COMMENT 'agent版本',
    `is_uninstalled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已被卸载',
    `deploy_date` datetime DEFAULT NULL COMMENT 'agent部署时间',
    `auto_deployment` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是自动部署的',
    `last_update_date` datetime DEFAULT NULL COMMENT '最近更新时间',
    `auto_updated` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是自动升级的',
    PRIMARY KEY (`id`),
    KEY `uuid` (`sidecar_id`)
    ) ENGINE=InnoDB  COMMENT='Sidecar管控的Agent信息表';

CREATE TABLE IF NOT EXISTS `operation_history` (
    `seq` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '操作序列',
    `op_name` varchar(20) NOT NULL DEFAULT '' COMMENT '操作名称',
    `op_time` datetime NOT NULL COMMENT '操作时间',
    `target` char(36) NOT NULL DEFAULT '' COMMENT '目标id（sidecar id）',
    `send_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '操作状态',
    `op_result` tinyint(1) DEFAULT '1' COMMENT '操作执行状态',
    `op_return_msg` mediumblob DEFAULT NULL COMMENT '操作返回内容',
    `finish_time` datetime DEFAULT NULL COMMENT '操作执行结束时间',
    PRIMARY KEY (`seq`)
    ) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `progress_history` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `seq` int(11) unsigned NOT NULL COMMENT '对应操作序列号',
    `ts` datetime NOT NULL COMMENT '事件时间',
    `progress` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '进度百分比',
    `sidecar_id` char(36) NOT NULL DEFAULT '' COMMENT 'sidecar id',
    `agent_id` char(36) DEFAULT '' COMMENT 'agent id',
    `msg` varchar(100) DEFAULT '' COMMENT '附带信息',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `sidecar_list` (
    `id` char(36) NOT NULL COMMENT 'Sidecar ID (UUID)',
    `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT 'Sidecar状态',
    `disabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否被屏蔽',
    `name` varchar(32) DEFAULT '' COMMENT 'Sidecar备注名',
    `version` varchar(32) DEFAULT '' COMMENT 'Sidecar版本',
    `host` varchar(255) DEFAULT '' COMMENT '主机域名或者ip',
    `os_type` varchar(16) DEFAULT '' COMMENT '目标系统类型,linux，windows等',
    `is_ecs` tinyint(1) DEFAULT '0' COMMENT '是否是ECS',
    `os_platform` varchar(64) DEFAULT '' COMMENT 'os完整的名称',
    `os_version` varchar(64) DEFAULT '' COMMENT 'os版本号',
    `cpu_serial` varchar(64) DEFAULT '' COMMENT 'cpu型号',
    `cpu_cores` tinyint(8) DEFAULT '0' COMMENT 'cpu内核数',
    `mem_size` bigint(20) unsigned DEFAULT '0' COMMENT '内存容量',
    `swap_size` bigint(20) unsigned DEFAULT '0' COMMENT '交换空间容量',
    `deploy_date` datetime DEFAULT NULL COMMENT 'Sidecar部署时间',
    `auto_deployment` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是自动部署的',
    `last_update_date` datetime DEFAULT NULL COMMENT '最近更新时间',
    `auto_updated` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是自动升级的',
    `server_host` varchar(255) DEFAULT '' COMMENT 'api server ip',
    `server_port` int(11) DEFAULT 0 COMMENT 'api server port',
    `ssh_host` varchar(255) DEFAULT '' COMMENT '安装用的ssh主机域名或者ip',
    `ssh_user` varchar(60) DEFAULT '' COMMENT 'ssh用户名',
    `ssh_password` varchar(100) DEFAULT '' COMMENT 'ssh密码',
    `ssh_port` int(11) DEFAULT '22' COMMENT 'ssh端口',
    `cpu_usage` decimal(6,2) DEFAULT '-1' COMMENT 'cpu使用率',
    `mem_usage` bigint(20) DEFAULT '-1' COMMENT '物理内存使用',
    `swap_usage` bigint(20) DEFAULT '-1' COMMENT '交换空间使用',
    `load1` float DEFAULT '-1' COMMENT 'cpu load1',
    `uptime` double DEFAULT '-1' COMMENT '系统启动时间',
    `disk_usage` text DEFAULT NULL COMMENT '各个硬盘使用率',
    `net_usage` text DEFAULT NULL COMMENT '各个网卡统计',
    `local_ip` varchar(255) DEFAULT '' COMMENT '主机ip',
    PRIMARY KEY (`id`),
    KEY `uuid` (`id`)
    ) ENGINE=InnoDB  COMMENT='Sidecar客户端';

CREATE TABLE IF NOT EXISTS `deploy_callback` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto',
    `time` int(11) unsigned NOT NULL DEFAULT '0',
    `client_id` varchar(64) DEFAULT '' COMMENT '客户端的UUID',
    `install_type` varchar(32) DEFAULT '' COMMENT '安装类型 [sidecar 等等]',
    `install_res` varchar(32) DEFAULT '' COMMENT '安装标识信息[success,failed 等等]',
    `msg` varchar(255) DEFAULT '' COMMENT '安装结论信息',
    `request_url` varchar(2000) DEFAULT '' COMMENT '访问的回调原始请求',
    `ip` varchar(32) DEFAULT '' COMMENT 'ip地址',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB  COMMENT='一键安装部署回调数据表';

CREATE TABLE IF NOT EXISTS `deploy_host` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    `sid` varchar(255) NULL DEFAULT '' COMMENT '主机sid',
    `hostname` varchar(255) NULL DEFAULT '' COMMENT '主机名',
    `ip` varchar(255) NULL DEFAULT '' COMMENT '主机ip',
    `status` int(11) NULL DEFAULT '0' COMMENT '1:管控安装成功,-1:管控安装失败,2:script安装成功,-2:script安装失败,3:主机初始化成功,-3:主机初始化失败',
    `steps`  int(11) NULL DEFAULT '0' COMMENT '0:默认;1:管控安装成功;2:script wrapper安装成功;3:主机初始化成功',
    `errorMsg` varchar(1024) NULL DEFAULT '' COMMENT '错误信息',
    `group` varchar(255) DEFAULT 'default' COMMENT '组信息',
    `isDeleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
    `updated` datetime NOT NULL COMMENT 'updated',
    `created` datetime NOT NULL COMMENT 'created',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_product_list` (
    `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `deploy_uuid` varchar(36) NOT NULL COMMENT '部署uuid',
    `parent_product_name` varchar(255) NOT NULL COMMENT '产品名称',
    `product_name` varchar(255) NOT NULL COMMENT '组件名称',
    `product_name_display` varchar(255) NOT NULL default '' COMMENT '组件显示名称',
    `product_version` varchar(255) NOT NULL COMMENT '产品版本',
    `product` text NOT NULL COMMENT '产品信息',
    `product_parsed` text NOT NULL COMMENT '已经解析的产品信息',
    `is_current_version` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否当前版本',
    `status` enum('undeployed', 'deploying', 'deployed', 'deploy fail', 'undeploying', 'undeploy fail') NOT NULL DEFAULT 'undeployed' COMMENT '产品状态',
    `alert_recover` tinyint(1) NOT NULL default '0' COMMENT '0，不恢复告警，1，恢复告警',
    `user_id` int(11) NOT NULL default 0 COMMENT '部署人id',
    `deploy_time` datetime COMMENT '部署时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create_time',
    `schema` text NOT NULL COMMENT '产品原始schema',
    PRIMARY KEY (`id`),
    UNIQUE KEY `product_name_version` (`product_name`, `product_version`)
    ) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_product_history` (
    `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `deploy_uuid` char(36) NOT NULL COMMENT '部署uuid',
    `parent_product_name` varchar(255) NOT NULL COMMENT '产品名称',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '组件名称',
    `product_name_display` varchar(255) NOT NULL DEFAULT '' COMMENT '组件显示名称',
    `product_version` varchar(255) NOT NULL DEFAULT '' COMMENT '产品版本',
    `status` char(32) NOT NULL DEFAULT '' COMMENT '产品状态',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create_time',
    `deploy_start_time` datetime COMMENT 'deploy_start_time',
    `deploy_end_time` datetime COMMENT 'deploy_end_time',
    `user_id` int(11) NOT NULL default 0 COMMENT '部署人id',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_schema_field_modify` (
    `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `product_name` varchar(50) NOT NULL DEFAULT '' COMMENT '组件名称',
    `service_name` varchar(50) NOT NULL DEFAULT '' COMMENT '服务名称',
    `field_path` varchar(255) NOT NULL DEFAULT '' COMMENT '字段路径',
    `field` text NOT NULL COMMENT '字段值',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update_time',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create_time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `names_path` (`product_name`, `service_name`, `field_path`)
    ) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_service_ip_list` (
    `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '组件名称',
    `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称',
    `ip_list` varchar(1024) NOT NULL DEFAULT '' COMMENT 'IP列表',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update_time',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create_time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `product_service_name` (`product_name`, `service_name`)
    ) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_instance_list` (
    `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `agent_id` char(36) NOT NULL DEFAULT '' COMMENT 'agent id',
    `sid` char(36) NOT NULL COMMENT '主机 ID (UUID)',
    `pid` int(11) UNSIGNED NOT NULL COMMENT 'product id',
    `ip` varchar(255) NOT NULL DEFAULT '' COMMENT '主机ip',
    `group` varchar(255) NOT NULL DEFAULT '' COMMENT '组名称',
    `prometheus_port` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'prometheus port',
    `service_name` varchar(255) NOT NULL COMMENT '服务名称',
    `service_name_display` varchar(255) NOT NULL default '' COMMENT '服务显示名称',
    `service_version` varchar(255) NOT NULL DEFAULT '' COMMENT '服务版本',
    `schema` longtext NOT NULL  COMMENT 'instance schema',
    `ha_role_cmd` varchar(255) NOT NULL COMMENT 'HA角色执行命令或脚本',
    `health_state` tinyint(2) NOT NULL DEFAULT -2 COMMENT '健康状态,0:不健康,1:健康,-1:未设置,-2:等待',
    `status` enum('installing', 'installed', 'install fail', 'uninstalling', 'uninstalled', 'uninstall fail', 'running', 'run fail', 'stopping', 'stopped', 'stop fail') NOT NULL DEFAULT 'installing' COMMENT '实例状态',
    `status_message` varchar(1024) NOT NULL DEFAULT '' COMMENT '状态详细信息',
    `heart_time` datetime DEFAULT NULL COMMENT '心跳更新时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `pid_service_name` (`pid`, `service_name`, `sid`)
    ) ENGINE=InnoDB  COMMENT='实例列表';

CREATE TABLE IF NOT EXISTS `deploy_instance_record` (
    `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `deploy_uuid` char(36) NOT NULL DEFAULT '' COMMENT '部署记录uuid',
    `instance_id` int(11) UNSIGNED NOT NULL COMMENT '实例id',
    `sid` char(36) NOT NULL DEFAULT '' COMMENT '主机 ID (UUID)',
    `ip` varchar(255) NOT NULL DEFAULT '' COMMENT '主机ip',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '组件名称',
    `product_name_display` varchar(255) NOT NULL DEFAULT '' COMMENT '组件显示名称',
    `product_version` varchar(255) NOT NULL DEFAULT '' COMMENT '产品版本',
    `group` varchar(255) NOT NULL DEFAULT '' COMMENT '组名称',
    `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称',
    `service_name_display` varchar(255) NOT NULL DEFAULT '' COMMENT '服务显示名称',
    `service_version` varchar(255) NOT NULL DEFAULT '' COMMENT '服务版本',
    `status` varchar(32) NOT NULL DEFAULT '' COMMENT '实例状态',
    `status_message` varchar(1024) NOT NULL DEFAULT '' COMMENT '状态详细信息',
    `progress` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '进度',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY (`deploy_uuid`)
    ) ENGINE=InnoDB  COMMENT='实例部署记录表';

CREATE TABLE IF NOT EXISTS deploy_service_ip_node
(
    id int(11) unsigned auto_increment comment 'id' primary key,
    product_name varchar(50) default '' not null comment '组件名称',
    service_name varchar(50) default '' not null comment '服务名称',
    ip varchar(20) default '' not null comment 'IP列表',
    node_id int default '0' not null comment '序号',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    constraint product_service_ip unique (product_name, service_name, ip)
    ) ENGINE=InnoDB  COMMENT='NodeId记录表';

CREATE TABLE IF NOT EXISTS `user_list` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(256) NOT NULL COMMENT '登录密码',
    `company` varchar(128) NOT NULL DEFAULT '' COMMENT '用户所属公司',
    `full_name` varchar(128) NOT NULL DEFAULT '' COMMENT '姓名',
    `email` varchar(255) NOT NULL COMMENT '邮箱地址',
    `phone` varchar(255) NOT NULL DEFAULT '' COMMENT '用户手机号',
    `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0 启动，1 禁用',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB  COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `role_list` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    `role_name` varchar(256) NOT NULL COMMENT '角色名称',
    `role_value` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'ADMIN(1), 运维(2)',
    `role_desc` varchar(256) NOT NULL DEFAULT '' COMMENT '角色描述',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0正常 1逻辑删除',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB  COMMENT='角色表';

CREATE TABLE IF NOT EXISTS `user_role` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `role_id` int(11) NOT NULL COMMENT '角色id',
    `user_id` int(11) NOT NULL COMMENT '用户id',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0正常 1逻辑删除',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB  COMMENT='用户角色关联表';

CREATE TABLE IF NOT EXISTS `deploy_unchecked_service` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `pid` int(11) UNSIGNED NOT NULL COMMENT 'product id',
    `unchecked_services` varchar(255) NOT NULL COMMENT '未勾选服务',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `pid` (`pid`)
    ) ENGINE=InnoDB  COMMENT='未勾选服务表';

insert into role_list(id,role_name,role_value,role_desc) values(-1,"system",-1,"system");
insert into role_list(id,role_name,role_value,role_desc) values(1,"管理员","1","管理员角色");
insert into role_list(id,role_name,role_value,role_desc) values(2,"运维人员","2","运维人员角色");
insert into user_list ( `phone`, `password`, `id`, `company`, `username`, `email`, `full_name`) values ( '11111111111', 'ca6590a271539cc89e2cc20bd6b58518', '1', 'dtstack', 'admin@dtstack.com', 'admin@dtstack.com', 'admin');
insert into user_role(id,role_id,user_id) values(1,1,1);

CREATE TABLE IF NOT EXISTS `deploy_instance_event` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `instance_id` int(11) UNSIGNED NOT NULL COMMENT 'instance id',
    `event_type` enum('install', 'uninstall', 'config update', 'start', 'stop', 'exec', 'error', 'unkown') NOT NULL DEFAULT 'unkown' COMMENT '事件类型',
    `content` text NOT NULL COMMENT '事件内容',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB  COMMENT='实例事件列表';



CREATE TABLE IF NOT EXISTS `deploy_instance_runtime_event` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `event_type` varchar(256) NOT NULL DEFAULT '' COMMENT '事件类型',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '产品名称',
    `parent_product_name` varchar(255) NOT NULL COMMENT '父产品名称',
    `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称',
    `host` varchar(255) NOT NULL DEFAULT '' COMMENT '主机ip',
    `content` text NOT NULL COMMENT '事件内容描述' ,
    `isDeleted` int(11) NOT NULL DEFAULT '0' COMMENT '是否已删除',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB  COMMENT='事件列表';

CREATE TABLE IF NOT EXISTS  `deploy_strategy_list` (
    `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '策略名称',
    `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '策略简介',
    `property` int(1) NOT NULL DEFAULT '0' COMMENT '0:服务类型、1:主机类型',
    `strategy_type` int(1) NOT NULL DEFAULT '0' COMMENT '0:脚本、1:编码',
    `deploy_status` int(1) NOT NULL DEFAULT '0' COMMENT '0:未发布、1:发布',
    `exe_status` int(1) NOT NULL DEFAULT '0' COMMENT '0:正常、1:异常',
    `error_message` text NOT NULL COMMENT '调度状态',
    `start_date` datetime NOT NULL COMMENT '生效日期',
    `end_date` datetime NOT NULL COMMENT '结束日期',
    `start_time` datetime NOT NULL COMMENT '开始时间',
    `end_time` datetime NOT NULL COMMENT '结束时间',
    `cron_period` int(11) NOT NULL DEFAULT '0' COMMENT '调度周期, 0:分钟、1:小时、2:天',
    `cron_interval` int(11) NOT NULL DEFAULT '1' COMMENT '调度间隔时间',
    `cron_time` datetime DEFAULT NULL COMMENT '具体调度时间',
    `params` text NOT NULL COMMENT '参数，逗号间隔',
    `time_out` bigint(10) NOT NULL DEFAULT '-1' COMMENT '超时设置，单位s',
    `is_deleted` int(1) NOT NULL DEFAULT '0' COMMENT '是否删除：0否、1是',
    `gmt_create` datetime NOT NULL COMMENT '创建日期',
    `gmt_modified` datetime NOT NULL COMMENT '最近更新日期',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8 COMMENT='策略表';


CREATE TABLE IF NOT EXISTS  `deploy_strategy_resource_list` (
    `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id` bigint(11) unsigned NOT NULL COMMENT '策略id',
    `content` text NOT NULL COMMENT '资源内容',
    `is_deleted` int(1) NOT NULL DEFAULT '0' COMMENT '是否删除：0否、1是',
    `gmt_create` datetime NOT NULL COMMENT '创建日期',
    `gmt_modified` datetime NOT NULL COMMENT '最近更新日期',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='策略资源表';


CREATE TABLE IF NOT EXISTS  `deploy_strategy_assign_list` (
    `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
    `strategy_id` bigint(11) unsigned NOT NULL COMMENT '策略id',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '产品名称',
    `parent_product_name` varchar(255) NOT NULL COMMENT '父亲产品名称',
    `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称,逗号间隔',
    `host` varchar(255) NOT NULL DEFAULT '' COMMENT '主机ip,逗号间隔',
    `is_deleted` int(1) NOT NULL DEFAULT '0' COMMENT '是否删除：0否、1是',
    `gmt_create` datetime NOT NULL COMMENT '创建日期',
    `gmt_modified` datetime NOT NULL COMMENT '最近更新日期',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB  COMMENT='策略分配表';


INSERT INTO deploy_instance_runtime_event(event_type,product_name,parent_product_name,service_name,content) VALUES('JavaHeapDump','','','','{}');
INSERT INTO deploy_instance_runtime_event(event_type,product_name,parent_product_name,service_name,content) VALUES('服务被动拉起','','','','{}');
INSERT INTO deploy_strategy_list VALUES (1,'JavaHeapDump','JavaHeapDump',1,0,1,0,'',now(),now(),now(),now(),0,5,NULL,'',10,0,now(),'0005-00-00 00:00:00');
INSERT INTO deploy_strategy_resource_list VALUES (1,1,'#!/bin/sh\n\ncurrent=`date \"+%Y-%m-%d %H:%M:%S\"`\ntimeStamp=`date -d \"$current\" +%s`\n\ntmp=\"/tmp/heapdumps_{SERVICENAME}_$timeStamp\"\n\nfind -L /opt/dtstack/{PRODUCTNAME}/{SERVICENAME}/*/heapdump.hprof -maxdepth 5 -size -4096M -type f -mmin -5 -print > $tmp 2>/dev/null\n\nif [ -f \"$tmp\" ];then\nfor i in `cat $tmp`\ndo\necho \'{\"file_name\":{\"desc\":\"JavaHeapDump文件名称\",\"value\":\"\'$i\'\"},\"product_name\":{\"desc\":\"所属组件\",\"value\":\"{PRODUCTNAME}\"},\"service_name\":{\"desc\":\"启动服务\",\"value\":\"{SERVICENAME}\"},\"host\":{\"desc\":\"主机IP\",\"value\":\"{HOSTIP}\"},\"generate_time\":{\"desc\":\"生成时间\",\"value\":\"\'$current\'\"},\"action\":{\"desc\":\"操作\",\"value\":\"下载\",\"instance\":\"{INSTANCEID}\",\"path\":\"\'$i\'\"}}\'\nbreak\ndone\nfi\n\nrm -f $tmp\n',0,now(),now());

INSERT INTO `deploy_strategy_list` VALUES (2,'服务被动拉起','服务被动拉起',1,0,1,0,'',now(),now(),now(),now(),1,1,NULL,'',10,0,now(),now());
INSERT INTO `deploy_strategy_resource_list` VALUES (2,2,'#!/bin/sh\n\ncurrent=`date \"+%Y-%m-%d %H:%M:%S\"`\ntimeStamp=`date -d \"$current\" +%s`\n\nagentId={AGENTID}\ninterval={INTERVAL}\n\n#agentId=cd3d4341-e9ad-446d-9932-8884a034d1cb\n#interval=28\n\nlogs=/opt/dtstack/easymanager/easyagent/logs/agent.log\n\ncdate=`date -d \"-$interval hour\" +\"%Y/%m/%d %H\"`\n\nuser=`whoami`\ns=0\n\nfor i in $(cat $logs |grep \"$cdate\"|grep \"exit(exit status 1\"|grep \"$agentId\"| awk \'{print$1\"#\"$2}\')\ndo\nret=`echo $(echo $(echo \"$i\"|sed \"s/\\//-/g\")| sed \"s/AGENT-DEBUG://g\" )|sed \"s/#/ /g\"`\nresults[s]=\'{\"start_time\":{\"desc\":\"启动时间\",\"value\":\"\'$ret\'\"},\"service_name\":{\"desc\":\"启动服务\",\"value\":\"{SERVICENAME}\"},\"host\":{\"desc\":\"启动主机\",\"value\":\"{HOSTIP}\"},\"product_name\":{\"desc\":\"所属组件\",\"value\":\"{PRODUCTNAME}\"},\"run_user\":{\"desc\":\"启动用户\",\"value\":\"\'$user\'\"}}\'\ns=$[$s+1];\ndone\n\nlen=${#results[@]}\nfor ((i=$len - 1;i>=0;i--))\ndo\n    echo ${results[$i]}\ndone',0,now(),now());


CREATE TABLE IF NOT EXISTS `deploy_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `sid` varchar(255) NULL DEFAULT '' COMMENT '主机sid',
  `hostname` varchar(255) NULL DEFAULT '' COMMENT '主机名',
  `ip` varchar(255) NULL DEFAULT '' COMMENT '主机ip',
  `status` int(11) NULL DEFAULT '6' COMMENT '6:K8S NODE初始化成功,-6:K8S NODE初始化失败',
  `steps`  int(11) NULL DEFAULT '6' COMMENT '6:K8S NODE初始化成功',
  `group` varchar(255) DEFAULT 'default' COMMENT '组信息，兼容deploy_host',
  `errorMsg` varchar(1024) NULL DEFAULT '' COMMENT '错误信息',
  `isDeleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
  `updated` datetime NOT NULL COMMENT 'updated',
  `created` datetime NOT NULL COMMENT 'created',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_cluster_image_store` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `clusterId` int(11) UNSIGNED NOT NULL COMMENT '集群 id',
  `name` varchar(64) NOT NULL COMMENT '仓库名称',
  `alias` varchar(64) NOT NULL COMMENT '仓库别名',
  `address` varchar(256) NOT NULL COMMENT '仓库地址',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `password` varchar(64) NOT NULL COMMENT '密码',
  `email` varchar(64) DEFAULT '' COMMENT '邮箱',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  COMMENT='k8s镜像仓库表';

CREATE TABLE IF NOT EXISTS `deploy_cluster_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(255) NOT NULL COMMENT '集群名',
  `type` varchar(255) NOT NULL COMMENT '集群类型 hosts/kubernetes',
  `mode` int(11) NOT NULL DEFAULT '0' COMMENT '0:自建,1:导入',
  `version` varchar(255) NULL DEFAULT '' COMMENT '集群版本，主机集群为空',
  `desc` varchar(255) NULL DEFAULT '' COMMENT '集群描述',
  `tags` varchar(1024) NULL DEFAULT '' COMMENT '集群标签',
  `configs` text DEFAULT NULL COMMENT '集群个性化配置',
  `yaml` text DEFAULT NULL COMMENT '集群配置详情',
  `status` int(11) NULL DEFAULT '0' COMMENT '0:Waiting,1:Pending,2:Running,-2:Error',
  `errorMsg` varchar(1024) NULL DEFAULT '' COMMENT '错误信息',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
  `create_user` varchar(255) NULL DEFAULT 'admin' COMMENT '创建人',
  `update_user` varchar(255) NULL DEFAULT 'admin' COMMENT '修改人',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_cluster_host_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `sid` varchar(255) NULL DEFAULT '' COMMENT '主机sid',
  `clusterId` int(11) UNSIGNED NOT NULL COMMENT '集群 id',
  `roles` varchar(1024) NULL DEFAULT '' COMMENT '角色，逗号间隔',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;


CREATE TABLE IF NOT EXISTS `deploy_cluster_k8s_available` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `mode` int(11) NOT NULL DEFAULT '0' COMMENT '0:自建,1:导入',
  `version` varchar(255) NULL DEFAULT '' COMMENT 'k8s集群版本',
  `properties` text NOT NULL COMMENT '版本所允许的配置信息，分号区分配置项，冒号区分配置名，逗号区分配置选项',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_cluster_k8s_only` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `clusterId` int(11) UNSIGNED NOT NULL COMMENT '集群 id',
  `kube_config` text DEFAULT NULL COMMENT '集群配置文件，主机集群为空',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;


INSERT INTO `deploy_cluster_list` (`id`, `name`, `type`, `mode`, `desc`, `tags`, `status`, `create_user`, `update_user`, `update_time`, `create_time`) VALUES ('1', 'dtstack', 'hosts', '0', '兼容EM2.0默认集群', '自动创建', '0', 'admin@dtstack.com', 'admin@dtstack.com', now(), now());
INSERT INTO `deploy_cluster_host_rel` (`id`, `sid`, `clusterId`,  `update_time`, `create_time`) select id, sid, '1','0000-00-00 00:00:00','0000-00-00 00:00:00' from deploy_host;
INSERT INTO `deploy_cluster_k8s_available` (`version`, `properties`) VALUES ('v1.16.3-rancher1-1', 'network_plugin:flannel');
ALTER TABLE sidecar_list ADD COLUMN `disk_usage_pct` decimal(6,2) DEFAULT '-1.00' COMMENT '磁盘使用率' AFTER `mem_usage`;


CREATE TABLE IF NOT EXISTS `deploy_cluster_product_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `pid` varchar(255) NULL DEFAULT '' COMMENT '产品 id',
  `clusterId` int(11) UNSIGNED NOT NULL COMMENT '集群 id',
  `deploy_uuid` varchar(36) NOT NULL COMMENT '部署uuid',
  `product_parsed` longtext NOT NULL COMMENT '已经解析的产品信息',
  `status` enum('undeployed','deploying','deployed','deploy fail','undeploying','undeploy fail') NOT NULL DEFAULT 'undeployed' COMMENT '产品状态',
  `alert_recover` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0，不恢复告警，1，恢复告警',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '部署人id',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
  `deploy_time` datetime COMMENT '部署时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

ALTER TABLE deploy_product_list ADD COLUMN `product_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '产品包类型,0,传统包, 1, k8s包' AFTER `schema`;
ALTER TABLE deploy_service_ip_list ADD COLUMN `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id' AFTER `id`;
ALTER TABLE deploy_service_ip_node ADD COLUMN `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id' AFTER `id`;
ALTER TABLE deploy_schema_field_modify ADD COLUMN `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id' AFTER `id`;
ALTER TABLE deploy_unchecked_service ADD COLUMN `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id' AFTER `id`;
ALTER TABLE deploy_product_history ADD COLUMN `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id' AFTER `id`;
ALTER TABLE deploy_instance_list ADD COLUMN `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id' AFTER `id`;

alter table deploy_service_ip_list drop index `product_service_name`;
alter table deploy_service_ip_list add unique key `cluster_product_service_name` (`cluster_id`,`product_name`,`service_name`);
alter table deploy_service_ip_node drop index `product_service_ip`;
alter table deploy_service_ip_node add unique key `cluster_product_service_ip` (`cluster_id`,`product_name`,`service_name`,`ip`);
alter table deploy_schema_field_modify drop index `names_path`;
alter table deploy_schema_field_modify add unique key `cluster_names_path` (`cluster_id`,`product_name`,`service_name`,`field_path`);
alter table deploy_unchecked_service drop index `pid`;
alter table deploy_unchecked_service add unique key `cluster_pid` (`cluster_id`,`pid`);
alter table deploy_instance_list drop index `pid_service_name`;
alter table deploy_instance_list add unique key `cluster_pid_service_name` (`cluster_id`,`pid`,`service_name`,`sid`);


CREATE TABLE IF NOT EXISTS `deploy_kube_base_product_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `pid` varchar(255) DEFAULT '' COMMENT 'pid',
  `clusterId` int(11) unsigned NOT NULL COMMENT 'cluster id',
  `namespace` varchar(255) DEFAULT '' COMMENT 'k8s部署模式，命名空间',
  `base_clusterId` varchar(36) NOT NULL COMMENT '依赖集群',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '1: 已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `pid_clusterId_name_base` (`pid`, `clusterId`, `namespace`)
) ENGINE=InnoDB AUTO_INCREMENT=35 ;


CREATE TABLE IF NOT EXISTS `deploy_kube_product_lock` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `pid` varchar(255) DEFAULT '' COMMENT 'pid',
  `clusterId` int(11) unsigned NOT NULL COMMENT 'cluster id',
  `namespace` varchar(255) DEFAULT '' COMMENT 'k8s部署模式，命名空间',
  `is_deploy` int(11) NOT NULL DEFAULT 0 COMMENT '1: deploying | 0: not deploy',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '1: 已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;


CREATE TABLE IF NOT EXISTS `deploy_cluster_kube_pod_list` (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
  `pid` int(11) UNSIGNED NOT NULL COMMENT 'product id',
  `clusterId` int(11) unsigned NOT NULL COMMENT 'cluster id',
  `namespace` varchar(255) DEFAULT '' COMMENT '命名空间',
  `product_name` varchar(255) NOT NULL COMMENT '产品名称',
  `product_version` varchar(255) NOT NULL DEFAULT '' COMMENT '产品版本',
  `service_name` varchar(255) NOT NULL COMMENT '服务名称',
  `service_version` varchar(255) NOT NULL DEFAULT '' COMMENT '服务版本',
  `pod_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'pod uid',
  `pod_name` varchar(255) NOT NULL DEFAULT '' COMMENT 'pod names',
  `pod_key` varchar(255) NOT NULL DEFAULT '' COMMENT 'pod informer cache key',
  `self_link` text NOT NULL COMMENT 'pod self link',
  `host_ip` varchar(255) NOT NULL DEFAULT '' COMMENT '主机ip',
  `pod_ip` varchar(255) NOT NULL DEFAULT '' COMMENT '主机ip',
  `phase` enum('Pending', 'Running', 'Succeeded', 'Failed', 'Unknown') NOT NULL DEFAULT 'Pending' COMMENT 'Pod状态',
  `message` varchar(1024) NOT NULL DEFAULT '' COMMENT '状态详细信息',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '1: 已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_kube_service_list` (
 `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
  `pid` int(11) UNSIGNED NOT NULL COMMENT 'product id',
  `clusterId` int(11) unsigned NOT NULL COMMENT 'cluster id',
  `namespace` varchar(255) DEFAULT '' COMMENT '命名空间',
  `product_name` varchar(255) NOT NULL COMMENT '产品名称',
  `product_version` varchar(255) NOT NULL DEFAULT '' COMMENT '产品版本',
  `service_name` varchar(255) NOT NULL COMMENT '服务名称',
  `service_version` varchar(255) NOT NULL DEFAULT '' COMMENT '服务版本',
  `cluster_ip` varchar(255) NOT NULL DEFAULT '' COMMENT 'cluster ip of the service',
  `type` enum('ClusterIP', 'NodePort', 'LoadBalancer', 'ExternalName') NOT NULL DEFAULT 'ClusterIP' COMMENT 'service type',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '1: 已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

ALTER TABLE deploy_cluster_image_store ADD COLUMN `is_default` tinyint(1) NOT NULL COMMENT '1:default image store' AFTER `clusterId`;
ALTER TABLE deploy_cluster_product_rel ADD COLUMN `namespace` varchar(255) NOT NULL COMMENT 'cluster namespace' AFTER `clusterId`;
ALTER TABLE deploy_product_history ADD COLUMN `product_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT ',0,, 1, k8s' AFTER `product_version`;
ALTER TABLE deploy_product_history ADD COLUMN `namespace` varchar(255) NOT NULL COMMENT 'cluster namespace' AFTER `cluster_id`;
ALTER TABLE deploy_instance_list ADD COLUMN `namespace` varchar(255) NOT NULL default '' COMMENT 'cluster namespace' AFTER `cluster_id`;
ALTER TABLE deploy_service_ip_list ADD COLUMN `namespace` varchar(255) NOT NULL default '' COMMENT 'cluster namespace' AFTER `cluster_id`;
ALTER TABLE deploy_service_ip_node ADD COLUMN `namespace` varchar(255) NOT NULL default '' COMMENT 'cluster namespace' AFTER `cluster_id`;
ALTER TABLE deploy_cluster_product_rel MODIFY COLUMN `pid` int(11);

alter table deploy_instance_list drop index `cluster_pid_service_name`;
alter table deploy_service_ip_list drop index `cluster_product_service_name`;
alter table deploy_service_ip_list add unique key `cluster_ns_product_service_name` (`cluster_id`,`namespace`,`product_name`,`service_name`);
alter table deploy_service_ip_node drop index `cluster_product_service_ip`;
alter table deploy_service_ip_node add unique key `cluster_ns_product_service_ip` (`cluster_id`,`namespace`,`product_name`,`service_name`,`ip`);

CREATE TABLE IF NOT EXISTS `safety_audit_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `operator` varchar(255) NOT NULL COMMENT '操作人',
  `module` varchar(255) NOT NULL COMMENT '操作模块',
  `operation` varchar(255) NOT NULL COMMENT '动作',
  `ip` varchar(255) NOT NULL COMMENT '来源ip',
  `content` text NOT NULL COMMENT '详细内容',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '1: 已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `audit_item_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `module` varchar(255) NOT NULL COMMENT '操作模块',
  `operation` varchar(255) NOT NULL COMMENT '动作',
  `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '1: 已删除',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB ;

UPDATE role_list SET role_list.role_name='Administrator', role_list.role_desc='超级管理员，具备产品所有操作权限' WHERE role_list.role_value='1';
UPDATE role_list SET role_list.role_name='Cluster Operator', role_list.role_desc='集群操作人员，一般指运维人员，具有安装部署、集群运维、监控告警功能操作权限' WHERE role_list.role_value='2';
INSERT INTO role_list(role_name, role_value, role_desc) VALUES('Cluster Reader', 4, '普通用户，只有集群的只读权限');

INSERT INTO audit_item_list(module, operation) VALUES('产品访问', '进入EM');
INSERT INTO audit_item_list(module, operation) VALUES('产品访问', '退出EM');
INSERT INTO audit_item_list(module, operation) VALUES('用户管理', '创建账号');
INSERT INTO audit_item_list(module, operation) VALUES('用户管理', '禁用账号');
INSERT INTO audit_item_list(module, operation) VALUES('用户管理', '启用账号');
INSERT INTO audit_item_list(module, operation) VALUES('用户管理', '移除账号');
INSERT INTO audit_item_list(module, operation) VALUES('用户管理', '重置密码');
INSERT INTO audit_item_list(module, operation) VALUES('集群管理', '创建集群');
INSERT INTO audit_item_list(module, operation) VALUES('集群管理', '编辑集群');
INSERT INTO audit_item_list(module, operation) VALUES('集群管理', '删除集群');
INSERT INTO audit_item_list(module, operation) VALUES('集群管理', '添加主机');
INSERT INTO audit_item_list(module, operation) VALUES('集群管理', '删除主机');
INSERT INTO audit_item_list(module, operation) VALUES('部署向导', '产品部署');
INSERT INTO audit_item_list(module, operation) VALUES('部署向导', '产品卸载');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '组件停止');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '组件启动');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '服务停止');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '服务启动');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '服务滚动重启');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '服务参数修改');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '服务参数新增');
INSERT INTO audit_item_list(module, operation) VALUES('集群运维', '配置下发');

CREATE TABLE IF NOT EXISTS `deploy_namespace_client` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `yaml` text,
  `namespace_id` int(11) DEFAULT NULL,
  `file_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 ;

CREATE TABLE IF NOT EXISTS `deploy_namespace_event` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `resource` varchar(255) DEFAULT NULL,
  `message` varchar(1500) DEFAULT NULL,
  `namespace_id` int(11) DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 ;

CREATE TABLE IF NOT EXISTS `deploy_namespace_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(255) NOT NULL,
  `namespace` varchar(255) NOT NULL,
  `registry_id` int(11) DEFAULT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `port` varchar(255) DEFAULT NULL,
  `cluster_id` int(11) NOT NULL,
  `status` varchar(255) DEFAULT NULL,
  `is_deleted` int(11) DEFAULT NULL,
  `user` varchar(255) DEFAULT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 ;

CREATE TABLE IF NOT EXISTS `import_init_moudle` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `service_account` varchar(1000) DEFAULT NULL,
  `role` varchar(1000) DEFAULT NULL,
  `role_binding` varchar(1000) DEFAULT NULL,
  `operator` varchar(1000) DEFAULT NULL,
  `is_deleted` int(2) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 ;

INSERT INTO `import_init_moudle` VALUES (1, '{\"apiVersion\":\"v1\",\"kind\":\"ServiceAccount\",\"metadata\":{\"name\":\"dtstack\",\"namespace\":\"{{.NAME_SPACE}}\"}}', '{\"apiVersion\":\"rbac.authorization.k8s.io/v1\",\"kind\":\"Role\",\"metadata\":{\"name\":\"dtstack-admin\",\"namespace\":\"{{.NAME_SPACE}}\"},\"rules\":[{\"apiGroups\":[\"*\"],\"resources\":[\"*\"],\"verbs\":[\"*\"]}]}', '{\"apiVersion\":\"rbac.authorization.k8s.io/v1\",\"kind\":\"RoleBinding\",\"metadata\":{\"name\":\"dtstack-admin-binding\",\"namespace\":\"{{.NAME_SPACE}}\"},\"roleRef\":{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"Role\",\"name\":\"dtstack-admin\"},\"subjects\":[{\"kind\":\"ServiceAccount\",\"name\":\"dtstack\",\"namespace\":\"{{.NAME_SPACE}}\"}]}', '{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"name\":\"mole-operator\",\"namespace\":\"{{.NAME_SPACE}}\"},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"name\":\"mole-operator\"}},\"template\":{\"metadata\":{\"labels\":{\"name\":\"mole-operator\"}},\"spec\":{\"containers\":[{\"command\":[\"mole-operator\"],\"env\":[{\"name\":\"WATCH_NAMESPACE\",\"value\":\"{{.NAME_SPACE}}\"},{\"name\":\"POD_NAME\",\"valueFrom\":{\"fieldRef\":{\"fieldPath\":\"metadata.name\"}}},{\"name\":\"OPERATOR_NAME\",\"value\":\"mole-operator\"}],\"image\":\"{{.REGISTRY}}/mole:v1.0.17\",\"imagePullPolicy\":\"Always\",\"name\":\"mole-operator\",\"resources\":{\"limits\":{\"cpu\":\"500m\",\"memory\":\"500Mi\"}}}],\"imagePullSecrets\":[{\"name\":\"{{.SECRET_NAME}}\"}]}}}}', 0);


CREATE TABLE IF NOT EXISTS `workload_definition` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `version` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `params` text COLLATE utf8mb4_bin,
  `latest` int(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7  COLLATE=utf8mb4_bin;

BEGIN;
INSERT INTO `workload_definition` VALUES (1, 'zookeeper', 'v1', '[{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.image\"},{\"key\":\"Replica\",\"ref\":\"spec.workloadpatrs.0.baseworkload.parameters.spec.replicas\"},{\"key\":\"ResourceRequest.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.requests.cpu\"},{\"key\":\"ResourceRequest.memory\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.requests.memory\"},{\"key\":\"ResourceLimit.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.limits.cpu\"},{\"key\":\"ResourceLimit.memory\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.limits.memory\"},{\"key\":\"Ports.0\",\"ref\":\"spec.workloadpatrs.0.steps.3.object.spec.ports.0.port\"},{\"key\":\"Environment.ZK_SERVERS\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.env.0.value\"},{\"key\":\"Environment.ZOO_INIT_LIMIT\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.env.1.value\"},{\"key\":\"Environment.ZOO_MAX_CLIENT_CNXNS\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.env.2.value\"},{\"key\":\"StorageClass\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.spec.storageClassName\"}]', 1);
INSERT INTO `workload_definition` VALUES (2, 'redis', 'v1', '[{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.3.object.image\"},{\"key\":\"Replica\",\"ref\":\"spec.workloadpatrs.0.baseworkload.parameters.spec.replicas\"},{\"key\":\"ResourceRequest.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.3.object.resources.requests.cpu\"},{\"key\":\"ResourceRequest.memory\",\"ref\":\"spec.workloadpatrs.0.steps.3.object.resources.requests.memory\"},{\"key\":\"ResourceLimit.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.3.object.resources.limits.cpu\"},{\"key\":\"ResourceLimit.memory\",\"ref\":\"spec.workloadpatrs.0.steps.3.object.resources.limits.memory\"},{\"key\":\"Ports.0\",\"ref\":\"spec.workloadpatrs.0.steps.5.object.spec.ports.0.port\"},{\"key\":\"ConfigPaths\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.data\"}]', 1);
INSERT INTO `workload_definition` VALUES (3, 'mysql', 'v1', '[{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.6.object.image\"},{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.4.object.image\"},{\"key\":\"Replica\",\"ref\":\"spec.workloadpatrs.0.baseworkload.parameters.spec.replicas\"},{\"key\":\"ResourceRequest.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.6.object.resources.requests.cpu\"},{\"key\":\"ResourceRequest.memory\",\"ref\":\"spec.workloadpatrs.0.steps.6.object.resources.requests.memory\"},{\"key\":\"ResourceLimit.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.6.object.resources.limits.cpu\"},{\"key\":\"ResourceLimit.memory\",\"ref\":\"spec.workloadpatrs.0.steps.6.object.resources.limits.memory\"},{\"key\":\"Ports.0\",\"ref\":\"spec.workloadpatrs.0.steps.7.object.spec.ports.0.port\"},{\"key\":\"Environment.MYSQL_ALLOW_EMPTY_PASSWORD\",\"ref\":\"spec.workloadpatrs.0.steps.6.object.env.0.value\"},{\"key\":\"ConfigPaths\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.data\"},{\"key\":\"StorageClass\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.spec.storageClassName\"}]', 1);
INSERT INTO `workload_definition` VALUES (4, 'pushgateway', 'v1', '[{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.image\"},{\"key\":\"Replica\",\"ref\":\"spec.workloadpatrs.0.baseworkload.parameters.spec.replicas\"},{\"key\":\"ResourceRequest.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.requests.cpu\"},{\"key\":\"ResourceRequest.memory\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.requests.memory\"},{\"key\":\"ResourceLimit.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.limits.cpu\"},{\"key\":\"ResourceLimit.memory\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.resources.limits.memory\"},{\"key\":\"Ports.0\",\"ref\":\"spec.workloadpatrs.0.steps.2.object.spec.ports.0.port\"}]', 1);
INSERT INTO `workload_definition` VALUES (5, 'prometheus', 'v1', '[{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.4.object.image\"},{\"key\":\"Replica\",\"ref\":\"spec.workloadpatrs.0.baseworkload.parameters.spec.replicas\"},{\"key\":\"ResourceRequest.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.4.object.resources.requests.cpu\"},{\"key\":\"ResourceRequest.memory\",\"ref\":\"spec.workloadpatrs.0.steps.4.object.resources.requests.memory\"},{\"key\":\"ResourceLimit.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.4.object.resources.limits.cpu\"},{\"key\":\"ResourceLimit.memory\",\"ref\":\"spec.workloadpatrs.0.steps.4.object.resources.limits.memory\"},{\"key\":\"Ports.0\",\"ref\":\"spec.workloadpatrs.0.steps.5.object.spec.ports.0.port\"},{\"key\":\"ConfigPaths\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.data\"},{\"key\":\"StorageClass\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.spec.storageClassName\"}]', 1);
INSERT INTO `workload_definition` VALUES (6, 'kafka', 'v1', '[{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.image\"},{\"key\":\"Replica\",\"ref\":\"spec.workloadpatrs.0.baseworkload.parameters.spec.replicas\"},{\"key\":\"ResourceRequest.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.resources.requests.cpu\"},{\"key\":\"ResourceRequest.memory\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.resources.requests.memory\"},{\"key\":\"ResourceLimit.cpu\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.resources.limits.cpu\"},{\"key\":\"ResourceLimit.memory\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.resources.limits.memory\"},{\"key\":\"Ports.0\",\"ref\":\"spec.workloadpatrs.0.steps.1.object.spec.ports.0.port\"},{\"key\":\"Environment.KAFKA_HEAP_OPTS\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.env.0.value\"},{\"key\":\"Environment.KAFKA_OPTS\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.env.1.value\"}]', 1);
COMMIT;
CREATE TABLE IF NOT EXISTS `workload_part` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `type` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `parameters` text COLLATE utf8mb4_bin,
  `workload_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7  COLLATE=utf8mb4_bin;

BEGIN;
INSERT INTO `workload_part` VALUES (1, 'leader', 'statefulset', '{\"spec\":{\"replicas\":3,\"serviceName\":\"@zk-hs\"}}', 1);
INSERT INTO `workload_part` VALUES (2, 'singlenode', 'deployment', '{\"spec\":{\"replicas\":1}}', 2);
INSERT INTO `workload_part` VALUES (3, 'master-slave', 'statefulset', '{\"spec\":{\"replicas\":2,\"serviceName\":\"@mysql-svc\"}}', 3);
INSERT INTO `workload_part` VALUES (4, 'singlenode', 'deployment', '{\"spec\":{\"replicas\":1}}', 4);
INSERT INTO `workload_part` VALUES (5, 'master', 'statefulset', '{\"spec\":{\"replicas\":1,\"serviceName\":\"@prometheus-svc\"}}', 5);
INSERT INTO `workload_part` VALUES (6, 'master', 'statefulset', '{\"spec\":{\"replicas\":2,\"serviceName\":\"@kafka-svc\"}}', 6);
COMMIT;

CREATE TABLE IF NOT EXISTS `workload_step` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `type` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `action` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `object` text COLLATE utf8mb4_bin,
  `workloadpart_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=32  COLLATE=utf8mb4_bin;

BEGIN;
INSERT INTO `workload_step` VALUES (1, 'zkdata', 'pvc', 'bound', '{\"spec\":{\"storageClassName\":\"local-storage-workload-zk\",\"accessModes\":[\"ReadWriteOnce\"],\"resources\":{\"requests\":{\"storage\":\"5G\"}}}}', 1);
INSERT INTO `workload_step` VALUES (2, 'zk-sts', 'container', 'bound', '{\"image\":\"zookeeper:3.4.14_02\",\"env\":[{\"name\":\"ZK_SERVERS\",\"value\":\"3\"},{\"name\":\"ZOO_INIT_LIMIT\",\"value\":\"10\"},{\"name\":\"ZOO_MAX_CLIENT_CNXNS\",\"value\":\"200\"},{\"name\":\"SERVER_JVMFLAGS\",\"value\":\"-javaagent:./dtstack/prometheus/jmx_prometheus_javaagent-0.14.0.jar=9505:./dtstack/prometheus/zookeeper.yml\"}],\"resources\":{\"requests\":{\"memory\":\"1Gi\",\"cpu\":\"1\"},\"limits\":{\"memory\":\"3Gi\",\"cpu\":\"3\"}},\"ports\":[{\"containerPort\":2888,\"name\":\"server\"},{\"containerPort\":2181,\"name\":\"client\"},{\"containerPort\":3888,\"name\":\"leader-election\"},{\"containerPort\":9505,\"name\":\"jmx-prom-agent\"}],\"command\":[\"sh\",\"-c\",\"/zookeeper/bin/dtstack/start-zookeeper-k8s.sh\"],\"readinessProbe\":{\"exec\":{\"command\":[\"/zookeeper/bin/dtstack/zookeeper-healthcheck.sh\",\"2181\"]},\"initialDelaySeconds\":10,\"timeoutSeconds\":5},\"livenessProbe\":{\"exec\":{\"command\":[\"/zookeeper/bin/dtstack/zookeeper-healthcheck.sh\",\"2181\"]},\"initialDelaySeconds\":10,\"timeoutSeconds\":5},\"volumeMounts\":[{\"name\":\"zkdata\",\"mountPath\":\"/data\",\"subPath\":\"zk_data\"},{\"name\":\"zkdata\",\"mountPath\":\"/datalog\",\"subPath\":\"zk_datalog\"},{\"name\":\"zkdata\",\"mountPath\":\"/logs\",\"subPath\":\"zk_logs\"}],\"securityContext\":{\"runAsUser\":0,\"runAsGroup\":0}}', 1);
INSERT INTO `workload_step` VALUES (3, 'zk-hs', 'service', 'createorupdate', '{\"spec\":{\"ports\":[{\"port\":2888,\"name\":\"server\"},{\"port\":3888,\"name\":\"leader-election\"},{\"port\":9505,\"name\":\"jmx-prom-agent\"}],\"clusterIP\":\"None\",\"selector\":{\"app\":\"@leader\"}}}', 1);
INSERT INTO `workload_step` VALUES (4, 'zk-cs', 'service', 'createorupdate', '{\"spec\":{\"ports\":[{\"port\":2181,\"name\":\"client\",\"targetPort\":2181}],\"selector\":{\"app\":\"@leader\"}}}', 1);
INSERT INTO `workload_step` VALUES (5, 'redis-cm', 'conf', 'createorupdate', '{\"data\":{\"redis.conf\":\"bind 127.0.0.1\\nport 6379\\ntcp-backlog 511\\ntimeout 0\\ntcp-keepalive 300\\ndaemonize no\\nsupervised no\\npidfile /var/run/redis_6379.pid\\nloglevel notice\\nlogfile \\\"\\\"\\ndatabases 16\\nsave 900 1\\nsave 300 10\\nsave 60 10000\\nstop-writes-on-bgsave-error yes\\nrdbcompression yes\\nrdbchecksum yes\\ndbfilename dump.rdb\\ndir ./\\nslave-serve-stale-data yes\\nslave-read-only yes\\nrepl-diskless-sync no\\nrepl-diskless-sync-delay 5\\nrepl-disable-tcp-nodelay no\\nslave-priority 100\\nappendonly no\\nappendfilename \\\"appendonly.aof\\\"\\nappendfsync everysec\\nno-appendfsync-on-rewrite no\\nauto-aof-rewrite-percentage 100\\nauto-aof-rewrite-min-size 64mb\\naof-load-truncated yes\\nlua-time-limit 5000\\nslowlog-log-slower-than 10000\\nslowlog-max-len 128\\nlatency-monitor-threshold 0\\nnotify-keyspace-events \\\"\\\"\\nhash-max-ziplist-entries 512\\nhash-max-ziplist-value 64\\nlist-max-ziplist-size -2\\nlist-compress-depth 0\\nzset-max-ziplist-entries 128\\nzset-max-ziplist-value 64\\nhll-sparse-max-bytes 3000\\nactiverehashing yes\\nclient-output-buffer-limit normal 0 0 0\\nclient-output-buffer-limit slave 256mb 64mb 60\\nclient-output-buffer-limit pubsub 32mb 8mb 60\\nhz 10\\naof-rewrite-incremental-fsync yes\\nrename-command FLUSHALL \\\"\\\"\\nrename-command FLUSHDB  \\\"\\\"\\nrename-command SHUTDOWN \\\"\\\"\\n\"}}', 2);
INSERT INTO `workload_step` VALUES (6, 'redisconfig-volume', 'volume', 'bound', '{\"configMap\":{\"name\":\"@redis-cm\"}}', 2);
INSERT INTO `workload_step` VALUES (7, 'redis-localtime', 'volume', 'bound', '{\"hostPath\":{\"path\":\"/etc/localtime\"}}', 2);
INSERT INTO `workload_step` VALUES (8, 'redis-dp', 'container', 'bound', '{\"image\":\"redis:3.2.12_02\",\"imagePullPolicy\":\"Always\",\"resources\":{\"requests\":{\"memory\":\"1Gi\",\"cpu\":\"1\"},\"limits\":{\"memory\":\"3Gi\",\"cpu\":\"3\"}},\"ports\":[{\"containerPort\":6379}],\"livenessProbe\":{\"exec\":{\"command\":[\"/usr/bin/dtstack/redis-liveness.sh\"]},\"initialDelaySeconds\":30,\"periodSeconds\":10,\"timeoutSeconds\":5},\"readinessProbe\":{\"exec\":{\"command\":[\"/usr/bin/dtstack/redis-readiness.sh\"]},\"initialDelaySeconds\":5,\"periodSeconds\":2,\"timeoutSeconds\":1},\"volumeMounts\":[{\"name\":\"redisconfig-volume\",\"mountPath\":\"/usr/local/etc/redis\"}]}', 2);
INSERT INTO `workload_step` VALUES (9, 'redis-exporter', 'container', 'bound', '{\"image\":\"redis_exporter:v1.14.0_01\",\"resources\":{\"requests\":{\"cpu\":\"100m\",\"memory\":\"100Mi\"}},\"ports\":[{\"containerPort\":9121}],\"volumeMounts\":[{\"name\":\"redis-localtime\",\"mountPath\":\"/etc/localtime\"}]}', 2);
INSERT INTO `workload_step` VALUES (10, 'redis-svc', 'service', 'createorupdate', '{\"metadata\":{\"annotations\":{\"prometheus.io/scrape\":\"true\",\"prometheus.io/port\":\"9121\"}},\"spec\":{\"ports\":[{\"name\":\"redis\",\"port\":6379,\"targetPort\":6379},{\"name\":\"prom\",\"port\":9121,\"targetPort\":9121}],\"selector\":{\"app\":\"@singlenode\"}}}', 2);
INSERT INTO `workload_step` VALUES (11, 'mysql-data', 'pvc', 'bound', '{\"spec\":{\"storageClassName\":\"local-storage-workload-mysql\",\"accessModes\":[\"ReadWriteOnce\"],\"resources\":{\"requests\":{\"storage\":\"5Gi\"}}}}', 3);
INSERT INTO `workload_step` VALUES (12, 'mysqlcm', 'conf', 'createorupdate', '{\"primary.cnf\":\"\\n# Apply this config only on the primary.\\n[mysqld]\\nlog-bin\\n\\n\",\"replica.cnf\":\"\\n# Apply this config only on replicas.\\n[mysqld]\\nsuper-read-only\\n\",\"hive.sql\":\"\\n# init hive sql.\\n\",\"createuser.sh\":\"\\n# create mysql user script.\\n\"}', 3);
INSERT INTO `workload_step` VALUES (13, 'mysql-conf', 'volume', 'bound', '{\"emptyDir\":{}}', 3);
INSERT INTO `workload_step` VALUES (14, 'mysql-cmvolume', 'volume', 'bound', '{\"configMap\":{\"name\":\"@mysqlcm\"}}', 3);
INSERT INTO `workload_step` VALUES (15, 'init-mysql', 'init-container', 'bound', '{\"image\":\"mysql:5.7.32_02\",\"imagePullPolicy\":\"Always\",\"env\":[{\"name\":\"MYSQL_ALLOW_EMPTY_PASSWORD\",\"value\":\"1\"}],\"command\":[\"bash\",\"-c\",\"/usr/bin/dtstack/init-mysql.sh\"],\"volumeMounts\":[{\"name\":\"mysql-conf\",\"mountPath\":\"/mnt/conf.d\"},{\"name\":\"mysql-cmvolume\",\"mountPath\":\"/mnt/config-map\"},{\"name\":\"mysql-data\",\"mountPath\":\"/var/lib/mysql\",\"subPath\":\"mysql\"}]}', 3);
INSERT INTO `workload_step` VALUES (16, 'clone-mysql', 'init-container', 'bound', '{\"image\":\"xtrabackup:2.4.21_02\",\"env\":[{\"name\":\"POD_NAMESPACE\",\"valueFrom\":{\"fieldRef\":{\"fieldPath\":\"metadata.namespace\"}}}],\"command\":[\"bash\",\"-c\",\"/usr/bin/dtstack/clone-mysql.sh\"],\"volumeMounts\":[{\"name\":\"mysql-data\",\"mountPath\":\"/var/lib/mysql\",\"subPath\":\"mysql\"},{\"name\":\"mysql-conf\",\"mountPath\":\"/etc/mysql/conf.d\"}]}', 3);
INSERT INTO `workload_step` VALUES (17, 'mysql-sts', 'container', 'bound', '{\"image\":\"mysql:5.7.32_02\",\"imagePullPolicy\":\"Always\",\"env\":[{\"name\":\"MYSQL_ALLOW_EMPTY_PASSWORD\",\"value\":\"1\"}],\"ports\":[{\"name\":\"mysql\",\"containerPort\":3306}],\"volumeMounts\":[{\"name\":\"mysql-data\",\"mountPath\":\"/var/lib/mysql\",\"subPath\":\"mysql\"},{\"name\":\"mysql-conf\",\"mountPath\":\"/etc/mysql/conf.d\"}],\"resources\":{\"requests\":{\"memory\":\"1Gi\",\"cpu\":\"1\"},\"limits\":{\"memory\":\"3Gi\",\"cpu\":\"3\"}},\"livenessProbe\":{\"exec\":{\"command\":[\"sh\",\"-c\",\"/usr/bin/dtstack/mysql-liveness.sh\"]},\"initialDelaySeconds\":30,\"periodSeconds\":10,\"timeoutSeconds\":5},\"readinessProbe\":{\"exec\":{\"command\":[\"sh\",\"-c\",\"/usr/bin/dtstack/mysql-liveness.sh\"]},\"initialDelaySeconds\":30,\"periodSeconds\":2,\"timeoutSeconds\":1}}', 3);
INSERT INTO `workload_step` VALUES (18, 'mysql-svc', 'service', 'createorupdate', '{\"spec\":{\"ports\":[{\"name\":\"mysql\",\"port\":3306,\"targetPort\":3306,\"protocol\":\"TCP\"},{\"name\":\"mysql-exporter\",\"port\":9104,\"protocol\":\"TCP\",\"targetPort\":9104}],\"clusterIP\":\"None\",\"selector\":{\"app\":\"@master-slave\"}}}', 3);
INSERT INTO `workload_step` VALUES (19, 'pushgatewaylocaltime', 'volume', 'bound', '{\"hostPath\":{\"path\":\"/etc/localtime\"}}', 4);
INSERT INTO `workload_step` VALUES (20, 'pushgateway-dp', 'container', 'bound', '{\"image\":\"pushgateway:v1.3.1_01\",\"readinessProbe\":{\"initialDelaySeconds\":10,\"timeoutSeconds\":5,\"httpGet\":{\"path\":\"/#/status\",\"port\":9091}},\"livenessProbe\":{\"initialDelaySeconds\":10,\"timeoutSeconds\":5,\"httpGet\":{\"path\":\"/#/status\",\"port\":9091}},\"resources\":{\"requests\":{\"memory\":\"1Gi\",\"cpu\":\"1\"},\"limits\":{\"memory\":\"3Gi\",\"cpu\":\"3\"}},\"ports\":[{\"containerPort\":9091}],\"volumeMounts\":[{\"name\":\"pushgatewaylocaltime\",\"mountPath\":\"/etc/localtime\"}]}', 4);
INSERT INTO `workload_step` VALUES (21, 'pushgateway-svc', 'service', 'createorupdate', '{\"metadata\":{\"annotations\":{\"prometheus.io/scrape\":\"true\",\"prometheus.io/port\":\"9091\"}},\"spec\":{\"ports\":[{\"name\":\"pushgateway\",\"port\":9091,\"targetPort\":9091}],\"selector\":{\"app\":\"@singlenode\"}}}', 4);
INSERT INTO `workload_step` VALUES (22, 'prometheus-data', 'pvc', 'bound', '{\"spec\":{\"storageClassName\":\"local-storage-workload-prometheus\",\"accessModes\":[\"ReadWriteOnce\"],\"resources\":{\"requests\":{\"storage\":\"5G\"}}}}', 5);
INSERT INTO `workload_step` VALUES (23, 'prometheus-cm', 'conf', 'createorupdate', '{\"data\":{\"prometheus.yml\":\"# my global config\\nglobal:\\n  scrape_interval:     15s\\n  evaluation_interval: 15s\\n  # scrape_timeout is set to the global default (10s).\\n\\nalerting:\\n  alertmanagers:\\n  - static_configs:\\n    - targets:\\n      # - alertmanager:9093\\n\\nrule_files:\\n  # - \\\"first_rules.yml\\\"\\n  # - \\\"second_rules.yml\\\"\\n\\nscrape_configs:\\n  - job_name: \'prometheus\'\\n    static_configs:\\n    - targets: [\'localhost:9090\']\\n  - job_name: \'pushgateway\'\\n    static_configs:\\n    - targets: [\'localhost:9091\']\\n\"}}', 5);
INSERT INTO `workload_step` VALUES (24, 'prome-localtime', 'volume', 'bound', '{\"hostPath\":{\"path\":\"/etc/localtime\"}}', 5);
INSERT INTO `workload_step` VALUES (25, 'prome-cmvolume', 'volume', 'bound', '{\"configMap\":{\"name\":\"@prometheus-cm\"}}', 5);
INSERT INTO `workload_step` VALUES (26, 'prometheus-sts', 'container', 'bound', '{\"image\":\"prometheus:v2.23.0_01\",\"resources\":{\"requests\":{\"memory\":\"1Gi\",\"cpu\":\"1\"},\"limits\":{\"memory\":\"3Gi\",\"cpu\":\"3\"}},\"ports\":[{\"containerPort\":9090,\"name\":\"prom-server\"}],\"command\":[\"/bin/prometheus\"],\"args\":[\"--config.file=/etc/prometheus/prometheus.yml\",\"--storage.tsdb.path=/data/prometheus\",\"--storage.tsdb.retention=24h\",\"--web.enable-admin-api\",\"--web.enable-lifecycle\"],\"readinessProbe\":{\"httpGet\":{\"path\":\"/-/ready\",\"port\":9090},\"initialDelaySeconds\":10,\"timeoutSeconds\":5},\"livenessProbe\":{\"httpGet\":{\"path\":\"/-/healthy\",\"port\":9090},\"initialDelaySeconds\":10,\"timeoutSeconds\":5},\"volumeMounts\":[{\"name\":\"prome-localtime\",\"mountPath\":\"/etc/localtime\"},{\"name\":\"prometheus-data\",\"mountPath\":\"/data/prometheus\"},{\"name\":\"prome-cmvolume\",\"mountPath\":\"/etc/prometheus\"}],\"securityContext\":{\"runAsUser\":0,\"runAsGroup\":0}}', 5);
INSERT INTO `workload_step` VALUES (27, 'prometheus-svc', 'service', 'createorupdate', '{\"spec\":{\"ports\":[{\"port\":9090,\"name\":\"prometheus\",\"targetPort\":9090}],\"selector\":{\"app\":\"@master\"}}}', 5);
INSERT INTO `workload_step` VALUES (28, 'kafka-sts', 'container', 'bound', '{\"image\":\"kafka:1.1.1_01\",\"resources\":{\"requests\":{\"memory\":\"1Gi\",\"cpu\":\"1\"},\"limits\":{\"memory\":\"3Gi\",\"cpu\":\"3\"}},\"ports\":[{\"containerPort\":9092,\"name\":\"kafka-server\"}],\"env\":[{\"name\":\"KAFKA_HEAP_OPTS\",\"value\":\"-Xmx256M -Xms256M\"},{\"name\":\"KAFKA_OPTS\",\"value\":\"-Dlogging.level=INFO\"},{\"name\":\"POD_NAMESPACE\",\"valueFrom\":{\"fieldRef\":{\"fieldPath\":\"metadata.namespace\"}}}],\"command\":[\"bash\",\"-c\",\"exec kafka-server-start.sh /opt/dtstack/kafka/config/server.properties --override broker.id=${HOSTNAME##*-} --override zookeeper.connect=dtbase-zookeeper-leader-0.dtbase-zookeeper-leader-zk-hs.$POD_NAMESPACE.svc.cluster.local:2181,dtbase-zookeeper-leader-1.dtbase-zookeeper-leader-zk-hs.$POD_NAMESPACE.svc.cluster.local:2181,dtbase-zookeeper-leader-2.dtbase-zookeeper-leader-zk-hs.$POD_NAMESPACE.svc.cluster.local:2181/streamtest --override listeners=PLAINTEXT://:9092\"],\"securityContext\":{\"runAsUser\":1000,\"runAsGroup\":1000}}', 6);
INSERT INTO `workload_step` VALUES (29, 'kafka-svc', 'service', 'createorupdate', '{\"spec\":{\"ports\":[{\"port\":9090,\"name\":\"prometheus\",\"targetPort\":9090}],\"selector\":{\"app\":\"@master\"}}}', 6);
INSERT INTO `workload_step` VALUES (30, 'mysql-xtrabackup', 'container', 'bound', '{\"image\":\"xtrabackup:2.4.21_02\",\"env\":[{\"name\":\"POD_NAMESPACE\",\"valueFrom\":{\"fieldRef\":{\"fieldPath\":\"metadata.namespace\"}}}],\"ports\":[{\"name\":\"xtrabackup\",\"containerPort\":3307}],\"command\":[\"bash\",\"-c\",\"/usr/bin/dtstack/xtrabackup.sh\"],\"volumeMounts\":[{\"name\":\"mysql-data\",\"mountPath\":\"/var/lib/mysql\",\"subPath\":\"mysql\"},{\"name\":\"mysql-conf\",\"mountPath\":\"/etc/mysql/conf.d\"}],\"resources\":{\"requests\":{\"memory\":\"100Mi\",\"cpu\":\"100m\"}}}', 3);
INSERT INTO `workload_step` VALUES (31, 'mysql-exporter', 'container', 'bound', '{\"image\":\"mysqld-exporter:v0.12.1_01\",\"env\":[{\"name\":\"DATA_SOURCE_NAME\",\"value\":\"monitor:Abc123Admini@(127.0.0.1:3306)/\"}],\"ports\":[{\"protocol\":\"TCP\",\"containerPort\":9104}]}', 3);
COMMIT;

UPDATE `import_init_moudle` SET operator = '{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"name\":\"mole-operator\",\"namespace\":\"{{.NAME_SPACE}}\"},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"name\":\"mole-operator\"}},\"template\":{\"metadata\":{\"labels\":{\"name\":\"mole-operator\"}},\"spec\":{\"containers\":[{\"command\":[\"mole-operator\"],\"env\":[{\"name\":\"WATCH_NAMESPACE\",\"value\":\"{{.NAME_SPACE}}\"},{\"name\":\"LOG_MEM_LIMIT\",\"value\":\"500Mi\"},{\"name\":\"LOG_CPU_LIMIT\",\"value\":\"100m\"},{\"name\":\"LOG_MEM_REQUEST\",\"value\":\"200Mi\"},{\"name\":\"LOG_CPU_REQUEST\",\"value\":\"10m\"},{\"name\":\"POD_NAME\",\"valueFrom\":{\"fieldRef\":{\"fieldPath\":\"metadata.name\"}}},{\"name\":\"OPERATOR_NAME\",\"value\":\"mole-operator\"}],\"image\":\"{{.REGISTRY}}/mole:v4.1.2\",\"imagePullPolicy\":\"Always\",\"name\":\"mole-operator\",\"resources\":{\"limits\":{\"cpu\":\"500m\",\"memory\":\"500Mi\"}}}],\"imagePullSecrets\":[{\"name\":\"{{.SECRET_NAME}}\"}]}}}}';

ALTER TABLE import_init_moudle ADD COLUMN `log_config` varchar(2000) DEFAULT " " AFTER `operator`;

UPDATE `import_init_moudle` SET log_config = '{\"apiVersion\":\"v1\",\"data\":{\"filebeat.yml\":\"logging.level: debug\\nfilebeat.inputs:\\n  - type: log\\n    tail_files: true\\n    index: \\\"log-${PRODUCT}-%{+yyyy-MM-dd}\\\"\\n    fields:\\n      namespace: ${NAMESPACE}\\n      serviceAccountName: ${SERVICE_ACCOUNT_NAME}\\n      product: ${PRODUCT}\\n      job: ${JOB}\\n      node: ${HOSTNAME}/${HOST_IP}\\n      pod_name: ${POD_NAME}\\n      pod_uid: ${POD_UID}\\n      pod_ip: ${POD_IP}\\n    tags: [ \\\"${PRODUCT}\\\",\\\"${JOB}\\\" ]\\n    paths: ${LOG_PATH}\\noutput.elasticsearch:\\n  hosts: [ \\\"${LOG_SERVER_ADDRESS}\\\" ]\\n  username: \\\"elastic\\\"\\n  password: \\\"dtstack\\\"\\n\",\"promtail.yaml\":\"client:\\n  backoff_config:\\n    max_period: 5m\\n    max_retries: 10\\n    min_period: 500ms\\n  batchsize: 1048576\\n  batchwait: 1s\\n  external_labels: {}\\n  timeout: 10s\\npositions:\\n  filename: /run/positions.yaml\\nserver:\\n  http_listen_port: 3101\\ntarget_config:\\n  sync_period: 10s\\nscrape_configs:\\n  - job_name: test\\n    static_configs:\\n      - labels:\\n          namespace: ${NAMESPACE}\\n          serviceAccountName: ${SERVICE_ACCOUNT_NAME}\\n          product: ${PRODUCT}\\n          job: ${JOB}\\n          node: ${HOSTNAME}/${HOST_IP}\\n          pod_name: ${POD_NAME}\\n          pod_uid: ${POD_UID}\\n          pod_ip: ${POD_IP}\\n          __path__: ${LOG_PATH}\\n\"},\"kind\":\"ConfigMap\",\"metadata\":{\"name\":\"log-configmap\",\"namespace\":\"{{.NAME_SPACE}}\"}}';

CREATE TABLE IF NOT EXISTS `deploy_instance_update_record` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `update_uuid` char(36) NOT NULL DEFAULT '' COMMENT '更新记录uuid',
    `instance_id` int(11) unsigned NOT NULL COMMENT '实例id',
    `sid` char(36) NOT NULL DEFAULT '' COMMENT '主机 ID (UUID)',
    `ip` varchar(255) NOT NULL DEFAULT '' COMMENT '主机ip',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '组件名称',
    `product_name_display` varchar(255) NOT NULL DEFAULT '' COMMENT '组件显示名称',
    `product_version` varchar(255) NOT NULL DEFAULT '' COMMENT '产品版本',
    `group` varchar(255) NOT NULL DEFAULT '' COMMENT '组名称',
    `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称',
    `service_name_display` varchar(255) NOT NULL DEFAULT '' COMMENT '服务显示名称',
    `service_version` varchar(255) NOT NULL DEFAULT '' COMMENT '服务版本',
    `status` varchar(32) NOT NULL DEFAULT '' COMMENT '实例状态',
    `status_message` varchar(1024) NOT NULL DEFAULT '' COMMENT '状态详细信息',
    `progress` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '进度',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `deploy_uuid` (`update_uuid`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1998  COMMENT='实例部署记录表';

CREATE TABLE IF NOT EXISTS `deploy_product_update_history` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id',
    `namespace` varchar(255) NOT NULL COMMENT 'cluster namespace',
    `update_uuid` char(36) NOT NULL COMMENT '更新uuid',
    `parent_product_name` varchar(255) NOT NULL COMMENT '产品名称',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '组件名称',
    `product_name_display` varchar(255) NOT NULL DEFAULT '' COMMENT '组件显示名称',
    `product_version` varchar(255) NOT NULL DEFAULT '' COMMENT '产品版本',
    `product_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT ',0,, 1, k8s',
    `status` char(32) NOT NULL DEFAULT '' COMMENT '产品状态',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create_time',
    `update_start_time` datetime DEFAULT NULL COMMENT 'deploy_start_time',
    `update_end_time` datetime DEFAULT NULL COMMENT 'deploy_end_time',
    `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '部署人id',
    `package_name` varchar(255) NOT NULL COMMENT '补丁包名称',
    `update_dir` varchar(255) NOT NULL COMMENT '目标目录',
    `backup_dir` varchar(255) NOT NULL COMMENT '备份目录',
    `product_id` int(11) NOT NULL COMMENT '产品包id',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=14 ;

alter table deploy_schema_field_modify modify field_path varchar(255) binary default '' not null comment '字段路径';

CREATE TABLE IF NOT EXISTS `inspect_report_template` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `cluster_id` int(11) NOT NULL COMMENT '集群id',
    `type` tinyint(1) NOT NULL COMMENT '图表所属类型，0为节点，1为应用',
    `module` varchar(255) NOT NULL COMMENT '模块',
    `metric` varchar(255) NOT NULL COMMENT '指标',
    `targets` text NOT NULL COMMENT '采集数据配置，数组类型，包含查询语句以及维度拼接模板',
    `unit` varchar(16) DEFAULT '' COMMENT 'y轴单位',
    `decimal` int(11) DEFAULT 0 COMMENT 'y轴小数点保留位数',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
    PRIMARY KEY(`id`)
    ) COMMENT '巡检报告图表自定义模板' ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_schema_multi_fields` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `cluster_id` int(11) NOT NULL COMMENT '集群id',
    `product_name` varchar(50) NOT NULL DEFAULT '' COMMENT '产品名称',
    `service_name` varchar(50) NOT NULL DEFAULT '' COMMENT '服务名称',
    `field_path` varchar(255) NOT NULL DEFAULT '' COMMENT '字段路径',
    `field` text NOT NULL COMMENT '字段值',
    `hosts` varchar(255) NOT NULL COMMENT '配置绑定主机ip，逗号连接',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY(`id`),
    UNIQUE KEY `namesz_path` (`product_name`, `service_name`, `field_path`, `hosts`)
    ) COMMENT 'schema配置多值记录' ENGINE=InnoDB ;

CREATE TABLE IF NOT EXISTS `deploy_notify_event` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `cluster_id` int(11) NOT NULL COMMENT '集群id',
    `type` tinyint(1) NOT NULL COMMENT '事件类型，0表示配置变更依赖组件重启',
    `product_name` varchar(255) DEFAULT NULL COMMENT '产品名称',
    `service_name` varchar(255) DEFAULT NULL COMMENT '服务名称',
    `depend_product_name` varchar(255) DEFAULT NULL COMMENT '被依赖的产品名称',
    `depend_service_name` varchar(255) DEFAULT NULL COMMENT '被依赖的服务名称',
    `host` varchar(1024) DEFAULT NULL COMMENT '主机ip',
    `is_deleted` tinyint(1) DEFAULT 0 COMMENT '是否删除',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY(`id`)
    ) COMMENT '通知事件' ENGINE=InnoDB ;

INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (1, 1, 0, 'System', 'System Load', '[{\"expr\":\"node_load1{}\",\"legend_format\":\"Load1-{{instance}}\"}]', '', 2, '2021-03-30 11:32:29', '2021-03-30 11:32:29', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (2, 1, 0, 'System', 'Disk Used(%)', '[{\"expr\":\"100-100*((node_filesystem_avail{device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"}/node_filesystem_size{device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"}))\",\"legend_format\":\"DiskUsed-{{instance}}{{mountpoint}}\"}]', '%', 0, '2021-03-30 13:56:00', '2021-03-30 13:56:00', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (3, 1, 0, 'System', 'Inode Used(%)', '[{\"expr\":\"100*((node_filesystem_files{device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"}-node_filesystem_files_free{device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"})/node_filesystem_files{device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"})\",\"legend_format\":\"InodeUsed-{{instance}}{{mountpoint}}\"}]', '%', 3, '2021-03-30 14:02:11', '2021-03-30 14:02:11', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (4, 1, 0, 'System', 'Mem Used(%)', '[{\"expr\":\"(1 - (sum by (instance) (node_memory_MemFree + node_memory_Buffers + node_memory_Cached)) / (sum by (instance) (node_memory_MemTotal)) ) * 100\",\"legend_format\":\"MemUsed-{{instance}}\"}]', '%', 0, '2021-03-30 14:02:52', '2021-03-30 14:02:52', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (5, 1, 0, 'System', 'Node_up', '[{\"expr\":\"up{job=\\\"node_exporter\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-03-30 14:07:10', '2021-03-30 14:07:10', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (6, 1, 0, 'MySQL', 'MySQL down-监控', '[{\"expr\":\"up{product_name=\\\"DTBase\\\",service_name=\\\"mysql\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-04-01 16:34:25', '2021-04-01 16:34:25', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (7, 1, 0, 'MySQL', 'MySQL Connection', '[{\"expr\":\"max(max_over_time(mysql_global_status_threads_connected{product_name=\\\"DTBase\\\",service_name=\\\"mysql\\\"}[1m])or mysql_global_status_threads_connected{product_name=\\\"DTBase\\\",service_name=\\\"mysql\\\"})\",\"legend_format\":\"Connections\"},{\"expr\":\"mysql_global_status_max_used_connections{product_name=\\\"DTBase\\\",service_name=\\\"mysql\\\"}\",\"legend_format\":\"MaxUsedConnections\"},{\"expr\":\"mysql_global_variables_max_connections{product_name=\\\"DTBase\\\",service_name=\\\"mysql\\\"}\",\"legend_format\":\"MaxConnections\"}]', '', 0, '2021-04-01 16:36:42', '2021-04-01 16:36:42', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (8, 1, 0, 'MySQL', 'MySQL Slow Queries', '[{\"expr\":\"rate(mysql_global_status_slow_queries{product_name=\\\"DTBase\\\",service_name=\\\"mysql\\\"}[1m])or irate(mysql_global_status_slow_queries{product_name=\\\"DTBase\\\",service_name=\\\"mysql\\\"}[5m])\",\"legend_format\":\"{{instance}}-SlowQueries\"}]', '', 1, '2021-04-01 17:06:57', '2021-04-01 17:06:57', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (9, 1, 0, 'Zookeeper', 'Zookeeper up', '[{\"expr\":\"up{product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-04-01 17:14:23', '2021-04-01 17:14:23', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (10, 1, 0, 'Zookeeper', 'Old Gen Used(%)', '[{\"expr\":\"100*jvm_memory_pool_bytes_used{pool=\\\"PS Old Gen\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"}/jvm_memory_pool_bytes_committed{pool=\\\"PS Old Gen\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"}\",\"legend_format\":\"OldGenUsed-{{instance}}\"}]', '%', 2, '2021-04-01 17:28:24', '2021-04-01 17:28:24', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (11, 1, 0, 'Zookeeper', 'Heap Used(%)', '[{\"expr\":\"100*jvm_memory_bytes_used{area=\\\"heap\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"}/jvm_memory_bytes_committed{area=\\\"heap\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"}\",\"legend_format\":\"HeapUsed-{{instance}}\"}]', '%', 2, '2021-04-01 17:45:34', '2021-04-01 17:45:34', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (12, 1, 0, 'Zookeeper', 'Full GC Count（2minutes）', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"}[2m]))\",\"legend_format\":\"FullGCCount-{{instance}}\"}]', '', 0, '2021-04-01 18:00:22', '2021-04-01 18:00:22', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (13, 1, 0, 'Zookeeper', 'Full GC Avg Time', '[{\"expr\":\"(jvm_gc_collection_seconds_sum{gc=\\\"PS MarkSweep\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"} > 0 or on() vector(0.001)) or (jvm_gc_collection_seconds_sum{gc=\\\"PS MarkSweep\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"} > bool 0)/(jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"} > 0 or on() vector(1)) or (jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"} > bool 0)\",\"legend_format\":\"FullGCAvgTime-{{instance}}\"}]', '', 3, '2021-04-01 18:04:04', '2021-04-01 18:04:04', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (14, 1, 0, 'Zookeeper', 'Full GC Count', '[{\"expr\":\"jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"DTBase\\\",service_name=\\\"zookeeper\\\"}\",\"legend_format\":\"FullGCCount-{{instance}}\"}]', '', 2, '2021-04-01 18:05:20', '2021-04-01 18:05:20', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (15, 1, 0, 'Redis', 'Redis_Up', '[{\"expr\":\"up{product_name=\\\"DTBase\\\",service_name=\\\"redis\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-04-01 19:25:28', '2021-04-01 19:25:28', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (16, 1, 0, 'NameNode', 'NameNode_up', '[{\"expr\":\"up{product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-04-02 15:08:07', '2021-04-02 15:08:07', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (17, 1, 0, 'NameNode', 'Old Gen Used(%)', '[{\"expr\":\"100*jvm_memory_pool_bytes_used{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}/jvm_memory_pool_bytes_committed{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"OldGenUsed-{{instance}}\"}]', '%', 2, '2021-04-02 15:10:27', '2021-04-02 15:10:27', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (18, 1, 0, 'NameNode', 'Heap Used(%)', '[{\"expr\":\"100*jvm_memory_bytes_used{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}/jvm_memory_bytes_committed{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"HeapUsed-{{instance}}\"}]', '%', 2, '2021-04-02 15:11:55', '2021-04-02 15:11:55', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (19, 1, 0, 'NameNode', 'Full GC Count(2 minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}[2m]))\",\"legend_format\":\"FullGCCount-{{instance}}\"}]', '', 0, '2021-04-02 15:16:06', '2021-04-02 15:16:06', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (20, 1, 0, 'NameNode', 'Full GC Avg Time', '[{\"expr\":\"jvm_gc_collection_seconds_sum{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}/jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"FullGCAvgTime-{{instance}}\"}]', '', 3, '2021-04-02 15:17:27', '2021-04-02 15:17:27', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (21, 1, 0, 'NameNode', 'Total/Used', '[{\"expr\":\"Hadoop_NameNode_Total{name=\\\"NameNodeInfo\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"Total-{{instance}}\"},{\"expr\":\"Hadoop_NameNode_Used{name=\\\"NameNodeInfo\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"Used-{{instance}}\"}]', '', 3, '2021-04-02 15:31:46', '2021-04-02 15:31:46', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (22, 1, 0, 'NameNode', 'TotalBlocks', '[{\"expr\":\"Hadoop_NameNode_TotalBlocks{name=\\\"NameNodeInfo\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"TotalBlocks-{{instance}}\"}]', '', 0, '2021-04-02 15:32:08', '2021-04-02 15:32:08', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (23, 1, 0, 'NameNode', 'Number Of MissingBlocks', '[{\"expr\":\"Hadoop_NameNode_NumberOfMissingBlocks{name=\\\"NameNodeInfo\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"NumberOfMissingBlocks-{{instance}}\"}]', '', 0, '2021-04-02 15:34:20', '2021-04-02 15:34:20', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (24, 1, 0, 'NameNode', 'TotalFiles', '[{\"expr\":\"Hadoop_NameNode_TotalFiles{name=\\\"FSNamesystem\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"TotalFiles-{{instance}}\"}]', '', 0, '2021-04-02 15:38:10', '2021-04-02 15:38:10', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (25, 1, 0, 'NameNode', 'CorruptBlocks', '[{\"expr\":\"Hadoop_NameNode_CorruptBlocks{name=\\\"FSNamesystem\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"CorruptBlocks-{{instance}}\"}]', '', 0, '2021-04-02 15:43:24', '2021-04-02 15:43:24', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (26, 1, 0, 'DataNode', 'DataNode_up', '[{\"expr\":\"up{product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-04-02 15:48:02', '2021-04-02 15:48:02', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (27, 1, 0, 'DataNode', 'Old Gen Used(%)', '[{\"expr\":\"100*jvm_memory_pool_bytes_used{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}/jvm_memory_pool_bytes_committed{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}\",\"legend_format\":\"OldGenUsed-{{instance}}\"}]', '%', 2, '2021-04-02 15:48:19', '2021-04-02 15:48:19', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (28, 1, 0, 'DataNode', 'Heap Used(%)', '[{\"expr\":\"100*jvm_memory_bytes_used{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}/jvm_memory_bytes_committed{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}\",\"legend_format\":\"HeapUsed-{{instance}}\"}]', '%', 2, '2021-04-02 15:49:53', '2021-04-02 15:49:53', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (29, 1, 0, 'DataNode', 'Full GC Count(2 minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}[2m]))\",\"legend_format\":\"FullGCCount-{{instance}}\"}]', '', 0, '2021-04-02 15:58:04', '2021-04-02 15:58:04', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (30, 1, 0, 'DataNode', 'Full GC Avg Time', '[{\"expr\":\"jvm_gc_collection_seconds_sum{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}/jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}\",\"legend_format\":\"FullGCAvgTime-{{instance}}\"}]', '', 3, '2021-04-02 16:03:55', '2021-04-02 16:03:55', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (31, 1, 0, 'ResourceManager', 'ResourceManager_up', '[{\"expr\":\"up{product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-04-02 17:07:18', '2021-04-02 17:07:18', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (32, 1, 0, 'ResourceManager', 'Old Gen Used(%)', '[{\"expr\":\"100*jvm_memory_pool_bytes_used{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}/jvm_memory_pool_bytes_committed{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}\",\"legend_format\":\"OldGenUsed-{{instance}}\"}]', '%', 2, '2021-04-02 17:10:48', '2021-04-02 17:10:48', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (33, 1, 0, 'ResourceManager', 'Heap Used(%)', '[{\"expr\":\"100*jvm_memory_bytes_used{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}/jvm_memory_bytes_committed{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}\",\"legend_format\":\"HeapUsed-{{instance}}\"}]', '%', 2, '2021-04-02 17:15:26', '2021-04-02 17:15:26', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (34, 1, 0, 'ResourceManager', 'Full GC Count(2 minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}[2m]))\",\"legend_format\":\"FullGCCount-{{instance}}\"}]', '', 0, '2021-04-13 17:04:11', '2021-04-13 17:04:11', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (35, 1, 0, 'ResourceManager', 'Full GC AvgTime', '[{\"expr\":\"jvm_gc_collection_seconds_sum{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}/jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}\",\"legend_format\":\"FullGCTimeAvg-{{instance}}\"}]', '', 3, '2021-04-13 17:07:03', '2021-04-13 17:07:03', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (36, 1, 0, 'ResourceManager', 'NumActiveNMs', '[{\"expr\":\"Hadoop_ResourceManager_NumActiveNMs{name=\\\"ClusterMetrics\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}\",\"legend_format\":\"NumActiveNMs-{{instance}}\"}]', '', 0, '2021-04-13 17:09:48', '2021-04-13 17:09:48', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (37, 1, 0, 'ResourceManager', 'NumUnhealthyNMs', '[{\"expr\":\"Hadoop_ResourceManager_NumUnhealthyNMs{name=\\\"ClusterMetrics\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}\",\"legend_format\":\"NumUnhealthyeNMs-{{instance}}\"}]', '', 0, '2021-04-13 17:12:18', '2021-04-13 17:12:18', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (38, 1, 0, 'ResourceManager', 'NumLostNMs', '[{\"expr\":\"Hadoop_ResourceManager_NumLostNMs{name=\\\"ClusterMetrics\\\",product_name=\\\"Hadoop\\\",service_name=\\\"yarn_resourcemanager\\\"}\",\"legend_format\":\"NumLostNMs-{{instance}}\"}]', '', 0, '2021-04-13 17:14:01', '2021-04-13 17:14:01', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (39, 1, 0, 'ThriftServer', 'Spark_ThriftServer_Up', '[{\"expr\":\"up{product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}\",\"legend_format\":\"{{instance}}\"}]', '', 0, '2021-04-07 10:02:53', '2021-04-07 10:02:53', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (40, 1, 0, 'ThriftServer', 'Old Gen Used(%)', '[{\"expr\":\"100*jvm_memory_pool_bytes_used{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}/jvm_memory_pool_bytes_committed{pool=\\\"PS Old Gen\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}\",\"legend_format\":\"OldGenUsed-{{instance}}\"}]', '%', 2, '2021-04-07 10:07:29', '2021-04-07 10:07:29', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (41, 1, 0, 'ThriftServer', 'Heap Used(%)', '[{\"expr\":\"100*jvm_memory_bytes_used{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}/jvm_memory_bytes_committed{area=\\\"heap\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}\",\"legend_format\":\"HeapUsed-{{instance}}\"}]', '%', 2, '2021-04-07 10:07:56', '2021-04-07 10:07:56', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (42, 1, 0, 'ThriftServer', 'Full GC Count(2 minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}[2m]))\",\"legend_format\":\"FullGCCount-{{instance}}\"}]', '', 0, '2021-04-07 10:08:30', '2021-04-07 10:08:30', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (43, 1, 0, 'ThriftServer', 'Full GC Avg Time', '[{\"expr\":\"jvm_gc_collection_seconds_sum{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}/jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}\",\"legend_format\":\"FullGCAvgTime-{{instance}}\"}]', 's', 3, '2021-04-07 10:08:57', '2021-04-07 10:08:57', 0);
INSERT INTO inspect_report_template(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES (44, 1, 0, 'ThriftServer', 'Full GC Count(total)', '[{\"expr\":\"jvm_gc_collection_seconds_count{gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"thriftserver\\\"}\",\"legend_format\":\"FullGCCount-{{instance}}\"}]', '', 0, '2021-04-07 10:09:19', '2021-04-07 10:09:19', 0);


-- 主机编排
alter table deploy_host
    add role_list varchar(255) null comment '角色列表';

create table host_role
(
    id         int auto_increment primary key,
    cluster_id int          null comment '集群 id',
    role_name  varchar(255) collate utf8_bin not null comment '角色名',
    role_type  int          not null comment '角色类型 1 默认 2 自定义',
    constraint host_role_list_cluster_id_role_name_uindex
        unique (cluster_id, role_name)
);
create table if not exists deploy_product_select_history
(
    cluster_id int          not null primary key,
    pid_list   varchar(255) null comment 'pid 清单'
    );

create table if not exists deploy_uuid
(
    id          int auto_increment primary key,
    uuid        varchar(255)                       not null,
    type        int                                not null comment '部署类型： 1 手动；2 自动；3 自动部署中的子产品的 uuid',
    parent_uuid varchar(255)                       null comment '自动部署 uuid',
    pid         varchar(255)                       null comment '产品包 id',
    create_time datetime default CURRENT_TIMESTAMP not null
    );


-- 执行完上面的建表语句后
--  用查询出的每个 id 替换以下插入语句中的每个 cluster_id_value 然后执行 insert 语句
insert into host_role (cluster_id, role_name, role_type) values (1,'web',1),(1,'manager',1),(1,'worker',1);

-- deploy_unchecked_service增加namespace字段和索引
ALTER TABLE `dtagent`.`deploy_unchecked_service`
    ADD COLUMN `namespace` varchar(255) NOT NULL COMMENT 'k8s多命名空间' AFTER `create_time`;

ALTER TABLE `dtagent`.`deploy_unchecked_service`
DROP INDEX `cluster_pid`,
ADD UNIQUE INDEX `cluster_pid`(`cluster_id`, `pid`, `namespace`) USING BTREE;


-- deploy_schema_field_modify增加namespace字段和索引
ALTER TABLE `dtagent`.`deploy_schema_field_modify`
    ADD COLUMN `namespace` varchar(255) NOT NULL COMMENT 'k8s多命名空间区分' AFTER `create_time`;

ALTER TABLE `dtagent`.`deploy_schema_field_modify`
DROP INDEX `cluster_names_path`,
ADD UNIQUE INDEX `cluster_names_path`(`cluster_id`, `product_name`, `service_name`, `field_path`, `namespace`) USING BTREE;

-- big image
BEGIN;
INSERT INTO `workload_definition` VALUES (7, 'plugin', 'v1', '[{\"key\":\"Image\",\"ref\":\"spec.workloadpatrs.0.steps.0.object.image\"}]', 1);
INSERT INTO `workload_part` VALUES (7, 'image-push', 'job', '{}', 7);
INSERT INTO `workload_step` VALUES (32, 'image', 'container', 'bound', '{\"image\":\"\",\"command\":[\"/bin/sleep\"],\"args\":[\"5\"],\"resources\":{\"limits\":{\"cpu\":\"50m\",\"memory\":\"10Mi\"},\"requests\":{\"cpu\":\"0m\",\"memory\":\"0Mi\"}}}', 7);
COMMIT;

ALTER TABLE deploy_notify_event ADD COLUMN product_stopped TINYINT DEFAULT '0' COMMENT '组件是否停止';
CREATE TABLE `deploy_switch_record` (
                                        `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
                                        `name` varchar(64) NOT NULL COMMENT '开关名称',
                                        `product_name` varchar(255) NOT NULL COMMENT '产品名称',
                                        `service_name` varchar(255) NOT NULL COMMENT '服务名称',
                                        `status` varchar(64) NOT NULL COMMENT '状态',
                                        `status_message` varchar(1024) DEFAULT NULL COMMENT '状态详细信息',
                                        `progress` tinyint(3) DEFAULT NULL COMMENT '进度',
                                        `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                        `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                                        `is_deleted` tinyint(1) DEFAULT '0' COMMENT '是否删除',
                                        `cluster_id` int(11) NOT NULL COMMENT '集群id',
                                        `switch_type` varchar(16) NOT NULL COMMENT '开关操作类型，on/off',
                                        PRIMARY KEY (`id`)
) ENGINE=InnoDB  COMMENT='开关记录表';

ALTER TABLE deploy_cluster_product_rel modify product_parsed LongText not null;

INSERT INTO `workload_step`(`name`, `type`, `action`, `object`, `workloadpart_id`) VALUES ('cs', 'service', 'createorupdate', '{\"spec\":{\"ports\":[{\"port\":3306,\"name\":\"mysql\"},{\"name\":\"mysql-exporter\",\"port\":9104,\"protocol\":\"TCP\",\"targetPort\":9104}],\"selector\":{\"app\":\"@master-slave\"}}}', 3);

-- 修复大字段
alter table deploy_instance_list modify column `schema` longtext not null;

-- 4.1.8
ALTER TABLE `dtagent`.`deploy_kube_base_product_list`
    CHANGE COLUMN `base_clusterId` `rely_namespace` varchar(255)  COLLATE utf8mb4_bin NOT NULL COMMENT '依赖namespace' AFTER `namespace`;

create table if not exists operation_list
(
    id               int auto_increment
    primary key,
    cluster_id       int                                 not null comment '集群 id',
    operation_id     varchar(255)                        null comment '操作 id',
    operation_type   int                                 not null comment '1. 产品包部署 2. 产品包升级 3. 产品包启动 4. 服务启动 5. 服务滚动重启 6. 主机初始化 7. Kerberos开启 8. Kerberos关闭',
    operation_status int                                 not null comment '1 进行中 2 正常 3 失败',
    object_type      int                                 not null comment '对象类型 1：产品包 2：服务 3：主机',
    object_value     varchar(255)                        not null comment '对象值，用于页面回显',
    create_time      timestamp default CURRENT_TIMESTAMP null,
    end_time         timestamp                           null comment '结束时间',
    duration         float                               null comment '持续时间 单位秒',
    update_time      timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint operation_list_operation_id_uindex
    unique (operation_id)
    )
    comment '操作清单表' ;


create table if not exists exec_shell_list
(
    id           int auto_increment
    primary key,
    cluster_id   int                                 not null comment '集群 id',
    exec_id      varchar(255)                        not null comment 'shell 执行 id ',
    operation_id varchar(255)                        not null comment '操作id',
    shell_type   int                                 not null comment '具体shell 类型 1 服务安装 2 服务启动 3 执行脚本',
    product_name varchar(255)                        null comment '所属产品包',
    service_name varchar(255)                        null comment '所属服务',
    sid          varchar(255)                        null comment '主机',
    seq          int                                 null comment 'exec seq',
    exec_status  int                                 null comment '1 进行中 2 正常 3 失败',
    create_time  timestamp default CURRENT_TIMESTAMP null,
    end_time     timestamp                           null comment '结束时间',
    duration     float                               null comment '持续时间 单位 秒',
    update_time  timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint exec_shell_list_exec_id_uindex
    unique (exec_id)
    )
    comment 'shell 执行记录表' ;

CREATE TABLE if not exists `deploy_inspect_report` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `name` varchar(64) NOT NULL COMMENT '巡检报告名称',
  `status` varchar(64) NOT NULL COMMENT '状态',
  `progress` tinyint(3) DEFAULT NULL COMMENT '进度',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint(1) DEFAULT '0' COMMENT '是否删除',
  `cluster_id` int(11) NOT NULL COMMENT '集群id',
  `file_path` varchar(255) DEFAULT NULL COMMENT '文件路径',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='巡检报告表';

alter table exec_shell_list add host_ip varchar(200) null comment 'ip' after sid;

CREATE TABLE if not exists `deploy_upload_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `upload_type` varchar(16) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '上传类型',
  `name` varchar(2048) COLLATE utf8mb4_bin NOT NULL COMMENT '链接地址',
  `progress` decimal(10,0) DEFAULT '0' COMMENT '进度',
  `status` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '状态',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `is_deleted` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

UPDATE `dtagent`.`import_init_moudle` SET `log_config` = '{\"apiVersion\":\"v1\",\"data\":{\"filebeat.yml\":\"logging.level: debug\\nfilebeat.inputs:\\n  - type: log\\n    tail_files: true\\n    index: \\\"log-${PRODUCT}-%{+yyyy-MM-dd}\\\"\\n    fields:\\n      namespace: ${NAMESPACE}\\n      serviceAccountName: ${SERVICE_ACCOUNT_NAME}\\n      product: ${PRODUCT}\\n      job: ${JOB}\\n      node: ${HOSTNAME}/${HOST_IP}\\n      pod_name: ${POD_NAME}\\n      pod_uid: ${POD_UID}\\n      pod_ip: ${POD_IP}\\n    tags: [ \\\"${PRODUCT}\\\",\\\"${JOB}\\\" ]\\n    paths: ${LOG_PATH}\\noutput.elasticsearch:\\n  hosts: [ \\\"${LOG_SERVER_ADDRESS}\\\" ]\\n  username: \\\"elastic\\\"\\n  password: \\\"dtstack\\\"\\n\",\"promtail.yaml\":\"client:\\n  backoff_config:\\n    max_period: 5m\\n    max_retries: 10\\n    min_period: 500ms\\n  batchsize: 1048576\\n  batchwait: 1s\\n  external_labels: {}\\n  timeout: 10s\\npositions:\\n  filename: /var/log/logs/positions.yaml\\nserver:\\n  http_listen_port: 3101\\ntarget_config:\\n  sync_period: 10s\\nscrape_configs:\\n  - job_name: test\\n    static_configs:\\n      - labels:\\n          namespace: ${NAMESPACE}\\n          serviceAccountName: ${SERVICE_ACCOUNT_NAME}\\n          product: ${PRODUCT}\\n          job: ${JOB}\\n          node: ${HOSTNAME}/${HOST_IP}\\n          pod_name: ${POD_NAME}\\n          pod_uid: ${POD_UID}\\n          pod_ip: ${POD_IP}\\n          __path__: ${LOG_PATH}\\n\"},\"kind\":\"ConfigMap\",\"metadata\":{\"name\":\"log-configmap\",\"namespace\":\"{{.NAME_SPACE}}\"}}' WHERE `id` = 1;
update workload_step set name="hs" where name="mysql-svc";

INSERT INTO `workload_step`(`name`, `type`, `action`, `object`, `workloadpart_id`) VALUES ('cs', 'service', 'createorupdate', '{\"spec\":{\"ports\":[{\"port\":3306,\"name\":\"mysql\"},{\"name\":\"mysql-exporter\",\"port\":9104,\"protocol\":\"TCP\",\"targetPort\":9104}],\"selector\":{\"app\":\"@master-slave\"}}}', 3);

CREATE TABLE IF NOT EXISTS `smoke_testing` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `cluster_id` int(11) NOT NULL COMMENT '集群id',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '产品名称',
    `operation_id` varchar(255) NOT NULL COMMENT '操作id',
    `exec_status` int NULL COMMENT '1 进行中,2 正常,3 失败',
    `report_url` varchar(255) NOT NULL COMMENT '报告地址',
    `create_time` timestamp DEFAULT CURRENT_TIMESTAMP NULL  COMMENT '开始时间',
    `end_time` timestamp NULL COMMENT '结束时间',
    PRIMARY KEY(`id`)
    ) COMMENT '冒烟测试记录表' ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `service_health_check` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `cluster_id` int(11) NOT NULL COMMENT '集群id',
    `product_name` varchar(255) NOT NULL COMMENT '产品名称',
    `pid` int(11) NOT NULL COMMENT 'product id',
    `service_name` varchar(255) NOT NULL COMMENT '服务名称',
    `agent_id` varchar(50) NOT NULL COMMENT 'agent id',
    `sid` varchar(50) NOT NULL COMMENT '主机id',
    `ip` varchar(50) NOT NULL COMMENT '主机ip',
    `script_name` varchar(255) NOT NULL COMMENT '脚本名称',
    `script_name_display` varchar(255) NOT NULL COMMENT '脚本显示名称',
    `auto_exec` tinyint(4) NOT NULL COMMENT '自动执行开关状态',
    `period` varchar(10) NOT NULL COMMENT '执行间隔时间',
    `retries` int(11) DEFAULT NULL COMMENT '执行重试次数',
    `exec_status` int(11) DEFAULT NULL COMMENT '0 未就绪,1 进行中,2 正常,3 失败',
    `error_message` varchar(1000) NOT NULL COMMENT '执行失败的错误信息',
    `start_time` timestamp NULL DEFAULT NULL COMMENT '执行开始时间',
    `end_time` timestamp NULL DEFAULT NULL COMMENT '执行结束时间',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='健康检查记录表';

ALTER TABLE `smoke_testing` ADD `test_script` VARCHAR(255) NOT NULL COMMENT '测试脚本' AFTER `operation_id`;
ALTER TABLE `smoke_testing` ADD `error_message` text NOT NULL COMMENT '错误信息' AFTER `report_url`;

CREATE TABLE `user_cluster_right` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `cluster_id` int(11) NOT NULL COMMENT '集群id',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0正常 1逻辑删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户集群权限关联表';

CREATE TABLE `product_backup_config` (
 `id` INT(11) NOT NULL AUTO_INCREMENT,
 `cluster_id` VARCHAR(255) COLLATE UTF8MB4_BIN NOT NULL COMMENT '集群id',
 `config_path` VARCHAR(255) COLLATE UTF8MB4_BIN NOT NULL COMMENT '备份路径',
 `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
 PRIMARY KEY (id)
) ENGINE=INNODB DEFAULT CHARSET=UTF8MB4 COLLATE = UTF8MB4_BIN COMMENT='组件备份路径记录表';

-- 平台配置
CREATE TABLE `sys_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `config_name` varchar(45) COLLATE utf8mb4_bin NOT NULL COMMENT '配置名称',
  `config_value` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '配置值',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0: 使用中;1: 已删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
-- 平台安全配置.密码传输加密方式
INSERT INTO `sys_config`(config_name,config_value) VALUES ('platformsecurity.login_encrypt','rsa');
-- 平台安全配置.强制用户修改初始密码
-- 0 关闭(default)，1 开启
INSERT INTO `sys_config`(config_name,config_value) VALUES ('platformsecurity.force_reset_password','0');
-- 平台安全配置.密码出错锁定开关
-- 0 关闭(default)，1 开启
INSERT INTO `sys_config`(config_name,config_value) VALUES ('platformsecurity.account_login_lock_switch','0');
-- 平台安全配置.密码出错次数
-- 3(default)
INSERT INTO `sys_config`(config_name,config_value) VALUES ('platformsecurity.account_login_limit_error','3');
-- 平台安全配置.锁定时长
-- 5(default)
INSERT INTO `sys_config`(config_name,config_value) VALUES ('platformsecurity.account_login_lock_time','5');
-- 平台安全配置.自动登出时长
-- 1440(default)
INSERT INTO `sys_config`(config_name,config_value) VALUES ('platformsecurity.account_logout_sleep_time','1440');
-- 用户.初始密码修改状态
ALTER TABLE `user_list` ADD COLUMN `reset_password_status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '初始密码修改状态：0 未修改，1 已修改' AFTER `status`;
UPDATE `user_list` SET `reset_password_status` = 1 WHERE 1 = 1;

CREATE TABLE `task_list` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `name` varchar(64) NOT NULL COMMENT '名称',
 `describe` varchar(256) NOT NULL COMMENT '描述',
 `spec` varchar(64) NOT NULL COMMENT 'cron表达式',
 `status` int NULL DEFAULT '0' COMMENT '定时状态: 0 关闭,1 开启',
 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
 `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0正常 1逻辑删除',
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='任务列表';

CREATE TABLE `task_host` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `task_id` int(11) NOT NULL COMMENT '任务id',
 `host_id` int(11) NOT NULL COMMENT '主机id',
 PRIMARY KEY (`id`),
 KEY `IDX_task_host_task_id` (`task_id`),
 KEY `IDX_task_host_host_id` (`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='任务主机关联表';

CREATE TABLE `task_log` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `task_id` varchar(11) NOT NULL COMMENT '任务id',
    `name` varchar(64) NOT NULL COMMENT '名称',
    `spec` varchar(64) NOT NULL COMMENT 'cron表达式',
    `ip` varchar(255) NOT NULL COMMENT '主机ip',
    `operation_id` varchar(255) NOT NULL COMMENT '操作id',
    `command` varchar(1024) NOT NULL COMMENT '执行命令',
    `exec_type` int NULL COMMENT '执行类型：0 定时执行，1 手动执行',
    `exec_status` int NULL COMMENT '执行状态: 0 未运行,1 运行中,2 正常,3 异常',
    `exec_result` longtext NOT NULL COMMENT '执行结果',
    `start_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `end_time` datetime NULL DEFAULT NULL COMMENT '结束时间',
    PRIMARY KEY (`id`),
    KEY `IDX_task_log_operation_id` (`operation_id`),
    KEY `IDX_task_log_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='任务日志表';

CREATE TABLE `deploy_backup_history` (
 `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
 `cluster_id` int(11) NOT NULL COMMENT '集群id',
 `db_name` varchar(255) NOT NULL COMMENT '数据库名称',
 `backup_sql` varchar(255) NOT NULL COMMENT '备份文件名称',
 `product_name` varchar(255) DEFAULT NULL COMMENT '触发此次备份的产品包名称',
 `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='备份历史记录表';

CREATE TABLE `deploy_upgrade_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `cluster_id` int(11) NOT NULL COMMENT '集群id',
  `product_name` varchar(64) NOT NULL COMMENT '产品名称',
  `source_version` varchar(64) NOT NULL COMMENT '源版本',
  `target_version` varchar(64) NOT NULL COMMENT '目标版本',
  `backup_name` varchar(64) DEFAULT '' COMMENT '备份名称，值为当前时间',
  `source_service_ip` text COMMENT '源版本服务编排信息',
  `source_config` text COMMENT '源版本服务配置信息',
  `source_multi_config` text COMMENT '源版本服务多配置信息',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `type` tinyint(4) NOT NULL COMMENT '类型，0为升级，1为回滚',
  `backup_sql` text COMMENT '备份SQL文件',
  `is_deleted` tinyint(4) DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='升级历史记录表';
