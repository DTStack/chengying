---
title: 监控告警
---

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

### 仪表盘
仪表盘用于展示集群各服务、主机的监控趋势，快速掌握集群运行的稳定性以及异常情况，及时运维。
#### 创建仪表盘
步骤一：创建仪表盘，填写仪表盘所属文件夹、名称、标签

![q1](/img/monitor/newBoard.png)
:::tip
仪表盘标签一般需要写产品名称、组件名称、服务名称，便于用标签筛选仪表盘。
:::

步骤二：在仪表盘中添加Panel
![q2](/img/monitor/addPanel.png)
仪表盘的Panel包含趋势图、单个指标、表格、文本、热点图、告警、饼图等，用户可根据需要自行选择。

步骤三：进入Panel编辑页面
![q3](/img/monitor/editPanel.png)
- 设置Panel基本信息
![q4](/img/monitor/panelBasic.png)
- 设置Panel查询条件
![q5](/img/monitor/panelCondition.png)
DataSource选择prometheus，下方写查询条件，按照Grafana的标准语句查询。
- 设置Panel坐标轴
![q6](/img/monitor/dimension.png)
- 设置Panel图例
![q7](/img/monitor/legend.png)
- 设置数据显示
![q8](/img/monitor/display.png)

#### 导入仪表盘
可通过JSON文件和上传仪表盘文件两种方式上传已有仪表盘。 
既然可以导入仪表盘，同样，仪表盘支持复制JSON语句，以及导出仪表盘。
![q9](/img/monitor/importDashboard.png)

### 告警
#### 告警内容设置
以配置Redis的监控告警为例：

步骤一：选择具体的仪表盘
![q10](/img/monitor/redisBoard.png)
步骤二：选择告警指标，进入编辑页面
![q11](/img/monitor/editRedis.png)
步骤三：切换到Alert配置页面，点击Create Alert
![q12](/img/monitor/createAlert.png)
步骤四：配置告警规则和告警通道，点击测试
![q13](/img/monitor/testAlert.png)
若无告警通道，在告警通道页面添加完成后再在此处设置。
![q14](/img/monitor/newAlertChannel.png)
步骤五：保存告警设置，同时用户可在告警规则页面查看配置的告警，支持对告警规则的启停、检索等
![q15](/img/monitor/saveAlert.png)
![q16](/img/monitor/alertList.png)

#### 告警通道
平台提供 “短信通道、邮件通道、钉钉通道、企业微信通道、自定义通道” 5种通道配置，用户根据需要选择合适通道并填写通道配置信息、消息模板、地址等完成通道配置。
![q17](/img/monitor/addChannel.png)

