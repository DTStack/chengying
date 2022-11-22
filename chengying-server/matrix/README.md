## Matrix 目录结构

```
├─agent                         // easyagent 生命周期管理
|
├─api                           // 接口定义
│  │  addons-api.go 			
│  │  agent-api.go              // 添加主机相关
│  │  api-schema.go 			
│  │  cluster-api.go            // 集群相关，包含主机集群，k8s 集群，镜像仓库操作
│  │  common-api.go             // 安全审计相关
│  │  group-api.go              // 组件启停，多个服务可能属于同个组
│  │  instance-api.go           // 实例相关
│  │  instance-record-api.go 	
│  │  product-api.go            // 产品包相关，包含产品包上传、传统集群，k8s集群的产品部署和卸载
│  │  role-api.go               // 角色权限相关
│  │  service-api.go            // 服务相关
│  │  user-api.go               // 用户相关
│  │
│  └─impl                       // 接口具体实现
|
├─asset                         // go-bindata 生成代码的存放位置
├─base                          // matrix 初始化相关
├─discover                      // Prometheus 服务发现，节点发现
├─event                         // 传统部署实例事件封装
├─grafana                       // 产品告警相关
├─group                         // 组件生命周期管理
├─harole						
├─host                          // easyagent 相关函数封装
├─instance                      // 实例生命周期管理
├─k8s
│  ├─deploy                     // k8s 部署实例函数封装
│  ├─docker                     // 镜像相关函数封装
│  ├─model                      // k8s 资源对象生成
│  ├─monitor                    // k8s 产品包部署后状态监控
│  ├─node                       // k8s 集群节点回调处理
│  ├─util
│  ├─web-socket                 // web-socket 相关函数封装
│  └─xke-service                // 自建 k8s 集群函数封装
│      └─driver
├─log                           // 日志工具封装
├─model                         // 数据库相关
├─service                       // 服务生命周期管理
├─strategy
├─templates                     // 脚本，yml 模板存放，通过 go-bindata 自动生成代码
└─util
```
