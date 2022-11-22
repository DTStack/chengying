Easyagent
============

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

[English](README_EN.md) | 中文

# 介绍
easyagent是在袋鼠云内部广泛使用的基础架构组件，最佳应用场景包括ELK体系[beats](https://github.com/elastic/beats)等数据采集器的管控和配置管理、数栈体系自动化部署等

# 基本原理
easyagent主要有sidecar和server两个组件，sidecar部署在主机端，sidecar和server间使用[gRPC](https://github.com/grpc/grpc-go)进行通信，使用[protobuf](https://github.com/gogo/protobuf)定义接口, 仅需sidecar到server单向网络可达，即可通过调用server端提供的REST API实现对目标主机进行服务安装、执行可执行文件等操作。整体架构如下图：
<div align=center>
  <img src=docs/images/easyagent.png width=700 />
</div>

# 快速开始

请点击[快速开始](docs/quickstart.md)

# API 文档

Please click [Api 文档](docs/server-api.md)

# 通信接口定义

请点击[通信接口定义](docs/protoc.md)

# Questions

请点击[Questions](docs/questions.md)

# 如何贡献

请点击[如何贡献](docs/contribution.md)

# License

easyagent is under the Apache 2.0 license. See the [LICENSE](http://www.apache.org/licenses/LICENSE-2.0) file for details.
