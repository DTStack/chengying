Easyagent
============

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

English | [中文](README.md)

# Introduction
Easyagent is an infrastructure component based on gRPC widely used in dtstack, applied to manage the life-cycle of service on the remote host. The best practice is to manage the life-cycle of the beats of ELK and Hadoop components.

# Fundamental
Easyagent is composed of two parts: sidecar and server, the former is deployed on the remote hosts and communicate with the latter with [gRPC](https://github.com/grpc/grpc-go) messages, which are defined with [protobuf](https://github.com/gogo/protobuf).
With one-way gRPC-connection between the sidecar and the server(from sidecar to server), upon the rest api supported by the easyagent server, other systems or components can control the remotes sidecar(host) and install/start/stop/config/uninstall the target application.
<div align=center>
  <img src=docs/images/easyagent.png width=700 />
</div>

# Quick Start

Please click [Quick Start](docs/quickstart.md)

# API Reference

Please click [Api Reference](docs/server-api.md)

# protoc messages

Please click [Quick Start](docs/quickstart.md)

# How to contribute

Please click [Contribution](docs/contribution.md)

# License

easyagent is under the Apache 2.0 license. See the [LICENSE](http://www.apache.org/licenses/LICENSE-2.0) file for details.
