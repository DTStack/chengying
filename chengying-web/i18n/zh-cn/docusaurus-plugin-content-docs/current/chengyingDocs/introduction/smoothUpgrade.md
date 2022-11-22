---
id: smoothUpgrade
title: 平滑升级
---

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)


## 特性

与普通升级不同的是，普通升级同一时间只能部署单个版本的组件，平滑升级的可以实现多版本共存

前提要求:

- 产品包中包含sql组件，服务名称后缀为Sql

- 产品包中包含数据库备份脚本，数据库回滚脚本

- chengying中包含多个版本

- 需要第三方数据库支持


## 特性

用户可在"部署中心-集群详情-已部署组件"，点击升级按钮，选择平滑升级

![q1](/img/upgrade/20220908135645.jpg)

普通升级只针对特定的产品包，要求如下:

- 产品包中包含sql组件，服务名称后缀为Sql

- 产品包中包含数据库备份脚本，数据库回滚脚本

- chengying中包含多个版本

- scheme 配置同普通升级

- 部署的组件中需要有服务部署在两台agent及以上

## 升级

### 备份数据库

选择目标组件-点击开始备份,查看备份结果
![q2](/img/upgrade/20220908144136.jpg)

### 配置数据库

配置第三方数据库

![q3](/img/upgrade/20220908145423.jpg)

### 开始部署

![q3](/img/upgrade/20220908142142.jpg)

升级结果

![q4](/img/upgrade/20220908142220.jpg)

## 回滚

### 还原数据库

"选择目标组件-选择备份还原库"

![q5](/img/upgrade/20220908150134.jpg)


### 开始回滚

"点击回滚"，等待回滚结果

![q6](/img/upgrade/20220908150249.jpg)




