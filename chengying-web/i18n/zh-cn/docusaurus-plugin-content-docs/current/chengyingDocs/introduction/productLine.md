---
id: productLine
title: 产品线自动部署
---



[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

## 特性

chengYing 可以根据用户自定义的产品线，根据服务的亲和性配置，自动部署产品包

前提要求：

- 主机设置了角色信息
- 产品包服务设置了亲和性


## 设置主机角色

### 选择主机
"部署中心-集群详情-主机资源"，点击编辑角色

![q1](/img/deploy/20220908151822.jpg)

### 角色编辑

分配已有角色，或添加
![q2](/img/deploy/20220908151937.jpg)

## 服务亲和性

参考[scheme亲和性配置](schema#服务内关键字段)

## 上传产品线

"组件管理-部署组件-选择集群"，点击下一步，上传产品线

![q3](/img/deploy/20220908151315.jpg)

### 产品线示例
```
{
    "product_line_name":"DTFront产品线",
    "product_line_version":"1.0",
    "product_serial":[{"id":1,"product_name":"DTBase","dependee":0},{"id":2,"product_name":"DTFront","dependee":1}]
}
```

上述产品线说明：
- DTBase的dependee为0，表示没有依赖
- DTFront的dependee为1，表示依赖id为1的组件包

:::tip
产品线的依赖不能有死循环，否则将会导致部署失败
:::

### 选择产品线

"选择产品线-点击下一步"
![q4](/img/deploy/20220908153825.jpg)

### 开始部署

执行部署，等待部署结果

![q5](/img/deploy/20220908154041.jpg)




