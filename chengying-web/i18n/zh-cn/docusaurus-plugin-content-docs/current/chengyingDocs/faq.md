---
title: FAQ
sidebar_position: 5
---
## FAQ

### Q: 主机下架时，报当前无Agent停止的主机，不支持下架
A: 主机下架需要去服务器停止agent，步骤如下：
1. 登录服务器，crontab -e查看crontab任务，删除已经存在的自动拉起agent的任务；
2. 执行脚本：sh /opt/dtstack/easymanager/easyagent/easyagent.sh stop;
3. 删除文件夹：rm -rf /opt/dtstack/easymanager/easyagent;
4. 等待几分钟页面就会显示agent已停止，可以下架。

### Q: 查看集群总览页面提示报错：Database err: sql:Scan error on column index 0, name id converting NULL to int is unsupported
A:  报错原因为其中某一台机器在初始化的时候添加了多次，导致sidecar_list和deploy_cluster_host_rel表中出现脏数据，需要手动执行SQL清除脏数据。
```aidl
select *
from deploy_cluster_host_rel
where is_deleted = 0
and clusterId = 1 // clusterId根据实际情况修改
and sid not in (select deploy_host.sid
from deploy_host
where deploy_host.isDeleted = 0
);
```
删除查询出来的记录即可。
