
## server configuration

```
# logger
log.dir: /tmp/em2.0
log.max-logger-size: 100
log.max-logger-backups: 3
log.days-to-keep: 1

# optional, publish event and metric to other system:
#  elasticsearch:
#    hosts: ["http://localhost:9200"]
#    username: ""
#    password: ""
#  influxdb:
#    hosts: ["http://localhost:8086"]
#    username: "admin"
#    password: "admin"
#  transfer:
#    server: localhost
#    port: 8790
#    concurrency: 4    # optional, default 1
#    timeout: 3s       # optional, default 3s
#    cert:             # optional
#    tls:              # optional
#    tls-skip-verify:  # optional
#  kafka:
#    hosts: ["121.43.166.210:9092"]
#    username: ""
#    password: ""
#    timeout: 3s
#  http:
#    host: "127.0.0.1:8864"
#    uri: /api/v2/instance/%s/%s #first %s means agent_id, second %s means type
#    timeout: 30s

# databases
mysqldb.host: 172.16.10.74
mysqldb.port: 3306
mysqldb.user: easyagent
mysqldb.password: dtstack
mysqldb.dbname: dtagent

# api
api.port: 8889
api.host:

# rpc
rpc.port: 8890
rpc.cert:
rpc.key:

```


## sidecar register mode configuration

```
# logger
log:
  dir: /tmp/
  max-logger-size: 100
  max-logger-backups: 3
  days-to-keep: 1

# rpc server
rpc:
  server: localhost
  port: 8890
  cert:
  tls:
  tls-skip-verify:

# uuid must be unique
# uuid is generated when install the sidecar witn curl cmd
easyagent:
  uuid: 87dacbba-5be4-4c85-9b05-f313b1ee80b9
  network: [eth0]
  monitor-interval: 30s

# http request triggered when sidecar register to server
# CallBack=url_encode(base64(url))
callback:
  - CallBack=aHR0cDovL21hdHJpeDo4ODY0L2FwaS92Mi9hZ2VudC9pbnN0YWxsL2NhbGxiYWNrP2FpZD0tMQ%3D%3D&TargetPath=%2Fopt%2Fdtstack%2Feasymanager%2Feasyagent
```

## sidecar local mode configuration

```
# logger
log:
  dir: /tmp/
  max-logger-size: 100
  max-logger-backups: 3
  days-to-keep: 1

# rpc server
rpc:
  server: localhost
  port: 8890
  cert:
  tls:
  tls-skip-verify:

# uuid must be unique
easyagent:
  uuid: 87dacbba-5be4-4c85-9b05-f313b1ee80b9
  network: [eth0]
  monitor-interval: 30s
  mode: local

# http request triggered when sidecar register to server
# CallBack=url_encode(base64(url))
callback:
  - CallBack=aHR0cDovL21hdHJpeDo4ODY0L2FwaS92Mi9hZ2VudC9pbnN0YWxsL2NhbGxiYWNrP2FpZD0tMQ%3D%3D&TargetPath=%2Fopt%2Fdtstack%2Feasymanager%2Feasyagent
```