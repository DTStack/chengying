---
id: quickstart
title: 快速开始
sidebar_position: 2
---
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

本文讨论如何快速部署承影。

## 环境要求

| 类别       | 要求                                                                    |
|:---------|:----------------------------------------------------------------------|
| 操作系统     | CentOS 7 以上                                                           |
| 时间同步     | 所有主机时区和时间必须一致                                                         |
| 磁盘目录     | 系统盘至少100G，数据盘挂载到 /data 目录                                             |
| 系统编码     | UTF-8                                                                 |
| 主机名      | 不重复                                                                   |
| 其他       | 要求所有主机网络互通<br />使用 root 部署 ChengYing 自身，具体产品可使用 root 或者带 sudo 权限的其他用户 |

## 快速启动

**请自行下载 ChengYing [最新版本的离线安装包](https://github.com/DTStack/chengying/releases/download/v1.0.0/chengying-release-1.0.tar.gz)**

```
# 解压安装包
$ tar zxvf chengying-release-1.0.tar.gz
$ cd chengying-release-1.0
# 运行安装脚本
$ sh install.sh local_ip			#local_ip为Chengying所在机器ip
```

**安装完成后，使用 docker ps 命令检查服务状态**
![](/img/quickstart/CheckStatus.png)

**若有异常，可以使用 docker-compose up -d 命令进行重新启动**

```
$ cd chengying-release-1.0/easymanager
$ docker-compose up -d
```

**登录**

```
地址：http://local_ip # local_ip为Chengying所在机器ip
用户名：admin@dtstack.com
密码：DT#passw0rd2019
```

## 源码构建

### Chengying-Server

**克隆代码仓库**

```
git clone https://github.com/DTStack/chengying-server.git
```

**安装依赖**

- Docker 18.09+ ([installation manual](https://docs.docker.com/install))
  
- Golang 1.14+ ([installation manual](https://golang.org/dl/))
- Iris ([installation manual](https://github.com/kataras/iris/))
- Gomonkey ([installation manual](https://github.com/wangqi811/gomonkey/))
- Gox ([installation manual](https://github.com/mitchellh/gox/))
- Go-bindata([go-bindata](https://github.com/go-bindata/go-bindata))

**初始化配置**

matrix 会默认加载配置文件 chengying-server/matrix/example-config.yml，请修改该配置文件

**本地调试**

```
$ cd matrix
$ go run main.go config.go -c example-config.yml --debug
```

**构建 matrix 容器镜像**

```
$ cd matrix
$ gox -os=linux -arch=amd64
$ mv matrix_linux_amd64 matrix
$ docker build -t matrix:latest .
```

### Chengying-Agent
**克隆代码仓库**
```
git clone https://github.com/DTStack/chengying-agent.git
```

**安装依赖**
- Golang 1.12+
- OS linux/unix/windows
- Protoc([protobuf](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.1))
- Go-bindata([go-bindata](https://github.com/go-bindata/go-bindata))

**构建**

支持三种操作系统linux/windows/darwin的构建
```
make all
```

### Chengying-Front
基于 webpack 的开发配置环境，可以使用 npm 管理项目
**克隆代码仓库**
```
git clone https://github.com/DTStack/chengying-front.git
```

**依赖**

推荐使用 yarn 管理依赖
```
yarn add package.name
```

**构建**
```
##开发构建
npm start 或者 npm run dev

##生产构建
npm run build
```


