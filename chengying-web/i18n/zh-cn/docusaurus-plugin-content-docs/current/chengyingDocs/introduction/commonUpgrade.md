---
id: commonUpgrade
title: 普通升级
---

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)


## 特性

用户可在"部署中心-集群详情-已部署组件"，点击升级按钮，选择普通升级

![q1](/img/upgrade/20220908135645.jpg)

前提要求:

- 产品包中包含sql组件，服务名称后缀为Sql

- 产品包中包含数据库备份脚本，数据库回滚脚本

- chengying中包含多个版本

scheme 配置如下

```
  dtuicSql:
    version: 
    instance:
      pseudo: true
      cmd: ./waiting.sh
      logs:
      - logs/*.log
      config_paths:
        - ./post_deploy.sh
      backup: backup.sh        // 备份脚本，放置于产品包目录下
      rollback: rollback.sh    // 回滚脚本，放置于产品包目录下
    group: UIC
    config:
      mysql_host: "${@mysql}"
      mysql_user: drpeco
      mysql_pass: DT@Stack#123
      mysql_port: 3306
      mysql_db: "dtuic"        // 数据库名称
```

## 升级

### 备份数据库

选择目标组件-点击开始备份,查看备份结果
![q2](/img/upgrade/20220908141949.jpg)

### 开始升级

"点击升级-执行部署"

![q3](/img/upgrade/20220908142142.jpg)

升级结果

![q4](/img/upgrade/20220908142220.jpg)

## 回滚

### 还原数据库
"选择目标组件-选择备份还原库-点击开始备份",查看备份结果

![q5](/img/upgrade/20220908142638.jpg)


### 开始回滚

"点击回滚"，等待回滚结果

![q6](/img/upgrade/20220908143252.jpg)



