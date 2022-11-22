# 参数配置组件

所需参数：

| 字段名称          | 值类型                        | 默认值   | 说明                                 |
| ----------------- | ----------------------------- | -------- | ------------------------------------ |
| defaultActiveKey  | string('resource' \| 'param') | resource | resourcef 为主机分配，param 为参数配置 |
| existIp?          | string[]                      | []       | 已存在的 host                         |
| hostList          | host[]                        | []       | 所有主机列表                         |
| serviceData       | Object                        | {}       | 当前服务的所有信息                   |
| saveResourceState | Function                      |          | host 发生更改的时候调用的方法         |
| resourceState     | Object                        |          | host 分配页面的 state                  |
|                   |                               |          |                                      |
|                   |                               |          |                                      |
|                   |                               |          |                                      |
|                   |                               |          |                                      |
|                   |                               |          |                                      |

