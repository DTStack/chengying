ALTER TABLE `deploy_cluster_product_rel`
    CHANGE COLUMN `status` `status` ENUM('undeployed', 'deploying', 'deployed', 'deploy fail', 'undeploying', 'undeploy fail', 'dbrollbacking', 'dbrollbacked', 'dbrollbackfail', 'rollbacking') NOT NULL DEFAULT 'undeployed' COMMENT '产品状态' ;


ALTER TABLE `deploy_backup_history`
    ADD COLUMN `product_version` VARCHAR(45) NOT NULL DEFAULT '' COMMENT '产品包版本' AFTER `product_name`,
    ADD COLUMN `status` VARCHAR(45) NOT NULL DEFAULT '' COMMENT '备份状态' AFTER `product_version`,
    ADD COLUMN `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`;


ALTER TABLE `deploy_backup_history`
    ADD COLUMN `exec_id` VARCHAR(50) NOT NULL DEFAULT '' AFTER `product_version`;


CREATE TABLE `rollback_exec_log` (
                                     `id` int(11) NOT NULL AUTO_INCREMENT,
                                     `product_name` varchar(45) COLLATE utf8mb4_bin DEFAULT NULL,
                                     `exec_id` varchar(45) COLLATE utf8mb4_bin DEFAULT NULL,
                                     `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

INSERT INTO `sys_config`(config_name,config_value) VALUES ('globalconfig.service_install_timeout_limit','30');

CREATE TABLE `deploy_smooth_upgrade_list` (
                                              `id` int(11) NOT NULL AUTO_INCREMENT,
                                              `product_name` varchar(255) NOT NULL COMMENT '产品名称',
                                              `service_name` varchar(255) NOT NULL COMMENT '服务名称',
                                              `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                              `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
                                              PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='平滑升级列表';

INSERT INTO `deploy_smooth_upgrade_list`(product_name,service_name) VALUES ('DTApi','ApiServer');

INSERT INTO `deploy_smooth_upgrade_list`(product_name,service_name) VALUES ('DTGateway','Gateway');

ALTER TABLE `deploy_upgrade_history` ADD `upgrade_mode` VARCHAR(32) NOT NULL COMMENT '升级模式' AFTER `backup_sql`;

CREATE TABLE `deploy_cluster_smooth_upgrade_product_rel` (
                                                             `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
                                                             `pid` int(11) DEFAULT NULL,
                                                             `clusterId` int(11) unsigned NOT NULL COMMENT '集群 id',
                                                             `namespace` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'cluster namespace',
                                                             `deploy_uuid` varchar(36) COLLATE utf8mb4_bin NOT NULL COMMENT '部署uuid',
                                                             `product_parsed` longtext COLLATE utf8mb4_bin NOT NULL COMMENT '已经解析的产品信息',
                                                             `status` enum('undeployed','deploying','deployed','deploy fail','undeploying','undeploy fail') COLLATE utf8mb4_bin NOT NULL DEFAULT 'undeployed' COMMENT '产品状态',
                                                             `alert_recover` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0，不恢复告警，1，恢复告警',
                                                             `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '部署人id',
                                                             `is_deleted` int(11) NOT NULL DEFAULT '0' COMMENT '0:未删除,1:已删除',
                                                             `deploy_time` datetime DEFAULT NULL COMMENT '部署时间',
                                                             `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                                                             `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                                             PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `deploy_mysql_ip_list` (
                                        `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                                        `cluster_id` int(11) unsigned NOT NULL COMMENT '集群id',
                                        `namespace` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'cluster namespace',
                                        `product_name` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '组件名称',
                                        `mysql_ip_list` varchar(1024) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'Mysql IP列表',
                                        `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update_time',
                                        `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create_time',
                                        PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

ALTER TABLE `task_list` ADD `exec_timeout` int(11) NOT NULL COMMENT '执行超时时间' AFTER `spec`;

ALTER TABLE `task_list` ADD `log_retention` int(11) NOT NULL COMMENT '执行历史保存周期' AFTER `exec_timeout`;

UPDATE task_list SET exec_timeout=60, update_time=NOW() WHERE exec_timeout=0 AND is_deleted=0;

UPDATE task_list SET log_retention=7, update_time=NOW() WHERE log_retention=0 AND is_deleted=0;

INSERT INTO `sys_config`(config_name,config_value) VALUES ('inspectconfig.fullGC_time','1');
INSERT INTO `sys_config`(config_name,config_value) VALUES ('inspectconfig.fullGC_frequency','0');
INSERT INTO `sys_config`(config_name,config_value) VALUES ('inspectconfig.dir_size','0');
INSERT INTO `sys_config`(config_name,config_value) VALUES ('inspectconfig.node_cpu_usage','0');
INSERT INTO `sys_config`(config_name,config_value) VALUES ('inspectconfig.node_mem_usage','0');
INSERT INTO `sys_config`(config_name,config_value) VALUES ('inspectconfig.node_disk_usage','0');
INSERT INTO `sys_config`(config_name,config_value) VALUES ('inspectconfig.node_inode_usage','0');

INSERT INTO `sys_config`(config_name,config_value) VALUES ('globalconfig.auto_test_timeout_limit','20');

INSERT INTO `inspect_report_template`(`id`, `cluster_id`, `type`, `module`, `metric`, `targets`, `unit`, `decimal`, `create_time`, `update_time`, `is_deleted`) VALUES
                                                                                                                                                                    (99, 1, 2, '{UPServiceName}', 'Full GC Count (2minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{cluster_name=\\\"{cluster_name}\\\",gc=\\\"G1 Old Generation\\\",product_name=\\\"{ProductName}\\\",service_name=\\\"{ServiceName}\\\"}[2m]))\",\"legend_format\":\"Full GC Count-{{instance}}\"}]', '', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (100, 1, 3, 'HDFS_NameNode', 'Full GC Count (2minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{cluster_name=\\\"{cluster_name}\\\",gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}[2m]))\",\"legend_format\":\"Full GC Count-{{instance}}\"}]', '', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (101, 1, 3, 'HDFS_DataNode', 'Full GC Count (2minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{cluster_name=\\\"{cluster_name}\\\",gc=\\\"PS MarkSweep\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_datanode\\\"}[2m]))\",\"legend_format\":\"Full GC Count-{{instance}}\"}]', '', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (102, 1, 3, 'hivemetastore', 'Full GC Count (2minutes)', '[{\"expr\":\"floor(delta(jvm_gc_collection_seconds_count{cluster_name=\\\"{cluster_name}\\\",gc=\\\"G1 Old Generation\\\",product_name=\\\"Hive\\\",service_name=\\\"hivemetastore\\\"}[2m]))\",\"legend_format\":\"Full GC Count-{{instance}}\"}]', '', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (103, 1, 3, 'HDFS 文件存储使用率', '', '[{\"expr\":\"100-100*((node_filesystem_avail{cluster_name=\\\"{cluster_name}\\\",device=~\\\"/dev/.*\\\",mountpoint=\\\"/data\\\"}/node_filesystem_size{cluster_name=\\\"{cluster_name}\\\",device=~\\\"/dev/.*\\\",mountpoint=\\\"/data\\\"}))\",\"legend_format\":\"instance-{{instance}}\"}]', '%', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (104, 1, 3, 'HDFS 小文件数量趋势', '', '[{\"expr\":\"Hadoop_NameNode_TotalFiles{cluster_name=\\\"{cluster_name}\\\",name=\\\"FSNamesystem\\\",product_name=\\\"Hadoop\\\",service_name=\\\"hdfs_namenode\\\"}\",\"legend_format\":\"TotalFiles-{{instance}}\"}]', '', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (105, 1, 4, 'CPU使用趋势', 'CPU Used(%)', '[{\"expr\":\"100 - ( avg(irate(node_cpu{cluster_name=\\\"{cluster_name}\\\",mode=\\\"idle\\\"}[2m])) by (instance) * 100 )\",\"legend_format\":\"instance-{{instance}}\"}]', '%', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (106, 1, 4, '内存使用趋势', 'Mem Used(%)', '[{\"expr\":\"(1 - (sum by (instance) (node_memory_MemFree{cluster_name=\\\"{cluster_name}\\\",} + node_memory_Buffers{cluster_name=\\\"{cluster_name}\\\",} + node_memory_Cached{cluster_name=\\\"{cluster_name}\\\",})) / (sum by (instance) (node_memory_MemTotal{cluster_name=\\\"{cluster_name}\\\",})) ) * 100\",\"legend_format\":\"instance-{{instance}}\"}]', '%', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (107, 1, 4, '磁盘使用率趋势', 'Disk Used(%)', '[{\"expr\":\"100-100*((node_filesystem_avail{cluster_name=\\\"{cluster_name}\\\",device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"}/node_filesystem_size{cluster_name=\\\"{cluster_name}\\\",device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"}))\",\"legend_format\":\"instance-{{instance}}{{mountpoint}}\"}]', '%', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (108, 1, 4, 'inode 使用率趋势', 'Inode Used(%)', '[{\"expr\":\"100*((node_filesystem_files{cluster_name=\\\"{cluster_name}\\\",device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"}-node_filesystem_files_free{cluster_name=\\\"{cluster_name}\\\",device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"})/node_filesystem_files{cluster_name=\\\"{cluster_name}\\\",device=~\\\"/dev/.*\\\",mountpoint!=\\\"/boot\\\"})\",\"legend_format\":\"instance-{{instance}}{{mountpoint}}\"}]', '%', 0, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0),
                                                                                                                                                                    (109, 1, 4, 'swap 使用量趋势', 'Swap Used', '[{\"expr\":\"((sum by(instance)(node_memory_SwapTotal{cluster_name=\\\"{cluster_name}\\\"}))-(sum by(instance)(node_memory_SwapFree{cluster_name=\\\"{cluster_name}\\\"}))-(sum by(instance)(node_memory_SwapCached{cluster_name=\\\"{cluster_name}\\\"})) )\",\"legend_format\":\"instance-{{instance}}\"}]', '', 2, '2022-05-18 11:32:29', '2022-05-18 11:32:29', 0);

CREATE TABLE `deploy_product_line_list` (
                                            `id` int(11) NOT NULL AUTO_INCREMENT,
                                            `product_line_name` varchar(255) NOT NULL COMMENT '产品线名称',
                                            `product_line_version` varchar(255) NOT NULL COMMENT '产品线版本',
                                            `product_serial` varchar(1024) NOT NULL COMMENT '产品系列',
                                            `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                            `update_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                                            `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
                                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='产品线列表';

CREATE TABLE `deploy_service_relations_list` (
                                                 `id` int(11) NOT NULL AUTO_INCREMENT,
                                                 `relations_type` varchar(255) NOT NULL COMMENT '关系类型',
                                                 `source_product_name` varchar(255) NOT NULL COMMENT '源产品名称',
                                                 `source_service_name` varchar(255) NOT NULL COMMENT '源服务名称',
                                                 `target_product_name` varchar(255) NOT NULL COMMENT '目标产品名称',
                                                 `target_service_name` varchar(255) NOT NULL COMMENT '目标服务名称',
                                                 `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                                 `update_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                                                 `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
                                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='服务关系列表';

INSERT INTO `deploy_service_relations_list` (relations_type, source_product_name, source_service_name, target_product_name, target_service_name, is_deleted) VALUES ('conflict', 'DTBase', 'mysql_slave', 'DTBase', 'mysql', 0);
INSERT INTO `deploy_service_relations_list` (relations_type, source_product_name, source_service_name, target_product_name, target_service_name, is_deleted) VALUES ('relyOn', 'DTBase', 'pushgateway', 'DTBase', 'prometheus', 0);

